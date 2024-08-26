// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package offramp

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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type InternalAny2EVMRampMessage struct {
	Header       InternalRampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     common.Address
	GasLimit     *big.Int
	TokenAmounts []InternalRampTokenAmount
}

type InternalExecutionReportSingleChain struct {
	SourceChainSelector uint64
	Messages            []InternalAny2EVMRampMessage
	OffchainTokenData   [][][]byte
	Proofs              [][32]byte
	ProofFlagBits       *big.Int
}

type InternalGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     *big.Int
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	GasPriceUpdates   []InternalGasPriceUpdate
}

type InternalRampMessageHeader struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type InternalRampTokenAmount struct {
	SourcePoolAddress []byte
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            *big.Int
	DestExecData      []byte
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type MultiOCR3BaseConfigInfo struct {
	ConfigDigest                   [32]byte
	F                              uint8
	N                              uint8
	IsSignatureVerificationEnabled bool
}

type MultiOCR3BaseOCRConfig struct {
	ConfigInfo   MultiOCR3BaseConfigInfo
	Signers      []common.Address
	Transmitters []common.Address
}

type MultiOCR3BaseOCRConfigArgs struct {
	ConfigDigest                   [32]byte
	OcrPluginType                  uint8
	F                              uint8
	IsSignatureVerificationEnabled bool
	Signers                        []common.Address
	Transmitters                   []common.Address
}

<<<<<<<< HEAD:core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp/evm_2_evm_multi_offramp.go
var EVM2EVMMultiOffRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyExecuted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"}],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"messageDestChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidMessageDestChainSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRoot\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidStaticConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeavesCannotBeEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"errorReason\",\"type\":\"bytes\"}],\"name\":\"MessageValidationError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"RootAlreadyCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleCommitReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelectorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"}],\"name\":\"StaticConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"isBlessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReportSingleChain[]\",\"name\":\"reports\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"gasLimitOverrides\",\"type\":\"uint256[][]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.UnblessedRoot[]\",\"name\":\"rootToReset\",\"type\":\"tuple[]\"}],\"name\":\"resetUnblessedRoots\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162006bab38038062006bab8339810160408190526200003591620008c7565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001fa565b5050466080525060208301516001600160a01b03161580620000ec575060408301516001600160a01b0316155b8062000103575060608301516001600160a01b0316155b1562000122576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160401b03166000036200014e5760405163c656089560e01b815260040160405180910390fd5b82516001600160401b0390811660a052602080850180516001600160a01b0390811660c05260408088018051831660e0526060808a01805185166101005283518b519098168852945184169587019590955251821690850152905116908201527f683eb52ee924eb817377cfa8f41f238f4bb7a877da5267869dfffbad85f564d89060800160405180910390a1620001e682620002a5565b620001f181620003c1565b50505062000c67565b336001600160a01b03821603620002545760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0316620002ce576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516060808b0180516001600160a01b039b8c166001600160c01b0319909a168a17600160a01b63ffffffff98891602176001600160c01b0316600160c01b948816949094026001600160e01b031693909317600160e01b93871693909302929092179098556080808b018051600580546001600160a01b031916918d169190911790558451988952955185169688019690965290518316918601919091525116938301939093529151909216908201527fa55bd56595c45f517e5967a3067f3dca684445a3080e7c04a4e0d5a40cda627d9060a00160405180910390a150565b60005b815181101562000666576000828281518110620003e557620003e562000a1d565b60200260200101519050600081602001519050806001600160401b0316600003620004235760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b03166200044c576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b0381166000908152600660205260408120600181018054919291620004789062000a33565b80601f0160208091040260200160405190810160405280929190818152602001828054620004a69062000a33565b8015620004f75780601f10620004cb57610100808354040283529160200191620004f7565b820191906000526020600020905b815481529060010190602001808311620004d957829003601f168201915b5050505050905060008460600151905081516000036200059e57805160000362000534576040516342bcdf7f60e11b815260040160405180910390fd5b6001830162000544828262000ac4565b508254600160a81b600160e81b031916600160a81b1783556040516001600160401b03851681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1620005d9565b8080519060200120828051906020012014620005d95760405163c39a620560e01b81526001600160401b038516600482015260240162000083565b604080860151845487516001600160a01b03166001600160a01b0319921515600160a01b02929092166001600160a81b031990911617178455516001600160401b038516907f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b906200064d90869062000b90565b60405180910390a25050505050806001019050620003c4565b5050565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715620006a557620006a56200066a565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620006d657620006d66200066a565b604052919050565b80516001600160401b0381168114620006f657600080fd5b919050565b6001600160a01b03811681146200071157600080fd5b50565b805163ffffffff81168114620006f657600080fd5b6000601f83601f8401126200073d57600080fd5b825160206001600160401b03808311156200075c576200075c6200066a565b8260051b6200076d838201620006ab565b93845286810183019383810190898611156200078857600080fd5b84890192505b85831015620008ba57825184811115620007a85760008081fd5b89016080601f19828d038101821315620007c25760008081fd5b620007cc62000680565b88840151620007db81620006fb565b81526040620007ec858201620006de565b8a8301526060808601518015158114620008065760008081fd5b838301529385015193898511156200081e5760008081fd5b84860195508f603f8701126200083657600094508485fd5b8a8601519450898511156200084f576200084f6200066a565b620008608b858f88011601620006ab565b93508484528f82868801011115620008785760008081fd5b60005b8581101562000898578681018301518582018d01528b016200087b565b5060009484018b0194909452509182015283525091840191908401906200078e565b9998505050505050505050565b6000806000838503610140811215620008df57600080fd5b6080811215620008ee57600080fd5b620008f862000680565b6200090386620006de565b815260208601516200091581620006fb565b602082015260408601516200092a81620006fb565b604082015260608601516200093f81620006fb565b6060820152935060a0607f19820112156200095957600080fd5b5060405160a081016001600160401b0380821183831017156200098057620009806200066a565b81604052608087015191506200099682620006fb565b818352620009a760a0880162000714565b6020840152620009ba60c0880162000714565b6040840152620009cd60e0880162000714565b60608401526101008701519150620009e582620006fb565b608083018290526101208701519294508083111562000a0357600080fd5b505062000a138682870162000729565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b600181811c9082168062000a4857607f821691505b60208210810362000a6957634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111562000abf576000816000526020600020601f850160051c8101602086101562000a9a5750805b601f850160051c820191505b8181101562000abb5782815560010162000aa6565b5050505b505050565b81516001600160401b0381111562000ae05762000ae06200066a565b62000af88162000af1845462000a33565b8462000a6f565b602080601f83116001811462000b30576000841562000b175750858301515b600019600386901b1c1916600185901b17855562000abb565b600085815260208120601f198616915b8281101562000b615788860151825594840194600190910190840162000b40565b508582101562000b805787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b602080825282546001600160a01b0381168383015260a081901c60ff161515604084015260a81c6001600160401b0316606083015260808083015260018084018054600093929190849062000be58162000a33565b8060a089015260c0600183166000811462000c09576001811462000c265762000c58565b60ff19841660c08b015260c083151560051b8b0101945062000c58565b85600052602060002060005b8481101562000c4f5781548c820185015290880190890162000c32565b8b0160c0019550505b50929998505050505050505050565b60805160a05160c05160e05161010051615ed562000cd6600039600081816102530152612cac0152600081816102240152612f860152600081816101f50152818161147b01526118d00152600081816101c501526128a1015260008181611e5e0152611eaa0152615ed56000f3fe608060405234801561001057600080fd5b50600436106101775760003560e01c806385572ffb116100d8578063d2a15d351161008c578063f2fde38b11610066578063f2fde38b1461059c578063f716f99f146105af578063ff888fb1146105c257600080fd5b8063d2a15d3514610556578063e9d68a8e14610569578063ece670b61461058957600080fd5b8063991a5018116100bd578063991a5018146104de578063c673e584146104f1578063ccd37ba31461051157600080fd5b806385572ffb146104b55780638da5cb5b146104c357600080fd5b80633f4b04aa1161012f5780637437ff9f116101145780637437ff9f1461038557806379ba50971461049a5780637d4eef60146104a257600080fd5b80633f4b04aa146103495780635e36480c1461036557600080fd5b8063181f5a7711610160578063181f5a77146102da5780632d04ab7614610323578063311cd5131461033657600080fd5b806304666f9c1461017c57806306285c6914610191575b600080fd5b61018f61018a366004614091565b6105e5565b005b61028360408051608081018252600080825260208201819052918101829052606081019190915260405180608001604052807f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b6040516102d19190815167ffffffffffffffff1681526020808301516001600160a01b0390811691830191909152604080840151821690830152606092830151169181019190915260800190565b60405180910390f35b6103166040518060400160405280601d81526020017f45564d3245564d4d756c74694f666652616d7020312e362e302d64657600000081525081565b6040516102d19190614200565b61018f6103313660046142ab565b6105f9565b61018f61034436600461435e565b6109fd565b60095460405167ffffffffffffffff90911681526020016102d1565b6103786103733660046143b2565b610a66565b6040516102d1919061440f565b61043d6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152506040805160a0810182526004546001600160a01b03808216835263ffffffff600160a01b830481166020850152780100000000000000000000000000000000000000000000000083048116948401949094527c01000000000000000000000000000000000000000000000000000000009091049092166060820152600554909116608082015290565b6040516102d19190600060a0820190506001600160a01b03808451168352602084015163ffffffff808216602086015280604087015116604086015280606087015116606086015250508060808501511660808401525092915050565b61018f610abc565b61018f6104b036600461497c565b610b7a565b61018f610177366004614aa7565b6000546040516001600160a01b0390911681526020016102d1565b61018f6104ec366004614af6565b610d1a565b6105046104ff366004614b7b565b610d2b565b6040516102d19190614bdb565b61054861051f366004614c50565b67ffffffffffffffff919091166000908152600860209081526040808320938352929052205490565b6040519081526020016102d1565b61018f610564366004614c7a565b610e89565b61057c610577366004614cef565b610f43565b6040516102d19190614d0a565b61018f610597366004614d58565b611062565b61018f6105aa366004614dbc565b6113d2565b61018f6105bd366004614e41565b6113e3565b6105d56105d0366004614f7f565b611425565b60405190151581526020016102d1565b6105ed6114e6565b6105f681611542565b50565b60006106078789018961511d565b8051515190915015158061062057508051602001515115155b156107205760095460208a01359067ffffffffffffffff808316911610156106df576009805467ffffffffffffffff191667ffffffffffffffff83161790556004805483516040517f3937306f0000000000000000000000000000000000000000000000000000000081526001600160a01b0390921692633937306f926106a8929101615385565b600060405180830381600087803b1580156106c257600080fd5b505af11580156106d6573d6000803e3d6000fd5b5050505061071e565b81602001515160000361071e576040517f2261116700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505b60005b8160200151518110156109465760008260200151828151811061074857610748615288565b6020026020010151905060008160000151905061076481611884565b600061076f82611986565b602084015151815491925067ffffffffffffffff9081167501000000000000000000000000000000000000000000909204161415806107c5575060208084015190810151905167ffffffffffffffff9182169116115b1561080e57825160208401516040517feefb0cac000000000000000000000000000000000000000000000000000000008152610805929190600401615398565b60405180910390fd5b60408301518061084a576040517f504570e300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b835167ffffffffffffffff166000908152600860209081526040808320848452909152902054156108bd5783516040517f32cf0cbf00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260248101829052604401610805565b60208085015101516108d09060016153e3565b82547fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000067ffffffffffffffff928316021790925592511660009081526008602090815260408083209483529390529190912042905550600101610723565b507f3a3950e13dd607cc37980db0ef14266c40d2bba9c01b2e44bfe549808883095d81604051610976919061540b565b60405180910390a16109f260008a8a8a8a8a8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808e0282810182019093528d82529093508d92508c9182918501908490808284376000920191909152508b92506119ed915050565b505050505050505050565b610a3d610a0c828401846154a8565b6040805160008082526020820190925290610a37565b6060815260200190600190039081610a225790505b50611d64565b604080516000808252602082019092529050610a606001858585858660006119ed565b50505050565b6000610a74600160046154dd565b6002610a81608085615506565b67ffffffffffffffff16610a95919061552d565b610a9f8585611e14565b901c166003811115610ab357610ab36143e5565b90505b92915050565b6001546001600160a01b03163314610b165760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610805565b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b82611e5b565b815181518114610bbe576040517f83e3f56400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610d0a576000848281518110610bdd57610bdd615288565b60200260200101519050600081602001515190506000858481518110610c0557610c05615288565b6020026020010151905080518214610c49576040517f83e3f56400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b82811015610cfb576000828281518110610c6857610c68615288565b6020026020010151905080600014610cf25784602001518281518110610c9057610c90615288565b602002602001015160800151811015610cf25784516040517fc8e9605100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024810183905260448101829052606401610805565b50600101610c4c565b50505050806001019050610bc1565b50610d158383611d64565b505050565b610d226114e6565b6105f681611edc565b610d6e6040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c082015294855291820180548451818402810184019095528085529293858301939092830182828015610e1757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610df9575b5050505050815260200160038201805480602002602001604051908101604052809291908181526020018280548015610e7957602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610e5b575b5050505050815250509050919050565b610e916114e6565b60005b81811015610d15576000838383818110610eb057610eb0615288565b905060400201803603810190610ec69190615544565b9050610ed58160200151611425565b610f3a57805167ffffffffffffffff1660009081526008602090815260408083208285018051855290835281842093909355915191519182527f202f1139a3e334b6056064c0e9b19fd07e44a88d8f6e5ded571b24cf8c371f12910160405180910390a15b50600101610e94565b604080516080808201835260008083526020808401829052838501829052606080850181905267ffffffffffffffff878116845260068352928690208651948501875280546001600160a01b0381168652600160a01b810460ff161515938601939093527501000000000000000000000000000000000000000000909204909216948301949094526001840180549394929391840191610fe29061557d565b80601f016020809104026020016040519081016040528092919081815260200182805461100e9061557d565b8015610e795780601f1061103057610100808354040283529160200191610e79565b820191906000526020600020905b81548152906001019060200180831161103e57505050919092525091949350505050565b33301461109b576040517f371a732800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60408051600080825260208201909252816110d8565b60408051808201909152600080825260208201528152602001906001900390816110b15790505b5060a0840151519091501561110b576111088360a001518460200151856060015186600001516020015186612089565b90505b6040805160a0810182528451518152845160209081015167ffffffffffffffff1681830152808601518351600094840192611147929101614200565b60408051601f19818403018152918152908252868101516020830152018390526005549091506001600160a01b03168015611254576040517f08d450a10000000000000000000000000000000000000000000000000000000081526001600160a01b038216906308d450a1906111c1908590600401615659565b600060405180830381600087803b1580156111db57600080fd5b505af19250505080156111ec575060015b611254573d80801561121a576040519150601f19603f3d011682016040523d82523d6000602084013e61121f565b606091505b50806040517f09c253250000000000000000000000000000000000000000000000000000000081526004016108059190614200565b60408501515115801561126957506080850151155b80611280575060608501516001600160a01b03163b155b806112c0575060608501516112be906001600160a01b03167f85572ffb00000000000000000000000000000000000000000000000000000000612167565b155b156112cc575050505050565b845160209081015167ffffffffffffffff16600090815260069091526040808220546080880151606089015192517f3cf9798300000000000000000000000000000000000000000000000000000000815284936001600160a01b0390931692633cf9798392611344928992611388929160040161566c565b6000604051808303816000875af1158015611363573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261138b91908101906156a8565b5091509150816113c957806040517f0a8d6e8c0000000000000000000000000000000000000000000000000000000081526004016108059190614200565b50505050505050565b6113da6114e6565b6105f681612183565b6113eb6114e6565b60005b81518110156114215761141982828151811061140c5761140c615288565b6020026020010151612239565b6001016113ee565b5050565b6040805180820182523081526020810183815291517f4d61677100000000000000000000000000000000000000000000000000000000815290516001600160a01b039081166004830152915160248201526000917f00000000000000000000000000000000000000000000000000000000000000001690634d61677190604401602060405180830381865afa1580156114c2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ab6919061573e565b6000546001600160a01b031633146115405760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610805565b565b60005b815181101561142157600082828151811061156257611562615288565b602002602001015190506000816020015190508067ffffffffffffffff166000036115b9576040517fc656089500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81516001600160a01b03166115fa576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526006602052604081206001810180549192916116259061557d565b80601f01602080910402602001604051908101604052809291908181526020018280546116519061557d565b801561169e5780601f106116735761010080835404028352916020019161169e565b820191906000526020600020905b81548152906001019060200180831161168157829003601f168201915b5050505050905060008460600151905081516000036117815780516000036116f2576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001830161170082826157a3565b5082547fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff16750100000000000000000000000000000000000000000017835560405167ffffffffffffffff851681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a16117d4565b80805190602001208280519060200120146117d4576040517fc39a620500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401610805565b604080860151845487516001600160a01b031673ffffffffffffffffffffffffffffffffffffffff19921515600160a01b02929092167fffffffffffffffffffffff000000000000000000000000000000000000000000909116171784555167ffffffffffffffff8516907f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b9061186c908690615863565b60405180910390a25050505050806001019050611545565b6040517f2cbc26bb000000000000000000000000000000000000000000000000000000008152608082901b77ffffffffffffffff000000000000000000000000000000001660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa15801561191f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611943919061573e565b156105f6576040517ffdbd6a7200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610805565b67ffffffffffffffff811660009081526006602052604081208054600160a01b900460ff16610ab6576040517fed053c5900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610805565b60ff87811660009081526002602090815260408083208151608081018352815481526001909101548086169382019390935261010083048516918101919091526201000090910490921615156060830152873590611a4c8760a4615931565b9050826060015115611a94578451611a6590602061552d565b8651611a7290602061552d565b611a7d9060a0615931565b611a879190615931565b611a919082615931565b90505b368114611ad6576040517f8e1192e100000000000000000000000000000000000000000000000000000000815260048101829052366024820152604401610805565b5081518114611b1e5781516040517f93df584c000000000000000000000000000000000000000000000000000000008152600481019190915260248101829052604401610805565b611b26611e5b565b60ff808a1660009081526003602090815260408083203384528252808320815180830190925280548086168352939491939092840191610100909104166002811115611b7457611b746143e5565b6002811115611b8557611b856143e5565b9052509050600281602001516002811115611ba257611ba26143e5565b148015611bf65750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff1681548110611bde57611bde615288565b6000918252602090912001546001600160a01b031633145b611c2c576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50816060015115611d0e576020820151611c47906001615944565b60ff16855114611c83576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8351855114611cbe576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008787604051611cd092919061595d565b604051908190038120611ce7918b9060200161596d565b604051602081830303815290604052805190602001209050611d0c8a8288888861257d565b505b6040805182815260208a81013567ffffffffffffffff169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b8151600003611d9e576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160408051600080825260208201909252911591905b8451811015611e0d57611e05858281518110611dd357611dd3615288565b602002602001015184611dff57858381518110611df257611df2615288565b6020026020010151612794565b83612794565b600101611db5565b5050505050565b67ffffffffffffffff8216600090815260076020526040812081611e39608085615981565b67ffffffffffffffff1681526020810191909152604001600020549392505050565b467f000000000000000000000000000000000000000000000000000000000000000014611540576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f00000000000000000000000000000000000000000000000000000000000000006004820152466024820152604401610805565b80516001600160a01b0316611f1d576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516004805460208085018051604080880180516060808b0180516001600160a01b039b8c167fffffffffffffffff000000000000000000000000000000000000000000000000909a168a17600160a01b63ffffffff988916021777ffffffffffffffffffffffffffffffffffffffffffffffff167801000000000000000000000000000000000000000000000000948816949094027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16939093177c010000000000000000000000000000000000000000000000000000000093871693909302929092179098556080808b0180516005805473ffffffffffffffffffffffffffffffffffffffff1916918d169190911790558451988952955185169688019690965290518316918601919091525116938301939093529151909216908201527fa55bd56595c45f517e5967a3067f3dca684445a3080e7c04a4e0d5a40cda627d9060a00160405180910390a150565b6060855167ffffffffffffffff8111156120a5576120a5613ea8565b6040519080825280602002602001820160405280156120ea57816020015b60408051808201909152600080825260208201528152602001906001900390816120c35790505b50905060005b865181101561215d5761213887828151811061210e5761210e615288565b602002602001015187878787868151811061212b5761212b615288565b6020026020010151612f25565b82828151811061214a5761214a615288565b60209081029190910101526001016120f0565b5095945050505050565b600061217283613334565b8015610ab35750610ab38383613398565b336001600160a01b038216036121db5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610805565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff16600003612264576000604051631b3fab5160e11b815260040161080591906159a8565b60208082015160ff808216600090815260029093526040832060018101549293909283921690036122d157606084015160018201805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff909216919091179055612326565b6060840151600182015460ff6201000090910416151590151514612326576040517f87f6037c00000000000000000000000000000000000000000000000000000000815260ff84166004820152602401610805565b60a08401518051601f60ff82161115612355576001604051631b3fab5160e11b815260040161080591906159a8565b6123bb85856003018054806020026020016040519081016040528092919081815260200182805480156123b157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612393575b5050505050613467565b8560600151156124ea5761242985856002018054806020026020016040519081016040528092919081815260200182805480156123b1576020028201919060005260206000209081546001600160a01b03168152600190910190602001808311612393575050505050613467565b608086015180516124439060028701906020840190613e02565b5080516001850180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff841690810291909117909155601f10156124a3576002604051631b3fab5160e11b815260040161080591906159a8565b60408801516124b39060036159c2565b60ff168160ff16116124db576003604051631b3fab5160e11b815260040161080591906159a8565b6124e7878360016134d0565b50505b6124f6858360026134d0565b815161250b9060038601906020850190613e02565b5060408681015160018501805460ff191660ff8316179055875180865560a089015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f54793612564938a939260028b019291906159de565b60405180910390a161257585613650565b505050505050565b612585613e74565b835160005b8181101561278a5760006001888684602081106125a9576125a9615288565b6125b691901a601b615944565b8985815181106125c8576125c8615288565b60200260200101518986815181106125e2576125e2615288565b602002602001015160405160008152602001604052604051612620949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612642573d6000803e3d6000fd5b505060408051601f1981015160ff808e166000908152600360209081528582206001600160a01b038516835281528582208587019096528554808416865293975090955092939284019161010090041660028111156126a3576126a36143e5565b60028111156126b4576126b46143e5565b90525090506001816020015160028111156126d1576126d16143e5565b14612708576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f811061271f5761271f615288565b60200201511561275b576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f811061277657612776615288565b91151560209092020152505060010161258a565b5050505050505050565b815161279f81611884565b60006127aa82611986565b60208501515190915060008190036127ed576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b846040015151811461282b576040517f57e0e08300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008167ffffffffffffffff81111561284657612846613ea8565b60405190808252806020026020018201604052801561286f578160200160208202803683370190505b50905060005b828110156129e45760008760200151828151811061289557612895615288565b602002602001015190507f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681600001516040015167ffffffffffffffff161461292857805160409081015190517f38432a2200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610805565b6129be8186600101805461293b9061557d565b80601f01602080910402602001604051908101604052809291908181526020018280546129679061557d565b80156129b45780601f10612989576101008083540402835291602001916129b4565b820191906000526020600020905b81548152906001019060200180831161299757829003601f168201915b505050505061366c565b8383815181106129d0576129d0615288565b602090810291909101015250600101612875565b5060006129fb858389606001518a6080015161378e565b905080600003612a43576040517f7dd17a7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86166004820152602401610805565b8551151560005b848110156109f257600089602001518281518110612a6a57612a6a615288565b602002602001015190506000612a8889836000015160600151610a66565b90506000816003811115612a9e57612a9e6143e5565b1480612abb57506003816003811115612ab957612ab96143e5565b145b612b12578151606001516040805167ffffffffffffffff808d16825290921660208301527f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c91015b60405180910390a15050612f1d565b8315612be257600454600090600160a01b900463ffffffff16612b3587426154dd565b1190508080612b5557506003826003811115612b5357612b536143e5565b145b612b97576040517fa9cfc86200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8b166004820152602401610805565b8a8481518110612ba957612ba9615288565b6020026020010151600014612bdc578a8481518110612bca57612bca615288565b60200260200101518360800181815250505b50612c43565b6000816003811115612bf657612bf66143e5565b14612c43578151606001516040805167ffffffffffffffff808d16825290921660208301527f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe651209101612b03565b81516080015167ffffffffffffffff1615612d31576000816003811115612c6c57612c6c6143e5565b03612d315781516080015160208301516040517fe0e03cae0000000000000000000000000000000000000000000000000000000081526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263e0e03cae92612ce3928e929190600401615a90565b6020604051808303816000875af1158015612d02573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d26919061573e565b612d31575050612f1d565b60008b604001518481518110612d4957612d49615288565b6020026020010151905080518360a001515114612dad578251606001516040517f1cfe6d8b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff808d1660048301529091166024820152604401610805565b612dc18a84600001516060015160016137e4565b600080612dce858461388c565b91509150612de58c866000015160600151846137e4565b8615612e55576003826003811115612dff57612dff6143e5565b03612e55576000846003811115612e1857612e186143e5565b14612e55578451516040517f2b11b8d900000000000000000000000000000000000000000000000000000000815261080591908390600401615abd565b6002826003811115612e6957612e696143e5565b14612ec3576003826003811115612e8257612e826143e5565b14612ec3578451606001516040517f926c5a3e000000000000000000000000000000000000000000000000000000008152610805918e918590600401615ad6565b8451805160609091015160405167ffffffffffffffff918216918f16907f8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df290612f0f9087908790615afc565b60405180910390a450505050505b600101612a4a565b60408051808201909152600080825260208201526000612f488760200151613956565b6040517fbbe4f6db0000000000000000000000000000000000000000000000000000000081526001600160a01b0380831660048301529192506000917f0000000000000000000000000000000000000000000000000000000000000000169063bbe4f6db90602401602060405180830381865afa158015612fcd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ff19190615b1c565b90506001600160a01b038116158061303957506130376001600160a01b0382167faff2afbf00000000000000000000000000000000000000000000000000000000612167565b155b1561307b576040517fae9b4ce90000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602401610805565b6000806131816040518061010001604052808b81526020018967ffffffffffffffff1681526020018a6001600160a01b031681526020018c606001518152602001866001600160a01b031681526020018c6000015181526020018c604001518152602001888152506040516024016130f39190615b39565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f3907753700000000000000000000000000000000000000000000000000000000179052600454859063ffffffff7c01000000000000000000000000000000000000000000000000000000009091041661138860846139fc565b5091509150816131bf57806040517fe1cd55090000000000000000000000000000000000000000000000000000000081526004016108059190614200565b80516020146132075780516040517f78ef8024000000000000000000000000000000000000000000000000000000008152602060048201526024810191909152604401610805565b60008180602001905181019061321d9190615c06565b6040516001600160a01b038b166024820152604481018290529091506132ca9060640160408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052600454879063ffffffff78010000000000000000000000000000000000000000000000009091041661138860846139fc565b5090935091508261330957816040517fe1cd55090000000000000000000000000000000000000000000000000000000081526004016108059190614200565b604080518082019091526001600160a01b039095168552602085015250919250505095945050505050565b6000613360827f01ffc9a700000000000000000000000000000000000000000000000000000000613398565b8015610ab65750613391827fffffffff00000000000000000000000000000000000000000000000000000000613398565b1592915050565b604080517fffffffff000000000000000000000000000000000000000000000000000000008316602480830191909152825180830390910181526044909101909152602080820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a700000000000000000000000000000000000000000000000000000000178152825160009392849283928392918391908a617530fa92503d91506000519050828015613450575060208210155b801561345c5750600081115b979650505050505050565b60005b8151811015610d155760ff83166000908152600360205260408120835190919084908490811061349c5761349c615288565b6020908102919091018101516001600160a01b03168252810191909152604001600020805461ffff1916905560010161346a565b60005b82518160ff161015610a60576000838260ff16815181106134f6576134f6615288565b6020026020010151905060006002811115613513576135136143e5565b60ff80871660009081526003602090815260408083206001600160a01b03871684529091529020546101009004166002811115613552576135526143e5565b14613573576004604051631b3fab5160e11b815260040161080591906159a8565b6001600160a01b0381166135b3576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052808360ff1681526020018460028111156135d9576135d96143e5565b905260ff80871660009081526003602090815260408083206001600160a01b0387168452825290912083518154931660ff198416811782559184015190929091839161ffff191617610100836002811115613636576136366143e5565b0217905550905050508061364990615c1f565b90506134d3565b60ff81166105f6576009805467ffffffffffffffff1916905550565b8151602080820151604092830151925160009384936136b2937f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f93909291889101615c3e565b60408051601f1981840301815290829052805160209182012086518051888401516060808b0151908401516080808d015195015195976136fb9794969395929491939101615c71565b604051602081830303815290604052805190602001208560400151805190602001208660a001516040516020016137329190615d68565b60408051601f198184030181528282528051602091820120908301969096528101939093526060830191909152608082015260a081019190915260c0015b60405160208183030381529060405280519060200120905092915050565b60008061379c858585613b22565b90506137a781611425565b6137b55760009150506137dc565b67ffffffffffffffff86166000908152600860209081526040808320938352929052205490505b949350505050565b600060026137f3608085615506565b67ffffffffffffffff16613807919061552d565b905060006138158585611e14565b905081613824600160046154dd565b901b19168183600381111561383b5761383b6143e5565b67ffffffffffffffff871660009081526007602052604081209190921b9290921791829161386a608088615981565b67ffffffffffffffff1681526020810191909152604001600020555050505050565b6040517fece670b6000000000000000000000000000000000000000000000000000000008152600090606090309063ece670b6906138d09087908790600401615dc8565b600060405180830381600087803b1580156138ea57600080fd5b505af19250505080156138fb575060015b61393a573d808015613929576040519150601f19603f3d011682016040523d82523d6000602084013e61392e565b606091505b5060039250905061394f565b50506040805160208101909152600081526002905b9250929050565b6000815160201461399557816040517f8d666f600000000000000000000000000000000000000000000000000000000081526004016108059190614200565b6000828060200190518101906139ab9190615c06565b90506001600160a01b038111806139c3575061040081105b15610ab657826040517f8d666f600000000000000000000000000000000000000000000000000000000081526004016108059190614200565b6000606060008361ffff1667ffffffffffffffff811115613a1f57613a1f613ea8565b6040519080825280601f01601f191660200182016040528015613a49576020820181803683370190505b509150863b613a7c577f0c3b563c0000000000000000000000000000000000000000000000000000000060005260046000fd5b5a85811015613aaf577fafa32a2c0000000000000000000000000000000000000000000000000000000060005260046000fd5b8590036040810481038710613ae8577f37c3be290000000000000000000000000000000000000000000000000000000060005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d84811115613b0b5750835b808352806000602085013e50955095509592505050565b8251825160009190818303613b63576040517f11a6b26400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101018211801590613b7757506101018111155b613b94576040516309bde33960e01b815260040160405180910390fd5b60001982820101610100811115613bbe576040516309bde33960e01b815260040160405180910390fd5b80600003613beb5786600081518110613bd957613bd9615288565b60200260200101519350505050613dba565b60008167ffffffffffffffff811115613c0657613c06613ea8565b604051908082528060200260200182016040528015613c2f578160200160208202803683370190505b50905060008080805b85811015613d595760006001821b8b811603613c935788851015613c7c578c5160018601958e918110613c6d57613c6d615288565b60200260200101519050613cb5565b8551600185019487918110613c6d57613c6d615288565b8b5160018401938d918110613caa57613caa615288565b602002602001015190505b600089861015613ce5578d5160018701968f918110613cd657613cd6615288565b60200260200101519050613d07565b8651600186019588918110613cfc57613cfc615288565b602002602001015190505b82851115613d28576040516309bde33960e01b815260040160405180910390fd5b613d328282613dc1565b878481518110613d4457613d44615288565b60209081029190910101525050600101613c38565b506001850382148015613d6b57508683145b8015613d7657508581145b613d93576040516309bde33960e01b815260040160405180910390fd5b836001860381518110613da857613da8615288565b60200260200101519750505050505050505b9392505050565b6000818310613dd957613dd48284613ddf565b610ab3565b610ab383835b604080516001602082015290810183905260608101829052600090608001613770565b828054828255906000526020600020908101928215613e64579160200282015b82811115613e64578251825473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03909116178255602090920191600190910190613e22565b50613e70929150613e93565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b80821115613e705760008155600101613e94565b634e487b7160e01b600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613ee157613ee1613ea8565b60405290565b60405160a0810167ffffffffffffffff81118282101715613ee157613ee1613ea8565b60405160c0810167ffffffffffffffff81118282101715613ee157613ee1613ea8565b6040805190810167ffffffffffffffff81118282101715613ee157613ee1613ea8565b6040516060810167ffffffffffffffff81118282101715613ee157613ee1613ea8565b604051601f8201601f1916810167ffffffffffffffff81118282101715613f9c57613f9c613ea8565b604052919050565b600067ffffffffffffffff821115613fbe57613fbe613ea8565b5060051b60200190565b6001600160a01b03811681146105f657600080fd5b803567ffffffffffffffff81168114613ff557600080fd5b919050565b80151581146105f657600080fd5b8035613ff581613ffa565b600067ffffffffffffffff82111561402d5761402d613ea8565b50601f01601f191660200190565b600082601f83011261404c57600080fd5b813561405f61405a82614013565b613f73565b81815284602083860101111561407457600080fd5b816020850160208301376000918101602001919091529392505050565b600060208083850312156140a457600080fd5b823567ffffffffffffffff808211156140bc57600080fd5b818501915085601f8301126140d057600080fd5b81356140de61405a82613fa4565b81815260059190911b830184019084810190888311156140fd57600080fd5b8585015b838110156141a3578035858111156141195760008081fd5b86016080818c03601f19018113156141315760008081fd5b614139613ebe565b8983013561414681613fc8565b81526040614155848201613fdd565b8b83015260608085013561416881613ffa565b8383015292840135928984111561418157600091508182fd5b61418f8f8d8688010161403b565b908301525085525050918601918601614101565b5098975050505050505050565b60005b838110156141cb5781810151838201526020016141b3565b50506000910152565b600081518084526141ec8160208601602086016141b0565b601f01601f19169290920160200192915050565b602081526000610ab360208301846141d4565b8060608101831015610ab657600080fd5b60008083601f84011261423657600080fd5b50813567ffffffffffffffff81111561424e57600080fd5b60208301915083602082850101111561394f57600080fd5b60008083601f84011261427857600080fd5b50813567ffffffffffffffff81111561429057600080fd5b6020830191508360208260051b850101111561394f57600080fd5b60008060008060008060008060e0898b0312156142c757600080fd5b6142d18a8a614213565b9750606089013567ffffffffffffffff808211156142ee57600080fd5b6142fa8c838d01614224565b909950975060808b013591508082111561431357600080fd5b61431f8c838d01614266565b909750955060a08b013591508082111561433857600080fd5b506143458b828c01614266565b999c989b50969995989497949560c00135949350505050565b60008060006080848603121561437357600080fd5b61437d8585614213565b9250606084013567ffffffffffffffff81111561439957600080fd5b6143a586828701614224565b9497909650939450505050565b600080604083850312156143c557600080fd5b6143ce83613fdd565b91506143dc60208401613fdd565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b6004811061440b5761440b6143e5565b9052565b60208101610ab682846143fb565b600060a0828403121561442f57600080fd5b614437613ee7565b90508135815261444960208301613fdd565b602082015261445a60408301613fdd565b604082015261446b60608301613fdd565b606082015261447c60808301613fdd565b608082015292915050565b8035613ff581613fc8565b600082601f8301126144a357600080fd5b813560206144b361405a83613fa4565b82815260059290921b840181019181810190868411156144d257600080fd5b8286015b848110156145a857803567ffffffffffffffff808211156144f75760008081fd5b8189019150608080601f19848d030112156145125760008081fd5b61451a613ebe565b878401358381111561452c5760008081fd5b61453a8d8a8388010161403b565b825250604080850135848111156145515760008081fd5b61455f8e8b8389010161403b565b8a84015250606080860135858111156145785760008081fd5b6145868f8c838a010161403b565b92840192909252949092013593810193909352505083529183019183016144d6565b509695505050505050565b600061014082840312156145c657600080fd5b6145ce613f0a565b90506145da838361441d565b815260a082013567ffffffffffffffff808211156145f757600080fd5b6146038583860161403b565b602084015260c084013591508082111561461c57600080fd5b6146288583860161403b565b604084015261463960e08501614487565b6060840152610100840135608084015261012084013591508082111561465e57600080fd5b5061466b84828501614492565b60a08301525092915050565b600082601f83011261468857600080fd5b8135602061469861405a83613fa4565b82815260059290921b840181019181810190868411156146b757600080fd5b8286015b848110156145a857803567ffffffffffffffff8111156146db5760008081fd5b6146e98986838b01016145b3565b8452509183019183016146bb565b600082601f83011261470857600080fd5b8135602061471861405a83613fa4565b82815260059290921b8401810191818101908684111561473757600080fd5b8286015b848110156145a857803567ffffffffffffffff81111561475b5760008081fd5b6147698986838b010161403b565b84525091830191830161473b565b600082601f83011261478857600080fd5b8135602061479861405a83613fa4565b82815260059290921b840181019181810190868411156147b757600080fd5b8286015b848110156145a857803567ffffffffffffffff8111156147db5760008081fd5b6147e98986838b01016146f7565b8452509183019183016147bb565b600082601f83011261480857600080fd5b8135602061481861405a83613fa4565b8083825260208201915060208460051b87010193508684111561483a57600080fd5b602086015b848110156145a8578035835291830191830161483f565b600082601f83011261486757600080fd5b8135602061487761405a83613fa4565b82815260059290921b8401810191818101908684111561489657600080fd5b8286015b848110156145a857803567ffffffffffffffff808211156148bb5760008081fd5b818901915060a080601f19848d030112156148d65760008081fd5b6148de613ee7565b6148e9888501613fdd565b8152604080850135848111156148ff5760008081fd5b61490d8e8b83890101614677565b8a84015250606080860135858111156149265760008081fd5b6149348f8c838a0101614777565b838501525060809150818601358581111561494f5760008081fd5b61495d8f8c838a01016147f7565b918401919091525091909301359083015250835291830191830161489a565b600080604080848603121561499057600080fd5b833567ffffffffffffffff808211156149a857600080fd5b6149b487838801614856565b94506020915081860135818111156149cb57600080fd5b8601601f810188136149dc57600080fd5b80356149ea61405a82613fa4565b81815260059190911b8201840190848101908a831115614a0957600080fd5b8584015b83811015614a9557803586811115614a255760008081fd5b8501603f81018d13614a375760008081fd5b87810135614a4761405a82613fa4565b81815260059190911b82018a0190898101908f831115614a675760008081fd5b928b01925b82841015614a855783358252928a0192908a0190614a6c565b8652505050918601918601614a0d565b50809750505050505050509250929050565b600060208284031215614ab957600080fd5b813567ffffffffffffffff811115614ad057600080fd5b820160a08185031215613dba57600080fd5b803563ffffffff81168114613ff557600080fd5b600060a08284031215614b0857600080fd5b614b10613ee7565b8235614b1b81613fc8565b8152614b2960208401614ae2565b6020820152614b3a60408401614ae2565b6040820152614b4b60608401614ae2565b60608201526080830135614b5e81613fc8565b60808201529392505050565b803560ff81168114613ff557600080fd5b600060208284031215614b8d57600080fd5b610ab382614b6a565b60008151808452602080850194506020840160005b83811015614bd05781516001600160a01b031687529582019590820190600101614bab565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a0840152614c2a60e0840182614b96565b90506040840151601f198483030160c0850152614c478282614b96565b95945050505050565b60008060408385031215614c6357600080fd5b614c6c83613fdd565b946020939093013593505050565b60008060208385031215614c8d57600080fd5b823567ffffffffffffffff80821115614ca557600080fd5b818501915085601f830112614cb957600080fd5b813581811115614cc857600080fd5b8660208260061b8501011115614cdd57600080fd5b60209290920196919550909350505050565b600060208284031215614d0157600080fd5b610ab382613fdd565b602081526001600160a01b03825116602082015260208201511515604082015267ffffffffffffffff6040830151166060820152600060608301516080808401526137dc60a08401826141d4565b60008060408385031215614d6b57600080fd5b823567ffffffffffffffff80821115614d8357600080fd5b614d8f868387016145b3565b93506020850135915080821115614da557600080fd5b50614db2858286016146f7565b9150509250929050565b600060208284031215614dce57600080fd5b8135613dba81613fc8565b600082601f830112614dea57600080fd5b81356020614dfa61405a83613fa4565b8083825260208201915060208460051b870101935086841115614e1c57600080fd5b602086015b848110156145a8578035614e3481613fc8565b8352918301918301614e21565b60006020808385031215614e5457600080fd5b823567ffffffffffffffff80821115614e6c57600080fd5b818501915085601f830112614e8057600080fd5b8135614e8e61405a82613fa4565b81815260059190911b83018401908481019088831115614ead57600080fd5b8585015b838110156141a357803585811115614ec857600080fd5b860160c0818c03601f19011215614edf5760008081fd5b614ee7613f0a565b8882013581526040614efa818401614b6a565b8a8301526060614f0b818501614b6a565b8284015260809150614f1e828501614008565b9083015260a08381013589811115614f365760008081fd5b614f448f8d83880101614dd9565b838501525060c0840135915088821115614f5e5760008081fd5b614f6c8e8c84870101614dd9565b9083015250845250918601918601614eb1565b600060208284031215614f9157600080fd5b5035919050565b80357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff81168114613ff557600080fd5b600082601f830112614fd557600080fd5b81356020614fe561405a83613fa4565b82815260069290921b8401810191818101908684111561500457600080fd5b8286015b848110156145a857604081890312156150215760008081fd5b615029613f2d565b61503282613fdd565b815261503f858301614f98565b81860152835291830191604001615008565b600082601f83011261506257600080fd5b8135602061507261405a83613fa4565b82815260079290921b8401810191818101908684111561509157600080fd5b8286015b848110156145a85780880360808112156150af5760008081fd5b6150b7613f50565b6150c083613fdd565b8152604080601f19840112156150d65760008081fd5b6150de613f2d565b92506150eb878501613fdd565b83526150f8818501613fdd565b8388015281870192909252606083013591810191909152835291830191608001615095565b6000602080838503121561513057600080fd5b823567ffffffffffffffff8082111561514857600080fd5b8185019150604080838803121561515e57600080fd5b615166613f2d565b83358381111561517557600080fd5b84016040818a03121561518757600080fd5b61518f613f2d565b81358581111561519e57600080fd5b8201601f81018b136151af57600080fd5b80356151bd61405a82613fa4565b81815260069190911b8201890190898101908d8311156151dc57600080fd5b928a01925b8284101561522c5787848f0312156151f95760008081fd5b615201613f2d565b843561520c81613fc8565b8152615219858d01614f98565b818d0152825292870192908a01906151e1565b84525050508187013593508484111561524457600080fd5b6152508a858401614fc4565b818801528252508385013591508282111561526a57600080fd5b61527688838601615051565b85820152809550505050505092915050565b634e487b7160e01b600052603260045260246000fd5b805160408084528151848201819052600092602091908201906060870190855b8181101561530a57835180516001600160a01b031684528501517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff168584015292840192918501916001016152be565b50508583015187820388850152805180835290840192506000918401905b80831015615379578351805167ffffffffffffffff1683528501517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1685830152928401926001929092019190850190615328565b50979650505050505050565b602081526000610ab3602083018461529e565b67ffffffffffffffff8316815260608101613dba6020830184805167ffffffffffffffff908116835260209182015116910152565b634e487b7160e01b600052601160045260246000fd5b67ffffffffffffffff818116838216019080821115615404576154046153cd565b5092915050565b60006020808352606084516040808487015261542a606087018361529e565b87850151878203601f19016040890152805180835290860193506000918601905b808310156141a357845167ffffffffffffffff81511683528781015161548a89850182805167ffffffffffffffff908116835260209182015116910152565b5084015182870152938601936001929092019160809091019061544b565b6000602082840312156154ba57600080fd5b813567ffffffffffffffff8111156154d157600080fd5b6137dc84828501614856565b81810381811115610ab657610ab66153cd565b634e487b7160e01b600052601260045260246000fd5b600067ffffffffffffffff80841680615521576155216154f0565b92169190910692915050565b8082028115828204841417610ab657610ab66153cd565b60006040828403121561555657600080fd5b61555e613f2d565b61556783613fdd565b8152602083013560208201528091505092915050565b600181811c9082168061559157607f821691505b6020821081036155b157634e487b7160e01b600052602260045260246000fd5b50919050565b805182526000602067ffffffffffffffff81840151168185015260408084015160a060408701526155eb60a08701826141d4565b90506060850151868203606088015261560482826141d4565b608087810151898303918a01919091528051808352908601935060009250908501905b8083101561537957835180516001600160a01b0316835286015186830152928501926001929092019190840190615627565b602081526000610ab360208301846155b7565b60808152600061567f60808301876155b7565b61ffff9590951660208301525060408101929092526001600160a01b0316606090910152919050565b6000806000606084860312156156bd57600080fd5b83516156c881613ffa565b602085015190935067ffffffffffffffff8111156156e557600080fd5b8401601f810186136156f657600080fd5b805161570461405a82614013565b81815287602083850101111561571957600080fd5b61572a8260208301602086016141b0565b809450505050604084015190509250925092565b60006020828403121561575057600080fd5b8151613dba81613ffa565b601f821115610d15576000816000526020600020601f850160051c810160208610156157845750805b601f850160051c820191505b8181101561257557828155600101615790565b815167ffffffffffffffff8111156157bd576157bd613ea8565b6157d1816157cb845461557d565b8461575b565b602080601f83116001811461580657600084156157ee5750858301515b600019600386901b1c1916600185901b178555612575565b600085815260208120601f198616915b8281101561583557888601518255948401946001909101908401615816565b50858210156158535787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000602080835283546001600160a01b038116602085015260ff8160a01c161515604085015267ffffffffffffffff8160a81c166060850152506001808501608080860152600081546158b58161557d565b8060a089015260c060018316600081146158d657600181146158f257615922565b60ff19841660c08b015260c083151560051b8b01019450615922565b85600052602060002060005b848110156159195781548c82018501529088019089016158fe565b8b0160c0019550505b50929998505050505050505050565b80820180821115610ab657610ab66153cd565b60ff8181168382160190811115610ab657610ab66153cd565b8183823760009101908152919050565b828152606082602083013760800192915050565b600067ffffffffffffffff8084168061599c5761599c6154f0565b92169190910492915050565b60208101600583106159bc576159bc6143e5565b91905290565b60ff8181168382160290811690818114615404576154046153cd565b600060a0820160ff881683526020878185015260a0604085015281875480845260c0860191508860005282600020935060005b81811015615a365784546001600160a01b031683526001948501949284019201615a11565b50508481036060860152865180825290820192508187019060005b81811015615a765782516001600160a01b031685529383019391830191600101615a51565b50505060ff851660808501525090505b9695505050505050565b600067ffffffffffffffff808616835280851660208401525060606040830152614c4760608301846141d4565b8281526040602082015260006137dc60408301846141d4565b67ffffffffffffffff848116825283166020820152606081016137dc60408301846143fb565b615b0681846143fb565b6040602082015260006137dc60408301846141d4565b600060208284031215615b2e57600080fd5b8151613dba81613fc8565b6020815260008251610100806020850152615b586101208501836141d4565b91506020850151615b75604086018267ffffffffffffffff169052565b5060408501516001600160a01b038116606086015250606085015160808501526080850151615baf60a08601826001600160a01b03169052565b5060a0850151601f19808685030160c0870152615bcc84836141d4565b935060c08701519150808685030160e0870152615be984836141d4565b935060e0870151915080868503018387015250615a8683826141d4565b600060208284031215615c1857600080fd5b5051919050565b600060ff821660ff8103615c3557615c356153cd565b60010192915050565b848152600067ffffffffffffffff808616602084015280851660408401525060806060830152615a8660808301846141d4565b86815260c060208201526000615c8a60c08301886141d4565b6001600160a01b039690961660408301525067ffffffffffffffff9384166060820152608081019290925290911660a09091015292915050565b600082825180855260208086019550808260051b84010181860160005b84811015615d5b57601f19868403018952815160808151818652615d07828701826141d4565b9150508582015185820387870152615d1f82826141d4565b91505060408083015186830382880152615d3983826141d4565b6060948501519790940196909652505098840198925090830190600101615ce1565b5090979650505050505050565b602081526000610ab36020830184615cc4565b60008282518085526020808601955060208260051b8401016020860160005b84811015615d5b57601f19868403018952615db68383516141d4565b98840198925090830190600101615d9a565b604081526000835180516040840152602081015167ffffffffffffffff80821660608601528060408401511660808601528060608401511660a08601528060808401511660c086015250505060208401516101408060e0850152615e306101808501836141d4565b915060408601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08086850301610100870152615e6d84836141d4565b935060608801519150615e8c6101208701836001600160a01b03169052565b60808801518387015260a0880151925080868503016101608701525050615eb38282615cc4565b9150508281036020840152614c478185615d7b56fea164736f6c6343000818000a",
========
type OffRampCommitReport struct {
	PriceUpdates InternalPriceUpdates
	MerkleRoots  []OffRampMerkleRoot
>>>>>>>> 377f0dbef9 (Rename ramps and rmn (#1323)):core/gethwrappers/ccip/generated/offramp/offramp.go
}

type OffRampDynamicConfig struct {
	PriceRegistry                           common.Address
	PermissionLessExecutionThresholdSeconds uint32
	MaxTokenTransferGas                     uint32
	MaxPoolReleaseOrMintGas                 uint32
	MessageValidator                        common.Address
}

type OffRampInterval struct {
	Min uint64
	Max uint64
}

type OffRampMerkleRoot struct {
	SourceChainSelector uint64
	Interval            OffRampInterval
	MerkleRoot          [32]byte
}

type OffRampSourceChainConfig struct {
	Router    common.Address
	IsEnabled bool
	MinSeqNr  uint64
	OnRamp    []byte
}

type OffRampSourceChainConfigArgs struct {
	Router              common.Address
	SourceChainSelector uint64
	IsEnabled           bool
	OnRamp              []byte
}

type OffRampStaticConfig struct {
	ChainSelector      uint64
	RmnProxy           common.Address
	TokenAdminRegistry common.Address
	NonceManager       common.Address
}

type OffRampUnblessedRoot struct {
	SourceChainSelector uint64
	MerkleRoot          [32]byte
}

var OffRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"}],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"messageDestChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidMessageDestChainSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRoot\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidStaticConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeavesCannotBeEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"errorReason\",\"type\":\"bytes\"}],\"name\":\"MessageValidationError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReleased\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePre\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePost\",\"type\":\"uint256\"}],\"name\":\"ReleaseOrMintBalanceMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"RootAlreadyCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleCommitReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelectorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"internalType\":\"structOffRamp.Interval\",\"name\":\"interval\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structOffRamp.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structOffRamp.CommitReport\",\"name\":\"report\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"}],\"name\":\"StaticConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"isBlessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReportSingleChain[]\",\"name\":\"reports\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"gasLimitOverrides\",\"type\":\"uint256[][]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structOffRamp.UnblessedRoot[]\",\"name\":\"rootToReset\",\"type\":\"tuple[]\"}],\"name\":\"resetUnblessedRoots\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxTokenTransferGas\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162006c1a38038062006c1a8339810160408190526200003591620008c7565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001fa565b5050466080525060208301516001600160a01b03161580620000ec575060408301516001600160a01b0316155b8062000103575060608301516001600160a01b0316155b1562000122576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160401b03166000036200014e5760405163c656089560e01b815260040160405180910390fd5b82516001600160401b0390811660a052602080850180516001600160a01b0390811660c05260408088018051831660e0526060808a01805185166101005283518b519098168852945184169587019590955251821690850152905116908201527f683eb52ee924eb817377cfa8f41f238f4bb7a877da5267869dfffbad85f564d89060800160405180910390a1620001e682620002a5565b620001f181620003c1565b50505062000c67565b336001600160a01b03821603620002545760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0316620002ce576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516060808b0180516001600160a01b039b8c166001600160c01b0319909a168a17600160a01b63ffffffff98891602176001600160c01b0316600160c01b948816949094026001600160e01b031693909317600160e01b93871693909302929092179098556080808b018051600580546001600160a01b031916918d169190911790558451988952955185169688019690965290518316918601919091525116938301939093529151909216908201527fa55bd56595c45f517e5967a3067f3dca684445a3080e7c04a4e0d5a40cda627d9060a00160405180910390a150565b60005b815181101562000666576000828281518110620003e557620003e562000a1d565b60200260200101519050600081602001519050806001600160401b0316600003620004235760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b03166200044c576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b0381166000908152600660205260408120600181018054919291620004789062000a33565b80601f0160208091040260200160405190810160405280929190818152602001828054620004a69062000a33565b8015620004f75780601f10620004cb57610100808354040283529160200191620004f7565b820191906000526020600020905b815481529060010190602001808311620004d957829003601f168201915b5050505050905060008460600151905081516000036200059e57805160000362000534576040516342bcdf7f60e11b815260040160405180910390fd5b6001830162000544828262000ac4565b508254600160a81b600160e81b031916600160a81b1783556040516001600160401b03851681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1620005d9565b8080519060200120828051906020012014620005d95760405163c39a620560e01b81526001600160401b038516600482015260240162000083565b604080860151845487516001600160a01b03166001600160a01b0319921515600160a01b02929092166001600160a81b031990911617178455516001600160401b038516907f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b906200064d90869062000b90565b60405180910390a25050505050806001019050620003c4565b5050565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715620006a557620006a56200066a565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620006d657620006d66200066a565b604052919050565b80516001600160401b0381168114620006f657600080fd5b919050565b6001600160a01b03811681146200071157600080fd5b50565b805163ffffffff81168114620006f657600080fd5b6000601f83601f8401126200073d57600080fd5b825160206001600160401b03808311156200075c576200075c6200066a565b8260051b6200076d838201620006ab565b93845286810183019383810190898611156200078857600080fd5b84890192505b85831015620008ba57825184811115620007a85760008081fd5b89016080601f19828d038101821315620007c25760008081fd5b620007cc62000680565b88840151620007db81620006fb565b81526040620007ec858201620006de565b8a8301526060808601518015158114620008065760008081fd5b838301529385015193898511156200081e5760008081fd5b84860195508f603f8701126200083657600094508485fd5b8a8601519450898511156200084f576200084f6200066a565b620008608b858f88011601620006ab565b93508484528f82868801011115620008785760008081fd5b60005b8581101562000898578681018301518582018d01528b016200087b565b5060009484018b0194909452509182015283525091840191908401906200078e565b9998505050505050505050565b6000806000838503610140811215620008df57600080fd5b6080811215620008ee57600080fd5b620008f862000680565b6200090386620006de565b815260208601516200091581620006fb565b602082015260408601516200092a81620006fb565b604082015260608601516200093f81620006fb565b6060820152935060a0607f19820112156200095957600080fd5b5060405160a081016001600160401b0380821183831017156200098057620009806200066a565b81604052608087015191506200099682620006fb565b818352620009a760a0880162000714565b6020840152620009ba60c0880162000714565b6040840152620009cd60e0880162000714565b60608401526101008701519150620009e582620006fb565b608083018290526101208701519294508083111562000a0357600080fd5b505062000a138682870162000729565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b600181811c9082168062000a4857607f821691505b60208210810362000a6957634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111562000abf576000816000526020600020601f850160051c8101602086101562000a9a5750805b601f850160051c820191505b8181101562000abb5782815560010162000aa6565b5050505b505050565b81516001600160401b0381111562000ae05762000ae06200066a565b62000af88162000af1845462000a33565b8462000a6f565b602080601f83116001811462000b30576000841562000b175750858301515b600019600386901b1c1916600185901b17855562000abb565b600085815260208120601f198616915b8281101562000b615788860151825594840194600190910190840162000b40565b508582101562000b805787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b602080825282546001600160a01b0381168383015260a081901c60ff161515604084015260a81c6001600160401b0316606083015260808083015260018084018054600093929190849062000be58162000a33565b8060a089015260c0600183166000811462000c09576001811462000c265762000c58565b60ff19841660c08b015260c083151560051b8b0101945062000c58565b85600052602060002060005b8481101562000c4f5781548c820185015290880190890162000c32565b8b0160c0019550505b50929998505050505050505050565b60805160a05160c05160e05161010051615f4462000cd6600039600081816102660152612a010152600081816102370152612ef40152600081816102080152818161142b015261196d0152600081816101d801526125f00152600081816117f3015261183f0152615f446000f3fe608060405234801561001057600080fd5b50600436106101775760003560e01c806385572ffb116100d8578063ccd37ba31161008c578063f2fde38b11610066578063f2fde38b14610583578063f716f99f14610596578063ff888fb1146105a957600080fd5b8063ccd37ba31461050b578063d2a15d3514610550578063e9d68a8e1461056357600080fd5b8063991a5018116100bd578063991a5018146104c5578063a80036b4146104d8578063c673e584146104eb57600080fd5b806385572ffb1461049c5780638da5cb5b146104aa57600080fd5b8063311cd5131161012f5780635e36480c116101145780635e36480c146103785780637437ff9f1461039857806379ba50971461049457600080fd5b8063311cd513146103495780633f4b04aa1461035c57600080fd5b806306285c691161016057806306285c69146101a4578063181f5a77146102ed5780632d04ab761461033657600080fd5b806304666f9c1461017c57806305d938b514610191575b600080fd5b61018f61018a3660046140af565b6105cc565b005b61018f61019f36600461473b565b6105e0565b61029660408051608081018252600080825260208201819052918101829052606081019190915260405180608001604052807f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b6040516102e49190815167ffffffffffffffff1681526020808301516001600160a01b0390811691830191909152604080840151821690830152606092830151169181019190915260800190565b60405180910390f35b6103296040518060400160405280601181526020017f4f666652616d7020312e362e302d64657600000000000000000000000000000081525081565b6040516102e491906148b6565b61018f610344366004614961565b610785565b61018f610357366004614a14565b610b5c565b60095460405167ffffffffffffffff90911681526020016102e4565b61038b610386366004614a68565b610bc5565b6040516102e49190614ac5565b6104376040805160a081018252600080825260208201819052918101829052606081018290526080810191909152506040805160a0810182526004546001600160a01b03808216835263ffffffff600160a01b83048116602085015278010000000000000000000000000000000000000000000000008304811694840194909452600160e01b9091049092166060820152600554909116608082015290565b6040516102e49190600060a0820190506001600160a01b03808451168352602084015163ffffffff808216602086015280604087015116604086015280606087015116606086015250508060808501511660808401525092915050565b61018f610c1b565b61018f610177366004614ad3565b6000546040516001600160a01b0390911681526020016102e4565b61018f6104d3366004614b22565b610cd9565b61018f6104e6366004614b96565b610cea565b6104fe6104f9366004614c03565b61105d565b6040516102e49190614c63565b610542610519366004614cd8565b67ffffffffffffffff919091166000908152600860209081526040808320938352929052205490565b6040519081526020016102e4565b61018f61055e366004614d02565b6111bb565b610576610571366004614d77565b611275565b6040516102e49190614d92565b61018f610591366004614de0565b611382565b61018f6105a4366004614e65565b611393565b6105bc6105b7366004614fa3565b6113d5565b60405190151581526020016102e4565b6105d4611496565b6105dd816114f2565b50565b6105e86117f0565b815181518114610624576040517f83e3f56400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8181101561077557600084828151811061064357610643614fbc565b6020026020010151905060008160200151519050600085848151811061066b5761066b614fbc565b60200260200101519050805182146106af576040517f83e3f56400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b828110156107665760008282815181106106ce576106ce614fbc565b602002602001015190508060001461075d57846020015182815181106106f6576106f6614fbc565b60200260200101516080015181101561075d5784516040517fc8e9605100000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260248101839052604481018290526064015b60405180910390fd5b506001016106b2565b50505050806001019050610627565b506107808383611871565b505050565b600061079387890189615142565b805151519091501515806107ac57508051602001515115155b156108ac5760095460208a01359067ffffffffffffffff8083169116101561086b576009805467ffffffffffffffff191667ffffffffffffffff83161790556004805483516040517f3937306f0000000000000000000000000000000000000000000000000000000081526001600160a01b0390921692633937306f9261083492910161536a565b600060405180830381600087803b15801561084e57600080fd5b505af1158015610862573d6000803e3d6000fd5b505050506108aa565b8160200151516000036108aa576040517f2261116700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505b60005b816020015151811015610aa5576000826020015182815181106108d4576108d4614fbc565b602002602001015190506000816000015190506108f081611921565b60006108fb82611a23565b602084015151815491925067ffffffffffffffff908116600160a81b9092041614158061093f575060208084015190810151905167ffffffffffffffff9182169116115b1561097f57825160208401516040517feefb0cac00000000000000000000000000000000000000000000000000000000815261075492919060040161537d565b6040830151806109bb576040517f504570e300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b835167ffffffffffffffff16600090815260086020908152604080832084845290915290205415610a2e5783516040517f32cf0cbf00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260248101829052604401610754565b6020808501510151610a419060016153c8565b82547fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff16600160a81b67ffffffffffffffff9283160217909255925116600090815260086020908152604080832094835293905291909120429055506001016108af565b507f3a3950e13dd607cc37980db0ef14266c40d2bba9c01b2e44bfe549808883095d81604051610ad591906153f0565b60405180910390a1610b5160008a8a8a8a8a8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808e0282810182019093528d82529093508d92508c9182918501908490808284376000920191909152508b9250611a8a915050565b505050505050505050565b610b9c610b6b8284018461548d565b6040805160008082526020820190925290610b96565b6060815260200190600190039081610b815790505b50611871565b604080516000808252602082019092529050610bbf600185858585866000611a8a565b50505050565b6000610bd3600160046154c2565b6002610be06080856154eb565b67ffffffffffffffff16610bf49190615512565b610bfe8585611e01565b901c166003811115610c1257610c12614a9b565b90505b92915050565b6001546001600160a01b03163314610c755760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610754565b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610ce1611496565b6105dd81611e48565b333014610d23576040517f371a732800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805160008082526020820190925281610d60565b6040805180820190915260008082526020820152815260200190600190039081610d395790505b5060a08501515190915015610d9457610d918460a00151856020015186606001518760000151602001518787611fae565b90505b6040805160a0810182528551518152855160209081015167ffffffffffffffff1681830152808701518351600094840192610dd09291016148b6565b60408051601f19818403018152918152908252878101516020830152018390526005549091506001600160a01b03168015610edd576040517f08d450a10000000000000000000000000000000000000000000000000000000081526001600160a01b038216906308d450a190610e4a9085906004016155cb565b600060405180830381600087803b158015610e6457600080fd5b505af1925050508015610e75575060015b610edd573d808015610ea3576040519150601f19603f3d011682016040523d82523d6000602084013e610ea8565b606091505b50806040517f09c2532500000000000000000000000000000000000000000000000000000000815260040161075491906148b6565b604086015151158015610ef257506080860151155b80610f09575060608601516001600160a01b03163b155b80610f4957506060860151610f47906001600160a01b03167f85572ffb000000000000000000000000000000000000000000000000000000006120cd565b155b15610f5657505050505050565b855160209081015167ffffffffffffffff1660009081526006909152604080822054608089015160608a015192517f3cf9798300000000000000000000000000000000000000000000000000000000815284936001600160a01b0390931692633cf9798392610fce92899261138892916004016155de565b6000604051808303816000875af1158015610fed573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611015919081019061561a565b50915091508161105357806040517f0a8d6e8c00000000000000000000000000000000000000000000000000000000815260040161075491906148b6565b5050505050505050565b6110a06040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c08201529485529182018054845181840281018401909552808552929385830193909283018282801561114957602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161112b575b50505050508152602001600382018054806020026020016040519081016040528092919081815260200182805480156111ab57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161118d575b5050505050815250509050919050565b6111c3611496565b60005b818110156107805760008383838181106111e2576111e2614fbc565b9050604002018036038101906111f891906156b0565b905061120781602001516113d5565b61126c57805167ffffffffffffffff1660009081526008602090815260408083208285018051855290835281842093909355915191519182527f202f1139a3e334b6056064c0e9b19fd07e44a88d8f6e5ded571b24cf8c371f12910160405180910390a15b506001016111c6565b604080516080808201835260008083526020808401829052838501829052606080850181905267ffffffffffffffff878116845260068352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b909204909216948301949094526001840180549394929391840191611302906156e9565b80601f016020809104026020016040519081016040528092919081815260200182805461132e906156e9565b80156111ab5780601f10611350576101008083540402835291602001916111ab565b820191906000526020600020905b81548152906001019060200180831161135e57505050919092525091949350505050565b61138a611496565b6105dd816120e9565b61139b611496565b60005b81518110156113d1576113c98282815181106113bc576113bc614fbc565b602002602001015161219f565b60010161139e565b5050565b6040805180820182523081526020810183815291517f4d61677100000000000000000000000000000000000000000000000000000000815290516001600160a01b039081166004830152915160248201526000917f00000000000000000000000000000000000000000000000000000000000000001690634d61677190604401602060405180830381865afa158015611472573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c159190615723565b6000546001600160a01b031633146114f05760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610754565b565b60005b81518110156113d157600082828151811061151257611512614fbc565b602002602001015190506000816020015190508067ffffffffffffffff16600003611569576040517fc656089500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81516001600160a01b0316611591576040516342bcdf7f60e11b815260040160405180910390fd5b67ffffffffffffffff811660009081526006602052604081206001810180549192916115bc906156e9565b80601f01602080910402602001604051908101604052809291908181526020018280546115e8906156e9565b80156116355780601f1061160a57610100808354040283529160200191611635565b820191906000526020600020905b81548152906001019060200180831161161857829003601f168201915b5050505050905060008460600151905081516000036116ed578051600003611670576040516342bcdf7f60e11b815260040160405180910390fd5b6001830161167e8282615788565b5082547fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff16600160a81b17835560405167ffffffffffffffff851681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1611740565b8080519060200120828051906020012014611740576040517fc39a620500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff85166004820152602401610754565b604080860151845487516001600160a01b031673ffffffffffffffffffffffffffffffffffffffff19921515600160a01b02929092167fffffffffffffffffffffff000000000000000000000000000000000000000000909116171784555167ffffffffffffffff8516907f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b906117d8908690615848565b60405180910390a250505050508060010190506114f5565b467f0000000000000000000000000000000000000000000000000000000000000000146114f0576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f00000000000000000000000000000000000000000000000000000000000000006004820152466024820152604401610754565b81516000036118ab576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160408051600080825260208201909252911591905b845181101561191a576119128582815181106118e0576118e0614fbc565b60200260200101518461190c578583815181106118ff576118ff614fbc565b60200260200101516124e3565b836124e3565b6001016118c2565b5050505050565b6040517f2cbc26bb000000000000000000000000000000000000000000000000000000008152608082901b77ffffffffffffffff000000000000000000000000000000001660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa1580156119bc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119e09190615723565b156105dd576040517ffdbd6a7200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610754565b67ffffffffffffffff811660009081526006602052604081208054600160a01b900460ff16610c15576040517fed053c5900000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610754565b60ff87811660009081526002602090815260408083208151608081018352815481526001909101548086169382019390935261010083048516918101919091526201000090910490921615156060830152873590611ae98760a4615916565b9050826060015115611b31578451611b02906020615512565b8651611b0f906020615512565b611b1a9060a0615916565b611b249190615916565b611b2e9082615916565b90505b368114611b73576040517f8e1192e100000000000000000000000000000000000000000000000000000000815260048101829052366024820152604401610754565b5081518114611bbb5781516040517f93df584c000000000000000000000000000000000000000000000000000000008152600481019190915260248101829052604401610754565b611bc36117f0565b60ff808a1660009081526003602090815260408083203384528252808320815180830190925280548086168352939491939092840191610100909104166002811115611c1157611c11614a9b565b6002811115611c2257611c22614a9b565b9052509050600281602001516002811115611c3f57611c3f614a9b565b148015611c935750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff1681548110611c7b57611c7b614fbc565b6000918252602090912001546001600160a01b031633145b611cc9576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50816060015115611dab576020820151611ce4906001615929565b60ff16855114611d20576040517f71253a2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8351855114611d5b576040517fa75d88af00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008787604051611d6d929190615942565b604051908190038120611d84918b90602001615952565b604051602081830303815290604052805190602001209050611da98a82888888612c86565b505b6040805182815260208a81013567ffffffffffffffff169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b67ffffffffffffffff8216600090815260076020526040812081611e26608085615966565b67ffffffffffffffff1681526020810191909152604001600020549392505050565b80516001600160a01b0316611e70576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516060808b0180516001600160a01b039b8c167fffffffffffffffff000000000000000000000000000000000000000000000000909a168a17600160a01b63ffffffff988916021777ffffffffffffffffffffffffffffffffffffffffffffffff167801000000000000000000000000000000000000000000000000948816949094026001600160e01b031693909317600160e01b93871693909302929092179098556080808b0180516005805473ffffffffffffffffffffffffffffffffffffffff1916918d169190911790558451988952955185169688019690965290518316918601919091525116938301939093529151909216908201527fa55bd56595c45f517e5967a3067f3dca684445a3080e7c04a4e0d5a40cda627d9060a00160405180910390a150565b6060865167ffffffffffffffff811115611fca57611fca613ec6565b60405190808252806020026020018201604052801561200f57816020015b6040805180820190915260008082526020820152815260200190600190039081611fe85790505b50905060005b87518110156120c15761209c88828151811061203357612033614fbc565b602002602001015188888888888781811061205057612050614fbc565b9050602002810190612062919061598d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612e9392505050565b8282815181106120ae576120ae614fbc565b6020908102919091010152600101612015565b505b9695505050505050565b60006120d883613238565b8015610c125750610c12838361329c565b336001600160a01b038216036121415760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610754565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff166000036121ca576000604051631b3fab5160e11b815260040161075491906159f2565b60208082015160ff8082166000908152600290935260408320600181015492939092839216900361223757606084015160018201805491151562010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffff90921691909117905561228c565b6060840151600182015460ff620100009091041615159015151461228c576040517f87f6037c00000000000000000000000000000000000000000000000000000000815260ff84166004820152602401610754565b60a08401518051601f60ff821611156122bb576001604051631b3fab5160e11b815260040161075491906159f2565b612321858560030180548060200260200160405190810160405280929190818152602001828054801561231757602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116122f9575b5050505050613357565b8560600151156124505761238f8585600201805480602002602001604051908101604052809291908181526020018280548015612317576020028201919060005260206000209081546001600160a01b031681526001909101906020018083116122f9575050505050613357565b608086015180516123a99060028701906020840190613e20565b5080516001850180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010060ff841690810291909117909155601f1015612409576002604051631b3fab5160e11b815260040161075491906159f2565b6040880151612419906003615a0c565b60ff168160ff1611612441576003604051631b3fab5160e11b815260040161075491906159f2565b61244d878360016133c0565b50505b61245c858360026133c0565b81516124719060038601906020850190613e20565b5060408681015160018501805460ff191660ff8316179055875180865560a089015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547936124ca938a939260028b01929190615a28565b60405180910390a16124db85613540565b505050505050565b81516124ee81611921565b60006124f982611a23565b602085015151909150600081900361253c576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b846040015151811461257a576040517f57e0e08300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008167ffffffffffffffff81111561259557612595613ec6565b6040519080825280602002602001820160405280156125be578160200160208202803683370190505b50905060005b82811015612733576000876020015182815181106125e4576125e4614fbc565b602002602001015190507f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681600001516040015167ffffffffffffffff161461267757805160409081015190517f38432a2200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610754565b61270d8186600101805461268a906156e9565b80601f01602080910402602001604051908101604052809291908181526020018280546126b6906156e9565b80156127035780601f106126d857610100808354040283529160200191612703565b820191906000526020600020905b8154815290600101906020018083116126e657829003601f168201915b505050505061355c565b83838151811061271f5761271f614fbc565b6020908102919091010152506001016125c4565b50600061274a858389606001518a6080015161367e565b905080600003612792576040517f7dd17a7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff86166004820152602401610754565b8551151560005b84811015610b515760005a905060008a6020015183815181106127be576127be614fbc565b6020026020010151905060006127dc8a836000015160600151610bc5565b905060008160038111156127f2576127f2614a9b565b148061280f5750600381600381111561280d5761280d614a9b565b145b612867578151606001516040805167ffffffffffffffff808e16825290921660208301527f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c91015b60405180910390a1505050612c7e565b841561293757600454600090600160a01b900463ffffffff1661288a88426154c2565b11905080806128aa575060038260038111156128a8576128a8614a9b565b145b6128ec576040517fa9cfc86200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8c166004820152602401610754565b8b85815181106128fe576128fe614fbc565b6020026020010151600014612931578b858151811061291f5761291f614fbc565b60200260200101518360800181815250505b50612998565b600081600381111561294b5761294b614a9b565b14612998578151606001516040805167ffffffffffffffff808e16825290921660208301527f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe651209101612857565b81516080015167ffffffffffffffff1615612a875760008160038111156129c1576129c1614a9b565b03612a875781516080015160208301516040517fe0e03cae0000000000000000000000000000000000000000000000000000000081526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263e0e03cae92612a38928f929190600401615ad4565b6020604051808303816000875af1158015612a57573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a7b9190615723565b612a8757505050612c7e565b60008c604001518581518110612a9f57612a9f614fbc565b6020026020010151905080518360a001515114612b03578251606001516040517f1cfe6d8b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff808e1660048301529091166024820152604401610754565b612b178b84600001516060015160016136d4565b600080612b24858461377c565b91509150612b3b8d866000015160600151846136d4565b8715612bab576003826003811115612b5557612b55614a9b565b03612bab576000846003811115612b6e57612b6e614a9b565b14612bab578451516040517f2b11b8d900000000000000000000000000000000000000000000000000000000815261075491908390600401615b01565b6002826003811115612bbf57612bbf614a9b565b14612c19576003826003811115612bd857612bd8614a9b565b14612c19578451606001516040517f926c5a3e000000000000000000000000000000000000000000000000000000008152610754918f918590600401615b1a565b8451805160609091015167ffffffffffffffff908116908f167fdc8ccbc35e0eebd81239bcd1971fcd53c7eb34034880142a0f43c809a458732f85855a612c60908d6154c2565b604051612c6f93929190615b40565b60405180910390a45050505050505b600101612799565b612c8e613e92565b835160005b81811015611053576000600188868460208110612cb257612cb2614fbc565b612cbf91901a601b615929565b898581518110612cd157612cd1614fbc565b6020026020010151898681518110612ceb57612ceb614fbc565b602002602001015160405160008152602001604052604051612d29949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612d4b573d6000803e3d6000fd5b505060408051601f1981015160ff808e166000908152600360209081528582206001600160a01b03851683528152858220858701909652855480841686529397509095509293928401916101009004166002811115612dac57612dac614a9b565b6002811115612dbd57612dbd614a9b565b9052509050600181602001516002811115612dda57612dda614a9b565b14612e11576040517fca31867a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051859060ff16601f8110612e2857612e28614fbc565b602002015115612e64576040517ff67bc7c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600185826000015160ff16601f8110612e7f57612e7f614fbc565b911515602090920201525050600101612c93565b60408051808201909152600080825260208201526000612eb68760200151613846565b6040517fbbe4f6db0000000000000000000000000000000000000000000000000000000081526001600160a01b0380831660048301529192506000917f0000000000000000000000000000000000000000000000000000000000000000169063bbe4f6db90602401602060405180830381865afa158015612f3b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f5f9190615b70565b90506001600160a01b0381161580612fa75750612fa56001600160a01b0382167faff2afbf000000000000000000000000000000000000000000000000000000006120cd565b155b15612fe9576040517fae9b4ce90000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602401610754565b600454600090819061300b9089908690600160e01b900463ffffffff166138ec565b9150915060008060006130d86040518061010001604052808e81526020018c67ffffffffffffffff1681526020018d6001600160a01b031681526020018f606001518152602001896001600160a01b031681526020018f6000015181526020018f6040015181526020018b8152506040516024016130899190615b8d565b60408051601f198184030181529190526020810180516001600160e01b03167f390775370000000000000000000000000000000000000000000000000000000017905287866113886084613a1a565b9250925092508261311757816040517fe1cd550900000000000000000000000000000000000000000000000000000000815260040161075491906148b6565b815160201461315f5781516040517f78ef8024000000000000000000000000000000000000000000000000000000008152602060048201526024810191909152604401610754565b6000828060200190518101906131759190615c5a565b9050866001600160a01b03168c6001600160a01b03161461320a5760006131a68d8a6131a1868a6154c2565b6138ec565b509050868110806131c05750816131bd88836154c2565b14155b15613208576040517fa966e21f000000000000000000000000000000000000000000000000000000008152600481018390526024810188905260448101829052606401610754565b505b604080518082019091526001600160a01b039098168852602088015250949550505050505095945050505050565b6000613264827f01ffc9a70000000000000000000000000000000000000000000000000000000061329c565b8015610c155750613295827fffffffff0000000000000000000000000000000000000000000000000000000061329c565b1592915050565b6040517fffffffff0000000000000000000000000000000000000000000000000000000082166024820152600090819060440160408051601f19818403018152919052602080820180516001600160e01b03167f01ffc9a700000000000000000000000000000000000000000000000000000000178152825192935060009283928392909183918a617530fa92503d91506000519050828015613340575060208210155b801561334c5750600081115b979650505050505050565b60005b81518110156107805760ff83166000908152600360205260408120835190919084908490811061338c5761338c614fbc565b6020908102919091018101516001600160a01b03168252810191909152604001600020805461ffff1916905560010161335a565b60005b82518160ff161015610bbf576000838260ff16815181106133e6576133e6614fbc565b602002602001015190506000600281111561340357613403614a9b565b60ff80871660009081526003602090815260408083206001600160a01b0387168452909152902054610100900416600281111561344257613442614a9b565b14613463576004604051631b3fab5160e11b815260040161075491906159f2565b6001600160a01b0381166134a3576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405180604001604052808360ff1681526020018460028111156134c9576134c9614a9b565b905260ff80871660009081526003602090815260408083206001600160a01b0387168452825290912083518154931660ff198416811782559184015190929091839161ffff19161761010083600281111561352657613526614a9b565b0217905550905050508061353990615c73565b90506133c3565b60ff81166105dd576009805467ffffffffffffffff1916905550565b8151602080820151604092830151925160009384936135a2937f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f93909291889101615c92565b60408051601f1981840301815290829052805160209182012086518051888401516060808b0151908401516080808d015195015195976135eb9794969395929491939101615cc5565b604051602081830303815290604052805190602001208560400151805190602001208660a001516040516020016136229190615dd7565b60408051601f198184030181528282528051602091820120908301969096528101939093526060830191909152608082015260a081019190915260c0015b60405160208183030381529060405280519060200120905092915050565b60008061368c858585613b40565b9050613697816113d5565b6136a55760009150506136cc565b67ffffffffffffffff86166000908152600860209081526040808320938352929052205490505b949350505050565b600060026136e36080856154eb565b67ffffffffffffffff166136f79190615512565b905060006137058585611e01565b905081613714600160046154c2565b901b19168183600381111561372b5761372b614a9b565b67ffffffffffffffff871660009081526007602052604081209190921b9290921791829161375a608088615966565b67ffffffffffffffff1681526020810191909152604001600020555050505050565b6040517fa80036b4000000000000000000000000000000000000000000000000000000008152600090606090309063a80036b4906137c09087908790600401615e37565b600060405180830381600087803b1580156137da57600080fd5b505af19250505080156137eb575060015b61382a573d808015613819576040519150601f19603f3d011682016040523d82523d6000602084013e61381e565b606091505b5060039250905061383f565b50506040805160208101909152600081526002905b9250929050565b6000815160201461388557816040517f8d666f6000000000000000000000000000000000000000000000000000000000815260040161075491906148b6565b60008280602001905181019061389b9190615c5a565b90506001600160a01b038111806138b3575061040081105b15610c1557826040517f8d666f6000000000000000000000000000000000000000000000000000000000815260040161075491906148b6565b60008060008060006139668860405160240161391791906001600160a01b0391909116815260200190565b60408051601f198184030181529190526020810180516001600160e01b03167f70a082310000000000000000000000000000000000000000000000000000000017905288886113886084613a1a565b925092509250826139a557816040517fe1cd550900000000000000000000000000000000000000000000000000000000815260040161075491906148b6565b60208251146139ed5781516040517f78ef8024000000000000000000000000000000000000000000000000000000008152602060048201526024810191909152604401610754565b81806020019051810190613a019190615c5a565b613a0b82886154c2565b94509450505050935093915050565b6000606060008361ffff1667ffffffffffffffff811115613a3d57613a3d613ec6565b6040519080825280601f01601f191660200182016040528015613a67576020820181803683370190505b509150863b613a9a577f0c3b563c0000000000000000000000000000000000000000000000000000000060005260046000fd5b5a85811015613acd577fafa32a2c0000000000000000000000000000000000000000000000000000000060005260046000fd5b8590036040810481038710613b06577f37c3be290000000000000000000000000000000000000000000000000000000060005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d84811115613b295750835b808352806000602085013e50955095509592505050565b8251825160009190818303613b81576040517f11a6b26400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101018211801590613b9557506101018111155b613bb2576040516309bde33960e01b815260040160405180910390fd5b60001982820101610100811115613bdc576040516309bde33960e01b815260040160405180910390fd5b80600003613c095786600081518110613bf757613bf7614fbc565b60200260200101519350505050613dd8565b60008167ffffffffffffffff811115613c2457613c24613ec6565b604051908082528060200260200182016040528015613c4d578160200160208202803683370190505b50905060008080805b85811015613d775760006001821b8b811603613cb15788851015613c9a578c5160018601958e918110613c8b57613c8b614fbc565b60200260200101519050613cd3565b8551600185019487918110613c8b57613c8b614fbc565b8b5160018401938d918110613cc857613cc8614fbc565b602002602001015190505b600089861015613d03578d5160018701968f918110613cf457613cf4614fbc565b60200260200101519050613d25565b8651600186019588918110613d1a57613d1a614fbc565b602002602001015190505b82851115613d46576040516309bde33960e01b815260040160405180910390fd5b613d508282613ddf565b878481518110613d6257613d62614fbc565b60209081029190910101525050600101613c56565b506001850382148015613d8957508683145b8015613d9457508581145b613db1576040516309bde33960e01b815260040160405180910390fd5b836001860381518110613dc657613dc6614fbc565b60200260200101519750505050505050505b9392505050565b6000818310613df757613df28284613dfd565b610c12565b610c1283835b604080516001602082015290810183905260608101829052600090608001613660565b828054828255906000526020600020908101928215613e82579160200282015b82811115613e82578251825473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03909116178255602090920191600190910190613e40565b50613e8e929150613eb1565b5090565b604051806103e00160405280601f906020820280368337509192915050565b5b80821115613e8e5760008155600101613eb2565b634e487b7160e01b600052604160045260246000fd5b6040516080810167ffffffffffffffff81118282101715613eff57613eff613ec6565b60405290565b60405160a0810167ffffffffffffffff81118282101715613eff57613eff613ec6565b60405160c0810167ffffffffffffffff81118282101715613eff57613eff613ec6565b6040805190810167ffffffffffffffff81118282101715613eff57613eff613ec6565b6040516060810167ffffffffffffffff81118282101715613eff57613eff613ec6565b604051601f8201601f1916810167ffffffffffffffff81118282101715613fba57613fba613ec6565b604052919050565b600067ffffffffffffffff821115613fdc57613fdc613ec6565b5060051b60200190565b6001600160a01b03811681146105dd57600080fd5b803567ffffffffffffffff8116811461401357600080fd5b919050565b80151581146105dd57600080fd5b803561401381614018565b600067ffffffffffffffff82111561404b5761404b613ec6565b50601f01601f191660200190565b600082601f83011261406a57600080fd5b813561407d61407882614031565b613f91565b81815284602083860101111561409257600080fd5b816020850160208301376000918101602001919091529392505050565b600060208083850312156140c257600080fd5b823567ffffffffffffffff808211156140da57600080fd5b818501915085601f8301126140ee57600080fd5b81356140fc61407882613fc2565b81815260059190911b8301840190848101908883111561411b57600080fd5b8585015b838110156141c1578035858111156141375760008081fd5b86016080818c03601f190181131561414f5760008081fd5b614157613edc565b8983013561416481613fe6565b81526040614173848201613ffb565b8b83015260608085013561418681614018565b8383015292840135928984111561419f57600091508182fd5b6141ad8f8d86880101614059565b90830152508552505091860191860161411f565b5098975050505050505050565b600060a082840312156141e057600080fd5b6141e8613f05565b9050813581526141fa60208301613ffb565b602082015261420b60408301613ffb565b604082015261421c60608301613ffb565b606082015261422d60808301613ffb565b608082015292915050565b803561401381613fe6565b600082601f83011261425457600080fd5b8135602061426461407883613fc2565b82815260059290921b8401810191818101908684111561428357600080fd5b8286015b848110156120c157803567ffffffffffffffff808211156142a85760008081fd5b818901915060a080601f19848d030112156142c35760008081fd5b6142cb613f05565b87840135838111156142dd5760008081fd5b6142eb8d8a83880101614059565b825250604080850135848111156143025760008081fd5b6143108e8b83890101614059565b8a84015250606080860135858111156143295760008081fd5b6143378f8c838a0101614059565b8385015250608091508186013581840152508285013592508383111561435d5760008081fd5b61436b8d8a85880101614059565b908201528652505050918301918301614287565b6000610140828403121561439257600080fd5b61439a613f28565b90506143a683836141ce565b815260a082013567ffffffffffffffff808211156143c357600080fd5b6143cf85838601614059565b602084015260c08401359150808211156143e857600080fd5b6143f485838601614059565b604084015261440560e08501614238565b6060840152610100840135608084015261012084013591508082111561442a57600080fd5b5061443784828501614243565b60a08301525092915050565b600082601f83011261445457600080fd5b8135602061446461407883613fc2565b82815260059290921b8401810191818101908684111561448357600080fd5b8286015b848110156120c157803567ffffffffffffffff8111156144a75760008081fd5b6144b58986838b010161437f565b845250918301918301614487565b600082601f8301126144d457600080fd5b813560206144e461407883613fc2565b82815260059290921b8401810191818101908684111561450357600080fd5b8286015b848110156120c157803567ffffffffffffffff8082111561452757600080fd5b818901915089603f83011261453b57600080fd5b8582013561454b61407882613fc2565b81815260059190911b830160400190878101908c83111561456b57600080fd5b604085015b838110156145a45780358581111561458757600080fd5b6145968f6040838a0101614059565b845250918901918901614570565b50875250505092840192508301614507565b600082601f8301126145c757600080fd5b813560206145d761407883613fc2565b8083825260208201915060208460051b8701019350868411156145f957600080fd5b602086015b848110156120c157803583529183019183016145fe565b600082601f83011261462657600080fd5b8135602061463661407883613fc2565b82815260059290921b8401810191818101908684111561465557600080fd5b8286015b848110156120c157803567ffffffffffffffff8082111561467a5760008081fd5b818901915060a080601f19848d030112156146955760008081fd5b61469d613f05565b6146a8888501613ffb565b8152604080850135848111156146be5760008081fd5b6146cc8e8b83890101614443565b8a84015250606080860135858111156146e55760008081fd5b6146f38f8c838a01016144c3565b838501525060809150818601358581111561470e5760008081fd5b61471c8f8c838a01016145b6565b9184019190915250919093013590830152508352918301918301614659565b600080604080848603121561474f57600080fd5b833567ffffffffffffffff8082111561476757600080fd5b61477387838801614615565b945060209150818601358181111561478a57600080fd5b8601601f8101881361479b57600080fd5b80356147a961407882613fc2565b81815260059190911b8201840190848101908a8311156147c857600080fd5b8584015b83811015614854578035868111156147e45760008081fd5b8501603f81018d136147f65760008081fd5b8781013561480661407882613fc2565b81815260059190911b82018a0190898101908f8311156148265760008081fd5b928b01925b828410156148445783358252928a0192908a019061482b565b86525050509186019186016147cc565b50809750505050505050509250929050565b60005b83811015614881578181015183820152602001614869565b50506000910152565b600081518084526148a2816020860160208601614866565b601f01601f19169290920160200192915050565b602081526000610c12602083018461488a565b8060608101831015610c1557600080fd5b60008083601f8401126148ec57600080fd5b50813567ffffffffffffffff81111561490457600080fd5b60208301915083602082850101111561383f57600080fd5b60008083601f84011261492e57600080fd5b50813567ffffffffffffffff81111561494657600080fd5b6020830191508360208260051b850101111561383f57600080fd5b60008060008060008060008060e0898b03121561497d57600080fd5b6149878a8a6148c9565b9750606089013567ffffffffffffffff808211156149a457600080fd5b6149b08c838d016148da565b909950975060808b01359150808211156149c957600080fd5b6149d58c838d0161491c565b909750955060a08b01359150808211156149ee57600080fd5b506149fb8b828c0161491c565b999c989b50969995989497949560c00135949350505050565b600080600060808486031215614a2957600080fd5b614a3385856148c9565b9250606084013567ffffffffffffffff811115614a4f57600080fd5b614a5b868287016148da565b9497909650939450505050565b60008060408385031215614a7b57600080fd5b614a8483613ffb565b9150614a9260208401613ffb565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b60048110614ac157614ac1614a9b565b9052565b60208101610c158284614ab1565b600060208284031215614ae557600080fd5b813567ffffffffffffffff811115614afc57600080fd5b820160a08185031215613dd857600080fd5b803563ffffffff8116811461401357600080fd5b600060a08284031215614b3457600080fd5b614b3c613f05565b8235614b4781613fe6565b8152614b5560208401614b0e565b6020820152614b6660408401614b0e565b6040820152614b7760608401614b0e565b60608201526080830135614b8a81613fe6565b60808201529392505050565b600080600060408486031215614bab57600080fd5b833567ffffffffffffffff80821115614bc357600080fd5b614bcf8783880161437f565b94506020860135915080821115614be557600080fd5b50614a5b8682870161491c565b803560ff8116811461401357600080fd5b600060208284031215614c1557600080fd5b610c1282614bf2565b60008151808452602080850194506020840160005b83811015614c585781516001600160a01b031687529582019590820190600101614c33565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a0840152614cb260e0840182614c1e565b90506040840151601f198483030160c0850152614ccf8282614c1e565b95945050505050565b60008060408385031215614ceb57600080fd5b614cf483613ffb565b946020939093013593505050565b60008060208385031215614d1557600080fd5b823567ffffffffffffffff80821115614d2d57600080fd5b818501915085601f830112614d4157600080fd5b813581811115614d5057600080fd5b8660208260061b8501011115614d6557600080fd5b60209290920196919550909350505050565b600060208284031215614d8957600080fd5b610c1282613ffb565b602081526001600160a01b03825116602082015260208201511515604082015267ffffffffffffffff6040830151166060820152600060608301516080808401526136cc60a084018261488a565b600060208284031215614df257600080fd5b8135613dd881613fe6565b600082601f830112614e0e57600080fd5b81356020614e1e61407883613fc2565b8083825260208201915060208460051b870101935086841115614e4057600080fd5b602086015b848110156120c1578035614e5881613fe6565b8352918301918301614e45565b60006020808385031215614e7857600080fd5b823567ffffffffffffffff80821115614e9057600080fd5b818501915085601f830112614ea457600080fd5b8135614eb261407882613fc2565b81815260059190911b83018401908481019088831115614ed157600080fd5b8585015b838110156141c157803585811115614eec57600080fd5b860160c0818c03601f19011215614f035760008081fd5b614f0b613f28565b8882013581526040614f1e818401614bf2565b8a8301526060614f2f818501614bf2565b8284015260809150614f42828501614026565b9083015260a08381013589811115614f5a5760008081fd5b614f688f8d83880101614dfd565b838501525060c0840135915088821115614f825760008081fd5b614f908e8c84870101614dfd565b9083015250845250918601918601614ed5565b600060208284031215614fb557600080fd5b5035919050565b634e487b7160e01b600052603260045260246000fd5b80356001600160e01b038116811461401357600080fd5b600082601f830112614ffa57600080fd5b8135602061500a61407883613fc2565b82815260069290921b8401810191818101908684111561502957600080fd5b8286015b848110156120c157604081890312156150465760008081fd5b61504e613f4b565b61505782613ffb565b8152615064858301614fd2565b8186015283529183019160400161502d565b600082601f83011261508757600080fd5b8135602061509761407883613fc2565b82815260079290921b840181019181810190868411156150b657600080fd5b8286015b848110156120c15780880360808112156150d45760008081fd5b6150dc613f6e565b6150e583613ffb565b8152604080601f19840112156150fb5760008081fd5b615103613f4b565b9250615110878501613ffb565b835261511d818501613ffb565b83880152818701929092526060830135918101919091528352918301916080016150ba565b6000602080838503121561515557600080fd5b823567ffffffffffffffff8082111561516d57600080fd5b8185019150604080838803121561518357600080fd5b61518b613f4b565b83358381111561519a57600080fd5b84016040818a0312156151ac57600080fd5b6151b4613f4b565b8135858111156151c357600080fd5b8201601f81018b136151d457600080fd5b80356151e261407882613fc2565b81815260069190911b8201890190898101908d83111561520157600080fd5b928a01925b828410156152515787848f03121561521e5760008081fd5b615226613f4b565b843561523181613fe6565b815261523e858d01614fd2565b818d0152825292870192908a0190615206565b84525050508187013593508484111561526957600080fd5b6152758a858401614fe9565b818801528252508385013591508282111561528f57600080fd5b61529b88838601615076565b85820152809550505050505092915050565b805160408084528151848201819052600092602091908201906060870190855b8181101561530457835180516001600160a01b031684528501516001600160e01b03168584015292840192918501916001016152cd565b50508583015187820388850152805180835290840192506000918401905b8083101561535e578351805167ffffffffffffffff1683528501516001600160e01b031685830152928401926001929092019190850190615322565b50979650505050505050565b602081526000610c1260208301846152ad565b67ffffffffffffffff8316815260608101613dd86020830184805167ffffffffffffffff908116835260209182015116910152565b634e487b7160e01b600052601160045260246000fd5b67ffffffffffffffff8181168382160190808211156153e9576153e96153b2565b5092915050565b60006020808352606084516040808487015261540f60608701836152ad565b87850151878203601f19016040890152805180835290860193506000918601905b808310156141c157845167ffffffffffffffff81511683528781015161546f89850182805167ffffffffffffffff908116835260209182015116910152565b50840151828701529386019360019290920191608090910190615430565b60006020828403121561549f57600080fd5b813567ffffffffffffffff8111156154b657600080fd5b6136cc84828501614615565b81810381811115610c1557610c156153b2565b634e487b7160e01b600052601260045260246000fd5b600067ffffffffffffffff80841680615506576155066154d5565b92169190910692915050565b8082028115828204841417610c1557610c156153b2565b805182526000602067ffffffffffffffff81840151168185015260408084015160a0604087015261555d60a087018261488a565b905060608501518682036060880152615576828261488a565b608087810151898303918a01919091528051808352908601935060009250908501905b8083101561535e57835180516001600160a01b0316835286015186830152928501926001929092019190840190615599565b602081526000610c126020830184615529565b6080815260006155f16080830187615529565b61ffff9590951660208301525060408101929092526001600160a01b0316606090910152919050565b60008060006060848603121561562f57600080fd5b835161563a81614018565b602085015190935067ffffffffffffffff81111561565757600080fd5b8401601f8101861361566857600080fd5b805161567661407882614031565b81815287602083850101111561568b57600080fd5b61569c826020830160208601614866565b809450505050604084015190509250925092565b6000604082840312156156c257600080fd5b6156ca613f4b565b6156d383613ffb565b8152602083013560208201528091505092915050565b600181811c908216806156fd57607f821691505b60208210810361571d57634e487b7160e01b600052602260045260246000fd5b50919050565b60006020828403121561573557600080fd5b8151613dd881614018565b601f821115610780576000816000526020600020601f850160051c810160208610156157695750805b601f850160051c820191505b818110156124db57828155600101615775565b815167ffffffffffffffff8111156157a2576157a2613ec6565b6157b6816157b084546156e9565b84615740565b602080601f8311600181146157eb57600084156157d35750858301515b600019600386901b1c1916600185901b1785556124db565b600085815260208120601f198616915b8281101561581a578886015182559484019460019091019084016157fb565b50858210156158385787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000602080835283546001600160a01b038116602085015260ff8160a01c161515604085015267ffffffffffffffff8160a81c1660608501525060018085016080808601526000815461589a816156e9565b8060a089015260c060018316600081146158bb57600181146158d757615907565b60ff19841660c08b015260c083151560051b8b01019450615907565b85600052602060002060005b848110156158fe5781548c82018501529088019089016158e3565b8b0160c0019550505b50929998505050505050505050565b80820180821115610c1557610c156153b2565b60ff8181168382160190811115610c1557610c156153b2565b8183823760009101908152919050565b828152606082602083013760800192915050565b600067ffffffffffffffff80841680615981576159816154d5565b92169190910492915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126159c257600080fd5b83018035915067ffffffffffffffff8211156159dd57600080fd5b60200191503681900382131561383f57600080fd5b6020810160058310615a0657615a06614a9b565b91905290565b60ff81811683821602908116908181146153e9576153e96153b2565b600060a0820160ff881683526020878185015260a0604085015281875480845260c0860191508860005282600020935060005b81811015615a805784546001600160a01b031683526001948501949284019201615a5b565b50508481036060860152865180825290820192508187019060005b81811015615ac05782516001600160a01b031685529383019391830191600101615a9b565b50505060ff851660808501525090506120c3565b600067ffffffffffffffff808616835280851660208401525060606040830152614ccf606083018461488a565b8281526040602082015260006136cc604083018461488a565b67ffffffffffffffff848116825283166020820152606081016136cc6040830184614ab1565b615b4a8185614ab1565b606060208201526000615b60606083018561488a565b9050826040830152949350505050565b600060208284031215615b8257600080fd5b8151613dd881613fe6565b6020815260008251610100806020850152615bac61012085018361488a565b91506020850151615bc9604086018267ffffffffffffffff169052565b5060408501516001600160a01b038116606086015250606085015160808501526080850151615c0360a08601826001600160a01b03169052565b5060a0850151601f19808685030160c0870152615c20848361488a565b935060c08701519150808685030160e0870152615c3d848361488a565b935060e08701519150808685030183870152506120c3838261488a565b600060208284031215615c6c57600080fd5b5051919050565b600060ff821660ff8103615c8957615c896153b2565b60010192915050565b848152600067ffffffffffffffff8086166020840152808516604084015250608060608301526120c3608083018461488a565b86815260c060208201526000615cde60c083018861488a565b6001600160a01b039690961660408301525067ffffffffffffffff9384166060820152608081019290925290911660a09091015292915050565b600082825180855260208086019550808260051b84010181860160005b84811015615dca57601f19868403018952815160a08151818652615d5b8287018261488a565b9150508582015185820387870152615d73828261488a565b91505060408083015186830382880152615d8d838261488a565b92505050606080830151818701525060808083015192508582038187015250615db6818361488a565b9a86019a9450505090830190600101615d35565b5090979650505050505050565b602081526000610c126020830184615d18565b60008282518085526020808601955060208260051b8401016020860160005b84811015615dca57601f19868403018952615e2583835161488a565b98840198925090830190600101615e09565b604081526000835180516040840152602081015167ffffffffffffffff80821660608601528060408401511660808601528060608401511660a08601528060808401511660c086015250505060208401516101408060e0850152615e9f61018085018361488a565b915060408601517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08086850301610100870152615edc848361488a565b935060608801519150615efb6101208701836001600160a01b03169052565b60808801518387015260a0880151925080868503016101608701525050615f228282615d18565b9150508281036020840152614ccf8185615dea56fea164736f6c6343000818000a",
}

var OffRampABI = OffRampMetaData.ABI

var OffRampBin = OffRampMetaData.Bin

func DeployOffRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OffRampStaticConfig, dynamicConfig OffRampDynamicConfig, sourceChainConfigs []OffRampSourceChainConfigArgs) (common.Address, *types.Transaction, *OffRamp, error) {
	parsed, err := OffRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OffRampBin), backend, staticConfig, dynamicConfig, sourceChainConfigs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffRamp{address: address, abi: *parsed, OffRampCaller: OffRampCaller{contract: contract}, OffRampTransactor: OffRampTransactor{contract: contract}, OffRampFilterer: OffRampFilterer{contract: contract}}, nil
}

type OffRamp struct {
	address common.Address
	abi     abi.ABI
	OffRampCaller
	OffRampTransactor
	OffRampFilterer
}

type OffRampCaller struct {
	contract *bind.BoundContract
}

type OffRampTransactor struct {
	contract *bind.BoundContract
}

type OffRampFilterer struct {
	contract *bind.BoundContract
}

type OffRampSession struct {
	Contract     *OffRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OffRampCallerSession struct {
	Contract *OffRampCaller
	CallOpts bind.CallOpts
}

type OffRampTransactorSession struct {
	Contract     *OffRampTransactor
	TransactOpts bind.TransactOpts
}

type OffRampRaw struct {
	Contract *OffRamp
}

type OffRampCallerRaw struct {
	Contract *OffRampCaller
}

type OffRampTransactorRaw struct {
	Contract *OffRampTransactor
}

func NewOffRamp(address common.Address, backend bind.ContractBackend) (*OffRamp, error) {
	abi, err := abi.JSON(strings.NewReader(OffRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOffRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffRamp{address: address, abi: abi, OffRampCaller: OffRampCaller{contract: contract}, OffRampTransactor: OffRampTransactor{contract: contract}, OffRampFilterer: OffRampFilterer{contract: contract}}, nil
}

func NewOffRampCaller(address common.Address, caller bind.ContractCaller) (*OffRampCaller, error) {
	contract, err := bindOffRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampCaller{contract: contract}, nil
}

func NewOffRampTransactor(address common.Address, transactor bind.ContractTransactor) (*OffRampTransactor, error) {
	contract, err := bindOffRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampTransactor{contract: contract}, nil
}

func NewOffRampFilterer(address common.Address, filterer bind.ContractFilterer) (*OffRampFilterer, error) {
	contract, err := bindOffRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffRampFilterer{contract: contract}, nil
}

func bindOffRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OffRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OffRamp *OffRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRamp.Contract.OffRampCaller.contract.Call(opts, result, method, params...)
}

func (_OffRamp *OffRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRamp.Contract.OffRampTransactor.contract.Transfer(opts)
}

func (_OffRamp *OffRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRamp.Contract.OffRampTransactor.contract.Transact(opts, method, params...)
}

func (_OffRamp *OffRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_OffRamp *OffRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRamp.Contract.contract.Transfer(opts)
}

func (_OffRamp *OffRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRamp.Contract.contract.Transact(opts, method, params...)
}

func (_OffRamp *OffRampCaller) CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "ccipReceive", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_OffRamp *OffRampSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRamp.Contract.CcipReceive(&_OffRamp.CallOpts, arg0)
}

func (_OffRamp *OffRampCallerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRamp.Contract.CcipReceive(&_OffRamp.CallOpts, arg0)
}

func (_OffRamp *OffRampCaller) GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OffRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampDynamicConfig)).(*OffRampDynamicConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRamp.Contract.GetDynamicConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRamp.Contract.GetDynamicConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getExecutionState", sourceChainSelector, sequenceNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_OffRamp *OffRampSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRamp.Contract.GetExecutionState(&_OffRamp.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRamp *OffRampCallerSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRamp.Contract.GetExecutionState(&_OffRamp.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRamp *OffRampCaller) GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getLatestPriceSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OffRamp *OffRampSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRamp.Contract.GetLatestPriceSequenceNumber(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRamp.Contract.GetLatestPriceSequenceNumber(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getMerkleRoot", sourceChainSelector, root)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OffRamp *OffRampSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRamp.Contract.GetMerkleRoot(&_OffRamp.CallOpts, sourceChainSelector, root)
}

func (_OffRamp *OffRampCallerSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRamp.Contract.GetMerkleRoot(&_OffRamp.CallOpts, sourceChainSelector, root)
}

func (_OffRamp *OffRampCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampSourceChainConfig)).(*OffRampSourceChainConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetSourceChainConfig(&_OffRamp.CallOpts, sourceChainSelector)
}

func (_OffRamp *OffRampCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetSourceChainConfig(&_OffRamp.CallOpts, sourceChainSelector)
}

func (_OffRamp *OffRampCaller) GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OffRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampStaticConfig)).(*OffRampStaticConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRamp.Contract.GetStaticConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRamp.Contract.GetStaticConfig(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) IsBlessed(opts *bind.CallOpts, root [32]byte) (bool, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "isBlessed", root)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OffRamp *OffRampSession) IsBlessed(root [32]byte) (bool, error) {
	return _OffRamp.Contract.IsBlessed(&_OffRamp.CallOpts, root)
}

func (_OffRamp *OffRampCallerSession) IsBlessed(root [32]byte) (bool, error) {
	return _OffRamp.Contract.IsBlessed(&_OffRamp.CallOpts, root)
}

func (_OffRamp *OffRampCaller) LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "latestConfigDetails", ocrPluginType)

	if err != nil {
		return *new(MultiOCR3BaseOCRConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOCRConfig)).(*MultiOCR3BaseOCRConfig)

	return out0, err

}

func (_OffRamp *OffRampSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRamp.Contract.LatestConfigDetails(&_OffRamp.CallOpts, ocrPluginType)
}

func (_OffRamp *OffRampCallerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRamp.Contract.LatestConfigDetails(&_OffRamp.CallOpts, ocrPluginType)
}

func (_OffRamp *OffRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OffRamp *OffRampSession) Owner() (common.Address, error) {
	return _OffRamp.Contract.Owner(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) Owner() (common.Address, error) {
	return _OffRamp.Contract.Owner(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OffRamp *OffRampSession) TypeAndVersion() (string, error) {
	return _OffRamp.Contract.TypeAndVersion(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) TypeAndVersion() (string, error) {
	return _OffRamp.Contract.TypeAndVersion(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "acceptOwnership")
}

func (_OffRamp *OffRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRamp.Contract.AcceptOwnership(&_OffRamp.TransactOpts)
}

func (_OffRamp *OffRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRamp.Contract.AcceptOwnership(&_OffRamp.TransactOpts)
}

func (_OffRamp *OffRampTransactor) ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "applySourceChainConfigUpdates", sourceChainConfigUpdates)
}

func (_OffRamp *OffRampSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.ApplySourceChainConfigUpdates(&_OffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRamp *OffRampTransactorSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.ApplySourceChainConfigUpdates(&_OffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRamp *OffRampTransactor) Commit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "commit", reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampSession) Commit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactorSession) Commit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactor) Execute(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "execute", reportContext, report)
}

func (_OffRamp *OffRampSession) Execute(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Execute(&_OffRamp.TransactOpts, reportContext, report)
}

func (_OffRamp *OffRampTransactorSession) Execute(reportContext [3][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Execute(&_OffRamp.TransactOpts, reportContext, report)
}

func (_OffRamp *OffRampTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData)
}

func (_OffRamp *OffRampSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData)
}

func (_OffRamp *OffRampTransactorSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData)
}

func (_OffRamp *OffRampTransactor) ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]*big.Int) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "manuallyExecute", reports, gasLimitOverrides)
}

func (_OffRamp *OffRampSession) ManuallyExecute(reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]*big.Int) (*types.Transaction, error) {
	return _OffRamp.Contract.ManuallyExecute(&_OffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRamp *OffRampTransactorSession) ManuallyExecute(reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]*big.Int) (*types.Transaction, error) {
	return _OffRamp.Contract.ManuallyExecute(&_OffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRamp *OffRampTransactor) ResetUnblessedRoots(opts *bind.TransactOpts, rootToReset []OffRampUnblessedRoot) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "resetUnblessedRoots", rootToReset)
}

func (_OffRamp *OffRampSession) ResetUnblessedRoots(rootToReset []OffRampUnblessedRoot) (*types.Transaction, error) {
	return _OffRamp.Contract.ResetUnblessedRoots(&_OffRamp.TransactOpts, rootToReset)
}

func (_OffRamp *OffRampTransactorSession) ResetUnblessedRoots(rootToReset []OffRampUnblessedRoot) (*types.Transaction, error) {
	return _OffRamp.Contract.ResetUnblessedRoots(&_OffRamp.TransactOpts, rootToReset)
}

func (_OffRamp *OffRampTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OffRamp *OffRampSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRamp.Contract.SetDynamicConfig(&_OffRamp.TransactOpts, dynamicConfig)
}

func (_OffRamp *OffRampTransactorSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRamp.Contract.SetDynamicConfig(&_OffRamp.TransactOpts, dynamicConfig)
}

func (_OffRamp *OffRampTransactor) SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "setOCR3Configs", ocrConfigArgs)
}

func (_OffRamp *OffRampSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.SetOCR3Configs(&_OffRamp.TransactOpts, ocrConfigArgs)
}

func (_OffRamp *OffRampTransactorSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRamp.Contract.SetOCR3Configs(&_OffRamp.TransactOpts, ocrConfigArgs)
}

func (_OffRamp *OffRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_OffRamp *OffRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRamp.Contract.TransferOwnership(&_OffRamp.TransactOpts, to)
}

func (_OffRamp *OffRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRamp.Contract.TransferOwnership(&_OffRamp.TransactOpts, to)
}

type OffRampAlreadyAttemptedIterator struct {
	Event *OffRampAlreadyAttempted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampAlreadyAttemptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampAlreadyAttempted)
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
		it.Event = new(OffRampAlreadyAttempted)
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

func (it *OffRampAlreadyAttemptedIterator) Error() error {
	return it.fail
}

func (it *OffRampAlreadyAttemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampAlreadyAttempted struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampAlreadyAttemptedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return &OffRampAlreadyAttemptedIterator{contract: _OffRamp.contract, event: "AlreadyAttempted", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampAlreadyAttempted) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampAlreadyAttempted)
				if err := _OffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseAlreadyAttempted(log types.Log) (*OffRampAlreadyAttempted, error) {
	event := new(OffRampAlreadyAttempted)
	if err := _OffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampCommitReportAcceptedIterator struct {
	Event *OffRampCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampCommitReportAccepted)
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
		it.Event = new(OffRampCommitReportAccepted)
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

func (it *OffRampCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *OffRampCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampCommitReportAccepted struct {
	Report OffRampCommitReport
	Raw    types.Log
}

func (_OffRamp *OffRampFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampCommitReportAcceptedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &OffRampCommitReportAcceptedIterator{contract: _OffRamp.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampCommitReportAccepted)
				if err := _OffRamp.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseCommitReportAccepted(log types.Log) (*OffRampCommitReportAccepted, error) {
	event := new(OffRampCommitReportAccepted)
	if err := _OffRamp.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampConfigSetIterator struct {
	Event *OffRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampConfigSet)
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
		it.Event = new(OffRampConfigSet)
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

func (it *OffRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampConfigSet struct {
	OcrPluginType uint8
	ConfigDigest  [32]byte
	Signers       []common.Address
	Transmitters  []common.Address
	F             uint8
	Raw           types.Log
}

func (_OffRamp *OffRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OffRampConfigSetIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampConfigSetIterator{contract: _OffRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseConfigSet(log types.Log) (*OffRampConfigSet, error) {
	event := new(OffRampConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampDynamicConfigSetIterator struct {
	Event *OffRampDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampDynamicConfigSet)
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
		it.Event = new(OffRampDynamicConfigSet)
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

func (it *OffRampDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampDynamicConfigSet struct {
	DynamicConfig OffRampDynamicConfig
	Raw           types.Log
}

func (_OffRamp *OffRampFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampDynamicConfigSetIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampDynamicConfigSetIterator{contract: _OffRamp.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampDynamicConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampDynamicConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseDynamicConfigSet(log types.Log) (*OffRampDynamicConfigSet, error) {
	event := new(OffRampDynamicConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampExecutionStateChangedIterator struct {
	Event *OffRampExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampExecutionStateChanged)
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
		it.Event = new(OffRampExecutionStateChanged)
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

func (it *OffRampExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *OffRampExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	State               uint8
	ReturnData          []byte
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampExecutionStateChangedIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &OffRampExecutionStateChangedIterator{contract: _OffRamp.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}
	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampExecutionStateChanged)
				if err := _OffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseExecutionStateChanged(log types.Log) (*OffRampExecutionStateChanged, error) {
	event := new(OffRampExecutionStateChanged)
	if err := _OffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampOwnershipTransferRequestedIterator struct {
	Event *OffRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampOwnershipTransferRequested)
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
		it.Event = new(OffRampOwnershipTransferRequested)
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

func (it *OffRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OffRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRamp *OffRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampOwnershipTransferRequestedIterator{contract: _OffRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampOwnershipTransferRequested)
				if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffRampOwnershipTransferRequested, error) {
	event := new(OffRampOwnershipTransferRequested)
	if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampOwnershipTransferredIterator struct {
	Event *OffRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampOwnershipTransferred)
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
		it.Event = new(OffRampOwnershipTransferred)
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

func (it *OffRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OffRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRamp *OffRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampOwnershipTransferredIterator{contract: _OffRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampOwnershipTransferred)
				if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseOwnershipTransferred(log types.Log) (*OffRampOwnershipTransferred, error) {
	event := new(OffRampOwnershipTransferred)
	if err := _OffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampRootRemovedIterator struct {
	Event *OffRampRootRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampRootRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampRootRemoved)
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
		it.Event = new(OffRampRootRemoved)
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

func (it *OffRampRootRemovedIterator) Error() error {
	return it.fail
}

func (it *OffRampRootRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampRootRemoved struct {
	Root [32]byte
	Raw  types.Log
}

func (_OffRamp *OffRampFilterer) FilterRootRemoved(opts *bind.FilterOpts) (*OffRampRootRemovedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return &OffRampRootRemovedIterator{contract: _OffRamp.contract, event: "RootRemoved", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampRootRemoved) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampRootRemoved)
				if err := _OffRamp.contract.UnpackLog(event, "RootRemoved", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseRootRemoved(log types.Log) (*OffRampRootRemoved, error) {
	event := new(OffRampRootRemoved)
	if err := _OffRamp.contract.UnpackLog(event, "RootRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampSkippedAlreadyExecutedMessageIterator struct {
	Event *OffRampSkippedAlreadyExecutedMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSkippedAlreadyExecutedMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSkippedAlreadyExecutedMessage)
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
		it.Event = new(OffRampSkippedAlreadyExecutedMessage)
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

func (it *OffRampSkippedAlreadyExecutedMessageIterator) Error() error {
	return it.fail
}

func (it *OffRampSkippedAlreadyExecutedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSkippedAlreadyExecutedMessage struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampSkippedAlreadyExecutedMessageIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return &OffRampSkippedAlreadyExecutedMessageIterator{contract: _OffRamp.contract, event: "SkippedAlreadyExecutedMessage", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampSkippedAlreadyExecutedMessage) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSkippedAlreadyExecutedMessage)
				if err := _OffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampSkippedAlreadyExecutedMessage, error) {
	event := new(OffRampSkippedAlreadyExecutedMessage)
	if err := _OffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampSourceChainConfigSetIterator struct {
	Event *OffRampSourceChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSourceChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSourceChainConfigSet)
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
		it.Event = new(OffRampSourceChainConfigSet)
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

func (it *OffRampSourceChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampSourceChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceConfig        OffRampSourceChainConfig
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampSourceChainConfigSetIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OffRampSourceChainConfigSetIterator{contract: _OffRamp.contract, event: "SourceChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSourceChainConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSourceChainConfigSet(log types.Log) (*OffRampSourceChainConfigSet, error) {
	event := new(OffRampSourceChainConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampSourceChainSelectorAddedIterator struct {
	Event *OffRampSourceChainSelectorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSourceChainSelectorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSourceChainSelectorAdded)
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
		it.Event = new(OffRampSourceChainSelectorAdded)
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

func (it *OffRampSourceChainSelectorAddedIterator) Error() error {
	return it.fail
}

func (it *OffRampSourceChainSelectorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSourceChainSelectorAdded struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampSourceChainSelectorAddedIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return &OffRampSourceChainSelectorAddedIterator{contract: _OffRamp.contract, event: "SourceChainSelectorAdded", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainSelectorAdded) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSourceChainSelectorAdded)
				if err := _OffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSourceChainSelectorAdded(log types.Log) (*OffRampSourceChainSelectorAdded, error) {
	event := new(OffRampSourceChainSelectorAdded)
	if err := _OffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampStaticConfigSetIterator struct {
	Event *OffRampStaticConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampStaticConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampStaticConfigSet)
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
		it.Event = new(OffRampStaticConfigSet)
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

func (it *OffRampStaticConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampStaticConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampStaticConfigSet struct {
	StaticConfig OffRampStaticConfig
	Raw          types.Log
}

func (_OffRamp *OffRampFilterer) FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampStaticConfigSetIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampStaticConfigSetIterator{contract: _OffRamp.contract, event: "StaticConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampStaticConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampStaticConfigSet)
				if err := _OffRamp.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseStaticConfigSet(log types.Log) (*OffRampStaticConfigSet, error) {
	event := new(OffRampStaticConfigSet)
	if err := _OffRamp.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampTransmittedIterator struct {
	Event *OffRampTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampTransmitted)
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
		it.Event = new(OffRampTransmitted)
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

func (it *OffRampTransmittedIterator) Error() error {
	return it.fail
}

func (it *OffRampTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampTransmitted struct {
	OcrPluginType  uint8
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_OffRamp *OffRampFilterer) FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampTransmittedIterator, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return &OffRampTransmittedIterator{contract: _OffRamp.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampTransmitted, ocrPluginType []uint8) (event.Subscription, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampTransmitted)
				if err := _OffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseTransmitted(log types.Log) (*OffRampTransmitted, error) {
	event := new(OffRampTransmitted)
	if err := _OffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OffRamp *OffRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OffRamp.abi.Events["AlreadyAttempted"].ID:
		return _OffRamp.ParseAlreadyAttempted(log)
	case _OffRamp.abi.Events["CommitReportAccepted"].ID:
		return _OffRamp.ParseCommitReportAccepted(log)
	case _OffRamp.abi.Events["ConfigSet"].ID:
		return _OffRamp.ParseConfigSet(log)
	case _OffRamp.abi.Events["DynamicConfigSet"].ID:
		return _OffRamp.ParseDynamicConfigSet(log)
	case _OffRamp.abi.Events["ExecutionStateChanged"].ID:
		return _OffRamp.ParseExecutionStateChanged(log)
	case _OffRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _OffRamp.ParseOwnershipTransferRequested(log)
	case _OffRamp.abi.Events["OwnershipTransferred"].ID:
		return _OffRamp.ParseOwnershipTransferred(log)
	case _OffRamp.abi.Events["RootRemoved"].ID:
		return _OffRamp.ParseRootRemoved(log)
	case _OffRamp.abi.Events["SkippedAlreadyExecutedMessage"].ID:
		return _OffRamp.ParseSkippedAlreadyExecutedMessage(log)
	case _OffRamp.abi.Events["SourceChainConfigSet"].ID:
		return _OffRamp.ParseSourceChainConfigSet(log)
	case _OffRamp.abi.Events["SourceChainSelectorAdded"].ID:
		return _OffRamp.ParseSourceChainSelectorAdded(log)
	case _OffRamp.abi.Events["StaticConfigSet"].ID:
		return _OffRamp.ParseStaticConfigSet(log)
	case _OffRamp.abi.Events["Transmitted"].ID:
		return _OffRamp.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OffRampAlreadyAttempted) Topic() common.Hash {
	return common.HexToHash("0x3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120")
}

func (OffRampCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x3a3950e13dd607cc37980db0ef14266c40d2bba9c01b2e44bfe549808883095d")
}

func (OffRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (OffRampDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0xa55bd56595c45f517e5967a3067f3dca684445a3080e7c04a4e0d5a40cda627d")
}

<<<<<<<< HEAD:core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp/evm_2_evm_multi_offramp.go
func (EVM2EVMMultiOffRampExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x8c324ce1367b83031769f6a813e3bb4c117aba2185789d66b98b791405be6df2")
========
func (OffRampExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0xdc8ccbc35e0eebd81239bcd1971fcd53c7eb34034880142a0f43c809a458732f")
>>>>>>>> 377f0dbef9 (Rename ramps and rmn (#1323)):core/gethwrappers/ccip/generated/offramp/offramp.go
}

func (OffRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OffRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OffRampRootRemoved) Topic() common.Hash {
	return common.HexToHash("0x202f1139a3e334b6056064c0e9b19fd07e44a88d8f6e5ded571b24cf8c371f12")
}

func (OffRampSkippedAlreadyExecutedMessage) Topic() common.Hash {
	return common.HexToHash("0x3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c")
}

func (OffRampSourceChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b")
}

func (OffRampSourceChainSelectorAdded) Topic() common.Hash {
	return common.HexToHash("0xf4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9")
}

func (OffRampStaticConfigSet) Topic() common.Hash {
	return common.HexToHash("0x683eb52ee924eb817377cfa8f41f238f4bb7a877da5267869dfffbad85f564d8")
}

func (OffRampTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_OffRamp *OffRamp) Address() common.Address {
	return _OffRamp.address
}

type OffRampInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error)

	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error)

	GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error)

	IsBlessed(opts *bind.CallOpts, root [32]byte) (bool, error)

	LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error)

	Commit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Execute(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error)

	ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte) (*types.Transaction, error)

	ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]*big.Int) (*types.Transaction, error)

	ResetUnblessedRoots(opts *bind.TransactOpts, rootToReset []OffRampUnblessedRoot) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error)

	SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampAlreadyAttemptedIterator, error)

	WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampAlreadyAttempted) (event.Subscription, error)

	ParseAlreadyAttempted(log types.Log) (*OffRampAlreadyAttempted, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*OffRampCommitReportAccepted, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OffRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OffRampConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampDynamicConfigSet) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*OffRampDynamicConfigSet, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*OffRampExecutionStateChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OffRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OffRampOwnershipTransferred, error)

	FilterRootRemoved(opts *bind.FilterOpts) (*OffRampRootRemovedIterator, error)

	WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampRootRemoved) (event.Subscription, error)

	ParseRootRemoved(log types.Log) (*OffRampRootRemoved, error)

	FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampSkippedAlreadyExecutedMessageIterator, error)

	WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampSkippedAlreadyExecutedMessage) (event.Subscription, error)

	ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampSkippedAlreadyExecutedMessage, error)

	FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampSourceChainConfigSetIterator, error)

	WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error)

	ParseSourceChainConfigSet(log types.Log) (*OffRampSourceChainConfigSet, error)

	FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampSourceChainSelectorAddedIterator, error)

	WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampSourceChainSelectorAdded) (event.Subscription, error)

	ParseSourceChainSelectorAdded(log types.Log) (*OffRampSourceChainSelectorAdded, error)

	FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampStaticConfigSetIterator, error)

	WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampStaticConfigSet) (event.Subscription, error)

	ParseStaticConfigSet(log types.Log) (*OffRampStaticConfigSet, error)

	FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampTransmitted, ocrPluginType []uint8) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*OffRampTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
