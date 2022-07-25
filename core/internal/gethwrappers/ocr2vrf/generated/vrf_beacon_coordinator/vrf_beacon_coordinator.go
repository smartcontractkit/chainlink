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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocksArg\",\"type\":\"uint256\"},{\"internalType\":\"contractDKG\",\"name\":\"keyProvider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"keyID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BeaconPeriodMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earliestAllowed\",\"type\":\"uint256\"}],\"name\":\"BlockTooRecent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"firstDelay\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"minDelay\",\"type\":\"uint16\"}],\"name\":\"ConfirmationDelayBlocksTooShort\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confirmationDelays\",\"type\":\"uint16[10]\"},{\"internalType\":\"uint8\",\"name\":\"violatingIndex\",\"type\":\"uint8\"}],\"name\":\"ConfirmationDelaysNotIncreasing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"separatorHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorTooOld\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"providedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"onchainHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorWrong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"keyProvider\",\"type\":\"address\"}],\"name\":\"KeyInfoMustComeFromProvider\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoWordsRequested\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confDelays\",\"type\":\"uint16[10]\"}],\"name\":\"NonZeroDelayAfterZeroDelay\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"expectedLength\",\"type\":\"uint256\"}],\"name\":\"OffchainConfigHasWrongLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"}],\"name\":\"RandomnessNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"ResponseMustBeRetrievedByRequester\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequestsReplaceContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySlotsReplaceContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"TooManyWords\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"UniverseHasEndedBangBangBang\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"occVersion\",\"type\":\"uint64\"}],\"name\":\"UnknownConfigVersion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"givenDelay\",\"type\":\"uint24\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"knownDelays\",\"type\":\"uint24[8]\"}],\"name\":\"UnknownConfirmationDelay\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconReport.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconReport.VRFOutput[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"recentBlockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structVRFBeaconReport.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"exposeType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"forgetConsumerSubscriptionID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"}],\"name\":\"getRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_StartSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structKeyDataStruct.KeyData\",\"name\":\"kd\",\"type\":\"tuple\"}],\"name\":\"keyGenerated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxErrorMsgLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDelay\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newKeyRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subID\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keyID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_provingKeyHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b50604051620057d3380380620057d383398101604081905262000034916200022f565b8181848681818181803380600081620000945760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c757620000c7816200016b565b5050506001600160a01b03166080526000829003620000f957604051632abc297960e01b815260040160405180910390fd5b60a082905260006200010c83436200027d565b905060008160a051620001209190620002b6565b90506200012e8143620002d0565b60c0525050601d80546001600160a01b0319166001600160a01b039990991698909817909755505050601e9290925550620002eb95505050505050565b336001600160a01b03821603620001c55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200008b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b03811681146200022c57600080fd5b50565b600080600080608085870312156200024657600080fd5b8451620002538162000216565b6020860151604087015191955093506200026d8162000216565b6060959095015193969295505050565b6000826200029b57634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b600082821015620002cb57620002cb620002a0565b500390565b60008219821115620002e657620002e6620002a0565b500190565b60805160a05160c05161546b62000368600039600061050f0152600081816104e8015281816106bb015281816132a2015281816132d1015281816133090152613bcb0152600081816102a70152818161150b015281816115c9015281816116b5015281816124740152818161289601526129b3015261546b6000f3fe608060405234801561001057600080fd5b506004361061021c5760003560e01c8063bbcdd0d811610125578063d09dc339116100ad578063e4902f821161007c578063e4902f821461057f578063eb5dcd6c146105a7578063f2fde38b146105ba578063f645dcb1146105cd578063fbffd2c1146105e057600080fd5b8063d09dc33914610531578063d57fc45a14610539578063dc92accf14610542578063e3d0e7121461056c57600080fd5b8063c4c92b37116100f4578063c4c92b37146104ae578063c63c4e9b146104bf578063cc31f7dd146104da578063cd0593df146104e3578063cf7e754a1461050a57600080fd5b8063bbcdd0d81461046e578063bf2732c714610477578063c10753291461048a578063c278e5b71461049d57600080fd5b80637a464944116101a85780639c849b30116101775780639c849b30146103f85780639e3616f41461040b578063afcb95d71461041e578063b121e14714610448578063b1dc65a41461045b57600080fd5b80637a4649441461039f57806381ff7048146103a75780638ac28d5a146103d45780638da5cb5b146103e757600080fd5b806329937268116101ef57806329937268146102e15780632f7527cc1461036057806355e487491461037a578063643dc1051461038457806379ba50971461039757600080fd5b80630b93e168146102215780630eafb25b1461024a578063181f5a771461026b5780631b6b6d23146102a2575b600080fd5b61023461022f366004614041565b6105f3565b6040516102419190614099565b60405180910390f35b61025d6102583660046140c1565b61078e565b604051908152602001610241565b6040805180820182526015815274565246426561636f6e20312e302e302d616c70686160581b602082015290516102419190614136565b6102c97f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610241565b610324600c54600160501b810463ffffffff90811692600160701b8304821692600160901b8104831692600160b01b82041691600160d01b90910462ffffff1690565b6040805163ffffffff9687168152948616602086015292851692840192909252909216606082015262ffffff909116608082015260a001610241565b610368600881565b60405160ff9091168152602001610241565b610382610893565b005b610382610392366004614173565b6108d8565b610382610abe565b61025d608081565b600d54600f54604080516000815264010000000090930463ffffffff166020840152820152606001610241565b6103826103e23660046140c1565b610b68565b6000546001600160a01b03166102c9565b610382610406366004614227565b610bda565b610382610419366004614292565b610dac565b600f546011546040805160008152602081019390935263ffffffff90911690820152606001610241565b6103826104563660046140c1565b610e3b565b610382610469366004614314565b610f17565b61025d6103e881565b61038261048536600461453e565b61139a565b610382610498366004614626565b611408565b6103826104ab366004614652565b50565b601c546001600160a01b03166102c9565b6104c7600381565b60405161ffff9091168152602001610241565b61025d601e5481565b61025d7f000000000000000000000000000000000000000000000000000000000000000081565b61025d7f000000000000000000000000000000000000000000000000000000000000000081565b61025d611693565b61025d601f5481565b6105556105503660046146bc565b61173f565b60405165ffffffffffff9091168152602001610241565b61038261057a366004614718565b61185c565b61059261058d3660046140c1565b611f85565b60405163ffffffff9091168152602001610241565b6103826105b5366004614805565b612034565b6103826105c83660046140c1565b61216c565b6105556105db36600461483e565b61217d565b6103826105ee3660046140c1565b61227d565b65ffffffffffff81166000818152600a602081815260408084208151608081018352815463ffffffff8116825262ffffff6401000000008204168286015261ffff600160381b820416938201939093526001600160a01b03600160481b84048116606083810191825298909752949093526001600160e81b0319909116905591511633146106b0576060810151604051638e30e82360e01b81526001600160a01b0390911660048201523360248201526044015b60405180910390fd5b80516000906106e6907f00000000000000000000000000000000000000000000000000000000000000009063ffffffff166148d9565b90506000826020015162ffffff16436106ff91906148f8565b905080821061072a576040516315ad27c360e01b8152600481018390524360248201526044016106a7565b6001600160401b03821115610755576040516302c6ef8160e11b8152600481018390526024016106a7565b60008281526007602090815260408083208287015162ffffff16845290915290205461078590869085908561228e565b95945050505050565b6001600160a01b03811660009081526012602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b031691810191909152906107f05750600092915050565b600c546020820151600091600160901b900463ffffffff169060169060ff16601f811061081f5761081f61490f565b600881049190910154600c54610852926007166004026101000a90910463ffffffff90811691600160301b900416614925565b63ffffffff1661086291906148d9565b61087090633b9aca006148d9565b905081604001516001600160601b03168161088b919061494a565b949350505050565b601d546001600160a01b03163381146108d05760405163292f4fb560e01b81523360048201526001600160a01b03821660248201526044016106a7565b506000601f55565b601c546001600160a01b03166108f66000546001600160a01b031690565b6001600160a01b0316336001600160a01b031614806109825750604051630d629b5f60e31b81526001600160a01b03821690636b14daf890610941903390600090369060040161498b565b602060405180830381865afa15801561095e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061098291906149b0565b6109ce5760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c60448201526064016106a7565b6109d6612462565b600c805467ffffffffffffffff60501b1916600160501b63ffffffff89811691820263ffffffff60701b191692909217600160701b8984169081029190911767ffffffffffffffff60901b1916600160901b89851690810263ffffffff60b01b191691909117600160b01b9489169485021762ffffff60d01b1916600160d01b62ffffff89169081029190911790955560408051938452602084019290925290820152606081019190915260808101919091527f0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f9060a00160405180910390a1505050505050565b6001546001600160a01b03163314610b115760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064016106a7565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b038181166000908152601a6020526040902054163314610bd15760405162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e20776974686472617700000000000000000060448201526064016106a7565b6104ab816127ea565b610be2612a0a565b828114610c315760405162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a6560448201526064016106a7565b60005b83811015610da5576000858583818110610c5057610c5061490f565b9050602002016020810190610c6591906140c1565b90506000848484818110610c7b57610c7b61490f565b9050602002016020810190610c9091906140c1565b6001600160a01b038084166000908152601a60205260409020549192501680158080610ccd5750826001600160a01b0316826001600160a01b0316145b610d0d5760405162461bcd60e51b81526020600482015260116024820152701c185e595948185b1c9958591e481cd95d607a1b60448201526064016106a7565b6001600160a01b038481166000908152601a6020526040902080546001600160a01b03191685831690811790915590831614610d8e57826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505050508080610d9d906149d2565b915050610c34565b5050505050565b610db4612a0a565b60005b81811015610e3657600060056000858585818110610dd757610dd761490f565b9050602002016020810190610dec91906140c1565b6001600160a01b031681526020810191909152604001600020805467ffffffffffffffff19166001600160401b039290921691909117905580610e2e816149d2565b915050610db7565b505050565b6001600160a01b038181166000908152601b6020526040902054163314610ea45760405162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e206163636570740060448201526064016106a7565b6001600160a01b038181166000818152601a602090815260408083208054336001600160a01b03198083168217909355601b909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a604080516101008082018352600c5460ff808216845291810464ffffffffff16602080850191909152600160301b820463ffffffff90811685870152600160501b830481166060860152600160701b830481166080860152600160901b8304811660a0860152600160b01b83041660c0850152600160d01b90910462ffffff1660e08401523360009081526012825293909320549394509092918c013591166110055760405162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d6974746572000000000000000060448201526064016106a7565b600f548b351461104f5760405162461bcd60e51b81526020600482015260156024820152740c6dedcccd2ce88d2cecae6e840dad2e6dac2e8c6d605b1b60448201526064016106a7565b61105d8a8a8a8a8a8a612a5f565b815161106a9060016149eb565b60ff1687146110bb5760405162461bcd60e51b815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e61747572657300000000000060448201526064016106a7565b86851461110a5760405162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e000060448201526064016106a7565b60008a8a60405161111c929190614a10565b604051908190038120611133918e90602001614a20565b60408051601f19818403018152828252805160209182012083830190925260008084529083018190529092509060005b8a8110156112cb5760006001858a84602081106111825761118261490f565b61118f91901a601b6149eb565b8f8f868181106111a1576111a161490f565b905060200201358e8e878181106111ba576111ba61490f565b90506020020135604051600081526020016040526040516111f7949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611219573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526013602090815290849020838501909452925460ff80821615158085526101009092041693830193909352909550925090506112a45760405162461bcd60e51b815260206004820152600f60248201526e39b4b3b730ba3ab9329032b93937b960891b60448201526064016106a7565b826020015160080260ff166001901b840193505080806112c3906149d2565b915050611163565b5081827e01010101010101010101010101010101010101010101010101010101010101161461132f5760405162461bcd60e51b815260206004820152601060248201526f323ab83634b1b0ba329039b4b3b732b960811b60448201526064016106a7565b506000915061137e9050838d836020020135848e8e8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612afc92505050565b905061138c83828633612f29565b505050505050505050505050565b601d546001600160a01b03163381146113d75760405163292f4fb560e01b81523360048201526001600160a01b03821660248201526044016106a7565b81516040516113e99190602001614a3c565b60408051601f198184030181529190528051602090910120601f555050565b6000546001600160a01b03163314806114925750601c54604051630d629b5f60e31b81526001600160a01b0390911690636b14daf890611451903390600090369060040161498b565b602060405180830381865afa15801561146e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061149291906149b0565b6114de5760405162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c60448201526064016106a7565b60006114e8613038565b6040516370a0823160e01b81523060048201529091506000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa158015611552573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115769190614a58565b9050818110156115bf5760405162461bcd60e51b8152602060048201526014602482015273696e73756666696369656e742062616c616e636560601b60448201526064016106a7565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663a9059cbb856116026115fc86866148f8565b87613202565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af115801561164d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061167191906149b0565b61168d5760405162461bcd60e51b81526004016106a790614a71565b50505050565b6040516370a0823160e01b815230600482015260009081906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa1580156116fc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117209190614a58565b9050600061172c613038565b90506117388183614a9d565b9250505090565b60008060008061174f878661321c565b92509250925065ffffffffffff83166000908152600a602090815260409182902084518154928601518487015160608801516001600160a01b0316600160481b027fffffff0000000000000000000000000000000000000000ffffffffffffffffff61ffff909216600160381b0291909116670100000000000000600160e81b031962ffffff9093166401000000000266ffffffffffffff1990961663ffffffff90941693909317949094171617919091179055516001600160401b038216907fc334d6f57be304c8192da2e39220c48e35f7e9afa16c541e68a6a859eff4dbc59061184790889062ffffff91909116815260200190565b60405180910390a250909150505b9392505050565b611864612a0a565b601f8911156118a85760405162461bcd60e51b815260206004820152601060248201526f746f6f206d616e79206f7261636c657360801b60448201526064016106a7565b8887146118f05760405162461bcd60e51b81526020600482015260166024820152750dee4c2c6d8ca40d8cadccee8d040dad2e6dac2e8c6d60531b60448201526064016106a7565b886118fc876003614adc565b60ff161061194c5760405162461bcd60e51b815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f2068696768000000000000000060448201526064016106a7565b6119588660ff166134d7565b6040805160e060208c02808301820190935260c082018c815260009383928f918f918291908601908490808284376000920191909152505050908252506040805160208c810282810182019093528c82529283019290918d918d91829185019084908082843760009201919091525050509082525060ff891660208083019190915260408051601f8a0183900483028101830182528981529201919089908990819084018382808284376000920191909152505050908252506001600160401b03861660208083019190915260408051601f870183900483028101830182528681529201919086908690819084018382808284376000920191909152505050915250600c805465ffffffffff00191690559050611a73612462565b60145460005b81811015611b2457600060148281548110611a9657611a9661490f565b6000918252602082200154601580546001600160a01b0390921693509084908110611ac357611ac361490f565b60009182526020808320909101546001600160a01b039485168352601382526040808420805461ffff1916905594168252601290529190912080546dffffffffffffffffffffffffffff191690555080611b1c816149d2565b915050611a79565b50611b3160146000613e58565b611b3d60156000613e58565b60005b825151811015611db6576013600084600001518381518110611b6457611b6461490f565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff1615611bd85760405162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e6572206164647265737300000000000000000060448201526064016106a7565b604080518082019091526001815260ff821660208201528351805160139160009185908110611c0957611c0961490f565b6020908102919091018101516001600160a01b03168252818101929092526040016000908120835181549484015161ffff1990951690151561ff0019161761010060ff90951694909402939093179092558401518051601292919084908110611c7457611c7461490f565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff1615611ce85760405162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d697474657220616464726573730000000060448201526064016106a7565b60405180606001604052806001151581526020018260ff16815260200160006001600160601b03168152506012600085602001518481518110611d2d57611d2d61490f565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020835181549385015194909201516001600160601b0316620100000262010000600160701b031960ff959095166101000261ff00199315159390931661ffff1990941693909317919091179290921617905580611dae816149d2565b915050611b40565b5081518051611dcd91601491602090910190613e76565b506020808301518051611de4926015920190613e76565b506040820151600c805460ff191660ff909216919091179055600d805467ffffffff0000000019811664010000000063ffffffff43811682029283179094558204831692600092611e3c929082169116176001614b05565b905080600d60006101000a81548163ffffffff021916908363ffffffff1602179055506000611e9046308463ffffffff16886000015189602001518a604001518b606001518c608001518d60a0015161351c565b905080600f600001819055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05838284886000015189602001518a604001518b606001518c608001518d60a00151604051611ef399989796959493929190614b66565b60405180910390a1600c54600160301b900463ffffffff1660005b865151811015611f685781601682601f8110611f2c57611f2c61490f565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff1602179055508080611f60906149d2565b915050611f0e565b50611f738b8b613577565b50505050505050505050505050505050565b6001600160a01b03811660009081526012602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b03169181019190915290611fe75750600092915050565b6016816020015160ff16601f81106120015761200161490f565b600881049190910154600c54611855926007166004026101000a90910463ffffffff90811691600160301b900416614925565b6001600160a01b038281166000908152601a602052604090205416331461209d5760405162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e2075706461746500000060448201526064016106a7565b6001600160a01b03811633036120f55760405162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016106a7565b6001600160a01b038083166000908152601b6020526040902080548383166001600160a01b031982168117909255909116908114610e36576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b612174612a0a565b6104ab81613585565b600080600061218c878761321c565b925050915060006040518060c001604052808465ffffffffffff1681526020018961ffff168152602001336001600160a01b031681526020018681526020018a6001600160401b031681526020018763ffffffff166001600160601b0316815250905081878a836040516020016122069493929190614bfb565b60408051601f19818403018152828252805160209182012065ffffffffffff871660009081526006909252919020557fa62e84e206cb87e2f6896795353c5358ff3d415d0bccc24e45c5fad83e17d03c906122689084908a908d908690614bfb565b60405180910390a15090979650505050505050565b612285612a0a565b6104ab8161362e565b6060826122c75760405163c7d41b1b60e01b815265ffffffffffff861660048201526001600160401b03831660248201526044016106a7565b6040805165ffffffffffff8716602080830191909152865163ffffffff168284015286015162ffffff166060808301919091529186015161ffff166080820152908501516001600160a01b031660a082015260c0810184905260009060e0016040516020818303038152906040528051906020012090506103e8856040015161ffff16111561237d576040808601519051634a90778560e01b815261ffff90911660048201526103e860248201526044016106a7565b6000856040015161ffff166001600160401b0381111561239f5761239f6143ca565b6040519080825280602002602001820160405280156123c8578160200160208202803683370190505b50905060005b866040015161ffff168161ffff16101561245757828160405160200161240b92919091825260f01b6001600160f01b031916602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff168151811061243a5761243a61490f565b60209081029190910101528061244f81614c9d565b9150506123ce565b509695505050505050565b600c54604080516103e08101918290527f000000000000000000000000000000000000000000000000000000000000000092600160301b900463ffffffff169160009190601690601f908285855b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116124b0579050505050505090506000601580548060200260200160405190810160405280929190818152602001828054801561254b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161252d575b5050505050905060005b81518110156127dc576000601260008484815181106125765761257661490f565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160029054906101000a90046001600160601b03166001600160601b031690506000601260008585815181106125d8576125d861490f565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160026101000a8154816001600160601b0302191690836001600160601b0316021790555060008483601f811061263b5761263b61490f565b6020020151600c5490870363ffffffff9081169250600160901b909104168102633b9aca0002820180156127d1576000601a60008787815181106126815761268161490f565b6020908102919091018101516001600160a01b03908116835290820192909252604090810160002054905163a9059cbb60e01b815290821660048201819052602482018590529250908a169063a9059cbb906044016020604051808303816000875af11580156126f5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061271991906149b0565b6127355760405162461bcd60e51b81526004016106a790614a71565b878786601f81106127485761274861490f565b602002019063ffffffff16908163ffffffff1681525050886001600160a01b0316816001600160a01b03168787815181106127855761278561490f565b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c856040516127c791815260200190565b60405180910390a4505b505050600101612555565b50610da5601683601f613edb565b6001600160a01b0381166000908152601260209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046001600160601b031692810192909252612847575050565b60006128528361078e565b90508015610e36576001600160a01b038381166000908152601a60205260409081902054905163a9059cbb60e01b81529082166004820181905260248201849052917f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af11580156128df573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061290391906149b0565b61291f5760405162461bcd60e51b81526004016106a790614a71565b600c60000160069054906101000a900463ffffffff166016846020015160ff16601f811061294f5761294f61490f565b6008810491909101805460079092166004026101000a63ffffffff8181021990931693909216919091029190911790556001600160a01b03848116600081815260126020908152604091829020805462010000600160701b031916905590518581527f0000000000000000000000000000000000000000000000000000000000000000841693851692917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c910160405180910390a450505050565b6000546001600160a01b03163314612a5d5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016106a7565b565b6000612a6c8260206148d9565b612a778560206148d9565b612a838861014461494a565b612a8d919061494a565b612a97919061494a565b612aa290600061494a565b9050368114612af35760405162461bcd60e51b815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d61746368000000000000000060448201526064016106a7565b50505050505050565b60008082806020019051810190612b139190614e8f565b64ffffffffff85166020880152604087018051919250612b3282615063565b63ffffffff1663ffffffff168152505085600c60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600e6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160000160126101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160000160166101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600001601a6101000a81548162ffffff021916908362ffffff160217905550905050600081604001516001600160401b031640905080826060015114612cd4576060820151604080840151905163aed0afe560e01b81526004810192909252602482018390526001600160401b031660448201526064016106a7565b6000808360000151516001600160401b03811115612cf457612cf46143ca565b604051908082528060200260200182016040528015612d3957816020015b6040805180820190915260008082526020820152815260200190600190039081612d125790505b50905060005b845151811015612e0957600085600001518281518110612d6157612d6161490f565b60200260200101519050612d7e81876040015188602001516136a4565b60408101515151151580612d9a57506040810151516020015115155b15612df657604051806040016040528082600001516001600160401b03168152602001826020015162ffffff16815250838381518110612ddc57612ddc61490f565b60200260200101819052508380612df290614c9d565b9450505b5080612e01816149d2565b915050612d3f565b5060008261ffff166001600160401b03811115612e2857612e286143ca565b604051908082528060200260200182016040528015612e6d57816020015b6040805180820190915260008082526020820152815260200190600190039081612e465790505b50905060005b8361ffff16811015612ec957828181518110612e9157612e9161490f565b6020026020010151828281518110612eab57612eab61490f565b60200260200101819052508080612ec1906149d2565b915050612e73565b50896040015163ffffffff167f7484067466b4f2452757769a8dc9a8b41497154367515673c79386f9f0b74f163387602001518c8c86604051612f1095949392919061507c565b60405180910390a2505050506020015195945050505050565b6000612f50633b9aca003a04866080015163ffffffff16876060015163ffffffff16613a7d565b90506010360260005a90506000612f798663ffffffff1685858b60e0015162ffffff1686613a9a565b90506000670de0b6b3a76400006001600160c01b03891683026001600160a01b03881660009081526012602052604090205460c08c01519290910492506201000090046001600160601b039081169163ffffffff16633b9aca000282840101908116821115612fee575050505050505061168d565b6001600160a01b038816600090815260126020526040902080546001600160601b03909216620100000262010000600160701b031990921691909117905550505050505050505050565b600080601580548060200260200160405190810160405280929190818152602001828054801561309157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613073575b50508351600c54604080516103e08101918290529697509195600160301b90910463ffffffff169450600093509150601690601f908285855b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116130ca5790505050505050905060005b8381101561315d578181601f811061312a5761312a61490f565b60200201516131399084614925565b6131499063ffffffff168761494a565b955080613155816149d2565b915050613110565b50600c5461317c90600160901b900463ffffffff16633b9aca006148d9565b61318690866148d9565b945060005b838110156131fa57601260008683815181106131a9576131a961490f565b6020908102919091018101516001600160a01b03168252810191909152604001600020546131e6906201000090046001600160601b03168761494a565b9550806131f2816149d2565b91505061318b565b505050505090565b600081831015613213575081613216565b50805b92915050565b604080516080810182526000808252602082018190529181018290526060810182905260006103e88561ffff16111561327657604051634a90778560e01b815261ffff861660048201526103e860248201526044016106a7565b8461ffff1660000361329b576040516308fad2a760e01b815260040160405180910390fd5b60006132c77f00000000000000000000000000000000000000000000000000000000000000004361512b565b90506000816132f67f00000000000000000000000000000000000000000000000000000000000000004361494a565b61330091906148f8565b9050600061332e7f00000000000000000000000000000000000000000000000000000000000000008361513f565b905063ffffffff8110613354576040516307b2a52360e41b815260040160405180910390fd5b6040805180820182526008805465ffffffffffff168252825161010081019384905284936000939291602084019160099084908288855b82829054906101000a900462ffffff1662ffffff168152602001906003019060208260020104928301926001038202915080841161338b57905050505091909252505081519192505065ffffffffffff808216106133fc57604051630568cab760e31b815260040160405180910390fd5b613407816001615153565b6008805465ffffffffffff191665ffffffffffff9290921691909117905560005b600881101561346e578a62ffffff168360200151826008811061344d5761344d61490f565b602002015162ffffff161461346e5780613466816149d2565b915050613428565b60088110613496576020830151604051630c4f769b60e41b81526106a7918d9160040161519c565b506040805160808101825263ffffffff909416845262ffffff8b16602085015261ffff8c169084015233606084015297509095509193505050509250925092565b806000106104ab5760405162461bcd60e51b815260206004820152601260248201527166206d75737420626520706f73697469766560701b60448201526064016106a7565b6000808a8a8a8a8a8a8a8a8a604051602001613540999897969594939291906151b6565b60408051601f1981840301815291905280516020909101206001600160f01b0316600160f01b179150509998505050505050505050565b6135818282613afe565b5050565b336001600160a01b038216036135dd5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016106a7565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b601c546001600160a01b03908116908216811461358157601c80546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912910160405180910390a15050565b82516001600160401b03808416911611156136e857825160405163012d824d60e01b81526001600160401b03808516600483015290911660248201526044016106a7565b60408301515151600090158015613706575060408401515160200151155b1561373e575082516001600160401b031660009081526007602090815260408083208287015162ffffff168452909152902054613798565b8360400151604051602001613753919061523f565b60408051601f19818403018152918152815160209283012086516001600160401b03166000908152600784528281208885015162ffffff168252909352912081905590505b6060840151516000816001600160401b038111156137b8576137b86143ca565b6040519080825280602002602001820160405280156137e1578160200160208202803683370190505b5090506000826001600160401b038111156137fe576137fe6143ca565b6040519080825280601f01601f191660200182016040528015613828576020820181803683370190505b5090506000836001600160401b03811115613845576138456143ca565b60405190808252806020026020018201604052801561387857816020015b60608152602001906001900390816138635790505b5090506000805b8581101561397b5760008a60600151828151811061389f5761389f61490f565b602090810291909101015190506000806138c38d600001518e602001518c86613bc1565b9150915081156139025780868661ffff16815181106138e4576138e461490f565b602002602001018190525084806138fa90614c9d565b955050613931565b600160f81b8785815181106139195761391961490f565b60200101906001600160f81b031916908160001a9053505b82515188518990869081106139485761394861490f565b602002602001019065ffffffffffff16908165ffffffffffff168152505050505080613973816149d2565b91505061387f565b5060608901515115613a725760008161ffff166001600160401b038111156139a5576139a56143ca565b6040519080825280602002602001820160405280156139d857816020015b60608152602001906001900390816139c35790505b50905060005b8261ffff16811015613a34578381815181106139fc576139fc61490f565b6020026020010151828281518110613a1657613a1661490f565b60200260200101819052508080613a2c906149d2565b9150506139de565b507f47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7858583604051613a6893929190615272565b60405180910390a1505b505050505050505050565b60008383811015613a9057600285850304015b6107858184613202565b600081861015613aec5760405162461bcd60e51b815260206004820181905260248201527f6c6566744761732063616e6e6f742065786365656420696e697469616c47617360448201526064016106a7565b50633b9aca0094039190910101020290565b610100818114613b2757828282604051635c9d52ef60e11b81526004016106a793929190615328565b613b2f613f72565b8181604051602001613b41919061534c565b6040516020818303038152906040525114613b5e57613b5e61535b565b6040805180820190915260085465ffffffffffff16815260208101613b8585870187615371565b905280516008805465ffffffffffff191665ffffffffffff9092169190911781556020820151613bb89060099083613f91565b5061168d915050565b6000606081613bf97f00000000000000000000000000000000000000000000000000000000000000006001600160401b03891661513f565b845160808101516040519293509091600091613c1d918b918b918690602001614bfb565b60408051601f198184030181529181528151602092830120845165ffffffffffff16600090815260069093529120549091508114613c8c5760016040518060400160405280601081526020016f756e6b6e6f776e2063616c6c6261636b60801b81525094509450505050613e03565b6040805160808101825263ffffffff8516815262ffffff8a1660208083019190915284015161ffff1681830152908301516001600160a01b031660608201528251600090613cdc90838b8e61228e565b6060808401518651918701516040519394509092600092635a47dd7160e01b92613d0b928791906024016153f8565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600b805466ff0000000000001916600160301b1790558b5160a0015191880151909250600091613d79916001600160601b039091169084613e0c565b600b805466ff0000000000001916905590508015613dcb575050935165ffffffffffff166000908152600660209081526040808320839055805191820190528181529097509550613e03945050505050565b60016040518060400160405280601081526020016f195e1958dd5d1a5bdb8819985a5b195960821b8152509950995050505050505050505b94509492505050565b60005a611388811015613e1e57600080fd5b611388810390508460408204820311613e3657600080fd5b50823b613e4257600080fd5b60008083516020850160008789f1949350505050565b50805460008255906000526020600020908101906104ab9190614018565b828054828255906000526020600020908101928215613ecb579160200282015b82811115613ecb57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613e96565b50613ed7929150614018565b5090565b600483019183908215613ecb5791602002820160005b83821115613f3557835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302613ef1565b8015613f655782816101000a81549063ffffffff0219169055600401602081600301049283019260010302613f35565b5050613ed7929150614018565b6040518061010001604052806008906020820280368337509192915050565b600183019183908215613ecb5791602002820160005b83821115613fe957835183826101000a81548162ffffff021916908362ffffff1602179055509260200192600301602081600201049283019260010302613fa7565b8015613f655782816101000a81549062ffffff0219169055600301602081600201049283019260010302613fe9565b5b80821115613ed75760008155600101614019565b65ffffffffffff811681146104ab57600080fd5b60006020828403121561405357600080fd5b81356118558161402d565b600081518084526020808501945080840160005b8381101561408e57815187529582019590820190600101614072565b509495945050505050565b602081526000611855602083018461405e565b6001600160a01b03811681146104ab57600080fd5b6000602082840312156140d357600080fd5b8135611855816140ac565b60005b838110156140f95781810151838201526020016140e1565b8381111561168d5750506000910152565b600081518084526141228160208601602086016140de565b601f01601f19169290920160200192915050565b602081526000611855602083018461410a565b803563ffffffff8116811461415d57600080fd5b919050565b62ffffff811681146104ab57600080fd5b600080600080600060a0868803121561418b57600080fd5b61419486614149565b94506141a260208701614149565b93506141b060408701614149565b92506141be60608701614149565b915060808601356141ce81614162565b809150509295509295909350565b60008083601f8401126141ee57600080fd5b5081356001600160401b0381111561420557600080fd5b6020830191508360208260051b850101111561422057600080fd5b9250929050565b6000806000806040858703121561423d57600080fd5b84356001600160401b038082111561425457600080fd5b614260888389016141dc565b9096509450602087013591508082111561427957600080fd5b50614286878288016141dc565b95989497509550505050565b600080602083850312156142a557600080fd5b82356001600160401b038111156142bb57600080fd5b6142c7858286016141dc565b90969095509350505050565b60008083601f8401126142e557600080fd5b5081356001600160401b038111156142fc57600080fd5b60208301915083602082850101111561422057600080fd5b60008060008060008060008060e0898b03121561433057600080fd5b606089018a81111561434157600080fd5b899850356001600160401b038082111561435a57600080fd5b6143668c838d016142d3565b909950975060808b013591508082111561437f57600080fd5b61438b8c838d016141dc565b909750955060a08b01359150808211156143a457600080fd5b506143b18b828c016141dc565b999c989b50969995989497949560c00135949350505050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715614402576144026143ca565b60405290565b60405160c081016001600160401b0381118282101715614402576144026143ca565b604051608081016001600160401b0381118282101715614402576144026143ca565b604051602081016001600160401b0381118282101715614402576144026143ca565b604051601f8201601f191681016001600160401b0381118282101715614496576144966143ca565b604052919050565b60006001600160401b038211156144b7576144b76143ca565b50601f01601f191660200190565b600082601f8301126144d657600080fd5b81356144e96144e48261449e565b61446e565b8181528460208386010111156144fe57600080fd5b816020850160208301376000918101602001919091529392505050565b60006001600160401b03821115614534576145346143ca565b5060051b60200190565b6000602080838503121561455157600080fd5b82356001600160401b038082111561456857600080fd5b908401906040828703121561457c57600080fd5b6145846143e0565b82358281111561459357600080fd5b61459f888286016144c5565b82525083830135828111156145b357600080fd5b80840193505086601f8401126145c857600080fd5b823591506145d86144e48361451b565b82815260059290921b830184019184810190888411156145f757600080fd5b938501935b83851015614615578435825293850193908501906145fc565b948201949094529695505050505050565b6000806040838503121561463957600080fd5b8235614644816140ac565b946020939093013593505050565b60006020828403121561466457600080fd5b81356001600160401b0381111561467a57600080fd5b82016080818503121561185557600080fd5b61ffff811681146104ab57600080fd5b6001600160401b03811681146104ab57600080fd5b803561415d8161469c565b6000806000606084860312156146d157600080fd5b83356146dc8161468c565b925060208401356146ec8161469c565b915060408401356146fc81614162565b809150509250925092565b803560ff8116811461415d57600080fd5b60008060008060008060008060008060c08b8d03121561473757600080fd5b8a356001600160401b038082111561474e57600080fd5b61475a8e838f016141dc565b909c509a5060208d013591508082111561477357600080fd5b61477f8e838f016141dc565b909a50985088915061479360408e01614707565b975060608d01359150808211156147a957600080fd5b6147b58e838f016142d3565b90975095508591506147c960808e016146b1565b945060a08d01359150808211156147df57600080fd5b506147ec8d828e016142d3565b915080935050809150509295989b9194979a5092959850565b6000806040838503121561481857600080fd5b8235614823816140ac565b91506020830135614833816140ac565b809150509250929050565b600080600080600060a0868803121561485657600080fd5b85356148618161469c565b945060208601356148718161468c565b9350604086013561488181614162565b925061488f60608701614149565b915060808601356001600160401b038111156148aa57600080fd5b6148b6888289016144c5565b9150509295509295909350565b634e487b7160e01b600052601160045260246000fd5b60008160001904831182151516156148f3576148f36148c3565b500290565b60008282101561490a5761490a6148c3565b500390565b634e487b7160e01b600052603260045260246000fd5b600063ffffffff83811690831681811015614942576149426148c3565b039392505050565b6000821982111561495d5761495d6148c3565b500190565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b03841681526040602082018190526000906107859083018486614962565b6000602082840312156149c257600080fd5b8151801515811461185557600080fd5b6000600182016149e4576149e46148c3565b5060010190565b600060ff821660ff84168060ff03821115614a0857614a086148c3565b019392505050565b8183823760009101908152919050565b8281526060826020830137600060809190910190815292915050565b60008251614a4e8184602087016140de565b9190910192915050565b600060208284031215614a6a57600080fd5b5051919050565b602080825260129082015271696e73756666696369656e742066756e647360701b604082015260600190565b60008083128015600160ff1b850184121615614abb57614abb6148c3565b6001600160ff1b0384018313811615614ad657614ad66148c3565b50500390565b600060ff821660ff84168160ff0481118215151615614afd57614afd6148c3565b029392505050565b600063ffffffff808316818516808303821115614b2457614b246148c3565b01949350505050565b600081518084526020808501945080840160005b8381101561408e5781516001600160a01b031687529582019590820190600101614b41565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152614b968184018a614b2d565b90508281036080840152614baa8189614b2d565b905060ff871660a084015282810360c0840152614bc7818761410a565b90506001600160401b03851660e0840152828103610100840152614beb818561410a565b9c9b505050505050505050505050565b60006001600160401b03808716835262ffffff8616602084015280851660408401526080606084015265ffffffffffff845116608084015261ffff60208501511660a084015260018060a01b0360408501511660c0840152606084015160c060e0850152614c6d61014085018261410a565b60808601519092166101008501525060a0909301516001600160601b031661012090920191909152509392505050565b600061ffff808316818103614cb457614cb46148c3565b6001019392505050565b805161415d8161469c565b600082601f830112614cda57600080fd5b8151614ce86144e48261449e565b818152846020838601011115614cfd57600080fd5b61088b8260208301602087016140de565b80516001600160601b038116811461415d57600080fd5b600082601f830112614d3657600080fd5b81516020614d466144e48361451b565b82815260059290921b84018101918181019086841115614d6557600080fd5b8286015b848110156124575780516001600160401b0380821115614d8857600080fd5b90880190601f196040838c0382011215614da157600080fd5b614da96143e0565b8784015183811115614dba57600080fd5b840160c0818e0384011215614dce57600080fd5b614dd6614408565b925088810151614de58161402d565b83526040810151614df58161468c565b838a01526060810151614e07816140ac565b6040840152608081015184811115614e1e57600080fd5b614e2c8e8b83850101614cc9565b606085015250614e3e60a08201614cbe565b6080840152614e4f60c08201614d0e565b60a084015250818152614e6460408501614d0e565b818901528652505050918301918301614d69565b80516001600160c01b038116811461415d57600080fd5b600060208284031215614ea157600080fd5b81516001600160401b0380821115614eb857600080fd5b9083019060808286031215614ecc57600080fd5b614ed461442a565b825182811115614ee357600080fd5b8301601f81018713614ef457600080fd5b8051614f026144e48261451b565b8082825260208201915060208360051b850101925089831115614f2457600080fd5b602084015b8381101561502457805187811115614f4057600080fd5b850160a0818d03601f19011215614f5657600080fd5b614f5e61442a565b6020820151614f6c8161469c565b81526040820151614f7c81614162565b60208201526040828e03605f19011215614f9557600080fd5b614f9d61444c565b8d607f840112614fac57600080fd5b614fb46143e0565b808f60a086011115614fc557600080fd5b606085015b60a08601811015614fe5578051835260209283019201614fca565b50825250604082015260a08201518981111561500057600080fd5b61500f8e602083860101614d25565b60608301525084525060209283019201614f29565b5084525061503791505060208401614e78565b602082015261504860408401614cbe565b60408201526060830151606082015280935050505092915050565b600063ffffffff808316818103614cb457614cb46148c3565b6001600160a01b03861681526001600160c01b038516602080830191909152604080830186905264ffffffffff8516606084015260a060808401819052845190840181905260009285810192909160c0860190855b8181101561510457855180516001600160401b0316845285015162ffffff168584015294840194918301916001016150d1565b50909b9a5050505050505050505050565b634e487b7160e01b600052601260045260246000fd5b60008261513a5761513a615115565b500690565b60008261514e5761514e615115565b500490565b600065ffffffffffff808316818516808303821115614b2457614b246148c3565b8060005b600881101561168d57815162ffffff16845260209384019390910190600101615178565b62ffffff8316815261012081016118556020830184615174565b8981526001600160a01b03891660208201526001600160401b038881166040830152610120606083018190526000916151f18483018b614b2d565b91508382036080850152615205828a614b2d565b915060ff881660a085015283820360c0850152615222828861410a565b90861660e08501528381036101008501529050614beb818561410a565b815160408201908260005b600281101561526957825182526020928301929091019060010161524a565b50505092915050565b606080825284519082018190526000906020906080840190828801845b828110156152b357815165ffffffffffff168452928401929084019060010161528f565b505050838103828501526152c7818761410a565b905083810360408501528085518083528383019150838160051b84010184880160005b8381101561531857601f1986840301855261530683835161410a565b948701949250908601906001016152ea565b50909a9950505050505050505050565b60408152600061533c604083018587614962565b9050826020830152949350505050565b61010081016132168284615174565b634e487b7160e01b600052600160045260246000fd5b600061010080838503121561538557600080fd5b83601f84011261539457600080fd5b6040518181018181106001600160401b03821117156153b5576153b56143ca565b6040529083019080858311156153ca57600080fd5b845b838110156153ed5780356153df81614162565b8252602091820191016153cc565b509095945050505050565b65ffffffffffff84168152606060208201526000615419606083018561405e565b828103604084015261542b818561410a565b969550505050505056fea2646970667358221220e338e0d0f8b461ac0f57baa583e1385fa9c296b5561c6d24fe064e29783fd75864736f6c634300080f0033",
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

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactor) GetRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.contract.Transact(opts, "getRandomness", requestID)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorSession) GetRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.GetRandomness(&_VRFBeaconCoordinator.TransactOpts, requestID)
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorTransactorSession) GetRandomness(requestID *big.Int) (*types.Transaction, error) {
	return _VRFBeaconCoordinator.Contract.GetRandomness(&_VRFBeaconCoordinator.TransactOpts, requestID)
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
	Transmitter       common.Address
	JuelsPerFeeCoin   *big.Int
	ConfigDigest      [32]byte
	EpochAndRound     *big.Int
	OutputsServed     []VRFBeaconReportOutputServed
	Raw               types.Log
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*VRFBeaconCoordinatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCoordinatorNewTransmissionIterator{contract: _VRFBeaconCoordinator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_VRFBeaconCoordinator *VRFBeaconCoordinatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _VRFBeaconCoordinator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
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
	return common.HexToHash("0x7484067466b4f2452757769a8dc9a8b41497154367515673c79386f9f0b74f16")
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

	GetRandomness(opts *bind.TransactOpts, requestID *big.Int) (*types.Transaction, error)

	KeyGenerated(opts *bind.TransactOpts, kd KeyDataStructKeyData) (*types.Transaction, error)

	NewKeyRequested(opts *bind.TransactOpts) (*types.Transaction, error)

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

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*VRFBeaconCoordinatorNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFBeaconCoordinatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error)

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
