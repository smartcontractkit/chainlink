// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keeper_registry_logic2_0

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

var KeeperRegistryLogicMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.Mode\",\"name\":\"mode\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnchainConfigNonEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"checkBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMode\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.Mode\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162006038380380620060388339810160408190526200003591620001ef565b838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c38162000126565b505050836002811115620000db57620000db62000251565b60e0816002811115620000f257620000f262000251565b60f81b9052506001600160601b0319606093841b811660805291831b821660a05290911b1660c05250620002679350505050565b6001600160a01b038116331415620001815760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ea57600080fd5b919050565b600080600080608085870312156200020657600080fd5b8451600381106200021657600080fd5b93506200022660208601620001d2565b92506200023660408601620001d2565b91506200024660608601620001d2565b905092959194509250565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160f81c615d306200030860003960008181610224015281816138870152818161394c015281816145b9015261477001526000818161026e01526141dc0152600081816103b701526142c501526000818161041e01528181610f7e0152818161126301528181611bf6015281816122690152818161259e01528181612a850152612b180152615d306000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80638dcf0fe711610104578063b121e147116100a2578063ca30e60311610071578063ca30e6031461041c578063eb5dcd6c14610442578063f2fde38b14610455578063f7d334ba1461046857600080fd5b8063b121e147146103db578063b148ab6b146103ee578063b79550be14610401578063c80480221461040957600080fd5b80639fab4386116100de5780639fab43861461037c578063a710b2211461038f578063a72aa27e146103a2578063b10b673c146103b557600080fd5b80638dcf0fe7146103435780638e86139b14610356578063948108f71461036957600080fd5b80636ded9eae1161017c5780638456cb591161014b5780638456cb59146102f757806385c1b0ba146102ff5780638765ecbe146103125780638da5cb5b1461032557600080fd5b80636ded9eae146102b3578063744bfe61146102d457806379ba5097146102e75780637d9b97e0146102ef57600080fd5b80633f4ba83a116101b85780633f4ba83a1461021a5780634b4fd03b146102225780635165f2f5146102595780636709d0e51461026c57600080fd5b8063187256e8146101df5780631a2af011146101f45780633b9cce5914610207575b600080fd5b6101f26101ed366004614e8d565b61048d565b005b6101f261020236600461521b565b6104fe565b6101f2610215366004614f6a565b610652565b6101f26108a8565b7f000000000000000000000000000000000000000000000000000000000000000060405161025091906157ef565b60405180910390f35b6101f2610267366004615202565b61090e565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610250565b6102c66102c1366004614ec8565b610a88565b604051908152602001610250565b6101f26102e236600461521b565b610c77565b6101f2611089565b6101f261118b565b6101f26112f5565b6101f261030d366004614fac565b611366565b6101f2610320366004615202565b611c81565b60005473ffffffffffffffffffffffffffffffffffffffff1661028e565b6101f261035136600461523e565b611e08565b6101f2610364366004615197565b611e6a565b6101f26103773660046152ad565b6120a6565b6101f261038a36600461523e565b612345565b6101f261039d366004614e5a565b6123f4565b6101f26103b036600461528a565b61267f565b7f000000000000000000000000000000000000000000000000000000000000000061028e565b6101f26103e9366004614e3f565b612761565b6101f26103fc366004615202565b612859565b6101f2612a4c565b6101f2610417366004615202565b612bb7565b7f000000000000000000000000000000000000000000000000000000000000000061028e565b6101f2610450366004614e5a565b612f6c565b6101f2610463366004614e3f565b6130cb565b61047b610476366004615202565b6130df565b60405161025096959493929190615660565b610495613736565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260166020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660018360038111156104f5576104f5615c11565b02179055505050565b610507826137b9565b73ffffffffffffffffffffffffffffffffffffffff8116331415610557576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166105a4576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff82811691161461064e5760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b61065a613736565b600b548114610695576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600b54811015610867576000600b82815481106106b7576106b7615c6f565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600c9092526040832054919350169085858581811061070157610701615c6f565b90506020020160208101906107169190614e3f565b905073ffffffffffffffffffffffffffffffffffffffff811615806107a9575073ffffffffffffffffffffffffffffffffffffffff82161580159061078757508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156107a9575073ffffffffffffffffffffffffffffffffffffffff81811614155b156107e0576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146108515773ffffffffffffffffffffffffffffffffffffffff8381166000908152600c6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061085f90615b56565b915050610698565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600b838360405161089c93929190615466565b60405180910390a15050565b6108b0613736565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b610917816137b9565b600081815260046020908152604091829020825160e081018452815463ffffffff8082168352640100000000820481169483019490945268010000000000000000810460ff1615159482018590526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c082015290610a19576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169055610a5860028361386c565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590610ad957506011546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314155b15610b10576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b2c6001610b1d613881565b610b279190615a8e565b613946565b601254604080516020810193909352309083015268010000000000000000900463ffffffff1660608201526080016040516020818303038152906040528051906020012060001c9050610bba8189898960008a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509250613a15915050565b6012805468010000000000000000900463ffffffff16906008610bdc83615b8f565b825463ffffffff9182166101009390930a9283029190920219909116179055506000818152601760205260409020610c15908484614963565b506040805163ffffffff8916815273ffffffffffffffffffffffffffffffffffffffff8816602082015282917fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012910160405180910390a2979650505050505050565b600f546f01000000000000000000000000000000900460ff1615610cc7576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116610d4e576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815160e081018352815463ffffffff8082168352640100000000820481168387015260ff6801000000000000000083041615158386015273ffffffffffffffffffffffffffffffffffffffff6901000000000000000000909204821660608401526001909301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490921660c082015286855260059093529220549091163314610e5b576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e63613881565b816020015163ffffffff161115610ea6576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546015546c010000000000000000000000009091046bffffffffffffffffffffffff1690610ee6908290615a8e565b60155560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b158015610fc457600080fd5b505af1158015610fd8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ffc9190615109565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461110f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611193613736565b6011546015546bffffffffffffffffffffffff909116906111b5908290615a8e565b601555601180547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b602060405180830381600087803b1580156112bd57600080fd5b505af11580156112d1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061064e9190615109565b6112fd613736565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff166e0100000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602001610904565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156113a2576113a2615c11565b141580156113ea5750600373ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156113e7576113e7615c11565b14155b15611421576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16611480576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816114b7576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff81111561150b5761150b615c9e565b60405190808252806020026020018201604052801561153e57816020015b60608152602001906001900390816115295790505b50905060008667ffffffffffffffff81111561155c5761155c615c9e565b604051908082528060200260200182016040528015611585578160200160208202803683370190505b50905060008767ffffffffffffffff8111156115a3576115a3615c9e565b60405190808252806020026020018201604052801561162857816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816115c15790505b50905060005b888110156119b35789898281811061164857611648615c6f565b60209081029290920135600081815260048452604090819020815160e081018352815463ffffffff8082168352640100000000820481169783019790975268010000000000000000810460ff16151593820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c0840152985090965061172a9050876137b9565b8582828151811061173d5761173d615c6f565b602002602001018190525060076000888152602001908152602001600020805461176690615b02565b80601f016020809104026020016040519081016040528092919081815260200182805461179290615b02565b80156117df5780601f106117b4576101008083540402835291602001916117df565b820191906000526020600020905b8154815290600101906020018083116117c257829003601f168201915b50505050508482815181106117f6576117f6615c6f565b60200260200101819052506005600088815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683828151811061184757611847615c6f565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015260a0860151611889906bffffffffffffffffffffffff168661594e565b600088815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff00000000000000000000000000000000000000000000000000000000169055600790915281209196506118fd9190614a05565b600087815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905561193c600288613e58565b5060a0860151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8a16602083015288917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a2806119ab81615b56565b91505061162e565b50836015546119c29190615a8e565b6015556040516000906119e1908b908b90859088908890602001615516565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260105490915073ffffffffffffffffffffffffffffffffffffffff808a1691638e86139b916c010000000000000000000000009091041663c71249ab60028c73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b8152600401602060405180830381600087803b158015611a9557600080fd5b505af1158015611aa9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611acd9190615320565b866040518463ffffffff1660e01b8152600401611aec9392919061583e565b60006040518083038186803b158015611b0457600080fd5b505afa158015611b18573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b5e91908101906151cd565b6040518263ffffffff1660e01b8152600401611b7a91906156fd565b600060405180830381600087803b158015611b9457600080fd5b505af1158015611ba8573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8b81166004830152602482018990527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b158015611c3c57600080fd5b505af1158015611c50573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c749190615109565b5050505050505050505050565b611c8a816137b9565b600081815260046020908152604091829020825160e081018452815463ffffffff8082168352640100000000820481169483019490945268010000000000000000810460ff16158015958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c082015290611d8e576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1668010000000000000000179055611dd8600283613e58565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b611e11836137b9565b6000838152601760205260409020611e2a908383614963565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508383604051611e5d9291906156b0565b60405180910390a2505050565b60023360009081526016602052604090205460ff166003811115611e9057611e90615c11565b14158015611ec2575060033360009081526016602052604090205460ff166003811115611ebf57611ebf615c11565b14155b15611ef9576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808080611f0a85870187615000565b935093509350935060005b845181101561209d57611fec858281518110611f3357611f33615c6f565b6020026020010151858381518110611f4d57611f4d615c6f565b602002602001015160600151868481518110611f6b57611f6b615c6f565b602002602001015160000151858581518110611f8957611f89615c6f565b6020026020010151888681518110611fa357611fa3615c6f565b602002602001015160a00151888781518110611fc157611fc1615c6f565b60200260200101518a8881518110611fdb57611fdb615c6f565b602002602001015160400151613a15565b848181518110611ffe57611ffe615c6f565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7185838151811061203957612039615c6f565b602002602001015160a00151336040516120839291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a28061209581615b56565b915050611f15565b50505050505050565b600082815260046020908152604091829020825160e081018452815463ffffffff80821683526401000000008204811694830185905268010000000000000000820460ff161515958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c082015291146121a8576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516121b8919061598b565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff9384160217905560155461221e9184169061594e565b6015556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd90606401602060405180830381600087803b1580156122c257600080fd5b505af11580156122d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122fa9190615109565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b61234e836137b9565b60125474010000000000000000000000000000000000000000900463ffffffff168111156123a8576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526007602052604090206123c1908383614963565b50827f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf8383604051611e5d9291906156b0565b73ffffffffffffffffffffffffffffffffffffffff8116612441576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c60205260409020541633146124a1576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f54600b546000916124d891859170010000000000000000000000000000000090046bffffffffffffffffffffffff1690613e64565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260086020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601554909150612542906bffffffffffffffffffffffff831690615a8e565b6015556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b1580156125e257600080fd5b505af11580156125f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061261a9190615109565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff1610806126a8575060125463ffffffff6401000000009091048116908216115b156126df576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6126e8826137b9565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8516908117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600d60205260409020541633146127c1576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600c602090815260408083208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217909355600d909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600081815260046020908152604091829020825160e081018452815463ffffffff80821683526401000000008204811694830185905268010000000000000000820460ff161515958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c0820152911461295b576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff1633146129b8576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b612a54613736565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b158015612adc57600080fd5b505afa158015612af0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b14919061517e565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb3360155484612b619190615a8e565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff909216600483015260248201526044016112a3565b6000818152600460209081526040808320815160e081018352815463ffffffff80821683526401000000008204811695830186905260ff6801000000000000000083041615159483019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c08201529291141590612ca560005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015612cfc5750808015612cfa5750612ced613881565b836020015163ffffffff16115b155b15612d33576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158015612d65575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15612d9c576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612da6613881565b905081612dbb57612db860328261594e565b90505b6000858152600460205260409020805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff90921691909117909155612e14906002908790613e5816565b5060105460808501516bffffffffffffffffffffffff9182169160009116821115612e79576080860151612e489083615aa5565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115612e79575060a08501515b808660a00151612e899190615aa5565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601154612eef9183911661598b565b601180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c6020526040902054163314612fcc576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811633141561301c576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600d602052604090205481169082161461064e5773ffffffffffffffffffffffffffffffffffffffff8281166000818152600d602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b6130d3613736565b6130dc8161408b565b50565b600060606000806000806130f1614181565b6000600f604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008a81526020019081526020016000206040518060e00160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff161515151581526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681525050905063ffffffff8016816020015163ffffffff161461340657505060408051602081019091526000808252965094506001935085915081905061372d565b80604001511561343557505060408051602081019091526000808252965094506002935085915081905061372d565b61343e826141b9565b825160125492965090945060009161347c9185917801000000000000000000000000000000000000000000000000900463ffffffff168888866143b5565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff1610156134c857600060405180602001604052806000815250600698509850985050505061372d565b5a60008b815260076020526040808220905192985090917f6e04ff0d000000000000000000000000000000000000000000000000000000009161350d91602401615710565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260608501516012549251919350600092839273ffffffffffffffffffffffffffffffffffffffff9092169163ffffffff909116906135c990869061544a565b60006040518083038160008787f1925050503d8060008114613607576040519150601f19603f3d011682016040523d82523d6000602084013e61360c565b606091505b50915091505a61361c908a615a8e565b98508161362c57600399506136c2565b80806020019051810190613640919061512d565b909c5090508b61366f5760006040518060200160405280600081525060049b509b509b5050505050505061372d565b6012548151780100000000000000000000000000000000000000000000000090910463ffffffff1610156136c25760006040518060200160405280600081525060059b509b509b5050505050505061372d565b604051806060016040528060016136d7613881565b6136e19190615a8e565b63ffffffff1681526020016136f96001610b1d613881565b8152602001828152506040516020016137129190615809565b6040516020818303038152906040529a50819b505050505050505b91939550919395565b60005473ffffffffffffffffffffffffffffffffffffffff1633146137b7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611106565b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613816576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902054640100000000900463ffffffff908116146130dc576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061387883836143fe565b90505b92915050565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156138b7576138b7615c11565b141561394157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561390457600080fd5b505afa158015613918573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061393c919061517e565b905090565b504390565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111561397c5761397c615c11565b1415613a0b576040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101839052606490632b407a829060240160206040518083038186803b1580156139d357600080fd5b505afa1580156139e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061387b919061517e565b504090565b919050565b600f546e010000000000000000000000000000900460ff1615613a64576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff86163b613ab2576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125482517401000000000000000000000000000000000000000090910463ffffffff161015613b0e576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8563ffffffff161080613b37575060125463ffffffff6401000000009091048116908616115b15613b6e576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000878152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1615613bd7576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518060e001604052808663ffffffff16815260200163ffffffff8016815260200182151581526020018773ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001846bffffffffffffffffffffffff168152602001600063ffffffff168152506004600089815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff160217905550905050836005600089815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550826bffffffffffffffffffffffff16601554613e20919061594e565b60155560008781526007602090815260409091208351613e4292850190614a3f565b50613e4e60028861386c565b5050505050505050565b6000613878838361444d565b73ffffffffffffffffffffffffffffffffffffffff831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e01000000000000000000000000000090920416606082018190528290613eef9086615aa5565b90506000613efd85836159cf565b90508083604001818151613f11919061598b565b6bffffffffffffffffffffffff9081169091528716606085015250613f368582615a63565b613f409083615aa5565b60118054600090613f609084906bffffffffffffffffffffffff1661598b565b825461010092830a6bffffffffffffffffffffffff81810219909216928216029190911790925573ffffffffffffffffffffffffffffffffffffffff999099166000908152600860209081526040918290208751815492890151938901516060909901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009093169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161760ff909316909b02919091177fffffffffffff000000000000000000000000000000000000000000000000ffff1662010000878416027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff16176e010000000000000000000000000000919092160217909755509095945050505050565b73ffffffffffffffffffffffffffffffffffffffff811633141561410b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611106565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b32156137b7576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561424057600080fd5b505afa158015614254573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061427891906152d0565b509450909250505060008113158061428f57508142105b806142b057508280156142b057506142a78242615a8e565b8463ffffffff16105b156142bf5760135495506142c3565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b15801561432957600080fd5b505afa15801561433d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061436191906152d0565b509450909250505060008113158061437857508142105b80614399575082801561439957506143908242615a8e565b8463ffffffff16105b156143a85760145494506143ac565b8094505b50505050915091565b6000806143c6868960000151614540565b90506000806143e18a8a63ffffffff16858a8a60018b614584565b90925090506143f0818361598b565b9a9950505050505050505050565b60008181526001830160205260408120546144455750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561387b565b50600061387b565b60008181526001830160205260408120548015614536576000614471600183615a8e565b855490915060009061448590600190615a8e565b90508181146144ea5760008660000182815481106144a5576144a5615c6f565b90600052602060002001549050808760000184815481106144c8576144c8615c6f565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806144fb576144fb615c40565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061387b565b600091505061387b565b600061455363ffffffff841660146159fa565b61455e836001615966565b61456d9060ff16611d4c6159fa565b61457a906201117061594e565b613878919061594e565b6000806000896080015161ffff168761459d91906159fa565b90508380156145ab5750803a105b156145b357503a5b600060027f000000000000000000000000000000000000000000000000000000000000000060028111156145e9576145e9615c11565b141561476c57604080516000815260208101909152851561464857600036604051806080016040528060488152602001615cdc6048913960405160200161463293929190615423565b60405160208183030381529060405290506146c4565b601254614678907801000000000000000000000000000000000000000000000000900463ffffffff166004615a37565b63ffffffff1667ffffffffffffffff81111561469657614696615c9e565b6040519080825280601f01601f1916602001820160405280156146c0576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e906147149084906004016156fd565b60206040518083038186803b15801561472c57600080fd5b505afa158015614740573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614764919061517e565b915050614828565b60017f000000000000000000000000000000000000000000000000000000000000000060028111156147a0576147a0615c11565b141561482857606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b1580156147ed57600080fd5b505afa158015614801573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614825919061517e565b90505b8461484457808b6080015161ffff1661484191906159fa565b90505b61485261ffff8716826159bb565b9050600087826148628c8e61594e565b61486c90866159fa565b614876919061594e565b61488890670de0b6b3a76400006159fa565b61489291906159bb565b905060008c6040015163ffffffff1664e8d4a510006148b191906159fa565b898e6020015163ffffffff16858f886148ca91906159fa565b6148d4919061594e565b6148e290633b9aca006159fa565b6148ec91906159fa565b6148f691906159bb565b614900919061594e565b90506b033b2e3c9fd0803ce8000000614919828461594e565b1115614951576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b82805461496f90615b02565b90600052602060002090601f01602090048101928261499157600085556149f5565b82601f106149c8578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008235161785556149f5565b828001600101855582156149f5579182015b828111156149f55782358255916020019190600101906149da565b50614a01929150614ab3565b5090565b508054614a1190615b02565b6000825580601f10614a21575050565b601f0160209004906000526020600020908101906130dc9190614ab3565b828054614a4b90615b02565b90600052602060002090601f016020900481019282614a6d57600085556149f5565b82601f10614a8657805160ff19168380011785556149f5565b828001600101855582156149f5579182015b828111156149f5578251825591602001919060010190614a98565b5b80821115614a015760008155600101614ab4565b803573ffffffffffffffffffffffffffffffffffffffff81168114613a1057600080fd5b60008083601f840112614afe57600080fd5b50813567ffffffffffffffff811115614b1657600080fd5b6020830191508360208260051b8501011115614b3157600080fd5b9250929050565b600082601f830112614b4957600080fd5b81356020614b5e614b59836158e4565b615895565b80838252828201915082860187848660051b8901011115614b7e57600080fd5b60005b85811015614ba457614b9282614ac8565b84529284019290840190600101614b81565b5090979650505050505050565b600082601f830112614bc257600080fd5b81356020614bd2614b59836158e4565b80838252828201915082860187848660051b8901011115614bf257600080fd5b60005b85811015614ba457813567ffffffffffffffff811115614c1457600080fd5b8801603f81018a13614c2557600080fd5b858101356040614c37614b5983615908565b8281528c82848601011115614c4b57600080fd5b828285018a8301376000928101890192909252508552509284019290840190600101614bf5565b600082601f830112614c8357600080fd5b81356020614c93614b59836158e4565b8281528181019085830160e080860288018501891015614cb257600080fd5b6000805b87811015614d575782848c031215614ccc578182fd5b614cd461586c565b614cdd85614df5565b8152614cea888601614df5565b88820152604080860135614cfd81615ccd565b908201526060614d0e868201614ac8565b908201526080614d1f868201614e23565b9082015260a0614d30868201614e23565b9082015260c0614d41868201614df5565b9082015286529486019492820192600101614cb6565b50929998505050505050505050565b60008083601f840112614d7857600080fd5b50813567ffffffffffffffff811115614d9057600080fd5b602083019150836020828501011115614b3157600080fd5b600082601f830112614db957600080fd5b8151614dc7614b5982615908565b818152846020838601011115614ddc57600080fd5b614ded826020830160208701615ad2565b949350505050565b803563ffffffff81168114613a1057600080fd5b805169ffffffffffffffffffff81168114613a1057600080fd5b80356bffffffffffffffffffffffff81168114613a1057600080fd5b600060208284031215614e5157600080fd5b61387882614ac8565b60008060408385031215614e6d57600080fd5b614e7683614ac8565b9150614e8460208401614ac8565b90509250929050565b60008060408385031215614ea057600080fd5b614ea983614ac8565b9150602083013560048110614ebd57600080fd5b809150509250929050565b600080600080600080600060a0888a031215614ee357600080fd5b614eec88614ac8565b9650614efa60208901614df5565b9550614f0860408901614ac8565b9450606088013567ffffffffffffffff80821115614f2557600080fd5b614f318b838c01614d66565b909650945060808a0135915080821115614f4a57600080fd5b50614f578a828b01614d66565b989b979a50959850939692959293505050565b60008060208385031215614f7d57600080fd5b823567ffffffffffffffff811115614f9457600080fd5b614fa085828601614aec565b90969095509350505050565b600080600060408486031215614fc157600080fd5b833567ffffffffffffffff811115614fd857600080fd5b614fe486828701614aec565b9094509250614ff7905060208501614ac8565b90509250925092565b6000806000806080858703121561501657600080fd5b843567ffffffffffffffff8082111561502e57600080fd5b818701915087601f83011261504257600080fd5b81356020615052614b59836158e4565b8083825282820191508286018c848660051b890101111561507257600080fd5b600096505b84871015615095578035835260019690960195918301918301615077565b50985050880135925050808211156150ac57600080fd5b6150b888838901614c72565b945060408701359150808211156150ce57600080fd5b6150da88838901614bb1565b935060608701359150808211156150f057600080fd5b506150fd87828801614b38565b91505092959194509250565b60006020828403121561511b57600080fd5b815161512681615ccd565b9392505050565b6000806040838503121561514057600080fd5b825161514b81615ccd565b602084015190925067ffffffffffffffff81111561516857600080fd5b61517485828601614da8565b9150509250929050565b60006020828403121561519057600080fd5b5051919050565b600080602083850312156151aa57600080fd5b823567ffffffffffffffff8111156151c157600080fd5b614fa085828601614d66565b6000602082840312156151df57600080fd5b815167ffffffffffffffff8111156151f657600080fd5b614ded84828501614da8565b60006020828403121561521457600080fd5b5035919050565b6000806040838503121561522e57600080fd5b82359150614e8460208401614ac8565b60008060006040848603121561525357600080fd5b83359250602084013567ffffffffffffffff81111561527157600080fd5b61527d86828701614d66565b9497909650939450505050565b6000806040838503121561529d57600080fd5b82359150614e8460208401614df5565b600080604083850312156152c057600080fd5b82359150614e8460208401614e23565b600080600080600060a086880312156152e857600080fd5b6152f186614e09565b945060208601519350604086015192506060860151915061531460808701614e09565b90509295509295909350565b60006020828403121561533257600080fd5b815160ff8116811461512657600080fd5b600081518084526020808501945080840160005b8381101561538957815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101615357565b509495945050505050565b6000815180845260208085019450848260051b860182860160005b85811015614ba45783830389526153c78383516153d9565b988501989250908401906001016153af565b600081518084526153f1816020860160208601615ad2565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b828482376000838201600081528351615440818360208801615ad2565b0195945050505050565b6000825161545c818460208701615ad2565b9190910192915050565b6000604082016040835280865480835260608501915087600052602092508260002060005b828110156154bd57815473ffffffffffffffffffffffffffffffffffffffff168452928401926001918201910161548b565b505050838103828501528481528590820160005b8681101561550a5773ffffffffffffffffffffffffffffffffffffffff6154f784614ac8565b16825291830191908301906001016154d1565b50979650505050505050565b60006080808352868184015260a07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff88111561555157600080fd5b8760051b808a838701378085019050818101600081526020838784030181880152818a5180845260c093508385019150828c01945060005b81811015615627578551805163ffffffff908116855285820151168585015260408082015115159085015260608082015173ffffffffffffffffffffffffffffffffffffffff1690850152888101516bffffffffffffffffffffffff168985015287810151615607898601826bffffffffffffffffffffffff169052565b5085015163ffffffff16838601529483019460e090920191600101615589565b5050878103604089015261563b818b615394565b9550505050505082810360608401526156548185615343565b98975050505050505050565b861515815260c06020820152600061567b60c08301886153d9565b90506007861061568d5761568d615c11565b8560408301528460608301528360808301528260a0830152979650505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60208152600061387860208301846153d9565b600060208083526000845481600182811c91508083168061573257607f831692505b858310811415615769577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b87860183815260200181801561578657600181146157b5576157e0565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008616825287820196506157e0565b60008b81526020902060005b868110156157da578154848201529085019089016157c1565b83019750505b50949998505050505050505050565b602081016003831061580357615803615c11565b91905290565b6020815263ffffffff82511660208201526020820151604082015260006040830151606080840152614ded60808401826153d9565b60ff8416815260ff8316602082015260606040820152600061586360608301846153d9565b95945050505050565b60405160e0810167ffffffffffffffff8111828210171561588f5761588f615c9e565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156158dc576158dc615c9e565b604052919050565b600067ffffffffffffffff8211156158fe576158fe615c9e565b5060051b60200190565b600067ffffffffffffffff82111561592257615922615c9e565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000821982111561596157615961615bb3565b500190565b600060ff821660ff84168060ff0382111561598357615983615bb3565b019392505050565b60006bffffffffffffffffffffffff8083168185168083038211156159b2576159b2615bb3565b01949350505050565b6000826159ca576159ca615be2565b500490565b60006bffffffffffffffffffffffff808416806159ee576159ee615be2565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615a3257615a32615bb3565b500290565b600063ffffffff80831681851681830481118215151615615a5a57615a5a615bb3565b02949350505050565b60006bffffffffffffffffffffffff80831681851681830481118215151615615a5a57615a5a615bb3565b600082821015615aa057615aa0615bb3565b500390565b60006bffffffffffffffffffffffff83811690831681811015615aca57615aca615bb3565b039392505050565b60005b83811015615aed578181015183820152602001615ad5565b83811115615afc576000848401525b50505050565b600181811c90821680615b1657607f821691505b60208210811415615b50577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615b8857615b88615bb3565b5060010190565b600063ffffffff80831681811415615ba957615ba9615bb3565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b80151581146130dc57600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
}

var KeeperRegistryLogicABI = KeeperRegistryLogicMetaData.ABI

var KeeperRegistryLogicBin = KeeperRegistryLogicMetaData.Bin

func DeployKeeperRegistryLogic(auth *bind.TransactOpts, backend bind.ContractBackend, mode uint8, link common.Address, linkNativeFeed common.Address, fastGasFeed common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogic, error) {
	parsed, err := KeeperRegistryLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogicBin), backend, mode, link, linkNativeFeed, fastGasFeed)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.GetLinkAddress(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.GetLinkAddress(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) GetMode() (uint8, error) {
	return _KeeperRegistryLogic.Contract.GetMode(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) GetMode() (uint8, error) {
	return _KeeperRegistryLogic.Contract.GetMode(&_KeeperRegistryLogic.CallOpts)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptOwnership(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptOwnership(&_KeeperRegistryLogic.TransactOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptPayeeship(&_KeeperRegistryLogic.TransactOpts, transmitter)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.AcceptPayeeship(&_KeeperRegistryLogic.TransactOpts, transmitter)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "checkUpkeep", id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.CheckUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.CheckUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "pauseUpkeep", id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.PauseUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.PauseUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.RegisterUpkeep(&_KeeperRegistryLogic.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.RegisterUpkeep(&_KeeperRegistryLogic.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "setPayees", payees)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetPayees(&_KeeperRegistryLogic.TransactOpts, payees)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetPayees(&_KeeperRegistryLogic.TransactOpts, payees)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogic.TransactOpts, id, config)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogic.TransactOpts, id, config)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferPayeeship(&_KeeperRegistryLogic.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.TransferPayeeship(&_KeeperRegistryLogic.TransactOpts, transmitter, proposed)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.UnpauseUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.UnpauseUpkeep(&_KeeperRegistryLogic.TransactOpts, id)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.UpdateCheckData(&_KeeperRegistryLogic.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicTransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic.Contract.UpdateCheckData(&_KeeperRegistryLogic.TransactOpts, id, newCheckData)
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

type KeeperRegistryLogicCancelledUpkeepReportIterator struct {
	Event *KeeperRegistryLogicCancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicCancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicCancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicCancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicCancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicCancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicCancelledUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicCancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicCancelledUpkeepReportIterator{contract: _KeeperRegistryLogic.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicCancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicCancelledUpkeepReport)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicCancelledUpkeepReport, error) {
	event := new(KeeperRegistryLogicCancelledUpkeepReport)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
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

type KeeperRegistryLogicInsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryLogicInsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicInsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicInsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicInsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicInsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicInsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicInsufficientFundsUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicInsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicInsufficientFundsUpkeepReportIterator{contract: _KeeperRegistryLogic.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicInsufficientFundsUpkeepReport)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicInsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryLogicInsufficientFundsUpkeepReport)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
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

type KeeperRegistryLogicPayeesUpdatedIterator struct {
	Event *KeeperRegistryLogicPayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicPayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicPayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicPayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicPayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicPayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicPayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicPayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPayeesUpdatedIterator{contract: _KeeperRegistryLogic.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicPayeesUpdated)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicPayeesUpdated, error) {
	event := new(KeeperRegistryLogicPayeesUpdated)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
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
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPayeeshipTransferRequestedIterator{contract: _KeeperRegistryLogic.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
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
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPayeeshipTransferredIterator{contract: _KeeperRegistryLogic.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
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
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicPaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicPaymentWithdrawnIterator{contract: _KeeperRegistryLogic.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
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

type KeeperRegistryLogicReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryLogicReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicReorgedUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicReorgedUpkeepReportIterator{contract: _KeeperRegistryLogic.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicReorgedUpkeepReport)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicReorgedUpkeepReport, error) {
	event := new(KeeperRegistryLogicReorgedUpkeepReport)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogicStaleUpkeepReportIterator struct {
	Event *KeeperRegistryLogicStaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicStaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicStaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicStaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicStaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicStaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicStaleUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicStaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicStaleUpkeepReportIterator{contract: _KeeperRegistryLogic.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicStaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicStaleUpkeepReport)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicStaleUpkeepReport, error) {
	event := new(KeeperRegistryLogicStaleUpkeepReport)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
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

type KeeperRegistryLogicUpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogicUpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogicUpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogicUpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogicUpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogicUpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogicUpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogicUpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepOffchainConfigSetIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogicUpkeepOffchainConfigSet)
				if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicUpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryLogicUpkeepOffchainConfigSet)
	if err := _KeeperRegistryLogic.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
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
	Id               *big.Int
	Success          bool
	CheckBlockNumber uint32
	GasUsed          *big.Int
	GasOverhead      *big.Int
	TotalPayment     *big.Int
	Raw              types.Log
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicUpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogicUpkeepPerformedIterator{contract: _KeeperRegistryLogic.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic *KeeperRegistryLogicFilterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogic.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
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
	case _KeeperRegistryLogic.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistryLogic.ParseCancelledUpkeepReport(log)
	case _KeeperRegistryLogic.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryLogic.ParseFundsAdded(log)
	case _KeeperRegistryLogic.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryLogic.ParseFundsWithdrawn(log)
	case _KeeperRegistryLogic.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistryLogic.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistryLogic.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistryLogic.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistryLogic.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryLogic.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryLogic.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryLogic.ParseOwnershipTransferred(log)
	case _KeeperRegistryLogic.abi.Events["Paused"].ID:
		return _KeeperRegistryLogic.ParsePaused(log)
	case _KeeperRegistryLogic.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistryLogic.ParsePayeesUpdated(log)
	case _KeeperRegistryLogic.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryLogic.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryLogic.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryLogic.ParsePayeeshipTransferred(log)
	case _KeeperRegistryLogic.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryLogic.ParsePaymentWithdrawn(log)
	case _KeeperRegistryLogic.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistryLogic.ParseReorgedUpkeepReport(log)
	case _KeeperRegistryLogic.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistryLogic.ParseStaleUpkeepReport(log)
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
	case _KeeperRegistryLogic.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistryLogic.ParseUpkeepOffchainConfigSet(log)
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

func (KeeperRegistryLogicCancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f")
}

func (KeeperRegistryLogicFundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogicFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogicInsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96")
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

func (KeeperRegistryLogicPayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
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

func (KeeperRegistryLogicReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13")
}

func (KeeperRegistryLogicStaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89")
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

func (KeeperRegistryLogicUpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryLogicUpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryLogicUpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420")
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
	GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkAddress(opts *bind.CallOpts) (common.Address, error)

	GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error)

	GetMode(opts *bind.CallOpts) (uint8, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error)

	SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)

	UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)

	WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicCancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicCancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogicCancelledUpkeepReport, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogicFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicFundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryLogicFundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicFundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicFundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogicFundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicInsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicInsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogicInsufficientFundsUpkeepReport, error)

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

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogicPayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogicPayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogicPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogicPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogicPayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogicPaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicPaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogicPaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogicReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicStaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicStaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogicStaleUpkeepReport, error)

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

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogicUpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogicUpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogicUpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogicUpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogicUpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

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
