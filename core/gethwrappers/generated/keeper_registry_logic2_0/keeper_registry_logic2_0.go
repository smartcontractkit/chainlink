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

var KeeperRegistryLogic20MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mode\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_0.Mode\"},{\"name\":\"link\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"linkNativeFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fastGasFeed\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"upkeepNeeded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"performData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"upkeepFailureReason\",\"type\":\"uint8\",\"internalType\":\"enumUpkeepFailureReason\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fastGasWei\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"linkNative\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getFastGasFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkNativeFeedAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMode\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_0.Mode\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateUpkeeps\",\"inputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"destination\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"receiveUpkeeps\",\"inputs\":[{\"name\":\"encodedUpkeeps\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recoverFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerUpkeep\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"checkData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPayees\",\"inputs\":[{\"name\":\"payees\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPeerRegistryMigrationPermission\",\"inputs\":[{\"name\":\"peer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permission\",\"type\":\"uint8\",\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepGasLimit\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpkeepOffchainConfig\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"config\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferUpkeepAdmin\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpauseUpkeep\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCheckData\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawOwnerFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPayment\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CancelledUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnerFundsWithdrawn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeesUpdated\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentWithdrawn\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReorgedUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleUpkeepReport\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepAdminTransferred\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCanceled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"atBlockHeight\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepCheckDataUpdated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"newCheckData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepGasLimitSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"gasLimit\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepMigrated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destination\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepOffchainConfigSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepPerformed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"},{\"name\":\"checkBlockNumber\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasOverhead\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalPayment\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepReceived\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startingBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"importedFrom\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepRegistered\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"executeGas\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpkeepUnpaused\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayHasNoEntries\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotCancel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckDataExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateEntry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasLimitOutsideRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfFaultyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectNumberOfSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MigrationNotPermitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAContract\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnchainConfigNonEmpty\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveSigners\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyActiveTransmitters\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByLINKToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByProposedPayee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySimulatedBackend\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyUnpausedUpkeep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ParameterLengthError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentGreaterThanAllLINK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RegistryPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TargetCheckReverted\",\"inputs\":[{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TranscoderNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotCanceled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UpkeepNotNeeded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueNotChanged\",\"inputs\":[]}]",
	Bin: "0x6101006040523480156200001257600080fd5b50604051620060e0380380620060e08339810160408190526200003591620001ef565b838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c38162000126565b505050836002811115620000db57620000db62000251565b60e0816002811115620000f257620000f262000251565b60f81b9052506001600160601b0319606093841b811660805291831b821660a05290911b1660c05250620002679350505050565b6001600160a01b038116331415620001815760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ea57600080fd5b919050565b600080600080608085870312156200020657600080fd5b8451600381106200021657600080fd5b93506200022660208601620001d2565b92506200023660408601620001d2565b91506200024660608601620001d2565b905092959194509250565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160f81c615dd86200030860003960008181610224015281816138850152818161394a01528181614668015261481f01526000818161026e015261428b0152600081816103b7015261437401526000818161041e01528181610f7e0152818161126301528181611bf4015281816122670152818161259c01528181612a830152612b160152615dd86000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80638dcf0fe711610104578063b121e147116100a2578063ca30e60311610071578063ca30e6031461041c578063eb5dcd6c14610442578063f2fde38b14610455578063f7d334ba1461046857600080fd5b8063b121e147146103db578063b148ab6b146103ee578063b79550be14610401578063c80480221461040957600080fd5b80639fab4386116100de5780639fab43861461037c578063a710b2211461038f578063a72aa27e146103a2578063b10b673c146103b557600080fd5b80638dcf0fe7146103435780638e86139b14610356578063948108f71461036957600080fd5b80636ded9eae1161017c5780638456cb591161014b5780638456cb59146102f757806385c1b0ba146102ff5780638765ecbe146103125780638da5cb5b1461032557600080fd5b80636ded9eae146102b3578063744bfe61146102d457806379ba5097146102e75780637d9b97e0146102ef57600080fd5b80633f4ba83a116101b85780633f4ba83a1461021a5780634b4fd03b146102225780635165f2f5146102595780636709d0e51461026c57600080fd5b8063187256e8146101df5780631a2af011146101f45780633b9cce5914610207575b600080fd5b6101f26101ed366004614f3c565b61048d565b005b6101f26102023660046152c3565b6104fe565b6101f2610215366004615019565b610652565b6101f26108a8565b7f00000000000000000000000000000000000000000000000000000000000000006040516102509190615897565b60405180910390f35b6101f26102673660046152aa565b61090e565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610250565b6102c66102c1366004614f77565b610a88565b604051908152602001610250565b6101f26102e23660046152c3565b610c77565b6101f2611089565b6101f261118b565b6101f26112f5565b6101f261030d36600461505b565b611366565b6101f26103203660046152aa565b611c7f565b60005473ffffffffffffffffffffffffffffffffffffffff1661028e565b6101f26103513660046152e6565b611e06565b6101f261036436600461523f565b611e68565b6101f2610377366004615355565b6120a4565b6101f261038a3660046152e6565b612343565b6101f261039d366004614f09565b6123f2565b6101f26103b0366004615332565b61267d565b7f000000000000000000000000000000000000000000000000000000000000000061028e565b6101f26103e9366004614eee565b61275f565b6101f26103fc3660046152aa565b612857565b6101f2612a4a565b6101f26104173660046152aa565b612bb5565b7f000000000000000000000000000000000000000000000000000000000000000061028e565b6101f2610450366004614f09565b612f6a565b6101f2610463366004614eee565b6130c9565b61047b6104763660046152aa565b6130dd565b60405161025096959493929190615708565b610495613734565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260166020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660018360038111156104f5576104f5615cb9565b02179055505050565b610507826137b7565b73ffffffffffffffffffffffffffffffffffffffff8116331415610557576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166105a4576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff82811691161461064e5760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b61065a613734565b600b548114610695576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600b54811015610867576000600b82815481106106b7576106b7615d17565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600c9092526040832054919350169085858581811061070157610701615d17565b90506020020160208101906107169190614eee565b905073ffffffffffffffffffffffffffffffffffffffff811615806107a9575073ffffffffffffffffffffffffffffffffffffffff82161580159061078757508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156107a9575073ffffffffffffffffffffffffffffffffffffffff81811614155b156107e0576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff818116146108515773ffffffffffffffffffffffffffffffffffffffff8381166000908152600c6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061085f90615bfe565b915050610698565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600b838360405161089c9392919061550e565b60405180910390a15050565b6108b0613734565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b610917816137b7565b600081815260046020908152604091829020825160e081018452815463ffffffff8082168352640100000000820481169483019490945268010000000000000000810460ff1615159482018590526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c082015290610a19576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169055610a5860028361386a565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590610ad957506011546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314155b15610b10576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b2c6001610b1d61387f565b610b279190615b36565b613944565b601254604080516020810193909352309083015268010000000000000000900463ffffffff1660608201526080016040516020818303038152906040528051906020012060001c9050610bba8189898960008a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509250613ac4915050565b6012805468010000000000000000900463ffffffff16906008610bdc83615c37565b825463ffffffff9182166101009390930a9283029190920219909116179055506000818152601760205260409020610c15908484614a12565b506040805163ffffffff8916815273ffffffffffffffffffffffffffffffffffffffff8816602082015282917fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012910160405180910390a2979650505050505050565b600f546f01000000000000000000000000000000900460ff1615610cc7576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116610d4e576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815160e081018352815463ffffffff8082168352640100000000820481168387015260ff6801000000000000000083041615158386015273ffffffffffffffffffffffffffffffffffffffff6901000000000000000000909204821660608401526001909301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490921660c082015286855260059093529220549091163314610e5b576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e6361387f565b816020015163ffffffff161115610ea6576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546015546c010000000000000000000000009091046bffffffffffffffffffffffff1690610ee6908290615b36565b60155560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b158015610fc457600080fd5b505af1158015610fd8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ffc91906151b8565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461110f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611193613734565b6011546015546bffffffffffffffffffffffff909116906111b5908290615b36565b601555601180547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b602060405180830381600087803b1580156112bd57600080fd5b505af11580156112d1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061064e91906151b8565b6112fd613734565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff166e0100000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602001610904565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156113a2576113a2615cb9565b141580156113ea5750600373ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156113e7576113e7615cb9565b14155b15611421576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16611480576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816114b7576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff81111561150b5761150b615d46565b60405190808252806020026020018201604052801561153e57816020015b60608152602001906001900390816115295790505b50905060008667ffffffffffffffff81111561155c5761155c615d46565b604051908082528060200260200182016040528015611585578160200160208202803683370190505b50905060008767ffffffffffffffff8111156115a3576115a3615d46565b60405190808252806020026020018201604052801561162857816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816115c15790505b50905060005b888110156119b35789898281811061164857611648615d17565b60209081029290920135600081815260048452604090819020815160e081018352815463ffffffff8082168352640100000000820481169783019790975268010000000000000000810460ff16151593820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c0840152985090965061172a9050876137b7565b8582828151811061173d5761173d615d17565b602002602001018190525060076000888152602001908152602001600020805461176690615baa565b80601f016020809104026020016040519081016040528092919081815260200182805461179290615baa565b80156117df5780601f106117b4576101008083540402835291602001916117df565b820191906000526020600020905b8154815290600101906020018083116117c257829003601f168201915b50505050508482815181106117f6576117f6615d17565b60200260200101819052506005600088815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683828151811061184757611847615d17565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015260a0860151611889906bffffffffffffffffffffffff16866159f6565b600088815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff00000000000000000000000000000000000000000000000000000000169055600790915281209196506118fd9190614ab4565b600087815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905561193c600288613f07565b5060a0860151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8a16602083015288917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a2806119ab81615bfe565b91505061162e565b50836015546119c29190615b36565b6015556040516000906119e1908b908b908590889088906020016155be565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260105490915073ffffffffffffffffffffffffffffffffffffffff808a1691638e86139b916c010000000000000000000000009091041663c71249ab60028c73ffffffffffffffffffffffffffffffffffffffff1663aab9edd66040518163ffffffff1660e01b815260040160206040518083038186803b158015611a9357600080fd5b505afa158015611aa7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611acb91906153c8565b866040518463ffffffff1660e01b8152600401611aea939291906158e6565b60006040518083038186803b158015611b0257600080fd5b505afa158015611b16573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b5c9190810190615275565b6040518263ffffffff1660e01b8152600401611b7891906157a5565b600060405180830381600087803b158015611b9257600080fd5b505af1158015611ba6573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8b81166004830152602482018990527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b158015611c3a57600080fd5b505af1158015611c4e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c7291906151b8565b5050505050505050505050565b611c88816137b7565b600081815260046020908152604091829020825160e081018452815463ffffffff8082168352640100000000820481169483019490945268010000000000000000810460ff16158015958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c082015290611d8c576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1668010000000000000000179055611dd6600283613f07565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b611e0f836137b7565b6000838152601760205260409020611e28908383614a12565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508383604051611e5b929190615758565b60405180910390a2505050565b60023360009081526016602052604090205460ff166003811115611e8e57611e8e615cb9565b14158015611ec0575060033360009081526016602052604090205460ff166003811115611ebd57611ebd615cb9565b14155b15611ef7576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808080611f08858701876150af565b935093509350935060005b845181101561209b57611fea858281518110611f3157611f31615d17565b6020026020010151858381518110611f4b57611f4b615d17565b602002602001015160600151868481518110611f6957611f69615d17565b602002602001015160000151858581518110611f8757611f87615d17565b6020026020010151888681518110611fa157611fa1615d17565b602002602001015160a00151888781518110611fbf57611fbf615d17565b60200260200101518a8881518110611fd957611fd9615d17565b602002602001015160400151613ac4565b848181518110611ffc57611ffc615d17565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7185838151811061203757612037615d17565b602002602001015160a00151336040516120819291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a28061209381615bfe565b915050611f13565b50505050505050565b600082815260046020908152604091829020825160e081018452815463ffffffff80821683526401000000008204811694830185905268010000000000000000820460ff161515958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c082015291146121a6576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a001516121b69190615a33565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff9384160217905560155461221c918416906159f6565b6015556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd90606401602060405180830381600087803b1580156122c057600080fd5b505af11580156122d4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122f891906151b8565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b61234c836137b7565b60125474010000000000000000000000000000000000000000900463ffffffff168111156123a6576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526007602052604090206123bf908383614a12565b50827f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf8383604051611e5b929190615758565b73ffffffffffffffffffffffffffffffffffffffff811661243f576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c602052604090205416331461249f576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f54600b546000916124d691859170010000000000000000000000000000000090046bffffffffffffffffffffffff1690613f13565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260086020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601554909150612540906bffffffffffffffffffffffff831690615b36565b6015556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b1580156125e057600080fd5b505af11580156125f4573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061261891906151b8565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff1610806126a6575060125463ffffffff6401000000009091048116908216115b156126dd576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6126e6826137b7565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8516908117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600d60205260409020541633146127bf576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600c602090815260408083208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217909355600d909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600081815260046020908152604091829020825160e081018452815463ffffffff80821683526401000000008204811694830185905268010000000000000000820460ff161515958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529114612959576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff1633146129b6576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b612a52613734565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b158015612ada57600080fd5b505afa158015612aee573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b129190615226565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb3360155484612b5f9190615b36565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff909216600483015260248201526044016112a3565b6000818152600460209081526040808320815160e081018352815463ffffffff80821683526401000000008204811695830186905260ff6801000000000000000083041615159483019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c08201529291141590612ca360005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015612cfa5750808015612cf85750612ceb61387f565b836020015163ffffffff16115b155b15612d31576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158015612d63575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15612d9a576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612da461387f565b905081612db957612db66032826159f6565b90505b6000858152600460205260409020805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff90921691909117909155612e12906002908790613f0716565b5060105460808501516bffffffffffffffffffffffff9182169160009116821115612e77576080860151612e469083615b4d565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115612e77575060a08501515b808660a00151612e879190615b4d565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601154612eed91839116615a33565b601180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c6020526040902054163314612fca576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff811633141561301a576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600d602052604090205481169082161461064e5773ffffffffffffffffffffffffffffffffffffffff8281166000818152600d602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b6130d1613734565b6130da8161413a565b50565b600060606000806000806130ef614230565b6000600f604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008a81526020019081526020016000206040518060e00160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff161515151581526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681525050905063ffffffff8016816020015163ffffffff161461340457505060408051602081019091526000808252965094506001935085915081905061372b565b80604001511561343357505060408051602081019091526000808252965094506002935085915081905061372b565b61343c82614268565b825160125492965090945060009161347a9185917801000000000000000000000000000000000000000000000000900463ffffffff16888886614464565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff1610156134c657600060405180602001604052806000815250600698509850985050505061372b565b5a60008b815260076020526040808220905192985090917f6e04ff0d000000000000000000000000000000000000000000000000000000009161350b916024016157b8565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260608501516012549251919350600092839273ffffffffffffffffffffffffffffffffffffffff9092169163ffffffff909116906135c79086906154f2565b60006040518083038160008787f1925050503d8060008114613605576040519150601f19603f3d011682016040523d82523d6000602084013e61360a565b606091505b50915091505a61361a908a615b36565b98508161362a57600399506136c0565b8080602001905181019061363e91906151d5565b909c5090508b61366d5760006040518060200160405280600081525060049b509b509b5050505050505061372b565b6012548151780100000000000000000000000000000000000000000000000090910463ffffffff1610156136c05760006040518060200160405280600081525060059b509b509b5050505050505061372b565b604051806060016040528060016136d561387f565b6136df9190615b36565b63ffffffff1681526020016136f76001610b1d61387f565b81526020018281525060405160200161371091906158b1565b6040516020818303038152906040529a50819b505050505050505b91939550919395565b60005473ffffffffffffffffffffffffffffffffffffffff1633146137b5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611106565b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613814576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902054640100000000900463ffffffff908116146130da576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061387683836144ad565b90505b92915050565b600060017f000000000000000000000000000000000000000000000000000000000000000060028111156138b5576138b5615cb9565b141561393f57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561390257600080fd5b505afa158015613916573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061393a9190615226565b905090565b504390565b600060017f0000000000000000000000000000000000000000000000000000000000000000600281111561397a5761397a615cb9565b1415613aba576000606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156139c957600080fd5b505afa1580156139dd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a019190615226565b90508083101580613a1c5750610100613a1a8483615b36565b115b15613a2a5750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815260048101849052606490632b407a829060240160206040518083038186803b158015613a7b57600080fd5b505afa158015613a8f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ab39190615226565b9392505050565b504090565b919050565b600f546e010000000000000000000000000000900460ff1615613b13576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff86163b613b61576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125482517401000000000000000000000000000000000000000090910463ffffffff161015613bbd576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8563ffffffff161080613be6575060125463ffffffff6401000000009091048116908616115b15613c1d576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000878152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1615613c86576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518060e001604052808663ffffffff16815260200163ffffffff8016815260200182151581526020018773ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001846bffffffffffffffffffffffff168152602001600063ffffffff168152506004600089815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff160217905550905050836005600089815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550826bffffffffffffffffffffffff16601554613ecf91906159f6565b60155560008781526007602090815260409091208351613ef192850190614aee565b50613efd60028861386a565b5050505050505050565b600061387683836144fc565b73ffffffffffffffffffffffffffffffffffffffff831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e01000000000000000000000000000090920416606082018190528290613f9e9086615b4d565b90506000613fac8583615a77565b90508083604001818151613fc09190615a33565b6bffffffffffffffffffffffff9081169091528716606085015250613fe58582615b0b565b613fef9083615b4d565b6011805460009061400f9084906bffffffffffffffffffffffff16615a33565b825461010092830a6bffffffffffffffffffffffff81810219909216928216029190911790925573ffffffffffffffffffffffffffffffffffffffff999099166000908152600860209081526040918290208751815492890151938901516060909901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009093169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161760ff909316909b02919091177fffffffffffff000000000000000000000000000000000000000000000000ffff1662010000878416027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff16176e010000000000000000000000000000919092160217909755509095945050505050565b73ffffffffffffffffffffffffffffffffffffffff81163314156141ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611106565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b32156137b5576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156142ef57600080fd5b505afa158015614303573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906143279190615378565b509450909250505060008113158061433e57508142105b8061435f575082801561435f57506143568242615b36565b8463ffffffff16105b1561436e576013549550614372565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156143d857600080fd5b505afa1580156143ec573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144109190615378565b509450909250505060008113158061442757508142105b806144485750828015614448575061443f8242615b36565b8463ffffffff16105b1561445757601454945061445b565b8094505b50505050915091565b6000806144758689600001516145ef565b90506000806144908a8a63ffffffff16858a8a60018b614633565b909250905061449f8183615a33565b9a9950505050505050505050565b60008181526001830160205260408120546144f457508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155613879565b506000613879565b600081815260018301602052604081205480156145e5576000614520600183615b36565b855490915060009061453490600190615b36565b905081811461459957600086600001828154811061455457614554615d17565b906000526020600020015490508087600001848154811061457757614577615d17565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806145aa576145aa615ce8565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050613879565b6000915050613879565b600061460263ffffffff84166014615aa2565b61460d836001615a0e565b61461c9060ff16611d4c615aa2565b61462990620111706159f6565b61387691906159f6565b6000806000896080015161ffff168761464c9190615aa2565b905083801561465a5750803a105b1561466257503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561469857614698615cb9565b141561481b5760408051600081526020810190915285156146f757600036604051806080016040528060488152602001615d84604891396040516020016146e1939291906154cb565b6040516020818303038152906040529050614773565b601254614727907801000000000000000000000000000000000000000000000000900463ffffffff166004615adf565b63ffffffff1667ffffffffffffffff81111561474557614745615d46565b6040519080825280601f01601f19166020018201604052801561476f576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e906147c39084906004016157a5565b60206040518083038186803b1580156147db57600080fd5b505afa1580156147ef573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906148139190615226565b9150506148d7565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111561484f5761484f615cb9565b14156148d757606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561489c57600080fd5b505afa1580156148b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906148d49190615226565b90505b846148f357808b6080015161ffff166148f09190615aa2565b90505b61490161ffff871682615a63565b9050600087826149118c8e6159f6565b61491b9086615aa2565b61492591906159f6565b61493790670de0b6b3a7640000615aa2565b6149419190615a63565b905060008c6040015163ffffffff1664e8d4a510006149609190615aa2565b898e6020015163ffffffff16858f886149799190615aa2565b61498391906159f6565b61499190633b9aca00615aa2565b61499b9190615aa2565b6149a59190615a63565b6149af91906159f6565b90506b033b2e3c9fd0803ce80000006149c882846159f6565b1115614a00576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b828054614a1e90615baa565b90600052602060002090601f016020900481019282614a405760008555614aa4565b82601f10614a77578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555614aa4565b82800160010185558215614aa4579182015b82811115614aa4578235825591602001919060010190614a89565b50614ab0929150614b62565b5090565b508054614ac090615baa565b6000825580601f10614ad0575050565b601f0160209004906000526020600020908101906130da9190614b62565b828054614afa90615baa565b90600052602060002090601f016020900481019282614b1c5760008555614aa4565b82601f10614b3557805160ff1916838001178555614aa4565b82800160010185558215614aa4579182015b82811115614aa4578251825591602001919060010190614b47565b5b80821115614ab05760008155600101614b63565b803573ffffffffffffffffffffffffffffffffffffffff81168114613abf57600080fd5b60008083601f840112614bad57600080fd5b50813567ffffffffffffffff811115614bc557600080fd5b6020830191508360208260051b8501011115614be057600080fd5b9250929050565b600082601f830112614bf857600080fd5b81356020614c0d614c088361598c565b61593d565b80838252828201915082860187848660051b8901011115614c2d57600080fd5b60005b85811015614c5357614c4182614b77565b84529284019290840190600101614c30565b5090979650505050505050565b600082601f830112614c7157600080fd5b81356020614c81614c088361598c565b80838252828201915082860187848660051b8901011115614ca157600080fd5b60005b85811015614c5357813567ffffffffffffffff811115614cc357600080fd5b8801603f81018a13614cd457600080fd5b858101356040614ce6614c08836159b0565b8281528c82848601011115614cfa57600080fd5b828285018a8301376000928101890192909252508552509284019290840190600101614ca4565b600082601f830112614d3257600080fd5b81356020614d42614c088361598c565b8281528181019085830160e080860288018501891015614d6157600080fd5b6000805b87811015614e065782848c031215614d7b578182fd5b614d83615914565b614d8c85614ea4565b8152614d99888601614ea4565b88820152604080860135614dac81615d75565b908201526060614dbd868201614b77565b908201526080614dce868201614ed2565b9082015260a0614ddf868201614ed2565b9082015260c0614df0868201614ea4565b9082015286529486019492820192600101614d65565b50929998505050505050505050565b60008083601f840112614e2757600080fd5b50813567ffffffffffffffff811115614e3f57600080fd5b602083019150836020828501011115614be057600080fd5b600082601f830112614e6857600080fd5b8151614e76614c08826159b0565b818152846020838601011115614e8b57600080fd5b614e9c826020830160208701615b7a565b949350505050565b803563ffffffff81168114613abf57600080fd5b805169ffffffffffffffffffff81168114613abf57600080fd5b80356bffffffffffffffffffffffff81168114613abf57600080fd5b600060208284031215614f0057600080fd5b61387682614b77565b60008060408385031215614f1c57600080fd5b614f2583614b77565b9150614f3360208401614b77565b90509250929050565b60008060408385031215614f4f57600080fd5b614f5883614b77565b9150602083013560048110614f6c57600080fd5b809150509250929050565b600080600080600080600060a0888a031215614f9257600080fd5b614f9b88614b77565b9650614fa960208901614ea4565b9550614fb760408901614b77565b9450606088013567ffffffffffffffff80821115614fd457600080fd5b614fe08b838c01614e15565b909650945060808a0135915080821115614ff957600080fd5b506150068a828b01614e15565b989b979a50959850939692959293505050565b6000806020838503121561502c57600080fd5b823567ffffffffffffffff81111561504357600080fd5b61504f85828601614b9b565b90969095509350505050565b60008060006040848603121561507057600080fd5b833567ffffffffffffffff81111561508757600080fd5b61509386828701614b9b565b90945092506150a6905060208501614b77565b90509250925092565b600080600080608085870312156150c557600080fd5b843567ffffffffffffffff808211156150dd57600080fd5b818701915087601f8301126150f157600080fd5b81356020615101614c088361598c565b8083825282820191508286018c848660051b890101111561512157600080fd5b600096505b84871015615144578035835260019690960195918301918301615126565b509850508801359250508082111561515b57600080fd5b61516788838901614d21565b9450604087013591508082111561517d57600080fd5b61518988838901614c60565b9350606087013591508082111561519f57600080fd5b506151ac87828801614be7565b91505092959194509250565b6000602082840312156151ca57600080fd5b8151613ab381615d75565b600080604083850312156151e857600080fd5b82516151f381615d75565b602084015190925067ffffffffffffffff81111561521057600080fd5b61521c85828601614e57565b9150509250929050565b60006020828403121561523857600080fd5b5051919050565b6000806020838503121561525257600080fd5b823567ffffffffffffffff81111561526957600080fd5b61504f85828601614e15565b60006020828403121561528757600080fd5b815167ffffffffffffffff81111561529e57600080fd5b614e9c84828501614e57565b6000602082840312156152bc57600080fd5b5035919050565b600080604083850312156152d657600080fd5b82359150614f3360208401614b77565b6000806000604084860312156152fb57600080fd5b83359250602084013567ffffffffffffffff81111561531957600080fd5b61532586828701614e15565b9497909650939450505050565b6000806040838503121561534557600080fd5b82359150614f3360208401614ea4565b6000806040838503121561536857600080fd5b82359150614f3360208401614ed2565b600080600080600060a0868803121561539057600080fd5b61539986614eb8565b94506020860151935060408601519250606086015191506153bc60808701614eb8565b90509295509295909350565b6000602082840312156153da57600080fd5b815160ff81168114613ab357600080fd5b600081518084526020808501945080840160005b8381101561543157815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016153ff565b509495945050505050565b6000815180845260208085019450848260051b860182860160005b85811015614c5357838303895261546f838351615481565b98850198925090840190600101615457565b60008151808452615499816020860160208601615b7a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8284823760008382016000815283516154e8818360208801615b7a565b0195945050505050565b60008251615504818460208701615b7a565b9190910192915050565b6000604082016040835280865480835260608501915087600052602092508260002060005b8281101561556557815473ffffffffffffffffffffffffffffffffffffffff1684529284019260019182019101615533565b505050838103828501528481528590820160005b868110156155b25773ffffffffffffffffffffffffffffffffffffffff61559f84614b77565b1682529183019190830190600101615579565b50979650505050505050565b60006080808352868184015260a07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8811156155f957600080fd5b8760051b808a838701378085019050818101600081526020838784030181880152818a5180845260c093508385019150828c01945060005b818110156156cf578551805163ffffffff908116855285820151168585015260408082015115159085015260608082015173ffffffffffffffffffffffffffffffffffffffff1690850152888101516bffffffffffffffffffffffff1689850152878101516156af898601826bffffffffffffffffffffffff169052565b5085015163ffffffff16838601529483019460e090920191600101615631565b505087810360408901526156e3818b61543c565b9550505050505082810360608401526156fc81856153eb565b98975050505050505050565b861515815260c06020820152600061572360c0830188615481565b90506007861061573557615735615cb9565b8560408301528460608301528360808301528260a0830152979650505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b6020815260006138766020830184615481565b600060208083526000845481600182811c9150808316806157da57607f831692505b858310811415615811577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b87860183815260200181801561582e576001811461585d57615888565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528782019650615888565b60008b81526020902060005b8681101561588257815484820152908501908901615869565b83019750505b50949998505050505050505050565b60208101600383106158ab576158ab615cb9565b91905290565b6020815263ffffffff82511660208201526020820151604082015260006040830151606080840152614e9c6080840182615481565b60ff8416815260ff8316602082015260606040820152600061590b6060830184615481565b95945050505050565b60405160e0810167ffffffffffffffff8111828210171561593757615937615d46565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561598457615984615d46565b604052919050565b600067ffffffffffffffff8211156159a6576159a6615d46565b5060051b60200190565b600067ffffffffffffffff8211156159ca576159ca615d46565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60008219821115615a0957615a09615c5b565b500190565b600060ff821660ff84168060ff03821115615a2b57615a2b615c5b565b019392505050565b60006bffffffffffffffffffffffff808316818516808303821115615a5a57615a5a615c5b565b01949350505050565b600082615a7257615a72615c8a565b500490565b60006bffffffffffffffffffffffff80841680615a9657615a96615c8a565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615ada57615ada615c5b565b500290565b600063ffffffff80831681851681830481118215151615615b0257615b02615c5b565b02949350505050565b60006bffffffffffffffffffffffff80831681851681830481118215151615615b0257615b02615c5b565b600082821015615b4857615b48615c5b565b500390565b60006bffffffffffffffffffffffff83811690831681811015615b7257615b72615c5b565b039392505050565b60005b83811015615b95578181015183820152602001615b7d565b83811115615ba4576000848401525b50505050565b600181811c90821680615bbe57607f821691505b60208210811415615bf8577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615c3057615c30615c5b565b5060010190565b600063ffffffff80831681811415615c5157615c51615c5b565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b80151581146130da57600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
}

var KeeperRegistryLogic20ABI = KeeperRegistryLogic20MetaData.ABI

var KeeperRegistryLogic20Bin = KeeperRegistryLogic20MetaData.Bin

func DeployKeeperRegistryLogic20(auth *bind.TransactOpts, backend bind.ContractBackend, mode uint8, link common.Address, linkNativeFeed common.Address, fastGasFeed common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogic20, error) {
	parsed, err := KeeperRegistryLogic20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogic20Bin), backend, mode, link, linkNativeFeed, fastGasFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeeperRegistryLogic20{address: address, abi: *parsed, KeeperRegistryLogic20Caller: KeeperRegistryLogic20Caller{contract: contract}, KeeperRegistryLogic20Transactor: KeeperRegistryLogic20Transactor{contract: contract}, KeeperRegistryLogic20Filterer: KeeperRegistryLogic20Filterer{contract: contract}}, nil
}

type KeeperRegistryLogic20 struct {
	address common.Address
	abi     abi.ABI
	KeeperRegistryLogic20Caller
	KeeperRegistryLogic20Transactor
	KeeperRegistryLogic20Filterer
}

type KeeperRegistryLogic20Caller struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogic20Transactor struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogic20Filterer struct {
	contract *bind.BoundContract
}

type KeeperRegistryLogic20Session struct {
	Contract     *KeeperRegistryLogic20
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogic20CallerSession struct {
	Contract *KeeperRegistryLogic20Caller
	CallOpts bind.CallOpts
}

type KeeperRegistryLogic20TransactorSession struct {
	Contract     *KeeperRegistryLogic20Transactor
	TransactOpts bind.TransactOpts
}

type KeeperRegistryLogic20Raw struct {
	Contract *KeeperRegistryLogic20
}

type KeeperRegistryLogic20CallerRaw struct {
	Contract *KeeperRegistryLogic20Caller
}

type KeeperRegistryLogic20TransactorRaw struct {
	Contract *KeeperRegistryLogic20Transactor
}

func NewKeeperRegistryLogic20(address common.Address, backend bind.ContractBackend) (*KeeperRegistryLogic20, error) {
	abi, err := abi.JSON(strings.NewReader(KeeperRegistryLogic20ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeeperRegistryLogic20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20{address: address, abi: abi, KeeperRegistryLogic20Caller: KeeperRegistryLogic20Caller{contract: contract}, KeeperRegistryLogic20Transactor: KeeperRegistryLogic20Transactor{contract: contract}, KeeperRegistryLogic20Filterer: KeeperRegistryLogic20Filterer{contract: contract}}, nil
}

func NewKeeperRegistryLogic20Caller(address common.Address, caller bind.ContractCaller) (*KeeperRegistryLogic20Caller, error) {
	contract, err := bindKeeperRegistryLogic20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20Caller{contract: contract}, nil
}

func NewKeeperRegistryLogic20Transactor(address common.Address, transactor bind.ContractTransactor) (*KeeperRegistryLogic20Transactor, error) {
	contract, err := bindKeeperRegistryLogic20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20Transactor{contract: contract}, nil
}

func NewKeeperRegistryLogic20Filterer(address common.Address, filterer bind.ContractFilterer) (*KeeperRegistryLogic20Filterer, error) {
	contract, err := bindKeeperRegistryLogic20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20Filterer{contract: contract}, nil
}

func bindKeeperRegistryLogic20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeeperRegistryLogic20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogic20.Contract.KeeperRegistryLogic20Caller.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.KeeperRegistryLogic20Transactor.contract.Transfer(opts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.KeeperRegistryLogic20Transactor.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeeperRegistryLogic20.Contract.contract.Call(opts, result, method, params...)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.contract.Transfer(opts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.contract.Transact(opts, method, params...)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Caller) GetFastGasFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic20.contract.Call(opts, &out, "getFastGasFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20CallerSession) GetFastGasFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.GetFastGasFeedAddress(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Caller) GetLinkAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic20.contract.Call(opts, &out, "getLinkAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.GetLinkAddress(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20CallerSession) GetLinkAddress() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.GetLinkAddress(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Caller) GetLinkNativeFeedAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic20.contract.Call(opts, &out, "getLinkNativeFeedAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20CallerSession) GetLinkNativeFeedAddress() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.GetLinkNativeFeedAddress(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Caller) GetMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogic20.contract.Call(opts, &out, "getMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) GetMode() (uint8, error) {
	return _KeeperRegistryLogic20.Contract.GetMode(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20CallerSession) GetMode() (uint8, error) {
	return _KeeperRegistryLogic20.Contract.GetMode(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeeperRegistryLogic20.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) Owner() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.Owner(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20CallerSession) Owner() (common.Address, error) {
	return _KeeperRegistryLogic20.Contract.Owner(&_KeeperRegistryLogic20.CallOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "acceptOwnership")
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AcceptOwnership(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AcceptOwnership(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AcceptPayeeship(&_KeeperRegistryLogic20.TransactOpts, transmitter)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AcceptPayeeship(&_KeeperRegistryLogic20.TransactOpts, transmitter)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) AcceptUpkeepAdmin(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "acceptUpkeepAdmin", id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) AcceptUpkeepAdmin(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AcceptUpkeepAdmin(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "addFunds", id, amount)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AddFunds(&_KeeperRegistryLogic20.TransactOpts, id, amount)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) AddFunds(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.AddFunds(&_KeeperRegistryLogic20.TransactOpts, id, amount)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "cancelUpkeep", id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.CancelUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) CancelUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.CancelUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) CheckUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "checkUpkeep", id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.CheckUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) CheckUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.CheckUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) MigrateUpkeeps(opts *bind.TransactOpts, ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "migrateUpkeeps", ids, destination)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.MigrateUpkeeps(&_KeeperRegistryLogic20.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) MigrateUpkeeps(ids []*big.Int, destination common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.MigrateUpkeeps(&_KeeperRegistryLogic20.TransactOpts, ids, destination)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "pause")
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.Pause(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) Pause() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.Pause(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) PauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "pauseUpkeep", id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.PauseUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) PauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.PauseUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) ReceiveUpkeeps(opts *bind.TransactOpts, encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "receiveUpkeeps", encodedUpkeeps)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.ReceiveUpkeeps(&_KeeperRegistryLogic20.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) ReceiveUpkeeps(encodedUpkeeps []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.ReceiveUpkeeps(&_KeeperRegistryLogic20.TransactOpts, encodedUpkeeps)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "recoverFunds")
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.RecoverFunds(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) RecoverFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.RecoverFunds(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "registerUpkeep", target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.RegisterUpkeep(&_KeeperRegistryLogic20.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) RegisterUpkeep(target common.Address, gasLimit uint32, admin common.Address, checkData []byte, offchainConfig []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.RegisterUpkeep(&_KeeperRegistryLogic20.TransactOpts, target, gasLimit, admin, checkData, offchainConfig)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) SetPayees(opts *bind.TransactOpts, payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "setPayees", payees)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetPayees(&_KeeperRegistryLogic20.TransactOpts, payees)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) SetPayees(payees []common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetPayees(&_KeeperRegistryLogic20.TransactOpts, payees)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) SetPeerRegistryMigrationPermission(opts *bind.TransactOpts, peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "setPeerRegistryMigrationPermission", peer, permission)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogic20.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) SetPeerRegistryMigrationPermission(peer common.Address, permission uint8) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetPeerRegistryMigrationPermission(&_KeeperRegistryLogic20.TransactOpts, peer, permission)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) SetUpkeepGasLimit(opts *bind.TransactOpts, id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "setUpkeepGasLimit", id, gasLimit)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogic20.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) SetUpkeepGasLimit(id *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetUpkeepGasLimit(&_KeeperRegistryLogic20.TransactOpts, id, gasLimit)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) SetUpkeepOffchainConfig(opts *bind.TransactOpts, id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "setUpkeepOffchainConfig", id, config)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogic20.TransactOpts, id, config)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) SetUpkeepOffchainConfig(id *big.Int, config []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.SetUpkeepOffchainConfig(&_KeeperRegistryLogic20.TransactOpts, id, config)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "transferOwnership", to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.TransferOwnership(&_KeeperRegistryLogic20.TransactOpts, to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.TransferOwnership(&_KeeperRegistryLogic20.TransactOpts, to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.TransferPayeeship(&_KeeperRegistryLogic20.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.TransferPayeeship(&_KeeperRegistryLogic20.TransactOpts, transmitter, proposed)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) TransferUpkeepAdmin(opts *bind.TransactOpts, id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "transferUpkeepAdmin", id, proposed)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogic20.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) TransferUpkeepAdmin(id *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.TransferUpkeepAdmin(&_KeeperRegistryLogic20.TransactOpts, id, proposed)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "unpause")
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.Unpause(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) Unpause() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.Unpause(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) UnpauseUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "unpauseUpkeep", id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.UnpauseUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) UnpauseUpkeep(id *big.Int) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.UnpauseUpkeep(&_KeeperRegistryLogic20.TransactOpts, id)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) UpdateCheckData(opts *bind.TransactOpts, id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "updateCheckData", id, newCheckData)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.UpdateCheckData(&_KeeperRegistryLogic20.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) UpdateCheckData(id *big.Int, newCheckData []byte) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.UpdateCheckData(&_KeeperRegistryLogic20.TransactOpts, id, newCheckData)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "withdrawFunds", id, to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.WithdrawFunds(&_KeeperRegistryLogic20.TransactOpts, id, to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) WithdrawFunds(id *big.Int, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.WithdrawFunds(&_KeeperRegistryLogic20.TransactOpts, id, to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) WithdrawOwnerFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "withdrawOwnerFunds")
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) WithdrawOwnerFunds() (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.WithdrawOwnerFunds(&_KeeperRegistryLogic20.TransactOpts)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Transactor) WithdrawPayment(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.contract.Transact(opts, "withdrawPayment", from, to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Session) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.WithdrawPayment(&_KeeperRegistryLogic20.TransactOpts, from, to)
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20TransactorSession) WithdrawPayment(from common.Address, to common.Address) (*types.Transaction, error) {
	return _KeeperRegistryLogic20.Contract.WithdrawPayment(&_KeeperRegistryLogic20.TransactOpts, from, to)
}

type KeeperRegistryLogic20CancelledUpkeepReportIterator struct {
	Event *KeeperRegistryLogic20CancelledUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20CancelledUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20CancelledUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20CancelledUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20CancelledUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20CancelledUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20CancelledUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20CancelledUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20CancelledUpkeepReportIterator{contract: _KeeperRegistryLogic20.contract, event: "CancelledUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20CancelledUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "CancelledUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20CancelledUpkeepReport)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogic20CancelledUpkeepReport, error) {
	event := new(KeeperRegistryLogic20CancelledUpkeepReport)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "CancelledUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20FundsAddedIterator struct {
	Event *KeeperRegistryLogic20FundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20FundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20FundsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20FundsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20FundsAddedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20FundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20FundsAdded struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogic20FundsAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20FundsAddedIterator{contract: _KeeperRegistryLogic20.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "FundsAdded", idRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20FundsAdded)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "FundsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseFundsAdded(log types.Log) (*KeeperRegistryLogic20FundsAdded, error) {
	event := new(KeeperRegistryLogic20FundsAdded)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20FundsWithdrawnIterator struct {
	Event *KeeperRegistryLogic20FundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20FundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20FundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20FundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20FundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20FundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20FundsWithdrawn struct {
	Id     *big.Int
	Amount *big.Int
	To     common.Address
	Raw    types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20FundsWithdrawnIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20FundsWithdrawnIterator{contract: _KeeperRegistryLogic20.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20FundsWithdrawn, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "FundsWithdrawn", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20FundsWithdrawn)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogic20FundsWithdrawn, error) {
	event := new(KeeperRegistryLogic20FundsWithdrawn)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20InsufficientFundsUpkeepReportIterator struct {
	Event *KeeperRegistryLogic20InsufficientFundsUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20InsufficientFundsUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20InsufficientFundsUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20InsufficientFundsUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20InsufficientFundsUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20InsufficientFundsUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20InsufficientFundsUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20InsufficientFundsUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20InsufficientFundsUpkeepReportIterator{contract: _KeeperRegistryLogic20.contract, event: "InsufficientFundsUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "InsufficientFundsUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20InsufficientFundsUpkeepReport)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogic20InsufficientFundsUpkeepReport, error) {
	event := new(KeeperRegistryLogic20InsufficientFundsUpkeepReport)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "InsufficientFundsUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20OwnerFundsWithdrawnIterator struct {
	Event *KeeperRegistryLogic20OwnerFundsWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20OwnerFundsWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20OwnerFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20OwnerFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20OwnerFundsWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20OwnerFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20OwnerFundsWithdrawn struct {
	Amount *big.Int
	Raw    types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogic20OwnerFundsWithdrawnIterator, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20OwnerFundsWithdrawnIterator{contract: _KeeperRegistryLogic20.contract, event: "OwnerFundsWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20OwnerFundsWithdrawn) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "OwnerFundsWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20OwnerFundsWithdrawn)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogic20OwnerFundsWithdrawn, error) {
	event := new(KeeperRegistryLogic20OwnerFundsWithdrawn)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "OwnerFundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20OwnershipTransferRequestedIterator struct {
	Event *KeeperRegistryLogic20OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20OwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20OwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogic20OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20OwnershipTransferRequestedIterator{contract: _KeeperRegistryLogic20.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20OwnershipTransferRequested)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogic20OwnershipTransferRequested, error) {
	event := new(KeeperRegistryLogic20OwnershipTransferRequested)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20OwnershipTransferredIterator struct {
	Event *KeeperRegistryLogic20OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogic20OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20OwnershipTransferredIterator{contract: _KeeperRegistryLogic20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20OwnershipTransferred)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogic20OwnershipTransferred, error) {
	event := new(KeeperRegistryLogic20OwnershipTransferred)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20PausedIterator struct {
	Event *KeeperRegistryLogic20Paused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20PausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20Paused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20Paused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20PausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20Paused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogic20PausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20PausedIterator{contract: _KeeperRegistryLogic20.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20Paused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20Paused)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParsePaused(log types.Log) (*KeeperRegistryLogic20Paused, error) {
	event := new(KeeperRegistryLogic20Paused)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20PayeesUpdatedIterator struct {
	Event *KeeperRegistryLogic20PayeesUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20PayeesUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20PayeesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20PayeesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20PayeesUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20PayeesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20PayeesUpdated struct {
	Transmitters []common.Address
	Payees       []common.Address
	Raw          types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogic20PayeesUpdatedIterator, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20PayeesUpdatedIterator{contract: _KeeperRegistryLogic20.contract, event: "PayeesUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PayeesUpdated) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "PayeesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20PayeesUpdated)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogic20PayeesUpdated, error) {
	event := new(KeeperRegistryLogic20PayeesUpdated)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PayeesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20PayeeshipTransferRequestedIterator struct {
	Event *KeeperRegistryLogic20PayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20PayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20PayeeshipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20PayeeshipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20PayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20PayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20PayeeshipTransferRequested struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogic20PayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20PayeeshipTransferRequestedIterator{contract: _KeeperRegistryLogic20.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20PayeeshipTransferRequested)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogic20PayeeshipTransferRequested, error) {
	event := new(KeeperRegistryLogic20PayeeshipTransferRequested)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20PayeeshipTransferredIterator struct {
	Event *KeeperRegistryLogic20PayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20PayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20PayeeshipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20PayeeshipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20PayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20PayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20PayeeshipTransferred struct {
	Transmitter common.Address
	From        common.Address
	To          common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogic20PayeeshipTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20PayeeshipTransferredIterator{contract: _KeeperRegistryLogic20.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20PayeeshipTransferred)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogic20PayeeshipTransferred, error) {
	event := new(KeeperRegistryLogic20PayeeshipTransferred)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20PaymentWithdrawnIterator struct {
	Event *KeeperRegistryLogic20PaymentWithdrawn

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20PaymentWithdrawnIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20PaymentWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20PaymentWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20PaymentWithdrawnIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20PaymentWithdrawn struct {
	Transmitter common.Address
	Amount      *big.Int
	To          common.Address
	Payee       common.Address
	Raw         types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogic20PaymentWithdrawnIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20PaymentWithdrawnIterator{contract: _KeeperRegistryLogic20.contract, event: "PaymentWithdrawn", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "PaymentWithdrawn", transmitterRule, amountRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20PaymentWithdrawn)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogic20PaymentWithdrawn, error) {
	event := new(KeeperRegistryLogic20PaymentWithdrawn)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "PaymentWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20ReorgedUpkeepReportIterator struct {
	Event *KeeperRegistryLogic20ReorgedUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20ReorgedUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20ReorgedUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20ReorgedUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20ReorgedUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20ReorgedUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20ReorgedUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20ReorgedUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20ReorgedUpkeepReportIterator{contract: _KeeperRegistryLogic20.contract, event: "ReorgedUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "ReorgedUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20ReorgedUpkeepReport)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogic20ReorgedUpkeepReport, error) {
	event := new(KeeperRegistryLogic20ReorgedUpkeepReport)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "ReorgedUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20StaleUpkeepReportIterator struct {
	Event *KeeperRegistryLogic20StaleUpkeepReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20StaleUpkeepReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20StaleUpkeepReport)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20StaleUpkeepReport)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20StaleUpkeepReportIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20StaleUpkeepReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20StaleUpkeepReport struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20StaleUpkeepReportIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20StaleUpkeepReportIterator{contract: _KeeperRegistryLogic20.contract, event: "StaleUpkeepReport", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20StaleUpkeepReport, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "StaleUpkeepReport", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20StaleUpkeepReport)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogic20StaleUpkeepReport, error) {
	event := new(KeeperRegistryLogic20StaleUpkeepReport)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "StaleUpkeepReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UnpausedIterator struct {
	Event *KeeperRegistryLogic20Unpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20Unpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20Unpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20Unpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogic20UnpausedIterator, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UnpausedIterator{contract: _KeeperRegistryLogic20.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20Unpaused) (event.Subscription, error) {

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20Unpaused)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUnpaused(log types.Log) (*KeeperRegistryLogic20Unpaused, error) {
	event := new(KeeperRegistryLogic20Unpaused)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepAdminTransferRequestedIterator struct {
	Event *KeeperRegistryLogic20UpkeepAdminTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepAdminTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepAdminTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepAdminTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepAdminTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepAdminTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepAdminTransferRequested struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogic20UpkeepAdminTransferRequestedIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepAdminTransferRequestedIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepAdminTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepAdminTransferRequested", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepAdminTransferRequested)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogic20UpkeepAdminTransferRequested, error) {
	event := new(KeeperRegistryLogic20UpkeepAdminTransferRequested)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepAdminTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepAdminTransferredIterator struct {
	Event *KeeperRegistryLogic20UpkeepAdminTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepAdminTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepAdminTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepAdminTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepAdminTransferredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepAdminTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepAdminTransferred struct {
	Id   *big.Int
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogic20UpkeepAdminTransferredIterator, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepAdminTransferredIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepAdminTransferred", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepAdminTransferred", idRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepAdminTransferred)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogic20UpkeepAdminTransferred, error) {
	event := new(KeeperRegistryLogic20UpkeepAdminTransferred)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepAdminTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepCanceledIterator struct {
	Event *KeeperRegistryLogic20UpkeepCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepCanceledIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepCanceled struct {
	Id            *big.Int
	AtBlockHeight uint64
	Raw           types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogic20UpkeepCanceledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepCanceledIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepCanceled", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var atBlockHeightRule []interface{}
	for _, atBlockHeightItem := range atBlockHeight {
		atBlockHeightRule = append(atBlockHeightRule, atBlockHeightItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepCanceled", idRule, atBlockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepCanceled)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogic20UpkeepCanceled, error) {
	event := new(KeeperRegistryLogic20UpkeepCanceled)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepCheckDataUpdatedIterator struct {
	Event *KeeperRegistryLogic20UpkeepCheckDataUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepCheckDataUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepCheckDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepCheckDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepCheckDataUpdatedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepCheckDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepCheckDataUpdated struct {
	Id           *big.Int
	NewCheckData []byte
	Raw          types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepCheckDataUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepCheckDataUpdatedIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepCheckDataUpdated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepCheckDataUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepCheckDataUpdated)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryLogic20UpkeepCheckDataUpdated, error) {
	event := new(KeeperRegistryLogic20UpkeepCheckDataUpdated)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepCheckDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepGasLimitSetIterator struct {
	Event *KeeperRegistryLogic20UpkeepGasLimitSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepGasLimitSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepGasLimitSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepGasLimitSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepGasLimitSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepGasLimitSet struct {
	Id       *big.Int
	GasLimit *big.Int
	Raw      types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepGasLimitSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepGasLimitSetIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepGasLimitSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepGasLimitSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepGasLimitSet)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogic20UpkeepGasLimitSet, error) {
	event := new(KeeperRegistryLogic20UpkeepGasLimitSet)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepMigratedIterator struct {
	Event *KeeperRegistryLogic20UpkeepMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepMigratedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepMigrated struct {
	Id               *big.Int
	RemainingBalance *big.Int
	Destination      common.Address
	Raw              types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepMigratedIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepMigrated", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepMigrated, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepMigrated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepMigrated)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogic20UpkeepMigrated, error) {
	event := new(KeeperRegistryLogic20UpkeepMigrated)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepOffchainConfigSetIterator struct {
	Event *KeeperRegistryLogic20UpkeepOffchainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepOffchainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepOffchainConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepOffchainConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepOffchainConfigSetIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepOffchainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepOffchainConfigSet struct {
	Id             *big.Int
	OffchainConfig []byte
	Raw            types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepOffchainConfigSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepOffchainConfigSetIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepOffchainConfigSet", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepOffchainConfigSet", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepOffchainConfigSet)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogic20UpkeepOffchainConfigSet, error) {
	event := new(KeeperRegistryLogic20UpkeepOffchainConfigSet)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepOffchainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepPausedIterator struct {
	Event *KeeperRegistryLogic20UpkeepPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepPausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepPaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepPausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepPausedIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepPaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepPaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepPaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepPaused)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogic20UpkeepPaused, error) {
	event := new(KeeperRegistryLogic20UpkeepPaused)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepPerformedIterator struct {
	Event *KeeperRegistryLogic20UpkeepPerformed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepPerformedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepPerformed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepPerformed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepPerformedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepPerformedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepPerformed struct {
	Id               *big.Int
	Success          bool
	CheckBlockNumber uint32
	GasUsed          *big.Int
	GasOverhead      *big.Int
	TotalPayment     *big.Int
	Raw              types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogic20UpkeepPerformedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepPerformedIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepPerformed", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var successRule []interface{}
	for _, successItem := range success {
		successRule = append(successRule, successItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepPerformed", idRule, successRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepPerformed)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogic20UpkeepPerformed, error) {
	event := new(KeeperRegistryLogic20UpkeepPerformed)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepPerformed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepReceivedIterator struct {
	Event *KeeperRegistryLogic20UpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepReceived struct {
	Id              *big.Int
	StartingBalance *big.Int
	ImportedFrom    common.Address
	Raw             types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepReceivedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepReceivedIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepReceived", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepReceived, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepReceived", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepReceived)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogic20UpkeepReceived, error) {
	event := new(KeeperRegistryLogic20UpkeepReceived)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepRegisteredIterator struct {
	Event *KeeperRegistryLogic20UpkeepRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepRegisteredIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepRegistered struct {
	Id         *big.Int
	ExecuteGas uint32
	Admin      common.Address
	Raw        types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepRegisteredIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepRegistered", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepRegistered)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogic20UpkeepRegistered, error) {
	event := new(KeeperRegistryLogic20UpkeepRegistered)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeeperRegistryLogic20UpkeepUnpausedIterator struct {
	Event *KeeperRegistryLogic20UpkeepUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeeperRegistryLogic20UpkeepUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeeperRegistryLogic20UpkeepUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(KeeperRegistryLogic20UpkeepUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeeperRegistryLogic20UpkeepUnpausedIterator) Error() error {
	return it.fail
}

func (it *KeeperRegistryLogic20UpkeepUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeeperRegistryLogic20UpkeepUnpaused struct {
	Id  *big.Int
	Raw types.Log
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepUnpausedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.FilterLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return &KeeperRegistryLogic20UpkeepUnpausedIterator{contract: _KeeperRegistryLogic20.contract, event: "UpkeepUnpaused", logs: logs, sub: sub}, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepUnpaused, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _KeeperRegistryLogic20.contract.WatchLogs(opts, "UpkeepUnpaused", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeeperRegistryLogic20UpkeepUnpaused)
				if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20Filterer) ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogic20UpkeepUnpaused, error) {
	event := new(KeeperRegistryLogic20UpkeepUnpaused)
	if err := _KeeperRegistryLogic20.contract.UnpackLog(event, "UpkeepUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeeperRegistryLogic20.abi.Events["CancelledUpkeepReport"].ID:
		return _KeeperRegistryLogic20.ParseCancelledUpkeepReport(log)
	case _KeeperRegistryLogic20.abi.Events["FundsAdded"].ID:
		return _KeeperRegistryLogic20.ParseFundsAdded(log)
	case _KeeperRegistryLogic20.abi.Events["FundsWithdrawn"].ID:
		return _KeeperRegistryLogic20.ParseFundsWithdrawn(log)
	case _KeeperRegistryLogic20.abi.Events["InsufficientFundsUpkeepReport"].ID:
		return _KeeperRegistryLogic20.ParseInsufficientFundsUpkeepReport(log)
	case _KeeperRegistryLogic20.abi.Events["OwnerFundsWithdrawn"].ID:
		return _KeeperRegistryLogic20.ParseOwnerFundsWithdrawn(log)
	case _KeeperRegistryLogic20.abi.Events["OwnershipTransferRequested"].ID:
		return _KeeperRegistryLogic20.ParseOwnershipTransferRequested(log)
	case _KeeperRegistryLogic20.abi.Events["OwnershipTransferred"].ID:
		return _KeeperRegistryLogic20.ParseOwnershipTransferred(log)
	case _KeeperRegistryLogic20.abi.Events["Paused"].ID:
		return _KeeperRegistryLogic20.ParsePaused(log)
	case _KeeperRegistryLogic20.abi.Events["PayeesUpdated"].ID:
		return _KeeperRegistryLogic20.ParsePayeesUpdated(log)
	case _KeeperRegistryLogic20.abi.Events["PayeeshipTransferRequested"].ID:
		return _KeeperRegistryLogic20.ParsePayeeshipTransferRequested(log)
	case _KeeperRegistryLogic20.abi.Events["PayeeshipTransferred"].ID:
		return _KeeperRegistryLogic20.ParsePayeeshipTransferred(log)
	case _KeeperRegistryLogic20.abi.Events["PaymentWithdrawn"].ID:
		return _KeeperRegistryLogic20.ParsePaymentWithdrawn(log)
	case _KeeperRegistryLogic20.abi.Events["ReorgedUpkeepReport"].ID:
		return _KeeperRegistryLogic20.ParseReorgedUpkeepReport(log)
	case _KeeperRegistryLogic20.abi.Events["StaleUpkeepReport"].ID:
		return _KeeperRegistryLogic20.ParseStaleUpkeepReport(log)
	case _KeeperRegistryLogic20.abi.Events["Unpaused"].ID:
		return _KeeperRegistryLogic20.ParseUnpaused(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepAdminTransferRequested"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepAdminTransferRequested(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepAdminTransferred"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepAdminTransferred(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepCanceled"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepCanceled(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepCheckDataUpdated"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepCheckDataUpdated(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepGasLimitSet"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepGasLimitSet(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepMigrated"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepMigrated(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepOffchainConfigSet"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepOffchainConfigSet(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepPaused"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepPaused(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepPerformed"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepPerformed(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepReceived"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepReceived(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepRegistered"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepRegistered(log)
	case _KeeperRegistryLogic20.abi.Events["UpkeepUnpaused"].ID:
		return _KeeperRegistryLogic20.ParseUpkeepUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeeperRegistryLogic20CancelledUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0xd84831b6a3a7fbd333f42fe7f9104a139da6cca4cc1507aef4ddad79b31d017f")
}

func (KeeperRegistryLogic20FundsAdded) Topic() common.Hash {
	return common.HexToHash("0xafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa734891506203")
}

func (KeeperRegistryLogic20FundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0xf3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318")
}

func (KeeperRegistryLogic20InsufficientFundsUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x7895fdfe292beab0842d5beccd078e85296b9e17a30eaee4c261a2696b84eb96")
}

func (KeeperRegistryLogic20OwnerFundsWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f1")
}

func (KeeperRegistryLogic20OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeeperRegistryLogic20OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (KeeperRegistryLogic20Paused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (KeeperRegistryLogic20PayeesUpdated) Topic() common.Hash {
	return common.HexToHash("0xa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725")
}

func (KeeperRegistryLogic20PayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (KeeperRegistryLogic20PayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (KeeperRegistryLogic20PaymentWithdrawn) Topic() common.Hash {
	return common.HexToHash("0x9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f40698")
}

func (KeeperRegistryLogic20ReorgedUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x561ff77e59394941a01a456497a9418dea82e2a39abb3ecebfb1cef7e0bfdc13")
}

func (KeeperRegistryLogic20StaleUpkeepReport) Topic() common.Hash {
	return common.HexToHash("0x5aa44821f7938098502bff537fbbdc9aaaa2fa655c10740646fce27e54987a89")
}

func (KeeperRegistryLogic20Unpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (KeeperRegistryLogic20UpkeepAdminTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b35")
}

func (KeeperRegistryLogic20UpkeepAdminTransferred) Topic() common.Hash {
	return common.HexToHash("0x5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c")
}

func (KeeperRegistryLogic20UpkeepCanceled) Topic() common.Hash {
	return common.HexToHash("0x91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f791181")
}

func (KeeperRegistryLogic20UpkeepCheckDataUpdated) Topic() common.Hash {
	return common.HexToHash("0x7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf")
}

func (KeeperRegistryLogic20UpkeepGasLimitSet) Topic() common.Hash {
	return common.HexToHash("0xc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c")
}

func (KeeperRegistryLogic20UpkeepMigrated) Topic() common.Hash {
	return common.HexToHash("0xb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff")
}

func (KeeperRegistryLogic20UpkeepOffchainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf4850")
}

func (KeeperRegistryLogic20UpkeepPaused) Topic() common.Hash {
	return common.HexToHash("0x8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f")
}

func (KeeperRegistryLogic20UpkeepPerformed) Topic() common.Hash {
	return common.HexToHash("0x29233ba1d7b302b8fe230ad0b81423aba5371b2a6f6b821228212385ee6a4420")
}

func (KeeperRegistryLogic20UpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a71")
}

func (KeeperRegistryLogic20UpkeepRegistered) Topic() common.Hash {
	return common.HexToHash("0xbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012")
}

func (KeeperRegistryLogic20UpkeepUnpaused) Topic() common.Hash {
	return common.HexToHash("0x7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a47456")
}

func (_KeeperRegistryLogic20 *KeeperRegistryLogic20) Address() common.Address {
	return _KeeperRegistryLogic20.address
}

type KeeperRegistryLogic20Interface interface {
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

	FilterCancelledUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20CancelledUpkeepReportIterator, error)

	WatchCancelledUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20CancelledUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseCancelledUpkeepReport(log types.Log) (*KeeperRegistryLogic20CancelledUpkeepReport, error)

	FilterFundsAdded(opts *bind.FilterOpts, id []*big.Int, from []common.Address) (*KeeperRegistryLogic20FundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20FundsAdded, id []*big.Int, from []common.Address) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*KeeperRegistryLogic20FundsAdded, error)

	FilterFundsWithdrawn(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20FundsWithdrawnIterator, error)

	WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20FundsWithdrawn, id []*big.Int) (event.Subscription, error)

	ParseFundsWithdrawn(log types.Log) (*KeeperRegistryLogic20FundsWithdrawn, error)

	FilterInsufficientFundsUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20InsufficientFundsUpkeepReportIterator, error)

	WatchInsufficientFundsUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20InsufficientFundsUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseInsufficientFundsUpkeepReport(log types.Log) (*KeeperRegistryLogic20InsufficientFundsUpkeepReport, error)

	FilterOwnerFundsWithdrawn(opts *bind.FilterOpts) (*KeeperRegistryLogic20OwnerFundsWithdrawnIterator, error)

	WatchOwnerFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20OwnerFundsWithdrawn) (event.Subscription, error)

	ParseOwnerFundsWithdrawn(log types.Log) (*KeeperRegistryLogic20OwnerFundsWithdrawn, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogic20OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeeperRegistryLogic20OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeeperRegistryLogic20OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeeperRegistryLogic20OwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*KeeperRegistryLogic20PausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20Paused) (event.Subscription, error)

	ParsePaused(log types.Log) (*KeeperRegistryLogic20Paused, error)

	FilterPayeesUpdated(opts *bind.FilterOpts) (*KeeperRegistryLogic20PayeesUpdatedIterator, error)

	WatchPayeesUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PayeesUpdated) (event.Subscription, error)

	ParsePayeesUpdated(log types.Log) (*KeeperRegistryLogic20PayeesUpdated, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogic20PayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PayeeshipTransferRequested, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*KeeperRegistryLogic20PayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, from []common.Address, to []common.Address) (*KeeperRegistryLogic20PayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PayeeshipTransferred, transmitter []common.Address, from []common.Address, to []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*KeeperRegistryLogic20PayeeshipTransferred, error)

	FilterPaymentWithdrawn(opts *bind.FilterOpts, transmitter []common.Address, amount []*big.Int, to []common.Address) (*KeeperRegistryLogic20PaymentWithdrawnIterator, error)

	WatchPaymentWithdrawn(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20PaymentWithdrawn, transmitter []common.Address, amount []*big.Int, to []common.Address) (event.Subscription, error)

	ParsePaymentWithdrawn(log types.Log) (*KeeperRegistryLogic20PaymentWithdrawn, error)

	FilterReorgedUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20ReorgedUpkeepReportIterator, error)

	WatchReorgedUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20ReorgedUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseReorgedUpkeepReport(log types.Log) (*KeeperRegistryLogic20ReorgedUpkeepReport, error)

	FilterStaleUpkeepReport(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20StaleUpkeepReportIterator, error)

	WatchStaleUpkeepReport(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20StaleUpkeepReport, id []*big.Int) (event.Subscription, error)

	ParseStaleUpkeepReport(log types.Log) (*KeeperRegistryLogic20StaleUpkeepReport, error)

	FilterUnpaused(opts *bind.FilterOpts) (*KeeperRegistryLogic20UnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20Unpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*KeeperRegistryLogic20Unpaused, error)

	FilterUpkeepAdminTransferRequested(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogic20UpkeepAdminTransferRequestedIterator, error)

	WatchUpkeepAdminTransferRequested(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepAdminTransferRequested, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferRequested(log types.Log) (*KeeperRegistryLogic20UpkeepAdminTransferRequested, error)

	FilterUpkeepAdminTransferred(opts *bind.FilterOpts, id []*big.Int, from []common.Address, to []common.Address) (*KeeperRegistryLogic20UpkeepAdminTransferredIterator, error)

	WatchUpkeepAdminTransferred(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepAdminTransferred, id []*big.Int, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseUpkeepAdminTransferred(log types.Log) (*KeeperRegistryLogic20UpkeepAdminTransferred, error)

	FilterUpkeepCanceled(opts *bind.FilterOpts, id []*big.Int, atBlockHeight []uint64) (*KeeperRegistryLogic20UpkeepCanceledIterator, error)

	WatchUpkeepCanceled(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepCanceled, id []*big.Int, atBlockHeight []uint64) (event.Subscription, error)

	ParseUpkeepCanceled(log types.Log) (*KeeperRegistryLogic20UpkeepCanceled, error)

	FilterUpkeepCheckDataUpdated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepCheckDataUpdatedIterator, error)

	WatchUpkeepCheckDataUpdated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepCheckDataUpdated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepCheckDataUpdated(log types.Log) (*KeeperRegistryLogic20UpkeepCheckDataUpdated, error)

	FilterUpkeepGasLimitSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepGasLimitSetIterator, error)

	WatchUpkeepGasLimitSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepGasLimitSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepGasLimitSet(log types.Log) (*KeeperRegistryLogic20UpkeepGasLimitSet, error)

	FilterUpkeepMigrated(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepMigratedIterator, error)

	WatchUpkeepMigrated(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepMigrated, id []*big.Int) (event.Subscription, error)

	ParseUpkeepMigrated(log types.Log) (*KeeperRegistryLogic20UpkeepMigrated, error)

	FilterUpkeepOffchainConfigSet(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepOffchainConfigSetIterator, error)

	WatchUpkeepOffchainConfigSet(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepOffchainConfigSet, id []*big.Int) (event.Subscription, error)

	ParseUpkeepOffchainConfigSet(log types.Log) (*KeeperRegistryLogic20UpkeepOffchainConfigSet, error)

	FilterUpkeepPaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepPausedIterator, error)

	WatchUpkeepPaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepPaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepPaused(log types.Log) (*KeeperRegistryLogic20UpkeepPaused, error)

	FilterUpkeepPerformed(opts *bind.FilterOpts, id []*big.Int, success []bool) (*KeeperRegistryLogic20UpkeepPerformedIterator, error)

	WatchUpkeepPerformed(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepPerformed, id []*big.Int, success []bool) (event.Subscription, error)

	ParseUpkeepPerformed(log types.Log) (*KeeperRegistryLogic20UpkeepPerformed, error)

	FilterUpkeepReceived(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepReceivedIterator, error)

	WatchUpkeepReceived(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepReceived, id []*big.Int) (event.Subscription, error)

	ParseUpkeepReceived(log types.Log) (*KeeperRegistryLogic20UpkeepReceived, error)

	FilterUpkeepRegistered(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepRegisteredIterator, error)

	WatchUpkeepRegistered(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepRegistered, id []*big.Int) (event.Subscription, error)

	ParseUpkeepRegistered(log types.Log) (*KeeperRegistryLogic20UpkeepRegistered, error)

	FilterUpkeepUnpaused(opts *bind.FilterOpts, id []*big.Int) (*KeeperRegistryLogic20UpkeepUnpausedIterator, error)

	WatchUpkeepUnpaused(opts *bind.WatchOpts, sink chan<- *KeeperRegistryLogic20UpkeepUnpaused, id []*big.Int) (event.Subscription, error)

	ParseUpkeepUnpaused(log types.Log) (*KeeperRegistryLogic20UpkeepUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
