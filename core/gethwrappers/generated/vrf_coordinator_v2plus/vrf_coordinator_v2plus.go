// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2plus

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

type VRFCoordinatorV2PlusRequestCommitment struct {
	BlockNum         uint64
	SubId            *big.Int
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
	ExtraArgs        []byte
}

type VRFProof struct {
	Pk            [2]*big.Int
	Gamma         [2]*big.Int
	C             *big.Int
	S             *big.Int
	Seed          *big.Int
	UWitness      common.Address
	CGammaWitness [2]*big.Int
	SHashWitness  [2]*big.Int
	ZInv          *big.Int
}

type VRFV2PlusClientRandomWordsRequest struct {
	KeyHash              [32]byte
	SubId                *big.Int
	RequestConfirmations uint16
	CallbackGasLimit     uint32
	NumWords             uint32
	ExtraArgs            []byte
}

var VRFCoordinatorV2PlusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendEther\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EthFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountEth\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldEthBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newEthBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithEth\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRICE_REGISTRY\",\"outputs\":[{\"internalType\":\"contractIVRFV2PlusPriceRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"USD_FEE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFCoordinatorV2Plus.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"ethBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdrawEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverEthFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalEthBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"name\":\"setLINK\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620061a5380380620061a58339810160408190526200003491620001a9565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e0565b5050506001600160601b0319606092831b8116608052911b1660a052620001e1565b6001600160a01b0381163314156200013b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001a457600080fd5b919050565b60008060408385031215620001bd57600080fd5b620001c8836200018c565b9150620001d8602084016200018c565b90509250929050565b60805160601c60a05160601c615f8a6200021b6000396000818161047501526112850152600081816105f30152613d380152615f8a6000f3fe6080604052600436106102f15760003560e01c80638da5cb5b1161018f578063bec4c08c116100e1578063dc311dd31161008a578063e95704bd11610064578063e95704bd14610921578063ee9d2d3814610951578063f2fde38b1461097e57600080fd5b8063dc311dd3146108bd578063e72f6e30146108ee578063e8509bff1461090e57600080fd5b8063d98e620e116100bb578063d98e620e14610847578063da2f261014610867578063dac83d291461089d57600080fd5b8063bec4c08c146107e7578063caf70c4a14610807578063cb6317971461082757600080fd5b8063a4c0ed3611610143578063aefb212f1161011d578063aefb212f14610785578063b2a7cac5146107b2578063bb0697a5146107d257600080fd5b8063a4c0ed3614610725578063a8cb447b14610745578063aa433aff1461076557600080fd5b80639d40a6fd116101745780639d40a6fd146106b6578063a02e0616146106f0578063a21a23e41461071057600080fd5b80638da5cb5b1461066a5780639b1c385e1461068857600080fd5b8063405b84fa1161024857806364d51a2a116101fc5780636f64f03f116101d65780636f64f03f1461061557806379ba50971461063557806386fe91c71461064a57600080fd5b806364d51a2a146105ac57806366316d8d146105c1578063689c4517146105e157600080fd5b806341af6c871161022d57806341af6c871461053c57806346d8d4861461056c5780635d06b4ab1461058c57600080fd5b8063405b84fa146104f157806340d6bb821461051157600080fd5b80630fe9842a116102aa5780631ed25726116102845780631ed2572614610463578063294daa4914610497578063330987b3146104b957600080fd5b80630fe9842a146103e357806315c48b84146104035780631b6b6d231461042b57600080fd5b8063088070f5116102db578063088070f51461034557806308821d58146103a35780630ae09540146103c357600080fd5b8062012291146102f657806304104edb14610323575b600080fd5b34801561030257600080fd5b5061030b61099e565b60405161031a93929190615ace565b60405180910390f35b34801561032f57600080fd5b5061034361033e3660046154ea565b610a1a565b005b34801561035157600080fd5b50600c5461037c9061ffff81169062010000810463ffffffff16906601000000000000900460ff1683565b6040805161ffff909416845263ffffffff909216602084015215159082015260600161031a565b3480156103af57600080fd5b506103436103be3660046155ff565b610bf5565b3480156103cf57600080fd5b506103436103de366004615830565b610da2565b3480156103ef57600080fd5b506103436103fe366004615806565b610e8c565b34801561040f57600080fd5b5061041860c881565b60405161ffff909116815260200161031a565b34801561043757600080fd5b5060025461044b906001600160a01b031681565b6040516001600160a01b03909116815260200161031a565b34801561046f57600080fd5b5061044b7f000000000000000000000000000000000000000000000000000000000000000081565b3480156104a357600080fd5b5060015b60405160ff909116815260200161031a565b3480156104c557600080fd5b506104d96104d43660046156d2565b610f87565b6040516001600160601b03909116815260200161031a565b3480156104fd57600080fd5b5061034361050c366004615830565b611584565b34801561051d57600080fd5b506105276101f481565b60405163ffffffff909116815260200161031a565b34801561054857600080fd5b5061055c610557366004615654565b61197d565b604051901515815260200161031a565b34801561057857600080fd5b50610343610587366004615507565b611b7f565b34801561059857600080fd5b506103436105a73660046154ea565b611d18565b3480156105b857600080fd5b50610418606481565b3480156105cd57600080fd5b506103436105dc366004615507565b611def565b3480156105ed57600080fd5b5061044b7f000000000000000000000000000000000000000000000000000000000000000081565b34801561062157600080fd5b50610343610630366004615540565b611f91565b34801561064157600080fd5b506103436120a9565b34801561065657600080fd5b506009546104d9906001600160601b031681565b34801561067657600080fd5b506000546001600160a01b031661044b565b34801561069457600080fd5b506106a86106a33660046157b0565b61215a565b60405190815260200161031a565b3480156106c257600080fd5b506006546106d79067ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161031a565b3480156106fc57600080fd5b5061034361070b3660046154ea565b6125b9565b34801561071c57600080fd5b506106a8612626565b34801561073157600080fd5b50610343610740366004615576565b6128c8565b34801561075157600080fd5b506103436107603660046154ea565b612a9d565b34801561077157600080fd5b50610343610780366004615654565b612bf3565b34801561079157600080fd5b506107a56107a0366004615855565b612c53565b60405161031a9190615a33565b3480156107be57600080fd5b506103436107cd366004615654565b612d6e565b3480156107de57600080fd5b506104a7600881565b3480156107f357600080fd5b50610343610802366004615830565b612ec1565b34801561081357600080fd5b506106a861082236600461561b565b61307a565b34801561083357600080fd5b50610343610842366004615830565b6130aa565b34801561085357600080fd5b506106a8610862366004615654565b6133e3565b34801561087357600080fd5b5061044b610882366004615654565b600d602052600090815260409020546001600160a01b031681565b3480156108a957600080fd5b506103436108b8366004615830565b613404565b3480156108c957600080fd5b506108dd6108d8366004615654565b613526565b60405161031a959493929190615c4c565b3480156108fa57600080fd5b506103436109093660046154ea565b61362b565b61034361091c366004615654565b613800565b34801561092d57600080fd5b506009546104d9906c0100000000000000000000000090046001600160601b031681565b34801561095d57600080fd5b506106a861096c366004615654565b600f6020526000908152604090205481565b34801561098a57600080fd5b506103436109993660046154ea565b61394b565b600c54600e805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff16939192839190830182828015610a0857602002820191906000526020600020905b8154815260200190600101908083116109f4575b50505050509050925092509250909192565b610a2261395c565b60105460005b81811015610baf57826001600160a01b031660108281548110610a4d57610a4d615f19565b6000918252602090912001546001600160a01b03161415610b9d576010610a75600184615deb565b81548110610a8557610a85615f19565b600091825260209091200154601080546001600160a01b039092169183908110610ab157610ab1615f19565b600091825260209091200180546001600160a01b0319166001600160a01b0392909216919091179055826010610ae8600185615deb565b81548110610af857610af8615f19565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506010805480610b3757610b37615f03565b6000828152602090819020600019908301810180546001600160a01b03191690559091019091556040516001600160a01b03851681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37910160405180910390a1505050565b80610ba781615e72565b915050610a28565b506040517f5428d4490000000000000000000000000000000000000000000000000000000081526001600160a01b03831660048201526024015b60405180910390fd5b50565b610bfd61395c565b604080518082018252600091610c2c91908490600290839083908082843760009201919091525061307a915050565b6000818152600d60205260409020549091506001600160a01b031680610c81576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101839052602401610be9565b6000828152600d6020526040812080546001600160a01b03191690555b600e54811015610d595782600e8281548110610cbc57610cbc615f19565b90600052602060002001541415610d4757600e805460009190610ce190600190615deb565b81548110610cf157610cf1615f19565b9060005260206000200154905080600e8381548110610d1257610d12615f19565b600091825260209091200155600e805480610d2f57610d2f615f03565b60019003818190600052602060002001600090559055505b80610d5181615e72565b915050610c9e565b50806001600160a01b03167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610d9591815260200190565b60405180910390a2505050565b60008281526004602052604090205482906001600160a01b031680610dda57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614610e0e57604051636c51fda960e11b81526001600160a01b0382166004820152602401610be9565b600c546601000000000000900460ff1615610e3c5760405163769dd35360e11b815260040160405180910390fd5b610e458461197d565b15610e7c576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e8684846139b8565b50505050565b610e9461395c565b60c861ffff83161115610ee7576040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff831660048201819052602482015260c86044820152606401610be9565b6040805160608101825261ffff841680825263ffffffff84166020808401829052600093850193909352600c80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000001683176201000083021766ff000000000000191690558351918252918101919091527f52b99d89f788ddd6c815037dcb52b920032f7ad04d1ed6655c5658021a2b1341910160405180910390a15050565b600c546000906601000000000000900460ff1615610fb85760405163769dd35360e11b815260040160405180910390fd5b60005a90506000610fc98585613b68565b90506000846060015163ffffffff1667ffffffffffffffff811115610ff057610ff0615f2f565b604051908082528060200260200182016040528015611019578160200160208202803683370190505b50905060005b856060015163ffffffff1681101561109957826040015181604051602001611051929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c82828151811061107c5761107c615f19565b60209081029190910101528061109181615e72565b91505061101f565b50602080830180516000908152600f90925260408083208390559051905182917f1fe543e300000000000000000000000000000000000000000000000000000000916110ea91908690602401615b41565b60408051601f198184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090941693909317909252600c805466ff000000000000191666010000000000001790559088015160808901519192506000916111829163ffffffff169084613e74565b600c805466ff00000000000019169055602089810151600090815260059091526040902054909150600160c01b900467ffffffffffffffff166111c6816001615d9d565b6020808b01516000908152600590915260408120805467ffffffffffffffff93909316600160c01b0277ffffffffffffffffffffffffffffffffffffffffffffffff9093169290921790915560a08a0151805161122590600190615deb565b8151811061123557611235615f19565b6020910101516040517fe6152d81000000000000000000000000000000000000000000000000000000008152600481018a90523a602482015260f89190911c6001146044820181905291506000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e6152d819060640160206040518083038186803b1580156112cf57600080fd5b505afa1580156112e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113079190615892565b90508115611422576020808c01516000908152600590915260409020546001600160601b038083166c0100000000000000000000000090920416101561136057604051631e9acf1760e31b815260040160405180910390fd5b60208b81015160009081526005909152604090208054829190600c906113a09084906c0100000000000000000000000090046001600160601b0316615e02565b82546101009290920a6001600160601b0381810219909316918316021790915589516000908152600d60209081526040808320546001600160a01b03168352600b9091528120805485945090926113f991859116615dc9565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555061150e565b6020808c01516000908152600590915260409020546001600160601b038083169116101561146357604051631e9acf1760e31b815260040160405180910390fd5b6020808c0151600090815260059091526040812080548392906114909084906001600160601b0316615e02565b82546101009290920a6001600160601b0381810219909316918316021790915589516000908152600d60209081526040808320546001600160a01b03168352600a9091528120805485945090926114e991859116615dc9565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b8a6020015188602001517f49580fdfd9497e1ed5c1b1cec0495087ae8e3f1267470ec2fb015db32e3d6aa78a60400151848860405161156b939291909283526001600160601b039190911660208301521515604082015260600190565b60405180910390a3985050505050505050505b92915050565b61158d81613ec2565b6115ce576040517f5428d4490000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602401610be9565b6000806000806115dd86613526565b945094505093509350336001600160a01b0316826001600160a01b0316146116475760405162461bcd60e51b815260206004820152601660248201527f4e6f7420737562736372697074696f6e206f776e6572000000000000000000006044820152606401610be9565b6116508661197d565b1561169d5760405162461bcd60e51b815260206004820152601660248201527f50656e64696e67207265717565737420657869737473000000000000000000006044820152606401610be9565b60006040518060c001604052806116b2600190565b60ff168152602001888152602001846001600160a01b03168152602001838152602001866001600160601b03168152602001856001600160601b031681525090506000816040516020016117069190615a59565b604051602081830303815290604052905061172088613f2c565b50506040517fce3f47190000000000000000000000000000000000000000000000000000000081526001600160a01b0388169063ce3f4719906001600160601b03881690611772908590600401615a46565b6000604051808303818588803b15801561178b57600080fd5b505af115801561179f573d6000803e3d6000fd5b505060025460405163a9059cbb60e01b81526001600160a01b038c811660048301526001600160601b038c166024830152909116935063a9059cbb92506044019050602060405180830381600087803b1580156117fb57600080fd5b505af115801561180f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118339190615637565b61187f5760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610be9565b60005b83518110156119305783818151811061189d5761189d615f19565b60209081029190910101516040517f8ea981170000000000000000000000000000000000000000000000000000000081526001600160a01b038a8116600483015290911690638ea9811790602401600060405180830381600087803b15801561190557600080fd5b505af1158015611919573d6000803e3d6000fd5b50505050808061192890615e72565b915050611882565b50604080516001600160a01b0389168152602081018a90527fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187910160405180910390a15050505050505050565b6000818152600460209081526040808320815160608101835281546001600160a01b0390811682526001830154168185015260028201805484518187028101870186528181528796939586019390929190830182828015611a0757602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116119e9575b505050505081525050905060005b816040015151811015611b755760005b600e54811015611b62576000611b2b600e8381548110611a4757611a47615f19565b906000526020600020015485604001518581518110611a6857611a68615f19565b6020026020010151886003600089604001518981518110611a8b57611a8b615f19565b6020908102919091018101516001600160a01b03908116835282820193909352604091820160009081208e825282528290205482518083018890529590931685830152606085019390935267ffffffffffffffff9091166080808501919091528151808503909101815260a08401825280519083012060c084019490945260e0808401859052815180850390910181526101009093019052815191012091565b506000818152600f602052604090205490915015611b4f5750600195945050505050565b5080611b5a81615e72565b915050611a25565b5080611b6d81615e72565b915050611a15565b5060009392505050565b600c546601000000000000900460ff1615611bad5760405163769dd35360e11b815260040160405180910390fd5b336000908152600b60205260409020546001600160601b0380831691161015611be957604051631e9acf1760e31b815260040160405180910390fd5b336000908152600b602052604081208054839290611c119084906001600160601b0316615e02565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550806009600c8282829054906101000a90046001600160601b0316611c599190615e02565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114611cd3576040519150601f19603f3d011682016040523d82523d6000602084013e611cd8565b606091505b5050905080611d13576040517fdcf35db000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050565b611d2061395c565b611d2981613ec2565b15611d6b576040517fac8a27ef0000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602401610be9565b601080546001810182556000919091527f1b6847dc741a1b0cd08d278845f9d819d87b734759afb55fe2de5cb82a9ae6720180546001600160a01b0319166001600160a01b0383169081179091556040519081527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259060200160405180910390a150565b600c546601000000000000900460ff1615611e1d5760405163769dd35360e11b815260040160405180910390fd5b336000908152600a60205260409020546001600160601b0380831691161015611e5957604051631e9acf1760e31b815260040160405180910390fd5b336000908152600a602052604081208054839290611e819084906001600160601b0316615e02565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600960008282829054906101000a90046001600160601b0316611ec99190615e02565b82546101009290920a6001600160601b0381810219909316918316021790915560025460405163a9059cbb60e01b81526001600160a01b03868116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b158015611f3857600080fd5b505af1158015611f4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f709190615637565b611f8d57604051631e9acf1760e31b815260040160405180910390fd5b5050565b611f9961395c565b604080518082018252600091611fc891908490600290839083908082843760009201919091525061307a915050565b6000818152600d60205260409020549091506001600160a01b03161561201d576040517f4a0b8fa700000000000000000000000000000000000000000000000000000000815260048101829052602401610be9565b6000818152600d6020908152604080832080546001600160a01b0319166001600160a01b038816908117909155600e805460018101825594527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610d95565b6001546001600160a01b031633146121035760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610be9565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600c546000906601000000000000900460ff161561218b5760405163769dd35360e11b815260040160405180910390fd5b6020808301356000908152600490915260409020546001600160a01b03166121c657604051630fb532db60e11b815260040160405180910390fd5b33600090815260036020908152604080832085830135845290915290205467ffffffffffffffff1680612231576040517f79bfd40100000000000000000000000000000000000000000000000000000000815260208401356004820152336024820152604401610be9565b600c5461ffff1661224860608501604086016157eb565b61ffff16108061226b575060c861226560608501604086016157eb565b61ffff16115b156122ca5761228060608401604085016157eb565b600c546040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff92831660048201529116602482015260c86044820152606401610be9565b600c5462010000900463ffffffff166122e96080850160608601615877565b63ffffffff161115612352576123056080840160608501615877565b600c546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff9283166004820152620100009091049091166024820152604401610be9565b6101f461236560a0850160808601615877565b63ffffffff1611156123c45761238160a0840160808501615877565b6040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff90911660048201526101f46024820152604401610be9565b60006123d1826001615d9d565b6040805186356020808301829052338385015280890135606084015267ffffffffffffffff85166080808501919091528451808503909101815260a0808501865281519183019190912060c085019390935260e0808501849052855180860390910181526101009094019094528251920191909120929350906000906124629061245d90890189615ca2565b614185565b9050600061246f8261424c565b90508361247a6142ea565b60208a013561248f60808c0160608d01615877565b61249f60a08d0160808e01615877565b33866040516020016124b79796959493929190615ba1565b60405160208183030381529060405280519060200120600f600086815260200190815260200160002081905550336001600160a01b0316886020013589600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e87878d604001602081019061252e91906157eb565b8e60600160208101906125419190615877565b8f60800160208101906125549190615877565b8960405161256796959493929190615b62565b60405180910390a450503360009081526003602090815260408083208983013584529091529020805467ffffffffffffffff191667ffffffffffffffff9490941693909317909255925050505b919050565b6125c161395c565b6002546001600160a01b031615612604576040517f2d118a6e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0392909216919091179055565b600c546000906601000000000000900460ff16156126575760405163769dd35360e11b815260040160405180910390fd5b600033612665600143615deb565b600654604051606093841b6bffffffffffffffffffffffff199081166020830152924060348201523090931b909116605483015260c01b7fffffffffffffffff00000000000000000000000000000000000000000000000016606882015260700160408051601f1981840301815291905280516020909101206006805491925067ffffffffffffffff9091169060006126fd83615e8d565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505060008067ffffffffffffffff81111561273f5761273f615f2f565b604051908082528060200260200182016040528015612768578160200160208202803683370190505b506040805160608082018352600080835260208084018281528486018381528984526005835286842095518654925191516001600160601b039182167fffffffffffffffff000000000000000000000000000000000000000000000000909416939093176c0100000000000000000000000091909216021777ffffffffffffffffffffffffffffffffffffffffffffffff16600160c01b67ffffffffffffffff9092169190910217909355835191820184523382528183018181528285018681528883526004855294909120825181546001600160a01b03199081166001600160a01b03928316178355925160018301805490941691161790915592518051949550909361287c926002850192019061520f565b5061288c91506007905083614383565b5060405133815282907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d9060200160405180910390a250905090565b600c546601000000000000900460ff16156128f65760405163769dd35360e11b815260040160405180910390fd5b6002546001600160a01b0316331461293a576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612974576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600061298282840184615654565b6000818152600460205260409020549091506001600160a01b03166129ba57604051630fb532db60e11b815260040160405180910390fd5b600081815260056020526040812080546001600160601b0316918691906129e18385615dc9565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600960008282829054906101000a90046001600160601b0316612a299190615dc9565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a828784612a7c9190615d85565b604080519283526020830191909152015b60405180910390a2505050505050565b612aa561395c565b60095447906c0100000000000000000000000090046001600160601b031681811115612b07576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610be9565b81811015611d13576000612b1b8284615deb565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d8060008114612b6a576040519150601f19603f3d011682016040523d82523d6000602084013e612b6f565b606091505b5050905080612baa576040517fdcf35db000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df377317910160405180910390a15050505050565b612bfb61395c565b6000818152600460205260409020546001600160a01b0316612c3057604051630fb532db60e11b815260040160405180910390fd5b600081815260046020526040902054610bf29082906001600160a01b03166139b8565b60606000612c61600761438f565b9050808410612c9c576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612ca88486615d85565b905081811180612cb6575083155b612cc05780612cc2565b815b90506000612cd08683615deb565b67ffffffffffffffff811115612ce857612ce8615f2f565b604051908082528060200260200182016040528015612d11578160200160208202803683370190505b50905060005b8151811015612d6457612d35612d2d8883615d85565b600790614399565b828281518110612d4757612d47615f19565b602090810291909101015280612d5c81615e72565b915050612d17565b5095945050505050565b600c546601000000000000900460ff1615612d9c5760405163769dd35360e11b815260040160405180910390fd5b6000818152600460205260409020546001600160a01b0316612dd157604051630fb532db60e11b815260040160405180910390fd5b6000818152600460205260409020600101546001600160a01b03163314612e44576000818152600460208190526040918290206001015491517fd084e9750000000000000000000000000000000000000000000000000000000081526001600160a01b0390921690820152602401610be9565b6000818152600460209081526040918290208054336001600160a01b0319808316821784556001909301805490931690925583516001600160a01b0390911680825292810191909152909183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691015b60405180910390a25050565b60008281526004602052604090205482906001600160a01b031680612ef957604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612f2d57604051636c51fda960e11b81526001600160a01b0382166004820152602401610be9565b600c546601000000000000900460ff1615612f5b5760405163769dd35360e11b815260040160405180910390fd5b60008481526004602052604090206002015460641415612fa7576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b038316600090815260036020908152604080832087845290915290205467ffffffffffffffff1615612fdf57610e86565b6001600160a01b03831660008181526003602090815260408083208884528252808320805467ffffffffffffffff19166001908117909155600483528184206002018054918201815584529282902090920180546001600160a01b03191684179055905191825285917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a250505050565b60008160405160200161308d9190615a25565b604051602081830303815290604052805190602001209050919050565b60008281526004602052604090205482906001600160a01b0316806130e257604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461311657604051636c51fda960e11b81526001600160a01b0382166004820152602401610be9565b600c546601000000000000900460ff16156131445760405163769dd35360e11b815260040160405180910390fd5b61314d8461197d565b15613184576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b038316600090815260036020908152604080832087845290915290205467ffffffffffffffff166131fa576040517f79bfd401000000000000000000000000000000000000000000000000000000008152600481018590526001600160a01b0384166024820152604401610be9565b60008481526004602090815260408083206002018054825181850281018501909352808352919290919083018282801561325d57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161323f575b505050505090506000600182516132749190615deb565b905060005b825181101561338057856001600160a01b031683828151811061329e5761329e615f19565b60200260200101516001600160a01b0316141561336e5760008383815181106132c9576132c9615f19565b6020026020010151905080600460008a815260200190815260200160002060020183815481106132fb576132fb615f19565b600091825260208083209190910180546001600160a01b0319166001600160a01b03949094169390931790925589815260049091526040902060020180548061334657613346615f03565b600082815260209020810160001990810180546001600160a01b031916905501905550613380565b8061337881615e72565b915050613279565b506001600160a01b03851660008181526003602090815260408083208a8452825291829020805467ffffffffffffffff19169055905191825287917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79101612a8d565b600e81815481106133f357600080fd5b600091825260209091200154905081565b60008281526004602052604090205482906001600160a01b03168061343c57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461347057604051636c51fda960e11b81526001600160a01b0382166004820152602401610be9565b600c546601000000000000900460ff161561349e5760405163769dd35360e11b815260040160405180910390fd5b6000848152600460205260409020600101546001600160a01b03848116911614610e865760008481526004602090815260409182902060010180546001600160a01b0319166001600160a01b03871690811790915582513381529182015285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1910161306c565b6000818152600460205260408120548190819081906060906001600160a01b031661356457604051630fb532db60e11b815260040160405180910390fd5b60008681526005602090815260408083205460048352928190208054600290910180548351818602810186019094528084526001600160601b03808716966c01000000000000000000000000810490911695600160c01b90910467ffffffffffffffff16946001600160a01b039094169391839183018282801561361157602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116135f3575b505050505090509450945094509450945091939590929450565b61363361395c565b6002546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000916001600160a01b0316906370a082319060240160206040518083038186803b15801561369057600080fd5b505afa1580156136a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136c8919061566d565b6009549091506001600160601b03168181111561371b576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610be9565b81811015611d1357600061372f8284615deb565b60025460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb90604401602060405180830381600087803b15801561377f57600080fd5b505af1158015613793573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137b79190615637565b50604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b600c546601000000000000900460ff161561382e5760405163769dd35360e11b815260040160405180910390fd5b6000818152600460205260409020546001600160a01b031661386357604051630fb532db60e11b815260040160405180910390fd5b600081815260056020526040902080546c0100000000000000000000000090046001600160601b0316903490600c61389b8385615dc9565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550346009600c8282829054906101000a90046001600160601b03166138e39190615dc9565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f3f1ddc3ab1bdb39001ad76ca51a0e6f57ce6627c69f251d1de41622847721cde8234846139369190615d85565b60408051928352602083019190915201612eb5565b61395361395c565b610bf2816143a5565b6000546001600160a01b031633146139b65760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610be9565b565b6000806139c484613f2c565b60025460405163a9059cbb60e01b81526001600160a01b0387811660048301526001600160601b0385166024830152939550919350919091169063a9059cbb90604401602060405180830381600087803b158015613a2157600080fd5b505af1158015613a35573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a599190615637565b613a7657604051631e9acf1760e31b815260040160405180910390fd5b6000836001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114613acc576040519150601f19603f3d011682016040523d82523d6000602084013e613ad1565b606091505b5050905080613b0c576040517fdcf35db000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516001600160a01b03861681526001600160601b038581166020830152841681830152905186917f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4919081900360600190a25050505050565b60408051606081018252600080825260208201819052918101919091526000613b94846000015161307a565b6000818152600d60205260409020549091506001600160a01b031680613be9576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101839052602401610be9565b6000828660800151604051602001613c0b929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600f90935291205490915080613c6a576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85516020808801516040808a015160608b015160808c015160a08d01519351613c99978a979096959101615bf8565b604051602081830303815290604052805190602001208114613ce7576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613cf6876000015161444f565b905080613e025786516040517fe9413d3800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e9413d389060240160206040518083038186803b158015613d8257600080fd5b505afa158015613d96573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613dba919061566d565b905080613e025786516040517f175dadad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610be9565b6000886080015182604051602001613e24929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c90506000613e4b8a83614556565b604080516060810182529889526020890196909652948701949094525093979650505050505050565b60005a611388811015613e8657600080fd5b611388810390508460408204820311613e9e57600080fd5b50823b613eaa57600080fd5b60008083516020850160008789f190505b9392505050565b6000805b601054811015613f2357826001600160a01b031660108281548110613eed57613eed615f19565b6000918252602090912001546001600160a01b03161415613f115750600192915050565b80613f1b81615e72565b915050613ec6565b50600092915050565b6000818152600460209081526040808320815160608101835281546001600160a01b03908116825260018301541681850152600282018054845181870281018701865281815287968796949594860193919290830182828015613fb857602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613f9a575b505050919092525050506000858152600560209081526040808320815160608101835290546001600160601b038082168084526c010000000000000000000000008304909116948301859052600160c01b90910467ffffffffffffffff16928201929092529096509094509192505b82604001515181101561409f57600360008460400151838151811061404e5761404e615f19565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208982529092529020805467ffffffffffffffff191690558061409781615e72565b915050614027565b50600085815260046020526040812080546001600160a01b031990811682556001820180549091169055906140d76002830182615274565b50506000858152600560205260408120556140f36007866145c1565b50600980548591906000906141129084906001600160601b0316615e02565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550826009600c8282829054906101000a90046001600160601b031661415a9190615e02565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050915091565b604080516020810190915260008152816141ae575060408051602081019091526000815261157e565b7f92fd1338000000000000000000000000000000000000000000000000000000006141d98385615e2a565b7fffffffff000000000000000000000000000000000000000000000000000000001614614232576040517f5247fdce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61423f8260048186615d5b565b810190613ebb9190615686565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161428591511515815260200190565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b60004661a4b18114806142ff575062066eed81145b1561437c5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561433e57600080fd5b505afa158015614352573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614376919061566d565b91505090565b4391505090565b6000613ebb83836145cd565b600061157e825490565b6000613ebb838361461c565b6001600160a01b0381163314156143fe5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610be9565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b1811480614464575062066eed81145b15614546576101008367ffffffffffffffff1661447f6142ea565b6144899190615deb565b11806144a657506144986142ea565b8367ffffffffffffffff1610155b156144b45750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a829060240160206040518083038186803b15801561450e57600080fd5b505afa158015614522573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ebb919061566d565b505067ffffffffffffffff164090565b600061458a8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151614646565b600383602001516040516020016145a2929190615b2d565b60408051601f1981840301815291905280516020909101209392505050565b6000613ebb8383614881565b60008181526001830160205260408120546146145750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561157e565b50600061157e565b600082600001828154811061463357614633615f19565b9060005260206000200154905092915050565b61464f89614974565b61469b5760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610be9565b6146a488614974565b6146f05760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610be9565b6146f983614974565b6147455760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610be9565b61474e82614974565b61479a5760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610be9565b6147a6878a8887614a4d565b6147f25760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610be9565b60006147fe8a87614b82565b90506000614811898b878b868989614be6565b90506000614822838d8d8a86614d12565b9050808a146148735760405162461bcd60e51b815260206004820152600d60248201527f696e76616c69642070726f6f66000000000000000000000000000000000000006044820152606401610be9565b505050505050505050505050565b6000818152600183016020526040812054801561496a5760006148a5600183615deb565b85549091506000906148b990600190615deb565b905081811461491e5760008660000182815481106148d9576148d9615f19565b90600052602060002001549050808760000184815481106148fc576148fc615f19565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061492f5761492f615f03565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061157e565b600091505061157e565b80516000906401000003d019116149cd5760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610be9565b60208201516401000003d01911614a265760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610be9565b60208201516401000003d019908009614a468360005b6020020151614d52565b1492915050565b60006001600160a01b038216614aa55760405162461bcd60e51b815260206004820152600b60248201527f626164207769746e6573730000000000000000000000000000000000000000006044820152606401610be9565b602084015160009060011615614abc57601c614abf565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa158015614b5a573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b614b8a615292565b614bb760018484604051602001614ba393929190615a04565b604051602081830303815290604052614d76565b90505b614bc381614974565b61157e578051604080516020810192909252614bdf9101614ba3565b9050614bba565b614bee615292565b825186516401000003d0199081900691061415614c4d5760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610be9565b614c58878988614dc4565b614ca45760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610be9565b614caf848685614dc4565b614cfb5760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610be9565b614d06868484614efe565b98975050505050505050565b600060028686868587604051602001614d30969594939291906159a5565b60408051601f1981840301815291905280516020909101209695505050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614d7e615292565b614d8782614fc5565b8152614d9c614d97826000614a3c565b615000565b6020820181905260029006600114156125b4576020810180516401000003d019039052919050565b600082614e135760405162461bcd60e51b815260206004820152600b60248201527f7a65726f207363616c61720000000000000000000000000000000000000000006044820152606401610be9565b83516020850151600090614e2990600290615eb5565b15614e3557601c614e38565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614eaa573d6000803e3d6000fd5b505050602060405103519050600086604051602001614ec99190615993565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614f06615292565b835160208086015185519186015160009384938493614f2793909190615020565b919450925090506401000003d019858209600114614f875760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610be9565b60405180604001604052806401000003d01980614fa657614fa6615eed565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d01981106125b457604080516020808201939093528151808203840181529082019091528051910120614fcd565b600061157e8260026150196401000003d0196001615d85565b901c615100565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000615060838385856151a2565b909850905061507188828e886151c6565b909850905061508288828c876151c6565b909850905060006150958d878b856151c6565b90985090506150a6888286866151a2565b90985090506150b788828e896151c6565b90985090508181146150ec576401000003d019818a0998506401000003d01982890997506401000003d01981830996506150f0565b8196505b5050505050509450945094915050565b60008061510b6152b0565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a082015261513d6152ce565b60208160c0846005600019fa9250826151985760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610be9565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215615264579160200282015b8281111561526457825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019061522f565b506152709291506152ec565b5090565b5080546000825590600052602060002090810190610bf291906152ec565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b8082111561527057600081556001016152ed565b80356125b481615f45565b806040810183101561157e57600080fd5b600082601f83011261532e57600080fd5b6040516040810181811067ffffffffffffffff8211171561535157615351615f2f565b806040525080838560408601111561536857600080fd5b60005b600281101561538a57813583526020928301929091019060010161536b565b509195945050505050565b600082601f8301126153a657600080fd5b813567ffffffffffffffff808211156153c1576153c1615f2f565b604051601f8301601f19908116603f011681019082821181831017156153e9576153e9615f2f565b8160405283815286602085880101111561540257600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060c0828403121561543457600080fd5b61543c615d0e565b9050813567ffffffffffffffff808216821461545757600080fd5b81835260208401356020840152615470604085016154d6565b6040840152615481606085016154d6565b606084015261549260808501615301565b608084015260a08401359150808211156154ab57600080fd5b506154b884828501615395565b60a08301525092915050565b803561ffff811681146125b457600080fd5b803563ffffffff811681146125b457600080fd5b6000602082840312156154fc57600080fd5b8135613ebb81615f45565b6000806040838503121561551a57600080fd5b823561552581615f45565b9150602083013561553581615f68565b809150509250929050565b6000806060838503121561555357600080fd5b823561555e81615f45565b915061556d846020850161530c565b90509250929050565b6000806000806060858703121561558c57600080fd5b843561559781615f45565b935060208501359250604085013567ffffffffffffffff808211156155bb57600080fd5b818701915087601f8301126155cf57600080fd5b8135818111156155de57600080fd5b8860208285010111156155f057600080fd5b95989497505060200194505050565b60006040828403121561561157600080fd5b613ebb838361530c565b60006040828403121561562d57600080fd5b613ebb838361531d565b60006020828403121561564957600080fd5b8151613ebb81615f5a565b60006020828403121561566657600080fd5b5035919050565b60006020828403121561567f57600080fd5b5051919050565b60006020828403121561569857600080fd5b6040516020810181811067ffffffffffffffff821117156156bb576156bb615f2f565b60405282356156c981615f5a565b81529392505050565b6000808284036101c08112156156e757600080fd5b6101a0808212156156f757600080fd5b6156ff615d37565b915061570b868661531d565b825261571a866040870161531d565b60208301526080850135604083015260a0850135606083015260c0850135608083015261574960e08601615301565b60a083015261010061575d8782880161531d565b60c084015261577087610140880161531d565b60e08401526101808601359083015290925083013567ffffffffffffffff81111561579a57600080fd5b6157a685828601615422565b9150509250929050565b6000602082840312156157c257600080fd5b813567ffffffffffffffff8111156157d957600080fd5b820160c08185031215613ebb57600080fd5b6000602082840312156157fd57600080fd5b613ebb826154c4565b6000806040838503121561581957600080fd5b615822836154c4565b915061556d602084016154d6565b6000806040838503121561584357600080fd5b82359150602083013561553581615f45565b6000806040838503121561586857600080fd5b50508035926020909101359150565b60006020828403121561588957600080fd5b613ebb826154d6565b6000602082840312156158a457600080fd5b8151613ebb81615f68565b600081518084526020808501945080840160005b838110156158e85781516001600160a01b0316875295820195908201906001016158c3565b509495945050505050565b8060005b6002811015610e865781518452602093840193909101906001016158f7565b600081518084526020808501945080840160005b838110156158e85781518752958201959082019060010161592a565b6000815180845260005b8181101561596c57602081850181015186830182015201615950565b8181111561597e576000602083870101525b50601f01601f19169290920160200192915050565b61599d81836158f3565b604001919050565b8681526159b560208201876158f3565b6159c260608201866158f3565b6159cf60a08201856158f3565b6159dc60e08201846158f3565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b838152615a1460208201846158f3565b606081019190915260800192915050565b6040810161157e82846158f3565b602081526000613ebb6020830184615916565b602081526000613ebb6020830184615946565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c06080840152615a9f60e08401826158af565b905060808401516001600160601b0380821660a08601528060a08701511660c086015250508091505092915050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b81811015615b1f57845183529383019391830191600101615b03565b509098975050505050505050565b82815260608101613ebb60208301846158f3565b828152604060208201526000615b5a6040830184615916565b949350505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a0830152614d0660c0830184615946565b878152866020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c0830152615beb60e0830184615946565b9998505050505050505050565b87815267ffffffffffffffff87166020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c0830152615beb60e0830184615946565b60006001600160601b03808816835280871660208401525067ffffffffffffffff851660408301526001600160a01b038416606083015260a06080830152615c9760a08301846158af565b979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112615cd757600080fd5b83018035915067ffffffffffffffff821115615cf257600080fd5b602001915036819003821315615d0757600080fd5b9250929050565b60405160c0810167ffffffffffffffff81118282101715615d3157615d31615f2f565b60405290565b604051610120810167ffffffffffffffff81118282101715615d3157615d31615f2f565b60008085851115615d6b57600080fd5b83861115615d7857600080fd5b5050820193919092039150565b60008219821115615d9857615d98615ed7565b500190565b600067ffffffffffffffff808316818516808303821115615dc057615dc0615ed7565b01949350505050565b60006001600160601b03808316818516808303821115615dc057615dc0615ed7565b600082821015615dfd57615dfd615ed7565b500390565b60006001600160601b0383811690831681811015615e2257615e22615ed7565b039392505050565b7fffffffff000000000000000000000000000000000000000000000000000000008135818116916004851015615e6a5780818660040360031b1b83161692505b505092915050565b6000600019821415615e8657615e86615ed7565b5060010190565b600067ffffffffffffffff80831681811415615eab57615eab615ed7565b6001019392505050565b600082615ed257634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114610bf257600080fd5b8015158114610bf257600080fd5b6001600160601b0381168114610bf257600080fdfea164736f6c6343000806000a",
}

var VRFCoordinatorV2PlusABI = VRFCoordinatorV2PlusMetaData.ABI

var VRFCoordinatorV2PlusBin = VRFCoordinatorV2PlusMetaData.Bin

func DeployVRFCoordinatorV2Plus(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address, priceRegistry common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV2Plus, error) {
	parsed, err := VRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2PlusBin), backend, blockhashStore, priceRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV2Plus{VRFCoordinatorV2PlusCaller: VRFCoordinatorV2PlusCaller{contract: contract}, VRFCoordinatorV2PlusTransactor: VRFCoordinatorV2PlusTransactor{contract: contract}, VRFCoordinatorV2PlusFilterer: VRFCoordinatorV2PlusFilterer{contract: contract}}, nil
}

type VRFCoordinatorV2Plus struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV2PlusCaller
	VRFCoordinatorV2PlusTransactor
	VRFCoordinatorV2PlusFilterer
}

type VRFCoordinatorV2PlusCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusSession struct {
	Contract     *VRFCoordinatorV2Plus
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2PlusCallerSession struct {
	Contract *VRFCoordinatorV2PlusCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV2PlusTransactorSession struct {
	Contract     *VRFCoordinatorV2PlusTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2PlusRaw struct {
	Contract *VRFCoordinatorV2Plus
}

type VRFCoordinatorV2PlusCallerRaw struct {
	Contract *VRFCoordinatorV2PlusCaller
}

type VRFCoordinatorV2PlusTransactorRaw struct {
	Contract *VRFCoordinatorV2PlusTransactor
}

func NewVRFCoordinatorV2Plus(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV2Plus, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV2PlusABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV2Plus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Plus{address: address, abi: abi, VRFCoordinatorV2PlusCaller: VRFCoordinatorV2PlusCaller{contract: contract}, VRFCoordinatorV2PlusTransactor: VRFCoordinatorV2PlusTransactor{contract: contract}, VRFCoordinatorV2PlusFilterer: VRFCoordinatorV2PlusFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorV2PlusCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV2PlusCaller, error) {
	contract, err := bindVRFCoordinatorV2Plus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusCaller{contract: contract}, nil
}

func NewVRFCoordinatorV2PlusTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV2PlusTransactor, error) {
	contract, err := bindVRFCoordinatorV2Plus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusTransactor{contract: contract}, nil
}

func NewVRFCoordinatorV2PlusFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV2PlusFilterer, error) {
	contract, err := bindVRFCoordinatorV2Plus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusFilterer{contract: contract}, nil
}

func bindVRFCoordinatorV2Plus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2Plus.Contract.VRFCoordinatorV2PlusCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.VRFCoordinatorV2PlusTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.VRFCoordinatorV2PlusTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2Plus.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.LINK(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.LINK(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2Plus.Contract.MAXCONSUMERS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2Plus.Contract.MAXCONSUMERS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2Plus.Contract.MAXNUMWORDS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2Plus.Contract.MAXNUMWORDS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2Plus.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2Plus.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) PRICEREGISTRY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "PRICE_REGISTRY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) PRICEREGISTRY() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.PRICEREGISTRY(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) PRICEREGISTRY() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.PRICEREGISTRY(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) USDFEEDECIMALS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "USD_FEE_DECIMALS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) USDFEEDECIMALS() (uint8, error) {
	return _VRFCoordinatorV2Plus.Contract.USDFEEDECIMALS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) USDFEEDECIMALS() (uint8, error) {
	return _VRFCoordinatorV2Plus.Contract.USDFEEDECIMALS(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "getActiveSubscriptionIds", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV2Plus.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV2Plus.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.GetRequestConfig(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.GetRequestConfig(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.EthBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV2Plus.Contract.GetSubscription(&_VRFCoordinatorV2Plus.CallOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV2Plus.Contract.GetSubscription(&_VRFCoordinatorV2Plus.CallOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.HashOfKey(&_VRFCoordinatorV2Plus.CallOpts, publicKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.HashOfKey(&_VRFCoordinatorV2Plus.CallOpts, publicKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) MigrationVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "migrationVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) MigrationVersion() (uint8, error) {
	return _VRFCoordinatorV2Plus.Contract.MigrationVersion(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) MigrationVersion() (uint8, error) {
	return _VRFCoordinatorV2Plus.Contract.MigrationVersion(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.Owner(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.Owner(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV2Plus.Contract.PendingRequestExists(&_VRFCoordinatorV2Plus.CallOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV2Plus.Contract.PendingRequestExists(&_VRFCoordinatorV2Plus.CallOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinimumRequestConfirmations = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.MaxGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ReentrancyLock = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV2Plus.Contract.SConfig(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV2Plus.Contract.SConfig(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SCurrentSubNonce(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_currentSubNonce")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV2Plus.Contract.SCurrentSubNonce(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV2Plus.Contract.SCurrentSubNonce(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_provingKeyHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.SProvingKeyHashes(&_VRFCoordinatorV2Plus.CallOpts, arg0)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.SProvingKeyHashes(&_VRFCoordinatorV2Plus.CallOpts, arg0)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_provingKeys", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SProvingKeys(arg0 [32]byte) (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.SProvingKeys(&_VRFCoordinatorV2Plus.CallOpts, arg0)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SProvingKeys(arg0 [32]byte) (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.SProvingKeys(&_VRFCoordinatorV2Plus.CallOpts, arg0)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_requestCommitments", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.SRequestCommitments(&_VRFCoordinatorV2Plus.CallOpts, arg0)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2Plus.Contract.SRequestCommitments(&_VRFCoordinatorV2Plus.CallOpts, arg0)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) STotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.STotalBalance(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.STotalBalance(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) STotalEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_totalEthBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) STotalEthBalance() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.STotalEthBalance(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) STotalEthBalance() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.STotalEthBalance(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.AcceptOwnership(&_VRFCoordinatorV2Plus.TransactOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.AcceptOwnership(&_VRFCoordinatorV2Plus.TransactOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2Plus.TransactOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2Plus.TransactOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.AddConsumer(&_VRFCoordinatorV2Plus.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.AddConsumer(&_VRFCoordinatorV2Plus.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.CancelSubscription(&_VRFCoordinatorV2Plus.TransactOpts, subId, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.CancelSubscription(&_VRFCoordinatorV2Plus.TransactOpts, subId, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.CreateSubscription(&_VRFCoordinatorV2Plus.TransactOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.CreateSubscription(&_VRFCoordinatorV2Plus.TransactOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "deregisterMigratableCoordinator", target)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorV2Plus.TransactOpts, target)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorV2Plus.TransactOpts, target)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.DeregisterProvingKey(&_VRFCoordinatorV2Plus.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.DeregisterProvingKey(&_VRFCoordinatorV2Plus.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2PlusRequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2PlusRequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.FulfillRandomWords(&_VRFCoordinatorV2Plus.TransactOpts, proof, rc)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2PlusRequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.FulfillRandomWords(&_VRFCoordinatorV2Plus.TransactOpts, proof, rc)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) FundSubscriptionWithEth(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "fundSubscriptionWithEth", subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) FundSubscriptionWithEth(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.FundSubscriptionWithEth(&_VRFCoordinatorV2Plus.TransactOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) FundSubscriptionWithEth(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.FundSubscriptionWithEth(&_VRFCoordinatorV2Plus.TransactOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "migrate", subId, newCoordinator)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.Migrate(&_VRFCoordinatorV2Plus.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.Migrate(&_VRFCoordinatorV2Plus.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OnTokenTransfer(&_VRFCoordinatorV2Plus.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OnTokenTransfer(&_VRFCoordinatorV2Plus.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OracleWithdraw(&_VRFCoordinatorV2Plus.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OracleWithdraw(&_VRFCoordinatorV2Plus.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) OracleWithdrawEth(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "oracleWithdrawEth", recipient, amount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) OracleWithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OracleWithdrawEth(&_VRFCoordinatorV2Plus.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) OracleWithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OracleWithdrawEth(&_VRFCoordinatorV2Plus.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2Plus.TransactOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2Plus.TransactOpts, subId)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) RecoverEthFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "recoverEthFunds", to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) RecoverEthFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RecoverEthFunds(&_VRFCoordinatorV2Plus.TransactOpts, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) RecoverEthFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RecoverEthFunds(&_VRFCoordinatorV2Plus.TransactOpts, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RecoverFunds(&_VRFCoordinatorV2Plus.TransactOpts, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RecoverFunds(&_VRFCoordinatorV2Plus.TransactOpts, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "registerMigratableCoordinator", target)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV2Plus.TransactOpts, target)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV2Plus.TransactOpts, target)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RegisterProvingKey(&_VRFCoordinatorV2Plus.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RegisterProvingKey(&_VRFCoordinatorV2Plus.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RemoveConsumer(&_VRFCoordinatorV2Plus.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RemoveConsumer(&_VRFCoordinatorV2Plus.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "requestRandomWords", req)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RequestRandomWords(&_VRFCoordinatorV2Plus.TransactOpts, req)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RequestRandomWords(&_VRFCoordinatorV2Plus.TransactOpts, req)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2Plus.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2Plus.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetConfig(&_VRFCoordinatorV2Plus.TransactOpts, minimumRequestConfirmations, maxGasLimit)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetConfig(&_VRFCoordinatorV2Plus.TransactOpts, minimumRequestConfirmations, maxGasLimit)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "setLINK", link)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SetLINK(link common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetLINK(&_VRFCoordinatorV2Plus.TransactOpts, link)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) SetLINK(link common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetLINK(&_VRFCoordinatorV2Plus.TransactOpts, link)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.TransferOwnership(&_VRFCoordinatorV2Plus.TransactOpts, to)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.TransferOwnership(&_VRFCoordinatorV2Plus.TransactOpts, to)
}

type VRFCoordinatorV2PlusConfigSetIterator struct {
	Event *VRFCoordinatorV2PlusConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	Raw                         types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusConfigSetIterator{contract: _VRFCoordinatorV2Plus.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusConfigSet)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV2PlusConfigSet, error) {
	event := new(VRFCoordinatorV2PlusConfigSet)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusCoordinatorDeregisteredIterator struct {
	Event *VRFCoordinatorV2PlusCoordinatorDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusCoordinatorDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusCoordinatorDeregistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusCoordinatorDeregistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusCoordinatorDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusCoordinatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusCoordinatorDeregistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusCoordinatorDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusCoordinatorDeregisteredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "CoordinatorDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusCoordinatorDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusCoordinatorDeregistered)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorV2PlusCoordinatorDeregistered, error) {
	event := new(VRFCoordinatorV2PlusCoordinatorDeregistered)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusCoordinatorRegisteredIterator struct {
	Event *VRFCoordinatorV2PlusCoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusCoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusCoordinatorRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusCoordinatorRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusCoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusCoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusCoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusCoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusCoordinatorRegisteredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusCoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusCoordinatorRegistered)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV2PlusCoordinatorRegistered, error) {
	event := new(VRFCoordinatorV2PlusCoordinatorRegistered)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusEthFundsRecoveredIterator struct {
	Event *VRFCoordinatorV2PlusEthFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusEthFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusEthFundsRecovered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusEthFundsRecovered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusEthFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusEthFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusEthFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterEthFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusEthFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "EthFundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusEthFundsRecoveredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "EthFundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchEthFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusEthFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "EthFundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusEthFundsRecovered)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "EthFundsRecovered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseEthFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusEthFundsRecovered, error) {
	event := new(VRFCoordinatorV2PlusEthFundsRecovered)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "EthFundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusFundsRecoveredIterator struct {
	Event *VRFCoordinatorV2PlusFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusFundsRecovered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusFundsRecovered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusFundsRecoveredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusFundsRecovered)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusFundsRecovered, error) {
	event := new(VRFCoordinatorV2PlusFundsRecovered)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusMigrationCompletedIterator struct {
	Event *VRFCoordinatorV2PlusMigrationCompleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusMigrationCompletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusMigrationCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusMigrationCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusMigrationCompletedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusMigrationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusMigrationCompleted struct {
	NewCoordinator common.Address
	SubId          *big.Int
	Raw            types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusMigrationCompletedIterator, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusMigrationCompletedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "MigrationCompleted", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusMigrationCompleted) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusMigrationCompleted)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV2PlusMigrationCompleted, error) {
	event := new(VRFCoordinatorV2PlusMigrationCompleted)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusOwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV2PlusOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusOwnershipTransferRequestedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusOwnershipTransferRequested)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2PlusOwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV2PlusOwnershipTransferRequested)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusOwnershipTransferredIterator struct {
	Event *VRFCoordinatorV2PlusOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusOwnershipTransferredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusOwnershipTransferred)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2PlusOwnershipTransferred, error) {
	event := new(VRFCoordinatorV2PlusOwnershipTransferred)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorV2PlusProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusProvingKeyDeregistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusProvingKeyDeregistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2PlusProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusProvingKeyDeregisteredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusProvingKeyDeregistered)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV2PlusProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorV2PlusProvingKeyDeregistered)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV2PlusProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusProvingKeyRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusProvingKeyRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2PlusProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusProvingKeyRegisteredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusProvingKeyRegistered)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2PlusProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV2PlusProvingKeyRegistered)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusRandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV2PlusRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusRandomWordsFulfilled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusRandomWordsFulfilled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusRandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	SubID      *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subID []*big.Int) (*VRFCoordinatorV2PlusRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule, subIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusRandomWordsFulfilledIterator{contract: _VRFCoordinatorV2Plus.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusRandomWordsFulfilled, requestId []*big.Int, subID []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule, subIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusRandomWordsFulfilled)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2PlusRandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV2PlusRandomWordsFulfilled)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusRandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV2PlusRandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusRandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusRandomWordsRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusRandomWordsRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusRandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusRandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       *big.Int
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	ExtraArgs                   []byte
	Sender                      common.Address
	Raw                         types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV2PlusRandomWordsRequestedIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusRandomWordsRequestedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusRandomWordsRequested)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2PlusRandomWordsRequested, error) {
	event := new(VRFCoordinatorV2PlusRandomWordsRequested)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionCanceled struct {
	SubId      *big.Int
	To         common.Address
	AmountLink *big.Int
	AmountEth  *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionCanceledIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionCanceled, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionCanceled)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2PlusSubscriptionCanceled, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionCanceled)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionConsumerAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionConsumerAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionConsumerAdded struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionConsumerAdded)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2PlusSubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionConsumerAdded)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionConsumerRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionConsumerRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionConsumerRemoved struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2PlusSubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionCreated struct {
	SubId *big.Int
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionCreatedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionCreated, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionCreated)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2PlusSubscriptionCreated, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionCreated)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionFundedIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionFunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionFunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionFunded struct {
	SubId      *big.Int
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionFundedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionFunded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionFunded)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2PlusSubscriptionFunded, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionFunded)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionFundedWithEthIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionFundedWithEth

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionFundedWithEthIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionFundedWithEth)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionFundedWithEth)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionFundedWithEthIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionFundedWithEthIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionFundedWithEth struct {
	SubId         *big.Int
	OldEthBalance *big.Int
	NewEthBalance *big.Int
	Raw           types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionFundedWithEth(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionFundedWithEthIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionFundedWithEth", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionFundedWithEthIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionFundedWithEth", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionFundedWithEth(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionFundedWithEth, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionFundedWithEth", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionFundedWithEth)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionFundedWithEth", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionFundedWithEth(log types.Log) (*VRFCoordinatorV2PlusSubscriptionFundedWithEth, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionFundedWithEth)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionFundedWithEth", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusSubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV2PlusSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusSubscriptionOwnerTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
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
		it.Event = new(VRFCoordinatorV2PlusSubscriptionOwnerTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorV2PlusSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusSubscriptionOwnerTransferred struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusSubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV2Plus.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2Plus.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusSubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV2PlusSubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV2Plus.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetSubscription struct {
	Balance    *big.Int
	EthBalance *big.Int
	ReqCount   uint64
	Owner      common.Address
	Consumers  []common.Address
}
type SConfig struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	ReentrancyLock              bool
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2Plus) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV2Plus.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV2Plus.ParseConfigSet(log)
	case _VRFCoordinatorV2Plus.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFCoordinatorV2Plus.ParseCoordinatorDeregistered(log)
	case _VRFCoordinatorV2Plus.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinatorV2Plus.ParseCoordinatorRegistered(log)
	case _VRFCoordinatorV2Plus.abi.Events["EthFundsRecovered"].ID:
		return _VRFCoordinatorV2Plus.ParseEthFundsRecovered(log)
	case _VRFCoordinatorV2Plus.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorV2Plus.ParseFundsRecovered(log)
	case _VRFCoordinatorV2Plus.abi.Events["MigrationCompleted"].ID:
		return _VRFCoordinatorV2Plus.ParseMigrationCompleted(log)
	case _VRFCoordinatorV2Plus.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorV2Plus.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorV2Plus.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorV2Plus.ParseOwnershipTransferred(log)
	case _VRFCoordinatorV2Plus.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorV2Plus.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorV2Plus.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorV2Plus.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorV2Plus.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorV2Plus.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorV2Plus.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorV2Plus.ParseRandomWordsRequested(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionCreated(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionFunded(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionFundedWithEth"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionFundedWithEth(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorV2Plus.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorV2Plus.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorV2PlusConfigSet) Topic() common.Hash {
	return common.HexToHash("0x52b99d89f788ddd6c815037dcb52b920032f7ad04d1ed6655c5658021a2b1341")
}

func (VRFCoordinatorV2PlusCoordinatorDeregistered) Topic() common.Hash {
	return common.HexToHash("0xf80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37")
}

func (VRFCoordinatorV2PlusCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFCoordinatorV2PlusEthFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df377317")
}

func (VRFCoordinatorV2PlusFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorV2PlusMigrationCompleted) Topic() common.Hash {
	return common.HexToHash("0xd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187")
}

func (VRFCoordinatorV2PlusOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorV2PlusOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorV2PlusProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d")
}

func (VRFCoordinatorV2PlusProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8")
}

func (VRFCoordinatorV2PlusRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x49580fdfd9497e1ed5c1b1cec0495087ae8e3f1267470ec2fb015db32e3d6aa7")
}

func (VRFCoordinatorV2PlusRandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0xeb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e")
}

func (VRFCoordinatorV2PlusSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0x8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4")
}

func (VRFCoordinatorV2PlusSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e1")
}

func (VRFCoordinatorV2PlusSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7")
}

func (VRFCoordinatorV2PlusSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d")
}

func (VRFCoordinatorV2PlusSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0x1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a")
}

func (VRFCoordinatorV2PlusSubscriptionFundedWithEth) Topic() common.Hash {
	return common.HexToHash("0x3f1ddc3ab1bdb39001ad76ca51a0e6f57ce6627c69f251d1de41622847721cde")
}

func (VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1")
}

func (VRFCoordinatorV2PlusSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0xd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386")
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2Plus) Address() common.Address {
	return _VRFCoordinatorV2Plus.address
}

type VRFCoordinatorV2PlusInterface interface {
	BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

	PRICEREGISTRY(opts *bind.CallOpts) (common.Address, error)

	USDFEEDECIMALS(opts *bind.CallOpts) (uint8, error)

	GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error)

	MigrationVersion(opts *bind.CallOpts) (uint8, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	SCurrentSubNonce(opts *bind.CallOpts) (uint64, error)

	SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error)

	SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	STotalBalance(opts *bind.CallOpts) (*big.Int, error)

	STotalEthBalance(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2PlusRequestCommitment) (*types.Transaction, error)

	FundSubscriptionWithEth(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OracleWithdrawEth(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	RecoverEthFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32) (*types.Transaction, error)

	SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorV2PlusConfigSet, error)

	FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusCoordinatorDeregisteredIterator, error)

	WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusCoordinatorDeregistered) (event.Subscription, error)

	ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorV2PlusCoordinatorDeregistered, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusCoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusCoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV2PlusCoordinatorRegistered, error)

	FilterEthFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusEthFundsRecoveredIterator, error)

	WatchEthFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusEthFundsRecovered) (event.Subscription, error)

	ParseEthFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusEthFundsRecovered, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusFundsRecovered, error)

	FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusMigrationCompletedIterator, error)

	WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusMigrationCompleted) (event.Subscription, error)

	ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV2PlusMigrationCompleted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2PlusOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2PlusOwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2PlusProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV2PlusProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2PlusProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusProvingKeyRegistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2PlusProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subID []*big.Int) (*VRFCoordinatorV2PlusRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusRandomWordsFulfilled, requestId []*big.Int, subID []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2PlusRandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV2PlusRandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2PlusRandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionCanceled, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2PlusSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2PlusSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2PlusSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionCreated, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2PlusSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionFunded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2PlusSubscriptionFunded, error)

	FilterSubscriptionFundedWithEth(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionFundedWithEthIterator, error)

	WatchSubscriptionFundedWithEth(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionFundedWithEth, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFundedWithEth(log types.Log) (*VRFCoordinatorV2PlusSubscriptionFundedWithEth, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2PlusSubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
