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
	FulfillmentFlatFeeLinkPPM uint32
	FulfillmentFlatFeeEthPPM  uint32
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendEther\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeEthPPM\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorV2Plus.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EthFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountEth\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldEthBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newEthBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithEth\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFCoordinatorV2Plus.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"ethBalance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdrawEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverEthFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeEthPPM\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalEthBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeEthPPM\",\"type\":\"uint32\"}],\"internalType\":\"structVRFCoordinatorV2Plus.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKETHFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200609738038062006097833981016040819052620000349162000188565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000dc565b505060016002555060601b6001600160601b031916608052620001ba565b6001600160a01b038116331415620001375760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019b57600080fd5b81516001600160a01b0381168114620001b357600080fd5b9392505050565b60805160601c615eb7620001e0600039600081816105f80152613aef0152615eb76000f3fe6080604052600436106102f15760003560e01c806386fe91c71161018f578063bec4c08c116100e1578063dc311dd31161008a578063e95704bd11610064578063e95704bd14610947578063ee9d2d381461096e578063f2fde38b1461099b57600080fd5b8063dc311dd3146108e4578063e72f6e3014610914578063e8509bff1461093457600080fd5b8063d98e620e116100bb578063d98e620e1461086e578063da2f26101461088e578063dac83d29146108c457600080fd5b8063bec4c08c1461080e578063caf70c4a1461082e578063cb6317971461084e57600080fd5b8063a4c0ed3611610143578063ad1783611161011d578063ad178361146107ae578063b08c8795146107ce578063b2a7cac5146107ee57600080fd5b8063a4c0ed361461074e578063a8cb447b1461076e578063aa433aff1461078e57600080fd5b80639b1c385e116101745780639b1c385e146106df5780639d40a6fd146106ff578063a21a23e41461073957600080fd5b806386fe91c7146106955780638da5cb5b146106c157600080fd5b806340d6bb821161024857806364d51a2a116101fc5780636b6feccc116101d65780636b6feccc1461061a5780636f64f03f1461066057806379ba50971461068057600080fd5b806364d51a2a146105b157806366316d8d146105c6578063689c4517146105e657600080fd5b806346d8d4861161022d57806346d8d4861461055157806357133e64146105715780635d06b4ab1461059157600080fd5b806340d6bb82146104f657806341af6c871461052157600080fd5b80630ae09540116102aa578063294daa4911610284578063294daa4914610482578063330987b31461049e578063405b84fa146104d657600080fd5b80630ae095401461040257806315c48b84146104225780631b6b6d231461044a57600080fd5b8063043bd6ae116102db578063043bd6ae14610345578063088070f51461036957806308821d58146103e257600080fd5b8062012291146102f657806304104edb14610323575b600080fd5b34801561030257600080fd5b5061030b6109bb565b60405161031a93929190615992565b60405180910390f35b34801561032f57600080fd5b5061034361033e3660046152f8565b610a37565b005b34801561035157600080fd5b5061035b600e5481565b60405190815260200161031a565b34801561037557600080fd5b50600f546103b19061ffff81169063ffffffff620100008204811691660100000000000081048216916a01000000000000000000009091041684565b6040805161ffff909516855263ffffffff93841660208601529183169184019190915216606082015260800161031a565b3480156103ee57600080fd5b506103436103fd366004615439565b610bf9565b34801561040e57600080fd5b5061034361041d3660046156f6565b610d8d565b34801561042e57600080fd5b5061043760c881565b60405161ffff909116815260200161031a565b34801561045657600080fd5b5060035461046a906001600160a01b031681565b6040516001600160a01b03909116815260200161031a565b34801561048e57600080fd5b506040516001815260200161031a565b3480156104aa57600080fd5b506104be6104b936600461550c565b610e79565b6040516001600160601b03909116815260200161031a565b3480156104e257600080fd5b506103436104f13660046156f6565b611338565b34801561050257600080fd5b5061050c6101f481565b60405163ffffffff909116815260200161031a565b34801561052d57600080fd5b5061054161053c36600461548e565b611717565b604051901515815260200161031a565b34801561055d57600080fd5b5061034361056c366004615315565b611919565b34801561057d57600080fd5b5061034361058c36600461534a565b611ab4565b34801561059d57600080fd5b506103436105ac3660046152f8565b611b2d565b3480156105bd57600080fd5b50610437606481565b3480156105d257600080fd5b506103436105e1366004615315565b611c04565b3480156105f257600080fd5b5061046a7f000000000000000000000000000000000000000000000000000000000000000081565b34801561062657600080fd5b506010546106439063ffffffff8082169164010000000090041682565b6040805163ffffffff93841681529290911660208301520161031a565b34801561066c57600080fd5b5061034361067b366004615383565b611dc0565b34801561068c57600080fd5b50610343611ed8565b3480156106a157600080fd5b506008546104be906801000000000000000090046001600160601b031681565b3480156106cd57600080fd5b506000546001600160a01b031661046a565b3480156106eb57600080fd5b5061035b6106fa366004615608565b611f89565b34801561070b57600080fd5b506008546107209067ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161031a565b34801561074557600080fd5b5061035b6123ca565b34801561075a57600080fd5b506103436107693660046153b0565b61262a565b34801561077a57600080fd5b506103436107893660046152f8565b612819565b34801561079a57600080fd5b506103436107a936600461548e565b612935565b3480156107ba57600080fd5b5060045461046a906001600160a01b031681565b3480156107da57600080fd5b506103436107e9366004615658565b612995565b3480156107fa57600080fd5b5061034361080936600461548e565b612b1f565b34801561081a57600080fd5b506103436108293660046156f6565b612c8c565b34801561083a57600080fd5b5061035b610849366004615455565b612e5f565b34801561085a57600080fd5b506103436108693660046156f6565b612e8f565b34801561087a57600080fd5b5061035b61088936600461548e565b6131ac565b34801561089a57600080fd5b5061046a6108a936600461548e565b600b602052600090815260409020546001600160a01b031681565b3480156108d057600080fd5b506103436108df3660046156f6565b6131cd565b3480156108f057600080fd5b506109046108ff36600461548e565b613305565b60405161031a9493929190615b90565b34801561092057600080fd5b5061034361092f3660046152f8565b6133e8565b61034361094236600461548e565b6135b0565b34801561095357600080fd5b506008546104be90600160a01b90046001600160601b031681565b34801561097a57600080fd5b5061035b61098936600461548e565b600d6020526000908152604090205481565b3480156109a757600080fd5b506103436109b63660046152f8565b613708565b600f54600c805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff16939192839190830182828015610a2557602002820191906000526020600020905b815481526020019060010190808311610a11575b50505050509050925092509250909192565b610a3f613719565b60115460005b81811015610bcc57826001600160a01b031660118281548110610a6a57610a6a615e3b565b6000918252602090912001546001600160a01b03161415610bba576011610a92600184615d33565b81548110610aa257610aa2615e3b565b600091825260209091200154601180546001600160a01b039092169183908110610ace57610ace615e3b565b600091825260209091200180546001600160a01b0319166001600160a01b0392909216919091179055826011610b05600185615d33565b81548110610b1557610b15615e3b565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506011805480610b5457610b54615e25565b6000828152602090819020600019908301810180546001600160a01b03191690559091019091556040516001600160a01b03851681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37910160405180910390a1505050565b80610bc481615da2565b915050610a45565b50604051635428d44960e01b81526001600160a01b03831660048201526024015b60405180910390fd5b50565b610c01613719565b604080518082018252600091610c30919084906002908390839080828437600092019190915250612e5f915050565b6000818152600b60205260409020549091506001600160a01b031680610c6c57604051631dfd6e1360e21b815260048101839052602401610bed565b6000828152600b6020526040812080546001600160a01b03191690555b600c54811015610d445782600c8281548110610ca757610ca7615e3b565b90600052602060002001541415610d3257600c805460009190610ccc90600190615d33565b81548110610cdc57610cdc615e3b565b9060005260206000200154905080600c8381548110610cfd57610cfd615e3b565b600091825260209091200155600c805480610d1a57610d1a615e25565b60019003818190600052602060002001600090559055505b80610d3c81615da2565b915050610c89565b50806001600160a01b03167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610d8091815260200190565b60405180910390a2505050565b60008281526006602052604090205482906001600160a01b031680610dc557604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614610df957604051636c51fda960e11b81526001600160a01b0382166004820152602401610bed565b600280541415610e395760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b60028055610e4684611717565b15610e6457604051631685ecdd60e31b815260040160405180910390fd5b610e6e8484613775565b505060016002555050565b6000600280541415610ebb5760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b6002805560005a90506000610ed0858561390c565b90506000610ee4608086016060870161571b565b63ffffffff1667ffffffffffffffff811115610f0257610f02615e51565b604051908082528060200260200182016040528015610f2b578160200160208202803683370190505b50905060005b610f41608087016060880161571b565b63ffffffff16811015610fb657826040015181604051602001610f6e929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c828281518110610f9957610f99615e3b565b602090810291909101015280610fae81615da2565b915050610f31565b50602080830180516000908152600d90925260408083208390559051905182917f1fe543e3000000000000000000000000000000000000000000000000000000009161100791908690602401615a05565b60408051601f198184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff166001600160e01b0319909416939093179092529150600090611083906110679060608b01908b0161571b565b63ffffffff1661107d60a08b0160808c016152f8565b84613c61565b9050600061109c61109760a08b018b615bd5565b613caf565b51600f549091506000906110c79089906a0100000000000000000000900463ffffffff163a85613d5e565b905081156111d0576020808b01356000908152600790915260409020546001600160601b03808316600160601b90920416101561111757604051631e9acf1760e31b815260040160405180910390fd5b60208a81013560009081526007909152604090208054829190600c9061114e908490600160601b90046001600160601b0316615d4a565b82546101009290920a6001600160601b0381810219909316918316021790915588516000908152600b60209081526040808320546001600160a01b03168352600a9091528120805485945090926111a791859116615cde565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506112bc565b6020808b01356000908152600790915260409020546001600160601b038083169116101561121157604051631e9acf1760e31b815260040160405180910390fd5b6020808b01356000908152600790915260408120805483929061123e9084906001600160601b0316615d4a565b82546101009290920a6001600160601b0381810219909316918316021790915588516000908152600b60209081526040808320546001600160a01b03168352600990915281208054859450909261129791859116615cde565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b60006112d76040518060200160405280851515815250613dae565b90508a6020013588602001517f7bad4a5a5566568e03ff9e9fdaedf6b8a760aa3ffc0abad348cf93787f04e17e8a6040015185858960405161131c9493929190615b55565b60405180910390a35060016002559a9950505050505050505050565b61134181613e34565b61136957604051635428d44960e01b81526001600160a01b0382166004820152602401610bed565b60008060008061137886613305565b9350935093509350336001600160a01b0316826001600160a01b0316146113e15760405162461bcd60e51b815260206004820152601660248201527f4e6f7420737562736372697074696f6e206f776e6572000000000000000000006044820152606401610bed565b6113ea86611717565b156114375760405162461bcd60e51b815260206004820152601660248201527f50656e64696e67207265717565737420657869737473000000000000000000006044820152606401610bed565b60006040518060c0016040528061144c600190565b60ff168152602001888152602001846001600160a01b03168152602001838152602001866001600160601b03168152602001856001600160601b031681525090506000816040516020016114a0919061591d565b60405160208183030381529060405290506114ba88613e9e565b50506040517fce3f47190000000000000000000000000000000000000000000000000000000081526001600160a01b0388169063ce3f4719906001600160601b0388169061150c90859060040161590a565b6000604051808303818588803b15801561152557600080fd5b505af1158015611539573d6000803e3d6000fd5b505060035460405163a9059cbb60e01b81526001600160a01b038c811660048301526001600160601b038c166024830152909116935063a9059cbb92506044019050602060405180830381600087803b15801561159557600080fd5b505af11580156115a9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115cd9190615471565b6116195760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610bed565b60005b83518110156116ca5783818151811061163757611637615e3b565b60209081029190910101516040517f8ea981170000000000000000000000000000000000000000000000000000000081526001600160a01b038a8116600483015290911690638ea9811790602401600060405180830381600087803b15801561169f57600080fd5b505af11580156116b3573d6000803e3d6000fd5b5050505080806116c290615da2565b91505061161c565b50604080516001600160a01b0389168152602081018a90527fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187910160405180910390a15050505050505050565b6000818152600660209081526040808320815160608101835281546001600160a01b03908116825260018301541681850152600282018054845181870281018701865281815287969395860193909291908301828280156117a157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611783575b505050505081525050905060005b81604001515181101561190f5760005b600c548110156118fc5760006118c5600c83815481106117e1576117e1615e3b565b90600052602060002001548560400151858151811061180257611802615e3b565b602002602001015188600560008960400151898151811061182557611825615e3b565b6020908102919091018101516001600160a01b03908116835282820193909352604091820160009081208e825282528290205482518083018890529590931685830152606085019390935267ffffffffffffffff9091166080808501919091528151808503909101815260a08401825280519083012060c084019490945260e0808401859052815180850390910181526101009093019052815191012091565b506000818152600d6020526040902054909150156118e95750600195945050505050565b50806118f481615da2565b9150506117bf565b508061190781615da2565b9150506117af565b5060009392505050565b6002805414156119595760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b60028055336000908152600a60205260409020546001600160601b038083169116101561199957604051631e9acf1760e31b815260040160405180910390fd5b336000908152600a6020526040812080548392906119c19084906001600160601b0316615d4a565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600860148282829054906101000a90046001600160601b0316611a099190615d4a565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114611a83576040519150601f19603f3d011682016040523d82523d6000602084013e611a88565b606091505b5050905080611aaa57604051630dcf35db60e41b815260040160405180910390fd5b5050600160025550565b611abc613719565b6003546001600160a01b031615611aff576040517f2d118a6e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380546001600160a01b039384166001600160a01b03199182161790915560048054929093169116179055565b611b35613719565b611b3e81613e34565b15611b80576040517fac8a27ef0000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602401610bed565b601180546001810182556000919091527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c680180546001600160a01b0319166001600160a01b0383169081179091556040519081527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259060200160405180910390a150565b600280541415611c445760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b60028055336000908152600960205260409020546001600160601b0380831691161015611c8457604051631e9acf1760e31b815260040160405180910390fd5b3360009081526009602052604081208054839290611cac9084906001600160601b0316615d4a565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550806008808282829054906101000a90046001600160601b0316611cf39190615d4a565b82546101009290920a6001600160601b0381810219909316918316021790915560035460405163a9059cbb60e01b81526001600160a01b03868116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b158015611d6257600080fd5b505af1158015611d76573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d9a9190615471565b611db757604051631e9acf1760e31b815260040160405180910390fd5b50506001600255565b611dc8613719565b604080518082018252600091611df7919084906002908390839080828437600092019190915250612e5f915050565b6000818152600b60205260409020549091506001600160a01b031615611e4c576040517f4a0b8fa700000000000000000000000000000000000000000000000000000000815260048101829052602401610bed565b6000818152600b6020908152604080832080546001600160a01b0319166001600160a01b038816908117909155600c805460018101825594527fdf6966c971051c3d54ec59162606531493a51404a002842f56009d7e5cf4a8c7909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610d80565b6001546001600160a01b03163314611f325760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610bed565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000600280541415611fcb5760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b600280556020808301356000908152600690915260409020546001600160a01b031661200a57604051630fb532db60e11b815260040160405180910390fd5b33600090815260056020908152604080832085830135845290915290205467ffffffffffffffff168061205c576040516379bfd40160e01b815260208401356004820152336024820152604401610bed565b600f5461ffff16612073606085016040860161563d565b61ffff161080612096575060c8612090606085016040860161563d565b61ffff16115b156120dc576120ab606084016040850161563d565b600f5460405163539c34bb60e11b815261ffff92831660048201529116602482015260c86044820152606401610bed565b600f5462010000900463ffffffff166120fb608085016060860161571b565b63ffffffff16111561216457612117608084016060850161571b565b600f546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff9283166004820152620100009091049091166024820152604401610bed565b6101f461217760a085016080860161571b565b63ffffffff1611156121d65761219360a084016080850161571b565b6040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff90911660048201526101f46024820152604401610bed565b60006121e3826001615cb2565b6040805186356020808301829052338385015280890135606084015267ffffffffffffffff85166080808501919091528451808503909101815260a0808501865281519183019190912060c085019390935260e08085018490528551808603909101815261010090940190945282519201919091209293509060009061226f9061109790890189615bd5565b9050600061227c82613dae565b9050836122876140ef565b60208a013561229c60808c0160608d0161571b565b6122ac60a08d0160808e0161571b565b33866040516020016122c49796959493929190615a92565b60405160208183030381529060405280519060200120600d600086815260200190815260200160002081905550336001600160a01b0316886020013589600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e87878d604001602081019061233b919061563d565b8e606001602081019061234e919061571b565b8f6080016020810190612361919061571b565b8960405161237496959493929190615a53565b60405180910390a4505033600090815260056020908152604080832098820135835297905295909520805467ffffffffffffffff191667ffffffffffffffff939093169290921790915560016002559392505050565b600060028054141561240c5760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b6002805560003361241e600143615d33565b600854604051606093841b6bffffffffffffffffffffffff199081166020830152924060348201523090931b909116605483015260c01b77ffffffffffffffffffffffffffffffffffffffffffffffff1916606882015260700160408051601f1981840301815291905280516020909101206008805491925067ffffffffffffffff9091169060006124af83615dbd565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505060008067ffffffffffffffff8111156124f1576124f1615e51565b60405190808252806020026020018201604052801561251a578160200160208202803683370190505b5060408051808201825260008082526020808301828152878352600782528483209351845491516001600160601b03908116600160601b0277ffffffffffffffffffffffffffffffffffffffffffffffff199093169116171790925582516060810184523381528083018281528185018681528884526006855294909220815181546001600160a01b039182166001600160a01b031991821617835593516001830180549190921694169390931790925592518051949550929390926125e7926002850192910190615123565b50506040513381528391507f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d9060200160405180910390a2509050600160025590565b60028054141561266a5760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b600280556003546001600160a01b031633146126b2576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146126ec576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006126fa8284018461548e565b6000818152600660205260409020549091506001600160a01b031661273257604051630fb532db60e11b815260040160405180910390fd5b600081815260076020526040812080546001600160601b0316918691906127598385615cde565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550846008808282829054906101000a90046001600160601b03166127a09190615cde565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846127f39190615c9a565b604080519283526020830191909152015b60405180910390a25050600160025550505050565b612821613719565b6008544790600160a01b90046001600160601b031681811115612861576040516354ced18160e11b81526004810182905260248101839052604401610bed565b818110156129305760006128758284615d33565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d80600081146128c4576040519150601f19603f3d011682016040523d82523d6000602084013e6128c9565b606091505b50509050806128eb57604051630dcf35db60e41b815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df377317910160405180910390a150505b505050565b61293d613719565b6000818152600660205260409020546001600160a01b031661297257604051630fb532db60e11b815260040160405180910390fd5b600081815260066020526040902054610bf69082906001600160a01b0316613775565b61299d613719565b60c861ffff871611156129d75760405163539c34bb60e11b815261ffff871660048201819052602482015260c86044820152606401610bed565b600082136129fb576040516321ea67b360e11b815260048101839052602401610bed565b604080516080808201835261ffff891680835263ffffffff89811660208086018290528a83168688018190528a84166060978801819052600f805465ffffffffffff19168717620100008602176dffffffffffffffff0000000000001916660100000000000084026dffffffff000000000000000000001916176a010000000000000000000083021790558951601080548c86015192881667ffffffffffffffff1990911617640100000000928816929092029190911790819055600e8c9055895196875286840194909452978501529483019590955291810186905283821660a08201529290911c1660c08201527f777357bb93f63d088f18112d3dba38457aec633eb8f1341e1d418380ad328e789060e00160405180910390a1505050505050565b600280541415612b5f5760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b600280556000818152600660205260409020546001600160a01b0316612b9857604051630fb532db60e11b815260040160405180910390fd5b6000818152600660205260409020600101546001600160a01b03163314612c0a57600081815260066020526040908190206001015490517fd084e9750000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610bed565b6000818152600660209081526040918290208054336001600160a01b0319808316821784556001909301805490931690925583516001600160a01b0390911680825292810191909152909183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691015b60405180910390a250506001600255565b60008281526006602052604090205482906001600160a01b031680612cc457604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612cf857604051636c51fda960e11b81526001600160a01b0382166004820152602401610bed565b600280541415612d385760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b60028080556000858152600660205260409020015460641415612d87576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b038316600090815260056020908152604080832087845290915290205467ffffffffffffffff1615612dbf57610e6e565b6001600160a01b03831660008181526005602090815260408083208884528252808320805467ffffffffffffffff19166001908117909155600683528184206002018054918201815584529282902090920180546001600160a01b03191684179055905191825285917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a2505060016002555050565b600081604051602001612e7291906158fc565b604051602081830303815290604052805190602001209050919050565b60008281526006602052604090205482906001600160a01b031680612ec757604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612efb57604051636c51fda960e11b81526001600160a01b0382166004820152602401610bed565b600280541415612f3b5760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b60028055612f4884611717565b15612f6657604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038316600090815260056020908152604080832087845290915290205467ffffffffffffffff16612fc3576040516379bfd40160e01b8152600481018590526001600160a01b0384166024820152604401610bed565b60008481526006602090815260408083206002018054825181850281018501909352808352919290919083018282801561302657602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613008575b5050505050905060006001825161303d9190615d33565b905060005b825181101561314957856001600160a01b031683828151811061306757613067615e3b565b60200260200101516001600160a01b0316141561313757600083838151811061309257613092615e3b565b6020026020010151905080600660008a815260200190815260200160002060020183815481106130c4576130c4615e3b565b600091825260208083209190910180546001600160a01b0319166001600160a01b03949094169390931790925589815260069091526040902060020180548061310f5761310f615e25565b600082815260209020810160001990810180546001600160a01b031916905501905550613149565b8061314181615da2565b915050613042565b506001600160a01b03851660008181526005602090815260408083208a8452825291829020805467ffffffffffffffff19169055905191825287917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79101612804565b600c81815481106131bc57600080fd5b600091825260209091200154905081565b60008281526006602052604090205482906001600160a01b03168061320557604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461323957604051636c51fda960e11b81526001600160a01b0382166004820152602401610bed565b6002805414156132795760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b600280556000848152600660205260409020600101546001600160a01b03848116911614610e6e5760008481526006602090815260409182902060010180546001600160a01b0319166001600160a01b03871690811790915582513381529182015285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a19101612e4c565b600081815260066020526040812054819081906060906001600160a01b031661334157604051630fb532db60e11b815260040160405180910390fd5b60008581526007602090815260408083205460068352928190208054600290910180548351818602810186019094528084526001600160601b0380871696600160601b900416946001600160a01b039093169391928391908301828280156133d257602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116133b4575b5050505050905093509350935093509193509193565b6133f0613719565b6003546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000916001600160a01b0316906370a082319060240160206040518083038186803b15801561344d57600080fd5b505afa158015613461573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061348591906154a7565b6008549091506801000000000000000090046001600160601b0316818111156134cb576040516354ced18160e11b81526004810182905260248101839052604401610bed565b818110156129305760006134df8284615d33565b60035460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb90604401602060405180830381600087803b15801561352f57600080fd5b505af1158015613543573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135679190615471565b50604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b6002805414156135f05760405162461bcd60e51b815260206004820152601f6024820152600080516020615e8b8339815191526044820152606401610bed565b600280556000818152600660205260409020546001600160a01b031661362957604051630fb532db60e11b815260040160405180910390fd5b60008181526007602052604090208054600160601b90046001600160601b0316903490600c6136588385615cde565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600860148282829054906101000a90046001600160601b03166136a09190615cde565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f3f1ddc3ab1bdb39001ad76ca51a0e6f57ce6627c69f251d1de41622847721cde8234846136f39190615c9a565b60408051928352602083019190915201612c7b565b613710613719565b610bf681614188565b6000546001600160a01b031633146137735760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610bed565b565b60008061378184613e9e565b60035460405163a9059cbb60e01b81526001600160a01b0387811660048301526001600160601b0385166024830152939550919350919091169063a9059cbb90604401602060405180830381600087803b1580156137de57600080fd5b505af11580156137f2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138169190615471565b61383357604051631e9acf1760e31b815260040160405180910390fd5b6000836001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114613889576040519150601f19603f3d011682016040523d82523d6000602084013e61388e565b606091505b50509050806138b057604051630dcf35db60e41b815260040160405180910390fd5b604080516001600160a01b03861681526001600160601b038581166020830152841681830152905186917f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4919081900360600190a25050505050565b6040805160608101825260008082526020820181905281830181905282518084018452919290916139559186906002908390839080828437600092019190915250612e5f915050565b6000818152600b60205260409020549091506001600160a01b03168061399157604051631dfd6e1360e21b815260048101839052602401610bed565b6000828660c001356040516020016139b3929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600d90935291205490915080613a12576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81613a206020880188615736565b6020880135613a3560608a0160408b0161571b565b613a4560808b0160608c0161571b565b613a5560a08c0160808d016152f8565b613a6260a08d018d615bd5565b604051602001613a79989796959493929190615adc565b604051602081830303815290604052805190602001208114613ac7576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613ade613ad96020890189615736565b614232565b905080613bec576001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663e9413d38613b2160208a018a615736565b6040516001600160e01b031960e084901b16815267ffffffffffffffff909116600482015260240160206040518083038186803b158015613b6157600080fd5b505afa158015613b75573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b9991906154a7565b905080613bec57613bad6020880188615736565b6040517f175dadad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610bed565b6040805160c08a0135602080830191909152818301849052825180830384018152606090920190925280519101206000613c34613c2e368c90038c018c615564565b83614339565b90506040518060600160405280888152602001868152602001828152509750505050505050505b92915050565b60005a611388811015613c7357600080fd5b611388810390508460408204820311613c8b57600080fd5b50823b613c9757600080fd5b60008083516020850160008789f190505b9392505050565b60408051602081019091526000815281613cd85750604080516020810190915260008152613c5b565b7f92fd133800000000000000000000000000000000000000000000000000000000613d038385615d72565b6001600160e01b03191614613d44576040517f5247fdce00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b613d518260048186615c70565b810190613ca891906154c0565b60008115613d8c57601054613d859086908690640100000000900463ffffffff16866143a4565b9050613da6565b601054613da3908690869063ffffffff168661440e565b90505b949350505050565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401613de791511515815260200190565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff166001600160e01b03199093169290921790915292915050565b6000805b601154811015613e9557826001600160a01b031660118281548110613e5f57613e5f615e3b565b6000918252602090912001546001600160a01b03161415613e835750600192915050565b80613e8d81615da2565b915050613e38565b50600092915050565b6000818152600660209081526040808320815160608101835281546001600160a01b03908116825260018301541681850152600282018054845181870281018701865281815287968796949594860193919290830182828015613f2a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613f0c575b5050509190925250505060008581526007602090815260408083208151808301909252546001600160601b03808216808452600160601b9092041692820183905296509094509192505b826040015151811015613fec576005600084604001518381518110613f9b57613f9b615e3b565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208982529092529020805467ffffffffffffffff1916905580613fe481615da2565b915050613f74565b50600085815260066020526040812080546001600160a01b031990811682556001820180549091169055906140246002830182615188565b50506000858152600760205260409020805477ffffffffffffffffffffffffffffffffffffffffffffffff1916905560088054859190819061407c9084906801000000000000000090046001600160601b0316615d4a565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555082600860148282829054906101000a90046001600160601b03166140c49190615d4a565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050915091565b60004661a4b1811480614104575062066eed81145b156141815760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561414357600080fd5b505afa158015614157573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061417b91906154a7565b91505090565b4391505090565b6001600160a01b0381163314156141e15760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610bed565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b1811480614247575062066eed81145b15614329576101008367ffffffffffffffff166142626140ef565b61426c9190615d33565b1180614289575061427b6140ef565b8367ffffffffffffffff1610155b156142975750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a829060240160206040518083038186803b1580156142f157600080fd5b505afa158015614305573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ca891906154a7565b505067ffffffffffffffff164090565b600061436d8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151614515565b600383602001516040516020016143859291906159f1565b60408051601f1981840301815291905280516020909101209392505050565b6000806143af614750565b905060005a6143be8888615c9a565b6143c89190615d33565b6143d29085615d14565b905060006143eb63ffffffff871664e8d4a51000615d14565b9050826143f88284615c9a565b6144029190615c9a565b98975050505050505050565b6000806144196147ac565b90506000811361443f576040516321ea67b360e11b815260048101829052602401610bed565b6000614449614750565b9050600082825a61445a8b8b615c9a565b6144649190615d33565b61446e9088615d14565b6144789190615c9a565b61448a90670de0b6b3a7640000615d14565b6144949190615d00565b905060006144ad63ffffffff881664e8d4a51000615d14565b90506144c5816b033b2e3c9fd0803ce8000000615d33565b8211156144fe576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6145088183615c9a565b9998505050505050505050565b61451e89614893565b61456a5760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610bed565b61457388614893565b6145bf5760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610bed565b6145c883614893565b6146145760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610bed565b61461d82614893565b6146695760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610bed565b614675878a888761496c565b6146c15760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610bed565b60006146cd8a87614aa1565b905060006146e0898b878b868989614b05565b905060006146f1838d8d8a86614c25565b9050808a146147425760405162461bcd60e51b815260206004820152600d60248201527f696e76616c69642070726f6f66000000000000000000000000000000000000006044820152606401610bed565b505050505050505050505050565b60004661a4b1811480614765575062066eed81145b156147a457606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b15801561414357600080fd5b600091505090565b600f5460048054604080517ffeaf968c00000000000000000000000000000000000000000000000000000000815290516000946601000000000000900463ffffffff169384151593869384936001600160a01b039092169263feaf968c928282019260a09290829003018186803b15801561482657600080fd5b505afa15801561483a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061485e9190615760565b50945090925084915050801561488257506148798242615d33565b8463ffffffff16105b15613da65750600e54949350505050565b80516000906401000003d019116148ec5760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610bed565b60208201516401000003d019116149455760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610bed565b60208201516401000003d0199080096149658360005b6020020151614c65565b1492915050565b60006001600160a01b0382166149c45760405162461bcd60e51b815260206004820152600b60248201527f626164207769746e6573730000000000000000000000000000000000000000006044820152606401610bed565b6020840151600090600116156149db57601c6149de565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa158015614a79573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b614aa96151a6565b614ad660018484604051602001614ac2939291906158db565b604051602081830303815290604052614c89565b90505b614ae281614893565b613c5b578051604080516020810192909252614afe9101614ac2565b9050614ad9565b614b0d6151a6565b825186516401000003d0199081900691061415614b6c5760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610bed565b614b77878988614cd8565b614bc35760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610bed565b614bce848685614cd8565b614c1a5760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610bed565b614402868484614e12565b600060028686868587604051602001614c439695949392919061587c565b60408051601f1981840301815291905280516020909101209695505050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614c916151a6565b614c9a82614ed9565b8152614caf614caa82600061495b565b614f14565b602082018190526002900660011415614cd3576020810180516401000003d0190390525b919050565b600082614d275760405162461bcd60e51b815260206004820152600b60248201527f7a65726f207363616c61720000000000000000000000000000000000000000006044820152606401610bed565b83516020850151600090614d3d90600290615de5565b15614d4957601c614d4c565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614dbe573d6000803e3d6000fd5b505050602060405103519050600086604051602001614ddd919061586a565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614e1a6151a6565b835160208086015185519186015160009384938493614e3b93909190614f34565b919450925090506401000003d019858209600114614e9b5760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610bed565b60405180604001604052806401000003d01980614eba57614eba615e0f565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d0198110614cd357604080516020808201939093528151808203840181529082019091528051910120614ee1565b6000613c5b826002614f2d6401000003d0196001615c9a565b901c615014565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614f74838385856150b6565b9098509050614f8588828e886150da565b9098509050614f9688828c876150da565b90985090506000614fa98d878b856150da565b9098509050614fba888286866150b6565b9098509050614fcb88828e896150da565b9098509050818114615000576401000003d019818a0998506401000003d01982890997506401000003d0198183099650615004565b8196505b5050505050509450945094915050565b60008061501f6151c4565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a08201526150516151e2565b60208160c0846005600019fa9250826150ac5760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610bed565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215615178579160200282015b8281111561517857825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190615143565b50615184929150615200565b5090565b5080546000825590600052602060002090810190610bf69190615200565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b808211156151845760008155600101615201565b8035614cd381615e67565b8060408101831015613c5b57600080fd5b600082601f83011261524257600080fd5b61524a615c23565b80838560408601111561525c57600080fd5b60005b600281101561527e57813584526020938401939091019060010161525f565b509095945050505050565b600060c0828403121561529b57600080fd5b50919050565b803561ffff81168114614cd357600080fd5b803563ffffffff81168114614cd357600080fd5b805169ffffffffffffffffffff81168114614cd357600080fd5b80356001600160601b0381168114614cd357600080fd5b60006020828403121561530a57600080fd5b8135613ca881615e67565b6000806040838503121561532857600080fd5b823561533381615e67565b9150615341602084016152e1565b90509250929050565b6000806040838503121561535d57600080fd5b823561536881615e67565b9150602083013561537881615e67565b809150509250929050565b6000806060838503121561539657600080fd5b82356153a181615e67565b91506153418460208501615220565b600080600080606085870312156153c657600080fd5b84356153d181615e67565b935060208501359250604085013567ffffffffffffffff808211156153f557600080fd5b818701915087601f83011261540957600080fd5b81358181111561541857600080fd5b88602082850101111561542a57600080fd5b95989497505060200194505050565b60006040828403121561544b57600080fd5b613ca88383615220565b60006040828403121561546757600080fd5b613ca88383615231565b60006020828403121561548357600080fd5b8151613ca881615e7c565b6000602082840312156154a057600080fd5b5035919050565b6000602082840312156154b957600080fd5b5051919050565b6000602082840312156154d257600080fd5b6040516020810181811067ffffffffffffffff821117156154f5576154f5615e51565b604052823561550381615e7c565b81529392505050565b6000808284036101c081121561552157600080fd5b6101a08082121561553157600080fd5b849350830135905067ffffffffffffffff81111561554e57600080fd5b61555a85828601615289565b9150509250929050565b60006101a0828403121561557757600080fd5b61557f615c4c565b6155898484615231565b81526155988460408501615231565b60208201526080830135604082015260a0830135606082015260c083013560808201526155c760e08401615215565b60a08201526101006155db85828601615231565b60c08301526155ee856101408601615231565b60e083015261018084013581830152508091505092915050565b60006020828403121561561a57600080fd5b813567ffffffffffffffff81111561563157600080fd5b613da684828501615289565b60006020828403121561564f57600080fd5b613ca8826152a1565b60008060008060008086880360e081121561567257600080fd5b61567b886152a1565b9650615689602089016152b3565b9550615697604089016152b3565b94506156a5606089016152b3565b9350608088013592506040609f19820112156156c057600080fd5b506156c9615c23565b6156d560a089016152b3565b81526156e360c089016152b3565b6020820152809150509295509295509295565b6000806040838503121561570957600080fd5b82359150602083013561537881615e67565b60006020828403121561572d57600080fd5b613ca8826152b3565b60006020828403121561574857600080fd5b813567ffffffffffffffff81168114613ca857600080fd5b600080600080600060a0868803121561577857600080fd5b615781866152c7565b94506020860151935060408601519250606086015191506157a4608087016152c7565b90509295509295909350565b600081518084526020808501945080840160005b838110156157e95781516001600160a01b0316875295820195908201906001016157c4565b509495945050505050565b8060005b60028110156158175781518452602093840193909101906001016157f8565b50505050565b6000815180845260005b8181101561584357602081850181015186830182015201615827565b81811115615855576000602083870101525b50601f01601f19169290920160200192915050565b61587481836157f4565b604001919050565b86815261588c60208201876157f4565b61589960608201866157f4565b6158a660a08201856157f4565b6158b360e08201846157f4565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b8381526158eb60208201846157f4565b606081019190915260800192915050565b60408101613c5b82846157f4565b602081526000613ca8602083018461581d565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c0608084015261596360e08401826157b0565b905060808401516001600160601b0380821660a08601528060a08701511660c086015250508091505092915050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b818110156159e3578451835293830193918301916001016159c7565b509098975050505050505050565b82815260608101613ca860208301846157f4565b6000604082018483526020604081850152818551808452606086019150828701935060005b81811015615a4657845183529383019391830191600101615a2a565b5090979650505050505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a083015261440260c083018461581d565b878152866020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c083015261450860e083018461581d565b88815267ffffffffffffffff88166020820152866040820152600063ffffffff80881660608401528087166080840152506001600160a01b03851660a083015260e060c08301528260e08301526101008385828501376000838501820152601f909301601f191690910190910198975050505050505050565b8481526001600160601b0384166020820152608060408201526000615b7d608083018561581d565b9050821515606083015295945050505050565b60006001600160601b0380871683528086166020840152506001600160a01b038416604083015260806060830152615bcb60808301846157b0565b9695505050505050565b6000808335601e19843603018112615bec57600080fd5b83018035915067ffffffffffffffff821115615c0757600080fd5b602001915036819003821315615c1c57600080fd5b9250929050565b6040805190810167ffffffffffffffff81118282101715615c4657615c46615e51565b60405290565b604051610120810167ffffffffffffffff81118282101715615c4657615c46615e51565b60008085851115615c8057600080fd5b83861115615c8d57600080fd5b5050820193919092039150565b60008219821115615cad57615cad615df9565b500190565b600067ffffffffffffffff808316818516808303821115615cd557615cd5615df9565b01949350505050565b60006001600160601b03808316818516808303821115615cd557615cd5615df9565b600082615d0f57615d0f615e0f565b500490565b6000816000190483118215151615615d2e57615d2e615df9565b500290565b600082821015615d4557615d45615df9565b500390565b60006001600160601b0383811690831681811015615d6a57615d6a615df9565b039392505050565b6001600160e01b03198135818116916004851015615d9a5780818660040360031b1b83161692505b505092915050565b6000600019821415615db657615db6615df9565b5060010190565b600067ffffffffffffffff80831681811415615ddb57615ddb615df9565b6001019392505050565b600082615df457615df4615e0f565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114610bf657600080fd5b8015158114610bf657600080fdfe5265656e7472616e637947756172643a207265656e7472616e742063616c6c00a164736f6c6343000806000a",
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
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)

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
	outstruct.StalenessSeconds = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[3], new(uint32)).(*uint32)

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

	outstruct.FulfillmentFlatFeeLinkPPM = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeEthPPM = *abi.ConvertType(out[1], new(uint32)).(*uint32)

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

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetConfig(&_VRFCoordinatorV2Plus.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetConfig(&_VRFCoordinatorV2Plus.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactor) SetLINKAndLINKETHFeed(opts *bind.TransactOpts, link common.Address, linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.contract.Transact(opts, "setLINKAndLINKETHFeed", link, linkEthFeed)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusSession) SetLINKAndLINKETHFeed(link common.Address, linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetLINKAndLINKETHFeed(&_VRFCoordinatorV2Plus.TransactOpts, link, linkEthFeed)
}

func (_VRFCoordinatorV2Plus *VRFCoordinatorV2PlusTransactorSession) SetLINKAndLINKETHFeed(link common.Address, linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2Plus.Contract.SetLINKAndLINKETHFeed(&_VRFCoordinatorV2Plus.TransactOpts, link, linkEthFeed)
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
	ExtraArgs  []byte
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
	Owner      common.Address
	Consumers  []common.Address
}
type SConfig struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
}
type SFeeConfig struct {
	FulfillmentFlatFeeLinkPPM uint32
	FulfillmentFlatFeeEthPPM  uint32
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
	return common.HexToHash("0x777357bb93f63d088f18112d3dba38457aec633eb8f1341e1d418380ad328e78")
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
	return common.HexToHash("0x7bad4a5a5566568e03ff9e9fdaedf6b8a760aa3ffc0abad348cf93787f04e17e")
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

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

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

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusFeeConfig) (*types.Transaction, error)

	SetLINKAndLINKETHFeed(opts *bind.TransactOpts, link common.Address, linkEthFeed common.Address) (*types.Transaction, error)

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
