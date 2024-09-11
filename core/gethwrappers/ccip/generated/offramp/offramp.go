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

type OffRampDynamicConfig struct {
	FeeQuoter                               common.Address
	PermissionLessExecutionThresholdSeconds uint32
	MessageValidator                        common.Address
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
	ChainSelector      uint64
	Rmn                common.Address
	TokenAdminRegistry common.Address
	NonceManager       common.Address
}

var OffRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNV2\",\"name\":\"rmn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenGasOverride\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionTokenGasOverride\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"messageDestChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidMessageDestChainSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRoot\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidStaticConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeavesCannotBeEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionGasAmountCountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"errorReason\",\"type\":\"bytes\"}],\"name\":\"MessageValidationError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReleased\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePre\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePost\",\"type\":\"uint256\"}],\"name\":\"ReleaseOrMintBalanceMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"RootAlreadyCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleCommitReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelectorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SkippedReportExecution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNV2\",\"name\":\"rmn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"}],\"name\":\"StaticConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNV2\",\"name\":\"rmn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destExecData\",\"type\":\"bytes\"}],\"internalType\":\"structInternal.RampTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReportSingleChain[]\",\"name\":\"reports\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageValidator\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b50604051620068603803806200686083398101604081905262000035916200085e565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001fa565b5050466080525060208301516001600160a01b03161580620000ec575060408301516001600160a01b0316155b8062000103575060608301516001600160a01b0316155b1562000122576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160401b03166000036200014e5760405163c656089560e01b815260040160405180910390fd5b82516001600160401b0390811660a052602080850180516001600160a01b0390811660c05260408088018051831660e0526060808a01805185166101005283518b519098168852945184169587019590955251821690850152905116908201527f683eb52ee924eb817377cfa8f41f238f4bb7a877da5267869dfffbad85f564d89060800160405180910390a1620001e682620002a5565b620001f1816200036d565b50505062000be5565b336001600160a01b03821603620002545760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0316620002ce576040516342bcdf7f60e11b815260040160405180910390fd5b805160048054602080850180516001600160a01b039586166001600160c01b03199094168417600160a01b63ffffffff928316021790945560408087018051600580546001600160a01b031916918916919091179055815194855291519094169183019190915251909216908201527fa1c15688cb2c24508e158f6942b9276c6f3028a85e1af8cf3fff0c3ff3d5fc8d9060600160405180910390a150565b60005b8151811015620006125760008282815181106200039157620003916200099b565b60200260200101519050600081602001519050806001600160401b0316600003620003cf5760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b0316620003f8576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b03811660009081526006602052604081206001810180549192916200042490620009b1565b80601f01602080910402602001604051908101604052809291908181526020018280546200045290620009b1565b8015620004a35780601f106200047757610100808354040283529160200191620004a3565b820191906000526020600020905b8154815290600101906020018083116200048557829003601f168201915b5050505050905060008460600151905081516000036200054a578051600003620004e0576040516342bcdf7f60e11b815260040160405180910390fd5b60018301620004f0828262000a42565b508254600160a81b600160e81b031916600160a81b1783556040516001600160401b03851681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a162000585565b8080519060200120828051906020012014620005855760405163c39a620560e01b81526001600160401b038516600482015260240162000083565b604080860151845487516001600160a01b03166001600160a01b0319921515600160a01b02929092166001600160a81b031990911617178455516001600160401b038516907f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b90620005f990869062000b0e565b60405180910390a2505050505080600101905062000370565b5050565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b038111828210171562000651576200065162000616565b60405290565b604051601f8201601f191681016001600160401b038111828210171562000682576200068262000616565b604052919050565b80516001600160401b0381168114620006a257600080fd5b919050565b6001600160a01b0381168114620006bd57600080fd5b50565b6000601f83601f840112620006d457600080fd5b825160206001600160401b0380831115620006f357620006f362000616565b8260051b6200070483820162000657565b93845286810183019383810190898611156200071f57600080fd5b84890192505b8583101562000851578251848111156200073f5760008081fd5b89016080601f19828d038101821315620007595760008081fd5b620007636200062c565b888401516200077281620006a7565b81526040620007838582016200068a565b8a83015260608086015180151581146200079d5760008081fd5b83830152938501519389851115620007b55760008081fd5b84860195508f603f870112620007cd57600094508485fd5b8a860151945089851115620007e657620007e662000616565b620007f78b858f8801160162000657565b93508484528f828688010111156200080f5760008081fd5b60005b858110156200082f578681018301518582018d01528b0162000812565b5060009484018b01949094525091820152835250918401919084019062000725565b9998505050505050505050565b60008060008385036101008112156200087657600080fd5b60808112156200088557600080fd5b6200088f6200062c565b6200089a866200068a565b81526020860151620008ac81620006a7565b60208201526040860151620008c181620006a7565b60408201526060860151620008d681620006a7565b606082810191909152909450607f1982011215620008f357600080fd5b50604051606081016001600160401b0380821183831017156200091a576200091a62000616565b81604052608087015191506200093082620006a7565b90825260a08601519063ffffffff821682146200094c57600080fd5b81602084015260c087015191506200096482620006a7565b6040830182905260e0870151929450808311156200098157600080fd5b50506200099186828701620006c0565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680620009c657607f821691505b602082108103620009e757634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111562000a3d576000816000526020600020601f850160051c8101602086101562000a185750805b601f850160051c820191505b8181101562000a395782815560010162000a24565b5050505b505050565b81516001600160401b0381111562000a5e5762000a5e62000616565b62000a768162000a6f8454620009b1565b84620009ed565b602080601f83116001811462000aae576000841562000a955750858301515b600019600386901b1c1916600185901b17855562000a39565b600085815260208120601f198616915b8281101562000adf5788860151825594840194600190910190840162000abe565b508582101562000afe5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b602080825282546001600160a01b0381168383015260a081901c60ff161515604084015260a81c6001600160401b0316606083015260808083015260018084018054600093929190849062000b6381620009b1565b8060a089015260c0600183166000811462000b87576001811462000ba45762000bd6565b60ff19841660c08b015260c083151560051b8b0101945062000bd6565b85600052602060002060005b8481101562000bcd5781548c820185015290880190890162000bb0565b8b0160c0019550505b50929998505050505050505050565b60805160a05160c05160e05161010051615c0562000c5b6000396000818161022c015261292c0152600081816101fd0152612bfe0152600081816101ce01528181610552015281816106ff015261237001526000818161019f01526125af015260008181611d980152611dcb0152615c056000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c806379ba5097116100cd578063d434969111610081578063e9d68a8e11610066578063e9d68a8e146104bc578063f2fde38b146104dc578063f716f99f146104ef57600080fd5b8063d434969114610496578063e62653f5146104a957600080fd5b80638da5cb5b116100b25780638da5cb5b14610417578063c673e58414610432578063ccd37ba31461045257600080fd5b806379ba50971461040157806385572ffb1461040957600080fd5b8063311cd513116101245780635e36480c116101095780635e36480c1461033c5780636f9e320f1461035c5780637437ff9f1461036f57600080fd5b8063311cd5131461030e5780633f4b04aa1461032157600080fd5b806304666f9c1461015657806306285c691461016b578063181f5a77146102b25780632d04ab76146102fb575b600080fd5b610169610164366004613bc5565b610502565b005b61025c60408051608081018252600080825260208201819052918101829052606081019190915260405180608001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160401b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b6040516102a9919081516001600160401b031681526020808301516001600160a01b0390811691830191909152604080840151821690830152606092830151169181019190915260800190565b60405180910390f35b6102ee6040518060400160405280601181526020017f4f666652616d7020312e362e302d64657600000000000000000000000000000081525081565b6040516102a99190613d33565b610169610309366004613de3565b610516565b61016961031c366004613e95565b6109c5565b6009546040516001600160401b0390911681526020016102a9565b61034f61034a366004613ee8565b610a2e565b6040516102a99190613f45565b61016961036a366004613f70565b610a83565b6103cb604080516060810182526000808252602082018190529181019190915250604080516060810182526004546001600160a01b038082168352600160a01b90910463ffffffff166020830152600554169181019190915290565b6040805182516001600160a01b03908116825260208085015163ffffffff169083015292820151909216908201526060016102a9565b610169610a94565b610169610151366004613fc4565b6000546040516001600160a01b0390911681526020016102a9565b61044561044036600461400f565b610b45565b6040516102a9919061406f565b6104886104603660046140e4565b6001600160401b03919091166000908152600860209081526040808320938352929052205490565b6040519081526020016102a9565b6101696104a4366004614380565b610ca3565b6101696104b7366004614708565b610f83565b6104cf6104ca36600461494f565b61122a565b6040516102a9919061496a565b6101696104ea3660046149b7565b611336565b6101696104fd366004614a3c565b611347565b61050a611389565b610513816113e5565b50565b600061052487890189614d91565b602081015151909150156105b6576020810151604080830151905162ccd5af60e81b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263ccd5af009261058592600401614fb9565b60006040518083038186803b15801561059d57600080fd5b505afa1580156105b1573d6000803e3d6000fd5b505050505b805151511515806105cc57508051602001515115155b156106985760095460208a0135906001600160401b0380831691161015610670576009805467ffffffffffffffff19166001600160401b038316179055600480548351604051633937306f60e01b81526001600160a01b0390921692633937306f926106399291016150ed565b600060405180830381600087803b15801561065357600080fd5b505af1158015610667573d6000803e3d6000fd5b50505050610696565b81602001515160000361069657604051632261116760e01b815260040160405180910390fd5b505b60005b816020015151811015610906576000826020015182815181106106c0576106c061501b565b60209081029190910101518051604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b166004820152919250906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632cbc26bb90602401602060405180830381865afa158015610746573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061076a9190615100565b1561079857604051637edeb53960e11b81526001600160401b03821660048201526024015b60405180910390fd5b60006107a38261168b565b604084015181549192506001600160401b03908116600160a81b909204161415806107e7575082606001516001600160401b031683604001516001600160401b0316115b1561082c57825160408085015160608601519151636af0786b60e11b81526001600160401b03938416600482015290831660248201529116604482015260640161078f565b60808301518061084f5760405163504570e360e01b815260040160405180910390fd5b83516001600160401b03166000908152600860209081526040808320848452909152902054156108a75783516040516332cf0cbf60e01b81526001600160401b0390911660048201526024810182905260440161078f565b60608401516108b7906001615133565b825467ffffffffffffffff60a81b1916600160a81b6001600160401b0392831602179092559251166000908152600860209081526040808320948352939052919091204290555060010161069b565b50602081015181516040517f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e49261093e92909161515a565b60405180910390a16109ba60008a8a8a8a8a8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808e0282810182019093528d82529093508d92508c9182918501908490808284376000920191909152508b92506116d7915050565b505050505050505050565b610a056109d48284018461517f565b60408051600080825260208201909252906109ff565b60608152602001906001900390816109ea5790505b506119d0565b604080516000808252602082019092529050610a286001858585858660006116d7565b50505050565b6000610a3c600160046151b3565b6002610a496080856151dc565b6001600160401b0316610a5c9190615202565b610a668585611a92565b901c166003811115610a7a57610a7a613f1b565b90505b92915050565b610a8b611389565b61051381611ad7565b6001546001600160a01b03163314610aee5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161078f565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610b886040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c082015294855291820180548451818402810184019095528085529293858301939092830182828015610c3157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610c13575b5050505050815260200160038201805480602002602001604051908101604052809291908181526020018280548015610c9357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610c75575b5050505050815250509050919050565b333014610cc3576040516306e34e6560e31b815260040160405180910390fd5b6040805160008082526020820190925281610d00565b6040805180820190915260008082526020820152815260200190600190039081610cd95790505b5060a08701515190915015610d3657610d338660a001518760200151886060015189600001516020015189898989611bb6565b90505b6040805160a081018252875151815287516020908101516001600160401b031681830152808901518351600094840192610d71929101613d33565b60408051601f19818403018152918152908252898101516020830152018390526005549091506001600160a01b03168015610e4c576040516308d450a160e01b81526001600160a01b038216906308d450a190610dd29085906004016152ba565b600060405180830381600087803b158015610dec57600080fd5b505af1925050508015610dfd575060015b610e4c573d808015610e2b576040519150601f19603f3d011682016040523d82523d6000602084013e610e30565b606091505b50806040516309c2532560e01b815260040161078f9190613d33565b604088015151158015610e6157506080880151155b80610e78575060608801516001600160a01b03163b155b80610e9f57506060880151610e9d906001600160a01b03166385572ffb60e01b611d79565b155b15610eac57505050610f7c565b87516020908101516001600160401b03166000908152600690915260408082205460808b015160608c01519251633cf9798360e01b815284936001600160a01b0390931692633cf9798392610f0a92899261138892916004016152cd565b6000604051808303816000875af1158015610f29573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610f519190810190615309565b509150915081610f7657806040516302a35ba360e21b815260040161078f9190613d33565b50505050505b5050505050565b610f8b611d95565b815181518114610fae576040516320f8fd5960e21b815260040160405180910390fd5b60005b8181101561121a576000848281518110610fcd57610fcd61501b565b60200260200101519050600081602001515190506000858481518110610ff557610ff561501b565b6020026020010151905080518214611020576040516320f8fd5960e21b815260040160405180910390fd5b60005b8281101561120b57600082828151811061103f5761103f61501b565b60200260200101516000015190506000856020015183815181106110655761106561501b565b60200260200101519050816000146110b95780608001518210156110b9578551815151604051633a98d46360e11b81526001600160401b03909216600483015260248201526044810183905260640161078f565b8383815181106110cb576110cb61501b565b602002602001015160200151518160a00151511461111857805180516060909101516040516370a193fd60e01b815260048101929092526001600160401b0316602482015260440161078f565b60005b8160a00151518110156111fd57600085858151811061113c5761113c61501b565b60200260200101516020015182815181106111595761115961501b565b602002602001015163ffffffff169050806000146111f45760008360a0015183815181106111895761118961501b565b6020026020010151608001518060200190518101906111a8919061539e565b90508063ffffffff168210156111f2578351516040516348e617b360e01b815260048101919091526024810184905263ffffffff821660448201526064810183905260840161078f565b505b5060010161111b565b505050806001019050611023565b50505050806001019050610fb1565b5061122583836119d0565b505050565b60408051608080820183526000808352602080840182905283850182905260608085018190526001600160401b03878116845260068352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b9092049092169483019490945260018401805493949293918401916112b6906153bb565b80601f01602080910402602001604051908101604052809291908181526020018280546112e2906153bb565b8015610c935780601f1061130457610100808354040283529160200191610c93565b820191906000526020600020905b81548152906001019060200180831161131257505050919092525091949350505050565b61133e611389565b61051381611dfd565b61134f611389565b60005b81518110156113855761137d8282815181106113705761137061501b565b6020026020010151611ea6565b600101611352565b5050565b6000546001600160a01b031633146113e35760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161078f565b565b60005b81518110156113855760008282815181106114055761140561501b565b60200260200101519050600081602001519050806001600160401b03166000036114425760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b031661146a576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b0381166000908152600660205260408120600181018054919291611494906153bb565b80601f01602080910402602001604051908101604052809291908181526020018280546114c0906153bb565b801561150d5780601f106114e25761010080835404028352916020019161150d565b820191906000526020600020905b8154815290600101906020018083116114f057829003601f168201915b5050505050905060008460600151905081516000036115b0578051600003611548576040516342bcdf7f60e11b815260040160405180910390fd5b600183016115568282615445565b50825467ffffffffffffffff60a81b1916600160a81b1783556040516001600160401b03851681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a16115e9565b80805190602001208280519060200120146115e95760405163c39a620560e01b81526001600160401b038516600482015260240161078f565b604080860151845487516001600160a01b03166001600160a01b0319921515600160a01b02929092167fffffffffffffffffffffff00000000000000000000000000000000000000000090911617178455516001600160401b038516907f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b90611673908690615504565b60405180910390a250505050508060010190506113e8565b6001600160401b03811660009081526006602052604081208054600160a01b900460ff16610a7d5760405163ed053c5960e01b81526001600160401b038416600482015260240161078f565b60ff878116600090815260026020908152604080832081516080810183528154815260019091015480861693820193909352610100830485169181019190915262010000909104909216151560608301528735906117368760a46155d1565b905082606001511561177e57845161174f906020615202565b865161175c906020615202565b6117679060a06155d1565b61177191906155d1565b61177b90826155d1565b90505b3681146117a757604051638e1192e160e01b81526004810182905236602482015260440161078f565b50815181146117d65781516040516324f7d61360e21b815260048101919091526024810182905260440161078f565b6117de611d95565b60ff808a166000908152600360209081526040808320338452825280832081518083019092528054808616835293949193909284019161010090910416600281111561182c5761182c613f1b565b600281111561183d5761183d613f1b565b905250905060028160200151600281111561185a5761185a613f1b565b1480156118ae5750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff16815481106118965761189661501b565b6000918252602090912001546001600160a01b031633145b6118cb57604051631b41e11d60e31b815260040160405180910390fd5b5081606001511561197b5760208201516118e69060016155e4565b60ff16855114611909576040516371253a2560e01b815260040160405180910390fd5b835185511461192b5760405163a75d88af60e01b815260040160405180910390fd5b6000878760405161193d9291906155fd565b604051908190038120611954918b9060200161560d565b6040516020818303038152906040528051906020012090506119798a82888888612185565b505b6040805182815260208a8101356001600160401b03169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b81516000036119f15760405162bf199760e01b815260040160405180910390fd5b80516040805160008082526020820190925291159181611a34565b604080518082019091526000815260606020820152815260200190600190039081611a0c5790505b50905060005b8451811015610f7c57611a8a858281518110611a5857611a5861501b565b602002602001015184611a8457858381518110611a7757611a7761501b565b6020026020010151612342565b83612342565b600101611a3a565b6001600160401b038216600090815260076020526040812081611ab6608085615621565b6001600160401b031681526020810191909152604001600020549392505050565b80516001600160a01b0316611aff576040516342bcdf7f60e11b815260040160405180910390fd5b805160048054602080850180516001600160a01b039586167fffffffffffffffff0000000000000000000000000000000000000000000000009094168417600160a01b63ffffffff928316021790945560408087018051600580546001600160a01b031916918916919091179055815194855291519094169183019190915251909216908201527fa1c15688cb2c24508e158f6942b9276c6f3028a85e1af8cf3fff0c3ff3d5fc8d9060600160405180910390a150565b606088516001600160401b03811115611bd157611bd16139e5565b604051908082528060200260200182016040528015611c1657816020015b6040805180820190915260008082526020820152815260200190600190039081611bef5790505b509050811560005b8a51811015611d6b5781611cc857848482818110611c3e57611c3e61501b565b9050602002016020810190611c539190615647565b63ffffffff1615611cc857848482818110611c7057611c7061501b565b9050602002016020810190611c859190615647565b6040805163ffffffff9092166020830152016040516020818303038152906040528b8281518110611cb857611cb861501b565b6020026020010151608001819052505b611d468b8281518110611cdd57611cdd61501b565b60200260200101518b8b8b8b8b87818110611cfa57611cfa61501b565b9050602002810190611d0c9190615664565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612bb692505050565b838281518110611d5857611d5861501b565b6020908102919091010152600101611c1e565b505098975050505050505050565b6000611d8483612eb4565b8015610a7a5750610a7a8383612eff565b467f0000000000000000000000000000000000000000000000000000000000000000146113e357604051630f01ce8560e01b81527f0000000000000000000000000000000000000000000000000000000000000000600482015246602482015260440161078f565b336001600160a01b03821603611e555760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161078f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff16600003611ed1576000604051631b3fab5160e11b815260040161078f91906156aa565b60208082015160ff80821660009081526002909352604083206001810154929390928392169003611f22576060840151600182018054911515620100000262ff000019909216919091179055611f5e565b6060840151600182015460ff6201000090910416151590151514611f5e576040516321fd80df60e21b815260ff8416600482015260240161078f565b60a084015180516101001015611f8a576001604051631b3fab5160e11b815260040161078f91906156aa565b611ff08484600301805480602002602001604051908101604052809291908181526020018280548015611fe657602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611fc8575b5050505050612fa1565b8460600151156120fa5761205e8484600201805480602002602001604051908101604052809291908181526020018280548015611fe6576020028201919060005260206000209081546001600160a01b03168152600190910190602001808311611fc8575050505050612fa1565b60808501518051610100101561208a576002604051631b3fab5160e11b815260040161078f91906156aa565b604086015161209a9060036156c4565b60ff168151116120c0576003604051631b3fab5160e11b815260040161078f91906156aa565b805160018401805461ff00191661010060ff8416021790556120eb906002860190602084019061396b565b506120f88582600161300a565b505b6121068482600261300a565b805161211b906003850190602084019061396b565b5060408581015160018401805460ff191660ff8316179055865180855560a088015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547936121749389939260028a019291906156e0565b60405180910390a1610f7c84613165565b8251600090815b818110156123385760006001888684602081106121ab576121ab61501b565b6121b891901a601b6155e4565b8985815181106121ca576121ca61501b565b60200260200101518986815181106121e4576121e461501b565b602002602001015160405160008152602001604052604051612222949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015612244573d6000803e3d6000fd5b505060408051601f1981015160ff808e166000908152600360209081528582206001600160a01b038516835281528582208587019096528554808416865293975090955092939284019161010090041660028111156122a5576122a5613f1b565b60028111156122b6576122b6613f1b565b90525090506001816020015160028111156122d3576122d3613f1b565b146122f157604051636518c33d60e11b815260040160405180910390fd5b8051600160ff9091161b85161561231b57604051633d9ef1f160e21b815260040160405180910390fd5b806000015160ff166001901b85179450505080600101905061218c565b5050505050505050565b81518151604051632cbc26bb60e01b8152608083901b67ffffffffffffffff60801b166004820152901515907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa1580156123bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123e39190615100565b1561245457801561241257604051637edeb53960e11b81526001600160401b038316600482015260240161078f565b6040516001600160401b03831681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d8949339060200160405180910390a150505050565b600061245f8361168b565b600101805461246d906153bb565b80601f0160208091040260200160405190810160405280929190818152602001828054612499906153bb565b80156124e65780601f106124bb576101008083540402835291602001916124e6565b820191906000526020600020905b8154815290600101906020018083116124c957829003601f168201915b505050602088015151929350505060008190036125155760405162bf199760e01b815260040160405180910390fd5b856040015151811461253a576040516357e0e08360e01b815260040160405180910390fd5b6000816001600160401b03811115612554576125546139e5565b60405190808252806020026020018201604052801561257d578160200160208202803683370190505b50905060005b8281101561269e576000886020015182815181106125a3576125a361501b565b602002602001015190507f00000000000000000000000000000000000000000000000000000000000000006001600160401b03168160000151604001516001600160401b03161461261a5780516040908101519051631c21951160e11b81526001600160401b03909116600482015260240161078f565b866001600160401b03168160000151602001516001600160401b03161461266e57805160200151604051636c95f1eb60e01b81526001600160401b03808a166004830152909116602482015260440161078f565b61267881866131bc565b83838151811061268a5761268a61501b565b602090810291909101015250600101612583565b5060006126b586838a606001518b608001516132de565b9050806000036126e357604051633ee8bd3f60e11b81526001600160401b038716600482015260240161078f565b60005b838110156109ba5760005a905060008a60200151838151811061270b5761270b61501b565b6020026020010151905060006127298a836000015160600151610a2e565b9050600081600381111561273f5761273f613f1b565b148061275c5750600381600381111561275a5761275a613f1b565b145b6127b257815160600151604080516001600160401b03808e16825290921660208301527f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c910160405180910390a1505050612bae565b60608915612891578b85815181106127cc576127cc61501b565b6020908102919091018101510151600454909150600090600160a01b900463ffffffff166127fa88426151b3565b119050808061281a5750600383600381111561281857612818613f1b565b145b612842576040516354e7e43160e11b81526001600160401b038d16600482015260240161078f565b8c86815181106128545761285461501b565b60200260200101516000015160001461288b578c86815181106128795761287961501b565b60209081029190910101515160808501525b506128fd565b60008260038111156128a5576128a5613f1b565b146128fd57825160600151604080516001600160401b03808f16825290921660208301527f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120910160405180910390a150505050612bae565b8251608001516001600160401b0316156129d657600082600381111561292557612925613f1b565b036129d6577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e0e03cae8c85600001516080015186602001516040518463ffffffff1660e01b815260040161298693929190615792565b6020604051808303816000875af11580156129a5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129c99190615100565b6129d65750505050612bae565b60008d6040015186815181106129ee576129ee61501b565b6020026020010151905080518460a001515114612a3857835160600151604051631cfe6d8b60e01b81526001600160401b03808f166004830152909116602482015260440161078f565b612a4c8c856000015160600151600161331b565b600080612a5a8684866133c0565b91509150612a718e8760000151606001518461331b565b8c15612ac8576003826003811115612a8b57612a8b613f1b565b03612ac8576000856003811115612aa457612aa4613f1b565b14612ac857855151604051632b11b8d960e01b815261078f919083906004016157be565b6002826003811115612adc57612adc613f1b565b14612b21576003826003811115612af557612af5613f1b565b14612b21578d866000015160600151836040516349362d1f60e11b815260040161078f939291906157d7565b8560000151600001518660000151606001516001600160401b03168f6001600160401b03167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b8d8c81518110612b7957612b7961501b565b602002602001015186865a612b8e908f6151b3565b604051612b9e94939291906157fc565b60405180910390a4505050505050505b6001016126e6565b60408051808201909152600080825260208201526000612bd98760200151613474565b60405163bbe4f6db60e01b81526001600160a01b0380831660048301529192506000917f0000000000000000000000000000000000000000000000000000000000000000169063bbe4f6db90602401602060405180830381865afa158015612c45573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c699190615833565b90506001600160a01b0381161580612c985750612c966001600160a01b03821663aff2afbf60e01b611d79565b155b15612cc15760405163ae9b4ce960e01b81526001600160a01b038216600482015260240161078f565b600080612cec88858c60800151806020019051810190612ce1919061539e565b63ffffffff166134e8565b915091506000806000612d9f6040518061010001604052808e81526020018c6001600160401b031681526020018d6001600160a01b031681526020018f606001518152602001896001600160a01b031681526020018f6000015181526020018f6040015181526020018b815250604051602401612d699190615850565b60408051601f198184030181529190526020810180516001600160e01b0316633907753760e01b179052878661138860846135cb565b92509250925082612dc5578160405163e1cd550960e01b815260040161078f9190613d33565b8151602014612df4578151604051631e3be00960e21b815260206004820152602481019190915260440161078f565b600082806020019051810190612e0a919061591c565b9050866001600160a01b03168c6001600160a01b031614612e86576000612e3b8d8a612e36868a6151b3565b6134e8565b50905086811080612e55575081612e5288836151b3565b14155b15612e845760405163a966e21f60e01b815260048101839052602481018890526044810182905260640161078f565b505b604080518082019091526001600160a01b039098168852602088015250949550505050505095945050505050565b6000612ec7826301ffc9a760e01b612eff565b8015610a7d5750612ef8827fffffffff00000000000000000000000000000000000000000000000000000000612eff565b1592915050565b6040517fffffffff0000000000000000000000000000000000000000000000000000000082166024820152600090819060440160408051601f19818403018152919052602080820180516001600160e01b03166301ffc9a760e01b178152825192935060009283928392909183918a617530fa92503d91506000519050828015612f8a575060208210155b8015612f965750600081115b979650505050505050565b60005b81518110156112255760ff831660009081526003602052604081208351909190849084908110612fd657612fd661501b565b6020908102919091018101516001600160a01b03168252810191909152604001600020805461ffff19169055600101612fa4565b60005b8251811015610a2857600083828151811061302a5761302a61501b565b602002602001015190506000600281111561304757613047613f1b565b60ff80871660009081526003602090815260408083206001600160a01b0387168452909152902054610100900416600281111561308657613086613f1b565b146130a7576004604051631b3fab5160e11b815260040161078f91906156aa565b6001600160a01b0381166130ce5760405163d6c62c9b60e01b815260040160405180910390fd5b60405180604001604052808360ff1681526020018460028111156130f4576130f4613f1b565b905260ff80871660009081526003602090815260408083206001600160a01b0387168452825290912083518154931660ff198416811782559184015190929091839161ffff19161761010083600281111561315157613151613f1b565b02179055509050505080600101905061300d565b60ff81166105135760ff808216600090815260026020526040902060010154620100009004166131a857604051631e8ed32560e21b815260040160405180910390fd5b6009805467ffffffffffffffff1916905550565b815160208082015160409283015192516000938493613202937f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f93909291889101615935565b60408051601f1981840301815290829052805160209182012086518051888401516060808b0151908401516080808d0151950151959761324b9794969395929491939101615967565b604051602081830303815290604052805190602001208560400151805190602001208660a001516040516020016132829190615a6b565b60408051601f198184030181528282528051602091820120908301969096528101939093526060830191909152608082015260a081019190915260c0015b60405160208183030381529060405280519060200120905092915050565b6000806132ec8585856136a5565b6001600160401b038716600090815260086020908152604080832093835292905220549150505b949350505050565b6000600261332a6080856151dc565b6001600160401b031661333d9190615202565b9050600061334b8585611a92565b90508161335a600160046151b3565b901b19168183600381111561337157613371613f1b565b6001600160401b03871660009081526007602052604081209190921b9290921791829161339f608088615621565b6001600160401b031681526020810191909152604001600020555050505050565b60405163d434969160e01b8152600090606090309063d4349691906133ed90889088908890600401615b02565b600060405180830381600087803b15801561340757600080fd5b505af1925050508015613418575060015b613457573d808015613446576040519150601f19603f3d011682016040523d82523d6000602084013e61344b565b606091505b5060039250905061346c565b50506040805160208101909152600081526002905b935093915050565b6000815160201461349a578160405163046b337b60e51b815260040161078f9190613d33565b6000828060200190518101906134b0919061591c565b90506001600160a01b038111806134c8575061040081105b15610a7d578260405163046b337b60e51b815260040161078f9190613d33565b60008060008060006135498860405160240161351391906001600160a01b0391909116815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166370a0823160e01b179052888861138860846135cb565b9250925092508261356f578160405163e1cd550960e01b815260040161078f9190613d33565b602082511461359e578151604051631e3be00960e21b815260206004820152602481019190915260440161078f565b818060200190518101906135b2919061591c565b6135bc82886151b3565b94509450505050935093915050565b6000606060008361ffff166001600160401b038111156135ed576135ed6139e5565b6040519080825280601f01601f191660200182016040528015613617576020820181803683370190505b509150863b6136315763030ed58f60e21b60005260046000fd5b5a8581101561364b57632be8ca8b60e21b60005260046000fd5b859003604081048103871061366b576337c3be2960e01b60005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d8481111561368e5750835b808352806000602085013e50955095509592505050565b82518251600091908183036136cd57604051630469ac9960e21b815260040160405180910390fd5b61010182118015906136e157506101018111155b6136fe576040516309bde33960e01b815260040160405180910390fd5b60001982820101610100811115613728576040516309bde33960e01b815260040160405180910390fd5b8060000361375557866000815181106137435761374361501b565b60200260200101519350505050613923565b6000816001600160401b0381111561376f5761376f6139e5565b604051908082528060200260200182016040528015613798578160200160208202803683370190505b50905060008080805b858110156138c25760006001821b8b8116036137fc57888510156137e5578c5160018601958e9181106137d6576137d661501b565b6020026020010151905061381e565b85516001850194879181106137d6576137d661501b565b8b5160018401938d9181106138135761381361501b565b602002602001015190505b60008986101561384e578d5160018701968f91811061383f5761383f61501b565b60200260200101519050613870565b86516001860195889181106138655761386561501b565b602002602001015190505b82851115613891576040516309bde33960e01b815260040160405180910390fd5b61389b828261392a565b8784815181106138ad576138ad61501b565b602090810291909101015250506001016137a1565b5060018503821480156138d457508683145b80156138df57508581145b6138fc576040516309bde33960e01b815260040160405180910390fd5b8360018603815181106139115761391161501b565b60200260200101519750505050505050505b9392505050565b60008183106139425761393d8284613948565b610a7a565b610a7a83835b6040805160016020820152908101839052606081018290526000906080016132c0565b8280548282559060005260206000209081019282156139c0579160200282015b828111156139c057825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019061398b565b506139cc9291506139d0565b5090565b5b808211156139cc57600081556001016139d1565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715613a1d57613a1d6139e5565b60405290565b604051606081016001600160401b0381118282101715613a1d57613a1d6139e5565b60405160a081016001600160401b0381118282101715613a1d57613a1d6139e5565b60405160c081016001600160401b0381118282101715613a1d57613a1d6139e5565b604080519081016001600160401b0381118282101715613a1d57613a1d6139e5565b604051601f8201601f191681016001600160401b0381118282101715613ad357613ad36139e5565b604052919050565b60006001600160401b03821115613af457613af46139e5565b5060051b60200190565b6001600160a01b038116811461051357600080fd5b80356001600160401b0381168114613b2a57600080fd5b919050565b801515811461051357600080fd5b8035613b2a81613b2f565b60006001600160401b03821115613b6157613b616139e5565b50601f01601f191660200190565b600082601f830112613b8057600080fd5b8135613b93613b8e82613b48565b613aab565b818152846020838601011115613ba857600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215613bd857600080fd5b82356001600160401b0380821115613bef57600080fd5b818501915085601f830112613c0357600080fd5b8135613c11613b8e82613adb565b81815260059190911b83018401908481019088831115613c3057600080fd5b8585015b83811015613cd657803585811115613c4c5760008081fd5b86016080818c03601f1901811315613c645760008081fd5b613c6c6139fb565b89830135613c7981613afe565b81526040613c88848201613b13565b8b830152606080850135613c9b81613b2f565b83830152928401359289841115613cb457600091508182fd5b613cc28f8d86880101613b6f565b908301525085525050918601918601613c34565b5098975050505050505050565b60005b83811015613cfe578181015183820152602001613ce6565b50506000910152565b60008151808452613d1f816020860160208601613ce3565b601f01601f19169290920160200192915050565b602081526000610a7a6020830184613d07565b8060608101831015610a7d57600080fd5b60008083601f840112613d6957600080fd5b5081356001600160401b03811115613d8057600080fd5b602083019150836020828501011115613d9857600080fd5b9250929050565b60008083601f840112613db157600080fd5b5081356001600160401b03811115613dc857600080fd5b6020830191508360208260051b8501011115613d9857600080fd5b60008060008060008060008060e0898b031215613dff57600080fd5b613e098a8a613d46565b975060608901356001600160401b0380821115613e2557600080fd5b613e318c838d01613d57565b909950975060808b0135915080821115613e4a57600080fd5b613e568c838d01613d9f565b909750955060a08b0135915080821115613e6f57600080fd5b50613e7c8b828c01613d9f565b999c989b50969995989497949560c00135949350505050565b600080600060808486031215613eaa57600080fd5b613eb48585613d46565b925060608401356001600160401b03811115613ecf57600080fd5b613edb86828701613d57565b9497909650939450505050565b60008060408385031215613efb57600080fd5b613f0483613b13565b9150613f1260208401613b13565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b60048110613f4157613f41613f1b565b9052565b60208101610a7d8284613f31565b8035613b2a81613afe565b63ffffffff8116811461051357600080fd5b600060608284031215613f8257600080fd5b613f8a613a23565b8235613f9581613afe565b81526020830135613fa581613f5e565b60208201526040830135613fb881613afe565b60408201529392505050565b600060208284031215613fd657600080fd5b81356001600160401b03811115613fec57600080fd5b820160a0818503121561392357600080fd5b803560ff81168114613b2a57600080fd5b60006020828403121561402157600080fd5b610a7a82613ffe565b60008151808452602080850194506020840160005b838110156140645781516001600160a01b03168752958201959082019060010161403f565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a08401526140be60e084018261402a565b90506040840151601f198483030160c08501526140db828261402a565b95945050505050565b600080604083850312156140f757600080fd5b61410083613b13565b946020939093013593505050565b600060a0828403121561412057600080fd5b614128613a45565b90508135815261413a60208301613b13565b602082015261414b60408301613b13565b604082015261415c60608301613b13565b606082015261416d60808301613b13565b608082015292915050565b600082601f83011261418957600080fd5b81356020614199613b8e83613adb565b82815260059290921b840181019181810190868411156141b857600080fd5b8286015b848110156142b25780356001600160401b03808211156141dc5760008081fd5b9088019060a0828b03601f19018113156141f65760008081fd5b6141fe613a45565b87840135838111156142105760008081fd5b61421e8d8a83880101613b6f565b825250604080850135848111156142355760008081fd5b6142438e8b83890101613b6f565b8a840152506060808601358581111561425c5760008081fd5b61426a8f8c838a0101613b6f565b838501525060809150818601358184015250828501359250838311156142905760008081fd5b61429e8d8a85880101613b6f565b9082015286525050509183019183016141bc565b509695505050505050565b600061014082840312156142d057600080fd5b6142d8613a67565b90506142e4838361410e565b815260a08201356001600160401b038082111561430057600080fd5b61430c85838601613b6f565b602084015260c084013591508082111561432557600080fd5b61433185838601613b6f565b604084015261434260e08501613f53565b6060840152610100840135608084015261012084013591508082111561436757600080fd5b5061437484828501614178565b60a08301525092915050565b60008060008060006060868803121561439857600080fd5b85356001600160401b03808211156143af57600080fd5b6143bb89838a016142bd565b965060208801359150808211156143d157600080fd5b6143dd89838a01613d9f565b909650945060408801359150808211156143f657600080fd5b5061440388828901613d9f565b969995985093965092949392505050565b600082601f83011261442557600080fd5b81356020614435613b8e83613adb565b82815260059290921b8401810191818101908684111561445457600080fd5b8286015b848110156142b25780356001600160401b038111156144775760008081fd5b6144858986838b01016142bd565b845250918301918301614458565b600082601f8301126144a457600080fd5b813560206144b4613b8e83613adb565b82815260059290921b840181019181810190868411156144d357600080fd5b8286015b848110156142b25780356001600160401b03808211156144f657600080fd5b818901915089603f83011261450a57600080fd5b8582013561451a613b8e82613adb565b81815260059190911b830160400190878101908c83111561453a57600080fd5b604085015b838110156145735780358581111561455657600080fd5b6145658f6040838a0101613b6f565b84525091890191890161453f565b508752505050928401925083016144d7565b600082601f83011261459657600080fd5b813560206145a6613b8e83613adb565b8083825260208201915060208460051b8701019350868411156145c857600080fd5b602086015b848110156142b257803583529183019183016145cd565b600082601f8301126145f557600080fd5b81356020614605613b8e83613adb565b82815260059290921b8401810191818101908684111561462457600080fd5b8286015b848110156142b25780356001600160401b03808211156146485760008081fd5b9088019060a0828b03601f19018113156146625760008081fd5b61466a613a45565b614675888501613b13565b81526040808501358481111561468b5760008081fd5b6146998e8b83890101614414565b8a84015250606080860135858111156146b25760008081fd5b6146c08f8c838a0101614493565b83850152506080915081860135858111156146db5760008081fd5b6146e98f8c838a0101614585565b9184019190915250919093013590830152508352918301918301614628565b6000806040838503121561471b57600080fd5b6001600160401b038335111561473057600080fd5b61473d84843585016145e4565b91506001600160401b036020840135111561475757600080fd5b6020830135830184601f82011261476d57600080fd5b61477a613b8e8235613adb565b81358082526020808301929160051b84010187101561479857600080fd5b602083015b6020843560051b850101811015614941576001600160401b03813511156147c357600080fd5b87603f8235860101126147d557600080fd5b6147e8613b8e6020833587010135613adb565b81358501602081810135808452908301929160059190911b016040018a101561481057600080fd5b604083358701015b83358701602081013560051b01604001811015614931576001600160401b038135111561484457600080fd5b833587018135016040818d03603f1901121561485f57600080fd5b614867613a89565b604082013581526001600160401b036060830135111561488657600080fd5b8c605f60608401358401011261489b57600080fd5b60406060830135830101356148b2613b8e82613adb565b808282526020820191508f60608460051b60608801358801010111156148d757600080fd5b6060808601358601015b60608460051b606088013588010101811015614911576149018135613f5e565b80358352602092830192016148e1565b508060208501525050508085525050602083019250602081019050614818565b508452506020928301920161479d565b508093505050509250929050565b60006020828403121561496157600080fd5b610a7a82613b13565b602081526001600160a01b0382511660208201526020820151151560408201526001600160401b0360408301511660608201526000606083015160808084015261331360a0840182613d07565b6000602082840312156149c957600080fd5b813561392381613afe565b600082601f8301126149e557600080fd5b813560206149f5613b8e83613adb565b8083825260208201915060208460051b870101935086841115614a1757600080fd5b602086015b848110156142b2578035614a2f81613afe565b8352918301918301614a1c565b60006020808385031215614a4f57600080fd5b82356001600160401b0380821115614a6657600080fd5b818501915085601f830112614a7a57600080fd5b8135614a88613b8e82613adb565b81815260059190911b83018401908481019088831115614aa757600080fd5b8585015b83811015613cd657803585811115614ac257600080fd5b860160c0818c03601f19011215614ad95760008081fd5b614ae1613a67565b8882013581526040614af4818401613ffe565b8a8301526060614b05818501613ffe565b8284015260809150614b18828501613b3d565b9083015260a08381013589811115614b305760008081fd5b614b3e8f8d838801016149d4565b838501525060c0840135915088821115614b585760008081fd5b614b668e8c848701016149d4565b9083015250845250918601918601614aab565b80356001600160e01b0381168114613b2a57600080fd5b600082601f830112614ba157600080fd5b81356020614bb1613b8e83613adb565b82815260069290921b84018101918181019086841115614bd057600080fd5b8286015b848110156142b25760408189031215614bed5760008081fd5b614bf5613a89565b614bfe82613b13565b8152614c0b858301614b79565b81860152835291830191604001614bd4565b600082601f830112614c2e57600080fd5b81356020614c3e613b8e83613adb565b82815260059290921b84018101918181019086841115614c5d57600080fd5b8286015b848110156142b25780356001600160401b0380821115614c815760008081fd5b9088019060a0828b03601f1901811315614c9b5760008081fd5b614ca3613a45565b614cae888501613b13565b815260408085013584811115614cc45760008081fd5b614cd28e8b83890101613b6f565b8a8401525060609350614ce6848601613b13565b908201526080614cf7858201613b13565b93820193909352920135908201528352918301918301614c61565b600082601f830112614d2357600080fd5b81356020614d33613b8e83613adb565b82815260069290921b84018101918181019086841115614d5257600080fd5b8286015b848110156142b25760408189031215614d6f5760008081fd5b614d77613a89565b813581528482013585820152835291830191604001614d56565b60006020808385031215614da457600080fd5b82356001600160401b0380821115614dbb57600080fd5b9084019060608287031215614dcf57600080fd5b614dd7613a23565b823582811115614de657600080fd5b83016040818903811315614df957600080fd5b614e01613a89565b823585811115614e1057600080fd5b8301601f81018b13614e2157600080fd5b8035614e2f613b8e82613adb565b81815260069190911b8201890190898101908d831115614e4e57600080fd5b928a01925b82841015614e9e5785848f031215614e6b5760008081fd5b614e73613a89565b8435614e7e81613afe565b8152614e8b858d01614b79565b818d0152825292850192908a0190614e53565b845250505082870135915084821115614eb657600080fd5b614ec28a838501614b90565b81880152835250508284013582811115614edb57600080fd5b614ee788828601614c1d565b85830152506040830135935081841115614f0057600080fd5b614f0c87858501614d12565b60408201529695505050505050565b600082825180855260208086019550808260051b84010181860160005b84811015614fac57601f19868403018952815160a06001600160401b03808351168652868301518288880152614f7083880182613d07565b60408581015184169089015260608086015190931692880192909252506080928301519290950191909152509783019790830190600101614f38565b5090979650505050505050565b60006040808352614fcd6040840186614f1b565b83810360208581019190915285518083528682019282019060005b8181101561500d57845180518452840151848401529383019391850191600101614fe8565b509098975050505050505050565b634e487b7160e01b600052603260045260246000fd5b805160408084528151848201819052600092602091908201906060870190855b8181101561508857835180516001600160a01b031684528501516001600160e01b0316858401529284019291850191600101615051565b50508583015187820388850152805180835290840192506000918401905b808310156150e157835180516001600160401b031683528501516001600160e01b0316858301529284019260019290920191908501906150a6565b50979650505050505050565b602081526000610a7a6020830184615031565b60006020828403121561511257600080fd5b815161392381613b2f565b634e487b7160e01b600052601160045260246000fd5b6001600160401b038181168382160190808211156151535761515361511d565b5092915050565b60408152600061516d6040830185614f1b565b82810360208401526140db8185615031565b60006020828403121561519157600080fd5b81356001600160401b038111156151a757600080fd5b613313848285016145e4565b81810381811115610a7d57610a7d61511d565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b03808416806151f6576151f66151c6565b92169190910692915050565b8082028115828204841417610a7d57610a7d61511d565b80518252600060206001600160401b0381840151168185015260408084015160a0604087015261524c60a0870182613d07565b9050606085015186820360608801526152658282613d07565b608087810151898303918a01919091528051808352908601935060009250908501905b808310156150e157835180516001600160a01b0316835286015186830152928501926001929092019190840190615288565b602081526000610a7a6020830184615219565b6080815260006152e06080830187615219565b61ffff9590951660208301525060408101929092526001600160a01b0316606090910152919050565b60008060006060848603121561531e57600080fd5b835161532981613b2f565b60208501519093506001600160401b0381111561534557600080fd5b8401601f8101861361535657600080fd5b8051615364613b8e82613b48565b81815287602083850101111561537957600080fd5b61538a826020830160208601613ce3565b809450505050604084015190509250925092565b6000602082840312156153b057600080fd5b815161392381613f5e565b600181811c908216806153cf57607f821691505b6020821081036153ef57634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115611225576000816000526020600020601f850160051c8101602086101561541e5750805b601f850160051c820191505b8181101561543d5782815560010161542a565b505050505050565b81516001600160401b0381111561545e5761545e6139e5565b6154728161546c84546153bb565b846153f5565b602080601f8311600181146154a7576000841561548f5750858301515b600019600386901b1c1916600185901b17855561543d565b600085815260208120601f198616915b828110156154d6578886015182559484019460019091019084016154b7565b50858210156154f45787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000602080835283546001600160a01b038116602085015260ff8160a01c16151560408501526001600160401b038160a81c16606085015250600180850160808086015260008154615555816153bb565b8060a089015260c060018316600081146155765760018114615592576155c2565b60ff19841660c08b015260c083151560051b8b010194506155c2565b85600052602060002060005b848110156155b95781548c820185015290880190890161559e565b8b0160c0019550505b50929998505050505050505050565b80820180821115610a7d57610a7d61511d565b60ff8181168382160190811115610a7d57610a7d61511d565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006001600160401b038084168061563b5761563b6151c6565b92169190910492915050565b60006020828403121561565957600080fd5b813561392381613f5e565b6000808335601e1984360301811261567b57600080fd5b8301803591506001600160401b0382111561569557600080fd5b602001915036819003821315613d9857600080fd5b60208101600583106156be576156be613f1b565b91905290565b60ff81811683821602908116908181146151535761515361511d565b600060a0820160ff881683526020878185015260a0604085015281875480845260c0860191508860005282600020935060005b818110156157385784546001600160a01b031683526001948501949284019201615713565b50508481036060860152865180825290820192508187019060005b818110156157785782516001600160a01b031685529383019391830191600101615753565b50505060ff851660808501525090505b9695505050505050565b60006001600160401b038086168352808516602084015250606060408301526140db6060830184613d07565b8281526040602082015260006133136040830184613d07565b6001600160401b03848116825283166020820152606081016133136040830184613f31565b84815261580c6020820185613f31565b6080604082015260006158226080830185613d07565b905082606083015295945050505050565b60006020828403121561584557600080fd5b815161392381613afe565b602081526000825161010080602085015261586f610120850183613d07565b9150602085015161588b60408601826001600160401b03169052565b5060408501516001600160a01b0381166060860152506060850151608085015260808501516158c560a08601826001600160a01b03169052565b5060a0850151601f19808685030160c08701526158e28483613d07565b935060c08701519150808685030160e08701526158ff8483613d07565b935060e08701519150808685030183870152506157888382613d07565b60006020828403121561592e57600080fd5b5051919050565b84815260006001600160401b038086166020840152808516604084015250608060608301526157886080830184613d07565b86815260c06020820152600061598060c0830188613d07565b6001600160a01b03969096166040830152506001600160401b039384166060820152608081019290925290911660a09091015292915050565b600082825180855260208086019550808260051b84010181860160005b84811015614fac57601f19868403018952815160a081518186526159fc82870182613d07565b9150508582015185820387870152615a148282613d07565b91505060408083015186830382880152615a2e8382613d07565b92505050606080830151818701525060808083015192508582038187015250615a578183613d07565b9a86019a94505050908301906001016159d6565b602081526000610a7a60208301846159b9565b60008282518085526020808601955060208260051b8401016020860160005b84811015614fac57601f19868403018952615ab9838351613d07565b98840198925090830190600101615a9d565b60008151808452602080850194506020840160005b8381101561406457815163ffffffff1687529582019590820190600101615ae0565b60608152600084518051606084015260208101516001600160401b0380821660808601528060408401511660a08601528060608401511660c08601528060808401511660e0860152505050602085015161014080610100850152615b6a6101a0850183613d07565b91506040870151605f198086850301610120870152615b898483613d07565b935060608901519150615ba6838701836001600160a01b03169052565b608089015161016087015260a0890151925080868503016101808701525050615bcf82826159b9565b9150508281036020840152615be48186615a7e565b905082810360408401526157888185615acb56fea164736f6c6343000818000a",
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

func (_OffRamp *OffRampTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampTransactorSession) ExecuteSingleMessage(message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _OffRamp.Contract.ExecuteSingleMessage(&_OffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_OffRamp *OffRampTransactor) ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.contract.Transact(opts, "manuallyExecute", reports, gasLimitOverrides)
}

func (_OffRamp *OffRampSession) ManuallyExecute(reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
	return _OffRamp.Contract.ManuallyExecute(&_OffRamp.TransactOpts, reports, gasLimitOverrides)
}

func (_OffRamp *OffRampTransactorSession) ManuallyExecute(reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error) {
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
	return common.HexToHash("0xa1c15688cb2c24508e158f6942b9276c6f3028a85e1af8cf3fff0c3ff3d5fc8d")
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

	LatestConfigDetails(opts *bind.CallOpts, ocrPluginType uint8) (MultiOCR3BaseOCRConfig, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []OffRampSourceChainConfigArgs) (*types.Transaction, error)

	Commit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Execute(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error)

	ExecuteSingleMessage(opts *bind.TransactOpts, message InternalAny2EVMRampMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error)

	ManuallyExecute(opts *bind.TransactOpts, reports []InternalExecutionReportSingleChain, gasLimitOverrides [][]OffRampGasLimitOverride) (*types.Transaction, error)

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
