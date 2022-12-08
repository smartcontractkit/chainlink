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
	ABI: "[{\"inputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.PaymentModel\",\"name\":\"paymentModel\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"fastGasFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ArrayHasNoEntries\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancel\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckDataExceedsLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateEntry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasLimitOutsideRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReport\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCheckDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxPerformDataSizeCanOnlyIncrease\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MigrationNotPermitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnchainConfigNonEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByLINKToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwnerOrRegistrar\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByProposedPayee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlySimulatedBackend\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpausedUpkeep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ParameterLengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentGreaterThanAllLINK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RegistryPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"TargetCheckReverted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TranscoderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotCanceled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpkeepNotNeeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueNotChanged\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"CancelledUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"FundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"InsufficientFundsUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"OwnerFundsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"PayeesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"PaymentWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ReorgedUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"StaleUpkeepReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"UpkeepAdminTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"atBlockHeight\",\"type\":\"uint64\"}],\"name\":\"UpkeepCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"UpkeepCheckDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"gasLimit\",\"type\":\"uint96\"}],\"name\":\"UpkeepGasLimitSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"UpkeepMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"UpkeepOffchainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"checkBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasOverhead\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"totalPayment\",\"type\":\"uint96\"}],\"name\":\"UpkeepPerformed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"importedFrom\",\"type\":\"address\"}],\"name\":\"UpkeepReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"executeGas\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"UpkeepRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"UpkeepUnpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"acceptUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"},{\"internalType\":\"enumUpkeepFailureReason\",\"name\":\"upkeepFailureReason\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fastGasWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"linkNative\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFastGasFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkNativeFeedAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPaymentModel\",\"outputs\":[{\"internalType\":\"enumKeeperRegistryBase2_0.PaymentModel\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"}],\"name\":\"migrateUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"pauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedUpkeeps\",\"type\":\"bytes\"}],\"name\":\"receiveUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"registerUpkeep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"peer\",\"type\":\"address\"},{\"internalType\":\"enumKeeperRegistryBase2_0.MigrationPermission\",\"name\":\"permission\",\"type\":\"uint8\"}],\"name\":\"setPeerRegistryMigrationPermission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"setUpkeepOffchainConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferUpkeepAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"unpauseUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"newCheckData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawOwnerFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162005ebf38038062005ebf8339810160408190526200003591620001ef565b838383833380600081620000905760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c357620000c38162000126565b505050836002811115620000db57620000db62000251565b60e0816002811115620000f257620000f262000251565b60f81b9052506001600160601b0319606093841b811660805291831b821660a05290911b1660c05250620002679350505050565b6001600160a01b038116331415620001815760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000087565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ea57600080fd5b919050565b600080600080608085870312156200020657600080fd5b8451600381106200021657600080fd5b93506200022660208601620001d2565b92506200023660408601620001d2565b91506200024660608601620001d2565b905092959194509250565b634e487b7160e01b600052602160045260246000fd5b60805160601c60a05160601c60c05160601c60e05160f81c615bc5620002fa600039600081816104250152818161443601526145ed015260008181610237015261405a01526000818161038501526141430152600081816103ec01528181610f610152818161124601528181611bdd015281816122500152818161258501528181612a6c0152612aff0152615bc56000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80638e86139b11610104578063b148ab6b116100a2578063eb5dcd6c11610071578063eb5dcd6c14610410578063f157014114610423578063f2fde38b14610451578063f7d334ba1461046457600080fd5b8063b148ab6b146103bc578063b79550be146103cf578063c8048022146103d7578063ca30e603146103ea57600080fd5b8063a710b221116100de578063a710b2211461035d578063a72aa27e14610370578063b10b673c14610383578063b121e147146103a957600080fd5b80638e86139b14610324578063948108f7146103375780639fab43861461034a57600080fd5b8063744bfe611161017c57806385c1b0ba1161014b57806385c1b0ba146102cd5780638765ecbe146102e05780638da5cb5b146102f35780638dcf0fe71461031157600080fd5b8063744bfe61146102a257806379ba5097146102b55780637d9b97e0146102bd5780638456cb59146102c557600080fd5b80633f4ba83a116101b85780633f4ba83a1461021a5780635165f2f5146102225780636709d0e5146102355780636ded9eae1461028157600080fd5b8063187256e8146101df5780631a2af011146101f45780633b9cce5914610207575b600080fd5b6101f26101ed366004614d0f565b610489565b005b6101f26102023660046150be565b6104fa565b6101f2610215366004614dec565b61064e565b6101f26108a4565b6101f261023036600461508c565b61090a565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61029461028f366004614d4a565b610a84565b604051908152602001610278565b6101f26102b03660046150be565b610c61565b6101f261106c565b6101f261116e565b6101f26112d8565b6101f26102db366004614e2e565b611349565b6101f26102ee36600461508c565b611c68565b60005473ffffffffffffffffffffffffffffffffffffffff16610257565b6101f261031f3660046150e1565b611def565b6101f2610332366004615000565b611e51565b6101f2610345366004615150565b61208d565b6101f26103583660046150e1565b61232c565b6101f261036b366004614cdc565b6123db565b6101f261037e36600461512d565b612666565b7f0000000000000000000000000000000000000000000000000000000000000000610257565b6101f26103b7366004614cc1565b612748565b6101f26103ca36600461508c565b612840565b6101f2612a33565b6101f26103e536600461508c565b612b9e565b7f0000000000000000000000000000000000000000000000000000000000000000610257565b6101f261041e366004614cdc565b612f41565b7f0000000000000000000000000000000000000000000000000000000000000000604051610278919061566f565b6101f261045f366004614cc1565b6130a0565b61047761047236600461508c565b6130b4565b604051610278969594939291906154e0565b610491613748565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260166020526040902080548291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660018360038111156104f1576104f1615a96565b02179055505050565b610503826137cb565b73ffffffffffffffffffffffffffffffffffffffff8116331415610553576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff81166105a0576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff82811691161461064a5760008281526006602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff851690811790915590519091339185917fb1cbb2c4b8480034c27e06da5f096b8233a8fd4497028593a41ff6df79726b3591a45b5050565b610656613748565b600b548114610691576040517fcf54c06a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b600b54811015610863576000600b82815481106106b3576106b3615af4565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff908116808452600c909252604083205491935016908585858181106106fd576106fd615af4565b90506020020160208101906107129190614cc1565b905073ffffffffffffffffffffffffffffffffffffffff811615806107a5575073ffffffffffffffffffffffffffffffffffffffff82161580159061078357508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156107a5575073ffffffffffffffffffffffffffffffffffffffff81811614155b156107dc576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181161461084d5773ffffffffffffffffffffffffffffffffffffffff8381166000908152600c6020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169183169190911790555b505050808061085b906159db565b915050610694565b507fa46de38886467c59be07a0675f14781206a5477d871628af46c2443822fcb725600b8383604051610898939291906152e6565b60405180910390a15050565b6108ac613748565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff1690556040513381527f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa906020015b60405180910390a1565b610913816137cb565b600081815260046020908152604091829020825160e081018452815463ffffffff8082168352640100000000820481169483019490945268010000000000000000810460ff1615159482018590526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c082015290610a15576040517f1b88a78400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff169055610a5460028361387e565b5060405182907f7bada562044eb163f6b4003c4553e4e62825344c0418eea087bed5ee05a4745690600090a25050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314801590610ad557506011546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff163314155b15610b0c576040517fd48b678b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b17600143615913565b6012546040805192406020840152309083015268010000000000000000900463ffffffff1660608201526080016040516020818303038152906040528051906020012060001c9050610ba48189898960008a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509250613893915050565b6012805468010000000000000000900463ffffffff16906008610bc683615a14565b825463ffffffff9182166101009390930a9283029190920219909116179055506000818152601760205260409020610bff9084846147e0565b506040805163ffffffff8916815273ffffffffffffffffffffffffffffffffffffffff8816602082015282917fbae366358c023f887e791d7a62f2e4316f1026bd77f6fb49501a917b3bc5d012910160405180910390a2979650505050505050565b600f546f01000000000000000000000000000000900460ff1615610cb1576040517f37ed32e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff166f0100000000000000000000000000000017905573ffffffffffffffffffffffffffffffffffffffff8116610d38576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000828152600460209081526040808320815160e081018352815463ffffffff8082168352640100000000820481168387015260ff6801000000000000000083041615158386015273ffffffffffffffffffffffffffffffffffffffff6901000000000000000000909204821660608401526001909301546bffffffffffffffffffffffff80821660808501526c0100000000000000000000000082041660a08401527801000000000000000000000000000000000000000000000000900490921660c082015286855260059093529220549091163314610e45576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b43816020015163ffffffff161115610e89576040517fff84e5dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000838152600460205260409020600101546015546c010000000000000000000000009091046bffffffffffffffffffffffff1690610ec9908290615913565b60155560008481526004602081905260409182902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16905590517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116928201929092526bffffffffffffffffffffffff831660248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b158015610fa757600080fd5b505af1158015610fbb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fdf9190614f8b565b50604080516bffffffffffffffffffffffff8316815273ffffffffffffffffffffffffffffffffffffffff8516602082015285917ff3b5906e5672f3e524854103bcafbbdba80dbdfeca2c35e116127b1060a68318910160405180910390a25050600f80547fffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffff1690555050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146110f2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611176613748565b6011546015546bffffffffffffffffffffffff90911690611198908290615913565b601555601180547fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690556040516bffffffffffffffffffffffff821681527f1d07d0b0be43d3e5fee41a80b579af370affee03fa595bf56d5d4c19328162f19060200160405180910390a16040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526bffffffffffffffffffffffff821660248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044015b602060405180830381600087803b1580156112a057600080fd5b505af11580156112b4573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061064a9190614f8b565b6112e0613748565b600f80547fffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffff166e0100000000000000000000000000001790556040513381527f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25890602001610900565b600173ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff16600381111561138557611385615a96565b141580156113cd5750600373ffffffffffffffffffffffffffffffffffffffff821660009081526016602052604090205460ff1660038111156113ca576113ca615a96565b14155b15611404576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6010546c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff16611463576040517fd12d7d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8161149a576040517f2c2fc94100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c081018290526000808567ffffffffffffffff8111156114ee576114ee615b23565b60405190808252806020026020018201604052801561152157816020015b606081526020019060019003908161150c5790505b50905060008667ffffffffffffffff81111561153f5761153f615b23565b604051908082528060200260200182016040528015611568578160200160208202803683370190505b50905060008767ffffffffffffffff81111561158657611586615b23565b60405190808252806020026020018201604052801561160b57816020015b6040805160e08101825260008082526020808301829052928201819052606082018190526080820181905260a0820181905260c082015282527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9092019101816115a45790505b50905060005b888110156119965789898281811061162b5761162b615af4565b60209081029290920135600081815260048452604090819020815160e081018352815463ffffffff8082168352640100000000820481169783019790975268010000000000000000810460ff16151593820193909352690100000000000000000090920473ffffffffffffffffffffffffffffffffffffffff166060830152600101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490931660c0840152985090965061170d9050876137cb565b8582828151811061172057611720615af4565b602002602001018190525060076000888152602001908152602001600020805461174990615987565b80601f016020809104026020016040519081016040528092919081815260200182805461177590615987565b80156117c25780601f10611797576101008083540402835291602001916117c2565b820191906000526020600020905b8154815290600101906020018083116117a557829003601f168201915b50505050508482815181106117d9576117d9615af4565b60200260200101819052506005600088815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683828151811061182a5761182a615af4565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015260a086015161186c906bffffffffffffffffffffffff16866157d3565b600088815260046020908152604080832080547fffffff000000000000000000000000000000000000000000000000000000000016815560010180547fffffffff00000000000000000000000000000000000000000000000000000000169055600790915281209196506118e09190614882565b600087815260066020526040902080547fffffffffffffffffffffffff000000000000000000000000000000000000000016905561191f600288613cd6565b5060a0860151604080516bffffffffffffffffffffffff909216825273ffffffffffffffffffffffffffffffffffffffff8a16602083015288917fb38647142fbb1ea4c000fc4569b37a4e9a9f6313317b84ee3e5326c1a6cd06ff910160405180910390a28061198e816159db565b915050611611565b50836015546119a59190615913565b6015556040516000906119c4908b908b90859088908890602001615396565b60405160208183030381529060405290508773ffffffffffffffffffffffffffffffffffffffff16638e86139b6010600001600c9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c71249ab60028c73ffffffffffffffffffffffffffffffffffffffff166348013d7b6040518163ffffffff1660e01b8152600401602060405180830381600087803b158015611a7c57600080fd5b505af1158015611a90573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ab4919061506b565b866040518463ffffffff1660e01b8152600401611ad393929190615682565b60006040518083038186803b158015611aeb57600080fd5b505afa158015611aff573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b459190810190615036565b6040518263ffffffff1660e01b8152600401611b61919061557d565b600060405180830381600087803b158015611b7b57600080fd5b505af1158015611b8f573d6000803e3d6000fd5b50506040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8b81166004830152602482018990527f000000000000000000000000000000000000000000000000000000000000000016925063a9059cbb9150604401602060405180830381600087803b158015611c2357600080fd5b505af1158015611c37573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c5b9190614f8b565b5050505050505050505050565b611c71816137cb565b600081815260046020908152604091829020825160e081018452815463ffffffff8082168352640100000000820481169483019490945268010000000000000000810460ff16158015958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a08301527801000000000000000000000000000000000000000000000000900490911660c082015290611d75576040517f514b6c2400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600082815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1668010000000000000000179055611dbf600283613cd6565b5060405182907f8ab10247ce168c27748e656ecf852b951fcaac790c18106b19aa0ae57a8b741f90600090a25050565b611df8836137cb565b6000838152601760205260409020611e119083836147e0565b50827f3e8740446213c8a77d40e08f79136ce3f347d13ed270a6ebdf57159e0faf48508383604051611e44929190615530565b60405180910390a2505050565b60023360009081526016602052604090205460ff166003811115611e7757611e77615a96565b14158015611ea9575060033360009081526016602052604090205460ff166003811115611ea657611ea6615a96565b14155b15611ee0576040517f0ebeec3c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808080611ef185870187614e82565b935093509350935060005b845181101561208457611fd3858281518110611f1a57611f1a615af4565b6020026020010151858381518110611f3457611f34615af4565b602002602001015160600151868481518110611f5257611f52615af4565b602002602001015160000151858581518110611f7057611f70615af4565b6020026020010151888681518110611f8a57611f8a615af4565b602002602001015160a00151888781518110611fa857611fa8615af4565b60200260200101518a8881518110611fc257611fc2615af4565b602002602001015160400151613893565b848181518110611fe557611fe5615af4565b60200260200101517f74931a144e43a50694897f241d973aecb5024c0e910f9bb80a163ea3c1cf5a7185838151811061202057612020615af4565b602002602001015160a001513360405161206a9291906bffffffffffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b60405180910390a28061207c816159db565b915050611efc565b50505050505050565b600082815260046020908152604091829020825160e081018452815463ffffffff80821683526401000000008204811694830185905268010000000000000000820460ff161515958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c0820152911461218f576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818160a0015161219f9190615810565b600084815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601554612205918416906157d3565b6015556040517f23b872dd0000000000000000000000000000000000000000000000000000000081523360048201523060248201526bffffffffffffffffffffffff831660448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd90606401602060405180830381600087803b1580156122a957600080fd5b505af11580156122bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122e19190614f8b565b506040516bffffffffffffffffffffffff83168152339084907fafd24114486da8ebfc32f3626dada8863652e187461aa74d4bfa7348915062039060200160405180910390a3505050565b612335836137cb565b60125474010000000000000000000000000000000000000000900463ffffffff1681111561238f576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008381526007602052604090206123a89083836147e0565b50827f7b778136e5211932b51a145badd01959415e79e051a933604b3d323f862dcabf8383604051611e44929190615530565b73ffffffffffffffffffffffffffffffffffffffff8116612428576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c6020526040902054163314612488576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600f54600b546000916124bf91859170010000000000000000000000000000000090046bffffffffffffffffffffffff1690613ce2565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260086020526040902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601554909150612529906bffffffffffffffffffffffff831690615913565b6015556040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526bffffffffffffffffffffffff831660248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90604401602060405180830381600087803b1580156125c957600080fd5b505af11580156125dd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126019190614f8b565b5060405133815273ffffffffffffffffffffffffffffffffffffffff808416916bffffffffffffffffffffffff8416918616907f9819093176a1851202c7bcfa46845809b4e47c261866550e94ed3775d2f406989060200160405180910390a4505050565b6108fc8163ffffffff16108061268f575060125463ffffffff6401000000009091048116908216115b156126c6576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6126cf826137cb565b60008281526004602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff8516908117909155915191825283917fc24c07e655ce79fba8a589778987d3c015bc6af1632bb20cf9182e02a65d972c910160405180910390a25050565b73ffffffffffffffffffffffffffffffffffffffff8181166000908152600d60205260409020541633146127a8576040517f6752e7aa00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8181166000818152600c602090815260408083208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217909355600d909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b600081815260046020908152604091829020825160e081018452815463ffffffff80821683526401000000008204811694830185905268010000000000000000820460ff161515958301959095526901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004831660c08201529114612942576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526006602052604090205473ffffffffffffffffffffffffffffffffffffffff16331461299f576040517f6352a85300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526005602090815260408083208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821790935560069094528285208054909216909155905173ffffffffffffffffffffffffffffffffffffffff90911692839186917f5cff4db96bef051785e999f44bfcd21c18823e034fb92dd376e3db4ce0feeb2c91a4505050565b612a3b613748565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b158015612ac357600080fd5b505afa158015612ad7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612afb91906150a5565b90507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb3360155484612b489190615913565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526024820152604401611286565b6000818152600460209081526040808320815160e081018352815463ffffffff80821683526401000000008204811695830186905260ff6801000000000000000083041615159483019490945273ffffffffffffffffffffffffffffffffffffffff69010000000000000000009091041660608201526001909101546bffffffffffffffffffffffff80821660808401526c0100000000000000000000000082041660a083015278010000000000000000000000000000000000000000000000009004821660c08201529291141590612c8c60005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16149050818015612cdc5750808015612cda575043836020015163ffffffff16115b155b15612d13576040517ffbc0357800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158015612d45575060008481526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314155b15612d7c576040517ffbdb8e5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b4381612d9057612d8d6032826157d3565b90505b6000858152600460205260409020805463ffffffff808416640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff90921691909117909155612de9906002908790613cd616565b5060105460808501516bffffffffffffffffffffffff9182169160009116821115612e4e576080860151612e1d908361592a565b90508560a001516bffffffffffffffffffffffff16816bffffffffffffffffffffffff161115612e4e575060a08501515b808660a00151612e5e919061592a565b600088815260046020526040902060010180547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff166c010000000000000000000000006bffffffffffffffffffffffff93841602179055601154612ec491839116615810565b601180547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff9290921691909117905560405167ffffffffffffffff84169088907f91cb3bb75cfbd718bbfccc56b7f53d92d7048ef4ca39a3b7b7c6d4af1f79118190600090a350505050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600c6020526040902054163314612fa1576040517fcebf515b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8116331415612ff1576040517f8c8728c700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8281166000908152600d602052604090205481169082161461064a5773ffffffffffffffffffffffffffffffffffffffff8281166000818152600d602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169486169485179055513392917f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836791a45050565b6130a8613748565b6130b181613f09565b50565b600060606000806000806130c6613fff565b6000600f604051806101200160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160059054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160099054906101000a900462ffffff1662ffffff1662ffffff16815260200160008201600c9054906101000a900461ffff1661ffff1661ffff16815260200160008201600e9054906101000a900460ff1615151515815260200160008201600f9054906101000a900460ff161515151581526020016000820160109054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160008201601c9054906101000a900463ffffffff1663ffffffff1663ffffffff168152505090506000600460008a81526020019081526020016000206040518060e00160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900460ff161515151581526020016000820160099054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff16815260200160018201600c9054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff166bffffffffffffffffffffffff1681526020016001820160189054906101000a900463ffffffff1663ffffffff1663ffffffff1681525050905063ffffffff8016816020015163ffffffff16146133db57505060408051602081019091526000808252965094506001935085915081905061373f565b80604001511561340a57505060408051602081019091526000808252965094506002935085915081905061373f565b61341382614037565b82516012549296509094506000916134519185917801000000000000000000000000000000000000000000000000900463ffffffff16888886614233565b9050806bffffffffffffffffffffffff168260a001516bffffffffffffffffffffffff16101561349d57600060405180602001604052806000815250600698509850985050505061373f565b5a60008b815260076020526040808220905192985090917f6e04ff0d00000000000000000000000000000000000000000000000000000000916134e291602401615590565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260608501516012549251919350600092839273ffffffffffffffffffffffffffffffffffffffff9092169163ffffffff9091169061359e9086906152ca565b60006040518083038160008787f1925050503d80600081146135dc576040519150601f19603f3d011682016040523d82523d6000602084013e6135e1565b606091505b50915091505a6135f1908a615913565b98508161361d5760006040518060200160405280600081525060039b509b509b5050505050505061373f565b6060818060200190518101906136339190614faf565b909d5090508c6136635760006040518060200160405280600081525060049c509c509c505050505050505061373f565b6012548151780100000000000000000000000000000000000000000000000090910463ffffffff1610156136b75760006040518060200160405280600081525060059c509c509c505050505050505061373f565b60405180606001604052806001436136cf9190615913565b63ffffffff1681526020016136e5600143615913565b408152602001828152506040516020016136ff91906156bc565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905260019d509b5060009a50505050505050505b91939550919395565b60005473ffffffffffffffffffffffffffffffffffffffff1633146137c9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016110e9565b565b60008181526005602052604090205473ffffffffffffffffffffffffffffffffffffffff163314613828576040517fa47c170600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081815260046020526040902054640100000000900463ffffffff908116146130b1576040517f9c0083a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061388a838361427c565b90505b92915050565b600f546e010000000000000000000000000000900460ff16156138e2576040517f24522f3400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff86163b613930576040517f09ee12d500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60125482517401000000000000000000000000000000000000000090910463ffffffff16101561398c576040517fae7235df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108fc8563ffffffff1610806139b5575060125463ffffffff6401000000009091048116908616115b156139ec576040517f14c237fb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000878152600460205260409020546901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1615613a55576040517f6e3b930b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518060e001604052808663ffffffff16815260200163ffffffff8016815260200182151581526020018773ffffffffffffffffffffffffffffffffffffffff16815260200160006bffffffffffffffffffffffff168152602001846bffffffffffffffffffffffff168152602001600063ffffffff168152506004600089815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548160ff02191690831515021790555060608201518160000160096101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160010160006101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060a082015181600101600c6101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060c08201518160010160186101000a81548163ffffffff021916908363ffffffff160217905550905050836005600089815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550826bffffffffffffffffffffffff16601554613c9e91906157d3565b60155560008781526007602090815260409091208351613cc0928501906148bc565b50613ccc60028861387e565b5050505050505050565b600061388a83836142cb565b73ffffffffffffffffffffffffffffffffffffffff831660009081526008602090815260408083208151608081018352905460ff80821615158352610100820416938201939093526bffffffffffffffffffffffff6201000084048116928201929092526e01000000000000000000000000000090920416606082018190528290613d6d908661592a565b90506000613d7b8583615854565b90508083604001818151613d8f9190615810565b6bffffffffffffffffffffffff9081169091528716606085015250613db485826158e8565b613dbe908361592a565b60118054600090613dde9084906bffffffffffffffffffffffff16615810565b825461010092830a6bffffffffffffffffffffffff81810219909216928216029190911790925573ffffffffffffffffffffffffffffffffffffffff999099166000908152600860209081526040918290208751815492890151938901516060909901517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009093169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161760ff909316909b02919091177fffffffffffff000000000000000000000000000000000000000000000000ffff1662010000878416027fffffffffffff000000000000000000000000ffffffffffffffffffffffffffff16176e010000000000000000000000000000919092160217909755509095945050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415613f89576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016110e9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b32156137c9576040517fb60ac5db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000836060015162ffffff1690506000808263ffffffff161190506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156140be57600080fd5b505afa1580156140d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906140f69190615173565b509450909250505060008113158061410d57508142105b8061412e575082801561412e57506141258242615913565b8463ffffffff16105b1561413d576013549550614141565b8095505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b1580156141a757600080fd5b505afa1580156141bb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906141df9190615173565b50945090925050506000811315806141f657508142105b806142175750828015614217575061420e8242615913565b8463ffffffff16105b1561422657601454945061422a565b8094505b50505050915091565b6000806142448689600001516143be565b905060008061425f8a8a63ffffffff16858a8a60018b614401565b909250905061426e8183615810565b9a9950505050505050505050565b60008181526001830160205260408120546142c35750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561388d565b50600061388d565b600081815260018301602052604081205480156143b45760006142ef600183615913565b855490915060009061430390600190615913565b905081811461436857600086600001828154811061432357614323615af4565b906000526020600020015490508087600001848154811061434657614346615af4565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061437957614379615ac5565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061388d565b600091505061388d565b60006143d163ffffffff8416601461587f565b6143dc8360016157eb565b6143eb9060ff16611d4c61587f565b6143f79061fde86157d3565b61388a91906157d3565b6000806000896080015161ffff168761441a919061587f565b90508380156144285750803a105b1561443057503a5b600060027f0000000000000000000000000000000000000000000000000000000000000000600281111561446657614466615a96565b14156145e95760408051600081526020810190915285156144c557600036604051806080016040528060488152602001615b71604891396040516020016144af939291906152a3565b6040516020818303038152906040529050614541565b6012546144f5907801000000000000000000000000000000000000000000000000900463ffffffff1660046158bc565b63ffffffff1667ffffffffffffffff81111561451357614513615b23565b6040519080825280601f01601f19166020018201604052801561453d576020820181803683370190505b5090505b6040517f49948e0e00000000000000000000000000000000000000000000000000000000815273420000000000000000000000000000000000000f906349948e0e9061459190849060040161557d565b60206040518083038186803b1580156145a957600080fd5b505afa1580156145bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906145e191906150a5565b9150506146a5565b60017f0000000000000000000000000000000000000000000000000000000000000000600281111561461d5761461d615a96565b14156146a557606c73ffffffffffffffffffffffffffffffffffffffff1663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561466a57600080fd5b505afa15801561467e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906146a291906150a5565b90505b846146c157808b6080015161ffff166146be919061587f565b90505b6146cf61ffff871682615840565b9050600087826146df8c8e6157d3565b6146e9908661587f565b6146f391906157d3565b61470590670de0b6b3a764000061587f565b61470f9190615840565b905060008c6040015163ffffffff1664e8d4a5100061472e919061587f565b898e6020015163ffffffff16858f88614747919061587f565b61475191906157d3565b61475f90633b9aca0061587f565b614769919061587f565b6147739190615840565b61477d91906157d3565b90506b033b2e3c9fd0803ce800000061479682846157d3565b11156147ce576040517f2ad7547a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b909c909b509950505050505050505050565b8280546147ec90615987565b90600052602060002090601f01602090048101928261480e5760008555614872565b82601f10614845578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555614872565b82800160010185558215614872579182015b82811115614872578235825591602001919060010190614857565b5061487e929150614930565b5090565b50805461488e90615987565b6000825580601f1061489e575050565b601f0160209004906000526020600020908101906130b19190614930565b8280546148c890615987565b90600052602060002090601f0160209004810192826148ea5760008555614872565b82601f1061490357805160ff1916838001178555614872565b82800160010185558215614872579182015b82811115614872578251825591602001919060010190614915565b5b8082111561487e5760008155600101614931565b803573ffffffffffffffffffffffffffffffffffffffff8116811461496957600080fd5b919050565b60008083601f84011261498057600080fd5b50813567ffffffffffffffff81111561499857600080fd5b6020830191508360208260051b85010111156149b357600080fd5b9250929050565b600082601f8301126149cb57600080fd5b813560206149e06149db83615769565b61571a565b80838252828201915082860187848660051b8901011115614a0057600080fd5b60005b85811015614a2657614a1482614945565b84529284019290840190600101614a03565b5090979650505050505050565b600082601f830112614a4457600080fd5b81356020614a546149db83615769565b80838252828201915082860187848660051b8901011115614a7457600080fd5b60005b85811015614a2657813567ffffffffffffffff811115614a9657600080fd5b8801603f81018a13614aa757600080fd5b858101356040614ab96149db8361578d565b8281528c82848601011115614acd57600080fd5b828285018a8301376000928101890192909252508552509284019290840190600101614a77565b600082601f830112614b0557600080fd5b81356020614b156149db83615769565b8281528181019085830160e080860288018501891015614b3457600080fd5b6000805b87811015614bd95782848c031215614b4e578182fd5b614b566156f1565b614b5f85614c77565b8152614b6c888601614c77565b88820152604080860135614b7f81615b62565b908201526060614b90868201614945565b908201526080614ba1868201614ca5565b9082015260a0614bb2868201614ca5565b9082015260c0614bc3868201614c77565b9082015286529486019492820192600101614b38565b50929998505050505050505050565b60008083601f840112614bfa57600080fd5b50813567ffffffffffffffff811115614c1257600080fd5b6020830191508360208285010111156149b357600080fd5b600082601f830112614c3b57600080fd5b8151614c496149db8261578d565b818152846020838601011115614c5e57600080fd5b614c6f826020830160208701615957565b949350505050565b803563ffffffff8116811461496957600080fd5b805169ffffffffffffffffffff8116811461496957600080fd5b80356bffffffffffffffffffffffff8116811461496957600080fd5b600060208284031215614cd357600080fd5b61388a82614945565b60008060408385031215614cef57600080fd5b614cf883614945565b9150614d0660208401614945565b90509250929050565b60008060408385031215614d2257600080fd5b614d2b83614945565b9150602083013560048110614d3f57600080fd5b809150509250929050565b600080600080600080600060a0888a031215614d6557600080fd5b614d6e88614945565b9650614d7c60208901614c77565b9550614d8a60408901614945565b9450606088013567ffffffffffffffff80821115614da757600080fd5b614db38b838c01614be8565b909650945060808a0135915080821115614dcc57600080fd5b50614dd98a828b01614be8565b989b979a50959850939692959293505050565b60008060208385031215614dff57600080fd5b823567ffffffffffffffff811115614e1657600080fd5b614e228582860161496e565b90969095509350505050565b600080600060408486031215614e4357600080fd5b833567ffffffffffffffff811115614e5a57600080fd5b614e668682870161496e565b9094509250614e79905060208501614945565b90509250925092565b60008060008060808587031215614e9857600080fd5b843567ffffffffffffffff80821115614eb057600080fd5b818701915087601f830112614ec457600080fd5b81356020614ed46149db83615769565b8083825282820191508286018c848660051b8901011115614ef457600080fd5b600096505b84871015614f17578035835260019690960195918301918301614ef9565b5098505088013592505080821115614f2e57600080fd5b614f3a88838901614af4565b94506040870135915080821115614f5057600080fd5b614f5c88838901614a33565b93506060870135915080821115614f7257600080fd5b50614f7f878288016149ba565b91505092959194509250565b600060208284031215614f9d57600080fd5b8151614fa881615b62565b9392505050565b60008060408385031215614fc257600080fd5b8251614fcd81615b62565b602084015190925067ffffffffffffffff811115614fea57600080fd5b614ff685828601614c2a565b9150509250929050565b6000806020838503121561501357600080fd5b823567ffffffffffffffff81111561502a57600080fd5b614e2285828601614be8565b60006020828403121561504857600080fd5b815167ffffffffffffffff81111561505f57600080fd5b614c6f84828501614c2a565b60006020828403121561507d57600080fd5b815160038110614fa857600080fd5b60006020828403121561509e57600080fd5b5035919050565b6000602082840312156150b757600080fd5b5051919050565b600080604083850312156150d157600080fd5b82359150614d0660208401614945565b6000806000604084860312156150f657600080fd5b83359250602084013567ffffffffffffffff81111561511457600080fd5b61512086828701614be8565b9497909650939450505050565b6000806040838503121561514057600080fd5b82359150614d0660208401614c77565b6000806040838503121561516357600080fd5b82359150614d0660208401614ca5565b600080600080600060a0868803121561518b57600080fd5b61519486614c8b565b94506020860151935060408601519250606086015191506151b760808701614c8b565b90509295509295909350565b600081518084526020808501945080840160005b8381101561520957815173ffffffffffffffffffffffffffffffffffffffff16875295820195908201906001016151d7565b509495945050505050565b6000815180845260208085019450848260051b860182860160005b85811015614a26578383038952615247838351615259565b9885019892509084019060010161522f565b60008151808452615271816020860160208601615957565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8284823760008382016000815283516152c0818360208801615957565b0195945050505050565b600082516152dc818460208701615957565b9190910192915050565b6000604082016040835280865480835260608501915087600052602092508260002060005b8281101561533d57815473ffffffffffffffffffffffffffffffffffffffff168452928401926001918201910161530b565b505050838103828501528481528590820160005b8681101561538a5773ffffffffffffffffffffffffffffffffffffffff61537784614945565b1682529183019190830190600101615351565b50979650505050505050565b60006080808352868184015260a07f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8811156153d157600080fd5b8760051b808a838701378085019050818101600081526020838784030181880152818a5180845260c093508385019150828c01945060005b818110156154a7578551805163ffffffff908116855285820151168585015260408082015115159085015260608082015173ffffffffffffffffffffffffffffffffffffffff1690850152888101516bffffffffffffffffffffffff168985015287810151615487898601826bffffffffffffffffffffffff169052565b5085015163ffffffff16838601529483019460e090920191600101615409565b505087810360408901526154bb818b615214565b9550505050505082810360608401526154d481856151c3565b98975050505050505050565b861515815260c0602082015260006154fb60c0830188615259565b90506007861061550d5761550d615a96565b8560408301528460608301528360808301528260a0830152979650505050505050565b60208152816020820152818360408301376000818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b60208152600061388a6020830184615259565b600060208083526000845481600182811c9150808316806155b257607f831692505b8583108114156155e9577f4e487b710000000000000000000000000000000000000000000000000000000085526022600452602485fd5b878601838152602001818015615606576001811461563557615660565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00861682528782019650615660565b60008b81526020902060005b8681101561565a57815484820152908501908901615641565b83019750505b50949998505050505050505050565b6020810161567c83615b52565b91905290565b61568b84615b52565b83815261569783615b52565b8260208201526060604082015260006156b36060830184615259565b95945050505050565b6020815263ffffffff82511660208201526020820151604082015260006040830151606080840152614c6f6080840182615259565b60405160e0810167ffffffffffffffff8111828210171561571457615714615b23565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561576157615761615b23565b604052919050565b600067ffffffffffffffff82111561578357615783615b23565b5060051b60200190565b600067ffffffffffffffff8211156157a7576157a7615b23565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082198211156157e6576157e6615a38565b500190565b600060ff821660ff84168060ff0382111561580857615808615a38565b019392505050565b60006bffffffffffffffffffffffff80831681851680830382111561583757615837615a38565b01949350505050565b60008261584f5761584f615a67565b500490565b60006bffffffffffffffffffffffff8084168061587357615873615a67565b92169190910492915050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156158b7576158b7615a38565b500290565b600063ffffffff808316818516818304811182151516156158df576158df615a38565b02949350505050565b60006bffffffffffffffffffffffff808316818516818304811182151516156158df576158df615a38565b60008282101561592557615925615a38565b500390565b60006bffffffffffffffffffffffff8381169083168181101561594f5761594f615a38565b039392505050565b60005b8381101561597257818101518382015260200161595a565b83811115615981576000848401525b50505050565b600181811c9082168061599b57607f821691505b602082108114156159d5577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615a0d57615a0d615a38565b5060010190565b600063ffffffff80831681811415615a2e57615a2e615a38565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600381106130b1576130b1615a96565b80151581146130b157600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
}

var KeeperRegistryLogicABI = KeeperRegistryLogicMetaData.ABI

var KeeperRegistryLogicBin = KeeperRegistryLogicMetaData.Bin

func DeployKeeperRegistryLogic(auth *bind.TransactOpts, backend bind.ContractBackend, paymentModel uint8, link common.Address, linkNativeFeed common.Address, fastGasFeed common.Address) (common.Address, *types.Transaction, *KeeperRegistryLogic, error) {
	parsed, err := KeeperRegistryLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeeperRegistryLogicBin), backend, paymentModel, link, linkNativeFeed, fastGasFeed)
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

func (_KeeperRegistryLogic *KeeperRegistryLogicCaller) GetPaymentModel(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _KeeperRegistryLogic.contract.Call(opts, &out, "getPaymentModel")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeeperRegistryLogic *KeeperRegistryLogicSession) GetPaymentModel() (uint8, error) {
	return _KeeperRegistryLogic.Contract.GetPaymentModel(&_KeeperRegistryLogic.CallOpts)
}

func (_KeeperRegistryLogic *KeeperRegistryLogicCallerSession) GetPaymentModel() (uint8, error) {
	return _KeeperRegistryLogic.Contract.GetPaymentModel(&_KeeperRegistryLogic.CallOpts)
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

	GetPaymentModel(opts *bind.CallOpts) (uint8, error)

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
