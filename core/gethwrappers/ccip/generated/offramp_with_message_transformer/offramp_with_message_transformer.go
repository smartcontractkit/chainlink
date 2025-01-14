// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package offramp_with_message_transformer

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
	TokenAmounts []InternalAny2EVMTokenTransfer
}

type InternalAny2EVMTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  common.Address
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

type InternalExecutionReport struct {
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

type InternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
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

type OffRampDynamicConfig struct {
	FeeQuoter                               common.Address
	PermissionLessExecutionThresholdSeconds uint32
	IsRMNVerificationDisabled               bool
	MessageInterceptor                      common.Address
}

type OffRampGasLimitOverride struct {
	ReceiverExecutionGasLimit *big.Int
	TokenGasOverrides         []uint32
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
	ChainSelector        uint64
	GasForCallExactCheck uint16
	RmnRemote            common.Address
	TokenAdminRegistry   common.Address
	NonceManager         common.Address
}

var OffRampWithMessageTransformerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applySourceChainConfigUpdates\",\"inputs\":[{\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"commit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeSingleMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllSourceChainConfigs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExecutionState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestPriceSequenceNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMerkleRoot\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageTransformer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSourceChainConfig\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"ocrConfig\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"components\":[{\"name\":\"configInfo\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"n\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyExecute\",\"inputs\":[{\"name\":\"reports\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\",\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"components\":[{\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMessageTransformer\",\"inputs\":[{\"name\":\"messageTransformerAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOCR3Configs\",\"inputs\":[{\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AlreadyAttempted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitReportAccepted\",\"inputs\":[{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"F\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DynamicConfigSet\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecutionStateChanged\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"state\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RootRemoved\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedAlreadyExecutedMessage\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedReportExecution\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainConfigSet\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sourceConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainSelectorAdded\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaticConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CanOnlySelfCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CommitOnRampMismatch\",\"inputs\":[{\"name\":\"reportOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"configOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyReport\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ExecutionError\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ForkedChain\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientGasForCallWithExact\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[{\"name\":\"errorType\",\"type\":\"uint8\",\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\"}]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInterval\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionGasLimit\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"newLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionTokenGasOverride\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oldLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverride\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidMessageDestChainSelector\",\"inputs\":[{\"name\":\"messageDestChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidNewState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newState\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}]},{\"type\":\"error\",\"name\":\"InvalidOnRampUpdate\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LeavesCannotBeEmpty\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionGasAmountCountMismatch\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ManualExecutionGasLimitMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionNotYetEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageTransformError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MessageValidationError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonUniqueSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotACompatiblePool\",\"inputs\":[{\"name\":\"notPool\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OracleCannotBeZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReceiverError\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ReleaseOrMintBalanceMismatch\",\"inputs\":[{\"name\":\"amountReleased\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePre\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RootAlreadyCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"RootNotCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignaturesOutOfRegistration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SourceChainNotEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SourceChainSelectorMismatch\",\"inputs\":[{\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StaleCommitReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaticConfigCannotBeChanged\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TokenDataMismatch\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TokenHandlingError\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnexpectedTokenData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongMessageLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WrongNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroChainSelectorNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461088f57616840803803809161001d82856108c5565b8339810190808203610160811261088f5760a0811261088f5760405160a081016001600160401b038111828210176108945760405261005b836108e8565b815260208301519061ffff8216820361088f57602081019182526040840151936001600160a01b038516850361088f576040820194855261009e606082016108fc565b946060830195865260806100b38184016108fc565b84820190815295609f19011261088f57604051936100d0856108aa565b6100dc60a084016108fc565b855260c08301519363ffffffff8516850361088f576020860194855261010460e08501610910565b966040870197885261011961010086016108fc565b606088019081526101208601519095906001600160401b03811161088f5781018b601f8201121561088f5780519b6001600160401b038d11610894578c60051b91604051809e6020850161016d90836108c5565b81526020019281016020019082821161088f5760208101935b82851061078f57505050505061014061019f91016108fc565b98331561077e57600180546001600160a01b031916331790554660805284516001600160a01b031615801561076c575b801561075a575b6107385782516001600160401b0316156107495782516001600160401b0390811660a090815286516001600160a01b0390811660c0528351811660e0528451811661010052865161ffff90811661012052604080519751909416875296519096166020860152955185169084015251831660608301525190911660808201527fb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b79190a182516001600160a01b031615610738579151600480548351865160ff60c01b90151560c01b1663ffffffff60a01b60a09290921b919091166001600160a01b039485166001600160c81b0319909316831717179091558351600580549184166001600160a01b031990921691909117905560408051918252925163ffffffff166020820152935115159184019190915290511660608201529091907fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee90608090a16000915b81518310156106805760009260208160051b8401015160018060401b036020820151169081156106715780516001600160a01b031615610662578186526008602052604086206060820151916001820192610399845461091d565b610603578254600160a81b600160e81b031916600160a81b1783556040518581527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb990602090a15b805180159081156105d8575b506105c9578051906001600160401b0382116105b55761040d855461091d565b601f8111610570575b50602090601f83116001146104f8579180916000805160206168208339815191529695949360019a9b9c926104ed575b5050600019600383901b1c191690881b1783555b60408101518254915160a089811b8a9003801960ff60a01b1990951693151590911b60ff60a01b1692909217929092169116178155610498846109da565b506104e26040519283926020845254888060a01b038116602085015260ff8160a01c1615156040850152888060401b039060a81c16606084015260808084015260a0830190610957565b0390a201919061033e565b015190503880610446565b858b52818b20919a601f198416905b8181106105585750916001999a9b8492600080516020616820833981519152989796958c951061053f575b505050811b01835561045a565b015160001960f88460031b161c19169055388080610532565b828d0151845560209c8d019c60019094019301610507565b858b5260208b20601f840160051c810191602085106105ab575b601f0160051c01905b8181106105a05750610416565b8b8155600101610593565b909150819061058a565b634e487b7160e01b8a52604160045260248afd5b6342bcdf7f60e11b8952600489fd5b9050602082012060405160208101908b8252602081526105f96040826108c5565b51902014386103ed565b825460a81c6001600160401b03166001141580610634575b156103e157632105803760e11b89526004859052602489fd5b5060405161064d816106468188610957565b03826108c5565b6020815191012081516020830120141561061b565b6342bcdf7f60e11b8652600486fd5b63c656089560e01b8652600486fd5b6001600160a01b0381161561073857600b8054600160401b600160e01b031916604092831b600160401b600160e01b031617905551615db29081610a6e823960805181613789015260a05181818161049f0152614757015260c0518181816104f501528181612db70152818161320b01526146f1015260e0518181816105240152614f340152610100518181816105530152614b1a0152610120518181816104c601528181612529015281816150270152615ae70152f35b6342bcdf7f60e11b60005260046000fd5b63c656089560e01b60005260046000fd5b5081516001600160a01b0316156101d6565b5080516001600160a01b0316156101cf565b639b15e16f60e01b60005260046000fd5b84516001600160401b03811161088f5782016080818603601f19011261088f57604051906107bc826108aa565b60208101516001600160a01b038116810361088f5782526107df604082016108e8565b60208301526107f060608201610910565b604083015260808101516001600160401b03811161088f57602091010185601f8201121561088f5780516001600160401b0381116108945760405191610840601f8301601f1916602001846108c5565b818352876020838301011161088f5760005b82811061087a5750509181600060208096949581960101526060820152815201940193610186565b80602080928401015182828701015201610852565b600080fd5b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761089457604052565b601f909101601f19168101906001600160401b0382119082101761089457604052565b51906001600160401b038216820361088f57565b51906001600160a01b038216820361088f57565b5190811515820361088f57565b90600182811c9216801561094d575b602083101461093757565b634e487b7160e01b600052602260045260246000fd5b91607f169161092c565b600092918154916109678361091d565b80835292600181169081156109bd575060011461098357505050565b60009081526020812093945091925b8383106109a3575060209250010190565b600181602092949394548385870101520191019190610992565b915050602093945060ff929192191683830152151560051b010190565b80600052600760205260406000205415600014610a675760065468010000000000000000811015610894576001810180600655811015610a51577ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0181905560065460009182526007602052604090912055600190565b634e487b7160e01b600052603260045260246000fd5b5060009056fe6080604052600436101561001257600080fd5b60003560e01c806304666f9c1461017757806306285c691461017257806315777ab21461016d578063181f5a77146101685780633f4b04aa146101635780635215505b1461015e5780635e36480c146101595780635e7bb0081461015457806360987c201461014f57806365b81aab1461014a5780637437ff9f1461014557806379ba5097146101405780637edf52f41461013b57806385572ffb146101365780638da5cb5b14610131578063c673e5841461012c578063ccd37ba314610127578063de5e0b9a14610122578063e9d68a8e1461011d578063f2fde38b14610118578063f58e03fc146101135763f716f99f1461010e57600080fd5b611990565b611873565b6117e8565b611749565b6116ad565b611625565b61157a565b611492565b61145c565b611296565b611216565b61116d565b6110f7565b61107c565b610e77565b610909565b6107c4565b6106b7565b610658565b6105d1565b61046c565b61034c565b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b038211176101ad57604052565b61017c565b60a081019081106001600160401b038211176101ad57604052565b604081019081106001600160401b038211176101ad57604052565b606081019081106001600160401b038211176101ad57604052565b60c081019081106001600160401b038211176101ad57604052565b90601f801991011681019081106001600160401b038211176101ad57604052565b6040519061024e60c08361021e565b565b6040519061024e60a08361021e565b6040519061024e6101008361021e565b6040519061024e60408361021e565b6001600160401b0381116101ad5760051b60200190565b6001600160a01b038116036102a657565b600080fd5b6001600160401b038116036102a657565b359061024e826102ab565b801515036102a657565b359061024e826102c7565b6001600160401b0381116101ad57601f01601f191660200190565b929192610303826102dc565b91610311604051938461021e565b8294818452818301116102a6578281602093846000960137010152565b9080601f830112156102a657816020610349933591016102f7565b90565b346102a65760203660031901126102a6576004356001600160401b0381116102a657366023820112156102a6578060040135906103888261027e565b90610396604051928361021e565b8282526024602083019360051b820101903682116102a65760248101935b8285106103c6576103c484611acb565b005b84356001600160401b0381116102a6578201608060231982360301126102a657604051916103f383610192565b602482013561040181610295565b83526044820135610411816102ab565b60208401526064820135610424816102c7565b60408401526084820135926001600160401b0384116102a65761045160209493602486953692010161032e565b60608201528152019401936103b4565b60009103126102a657565b346102a65760003660031901126102a657610485611d77565b506105cd604051610495816101b2565b6001600160401b037f000000000000000000000000000000000000000000000000000000000000000016815261ffff7f00000000000000000000000000000000000000000000000000000000000000001660208201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660408201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660608201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660808201526040519182918291909160806001600160a01b038160a08401956001600160401b03815116855261ffff6020820151166020860152826040820151166040860152826060820151166060860152015116910152565b0390f35b346102a65760003660031901126102a65760206001600160a01b03600b5460401c16604051908152f35b6040519061060a60208361021e565b60008252565b60005b8381106106235750506000910152565b8181015183820152602001610613565b9060209161064c81518092818552858086019101610610565b601f01601f1916010190565b346102a65760003660031901126102a6576105cd604080519061067b818361021e565b601182527f4f666652616d7020312e362e302d646576000000000000000000000000000000602083015251918291602083526020830190610633565b346102a65760003660031901126102a65760206001600160401b03600b5416604051908152f35b9060806060610349936001600160a01b0381511684526020810151151560208501526001600160401b0360408201511660408501520151918160608201520190610633565b6040810160408252825180915260206060830193019060005b8181106107a5575050506020818303910152815180825260208201916020808360051b8301019401926000915b83831061077857505050505090565b9091929394602080610796600193601f1986820301875289516106de565b97019301930191939290610769565b82516001600160401b031685526020948501949092019160010161073c565b346102a65760003660031901126102a6576006546107e18161027e565b906107ef604051928361021e565b808252601f196107fe8261027e565b0160005b8181106108c057505061081481611dc9565b9060005b8181106108305750506105cd60405192839283610723565b8061086661084e6108426001946145d8565b6001600160401b031690565b6108588387611e23565b906001600160401b03169052565b6108a461089f6108866108798488611e23565b516001600160401b031690565b6001600160401b03166000526008602052604060002090565b611f0f565b6108ae8287611e23565b526108b98186611e23565b5001610818565b6020906108cb611da2565b82828701015201610802565b634e487b7160e01b600052602160045260246000fd5b600411156108f757565b6108d7565b9060048210156108f75752565b346102a65760403660031901126102a657602061093d60043561092b816102ab565b60243590610938826102ab565b611fa7565b61094a60405180926108fc565bf35b91908260a09103126102a657604051610964816101b2565b608080829480358452602081013561097b816102ab565b6020850152604081013561098e816102ab565b604085015260608101356109a1816102ab565b60608501520135916109b2836102ab565b0152565b359061024e82610295565b63ffffffff8116036102a657565b359061024e826109c1565b81601f820112156102a6578035906109f18261027e565b926109ff604051948561021e565b82845260208085019360051b830101918183116102a65760208101935b838510610a2b57505050505090565b84356001600160401b0381116102a657820160a0818503601f1901126102a65760405191610a58836101b2565b60208201356001600160401b0381116102a657856020610a7a9285010161032e565b83526040820135610a8a81610295565b6020840152610a9b606083016109cf565b60408401526080820135926001600160401b0384116102a65760a083610ac888602080988198010161032e565b606084015201356080820152815201940193610a1c565b919091610140818403126102a657610af561023f565b92610b00818361094c565b845260a08201356001600160401b0381116102a65781610b2191840161032e565b602085015260c08201356001600160401b0381116102a65781610b4591840161032e565b6040850152610b5660e083016109b6565b606085015261010082013560808501526101208201356001600160401b0381116102a657610b8492016109da565b60a0830152565b9080601f830112156102a6578135610ba28161027e565b92610bb0604051948561021e565b81845260208085019260051b820101918383116102a65760208201905b838210610bdc57505050505090565b81356001600160401b0381116102a657602091610bfe87848094880101610adf565b815201910190610bcd565b81601f820112156102a657803590610c208261027e565b92610c2e604051948561021e565b82845260208085019360051b830101918183116102a65760208101935b838510610c5a57505050505090565b84356001600160401b0381116102a657820183603f820112156102a6576020810135610c858161027e565b91610c93604051938461021e565b8183526020808085019360051b83010101918683116102a65760408201905b838210610ccc575050509082525060209485019401610c4b565b81356001600160401b0381116102a657602091610cf08a848080958901010161032e565b815201910190610cb2565b929190610d078161027e565b93610d15604051958661021e565b602085838152019160051b81019283116102a657905b828210610d3757505050565b8135815260209182019101610d2b565b9080601f830112156102a65781602061034993359101610cfb565b81601f820112156102a657803590610d798261027e565b92610d87604051948561021e565b82845260208085019360051b830101918183116102a65760208101935b838510610db357505050505090565b84356001600160401b0381116102a657820160a0818503601f1901126102a657610ddb610250565b91610de8602083016102bc565b835260408201356001600160401b0381116102a657856020610e0c92850101610b8b565b602084015260608201356001600160401b0381116102a657856020610e3392850101610c09565b60408401526080820135926001600160401b0384116102a65760a083610e60886020809881980101610d47565b606084015201356080820152815201940193610da4565b346102a65760403660031901126102a6576004356001600160401b0381116102a657610ea7903690600401610d62565b6024356001600160401b0381116102a657366023820112156102a657806004013591610ed28361027e565b91610ee0604051938461021e565b8383526024602084019460051b820101903682116102a65760248101945b828610610f0f576103c48585611fef565b85356001600160401b0381116102a6578201366043820112156102a6576024810135610f3a8161027e565b91610f48604051938461021e565b818352602060248185019360051b83010101903682116102a65760448101925b828410610f82575050509082525060209586019501610efe565b83356001600160401b0381116102a6576024908301016040601f1982360301126102a65760405190610fb3826101cd565b6020810135825260408101356001600160401b0381116102a657602091010136601f820112156102a657803590610fe98261027e565b91610ff7604051938461021e565b80835260208084019160051b830101913683116102a657602001905b8282106110325750505091816020938480940152815201930192610f68565b602080918335611041816109c1565b815201910190611013565b9181601f840112156102a6578235916001600160401b0383116102a6576020808501948460051b0101116102a657565b346102a65760603660031901126102a6576004356001600160401b0381116102a6576110ac903690600401610adf565b6024356001600160401b0381116102a6576110cb90369060040161104c565b91604435926001600160401b0384116102a6576110ef6103c494369060040161104c565b939092612402565b346102a65760203660031901126102a65760043561111481610295565b61111c613586565b7fffffffff0000000000000000000000000000000000000000ffffffffffffffff7bffffffffffffffffffffffffffffffffffffffff0000000000000000600b549260401b16911617600b55600080f35b346102a65760003660031901126102a6576111866126da565b506105cd60405161119681610192565b60ff6004546001600160a01b038116835263ffffffff8160a01c16602084015260c01c16151560408201526001600160a01b036005541660608201526040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b346102a65760003660031901126102a6576000546001600160a01b0381163303611285576001600160a01b0319600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b63015aa1e360e11b60005260046000fd5b346102a65760803660031901126102a65760006040516112b581610192565b6004356112c181610295565b81526024356112cf816109c1565b60208201526044356112e0816102c7565b60408201526064356112f181610295565b60608201526112fe613586565b6001600160a01b038151161561144d576114478161135d6001600160a01b037fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9451166001600160a01b03166001600160a01b03196004541617600455565b60208101516004547fffffffffffffff0000000000ffffffffffffffffffffffffffffffffffffffff77ffffffff000000000000000000000000000000000000000078ff0000000000000000000000000000000000000000000000006040860151151560c01b169360a01b16911617176004556114036113e760608301516001600160a01b031690565b6001600160a01b03166001600160a01b03196005541617600555565b6040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b0390a180f35b6342bcdf7f60e11b8252600482fd5b346102a65760203660031901126102a6576004356001600160401b0381116102a65760a09060031990360301126102a657600080fd5b346102a65760003660031901126102a65760206001600160a01b0360015416604051908152f35b6004359060ff821682036102a657565b359060ff821682036102a657565b906020808351928381520192019060005b8181106114f55750505090565b82516001600160a01b03168452602093840193909201916001016114e8565b906103499160208152606082518051602084015260ff602082015116604084015260ff604082015116828401520151151560808201526040611565602084015160c060a085015260e08401906114d7565b9201519060c0601f19828503019101526114d7565b346102a65760203660031901126102a65760ff6115956114b9565b6060604080516115a4816101e8565b6115ac6126da565b815282602082015201521660005260026020526105cd60406000206003611614604051926115d9846101e8565b6115e2816126ff565b84526040516115ff816115f88160028601612738565b038261021e565b60208501526115f86040518094819301612738565b604082015260405191829182611514565b346102a65760403660031901126102a657600435611642816102ab565b6001600160401b036024359116600052600a6020526040600020906000526020526020604060002054604051908152f35b906004916044116102a657565b9181601f840112156102a6578235916001600160401b0383116102a657602083818601950101116102a657565b346102a65760c03660031901126102a6576116c736611673565b6044356001600160401b0381116102a6576116e6903690600401611680565b6064929192356001600160401b0381116102a65761170890369060040161104c565b60843594916001600160401b0386116102a65761172c6103c496369060040161104c565b94909360a43596612d72565b9060206103499281815201906106de565b346102a65760203660031901126102a6576001600160401b0360043561176e816102ab565b611776611da2565b501660005260086020526105cd604060002060016117d76040519261179a84610192565b6001600160401b0381546001600160a01b038116865260ff8160a01c161515602087015260a81c1660408501526115f86040518094819301611e71565b606082015260405191829182611738565b346102a65760203660031901126102a6576001600160a01b0360043561180d81610295565b611815613586565b1633811461186257806001600160a01b031960005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b636d6c4ee560e11b60005260046000fd5b346102a65760603660031901126102a65761188d36611673565b6044356001600160401b0381116102a6576118ac903690600401611680565b918282016020838203126102a6578235906001600160401b0382116102a6576118d6918401610d62565b6040519060206118e6818461021e565b60008352601f19810160005b81811061191a575050506103c4949161190a916137ca565b611912613281565b928392613fa0565b606085820184015282016118f2565b9080601f830112156102a65781356119408161027e565b9261194e604051948561021e565b81845260208085019260051b8201019283116102a657602001905b8282106119765750505090565b60208091833561198581610295565b815201910190611969565b346102a65760203660031901126102a6576004356001600160401b0381116102a657366023820112156102a6578060040135906119cc8261027e565b906119da604051928361021e565b8282526024602083019360051b820101903682116102a65760248101935b828510611a08576103c48461329d565b84356001600160401b0381116102a657820160c060231982360301126102a657611a3061023f565b9160248201358352611a44604483016114c9565b6020840152611a55606483016114c9565b6040840152611a66608483016102d1565b606084015260a48201356001600160401b0381116102a657611a8e9060243691850101611929565b608084015260c4820135926001600160401b0384116102a657611abb602094936024869536920101611929565b60a08201528152019401936119f8565b611ad3613586565b60005b8151811015611d7357611ae98183611e23565b5190611aff60208301516001600160401b031690565b916001600160401b038316908115611d6257611b34611b28611b2883516001600160a01b031690565b6001600160a01b031690565b15611cc957611b56846001600160401b03166000526008602052604060002090565b906060810151916001810195611b6c8754611e37565b611cf057611bdf7ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb991611bc584750100000000000000000000000000000000000000000067ffffffffffffffff60a81b19825416179055565b6040516001600160401b0390911681529081906020820190565b0390a15b82518015908115611cda575b50611cc957611caa611c8e611cc093611c2b7f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b9660019a613628565b611c81611c3b6040830151151590565b85547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690151560a01b74ff000000000000000000000000000000000000000016178555565b516001600160a01b031690565b82906001600160a01b03166001600160a01b0319825416179055565b611cb384615600565b50604051918291826136f9565b0390a201611ad6565b6342bcdf7f60e11b60005260046000fd5b90506020840120611ce96135ab565b1438611bef565b60016001600160401b03611d0f84546001600160401b039060a81c1690565b16141580611d43575b611d225750611be3565b632105803760e11b6000526001600160401b031660045260246000fd5b6000fd5b50611d4d87611ef4565b60208151910120845160208601201415611d18565b63c656089560e01b60005260046000fd5b5050565b60405190611d84826101b2565b60006080838281528260208201528260408201528260608201520152565b60405190611daf82610192565b606080836000815260006020820152600060408201520152565b90611dd38261027e565b611de0604051918261021e565b8281528092611df1601f199161027e565b0190602036910137565b634e487b7160e01b600052603260045260246000fd5b805115611e1e5760200190565b611dfb565b8051821015611e1e5760209160051b010190565b90600182811c92168015611e67575b6020831014611e5157565b634e487b7160e01b600052602260045260246000fd5b91607f1691611e46565b60009291815491611e8183611e37565b8083529260018116908115611ed75750600114611e9d57505050565b60009081526020812093945091925b838310611ebd575060209250010190565b600181602092949394548385870101520191019190611eac565b915050602093945060ff929192191683830152151560051b010190565b9061024e611f089260405193848092611e71565b038361021e565b9060016060604051611f2081610192565b6109b281956001600160401b0381546001600160a01b038116855260ff8160a01c161515602086015260a81c166040840152611f626040518096819301611e71565b038461021e565b634e487b7160e01b600052601160045260246000fd5b908160051b9180830460201490151715611f9557565b611f69565b91908203918211611f9557565b611fb382607f92613743565b9116906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611f95576003911c1660048110156108f75790565b611ff7613787565b8051825181036121ee5760005b8181106120175750509061024e916137ca565b6120218184611e23565b5160208101908151516120348488611e23565b5192835182036121ee5790916000925b808410612058575050505050600101612004565b9194939861206a848b98939598611e23565b515198612078888851611e23565b5199806121a5575b5060a08a01988b60206120968b8d515193611e23565b51015151036121685760005b8a5151811015612153576120de6120d56120cb8f60206120c38f8793611e23565b510151611e23565b5163ffffffff1690565b63ffffffff1690565b8b816120ef575b50506001016120a2565b6120d560406121028561210e9451611e23565b51015163ffffffff1690565b9081811061211d57508b6120e5565b8d51516040516348e617b360e01b81526004810191909152602481019390935260448301919091526064820152608490fd5b0390fd5b50985098509893949095600101929091612044565b611d3f8b51612183606082519201516001600160401b031690565b6370a193fd60e01b6000526004919091526001600160401b0316602452604490565b60808b015181101561208057611d3f908b6121c788516001600160401b031690565b905151633a98d46360e11b6000526001600160401b03909116600452602452604452606490565b6320f8fd5960e21b60005260046000fd5b6040519061220c826101cd565b60006020838281520152565b6040519061222760208361021e565b600080835282815b82811061223b57505050565b6020906122466121ff565b8282850101520161222f565b805182526001600160401b0360208201511660208301526080612299612287604084015160a0604087015260a0860190610633565b60608401518582036060870152610633565b9101519160808183039101526020808351928381520192019060005b8181106122c25750505090565b825180516001600160a01b0316855260209081015181860152604090940193909201916001016122b5565b906020610349928181520190612252565b6040513d6000823e3d90fd5b3d15612335573d9061231b826102dc565b91612329604051938461021e565b82523d6000602084013e565b606090565b906020610349928181520190610633565b81601f820112156102a6578051612361816102dc565b9261236f604051948561021e565b818452602082840101116102a6576103499160208085019101610610565b90916060828403126102a65781516123a4816102c7565b9260208301516001600160401b0381116102a6576040916123c691850161234b565b92015190565b9293606092959461ffff6123f06001600160a01b0394608088526080880190612252565b97166020860152604085015216910152565b93919290933033036126c95761241790613c1a565b92612420612218565b9460a08501518051612682575b505050505080519161244b602084519401516001600160401b031690565b906020830151916040840192612478845192612465610250565b9788526001600160401b03166020880152565b6040860152606085015260808401526001600160a01b036124a16005546001600160a01b031690565b1680612605575b50515115806125f9575b80156125e3575b80156125ba575b611d735761255291816124f7611b286124ea610886602060009751016001600160401b0390511690565b546001600160a01b031690565b9083612512606060808401519301516001600160a01b031690565b604051633cf9798360e01b815296879586948593917f000000000000000000000000000000000000000000000000000000000000000090600486016123cc565b03925af19081156125b55760009060009261258e575b50156125715750565b6040516302a35ba360e21b815290819061214f906004830161233a565b90506125ad91503d806000833e6125a5818361021e565b81019061238d565b509038612568565b6122fe565b506125de6125da6125d560608401516001600160a01b031690565b613e61565b1590565b6124c0565b5060608101516001600160a01b03163b156124b9565b506080810151156124b2565b803b156102a657600060405180926308d450a160e01b825281838161262d8a600483016122ed565b03925af19081612667575b506126615761214f61264861230a565b6040516309c2532560e01b81529182916004830161233a565b386124a8565b80612676600061267c9361021e565b80610461565b38612638565b85965060206126be9601516126a160608901516001600160a01b031690565b906126b860208a51016001600160401b0390511690565b92613d48565b90388080808061242d565b6306e34e6560e31b60005260046000fd5b604051906126e782610192565b60006060838281528260208201528260408201520152565b9060405161270c81610192565b606060ff600183958054855201548181166020850152818160081c16604085015260101c161515910152565b906020825491828152019160005260206000209060005b81811061275c5750505090565b82546001600160a01b031684526020909301926001928301920161274f565b9061024e611f089260405193848092612738565b35906001600160e01b03821682036102a657565b81601f820112156102a6578035906127ba8261027e565b926127c8604051948561021e565b82845260208085019360061b830101918183116102a657602001925b8284106127f2575050505090565b6040848303126102a6576020604091825161280c816101cd565b8635612817816102ab565b815261282483880161278f565b838201528152019301926127e4565b81601f820112156102a65780359061284a8261027e565b92612858604051948561021e565b82845260208085019360051b830101918183116102a65760208101935b83851061288457505050505090565b84356001600160401b0381116102a657820160a0818503601f1901126102a657604051916128b1836101b2565b60208201356128bf816102ab565b83526040820135926001600160401b0384116102a65760a0836128e988602080988198010161032e565b8584015260608101356128fb816102ab565b6040840152608081013561290e816102ab565b606084015201356080820152815201940193612875565b81601f820112156102a65780359061293c8261027e565b9261294a604051948561021e565b82845260208085019360061b830101918183116102a657602001925b828410612974575050505090565b6040848303126102a6576020604091825161298e816101cd565b863581528287013583820152815201930192612966565b6020818303126102a6578035906001600160401b0382116102a657016060818303126102a657604051916129d8836101e8565b81356001600160401b0381116102a65782016040818303126102a65760405190612a01826101cd565b80356001600160401b0381116102a657810183601f820112156102a6578035612a298161027e565b91612a37604051938461021e565b81835260208084019260061b820101908682116102a657602001915b818310612acf5750505082526020810135906001600160401b0382116102a657612a7f918491016127a3565b6020820152835260208201356001600160401b0381116102a65781612aa5918401612833565b602084015260408201356001600160401b0381116102a657612ac79201612925565b604082015290565b6040838803126102a65760206040918251612ae9816101cd565b8535612af481610295565b8152612b0183870161278f565b83820152815201920191612a53565b9080602083519182815201916020808360051b8301019401926000915b838310612b3c57505050505090565b9091929394602080600192601f198582030186528851906001600160401b038251168152608080612b7a8585015160a08786015260a0850190610633565b936001600160401b0360408201511660408501526001600160401b036060820151166060850152015191015297019301930191939290612b2d565b916001600160a01b03612bd692168352606060208401526060830190612b10565b9060408183039101526020808351928381520192019060005b818110612bfc5750505090565b8251805185526020908101518186015260409094019390920191600101612bef565b906020808351928381520192019060005b818110612c3c5750505090565b825180516001600160401b031685526020908101516001600160e01b03168186015260409094019390920191600101612c2f565b9190604081019083519160408252825180915260206060830193019060005b818110612cb057505050602061034993940151906020818403910152612c1e565b825180516001600160a01b031686526020908101516001600160e01b03168187015260409095019490920191600101612c8f565b906020610349928181520190612c70565b908160209103126102a65751610349816102c7565b9091612d2161034993604084526040840190610633565b916020818403910152611e71565b6001600160401b036001911601906001600160401b038211611f9557565b9091612d6461034993604084526040840190612b10565b916020818403910152612c70565b929693959190979497612d87828201826129a5565b98612d9b6125da60045460ff9060c01c1690565b6131ef575b895180515115908115916131e0575b50613107575b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316999860208a019860005b8a5180518210156130a55781612dfe91611e23565b518d612e1182516001600160401b031690565b604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b1660048201529091602090829060249082905afa9081156125b557600091613077575b5061305a57612e6090613eaf565b60208201805160208151910120906001830191612e7c83611ef4565b602081519101200361303d575050805460408301516001600160401b039081169160a81c168114801590613015575b612fc357506080820151908115612fb257612efc82612eed612ed486516001600160401b031690565b6001600160401b0316600052600a602052604060002090565b90600052602052604060002090565b54612f7e578291612f62612f7792612f29612f2460606001999801516001600160401b031690565b612d2f565b67ffffffffffffffff60a81b197cffffffffffffffff00000000000000000000000000000000000000000083549260a81b169116179055565b612eed612ed44294516001600160401b031690565b5501612de9565b50612f93611d3f92516001600160401b031690565b6332cf0cbf60e01b6000526001600160401b0316600452602452604490565b63504570e360e01b60005260046000fd5b82611d3f91612fed6060612fde84516001600160401b031690565b9301516001600160401b031690565b636af0786b60e11b6000526001600160401b0392831660045290821660245216604452606490565b5061302d61084260608501516001600160401b031690565b6001600160401b03821611612eab565b5161214f60405192839263b80d8fa960e01b845260048401612d0a565b637edeb53960e11b6000526001600160401b031660045260246000fd5b613098915060203d811161309e575b613090818361021e565b810190612cf5565b38612e52565b503d613086565b50506131019496989b507f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e461024e9b6130f9949597999b519051906130ef60405192839283612d4d565b0390a13691610cfb565b943691610cfb565b9361429a565b61311c602086015b356001600160401b031690565b600b546001600160401b03828116911610156131c457613152906001600160401b03166001600160401b0319600b541617600b55565b61316a611b28611b286004546001600160a01b031690565b8a5190803b156102a657604051633937306f60e01b815291600091839182908490829061319a9060048301612ce4565b03925af180156125b5576131af575b50612db5565b8061267660006131be9361021e565b386131a9565b5060208a015151612db557632261116760e01b60005260046000fd5b60209150015151151538612daf565b60208a01518051613201575b50612da0565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169060408c0151823b156102a657604051633854844f60e11b81529260009284928391829161325d913060048501612bb5565b03915afa80156125b557156131fb5780612676600061327b9361021e565b386131fb565b6040519061329060208361021e565b6000808352366020840137565b6132a5613586565b60005b8151811015611d73576132bb8183611e23565b51906040820160ff6132ce825160ff1690565b161561357057602083015160ff16926132f48460ff166000526002602052604060002090565b916001830191825461330f6133098260ff1690565b60ff1690565b613535575061333c6133246060830151151590565b845462ff0000191690151560101b62ff000016178455565b60a081019182516101008151116134dd5780511561351f576003860161336a6133648261277b565b8a6153ae565b60608401516133fa575b947fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547946002946133d66133c66133f49a966133bf8760019f9c6133ba6133ec9a8f61550f565b6144db565b5160ff1690565b845460ff191660ff821617909455565b5190818555519060405195869501908886614561565b0390a1615591565b016132a8565b979460028793959701966134166134108961277b565b886153ae565b60808501519461010086511161350957855161343e6133096134398a5160ff1690565b6144c7565b10156134f35785518451116134dd576133d66133c67fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547986133bf8760019f6133ba6133f49f9a8f6134c560029f6134bf6133ec9f8f906133ba84926134a4845160ff1690565b908054909161ff001990911660089190911b61ff0016179055565b82615442565b505050979c9f50975050969a50505094509450613374565b631b3fab5160e11b600052600160045260246000fd5b631b3fab5160e11b600052600360045260246000fd5b631b3fab5160e11b600052600260045260246000fd5b631b3fab5160e11b600052600560045260246000fd5b60101c60ff1661355061354b6060840151151590565b151590565b9015151461333c576321fd80df60e21b60005260ff861660045260246000fd5b631b3fab5160e11b600090815260045260246000fd5b6001600160a01b0360015416330361359a57565b6315ae3a6f60e11b60005260046000fd5b604051602081019060008252602081526135c660408261021e565b51902090565b8181106135d7575050565b600081556001016135cc565b9190601f81116135f257505050565b61024e926000526020600020906020601f840160051c8301931061361e575b601f0160051c01906135cc565b9091508190613611565b91909182516001600160401b0381116101ad5761364f816136498454611e37565b846135e3565b6020601f8211600114613690578190613681939495600092613685575b50508160011b916000199060031b1c19161790565b9055565b01519050388061366c565b601f198216906136a584600052602060002090565b9160005b8181106136e1575095836001959697106136c8575b505050811b019055565b015160001960f88460031b161c191690553880806136be565b9192602060018192868b0151815501940192016136a9565b90600160a061034993602081526001600160401b0384546001600160a01b038116602084015260ff81851c161515604084015260a81c166060820152608080820152019101611e71565b906001600160401b03613783921660005260096020526701ffffffffffffff60406000209160071c166001600160401b0316600052602052604060002090565b5490565b7f00000000000000000000000000000000000000000000000000000000000000004681036137b25750565b630f01ce8560e01b6000526004524660245260446000fd5b91909180511561386c5782511592602091604051926137e9818561021e565b60008452601f19810160005b8181106138485750505060005b8151811015613840578061382961381b60019385611e23565b51881561382f5786906146a0565b01613802565b6138398387611e23565b51906146a0565b505050509050565b8290604051613856816101cd565b60008152606083820152828289010152016137f5565b63c2e5347d60e01b60005260046000fd5b91908260a09103126102a657604051613895816101b2565b60808082948051845260208101516138ac816102ab565b602085015260408101516138bf816102ab565b604085015260608101516138d2816102ab565b60608501520151916109b2836102ab565b519061024e82610295565b519061024e826109c1565b81601f820112156102a6578051906139108261027e565b9261391e604051948561021e565b82845260208085019360051b830101918183116102a65760208101935b83851061394a57505050505090565b84516001600160401b0381116102a657820160a0818503601f1901126102a65760405191613977836101b2565b60208201516001600160401b0381116102a6578560206139999285010161234b565b835260408201516139a981610295565b60208401526139ba606083016138ee565b60408401526080820151926001600160401b0384116102a65760a0836139e788602080988198010161234b565b60608401520151608082015281520194019361393b565b6020818303126102a6578051906001600160401b0382116102a65701610140818303126102a657613a2d61023f565b91613a38818361387d565b835260a08201516001600160401b0381116102a65781613a5991840161234b565b602084015260c08201516001600160401b0381116102a65781613a7d91840161234b565b6040840152613a8e60e083016138e3565b606084015261010082015160808401526101208201516001600160401b0381116102a657613abc92016138f9565b60a082015290565b9080602083519182815201916020808360051b8301019401926000915b838310613af057505050505090565b9091929394602080600192601f19858203018652885190608080613b53613b20855160a0865260a0860190610633565b6001600160a01b0387870151168786015263ffffffff604087015116604086015260608601518582036060870152610633565b93015191015297019301930191939290613ae1565b610349916001600160401b036080835180518452826020820151166020850152826040820151166040850152826060820151166060850152015116608082015260a0613bd9613bc7602085015161014084860152610140850190610633565b604085015184820360c0860152610633565b60608401516001600160a01b031660e0840152926080810151610100840152015190610120818403910152613ac4565b906020610349928181520190613b68565b6000613c928192604051613c2d81610203565b613c35611d77565b81526060602082015260606040820152836060820152836080820152606060a082015250613c75611b28611b28600b546001600160a01b039060401c1690565b90604051948580948193634546c6e560e01b835260048301613c09565b03925af160009181613cc8575b506103495761214f613caf61230a565b60405163828ebdfb60e01b81529182916004830161233a565b613ce69192503d806000833e613cde818361021e565b8101906139fe565b9038613c9f565b9190811015611e1e5760051b0190565b35610349816109c1565b9190811015611e1e5760051b81013590601e19813603018212156102a65701908135916001600160401b0383116102a65760200182360381136102a6579190565b90929491939796815196613d5b8861027e565b97613d69604051998a61021e565b808952613d78601f199161027e565b0160005b818110613e4a57505060005b8351811015613e3d5780613dcf8c8a8a8a613dc9613dc2878d613dbb828f8f9d8f9e60019f81613deb575b505050611e23565b5197613d07565b36916102f7565b93614ee5565b613dd9828c611e23565b52613de4818b611e23565b5001613d88565b63ffffffff613e03613dfe858585613ced565b613cfd565b1615613db357613e3392613e1a92613dfe92613ced565b6040613e268585611e23565b51019063ffffffff169052565b8f8f908391613db3565b5096985050505050505050565b602090613e556121ff565b82828d01015201613d7c565b613e726385572ffb60e01b82615248565b9081613e8c575b81613e82575090565b610349915061521a565b9050613e978161519f565b1590613e79565b613e7263aff2afbf60e01b82615248565b6001600160401b031680600052600860205260406000209060ff825460a01c1615613ed8575090565b63ed053c5960e01b60005260045260246000fd5b6084019081608411611f9557565b60a001908160a011611f9557565b91908201809211611f9557565b600311156108f757565b60038210156108f75752565b9061024e604051613f3b816101cd565b602060ff829554818116845260081c169101613f1f565b8054821015611e1e5760005260206000200190600090565b60ff60019116019060ff8211611f9557565b60ff601b9116019060ff8211611f9557565b90606092604091835260208301370190565b6001600052600260205293613fd47fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e06126ff565b93853594613fe185613eec565b6060820190613ff08251151590565b61426c575b8036036142545750815187810361423b575061400f613787565b6001600052600360205261405e6140597fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c5b336001600160a01b0316600052602052604060002090565b613f2b565b6002602082015161406e81613f15565b61407781613f15565b1490816141d3575b50156141a7575b516140de575b50505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0906140c261310f60019460200190565b604080519283526001600160401b0391909116602083015290a2565b6140ff6133096140fa602085969799989a955194015160ff1690565b613f6a565b036141965781518351036141855761417d60006140c29461310f946141497f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09960019b36916102f7565b602081519101206040516141748161416689602083019586613f8e565b03601f19810183528261021e565b5190208a615278565b94839461408c565b63a75d88af60e01b60005260046000fd5b6371253a2560e01b60005260046000fd5b72c11c11c11c11c11c11c11c11c11c11c11c11c133031561408657631b41e11d60e31b60005260046000fd5b600160005260026020526142339150611b28906142209061421a60037fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e05b01915160ff1690565b90613f52565b90546001600160a01b039160031b1c1690565b33143861407f565b6324f7d61360e21b600052600452602487905260446000fd5b638e1192e160e01b6000526004523660245260446000fd5b6142959061428f6142856142808751611f7f565b613efa565b61428f8851611f7f565b90613f08565b613ff5565b600080526002602052949093909290916142d37fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b6126ff565b948635956142e083613eec565b60608201906142ef8251151590565b6144a4575b8036036142545750815188810361448b575061430e613787565b6000805260036020526143436140597f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff614041565b6002602082015161435381613f15565b61435c81613f15565b149081614442575b5015614416575b516143a8575b5050505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0906140c261310f60009460200190565b6143c46133096140fa602087989a999b96975194015160ff1690565b03614196578351865103614185576000967f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0966140c29561414961440d9461310f9736916102f7565b94839438614371565b72c11c11c11c11c11c11c11c11c11c11c11c11c133031561436b57631b41e11d60e31b60005260046000fd5b6000805260026020526144839150611b28906142209061421a60037fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b614211565b331438614364565b6324f7d61360e21b600052600452602488905260446000fd5b6144c29061428f6144b86142808951611f7f565b61428f8a51611f7f565b6142f4565b60ff166003029060ff8216918203611f9557565b8151916001600160401b0383116101ad576801000000000000000083116101ad576020908254848455808510614544575b500190600052602060002060005b8381106145275750505050565b60019060206001600160a01b03855116940193818401550161451a565b61455b9084600052858460002091820191016135cc565b3861450c565b95949392909160ff61458693168752602087015260a0604087015260a0860190612738565b84810360608601526020808351928381520192019060005b8181106145b95750505090608061024e9294019060ff169052565b82516001600160a01b031684526020938401939092019160010161459e565b600654811015611e1e5760066000527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f015490565b6001600160401b036103499493816060941683521660208201528160408201520190610633565b604090610349939281528160208201520190610633565b9291906001600160401b039081606495166004521660245260048110156108f757604452565b94939261468a60609361469b93885260208801906108fc565b608060408701526080860190610633565b930152565b906146b282516001600160401b031690565b8151604051632cbc26bb60e01b815267ffffffffffffffff60801b608084901b1660048201529015159391906001600160401b038216906020816024817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9081156125b557600091614dce575b50614d8c576020830191825151948515614d5c57604085018051518703614d4b5761475487611dc9565b957f000000000000000000000000000000000000000000000000000000000000000061478a600161478487613eaf565b01611ef4565b602081519101206040516147ea816141666020820194868b876001600160401b036060929594938160808401977f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f85521660208401521660408201520152565b519020906001600160401b031660005b8a8110614cb357505050806080606061481a93015191015190888661579d565b978815614c955760005b8881106148375750505050505050505050565b5a614843828951611e23565b5180516060015161485d906001600160401b031688611fa7565b614866816108ed565b8015908d8283159384614c82575b15614c3f5760608815614bc2575061489b6020614891898d611e23565b5101519242611f9a565b6004546148b09060a01c63ffffffff166120d5565b108015614baf575b15614b91576148c7878b611e23565b5151614b7b575b8451608001516148e6906001600160401b0316610842565b614ac3575b506148f7868951611e23565b5160a085015151815103614a87579361495c9695938c938f9661493c8e958c9261493661493060608951016001600160401b0390511690565b896157cf565b866159d0565b9a90809661495660608851016001600160401b0390511690565b90615857565b614a35575b505061496c826108ed565b600282036149ed575b6001966149e37f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b936001600160401b039351926149d46149cb8b6149c360608801516001600160401b031690565b96519b611e23565b51985a90611f9a565b91604051958695169885614671565b0390a45b01614824565b915091939492506149fd826108ed565b60038203614a11578b929493918a91614975565b51606001516349362d1f60e11b600052611d3f91906001600160401b03168961464b565b614a3e846108ed565b60038403614961579092949550614a569193506108ed565b614a66578b92918a913880614961565b5151604051632b11b8d960e01b815290819061214f90879060048401614634565b611d3f8b614aa160608851016001600160401b0390511690565b631cfe6d8b60e01b6000526001600160401b0391821660045216602452604490565b614acc836108ed565b614ad7575b386148eb565b8351608001516001600160401b0316602080860151918c614b0c60405194859384936370701e5760e11b85526004850161460d565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af19081156125b557600091614b5d575b50614ad15750505050506001906149e7565b614b75915060203d811161309e57613090818361021e565b38614b4b565b614b85878b611e23565b515160808601526148ce565b6354e7e43160e11b6000526001600160401b038b1660045260246000fd5b50614bb9836108ed565b600383146148b8565b915083614bce846108ed565b156148ce57506001959450614c379250614c1591507f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe651209351016001600160401b0390511690565b604080516001600160401b03808c168252909216602083015290918291820190565b0390a16149e7565b505050506001929150614c37614c1560607f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c9351016001600160401b0390511690565b50614c8c836108ed565b60038314614874565b633ee8bd3f60e11b6000526001600160401b03841660045260246000fd5b614cbe818a51611e23565b518051604001516001600160401b0316838103614d2e57508051602001516001600160401b0316898103614d0b575090614cfa84600193615695565b614d04828d611e23565b52016147fa565b636c95f1eb60e01b6000526001600160401b03808a166004521660245260446000fd5b631c21951160e11b6000526001600160401b031660045260246000fd5b6357e0e08360e01b60005260046000fd5b611d3f614d7086516001600160401b031690565b63676cf24b60e11b6000526001600160401b0316600452602490565b509291505061305a576040516001600160401b039190911681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d89493390602090a1565b614de7915060203d60201161309e57613090818361021e565b3861472a565b908160209103126102a6575161034981610295565b90610349916020815260e0614ea0614e8b614e2b85516101006020870152610120860190610633565b60208601516001600160401b0316604086015260408601516001600160a01b0316606086015260608601516080860152614e75608087015160a08701906001600160a01b03169052565b60a0860151858203601f190160c0870152610633565b60c0850151848203601f190184860152610633565b92015190610100601f1982850301910152610633565b6040906001600160a01b0361034994931681528160208201520190610633565b908160209103126102a6575190565b91939293614ef16121ff565b5060208301516001600160a01b031660405163bbe4f6db60e01b81526001600160a01b038216600482015290959092602084806024810103816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9384156125b55760009461516e575b506001600160a01b038416958615801561515c575b61513e5761502361504c9261416692614fa7614fa06120d560408c015163ffffffff1690565b8c89615aaf565b9690996080810151614fd56060835193015193614fc261025f565b9687526001600160401b03166020870152565b6001600160a01b038a16604086015260608501526001600160a01b038d16608085015260a084015260c083015260e0820152604051633907753760e01b602082015292839160248301614e02565b82857f000000000000000000000000000000000000000000000000000000000000000092615b3d565b949091156151225750805160208103615109575090615075826020808a95518301019101614ed6565b956001600160a01b038416036150ad575b50505050506150a561509661026f565b6001600160a01b039093168352565b602082015290565b6150c0936150ba91611f9a565b91615aaf565b509080821080156150f6575b6150d857808481615086565b63a966e21f60e01b6000908152600493909352602452604452606490fd5b50826151028284611f9a565b14156150cc565b631e3be00960e21b600052602060045260245260446000fd5b61214f604051928392634ff17cad60e11b845260048401614eb6565b63ae9b4ce960e01b6000526001600160a01b03851660045260246000fd5b506151696125da86613e9e565b614f7a565b61519191945060203d602011615198575b615189818361021e565b810190614ded565b9238614f65565b503d61517f565b60405160208101916301ffc9a760e01b835263ffffffff60e01b6024830152602482526151cd60448361021e565b6179185a10615209576020926000925191617530fa6000513d826151fd575b50816151f6575090565b9050151590565b602011159150386151ec565b63753fa58960e11b60005260046000fd5b60405160208101916301ffc9a760e01b83526301ffc9a760e01b6024830152602482526151cd60448361021e565b6040519060208201926301ffc9a760e01b845263ffffffff60e01b166024830152602482526151cd60448361021e565b919390926000948051946000965b868810615297575050505050505050565b6020881015611e1e57602060006152af878b1a613f7c565b6152b98b87611e23565b51906152f06152c88d8a611e23565b5160405193849389859094939260ff6060936080840197845216602083015260408201520152565b838052039060015afa156125b55761533661405960005161531e8960ff166000526003602052604060002090565b906001600160a01b0316600052602052604060002090565b906001602083015161534781613f15565b61535081613f15565b0361539d5761536d615363835160ff1690565b60ff600191161b90565b811661538c576153836153636001935160ff1690565b17970196615286565b633d9ef1f160e21b60005260046000fd5b636518c33d60e11b60005260046000fd5b91909160005b83518110156154075760019060ff831660005260036020526000615400604082206001600160a01b036153e7858a611e23565b51166001600160a01b0316600052602052604060002090565b55016153b4565b50509050565b8151815460ff191660ff91909116178155906020015160038110156108f757815461ff00191660089190911b61ff0016179055565b919060005b81518110156154075761545d611c818284611e23565b9061548661547c8361531e8860ff166000526003602052604060002090565b5460081c60ff1690565b61548f81613f15565b6154fa576001600160a01b038216156154e9576154e36001926154de6154b361026f565b60ff85168152916154c78660208501613f1f565b61531e8960ff166000526003602052604060002090565b61540d565b01615447565b63d6c62c9b60e01b60005260046000fd5b631b3fab5160e11b6000526004805260246000fd5b919060005b81518110156154075761552a611c818284611e23565b9061554961547c8361531e8860ff166000526003602052604060002090565b61555281613f15565b6154fa576001600160a01b038216156154e95761558b6001926154de61557661026f565b60ff85168152916154c7600260208501613f1f565b01615514565b60ff1680600052600260205260ff60016040600020015460101c169080156000146155df5750156155ce576001600160401b0319600b5416600b55565b6317bd8dd160e11b60005260046000fd5b6001146155e95750565b6155ef57565b6307b8c74d60e51b60005260046000fd5b8060005260076020526040600020541560001461567e57600654680100000000000000008110156101ad57600181016006556000600654821015611e1e57600690527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01819055600654906000526007602052604060002055600190565b50600090565b906020610349928181520190613ac4565b6135c681518051906157296156b460608601516001600160a01b031690565b6141666156cb60608501516001600160401b031690565b936156e46080808a01519201516001600160401b031690565b90604051958694602086019889936001600160401b036080946001600160a01b0382959998949960a089019a8952166020880152166040860152606085015216910152565b5190206141666020840151602081519101209360a060408201516020815191012091015160405161576281614166602082019485615684565b51902090604051958694602086019889919260a093969594919660c08401976000855260208501526040840152606083015260808201520152565b926001600160401b03926157b092615bfa565b9116600052600a60205260406000209060005260205260406000205490565b607f8216906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611f9557615854916001600160401b036158128584613743565b921660005260096020526701ffffffffffffff60406000209460071c169160036001831b921b19161792906001600160401b0316600052602052604060002090565b55565b9091607f83166801fffffffffffffffe6001600160401b0382169160011b169080820460021490151715611f955761588f8484613743565b60048310156108f7576001600160401b036158549416600052600960205260036701ffffffffffffff60406000209660071c1693831b921b19161792906001600160401b0316600052602052604060002090565b906158f690606083526060830190613b68565b8181036020830152825180825260208201916020808360051b8301019501926000915b83831061596657505050505060408183039101526020808351928381520192019060005b81811061594a5750505090565b825163ffffffff1684526020938401939092019160010161593d565b9091929395602080615984600193601f198682030187528a51610633565b98019301930191939290615919565b80516020909101516001600160e01b03198116929190600482106159b5575050565b6001600160e01b031960049290920360031b82901b16169150565b90303b156102a6576000916159f96040519485938493630304c3e160e51b8552600485016158e3565b038183305af19081615a9a575b50615a8f57615a1361230a565b9072c11c11c11c11c11c11c11c11c11c11c11c11c13314615a35575b60039190565b615a4e615a4183615993565b6001600160e01b03191690565b6337c3be2960e01b148015615a74575b15615a2f57631d1ccf9f60e01b60005260046000fd5b50615a81615a4183615993565b632be8ca8b60e21b14615a5e565b6002906103496105fb565b806126766000615aa99361021e565b38615a06565b6040516370a0823160e01b60208201526001600160a01b039091166024820152919291615b0c90615ae38160448101614166565b84837f000000000000000000000000000000000000000000000000000000000000000092615b3d565b929091156151225750805160208103615109575090615b378260208061034995518301019101614ed6565b93611f9a565b939193615b4a60846102dc565b94615b58604051968761021e565b60848652615b6660846102dc565b602087019590601f1901368737833b15615be9575a90808210615bd8578291038060061c90031115615bc7576000918291825a9560208451940192f1905a9003923d9060848211615bbe575b6000908287523e929190565b60849150615bb2565b6337c3be2960e01b60005260046000fd5b632be8ca8b60e21b60005260046000fd5b63030ed58f60e21b60005260046000fd5b8051928251908415615d565761010185111580615d4a575b15615c7957818501946000198601956101008711615c79578615615d3a57615c3987611dc9565b9660009586978795885b848110615c9e575050505050600119018095149384615c94575b505082615c8a575b505015615c7957615c7591611e23565b5190565b6309bde33960e01b60005260046000fd5b1490503880615c65565b1492503880615c5d565b6001811b82811603615d2c57868a1015615d1757615cc060018b019a85611e23565b51905b8c888c1015615d035750615cdb60018c019b86611e23565b515b818d11615c7957615cfc828f92615cf690600196615d67565b92611e23565b5201615c43565b60018d019c615d1191611e23565b51615cdd565b615d2560018c019b8d611e23565b5190615cc3565b615d25600189019884611e23565b505050509050615c759150611e11565b50610101821115615c12565b630469ac9960e21b60005260046000fd5b81811015615d79579061034991615d7e565b610349915b906040519060208201926001845260408301526060820152606081526135c660808261021e56fea164736f6c634300081a000a49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b",
}

var OffRampWithMessageTransformerABI = OffRampWithMessageTransformerMetaData.ABI

var OffRampWithMessageTransformerBin = OffRampWithMessageTransformerMetaData.Bin

func DeployOffRampWithMessageTransformer(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig OffRampStaticConfig, dynamicConfig OffRampDynamicConfig, sourceChainConfigs []OffRampSourceChainConfigArgs, messageTransformerAddr common.Address) (common.Address, *types.Transaction, *OffRampWithMessageTransformer, error) {
	parsed, err := OffRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OffRampWithMessageTransformerBin), backend, staticConfig, dynamicConfig, sourceChainConfigs, messageTransformerAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffRampWithMessageTransformer{address: address, abi: *parsed, OffRampWithMessageTransformerCaller: OffRampWithMessageTransformerCaller{contract: contract}, OffRampWithMessageTransformerTransactor: OffRampWithMessageTransformerTransactor{contract: contract}, OffRampWithMessageTransformerFilterer: OffRampWithMessageTransformerFilterer{contract: contract}}, nil
}

type OffRampWithMessageTransformer struct {
	address common.Address
	abi     abi.ABI
	OffRampWithMessageTransformerCaller
	OffRampWithMessageTransformerTransactor
	OffRampWithMessageTransformerFilterer
}

type OffRampWithMessageTransformerCaller struct {
	contract *bind.BoundContract
}

type OffRampWithMessageTransformerTransactor struct {
	contract *bind.BoundContract
}

type OffRampWithMessageTransformerFilterer struct {
	contract *bind.BoundContract
}

type OffRampWithMessageTransformerSession struct {
	Contract     *OffRampWithMessageTransformer
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OffRampWithMessageTransformerCallerSession struct {
	Contract *OffRampWithMessageTransformerCaller
	CallOpts bind.CallOpts
}

type OffRampWithMessageTransformerTransactorSession struct {
	Contract     *OffRampWithMessageTransformerTransactor
	TransactOpts bind.TransactOpts
}

type OffRampWithMessageTransformerRaw struct {
	Contract *OffRampWithMessageTransformer
}

type OffRampWithMessageTransformerCallerRaw struct {
	Contract *OffRampWithMessageTransformerCaller
}

type OffRampWithMessageTransformerTransactorRaw struct {
	Contract *OffRampWithMessageTransformerTransactor
}

func NewOffRampWithMessageTransformer(address common.Address, backend bind.ContractBackend) (*OffRampWithMessageTransformer, error) {
	abi, err := abi.JSON(strings.NewReader(OffRampWithMessageTransformerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOffRampWithMessageTransformer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformer{address: address, abi: abi, OffRampWithMessageTransformerCaller: OffRampWithMessageTransformerCaller{contract: contract}, OffRampWithMessageTransformerTransactor: OffRampWithMessageTransformerTransactor{contract: contract}, OffRampWithMessageTransformerFilterer: OffRampWithMessageTransformerFilterer{contract: contract}}, nil
}

func NewOffRampWithMessageTransformerCaller(address common.Address, caller bind.ContractCaller) (*OffRampWithMessageTransformerCaller, error) {
	contract, err := bindOffRampWithMessageTransformer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerCaller{contract: contract}, nil
}

func NewOffRampWithMessageTransformerTransactor(address common.Address, transactor bind.ContractTransactor) (*OffRampWithMessageTransformerTransactor, error) {
	contract, err := bindOffRampWithMessageTransformer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerTransactor{contract: contract}, nil
}

func NewOffRampWithMessageTransformerFilterer(address common.Address, filterer bind.ContractFilterer) (*OffRampWithMessageTransformerFilterer, error) {
	contract, err := bindOffRampWithMessageTransformer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerFilterer{contract: contract}, nil
}

func bindOffRampWithMessageTransformer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OffRampWithMessageTransformerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRampWithMessageTransformer.Contract.OffRampWithMessageTransformerCaller.contract.Call(opts, result, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.OffRampWithMessageTransformerTransactor.contract.Transfer(opts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.OffRampWithMessageTransformerTransactor.contract.Transact(opts, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRampWithMessageTransformer.Contract.contract.Call(opts, result, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.contract.Transfer(opts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.contract.Transact(opts, method, params...)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "ccipReceive", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRampWithMessageTransformer.Contract.CcipReceive(&_OffRampWithMessageTransformer.CallOpts, arg0)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _OffRampWithMessageTransformer.Contract.CcipReceive(&_OffRampWithMessageTransformer.CallOpts, arg0)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getAllSourceChainConfigs")

	if err != nil {
		return *new([]uint64), *new([]OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)
	out1 := *abi.ConvertType(out[1], new([]OffRampSourceChainConfig)).(*[]OffRampSourceChainConfig)

	return out0, out1, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetAllSourceChainConfigs(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetAllSourceChainConfigs(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(OffRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampDynamicConfig)).(*OffRampDynamicConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetDynamicConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetDynamicConfig() (OffRampDynamicConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetDynamicConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getExecutionState", sourceChainSelector, sequenceNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRampWithMessageTransformer.Contract.GetExecutionState(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetExecutionState(sourceChainSelector uint64, sequenceNumber uint64) (uint8, error) {
	return _OffRampWithMessageTransformer.Contract.GetExecutionState(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, sequenceNumber)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getLatestPriceSequenceNumber")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRampWithMessageTransformer.Contract.GetLatestPriceSequenceNumber(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetLatestPriceSequenceNumber() (uint64, error) {
	return _OffRampWithMessageTransformer.Contract.GetLatestPriceSequenceNumber(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getMerkleRoot", sourceChainSelector, root)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRampWithMessageTransformer.Contract.GetMerkleRoot(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, root)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetMerkleRoot(sourceChainSelector uint64, root [32]byte) (*big.Int, error) {
	return _OffRampWithMessageTransformer.Contract.GetMerkleRoot(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector, root)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetMessageTransformer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getMessageTransformer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetMessageTransformer() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.GetMessageTransformer(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetMessageTransformer() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.GetMessageTransformer(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampSourceChainConfig)).(*OffRampSourceChainConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetSourceChainConfig(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (OffRampSourceChainConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetSourceChainConfig(&_OffRampWithMessageTransformer.CallOpts, sourceChainSelector)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(OffRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(OffRampStaticConfig)).(*OffRampStaticConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetStaticConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) GetStaticConfig() (OffRampStaticConfig, error) {
	return _OffRampWithMessageTransformer.Contract.GetStaticConfig(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "latestConfigDetails", ocrPluginType)

	if err != nil {
		return *new(MultiOCR3BaseOCRConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(MultiOCR3BaseOCRConfig)).(*MultiOCR3BaseOCRConfig)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRampWithMessageTransformer.Contract.LatestConfigDetails(&_OffRampWithMessageTransformer.CallOpts, ocrPluginType)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) LatestConfigDetails(ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error) {
	return _OffRampWithMessageTransformer.Contract.LatestConfigDetails(&_OffRampWithMessageTransformer.CallOpts, ocrPluginType)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) Owner() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.Owner(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) Owner() (common.Address, error) {
	return _OffRampWithMessageTransformer.Contract.Owner(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OffRampWithMessageTransformer.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) TypeAndVersion() (string, error) {
	return _OffRampWithMessageTransformer.Contract.TypeAndVersion(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerCallerSession) TypeAndVersion() (string, error) {
	return _OffRampWithMessageTransformer.Contract.TypeAndVersion(&_OffRampWithMessageTransformer.CallOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "acceptOwnership")
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.AcceptOwnership(&_OffRampWithMessageTransformer.TransactOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.AcceptOwnership(&_OffRampWithMessageTransformer.TransactOpts)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "applySourceChainConfigUpdates", sourceChainConfigUpdates)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ApplySourceChainConfigUpdates(&_OffRampWithMessageTransformer.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ApplySourceChainConfigUpdates(&_OffRampWithMessageTransformer.TransactOpts, sourceChainConfigUpdates)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) Commit(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "commit", reportContext, report, rs, ss, rawVs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Commit(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Commit(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) Execute(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "execute", reportContext, report)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Execute(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.Execute(&_OffRampWithMessageTransformer.TransactOpts, reportContext, report)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData, tokenGasOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ExecuteSingleMessage(&_OffRampWithMessageTransformer.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ExecuteSingleMessage(&_OffRampWithMessageTransformer.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "manuallyExecute", reports, gasLimitOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ManuallyExecute(&_OffRampWithMessageTransformer.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.ManuallyExecute(&_OffRampWithMessageTransformer.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "setDynamicConfig", dynamicConfig)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetDynamicConfig(&_OffRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) SetDynamicConfig(dynamicConfig OffRampDynamicConfig) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetDynamicConfig(&_OffRampWithMessageTransformer.TransactOpts, dynamicConfig)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) SetMessageTransformer(opts *bind.TransactOpts, messageTransformerAddr common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "setMessageTransformer", messageTransformerAddr)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) SetMessageTransformer(messageTransformerAddr common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetMessageTransformer(&_OffRampWithMessageTransformer.TransactOpts, messageTransformerAddr)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) SetMessageTransformer(messageTransformerAddr common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetMessageTransformer(&_OffRampWithMessageTransformer.TransactOpts, messageTransformerAddr)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "setOCR3Configs", ocrConfigArgs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetOCR3Configs(&_OffRampWithMessageTransformer.TransactOpts, ocrConfigArgs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) SetOCR3Configs(ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.SetOCR3Configs(&_OffRampWithMessageTransformer.TransactOpts, ocrConfigArgs)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.contract.Transact(opts, "transferOwnership", to)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.TransferOwnership(&_OffRampWithMessageTransformer.TransactOpts, to)
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRampWithMessageTransformer.Contract.TransferOwnership(&_OffRampWithMessageTransformer.TransactOpts, to)
}

type OffRampWithMessageTransformerAlreadyAttemptedIterator struct {
	Event *OffRampWithMessageTransformerAlreadyAttempted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerAlreadyAttemptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerAlreadyAttempted)
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
		it.Event = new(OffRampWithMessageTransformerAlreadyAttempted)
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

func (it *OffRampWithMessageTransformerAlreadyAttemptedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerAlreadyAttemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerAlreadyAttempted struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerAlreadyAttemptedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerAlreadyAttemptedIterator{contract: _OffRampWithMessageTransformer.contract, event: "AlreadyAttempted", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerAlreadyAttempted) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerAlreadyAttempted)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseAlreadyAttempted(log types.Log) (*OffRampWithMessageTransformerAlreadyAttempted, error) {
	event := new(OffRampWithMessageTransformerAlreadyAttempted)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerCommitReportAcceptedIterator struct {
	Event *OffRampWithMessageTransformerCommitReportAccepted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerCommitReportAcceptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerCommitReportAccepted)
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
		it.Event = new(OffRampWithMessageTransformerCommitReportAccepted)
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

func (it *OffRampWithMessageTransformerCommitReportAcceptedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerCommitReportAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerCommitReportAccepted struct {
	MerkleRoots  []InternalMerkleRoot
	PriceUpdates InternalPriceUpdates
	Raw          types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerCommitReportAcceptedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerCommitReportAcceptedIterator{contract: _OffRampWithMessageTransformer.contract, event: "CommitReportAccepted", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerCommitReportAccepted) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "CommitReportAccepted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerCommitReportAccepted)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseCommitReportAccepted(log types.Log) (*OffRampWithMessageTransformerCommitReportAccepted, error) {
	event := new(OffRampWithMessageTransformerCommitReportAccepted)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "CommitReportAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerConfigSetIterator struct {
	Event *OffRampWithMessageTransformerConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerConfigSet)
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

func (it *OffRampWithMessageTransformerConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerConfigSet struct {
	OcrPluginType uint8
	ConfigDigest  [32]byte
	Signers       []common.Address
	Transmitters  []common.Address
	F             uint8
	Raw           types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerConfigSetIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseConfigSet(log types.Log) (*OffRampWithMessageTransformerConfigSet, error) {
	event := new(OffRampWithMessageTransformerConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerDynamicConfigSetIterator struct {
	Event *OffRampWithMessageTransformerDynamicConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerDynamicConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerDynamicConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerDynamicConfigSet)
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

func (it *OffRampWithMessageTransformerDynamicConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerDynamicConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerDynamicConfigSet struct {
	DynamicConfig OffRampDynamicConfig
	Raw           types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerDynamicConfigSetIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerDynamicConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "DynamicConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerDynamicConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "DynamicConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerDynamicConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseDynamicConfigSet(log types.Log) (*OffRampWithMessageTransformerDynamicConfigSet, error) {
	event := new(OffRampWithMessageTransformerDynamicConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "DynamicConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerExecutionStateChangedIterator struct {
	Event *OffRampWithMessageTransformerExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerExecutionStateChanged)
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
		it.Event = new(OffRampWithMessageTransformerExecutionStateChanged)
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

func (it *OffRampWithMessageTransformerExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           [32]byte
	MessageHash         [32]byte
	State               uint8
	ReturnData          []byte
	GasUsed             *big.Int
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampWithMessageTransformerExecutionStateChangedIterator, error) {

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

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerExecutionStateChangedIterator{contract: _OffRampWithMessageTransformer.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "ExecutionStateChanged", sourceChainSelectorRule, sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerExecutionStateChanged)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseExecutionStateChanged(log types.Log) (*OffRampWithMessageTransformerExecutionStateChanged, error) {
	event := new(OffRampWithMessageTransformerExecutionStateChanged)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerOwnershipTransferRequestedIterator struct {
	Event *OffRampWithMessageTransformerOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerOwnershipTransferRequested)
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
		it.Event = new(OffRampWithMessageTransformerOwnershipTransferRequested)
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

func (it *OffRampWithMessageTransformerOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerOwnershipTransferRequestedIterator{contract: _OffRampWithMessageTransformer.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerOwnershipTransferRequested)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferRequested, error) {
	event := new(OffRampWithMessageTransformerOwnershipTransferRequested)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerOwnershipTransferredIterator struct {
	Event *OffRampWithMessageTransformerOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerOwnershipTransferred)
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
		it.Event = new(OffRampWithMessageTransformerOwnershipTransferred)
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

func (it *OffRampWithMessageTransformerOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerOwnershipTransferredIterator{contract: _OffRampWithMessageTransformer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerOwnershipTransferred)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseOwnershipTransferred(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferred, error) {
	event := new(OffRampWithMessageTransformerOwnershipTransferred)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerRootRemovedIterator struct {
	Event *OffRampWithMessageTransformerRootRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerRootRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerRootRemoved)
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
		it.Event = new(OffRampWithMessageTransformerRootRemoved)
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

func (it *OffRampWithMessageTransformerRootRemovedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerRootRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerRootRemoved struct {
	Root [32]byte
	Raw  types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterRootRemoved(opts *bind.FilterOpts) (*OffRampWithMessageTransformerRootRemovedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerRootRemovedIterator{contract: _OffRampWithMessageTransformer.contract, event: "RootRemoved", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerRootRemoved) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "RootRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerRootRemoved)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "RootRemoved", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseRootRemoved(log types.Log) (*OffRampWithMessageTransformerRootRemoved, error) {
	event := new(OffRampWithMessageTransformerRootRemoved)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "RootRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator struct {
	Event *OffRampWithMessageTransformerSkippedAlreadyExecutedMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
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
		it.Event = new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
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

func (it *OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSkippedAlreadyExecutedMessage struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator{contract: _OffRampWithMessageTransformer.contract, event: "SkippedAlreadyExecutedMessage", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedAlreadyExecutedMessage) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SkippedAlreadyExecutedMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessage, error) {
	event := new(OffRampWithMessageTransformerSkippedAlreadyExecutedMessage)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSkippedReportExecutionIterator struct {
	Event *OffRampWithMessageTransformerSkippedReportExecution

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSkippedReportExecutionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSkippedReportExecution)
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
		it.Event = new(OffRampWithMessageTransformerSkippedReportExecution)
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

func (it *OffRampWithMessageTransformerSkippedReportExecutionIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSkippedReportExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSkippedReportExecution struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedReportExecutionIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSkippedReportExecutionIterator{contract: _OffRampWithMessageTransformer.contract, event: "SkippedReportExecution", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedReportExecution) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSkippedReportExecution)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSkippedReportExecution(log types.Log) (*OffRampWithMessageTransformerSkippedReportExecution, error) {
	event := new(OffRampWithMessageTransformerSkippedReportExecution)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSourceChainConfigSetIterator struct {
	Event *OffRampWithMessageTransformerSourceChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSourceChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSourceChainConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerSourceChainConfigSet)
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

func (it *OffRampWithMessageTransformerSourceChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSourceChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceConfig        OffRampSourceChainConfig
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampWithMessageTransformerSourceChainConfigSetIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSourceChainConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "SourceChainConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSourceChainConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSourceChainConfigSet(log types.Log) (*OffRampWithMessageTransformerSourceChainConfigSet, error) {
	event := new(OffRampWithMessageTransformerSourceChainConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerSourceChainSelectorAddedIterator struct {
	Event *OffRampWithMessageTransformerSourceChainSelectorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerSourceChainSelectorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerSourceChainSelectorAdded)
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
		it.Event = new(OffRampWithMessageTransformerSourceChainSelectorAdded)
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

func (it *OffRampWithMessageTransformerSourceChainSelectorAddedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerSourceChainSelectorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerSourceChainSelectorAdded struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSourceChainSelectorAddedIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerSourceChainSelectorAddedIterator{contract: _OffRampWithMessageTransformer.contract, event: "SourceChainSelectorAdded", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainSelectorAdded) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerSourceChainSelectorAdded)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseSourceChainSelectorAdded(log types.Log) (*OffRampWithMessageTransformerSourceChainSelectorAdded, error) {
	event := new(OffRampWithMessageTransformerSourceChainSelectorAdded)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerStaticConfigSetIterator struct {
	Event *OffRampWithMessageTransformerStaticConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerStaticConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerStaticConfigSet)
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
		it.Event = new(OffRampWithMessageTransformerStaticConfigSet)
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

func (it *OffRampWithMessageTransformerStaticConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerStaticConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerStaticConfigSet struct {
	StaticConfig OffRampStaticConfig
	Raw          types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerStaticConfigSetIterator, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerStaticConfigSetIterator{contract: _OffRampWithMessageTransformer.contract, event: "StaticConfigSet", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerStaticConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "StaticConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerStaticConfigSet)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseStaticConfigSet(log types.Log) (*OffRampWithMessageTransformerStaticConfigSet, error) {
	event := new(OffRampWithMessageTransformerStaticConfigSet)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "StaticConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampWithMessageTransformerTransmittedIterator struct {
	Event *OffRampWithMessageTransformerTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampWithMessageTransformerTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampWithMessageTransformerTransmitted)
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
		it.Event = new(OffRampWithMessageTransformerTransmitted)
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

func (it *OffRampWithMessageTransformerTransmittedIterator) Error() error {
	return it.fail
}

func (it *OffRampWithMessageTransformerTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampWithMessageTransformerTransmitted struct {
	OcrPluginType  uint8
	ConfigDigest   [32]byte
	SequenceNumber uint64
	Raw            types.Log
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampWithMessageTransformerTransmittedIterator, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.FilterLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return &OffRampWithMessageTransformerTransmittedIterator{contract: _OffRampWithMessageTransformer.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerTransmitted, ocrPluginType []uint8) (event.Subscription, error) {

	var ocrPluginTypeRule []interface{}
	for _, ocrPluginTypeItem := range ocrPluginType {
		ocrPluginTypeRule = append(ocrPluginTypeRule, ocrPluginTypeItem)
	}

	logs, sub, err := _OffRampWithMessageTransformer.contract.WatchLogs(opts, "Transmitted", ocrPluginTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampWithMessageTransformerTransmitted)
				if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformerFilterer) ParseTransmitted(log types.Log) (*OffRampWithMessageTransformerTransmitted, error) {
	event := new(OffRampWithMessageTransformerTransmitted)
	if err := _OffRampWithMessageTransformer.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformer) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OffRampWithMessageTransformer.abi.Events["AlreadyAttempted"].ID:
		return _OffRampWithMessageTransformer.ParseAlreadyAttempted(log)
	case _OffRampWithMessageTransformer.abi.Events["CommitReportAccepted"].ID:
		return _OffRampWithMessageTransformer.ParseCommitReportAccepted(log)
	case _OffRampWithMessageTransformer.abi.Events["ConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["DynamicConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseDynamicConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["ExecutionStateChanged"].ID:
		return _OffRampWithMessageTransformer.ParseExecutionStateChanged(log)
	case _OffRampWithMessageTransformer.abi.Events["OwnershipTransferRequested"].ID:
		return _OffRampWithMessageTransformer.ParseOwnershipTransferRequested(log)
	case _OffRampWithMessageTransformer.abi.Events["OwnershipTransferred"].ID:
		return _OffRampWithMessageTransformer.ParseOwnershipTransferred(log)
	case _OffRampWithMessageTransformer.abi.Events["RootRemoved"].ID:
		return _OffRampWithMessageTransformer.ParseRootRemoved(log)
	case _OffRampWithMessageTransformer.abi.Events["SkippedAlreadyExecutedMessage"].ID:
		return _OffRampWithMessageTransformer.ParseSkippedAlreadyExecutedMessage(log)
	case _OffRampWithMessageTransformer.abi.Events["SkippedReportExecution"].ID:
		return _OffRampWithMessageTransformer.ParseSkippedReportExecution(log)
	case _OffRampWithMessageTransformer.abi.Events["SourceChainConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseSourceChainConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["SourceChainSelectorAdded"].ID:
		return _OffRampWithMessageTransformer.ParseSourceChainSelectorAdded(log)
	case _OffRampWithMessageTransformer.abi.Events["StaticConfigSet"].ID:
		return _OffRampWithMessageTransformer.ParseStaticConfigSet(log)
	case _OffRampWithMessageTransformer.abi.Events["Transmitted"].ID:
		return _OffRampWithMessageTransformer.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OffRampWithMessageTransformerAlreadyAttempted) Topic() common.Hash {
	return common.HexToHash("0x3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120")
}

func (OffRampWithMessageTransformerCommitReportAccepted) Topic() common.Hash {
	return common.HexToHash("0x35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4")
}

func (OffRampWithMessageTransformerConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (OffRampWithMessageTransformerDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0xcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee")
}

func (OffRampWithMessageTransformerExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b")
}

func (OffRampWithMessageTransformerOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OffRampWithMessageTransformerOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OffRampWithMessageTransformerRootRemoved) Topic() common.Hash {
	return common.HexToHash("0x202f1139a3e334b6056064c0e9b19fd07e44a88d8f6e5ded571b24cf8c371f12")
}

func (OffRampWithMessageTransformerSkippedAlreadyExecutedMessage) Topic() common.Hash {
	return common.HexToHash("0x3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c")
}

func (OffRampWithMessageTransformerSkippedReportExecution) Topic() common.Hash {
	return common.HexToHash("0xaab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d894933")
}

func (OffRampWithMessageTransformerSourceChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b")
}

func (OffRampWithMessageTransformerSourceChainSelectorAdded) Topic() common.Hash {
	return common.HexToHash("0xf4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9")
}

func (OffRampWithMessageTransformerStaticConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b7")
}

func (OffRampWithMessageTransformerTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_OffRampWithMessageTransformer *OffRampWithMessageTransformer) Address() common.Address {
	return _OffRampWithMessageTransformer.address
}

type OffRampWithMessageTransformerInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error)

	GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error)

	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error)

	GetMessageTransformer(opts *bind.CallOpts) (common.Address, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (OffRampSourceChainConfig, error)

	GetStaticConfig(opts *bind.CallOpts) (OffRampStaticConfig, error)

	LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error)

	Commit(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Execute(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error)

	ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error)

	ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error)

	SetDynamicConfig(opts *bind.TransactOpts, dynamicConfig OffRampDynamicConfig) (*types.Transaction, error)

	SetMessageTransformer(opts *bind.TransactOpts, messageTransformerAddr common.Address) (*types.Transaction, error)

	SetOCR3Configs(opts *bind.TransactOpts, ocrConfigArgs []MultiOCR3BaseOCRConfigArgs) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAlreadyAttempted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerAlreadyAttemptedIterator, error)

	WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerAlreadyAttempted) (event.Subscription, error)

	ParseAlreadyAttempted(log types.Log) (*OffRampWithMessageTransformerAlreadyAttempted, error)

	FilterCommitReportAccepted(opts *bind.FilterOpts) (*OffRampWithMessageTransformerCommitReportAcceptedIterator, error)

	WatchCommitReportAccepted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerCommitReportAccepted) (event.Subscription, error)

	ParseCommitReportAccepted(log types.Log) (*OffRampWithMessageTransformerCommitReportAccepted, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OffRampWithMessageTransformerConfigSet, error)

	FilterDynamicConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerDynamicConfigSetIterator, error)

	WatchDynamicConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerDynamicConfigSet) (event.Subscription, error)

	ParseDynamicConfigSet(log types.Log) (*OffRampWithMessageTransformerDynamicConfigSet, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (*OffRampWithMessageTransformerExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerExecutionStateChanged, sourceChainSelector []uint64, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*OffRampWithMessageTransformerExecutionStateChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampWithMessageTransformerOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OffRampWithMessageTransformerOwnershipTransferred, error)

	FilterRootRemoved(opts *bind.FilterOpts) (*OffRampWithMessageTransformerRootRemovedIterator, error)

	WatchRootRemoved(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerRootRemoved) (event.Subscription, error)

	ParseRootRemoved(log types.Log) (*OffRampWithMessageTransformerRootRemoved, error)

	FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessageIterator, error)

	WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedAlreadyExecutedMessage) (event.Subscription, error)

	ParseSkippedAlreadyExecutedMessage(log types.Log) (*OffRampWithMessageTransformerSkippedAlreadyExecutedMessage, error)

	FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSkippedReportExecutionIterator, error)

	WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSkippedReportExecution) (event.Subscription, error)

	ParseSkippedReportExecution(log types.Log) (*OffRampWithMessageTransformerSkippedReportExecution, error)

	FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*OffRampWithMessageTransformerSourceChainConfigSetIterator, error)

	WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error)

	ParseSourceChainConfigSet(log types.Log) (*OffRampWithMessageTransformerSourceChainConfigSet, error)

	FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*OffRampWithMessageTransformerSourceChainSelectorAddedIterator, error)

	WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerSourceChainSelectorAdded) (event.Subscription, error)

	ParseSourceChainSelectorAdded(log types.Log) (*OffRampWithMessageTransformerSourceChainSelectorAdded, error)

	FilterStaticConfigSet(opts *bind.FilterOpts) (*OffRampWithMessageTransformerStaticConfigSetIterator, error)

	WatchStaticConfigSet(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerStaticConfigSet) (event.Subscription, error)

	ParseStaticConfigSet(log types.Log) (*OffRampWithMessageTransformerStaticConfigSet, error)

	FilterTransmitted(opts *bind.FilterOpts, ocrPluginType []uint8) (*OffRampWithMessageTransformerTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OffRampWithMessageTransformerTransmitted, ocrPluginType []uint8) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*OffRampWithMessageTransformerTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
