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

var OffRampMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applySourceChainConfigUpdates\",\"inputs\":[{\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"commit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeSingleMessage\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structInternal.Any2EVMRampMessage\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllSourceChainConfigs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDynamicConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExecutionState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestPriceSequenceNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMerkleRoot\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSourceChainConfig\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStaticConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"ocrConfig\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"components\":[{\"name\":\"configInfo\",\"type\":\"tuple\",\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"n\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyExecute\",\"inputs\":[{\"name\":\"reports\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.ExecutionReport[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"components\":[{\"name\":\"header\",\"type\":\"tuple\",\"internalType\":\"structInternal.RampMessageHeader\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"components\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"destGasAmount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"proofs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"proofFlagBits\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\",\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"components\":[{\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDynamicConfig\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOCR3Configs\",\"inputs\":[{\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\",\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"components\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"F\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AlreadyAttempted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitReportAccepted\",\"inputs\":[{\"name\":\"merkleRoots\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"components\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRampAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"priceUpdates\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"components\":[{\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"components\":[{\"name\":\"sourceToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdPerToken\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]},{\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structInternal.GasPriceUpdate[]\",\"components\":[{\"name\":\"destChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"usdPerUnitGas\",\"type\":\"uint224\",\"internalType\":\"uint224\"}]}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"F\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DynamicConfigSet\",\"inputs\":[{\"name\":\"dynamicConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"components\":[{\"name\":\"feeQuoter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"messageInterceptor\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecutionStateChanged\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"messageHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"state\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\"},{\"name\":\"returnData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"gasUsed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RootRemoved\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedAlreadyExecutedMessage\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkippedReportExecution\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainConfigSet\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sourceConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"components\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"contractIRouter\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"minSeqNr\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SourceChainSelectorAdded\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaticConfigSet\",\"inputs\":[{\"name\":\"staticConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"components\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"rmnRemote\",\"type\":\"address\",\"internalType\":\"contractIRMNRemote\"},{\"name\":\"tokenAdminRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonceManager\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CanOnlySelfCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CommitOnRampMismatch\",\"inputs\":[{\"name\":\"reportOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"configOnRamp\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"EmptyBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyReport\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ExecutionError\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ForkedChain\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientGasToCompleteTx\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[{\"name\":\"errorType\",\"type\":\"uint8\",\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\"}]},{\"type\":\"error\",\"name\":\"InvalidDataLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInterval\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"min\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"max\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionGasLimit\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"newLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidManualExecutionTokenGasOverride\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oldLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenGasOverride\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidMessageDestChainSelector\",\"inputs\":[{\"name\":\"messageDestChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidNewState\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newState\",\"type\":\"uint8\",\"internalType\":\"enumInternal.MessageExecutionState\"}]},{\"type\":\"error\",\"name\":\"InvalidOnRampUpdate\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LeavesCannotBeEmpty\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionGasAmountCountMismatch\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ManualExecutionGasLimitMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ManualExecutionNotYetEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"MessageValidationError\",\"inputs\":[{\"name\":\"errorReason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonUniqueSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotACompatiblePool\",\"inputs\":[{\"name\":\"notPool\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OracleCannotBeZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReceiverError\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ReleaseOrMintBalanceMismatch\",\"inputs\":[{\"name\":\"amountReleased\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePre\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balancePost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RootAlreadyCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"RootNotCommitted\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignaturesOutOfRegistration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SourceChainNotEnabled\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SourceChainSelectorMismatch\",\"inputs\":[{\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StaleCommitReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaticConfigCannotBeChanged\",\"inputs\":[{\"name\":\"ocrPluginType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"TokenDataMismatch\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sequenceNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TokenHandlingError\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnexpectedTokenData\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongMessageLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WrongNumberOfSignatures\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroChainSelectorNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610140806040523461084857616424803803809161001d828561087e565b833981019080820361014081126108485760a08112610848576040519060a082016001600160401b0381118382101761084d5760405261005c836108a1565b825260208301519261ffff84168403610848576020830193845260408101516001600160a01b0381168103610848576040840190815261009e606083016108b5565b946060850195865260806100b38185016108b5565b86820190815294609f19011261084857604051946100d086610863565b6100dc60a085016108b5565b865260c08401519463ffffffff86168603610848576020870195865261010460e086016108c9565b976040880198895261011961010087016108b5565b6060890190815261012087015190966001600160401b03821161084857018a601f820112156108485780519a6001600160401b038c1161084d578b60051b916020806040519e8f9061016d8388018361087e565b81520193820101908282116108485760208101935b828510610748575050505050331561073757600180546001600160a01b031916331790554660805284516001600160a01b0316158015610725575b8015610713575b6106f15782516001600160401b0316156107025782516001600160401b0390811660a090815286516001600160a01b0390811660c0528351811660e0528451811661010052865161ffff90811661012052604080519751909416875296519096166020860152955185169084015251831660608301525190911660808201527fb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b79190a182516001600160a01b0316156106f1579151600480548351865160ff60c01b90151560c01b1663ffffffff60a01b60a09290921b919091166001600160a01b039485166001600160c81b0319909316831717179091558351600580549184166001600160a01b031990921691909117905560408051918252925163ffffffff166020820152935115159184019190915290511660608201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee90608090a160005b815181101561066b576020600582901b8301810151908101516001600160401b031690600090821561065c5780516001600160a01b03161561064d57828252600860205260408220906060810151600183019361038585546108d6565b6105ee578354600160a81b600160e81b031916600160a81b1784556040518681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb990602090a15b815180159081156105c3575b506105b4578151916001600160401b0383116105a0576103f986546108d6565b601f811161055b575b50602091601f84116001146104e257926001989796949281926000805160206164048339815191529795926104d7575b5050600019600383901b1c191690881b1783555b60408101518254915160a089811b8a9003801960ff60a01b1990951693151590911b60ff60a01b169290921792909216911617815561048484610993565b506104ce6040519283926020845254888060a01b038116602085015260ff8160a01c1615156040850152888060401b039060a81c16606084015260808084015260a0830190610910565b0390a201610328565b015190503880610432565b9190601f198416878452828420935b8181106105435750926001999897959392859260008051602061640483398151915298968c951061052a575b505050811b018355610446565b015160001960f88460031b161c1916905538808061051d565b929360206001819287860151815501950193016104f1565b86835260208320601f850160051c81019160208610610596575b601f0160051c01905b81811061058b5750610402565b83815560010161057e565b9091508190610575565b634e487b7160e01b82526041600452602482fd5b6342bcdf7f60e11b8152600490fd5b905060208301206040516020810190838252602081526105e460408261087e565b51902014386103d9565b835460a81c6001600160401b0316600114158061061f575b156103cd57632105803760e11b81526004869052602490fd5b50604051610638816106318189610910565b038261087e565b60208151910120825160208401201415610606565b6342bcdf7f60e11b8252600482fd5b63c656089560e01b8252600482fd5b6040516159dd9081610a2782396080518161368c015260a05181818161048e01526141ea015260c0518181816104e401528181612cba0152818161310e0152614184015260e05181818161051301526149c701526101005181818161054201526145ad0152610120518181816104b50152818161243401528181614aba01526157120152f35b6342bcdf7f60e11b60005260046000fd5b63c656089560e01b60005260046000fd5b5081516001600160a01b0316156101c4565b5080516001600160a01b0316156101bd565b639b15e16f60e01b60005260046000fd5b84516001600160401b0381116108485782016080818603601f190112610848576040519061077582610863565b60208101516001600160a01b0381168103610848578252610798604082016108a1565b60208301526107a9606082016108c9565b604083015260808101516001600160401b03811161084857602091010185601f820112156108485780516001600160401b03811161084d57604051916107f9601f8301601f19166020018461087e565b81835287602083830101116108485760005b8281106108335750509181600060208096949581960101526060820152815201940193610182565b8060208092840101518282870101520161080b565b600080fd5b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761084d57604052565b601f909101601f19168101906001600160401b0382119082101761084d57604052565b51906001600160401b038216820361084857565b51906001600160a01b038216820361084857565b5190811515820361084857565b90600182811c92168015610906575b60208310146108f057565b634e487b7160e01b600052602260045260246000fd5b91607f16916108e5565b60009291815491610920836108d6565b8083529260018116908115610976575060011461093c57505050565b60009081526020812093945091925b83831061095c575060209250010190565b60018160209294939454838587010152019101919061094b565b915050602093945060ff929192191683830152151560051b010190565b80600052600760205260406000205415600014610a20576006546801000000000000000081101561084d576001810180600655811015610a0a577ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0181905560065460009182526007602052604090912055600190565b634e487b7160e01b600052603260045260246000fd5b5060009056fe6080604052600436101561001257600080fd5b60003560e01c806304666f9c1461015757806306285c6914610152578063181f5a771461014d5780633f4b04aa146101485780635215505b146101435780635e36480c1461013e5780635e7bb0081461013957806360987c20146101345780637437ff9f1461012f57806379ba50971461012a5780637edf52f41461012557806385572ffb146101205780638da5cb5b1461011b578063c673e58414610116578063ccd37ba314610111578063de5e0b9a1461010c578063e9d68a8e14610107578063f2fde38b14610102578063f58e03fc146100fd5763f716f99f146100f857600080fd5b6118d4565b6117b7565b61172c565b611691565b6115f5565b611571565b6114c6565b6113de565b6113a8565b6111e2565b611162565b6110b9565b61103e565b610e39565b6108ce565b610789565b61067c565b61061d565b61043d565b61031f565b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761018d57604052565b61015c565b60a081019081106001600160401b0382111761018d57604052565b604081019081106001600160401b0382111761018d57604052565b606081019081106001600160401b0382111761018d57604052565b90601f801991011681019081106001600160401b0382111761018d57604052565b6040519061021360c0836101e3565b565b6040519061021360a0836101e3565b60405190610213610100836101e3565b604051906102136040836101e3565b6001600160401b03811161018d5760051b60200190565b6001600160a01b0381160361026b57565b600080fd5b600435906001600160401b038216820361026b57565b35906001600160401b038216820361026b57565b8015150361026b57565b35906102138261029a565b6001600160401b03811161018d57601f01601f191660200190565b9291926102d6826102af565b916102e460405193846101e3565b82948184528183011161026b578281602093846000960137010152565b9080601f8301121561026b5781602061031c933591016102ca565b90565b3461026b57602036600319011261026b576004356001600160401b03811161026b573660238201121561026b5780600401359061035b82610243565b9061036960405192836101e3565b8282526024602083019360051b8201019036821161026b5760248101935b8285106103995761039784611a0f565b005b84356001600160401b03811161026b5782016080602319823603011261026b57604051916103c683610172565b60248201356103d48161025a565b83526103e260448301610286565b602084015260648201356103f58161029a565b60408401526084820135926001600160401b03841161026b57610422602094936024869536920101610301565b6060820152815201940193610387565b600091031261026b57565b3461026b57600036600319011261026b576000608060405161045e81610192565b82815282602082015282604082015282606082015201526105bc60405161048481610192565b6001600160401b037f000000000000000000000000000000000000000000000000000000000000000016815261ffff7f00000000000000000000000000000000000000000000000000000000000000001660208201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660408201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660608201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660808201526040519182918291909160806001600160a01b038160a08401956001600160401b03815116855261ffff6020820151166020860152826040820151166040860152826060820151166060860152015116910152565b0390f35b604051906105cf6020836101e3565b60008252565b60005b8381106105e85750506000910152565b81810151838201526020016105d8565b90602091610611815180928185528580860191016105d5565b601f01601f1916010190565b3461026b57600036600319011261026b576105bc604080519061064081836101e3565b601182527f4f666652616d7020312e362e302d6465760000000000000000000000000000006020830152519182916020835260208301906105f8565b3461026b57600036600319011261026b5760206001600160401b03600b5416604051908152f35b906080606061031c936001600160a01b0381511684526020810151151560208501526001600160401b03604082015116604085015201519181606082015201906105f8565b6040810160408252825180915260206060830193019060005b81811061076a575050506020818303910152815180825260208201916020808360051b8301019401926000915b83831061073d57505050505090565b909192939460208061075b600193601f1986820301875289516106a3565b9701930193019193929061072e565b82516001600160401b0316855260209485019490920191600101610701565b3461026b57600036600319011261026b576006546107a681610243565b906107b460405192836101e3565b808252601f196107c382610243565b0160005b8181106108855750506107d981611ce2565b9060005b8181106107f55750506105bc604051928392836106e8565b8061082b61081361080760019461406b565b6001600160401b031690565b61081d8387611d3c565b906001600160401b03169052565b61086961086461084b61083e8488611d3c565b516001600160401b031690565b6001600160401b03166000526008602052604060002090565b611e28565b6108738287611d3c565b5261087e8186611d3c565b50016107dd565b602090610890611cbb565b828287010152016107c7565b634e487b7160e01b600052602160045260246000fd5b600411156108bc57565b61089c565b9060048210156108bc5752565b3461026b57604036600319011261026b576108e7610270565b602435906001600160401b038216820361026b5760209161090791611ec4565b61091460405180926108c1565bf35b91908260a091031261026b5760405161092e81610192565b60806109738183958035855261094660208201610286565b602086015261095760408201610286565b604086015261096860608201610286565b606086015201610286565b910152565b35906102138261025a565b63ffffffff81160361026b57565b359061021382610983565b81601f8201121561026b578035906109b382610243565b926109c160405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b8385106109ed57505050505090565b84356001600160401b03811161026b57820160a0818503601f19011261026b5760405191610a1a83610192565b60208201356001600160401b03811161026b57856020610a3c92850101610301565b83526040820135610a4c8161025a565b6020840152610a5d60608301610991565b60408401526080820135926001600160401b03841161026b5760a083610a8a886020809881980101610301565b6060840152013560808201528152019401936109de565b9190916101408184031261026b57610ab7610204565b92610ac28183610916565b845260a08201356001600160401b03811161026b5781610ae3918401610301565b602085015260c08201356001600160401b03811161026b5781610b07918401610301565b6040850152610b1860e08301610978565b606085015261010082013560808501526101208201356001600160401b03811161026b57610b46920161099c565b60a0830152565b9080601f8301121561026b578135610b6481610243565b92610b7260405194856101e3565b81845260208085019260051b8201019183831161026b5760208201905b838210610b9e57505050505090565b81356001600160401b03811161026b57602091610bc087848094880101610aa1565b815201910190610b8f565b81601f8201121561026b57803590610be282610243565b92610bf060405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b838510610c1c57505050505090565b84356001600160401b03811161026b57820183603f8201121561026b576020810135610c4781610243565b91610c5560405193846101e3565b8183526020808085019360051b830101019186831161026b5760408201905b838210610c8e575050509082525060209485019401610c0d565b81356001600160401b03811161026b57602091610cb28a8480809589010101610301565b815201910190610c74565b929190610cc981610243565b93610cd760405195866101e3565b602085838152019160051b810192831161026b57905b828210610cf957505050565b8135815260209182019101610ced565b9080601f8301121561026b5781602061031c93359101610cbd565b81601f8201121561026b57803590610d3b82610243565b92610d4960405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b838510610d7557505050505090565b84356001600160401b03811161026b57820160a0818503601f19011261026b57610d9d610215565b91610daa60208301610286565b835260408201356001600160401b03811161026b57856020610dce92850101610b4d565b602084015260608201356001600160401b03811161026b57856020610df592850101610bcb565b60408401526080820135926001600160401b03841161026b5760a083610e22886020809881980101610d09565b606084015201356080820152815201940193610d66565b3461026b57604036600319011261026b576004356001600160401b03811161026b57610e69903690600401610d24565b6024356001600160401b03811161026b573660238201121561026b57806004013591610e9483610243565b91610ea260405193846101e3565b8383526024602084019460051b8201019036821161026b5760248101945b828610610ed1576103978585611f0c565b85356001600160401b03811161026b5782013660438201121561026b576024810135610efc81610243565b91610f0a60405193846101e3565b818352602060248185019360051b830101019036821161026b5760448101925b828410610f44575050509082525060209586019501610ec0565b83356001600160401b03811161026b576024908301016040601f19823603011261026b5760405190610f75826101ad565b6020810135825260408101356001600160401b03811161026b57602091010136601f8201121561026b57803590610fab82610243565b91610fb960405193846101e3565b80835260208084019160051b8301019136831161026b57602001905b828210610ff45750505091816020938480940152815201930192610f2a565b60208091833561100381610983565b815201910190610fd5565b9181601f8401121561026b578235916001600160401b03831161026b576020808501948460051b01011161026b57565b3461026b57606036600319011261026b576004356001600160401b03811161026b5761106e903690600401610aa1565b6024356001600160401b03811161026b5761108d90369060040161100e565b91604435926001600160401b03841161026b576110b161039794369060040161100e565b939092612318565b3461026b57600036600319011261026b576110d26125e5565b506105bc6040516110e281610172565b60ff6004546001600160a01b038116835263ffffffff8160a01c16602084015260c01c16151560408201526001600160a01b036005541660608201526040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b3461026b57600036600319011261026b576000546001600160a01b03811633036111d1576001600160a01b0319600154913382841617600155166000556001600160a01b033391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b63015aa1e360e11b60005260046000fd5b3461026b57608036600319011261026b57600060405161120181610172565b60043561120d8161025a565b815260243561121b81610983565b602082015260443561122c8161029a565b604082015260643561123d8161025a565b606082015261124a613489565b6001600160a01b038151161561139957611393816112a96001600160a01b037fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9451166001600160a01b03166001600160a01b03196004541617600455565b60208101516004547fffffffffffffff0000000000ffffffffffffffffffffffffffffffffffffffff77ffffffff000000000000000000000000000000000000000078ff0000000000000000000000000000000000000000000000006040860151151560c01b169360a01b169116171760045561134f61133360608301516001600160a01b031690565b6001600160a01b03166001600160a01b03196005541617600555565b6040519182918291909160606001600160a01b0381608084019582815116855263ffffffff6020820151166020860152604081015115156040860152015116910152565b0390a180f35b6342bcdf7f60e11b8252600482fd5b3461026b57602036600319011261026b576004356001600160401b03811161026b5760a090600319903603011261026b57600080fd5b3461026b57600036600319011261026b5760206001600160a01b0360015416604051908152f35b6004359060ff8216820361026b57565b359060ff8216820361026b57565b906020808351928381520192019060005b8181106114415750505090565b82516001600160a01b0316845260209384019390920191600101611434565b9061031c9160208152606082518051602084015260ff602082015116604084015260ff6040820151168284015201511515608082015260406114b1602084015160c060a085015260e0840190611423565b9201519060c0601f1982850301910152611423565b3461026b57602036600319011261026b5760ff6114e1611405565b6060604080516114f0816101c8565b6114f86125e5565b815282602082015201521660005260026020526105bc6040600020600361156060405192611525846101c8565b61152e8161260a565b845260405161154b816115448160028601612643565b03826101e3565b60208501526115446040518094819301612643565b604082015260405191829182611460565b3461026b57604036600319011261026b5761158a610270565b6001600160401b036024359116600052600a6020526040600020906000526020526020604060002054604051908152f35b9060049160441161026b57565b9181601f8401121561026b578235916001600160401b03831161026b576020838186019501011161026b57565b3461026b5760c036600319011261026b5761160f366115bb565b6044356001600160401b03811161026b5761162e9036906004016115c8565b6064929192356001600160401b03811161026b5761165090369060040161100e565b60843594916001600160401b03861161026b5761167461039796369060040161100e565b94909360a43596612c75565b90602061031c9281815201906106a3565b3461026b57602036600319011261026b576001600160401b036116b2610270565b6116ba611cbb565b501660005260086020526105bc6040600020600161171b604051926116de84610172565b6001600160401b0381546001600160a01b038116865260ff8160a01c161515602087015260a81c1660408501526115446040518094819301611d8a565b606082015260405191829182611680565b3461026b57602036600319011261026b576001600160a01b036004356117518161025a565b611759613489565b163381146117a657806001600160a01b031960005416176000556001600160a01b03600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b636d6c4ee560e11b60005260046000fd5b3461026b57606036600319011261026b576117d1366115bb565b6044356001600160401b03811161026b576117f09036906004016115c8565b9182820160208382031261026b578235906001600160401b03821161026b5761181a918401610d24565b60405190602061182a81846101e3565b60008352601f19810160005b81811061185e57505050610397949161184e916136cd565b611856613184565b928392613a33565b60608582018401528201611836565b9080601f8301121561026b57813561188481610243565b9261189260405194856101e3565b81845260208085019260051b82010192831161026b57602001905b8282106118ba5750505090565b6020809183356118c98161025a565b8152019101906118ad565b3461026b57602036600319011261026b576004356001600160401b03811161026b573660238201121561026b5780600401359061191082610243565b9061191e60405192836101e3565b8282526024602083019360051b8201019036821161026b5760248101935b82851061194c57610397846131a0565b84356001600160401b03811161026b57820160c0602319823603011261026b57611974610204565b916024820135835261198860448301611415565b602084015261199960648301611415565b60408401526119aa608483016102a4565b606084015260a48201356001600160401b03811161026b576119d2906024369185010161186d565b608084015260c4820135926001600160401b03841161026b576119ff60209493602486953692010161186d565b60a082015281520194019361193c565b611a17613489565b60005b8151811015611cb757611a2d8183611d3c565b5190611a4360208301516001600160401b031690565b916001600160401b038316908115611ca657611a78611a6c611a6c83516001600160a01b031690565b6001600160a01b031690565b15611c0d57611a9a846001600160401b03166000526008602052604060002090565b906060810151916001810195611ab08754611d50565b611c3457611b237ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb991611b0984750100000000000000000000000000000000000000000067ffffffffffffffff60a81b19825416179055565b6040516001600160401b0390911681529081906020820190565b0390a15b82518015908115611c1e575b50611c0d57611bee611bd2611c0493611b6f7f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b9660019a61352b565b611bc5611b7f6040830151151590565b85547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690151560a01b74ff000000000000000000000000000000000000000016178555565b516001600160a01b031690565b82906001600160a01b03166001600160a01b0319825416179055565b611bf784615093565b50604051918291826135fc565b0390a201611a1a565b6342bcdf7f60e11b60005260046000fd5b90506020840120611c2d6134ae565b1438611b33565b60016001600160401b03611c5384546001600160401b039060a81c1690565b16141580611c87575b611c665750611b27565b632105803760e11b6000526001600160401b031660045260246000fd5b6000fd5b50611c9187611e0d565b60208151910120845160208601201415611c5c565b63c656089560e01b60005260046000fd5b5050565b60405190611cc882610172565b606080836000815260006020820152600060408201520152565b90611cec82610243565b611cf960405191826101e3565b8281528092611d0a601f1991610243565b0190602036910137565b634e487b7160e01b600052603260045260246000fd5b805115611d375760200190565b611d14565b8051821015611d375760209160051b010190565b90600182811c92168015611d80575b6020831014611d6a57565b634e487b7160e01b600052602260045260246000fd5b91607f1691611d5f565b60009291815491611d9a83611d50565b8083529260018116908115611df05750600114611db657505050565b60009081526020812093945091925b838310611dd6575060209250010190565b600181602092949394548385870101520191019190611dc5565b915050602093945060ff929192191683830152151560051b010190565b90610213611e219260405193848092611d8a565b03836101e3565b9060016060604051611e3981610172565b611e8281956001600160401b0381546001600160a01b038116855260ff8160a01c161515602086015260a81c166040840152611e7b6040518096819301611d8a565b03846101e3565b0152565b634e487b7160e01b600052601160045260246000fd5b908160051b9180830460201490151715611eb257565b611e86565b91908203918211611eb257565b611ed082607f92613646565b9116906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611eb2576003911c1660048110156108bc5790565b611f1461368a565b80518251810361210b5760005b818110611f3457505090610213916136cd565b611f3e8184611d3c565b516020810190815151611f518488611d3c565b51928351820361210b5790916000925b808410611f75575050505050600101611f21565b91949398611f87848b98939598611d3c565b515198611f95888851611d3c565b5199806120c2575b5060a08a01988b6020611fb38b8d515193611d3c565b51015151036120855760005b8a515181101561207057611ffb611ff2611fe88f6020611fe08f8793611d3c565b510151611d3c565b5163ffffffff1690565b63ffffffff1690565b8b8161200c575b5050600101611fbf565b611ff2604061201f8561202b9451611d3c565b51015163ffffffff1690565b9081811061203a57508b612002565b8d51516040516348e617b360e01b81526004810191909152602481019390935260448301919091526064820152608490fd5b0390fd5b50985098509893949095600101929091611f61565b611c838b516120a0606082519201516001600160401b031690565b6370a193fd60e01b6000526004919091526001600160401b0316602452604490565b60808b0151811015611f9d57611c83908b6120e488516001600160401b031690565b905151633a98d46360e11b6000526001600160401b03909116600452602452604452606490565b6320f8fd5960e21b60005260046000fd5b60405190612129826101ad565b60006020838281520152565b604051906121446020836101e3565b600080835282815b82811061215857505050565b60209061216361211c565b8282850101520161214c565b805182526001600160401b03602082015116602083015260806121b66121a4604084015160a0604087015260a08601906105f8565b606084015185820360608701526105f8565b9101519160808183039101526020808351928381520192019060005b8181106121df5750505090565b825180516001600160a01b0316855260209081015181860152604090940193909201916001016121d2565b90602061031c92818152019061216f565b6040513d6000823e3d90fd5b3d15612252573d90612238826102af565b9161224660405193846101e3565b82523d6000602084013e565b606090565b90602061031c9281815201906105f8565b909160608284031261026b57815161227f8161029a565b9260208301516001600160401b03811161026b5783019080601f8301121561026b578151916122ad836102af565b916122bb60405193846101e3565b8383526020848301011161026b576040926122dc91602080850191016105d5565b92015190565b9293606092959461ffff6123066001600160a01b039460808852608088019061216f565b97166020860152604085015216910152565b929093913033036125d45761232b612135565b9460a0850151805161258d575b5050505050805191612356602084519401516001600160401b031690565b906020830151916040840192612383845192612370610215565b9788526001600160401b03166020880152565b6040860152606085015260808401526001600160a01b036123ac6005546001600160a01b031690565b1680612510575b5051511580612504575b80156124ee575b80156124c5575b611cb75761245d9181612402611a6c6123f561084b602060009751016001600160401b0390511690565b546001600160a01b031690565b908361241d606060808401519301516001600160a01b031690565b604051633cf9798360e01b815296879586948593917f000000000000000000000000000000000000000000000000000000000000000090600486016122e2565b03925af19081156124c057600090600092612499575b501561247c5750565b6040516302a35ba360e21b815290819061206c9060048301612257565b90506124b891503d806000833e6124b081836101e3565b810190612268565b509038612473565b61221b565b506124e96124e56124e060608401516001600160a01b031690565b6138f4565b1590565b6123cb565b5060608101516001600160a01b03163b156123c4565b506080810151156123bd565b803b1561026b57600060405180926308d450a160e01b82528183816125388a6004830161220a565b03925af19081612572575b5061256c5761206c612553612227565b6040516309c2532560e01b815291829160048301612257565b386123b3565b806125816000612587936101e3565b80610432565b38612543565b85965060206125c99601516125ac60608901516001600160a01b031690565b906125c360208a51016001600160401b0390511690565b926137db565b903880808080612338565b6306e34e6560e31b60005260046000fd5b604051906125f282610172565b60006060838281528260208201528260408201520152565b9060405161261781610172565b606060ff600183958054855201548181166020850152818160081c16604085015260101c161515910152565b906020825491828152019160005260206000209060005b8181106126675750505090565b82546001600160a01b031684526020909301926001928301920161265a565b90610213611e219260405193848092612643565b35906001600160e01b038216820361026b57565b81601f8201121561026b578035906126c582610243565b926126d360405194856101e3565b82845260208085019360061b8301019181831161026b57602001925b8284106126fd575050505090565b60408483031261026b5760206040918251612717816101ad565b61272087610286565b815261272d83880161269a565b838201528152019301926126ef565b81601f8201121561026b5780359061275382610243565b9261276160405194856101e3565b82845260208085019360051b8301019181831161026b5760208101935b83851061278d57505050505090565b84356001600160401b03811161026b57820160a0818503601f19011261026b57604051916127ba83610192565b6127c660208301610286565b83526040820135926001600160401b03841161026b5760a0836127f0886020809881980101610301565b8584015261280060608201610286565b604084015261281160808201610286565b60608401520135608082015281520194019361277e565b81601f8201121561026b5780359061283f82610243565b9261284d60405194856101e3565b82845260208085019360061b8301019181831161026b57602001925b828410612877575050505090565b60408483031261026b5760206040918251612891816101ad565b863581528287013583820152815201930192612869565b60208183031261026b578035906001600160401b03821161026b570160608183031261026b57604051916128db836101c8565b81356001600160401b03811161026b57820160408183031261026b5760405190612904826101ad565b80356001600160401b03811161026b57810183601f8201121561026b57803561292c81610243565b9161293a60405193846101e3565b81835260208084019260061b8201019086821161026b57602001915b8183106129d25750505082526020810135906001600160401b03821161026b57612982918491016126ae565b6020820152835260208201356001600160401b03811161026b57816129a891840161273c565b602084015260408201356001600160401b03811161026b576129ca9201612828565b604082015290565b60408388031261026b57602060409182516129ec816101ad565b85356129f78161025a565b8152612a0483870161269a565b83820152815201920191612956565b9080602083519182815201916020808360051b8301019401926000915b838310612a3f57505050505090565b9091929394602080600192601f198582030186528851906001600160401b038251168152608080612a7d8585015160a08786015260a08501906105f8565b936001600160401b0360408201511660408501526001600160401b036060820151166060850152015191015297019301930191939290612a30565b916001600160a01b03612ad992168352606060208401526060830190612a13565b9060408183039101526020808351928381520192019060005b818110612aff5750505090565b8251805185526020908101518186015260409094019390920191600101612af2565b906020808351928381520192019060005b818110612b3f5750505090565b825180516001600160401b031685526020908101516001600160e01b03168186015260409094019390920191600101612b32565b9190604081019083519160408252825180915260206060830193019060005b818110612bb357505050602061031c93940151906020818403910152612b21565b825180516001600160a01b031686526020908101516001600160e01b03168187015260409095019490920191600101612b92565b90602061031c928181520190612b73565b9081602091031261026b575161031c8161029a565b9091612c2461031c936040845260408401906105f8565b916020818403910152611d8a565b6001600160401b036001911601906001600160401b038211611eb257565b9091612c6761031c93604084526040840190612a13565b916020818403910152612b73565b929693959190979497612c8a828201826128a8565b98612c9e6124e560045460ff9060c01c1690565b6130f2575b895180515115908115916130e3575b5061300a575b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316999860208a019860005b8a518051821015612fa85781612d0191611d3c565b518d612d1482516001600160401b031690565b604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b1660048201529091602090829060249082905afa9081156124c057600091612f7a575b50612f5d57612d6390613942565b60208201805160208151910120906001830191612d7f83611e0d565b6020815191012003612f40575050805460408301516001600160401b039081169160a81c168114801590612f18575b612ec657506080820151908115612eb557612dff82612df0612dd786516001600160401b031690565b6001600160401b0316600052600a602052604060002090565b90600052602052604060002090565b54612e81578291612e65612e7a92612e2c612e2760606001999801516001600160401b031690565b612c32565b67ffffffffffffffff60a81b197cffffffffffffffff00000000000000000000000000000000000000000083549260a81b169116179055565b612df0612dd74294516001600160401b031690565b5501612cec565b50612e96611c8392516001600160401b031690565b6332cf0cbf60e01b6000526001600160401b0316600452602452604490565b63504570e360e01b60005260046000fd5b82611c8391612ef06060612ee184516001600160401b031690565b9301516001600160401b031690565b636af0786b60e11b6000526001600160401b0392831660045290821660245216604452606490565b50612f3061080760608501516001600160401b031690565b6001600160401b03821611612dae565b5161206c60405192839263b80d8fa960e01b845260048401612c0d565b637edeb53960e11b6000526001600160401b031660045260246000fd5b612f9b915060203d8111612fa1575b612f9381836101e3565b810190612bf8565b38612d55565b503d612f89565b50506130049496989b507f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e46102139b612ffc949597999b51905190612ff260405192839283612c50565b0390a13691610cbd565b943691610cbd565b93613d2d565b61301f602086015b356001600160401b031690565b600b546001600160401b03828116911610156130c757613055906001600160401b03166001600160401b0319600b541617600b55565b61306d611a6c611a6c6004546001600160a01b031690565b8a5190803b1561026b57604051633937306f60e01b815291600091839182908490829061309d9060048301612be7565b03925af180156124c0576130b2575b50612cb8565b8061258160006130c1936101e3565b386130ac565b5060208a015151612cb857632261116760e01b60005260046000fd5b60209150015151151538612cb2565b60208a01518051613104575b50612ca3565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169060408c0151823b1561026b57604051633854844f60e11b815292600092849283918291613160913060048501612ab8565b03915afa80156124c057156130fe5780612581600061317e936101e3565b386130fe565b604051906131936020836101e3565b6000808352366020840137565b6131a8613489565b60005b8151811015611cb7576131be8183611d3c565b51906040820160ff6131d1825160ff1690565b161561347357602083015160ff16926131f78460ff166000526002602052604060002090565b916001830191825461321261320c8260ff1690565b60ff1690565b613438575061323f6132276060830151151590565b845462ff0000191690151560101b62ff000016178455565b60a081019182516101008151116133e057805115613422576003860161326d61326782612686565b8a614e41565b60608401516132fd575b947fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547946002946132d96132c96132f79a966132c28760019f9c6132bd6132ef9a8f614fa2565b613f6e565b5160ff1690565b845460ff191660ff821617909455565b5190818555519060405195869501908886613ff4565b0390a1615024565b016131ab565b9794600287939597019661331961331389612686565b88614e41565b60808501519461010086511161340c57855161334161320c61333c8a5160ff1690565b613f5a565b10156133f65785518451116133e0576132d96132c97fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547986132c28760019f6132bd6132f79f9a8f6133c860029f6133c26132ef9f8f906132bd84926133a7845160ff1690565b908054909161ff001990911660089190911b61ff0016179055565b82614ed5565b505050979c9f50975050969a50505094509450613277565b631b3fab5160e11b600052600160045260246000fd5b631b3fab5160e11b600052600360045260246000fd5b631b3fab5160e11b600052600260045260246000fd5b631b3fab5160e11b600052600560045260246000fd5b60101c60ff1661345361344e6060840151151590565b151590565b9015151461323f576321fd80df60e21b60005260ff861660045260246000fd5b631b3fab5160e11b600090815260045260246000fd5b6001600160a01b0360015416330361349d57565b6315ae3a6f60e11b60005260046000fd5b604051602081019060008252602081526134c96040826101e3565b51902090565b8181106134da575050565b600081556001016134cf565b9190601f81116134f557505050565b610213926000526020600020906020601f840160051c83019310613521575b601f0160051c01906134cf565b9091508190613514565b91909182516001600160401b03811161018d576135528161354c8454611d50565b846134e6565b6020601f8211600114613593578190613584939495600092613588575b50508160011b916000199060031b1c19161790565b9055565b01519050388061356f565b601f198216906135a884600052602060002090565b9160005b8181106135e4575095836001959697106135cb575b505050811b019055565b015160001960f88460031b161c191690553880806135c1565b9192602060018192868b0151815501940192016135ac565b90600160a061031c93602081526001600160401b0384546001600160a01b038116602084015260ff81851c161515604084015260a81c166060820152608080820152019101611d8a565b906001600160401b03613686921660005260096020526701ffffffffffffff60406000209160071c166001600160401b0316600052602052604060002090565b5490565b7f00000000000000000000000000000000000000000000000000000000000000004681036136b55750565b630f01ce8560e01b6000526004524660245260446000fd5b91909180511561376f5782511592602091604051926136ec81856101e3565b60008452601f19810160005b81811061374b5750505060005b8151811015613743578061372c61371e60019385611d3c565b518815613732578690614133565b01613705565b61373c8387611d3c565b5190614133565b505050509050565b8290604051613759816101ad565b60008152606083820152828289010152016136f8565b63c2e5347d60e01b60005260046000fd5b9190811015611d375760051b0190565b3561031c81610983565b9190811015611d375760051b81013590601e198136030182121561026b5701908135916001600160401b03831161026b57602001823603811361026b579190565b909294919397968151966137ee88610243565b976137fc604051998a6101e3565b80895261380b601f1991610243565b0160005b8181106138dd57505060005b83518110156138d057806138628c8a8a8a61385c613855878d61384e828f8f9d8f9e60019f8161387e575b505050611d3c565b519761379a565b36916102ca565b93614978565b61386c828c611d3c565b52613877818b611d3c565b500161381b565b63ffffffff613896613891858585613780565b613790565b1615613846576138c6926138ad9261389192613780565b60406138b98585611d3c565b51019063ffffffff169052565b8f8f908391613846565b5096985050505050505050565b6020906138e861211c565b82828d0101520161380f565b6139056385572ffb60e01b82614cdb565b908161391f575b81613915575090565b61031c9150614cad565b905061392a81614c32565b159061390c565b61390563aff2afbf60e01b82614cdb565b6001600160401b031680600052600860205260406000209060ff825460a01c161561396b575090565b63ed053c5960e01b60005260045260246000fd5b6084019081608411611eb257565b60a001908160a011611eb257565b91908201809211611eb257565b600311156108bc57565b60038210156108bc5752565b906102136040516139ce816101ad565b602060ff829554818116845260081c1691016139b2565b8054821015611d375760005260206000200190600090565b60ff60019116019060ff8211611eb257565b60ff601b9116019060ff8211611eb257565b90606092604091835260208301370190565b6001600052600260205293613a677fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e061260a565b93853594613a748561397f565b6060820190613a838251151590565b613cff575b803603613ce757508151878103613cce5750613aa261368a565b60016000526003602052613af1613aec7fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c5b336001600160a01b0316600052602052604060002090565b6139be565b60026020820151613b01816139a8565b613b0a816139a8565b149081613c66575b5015613c3a575b51613b71575b50505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef090613b5561301260019460200190565b604080519283526001600160401b0391909116602083015290a2565b613b9261320c613b8d602085969799989a955194015160ff1690565b6139fd565b03613c29578151835103613c1857613c106000613b559461301294613bdc7f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef09960019b36916102ca565b60208151910120604051613c0781613bf989602083019586613a21565b03601f1981018352826101e3565b5190208a614d0b565b948394613b1f565b63a75d88af60e01b60005260046000fd5b6371253a2560e01b60005260046000fd5b72c11c11c11c11c11c11c11c11c11c11c11c11c1330315613b1957631b41e11d60e31b60005260046000fd5b60016000526002602052613cc69150611a6c90613cb390613cad60037fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e05b01915160ff1690565b906139e5565b90546001600160a01b039160031b1c1690565b331438613b12565b6324f7d61360e21b600052600452602487905260446000fd5b638e1192e160e01b6000526004523660245260446000fd5b613d2890613d22613d18613d138751611e9c565b61398d565b613d228851611e9c565b9061399b565b613a88565b60008052600260205294909390929091613d667fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b61260a565b94863595613d738361397f565b6060820190613d828251151590565b613f37575b803603613ce757508151888103613f1e5750613da161368a565b600080526003602052613dd6613aec7f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff613ad4565b60026020820151613de6816139a8565b613def816139a8565b149081613ed5575b5015613ea9575b51613e3b575b5050505050507f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef090613b5561301260009460200190565b613e5761320c613b8d602087989a999b96975194015160ff1690565b03613c29578351865103613c18576000967f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef096613b5595613bdc613ea0946130129736916102ca565b94839438613e04565b72c11c11c11c11c11c11c11c11c11c11c11c11c1330315613dfe57631b41e11d60e31b60005260046000fd5b600080526002602052613f169150611a6c90613cb390613cad60037fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b613ca4565b331438613df7565b6324f7d61360e21b600052600452602488905260446000fd5b613f5590613d22613f4b613d138951611e9c565b613d228a51611e9c565b613d87565b60ff166003029060ff8216918203611eb257565b8151916001600160401b03831161018d5768010000000000000000831161018d576020908254848455808510613fd7575b500190600052602060002060005b838110613fba5750505050565b60019060206001600160a01b038551169401938184015501613fad565b613fee9084600052858460002091820191016134cf565b38613f9f565b95949392909160ff61401993168752602087015260a0604087015260a0860190612643565b84810360608601526020808351928381520192019060005b81811061404c575050509060806102139294019060ff169052565b82516001600160a01b0316845260209384019390920191600101614031565b600654811015611d375760066000527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f015490565b6001600160401b0361031c94938160609416835216602082015281604082015201906105f8565b60409061031c9392815281602082015201906105f8565b9291906001600160401b039081606495166004521660245260048110156108bc57604452565b94939261411d60609361412e93885260208801906108c1565b6080604087015260808601906105f8565b930152565b9061414582516001600160401b031690565b8151604051632cbc26bb60e01b815267ffffffffffffffff60801b608084901b1660048201529015159391906001600160401b038216906020816024817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9081156124c057600091614861575b5061481f5760208301918251519485156147ef576040850180515187036147de576141e787611ce2565b957f000000000000000000000000000000000000000000000000000000000000000061421d600161421787613942565b01611e0d565b6020815191012060405161427d81613bf96020820194868b876001600160401b036060929594938160808401977f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f85521660208401521660408201520152565b519020906001600160401b031660005b8a81106147465750505080608060606142ad9301519101519088866152d4565b9788156147285760005b8881106142ca5750505050505050505050565b5a6142d6828951611d3c565b518051606001516142f0906001600160401b031688611ec4565b6142f9816108b2565b8015908d8283159384614715575b156146d25760608815614655575061432e6020614324898d611d3c565b5101519242611eb7565b6004546143439060a01c63ffffffff16611ff2565b108015614642575b156146245761435a878b611d3c565b515161460e575b845160800151614379906001600160401b0316610807565b614556575b5061438a868951611d3c565b5160a08501515181510361451a57936143ef9695938c938f966143cf8e958c926143c96143c360608951016001600160401b0390511690565b89615306565b866155c1565b9a9080966143e960608851016001600160401b0390511690565b9061538e565b6144c8575b50506143ff826108b2565b60028203614480575b6001966144767f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b936001600160401b0393519261446761445e8b61445660608801516001600160401b031690565b96519b611d3c565b51985a90611eb7565b91604051958695169885614104565b0390a45b016142b7565b91509193949250614490826108b2565b600382036144a4578b929493918a91614408565b51606001516349362d1f60e11b600052611c8391906001600160401b0316896140de565b6144d1846108b2565b600384036143f45790929495506144e99193506108b2565b6144f9578b92918a9138806143f4565b5151604051632b11b8d960e01b815290819061206c908790600484016140c7565b611c838b61453460608851016001600160401b0390511690565b631cfe6d8b60e01b6000526001600160401b0391821660045216602452604490565b61455f836108b2565b61456a575b3861437e565b8351608001516001600160401b0316602080860151918c61459f60405194859384936370701e5760e11b8552600485016140a0565b038160006001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165af19081156124c0576000916145f0575b5061456457505050505060019061447a565b614608915060203d8111612fa157612f9381836101e3565b386145de565b614618878b611d3c565b51516080860152614361565b6354e7e43160e11b6000526001600160401b038b1660045260246000fd5b5061464c836108b2565b6003831461434b565b915083614661846108b2565b15614361575060019594506146ca92506146a891507f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe651209351016001600160401b0390511690565b604080516001600160401b03808c168252909216602083015290918291820190565b0390a161447a565b5050505060019291506146ca6146a860607f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c9351016001600160401b0390511690565b5061471f836108b2565b60038314614307565b633ee8bd3f60e11b6000526001600160401b03841660045260246000fd5b614751818a51611d3c565b518051604001516001600160401b03168381036147c157508051602001516001600160401b031689810361479e57509061478d846001936151cc565b614797828d611d3c565b520161428d565b636c95f1eb60e01b6000526001600160401b03808a166004521660245260446000fd5b631c21951160e11b6000526001600160401b031660045260246000fd5b6357e0e08360e01b60005260046000fd5b611c8361480386516001600160401b031690565b63676cf24b60e11b6000526001600160401b0316600452602490565b5092915050612f5d576040516001600160401b039190911681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d89493390602090a1565b61487a915060203d602011612fa157612f9381836101e3565b386141bd565b9081602091031261026b575161031c8161025a565b9061031c916020815260e061493361491e6148be855161010060208701526101208601906105f8565b60208601516001600160401b0316604086015260408601516001600160a01b0316606086015260608601516080860152614908608087015160a08701906001600160a01b03169052565b60a0860151858203601f190160c08701526105f8565b60c0850151848203601f1901848601526105f8565b92015190610100601f19828503019101526105f8565b6040906001600160a01b0361031c949316815281602082015201906105f8565b9081602091031261026b575190565b9193929361498461211c565b5060208301516001600160a01b031660405163bbe4f6db60e01b81526001600160a01b038216600482015290959092602084806024810103816001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000165afa9384156124c057600094614c01575b506001600160a01b0384169586158015614bef575b614bd157614ab6614adf92613bf992614a3a614a33611ff260408c015163ffffffff1690565b8c896156da565b9690996080810151614a686060835193015193614a55610224565b9687526001600160401b03166020870152565b6001600160a01b038a16604086015260608501526001600160a01b038d16608085015260a084015260c083015260e0820152604051633907753760e01b602082015292839160248301614895565b82857f000000000000000000000000000000000000000000000000000000000000000092615768565b94909115614bb55750805160208103614b9c575090614b08826020808a95518301019101614969565b956001600160a01b03841603614b40575b5050505050614b38614b29610234565b6001600160a01b039093168352565b602082015290565b614b5393614b4d91611eb7565b916156da565b50908082108015614b89575b614b6b57808481614b19565b63a966e21f60e01b6000908152600493909352602452604452606490fd5b5082614b958284611eb7565b1415614b5f565b631e3be00960e21b600052602060045260245260446000fd5b61206c604051928392634ff17cad60e11b845260048401614949565b63ae9b4ce960e01b6000526001600160a01b03851660045260246000fd5b50614bfc6124e586613931565b614a0d565b614c2491945060203d602011614c2b575b614c1c81836101e3565b810190614880565b92386149f8565b503d614c12565b60405160208101916301ffc9a760e01b835263ffffffff60e01b602483015260248252614c606044836101e3565b6179185a10614c9c576020926000925191617530fa6000513d82614c90575b5081614c89575090565b9050151590565b60201115915038614c7f565b63753fa58960e11b60005260046000fd5b60405160208101916301ffc9a760e01b83526301ffc9a760e01b602483015260248252614c606044836101e3565b6040519060208201926301ffc9a760e01b845263ffffffff60e01b16602483015260248252614c606044836101e3565b919390926000948051946000965b868810614d2a575050505050505050565b6020881015611d375760206000614d42878b1a613a0f565b614d4c8b87611d3c565b5190614d83614d5b8d8a611d3c565b5160405193849389859094939260ff6060936080840197845216602083015260408201520152565b838052039060015afa156124c057614dc9613aec600051614db18960ff166000526003602052604060002090565b906001600160a01b0316600052602052604060002090565b9060016020830151614dda816139a8565b614de3816139a8565b03614e3057614e00614df6835160ff1690565b60ff600191161b90565b8116614e1f57614e16614df66001935160ff1690565b17970196614d19565b633d9ef1f160e21b60005260046000fd5b636518c33d60e11b60005260046000fd5b91909160005b8351811015614e9a5760019060ff831660005260036020526000614e93604082206001600160a01b03614e7a858a611d3c565b51166001600160a01b0316600052602052604060002090565b5501614e47565b50509050565b8151815460ff191660ff91909116178155906020015160038110156108bc57815461ff00191660089190911b61ff0016179055565b919060005b8151811015614e9a57614ef0611bc58284611d3c565b90614f19614f0f83614db18860ff166000526003602052604060002090565b5460081c60ff1690565b614f22816139a8565b614f8d576001600160a01b03821615614f7c57614f76600192614f71614f46610234565b60ff8516815291614f5a86602085016139b2565b614db18960ff166000526003602052604060002090565b614ea0565b01614eda565b63d6c62c9b60e01b60005260046000fd5b631b3fab5160e11b6000526004805260246000fd5b919060005b8151811015614e9a57614fbd611bc58284611d3c565b90614fdc614f0f83614db18860ff166000526003602052604060002090565b614fe5816139a8565b614f8d576001600160a01b03821615614f7c5761501e600192614f71615009610234565b60ff8516815291614f5a6002602085016139b2565b01614fa7565b60ff1680600052600260205260ff60016040600020015460101c16908015600014615072575015615061576001600160401b0319600b5416600b55565b6317bd8dd160e11b60005260046000fd5b60011461507c5750565b61508257565b6307b8c74d60e51b60005260046000fd5b80600052600760205260406000205415600014615111576006546801000000000000000081101561018d57600181016006556000600654821015611d3757600690527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01819055600654906000526007602052604060002055600190565b50600090565b9080602083519182815201916020808360051b8301019401926000915b83831061514357505050505090565b9091929394602080600192601f198582030186528851906080806151a6615173855160a0865260a08601906105f8565b6001600160a01b0387870151168786015263ffffffff6040870151166040860152606086015185820360608701526105f8565b93015191015297019301930191939290615134565b90602061031c928181520190615117565b6134c981518051906152606151eb60608601516001600160a01b031690565b613bf961520260608501516001600160401b031690565b9361521b6080808a01519201516001600160401b031690565b90604051958694602086019889936001600160401b036080946001600160a01b0382959998949960a089019a8952166020880152166040860152606085015216910152565b519020613bf96020840151602081519101209360a060408201516020815191012091015160405161529981613bf96020820194856151bb565b51902090604051958694602086019889919260a093969594919660c08401976000855260208501526040840152606083015260808201520152565b926001600160401b03926152e792615825565b9116600052600a60205260406000209060005260205260406000205490565b607f8216906801fffffffffffffffe6001600160401b0383169260011b169180830460021490151715611eb25761538b916001600160401b036153498584613646565b921660005260096020526701ffffffffffffff60406000209460071c169160036001831b921b19161792906001600160401b0316600052602052604060002090565b55565b9091607f83166801fffffffffffffffe6001600160401b0382169160011b169080820460021490151715611eb2576153c68484613646565b60048310156108bc576001600160401b0361538b9416600052600960205260036701ffffffffffffff60406000209660071c1693831b921b19161792906001600160401b0316600052602052604060002090565b9080602083519182815201916020808360051b8301019401926000915b83831061544657505050505090565b9091929394602080615464600193601f1986820301875289516105f8565b97019301930191939290615437565b906020808351928381520192019060005b8181106154915750505090565b825163ffffffff16845260209384019390920191600101615484565b916155769061556861031c9593606086526001600160401b0360808251805160608a015282602082015116828a01528260408201511660a08a01528260608201511660c08a015201511660e087015260a061553461551d60208401516101406101008b01526101a08a01906105f8565b6040840151898203605f19016101208b01526105f8565b60608301516001600160a01b03166101408901529160808101516101608901520151868203605f1901610180880152615117565b90848203602086015261541a565b916040818403910152615473565b80516020909101516001600160e01b03198116929190600482106155a6575050565b6001600160e01b031960049290920360031b82901b16169150565b90303b1561026b576000916155ea6040519485938493630304c3e160e51b8552600485016154ad565b038183305af190816156c5575b506156ba57615604612227565b9072c11c11c11c11c11c11c11c11c11c11c11c11c13314615626575b60039190565b61563f61563283615584565b6001600160e01b03191690565b6337c3be2960e01b14801561569f575b8015615684575b1561562057611c8361566783615584565b632882569d60e01b6000526001600160e01b031916600452602490565b5061569161563283615584565b63753fa58960e11b14615656565b506156ac61563283615584565b632be8ca8b60e21b1461564f565b60029061031c6105c0565b8061258160006156d4936101e3565b386155f7565b6040516370a0823160e01b60208201526001600160a01b0390911660248201529192916157379061570e8160448101613bf9565b84837f000000000000000000000000000000000000000000000000000000000000000092615768565b92909115614bb55750805160208103614b9c5750906157628260208061031c95518301019101614969565b93611eb7565b93919361577560846102af565b9461578360405196876101e3565b6084865261579160846102af565b602087019590601f1901368737833b15615814575a90808210615803578291038060061c900311156157f2576000918291825a9560208451940192f1905a9003923d90608482116157e9575b6000908287523e929190565b608491506157dd565b6337c3be2960e01b60005260046000fd5b632be8ca8b60e21b60005260046000fd5b63030ed58f60e21b60005260046000fd5b80519282519084156159815761010185111580615975575b156158a4578185019460001986019561010087116158a45786156159655761586487611ce2565b9660009586978795885b8481106158c95750505050506001190180951493846158bf575b5050826158b5575b5050156158a4576158a091611d3c565b5190565b6309bde33960e01b60005260046000fd5b1490503880615890565b1492503880615888565b6001811b8281160361595757868a1015615942576158eb60018b019a85611d3c565b51905b8c888c101561592e575061590660018c019b86611d3c565b515b818d116158a457615927828f9261592190600196615992565b92611d3c565b520161586e565b60018d019c61593c91611d3c565b51615908565b61595060018c019b8d611d3c565b51906158ee565b615950600189019884611d3c565b5050505090506158a09150611d2a565b5061010182111561583d565b630469ac9960e21b60005260046000fd5b818110156159a4579061031c916159a9565b61031c915b906040519060208201926001845260408301526060820152606081526134c96080826101e356fea164736f6c634300081a000a49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b",
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

func (_OffRamp *OffRampCaller) GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error) {
	var out []interface{}
	err := _OffRamp.contract.Call(opts, &out, "getAllSourceChainConfigs")

	if err != nil {
		return *new([]uint64), *new([]OffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)
	out1 := *abi.ConvertType(out[1], new([]OffRampSourceChainConfig)).(*[]OffRampSourceChainConfig)

	return out0, out1, err

}

func (_OffRamp *OffRampSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetAllSourceChainConfigs(&_OffRamp.CallOpts)
}

func (_OffRamp *OffRampCallerSession) GetAllSourceChainConfigs() ([]uint64, []OffRampSourceChainConfig, error) {
	return _OffRamp.Contract.GetAllSourceChainConfigs(&_OffRamp.CallOpts)
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

func (_OffRamp *OffRampTransactor) Commit(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "commit", reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactorSession) Commit(reportContext [2][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Commit(&_OffRamp.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OffRamp *OffRampTransactor) Execute(opts *bind.TransactOpts, reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "execute", reportContext, report)
}

func (_OffRamp *OffRampSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Execute(&_OffRamp.TransactOpts, reportContext, report)
}

func (_OffRamp *OffRampTransactorSession) Execute(reportContext [2][32]byte, report []byte) (*types.Transaction, error) {
	return _OffRamp.Contract.Execute(&_OffRamp.TransactOpts, reportContext, report)
}

func (_OffRamp *OffRampTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampTransactorSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampTransactor) ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "manuallyExecute", reports, gasLimitOverrides)
}

func (_OffRamp *OffRampSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.Contract.ManuallyExecute(&_OffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRamp *OffRampTransactorSession) ManuallyExecute(reports []InternalExecutionReport, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.Contract.ManuallyExecute(&_OffRamp.TransactOpts, reports, gasLimitOverrides)
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
	MerkleRoots  []InternalMerkleRoot
	PriceUpdates InternalPriceUpdates
	Raw          types.Log
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
	MessageHash         [32]byte
	State               uint8
	ReturnData          []byte
	GasUsed             *big.Int
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

type OffRampSkippedReportExecutionIterator struct {
	Event *OffRampSkippedReportExecution

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampSkippedReportExecutionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampSkippedReportExecution)
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
		it.Event = new(OffRampSkippedReportExecution)
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

func (it *OffRampSkippedReportExecutionIterator) Error() error {
	return it.fail
}

func (it *OffRampSkippedReportExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampSkippedReportExecution struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_OffRamp *OffRampFilterer) FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampSkippedReportExecutionIterator, error) {

	logs, sub, err := _OffRamp.contract.FilterLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return &OffRampSkippedReportExecutionIterator{contract: _OffRamp.contract, event: "SkippedReportExecution", logs: logs, sub: sub}, nil
}

func (_OffRamp *OffRampFilterer) WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampSkippedReportExecution) (event.Subscription, error) {

	logs, sub, err := _OffRamp.contract.WatchLogs(opts, "SkippedReportExecution")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampSkippedReportExecution)
				if err := _OffRamp.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
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

func (_OffRamp *OffRampFilterer) ParseSkippedReportExecution(log types.Log) (*OffRampSkippedReportExecution, error) {
	event := new(OffRampSkippedReportExecution)
	if err := _OffRamp.contract.UnpackLog(event, "SkippedReportExecution", log); err != nil {
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
	case _OffRamp.abi.Events["SkippedReportExecution"].ID:
		return _OffRamp.ParseSkippedReportExecution(log)
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
	return common.HexToHash("0x35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4")
}

func (OffRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0xab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547")
}

func (OffRampDynamicConfigSet) Topic() common.Hash {
	return common.HexToHash("0xcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee")
}

func (OffRampExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b")
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

func (OffRampSkippedReportExecution) Topic() common.Hash {
	return common.HexToHash("0xaab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d894933")
}

func (OffRampSourceChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0x49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b")
}

func (OffRampSourceChainSelectorAdded) Topic() common.Hash {
	return common.HexToHash("0xf4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9")
}

func (OffRampStaticConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b7")
}

func (OffRampTransmitted) Topic() common.Hash {
	return common.HexToHash("0x198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0")
}

func (_OffRamp *OffRamp) Address() common.Address {
	return _OffRamp.address
}

type OffRampInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	GetAllSourceChainConfigs(opts *bind.CallOpts) ([]uint64, []OffRampSourceChainConfig, error)

	GetDynamicConfig(opts *bind.CallOpts) (OffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (uint8, error)

	GetLatestPriceSequenceNumber(opts *bind.CallOpts) (uint64, error)

	GetMerkleRoot(opts *bind.CallOpts, sourceChainSelector uint64, root [32]byte) (*big.Int, error)

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

	FilterSkippedReportExecution(opts *bind.FilterOpts) (*OffRampSkippedReportExecutionIterator, error)

	WatchSkippedReportExecution(opts *bind.WatchOpts, sink chan<- *OffRampSkippedReportExecution) (event.Subscription, error)

	ParseSkippedReportExecution(log types.Log) (*OffRampSkippedReportExecution, error)

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
