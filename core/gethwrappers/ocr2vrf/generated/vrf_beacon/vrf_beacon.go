// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_beacon

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

type ECCArithmeticG1Point struct {
	P [2]*big.Int
}

type KeyDataStructKeyData struct {
	PublicKey []byte
	Hashes    [][32]byte
}

type VRFBeaconReportReport struct {
	Outputs            []VRFBeaconTypesVRFOutput
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	RecentBlockHeight  uint64
	RecentBlockHash    [32]byte
}

type VRFBeaconTypesCallback struct {
	RequestID      *big.Int
	NumWords       uint16
	Requester      common.Address
	Arguments      []byte
	GasAllowance   *big.Int
	SubID          *big.Int
	GasPrice       *big.Int
	WeiPerUnitLink *big.Int
}

type VRFBeaconTypesCostedCallback struct {
	Callback VRFBeaconTypesCallback
	Price    *big.Int
}

type VRFBeaconTypesOutputServed struct {
	Height            uint64
	ConfirmationDelay *big.Int
	ProofG1X          *big.Int
	ProofG1Y          *big.Int
}

type VRFBeaconTypesVRFOutput struct {
	BlockHeight       uint64
	ConfirmationDelay *big.Int
	VrfOutput         ECCArithmeticG1Point
	Callbacks         []VRFBeaconTypesCostedCallback
}

var VRFBeaconMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"contractIVRFCoordinatorProducerAPI\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"contractDKG\",\"name\":\"keyProvider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"keyID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualLength\",\"type\":\"uint256\"}],\"name\":\"CalldataLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotAcceptPayeeship\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"providedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"onchainHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorWrong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numTransmitters\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numPayees\",\"type\":\"uint256\"}],\"name\":\"IncorrectNumberOfPayees\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expectedNumSignatures\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"actualBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"keyProvider\",\"type\":\"address\"}],\"name\":\"KeyInfoMustComeFromProvider\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeftGasExceedsInitialGas\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeOwnerOrBillingAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"numFaultyOracles\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"}],\"name\":\"NumberOfFaultyOraclesTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"expectedLength\",\"type\":\"uint256\"}],\"name\":\"OnchainConfigHasWrongLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCurrentPayee\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"existingPayee\",\"type\":\"address\"}],\"name\":\"PayeeAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"repeatedSignerAddress\",\"type\":\"address\"}],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"repeatedTransmitterAddress\",\"type\":\"address\"}],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"reportHash\",\"type\":\"bytes32\"}],\"name\":\"SeenReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numTransmitters\",\"type\":\"uint256\"}],\"name\":\"SignersTransmittersMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"maxOracles\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"providedOracles\",\"type\":\"uint256\"}],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"ocrVersion\",\"type\":\"uint64\"}],\"name\":\"UnknownConfigVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maximumGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"observationPayment\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"transmissionPayment\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconTypes.VRFOutput[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"recentBlockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structVRFBeaconReport.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"exposeType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"maximumGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"observationPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"transmissionPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_coordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorProducerAPI\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structKeyDataStruct.KeyData\",\"name\":\"kd\",\"type\":\"tuple\"}],\"name\":\"keyGenerated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newKeyRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keyID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keyProvider\",\"outputs\":[{\"internalType\":\"contractDKG\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_provingKeyHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"maximumGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"observationPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"transmissionPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162004796380380620047968339810160408190526200003491620001c7565b8181858581813380600081620000915760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c457620000c48162000103565b5050506001600160a01b03918216608052811660a052601380546001600160a01b03191695909116949094179093555060c05250620002219350505050565b336001600160a01b038216036200015d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000088565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620001c457600080fd5b50565b60008060008060808587031215620001de57600080fd5b8451620001eb81620001ae565b6020860151909450620001fe81620001ae565b60408601519093506200021181620001ae565b6060959095015193969295505050565b60805160a05160c0516144b9620002dd60003960006103d50152600081816102f001528181610e3601528181610ef901528181610fac0152818161107f0152818161110101528181611c0f01528181611cd201528181611e32015281816121b80152818161228a015281816126ab0152818161276e01528181612c3e015261304801526000818161029901528181610eca0152818161105001528181611ca301528181611eb10152818161273f0152612b6701526144b96000f3fe608060405234801561001057600080fd5b506004361061018f5760003560e01c8063afcb95d7116100e4578063d09dc33911610092578063d09dc339146103f7578063d57fc45a146103ff578063e3d0e71214610408578063e4902f821461041b578063e53bbc9a14610443578063eb5dcd6c14610456578063f2fde38b14610469578063fbffd2c11461047c57600080fd5b8063afcb95d714610349578063b121e14714610373578063b1dc65a414610386578063bf2732c714610399578063c1075329146103ac578063c4c92b37146103bf578063cc31f7dd146103d057600080fd5b806379ba50971161014157806379ba50971461028c5780637d253aff1461029457806381ff7048146102bb5780638a1b1772146102eb5780638ac28d5a146103125780638da5cb5b146103255780639c849b301461033657600080fd5b80630eafb25b1461019457806310c29dbc146101ba578063181f5a77146101cd57806329937268146101fe5780632f7527cc1461024a57806355e48749146102645780635f27026f1461026c575b600080fd5b6101a76101a23660046131f5565b61048f565b6040519081526020015b60405180910390f35b6101cb6101c8366004613212565b50565b005b604080518082018252600f81526e0565246426561636f6e20312e302e3608c1b602082015290516101b1919061329c565b6002546003546040516101b192600160281b81046001600160401b0390811693600160681b8304821693600160a81b90930482169291811691600160401b90910462ffffff16906132af565b610252600881565b60405160ff90911681526020016101b1565b6101cb610587565b60135461027f906001600160a01b031681565b6040516101b191906132e6565b6101cb6105c8565b61027f7f000000000000000000000000000000000000000000000000000000000000000081565b6004546005546040805163ffffffff808516825264010000000090940490931660208401528201526060016101b1565b61027f7f000000000000000000000000000000000000000000000000000000000000000081565b6101cb6103203660046131f5565b610672565b6000546001600160a01b031661027f565b6101cb610344366004613345565b6106b5565b6005546006546040805160008152602081019390935263ffffffff909116908201526060016101b1565b6101cb6103813660046131f5565b610846565b6101cb6103943660046133f1565b6108f3565b6101cb6103a7366004613653565b610d19565b6101cb6103ba36600461371f565b610d7a565b6012546001600160a01b031661027f565b6101a77f000000000000000000000000000000000000000000000000000000000000000081565b6101a761104b565b6101a760145481565b6101cb610416366004613781565b6111ac565b61042e6104293660046131f5565b61186d565b60405163ffffffff90911681526020016101b1565b6101cb61045136600461387f565b611921565b6101cb6104643660046138f0565b611aa4565b6101cb6104773660046131f5565b611b7f565b6101cb61048a3660046131f5565b611b90565b6001600160a01b03811660009081526008602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b031691810191909152906104f15750600092915050565b6002546020820151600091600160a81b90046001600160401b031690600c9060ff16601f811061052357610523613929565b6008810491909101546002546105549260071660040261010090810a90920463ffffffff9081169290910416613955565b63ffffffff166105649190613979565b905081604001516001600160601b03168161057f9190613990565b949350505050565b6013546001600160a01b03163381146105c057338160405163292f4fb560e01b81526004016105b79291906139a3565b60405180910390fd5b506000601455565b6001546001600160a01b0316331461061b5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064016105b7565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b038181166000908152601060205260409020541633146106ac57604051633738e30960e21b815260040160405180910390fd5b6101c881611ba1565b6106bd611f30565b8281146106e7576040516336d2045960e01b815260048101849052602481018290526044016105b7565b60005b8381101561083f57600085858381811061070657610706613929565b905060200201602081019061071b91906131f5565b9050600084848481811061073157610731613929565b905060200201602081019061074691906131f5565b6001600160a01b038084166000908152601060205260409020549192501680158015816107855750826001600160a01b0316826001600160a01b031614155b156107a75783826040516375ef8bab60e11b81526004016105b79291906139a3565b6001600160a01b03848116600090815260106020526040902080546001600160a01b0319168583169081179091559083161461082857826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505050508080610837906139bd565b9150506106ea565b5050505050565b6001600160a01b0381811660009081526011602052604090205416331461088057604051639d12ec4f60e01b815260040160405180910390fd5b6001600160a01b0381811660008181526010602090815260408083208054336001600160a01b031980831682179093556011909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a6040805160e08101825260025460ff8082168352610100820463ffffffff16602080850191909152600160281b83046001600160401b0390811685870152600160681b840481166060860152600160a81b9093048316608085015260035492831660a0850152600160401b90920462ffffff1660c08401523360009081526008835293909320549394509092908c013591166109a757336040516358e0fb4760e11b81526004016105b791906132e6565b6005548b35146109d8576005546040516324f7d61360e21b815260048101919091528b3560248201526044016105b7565b6109e68a8a8a8a8a8a611f85565b81516109f39060016139d6565b60ff1687141580610a045750868514155b15610a43578151610a169060016139d6565b60405163fc33647560e01b815260ff909116600482015260248101889052604481018690526064016105b7565b60008a8a604051610a559291906139ef565b604051908190038120610a6c918e906020016139ff565b60408051601f19818403018152828252805160209182012083830190925260008084529083018190529092509060005b8a811015610be55760006001858a8460208110610abb57610abb613929565b610ac891901a601b6139d6565b8f8f86818110610ada57610ada613929565b905060200201358e8e87818110610af357610af3613929565b9050602002013560405160008152602001604052604051610b30949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610b52573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526009602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955092509050610bbe578060405163107dba7760e11b81526004016105b791906132e6565b826020015160080260ff166001901b84019350508080610bdd906139bd565b915050610a9c565b5081827e010101010101010101010101010101010101010101010101010101010101011614610c2757604051636081df1760e11b815260040160405180910390fd5b50505060008a8a604051610c3c9291906139ef565b604080519182900390912060008181526007602052919091205490915060ff1615610c7d57604051630b8d39d560e01b8152600481018290526024016105b7565b6000908152600760205260408120805460ff1916600117905580610cdf848e836020020135858f8f8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611ffc92505050565b6006805463ffffffff191663ffffffff600888901c161790559092509050610d0a8483838833612392565b50505050505050505050505050565b6013546001600160a01b0316338114610d4957338160405163292f4fb560e01b81526004016105b79291906139a3565b8151604051610d5b9190602001613a13565b60408051601f1981840301815291905280516020909101206014555050565b6000546001600160a01b03163314801590610e085750601254604051630d629b5f60e31b81526001600160a01b0390911690636b14daf890610dc59033906000903690600401613a58565b602060405180830381865afa158015610de2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e069190613a7d565b155b15610e2657604051631809d98560e31b815260040160405180910390fd5b6000610e306124a3565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663597d2f3c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e92573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610eb69190613a9f565b90506000610ec48284613990565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a082317f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b8152600401610f3491906132e6565b602060405180830381865afa158015610f51573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f759190613a9f565b905081811015610fa25760405163cf47918160e01b815260048101829052602481018390526044016105b7565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663f99b1d6887610fe5610fdf8686613ab8565b89612660565b6040516001600160e01b031960e085901b1681526001600160a01b0390921660048301526024820152604401600060405180830381600087803b15801561102b57600080fd5b505af115801561103f573d6000803e3d6000fd5b50505050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a082317f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b81526004016110ba91906132e6565b602060405180830381865afa1580156110d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110fb9190613a9f565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663597d2f3c6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561115d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111819190613a9f565b9050600061118d6124a3565b90508161119a8285613acb565b6111a49190613acb565b935050505090565b888787601f8311156111db57604051631013f88560e31b8152601f6004820152602481018490526044016105b7565b81831461120557604051632622820160e21b815260048101849052602481018390526044016105b7565b611210816003613aeb565b60ff16831161123e57604051631fb53b6f60e31b815260ff82166004820152602481018490526044016105b7565b61124a8160ff1661267a565b611252611f30565b60006040518060c001604052808f8f80806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505081526020018d8d8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050509082525060ff8c1660208083019190915260408051601f8d0183900483028101830182528c8152920191908c908c90819084018382808284376000920191909152505050908252506001600160401b03891660208083019190915260408051601f8a0183900483028101830182528981529201919089908990819084018382808284376000920191909152505050915250905061137361269b565b600a5460005b8181101561141d576000600a828154811061139657611396613929565b6000918252602082200154600b80546001600160a01b03909216935090849081106113c3576113c3613929565b60009182526020808320909101546001600160a01b039485168352600982526040808420805461ffff1916905594168252600890529190912080546001600160701b03191690555080611415816139bd565b915050611379565b5061142a600a60006130b5565b611436600b60006130b5565b60005b82515181101561168d57600960008460000151838151811061145d5761145d613929565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156114bf57825180518290811061149d5761149d613929565b6020026020010151604051633a28fc1f60e11b81526004016105b791906132e6565b604080518082019091526001815260ff8216602082015283518051600991600091859081106114f0576114f0613929565b6020908102919091018101516001600160a01b03168252818101929092526040016000908120835181549484015161ffff1990951690151561ff0019161761010060ff9095169490940293909317909255840151805160089291908490811061155b5761155b613929565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156115bf578260200151818151811061159d5761159d613929565b602002602001015160405163e8d2989960e01b81526004016105b791906132e6565b60405180606001604052806001151581526020018260ff16815260200160006001600160601b0316815250600860008560200151848151811061160457611604613929565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020835181549385015194909201516001600160601b0316620100000262010000600160701b031960ff959095166101000261ff00199315159390931661ffff1990941693909317919091179290921617905580611685816139bd565b915050611439565b50815180516116a491600a916020909101906130d3565b5060208083015180516116bb92600b9201906130d3565b5060408201516002805460ff191660ff909216919091179055600454640100000000900463ffffffff166116ed612cb5565b6004805463ffffffff9283166401000000000267ffffffff0000000019821681179092556000926117249281169116176001613b07565b905080600460006101000a81548163ffffffff021916908363ffffffff160217905550600061177846308463ffffffff16886000015189602001518a604001518b606001518c608001518d60a00151612d3f565b9050806005819055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e058360055484886000015189602001518a604001518b606001518c608001518d60a001516040516117da99989796959493929190613b68565b60405180910390a1600254610100900463ffffffff1660005b86515181101561184d5781600c82601f811061181157611811613929565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff1602179055508080611845906139bd565b9150506117f3565b506118588e8e612d9a565b50505050505050505050505050505050505050565b6001600160a01b03811660009081526008602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b031691810191909152906118cf5750600092915050565b600c816020015160ff16601f81106118e9576118e9613929565b60088104919091015460025461191a9260071660040261010090810a90920463ffffffff9081169290910416613955565b9392505050565b6000546001600160a01b031633148015906119af5750601254604051630d629b5f60e31b81526001600160a01b0390911690636b14daf89061196c9033906000903690600401613a58565b602060405180830381865afa158015611989573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119ad9190613a7d565b155b156119cd57604051631809d98560e31b815260040160405180910390fd5b6119d561269b565b6002805465010000000000600160a81b031916600160281b6001600160401b038881169190910267ffffffffffffffff60681b191691909117600160681b878316021767ffffffffffffffff60a81b1916600160a81b8683160217909155600380549184166affffffffffffffffffffff1990921691909117600160401b62ffffff8416021790556040517f49275ddcdfc9c0519b3d094308c8bf675f06070a754ce90c152163cb6e66e8a090611a9590879087908790879087906132af565b60405180910390a15050505050565b6001600160a01b03828116600090815260106020526040902054163314611ade57604051633738e30960e21b815260040160405180910390fd5b6001600160a01b0381163303611b0757604051631670f44760e31b815260040160405180910390fd5b6001600160a01b03808316600090815260116020526040902080548383166001600160a01b031982168117909255909116908114611b7a576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b611b87611f30565b6101c881612da8565b611b98611f30565b6101c881612e4b565b6001600160a01b0381166000908152600860209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046001600160601b031692810192909252611bfe575050565b6000611c098361048f565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663597d2f3c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611c6b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c8f9190613a9f565b90506000611c9d8284613990565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a082317f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b8152600401611d0d91906132e6565b602060405180830381865afa158015611d2a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d4e9190613a9f565b905081811015611d7b5760405163cf47918160e01b815260048101829052602481018390526044016105b7565b8315611f28576001600160a01b038681166000908152601060209081526040909120546002549188015192169161010090910463ffffffff1690600c9060ff16601f8110611dcb57611dcb613929565b6008808204909201805463ffffffff9485166004600790941684026101000a90810295021916939093179092556001600160a01b03808a166000908152602092909252604091829020805462010000600160701b03191690559051631f3363ad60e31b81527f00000000000000000000000000000000000000000000000000000000000000009091169163f99b1d6891611e7d9185918a91016001600160a01b03929092168252602082015260400190565b600060405180830381600087803b158015611e9757600080fd5b505af1158015611eab573d6000803e3d6000fd5b505050507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316816001600160a01b0316886001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c88604051611f1e91815260200190565b60405180910390a4505b505050505050565b6000546001600160a01b03163314611f835760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016105b7565b565b6000611f92826020613979565b611f9d856020613979565b611fa988610144613990565b611fb39190613990565b611fbd9190613990565b611fc8906000613990565b9050368114611ff357604051637bdca78560e11b8152600481018290523660248201526044016105b7565b50505050505050565b6000806000838060200190518101906120159190613e0a565b60208801805191925061202782613fef565b63ffffffff1663ffffffff168152505086600260008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160056101000a8154816001600160401b0302191690836001600160401b03160217905550606082015181600001600d6101000a8154816001600160401b0302191690836001600160401b0316021790555060808201518160000160156101000a8154816001600160401b0302191690836001600160401b0316021790555060a08201518160010160006101000a8154816001600160401b0302191690836001600160401b0316021790555060c08201518160010160086101000a81548162ffffff021916908362ffffff160217905550905050600061216d8260600151612ebb565b9050808260800151146121b6576080820151606083015160405163aed0afe560e01b81526004810192909252602482018390526001600160401b031660448201526064016105b7565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663483af703836000015184602001518560400151866060015187608001516040518663ffffffff1660e01b815260040161221e959493929190614116565b6000604051808303816000875af115801561223d573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526122659190810190614216565b5060408281015190516302fa566360e11b81526001600160401b0390911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906305f4acc690602401600060405180830381600087803b1580156122d657600080fd5b505af11580156122ea573d6000803e3d6000fd5b505050508564ffffffffff16886020015163ffffffff167f27bf3f1077f091da6885751ba10f5775d06657fd59e47a6ab1f7635e5a115afe33856020015186604001518c60405161237194939291906001600160a01b039490941684526001600160c01b039290921660208401526001600160401b03166040830152606082015260800190565b60405180910390a38160200151826040015193509350505094509492505050565b60006123bd3a6001600160401b038616156123ad57856123b3565b87606001515b8860400151612f7b565b90506010360260005a905060006123e68663ffffffff1685858c60c0015162ffffff1686612fc9565b90506000670de0b6b3a76400006001600160c01b038a1683026001600160a01b03881660009081526008602052604090205460a08d01519290910492506201000090046001600160601b03908116916001600160401b031682840101908116821115612458575050505050505061083f565b6001600160a01b038816600090815260086020526040902080546001600160601b03909216620100000262010000600160701b03199092169190911790555050505050505050505050565b600080600b8054806020026020016040519081016040528092919081815260200182805480156124fc57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116124de575b50508351600254604080516103e0810191829052969750919561010090910463ffffffff169450600093509150600c90601f908285855b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116125335790505050505050905060005b838110156125c6578181601f811061259357612593613929565b60200201516125a29084613955565b6125b29063ffffffff1687613990565b9550806125be816139bd565b915050612579565b506002546125e490600160a81b90046001600160401b031686613979565b945060005b83811015612658576008600086838151811061260757612607613929565b6020908102919091018101516001600160a01b0316825281019190915260400160002054612644906201000090046001600160601b031687613990565b955080612650816139bd565b9150506125e9565b505050505090565b600081831015612671575081612674565b50805b92915050565b806000036101c8576040516373bedd2b60e11b815260040160405180910390fd5b60006126a56124a3565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663597d2f3c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612707573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061272b9190613a9f565b905060006127398284613990565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a082317f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b81526004016127a991906132e6565b602060405180830381865afa1580156127c6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127ea9190613a9f565b9050818110156128175760405163cf47918160e01b815260048101829052602481018390526044016105b7565b600254604080516103e081019182905261010090920463ffffffff1691600091600c90601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411612841579050505050505090506000600b8054806020026020016040519081016040528092919081815260200182805480156128dc57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116128be575b505050505090506000815190506000816001600160401b03811115612903576129036134a7565b60405190808252806020026020018201604052801561292c578160200160208202803683370190505b5090506000826001600160401b03811115612949576129496134a7565b604051908082528060200260200182016040528015612972578160200160208202803683370190505b5090506000805b84811015612c025760006008600088848151811061299957612999613929565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160029054906101000a90046001600160601b03166001600160601b031690506000600860008985815181106129fb576129fb613929565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160026101000a8154816001600160601b0302191690836001600160601b0316021790555060008883601f8110612a5e57612a5e613929565b6020020151600254908b0363ffffffff169150600160a81b90046001600160401b0316810282018015612bf7576000601060008b8781518110612aa357612aa3613929565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b0316905080888781518110612af457612af4613929565b60200260200101906001600160a01b031690816001600160a01b03168152505081878781518110612b2757612b27613929565b6020026020010181815250508b8b86601f8110612b4657612b46613929565b602002019063ffffffff16908163ffffffff168152505085806001019650507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316816001600160a01b03168b8781518110612bab57612bab613929565b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c85604051612bed91815260200190565b60405180910390a4505b505050600101612979565b5081518114612c12578082528083525b612c1f600c87601f613138565b50815115612ca8576040516373433a2f60e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906373433a2f90612c759086908690600401614307565b600060405180830381600087803b158015612c8f57600080fd5b505af1158015612ca3573d6000803e3d6000fd5b505050505b5050505050505050505050565b60004661a4b1811480612cca575062066eed81145b15612d385760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612d0e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d329190613a9f565b91505090565b4391505090565b6000808a8a8a8a8a8a8a8a8a604051602001612d639998979695949392919061435e565b60408051601f1981840301815291905280516020909101206001600160f01b0316600160f01b179150509998505050505050505050565b612da48282612ff8565b5050565b336001600160a01b03821603612dfa5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b60448201526064016105b7565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6012546001600160a01b039081169082168114612da457601280546001600160a01b0319166001600160a01b0384161790556040517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891290612eaf90839085906139a3565b60405180910390a15050565b60004661a4b1811480612ed0575062066eed81145b15612f6c57610100836001600160401b0316612eea612cb5565b612ef49190613ab8565b1115612f035750600092915050565b6040516315a03d4160e11b81526001600160401b0384166004820152606490632b407a8290602401602060405180830381865afa158015612f48573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061191a9190613a9f565b50506001600160401b03164090565b6000836001600160401b038416811015612fad57600285856001600160401b03160381612faa57612faa6142f1565b04015b612fc081846001600160401b0316612660565b95945050505050565b600081861015612fec57604051633fef97df60e01b815260040160405180910390fd5b50909303019091010290565b6101008181146130215782828260405163418a179b60e01b81526004016105b7939291906143e7565b600061302f8385018561440b565b9050604051638eef585f60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690638eef585f9061307d908490600401614475565b600060405180830381600087803b15801561309757600080fd5b505af11580156130ab573d6000803e3d6000fd5b5050505050505050565b50805460008255906000526020600020908101906101c891906131cb565b828054828255906000526020600020908101928215613128579160200282015b8281111561312857825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906130f3565b506131349291506131cb565b5090565b6004830191839082156131285791602002820160005b8382111561319257835183826101000a81548163ffffffff021916908363ffffffff160217905550926020019260040160208160030104928301926001030261314e565b80156131c25782816101000a81549063ffffffff0219169055600401602081600301049283019260010302613192565b50506131349291505b5b8082111561313457600081556001016131cc565b6001600160a01b03811681146101c857600080fd5b60006020828403121561320757600080fd5b813561191a816131e0565b60006020828403121561322457600080fd5b81356001600160401b0381111561323a57600080fd5b820160a0818503121561191a57600080fd5b60005b8381101561326757818101518382015260200161324f565b50506000910152565b6000815180845261328881602086016020860161324c565b601f01601f19169290920160200192915050565b60208152600061191a6020830184613270565b6001600160401b03958616815293851660208501529184166040840152909216606082015262ffffff909116608082015260a00190565b6001600160a01b0391909116815260200190565b60008083601f84011261330c57600080fd5b5081356001600160401b0381111561332357600080fd5b6020830191508360208260051b850101111561333e57600080fd5b9250929050565b6000806000806040858703121561335b57600080fd5b84356001600160401b038082111561337257600080fd5b61337e888389016132fa565b9096509450602087013591508082111561339757600080fd5b506133a4878288016132fa565b95989497509550505050565b60008083601f8401126133c257600080fd5b5081356001600160401b038111156133d957600080fd5b60208301915083602082850101111561333e57600080fd5b60008060008060008060008060e0898b03121561340d57600080fd5b606089018a81111561341e57600080fd5b899850356001600160401b038082111561343757600080fd5b6134438c838d016133b0565b909950975060808b013591508082111561345c57600080fd5b6134688c838d016132fa565b909750955060a08b013591508082111561348157600080fd5b5061348e8b828c016132fa565b999c989b50969995989497949560c00135949350505050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156134df576134df6134a7565b60405290565b60405161010081016001600160401b03811182821017156134df576134df6134a7565b60405160a081016001600160401b03811182821017156134df576134df6134a7565b604051608081016001600160401b03811182821017156134df576134df6134a7565b604051602081016001600160401b03811182821017156134df576134df6134a7565b604051601f8201601f191681016001600160401b0381118282101715613596576135966134a7565b604052919050565b60006001600160401b038211156135b7576135b76134a7565b50601f01601f191660200190565b60006001600160401b038211156135de576135de6134a7565b5060051b60200190565b600082601f8301126135f957600080fd5b8135602061360e613609836135c5565b61356e565b82815260059290921b8401810191818101908684111561362d57600080fd5b8286015b848110156136485780358352918301918301613631565b509695505050505050565b6000602080838503121561366657600080fd5b82356001600160401b038082111561367d57600080fd5b908401906040828703121561369157600080fd5b6136996134bd565b8235828111156136a857600080fd5b8301601f810188136136b957600080fd5b80356136c76136098261359e565b81815289878385010111156136db57600080fd5b81878401888301376000878383010152808452505050838301358281111561370257600080fd5b61370e888286016135e8565b948201949094529695505050505050565b6000806040838503121561373257600080fd5b823561373d816131e0565b946020939093013593505050565b803560ff8116811461375c57600080fd5b919050565b6001600160401b03811681146101c857600080fd5b803561375c81613761565b60008060008060008060008060008060c08b8d0312156137a057600080fd5b8a356001600160401b03808211156137b757600080fd5b6137c38e838f016132fa565b909c509a5060208d01359150808211156137dc57600080fd5b6137e88e838f016132fa565b909a5098508891506137fc60408e0161374b565b975060608d013591508082111561381257600080fd5b61381e8e838f016133b0565b909750955085915061383260808e01613776565b945060a08d013591508082111561384857600080fd5b506138558d828e016133b0565b915080935050809150509295989b9194979a5092959850565b62ffffff811681146101c857600080fd5b600080600080600060a0868803121561389757600080fd5b85356138a281613761565b945060208601356138b281613761565b935060408601356138c281613761565b925060608601356138d281613761565b915060808601356138e28161386e565b809150509295509295909350565b6000806040838503121561390357600080fd5b823561390e816131e0565b9150602083013561391e816131e0565b809150509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b63ffffffff8281168282160390808211156139725761397261393f565b5092915050565b80820281158282048414176126745761267461393f565b808201808211156126745761267461393f565b6001600160a01b0392831681529116602082015260400190565b6000600182016139cf576139cf61393f565b5060010190565b60ff81811683821601908111156126745761267461393f565b8183823760009101908152919050565b828152606082602083013760800192915050565b60008251613a2581846020870161324c565b9190910192915050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0384168152604060208201819052600090612fc09083018486613a2f565b600060208284031215613a8f57600080fd5b8151801515811461191a57600080fd5b600060208284031215613ab157600080fd5b5051919050565b818103818111156126745761267461393f565b81810360008312801583831316838312821617156139725761397261393f565b60ff81811683821602908116908181146139725761397261393f565b63ffffffff8181168382160190808211156139725761397261393f565b600081518084526020808501945080840160005b83811015613b5d5781516001600160a01b031687529582019590820190600101613b38565b509495945050505050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152613b988184018a613b24565b90508281036080840152613bac8189613b24565b905060ff871660a084015282810360c0840152613bc98187613270565b90506001600160401b03851660e0840152828103610100840152613bed8185613270565b9c9b505050505050505050505050565b805161375c81613761565b805165ffffffffffff8116811461375c57600080fd5b805161ffff8116811461375c57600080fd5b805161375c816131e0565b600082601f830112613c4c57600080fd5b8151613c5a6136098261359e565b818152846020838601011115613c6f57600080fd5b61057f82602083016020870161324c565b80516001600160601b038116811461375c57600080fd5b600082601f830112613ca857600080fd5b81516020613cb8613609836135c5565b82815260059290921b84018101918181019086841115613cd757600080fd5b8286015b848110156136485780516001600160401b0380821115613cfa57600080fd5b90880190601f196040838c0382011215613d1357600080fd5b613d1b6134bd565b8784015183811115613d2c57600080fd5b8401610100818e0384011215613d4157600080fd5b613d496134e5565b9250613d56898201613c08565b8352613d6460408201613c1e565b89840152613d7460608201613c30565b6040840152608081015184811115613d8b57600080fd5b613d998e8b83850101613c3b565b606085015250613dab60a08201613c80565b608084015260c081015160a084015260e081015160c084015261010081015160e084015250818152613ddf60408501613c80565b818901528652505050918301918301613cdb565b80516001600160c01b038116811461375c57600080fd5b600060208284031215613e1c57600080fd5b81516001600160401b0380821115613e3357600080fd5b9083019060a08286031215613e4757600080fd5b613e4f613508565b825182811115613e5e57600080fd5b8301601f81018713613e6f57600080fd5b8051613e7d613609826135c5565b8082825260208201915060208360051b850101925089831115613e9f57600080fd5b602084015b83811015613f9f57805187811115613ebb57600080fd5b850160a0818d03601f19011215613ed157600080fd5b613ed961352a565b6020820151613ee781613761565b81526040820151613ef78161386e565b60208201526040828e03605f19011215613f1057600080fd5b613f1861354c565b8d607f840112613f2757600080fd5b613f2f6134bd565b808f60a086011115613f4057600080fd5b606085015b60a08601811015613f60578051835260209283019201613f45565b50825250604082015260a082015189811115613f7b57600080fd5b613f8a8e602083860101613c97565b60608301525084525060209283019201613ea4565b50845250613fb291505060208401613df3565b6020820152613fc360408401613bfd565b6040820152613fd460608401613bfd565b60608201526080830151608082015280935050505092915050565b600063ffffffff8083168181036140085761400861393f565b6001019392505050565b600081518084526020808501808196508360051b8101915082860160005b8581101561410957828403895281516040815181875265ffffffffffff81511682880152878101516060614069818a018361ffff169052565b928201519260809150614086898301856001600160a01b03169052565b8083015193505061010060a081818b01526140a56101408b0186613270565b9284015192945060c06140c28b8201856001600160601b03169052565b9084015160e08b81019190915290840151918a01919091529091015161012088015250908601516001600160601b0316948601949094529784019790840190600101614030565b5091979650505050505050565b600060a080830181845280895180835260c08601915060c08160051b87010192506020808c016000805b848110156141c65789870360bf19018652825180516001600160401b031688528481015162ffffff1685890152604080820151519084908a015b600282101561419957825181529187019160019190910190870161417a565b50505060600151608088018990526141b3888a0182614012565b9750509483019491830191600101614140565b5050508395506141e08188018c6001600160c01b03169052565b50505050506141fa60408301866001600160401b03169052565b6001600160401b03939093166060820152608001529392505050565b6000602080838503121561422957600080fd5b82516001600160401b0381111561423f57600080fd5b8301601f8101851361425057600080fd5b805161425e613609826135c5565b81815260079190911b8201830190838101908783111561427d57600080fd5b928401925b828410156142e6576080848903121561429b5760008081fd5b6142a361352a565b84516142ae81613761565b8152848601516142bd8161386e565b818701526040858101519082015260608086015190820152825260809093019290840190614282565b979650505050505050565b634e487b7160e01b600052601260045260246000fd5b60408152600061431a6040830185613b24565b82810360208481019190915284518083528582019282019060005b8181101561435157845183529383019391830191600101614335565b5090979650505050505050565b8981526001600160a01b03891660208201526001600160401b038881166040830152610120606083018190526000916143998483018b613b24565b915083820360808501526143ad828a613b24565b915060ff881660a085015283820360c08501526143ca8288613270565b90861660e08501528381036101008501529050613bed8185613270565b6040815260006143fb604083018587613a2f565b9050826020830152949350505050565b600061010080838503121561441f57600080fd5b83601f84011261442e57600080fd5b6144366134e5565b90830190808583111561444857600080fd5b845b8381101561446b57803561445d8161386e565b83526020928301920161444a565b5095945050505050565b6101008101818360005b60088110156144a357815162ffffff1683526020928301929091019060010161447f565b5050509291505056fea164736f6c6343000813000a",
}

var VRFBeaconABI = VRFBeaconMetaData.ABI

var VRFBeaconBin = VRFBeaconMetaData.Bin

func DeployVRFBeacon(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, coordinator common.Address, keyProvider common.Address, keyID [32]byte) (common.Address, *types.Transaction, *VRFBeacon, error) {
	parsed, err := VRFBeaconMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFBeaconBin), backend, link, coordinator, keyProvider, keyID)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFBeacon{VRFBeaconCaller: VRFBeaconCaller{contract: contract}, VRFBeaconTransactor: VRFBeaconTransactor{contract: contract}, VRFBeaconFilterer: VRFBeaconFilterer{contract: contract}}, nil
}

type VRFBeacon struct {
	address common.Address
	abi     abi.ABI
	VRFBeaconCaller
	VRFBeaconTransactor
	VRFBeaconFilterer
}

type VRFBeaconCaller struct {
	contract *bind.BoundContract
}

type VRFBeaconTransactor struct {
	contract *bind.BoundContract
}

type VRFBeaconFilterer struct {
	contract *bind.BoundContract
}

type VRFBeaconSession struct {
	Contract     *VRFBeacon
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFBeaconCallerSession struct {
	Contract *VRFBeaconCaller
	CallOpts bind.CallOpts
}

type VRFBeaconTransactorSession struct {
	Contract     *VRFBeaconTransactor
	TransactOpts bind.TransactOpts
}

type VRFBeaconRaw struct {
	Contract *VRFBeacon
}

type VRFBeaconCallerRaw struct {
	Contract *VRFBeaconCaller
}

type VRFBeaconTransactorRaw struct {
	Contract *VRFBeaconTransactor
}

func NewVRFBeacon(address common.Address, backend bind.ContractBackend) (*VRFBeacon, error) {
	abi, err := abi.JSON(strings.NewReader(VRFBeaconABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFBeacon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFBeacon{address: address, abi: abi, VRFBeaconCaller: VRFBeaconCaller{contract: contract}, VRFBeaconTransactor: VRFBeaconTransactor{contract: contract}, VRFBeaconFilterer: VRFBeaconFilterer{contract: contract}}, nil
}

func NewVRFBeaconCaller(address common.Address, caller bind.ContractCaller) (*VRFBeaconCaller, error) {
	contract, err := bindVRFBeacon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconCaller{contract: contract}, nil
}

func NewVRFBeaconTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFBeaconTransactor, error) {
	contract, err := bindVRFBeacon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconTransactor{contract: contract}, nil
}

func NewVRFBeaconFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFBeaconFilterer, error) {
	contract, err := bindVRFBeacon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconFilterer{contract: contract}, nil
}

func bindVRFBeacon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFBeaconMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFBeacon *VRFBeaconRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFBeacon.Contract.VRFBeaconCaller.contract.Call(opts, result, method, params...)
}

func (_VRFBeacon *VRFBeaconRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeacon.Contract.VRFBeaconTransactor.contract.Transfer(opts)
}

func (_VRFBeacon *VRFBeaconRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFBeacon.Contract.VRFBeaconTransactor.contract.Transact(opts, method, params...)
}

func (_VRFBeacon *VRFBeaconCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFBeacon.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFBeacon *VRFBeaconTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeacon.Contract.contract.Transfer(opts)
}

func (_VRFBeacon *VRFBeaconTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFBeacon.Contract.contract.Transact(opts, method, params...)
}

func (_VRFBeacon *VRFBeaconCaller) NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "NUM_CONF_DELAYS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFBeacon.Contract.NUMCONFDELAYS(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFBeacon.Contract.NUMCONFDELAYS(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) GetBilling(opts *bind.CallOpts) (GetBilling,

	error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "getBilling")

	outstruct := new(GetBilling)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPrice = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.ReasonableGasPrice = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.ObservationPayment = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.TransmissionPayment = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.AccountingGas = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFBeacon *VRFBeaconSession) GetBilling() (GetBilling,

	error) {
	return _VRFBeacon.Contract.GetBilling(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) GetBilling() (GetBilling,

	error) {
	return _VRFBeacon.Contract.GetBilling(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) GetBillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "getBillingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) GetBillingAccessController() (common.Address, error) {
	return _VRFBeacon.Contract.GetBillingAccessController(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) GetBillingAccessController() (common.Address, error) {
	return _VRFBeacon.Contract.GetBillingAccessController(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) ICoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "i_coordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) ICoordinator() (common.Address, error) {
	return _VRFBeacon.Contract.ICoordinator(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) ICoordinator() (common.Address, error) {
	return _VRFBeacon.Contract.ICoordinator(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) ILink(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "i_link")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) ILink() (common.Address, error) {
	return _VRFBeacon.Contract.ILink(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) ILink() (common.Address, error) {
	return _VRFBeacon.Contract.ILink(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_VRFBeacon *VRFBeaconSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _VRFBeacon.Contract.LatestConfigDetails(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _VRFBeacon.Contract.LatestConfigDetails(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFBeacon *VRFBeaconSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _VRFBeacon.Contract.LatestConfigDigestAndEpoch(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _VRFBeacon.Contract.LatestConfigDigestAndEpoch(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) LinkAvailableForPayment() (*big.Int, error) {
	return _VRFBeacon.Contract.LinkAvailableForPayment(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _VRFBeacon.Contract.LinkAvailableForPayment(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "oracleObservationCount", transmitterAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _VRFBeacon.Contract.OracleObservationCount(&_VRFBeacon.CallOpts, transmitterAddress)
}

func (_VRFBeacon *VRFBeaconCallerSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _VRFBeacon.Contract.OracleObservationCount(&_VRFBeacon.CallOpts, transmitterAddress)
}

func (_VRFBeacon *VRFBeaconCaller) OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "owedPayment", transmitterAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _VRFBeacon.Contract.OwedPayment(&_VRFBeacon.CallOpts, transmitterAddress)
}

func (_VRFBeacon *VRFBeaconCallerSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _VRFBeacon.Contract.OwedPayment(&_VRFBeacon.CallOpts, transmitterAddress)
}

func (_VRFBeacon *VRFBeaconCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) Owner() (common.Address, error) {
	return _VRFBeacon.Contract.Owner(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) Owner() (common.Address, error) {
	return _VRFBeacon.Contract.Owner(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) SKeyID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "s_keyID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) SKeyID() ([32]byte, error) {
	return _VRFBeacon.Contract.SKeyID(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) SKeyID() ([32]byte, error) {
	return _VRFBeacon.Contract.SKeyID(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) SKeyProvider(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "s_keyProvider")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) SKeyProvider() (common.Address, error) {
	return _VRFBeacon.Contract.SKeyProvider(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) SKeyProvider() (common.Address, error) {
	return _VRFBeacon.Contract.SKeyProvider(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "s_provingKeyHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) SProvingKeyHash() ([32]byte, error) {
	return _VRFBeacon.Contract.SProvingKeyHash(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) SProvingKeyHash() ([32]byte, error) {
	return _VRFBeacon.Contract.SProvingKeyHash(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFBeacon.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFBeacon *VRFBeaconSession) TypeAndVersion() (string, error) {
	return _VRFBeacon.Contract.TypeAndVersion(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconCallerSession) TypeAndVersion() (string, error) {
	return _VRFBeacon.Contract.TypeAndVersion(&_VRFBeacon.CallOpts)
}

func (_VRFBeacon *VRFBeaconTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "acceptOwnership")
}

func (_VRFBeacon *VRFBeaconSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFBeacon.Contract.AcceptOwnership(&_VRFBeacon.TransactOpts)
}

func (_VRFBeacon *VRFBeaconTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFBeacon.Contract.AcceptOwnership(&_VRFBeacon.TransactOpts)
}

func (_VRFBeacon *VRFBeaconTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_VRFBeacon *VRFBeaconSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.AcceptPayeeship(&_VRFBeacon.TransactOpts, transmitter)
}

func (_VRFBeacon *VRFBeaconTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.AcceptPayeeship(&_VRFBeacon.TransactOpts, transmitter)
}

func (_VRFBeacon *VRFBeaconTransactor) ExposeType(opts *bind.TransactOpts, arg0 VRFBeaconReportReport) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "exposeType", arg0)
}

func (_VRFBeacon *VRFBeaconSession) ExposeType(arg0 VRFBeaconReportReport) (*types.Transaction, error) {
	return _VRFBeacon.Contract.ExposeType(&_VRFBeacon.TransactOpts, arg0)
}

func (_VRFBeacon *VRFBeaconTransactorSession) ExposeType(arg0 VRFBeaconReportReport) (*types.Transaction, error) {
	return _VRFBeacon.Contract.ExposeType(&_VRFBeacon.TransactOpts, arg0)
}

func (_VRFBeacon *VRFBeaconTransactor) KeyGenerated(opts *bind.TransactOpts, kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "keyGenerated", kd)
}

func (_VRFBeacon *VRFBeaconSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _VRFBeacon.Contract.KeyGenerated(&_VRFBeacon.TransactOpts, kd)
}

func (_VRFBeacon *VRFBeaconTransactorSession) KeyGenerated(kd KeyDataStructKeyData) (*types.Transaction, error) {
	return _VRFBeacon.Contract.KeyGenerated(&_VRFBeacon.TransactOpts, kd)
}

func (_VRFBeacon *VRFBeaconTransactor) NewKeyRequested(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "newKeyRequested")
}

func (_VRFBeacon *VRFBeaconSession) NewKeyRequested() (*types.Transaction, error) {
	return _VRFBeacon.Contract.NewKeyRequested(&_VRFBeacon.TransactOpts)
}

func (_VRFBeacon *VRFBeaconTransactorSession) NewKeyRequested() (*types.Transaction, error) {
	return _VRFBeacon.Contract.NewKeyRequested(&_VRFBeacon.TransactOpts)
}

func (_VRFBeacon *VRFBeaconTransactor) SetBilling(opts *bind.TransactOpts, maximumGasPrice uint64, reasonableGasPrice uint64, observationPayment uint64, transmissionPayment uint64, accountingGas *big.Int) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "setBilling", maximumGasPrice, reasonableGasPrice, observationPayment, transmissionPayment, accountingGas)
}

func (_VRFBeacon *VRFBeaconSession) SetBilling(maximumGasPrice uint64, reasonableGasPrice uint64, observationPayment uint64, transmissionPayment uint64, accountingGas *big.Int) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetBilling(&_VRFBeacon.TransactOpts, maximumGasPrice, reasonableGasPrice, observationPayment, transmissionPayment, accountingGas)
}

func (_VRFBeacon *VRFBeaconTransactorSession) SetBilling(maximumGasPrice uint64, reasonableGasPrice uint64, observationPayment uint64, transmissionPayment uint64, accountingGas *big.Int) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetBilling(&_VRFBeacon.TransactOpts, maximumGasPrice, reasonableGasPrice, observationPayment, transmissionPayment, accountingGas)
}

func (_VRFBeacon *VRFBeaconTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

func (_VRFBeacon *VRFBeaconSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetBillingAccessController(&_VRFBeacon.TransactOpts, _billingAccessController)
}

func (_VRFBeacon *VRFBeaconTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetBillingAccessController(&_VRFBeacon.TransactOpts, _billingAccessController)
}

func (_VRFBeacon *VRFBeaconTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_VRFBeacon *VRFBeaconSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetConfig(&_VRFBeacon.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_VRFBeacon *VRFBeaconTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetConfig(&_VRFBeacon.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_VRFBeacon *VRFBeaconTransactor) SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "setPayees", transmitters, payees)
}

func (_VRFBeacon *VRFBeaconSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetPayees(&_VRFBeacon.TransactOpts, transmitters, payees)
}

func (_VRFBeacon *VRFBeaconTransactorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.SetPayees(&_VRFBeacon.TransactOpts, transmitters, payees)
}

func (_VRFBeacon *VRFBeaconTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFBeacon *VRFBeaconSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.TransferOwnership(&_VRFBeacon.TransactOpts, to)
}

func (_VRFBeacon *VRFBeaconTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.TransferOwnership(&_VRFBeacon.TransactOpts, to)
}

func (_VRFBeacon *VRFBeaconTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_VRFBeacon *VRFBeaconSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.TransferPayeeship(&_VRFBeacon.TransactOpts, transmitter, proposed)
}

func (_VRFBeacon *VRFBeaconTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.TransferPayeeship(&_VRFBeacon.TransactOpts, transmitter, proposed)
}

func (_VRFBeacon *VRFBeaconTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_VRFBeacon *VRFBeaconSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _VRFBeacon.Contract.Transmit(&_VRFBeacon.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_VRFBeacon *VRFBeaconTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _VRFBeacon.Contract.Transmit(&_VRFBeacon.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_VRFBeacon *VRFBeaconTransactor) WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "withdrawFunds", recipient, amount)
}

func (_VRFBeacon *VRFBeaconSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFBeacon.Contract.WithdrawFunds(&_VRFBeacon.TransactOpts, recipient, amount)
}

func (_VRFBeacon *VRFBeaconTransactorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFBeacon.Contract.WithdrawFunds(&_VRFBeacon.TransactOpts, recipient, amount)
}

func (_VRFBeacon *VRFBeaconTransactor) WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeacon.contract.Transact(opts, "withdrawPayment", transmitter)
}

func (_VRFBeacon *VRFBeaconSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.WithdrawPayment(&_VRFBeacon.TransactOpts, transmitter)
}

func (_VRFBeacon *VRFBeaconTransactorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _VRFBeacon.Contract.WithdrawPayment(&_VRFBeacon.TransactOpts, transmitter)
}

type VRFBeaconBillingAccessControllerSetIterator struct {
	Event *VRFBeaconBillingAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconBillingAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconBillingAccessControllerSet)
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
		it.Event = new(VRFBeaconBillingAccessControllerSet)
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

func (it *VRFBeaconBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*VRFBeaconBillingAccessControllerSetIterator, error) {

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconBillingAccessControllerSetIterator{contract: _VRFBeacon.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconBillingAccessControllerSet)
				if err := _VRFBeacon.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseBillingAccessControllerSet(log types.Log) (*VRFBeaconBillingAccessControllerSet, error) {
	event := new(VRFBeaconBillingAccessControllerSet)
	if err := _VRFBeacon.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconBillingSetIterator struct {
	Event *VRFBeaconBillingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconBillingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconBillingSet)
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
		it.Event = new(VRFBeaconBillingSet)
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

func (it *VRFBeaconBillingSetIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconBillingSet struct {
	MaximumGasPrice     uint64
	ReasonableGasPrice  uint64
	ObservationPayment  uint64
	TransmissionPayment uint64
	AccountingGas       *big.Int
	Raw                 types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterBillingSet(opts *bind.FilterOpts) (*VRFBeaconBillingSetIterator, error) {

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconBillingSetIterator{contract: _VRFBeacon.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconBillingSet) (event.Subscription, error) {

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconBillingSet)
				if err := _VRFBeacon.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseBillingSet(log types.Log) (*VRFBeaconBillingSet, error) {
	event := new(VRFBeaconBillingSet)
	if err := _VRFBeacon.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconConfigSetIterator struct {
	Event *VRFBeaconConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconConfigSet)
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
		it.Event = new(VRFBeaconConfigSet)
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

func (it *VRFBeaconConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconConfigSet struct {
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

func (_VRFBeacon *VRFBeaconFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFBeaconConfigSetIterator, error) {

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconConfigSetIterator{contract: _VRFBeacon.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconConfigSet)
				if err := _VRFBeacon.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseConfigSet(log types.Log) (*VRFBeaconConfigSet, error) {
	event := new(VRFBeaconConfigSet)
	if err := _VRFBeacon.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconNewTransmissionIterator struct {
	Event *VRFBeaconNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconNewTransmission)
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
		it.Event = new(VRFBeaconNewTransmission)
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

func (it *VRFBeaconNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconNewTransmission struct {
	AggregatorRoundId  uint32
	EpochAndRound      *big.Int
	Transmitter        common.Address
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	ConfigDigest       [32]byte
	Raw                types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFBeaconNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconNewTransmissionIterator{contract: _VRFBeacon.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFBeaconNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconNewTransmission)
				if err := _VRFBeacon.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseNewTransmission(log types.Log) (*VRFBeaconNewTransmission, error) {
	event := new(VRFBeaconNewTransmission)
	if err := _VRFBeacon.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconOraclePaidIterator struct {
	Event *VRFBeaconOraclePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconOraclePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconOraclePaid)
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
		it.Event = new(VRFBeaconOraclePaid)
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

func (it *VRFBeaconOraclePaidIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*VRFBeaconOraclePaidIterator, error) {

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

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconOraclePaidIterator{contract: _VRFBeacon.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *VRFBeaconOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconOraclePaid)
				if err := _VRFBeacon.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseOraclePaid(log types.Log) (*VRFBeaconOraclePaid, error) {
	event := new(VRFBeaconOraclePaid)
	if err := _VRFBeacon.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconOutputsServedIterator struct {
	Event *VRFBeaconOutputsServed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconOutputsServedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconOutputsServed)
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
		it.Event = new(VRFBeaconOutputsServed)
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

func (it *VRFBeaconOutputsServedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconOutputsServedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconOutputsServed struct {
	RecentBlockHeight  uint64
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	OutputsServed      []VRFBeaconTypesOutputServed
	Raw                types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterOutputsServed(opts *bind.FilterOpts) (*VRFBeaconOutputsServedIterator, error) {

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconOutputsServedIterator{contract: _VRFBeacon.contract, event: "OutputsServed", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *VRFBeaconOutputsServed) (event.Subscription, error) {

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconOutputsServed)
				if err := _VRFBeacon.contract.UnpackLog(event, "OutputsServed", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseOutputsServed(log types.Log) (*VRFBeaconOutputsServed, error) {
	event := new(VRFBeaconOutputsServed)
	if err := _VRFBeacon.contract.UnpackLog(event, "OutputsServed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconOwnershipTransferRequestedIterator struct {
	Event *VRFBeaconOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconOwnershipTransferRequested)
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
		it.Event = new(VRFBeaconOwnershipTransferRequested)
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

func (it *VRFBeaconOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconOwnershipTransferRequestedIterator{contract: _VRFBeacon.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconOwnershipTransferRequested)
				if err := _VRFBeacon.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFBeaconOwnershipTransferRequested, error) {
	event := new(VRFBeaconOwnershipTransferRequested)
	if err := _VRFBeacon.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconOwnershipTransferredIterator struct {
	Event *VRFBeaconOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconOwnershipTransferred)
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
		it.Event = new(VRFBeaconOwnershipTransferred)
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

func (it *VRFBeaconOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconOwnershipTransferredIterator{contract: _VRFBeacon.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconOwnershipTransferred)
				if err := _VRFBeacon.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseOwnershipTransferred(log types.Log) (*VRFBeaconOwnershipTransferred, error) {
	event := new(VRFBeaconOwnershipTransferred)
	if err := _VRFBeacon.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconPayeeshipTransferRequestedIterator struct {
	Event *VRFBeaconPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconPayeeshipTransferRequested)
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
		it.Event = new(VRFBeaconPayeeshipTransferRequested)
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

func (it *VRFBeaconPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*VRFBeaconPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconPayeeshipTransferRequestedIterator{contract: _VRFBeacon.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconPayeeshipTransferRequested)
				if err := _VRFBeacon.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParsePayeeshipTransferRequested(log types.Log) (*VRFBeaconPayeeshipTransferRequested, error) {
	event := new(VRFBeaconPayeeshipTransferRequested)
	if err := _VRFBeacon.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconPayeeshipTransferredIterator struct {
	Event *VRFBeaconPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconPayeeshipTransferred)
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
		it.Event = new(VRFBeaconPayeeshipTransferred)
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

func (it *VRFBeaconPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*VRFBeaconPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconPayeeshipTransferredIterator{contract: _VRFBeacon.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconPayeeshipTransferred)
				if err := _VRFBeacon.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParsePayeeshipTransferred(log types.Log) (*VRFBeaconPayeeshipTransferred, error) {
	event := new(VRFBeaconPayeeshipTransferred)
	if err := _VRFBeacon.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconRandomWordsFulfilledIterator struct {
	Event *VRFBeaconRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconRandomWordsFulfilled)
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
		it.Event = new(VRFBeaconRandomWordsFulfilled)
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

func (it *VRFBeaconRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconRandomWordsFulfilled struct {
	RequestIDs            []*big.Int
	SuccessfulFulfillment []byte
	TruncatedErrorData    [][]byte
	Raw                   types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFBeaconRandomWordsFulfilledIterator, error) {

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFBeaconRandomWordsFulfilledIterator{contract: _VRFBeacon.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFBeaconRandomWordsFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconRandomWordsFulfilled)
				if err := _VRFBeacon.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFBeaconRandomWordsFulfilled, error) {
	event := new(VRFBeaconRandomWordsFulfilled)
	if err := _VRFBeacon.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconRandomnessFulfillmentRequestedIterator struct {
	Event *VRFBeaconRandomnessFulfillmentRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconRandomnessFulfillmentRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconRandomnessFulfillmentRequested)
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
		it.Event = new(VRFBeaconRandomnessFulfillmentRequested)
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

func (it *VRFBeaconRandomnessFulfillmentRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconRandomnessFulfillmentRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconRandomnessFulfillmentRequested struct {
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  *big.Int
	NumWords               uint16
	GasAllowance           uint32
	GasPrice               *big.Int
	WeiPerUnitLink         *big.Int
	Arguments              []byte
	Raw                    types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFBeaconRandomnessFulfillmentRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconRandomnessFulfillmentRequestedIterator{contract: _VRFBeacon.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconRandomnessFulfillmentRequested)
				if err := _VRFBeacon.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseRandomnessFulfillmentRequested(log types.Log) (*VRFBeaconRandomnessFulfillmentRequested, error) {
	event := new(VRFBeaconRandomnessFulfillmentRequested)
	if err := _VRFBeacon.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFBeaconRandomnessRequestedIterator struct {
	Event *VRFBeaconRandomnessRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFBeaconRandomnessRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFBeaconRandomnessRequested)
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
		it.Event = new(VRFBeaconRandomnessRequested)
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

func (it *VRFBeaconRandomnessRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFBeaconRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFBeaconRandomnessRequested struct {
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  *big.Int
	NumWords               uint16
	Raw                    types.Log
}

func (_VRFBeacon *VRFBeaconFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFBeaconRandomnessRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFBeacon.contract.FilterLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFBeaconRandomnessRequestedIterator{contract: _VRFBeacon.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_VRFBeacon *VRFBeaconFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFBeacon.contract.WatchLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFBeaconRandomnessRequested)
				if err := _VRFBeacon.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
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

func (_VRFBeacon *VRFBeaconFilterer) ParseRandomnessRequested(log types.Log) (*VRFBeaconRandomnessRequested, error) {
	event := new(VRFBeaconRandomnessRequested)
	if err := _VRFBeacon.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetBilling struct {
	MaximumGasPrice     uint64
	ReasonableGasPrice  uint64
	ObservationPayment  uint64
	TransmissionPayment uint64
	AccountingGas       *big.Int
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

func (_VRFBeacon *VRFBeacon) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFBeacon.abi.Events["BillingAccessControllerSet"].ID:
		return _VRFBeacon.ParseBillingAccessControllerSet(log)
	case _VRFBeacon.abi.Events["BillingSet"].ID:
		return _VRFBeacon.ParseBillingSet(log)
	case _VRFBeacon.abi.Events["ConfigSet"].ID:
		return _VRFBeacon.ParseConfigSet(log)
	case _VRFBeacon.abi.Events["NewTransmission"].ID:
		return _VRFBeacon.ParseNewTransmission(log)
	case _VRFBeacon.abi.Events["OraclePaid"].ID:
		return _VRFBeacon.ParseOraclePaid(log)
	case _VRFBeacon.abi.Events["OutputsServed"].ID:
		return _VRFBeacon.ParseOutputsServed(log)
	case _VRFBeacon.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFBeacon.ParseOwnershipTransferRequested(log)
	case _VRFBeacon.abi.Events["OwnershipTransferred"].ID:
		return _VRFBeacon.ParseOwnershipTransferred(log)
	case _VRFBeacon.abi.Events["PayeeshipTransferRequested"].ID:
		return _VRFBeacon.ParsePayeeshipTransferRequested(log)
	case _VRFBeacon.abi.Events["PayeeshipTransferred"].ID:
		return _VRFBeacon.ParsePayeeshipTransferred(log)
	case _VRFBeacon.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFBeacon.ParseRandomWordsFulfilled(log)
	case _VRFBeacon.abi.Events["RandomnessFulfillmentRequested"].ID:
		return _VRFBeacon.ParseRandomnessFulfillmentRequested(log)
	case _VRFBeacon.abi.Events["RandomnessRequested"].ID:
		return _VRFBeacon.ParseRandomnessRequested(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFBeaconBillingAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912")
}

func (VRFBeaconBillingSet) Topic() common.Hash {
	return common.HexToHash("0x49275ddcdfc9c0519b3d094308c8bf675f06070a754ce90c152163cb6e66e8a0")
}

func (VRFBeaconConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (VRFBeaconNewTransmission) Topic() common.Hash {
	return common.HexToHash("0x27bf3f1077f091da6885751ba10f5775d06657fd59e47a6ab1f7635e5a115afe")
}

func (VRFBeaconOraclePaid) Topic() common.Hash {
	return common.HexToHash("0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c")
}

func (VRFBeaconOutputsServed) Topic() common.Hash {
	return common.HexToHash("0xf10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c44")
}

func (VRFBeaconOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFBeaconOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFBeaconPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (VRFBeaconPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (VRFBeaconRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7")
}

func (VRFBeaconRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0x24f0e469e0097d1e8d9975137f9f4dd17d2c1481b3a2f25f2382f51287eda1dc")
}

func (VRFBeaconRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xc3b31df4232b05afd212fc28027dae6fd6a81618c2a3116182cb57c7f0a3fd0a")
}

func (_VRFBeacon *VRFBeacon) Address() common.Address {
	return _VRFBeacon.address
}

type VRFBeaconInterface interface {
	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	GetBilling(opts *bind.CallOpts) (GetBilling,

		error)

	GetBillingAccessController(opts *bind.CallOpts) (common.Address, error)

	ICoordinator(opts *bind.CallOpts) (common.Address, error)

	ILink(opts *bind.CallOpts) (common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error)

	OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SKeyID(opts *bind.CallOpts) ([32]byte, error)

	SKeyProvider(opts *bind.CallOpts) (common.Address, error)

	SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	ExposeType(opts *bind.TransactOpts, arg0 VRFBeaconReportReport) (*types.Transaction, error)

	KeyGenerated(opts *bind.TransactOpts, kd KeyDataStructKeyData) (*types.Transaction, error)

	NewKeyRequested(opts *bind.TransactOpts) (*types.Transaction, error)

	SetBilling(opts *bind.TransactOpts, maximumGasPrice uint64, reasonableGasPrice uint64, observationPayment uint64, transmissionPayment uint64, accountingGas *big.Int) (*types.Transaction, error)

	SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*VRFBeaconBillingAccessControllerSetIterator, error)

	WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconBillingAccessControllerSet) (event.Subscription, error)

	ParseBillingAccessControllerSet(log types.Log) (*VRFBeaconBillingAccessControllerSet, error)

	FilterBillingSet(opts *bind.FilterOpts) (*VRFBeaconBillingSetIterator, error)

	WatchBillingSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconBillingSet) (event.Subscription, error)

	ParseBillingSet(log types.Log) (*VRFBeaconBillingSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFBeaconConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFBeaconConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFBeaconConfigSet, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFBeaconNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFBeaconNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*VRFBeaconNewTransmission, error)

	FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*VRFBeaconOraclePaidIterator, error)

	WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *VRFBeaconOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error)

	ParseOraclePaid(log types.Log) (*VRFBeaconOraclePaid, error)

	FilterOutputsServed(opts *bind.FilterOpts) (*VRFBeaconOutputsServedIterator, error)

	WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *VRFBeaconOutputsServed) (event.Subscription, error)

	ParseOutputsServed(log types.Log) (*VRFBeaconOutputsServed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFBeaconOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFBeaconOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFBeaconOwnershipTransferred, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*VRFBeaconPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*VRFBeaconPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*VRFBeaconPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *VRFBeaconPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*VRFBeaconPayeeshipTransferred, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFBeaconRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFBeaconRandomWordsFulfilled) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFBeaconRandomWordsFulfilled, error)

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFBeaconRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*VRFBeaconRandomnessFulfillmentRequested, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFBeaconRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFBeaconRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessRequested(log types.Log) (*VRFBeaconRandomnessRequested, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
