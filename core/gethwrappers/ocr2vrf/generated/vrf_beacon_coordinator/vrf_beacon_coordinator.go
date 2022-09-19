// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_beacon_coordinator

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

type ECCArithmeticG1Point struct {
	P [2]*big.Int
}

type KeyDataStructKeyData struct {
	PublicKey []byte
	Hashes    [][32]byte
}

type VRFBeaconReportOutputServed struct {
	Height            uint64
	ConfirmationDelay *big.Int
}

type VRFBeaconReportReport struct {
	Outputs           []VRFBeaconReportVRFOutput
	JuelsPerFeeCoin   *big.Int
	RecentBlockHeight uint64
	RecentBlockHash   [32]byte
}

type VRFBeaconReportVRFOutput struct {
	BlockHeight       uint64
	ConfirmationDelay *big.Int
	VrfOutput         ECCArithmeticG1Point
	Callbacks         []VRFBeaconTypesCostedCallback
}

type VRFBeaconTypesCallback struct {
	RequestID    *big.Int
	NumWords     uint16
	Requester    common.Address
	Arguments    []byte
	SubID        uint64
	GasAllowance *big.Int
}

type VRFBeaconTypesCostedCallback struct {
	Callback VRFBeaconTypesCallback
	Price    *big.Int
}

var VRFBeaconCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocksArg\",\"type\":\"uint256\"},{\"internalType\":\"contractDKG\",\"name\":\"keyProvider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"keyID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BeaconPeriodMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earliestAllowed\",\"type\":\"uint256\"}],\"name\":\"BlockTooRecent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"firstDelay\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"minDelay\",\"type\":\"uint16\"}],\"name\":\"ConfirmationDelayBlocksTooShort\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confirmationDelays\",\"type\":\"uint16[10]\"},{\"internalType\":\"uint8\",\"name\":\"violatingIndex\",\"type\":\"uint8\"}],\"name\":\"ConfirmationDelaysNotIncreasing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"separatorHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorTooOld\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"providedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"onchainHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorWrong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"keyProvider\",\"type\":\"address\"}],\"name\":\"KeyInfoMustComeFromProvider\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoWordsRequested\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confDelays\",\"type\":\"uint16[10]\"}],\"name\":\"NonZeroDelayAfterZeroDelay\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"expectedLength\",\"type\":\"uint256\"}],\"name\":\"OffchainConfigHasWrongLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"}],\"name\":\"RandomnessNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"ResponseMustBeRetrievedByRequester\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequestsReplaceContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySlotsReplaceContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"TooManyWords\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"UniverseHasEndedBangBangBang\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"occVersion\",\"type\":\"uint64\"}],\"name\":\"UnknownConfigVersion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"givenDelay\",\"type\":\"uint24\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"knownDelays\",\"type\":\"uint24[8]\"}],\"name\":\"UnknownConfirmationDelay\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconReport.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconReport.VRFOutput[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"recentBlockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structVRFBeaconReport.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"exposeType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"forgetConsumerSubscriptionID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfirmationDelays\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"\",\"type\":\"uint24[8]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_StartSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structKeyDataStruct.KeyData\",\"name\":\"kd\",\"type\":\"tuple\"}],\"name\":\"keyGenerated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxErrorMsgLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDelay\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newKeyRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keyID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_provingKeyHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200583c3803806200583c83398101604081905262000034916200022f565b8181848681818181803380600081620000945760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c757620000c7816200016b565b5050506001600160a01b03166080526000829003620000f957604051632abc297960e01b815260040160405180910390fd5b60a082905260006200010c83436200027d565b905060008160a051620001209190620002b6565b90506200012e8143620002d0565b60c0525050601d80546001600160a01b0319166001600160a01b039990991698909817909755505050601e9290925550620002eb95505050505050565b336001600160a01b03821603620001c55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200022c57600080fd5b50565b600080600080608085870312156200024657600080fd5b8451620002538162000216565b6020860151604087015191955093506200026d8162000216565b6060959095015193969295505050565b6000826200029b57634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b600082821015620002cb57620002cb620002a0565b500390565b60008219821115620002e657620002e6620002a0565b500190565b60805160a05160c0516154d462000368600039600061052b01526000818161050401528181610a070152818161332901528181613358015281816133900152613c5201526000818161028e0152818161158c0152818161164a0152818161173601528181612321015281816129170152612a3401526154d46000f3fe608060405234801561001057600080fd5b50600436106102275760003560e01c8063b1dc65a411610130578063cf7e754a116100b8578063e4902f821161007c578063e4902f821461059b578063eb5dcd6c146105c3578063f2fde38b146105d6578063f645dcb1146105e9578063fbffd2c1146105fc57600080fd5b8063cf7e754a14610526578063d09dc3391461054d578063d57fc45a14610555578063dc92accf1461055e578063e3d0e7121461058857600080fd5b8063c278e5b7116100ff578063c278e5b7146104b9578063c4c92b37146104ca578063c63c4e9b146104db578063cc31f7dd146104f6578063cd0593df146104ff57600080fd5b8063b1dc65a414610477578063bbcdd0d81461048a578063bf2732c714610493578063c1075329146104a657600080fd5b80637a464944116101b35780638da5cb5b116101825780638da5cb5b146104035780639c849b30146104145780639e3616f414610427578063afcb95d71461043a578063b121e1471461046457600080fd5b80637a464944146103a657806381ff7048146103ae57806385c64e11146103db5780638ac28d5a146103f057600080fd5b80632f7527cc116101fa5780632f7527cc1461034757806355e4874914610361578063643dc1051461036b57806374d846111461037e57806379ba50971461039e57600080fd5b80630eafb25b1461022c578063181f5a77146102525780631b6b6d231461028957806329937268146102c8575b600080fd5b61023f61023a3660046140c9565b61060f565b6040519081526020015b60405180910390f35b6040805180820182526015815274565246426561636f6e20312e302e302d616c70686160581b60208201529051610249919061413e565b6102b07f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610249565b61030b600c54600160501b810463ffffffff90811692600160701b8304821692600160901b8104831692600160b01b82041691600160d01b90910462ffffff1690565b6040805163ffffffff9687168152948616602086015292851692840192909252909216606082015262ffffff909116608082015260a001610249565b61034f600881565b60405160ff9091168152602001610249565b610369610714565b005b61036961037936600461417b565b61075e565b61039161038c3660046141f8565b610944565b6040516102499190614250565b610369610ada565b61023f608081565b600d54600f54604080516000815264010000000090930463ffffffff166020840152820152606001610249565b6103e3610b84565b604051610249919061428b565b6103696103fe3660046140c9565b610be9565b6000546001600160a01b03166102b0565b6103696104223660046142e5565b610c5b565b610369610435366004614350565b610e2d565b600f546011546040805160008152602081019390935263ffffffff90911690820152606001610249565b6103696104723660046140c9565b610ebc565b6103696104853660046143d2565b610f98565b61023f6103e881565b6103696104a13660046145fc565b61141b565b6103696104b43660046146e4565b611489565b6103696104c7366004614710565b50565b601c546001600160a01b03166102b0565b6104e3600381565b60405161ffff9091168152602001610249565b61023f601e5481565b61023f7f000000000000000000000000000000000000000000000000000000000000000081565b61023f7f000000000000000000000000000000000000000000000000000000000000000081565b61023f611714565b61023f601f5481565b61057161056c36600461477a565b6117c0565b60405165ffffffffffff9091168152602001610249565b6103696105963660046147d6565b6118dd565b6105ae6105a93660046140c9565b612006565b60405163ffffffff9091168152602001610249565b6103696105d13660046148c3565b6120b5565b6103696105e43660046140c9565b6121ed565b6105716105f73660046148fc565b6121fe565b61036961060a3660046140c9565b6122fe565b6001600160a01b03811660009081526012602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b031691810191909152906106715750600092915050565b600c546020820151600091600160901b900463ffffffff169060169060ff16601f81106106a0576106a0614981565b600881049190910154600c546106d3926007166004026101000a90910463ffffffff90811691600160301b9004166149ad565b63ffffffff166106e391906149d2565b6106f190633b9aca006149d2565b905081604001516001600160601b03168161070c91906149f1565b949350505050565b601d546001600160a01b03163381146107565760405163292f4fb560e01b81523360048201526001600160a01b03821660248201526044015b60405180910390fd5b506000601f55565b601c546001600160a01b031661077c6000546001600160a01b031690565b6001600160a01b0316336001600160a01b031614806108085750604051630d629b5f60e31b81526001600160a01b03821690636b14daf8906107c79033906000903690600401614a32565b602060405180830381865afa1580156107e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108089190614a57565b6108545760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015260640161074d565b61085c61230f565b600c805467ffffffffffffffff60501b1916600160501b63ffffffff89811691820263ffffffff60701b191692909217600160701b8984169081029190911767ffffffffffffffff60901b1916600160901b89851690810263ffffffff60b01b191691909117600160b01b9489169485021762ffffff60d01b1916600160d01b62ffffff89169081029190911790955560408051938452602084019290925290820152606081019190915260808101919091527f0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f9060a00160405180910390a1505050505050565b65ffffffffffff81166000818152600a602081815260408084208151608081018352815463ffffffff8116825262ffffff6401000000008204168286015261ffff600160381b820416938201939093526001600160a01b03600160481b84048116606083810191825298909752949093526001600160e81b0319909116905591511633146109fc576060810151604051638e30e82360e01b81526001600160a01b03909116600482015233602482015260440161074d565b8051600090610a32907f00000000000000000000000000000000000000000000000000000000000000009063ffffffff166149d2565b90506000826020015162ffffff1643610a4b9190614a79565b9050808210610a76576040516315ad27c360e01b81526004810183905243602482015260440161074d565b6001600160401b03821115610aa1576040516302c6ef8160e11b81526004810183905260240161074d565b60008281526007602090815260408083208287015162ffffff168452909152902054610ad1908690859085612697565b95945050505050565b6001546001600160a01b03163314610b2d5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015260640161074d565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b8c613edf565b6040805161010081019182905290600990600890826000855b82829054906101000a900462ffffff1662ffffff1681526020019060030190602082600201049283019260010382029150808411610ba55790505050505050905090565b6001600160a01b038181166000908152601a6020526040902054163314610c525760405162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015260640161074d565b6104c78161286b565b610c63612a8b565b828114610cb25760405162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015260640161074d565b60005b83811015610e26576000858583818110610cd157610cd1614981565b9050602002016020810190610ce691906140c9565b90506000848484818110610cfc57610cfc614981565b9050602002016020810190610d1191906140c9565b6001600160a01b038084166000908152601a60205260409020549192501680158080610d4e5750826001600160a01b0316826001600160a01b0316145b610d8e5760405162461bcd60e51b81526020600482015260116024820152701c185e595948185b1c9958591e481cd95d607a1b604482015260640161074d565b6001600160a01b038481166000908152601a6020526040902080546001600160a01b03191685831690811790915590831614610e0f57826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505050508080610e1e90614a90565b915050610cb5565b5050505050565b610e35612a8b565b60005b81811015610eb757600060056000858585818110610e5857610e58614981565b9050602002016020810190610e6d91906140c9565b6001600160a01b031681526020810191909152604001600020805467ffffffffffffffff19166001600160401b039290921691909117905580610eaf81614a90565b915050610e38565b505050565b6001600160a01b038181166000908152601b6020526040902054163314610f255760405162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015260640161074d565b6001600160a01b038181166000818152601a602090815260408083208054336001600160a01b03198083168217909355601b909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a604080516101008082018352600c5460ff808216845291810464ffffffffff16602080850191909152600160301b820463ffffffff90811685870152600160501b830481166060860152600160701b830481166080860152600160901b8304811660a0860152600160b01b83041660c0850152600160d01b90910462ffffff1660e08401523360009081526012825293909320549394509092918c013591166110865760405162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015260640161074d565b600f548b35146110d05760405162461bcd60e51b81526020600482015260156024820152740c6dedcccd2ce88d2cecae6e840dad2e6dac2e8c6d605b1b604482015260640161074d565b6110de8a8a8a8a8a8a612ae0565b81516110eb906001614aa9565b60ff16871461113c5760405162461bcd60e51b815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e617475726573000000000000604482015260640161074d565b86851461118b5760405162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015260640161074d565b60008a8a60405161119d929190614ace565b6040519081900381206111b4918e90602001614ade565b60408051601f19818403018152828252805160209182012083830190925260008084529083018190529092509060005b8a81101561134c5760006001858a846020811061120357611203614981565b61121091901a601b614aa9565b8f8f8681811061122257611222614981565b905060200201358e8e8781811061123b5761123b614981565b9050602002013560405160008152602001604052604051611278949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561129a573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526013602090815290849020838501909452925460ff80821615158085526101009092041693830193909352909550925090506113255760405162461bcd60e51b815260206004820152600f60248201526e39b4b3b730ba3ab9329032b93937b960891b604482015260640161074d565b826020015160080260ff166001901b8401935050808061134490614a90565b9150506111e4565b5081827e0101010101010101010101010101010101010101010101010101010101010116146113b05760405162461bcd60e51b815260206004820152601060248201526f323ab83634b1b0ba329039b4b3b732b960811b604482015260640161074d565b50600091506113ff9050838d836020020135848e8e8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612b7d92505050565b905061140d83828633612fb0565b505050505050505050505050565b601d546001600160a01b03163381146114585760405163292f4fb560e01b81523360048201526001600160a01b038216602482015260440161074d565b815160405161146a9190602001614afa565b60408051601f198184030181529190528051602090910120601f555050565b6000546001600160a01b03163314806115135750601c54604051630d629b5f60e31b81526001600160a01b0390911690636b14daf8906114d29033906000903690600401614a32565b602060405180830381865afa1580156114ef573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115139190614a57565b61155f5760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015260640161074d565b60006115696130bf565b6040516370a0823160e01b81523060048201529091506000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa1580156115d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115f79190614b16565b9050818110156116405760405162461bcd60e51b8152602060048201526014602482015273696e73756666696369656e742062616c616e636560601b604482015260640161074d565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663a9059cbb8561168361167d8686614a79565b87613289565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af11580156116ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116f29190614a57565b61170e5760405162461bcd60e51b815260040161074d90614b2f565b50505050565b6040516370a0823160e01b815230600482015260009081906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa15801561177d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117a19190614b16565b905060006117ad6130bf565b90506117b98183614b5b565b9250505090565b6000806000806117d087866132a3565b92509250925065ffffffffffff83166000908152600a602090815260409182902084518154928601518487015160608801516001600160a01b0316600160481b027fffffff0000000000000000000000000000000000000000ffffffffffffffffff61ffff909216600160381b0291909116670100000000000000600160e81b031962ffffff9093166401000000000266ffffffffffffff1990961663ffffffff90941693909317949094171617919091179055516001600160401b038216907fc334d6f57be304c8192da2e39220c48e35f7e9afa16c541e68a6a859eff4dbc5906118c890889062ffffff91909116815260200190565b60405180910390a250909150505b9392505050565b6118e5612a8b565b601f8911156119295760405162461bcd60e51b815260206004820152601060248201526f746f6f206d616e79206f7261636c657360801b604482015260640161074d565b8887146119715760405162461bcd60e51b81526020600482015260166024820152750dee4c2c6d8ca40d8cadccee8d040dad2e6dac2e8c6d60531b604482015260640161074d565b8861197d876003614b9a565b60ff16106119cd5760405162461bcd60e51b815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f20686967680000000000000000604482015260640161074d565b6119d98660ff1661355e565b6040805160e060208c02808301820190935260c082018c815260009383928f918f918291908601908490808284376000920191909152505050908252506040805160208c810282810182019093528c82529283019290918d918d91829185019084908082843760009201919091525050509082525060ff891660208083019190915260408051601f8a0183900483028101830182528981529201919089908990819084018382808284376000920191909152505050908252506001600160401b03861660208083019190915260408051601f870183900483028101830182528681529201919086908690819084018382808284376000920191909152505050915250600c805465ffffffffff00191690559050611af461230f565b60145460005b81811015611ba557600060148281548110611b1757611b17614981565b6000918252602082200154601580546001600160a01b0390921693509084908110611b4457611b44614981565b60009182526020808320909101546001600160a01b039485168352601382526040808420805461ffff1916905594168252601290529190912080546dffffffffffffffffffffffffffff191690555080611b9d81614a90565b915050611afa565b50611bb260146000613efe565b611bbe60156000613efe565b60005b825151811015611e37576013600084600001518381518110611be557611be5614981565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff1615611c595760405162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015260640161074d565b604080518082019091526001815260ff821660208201528351805160139160009185908110611c8a57611c8a614981565b6020908102919091018101516001600160a01b03168252818101929092526040016000908120835181549484015161ffff1990951690151561ff0019161761010060ff90951694909402939093179092558401518051601292919084908110611cf557611cf5614981565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff1615611d695760405162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015260640161074d565b60405180606001604052806001151581526020018260ff16815260200160006001600160601b03168152506012600085602001518481518110611dae57611dae614981565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020835181549385015194909201516001600160601b0316620100000262010000600160701b031960ff959095166101000261ff00199315159390931661ffff1990941693909317919091179290921617905580611e2f81614a90565b915050611bc1565b5081518051611e4e91601491602090910190613f1c565b506020808301518051611e65926015920190613f1c565b506040820151600c805460ff191660ff909216919091179055600d805467ffffffff0000000019811664010000000063ffffffff43811682029283179094558204831692600092611ebd929082169116176001614bc3565b905080600d60006101000a81548163ffffffff021916908363ffffffff1602179055506000611f1146308463ffffffff16886000015189602001518a604001518b606001518c608001518d60a001516135a3565b905080600f600001819055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05838284886000015189602001518a604001518b606001518c608001518d60a00151604051611f7499989796959493929190614c24565b60405180910390a1600c54600160301b900463ffffffff1660005b865151811015611fe95781601682601f8110611fad57611fad614981565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff1602179055508080611fe190614a90565b915050611f8f565b50611ff48b8b6135fe565b50505050505050505050505050505050565b6001600160a01b03811660009081526012602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b031691810191909152906120685750600092915050565b6016816020015160ff16601f811061208257612082614981565b600881049190910154600c546118d6926007166004026101000a90910463ffffffff90811691600160301b9004166149ad565b6001600160a01b038281166000908152601a602052604090205416331461211e5760405162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015260640161074d565b6001600160a01b03811633036121765760405162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161074d565b6001600160a01b038083166000908152601b6020526040902080548383166001600160a01b031982168117909255909116908114610eb7576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b6121f5612a8b565b6104c78161360c565b600080600061220d87876132a3565b925050915060006040518060c001604052808465ffffffffffff1681526020018961ffff168152602001336001600160a01b031681526020018681526020018a6001600160401b031681526020018763ffffffff166001600160601b0316815250905081878a836040516020016122879493929190614cb9565b60408051601f19818403018152828252805160209182012065ffffffffffff871660009081526006909252919020557fa62e84e206cb87e2f6896795353c5358ff3d415d0bccc24e45c5fad83e17d03c906122e99084908a908d908690614cb9565b60405180910390a15090979650505050505050565b612306612a8b565b6104c7816136b5565b600c54604080516103e08101918290527f000000000000000000000000000000000000000000000000000000000000000092600160301b900463ffffffff169160009190601690601f908285855b82829054906101000a900463ffffffff1663ffffffff168152602001906004019060208260030104928301926001038202915080841161235d57905050505050509050600060158054806020026020016040519081016040528092919081815260200182805480156123f857602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116123da575b5050505050905060005b81518110156126895760006012600084848151811061242357612423614981565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160029054906101000a90046001600160601b03166001600160601b0316905060006012600085858151811061248557612485614981565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160026101000a8154816001600160601b0302191690836001600160601b0316021790555060008483601f81106124e8576124e8614981565b6020020151600c5490870363ffffffff9081169250600160901b909104168102633b9aca00028201801561267e576000601a600087878151811061252e5761252e614981565b6020908102919091018101516001600160a01b03908116835290820192909252604090810160002054905163a9059cbb60e01b815290821660048201819052602482018590529250908a169063a9059cbb906044016020604051808303816000875af11580156125a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125c69190614a57565b6125e25760405162461bcd60e51b815260040161074d90614b2f565b878786601f81106125f5576125f5614981565b602002019063ffffffff16908163ffffffff1681525050886001600160a01b0316816001600160a01b031687878151811061263257612632614981565b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c8560405161267491815260200190565b60405180910390a4505b505050600101612402565b50610e26601683601f613f81565b6060826126d05760405163c7d41b1b60e01b815265ffffffffffff861660048201526001600160401b038316602482015260440161074d565b6040805165ffffffffffff8716602080830191909152865163ffffffff168284015286015162ffffff166060808301919091529186015161ffff166080820152908501516001600160a01b031660a082015260c0810184905260009060e0016040516020818303038152906040528051906020012090506103e8856040015161ffff161115612786576040808601519051634a90778560e01b815261ffff90911660048201526103e8602482015260440161074d565b6000856040015161ffff166001600160401b038111156127a8576127a8614488565b6040519080825280602002602001820160405280156127d1578160200160208202803683370190505b50905060005b866040015161ffff168161ffff16101561286057828160405160200161281492919091825260f01b6001600160f01b031916602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff168151811061284357612843614981565b60209081029190910101528061285881614d5b565b9150506127d7565b509695505050505050565b6001600160a01b0381166000908152601260209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046001600160601b0316928101929092526128c8575050565b60006128d38361060f565b90508015610eb7576001600160a01b038381166000908152601a60205260409081902054905163a9059cbb60e01b81529082166004820181905260248201849052917f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015612960573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129849190614a57565b6129a05760405162461bcd60e51b815260040161074d90614b2f565b600c60000160069054906101000a900463ffffffff166016846020015160ff16601f81106129d0576129d0614981565b6008810491909101805460079092166004026101000a63ffffffff8181021990931693909216919091029190911790556001600160a01b03848116600081815260126020908152604091829020805462010000600160701b031916905590518581527f0000000000000000000000000000000000000000000000000000000000000000841693851692917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c910160405180910390a450505050565b6000546001600160a01b03163314612ade5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015260640161074d565b565b6000612aed8260206149d2565b612af88560206149d2565b612b04886101446149f1565b612b0e91906149f1565b612b1891906149f1565b612b239060006149f1565b9050368114612b745760405162461bcd60e51b815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d617463680000000000000000604482015260640161074d565b50505050505050565b60008082806020019051810190612b949190614f4d565b64ffffffffff85166020880152604087018051919250612bb382615121565b63ffffffff1663ffffffff168152505085600c60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600e6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160000160126101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160000160166101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600001601a6101000a81548162ffffff021916908362ffffff160217905550905050600081604001516001600160401b031640905080826060015114612d55576060820151604080840151905163aed0afe560e01b81526004810192909252602482018390526001600160401b0316604482015260640161074d565b6000808360000151516001600160401b03811115612d7557612d75614488565b604051908082528060200260200182016040528015612dba57816020015b6040805180820190915260008082526020820152815260200190600190039081612d935790505b50905060005b845151811015612e8a57600085600001518281518110612de257612de2614981565b60200260200101519050612dff818760400151886020015161372b565b60408101515151151580612e1b57506040810151516020015115155b15612e7757604051806040016040528082600001516001600160401b03168152602001826020015162ffffff16815250838381518110612e5d57612e5d614981565b60200260200101819052508380612e7390614d5b565b9450505b5080612e8281614a90565b915050612dc0565b5060008261ffff166001600160401b03811115612ea957612ea9614488565b604051908082528060200260200182016040528015612eee57816020015b6040805180820190915260008082526020820152815260200190600190039081612ec75790505b50905060005b8361ffff16811015612f4a57828181518110612f1257612f12614981565b6020026020010151828281518110612f2c57612f2c614981565b60200260200101819052508080612f4290614a90565b915050612ef4565b508764ffffffffff168a6040015163ffffffff167fe0c90b8e55243fcba0f8b68b201983b97f7a3d5aebd6dfa1a4082a07925cc7443388602001518d86604051612f97949392919061513a565b60405180910390a3505050506020015195945050505050565b6000612fd7633b9aca003a04866080015163ffffffff16876060015163ffffffff16613b04565b90506010360260005a905060006130008663ffffffff1685858b60e0015162ffffff1686613b21565b90506000670de0b6b3a76400006001600160c01b03891683026001600160a01b03881660009081526012602052604090205460c08c01519290910492506201000090046001600160601b039081169163ffffffff16633b9aca000282840101908116821115613075575050505050505061170e565b6001600160a01b038816600090815260126020526040902080546001600160601b03909216620100000262010000600160701b031990921691909117905550505050505050505050565b600080601580548060200260200160405190810160405280929190818152602001828054801561311857602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116130fa575b50508351600c54604080516103e08101918290529697509195600160301b90910463ffffffff169450600093509150601690601f908285855b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116131515790505050505050905060005b838110156131e4578181601f81106131b1576131b1614981565b60200201516131c090846149ad565b6131d09063ffffffff16876149f1565b9550806131dc81614a90565b915050613197565b50600c5461320390600160901b900463ffffffff16633b9aca006149d2565b61320d90866149d2565b945060005b83811015613281576012600086838151811061323057613230614981565b6020908102919091018101516001600160a01b031682528101919091526040016000205461326d906201000090046001600160601b0316876149f1565b95508061327981614a90565b915050613212565b505050505090565b60008183101561329a57508161329d565b50805b92915050565b604080516080810182526000808252602082018190529181018290526060810182905260006103e88561ffff1611156132fd57604051634a90778560e01b815261ffff861660048201526103e8602482015260440161074d565b8461ffff16600003613322576040516308fad2a760e01b815260040160405180910390fd5b600061334e7f0000000000000000000000000000000000000000000000000000000000000000436151db565b905060008161337d7f0000000000000000000000000000000000000000000000000000000000000000436149f1565b6133879190614a79565b905060006133b57f0000000000000000000000000000000000000000000000000000000000000000836151ef565b905063ffffffff81106133db576040516307b2a52360e41b815260040160405180910390fd5b6040805180820182526008805465ffffffffffff168252825161010081019384905284936000939291602084019160099084908288855b82829054906101000a900462ffffff1662ffffff168152602001906003019060208260020104928301926001038202915080841161341257905050505091909252505081519192505065ffffffffffff8082161061348357604051630568cab760e31b815260040160405180910390fd5b61348e816001615203565b6008805465ffffffffffff191665ffffffffffff9290921691909117905560005b60088110156134f5578a62ffffff16836020015182600881106134d4576134d4614981565b602002015162ffffff16146134f557806134ed81614a90565b9150506134af565b6008811061351d576020830151604051630c4f769b60e41b815261074d918d91600401615224565b506040805160808101825263ffffffff909416845262ffffff8b16602085015261ffff8c169084015233606084015297509095509193505050509250925092565b806000106104c75760405162461bcd60e51b815260206004820152601260248201527166206d75737420626520706f73697469766560701b604482015260640161074d565b6000808a8a8a8a8a8a8a8a8a6040516020016135c79998979695949392919061523e565b60408051601f1981840301815291905280516020909101206001600160f01b0316600160f01b179150509998505050505050505050565b6136088282613b85565b5050565b336001600160a01b038216036136645760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161074d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b601c546001600160a01b03908116908216811461360857601c80546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912910160405180910390a15050565b82516001600160401b038084169116111561376f57825160405163012d824d60e01b81526001600160401b038085166004830152909116602482015260440161074d565b6040830151515160009015801561378d575060408401515160200151155b156137c5575082516001600160401b031660009081526007602090815260408083208287015162ffffff16845290915290205461381f565b83604001516040516020016137da91906152c7565b60408051601f19818403018152918152815160209283012086516001600160401b03166000908152600784528281208885015162ffffff168252909352912081905590505b6060840151516000816001600160401b0381111561383f5761383f614488565b604051908082528060200260200182016040528015613868578160200160208202803683370190505b5090506000826001600160401b0381111561388557613885614488565b6040519080825280601f01601f1916602001820160405280156138af576020820181803683370190505b5090506000836001600160401b038111156138cc576138cc614488565b6040519080825280602002602001820160405280156138ff57816020015b60608152602001906001900390816138ea5790505b5090506000805b85811015613a025760008a60600151828151811061392657613926614981565b6020908102919091010151905060008061394a8d600001518e602001518c86613c48565b9150915081156139895780868661ffff168151811061396b5761396b614981565b6020026020010181905250848061398190614d5b565b9550506139b8565b600160f81b8785815181106139a0576139a0614981565b60200101906001600160f81b031916908160001a9053505b82515188518990869081106139cf576139cf614981565b602002602001019065ffffffffffff16908165ffffffffffff1681525050505050806139fa81614a90565b915050613906565b5060608901515115613af95760008161ffff166001600160401b03811115613a2c57613a2c614488565b604051908082528060200260200182016040528015613a5f57816020015b6060815260200190600190039081613a4a5790505b50905060005b8261ffff16811015613abb57838181518110613a8357613a83614981565b6020026020010151828281518110613a9d57613a9d614981565b60200260200101819052508080613ab390614a90565b915050613a65565b507f47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7858583604051613aef939291906152fa565b60405180910390a1505b505050505050505050565b60008383811015613b1757600285850304015b610ad18184613289565b600081861015613b735760405162461bcd60e51b815260206004820181905260248201527f6c6566744761732063616e6e6f742065786365656420696e697469616c476173604482015260640161074d565b50633b9aca0094039190910101020290565b610100818114613bae57828282604051635c9d52ef60e11b815260040161074d939291906153a0565b613bb6613edf565b8181604051602001613bc8919061428b565b6040516020818303038152906040525114613be557613be56153c4565b6040805180820190915260085465ffffffffffff16815260208101613c0c858701876153da565b905280516008805465ffffffffffff191665ffffffffffff9092169190911781556020820151613c3f9060099083614018565b5061170e915050565b6000606081613c807f00000000000000000000000000000000000000000000000000000000000000006001600160401b0389166151ef565b845160808101516040519293509091600091613ca4918b918b918690602001614cb9565b60408051601f198184030181529181528151602092830120845165ffffffffffff16600090815260069093529120549091508114613d135760016040518060400160405280601081526020016f756e6b6e6f776e2063616c6c6261636b60801b81525094509450505050613e8a565b6040805160808101825263ffffffff8516815262ffffff8a1660208083019190915284015161ffff1681830152908301516001600160a01b031660608201528251600090613d6390838b8e612697565b6060808401518651918701516040519394509092600092635a47dd7160e01b92613d9292879190602401615461565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600b805466ff0000000000001916600160301b1790558b5160a0015191880151909250600091613e00916001600160601b039091169084613e93565b600b805466ff0000000000001916905590508015613e52575050935165ffffffffffff166000908152600660209081526040808320839055805191820190528181529097509550613e8a945050505050565b60016040518060400160405280601081526020016f195e1958dd5d1a5bdb8819985a5b195960821b8152509950995050505050505050505b94509492505050565b60005a611388811015613ea557600080fd5b611388810390508460408204820311613ebd57600080fd5b50823b613ec957600080fd5b60008083516020850160008789f1949350505050565b6040518061010001604052806008906020820280368337509192915050565b50805460008255906000526020600020908101906104c7919061409f565b828054828255906000526020600020908101928215613f71579160200282015b82811115613f7157825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613f3c565b50613f7d92915061409f565b5090565b600483019183908215613f715791602002820160005b83821115613fdb57835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302613f97565b801561400b5782816101000a81549063ffffffff0219169055600401602081600301049283019260010302613fdb565b5050613f7d92915061409f565b600183019183908215613f715791602002820160005b8382111561407057835183826101000a81548162ffffff021916908362ffffff160217905550926020019260030160208160020104928301926001030261402e565b801561400b5782816101000a81549062ffffff0219169055600301602081600201049283019260010302614070565b5b80821115613f7d57600081556001016140a0565b6001600160a01b03811681146104c757600080fd5b6000602082840312156140db57600080fd5b81356118d6816140b4565b60005b838110156141015781810151838201526020016140e9565b8381111561170e5750506000910152565b6000815180845261412a8160208601602086016140e6565b601f01601f19169290920160200192915050565b6020815260006118d66020830184614112565b803563ffffffff8116811461416557600080fd5b919050565b62ffffff811681146104c757600080fd5b600080600080600060a0868803121561419357600080fd5b61419c86614151565b94506141aa60208701614151565b93506141b860408701614151565b92506141c660608701614151565b915060808601356141d68161416a565b809150509295509295909350565b65ffffffffffff811681146104c757600080fd5b60006020828403121561420a57600080fd5b81356118d6816141e4565b600081518084526020808501945080840160005b8381101561424557815187529582019590820190600101614229565b509495945050505050565b6020815260006118d66020830184614215565b8060005b600881101561170e57815162ffffff16845260209384019390910190600101614267565b610100810161329d8284614263565b60008083601f8401126142ac57600080fd5b5081356001600160401b038111156142c357600080fd5b6020830191508360208260051b85010111156142de57600080fd5b9250929050565b600080600080604085870312156142fb57600080fd5b84356001600160401b038082111561431257600080fd5b61431e8883890161429a565b9096509450602087013591508082111561433757600080fd5b506143448782880161429a565b95989497509550505050565b6000806020838503121561436357600080fd5b82356001600160401b0381111561437957600080fd5b6143858582860161429a565b90969095509350505050565b60008083601f8401126143a357600080fd5b5081356001600160401b038111156143ba57600080fd5b6020830191508360208285010111156142de57600080fd5b60008060008060008060008060e0898b0312156143ee57600080fd5b606089018a8111156143ff57600080fd5b899850356001600160401b038082111561441857600080fd5b6144248c838d01614391565b909950975060808b013591508082111561443d57600080fd5b6144498c838d0161429a565b909750955060a08b013591508082111561446257600080fd5b5061446f8b828c0161429a565b999c989b50969995989497949560c00135949350505050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156144c0576144c0614488565b60405290565b60405160c081016001600160401b03811182821017156144c0576144c0614488565b604051608081016001600160401b03811182821017156144c0576144c0614488565b604051602081016001600160401b03811182821017156144c0576144c0614488565b604051601f8201601f191681016001600160401b038111828210171561455457614554614488565b604052919050565b60006001600160401b0382111561457557614575614488565b50601f01601f191660200190565b600082601f83011261459457600080fd5b81356145a76145a28261455c565b61452c565b8181528460208386010111156145bc57600080fd5b816020850160208301376000918101602001919091529392505050565b60006001600160401b038211156145f2576145f2614488565b5060051b60200190565b6000602080838503121561460f57600080fd5b82356001600160401b038082111561462657600080fd5b908401906040828703121561463a57600080fd5b61464261449e565b82358281111561465157600080fd5b61465d88828601614583565b825250838301358281111561467157600080fd5b80840193505086601f84011261468657600080fd5b823591506146966145a2836145d9565b82815260059290921b830184019184810190888411156146b557600080fd5b938501935b838510156146d3578435825293850193908501906146ba565b948201949094529695505050505050565b600080604083850312156146f757600080fd5b8235614702816140b4565b946020939093013593505050565b60006020828403121561472257600080fd5b81356001600160401b0381111561473857600080fd5b8201608081850312156118d657600080fd5b61ffff811681146104c757600080fd5b6001600160401b03811681146104c757600080fd5b80356141658161475a565b60008060006060848603121561478f57600080fd5b833561479a8161474a565b925060208401356147aa8161475a565b915060408401356147ba8161416a565b809150509250925092565b803560ff8116811461416557600080fd5b60008060008060008060008060008060c08b8d0312156147f557600080fd5b8a356001600160401b038082111561480c57600080fd5b6148188e838f0161429a565b909c509a5060208d013591508082111561483157600080fd5b61483d8e838f0161429a565b909a50985088915061485160408e016147c5565b975060608d013591508082111561486757600080fd5b6148738e838f01614391565b909750955085915061488760808e0161476f565b945060a08d013591508082111561489d57600080fd5b506148aa8d828e01614391565b915080935050809150509295989b9194979a5092959850565b600080604083850312156148d657600080fd5b82356148e1816140b4565b915060208301356148f1816140b4565b809150509250929050565b600080600080600060a0868803121561491457600080fd5b853561491f8161475a565b9450602086013561492f8161474a565b9350604086013561493f8161416a565b925061494d60608701614151565b915060808601356001600160401b0381111561496857600080fd5b61497488828901614583565b9150509295509295909350565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600063ffffffff838116908316818110156149ca576149ca614997565b039392505050565b60008160001904831182151516156149ec576149ec614997565b500290565b60008219821115614a0457614a04614997565b500190565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0384168152604060208201819052600090610ad19083018486614a09565b600060208284031215614a6957600080fd5b815180151581146118d657600080fd5b600082821015614a8b57614a8b614997565b500390565b600060018201614aa257614aa2614997565b5060010190565b600060ff821660ff84168060ff03821115614ac657614ac6614997565b019392505050565b8183823760009101908152919050565b8281526060826020830137600060809190910190815292915050565b60008251614b0c8184602087016140e6565b9190910192915050565b600060208284031215614b2857600080fd5b5051919050565b602080825260129082015271696e73756666696369656e742066756e647360701b604082015260600190565b60008083128015600160ff1b850184121615614b7957614b79614997565b6001600160ff1b0384018313811615614b9457614b94614997565b50500390565b600060ff821660ff84168160ff0481118215151615614bbb57614bbb614997565b029392505050565b600063ffffffff808316818516808303821115614be257614be2614997565b01949350505050565b600081518084526020808501945080840160005b838110156142455781516001600160a01b031687529582019590820190600101614bff565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614c548184018a614beb565b90508281036080840152614c688189614beb565b905060ff871660a084015282810360c0840152614c858187614112565b90506001600160401b03851660e0840152828103610100840152614ca98185614112565b9c9b505050505050505050505050565b60006001600160401b03808716835262ffffff8616602084015280851660408401526080606084015265ffffffffffff845116608084015261ffff60208501511660a084015260018060a01b0360408501511660c0840152606084015160c060e0850152614d2b610140850182614112565b60808601519092166101008501525060a0909301516001600160601b031661012090920191909152509392505050565b600061ffff808316818103614d7257614d72614997565b6001019392505050565b80516141658161475a565b600082601f830112614d9857600080fd5b8151614da66145a28261455c565b818152846020838601011115614dbb57600080fd5b61070c8260208301602087016140e6565b80516001600160601b038116811461416557600080fd5b600082601f830112614df457600080fd5b81516020614e046145a2836145d9565b82815260059290921b84018101918181019086841115614e2357600080fd5b8286015b848110156128605780516001600160401b0380821115614e4657600080fd5b90880190601f196040838c0382011215614e5f57600080fd5b614e6761449e565b8784015183811115614e7857600080fd5b840160c0818e0384011215614e8c57600080fd5b614e946144c6565b925088810151614ea3816141e4565b83526040810151614eb38161474a565b838a01526060810151614ec5816140b4565b6040840152608081015184811115614edc57600080fd5b614eea8e8b83850101614d87565b606085015250614efc60a08201614d7c565b6080840152614f0d60c08201614dcc565b60a084015250818152614f2260408501614dcc565b818901528652505050918301918301614e27565b80516001600160c01b038116811461416557600080fd5b600060208284031215614f5f57600080fd5b81516001600160401b0380821115614f7657600080fd5b9083019060808286031215614f8a57600080fd5b614f926144e8565b825182811115614fa157600080fd5b8301601f81018713614fb257600080fd5b8051614fc06145a2826145d9565b8082825260208201915060208360051b850101925089831115614fe257600080fd5b602084015b838110156150e257805187811115614ffe57600080fd5b850160a0818d03601f1901121561501457600080fd5b61501c6144e8565b602082015161502a8161475a565b8152604082015161503a8161416a565b60208201526040828e03605f1901121561505357600080fd5b61505b61450a565b8d607f84011261506a57600080fd5b61507261449e565b808f60a08601111561508357600080fd5b606085015b60a086018110156150a3578051835260209283019201615088565b50825250604082015260a0820151898111156150be57600080fd5b6150cd8e602083860101614de3565b60608301525084525060209283019201614fe7565b508452506150f591505060208401614f36565b602082015261510660408401614d7c565b60408201526060830151606082015280935050505092915050565b600063ffffffff808316818103614d7257614d72614997565b6001600160a01b03851681526001600160c01b0384166020808301919091526040808301859052608060608401819052845190840181905260009285810192909160a0860190855b818110156151b557855180516001600160401b0316845285015162ffffff16858401529484019491830191600101615182565b50909a9950505050505050505050565b634e487b7160e01b600052601260045260246000fd5b6000826151ea576151ea6151c5565b500690565b6000826151fe576151fe6151c5565b500490565b600065ffffffffffff808316818516808303821115614be257614be2614997565b62ffffff8316815261012081016118d66020830184614263565b8981526001600160a01b03891660208201526001600160401b038881166040830152610120606083018190526000916152798483018b614beb565b9150838203608085015261528d828a614beb565b915060ff881660a085015283820360c08501526152aa8288614112565b90861660e08501528381036101008501529050614ca98185614112565b815160408201908260005b60028110156152f15782518252602092830192909101906001016152d2565b50505092915050565b606080825284519082018190526000906020906080840190828801845b8281101561533b57815165ffffffffffff1684529284019290840190600101615317565b5050508381038285015261534f8187614112565b905083810360408501528085518083528383019150838160051b84010184880160005b838110156151b557601f1986840301855261538e838351614112565b94870194925090860190600101615372565b6040815260006153b4604083018587614a09565b9050826020830152949350505050565b634e487b7160e01b600052600160045260246000fd5b60006101008083850312156153ee57600080fd5b83601f8401126153fd57600080fd5b6040518181018181106001600160401b038211171561541e5761541e614488565b60405290830190808583111561543357600080fd5b845b838110156154565780356154488161416a565b825260209182019101615435565b509095945050505050565b65ffffffffffff841681526060602082015260006154826060830185614215565b82810360408401526154948185614112565b969550505050505056fea2646970667358221220ec5b48ff6e748a98f7b43fbe440717b70955673010274510fc08bd2c985dfb5164736f6c634300080f0033",
}

var VRFBeaconCoordinatorABI = VRFBeaconCoordinatorMetaData.ABI

var VRFBeaconCoordinatorBin = VRFBeaconCoordinatorMetaData.Bin

func DeployVRFBeaconCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, beaconPeriodBlocksArg *big.Int, keyProvider common.Address, keyID [32]byte) (common.Address, *types.Transaction, *VRFBeaconCoordinator, error) {
	parsed, err := VRFBeaconCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFBeaconCoordinatorBin), backend, link, beaconPeriodBlocksArg, keyProvider, keyID)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFBeaconCoordinator{VRFBeaconCoordinatorCaller: VRFBeaconCoordinatorCaller{contract: contract}, VRFBeaconCoordinatorTransactor: VRFBeaconCoordinatorTransactor{contract: contract}, VRFBeaconCoordinatorFilterer: VRFBeaconCoordinatorFilterer{contract: contract}}, nil
}

type VRFBeaconCoordinator struct {
	address common.Address
	abi     abi.ABI
	VRFBeaconCoordinatorCaller
	VRFBeaconCoordinatorTransactor
	VRFBeaconCoordinatorFilterer
}

type VRFBeaconCoordinatorCaller struct {
	contract *bind.BoundContract
}

type VRFBeaconCoordinatorTransactor struct {
	contract *bind.BoundContract
}

type VRFBeaconCoordinatorFilterer struct {
	contract *bind.BoundContract
}

type VRFBeaconCoordinatorSession struct {
	Contract     *VRFBeaconCoordinator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFBeaconCoordinatorCallerSession struct {
	Contract *VRFBeaconCoordinatorCaller
	CallOpts bind.CallOpts
}

type VRFBeaconCoordinatorTransactorSession struct {
	Contract     *VRFBeaconCoordinatorTransactor
	TransactOpts bind.TransactOpts
}

type VRFBeaconCoordinatorRaw struct {
	Contract *VRFBeaconCoordinator
}

type VRFBeaconCoordinatorCallerRaw struct {
	Contract *VRFBeaconCoordinatorCaller
}

type VRFBeaconCoordinatorTransactorRaw struct {
	Contract *VRFBeaconCoordinatorTransactor
}

func NewVRFBeaconCoordinator(address common.Address, backend bind.ContractBackend) (*VRFBeaconCoordinator, error) {
	abi, err := abi.JSON(strings.NewReader(VRFBeaconCoordinatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFBeaconCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinator{address: address, abi: abi, VRFBeaconCoordinatorCaller: VRFBeaconCoordinatorCaller{contract: contract}, VRFBeaconCoordinatorTransactor: VRFBeaconCoordinatorTransactor{contract: contract}, VRFBeaconCoordinatorFilterer: VRFBeaconCoordinatorFilterer{contract: contract}}, nil
}

func NewVRFBeaconCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*VRFBeaconCoordinatorCaller, error) {
	contract, err := bindVRFBeaconCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorCaller{contract: contract}, nil
}

func NewVRFBeaconCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFBeaconCoordinatorTransactor, error) {
	contract, err := bindVRFBeaconCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorTransactor{contract: contract}, nil
}

func NewVRFBeaconCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFBeaconCoordinatorFilterer, error) {
	contract, err := bindVRFBeaconCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorFilterer{contract: contract}, nil
}

func bindVRFBeaconCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFBeaconCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFBeaconCoordinator.Contract.VRFBeaconCoordinatorCaller.contract.Call(opts, result, method, params...)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.VRFBeaconCoordinatorTransactor.contract.Transfer(opts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.VRFBeaconCoordinatorTransactor.contract.Transact(opts, method, params...)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFBeaconCoordinator.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.contract.Transfer(opts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.contract.Transact(opts, method, params...)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) LINK() (common.Address, error) {
	return _VRFBeaconCoordinator.Contract.LINK(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) LINK() (common.Address, error) {
	return _VRFBeaconCoordinator.Contract.LINK(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "NUM_CONF_DELAYS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFBeaconCoordinator.Contract.NUMCONFDELAYS(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFBeaconCoordinator.Contract.NUMCONFDELAYS(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) GetBilling(opts *bind.CallOpts) (GetBilling,

	error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "getBilling")

	outstruct := new(GetBilling)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPriceGwei = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ReasonableGasPriceGwei = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ObservationPaymentGjuels = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.TransmissionPaymentGjuels = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.AccountingGas = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) GetBilling() (GetBilling,

	error) {
	return _VRFBeaconCoordinator.Contract.GetBilling(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) GetBilling() (GetBilling,

	error) {
	return _VRFBeaconCoordinator.Contract.GetBilling(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) GetBillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "getBillingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) GetBillingAccessController() (common.Address, error) {
	return _VRFBeaconCoordinator.Contract.GetBillingAccessController(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) GetBillingAccessController() (common.Address, error) {
	return _VRFBeaconCoordinator.Contract.GetBillingAccessController(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "getConfirmationDelays")

	if err != nil {
		return *new([8]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([8]*big.Int)).(*[8]*big.Int)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) GetConfirmationDelays() ([8]*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.GetConfirmationDelays(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) GetConfirmationDelays() ([8]*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.GetConfirmationDelays(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) IStartSlot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "i_StartSlot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) IStartSlot() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.IStartSlot(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) IStartSlot() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.IStartSlot(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "i_beaconPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.IBeaconPeriodBlocks(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.IBeaconPeriodBlocks(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _VRFBeaconCoordinator.Contract.LatestConfigDetails(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _VRFBeaconCoordinator.Contract.LatestConfigDetails(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _VRFBeaconCoordinator.Contract.LatestConfigDigestAndEpoch(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _VRFBeaconCoordinator.Contract.LatestConfigDigestAndEpoch(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.LinkAvailableForPayment(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.LinkAvailableForPayment(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) MaxErrorMsgLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "maxErrorMsgLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) MaxErrorMsgLength() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.MaxErrorMsgLength(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) MaxErrorMsgLength() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.MaxErrorMsgLength(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) MaxNumWords(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "maxNumWords")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) MaxNumWords() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.MaxNumWords(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) MaxNumWords() (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.MaxNumWords(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) MinDelay(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "minDelay")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) MinDelay() (uint16, error) {
	return _VRFBeaconCoordinator.Contract.MinDelay(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) MinDelay() (uint16, error) {
	return _VRFBeaconCoordinator.Contract.MinDelay(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "oracleObservationCount", transmitterAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _VRFBeaconCoordinator.Contract.OracleObservationCount(&_VRFBeaconCoordinator.CallOpts, transmitterAddress)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _VRFBeaconCoordinator.Contract.OracleObservationCount(&_VRFBeaconCoordinator.CallOpts, transmitterAddress)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "owedPayment", transmitterAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.OwedPayment(&_VRFBeaconCoordinator.CallOpts, transmitterAddress)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _VRFBeaconCoordinator.Contract.OwedPayment(&_VRFBeaconCoordinator.CallOpts, transmitterAddress)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) Owner() (common.Address, error) {
	return _VRFBeaconCoordinator.Contract.Owner(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) Owner() (common.Address, error) {
	return _VRFBeaconCoordinator.Contract.Owner(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) SKeyID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "s_keyID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) SKeyID() ([32]byte, error) {
	return _VRFBeaconCoordinator.Contract.SKeyID(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) SKeyID() ([32]byte, error) {
	return _VRFBeaconCoordinator.Contract.SKeyID(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "s_provingKeyHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) SProvingKeyHash() ([32]byte, error) {
	return _VRFBeaconCoordinator.Contract.SProvingKeyHash(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) SProvingKeyHash() ([32]byte, error) {
	return _VRFBeaconCoordinator.Contract.SProvingKeyHash(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFBeaconCoordinator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) TypeAndVersion() (string, error) {
	return _VRFBeaconCoordinator.Contract.TypeAndVersion(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorCallerSession) TypeAndVersion() (string, error) {
	return _VRFBeaconCoordinator.Contract.TypeAndVersion(&_VRFBeaconCoordinator.CallOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "acceptOwnership")
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.AcceptOwnership(&_VRFBeaconCoordinator.TransactOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.AcceptOwnership(&_VRFBeaconCoordinator.TransactOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.AcceptPayeeship(&_VRFBeaconCoordinator.TransactOpts, transmitter)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.AcceptPayeeship(&_VRFBeaconCoordinator.TransactOpts, transmitter)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) ExposeType(opts *bind.TransactOpts, arg0 VRFBeaconReportReport) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "exposeType", arg0)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) ExposeType(arg0 VRFBeaconReportReport) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.ExposeType(&_VRFBeaconCoordinator.TransactOpts, arg0)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) ExposeType(arg0 VRFBeaconReportReport) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.ExposeType(&_VRFBeaconCoordinator.TransactOpts, arg0)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) ForgetConsumerSubscriptionID(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "forgetConsumerSubscriptionID", consumers)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) ForgetConsumerSubscriptionID(consumers []common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.ForgetConsumerSubscriptionID(&_VRFBeaconCoordinator.TransactOpts, consumers)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) ForgetConsumerSubscriptionID(consumers []common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.ForgetConsumerSubscriptionID(&_VRFBeaconCoordinator.TransactOpts, consumers)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) KeyGenerated(opts *bind.TransactOpts, kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "keyGenerated", kd)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.KeyGenerated(&_VRFBeaconCoordinator.TransactOpts, kd)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.KeyGenerated(&_VRFBeaconCoordinator.TransactOpts, kd)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) NewKeyRequested(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "newKeyRequested")
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) NewKeyRequested() (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.NewKeyRequested(&_VRFBeaconCoordinator.TransactOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) NewKeyRequested() (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.NewKeyRequested(&_VRFBeaconCoordinator.TransactOpts)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) RedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "redeemRandomness", requestID)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) RedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.RedeemRandomness(&_VRFBeaconCoordinator.TransactOpts, requestID)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) RedeemRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.RedeemRandomness(&_VRFBeaconCoordinator.TransactOpts, requestID)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) RequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "requestRandomness", numWords, subID, confirmationDelayArg)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) RequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.RequestRandomness(&_VRFBeaconCoordinator.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) RequestRandomness(numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.RequestRandomness(&_VRFBeaconCoordinator.TransactOpts, numWords, subID, confirmationDelayArg)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) RequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "requestRandomnessFulfillment", subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) RequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.RequestRandomnessFulfillment(&_VRFBeaconCoordinator.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) RequestRandomnessFulfillment(subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.RequestRandomnessFulfillment(&_VRFBeaconCoordinator.TransactOpts, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "setBilling", maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetBilling(&_VRFBeaconCoordinator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetBilling(&_VRFBeaconCoordinator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetBillingAccessController(&_VRFBeaconCoordinator.TransactOpts, _billingAccessController)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetBillingAccessController(&_VRFBeaconCoordinator.TransactOpts, _billingAccessController)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetConfig(&_VRFBeaconCoordinator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetConfig(&_VRFBeaconCoordinator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "setPayees", transmitters, payees)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetPayees(&_VRFBeaconCoordinator.TransactOpts, transmitters, payees)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.SetPayees(&_VRFBeaconCoordinator.TransactOpts, transmitters, payees)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.TransferOwnership(&_VRFBeaconCoordinator.TransactOpts, to)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.TransferOwnership(&_VRFBeaconCoordinator.TransactOpts, to)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.TransferPayeeship(&_VRFBeaconCoordinator.TransactOpts, transmitter, proposed)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.TransferPayeeship(&_VRFBeaconCoordinator.TransactOpts, transmitter, proposed)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.Transmit(&_VRFBeaconCoordinator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.Transmit(&_VRFBeaconCoordinator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "withdrawFunds", recipient, amount)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.WithdrawFunds(&_VRFBeaconCoordinator.TransactOpts, recipient, amount)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.WithdrawFunds(&_VRFBeaconCoordinator.TransactOpts, recipient, amount)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "withdrawPayment", transmitter)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.WithdrawPayment(&_VRFBeaconCoordinator.TransactOpts, transmitter)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.WithdrawPayment(&_VRFBeaconCoordinator.TransactOpts, transmitter)
}

type VRFBeaconCoordinatorBillingAccessControllerSetIterator struct {
	Event *VRFBeaconCoordinatorBillingAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorBillingAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorBillingAccessControllerSet)
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
		it.Event = new(VRFBeaconCoordinatorBillingAccessControllerSet)
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

func (it *VRFBeaconCoordinatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*VRFBeaconCoordinatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorBillingAccessControllerSetIterator{contract: _VRFBeaconCoordinator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorBillingAccessControllerSet)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*VRFBeaconCoordinatorBillingAccessControllerSet, error) {
	event := new(VRFBeaconCoordinatorBillingAccessControllerSet)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorBillingSetIterator struct {
	Event *VRFBeaconCoordinatorBillingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorBillingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorBillingSet)
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
		it.Event = new(VRFBeaconCoordinatorBillingSet)
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

func (it *VRFBeaconCoordinatorBillingSetIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorBillingSet struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
	Raw                       types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*VRFBeaconCoordinatorBillingSetIterator, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorBillingSetIterator{contract: _VRFBeaconCoordinator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorBillingSet)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseBillingSet(log types.Log) (*VRFBeaconCoordinatorBillingSet, error) {
	event := new(VRFBeaconCoordinatorBillingSet)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorConfigSetIterator struct {
	Event *VRFBeaconCoordinatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorConfigSet)
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
		it.Event = new(VRFBeaconCoordinatorConfigSet)
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

func (it *VRFBeaconCoordinatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFBeaconCoordinatorConfigSetIterator, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorConfigSetIterator{contract: _VRFBeaconCoordinator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorConfigSet)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseConfigSet(log types.Log) (*VRFBeaconCoordinatorConfigSet, error) {
	event := new(VRFBeaconCoordinatorConfigSet)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorNewTransmissionIterator struct {
	Event *VRFBeaconCoordinatorNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorNewTransmission)
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
		it.Event = new(VRFBeaconCoordinatorNewTransmission)
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

func (it *VRFBeaconCoordinatorNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorNewTransmission struct {
	AggregatorRoundId uint32
	EpochAndRound     *big.Int
	Transmitter       common.Address
	JuelsPerFeeCoin   *big.Int
	ConfigDigest      [32]byte
	OutputsServed     []VRFBeaconReportOutputServed
	Raw               types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFBeaconCoordinatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorNewTransmissionIterator{contract: _VRFBeaconCoordinator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorNewTransmission)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseNewTransmission(log types.Log) (*VRFBeaconCoordinatorNewTransmission, error) {
	event := new(VRFBeaconCoordinatorNewTransmission)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorOraclePaidIterator struct {
	Event *VRFBeaconCoordinatorOraclePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorOraclePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorOraclePaid)
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
		it.Event = new(VRFBeaconCoordinatorOraclePaid)
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

func (it *VRFBeaconCoordinatorOraclePaidIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*VRFBeaconCoordinatorOraclePaidIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorOraclePaidIterator{contract: _VRFBeaconCoordinator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorOraclePaid)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseOraclePaid(log types.Log) (*VRFBeaconCoordinatorOraclePaid, error) {
	event := new(VRFBeaconCoordinatorOraclePaid)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorOwnershipTransferRequestedIterator struct {
	Event *VRFBeaconCoordinatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorOwnershipTransferRequested)
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
		it.Event = new(VRFBeaconCoordinatorOwnershipTransferRequested)
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

func (it *VRFBeaconCoordinatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconCoordinatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorOwnershipTransferRequestedIterator{contract: _VRFBeaconCoordinator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorOwnershipTransferRequested)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFBeaconCoordinatorOwnershipTransferRequested, error) {
	event := new(VRFBeaconCoordinatorOwnershipTransferRequested)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorOwnershipTransferredIterator struct {
	Event *VRFBeaconCoordinatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorOwnershipTransferred)
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
		it.Event = new(VRFBeaconCoordinatorOwnershipTransferred)
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

func (it *VRFBeaconCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconCoordinatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorOwnershipTransferredIterator{contract: _VRFBeaconCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorOwnershipTransferred)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*VRFBeaconCoordinatorOwnershipTransferred, error) {
	event := new(VRFBeaconCoordinatorOwnershipTransferred)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorPayeeshipTransferRequestedIterator struct {
	Event *VRFBeaconCoordinatorPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorPayeeshipTransferRequested)
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
		it.Event = new(VRFBeaconCoordinatorPayeeshipTransferRequested)
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

func (it *VRFBeaconCoordinatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*VRFBeaconCoordinatorPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorPayeeshipTransferRequestedIterator{contract: _VRFBeaconCoordinator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorPayeeshipTransferRequested)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*VRFBeaconCoordinatorPayeeshipTransferRequested, error) {
	event := new(VRFBeaconCoordinatorPayeeshipTransferRequested)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorPayeeshipTransferredIterator struct {
	Event *VRFBeaconCoordinatorPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorPayeeshipTransferred)
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
		it.Event = new(VRFBeaconCoordinatorPayeeshipTransferred)
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

func (it *VRFBeaconCoordinatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*VRFBeaconCoordinatorPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorPayeeshipTransferredIterator{contract: _VRFBeaconCoordinator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorPayeeshipTransferred)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParsePayeeshipTransferred(log types.Log) (*VRFBeaconCoordinatorPayeeshipTransferred, error) {
	event := new(VRFBeaconCoordinatorPayeeshipTransferred)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorRandomWordsFulfilledIterator struct {
	Event *VRFBeaconCoordinatorRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorRandomWordsFulfilled)
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
		it.Event = new(VRFBeaconCoordinatorRandomWordsFulfilled)
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

func (it *VRFBeaconCoordinatorRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorRandomWordsFulfilled struct {
	RequestIDs            []*big.Int
	SuccessfulFulfillment []byte
	TruncatedErrorData    [][]byte
	Raw                   types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFBeaconCoordinatorRandomWordsFulfilledIterator, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorRandomWordsFulfilledIterator{contract: _VRFBeaconCoordinator.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorRandomWordsFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorRandomWordsFulfilled)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFBeaconCoordinatorRandomWordsFulfilled, error) {
	event := new(VRFBeaconCoordinatorRandomWordsFulfilled)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorRandomnessFulfillmentRequestedIterator struct {
	Event *VRFBeaconCoordinatorRandomnessFulfillmentRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorRandomnessFulfillmentRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorRandomnessFulfillmentRequested)
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
		it.Event = new(VRFBeaconCoordinatorRandomnessFulfillmentRequested)
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

func (it *VRFBeaconCoordinatorRandomnessFulfillmentRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorRandomnessFulfillmentRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorRandomnessFulfillmentRequested struct {
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  uint64
	Callback               VRFBeaconTypesCallback
	Raw                    types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts) (*VRFBeaconCoordinatorRandomnessFulfillmentRequestedIterator, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "RandomnessFulfillmentRequested")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorRandomnessFulfillmentRequestedIterator{contract: _VRFBeaconCoordinator.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorRandomnessFulfillmentRequested) (event.Subscription, error) {

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "RandomnessFulfillmentRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorRandomnessFulfillmentRequested)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseRandomnessFulfillmentRequested(log types.Log) (*VRFBeaconCoordinatorRandomnessFulfillmentRequested, error) {
	event := new(VRFBeaconCoordinatorRandomnessFulfillmentRequested)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconCoordinatorRandomnessRequestedIterator struct {
	Event *VRFBeaconCoordinatorRandomnessRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconCoordinatorRandomnessRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconCoordinatorRandomnessRequested)
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
		it.Event = new(VRFBeaconCoordinatorRandomnessRequested)
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

func (it *VRFBeaconCoordinatorRandomnessRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconCoordinatorRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconCoordinatorRandomnessRequested struct {
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	Raw                    types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, nextBeaconOutputHeight []uint64) (*VRFBeaconCoordinatorRandomnessRequestedIterator, error) {

	var nextBeaconOutputHeightRule []interface{}
	for _, nextBeaconOutputHeightItem := range nextBeaconOutputHeight {
		nextBeaconOutputHeightRule = append(nextBeaconOutputHeightRule, nextBeaconOutputHeightItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "RandomnessRequested", nextBeaconOutputHeightRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorRandomnessRequestedIterator{contract: _VRFBeaconCoordinator.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorRandomnessRequested, nextBeaconOutputHeight []uint64) (event.Subscription, error) {

	var nextBeaconOutputHeightRule []interface{}
	for _, nextBeaconOutputHeightItem := range nextBeaconOutputHeight {
		nextBeaconOutputHeightRule = append(nextBeaconOutputHeightRule, nextBeaconOutputHeightItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "RandomnessRequested", nextBeaconOutputHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconCoordinatorRandomnessRequested)
				if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) ParseRandomnessRequested(log types.Log) (*VRFBeaconCoordinatorRandomnessRequested, error) {
	event := new(VRFBeaconCoordinatorRandomnessRequested)
	if err := _VRFBeaconCoordinator.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetBilling struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}
type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestConfigDigestAndEpoch struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFBeaconCoordinator.abi.Events["BillingAccessControllerSet"].ID:
		return _VRFBeaconCoordinator.ParseBillingAccessControllerSet(log)
	case _VRFBeaconCoordinator.abi.Events["BillingSet"].ID:
		return _VRFBeaconCoordinator.ParseBillingSet(log)
	case _VRFBeaconCoordinator.abi.Events["ConfigSet"].ID:
		return _VRFBeaconCoordinator.ParseConfigSet(log)
	case _VRFBeaconCoordinator.abi.Events["NewTransmission"].ID:
		return _VRFBeaconCoordinator.ParseNewTransmission(log)
	case _VRFBeaconCoordinator.abi.Events["OraclePaid"].ID:
		return _VRFBeaconCoordinator.ParseOraclePaid(log)
	case _VRFBeaconCoordinator.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFBeaconCoordinator.ParseOwnershipTransferRequested(log)
	case _VRFBeaconCoordinator.abi.Events["OwnershipTransferred"].ID:
		return _VRFBeaconCoordinator.ParseOwnershipTransferred(log)
	case _VRFBeaconCoordinator.abi.Events["PayeeshipTransferRequested"].ID:
		return _VRFBeaconCoordinator.ParsePayeeshipTransferRequested(log)
	case _VRFBeaconCoordinator.abi.Events["PayeeshipTransferred"].ID:
		return _VRFBeaconCoordinator.ParsePayeeshipTransferred(log)
	case _VRFBeaconCoordinator.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFBeaconCoordinator.ParseRandomWordsFulfilled(log)
	case _VRFBeaconCoordinator.abi.Events["RandomnessFulfillmentRequested"].ID:
		return _VRFBeaconCoordinator.ParseRandomnessFulfillmentRequested(log)
	case _VRFBeaconCoordinator.abi.Events["RandomnessRequested"].ID:
		return _VRFBeaconCoordinator.ParseRandomnessRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFBeaconCoordinatorBillingAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912")
}

func (VRFBeaconCoordinatorBillingSet) Topic() common.Hash {
	return common.HexToHash("0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f")
}

func (VRFBeaconCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (VRFBeaconCoordinatorNewTransmission) Topic() common.Hash {
	return common.HexToHash("0xe0c90b8e55243fcba0f8b68b201983b97f7a3d5aebd6dfa1a4082a07925cc744")
}

func (VRFBeaconCoordinatorOraclePaid) Topic() common.Hash {
	return common.HexToHash("0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c")
}

func (VRFBeaconCoordinatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFBeaconCoordinatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFBeaconCoordinatorPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (VRFBeaconCoordinatorPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (VRFBeaconCoordinatorRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7")
}

func (VRFBeaconCoordinatorRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0xa62e84e206cb87e2f6896795353c5358ff3d415d0bccc24e45c5fad83e17d03c")
}

func (VRFBeaconCoordinatorRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xc334d6f57be304c8192da2e39220c48e35f7e9afa16c541e68a6a859eff4dbc5")
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinator) Address() common.Address {
	return _VRFBeaconCoordinator.address
}

type VRFBeaconCoordinatorInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	GetBilling(opts *bind.CallOpts) (GetBilling,

		error)

	GetBillingAccessController(opts *bind.CallOpts) (common.Address, error)

	GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error)

	IStartSlot(opts *bind.CallOpts) (*big.Int, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	MaxErrorMsgLength(opts *bind.CallOpts) (*big.Int, error)

	MaxNumWords(opts *bind.CallOpts) (*big.Int, error)

	MinDelay(opts *bind.CallOpts) (uint16, error)

	OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error)

	OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SKeyID(opts *bind.CallOpts) ([32]byte, error)

	SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	ExposeType(opts *bind.TransactOpts, arg0 VRFBeaconReportReport) (*types.Transaction, error)

	ForgetConsumerSubscriptionID(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	KeyGenerated(opts *bind.TransactOpts, kd KeyDataStructKeyData) (*types.Transaction, error)

	NewKeyRequested(opts *bind.TransactOpts) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*types.Transaction, error)

	SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error)

	SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*VRFBeaconCoordinatorBillingAccessControllerSetIterator, error)

	WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorBillingAccessControllerSet) (event.Subscription, error)

	ParseBillingAccessControllerSet(log types.Log) (*VRFBeaconCoordinatorBillingAccessControllerSet, error)

	FilterBillingSet(opts *bind.FilterOpts) (*VRFBeaconCoordinatorBillingSetIterator, error)

	WatchBillingSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorBillingSet) (event.Subscription, error)

	ParseBillingSet(log types.Log) (*VRFBeaconCoordinatorBillingSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFBeaconCoordinatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFBeaconCoordinatorConfigSet, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFBeaconCoordinatorNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*VRFBeaconCoordinatorNewTransmission, error)

	FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*VRFBeaconCoordinatorOraclePaidIterator, error)

	WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error)

	ParseOraclePaid(log types.Log) (*VRFBeaconCoordinatorOraclePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconCoordinatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFBeaconCoordinatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconCoordinatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFBeaconCoordinatorOwnershipTransferred, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*VRFBeaconCoordinatorPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*VRFBeaconCoordinatorPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*VRFBeaconCoordinatorPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*VRFBeaconCoordinatorPayeeshipTransferred, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFBeaconCoordinatorRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorRandomWordsFulfilled) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFBeaconCoordinatorRandomWordsFulfilled, error)

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts) (*VRFBeaconCoordinatorRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorRandomnessFulfillmentRequested) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*VRFBeaconCoordinatorRandomnessFulfillmentRequested, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, nextBeaconOutputHeight []uint64) (*VRFBeaconCoordinatorRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorRandomnessRequested, nextBeaconOutputHeight []uint64) (event.Subscription, error)

	ParseRandomnessRequested(log types.Log) (*VRFBeaconCoordinatorRandomnessRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
