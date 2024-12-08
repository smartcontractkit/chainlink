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
	ChainSelector      uint64
	RmnRemote          common.Address
	TokenAdminRegistry common.Address
	NonceManager       common.Address
}

var OffRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reportOnRamp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"configOnRamp\",\"type\":\"bytes\"}],\"name\":\"CommitOnRampMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyBatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientGasForCallWithExact\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenGasOverride\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionTokenGasOverride\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"messageDestChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidMessageDestChainSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidOnRampUpdate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeavesCannotBeEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionGasAmountCountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"errorReason\",\"type\":\"bytes\"}],\"name\":\"MessageValidationError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReleased\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePre\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePost\",\"type\":\"uint256\"}],\"name\":\"ReleaseOrMintBalanceMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"RootAlreadyCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationNotAllowedInExecutionPlugin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationRequiredInCommitPlugin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleCommitReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelectorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SkippedReportExecution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"}],\"name\":\"StaticConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[2]\",\"name\":\"reportContext\",\"type\":\"bytes32[2]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllSourceChainConfigs\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReport[]\",\"name\":\"reports\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isRMNVerificationDisabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101406040523480156200001257600080fd5b5060405162006cf338038062006cf3833981016040819052620000359162000886565b336000816200005757604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008a576200008a81620001ca565b505046608052600160a05260208301516001600160a01b03161580620000bb575060408301516001600160a01b0316155b80620000d2575060608301516001600160a01b0316155b15620000f1576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160401b03166000036200011d5760405163c656089560e01b815260040160405180910390fd5b82516001600160401b0390811660c052602080850180516001600160a01b0390811660e052604080880180518316610100526060808a01805185166101205283518b519098168852945184169587019590955251821690850152905116908201527f683eb52ee924eb817377cfa8f41f238f4bb7a877da5267869dfffbad85f564d89060800160405180910390a1620001b68262000244565b620001c18162000332565b50505062000c78565b336001600160a01b03821603620001f457604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03166200026d576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516001600160a01b039889166001600160c01b03199097168717600160a01b63ffffffff958616021760ff60c01b1916600160c01b911515919091021790965560608089018051600580546001600160a01b031916918b169190911790558251968752935190921693850193909352935115159183019190915251909216908201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9060800160405180910390a150565b60005b8151811015620005c7576000828281518110620003565762000356620009b0565b60200260200101519050600081602001519050806001600160401b0316600003620003945760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b0316620003bd576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b03811660009081526008602052604090206060830151600182018054620003eb90620009c6565b90506000036200044e578154600160a81b600160e81b031916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1620004bf565b8154600160a81b90046001600160401b0316600114801590620004915750805160208201206040516200048690600185019062000a02565b604051809103902014155b15620004bf57604051632105803760e11b81526001600160401b038416600482015260240160405180910390fd5b80511580620004f55750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b1562000514576040516342bcdf7f60e11b815260040160405180910390fd5b6001820162000524828262000ad5565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b02929092166001600160a81b0319909116171782556200057360066001600160401b038516620005cb565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b83604051620005af919062000ba1565b60405180910390a25050505080600101905062000335565b5050565b6000620005d98383620005e2565b90505b92915050565b60008181526001830160205260408120546200062b57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620005dc565b506000620005dc565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b03811182821017156200066f576200066f62000634565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620006a057620006a062000634565b604052919050565b80516001600160401b0381168114620006c057600080fd5b919050565b6001600160a01b0381168114620006db57600080fd5b50565b80518015158114620006c057600080fd5b6000601f83601f8401126200070357600080fd5b825160206001600160401b038083111562000722576200072262000634565b8260051b6200073383820162000675565b93845286810183019383810190898611156200074e57600080fd5b84890192505b8583101562000879578251848111156200076e5760008081fd5b89016080601f19828d038101821315620007885760008081fd5b620007926200064a565b88840151620007a181620006c5565b81526040620007b2858201620006a8565b8a8301526060620007c5818701620006de565b83830152938501519389851115620007dd5760008081fd5b84860195508f603f870112620007f557600094508485fd5b8a8601519450898511156200080e576200080e62000634565b6200081f8b858f8801160162000675565b93508484528f82868801011115620008375760008081fd5b60005b8581101562000857578681018301518582018d01528b016200083a565b5060009484018b01949094525091820152835250918401919084019062000754565b9998505050505050505050565b60008060008385036101208112156200089e57600080fd5b6080811215620008ad57600080fd5b620008b76200064a565b620008c286620006a8565b81526020860151620008d481620006c5565b60208201526040860151620008e981620006c5565b60408201526060860151620008fe81620006c5565b606082015293506080607f19820112156200091857600080fd5b50620009236200064a565b60808501516200093381620006c5565b815260a085015163ffffffff811681146200094d57600080fd5b60208201526200096060c08601620006de565b604082015260e08501516200097581620006c5565b60608201526101008501519092506001600160401b038111156200099857600080fd5b620009a686828701620006ef565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680620009db57607f821691505b602082108103620009fc57634e487b7160e01b600052602260045260246000fd5b50919050565b600080835462000a1281620009c6565b6001828116801562000a2d576001811462000a435762000a74565b60ff198416875282151583028701945062000a74565b8760005260208060002060005b8581101562000a6b5781548a82015290840190820162000a50565b50505082870194505b50929695505050505050565b601f82111562000ad0576000816000526020600020601f850160051c8101602086101562000aab5750805b601f850160051c820191505b8181101562000acc5782815560010162000ab7565b5050505b505050565b81516001600160401b0381111562000af15762000af162000634565b62000b098162000b028454620009c6565b8462000a80565b602080601f83116001811462000b41576000841562000b285750858301515b600019600386901b1c1916600185901b17855562000acc565b600085815260208120601f198616915b8281101562000b725788860151825594840194600190910190840162000b51565b508582101562000b915787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b602080825282546001600160a01b0381168383015260a081901c60ff161515604084015260a81c6001600160401b0316606083015260808083015260018084018054600093929190849062000bf681620009c6565b8060a089015260c0600183166000811462000c1a576001811462000c375762000c69565b60ff19841660c08b015260c083151560051b8b0101945062000c69565b85600052602060002060005b8481101562000c605781548c820185015290880190890162000c43565b8b0160c0019550505b50929998505050505050505050565b60805160a05160c05160e0516101005161012051615fec62000d07600039600081816102070152612a7a0152600081816101d80152612d220152600081816101a901528181610f7501528181611125015261247a01526000818161017a0152818161262501526126dc015260008181611f30015261374f01526000818161190401526119370152615fec6000f3fe608060405234801561001057600080fd5b506004361061012c5760003560e01c80637edf52f4116100ad578063de5e0b9a11610071578063de5e0b9a146104b2578063e9d68a8e146104c5578063f2fde38b146104e5578063f58e03fc146104f8578063f716f99f1461050b57600080fd5b80637edf52f41461041257806385572ffb146104255780638da5cb5b14610433578063c673e5841461044e578063ccd37ba31461046e57600080fd5b80635e36480c116100f45780635e36480c146103075780635e7bb0081461032757806360987c201461033a5780637437ff9f1461034d57806379ba50971461040a57600080fd5b806304666f9c1461013157806306285c6914610146578063181f5a771461028d5780633f4b04aa146102d65780635215505b146102f1575b600080fd5b61014461013f366004613ed9565b61051e565b005b61023760408051608081018252600080825260208201819052918101829052606081019190915260405180608001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160401b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b604051610284919081516001600160401b031681526020808301516001600160a01b0390811691830191909152604080840151821690830152606092830151169181019190915260800190565b60405180910390f35b6102c96040518060400160405280601181526020017f4f666652616d7020312e362e302d64657600000000000000000000000000000081525081565b6040516102849190614047565b600b546040516001600160401b039091168152602001610284565b6102f9610532565b6040516102849291906140a1565b61031a610315366004614142565b61078d565b604051610284919061419f565b610144610335366004614708565b6107e2565b610144610348366004614997565b610a76565b6103c360408051608081018252600080825260208201819052918101829052606081019190915250604080516080810182526004546001600160a01b038082168352600160a01b820463ffffffff166020840152600160c01b90910460ff16151592820192909252600554909116606082015290565b604051610284919081516001600160a01b03908116825260208084015163ffffffff1690830152604080840151151590830152606092830151169181019190915260800190565b610144610d33565b610144610420366004614a2b565b610db6565b61014461012c366004614a90565b6001546040516001600160a01b039091168152602001610284565b61046161045c366004614adb565b610dc7565b6040516102849190614b3b565b6104a461047c366004614bb0565b6001600160401b03919091166000908152600a60209081526040808320938352929052205490565b604051908152602001610284565b6101446104c0366004614c2c565b610f25565b6104d86104d3366004614cde565b611428565b6040516102849190614cf9565b6101446104f3366004614d0c565b611534565b610144610506366004614d29565b611545565b610144610519366004614de4565b6115ae565b6105266115f0565b61052f8161161d565b50565b606080600061054160066118a6565b6001600160401b0381111561055857610558613cf9565b6040519080825280602002602001820160405280156105a957816020015b60408051608081018252600080825260208083018290529282015260608082015282526000199092019101816105765790505b50905060006105b860066118a6565b6001600160401b038111156105cf576105cf613cf9565b6040519080825280602002602001820160405280156105f8578160200160208202803683370190505b50905060005b61060860066118a6565b8110156107845761061a6006826118b0565b82828151811061062c5761062c614f21565b60200260200101906001600160401b031690816001600160401b0316815250506008600083838151811061066257610662614f21565b6020908102919091018101516001600160401b039081168352828201939093526040918201600020825160808101845281546001600160a01b038116825260ff600160a01b820416151593820193909352600160a81b909204909316918101919091526001820180549192916060840191906106dd90614f37565b80601f016020809104026020016040519081016040528092919081815260200182805461070990614f37565b80156107565780601f1061072b57610100808354040283529160200191610756565b820191906000526020600020905b81548152906001019060200180831161073957829003601f168201915b50505050508152505083828151811061077157610771614f21565b60209081029190910101526001016105fe565b50939092509050565b600061079b60016004614f87565b60026107a8608085614fb0565b6001600160401b03166107bb9190614fd6565b6107c585856118bc565b901c1660038111156107d9576107d9614175565b90505b92915050565b6107ea611901565b81518151811461080d576040516320f8fd5960e21b815260040160405180910390fd5b60005b81811015610a6657600084828151811061082c5761082c614f21565b6020026020010151905060008160200151519050600085848151811061085457610854614f21565b602002602001015190508051821461087f576040516320f8fd5960e21b815260040160405180910390fd5b60005b82811015610a5757600082828151811061089e5761089e614f21565b60200260200101516000015190506000856020015183815181106108c4576108c4614f21565b602002602001015190508160001461091d57806080015182101561091d578551815151604051633a98d46360e11b81526001600160401b0390921660048301526024820152604481018390526064015b60405180910390fd5b83838151811061092f5761092f614f21565b602002602001015160200151518160a00151511461097c57805180516060909101516040516370a193fd60e01b815260048101929092526001600160401b03166024820152604401610914565b60005b8160a0015151811015610a495760008585815181106109a0576109a0614f21565b60200260200101516020015182815181106109bd576109bd614f21565b602002602001015163ffffffff16905080600014610a405760008360a0015183815181106109ed576109ed614f21565b60200260200101516040015163ffffffff16905080821015610a3e578351516040516348e617b360e01b81526004810191909152602481018490526044810182905260648101839052608401610914565b505b5060010161097f565b505050806001019050610882565b50505050806001019050610810565b50610a718383611969565b505050565b333014610a96576040516306e34e6560e31b815260040160405180910390fd5b6040805160008082526020820190925281610ad3565b6040805180820190915260008082526020820152815260200190600190039081610aac5790505b5060a08701515190915015610b0957610b068660a001518760200151886060015189600001516020015189898989611a2c565b90505b6040805160a081018252875151815287516020908101516001600160401b03168183015288015181830152908701516060820152608081018290526005546001600160a01b03168015610bfc576040516308d450a160e01b81526001600160a01b038216906308d450a190610b8290859060040161509a565b600060405180830381600087803b158015610b9c57600080fd5b505af1925050508015610bad575060015b610bfc573d808015610bdb576040519150601f19603f3d011682016040523d82523d6000602084013e610be0565b606091505b50806040516309c2532560e01b81526004016109149190614047565b604088015151158015610c1157506080880151155b80610c28575060608801516001600160a01b03163b155b80610c4f57506060880151610c4d906001600160a01b03166385572ffb60e01b611bdd565b155b15610c5c57505050610d2c565b87516020908101516001600160401b03166000908152600890915260408082205460808b015160608c01519251633cf9798360e01b815284936001600160a01b0390931692633cf9798392610cba92899261138892916004016150ad565b6000604051808303816000875af1158015610cd9573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610d0191908101906150e9565b509150915081610d2657806040516302a35ba360e21b81526004016109149190614047565b50505050505b5050505050565b6000546001600160a01b03163314610d5e5760405163015aa1e360e11b815260040160405180910390fd5b600180546001600160a01b0319808216339081179093556000805490911681556040516001600160a01b03909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610dbe6115f0565b61052f81611bf9565b610e0a6040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c082015294855291820180548451818402810184019095528085529293858301939092830182828015610eb357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610e95575b5050505050815260200160038201805480602002602001604051908101604052809291908181526020018280548015610f1557602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610ef7575b5050505050815250509050919050565b6000610f3387890189615396565b6004805491925090600160c01b900460ff16610fdd5760208201515115610fdd5760208201516040808401519051633854844f60e11b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016926370a9089e92610fac92309291906004016155be565b60006040518083038186803b158015610fc457600080fd5b505afa158015610fd8573d6000803e3d6000fd5b505050505b81515151151580610ff357508151602001515115155b156110be57600b5460208b0135906001600160401b038083169116101561109657600b805467ffffffffffffffff19166001600160401b03831617905581548351604051633937306f60e01b81526001600160a01b0390921691633937306f9161105f916004016156d1565b600060405180830381600087803b15801561107957600080fd5b505af115801561108d573d6000803e3d6000fd5b505050506110bc565b8260200151516000036110bc57604051632261116760e01b815260040160405180910390fd5b505b60005b826020015151811015611374576000836020015182815181106110e6576110e6614f21565b60209081029190910101518051604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b166004820152919250906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632cbc26bb90602401602060405180830381865afa15801561116c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061119091906156e4565b156111b957604051637edeb53960e11b81526001600160401b0382166004820152602401610914565b60006111c482611cfe565b9050806001016040516111d79190615701565b6040518091039020836020015180519060200120146112145782602001518160010160405163b80d8fa960e01b81526004016109149291906157f4565b60408301518154600160a81b90046001600160401b039081169116141580611255575082606001516001600160401b031683604001516001600160401b0316115b1561129a57825160408085015160608601519151636af0786b60e11b81526001600160401b039384166004820152908316602482015291166044820152606401610914565b6080830151806112bd5760405163504570e360e01b815260040160405180910390fd5b83516001600160401b03166000908152600a60209081526040808320848452909152902054156113155783516040516332cf0cbf60e01b81526001600160401b03909116600482015260248101829052604401610914565b6060840151611325906001615819565b825467ffffffffffffffff60a81b1916600160a81b6001600160401b0392831602179092559251166000908152600a6020908152604080832094835293905291909120429055506001016110c1565b50602082015182516040517f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4926113ac929091615840565b60405180910390a1610d2660008b8b8b8b8b8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808f0282810182019093528e82529093508e92508d9182918501908490808284376000920191909152508c9250611d4a915050565b60408051608080820183526000808352602080840182905283850182905260608085018190526001600160401b03878116845260088352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b9092049092169483019490945260018401805493949293918401916114b490614f37565b80601f01602080910402602001604051908101604052809291908181526020018280546114e090614f37565b8015610f155780601f1061150257610100808354040283529160200191610f15565b820191906000526020600020905b81548152906001019060200180831161151057505050919092525091949350505050565b61153c6115f0565b61052f81612073565b61158561155482840184615865565b604080516000808252602082019092529061157f565b606081526020019060019003908161156a5790505b50611969565b6040805160008082526020820190925290506115a8600185858585866000611d4a565b50505050565b6115b66115f0565b60005b81518110156115ec576115e48282815181106115d7576115d7614f21565b60200260200101516120ec565b6001016115b9565b5050565b6001546001600160a01b0316331461161b576040516315ae3a6f60e11b815260040160405180910390fd5b565b60005b81518110156115ec57600082828151811061163d5761163d614f21565b60200260200101519050600081602001519050806001600160401b031660000361167a5760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b03166116a2576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b038116600090815260086020526040902060608301516001820180546116ce90614f37565b905060000361173057815467ffffffffffffffff60a81b1916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1611799565b8154600160a81b90046001600160401b0316600114801590611770575080516020820120604051611765906001850190615701565b604051809103902014155b1561179957604051632105803760e11b81526001600160401b0384166004820152602401610914565b805115806117ce5750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b156117ec576040516342bcdf7f60e11b815260040160405180910390fd5b600182016117fa82826158e9565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b029290921674ffffffffffffffffffffffffffffffffffffffffff199091161717825561185560066001600160401b038516612416565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b8360405161188f91906159a8565b60405180910390a250505050806001019050611620565b60006107dc825490565b60006107d98383612422565b6001600160401b0382166000908152600960205260408120816118e06080856159f6565b6001600160401b031681526020810191909152604001600020549392505050565b467f00000000000000000000000000000000000000000000000000000000000000001461161b57604051630f01ce8560e01b81527f00000000000000000000000000000000000000000000000000000000000000006004820152466024820152604401610914565b815160000361198b5760405163c2e5347d60e01b815260040160405180910390fd5b805160408051600080825260208201909252911591816119ce565b6040805180820190915260008152606060208201528152602001906001900390816119a65790505b50905060005b8451811015610d2c57611a248582815181106119f2576119f2614f21565b602002602001015184611a1e57858381518110611a1157611a11614f21565b602002602001015161244c565b8361244c565b6001016119d4565b606088516001600160401b03811115611a4757611a47613cf9565b604051908082528060200260200182016040528015611a8c57816020015b6040805180820190915260008082526020820152815260200190600190039081611a655790505b509050811560005b8a51811015611bcf5781611b2c57848482818110611ab457611ab4614f21565b9050602002016020810190611ac99190615a1c565b63ffffffff1615611b2c57848482818110611ae657611ae6614f21565b9050602002016020810190611afb9190615a1c565b8b8281518110611b0d57611b0d614f21565b60200260200101516040019063ffffffff16908163ffffffff16815250505b611baa8b8281518110611b4157611b41614f21565b60200260200101518b8b8b8b8b87818110611b5e57611b5e614f21565b9050602002810190611b709190615a37565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612ce792505050565b838281518110611bbc57611bbc614f21565b6020908102919091010152600101611a94565b505098975050505050505050565b6000611be883612fc9565b80156107d957506107d98383612ffc565b80516001600160a01b0316611c21576040516342bcdf7f60e11b815260040160405180910390fd5b80516004805460208085018051604080880180516001600160a01b039889167fffffffffffffffff0000000000000000000000000000000000000000000000009097168717600160a01b63ffffffff958616021760ff60c01b1916600160c01b911515919091021790965560608089018051600580546001600160a01b031916918b169190911790558251968752935190921693850193909352935115159183019190915251909216908201527fcbb53bda7106a610de67df506ac86b65c44d5afac0fd2b11070dc2d61a6f2dee9060800160405180910390a150565b6001600160401b03811660009081526008602052604081208054600160a01b900460ff166107dc5760405163ed053c5960e01b81526001600160401b0384166004820152602401610914565b60ff87811660009081526002602090815260408083208151608081018352815481526001909101548086169382019390935261010083048516918101919091526201000090910490921615156060830152873590611da9876084615a7d565b9050826060015115611df1578451611dc2906020614fd6565b8651611dcf906020614fd6565b611dda9060a0615a7d565b611de49190615a7d565b611dee9082615a7d565b90505b368114611e1a57604051638e1192e160e01b815260048101829052366024820152604401610914565b5081518114611e495781516040516324f7d61360e21b8152600481019190915260248101829052604401610914565b611e51611901565b60ff808a1660009081526003602090815260408083203384528252808320815180830190925280548086168352939491939092840191610100909104166002811115611e9f57611e9f614175565b6002811115611eb057611eb0614175565b9052509050600281602001516002811115611ecd57611ecd614175565b148015611f215750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff1681548110611f0957611f09614f21565b6000918252602090912001546001600160a01b031633145b611f6e57336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611f6e57604051631b41e11d60e31b815260040160405180910390fd5b5081606001511561201e576020820151611f89906001615a90565b60ff16855114611fac576040516371253a2560e01b815260040160405180910390fd5b8351855114611fce5760405163a75d88af60e01b815260040160405180910390fd5b60008787604051611fe0929190615aa9565b604051908190038120611ff7918b90602001615ab9565b60405160208183030381529060405280519060200120905061201c8a82888888613086565b505b6040805182815260208a8101356001600160401b03169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b336001600160a01b0382160361209c57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff16600003612117576000604051631b3fab5160e11b81526004016109149190615acd565b60208082015160ff80821660009081526002909352604083206001810154929390928392169003612168576060840151600182018054911515620100000262ff0000199092169190911790556121a4565b6060840151600182015460ff62010000909104161515901515146121a4576040516321fd80df60e21b815260ff84166004820152602401610914565b60a0840151805161010010156121d0576001604051631b3fab5160e11b81526004016109149190615acd565b80516000036121f5576005604051631b3fab5160e11b81526004016109149190615acd565b61225b848460030180548060200260200160405190810160405280929190818152602001828054801561225157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612233575b5050505050613239565b84606001511561238b576122c98484600201805480602002602001604051908101604052809291908181526020018280548015612251576020028201919060005260206000209081546001600160a01b03168152600190910190602001808311612233575050505050613239565b6080850151805161010010156122f5576002604051631b3fab5160e11b81526004016109149190615acd565b6040860151612305906003615ae7565b60ff1681511161232b576003604051631b3fab5160e11b81526004016109149190615acd565b815181511015612351576001604051631b3fab5160e11b81526004016109149190615acd565b805160018401805461ff00191661010060ff84160217905561237c9060028601906020840190613c7f565b50612389858260016132a2565b505b612397848260026132a2565b80516123ac9060038501906020840190613c7f565b5060408581015160018401805460ff191660ff8316179055865180855560a088015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f547936124059389939260028a01929190615b03565b60405180910390a1610d2c846133fd565b60006107d98383613480565b600082600001828154811061243957612439614f21565b9060005260206000200154905092915050565b81518151604051632cbc26bb60e01b8152608083901b67ffffffffffffffff60801b166004820152901515907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa1580156124c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124ed91906156e4565b1561255e57801561251c57604051637edeb53960e11b81526001600160401b0383166004820152602401610914565b6040516001600160401b03831681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d8949339060200160405180910390a150505050565b602084015151600081900361259457845160405163676cf24b60e11b81526001600160401b039091166004820152602401610914565b84604001515181146125b9576040516357e0e08360e01b815260040160405180910390fd5b6000816001600160401b038111156125d3576125d3613cf9565b6040519080825280602002602001820160405280156125fc578160200160208202803683370190505b50905060007f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f857f000000000000000000000000000000000000000000000000000000000000000061264d88611cfe565b60010160405161265d9190615701565b604051908190038120612695949392916020019384526001600160401b03928316602085015291166040830152606082015260800190565b60405160208183030381529060405280519060200120905060005b838110156127cb576000886020015182815181106126d0576126d0614f21565b602002602001015190507f00000000000000000000000000000000000000000000000000000000000000006001600160401b03168160000151604001516001600160401b0316146127475780516040908101519051631c21951160e11b81526001600160401b039091166004820152602401610914565b866001600160401b03168160000151602001516001600160401b03161461279b57805160200151604051636c95f1eb60e01b81526001600160401b03808a1660048301529091166024820152604401610914565b6127a581846134cf565b8483815181106127b7576127b7614f21565b6020908102919091010152506001016126b0565b505060006127e3858389606001518a608001516135d7565b90508060000361281157604051633ee8bd3f60e11b81526001600160401b0386166004820152602401610914565b60005b83811015612cdd5760005a905060008960200151838151811061283957612839614f21565b6020026020010151905060006128578983600001516060015161078d565b9050600081600381111561286d5761286d614175565b148061288a5750600381600381111561288857612888614175565b145b6128e057815160600151604080516001600160401b03808d16825290921660208301527f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c910160405180910390a1505050612cd5565b606088156129bf578a85815181106128fa576128fa614f21565b6020908102919091018101510151600454909150600090600160a01b900463ffffffff166129288842614f87565b11905080806129485750600383600381111561294657612946614175565b145b612970576040516354e7e43160e11b81526001600160401b038c166004820152602401610914565b8b868151811061298257612982614f21565b6020026020010151600001516000146129b9578b86815181106129a7576129a7614f21565b60209081029190910101515160808501525b50612a2b565b60008260038111156129d3576129d3614175565b14612a2b57825160600151604080516001600160401b03808e16825290921660208301527f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120910160405180910390a150505050612cd5565b8251608001516001600160401b031615612b01576000826003811115612a5357612a53614175565b03612b015782516080015160208401516040516370701e5760e11b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263e0e03cae92612ab1928f929190600401615bb5565b6020604051808303816000875af1158015612ad0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612af491906156e4565b612b015750505050612cd5565b60008c604001518681518110612b1957612b19614f21565b6020026020010151905080518460a001515114612b6357835160600151604051631cfe6d8b60e01b81526001600160401b03808e1660048301529091166024820152604401610914565b612b778b8560000151606001516001613614565b600080612b858684866136b9565b91509150612b9c8d87600001516060015184613614565b8b15612bf3576003826003811115612bb657612bb6614175565b03612bf3576000856003811115612bcf57612bcf614175565b14612bf357855151604051632b11b8d960e01b815261091491908390600401615be1565b6002826003811115612c0757612c07614175565b14612c48576003826003811115612c2057612c20614175565b14612c48578551606001516040516349362d1f60e11b8152610914918f918590600401615bfa565b8560000151600001518660000151606001516001600160401b03168e6001600160401b03167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b8d8c81518110612ca057612ca0614f21565b602002602001015186865a612cb5908f614f87565b604051612cc59493929190615c1f565b60405180910390a4505050505050505b600101612814565b5050505050505050565b6040805180820190915260008082526020820152602086015160405163bbe4f6db60e01b81526001600160a01b0380831660048301526000917f00000000000000000000000000000000000000000000000000000000000000009091169063bbe4f6db90602401602060405180830381865afa158015612d6b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d8f9190615c56565b90506001600160a01b0381161580612dbe5750612dbc6001600160a01b03821663aff2afbf60e01b611bdd565b155b15612de75760405163ae9b4ce960e01b81526001600160a01b0382166004820152602401610914565b600080612dff88858c6040015163ffffffff166137fa565b915091506000806000612eb26040518061010001604052808e81526020018c6001600160401b031681526020018d6001600160a01b031681526020018f608001518152602001896001600160a01b031681526020018f6000015181526020018f6060015181526020018b815250604051602401612e7c9190615c73565b60408051601f198184030181529190526020810180516001600160e01b0316633907753760e01b179052878661138860846138df565b92509250925082612eda578582604051634ff17cad60e11b8152600401610914929190615d3f565b8151602014612f09578151604051631e3be00960e21b8152602060048201526024810191909152604401610914565b600082806020019051810190612f1f9190615d61565b9050866001600160a01b03168c6001600160a01b031614612f9b576000612f508d8a612f4b868a614f87565b6137fa565b50905086811080612f6a575081612f678883614f87565b14155b15612f995760405163a966e21f60e01b8152600481018390526024810188905260448101829052606401610914565b505b604080518082019091526001600160a01b039098168852602088015250949550505050505095945050505050565b6000612fdc826301ffc9a760e01b612ffc565b80156107dc5750612ff5826001600160e01b0319612ffc565b1592915050565b6040516001600160e01b031982166024820152600090819060440160408051601f19818403018152919052602080820180516001600160e01b03166301ffc9a760e01b178152825192935060009283928392909183918a617530fa92503d9150600051905082801561306f575060208210155b801561307b5750600081115b979650505050505050565b8251600090815b81811015612cdd5760006001888684602081106130ac576130ac614f21565b6130b991901a601b615a90565b8985815181106130cb576130cb614f21565b60200260200101518986815181106130e5576130e5614f21565b602002602001015160405160008152602001604052604051613123949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015613145573d6000803e3d6000fd5b505060408051601f1981015160ff808e166000908152600360209081528582206001600160a01b038516835281528582208587019096528554808416865293975090955092939284019161010090041660028111156131a6576131a6614175565b60028111156131b7576131b7614175565b90525090506001816020015160028111156131d4576131d4614175565b146131f257604051636518c33d60e11b815260040160405180910390fd5b8051600160ff9091161b85161561321c57604051633d9ef1f160e21b815260040160405180910390fd5b806000015160ff166001901b85179450505080600101905061308d565b60005b8151811015610a715760ff83166000908152600360205260408120835190919084908490811061326e5761326e614f21565b6020908102919091018101516001600160a01b03168252810191909152604001600020805461ffff1916905560010161323c565b60005b82518110156115a85760008382815181106132c2576132c2614f21565b60200260200101519050600060028111156132df576132df614175565b60ff80871660009081526003602090815260408083206001600160a01b0387168452909152902054610100900416600281111561331e5761331e614175565b1461333f576004604051631b3fab5160e11b81526004016109149190615acd565b6001600160a01b0381166133665760405163d6c62c9b60e01b815260040160405180910390fd5b60405180604001604052808360ff16815260200184600281111561338c5761338c614175565b905260ff80871660009081526003602090815260408083206001600160a01b0387168452825290912083518154931660ff198416811782559184015190929091839161ffff1916176101008360028111156133e9576133e9614175565b0217905550905050508060010190506132a5565b60ff818116600081815260026020526040902060010154620100009004909116906134555780613440576040516317bd8dd160e11b815260040160405180910390fd5b600b805467ffffffffffffffff191690555050565b60001960ff8316016115ec5780156115ec576040516307b8c74d60e51b815260040160405180910390fd5b60008181526001830160205260408120546134c7575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556107dc565b5060006107dc565b81518051606080850151908301516080808701519401516040516000958695889561353395919490939192916020019485526001600160a01b039390931660208501526001600160401b039182166040850152606084015216608082015260a00190565b604051602081830303815290604052805190602001208560200151805190602001208660400151805190602001208760a001516040516020016135769190615e1b565b60408051601f198184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e0015b60405160208183030381529060405280519060200120905092915050565b6000806135e58585856139b9565b6001600160401b0387166000908152600a6020908152604080832093835292905220549150505b949350505050565b60006002613623608085614fb0565b6001600160401b03166136369190614fd6565b9050600061364485856118bc565b90508161365360016004614f87565b901b19168183600381111561366a5761366a614175565b6001600160401b03871660009081526009602052604081209190921b929092179182916136986080886159f6565b6001600160401b031681526020810191909152604001600020555050505050565b604051630304c3e160e51b815260009060609030906360987c20906136e690889088908890600401615eb2565b600060405180830381600087803b15801561370057600080fd5b505af1925050508015613711575060015b6137dd573d80801561373f576040519150601f19603f3d011682016040523d82523d6000602084013e613744565b606091505b506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001633036137d25761377e81615fa8565b6001600160e01b0319166337c3be2960e01b14806137b457506137a081615fa8565b6001600160e01b031916632be8ca8b60e21b145b156137d257604051631d1ccf9f60e01b815260040160405180910390fd5b6003925090506137f2565b50506040805160208101909152600081526002905b935093915050565b600080600080600061385b8860405160240161382591906001600160a01b0391909116815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166370a0823160e01b179052888861138860846138df565b92509250925082613883578682604051634ff17cad60e11b8152600401610914929190615d3f565b60208251146138b2578151604051631e3be00960e21b8152602060048201526024810191909152604401610914565b818060200190518101906138c69190615d61565b6138d08288614f87565b94509450505050935093915050565b6000606060008361ffff166001600160401b0381111561390157613901613cf9565b6040519080825280601f01601f19166020018201604052801561392b576020820181803683370190505b509150863b6139455763030ed58f60e21b60005260046000fd5b5a8581101561395f57632be8ca8b60e21b60005260046000fd5b859003604081048103871061397f576337c3be2960e01b60005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d848111156139a25750835b808352806000602085013e50955095509592505050565b82518251600091908183036139e157604051630469ac9960e21b815260040160405180910390fd5b61010182118015906139f557506101018111155b613a12576040516309bde33960e01b815260040160405180910390fd5b60001982820101610100811115613a3c576040516309bde33960e01b815260040160405180910390fd5b80600003613a695786600081518110613a5757613a57614f21565b60200260200101519350505050613c37565b6000816001600160401b03811115613a8357613a83613cf9565b604051908082528060200260200182016040528015613aac578160200160208202803683370190505b50905060008080805b85811015613bd65760006001821b8b811603613b105788851015613af9578c5160018601958e918110613aea57613aea614f21565b60200260200101519050613b32565b8551600185019487918110613aea57613aea614f21565b8b5160018401938d918110613b2757613b27614f21565b602002602001015190505b600089861015613b62578d5160018701968f918110613b5357613b53614f21565b60200260200101519050613b84565b8651600186019588918110613b7957613b79614f21565b602002602001015190505b82851115613ba5576040516309bde33960e01b815260040160405180910390fd5b613baf8282613c3e565b878481518110613bc157613bc1614f21565b60209081029190910101525050600101613ab5565b506001850382148015613be857508683145b8015613bf357508581145b613c10576040516309bde33960e01b815260040160405180910390fd5b836001860381518110613c2557613c25614f21565b60200260200101519750505050505050505b9392505050565b6000818310613c5657613c518284613c5c565b6107d9565b6107d983835b6040805160016020820152908101839052606081018290526000906080016135b9565b828054828255906000526020600020908101928215613cd4579160200282015b82811115613cd457825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613c9f565b50613ce0929150613ce4565b5090565b5b80821115613ce05760008155600101613ce5565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715613d3157613d31613cf9565b60405290565b60405160a081016001600160401b0381118282101715613d3157613d31613cf9565b60405160c081016001600160401b0381118282101715613d3157613d31613cf9565b604080519081016001600160401b0381118282101715613d3157613d31613cf9565b604051606081016001600160401b0381118282101715613d3157613d31613cf9565b604051601f8201601f191681016001600160401b0381118282101715613de757613de7613cf9565b604052919050565b60006001600160401b03821115613e0857613e08613cf9565b5060051b60200190565b6001600160a01b038116811461052f57600080fd5b80356001600160401b0381168114613e3e57600080fd5b919050565b801515811461052f57600080fd5b8035613e3e81613e43565b60006001600160401b03821115613e7557613e75613cf9565b50601f01601f191660200190565b600082601f830112613e9457600080fd5b8135613ea7613ea282613e5c565b613dbf565b818152846020838601011115613ebc57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215613eec57600080fd5b82356001600160401b0380821115613f0357600080fd5b818501915085601f830112613f1757600080fd5b8135613f25613ea282613def565b81815260059190911b83018401908481019088831115613f4457600080fd5b8585015b83811015613fea57803585811115613f605760008081fd5b86016080818c03601f1901811315613f785760008081fd5b613f80613d0f565b89830135613f8d81613e12565b81526040613f9c848201613e27565b8b830152606080850135613faf81613e43565b83830152928401359289841115613fc857600091508182fd5b613fd68f8d86880101613e83565b908301525085525050918601918601613f48565b5098975050505050505050565b60005b83811015614012578181015183820152602001613ffa565b50506000910152565b60008151808452614033816020860160208601613ff7565b601f01601f19169290920160200192915050565b6020815260006107d9602083018461401b565b6001600160a01b0381511682526020810151151560208301526001600160401b036040820151166040830152600060608201516080606085015261360c608085018261401b565b604080825283519082018190526000906020906060840190828701845b828110156140e35781516001600160401b0316845292840192908401906001016140be565b50505083810382850152845180825282820190600581901b8301840187850160005b8381101561413357601f1986840301855261412183835161405a565b94870194925090860190600101614105565b50909998505050505050505050565b6000806040838503121561415557600080fd5b61415e83613e27565b915061416c60208401613e27565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b6004811061419b5761419b614175565b9052565b602081016107dc828461418b565b600060a082840312156141bf57600080fd5b6141c7613d37565b9050813581526141d960208301613e27565b60208201526141ea60408301613e27565b60408201526141fb60608301613e27565b606082015261420c60808301613e27565b608082015292915050565b8035613e3e81613e12565b803563ffffffff81168114613e3e57600080fd5b600082601f83011261424757600080fd5b81356020614257613ea283613def565b82815260059290921b8401810191818101908684111561427657600080fd5b8286015b848110156143465780356001600160401b038082111561429a5760008081fd5b9088019060a0828b03601f19018113156142b45760008081fd5b6142bc613d37565b87840135838111156142ce5760008081fd5b6142dc8d8a83880101613e83565b8252506040808501356142ee81613e12565b828a015260606142ff868201614222565b828401526080915081860135858111156143195760008081fd5b6143278f8c838a0101613e83565b918401919091525091909301359083015250835291830191830161427a565b509695505050505050565b6000610140828403121561436457600080fd5b61436c613d59565b905061437883836141ad565b815260a08201356001600160401b038082111561439457600080fd5b6143a085838601613e83565b602084015260c08401359150808211156143b957600080fd5b6143c585838601613e83565b60408401526143d660e08501614217565b606084015261010084013560808401526101208401359150808211156143fb57600080fd5b5061440884828501614236565b60a08301525092915050565b600082601f83011261442557600080fd5b81356020614435613ea283613def565b82815260059290921b8401810191818101908684111561445457600080fd5b8286015b848110156143465780356001600160401b038111156144775760008081fd5b6144858986838b0101614351565b845250918301918301614458565b600082601f8301126144a457600080fd5b813560206144b4613ea283613def565b82815260059290921b840181019181810190868411156144d357600080fd5b8286015b848110156143465780356001600160401b03808211156144f657600080fd5b818901915089603f83011261450a57600080fd5b8582013561451a613ea282613def565b81815260059190911b830160400190878101908c83111561453a57600080fd5b604085015b838110156145735780358581111561455657600080fd5b6145658f6040838a0101613e83565b84525091890191890161453f565b508752505050928401925083016144d7565b600082601f83011261459657600080fd5b813560206145a6613ea283613def565b8083825260208201915060208460051b8701019350868411156145c857600080fd5b602086015b8481101561434657803583529183019183016145cd565b600082601f8301126145f557600080fd5b81356020614605613ea283613def565b82815260059290921b8401810191818101908684111561462457600080fd5b8286015b848110156143465780356001600160401b03808211156146485760008081fd5b9088019060a0828b03601f19018113156146625760008081fd5b61466a613d37565b614675888501613e27565b81526040808501358481111561468b5760008081fd5b6146998e8b83890101614414565b8a84015250606080860135858111156146b25760008081fd5b6146c08f8c838a0101614493565b83850152506080915081860135858111156146db5760008081fd5b6146e98f8c838a0101614585565b9184019190915250919093013590830152508352918301918301614628565b6000806040838503121561471b57600080fd5b6001600160401b038335111561473057600080fd5b61473d84843585016145e4565b91506001600160401b036020840135111561475757600080fd5b6020830135830184601f82011261476d57600080fd5b61477a613ea28235613def565b81358082526020808301929160051b84010187101561479857600080fd5b602083015b6020843560051b85010181101561493e576001600160401b03813511156147c357600080fd5b87603f8235860101126147d557600080fd5b6147e8613ea26020833587010135613def565b81358501602081810135808452908301929160059190911b016040018a101561481057600080fd5b604083358701015b83358701602081013560051b0160400181101561492e576001600160401b038135111561484457600080fd5b833587018135016040818d03603f1901121561485f57600080fd5b614867613d7b565b604082013581526001600160401b036060830135111561488657600080fd5b8c605f60608401358401011261489b57600080fd5b60406060830135830101356148b2613ea282613def565b808282526020820191508f60608460051b60608801358801010111156148d757600080fd5b6060808601358601015b60608460051b60608801358801010181101561490e5761490081614222565b8352602092830192016148e1565b508060208501525050508085525050602083019250602081019050614818565b508452506020928301920161479d565b508093505050509250929050565b60008083601f84011261495e57600080fd5b5081356001600160401b0381111561497557600080fd5b6020830191508360208260051b850101111561499057600080fd5b9250929050565b6000806000806000606086880312156149af57600080fd5b85356001600160401b03808211156149c657600080fd5b6149d289838a01614351565b965060208801359150808211156149e857600080fd5b6149f489838a0161494c565b90965094506040880135915080821115614a0d57600080fd5b50614a1a8882890161494c565b969995985093965092949392505050565b600060808284031215614a3d57600080fd5b614a45613d0f565b8235614a5081613e12565b8152614a5e60208401614222565b60208201526040830135614a7181613e43565b60408201526060830135614a8481613e12565b60608201529392505050565b600060208284031215614aa257600080fd5b81356001600160401b03811115614ab857600080fd5b820160a08185031215613c3757600080fd5b803560ff81168114613e3e57600080fd5b600060208284031215614aed57600080fd5b6107d982614aca565b60008151808452602080850194506020840160005b83811015614b305781516001600160a01b031687529582019590820190600101614b0b565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a0840152614b8a60e0840182614af6565b90506040840151601f198483030160c0850152614ba78282614af6565b95945050505050565b60008060408385031215614bc357600080fd5b614bcc83613e27565b946020939093013593505050565b80604081018310156107dc57600080fd5b60008083601f840112614bfd57600080fd5b5081356001600160401b03811115614c1457600080fd5b60208301915083602082850101111561499057600080fd5b60008060008060008060008060c0898b031215614c4857600080fd5b614c528a8a614bda565b975060408901356001600160401b0380821115614c6e57600080fd5b614c7a8c838d01614beb565b909950975060608b0135915080821115614c9357600080fd5b614c9f8c838d0161494c565b909750955060808b0135915080821115614cb857600080fd5b50614cc58b828c0161494c565b999c989b50969995989497949560a00135949350505050565b600060208284031215614cf057600080fd5b6107d982613e27565b6020815260006107d9602083018461405a565b600060208284031215614d1e57600080fd5b8135613c3781613e12565b600080600060608486031215614d3e57600080fd5b614d488585614bda565b925060408401356001600160401b03811115614d6357600080fd5b614d6f86828701614beb565b9497909650939450505050565b600082601f830112614d8d57600080fd5b81356020614d9d613ea283613def565b8083825260208201915060208460051b870101935086841115614dbf57600080fd5b602086015b84811015614346578035614dd781613e12565b8352918301918301614dc4565b60006020808385031215614df757600080fd5b82356001600160401b0380821115614e0e57600080fd5b818501915085601f830112614e2257600080fd5b8135614e30613ea282613def565b81815260059190911b83018401908481019088831115614e4f57600080fd5b8585015b83811015613fea57803585811115614e6a57600080fd5b860160c0818c03601f19011215614e815760008081fd5b614e89613d59565b8882013581526040614e9c818401614aca565b8a8301526060614ead818501614aca565b8284015260809150614ec0828501613e51565b9083015260a08381013589811115614ed85760008081fd5b614ee68f8d83880101614d7c565b838501525060c0840135915088821115614f005760008081fd5b614f0e8e8c84870101614d7c565b9083015250845250918601918601614e53565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680614f4b57607f821691505b602082108103614f6b57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b818103818111156107dc576107dc614f71565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b0380841680614fca57614fca614f9a565b92169190910692915050565b80820281158282048414176107dc576107dc614f71565b80518252600060206001600160401b0381840151168185015260408084015160a0604087015261502060a087018261401b565b905060608501518682036060880152615039828261401b565b608087810151898303918a01919091528051808352908601935060009250908501905b8083101561508e57835180516001600160a01b031683528601518683015292850192600192909201919084019061505c565b50979650505050505050565b6020815260006107d96020830184614fed565b6080815260006150c06080830187614fed565b61ffff9590951660208301525060408101929092526001600160a01b0316606090910152919050565b6000806000606084860312156150fe57600080fd5b835161510981613e43565b60208501519093506001600160401b0381111561512557600080fd5b8401601f8101861361513657600080fd5b8051615144613ea282613e5c565b81815287602083850101111561515957600080fd5b61516a826020830160208601613ff7565b809450505050604084015190509250925092565b80356001600160e01b0381168114613e3e57600080fd5b600082601f8301126151a657600080fd5b813560206151b6613ea283613def565b82815260069290921b840181019181810190868411156151d557600080fd5b8286015b8481101561434657604081890312156151f25760008081fd5b6151fa613d7b565b61520382613e27565b815261521085830161517e565b818601528352918301916040016151d9565b600082601f83011261523357600080fd5b81356020615243613ea283613def565b82815260059290921b8401810191818101908684111561526257600080fd5b8286015b848110156143465780356001600160401b03808211156152865760008081fd5b9088019060a0828b03601f19018113156152a05760008081fd5b6152a8613d37565b6152b3888501613e27565b8152604080850135848111156152c95760008081fd5b6152d78e8b83890101613e83565b8a84015250606093506152eb848601613e27565b9082015260806152fc858201613e27565b93820193909352920135908201528352918301918301615266565b600082601f83011261532857600080fd5b81356020615338613ea283613def565b82815260069290921b8401810191818101908684111561535757600080fd5b8286015b8481101561434657604081890312156153745760008081fd5b61537c613d7b565b81358152848201358582015283529183019160400161535b565b600060208083850312156153a957600080fd5b82356001600160401b03808211156153c057600080fd5b90840190606082870312156153d457600080fd5b6153dc613d9d565b8235828111156153eb57600080fd5b830160408189038113156153fe57600080fd5b615406613d7b565b82358581111561541557600080fd5b8301601f81018b1361542657600080fd5b8035615434613ea282613def565b81815260069190911b8201890190898101908d83111561545357600080fd5b928a01925b828410156154a35785848f0312156154705760008081fd5b615478613d7b565b843561548381613e12565b8152615490858d0161517e565b818d0152825292850192908a0190615458565b8452505050828701359150848211156154bb57600080fd5b6154c78a838501615195565b818801528352505082840135828111156154e057600080fd5b6154ec88828601615222565b8583015250604083013593508184111561550557600080fd5b61551187858501615317565b60408201529695505050505050565b600082825180855260208086019550808260051b84010181860160005b848110156155b157601f19868403018952815160a06001600160401b038083511686528683015182888801526155758388018261401b565b6040858101518416908901526060808601519093169288019290925250608092830151929095019190915250978301979083019060010161553d565b5090979650505050505050565b6001600160a01b0384168152600060206060818401526155e16060840186615520565b83810360408581019190915285518083528387019284019060005b81811015614133578451805184528601518684015293850193918301916001016155fc565b805160408084528151848201819052600092602091908201906060870190855b8181101561567857835180516001600160a01b031684528501516001600160e01b0316858401529284019291850191600101615641565b50508583015187820388850152805180835290840192506000918401905b8083101561508e57835180516001600160401b031683528501516001600160e01b031685830152928401926001929092019190850190615696565b6020815260006107d96020830184615621565b6000602082840312156156f657600080fd5b8151613c3781613e43565b600080835461570f81614f37565b60018281168015615727576001811461573c5761576b565b60ff198416875282151583028701945061576b565b8760005260208060002060005b858110156157625781548a820152908401908201615749565b50505082870194505b50929695505050505050565b6000815461578481614f37565b8085526020600183811680156157a157600181146157bb576157e9565b60ff1985168884015283151560051b8801830195506157e9565b866000528260002060005b858110156157e15781548a82018601529083019084016157c6565b890184019650505b505050505092915050565b604081526000615807604083018561401b565b8281036020840152614ba78185615777565b6001600160401b0381811683821601908082111561583957615839614f71565b5092915050565b6040815260006158536040830185615520565b8281036020840152614ba78185615621565b60006020828403121561587757600080fd5b81356001600160401b0381111561588d57600080fd5b61360c848285016145e4565b601f821115610a71576000816000526020600020601f850160051c810160208610156158c25750805b601f850160051c820191505b818110156158e1578281556001016158ce565b505050505050565b81516001600160401b0381111561590257615902613cf9565b615916816159108454614f37565b84615899565b602080601f83116001811461594b57600084156159335750858301515b600019600386901b1c1916600185901b1785556158e1565b600085815260208120601f198616915b8281101561597a5788860151825594840194600190910190840161595b565b50858210156159985787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60208152600082546001600160a01b038116602084015260ff8160a01c16151560408401526001600160401b038160a81c166060840152506080808301526107d960a0830160018501615777565b60006001600160401b0380841680615a1057615a10614f9a565b92169190910492915050565b600060208284031215615a2e57600080fd5b6107d982614222565b6000808335601e19843603018112615a4e57600080fd5b8301803591506001600160401b03821115615a6857600080fd5b60200191503681900382131561499057600080fd5b808201808211156107dc576107dc614f71565b60ff81811683821601908111156107dc576107dc614f71565b8183823760009101908152919050565b828152604082602083013760600192915050565b6020810160068310615ae157615ae1614175565b91905290565b60ff818116838216029081169081811461583957615839614f71565b600060a0820160ff881683526020878185015260a0604085015281875480845260c0860191508860005282600020935060005b81811015615b5b5784546001600160a01b031683526001948501949284019201615b36565b50508481036060860152865180825290820192508187019060005b81811015615b9b5782516001600160a01b031685529383019391830191600101615b76565b50505060ff851660808501525090505b9695505050505050565b60006001600160401b03808616835280851660208401525060606040830152614ba7606083018461401b565b82815260406020820152600061360c604083018461401b565b6001600160401b038481168252831660208201526060810161360c604083018461418b565b848152615c2f602082018561418b565b608060408201526000615c45608083018561401b565b905082606083015295945050505050565b600060208284031215615c6857600080fd5b8151613c3781613e12565b6020815260008251610100806020850152615c9261012085018361401b565b91506020850151615cae60408601826001600160401b03169052565b5060408501516001600160a01b038116606086015250606085015160808501526080850151615ce860a08601826001600160a01b03169052565b5060a0850151601f19808685030160c0870152615d05848361401b565b935060c08701519150808685030160e0870152615d22848361401b565b935060e0870151915080868503018387015250615bab838261401b565b6001600160a01b038316815260406020820152600061360c604083018461401b565b600060208284031215615d7357600080fd5b5051919050565b600082825180855260208086019550808260051b84010181860160005b848110156155b157601f19868403018952815160a08151818652615dbd8287018261401b565b9150506001600160a01b03868301511686860152604063ffffffff8184015116818701525060608083015186830382880152615df9838261401b565b6080948501519790940196909652505098840198925090830190600101615d97565b6020815260006107d96020830184615d7a565b60008282518085526020808601955060208260051b8401016020860160005b848110156155b157601f19868403018952615e6983835161401b565b98840198925090830190600101615e4d565b60008151808452602080850194506020840160005b83811015614b3057815163ffffffff1687529582019590820190600101615e90565b60608152600084518051606084015260208101516001600160401b0380821660808601528060408401511660a08601528060608401511660c08601528060808401511660e0860152505050602085015161014080610100850152615f1a6101a085018361401b565b91506040870151605f198086850301610120870152615f39848361401b565b935060608901519150615f56838701836001600160a01b03169052565b608089015161016087015260a0890151925080868503016101808701525050615f7f8282615d7a565b9150508281036020840152615f948186615e2e565b90508281036040840152615bab8185615e7b565b805160208201516001600160e01b03198082169291906004831015615fd75780818460040360031b1b83161693505b50505091905056fea164736f6c6343000818000a",
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
