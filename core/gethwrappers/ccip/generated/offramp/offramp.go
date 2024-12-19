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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reportOnRamp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"configOnRamp\",\"type\":\"bytes\"}],\"name\":\"CommitOnRampMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyBatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenGasOverride\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionTokenGasOverride\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"messageDestChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidMessageDestChainSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidOnRampUpdate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeavesCannotBeEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionGasAmountCountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"errorReason\",\"type\":\"bytes\"}],\"name\":\"MessageValidationError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReleased\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePre\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePost\",\"type\":\"uint256\"}],\"name\":\"ReleaseOrMintBalanceMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"RootAlreadyCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleCommitReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelectorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SkippedReportExecution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"}],\"name\":\"StaticConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllSourceChainConfigs\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReport[]\",\"name\":\"reports\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101406040523480156200001257600080fd5b5060405162006d1138038062006d11833981016040819052620000359162000936565b336000816200005757604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008a576200008a81620001dc565b50504660805260408301516001600160a01b03161580620000b6575060608301516001600160a01b0316155b80620000cd575060808301516001600160a01b0316155b15620000ec576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160401b0316600003620001185760405163c656089560e01b815260040160405180910390fd5b82516001600160401b0390811660a0908152604080860180516001600160a01b0390811660c05260608089018051831660e0526080808b0180518516610100526020808d01805161ffff9081166101205289518f51909c168c52905116908a0152945184169588019590955251821690860152905116908301527fb0fa1fb01508c5097c502ad056fd77018870c9be9a86d9e56b6b471862d7c5b7910160405180910390a1620001c88262000256565b620001d38162000344565b50505062000cde565b336001600160a01b038216036200020657604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03166200027f576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516001600160a01b039889166001600160c01b03199097168717600160a01b63ffffffff958616021760ff60c01b1916600160c01b911515919091021790965560608089018051600580546001600160a01b031916918b169190911790558251968752935190921693850193909352935115159183019190915251909216908201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9060800160405180910390a150565b60005b8151811015620005d957600082828151811062000368576200036862000a16565b60200260200101519050600081602001519050806001600160401b0316600003620003a65760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b0316620003cf576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b03811660009081526008602052604090206060830151600182018054620003fd9062000a2c565b905060000362000460578154600160a81b600160e81b031916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1620004d1565b8154600160a81b90046001600160401b0316600114801590620004a35750805160208201206040516200049890600185019062000a68565b604051809103902014155b15620004d157604051632105803760e11b81526001600160401b038416600482015260240160405180910390fd5b80511580620005075750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b1562000526576040516342bcdf7f60e11b815260040160405180910390fd5b6001820162000536828262000b3b565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b02929092166001600160a81b0319909116171782556200058560066001600160401b038516620005dd565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b83604051620005c1919062000c07565b60405180910390a25050505080600101905062000347565b5050565b6000620005eb8383620005f4565b90505b92915050565b60008181526001830160205260408120546200063d57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620005ee565b506000620005ee565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b038111828210171562000681576200068162000646565b60405290565b60405160a081016001600160401b038111828210171562000681576200068162000646565b604051601f8201601f191681016001600160401b0381118282101715620006d757620006d762000646565b604052919050565b80516001600160401b0381168114620006f757600080fd5b919050565b6001600160a01b03811681146200071257600080fd5b50565b80518015158114620006f757600080fd5b6000608082840312156200073957600080fd5b620007436200065c565b905081516200075281620006fc565b8152602082015163ffffffff811681146200076c57600080fd5b60208201526200077f6040830162000715565b604082015260608201516200079481620006fc565b606082015292915050565b6000601f83601f840112620007b357600080fd5b825160206001600160401b0380831115620007d257620007d262000646565b8260051b620007e3838201620006ac565b9384528681018301938381019089861115620007fe57600080fd5b84890192505b8583101562000929578251848111156200081e5760008081fd5b89016080601f19828d038101821315620008385760008081fd5b620008426200065c565b888401516200085181620006fc565b8152604062000862858201620006df565b8a83015260606200087581870162000715565b838301529385015193898511156200088d5760008081fd5b84860195508f603f870112620008a557600094508485fd5b8a860151945089851115620008be57620008be62000646565b620008cf8b858f88011601620006ac565b93508484528f82868801011115620008e75760008081fd5b60005b8581101562000907578681018301518582018d01528b01620008ea565b5060009484018b01949094525091820152835250918401919084019062000804565b9998505050505050505050565b60008060008385036101408112156200094e57600080fd5b60a08112156200095d57600080fd5b506200096862000687565b6200097385620006df565b8152602085015161ffff811681146200098b57600080fd5b60208201526040850151620009a081620006fc565b60408201526060850151620009b581620006fc565b60608201526080850151620009ca81620006fc565b60808201529250620009e08560a0860162000726565b6101208501519092506001600160401b03811115620009fe57600080fd5b62000a0c868287016200079f565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b600181811c9082168062000a4157607f821691505b60208210810362000a6257634e487b7160e01b600052602260045260246000fd5b50919050565b600080835462000a788162000a2c565b6001828116801562000a93576001811462000aa95762000ada565b60ff198416875282151583028701945062000ada565b8760005260208060002060005b8581101562000ad15781548a82015290840190820162000ab6565b50505082870194505b50929695505050505050565b601f82111562000b36576000816000526020600020601f850160051c8101602086101562000b115750805b601f850160051c820191505b8181101562000b325782815560010162000b1d565b5050505b505050565b81516001600160401b0381111562000b575762000b5762000646565b62000b6f8162000b68845462000a2c565b8462000ae6565b602080601f83116001811462000ba7576000841562000b8e5750858301515b600019600386901b1c1916600185901b17855562000b32565b600085815260208120601f198616915b8281101562000bd85788860151825594840194600190910190840162000bb7565b508582101562000bf75787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b602080825282546001600160a01b0381168383015260a081901c60ff161515604084015260a81c6001600160401b0316606083015260808083015260018084018054600093929190849062000c5c8162000a2c565b8060a089015260c0600183166000811462000c80576001811462000c9d5762000ccf565b60ff19841660c08b015260c083151560051b8b0101945062000ccf565b85600052602060002060005b8481101562000cc65781548c820185015290880190890162000ca9565b8b0160c0019550505b50929998505050505050505050565b60805160a05160c05160e0516101005161012051615f9662000d7b600039600081816101b001528181610ce801528181612ef401526138160152600081816102380152612ac40152600081816102090152612d6c0152600081816101da01528181610fcc0152818161117c01526124c40152600081816101810152818161266f015261272601526000818161195b015261198e0152615f966000f3fe608060405234801561001057600080fd5b506004361061012c5760003560e01c80637edf52f4116100ad578063de5e0b9a11610071578063de5e0b9a146104eb578063e9d68a8e146104fe578063f2fde38b1461051e578063f58e03fc14610531578063f716f99f1461054457600080fd5b80637edf52f41461044b57806385572ffb1461045e5780638da5cb5b1461046c578063c673e58414610487578063ccd37ba3146104a757600080fd5b80635e36480c116100f45780635e36480c146103405780635e7bb0081461036057806360987c20146103735780637437ff9f1461038657806379ba50971461044357600080fd5b806304666f9c1461013157806306285c6914610146578063181f5a77146102c65780633f4b04aa1461030f5780635215505b1461032a575b600080fd5b61014461013f366004613eba565b610557565b005b6102686040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526040518060a001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160401b031681526020017f000000000000000000000000000000000000000000000000000000000000000061ffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b6040805182516001600160401b0316815260208084015161ffff1690820152828201516001600160a01b03908116928201929092526060808401518316908201526080928301519091169181019190915260a0015b60405180910390f35b6103026040518060400160405280601181526020017f4f666652616d7020312e362e302d64657600000000000000000000000000000081525081565b6040516102bd9190614028565b600b546040516001600160401b0390911681526020016102bd565b61033261056b565b6040516102bd929190614082565b61035361034e366004614123565b6107c6565b6040516102bd9190614180565b61014461036e3660046146e9565b61081b565b610144610381366004614978565b610aaf565b6103fc60408051608081018252600080825260208201819052918101829052606081019190915250604080516080810182526004546001600160a01b038082168352600160a01b820463ffffffff166020840152600160c01b90910460ff16151592820192909252600554909116606082015290565b6040516102bd919081516001600160a01b03908116825260208084015163ffffffff1690830152604080840151151590830152606092830151169181019190915260800190565b610144610d8a565b610144610459366004614a0c565b610e0d565b61014461012c366004614a71565b6001546040516001600160a01b0390911681526020016102bd565b61049a610495366004614abc565b610e1e565b6040516102bd9190614b1c565b6104dd6104b5366004614b91565b6001600160401b03919091166000908152600a60209081526040808320938352929052205490565b6040519081526020016102bd565b6101446104f9366004614c0d565b610f7c565b61051161050c366004614cbf565b61147f565b6040516102bd9190614cda565b61014461052c366004614ced565b61158b565b61014461053f366004614d0a565b61159c565b610144610552366004614dc5565b611605565b61055f611647565b61056881611674565b50565b606080600061057a60066118fd565b6001600160401b0381111561059157610591613cda565b6040519080825280602002602001820160405280156105e257816020015b60408051608081018252600080825260208083018290529282015260608082015282526000199092019101816105af5790505b50905060006105f160066118fd565b6001600160401b0381111561060857610608613cda565b604051908082528060200260200182016040528015610631578160200160208202803683370190505b50905060005b61064160066118fd565b8110156107bd57610653600682611907565b82828151811061066557610665614f02565b60200260200101906001600160401b031690816001600160401b0316815250506008600083838151811061069b5761069b614f02565b6020908102919091018101516001600160401b039081168352828201939093526040918201600020825160808101845281546001600160a01b038116825260ff600160a01b820416151593820193909352600160a81b9092049093169181019190915260018201805491929160608401919061071690614f18565b80601f016020809104026020016040519081016040528092919081815260200182805461074290614f18565b801561078f5780601f106107645761010080835404028352916020019161078f565b820191906000526020600020905b81548152906001019060200180831161077257829003601f168201915b5050505050815250508382815181106107aa576107aa614f02565b6020908102919091010152600101610637565b50939092509050565b60006107d460016004614f68565b60026107e1608085614f91565b6001600160401b03166107f49190614fb7565b6107fe8585611913565b901c16600381111561081257610812614156565b90505b92915050565b610823611958565b815181518114610846576040516320f8fd5960e21b815260040160405180910390fd5b60005b81811015610a9f57600084828151811061086557610865614f02565b6020026020010151905060008160200151519050600085848151811061088d5761088d614f02565b60200260200101519050805182146108b8576040516320f8fd5960e21b815260040160405180910390fd5b60005b82811015610a905760008282815181106108d7576108d7614f02565b60200260200101516000015190506000856020015183815181106108fd576108fd614f02565b6020026020010151905081600014610956578060800151821015610956578551815151604051633a98d46360e11b81526001600160401b0390921660048301526024820152604481018390526064015b60405180910390fd5b83838151811061096857610968614f02565b602002602001015160200151518160a0015151146109b557805180516060909101516040516370a193fd60e01b815260048101929092526001600160401b0316602482015260440161094d565b60005b8160a0015151811015610a825760008585815181106109d9576109d9614f02565b60200260200101516020015182815181106109f6576109f6614f02565b602002602001015163ffffffff16905080600014610a795760008360a001518381518110610a2657610a26614f02565b60200260200101516040015163ffffffff16905080821015610a77578351516040516348e617b360e01b8152600481019190915260248101849052604481018290526064810183905260840161094d565b505b506001016109b8565b5050508060010190506108bb565b50505050806001019050610849565b50610aaa83836119c0565b505050565b333014610acf576040516306e34e6560e31b815260040160405180910390fd5b6040805160008082526020820190925281610b0c565b6040805180820190915260008082526020820152815260200190600190039081610ae55790505b5060a08701515190915015610b4257610b3f8660a001518760200151886060015189600001516020015189898989611a83565b90505b6040805160a081018252875151815287516020908101516001600160401b03168183015288015181830152908701516060820152608081018290526005546001600160a01b03168015610c35576040516308d450a160e01b81526001600160a01b038216906308d450a190610bbb90859060040161507b565b600060405180830381600087803b158015610bd557600080fd5b505af1925050508015610be6575060015b610c35573d808015610c14576040519150601f19603f3d011682016040523d82523d6000602084013e610c19565b606091505b50806040516309c2532560e01b815260040161094d9190614028565b604088015151158015610c4a57506080880151155b80610c61575060608801516001600160a01b03163b155b80610c8857506060880151610c86906001600160a01b03166385572ffb60e01b611c34565b155b15610c9557505050610d83565b87516020908101516001600160401b03166000908152600890915260408082205460808b015160608c01519251633cf9798360e01b815284936001600160a01b0390931692633cf9798392610d119289927f0000000000000000000000000000000000000000000000000000000000000000929160040161508e565b6000604051808303816000875af1158015610d30573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610d5891908101906150ca565b509150915081610d7d57806040516302a35ba360e21b815260040161094d9190614028565b50505050505b5050505050565b6000546001600160a01b03163314610db55760405163015aa1e360e11b815260040160405180910390fd5b600180546001600160a01b0319808216339081179093556000805490911681556040516001600160a01b03909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610e15611647565b61056881611c73565b610e616040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c082015294855291820180548451818402810184019095528085529293858301939092830182828015610f0a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610eec575b5050505050815260200160038201805480602002602001604051908101604052809291908181526020018280548015610f6c57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610f4e575b5050505050815250509050919050565b6000610f8a87890189615377565b6004805491925090600160c01b900460ff1661103457602082015151156110345760208201516040808401519051633854844f60e11b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016926370a9089e92611003923092919060040161559f565b60006040518083038186803b15801561101b57600080fd5b505afa15801561102f573d6000803e3d6000fd5b505050505b8151515115158061104a57508151602001515115155b1561111557600b5460208b0135906001600160401b03808316911610156110ed57600b805467ffffffffffffffff19166001600160401b03831617905581548351604051633937306f60e01b81526001600160a01b0390921691633937306f916110b6916004016156b2565b600060405180830381600087803b1580156110d057600080fd5b505af11580156110e4573d6000803e3d6000fd5b50505050611113565b82602001515160000361111357604051632261116760e01b815260040160405180910390fd5b505b60005b8260200151518110156113cb5760008360200151828151811061113d5761113d614f02565b60209081029190910101518051604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b166004820152919250906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632cbc26bb90602401602060405180830381865afa1580156111c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111e791906156c5565b1561121057604051637edeb53960e11b81526001600160401b038216600482015260240161094d565b600061121b82611d78565b90508060010160405161122e91906156e2565b60405180910390208360200151805190602001201461126b5782602001518160010160405163b80d8fa960e01b815260040161094d9291906157d5565b60408301518154600160a81b90046001600160401b0390811691161415806112ac575082606001516001600160401b031683604001516001600160401b0316115b156112f157825160408085015160608601519151636af0786b60e11b81526001600160401b03938416600482015290831660248201529116604482015260640161094d565b6080830151806113145760405163504570e360e01b815260040160405180910390fd5b83516001600160401b03166000908152600a602090815260408083208484529091529020541561136c5783516040516332cf0cbf60e01b81526001600160401b0390911660048201526024810182905260440161094d565b606084015161137c9060016157fa565b825467ffffffffffffffff60a81b1916600160a81b6001600160401b0392831602179092559251166000908152600a602090815260408083209483529390529190912042905550600101611118565b50602082015182516040517f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e492611403929091615821565b60405180910390a1610d7d60008b8b8b8b8b8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808f0282810182019093528e82529093508e92508d9182918501908490808284376000920191909152508c9250611dc4915050565b60408051608080820183526000808352602080840182905283850182905260608085018190526001600160401b03878116845260088352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b90920490921694830194909452600184018054939492939184019161150b90614f18565b80601f016020809104026020016040519081016040528092919081815260200182805461153790614f18565b8015610f6c5780601f1061155957610100808354040283529160200191610f6c565b820191906000526020600020905b81548152906001019060200180831161156757505050919092525091949350505050565b611593611647565b610568816120bd565b6115dc6115ab82840184615846565b60408051600080825260208201909252906115d6565b60608152602001906001900390816115c15790505b506119c0565b6040805160008082526020820190925290506115ff600185858585866000611dc4565b50505050565b61160d611647565b60005b81518110156116435761163b82828151811061162e5761162e614f02565b6020026020010151612136565b600101611610565b5050565b6001546001600160a01b03163314611672576040516315ae3a6f60e11b815260040160405180910390fd5b565b60005b815181101561164357600082828151811061169457611694614f02565b60200260200101519050600081602001519050806001600160401b03166000036116d15760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b03166116f9576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b0381166000908152600860205260409020606083015160018201805461172590614f18565b905060000361178757815467ffffffffffffffff60a81b1916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a16117f0565b8154600160a81b90046001600160401b03166001148015906117c75750805160208201206040516117bc9060018501906156e2565b604051809103902014155b156117f057604051632105803760e11b81526001600160401b038416600482015260240161094d565b805115806118255750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b15611843576040516342bcdf7f60e11b815260040160405180910390fd5b6001820161185182826158ca565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b029290921674ffffffffffffffffffffffffffffffffffffffffff19909116171782556118ac60066001600160401b038516612460565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b836040516118e69190615989565b60405180910390a250505050806001019050611677565b6000610815825490565b6000610812838361246c565b6001600160401b0382166000908152600960205260408120816119376080856159d7565b6001600160401b031681526020810191909152604001600020549392505050565b467f00000000000000000000000000000000000000000000000000000000000000001461167257604051630f01ce8560e01b81527f0000000000000000000000000000000000000000000000000000000000000000600482015246602482015260440161094d565b81516000036119e25760405163c2e5347d60e01b815260040160405180910390fd5b80516040805160008082526020820190925291159181611a25565b6040805180820190915260008152606060208201528152602001906001900390816119fd5790505b50905060005b8451811015610d8357611a7b858281518110611a4957611a49614f02565b602002602001015184611a7557858381518110611a6857611a68614f02565b6020026020010151612496565b83612496565b600101611a2b565b606088516001600160401b03811115611a9e57611a9e613cda565b604051908082528060200260200182016040528015611ae357816020015b6040805180820190915260008082526020820152815260200190600190039081611abc5790505b509050811560005b8a51811015611c265781611b8357848482818110611b0b57611b0b614f02565b9050602002016020810190611b2091906159fd565b63ffffffff1615611b8357848482818110611b3d57611b3d614f02565b9050602002016020810190611b5291906159fd565b8b8281518110611b6457611b64614f02565b60200260200101516040019063ffffffff16908163ffffffff16815250505b611c018b8281518110611b9857611b98614f02565b60200260200101518b8b8b8b8b87818110611bb557611bb5614f02565b9050602002810190611bc79190615a18565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612d3192505050565b838281518110611c1357611c13614f02565b6020908102919091010152600101611aeb565b505098975050505050505050565b6000611c408383613031565b8015611c5b5750611c59836001600160e01b0319613031565b155b80156108125750610812836301ffc9a760e01b613031565b80516001600160a01b0316611c9b576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516001600160a01b039889167fffffffffffffffff0000000000000000000000000000000000000000000000009097168717600160a01b63ffffffff958616021760ff60c01b1916600160c01b911515919091021790965560608089018051600580546001600160a01b031916918b169190911790558251968752935190921693850193909352935115159183019190915251909216908201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9060800160405180910390a150565b6001600160401b03811660009081526008602052604081208054600160a01b900460ff166108155760405163ed053c5960e01b81526001600160401b038416600482015260240161094d565b60ff87811660009081526002602090815260408083208151608081018352815481526001909101548086169382019390935261010083048516918101919091526201000090910490921615156060830152873590611e23876084615a5e565b9050826060015115611e6b578451611e3c906020614fb7565b8651611e49906020614fb7565b611e549060a0615a5e565b611e5e9190615a5e565b611e689082615a5e565b90505b368114611e9457604051638e1192e160e01b81526004810182905236602482015260440161094d565b5081518114611ec35781516040516324f7d61360e21b815260048101919091526024810182905260440161094d565b611ecb611958565b60ff808a1660009081526003602090815260408083203384528252808320815180830190925280548086168352939491939092840191610100909104166002811115611f1957611f19614156565b6002811115611f2a57611f2a614156565b9052509050600281602001516002811115611f4757611f47614156565b148015611f9b5750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff1681548110611f8357611f83614f02565b6000918252602090912001546001600160a01b031633145b611fb857604051631b41e11d60e31b815260040160405180910390fd5b50816060015115612068576020820151611fd3906001615a71565b60ff16855114611ff6576040516371253a2560e01b815260040160405180910390fd5b83518551146120185760405163a75d88af60e01b815260040160405180910390fd5b6000878760405161202a929190615a8a565b604051908190038120612041918b90602001615a9a565b6040516020818303038152906040528051906020012090506120668a828888886130d6565b505b6040805182815260208a8101356001600160401b03169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b336001600160a01b038216036120e657604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff16600003612161576000604051631b3fab5160e11b815260040161094d9190615aae565b60208082015160ff808216600090815260029093526040832060018101549293909283921690036121b2576060840151600182018054911515620100000262ff0000199092169190911790556121ee565b6060840151600182015460ff62010000909104161515901515146121ee576040516321fd80df60e21b815260ff8416600482015260240161094d565b60a08401518051610100101561221a576001604051631b3fab5160e11b815260040161094d9190615aae565b805160000361223f576005604051631b3fab5160e11b815260040161094d9190615aae565b6122a5848460030180548060200260200160405190810160405280929190818152602001828054801561229b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161227d575b5050505050613289565b8460600151156123d557612313848460020180548060200260200160405190810160405280929190818152602001828054801561229b576020028201919060005260206000209081546001600160a01b0316815260019091019060200180831161227d575050505050613289565b60808501518051610100101561233f576002604051631b3fab5160e11b815260040161094d9190615aae565b604086015161234f906003615ac8565b60ff16815111612375576003604051631b3fab5160e11b815260040161094d9190615aae565b81518151101561239b576001604051631b3fab5160e11b815260040161094d9190615aae565b805160018401805461ff00191661010060ff8416021790556123c69060028601906020840190613c60565b506123d3858260016132f2565b505b6123e1848260026132f2565b80516123f69060038501906020840190613c60565b5060408581015160018401805460ff191660ff8316179055865180855560a088015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f5479361244f9389939260028a01929190615ae4565b60405180910390a1610d838461344d565b600061081283836134d0565b600082600001828154811061248357612483614f02565b9060005260206000200154905092915050565b81518151604051632cbc26bb60e01b8152608083901b67ffffffffffffffff60801b166004820152901515907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa158015612513573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061253791906156c5565b156125a857801561256657604051637edeb53960e11b81526001600160401b038316600482015260240161094d565b6040516001600160401b03831681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d8949339060200160405180910390a150505050565b60208401515160008190036125de57845160405163676cf24b60e11b81526001600160401b03909116600482015260240161094d565b8460400151518114612603576040516357e0e08360e01b815260040160405180910390fd5b6000816001600160401b0381111561261d5761261d613cda565b604051908082528060200260200182016040528015612646578160200160208202803683370190505b50905060007f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f857f000000000000000000000000000000000000000000000000000000000000000061269788611d78565b6001016040516126a791906156e2565b6040519081900381206126df949392916020019384526001600160401b03928316602085015291166040830152606082015260800190565b60405160208183030381529060405280519060200120905060005b838110156128155760008860200151828151811061271a5761271a614f02565b602002602001015190507f00000000000000000000000000000000000000000000000000000000000000006001600160401b03168160000151604001516001600160401b0316146127915780516040908101519051631c21951160e11b81526001600160401b03909116600482015260240161094d565b866001600160401b03168160000151602001516001600160401b0316146127e557805160200151604051636c95f1eb60e01b81526001600160401b03808a166004830152909116602482015260440161094d565b6127ef818461351f565b84838151811061280157612801614f02565b6020908102919091010152506001016126fa565b5050600061282d858389606001518a60800151613627565b90508060000361285b57604051633ee8bd3f60e11b81526001600160401b038616600482015260240161094d565b60005b83811015612d275760005a905060008960200151838151811061288357612883614f02565b6020026020010151905060006128a1898360000151606001516107c6565b905060008160038111156128b7576128b7614156565b14806128d4575060038160038111156128d2576128d2614156565b145b61292a57815160600151604080516001600160401b03808d16825290921660208301527f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c910160405180910390a1505050612d1f565b60608815612a09578a858151811061294457612944614f02565b6020908102919091018101510151600454909150600090600160a01b900463ffffffff166129728842614f68565b11905080806129925750600383600381111561299057612990614156565b145b6129ba576040516354e7e43160e11b81526001600160401b038c16600482015260240161094d565b8b86815181106129cc576129cc614f02565b602002602001015160000151600014612a03578b86815181106129f1576129f1614f02565b60209081029190910101515160808501525b50612a75565b6000826003811115612a1d57612a1d614156565b14612a7557825160600151604080516001600160401b03808e16825290921660208301527f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120910160405180910390a150505050612d1f565b8251608001516001600160401b031615612b4b576000826003811115612a9d57612a9d614156565b03612b4b5782516080015160208401516040516370701e5760e11b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263e0e03cae92612afb928f929190600401615b96565b6020604051808303816000875af1158015612b1a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b3e91906156c5565b612b4b5750505050612d1f565b60008c604001518681518110612b6357612b63614f02565b6020026020010151905080518460a001515114612bad57835160600151604051631cfe6d8b60e01b81526001600160401b03808e166004830152909116602482015260440161094d565b612bc18b8560000151606001516001613664565b600080612bcf868486613709565b91509150612be68d87600001516060015184613664565b8b15612c3d576003826003811115612c0057612c00614156565b03612c3d576000856003811115612c1957612c19614156565b14612c3d57855151604051632b11b8d960e01b815261094d91908390600401615bc2565b6002826003811115612c5157612c51614156565b14612c92576003826003811115612c6a57612c6a614156565b14612c92578551606001516040516349362d1f60e11b815261094d918f918590600401615bdb565b8560000151600001518660000151606001516001600160401b03168e6001600160401b03167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b8d8c81518110612cea57612cea614f02565b602002602001015186865a612cff908f614f68565b604051612d0f9493929190615c00565b60405180910390a4505050505050505b60010161285e565b5050505050505050565b6040805180820190915260008082526020820152602086015160405163bbe4f6db60e01b81526001600160a01b0380831660048301526000917f00000000000000000000000000000000000000000000000000000000000000009091169063bbe4f6db90602401602060405180830381865afa158015612db5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612dd99190615c37565b90506001600160a01b0381161580612e085750612e066001600160a01b03821663aff2afbf60e01b611c34565b155b15612e315760405163ae9b4ce960e01b81526001600160a01b038216600482015260240161094d565b600080612e4988858c6040015163ffffffff166137bd565b915091506000806000612f1a6040518061010001604052808e81526020018c6001600160401b031681526020018d6001600160a01b031681526020018f608001518152602001896001600160a01b031681526020018f6000015181526020018f6060015181526020018b815250604051602401612ec69190615c54565b60408051601f198184030181529190526020810180516001600160e01b0316633907753760e01b17905287867f000000000000000000000000000000000000000000000000000000000000000060846138c0565b92509250925082612f42578582604051634ff17cad60e11b815260040161094d929190615d20565b8151602014612f71578151604051631e3be00960e21b815260206004820152602481019190915260440161094d565b600082806020019051810190612f879190615d42565b9050866001600160a01b03168c6001600160a01b031614613003576000612fb88d8a612fb3868a614f68565b6137bd565b50905086811080612fd2575081612fcf8883614f68565b14155b156130015760405163a966e21f60e01b815260048101839052602481018890526044810182905260640161094d565b505b604080518082019091526001600160a01b039098168852602088015250949550505050505095945050505050565b604080516001600160e01b031983166024808301919091528251808303909101815260449091019091526020810180516001600160e01b03166301ffc9a760e01b17905260009081808063753fa58960e11b6179185a1015613097578060005260046000fd5b602060008651602088018b617530fa93503d925060005191508380156130be575060208310155b80156130ca5750600082115b98975050505050505050565b8251600090815b81811015612d275760006001888684602081106130fc576130fc614f02565b61310991901a601b615a71565b89858151811061311b5761311b614f02565b602002602001015189868151811061313557613135614f02565b602002602001015160405160008152602001604052604051613173949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015613195573d6000803e3d6000fd5b505060408051601f1981015160ff808e166000908152600360209081528582206001600160a01b038516835281528582208587019096528554808416865293975090955092939284019161010090041660028111156131f6576131f6614156565b600281111561320757613207614156565b905250905060018160200151600281111561322457613224614156565b1461324257604051636518c33d60e11b815260040160405180910390fd5b8051600160ff9091161b85161561326c57604051633d9ef1f160e21b815260040160405180910390fd5b806000015160ff166001901b8517945050508060010190506130dd565b60005b8151811015610aaa5760ff8316600090815260036020526040812083519091908490849081106132be576132be614f02565b6020908102919091018101516001600160a01b03168252810191909152604001600020805461ffff1916905560010161328c565b60005b82518110156115ff57600083828151811061331257613312614f02565b602002602001015190506000600281111561332f5761332f614156565b60ff80871660009081526003602090815260408083206001600160a01b0387168452909152902054610100900416600281111561336e5761336e614156565b1461338f576004604051631b3fab5160e11b815260040161094d9190615aae565b6001600160a01b0381166133b65760405163d6c62c9b60e01b815260040160405180910390fd5b60405180604001604052808360ff1681526020018460028111156133dc576133dc614156565b905260ff80871660009081526003602090815260408083206001600160a01b0387168452825290912083518154931660ff198416811782559184015190929091839161ffff19161761010083600281111561343957613439614156565b0217905550905050508060010190506132f5565b60ff818116600081815260026020526040902060010154620100009004909116906134a55780613490576040516317bd8dd160e11b815260040160405180910390fd5b600b805467ffffffffffffffff191690555050565b60001960ff831601611643578015611643576040516307b8c74d60e51b815260040160405180910390fd5b600081815260018301602052604081205461351757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610815565b506000610815565b81518051606080850151908301516080808701519401516040516000958695889561358395919490939192916020019485526001600160a01b039390931660208501526001600160401b039182166040850152606084015216608082015260a00190565b604051602081830303815290604052805190602001208560200151805190602001208660400151805190602001208760a001516040516020016135c69190615dfc565b60408051601f198184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e0015b60405160208183030381529060405280519060200120905092915050565b60008061363585858561399a565b6001600160401b0387166000908152600a6020908152604080832093835292905220549150505b949350505050565b60006002613673608085614f91565b6001600160401b03166136869190614fb7565b905060006136948585611913565b9050816136a360016004614f68565b901b1916818360038111156136ba576136ba614156565b6001600160401b03871660009081526009602052604081209190921b929092179182916136e86080886159d7565b6001600160401b031681526020810191909152604001600020555050505050565b604051630304c3e160e51b815260009060609030906360987c209061373690889088908890600401615e93565b600060405180830381600087803b15801561375057600080fd5b505af1925050508015613761575060015b6137a0573d80801561378f576040519150601f19603f3d011682016040523d82523d6000602084013e613794565b606091505b506003925090506137b5565b50506040805160208101909152600081526002905b935093915050565b600080600080600061383c886040516024016137e891906001600160a01b0391909116815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166370a0823160e01b17905288887f000000000000000000000000000000000000000000000000000000000000000060846138c0565b92509250925082613864578682604051634ff17cad60e11b815260040161094d929190615d20565b6020825114613893578151604051631e3be00960e21b815260206004820152602481019190915260440161094d565b818060200190518101906138a79190615d42565b6138b18288614f68565b94509450505050935093915050565b6000606060008361ffff166001600160401b038111156138e2576138e2613cda565b6040519080825280601f01601f19166020018201604052801561390c576020820181803683370190505b509150863b6139265763030ed58f60e21b60005260046000fd5b5a8581101561394057632be8ca8b60e21b60005260046000fd5b8590036040810481038710613960576337c3be2960e01b60005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d848111156139835750835b808352806000602085013e50955095509592505050565b82518251600091908183036139c257604051630469ac9960e21b815260040160405180910390fd5b61010182118015906139d657506101018111155b6139f3576040516309bde33960e01b815260040160405180910390fd5b60001982820101610100811115613a1d576040516309bde33960e01b815260040160405180910390fd5b80600003613a4a5786600081518110613a3857613a38614f02565b60200260200101519350505050613c18565b6000816001600160401b03811115613a6457613a64613cda565b604051908082528060200260200182016040528015613a8d578160200160208202803683370190505b50905060008080805b85811015613bb75760006001821b8b811603613af15788851015613ada578c5160018601958e918110613acb57613acb614f02565b60200260200101519050613b13565b8551600185019487918110613acb57613acb614f02565b8b5160018401938d918110613b0857613b08614f02565b602002602001015190505b600089861015613b43578d5160018701968f918110613b3457613b34614f02565b60200260200101519050613b65565b8651600186019588918110613b5a57613b5a614f02565b602002602001015190505b82851115613b86576040516309bde33960e01b815260040160405180910390fd5b613b908282613c1f565b878481518110613ba257613ba2614f02565b60209081029190910101525050600101613a96565b506001850382148015613bc957508683145b8015613bd457508581145b613bf1576040516309bde33960e01b815260040160405180910390fd5b836001860381518110613c0657613c06614f02565b60200260200101519750505050505050505b9392505050565b6000818310613c3757613c328284613c3d565b610812565b61081283835b604080516001602082015290810183905260608101829052600090608001613609565b828054828255906000526020600020908101928215613cb5579160200282015b82811115613cb557825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613c80565b50613cc1929150613cc5565b5090565b5b80821115613cc15760008155600101613cc6565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715613d1257613d12613cda565b60405290565b60405160a081016001600160401b0381118282101715613d1257613d12613cda565b60405160c081016001600160401b0381118282101715613d1257613d12613cda565b604080519081016001600160401b0381118282101715613d1257613d12613cda565b604051606081016001600160401b0381118282101715613d1257613d12613cda565b604051601f8201601f191681016001600160401b0381118282101715613dc857613dc8613cda565b604052919050565b60006001600160401b03821115613de957613de9613cda565b5060051b60200190565b6001600160a01b038116811461056857600080fd5b80356001600160401b0381168114613e1f57600080fd5b919050565b801515811461056857600080fd5b8035613e1f81613e24565b60006001600160401b03821115613e5657613e56613cda565b50601f01601f191660200190565b600082601f830112613e7557600080fd5b8135613e88613e8382613e3d565b613da0565b818152846020838601011115613e9d57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215613ecd57600080fd5b82356001600160401b0380821115613ee457600080fd5b818501915085601f830112613ef857600080fd5b8135613f06613e8382613dd0565b81815260059190911b83018401908481019088831115613f2557600080fd5b8585015b83811015613fcb57803585811115613f415760008081fd5b86016080818c03601f1901811315613f595760008081fd5b613f61613cf0565b89830135613f6e81613df3565b81526040613f7d848201613e08565b8b830152606080850135613f9081613e24565b83830152928401359289841115613fa957600091508182fd5b613fb78f8d86880101613e64565b908301525085525050918601918601613f29565b5098975050505050505050565b60005b83811015613ff3578181015183820152602001613fdb565b50506000910152565b60008151808452614014816020860160208601613fd8565b601f01601f19169290920160200192915050565b6020815260006108126020830184613ffc565b6001600160a01b0381511682526020810151151560208301526001600160401b036040820151166040830152600060608201516080606085015261365c6080850182613ffc565b604080825283519082018190526000906020906060840190828701845b828110156140c45781516001600160401b03168452928401929084019060010161409f565b50505083810382850152845180825282820190600581901b8301840187850160005b8381101561411457601f1986840301855261410283835161403b565b948701949250908601906001016140e6565b50909998505050505050505050565b6000806040838503121561413657600080fd5b61413f83613e08565b915061414d60208401613e08565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b6004811061417c5761417c614156565b9052565b60208101610815828461416c565b600060a082840312156141a057600080fd5b6141a8613d18565b9050813581526141ba60208301613e08565b60208201526141cb60408301613e08565b60408201526141dc60608301613e08565b60608201526141ed60808301613e08565b608082015292915050565b8035613e1f81613df3565b803563ffffffff81168114613e1f57600080fd5b600082601f83011261422857600080fd5b81356020614238613e8383613dd0565b82815260059290921b8401810191818101908684111561425757600080fd5b8286015b848110156143275780356001600160401b038082111561427b5760008081fd5b9088019060a0828b03601f19018113156142955760008081fd5b61429d613d18565b87840135838111156142af5760008081fd5b6142bd8d8a83880101613e64565b8252506040808501356142cf81613df3565b828a015260606142e0868201614203565b828401526080915081860135858111156142fa5760008081fd5b6143088f8c838a0101613e64565b918401919091525091909301359083015250835291830191830161425b565b509695505050505050565b6000610140828403121561434557600080fd5b61434d613d3a565b9050614359838361418e565b815260a08201356001600160401b038082111561437557600080fd5b61438185838601613e64565b602084015260c084013591508082111561439a57600080fd5b6143a685838601613e64565b60408401526143b760e085016141f8565b606084015261010084013560808401526101208401359150808211156143dc57600080fd5b506143e984828501614217565b60a08301525092915050565b600082601f83011261440657600080fd5b81356020614416613e8383613dd0565b82815260059290921b8401810191818101908684111561443557600080fd5b8286015b848110156143275780356001600160401b038111156144585760008081fd5b6144668986838b0101614332565b845250918301918301614439565b600082601f83011261448557600080fd5b81356020614495613e8383613dd0565b82815260059290921b840181019181810190868411156144b457600080fd5b8286015b848110156143275780356001600160401b03808211156144d757600080fd5b818901915089603f8301126144eb57600080fd5b858201356144fb613e8382613dd0565b81815260059190911b830160400190878101908c83111561451b57600080fd5b604085015b838110156145545780358581111561453757600080fd5b6145468f6040838a0101613e64565b845250918901918901614520565b508752505050928401925083016144b8565b600082601f83011261457757600080fd5b81356020614587613e8383613dd0565b8083825260208201915060208460051b8701019350868411156145a957600080fd5b602086015b8481101561432757803583529183019183016145ae565b600082601f8301126145d657600080fd5b813560206145e6613e8383613dd0565b82815260059290921b8401810191818101908684111561460557600080fd5b8286015b848110156143275780356001600160401b03808211156146295760008081fd5b9088019060a0828b03601f19018113156146435760008081fd5b61464b613d18565b614656888501613e08565b81526040808501358481111561466c5760008081fd5b61467a8e8b838901016143f5565b8a84015250606080860135858111156146935760008081fd5b6146a18f8c838a0101614474565b83850152506080915081860135858111156146bc5760008081fd5b6146ca8f8c838a0101614566565b9184019190915250919093013590830152508352918301918301614609565b600080604083850312156146fc57600080fd5b6001600160401b038335111561471157600080fd5b61471e84843585016145c5565b91506001600160401b036020840135111561473857600080fd5b6020830135830184601f82011261474e57600080fd5b61475b613e838235613dd0565b81358082526020808301929160051b84010187101561477957600080fd5b602083015b6020843560051b85010181101561491f576001600160401b03813511156147a457600080fd5b87603f8235860101126147b657600080fd5b6147c9613e836020833587010135613dd0565b81358501602081810135808452908301929160059190911b016040018a10156147f157600080fd5b604083358701015b83358701602081013560051b0160400181101561490f576001600160401b038135111561482557600080fd5b833587018135016040818d03603f1901121561484057600080fd5b614848613d5c565b604082013581526001600160401b036060830135111561486757600080fd5b8c605f60608401358401011261487c57600080fd5b6040606083013583010135614893613e8382613dd0565b808282526020820191508f60608460051b60608801358801010111156148b857600080fd5b6060808601358601015b60608460051b6060880135880101018110156148ef576148e181614203565b8352602092830192016148c2565b5080602085015250505080855250506020830192506020810190506147f9565b508452506020928301920161477e565b508093505050509250929050565b60008083601f84011261493f57600080fd5b5081356001600160401b0381111561495657600080fd5b6020830191508360208260051b850101111561497157600080fd5b9250929050565b60008060008060006060868803121561499057600080fd5b85356001600160401b03808211156149a757600080fd5b6149b389838a01614332565b965060208801359150808211156149c957600080fd5b6149d589838a0161492d565b909650945060408801359150808211156149ee57600080fd5b506149fb8882890161492d565b969995985093965092949392505050565b600060808284031215614a1e57600080fd5b614a26613cf0565b8235614a3181613df3565b8152614a3f60208401614203565b60208201526040830135614a5281613e24565b60408201526060830135614a6581613df3565b60608201529392505050565b600060208284031215614a8357600080fd5b81356001600160401b03811115614a9957600080fd5b820160a08185031215613c1857600080fd5b803560ff81168114613e1f57600080fd5b600060208284031215614ace57600080fd5b61081282614aab565b60008151808452602080850194506020840160005b83811015614b115781516001600160a01b031687529582019590820190600101614aec565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a0840152614b6b60e0840182614ad7565b90506040840151601f198483030160c0850152614b888282614ad7565b95945050505050565b60008060408385031215614ba457600080fd5b614bad83613e08565b946020939093013593505050565b806040810183101561081557600080fd5b60008083601f840112614bde57600080fd5b5081356001600160401b03811115614bf557600080fd5b60208301915083602082850101111561497157600080fd5b60008060008060008060008060c0898b031215614c2957600080fd5b614c338a8a614bbb565b975060408901356001600160401b0380821115614c4f57600080fd5b614c5b8c838d01614bcc565b909950975060608b0135915080821115614c7457600080fd5b614c808c838d0161492d565b909750955060808b0135915080821115614c9957600080fd5b50614ca68b828c0161492d565b999c989b50969995989497949560a00135949350505050565b600060208284031215614cd157600080fd5b61081282613e08565b602081526000610812602083018461403b565b600060208284031215614cff57600080fd5b8135613c1881613df3565b600080600060608486031215614d1f57600080fd5b614d298585614bbb565b925060408401356001600160401b03811115614d4457600080fd5b614d5086828701614bcc565b9497909650939450505050565b600082601f830112614d6e57600080fd5b81356020614d7e613e8383613dd0565b8083825260208201915060208460051b870101935086841115614da057600080fd5b602086015b84811015614327578035614db881613df3565b8352918301918301614da5565b60006020808385031215614dd857600080fd5b82356001600160401b0380821115614def57600080fd5b818501915085601f830112614e0357600080fd5b8135614e11613e8382613dd0565b81815260059190911b83018401908481019088831115614e3057600080fd5b8585015b83811015613fcb57803585811115614e4b57600080fd5b860160c0818c03601f19011215614e625760008081fd5b614e6a613d3a565b8882013581526040614e7d818401614aab565b8a8301526060614e8e818501614aab565b8284015260809150614ea1828501613e32565b9083015260a08381013589811115614eb95760008081fd5b614ec78f8d83880101614d5d565b838501525060c0840135915088821115614ee15760008081fd5b614eef8e8c84870101614d5d565b9083015250845250918601918601614e34565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680614f2c57607f821691505b602082108103614f4c57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8181038181111561081557610815614f52565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b0380841680614fab57614fab614f7b565b92169190910692915050565b808202811582820484141761081557610815614f52565b80518252600060206001600160401b0381840151168185015260408084015160a0604087015261500160a0870182613ffc565b90506060850151868203606088015261501a8282613ffc565b608087810151898303918a01919091528051808352908601935060009250908501905b8083101561506f57835180516001600160a01b031683528601518683015292850192600192909201919084019061503d565b50979650505050505050565b6020815260006108126020830184614fce565b6080815260006150a16080830187614fce565b61ffff9590951660208301525060408101929092526001600160a01b0316606090910152919050565b6000806000606084860312156150df57600080fd5b83516150ea81613e24565b60208501519093506001600160401b0381111561510657600080fd5b8401601f8101861361511757600080fd5b8051615125613e8382613e3d565b81815287602083850101111561513a57600080fd5b61514b826020830160208601613fd8565b809450505050604084015190509250925092565b80356001600160e01b0381168114613e1f57600080fd5b600082601f83011261518757600080fd5b81356020615197613e8383613dd0565b82815260069290921b840181019181810190868411156151b657600080fd5b8286015b8481101561432757604081890312156151d35760008081fd5b6151db613d5c565b6151e482613e08565b81526151f185830161515f565b818601528352918301916040016151ba565b600082601f83011261521457600080fd5b81356020615224613e8383613dd0565b82815260059290921b8401810191818101908684111561524357600080fd5b8286015b848110156143275780356001600160401b03808211156152675760008081fd5b9088019060a0828b03601f19018113156152815760008081fd5b615289613d18565b615294888501613e08565b8152604080850135848111156152aa5760008081fd5b6152b88e8b83890101613e64565b8a84015250606093506152cc848601613e08565b9082015260806152dd858201613e08565b93820193909352920135908201528352918301918301615247565b600082601f83011261530957600080fd5b81356020615319613e8383613dd0565b82815260069290921b8401810191818101908684111561533857600080fd5b8286015b8481101561432757604081890312156153555760008081fd5b61535d613d5c565b81358152848201358582015283529183019160400161533c565b6000602080838503121561538a57600080fd5b82356001600160401b03808211156153a157600080fd5b90840190606082870312156153b557600080fd5b6153bd613d7e565b8235828111156153cc57600080fd5b830160408189038113156153df57600080fd5b6153e7613d5c565b8235858111156153f657600080fd5b8301601f81018b1361540757600080fd5b8035615415613e8382613dd0565b81815260069190911b8201890190898101908d83111561543457600080fd5b928a01925b828410156154845785848f0312156154515760008081fd5b615459613d5c565b843561546481613df3565b8152615471858d0161515f565b818d0152825292850192908a0190615439565b84525050508287013591508482111561549c57600080fd5b6154a88a838501615176565b818801528352505082840135828111156154c157600080fd5b6154cd88828601615203565b858301525060408301359350818411156154e657600080fd5b6154f2878585016152f8565b60408201529695505050505050565b600082825180855260208086019550808260051b84010181860160005b8481101561559257601f19868403018952815160a06001600160401b0380835116865286830151828888015261555683880182613ffc565b6040858101518416908901526060808601519093169288019290925250608092830151929095019190915250978301979083019060010161551e565b5090979650505050505050565b6001600160a01b0384168152600060206060818401526155c26060840186615501565b83810360408581019190915285518083528387019284019060005b81811015614114578451805184528601518684015293850193918301916001016155dd565b805160408084528151848201819052600092602091908201906060870190855b8181101561565957835180516001600160a01b031684528501516001600160e01b0316858401529284019291850191600101615622565b50508583015187820388850152805180835290840192506000918401905b8083101561506f57835180516001600160401b031683528501516001600160e01b031685830152928401926001929092019190850190615677565b6020815260006108126020830184615602565b6000602082840312156156d757600080fd5b8151613c1881613e24565b60008083546156f081614f18565b60018281168015615708576001811461571d5761574c565b60ff198416875282151583028701945061574c565b8760005260208060002060005b858110156157435781548a82015290840190820161572a565b50505082870194505b50929695505050505050565b6000815461576581614f18565b808552602060018381168015615782576001811461579c576157ca565b60ff1985168884015283151560051b8801830195506157ca565b866000528260002060005b858110156157c25781548a82018601529083019084016157a7565b890184019650505b505050505092915050565b6040815260006157e86040830185613ffc565b8281036020840152614b888185615758565b6001600160401b0381811683821601908082111561581a5761581a614f52565b5092915050565b6040815260006158346040830185615501565b8281036020840152614b888185615602565b60006020828403121561585857600080fd5b81356001600160401b0381111561586e57600080fd5b61365c848285016145c5565b601f821115610aaa576000816000526020600020601f850160051c810160208610156158a35750805b601f850160051c820191505b818110156158c2578281556001016158af565b505050505050565b81516001600160401b038111156158e3576158e3613cda565b6158f7816158f18454614f18565b8461587a565b602080601f83116001811461592c57600084156159145750858301515b600019600386901b1c1916600185901b1785556158c2565b600085815260208120601f198616915b8281101561595b5788860151825594840194600190910190840161593c565b50858210156159795787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60208152600082546001600160a01b038116602084015260ff8160a01c16151560408401526001600160401b038160a81c1660608401525060808083015261081260a0830160018501615758565b60006001600160401b03808416806159f1576159f1614f7b565b92169190910492915050565b600060208284031215615a0f57600080fd5b61081282614203565b6000808335601e19843603018112615a2f57600080fd5b8301803591506001600160401b03821115615a4957600080fd5b60200191503681900382131561497157600080fd5b8082018082111561081557610815614f52565b60ff818116838216019081111561081557610815614f52565b8183823760009101908152919050565b828152604082602083013760600192915050565b6020810160068310615ac257615ac2614156565b91905290565b60ff818116838216029081169081811461581a5761581a614f52565b600060a0820160ff881683526020878185015260a0604085015281875480845260c0860191508860005282600020935060005b81811015615b3c5784546001600160a01b031683526001948501949284019201615b17565b50508481036060860152865180825290820192508187019060005b81811015615b7c5782516001600160a01b031685529383019391830191600101615b57565b50505060ff851660808501525090505b9695505050505050565b60006001600160401b03808616835280851660208401525060606040830152614b886060830184613ffc565b82815260406020820152600061365c6040830184613ffc565b6001600160401b038481168252831660208201526060810161365c604083018461416c565b848152615c10602082018561416c565b608060408201526000615c266080830185613ffc565b905082606083015295945050505050565b600060208284031215615c4957600080fd5b8151613c1881613df3565b6020815260008251610100806020850152615c73610120850183613ffc565b91506020850151615c8f60408601826001600160401b03169052565b5060408501516001600160a01b038116606086015250606085015160808501526080850151615cc960a08601826001600160a01b03169052565b5060a0850151601f19808685030160c0870152615ce68483613ffc565b935060c08701519150808685030160e0870152615d038483613ffc565b935060e0870151915080868503018387015250615b8c8382613ffc565b6001600160a01b038316815260406020820152600061365c6040830184613ffc565b600060208284031215615d5457600080fd5b5051919050565b600082825180855260208086019550808260051b84010181860160005b8481101561559257601f19868403018952815160a08151818652615d9e82870182613ffc565b9150506001600160a01b03868301511686860152604063ffffffff8184015116818701525060608083015186830382880152615dda8382613ffc565b6080948501519790940196909652505098840198925090830190600101615d78565b6020815260006108126020830184615d5b565b60008282518085526020808601955060208260051b8401016020860160005b8481101561559257601f19868403018952615e4a838351613ffc565b98840198925090830190600101615e2e565b60008151808452602080850194506020840160005b83811015614b1157815163ffffffff1687529582019590820190600101615e71565b60608152600084518051606084015260208101516001600160401b0380821660808601528060408401511660a08601528060608401511660c08601528060808401511660e0860152505050602085015161014080610100850152615efb6101a0850183613ffc565b91506040870151605f198086850301610120870152615f1a8483613ffc565b935060608901519150615f37838701836001600160a01b03169052565b608089015161016087015260a0890151925080868503016101808701525050615f608282615d5b565b9150508281036020840152615f758186615e0f565b90508281036040840152615b8c8185615e5c56fea164736f6c6343000818000a",
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
