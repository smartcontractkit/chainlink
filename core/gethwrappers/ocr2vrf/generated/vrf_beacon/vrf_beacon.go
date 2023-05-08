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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"contractIVRFCoordinatorProducerAPI\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"contractDKG\",\"name\":\"keyProvider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"keyID\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualLength\",\"type\":\"uint256\"}],\"name\":\"CalldataLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotAcceptPayeeship\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"providedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"onchainHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorWrong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectNumberOfFaultyOracles\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numTransmitters\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numPayees\",\"type\":\"uint256\"}],\"name\":\"IncorrectNumberOfPayees\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expectedNumSignatures\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"rsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ssLength\",\"type\":\"uint256\"}],\"name\":\"IncorrectNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"actualBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayee\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"keyProvider\",\"type\":\"address\"}],\"name\":\"KeyInfoMustComeFromProvider\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeftGasExceedsInitialGas\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeOwnerOrBillingAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"numFaultyOracles\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"}],\"name\":\"NumberOfFaultyOraclesTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"expectedLength\",\"type\":\"uint256\"}],\"name\":\"OffchainConfigHasWrongLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"OnlyActiveSigners\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"OnlyActiveTransmitters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCurrentPayee\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"existingPayee\",\"type\":\"address\"}],\"name\":\"PayeeAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"repeatedSignerAddress\",\"type\":\"address\"}],\"name\":\"RepeatedSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"repeatedTransmitterAddress\",\"type\":\"address\"}],\"name\":\"RepeatedTransmitter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numSigners\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numTransmitters\",\"type\":\"uint256\"}],\"name\":\"SignersTransmittersMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"maxOracles\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"providedOracles\",\"type\":\"uint256\"}],\"name\":\"TooManyOracles\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"occVersion\",\"type\":\"uint64\"}],\"name\":\"UnknownConfigVersion\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maximumGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"observationPayment\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"transmissionPayment\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconTypes.VRFOutput[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"recentBlockHash\",\"type\":\"bytes32\"}],\"internalType\":\"structVRFBeaconReport.Report\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"exposeType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"maximumGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"observationPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"transmissionPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_coordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorProducerAPI\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"hashes\",\"type\":\"bytes32[]\"}],\"internalType\":\"structKeyDataStruct.KeyData\",\"name\":\"kd\",\"type\":\"tuple\"}],\"name\":\"keyGenerated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newKeyRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keyID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_keyProvider\",\"outputs\":[{\"internalType\":\"contractDKG\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_provingKeyHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"maximumGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"observationPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"transmissionPayment\",\"type\":\"uint64\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162004735380380620047358339810160408190526200003491620001c7565b8181858581813380600081620000915760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c457620000c48162000103565b5050506001600160a01b03918216608052811660a052601380546001600160a01b03191695909116949094179093555060145550620002219350505050565b336001600160a01b038216036200015d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000088565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620001c457600080fd5b50565b60008060008060808587031215620001de57600080fd5b8451620001eb81620001ae565b6020860151909450620001fe81620001ae565b60408601519093506200021181620001ae565b6060959095015193969295505050565b60805160a0516144a86200028d6000396000818161036901528181610f4901528181611019015281816110d401528181611dfc015281816121f6015281816122e201528181612b56015261300401526000818161031201528181611e7b015261272501526144a86000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c8063afcb95d711610104578063d09dc339116100a2578063e53bbc9a11610071578063e53bbc9a1461049e578063eb5dcd6c146104b1578063f2fde38b146104c4578063fbffd2c1146104d757600080fd5b8063d09dc33914610452578063d57fc45a1461045a578063e3d0e71214610463578063e4902f821461047657600080fd5b8063bf2732c7116100de578063bf2732c714610412578063c107532914610425578063c4c92b3714610438578063cc31f7dd1461044957600080fd5b8063afcb95d7146103c2578063b121e147146103ec578063b1dc65a4146103ff57600080fd5b806379ba5097116101715780638a1b17721161014b5780638a1b1772146103645780638ac28d5a1461038b5780638da5cb5b1461039e5780639c849b30146103af57600080fd5b806379ba5097146103055780637d253aff1461030d57806381ff70481461033457600080fd5b806329937268116101ad578063299372681461024c5780632f7527cc146102b857806355e48749146102d25780635f27026f146102da57600080fd5b80630eafb25b146101d457806310c29dbc146101fa578063181f5a771461020d575b600080fd5b6101e76101e2366004613179565b6104ea565b6040519081526020015b60405180910390f35b61020b610208366004613196565b50565b005b604080518082018252600f81527f565246426561636f6e20312e302e300000000000000000000000000000000000602082015290516101f1919061322d565b600254600354604080516a0100000000000000000000840467ffffffffffffffff9081168252600160901b90940484166020820152838316918101919091526801000000000000000082049092166060830152600160801b900462ffffff16608082015260a0016101f1565b6102c0600881565b60405160ff90911681526020016101f1565b61020b6105e1565b6013546102ed906001600160a01b031681565b6040516001600160a01b0390911681526020016101f1565b61020b61062b565b6102ed7f000000000000000000000000000000000000000000000000000000000000000081565b6004546005546040805163ffffffff808516825264010000000090940490931660208401528201526060016101f1565b6102ed7f000000000000000000000000000000000000000000000000000000000000000081565b61020b610399366004613179565b6106dc565b6000546001600160a01b03166102ed565b61020b6103bd36600461328c565b61071f565b6005546007546040805160008152602081019390935263ffffffff909116908201526060016101f1565b61020b6103fa366004613179565b6108f1565b61020b61040d36600461333a565b6109b7565b61020b6104203660046135a5565b610e1f565b61020b610433366004613672565b610e8d565b6012546001600160a01b03166102ed565b6101e760145481565b6101e76110cf565b6101e760155481565b61020b6104713660046136d5565b611173565b610489610484366004613179565b6118f1565b60405163ffffffff90911681526020016101f1565b61020b6104ac3660046137d4565b6119aa565b61020b6104bf366004613845565b611ba2565b61020b6104d2366004613179565b611c96565b61020b6104e5366004613179565b611ca7565b6001600160a01b03811660009081526008602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b0316918101919091529061054c5750600092915050565b600354602082015160009167ffffffffffffffff1690600c9060ff16601f81106105785761057861387e565b6008810491909101546002546105ae926007166004026101000a90910463ffffffff9081169166010000000000009004166138aa565b63ffffffff166105be91906138cf565b905081604001516001600160601b0316816105d991906138ee565b949350505050565b6013546001600160a01b03163381146106235760405163292f4fb560e01b81523360048201526001600160a01b03821660248201526044015b60405180910390fd5b506000601555565b6001546001600160a01b031633146106855760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161061a565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b0381811660009081526010602052604090205416331461071657604051633738e30960e21b815260040160405180910390fd5b61020881611cb8565b610727611ef6565b82811461076a576040517f36d20459000000000000000000000000000000000000000000000000000000008152600481018490526024810182905260440161061a565b60005b838110156108ea5760008585838181106107895761078961387e565b905060200201602081019061079e9190613179565b905060008484848181106107b4576107b461387e565b90506020020160208101906107c99190613179565b6001600160a01b038084166000908152601060205260409020549192501680158015816108085750826001600160a01b0316826001600160a01b031614155b15610852576040517febdf17560000000000000000000000000000000000000000000000000000000081526001600160a01b0380861660048301528316602482015260440161061a565b6001600160a01b03848116600090815260106020526040902080546001600160a01b031916858316908117909155908316146108d357826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050505080806108e290613906565b91505061076d565b5050505050565b6001600160a01b03818116600090815260116020526040902054163314610944576040517f9d12ec4f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b0381811660008181526010602090815260408083208054336001600160a01b031980831682179093556011909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a60408051610100808201835260025460ff808216845291810464ffffffffff166020808501919091526601000000000000820463ffffffff16848601526a0100000000000000000000820467ffffffffffffffff9081166060860152600160901b9092048216608085015260035480831660a086015268010000000000000000810490921660c0850152600160801b90910462ffffff1660e08401523360009081526008825293909320549394509092918c01359116610aa8576040517fb1c1f68e00000000000000000000000000000000000000000000000000000000815233600482015260240161061a565b6005548b3514610af2576005546040517f93df584c00000000000000000000000000000000000000000000000000000000815260048101919091528b35602482015260440161061a565b610b008a8a8a8a8a8a611f52565b8151610b0d90600161391f565b60ff1687141580610b1e5750868514155b15610b76578151610b3090600161391f565b6040517ffc33647500000000000000000000000000000000000000000000000000000000815260ff9091166004820152602481018890526044810186905260640161061a565b60008a8a604051610b88929190613944565b604051908190038120610b9f918e90602001613954565b60408051601f19818403018152828252805160209182012083830190925260008084529083018190529092509060005b8a811015610d3a5760006001858a8460208110610bee57610bee61387e565b610bfb91901a601b61391f565b8f8f86818110610c0d57610c0d61387e565b905060200201358e8e87818110610c2657610c2661387e565b9050602002013560405160008152602001604052604051610c63949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015610c85573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526009602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955092509050610d13576040517f20fb74ee0000000000000000000000000000000000000000000000000000000081526001600160a01b038216600482015260240161061a565b826020015160080260ff166001901b84019350508080610d3290613906565b915050610bcf565b5081827e010101010101010101010101010101010101010101010101010101010101011614610d95576040517fc103be2e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5060009150819050610de5848e836020020135858f8f8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611fe292505050565b6007805463ffffffff191663ffffffff600888901c161790559092509050610e1084838388336123eb565b50505050505050505050505050565b6013546001600160a01b0316338114610e5c5760405163292f4fb560e01b81523360048201526001600160a01b038216602482015260440161061a565b8151604051610e6e9190602001613970565b60408051601f1981840301815291905280516020909101206015555050565b6000546001600160a01b03163314801590610f1b5750601254604051630d629b5f60e31b81526001600160a01b0390911690636b14daf890610ed890339060009036906004016139b5565b602060405180830381865afa158015610ef5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f1991906139d8565b155b15610f3957604051631809d98560e31b815260040160405180910390fd5b6000610f43612503565b905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166345ccbb8b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610fa5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fc991906139fa565b90508181101561100f576040517fcf479181000000000000000000000000000000000000000000000000000000008152600481018290526024810183905260440161061a565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663f99b1d688561105261104c8686613a13565b876126bf565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b1681526001600160a01b03909216600483015260248201526044015b600060405180830381600087803b1580156110b157600080fd5b505af11580156110c5573d6000803e3d6000fd5b5050505050505050565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166345ccbb8b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611130573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061115491906139fa565b90506000611160612503565b905061116c8183613a2a565b9250505090565b888787601f8311156111bb576040517f809fc428000000000000000000000000000000000000000000000000000000008152601f60048201526024810184905260440161061a565b8183146111fe576040517f988a0804000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260440161061a565b611209816003613a9e565b60ff168311611250576040517ffda9db7800000000000000000000000000000000000000000000000000000000815260ff821660048201526024810184905260440161061a565b61125c8160ff166126d9565b611264611ef6565b60006040518060c001604052808f8f80806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505081526020018d8d8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050509082525060ff8c1660208083019190915260408051601f8d0183900483028101830182528c8152920191908c908c908190840183828082843760009201919091525050509082525067ffffffffffffffff891660208083019190915260408051601f8a01839004830281018301825289815292019190899089908190840183828082843760009201919091525050509152506002805465ffffffffff00191690559050611395612713565b600a5460005b81811015611446576000600a82815481106113b8576113b861387e565b6000918252602082200154600b80546001600160a01b03909216935090849081106113e5576113e561387e565b60009182526020808320909101546001600160a01b039485168352600982526040808420805461ffff1916905594168252600890529190912080546dffffffffffffffffffffffffffff19169055508061143e81613906565b91505061139b565b50611453600a6000613039565b61145f600b6000613039565b60005b82515181101561170b5760096000846000015183815181106114865761148661387e565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16156115105782518051829081106114c6576114c661387e565b60200260200101516040517f7451f83e00000000000000000000000000000000000000000000000000000000815260040161061a91906001600160a01b0391909116815260200190565b604080518082019091526001815260ff8216602082015283518051600991600091859081106115415761154161387e565b6020908102919091018101516001600160a01b03168252818101929092526040016000908120835181549484015161ffff1990951690151561ff0019161761010060ff909516949094029390931790925584015180516008929190849081106115ac576115ac61387e565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff161561163857826020015181815181106115ee576115ee61387e565b60200260200101516040517fe8d2989900000000000000000000000000000000000000000000000000000000815260040161061a91906001600160a01b0391909116815260200190565b60405180606001604052806001151581526020018260ff16815260200160006001600160601b0316815250600860008560200151848151811061167d5761167d61387e565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020835181549385015194909201516001600160601b031662010000026dffffffffffffffffffffffff00001960ff959095166101000261ff00199315159390931661ffff199094169390931791909117929092161790558061170381613906565b915050611462565b508151805161172291600a91602090910190613057565b50602080830151805161173992600b920190613057565b5060408201516002805460ff191660ff909216919091179055600454640100000000900463ffffffff1661176b612bc9565b6004805463ffffffff9283166401000000000267ffffffff0000000019821681179092556000926117a29281169116176001613ac7565b905080600460006101000a81548163ffffffff021916908363ffffffff16021790555060006117f646308463ffffffff16886000015189602001518a604001518b606001518c608001518d60a00151612c53565b9050806005600001819055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05838284886000015189602001518a604001518b606001518c608001518d60a0015160405161185999989796959493929190613b33565b60405180910390a16002546601000000000000900463ffffffff1660005b8651518110156118d15781600c82601f81106118955761189561387e565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff16021790555080806118c990613906565b915050611877565b506118dc8e8e612ce0565b50505050505050505050505050505050505050565b6001600160a01b03811660009081526008602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b031691810191909152906119535750600092915050565b600c816020015160ff16601f811061196d5761196d61387e565b6008810491909101546002546119a3926007166004026101000a90910463ffffffff9081169166010000000000009004166138aa565b9392505050565b6000546001600160a01b03163314801590611a385750601254604051630d629b5f60e31b81526001600160a01b0390911690636b14daf8906119f590339060009036906004016139b5565b602060405180830381865afa158015611a12573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a3691906139d8565b155b15611a5657604051631809d98560e31b815260040160405180910390fd5b611a5e612713565b600280547fffffffffffff00000000000000000000000000000000ffffffffffffffffffff166a010000000000000000000067ffffffffffffffff8881169182027fffffffffffff0000000000000000ffffffffffffffffffffffffffffffffffff1692909217600160901b88841690810291909117909355600380548784167fffffffffffffffffffffffffffffffff00000000000000000000000000000000909116811768010000000000000000948816948502177fffffffffffffffffffffffffff000000ffffffffffffffffffffffffffffffff16600160801b62ffffff881690810291909117909255604080519384526020840195909552828501526060820192909252608081019190915290517f49275ddcdfc9c0519b3d094308c8bf675f06070a754ce90c152163cb6e66e8a09181900360a00190a15050505050565b6001600160a01b03828116600090815260106020526040902054163314611bdc57604051633738e30960e21b815260040160405180910390fd5b6001600160a01b0381163303611c1e576040517fb387a23800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03808316600090815260116020526040902080548383166001600160a01b031982168117909255909116908114611c91576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a45b505050565b611c9e611ef6565b61020881612cee565b611caf611ef6565b61020881612d97565b6001600160a01b0381166000908152600860209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046001600160601b031692810192909252611d15575050565b6000611d20836104ea565b90508015611c91576001600160a01b0383811660009081526010602090815260409091205460025491850151921691660100000000000090910463ffffffff1690600c9060ff16601f8110611d7757611d7761387e565b6008808204909201805463ffffffff9485166004600790941684026101000a90810295021916939093179092556001600160a01b03808716600090815260209290925260409182902080546dffffffffffffffffffffffff00001916905590517ff99b1d680000000000000000000000000000000000000000000000000000000081527f00000000000000000000000000000000000000000000000000000000000000009091169163f99b1d6891611e479185918791016001600160a01b03929092168252602082015260400190565b600060405180830381600087803b158015611e6157600080fd5b505af1158015611e75573d6000803e3d6000fd5b505050507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316816001600160a01b0316856001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c85604051611ee891815260200190565b60405180910390a450505050565b6000546001600160a01b03163314611f505760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161061a565b565b6000611f5f8260206138cf565b611f6a8560206138cf565b611f76886101446138ee565b611f8091906138ee565b611f8a91906138ee565b611f959060006138ee565b9050368114611fd9576040517ff7b94f0a0000000000000000000000000000000000000000000000000000000081526004810182905236602482015260440161061a565b50505050505050565b600080600083806020019051810190611ffb9190613dd7565b64ffffffffff8616602089015260408801805191925061201a82613fbd565b63ffffffff1663ffffffff168152505086600260008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060808201518160000160126101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060a08201518160010160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060c08201518160010160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060e08201518160010160106101000a81548162ffffff021916908362ffffff16021790555090505060006121918260600151612e0d565b9050808260800151146121f457608082015160608301516040517faed0afe500000000000000000000000000000000000000000000000000000000815260048101929092526024820183905267ffffffffffffffff16604482015260640161061a565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663483af703836000015184602001518560400151866060015187608001516040518663ffffffff1660e01b815260040161225c9594939291906140e4565b6000604051808303816000875af115801561227b573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526122a391908101906141e7565b5060408281015190517f05f4acc600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906305f4acc690602401600060405180830381600087803b15801561232e57600080fd5b505af1158015612342573d6000803e3d6000fd5b505050508564ffffffffff16886040015163ffffffff167f27bf3f1077f091da6885751ba10f5775d06657fd59e47a6ab1f7635e5a115afe33856020015186604001518c6040516123ca94939291906001600160a01b039490941684526001600160c01b0392909216602084015267ffffffffffffffff166040830152606082015260800190565b60405180910390a38160200151826040015193509350505094509492505050565b60006124173a67ffffffffffffffff861615612407578561240d565b87608001515b8860600151612ee9565b90506010360260005a905060006124408663ffffffff1685858c60e0015162ffffff1686612f3a565b90506000670de0b6b3a76400006001600160c01b038a1683026001600160a01b03881660009081526008602052604090205460c08d01519290910492506201000090046001600160601b039081169167ffffffffffffffff16828401019081168211156124b357505050505050506108ea565b6001600160a01b038816600090815260086020526040902080546001600160601b0390921662010000026dffffffffffffffffffffffff0000199092169190911790555050505050505050505050565b600080600b80548060200260200160405190810160405280929190818152602001828054801561255c57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161253e575b50508351600254604080516103e08101918290529697509195660100000000000090910463ffffffff169450600093509150600c90601f908285855b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116125985790505050505050905060005b8381101561262b578181601f81106125f8576125f861387e565b602002015161260790846138aa565b6126179063ffffffff16876138ee565b95508061262381613906565b9150506125de565b506003546126439067ffffffffffffffff16866138cf565b945060005b838110156126b757600860008683815181106126665761266661387e565b6020908102919091018101516001600160a01b03168252810191909152604001600020546126a3906201000090046001600160601b0316876138ee565b9550806126af81613906565b915050612648565b505050505090565b6000818310156126d05750816126d3565b50805b92915050565b80600003610208576040517fe77dba5600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600254604080516103e08101918290527f0000000000000000000000000000000000000000000000000000000000000000926601000000000000900463ffffffff169160009190600c90601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411612764579050505050505090506000600b8054806020026020016040519081016040528092919081815260200182805480156127ff57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116127e1575b5050505050905060008151905060008167ffffffffffffffff811115612827576128276133f1565b604051908082528060200260200182016040528015612850578160200160208202803683370190505b50905060008267ffffffffffffffff81111561286e5761286e6133f1565b604051908082528060200260200182016040528015612897578160200160208202803683370190505b5090506000805b84811015612b01576000600860008884815181106128be576128be61387e565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160029054906101000a90046001600160601b03166001600160601b031690506000600860008985815181106129205761292061387e565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160026101000a8154816001600160601b0302191690836001600160601b0316021790555060008883601f81106129835761298361387e565b6020020151600354908b0363ffffffff16915067ffffffffffffffff16810282018015612af6576000601060008b87815181106129c2576129c261387e565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b0316905080888781518110612a1357612a1361387e565b60200260200101906001600160a01b031690816001600160a01b03168152505081878781518110612a4657612a4661387e565b6020026020010181815250508b8b86601f8110612a6557612a6561387e565b602002019063ffffffff16908163ffffffff168152505085806001019650508c6001600160a01b0316816001600160a01b03168b8781518110612aaa57612aaa61387e565b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c85604051612aec91815260200190565b60405180910390a4505b50505060010161289e565b5081518114612b11578082528083525b612b1e600c87601f6130bc565b508151156110c5576040517f73433a2f0000000000000000000000000000000000000000000000000000000081526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906373433a2f90612b8d90869086906004016142d9565b600060405180830381600087803b158015612ba757600080fd5b505af1158015612bbb573d6000803e3d6000fd5b505050505050505050505050565b60004661a4b1811480612bde575062066eed81145b15612c4c5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612c22573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c4691906139fa565b91505090565b4391505090565b6000808a8a8a8a8a8a8a8a8a604051602001612c7799989796959493929190614330565b60408051601f1981840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b612cea8282612f82565b5050565b336001600160a01b03821603612d465760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161061a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6012546001600160a01b039081169082168114612cea57601280546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912910160405180910390a15050565b60004661a4b1811480612e22575062066eed81145b15612ed9576101008367ffffffffffffffff16612e3d612bc9565b612e479190613a13565b1115612e565750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a8290602401602060405180830381865afa158015612eb5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119a391906139fa565b505067ffffffffffffffff164090565b60008367ffffffffffffffff8416811015612f1d576002858567ffffffffffffffff160381612f1a57612f1a6142c3565b04015b612f31818467ffffffffffffffff166126bf565b95945050505050565b600081861015612f76576040517f3fef97df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50909303019091010290565b610100818114612fc4578282826040517fb93aa5de00000000000000000000000000000000000000000000000000000000815260040161061a939291906143b8565b6000612fd2838501856143dc565b90506040517f8eef585f0000000000000000000000000000000000000000000000000000000081526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690638eef585f90611097908490600401614464565b5080546000825590600052602060002090810190610208919061314f565b8280548282559060005260206000209081019282156130ac579160200282015b828111156130ac57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613077565b506130b892915061314f565b5090565b6004830191839082156130ac5791602002820160005b8382111561311657835183826101000a81548163ffffffff021916908363ffffffff16021790555092602001926004016020816003010492830192600103026130d2565b80156131465782816101000a81549063ffffffff0219169055600401602081600301049283019260010302613116565b50506130b89291505b5b808211156130b85760008155600101613150565b6001600160a01b038116811461020857600080fd5b60006020828403121561318b57600080fd5b81356119a381613164565b6000602082840312156131a857600080fd5b813567ffffffffffffffff8111156131bf57600080fd5b820160a081850312156119a357600080fd5b60005b838110156131ec5781810151838201526020016131d4565b838111156131fb576000848401525b50505050565b600081518084526132198160208601602086016131d1565b601f01601f19169290920160200192915050565b6020815260006119a36020830184613201565b60008083601f84011261325257600080fd5b50813567ffffffffffffffff81111561326a57600080fd5b6020830191508360208260051b850101111561328557600080fd5b9250929050565b600080600080604085870312156132a257600080fd5b843567ffffffffffffffff808211156132ba57600080fd5b6132c688838901613240565b909650945060208701359150808211156132df57600080fd5b506132ec87828801613240565b95989497509550505050565b60008083601f84011261330a57600080fd5b50813567ffffffffffffffff81111561332257600080fd5b60208301915083602082850101111561328557600080fd5b60008060008060008060008060e0898b03121561335657600080fd5b606089018a81111561336757600080fd5b8998503567ffffffffffffffff8082111561338157600080fd5b61338d8c838d016132f8565b909950975060808b01359150808211156133a657600080fd5b6133b28c838d01613240565b909750955060a08b01359150808211156133cb57600080fd5b506133d88b828c01613240565b999c989b50969995989497949560c00135949350505050565b634e487b7160e01b600052604160045260246000fd5b6040805190810167ffffffffffffffff8111828210171561342a5761342a6133f1565b60405290565b604051610100810167ffffffffffffffff8111828210171561342a5761342a6133f1565b60405160a0810167ffffffffffffffff8111828210171561342a5761342a6133f1565b6040516080810167ffffffffffffffff8111828210171561342a5761342a6133f1565b6040516020810167ffffffffffffffff8111828210171561342a5761342a6133f1565b604051601f8201601f1916810167ffffffffffffffff811182821017156134e6576134e66133f1565b604052919050565b600067ffffffffffffffff821115613508576135086133f1565b50601f01601f191660200190565b600067ffffffffffffffff821115613530576135306133f1565b5060051b60200190565b600082601f83011261354b57600080fd5b8135602061356061355b83613516565b6134bd565b82815260059290921b8401810191818101908684111561357f57600080fd5b8286015b8481101561359a5780358352918301918301613583565b509695505050505050565b600060208083850312156135b857600080fd5b823567ffffffffffffffff808211156135d057600080fd5b90840190604082870312156135e457600080fd5b6135ec613407565b8235828111156135fb57600080fd5b8301601f8101881361360c57600080fd5b803561361a61355b826134ee565b818152898783850101111561362e57600080fd5b81878401888301376000878383010152808452505050838301358281111561365557600080fd5b6136618882860161353a565b948201949094529695505050505050565b6000806040838503121561368557600080fd5b823561369081613164565b946020939093013593505050565b803560ff811681146136af57600080fd5b919050565b67ffffffffffffffff8116811461020857600080fd5b80356136af816136b4565b60008060008060008060008060008060c08b8d0312156136f457600080fd5b8a3567ffffffffffffffff8082111561370c57600080fd5b6137188e838f01613240565b909c509a5060208d013591508082111561373157600080fd5b61373d8e838f01613240565b909a50985088915061375160408e0161369e565b975060608d013591508082111561376757600080fd5b6137738e838f016132f8565b909750955085915061378760808e016136ca565b945060a08d013591508082111561379d57600080fd5b506137aa8d828e016132f8565b915080935050809150509295989b9194979a5092959850565b62ffffff8116811461020857600080fd5b600080600080600060a086880312156137ec57600080fd5b85356137f7816136b4565b94506020860135613807816136b4565b93506040860135613817816136b4565b92506060860135613827816136b4565b91506080860135613837816137c3565b809150509295509295909350565b6000806040838503121561385857600080fd5b823561386381613164565b9150602083013561387381613164565b809150509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600063ffffffff838116908316818110156138c7576138c7613894565b039392505050565b60008160001904831182151516156138e9576138e9613894565b500290565b6000821982111561390157613901613894565b500190565b60006001820161391857613918613894565b5060010190565b600060ff821660ff84168060ff0382111561393c5761393c613894565b019392505050565b8183823760009101908152919050565b8281526060826020830137600060809190910190815292915050565b600082516139828184602087016131d1565b9190910192915050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0384168152604060208201526000612f3160408301848661398c565b6000602082840312156139ea57600080fd5b815180151581146119a357600080fd5b600060208284031215613a0c57600080fd5b5051919050565b600082821015613a2557613a25613894565b500390565b6000808312837f800000000000000000000000000000000000000000000000000000000000000001831281151615613a6457613a64613894565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff018313811615613a9857613a98613894565b50500390565b600060ff821660ff84168160ff0481118215151615613abf57613abf613894565b029392505050565b600063ffffffff808316818516808303821115613ae657613ae6613894565b01949350505050565b600081518084526020808501945080840160005b83811015613b285781516001600160a01b031687529582019590820190600101613b03565b509495945050505050565b600061012063ffffffff808d1684528b6020850152808b16604085015250806060840152613b638184018a613aef565b90508281036080840152613b778189613aef565b905060ff871660a084015282810360c0840152613b948187613201565b905067ffffffffffffffff851660e0840152828103610100840152613bb98185613201565b9c9b505050505050505050505050565b80516136af816136b4565b805165ffffffffffff811681146136af57600080fd5b805161ffff811681146136af57600080fd5b80516136af81613164565b600082601f830112613c1857600080fd5b8151613c2661355b826134ee565b818152846020838601011115613c3b57600080fd5b6105d98260208301602087016131d1565b80516001600160601b03811681146136af57600080fd5b600082601f830112613c7457600080fd5b81516020613c8461355b83613516565b82815260059290921b84018101918181019086841115613ca357600080fd5b8286015b8481101561359a57805167ffffffffffffffff80821115613cc757600080fd5b90880190601f196040838c0382011215613ce057600080fd5b613ce8613407565b8784015183811115613cf957600080fd5b8401610100818e0384011215613d0e57600080fd5b613d16613430565b9250613d23898201613bd4565b8352613d3160408201613bea565b89840152613d4160608201613bfc565b6040840152608081015184811115613d5857600080fd5b613d668e8b83850101613c07565b606085015250613d7860a08201613c4c565b608084015260c081015160a084015260e081015160c084015261010081015160e084015250818152613dac60408501613c4c565b818901528652505050918301918301613ca7565b80516001600160c01b03811681146136af57600080fd5b600060208284031215613de957600080fd5b815167ffffffffffffffff80821115613e0157600080fd5b9083019060a08286031215613e1557600080fd5b613e1d613454565b825182811115613e2c57600080fd5b8301601f81018713613e3d57600080fd5b8051613e4b61355b82613516565b8082825260208201915060208360051b850101925089831115613e6d57600080fd5b602084015b83811015613f6d57805187811115613e8957600080fd5b850160a0818d03601f19011215613e9f57600080fd5b613ea7613477565b6020820151613eb5816136b4565b81526040820151613ec5816137c3565b60208201526040828e03605f19011215613ede57600080fd5b613ee661349a565b8d607f840112613ef557600080fd5b613efd613407565b808f60a086011115613f0e57600080fd5b606085015b60a08601811015613f2e578051835260209283019201613f13565b50825250604082015260a082015189811115613f4957600080fd5b613f588e602083860101613c63565b60608301525084525060209283019201613e72565b50845250613f8091505060208401613dc0565b6020820152613f9160408401613bc9565b6040820152613fa260608401613bc9565b60608201526080830151608082015280935050505092915050565b600063ffffffff808316818103613fd657613fd6613894565b6001019392505050565b600081518084526020808501808196508360051b8101915082860160005b858110156140d757828403895281516040815181875265ffffffffffff81511682880152878101516060614037818a018361ffff169052565b928201519260809150614054898301856001600160a01b03169052565b8083015193505061010060a081818b01526140736101408b0186613201565b9284015192945060c06140908b8201856001600160601b03169052565b9084015160e08b81019190915290840151918a01919091529091015161012088015250908601516001600160601b0316948601949094529784019790840190600101613ffe565b5091979650505050505050565b600060a080830181845280895180835260c08601915060c08160051b87010192506020808c016000805b848110156141955789870360bf190186528251805167ffffffffffffffff1688528481015162ffffff1685890152604080820151519084908a015b6002821015614168578251815291870191600191909101908701614149565b5050506060015160808801899052614182888a0182613fe0565b975050948301949183019160010161410e565b5050508395506141af8188018c6001600160c01b03169052565b50505050506141ca604083018667ffffffffffffffff169052565b67ffffffffffffffff939093166060820152608001529392505050565b600060208083850312156141fa57600080fd5b825167ffffffffffffffff81111561421157600080fd5b8301601f8101851361422257600080fd5b805161423061355b82613516565b81815260079190911b8201830190838101908783111561424f57600080fd5b928401925b828410156142b8576080848903121561426d5760008081fd5b614275613477565b8451614280816136b4565b81528486015161428f816137c3565b818701526040858101519082015260608086015190820152825260809093019290840190614254565b979650505050505050565b634e487b7160e01b600052601260045260246000fd5b6040815260006142ec6040830185613aef565b82810360208481019190915284518083528582019282019060005b8181101561432357845183529383019391830191600101614307565b5090979650505050505050565b60006101208b83526001600160a01b038b16602084015267ffffffffffffffff808b16604085015281606085015261436a8285018b613aef565b9150838203608085015261437e828a613aef565b915060ff881660a085015283820360c085015261439b8288613201565b90861660e08501528381036101008501529050613bb98185613201565b6040815260006143cc60408301858761398c565b9050826020830152949350505050565b60006101008083850312156143f057600080fd5b83601f8401126143ff57600080fd5b60405181810181811067ffffffffffffffff82111715614421576144216133f1565b60405290830190808583111561443657600080fd5b845b8381101561445957803561444b816137c3565b825260209182019101614438565b509095945050505050565b6101008101818360005b600881101561449257815162ffffff1683526020928301929091019060010161446e565b5050509291505056fea164736f6c634300080f000a",
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
