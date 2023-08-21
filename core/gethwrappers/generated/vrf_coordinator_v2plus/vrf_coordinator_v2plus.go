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

type VRFCoordinatorV2PlusFeeConfig struct {
	FulfillmentFlatFeeLinkUSD *big.Int
	FulfillmentFlatFeeEthUSD  *big.Int
}

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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendEther\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"ethUSD\",\"type\":\"int256\"}],\"name\":\"InvalidEthUSDPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkUSD\",\"type\":\"int256\"}],\"name\":\"InvalidLinkUSDPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feed\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"}],\"name\":\"InvalidUSDPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitEth\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorV2Plus.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EthFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountEth\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldEthBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newEthBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithEth\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ETH_USD_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_USD_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"USD_FEE_DECIMALS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFCoordinatorV2Plus.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"ethBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdrawEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverEthFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackUSDPerUnitEth\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackUSDPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeConfig\",\"outputs\":[{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalEthBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitEth\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"fallbackUSDPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeLinkUSD\",\"type\":\"uint40\"},{\"internalType\":\"uint40\",\"name\":\"fulfillmentFlatFeeEthUSD\",\"type\":\"uint40\"}],\"internalType\":\"structVRFCoordinatorV2Plus.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"ethUSDFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620063bd380380620063bd833981016040819052620000349162000183565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d7565b50505060601b6001600160601b031916608052620001b5565b6001600160a01b038116331415620001325760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019657600080fd5b81516001600160a01b0381168114620001ae57600080fd5b9392505050565b60805160601c6161e2620001db600039600081816105a201526137d501526161e26000f3fe6080604052600436106102435760003560e01c80620122911461024857806304104edb14610275578063043bd6ae14610297578063088070f5146102bb57806308821d581461033b5780630ae095401461035b57806315c48b841461037b5780631b6b6d23146103a3578063294daa49146103d0578063330987b3146103f25780633dc2c7eb1461042a578063405b84fa1461044a57806340d6bb821461046a57806341af6c871461049557806346d8d486146104c55780634b326d93146104e557806359392b6d14610505578063595cbbfa1461051b5780635d06b4ab1461053b57806364d51a2a1461055b57806366316d8d14610570578063689c4517146105905780636b6feccc146105c45780636f64f03f1461060b57806379ba50971461062b578063835c0dfc1461064057806386fe91c7146106565780638da5cb5b1461067657806394d8572a146106945780639b1c385e146106b45780639d40a6fd146106d4578063a21a23e414610701578063a4c0ed3614610716578063a8cb447b14610736578063aa433aff14610756578063ad17836114610776578063aefb212f14610796578063b2a7cac5146107c3578063bb0697a5146107e3578063bec4c08c146107f8578063caf70c4a14610818578063cb63179714610838578063d98e620e14610858578063da2f261014610878578063dac83d29146108ae578063dc311dd3146108ce578063e72f6e30146108ff578063e8509bff1461091f578063e95704bd14610932578063ee9d2d3814610959578063f2fde38b14610986575b600080fd5b34801561025457600080fd5b5061025d6109a6565b60405161026c93929190615bf2565b60405180910390f35b34801561028157600080fd5b5061029561029036600461543b565b610a22565b005b3480156102a357600080fd5b506102ad60135481565b60405190815260200161026c565b3480156102c757600080fd5b50600f546103039061ffff81169063ffffffff62010000820481169160ff600160301b82041691600160381b8204811691600160581b90041685565b6040805161ffff909616865263ffffffff9485166020870152921515928501929092528216606084015216608082015260a00161026c565b34801561034757600080fd5b5061029561035636600461558d565b610bd8565b34801561036757600080fd5b50610295610376366004615845565b610d6c565b34801561038757600080fd5b5061039060c881565b60405161ffff909116815260200161026c565b3480156103af57600080fd5b506002546103c3906001600160a01b031681565b60405161026c9190615a96565b3480156103dc57600080fd5b5060015b60405160ff909116815260200161026c565b3480156103fe57600080fd5b5061041261040d36600461565f565b610e31565b6040516001600160601b03909116815260200161026c565b34801561043657600080fd5b506004546103c3906001600160a01b031681565b34801561045657600080fd5b50610295610465366004615845565b61130c565b34801561047657600080fd5b506104806101f481565b60405163ffffffff909116815260200161026c565b3480156104a157600080fd5b506104b56104b03660046155e2565b611682565b604051901515815260200161026c565b3480156104d157600080fd5b506102956104e0366004615458565b611823565b3480156104f157600080fd5b50610295610500366004615791565b61199f565b34801561051157600080fd5b506102ad60155481565b34801561052757600080fd5b5061029561053636600461548d565b611b6d565b34801561054757600080fd5b5061029561055636600461543b565b611bde565b34801561056757600080fd5b50610390606481565b34801561057c57600080fd5b5061029561058b366004615458565b611c95565b34801561059c57600080fd5b506103c37f000000000000000000000000000000000000000000000000000000000000000081565b3480156105d057600080fd5b506016546105ed9064ffffffffff80821691600160281b90041682565b6040805164ffffffffff93841681529290911660208301520161026c565b34801561061757600080fd5b506102956106263660046154d8565b611e34565b34801561063757600080fd5b50610295611f33565b34801561064c57600080fd5b506102ad60145481565b34801561066257600080fd5b50600c54610412906001600160601b031681565b34801561068257600080fd5b506000546001600160a01b03166103c3565b3480156106a057600080fd5b506005546103c3906001600160a01b031681565b3480156106c057600080fd5b506102ad6106cf36600461573c565b611fdd565b3480156106e057600080fd5b506009546106f4906001600160401b031681565b60405161026c9190615d8b565b34801561070d57600080fd5b506102ad61234d565b34801561072257600080fd5b50610295610731366004615505565b61259b565b34801561074257600080fd5b5061029561075136600461543b565b612738565b34801561076257600080fd5b506102956107713660046155e2565b612844565b34801561078257600080fd5b506003546103c3906001600160a01b031681565b3480156107a257600080fd5b506107b66107b1366004615875565b6128a4565b60405161026c9190615b0d565b3480156107cf57600080fd5b506102956107de3660046155e2565b6129a5565b3480156107ef57600080fd5b506103e0600881565b34801561080457600080fd5b50610295610813366004615845565b612ad5565b34801561082457600080fd5b506102ad6108333660046155a9565b612c6c565b34801561084457600080fd5b50610295610853366004615845565b612c9c565b34801561086457600080fd5b506102ad6108733660046155e2565b612f89565b34801561088457600080fd5b506103c36108933660046155e2565b6010602052600090815260409020546001600160a01b031681565b3480156108ba57600080fd5b506102956108c9366004615845565b612faa565b3480156108da57600080fd5b506108ee6108e93660046155e2565b6130ba565b60405161026c959493929190615d9f565b34801561090b57600080fd5b5061029561091a36600461543b565b6131b5565b61029561092d3660046155e2565b61334b565b34801561093e57600080fd5b50600c5461041290600160601b90046001600160601b031681565b34801561096557600080fd5b506102ad6109743660046155e2565b60126020526000908152604090205481565b34801561099257600080fd5b506102956109a136600461543b565b613482565b600f546011805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff16939192839190830182828015610a1057602002820191906000526020600020905b8154815260200190600101908083116109fc575b50505050509050925092509250909192565b610a2a613493565b60175460005b81811015610bb057826001600160a01b031660178281548110610a5557610a55616186565b6000918252602090912001546001600160a01b03161415610b9e576017610a7d60018461605e565b81548110610a8d57610a8d616186565b600091825260209091200154601780546001600160a01b039092169183908110610ab957610ab9616186565b600091825260209091200180546001600160a01b0319166001600160a01b0392909216919091179055826017610af060018561605e565b81548110610b0057610b00616186565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506017805480610b3f57610b3f616170565b600082815260209020810160001990810180546001600160a01b03191690550190556040517ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3790610b91908590615a96565b60405180910390a1505050565b80610ba8816160ee565b915050610a30565b5081604051635428d44960e01b8152600401610bcc9190615a96565b60405180910390fd5b50565b610be0613493565b604080518082018252600091610c0f919084906002908390839080828437600092019190915250612c6c915050565b6000818152601060205260409020549091506001600160a01b031680610c4b57604051631dfd6e1360e21b815260048101839052602401610bcc565b600082815260106020526040812080546001600160a01b03191690555b601154811015610d23578260118281548110610c8657610c86616186565b90600052602060002001541415610d11576011805460009190610cab9060019061605e565b81548110610cbb57610cbb616186565b906000526020600020015490508060118381548110610cdc57610cdc616186565b6000918252602090912001556011805480610cf957610cf9616170565b60019003818190600052602060002001600090559055505b80610d1b816160ee565b915050610c68565b50806001600160a01b03167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610d5f91815260200190565b60405180910390a2505050565b60008281526007602052604090205482906001600160a01b031680610da457604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614610dcf5780604051636c51fda960e11b8152600401610bcc9190615a96565b600f54600160301b900460ff1615610dfa5760405163769dd35360e11b815260040160405180910390fd5b610e0384611682565b15610e2157604051631685ecdd60e31b815260040160405180910390fd5b610e2b84846134e8565b50505050565b600f54600090600160301b900460ff1615610e5f5760405163769dd35360e11b815260040160405180910390fd5b60005a90506000610e70858561367b565b90506000846060015163ffffffff166001600160401b03811115610e9657610e9661619c565b604051908082528060200260200182016040528015610ebf578160200160208202803683370190505b50905060005b856060015163ffffffff16811015610f3657826040015181604051602001610eee929190615b20565b6040516020818303038152906040528051906020012060001c828281518110610f1957610f19616186565b602090810291909101015280610f2e816160ee565b915050610ec5565b5060208083018051600090815260129092526040808320839055905190518291631fe543e360e01b91610f6e91908690602401615c7c565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600f805460ff60301b1916600160301b179055908801516080890151919250600091610fd39163ffffffff1690846138ee565b600f805460ff60301b19169055602089810151600090815260089091526040902054909150600160c01b90046001600160401b0316611013816001615ef0565b6020808b0151600090815260089091526040812080546001600160401b0393909316600160c01b026001600160c01b039093169290921790915560a08a015180516110609060019061605e565b8151811061107057611070616186565b602091010151600f5460f89190911c60011491506000906110a1908a90600160581b900463ffffffff163a8561393c565b905081156111aa576020808c01516000908152600890915260409020546001600160601b03808316600160601b9092041610156110f157604051631e9acf1760e31b815260040160405180910390fd5b60208b81015160009081526008909152604090208054829190600c90611128908490600160601b90046001600160601b0316616098565b82546101009290920a6001600160601b0381810219909316918316021790915589516000908152601060209081526040808320546001600160a01b03168352600e90915281208054859450909261118191859116615f1b565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550611296565b6020808c01516000908152600890915260409020546001600160601b03808316911610156111eb57604051631e9acf1760e31b815260040160405180910390fd5b6020808c0151600090815260089091526040812080548392906112189084906001600160601b0316616098565b82546101009290920a6001600160601b0381810219909316918316021790915589516000908152601060209081526040808320546001600160a01b03168352600d90915281208054859450909261127191859116615f1b565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b8a6020015188602001517f49580fdfd9497e1ed5c1b1cec0495087ae8e3f1267470ec2fb015db32e3d6aa78a6040015184886040516112f3939291909283526001600160601b039190911660208301521515604082015260600190565b60405180910390a3985050505050505050505b92915050565b6113158161398d565b6113345780604051635428d44960e01b8152600401610bcc9190615a96565b600080600080611343866130ba565b945094505093509350336001600160a01b0316826001600160a01b0316146113a65760405162461bcd60e51b81526020600482015260166024820152752737ba1039bab139b1b934b83a34b7b71037bbb732b960511b6044820152606401610bcc565b6113af86611682565b156113f55760405162461bcd60e51b815260206004820152601660248201527550656e64696e6720726571756573742065786973747360501b6044820152606401610bcc565b60006040518060c0016040528061140a600190565b60ff168152602001888152602001846001600160a01b03168152602001838152602001866001600160601b03168152602001856001600160601b0316815250905060008160405160200161145e9190615b5f565b6040516020818303038152906040529050611478886139f7565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b038816906114b1908590600401615b4c565b6000604051808303818588803b1580156114ca57600080fd5b505af11580156114de573d6000803e3d6000fd5b505060025460405163a9059cbb60e01b81526001600160a01b03909116935063a9059cbb925061151591508a908a90600401615add565b602060405180830381600087803b15801561152f57600080fd5b505af1158015611543573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061156791906155c5565b6115a85760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b6044820152606401610bcc565b60005b835181101561163e578381815181106115c6576115c6616186565b60200260200101516001600160a01b0316638ea98117896040518263ffffffff1660e01b81526004016115f99190615a96565b600060405180830381600087803b15801561161357600080fd5b505af1158015611627573d6000803e3d6000fd5b505050508080611636906160ee565b9150506115ab565b507fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be41878789604051611670929190615aaa565b60405180910390a15050505050505050565b6000818152600760209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561170c57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116116ee575b505050505081525050905060005b8160400151518110156118195760005b6011548110156118065760006117cf6011838154811061174c5761174c616186565b90600052602060002001548560400151858151811061176d5761176d616186565b602002602001015188600660008960400151898151811061179057611790616186565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208d82529092529020546001600160401b0316613c44565b50600081815260126020526040902054909150156117f35750600195945050505050565b50806117fe816160ee565b91505061172a565b5080611811816160ee565b91505061171a565b5060009392505050565b600f54600160301b900460ff161561184e5760405163769dd35360e11b815260040160405180910390fd5b336000908152600e60205260409020546001600160601b038083169116101561188a57604051631e9acf1760e31b815260040160405180910390fd5b336000908152600e6020526040812080548392906118b29084906001600160601b0316616098565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600c808282829054906101000a90046001600160601b03166118f99190616098565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114611973576040519150601f19603f3d011682016040523d82523d6000602084013e611978565b606091505b505090508061199a57604051630dcf35db60e41b815260040160405180910390fd5b505050565b6119a7613493565b60c861ffff891611156119d457878860c860405163539c34bb60e11b8152600401610bcc93929190615bd4565b600084136119f8576040516321ea67b360e11b815260048101859052602401610bcc565b60008313611a1c5760405163599d67e360e01b815260048101849052602401610bcc565b60008213611a40576040516325b2499f60e01b815260048101839052602401610bcc565b6040805160a0808201835261ffff8b1680835263ffffffff8b811660208086018290526000868801528c83166060808801829052938d166080978801819052600f805465ffffffffffff1916871762010000860217600160301b600160781b031916600160381b840263ffffffff60581b191617600160581b83021790558951601680548c86015164ffffffffff9384166001600160501b031990921691909117600160281b91841691909102179081905560138f905560148e905560158d90558a5197885293870194909452978501529183019590955292810188905290810186905260c0810185905281831660e082015260289190911c9091166101008201527f306a7924a41e68523d99ac4f68693ab16b5dba36ff53896a0977971a564cc1029061012001611670565b611b75613493565b6002546001600160a01b031615611b9f57604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039485166001600160a01b031991821617909155600380549385169382169390931790925560058054919093169116179055565b611be6613493565b611bef8161398d565b15611c0f578060405163ac8a27ef60e01b8152600401610bcc9190615a96565b601780546001810182556000919091527fc624b66cc0138b8fabc209247f72d758e1cf3343756d543badbf24212bed8c150180546001600160a01b0319166001600160a01b0383161790556040517fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af0162590611c8a908390615a96565b60405180910390a150565b600f54600160301b900460ff1615611cc05760405163769dd35360e11b815260040160405180910390fd5b336000908152600d60205260409020546001600160601b0380831691161015611cfc57604051631e9acf1760e31b815260040160405180910390fd5b336000908152600d602052604081208054839290611d249084906001600160601b0316616098565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600c60008282829054906101000a90046001600160601b0316611d6c9190616098565b82546001600160601b039182166101009390930a92830291909202199091161790555060025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb90611dc19085908590600401615add565b602060405180830381600087803b158015611ddb57600080fd5b505af1158015611def573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e1391906155c5565b611e3057604051631e9acf1760e31b815260040160405180910390fd5b5050565b611e3c613493565b604080518082018252600091611e6b919084906002908390839080828437600092019190915250612c6c915050565b6000818152601060205260409020549091506001600160a01b031615611ea757604051634a0b8fa760e01b815260048101829052602401610bcc565b600081815260106020908152604080832080546001600160a01b0319166001600160a01b0388169081179091556011805460018101825594527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c68909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610d5f565b6001546001600160a01b03163314611f865760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610bcc565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600f54600090600160301b900460ff161561200b5760405163769dd35360e11b815260040160405180910390fd5b6020808301356000908152600790915260409020546001600160a01b031661204657604051630fb532db60e11b815260040160405180910390fd5b3360009081526006602090815260408083208583013584529091529020546001600160401b031680612093578260200135336040516379bfd40160e01b8152600401610bcc929190615c51565b600f5461ffff166120aa6060850160408601615776565b61ffff1610806120cd575060c86120c76060850160408601615776565b61ffff16115b15612107576120e26060840160408501615776565b600f5460405163539c34bb60e11b8152610bcc929161ffff169060c890600401615bd4565b600f5462010000900463ffffffff166121266080850160608601615897565b63ffffffff16111561216c576121426080840160608501615897565b600f54604051637aebf00f60e11b8152610bcc929162010000900463ffffffff1690600401615d74565b6101f461217f60a0850160808601615897565b63ffffffff1611156121b95761219b60a0840160808501615897565b6101f46040516311ce1afb60e21b8152600401610bcc929190615d74565b60006121c6826001615ef0565b90506000806121dc863533602089013586613c44565b909250905060006121f86121f360a0890189615df4565b613ccd565b9050600061220582613d4a565b905083612210613dbb565b60208a013561222560808c0160608d01615897565b61223560a08d0160808e01615897565b338660405160200161224d9796959493929190615cd4565b604051602081830303815290604052805190602001206012600086815260200190815260200160002081905550336001600160a01b0316886020013589600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e87878d60400160208101906122c49190615776565b8e60600160208101906122d79190615897565b8f60800160208101906122ea9190615897565b896040516122fd96959493929190615c95565b60405180910390a45050336000908152600660209081526040808320898301358452909152902080546001600160401b0319166001600160401b039490941693909317909255925050505b919050565b600f54600090600160301b900460ff161561237b5760405163769dd35360e11b815260040160405180910390fd5b60003361238960014361605e565b600954604051606093841b6001600160601b03199081166020830152924060348201523090931b909116605483015260c01b6001600160c01b031916606882015260700160408051601f198184030181529190528051602090910120600980549192506001600160401b0390911690600061240383616109565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550506000806001600160401b038111156124425761244261619c565b60405190808252806020026020018201604052801561246b578160200160208202803683370190505b506040805160608082018352600080835260208084018281528486018381528984526008835286842095518654925191516001600160601b039182166001600160c01b031990941693909317600160601b9190921602176001600160c01b0316600160c01b6001600160401b039092169190910217909355835191820184523382528183018181528285018681528883526007855294909120825181546001600160a01b03199081166001600160a01b03928316178355925160018301805490941691161790915592518051949550909361254c926002850192019061513f565b5061255c9150600a905083613e54565b50817f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d3360405161258d9190615a96565b60405180910390a250905090565b600f54600160301b900460ff16156125c65760405163769dd35360e11b815260040160405180910390fd5b6002546001600160a01b031633146125f1576040516344b0e3c360e01b815260040160405180910390fd5b6020811461261257604051638129bbcd60e01b815260040160405180910390fd5b6000612620828401846155e2565b6000818152600760205260409020549091506001600160a01b031661265857604051630fb532db60e11b815260040160405180910390fd5b600081815260086020526040812080546001600160601b03169186919061267f8385615f1b565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600c60008282829054906101000a90046001600160601b03166126c79190615f1b565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a82878461271a9190615ed8565b604051612728929190615b20565b60405180910390a2505050505050565b612740613493565b600c544790600160601b90046001600160601b03168181111561277a5780826040516354ced18160e11b8152600401610bcc929190615b20565b8181101561199a57600061278e828461605e565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d80600081146127dd576040519150601f19603f3d011682016040523d82523d6000602084013e6127e2565b606091505b505090508061280457604051630dcf35db60e41b815260040160405180910390fd5b7f879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df3773178583604051612835929190615aaa565b60405180910390a15050505050565b61284c613493565b6000818152600760205260409020546001600160a01b031661288157604051630fb532db60e11b815260040160405180910390fd5b600081815260076020526040902054610bd59082906001600160a01b03166134e8565b606060006128b2600a613e60565b90508084106128d457604051631390f2a160e01b815260040160405180910390fd5b60006128e08486615ed8565b9050818111806128ee575083155b6128f857806128fa565b815b90506000612908868361605e565b6001600160401b0381111561291f5761291f61619c565b604051908082528060200260200182016040528015612948578160200160208202803683370190505b50905060005b815181101561299b5761296c6129648883615ed8565b600a90613e6a565b82828151811061297e5761297e616186565b602090810291909101015280612993816160ee565b91505061294e565b5095945050505050565b600f54600160301b900460ff16156129d05760405163769dd35360e11b815260040160405180910390fd5b6000818152600760205260409020546001600160a01b0316612a0557604051630fb532db60e11b815260040160405180910390fd5b6000818152600760205260409020600101546001600160a01b03163314612a5c576000818152600760205260409081902060010154905163d084e97560e01b8152610bcc916001600160a01b031690600401615a96565b6000818152600760205260409081902080546001600160a01b031980821633908117845560019093018054909116905591516001600160a01b039092169183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691612ac9918591615ac3565b60405180910390a25050565b60008281526007602052604090205482906001600160a01b031680612b0d57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612b385780604051636c51fda960e11b8152600401610bcc9190615a96565b600f54600160301b900460ff1615612b635760405163769dd35360e11b815260040160405180910390fd5b60008481526007602052604090206002015460641415612b96576040516305a48e0f60e01b815260040160405180910390fd5b6001600160a01b03831660009081526006602090815260408083208784529091529020546001600160401b031615612bcd57610e2b565b6001600160a01b0383166000818152600660209081526040808320888452825280832080546001600160401b031916600190811790915560078352818420600201805491820181558452919092200180546001600160a01b0319169092179091555184907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e190612c5e908690615a96565b60405180910390a250505050565b600081604051602001612c7f9190615aff565b604051602081830303815290604052805190602001209050919050565b60008281526007602052604090205482906001600160a01b031680612cd457604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612cff5780604051636c51fda960e11b8152600401610bcc9190615a96565b600f54600160301b900460ff1615612d2a5760405163769dd35360e11b815260040160405180910390fd5b612d3384611682565b15612d5157604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b03831660009081526006602090815260408083208784529091529020546001600160401b0316612d9f5783836040516379bfd40160e01b8152600401610bcc929190615c51565b600084815260076020908152604080832060020180548251818502810185019093528083529192909190830182828015612e0257602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612de4575b50505050509050600060018251612e19919061605e565b905060005b8251811015612f2557856001600160a01b0316838281518110612e4357612e43616186565b60200260200101516001600160a01b03161415612f13576000838381518110612e6e57612e6e616186565b6020026020010151905080600760008a81526020019081526020016000206002018381548110612ea057612ea0616186565b600091825260208083209190910180546001600160a01b0319166001600160a01b039490941693909317909255898152600790915260409020600201805480612eeb57612eeb616170565b600082815260209020810160001990810180546001600160a01b031916905501905550612f25565b80612f1d816160ee565b915050612e1e565b506001600160a01b03851660009081526006602090815260408083208984529091529081902080546001600160401b03191690555186907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a790612728908890615a96565b60118181548110612f9957600080fd5b600091825260209091200154905081565b60008281526007602052604090205482906001600160a01b031680612fe257604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461300d5780604051636c51fda960e11b8152600401610bcc9190615a96565b600f54600160301b900460ff16156130385760405163769dd35360e11b815260040160405180910390fd5b6000848152600760205260409020600101546001600160a01b03848116911614610e2b576000848152600760205260409081902060010180546001600160a01b0319166001600160a01b0386161790555184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a190612c5e9033908790615ac3565b6000818152600760205260408120548190819081906060906001600160a01b03166130f857604051630fb532db60e11b815260040160405180910390fd5b60008681526008602090815260408083205460078352928190208054600290910180548351818602810186019094528084526001600160601b0380871696600160601b810490911695600160c01b9091046001600160401b0316946001600160a01b039094169391839183018282801561319b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161317d575b505050505090509450945094509450945091939590929450565b6131bd613493565b6002546040516370a0823160e01b81526000916001600160a01b0316906370a08231906131ee903090600401615a96565b60206040518083038186803b15801561320657600080fd5b505afa15801561321a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061323e91906155fb565b600c549091506001600160601b0316818111156132725780826040516354ced18160e11b8152600401610bcc929190615b20565b8181101561199a576000613286828461605e565b60025460405163a9059cbb60e01b81529192506001600160a01b03169063a9059cbb906132b99087908590600401615aaa565b602060405180830381600087803b1580156132d357600080fd5b505af11580156132e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061330b91906155c5565b507f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600848260405161333d929190615aaa565b60405180910390a150505050565b600f54600160301b900460ff16156133765760405163769dd35360e11b815260040160405180910390fd5b6000818152600760205260409020546001600160a01b03166133ab57604051630fb532db60e11b815260040160405180910390fd5b60008181526008602052604090208054600160601b90046001600160601b0316903490600c6133da8385615f1b565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600c808282829054906101000a90046001600160601b03166134219190615f1b565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f3f1ddc3ab1bdb39001ad76ca51a0e6f57ce6627c69f251d1de41622847721cde8234846134749190615ed8565b604051612ac9929190615b20565b61348a613493565b610bd581613e76565b6000546001600160a01b031633146134e65760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610bcc565b565b6000806134f4846139f7565b60025460405163a9059cbb60e01b81529294509092506001600160a01b03169063a9059cbb906135339086906001600160601b03871690600401615aaa565b602060405180830381600087803b15801561354d57600080fd5b505af1158015613561573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061358591906155c5565b6135a257604051631e9acf1760e31b815260040160405180910390fd5b6000836001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d80600081146135f8576040519150601f19603f3d011682016040523d82523d6000602084013e6135fd565b606091505b505090508061361f57604051630dcf35db60e41b815260040160405180910390fd5b604080516001600160a01b03861681526001600160601b038581166020830152841681830152905186917f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4919081900360600190a25050505050565b604080516060810182526000808252602082018190529181019190915260006136a78460000151612c6c565b6000818152601060205260409020549091506001600160a01b0316806136e357604051631dfd6e1360e21b815260048101839052602401610bcc565b60008286608001516040516020016136fc929190615b20565b60408051601f198184030181529181528151602092830120600081815260129093529120549091508061374257604051631b44092560e11b815260040160405180910390fd5b85516020808801516040808a015160608b015160808c015160a08d01519351613771978a979096959101615d20565b6040516020818303038152906040528051906020012081146137a65760405163354a450b60e21b815260040160405180910390fd5b60006137b58760000151613f1a565b90508061387c578651604051631d2827a760e31b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163e9413d38916138099190600401615d8b565b60206040518083038186803b15801561382157600080fd5b505afa158015613835573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061385991906155fb565b90508061387c57865160405163175dadad60e01b8152610bcc9190600401615d8b565b600088608001518260405160200161389e929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c905060006138c58a83614000565b604080516060810182529889526020890196909652948701949094525093979650505050505050565b60005a61138881101561390057600080fd5b61138881039050846040820482031161391857600080fd5b50823b61392457600080fd5b60008083516020850160008789f190505b9392505050565b6000811561396a576016546139639086908690600160281b900464ffffffffff168661406b565b9050613985565b601654613982908690869064ffffffffff16866140d9565b90505b949350505050565b6000805b6017548110156139ee57826001600160a01b0316601782815481106139b8576139b8616186565b6000918252602090912001546001600160a01b031614156139dc5750600192915050565b806139e6816160ee565b915050613991565b50600092915050565b6000818152600760209081526040808320815160608101835281546001600160a01b03908116825260018301541681850152600282018054845181870281018701865281815287968796949594860193919290830182828015613a8357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613a65575b505050919092525050506000858152600860209081526040808320815160608101835290546001600160601b03808216808452600160601b8304909116948301859052600160c01b9091046001600160401b0316928201929092529096509094509192505b826040015151811015613b5f576006600084604001518381518110613b0f57613b0f616186565b6020908102919091018101516001600160a01b031682528181019290925260409081016000908120898252909252902080546001600160401b031916905580613b57816160ee565b915050613ae8565b50600085815260076020526040812080546001600160a01b03199081168255600182018054909116905590613b9760028301826151a4565b5050600085815260086020526040812055613bb3600a866141ca565b50600c8054859190600090613bd29084906001600160601b0316616098565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555082600c808282829054906101000a90046001600160601b0316613c199190616098565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050915091565b60408051602081018690526001600160a01b03851691810191909152606081018390526001600160401b03821660808201526000908190819060a00160408051601f198184030181529082905280516020918201209250613ca9918991849101615b20565b60408051808303601f19018152919052805160209091012097909650945050505050565b60408051602081019091526000815281613cf65750604080516020810190915260008152611306565b63125fa26760e31b613d0883856160c0565b6001600160e01b03191614613d3057604051632923fee760e11b815260040160405180910390fd5b613d3d8260048186615eae565b8101906139359190615614565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401613d8391511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b60004661a4b1811480613dd0575062066eed81145b15613e4d5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613e0f57600080fd5b505afa158015613e23573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e4791906155fb565b91505090565b4391505090565b600061393583836141d6565b6000611306825490565b60006139358383614225565b6001600160a01b038116331415613ec95760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610bcc565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b1811480613f2f575062066eed81145b15613ff157610100836001600160401b0316613f49613dbb565b613f53919061605e565b1180613f6f5750613f62613dbb565b836001600160401b031610155b15613f7d5750600092915050565b6040516315a03d4160e11b8152606490632b407a8290613fa1908690600401615d8b565b60206040518083038186803b158015613fb957600080fd5b505afa158015613fcd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061393591906155fb565b50506001600160401b03164090565b60006140348360000151846020015185604001518660600151868860a001518960c001518a60e001518b610100015161424f565b6003836020015160405160200161404c929190615c68565b60408051601f1981840301815291905280516020909101209392505050565b60008061407661446a565b905060005a6140858888615ed8565b61408f919061605e565b614099908561603f565b6005549091506000906140b69087906001600160a01b03166144c6565b9050826140c38284615ed8565b6140cd9190615ed8565b98975050505050505050565b6000806140e46145e6565b90506000811361410a576040516321ea67b360e11b815260048101829052602401610bcc565b600061411461446a565b9050600082825a6141258b8b615ed8565b61412f919061605e565b614139908861603f565b6141439190615ed8565b61415590670de0b6b3a764000061603f565b61415f9190615f3d565b60035490915060009061417c9088906001600160a01b03166144c6565b905061419381676765c793fa10079d601b1b61605e565b8211156141b35760405163e80fa38160e01b815260040160405180910390fd5b6141bd8183615ed8565b9998505050505050505050565b600061393583836146b1565b600081815260018301602052604081205461421d57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611306565b506000611306565b600082600001828154811061423c5761423c616186565b9060005260206000200154905092915050565b614258896147a4565b6142a15760405162461bcd60e51b815260206004820152601a6024820152797075626c6963206b6579206973206e6f74206f6e20637572766560301b6044820152606401610bcc565b6142aa886147a4565b6142ee5760405162461bcd60e51b815260206004820152601560248201527467616d6d61206973206e6f74206f6e20637572766560581b6044820152606401610bcc565b6142f7836147a4565b6143435760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610bcc565b61434c826147a4565b6143975760405162461bcd60e51b815260206004820152601c60248201527b73486173685769746e657373206973206e6f74206f6e20637572766560201b6044820152606401610bcc565b6143a3878a8887614867565b6143eb5760405162461bcd60e51b81526020600482015260196024820152786164647228632a706b2b732a6729213d5f755769746e65737360381b6044820152606401610bcc565b60006143f78a8761497b565b9050600061440a898b878b8689896149df565b9050600061441b838d8d8a86614af2565b9050808a1461445c5760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610bcc565b505050505050505050505050565b60004661a4b181148061447f575062066eed81145b156144be57606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b158015613e0f57600080fd5b600091505090565b600064ffffffffff83166144dc57506000611306565b6000806144e884614b32565b90925090506000821361451257838260405163c3388fe760e01b8152600401610bcc929190615aaa565b600860ff8216101561456d57600061452b826008616075565b905061453881600a615f94565b614542908461603f565b61455b64ffffffffff8816670de0b6b3a764000061603f565b6145659190615f3d565b9350506145de565b600860ff821611156145b7576000614586600883616075565b90508261459482600a615f94565b6145ad64ffffffffff8916670de0b6b3a764000061603f565b61455b919061603f565b816145d164ffffffffff8716670de0b6b3a764000061603f565b6145db9190615f3d565b92505b505092915050565b600f5460035460408051633fabe5a360e21b81529051600093600160381b900463ffffffff169283151592859283926001600160a01b03169163feaf968c9160048083019260a0929190829003018186803b15801561464457600080fd5b505afa158015614658573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061467c91906158b2565b5094509092508491505080156146a05750614697824261605e565b8463ffffffff16105b156139855750601354949350505050565b6000818152600183016020526040812054801561479a5760006146d560018361605e565b85549091506000906146e99060019061605e565b905081811461474e57600086600001828154811061470957614709616186565b906000526020600020015490508087600001848154811061472c5761472c616186565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061475f5761475f616170565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611306565b6000915050611306565b80516000906401000003d019116147f25760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420782d6f7264696e61746560701b6044820152606401610bcc565b60208201516401000003d019116148405760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420792d6f7264696e61746560701b6044820152606401610bcc565b60208201516401000003d0199080096148608360005b6020020151614cb0565b1492915050565b60006001600160a01b0382166148ad5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610bcc565b6020840151600090600116156148c457601c6148c7565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe199182039250600091908909875160408051600080825260209091019182905292935060019161493191869188918790615b2e565b6020604051602081039080840390855afa158015614953573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b6149836151c2565b6149b06001848460405160200161499c93929190615a75565b604051602081830303815290604052614cd4565b90505b6149bc816147a4565b6113065780516040805160208101929092526149d8910161499c565b90506149b3565b6149e76151c2565b825186516401000003d0199081900691061415614a465760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610bcc565b614a51878988614d22565b614a965760405162461bcd60e51b8152602060048201526016602482015275119a5c9cdd081b5d5b0818da1958dac819985a5b195960521b6044820152606401610bcc565b614aa1848685614d22565b614ae75760405162461bcd60e51b815260206004820152601760248201527614d958dbdb99081b5d5b0818da1958dac819985a5b1959604a1b6044820152606401610bcc565b6140cd868484614e3d565b600060028686868587604051602001614b1096959493929190615a1b565b60408051601f1981840301815291905280516020909101209695505050505050565b6000806000600f60000160079054906101000a900463ffffffff1690506000808263ffffffff161190506000856001600160a01b031663feaf968c6040518163ffffffff1660e01b815260040160a06040518083038186803b158015614b9757600080fd5b505afa158015614bab573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614bcf91906158b2565b50919750909250839150508015614bf45750614beb814261605e565b8363ffffffff16105b15614c35576005546001600160a01b0387811691161415614c19576014549450614c35565b6003546001600160a01b03878116911614156102435760135494505b856001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b158015614c6e57600080fd5b505afa158015614c82573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614ca69190615902565b9350505050915091565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614cdc6151c2565b614ce582614f00565b8152614cfa614cf5826000614856565b614f3b565b602082018190526002900660011415612348576020810180516401000003d019039052919050565b600082614d5f5760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610bcc565b83516020850151600090614d7590600290616130565b15614d8157601c614d84565b601b5b9050600070014551231950b75fc4402da1732fc9bebe19838709604080516000808252602090910191829052919250600190614dc7908390869088908790615b2e565b6020604051602081039080840390855afa158015614de9573d6000803e3d6000fd5b505050602060405103519050600086604051602001614e089190615a09565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614e456151c2565b835160208086015185519186015160009384938493614e6693909190614f5b565b919450925090506401000003d019858209600114614ec25760405162461bcd60e51b815260206004820152601960248201527834b73b2d1036bab9ba1031329034b73b32b939b29037b3103d60391b6044820152606401610bcc565b60405180604001604052806401000003d01980614ee157614ee161615a565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d019811061234857604080516020808201939093528151808203840181529082019091528051910120614f08565b6000611306826002614f546401000003d0196001615ed8565b901c61503b565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614f9b838385856150d2565b9098509050614fac88828e886150f6565b9098509050614fbd88828c876150f6565b90985090506000614fd08d878b856150f6565b9098509050614fe1888286866150d2565b9098509050614ff288828e896150f6565b9098509050818114615027576401000003d019818a0998506401000003d01982890997506401000003d019818309965061502b565b8196505b5050505050509450945094915050565b6000806150466151e0565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a08201526150786151fe565b60208160c0846005600019fa9250826150c85760405162461bcd60e51b81526020600482015260126024820152716269674d6f64457870206661696c7572652160701b6044820152606401610bcc565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215615194579160200282015b8281111561519457825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019061515f565b506151a092915061521c565b5090565b5080546000825590600052602060002090810190610bd5919061521c565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b808211156151a0576000815560010161521d565b8035612348816161b2565b806040810183101561130657600080fd5b600082601f83011261525e57600080fd5b615266615e41565b80838560408601111561527857600080fd5b60005b600281101561529a57813584526020938401939091019060010161527b565b509095945050505050565b600082601f8301126152b657600080fd5b81356001600160401b03808211156152d0576152d061619c565b604051601f8301601f19908116603f011681019082821181831017156152f8576152f861619c565b8160405283815286602085880101111561531157600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060c0828403121561534357600080fd5b61534b615e69565b905081356001600160401b03808216821461536557600080fd5b8183526020840135602084015261537e604085016153e4565b604084015261538f606085016153e4565b60608401526153a060808501615231565b608084015260a08401359150808211156153b957600080fd5b506153c6848285016152a5565b60a08301525092915050565b803561ffff8116811461234857600080fd5b803563ffffffff8116811461234857600080fd5b803564ffffffffff8116811461234857600080fd5b80516001600160501b038116811461234857600080fd5b80356001600160601b038116811461234857600080fd5b60006020828403121561544d57600080fd5b8135613935816161b2565b6000806040838503121561546b57600080fd5b8235615476816161b2565b915061548460208401615424565b90509250929050565b6000806000606084860312156154a257600080fd5b83356154ad816161b2565b925060208401356154bd816161b2565b915060408401356154cd816161b2565b809150509250925092565b600080606083850312156154eb57600080fd5b82356154f6816161b2565b9150615484846020850161523c565b6000806000806060858703121561551b57600080fd5b8435615526816161b2565b93506020850135925060408501356001600160401b038082111561554957600080fd5b818701915087601f83011261555d57600080fd5b81358181111561556c57600080fd5b88602082850101111561557e57600080fd5b95989497505060200194505050565b60006040828403121561559f57600080fd5b613935838361523c565b6000604082840312156155bb57600080fd5b613935838361524d565b6000602082840312156155d757600080fd5b8151613935816161c7565b6000602082840312156155f457600080fd5b5035919050565b60006020828403121561560d57600080fd5b5051919050565b60006020828403121561562657600080fd5b604051602081016001600160401b03811182821017156156485761564861619c565b6040528235615656816161c7565b81529392505050565b6000808284036101c081121561567457600080fd5b6101a08082121561568457600080fd5b61568c615e8b565b9150615698868661524d565b82526156a7866040870161524d565b60208301526080850135604083015260a0850135606083015260c085013560808301526156d660e08601615231565b60a08301526101006156ea8782880161524d565b60c08401526156fd87610140880161524d565b60e0840152610180860135908301529092508301356001600160401b0381111561572657600080fd5b61573285828601615331565b9150509250929050565b60006020828403121561574e57600080fd5b81356001600160401b0381111561576457600080fd5b820160c0818503121561393557600080fd5b60006020828403121561578857600080fd5b613935826153d2565b600080600080600080600080888a036101208112156157af57600080fd5b6157b88a6153d2565b98506157c660208b016153e4565b97506157d460408b016153e4565b96506157e260608b016153e4565b955060808a0135945060a08a0135935060c08a01359250604060df198201121561580b57600080fd5b50615814615e41565b61582060e08b016153f8565b815261582f6101008b016153f8565b6020820152809150509295985092959890939650565b6000806040838503121561585857600080fd5b82359150602083013561586a816161b2565b809150509250929050565b6000806040838503121561588857600080fd5b50508035926020909101359150565b6000602082840312156158a957600080fd5b613935826153e4565b600080600080600060a086880312156158ca57600080fd5b6158d38661540d565b94506020860151935060408601519250606086015191506158f66080870161540d565b90509295509295909350565b60006020828403121561591457600080fd5b815160ff8116811461393557600080fd5b600081518084526020808501945080840160005b8381101561595e5781516001600160a01b031687529582019590820190600101615939565b509495945050505050565b8060005b6002811015610e2b57815184526020938401939091019060010161596d565b600081518084526020808501945080840160005b8381101561595e578151875295820195908201906001016159a0565b6000815180845260005b818110156159e2576020818501810151868301820152016159c6565b818111156159f4576000602083870101525b50601f01601f19169290920160200192915050565b615a138183615969565b604001919050565b868152615a2b6020820187615969565b615a386060820186615969565b615a4560a0820185615969565b615a5260e0820184615969565b60609190911b6001600160601b0319166101208201526101340195945050505050565b838152615a856020820184615969565b606081019190915260800192915050565b6001600160a01b0391909116815260200190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b039290921682526001600160601b0316602082015260400190565b604081016113068284615969565b602081526000613935602083018461598c565b918252602082015260400190565b93845260ff9290921660208401526040830152606082015260800190565b60208152600061393560208301846159bc565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160c06080840152615ba460e0840182615925565b60808501516001600160601b0390811660a0868101919091529095015190941660c0909301929092525090919050565b61ffff93841681529183166020830152909116604082015260600190565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b81811015615c4357845183529383019391830191600101615c27565b509098975050505050505050565b9182526001600160a01b0316602082015260400190565b828152606081016139356020830184615969565b828152604060208201526000613985604083018461598c565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a08301526140cd60c08301846159bc565b878152602081018790526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c082018190526000906141bd908301846159bc565b8781526001600160401b03871660208201526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c082018190526000906141bd908301846159bc565b63ffffffff92831681529116602082015260400190565b6001600160401b0391909116815260200190565b6001600160601b038681168252851660208201526001600160401b03841660408201526001600160a01b038316606082015260a060808201819052600090615de990830184615925565b979650505050505050565b6000808335601e19843603018112615e0b57600080fd5b8301803591506001600160401b03821115615e2557600080fd5b602001915036819003821315615e3a57600080fd5b9250929050565b604080519081016001600160401b0381118282101715615e6357615e6361619c565b60405290565b60405160c081016001600160401b0381118282101715615e6357615e6361619c565b60405161012081016001600160401b0381118282101715615e6357615e6361619c565b60008085851115615ebe57600080fd5b83861115615ecb57600080fd5b5050820193919092039150565b60008219821115615eeb57615eeb616144565b500190565b60006001600160401b03828116848216808303821115615f1257615f12616144565b01949350505050565b60006001600160601b03828116848216808303821115615f1257615f12616144565b600082615f4c57615f4c61615a565b500490565b600181815b80851115615f8c578160001904821115615f7257615f72616144565b80851615615f7f57918102915b93841c9390800290615f56565b509250929050565b600061393560ff841683600082615fad57506001611306565b81615fba57506000611306565b8160018114615fd05760028114615fda57615ff6565b6001915050611306565b60ff841115615feb57615feb616144565b50506001821b611306565b5060208310610133831016604e8410600b8410161715616019575081810a611306565b6160238383615f51565b806000190482111561603757616037616144565b029392505050565b600081600019048311821515161561605957616059616144565b500290565b60008282101561607057616070616144565b500390565b600060ff821660ff84168082101561608f5761608f616144565b90039392505050565b60006001600160601b03838116908316818110156160b8576160b8616144565b039392505050565b6001600160e01b031981358181169160048510156145de5760049490940360031b84901b1690921692915050565b600060001982141561610257616102616144565b5060010190565b60006001600160401b038281168082141561612657616126616144565b6001019392505050565b60008261613f5761613f61615a565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114610bd557600080fd5b8015158114610bd557600080fdfea164736f6c6343000806000a",
}

var VRFCoordinatorV2PlusABI = VRFCoordinatorV2PlusMetaData.ABI

var VRFCoordinatorV2PlusBin = VRFCoordinatorV2PlusMetaData.Bin

func DeployVRFCoordinatorV2Plus(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV2Plus, error) {
	parsed, err := VRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2PlusBin), backend, blockhashStore)
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

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) ETHUSDFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "ETH_USD_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) ETHUSDFEED() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.ETHUSDFEED(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) ETHUSDFEED() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.ETHUSDFEED(&_VRFCoordinatorV2Plus.CallOpts)
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

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.LINKETHFEED(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.LINKETHFEED(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) LINKUSDFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "LINK_USD_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) LINKUSDFEED() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.LINKUSDFEED(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) LINKUSDFEED() (common.Address, error) {
	return _VRFCoordinatorV2Plus.Contract.LINKUSDFEED(&_VRFCoordinatorV2Plus.CallOpts)
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
	outstruct.StalenessSeconds = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[4], new(uint32)).(*uint32)

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

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SFallbackUSDPerUnitEth(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_fallbackUSDPerUnitEth")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SFallbackUSDPerUnitEth() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.SFallbackUSDPerUnitEth(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SFallbackUSDPerUnitEth() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.SFallbackUSDPerUnitEth(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SFallbackUSDPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_fallbackUSDPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SFallbackUSDPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.SFallbackUSDPerUnitLink(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SFallbackUSDPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.SFallbackUSDPerUnitLink(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_fallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2Plus.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCaller) SFeeConfig(opts *bind.CallOpts) (SFeeConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2Plus.contract.Call(opts, &out, "s_feeConfig")

	outstruct := new(SFeeConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.FulfillmentFlatFeeLinkUSD = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FulfillmentFlatFeeEthUSD = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SFeeConfig() (SFeeConfig,

	error) {
	return _VRFCoordinatorV2Plus.Contract.SFeeConfig(&_VRFCoordinatorV2Plus.CallOpts)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusCallerSession) SFeeConfig() (SFeeConfig,

	error) {
	return _VRFCoordinatorV2Plus.Contract.SFeeConfig(&_VRFCoordinatorV2Plus.CallOpts)
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

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fallbackUSDPerUnitEth, fallbackUSDPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetConfig(&_VRFCoordinatorV2Plus.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fallbackUSDPerUnitEth, fallbackUSDPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetConfig(&_VRFCoordinatorV2Plus.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fallbackUSDPerUnitEth, fallbackUSDPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) SetLINKAndFeeds(opts *bind.TransactOpts, link common.Address, linkEthFeed common.Address, ethUSDFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "setLINKAndFeeds", link, linkEthFeed, ethUSDFeed)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SetLINKAndFeeds(link common.Address, linkEthFeed common.Address, ethUSDFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetLINKAndFeeds(&_VRFCoordinatorV2Plus.TransactOpts, link, linkEthFeed, ethUSDFeed)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) SetLINKAndFeeds(link common.Address, linkEthFeed common.Address, ethUSDFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetLINKAndFeeds(&_VRFCoordinatorV2Plus.TransactOpts, link, linkEthFeed, ethUSDFeed)
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
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	FallbackUSDPerUnitEth       *big.Int
	FallbackUSDPerUnitLink      *big.Int
	FeeConfig                   VRFCoordinatorV2PlusFeeConfig
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
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
}
type SFeeConfig struct {
	FulfillmentFlatFeeLinkUSD *big.Int
	FulfillmentFlatFeeEthUSD  *big.Int
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
	return common.HexToHash("0x306a7924a41e68523d99ac4f68693ab16b5dba36ff53896a0977971a564cc102")
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

	ETHUSDFEED(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	LINKUSDFEED(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

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

	SFallbackUSDPerUnitEth(opts *bind.CallOpts) (*big.Int, error)

	SFallbackUSDPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SFeeConfig(opts *bind.CallOpts) (SFeeConfig,

		error)

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

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fallbackUSDPerUnitEth *big.Int, fallbackUSDPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error)

	SetLINKAndFeeds(opts *bind.TransactOpts, link common.Address, linkEthFeed common.Address, ethUSDFeed common.Address) (*types.Transaction, error)

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
