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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"reportOnRamp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"configOnRamp\",\"type\":\"bytes\"}],\"name\":\"CommitOnRampMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyBatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumMultiOCR3Base.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"min\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"max\",\"type\":\"uint64\"}],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenGasOverride\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionTokenGasOverride\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"messageDestChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidMessageDestChainSelector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidOnRampUpdate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LeavesCannotBeEmpty\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionGasAmountCountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"errorReason\",\"type\":\"bytes\"}],\"name\":\"MessageValidationError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonUniqueSignatures\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReleased\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePre\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePost\",\"type\":\"uint256\"}],\"name\":\"ReleaseOrMintBalanceMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"RootAlreadyCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureVerificationDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesOutOfRegistration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainNotEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportSourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"messageSourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StaleCommitReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"StaticConfigCannotBeChanged\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongNumberOfSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroChainSelectorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRampAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structInternal.MerkleRoot[]\",\"name\":\"merkleRoots\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"uint224\",\"name\":\"usdPerToken\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.TokenPriceUpdate[]\",\"name\":\"tokenPriceUpdates\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint224\",\"name\":\"usdPerUnitGas\",\"type\":\"uint224\"}],\"internalType\":\"structInternal.GasPriceUpdate[]\",\"name\":\"gasPriceUpdates\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structInternal.PriceUpdates\",\"name\":\"priceUpdates\",\"type\":\"tuple\"}],\"name\":\"CommitReportAccepted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"DynamicConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SkippedReportExecution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"}],\"name\":\"StaticConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllSourceChainConfigs\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPriceSequenceNumber\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"contractIRouter\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"minSeqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"onRamp\",\"type\":\"bytes\"}],\"internalType\":\"structOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"contractIRMNRemote\",\"name\":\"rmnRemote\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nonceManager\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"}],\"name\":\"latestConfigDetails\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"n\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"}],\"internalType\":\"structMultiOCR3Base.ConfigInfo\",\"name\":\"configInfo\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfig\",\"name\":\"ocrConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structInternal.RampMessageHeader\",\"name\":\"header\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"destTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"destGasAmount\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.Any2EVMTokenTransfer[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structInternal.Any2EVMRampMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReport[]\",\"name\":\"reports\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"internalType\":\"structOffRamp.GasLimitOverride[][]\",\"name\":\"gasLimitOverrides\",\"type\":\"tuple[][]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"feeQuoter\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"messageInterceptor\",\"type\":\"address\"}],\"internalType\":\"structOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"setDynamicConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"ocrPluginType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"F\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isSignatureVerificationEnabled\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"}],\"internalType\":\"structMultiOCR3Base.OCRConfigArgs[]\",\"name\":\"ocrConfigArgs\",\"type\":\"tuple[]\"}],\"name\":\"setOCR3Configs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162006c3338038062006c33833981016040819052620000359162000885565b33806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620001fa565b5050466080525060208301516001600160a01b03161580620000ec575060408301516001600160a01b0316155b8062000103575060608301516001600160a01b0316155b1562000122576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160401b03166000036200014e5760405163c656089560e01b815260040160405180910390fd5b82516001600160401b0390811660a052602080850180516001600160a01b0390811660c05260408088018051831660e0526060808a01805185166101005283518b519098168852945184169587019590955251821690850152905116908201527f683eb52ee924eb817377cfa8f41f238f4bb7a877da5267869dfffbad85f564d89060800160405180910390a1620001e682620002a5565b620001f1816200036d565b50505062000c0c565b336001600160a01b03821603620002545760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0316620002ce576040516342bcdf7f60e11b815260040160405180910390fd5b805160048054602080850180516001600160a01b039586166001600160c01b03199094168417600160a01b63ffffffff928316021790945560408087018051600580546001600160a01b031916918916919091179055815194855291519094169183019190915251909216908201527fa1c15688cb2c24508e158f6942b9276c6f3028a85e1af8cf3fff0c3ff3d5fc8d9060600160405180910390a150565b60005b8151811015620005d0576000828281518110620003915762000391620009c2565b60200260200101519050600081602001519050806001600160401b0316600003620003cf5760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b0316620003f8576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b038116600090815260086020526040902060608301516001820180546200042690620009d8565b905060000362000489578154600160a81b600160e81b031916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a1620004c8565b8154600160a81b90046001600160401b0316600114620004c857604051632105803760e11b81526001600160401b038416600482015260240162000083565b80511580620004fe5750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b156200051d576040516342bcdf7f60e11b815260040160405180910390fd5b600182016200052d828262000a69565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b02929092166001600160a81b0319909116171782556200057c60066001600160401b038516620005d4565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b83604051620005b8919062000b35565b60405180910390a25050505080600101905062000370565b5050565b6000620005e28383620005eb565b90505b92915050565b60008181526001830160205260408120546200063457508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620005e5565b506000620005e5565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b03811182821017156200067857620006786200063d565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620006a957620006a96200063d565b604052919050565b80516001600160401b0381168114620006c957600080fd5b919050565b6001600160a01b0381168114620006e457600080fd5b50565b6000601f83601f840112620006fb57600080fd5b825160206001600160401b03808311156200071a576200071a6200063d565b8260051b6200072b8382016200067e565b93845286810183019383810190898611156200074657600080fd5b84890192505b858310156200087857825184811115620007665760008081fd5b89016080601f19828d038101821315620007805760008081fd5b6200078a62000653565b888401516200079981620006ce565b81526040620007aa858201620006b1565b8a8301526060808601518015158114620007c45760008081fd5b83830152938501519389851115620007dc5760008081fd5b84860195508f603f870112620007f457600094508485fd5b8a8601519450898511156200080d576200080d6200063d565b6200081e8b858f880116016200067e565b93508484528f82868801011115620008365760008081fd5b60005b8581101562000856578681018301518582018d01528b0162000839565b5060009484018b0194909452509182015283525091840191908401906200074c565b9998505050505050505050565b60008060008385036101008112156200089d57600080fd5b6080811215620008ac57600080fd5b620008b662000653565b620008c186620006b1565b81526020860151620008d381620006ce565b60208201526040860151620008e881620006ce565b60408201526060860151620008fd81620006ce565b606082810191909152909450607f19820112156200091a57600080fd5b50604051606081016001600160401b0380821183831017156200094157620009416200063d565b81604052608087015191506200095782620006ce565b90825260a08601519063ffffffff821682146200097357600080fd5b81602084015260c087015191506200098b82620006ce565b6040830182905260e087015192945080831115620009a857600080fd5b5050620009b886828701620006e7565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680620009ed57607f821691505b60208210810362000a0e57634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111562000a64576000816000526020600020601f850160051c8101602086101562000a3f5750805b601f850160051c820191505b8181101562000a605782815560010162000a4b565b5050505b505050565b81516001600160401b0381111562000a855762000a856200063d565b62000a9d8162000a968454620009d8565b8462000a14565b602080601f83116001811462000ad5576000841562000abc5750858301515b600019600386901b1c1916600185901b17855562000a60565b600085815260208120601f198616915b8281101562000b065788860151825594840194600190910190840162000ae5565b508582101562000b255787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b602080825282546001600160a01b0381168383015260a081901c60ff161515604084015260a81c6001600160401b0316606083015260808083015260018084018054600093929190849062000b8a81620009d8565b8060a089015260c0600183166000811462000bae576001811462000bcb5762000bfd565b60ff19841660c08b015260c083151560051b8b0101945062000bfd565b85600052602060002060005b8481101562000bf45781548c820185015290880190890162000bd7565b8b0160c0019550505b50929998505050505050505050565b60805160a05160c05160e05161010051615fb162000c82600039600081816102470152612c840152600081816102180152612f730152600081816101e9015281816105890152818161073b01526126360152600081816101ba0152612884015260008181611d5b0152611d8e0152615fb16000f3fe608060405234801561001057600080fd5b506004361061016c5760003560e01c80636f9e320f116100cd578063c673e58411610081578063e9d68a8e11610066578063e9d68a8e146104ed578063f2fde38b1461050d578063f716f99f1461052057600080fd5b8063c673e58414610489578063ccd37ba3146104a957600080fd5b806379ba5097116100b257806379ba50971461045857806385572ffb146104605780638da5cb5b1461046e57600080fd5b80636f9e320f146103b35780637437ff9f146103c657600080fd5b80633f4b04aa116101245780635e36480c116101095780635e36480c1461036d5780635e7bb0081461038d57806360987c20146103a057600080fd5b80633f4b04aa1461033c5780635215505b1461035757600080fd5b8063181f5a7711610155578063181f5a77146102cd5780632d04ab7614610316578063311cd5131461032957600080fd5b806304666f9c1461017157806306285c6914610186575b600080fd5b61018461017f366004613ec8565b610533565b005b61027760408051608081018252600080825260208201819052918101829052606081019190915260405180608001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160401b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b6040516102c4919081516001600160401b031681526020808301516001600160a01b0390811691830191909152604080840151821690830152606092830151169181019190915260800190565b60405180910390f35b6103096040518060400160405280601181526020017f4f666652616d7020312e362e302d64657600000000000000000000000000000081525081565b6040516102c49190614036565b6101846103243660046140e6565b610547565b610184610337366004614198565b610a4e565b600b546040516001600160401b0390911681526020016102c4565b61035f610ab7565b6040516102c4929190614232565b61038061037b3660046142d3565b610d12565b6040516102c49190614330565b61018461039b366004614899565b610d67565b6101846103ae366004614add565b610ff6565b6101846103c1366004614b71565b6112b3565b610422604080516060810182526000808252602082018190529181019190915250604080516060810182526004546001600160a01b038082168352600160a01b90910463ffffffff166020830152600554169181019190915290565b6040805182516001600160a01b03908116825260208085015163ffffffff169083015292820151909216908201526060016102c4565b6101846112c4565b61018461016c366004614be0565b6000546040516001600160a01b0390911681526020016102c4565b61049c610497366004614c2b565b611375565b6040516102c49190614c8b565b6104df6104b7366004614d00565b6001600160401b03919091166000908152600a60209081526040808320938352929052205490565b6040519081526020016102c4565b6105006104fb366004614d2a565b6114d3565b6040516102c49190614d45565b61018461051b366004614d58565b6115df565b61018461052e366004614ddd565b6115f0565b61053b611632565b6105448161168e565b50565b600061055587890189615132565b602081015151909150156105f257602081015160408083015160608401519151638d8741cb60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001693638d8741cb936105c1933093909190600401615367565b60006040518083038186803b1580156105d957600080fd5b505afa1580156105ed573d6000803e3d6000fd5b505050505b8051515115158061060857508051602001515115155b156106d457600b5460208a0135906001600160401b03808316911610156106ac57600b805467ffffffffffffffff19166001600160401b038316179055600480548351604051633937306f60e01b81526001600160a01b0390921692633937306f926106759291016154b4565b600060405180830381600087803b15801561068f57600080fd5b505af11580156106a3573d6000803e3d6000fd5b505050506106d2565b8160200151516000036106d257604051632261116760e01b815260040160405180910390fd5b505b60005b81602001515181101561098f576000826020015182815181106106fc576106fc6153e2565b60209081029190910101518051604051632cbc26bb60e01b815267ffffffffffffffff60801b608083901b166004820152919250906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632cbc26bb90602401602060405180830381865afa158015610782573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107a691906154c7565b156107d457604051637edeb53960e11b81526001600160401b03821660048201526024015b60405180910390fd5b60006107df826118f5565b9050806001016040516107f2919061551e565b60405180910390208360200151805190602001201461082f5782602001518160010160405163b80d8fa960e01b81526004016107cb929190615611565b60408301518154600160a81b90046001600160401b039081169116141580610870575082606001516001600160401b031683604001516001600160401b0316115b156108b557825160408085015160608601519151636af0786b60e11b81526001600160401b0393841660048201529083166024820152911660448201526064016107cb565b6080830151806108d85760405163504570e360e01b815260040160405180910390fd5b83516001600160401b03166000908152600a60209081526040808320848452909152902054156109305783516040516332cf0cbf60e01b81526001600160401b039091166004820152602481018290526044016107cb565b606084015161094090600161564c565b825467ffffffffffffffff60a81b1916600160a81b6001600160401b0392831602179092559251166000908152600a6020908152604080832094835293905291909120429055506001016106d7565b50602081015181516040517f35c02761bcd3ef995c6a601a1981f4ed3934dcbe5041e24e286c89f5531d17e4926109c7929091615673565b60405180910390a1610a4360008a8a8a8a8a8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808e0282810182019093528d82529093508d92508c9182918501908490808284376000920191909152508b9250611941915050565b505050505050505050565b610a8e610a5d82840184615698565b6040805160008082526020820190925290610a88565b6060815260200190600190039081610a735790505b50611c3a565b604080516000808252602082019092529050610ab1600185858585866000611941565b50505050565b6060806000610ac66006611cfd565b6001600160401b03811115610add57610add613d0a565b604051908082528060200260200182016040528015610b2e57816020015b6040805160808101825260008082526020808301829052928201526060808201528252600019909201910181610afb5790505b5090506000610b3d6006611cfd565b6001600160401b03811115610b5457610b54613d0a565b604051908082528060200260200182016040528015610b7d578160200160208202803683370190505b50905060005b610b8d6006611cfd565b811015610d0957610b9f600682611d07565b828281518110610bb157610bb16153e2565b60200260200101906001600160401b031690816001600160401b03168152505060086000838381518110610be757610be76153e2565b6020908102919091018101516001600160401b039081168352828201939093526040918201600020825160808101845281546001600160a01b038116825260ff600160a01b820416151593820193909352600160a81b90920490931691810191909152600182018054919291606084019190610c62906154e4565b80601f0160208091040260200160405190810160405280929190818152602001828054610c8e906154e4565b8015610cdb5780601f10610cb057610100808354040283529160200191610cdb565b820191906000526020600020905b815481529060010190602001808311610cbe57829003601f168201915b505050505081525050838281518110610cf657610cf66153e2565b6020908102919091010152600101610b83565b50939092509050565b6000610d20600160046156cc565b6002610d2d6080856156f5565b6001600160401b0316610d40919061571b565b610d4a8585611d13565b901c166003811115610d5e57610d5e614306565b90505b92915050565b610d6f611d58565b815181518114610d92576040516320f8fd5960e21b815260040160405180910390fd5b60005b81811015610fe6576000848281518110610db157610db16153e2565b60200260200101519050600081602001515190506000858481518110610dd957610dd96153e2565b6020026020010151905080518214610e04576040516320f8fd5960e21b815260040160405180910390fd5b60005b82811015610fd7576000828281518110610e2357610e236153e2565b6020026020010151600001519050600085602001518381518110610e4957610e496153e2565b6020026020010151905081600014610e9d578060800151821015610e9d578551815151604051633a98d46360e11b81526001600160401b0390921660048301526024820152604481018390526064016107cb565b838381518110610eaf57610eaf6153e2565b602002602001015160200151518160a001515114610efc57805180516060909101516040516370a193fd60e01b815260048101929092526001600160401b031660248201526044016107cb565b60005b8160a0015151811015610fc9576000858581518110610f2057610f206153e2565b6020026020010151602001518281518110610f3d57610f3d6153e2565b602002602001015163ffffffff16905080600014610fc05760008360a001518381518110610f6d57610f6d6153e2565b60200260200101516040015163ffffffff16905080821015610fbe578351516040516348e617b360e01b815260048101919091526024810184905260448101829052606481018390526084016107cb565b505b50600101610eff565b505050806001019050610e07565b50505050806001019050610d95565b50610ff18383611c3a565b505050565b333014611016576040516306e34e6560e31b815260040160405180910390fd5b6040805160008082526020820190925281611053565b604080518082019091526000808252602082015281526020019060019003908161102c5790505b5060a08701515190915015611089576110868660a001518760200151886060015189600001516020015189898989611dc0565b90505b6040805160a081018252875151815287516020908101516001600160401b03168183015288015181830152908701516060820152608081018290526005546001600160a01b0316801561117c576040516308d450a160e01b81526001600160a01b038216906308d450a1906111029085906004016157d3565b600060405180830381600087803b15801561111c57600080fd5b505af192505050801561112d575060015b61117c573d80801561115b576040519150601f19603f3d011682016040523d82523d6000602084013e611160565b606091505b50806040516309c2532560e01b81526004016107cb9190614036565b60408801515115801561119157506080880151155b806111a8575060608801516001600160a01b03163b155b806111cf575060608801516111cd906001600160a01b03166385572ffb60e01b611f71565b155b156111dc575050506112ac565b87516020908101516001600160401b03166000908152600890915260408082205460808b015160608c01519251633cf9798360e01b815284936001600160a01b0390931692633cf979839261123a92899261138892916004016157e6565b6000604051808303816000875af1158015611259573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526112819190810190615822565b5091509150816112a657806040516302a35ba360e21b81526004016107cb9190614036565b50505050505b5050505050565b6112bb611632565b61054481611f8d565b6001546001600160a01b0316331461131e5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016107cb565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6113b86040805160e081019091526000606082018181526080830182905260a0830182905260c08301919091528190815260200160608152602001606081525090565b60ff808316600090815260026020818152604092839020835160e081018552815460608201908152600183015480881660808401526101008104881660a0840152620100009004909616151560c08201529485529182018054845181840281018401909552808552929385830193909283018282801561146157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611443575b50505050508152602001600382018054806020026020016040519081016040528092919081815260200182805480156114c357602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116114a5575b5050505050815250509050919050565b60408051608080820183526000808352602080840182905283850182905260608085018190526001600160401b03878116845260088352928690208651948501875280546001600160a01b0381168652600160a01b810460ff16151593860193909352600160a81b90920490921694830194909452600184018054939492939184019161155f906154e4565b80601f016020809104026020016040519081016040528092919081815260200182805461158b906154e4565b80156114c35780601f106115ad576101008083540402835291602001916114c3565b820191906000526020600020905b8154815290600101906020018083116115bb57505050919092525091949350505050565b6115e7611632565b6105448161206c565b6115f8611632565b60005b815181101561162e57611626828281518110611619576116196153e2565b6020026020010151612115565b6001016115fb565b5050565b6000546001600160a01b0316331461168c5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016107cb565b565b60005b815181101561162e5760008282815181106116ae576116ae6153e2565b60200260200101519050600081602001519050806001600160401b03166000036116eb5760405163c656089560e01b815260040160405180910390fd5b81516001600160a01b0316611713576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b0381166000908152600860205260409020606083015160018201805461173f906154e4565b90506000036117a157815467ffffffffffffffff60a81b1916600160a81b1782556040516001600160401b03841681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a16117de565b8154600160a81b90046001600160401b03166001146117de57604051632105803760e11b81526001600160401b03841660048201526024016107cb565b805115806118135750604080516000602082015201604051602081830303815290604052805190602001208180519060200120145b15611831576040516342bcdf7f60e11b815260040160405180910390fd5b6001820161183f8282615907565b506040840151825485516001600160a01b03166001600160a01b0319921515600160a01b02929092167fffffffffffffffffffffff000000000000000000000000000000000000000000909116171782556118a460066001600160401b03851661243f565b50826001600160401b03167f49f51971edd25182e97182d6ea372a0488ce2ab639f6a3a7ab4df0d2636fe56b836040516118de91906159c6565b60405180910390a250505050806001019050611691565b6001600160401b03811660009081526008602052604081208054600160a01b900460ff16610d615760405163ed053c5960e01b81526001600160401b03841660048201526024016107cb565b60ff878116600090815260026020908152604080832081516080810183528154815260019091015480861693820193909352610100830485169181019190915262010000909104909216151560608301528735906119a08760a4615a14565b90508260600151156119e85784516119b990602061571b565b86516119c690602061571b565b6119d19060a0615a14565b6119db9190615a14565b6119e59082615a14565b90505b368114611a1157604051638e1192e160e01b8152600481018290523660248201526044016107cb565b5081518114611a405781516040516324f7d61360e21b81526004810191909152602481018290526044016107cb565b611a48611d58565b60ff808a1660009081526003602090815260408083203384528252808320815180830190925280548086168352939491939092840191610100909104166002811115611a9657611a96614306565b6002811115611aa757611aa7614306565b9052509050600281602001516002811115611ac457611ac4614306565b148015611b185750600260008b60ff1660ff168152602001908152602001600020600301816000015160ff1681548110611b0057611b006153e2565b6000918252602090912001546001600160a01b031633145b611b3557604051631b41e11d60e31b815260040160405180910390fd5b50816060015115611be5576020820151611b50906001615a27565b60ff16855114611b73576040516371253a2560e01b815260040160405180910390fd5b8351855114611b955760405163a75d88af60e01b815260040160405180910390fd5b60008787604051611ba7929190615a40565b604051908190038120611bbe918b90602001615a50565b604051602081830303815290604052805190602001209050611be38a8288888861244b565b505b6040805182815260208a8101356001600160401b03169082015260ff8b16917f198d6990ef96613a9026203077e422916918b03ff47f0be6bee7b02d8e139ef0910160405180910390a2505050505050505050565b8151600003611c5c5760405163c2e5347d60e01b815260040160405180910390fd5b80516040805160008082526020820190925291159181611c9f565b604080518082019091526000815260606020820152815260200190600190039081611c775790505b50905060005b84518110156112ac57611cf5858281518110611cc357611cc36153e2565b602002602001015184611cef57858381518110611ce257611ce26153e2565b6020026020010151612608565b83612608565b600101611ca5565b6000610d61825490565b6000610d5e8383612f0e565b6001600160401b038216600090815260096020526040812081611d37608085615a64565b6001600160401b031681526020810191909152604001600020549392505050565b467f00000000000000000000000000000000000000000000000000000000000000001461168c57604051630f01ce8560e01b81527f000000000000000000000000000000000000000000000000000000000000000060048201524660248201526044016107cb565b606088516001600160401b03811115611ddb57611ddb613d0a565b604051908082528060200260200182016040528015611e2057816020015b6040805180820190915260008082526020820152815260200190600190039081611df95790505b509050811560005b8a51811015611f635781611ec057848482818110611e4857611e486153e2565b9050602002016020810190611e5d9190615a8a565b63ffffffff1615611ec057848482818110611e7a57611e7a6153e2565b9050602002016020810190611e8f9190615a8a565b8b8281518110611ea157611ea16153e2565b60200260200101516040019063ffffffff16908163ffffffff16815250505b611f3e8b8281518110611ed557611ed56153e2565b60200260200101518b8b8b8b8b87818110611ef257611ef26153e2565b9050602002810190611f049190615aa5565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612f3892505050565b838281518110611f5057611f506153e2565b6020908102919091010152600101611e28565b505098975050505050505050565b6000611f7c83613218565b8015610d5e5750610d5e8383613263565b80516001600160a01b0316611fb5576040516342bcdf7f60e11b815260040160405180910390fd5b805160048054602080850180516001600160a01b039586167fffffffffffffffff0000000000000000000000000000000000000000000000009094168417600160a01b63ffffffff928316021790945560408087018051600580546001600160a01b031916918916919091179055815194855291519094169183019190915251909216908201527fa1c15688cb2c24508e158f6942b9276c6f3028a85e1af8cf3fff0c3ff3d5fc8d9060600160405180910390a150565b336001600160a01b038216036120c45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016107cb565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b806040015160ff16600003612140576000604051631b3fab5160e11b81526004016107cb9190615aeb565b60208082015160ff80821660009081526002909352604083206001810154929390928392169003612191576060840151600182018054911515620100000262ff0000199092169190911790556121cd565b6060840151600182015460ff62010000909104161515901515146121cd576040516321fd80df60e21b815260ff841660048201526024016107cb565b60a0840151805161010010156121f9576001604051631b3fab5160e11b81526004016107cb9190615aeb565b805160000361221e576005604051631b3fab5160e11b81526004016107cb9190615aeb565b612284848460030180548060200260200160405190810160405280929190818152602001828054801561227a57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161225c575b5050505050613305565b8460600151156123b4576122f2848460020180548060200260200160405190810160405280929190818152602001828054801561227a576020028201919060005260206000209081546001600160a01b0316815260019091019060200180831161225c575050505050613305565b60808501518051610100101561231e576002604051631b3fab5160e11b81526004016107cb9190615aeb565b604086015161232e906003615b05565b60ff16815111612354576003604051631b3fab5160e11b81526004016107cb9190615aeb565b81518151101561237a576001604051631b3fab5160e11b81526004016107cb9190615aeb565b805160018401805461ff00191661010060ff8416021790556123a59060028601906020840190613c90565b506123b28582600161336e565b505b6123c08482600261336e565b80516123d59060038501906020840190613c90565b5060408581015160018401805460ff191660ff8316179055865180855560a088015192517fab8b1b57514019638d7b5ce9c638fe71366fe8e2be1c40a7a80f1733d0e9f5479361242e9389939260028a01929190615b21565b60405180910390a16112ac846134c9565b6000610d5e8383613520565b8251600090815b818110156125fe576000600188868460208110612471576124716153e2565b61247e91901a601b615a27565b898581518110612490576124906153e2565b60200260200101518986815181106124aa576124aa6153e2565b6020026020010151604051600081526020016040526040516124e8949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561250a573d6000803e3d6000fd5b505060408051601f1981015160ff808e166000908152600360209081528582206001600160a01b0385168352815285822085870190965285548084168652939750909550929392840191610100900416600281111561256b5761256b614306565b600281111561257c5761257c614306565b905250905060018160200151600281111561259957612599614306565b146125b757604051636518c33d60e11b815260040160405180910390fd5b8051600160ff9091161b8516156125e157604051633d9ef1f160e21b815260040160405180910390fd5b806000015160ff166001901b851794505050806001019050612452565b5050505050505050565b81518151604051632cbc26bb60e01b8152608083901b67ffffffffffffffff60801b166004820152901515907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa158015612685573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126a991906154c7565b1561271a5780156126d857604051637edeb53960e11b81526001600160401b03831660048201526024016107cb565b6040516001600160401b03831681527faab522ed53d887e56ed53dd37398a01aeef6a58e0fa77c2173beb9512d8949339060200160405180910390a150505050565b6000612725836118f5565b6001018054612733906154e4565b80601f016020809104026020016040519081016040528092919081815260200182805461275f906154e4565b80156127ac5780601f10612781576101008083540402835291602001916127ac565b820191906000526020600020905b81548152906001019060200180831161278f57829003601f168201915b505050602088015151929350505060008190036127ea57855160405163676cf24b60e11b81526001600160401b0390911660048201526024016107cb565b856040015151811461280f576040516357e0e08360e01b815260040160405180910390fd5b6000816001600160401b0381111561282957612829613d0a565b604051908082528060200260200182016040528015612852578160200160208202803683370190505b50905060005b828110156129f657600088602001518281518110612878576128786153e2565b602002602001015190507f00000000000000000000000000000000000000000000000000000000000000006001600160401b03168160000151604001516001600160401b0316146128ef5780516040908101519051631c21951160e11b81526001600160401b0390911660048201526024016107cb565b866001600160401b03168160000151602001516001600160401b03161461294357805160200151604051636c95f1eb60e01b81526001600160401b03808a16600483015290911660248201526044016107cb565b6129d0817f2425b0b9f9054c76ff151b0a175b18f37a4a4e82013a72e9f15c9caa095ed21f83600001516020015184600001516040015189805190602001206040516020016129b594939291909384526001600160401b03928316602085015291166040830152606082015260800190565b6040516020818303038152906040528051906020012061356f565b8383815181106129e2576129e26153e2565b602090810291909101015250600101612858565b506000612a0d86838a606001518b60800151613677565b905080600003612a3b57604051633ee8bd3f60e11b81526001600160401b03871660048201526024016107cb565b60005b83811015610a435760005a905060008a602001518381518110612a6357612a636153e2565b602002602001015190506000612a818a836000015160600151610d12565b90506000816003811115612a9757612a97614306565b1480612ab457506003816003811115612ab257612ab2614306565b145b612b0a57815160600151604080516001600160401b03808e16825290921660208301527f3b575419319662b2a6f5e2467d84521517a3382b908eb3d557bb3fdb0c50e23c910160405180910390a1505050612f06565b60608915612be9578b8581518110612b2457612b246153e2565b6020908102919091018101510151600454909150600090600160a01b900463ffffffff16612b5288426156cc565b1190508080612b7257506003836003811115612b7057612b70614306565b145b612b9a576040516354e7e43160e11b81526001600160401b038d1660048201526024016107cb565b8c8681518110612bac57612bac6153e2565b602002602001015160000151600014612be3578c8681518110612bd157612bd16153e2565b60209081029190910101515160808501525b50612c55565b6000826003811115612bfd57612bfd614306565b14612c5557825160600151604080516001600160401b03808f16825290921660208301527f3ef2a99c550a751d4b0b261268f05a803dfb049ab43616a1ffb388f61fe65120910160405180910390a150505050612f06565b8251608001516001600160401b031615612d2e576000826003811115612c7d57612c7d614306565b03612d2e577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e0e03cae8c85600001516080015186602001516040518463ffffffff1660e01b8152600401612cde93929190615bd3565b6020604051808303816000875af1158015612cfd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d2191906154c7565b612d2e5750505050612f06565b60008d604001518681518110612d4657612d466153e2565b6020026020010151905080518460a001515114612d9057835160600151604051631cfe6d8b60e01b81526001600160401b03808f16600483015290911660248201526044016107cb565b612da48c85600001516060015160016136b4565b600080612db2868486613759565b91509150612dc98e876000015160600151846136b4565b8c15612e20576003826003811115612de357612de3614306565b03612e20576000856003811115612dfc57612dfc614306565b14612e2057855151604051632b11b8d960e01b81526107cb91908390600401615bff565b6002826003811115612e3457612e34614306565b14612e79576003826003811115612e4d57612e4d614306565b14612e79578d866000015160600151836040516349362d1f60e11b81526004016107cb93929190615c18565b8560000151600001518660000151606001516001600160401b03168f6001600160401b03167f05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b8d8c81518110612ed157612ed16153e2565b602002602001015186865a612ee6908f6156cc565b604051612ef69493929190615c3d565b60405180910390a4505050505050505b600101612a3e565b6000826000018281548110612f2557612f256153e2565b9060005260206000200154905092915050565b6040805180820190915260008082526020820152602086015160405163bbe4f6db60e01b81526001600160a01b0380831660048301526000917f00000000000000000000000000000000000000000000000000000000000000009091169063bbe4f6db90602401602060405180830381865afa158015612fbc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612fe09190615c74565b90506001600160a01b038116158061300f575061300d6001600160a01b03821663aff2afbf60e01b611f71565b155b156130385760405163ae9b4ce960e01b81526001600160a01b03821660048201526024016107cb565b60008061305088858c6040015163ffffffff1661380d565b9150915060008060006131036040518061010001604052808e81526020018c6001600160401b031681526020018d6001600160a01b031681526020018f608001518152602001896001600160a01b031681526020018f6000015181526020018f6060015181526020018b8152506040516024016130cd9190615c91565b60408051601f198184030181529190526020810180516001600160e01b0316633907753760e01b179052878661138860846138f0565b92509250925082613129578160405163e1cd550960e01b81526004016107cb9190614036565b8151602014613158578151604051631e3be00960e21b81526020600482015260248101919091526044016107cb565b60008280602001905181019061316e9190615d5d565b9050866001600160a01b03168c6001600160a01b0316146131ea57600061319f8d8a61319a868a6156cc565b61380d565b509050868110806131b95750816131b688836156cc565b14155b156131e85760405163a966e21f60e01b81526004810183905260248101889052604481018290526064016107cb565b505b604080518082019091526001600160a01b039098168852602088015250949550505050505095945050505050565b600061322b826301ffc9a760e01b613263565b8015610d61575061325c827fffffffff00000000000000000000000000000000000000000000000000000000613263565b1592915050565b6040517fffffffff0000000000000000000000000000000000000000000000000000000082166024820152600090819060440160408051601f19818403018152919052602080820180516001600160e01b03166301ffc9a760e01b178152825192935060009283928392909183918a617530fa92503d915060005190508280156132ee575060208210155b80156132fa5750600081115b979650505050505050565b60005b8151811015610ff15760ff83166000908152600360205260408120835190919084908490811061333a5761333a6153e2565b6020908102919091018101516001600160a01b03168252810191909152604001600020805461ffff19169055600101613308565b60005b8251811015610ab157600083828151811061338e5761338e6153e2565b60200260200101519050600060028111156133ab576133ab614306565b60ff80871660009081526003602090815260408083206001600160a01b038716845290915290205461010090041660028111156133ea576133ea614306565b1461340b576004604051631b3fab5160e11b81526004016107cb9190615aeb565b6001600160a01b0381166134325760405163d6c62c9b60e01b815260040160405180910390fd5b60405180604001604052808360ff16815260200184600281111561345857613458614306565b905260ff80871660009081526003602090815260408083206001600160a01b0387168452825290912083518154931660ff198416811782559184015190929091839161ffff1916176101008360028111156134b5576134b5614306565b021790555090505050806001019050613371565b60ff81166105445760ff8082166000908152600260205260409020600101546201000090041661350c57604051631e8ed32560e21b815260040160405180910390fd5b600b805467ffffffffffffffff1916905550565b600081815260018301602052604081205461356757508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610d61565b506000610d61565b8151805160608085015190830151608080870151940151604051600095869588956135d395919490939192916020019485526001600160a01b039390931660208501526001600160401b039182166040850152606084015216608082015260a00190565b604051602081830303815290604052805190602001208560200151805190602001208660400151805190602001208760a001516040516020016136169190615e17565b60408051601f198184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e0015b60405160208183030381529060405280519060200120905092915050565b6000806136858585856139ca565b6001600160401b0387166000908152600a6020908152604080832093835292905220549150505b949350505050565b600060026136c36080856156f5565b6001600160401b03166136d6919061571b565b905060006136e48585611d13565b9050816136f3600160046156cc565b901b19168183600381111561370a5761370a614306565b6001600160401b03871660009081526009602052604081209190921b92909217918291613738608088615a64565b6001600160401b031681526020810191909152604001600020555050505050565b604051630304c3e160e51b815260009060609030906360987c209061378690889088908890600401615eae565b600060405180830381600087803b1580156137a057600080fd5b505af19250505080156137b1575060015b6137f0573d8080156137df576040519150601f19603f3d011682016040523d82523d6000602084013e6137e4565b606091505b50600392509050613805565b50506040805160208101909152600081526002905b935093915050565b600080600080600061386e8860405160240161383891906001600160a01b0391909116815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166370a0823160e01b179052888861138860846138f0565b92509250925082613894578160405163e1cd550960e01b81526004016107cb9190614036565b60208251146138c3578151604051631e3be00960e21b81526020600482015260248101919091526044016107cb565b818060200190518101906138d79190615d5d565b6138e182886156cc565b94509450505050935093915050565b6000606060008361ffff166001600160401b0381111561391257613912613d0a565b6040519080825280601f01601f19166020018201604052801561393c576020820181803683370190505b509150863b6139565763030ed58f60e21b60005260046000fd5b5a8581101561397057632be8ca8b60e21b60005260046000fd5b8590036040810481038710613990576337c3be2960e01b60005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d848111156139b35750835b808352806000602085013e50955095509592505050565b82518251600091908183036139f257604051630469ac9960e21b815260040160405180910390fd5b6101018211801590613a0657506101018111155b613a23576040516309bde33960e01b815260040160405180910390fd5b60001982820101610100811115613a4d576040516309bde33960e01b815260040160405180910390fd5b80600003613a7a5786600081518110613a6857613a686153e2565b60200260200101519350505050613c48565b6000816001600160401b03811115613a9457613a94613d0a565b604051908082528060200260200182016040528015613abd578160200160208202803683370190505b50905060008080805b85811015613be75760006001821b8b811603613b215788851015613b0a578c5160018601958e918110613afb57613afb6153e2565b60200260200101519050613b43565b8551600185019487918110613afb57613afb6153e2565b8b5160018401938d918110613b3857613b386153e2565b602002602001015190505b600089861015613b73578d5160018701968f918110613b6457613b646153e2565b60200260200101519050613b95565b8651600186019588918110613b8a57613b8a6153e2565b602002602001015190505b82851115613bb6576040516309bde33960e01b815260040160405180910390fd5b613bc08282613c4f565b878481518110613bd257613bd26153e2565b60209081029190910101525050600101613ac6565b506001850382148015613bf957508683145b8015613c0457508581145b613c21576040516309bde33960e01b815260040160405180910390fd5b836001860381518110613c3657613c366153e2565b60200260200101519750505050505050505b9392505050565b6000818310613c6757613c628284613c6d565b610d5e565b610d5e83835b604080516001602082015290810183905260608101829052600090608001613659565b828054828255906000526020600020908101928215613ce5579160200282015b82811115613ce557825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613cb0565b50613cf1929150613cf5565b5090565b5b80821115613cf15760008155600101613cf6565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b0381118282101715613d4257613d42613d0a565b60405290565b60405160a081016001600160401b0381118282101715613d4257613d42613d0a565b60405160c081016001600160401b0381118282101715613d4257613d42613d0a565b604080519081016001600160401b0381118282101715613d4257613d42613d0a565b604051601f8201601f191681016001600160401b0381118282101715613dd657613dd6613d0a565b604052919050565b60006001600160401b03821115613df757613df7613d0a565b5060051b60200190565b6001600160a01b038116811461054457600080fd5b80356001600160401b0381168114613e2d57600080fd5b919050565b801515811461054457600080fd5b8035613e2d81613e32565b60006001600160401b03821115613e6457613e64613d0a565b50601f01601f191660200190565b600082601f830112613e8357600080fd5b8135613e96613e9182613e4b565b613dae565b818152846020838601011115613eab57600080fd5b816020850160208301376000918101602001919091529392505050565b60006020808385031215613edb57600080fd5b82356001600160401b0380821115613ef257600080fd5b818501915085601f830112613f0657600080fd5b8135613f14613e9182613dde565b81815260059190911b83018401908481019088831115613f3357600080fd5b8585015b83811015613fd957803585811115613f4f5760008081fd5b86016080818c03601f1901811315613f675760008081fd5b613f6f613d20565b89830135613f7c81613e01565b81526040613f8b848201613e16565b8b830152606080850135613f9e81613e32565b83830152928401359289841115613fb757600091508182fd5b613fc58f8d86880101613e72565b908301525085525050918601918601613f37565b5098975050505050505050565b60005b83811015614001578181015183820152602001613fe9565b50506000910152565b60008151808452614022816020860160208601613fe6565b601f01601f19169290920160200192915050565b602081526000610d5e602083018461400a565b8060608101831015610d6157600080fd5b60008083601f84011261406c57600080fd5b5081356001600160401b0381111561408357600080fd5b60208301915083602082850101111561409b57600080fd5b9250929050565b60008083601f8401126140b457600080fd5b5081356001600160401b038111156140cb57600080fd5b6020830191508360208260051b850101111561409b57600080fd5b60008060008060008060008060e0898b03121561410257600080fd5b61410c8a8a614049565b975060608901356001600160401b038082111561412857600080fd5b6141348c838d0161405a565b909950975060808b013591508082111561414d57600080fd5b6141598c838d016140a2565b909750955060a08b013591508082111561417257600080fd5b5061417f8b828c016140a2565b999c989b50969995989497949560c00135949350505050565b6000806000608084860312156141ad57600080fd5b6141b78585614049565b925060608401356001600160401b038111156141d257600080fd5b6141de8682870161405a565b9497909650939450505050565b6001600160a01b0381511682526020810151151560208301526001600160401b03604082015116604083015260006060820151608060608501526136ac608085018261400a565b604080825283519082018190526000906020906060840190828701845b828110156142745781516001600160401b03168452928401929084019060010161424f565b50505083810382850152845180825282820190600581901b8301840187850160005b838110156142c457601f198684030185526142b28383516141eb565b94870194925090860190600101614296565b50909998505050505050505050565b600080604083850312156142e657600080fd5b6142ef83613e16565b91506142fd60208401613e16565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b6004811061432c5761432c614306565b9052565b60208101610d61828461431c565b600060a0828403121561435057600080fd5b614358613d48565b90508135815261436a60208301613e16565b602082015261437b60408301613e16565b604082015261438c60608301613e16565b606082015261439d60808301613e16565b608082015292915050565b8035613e2d81613e01565b803563ffffffff81168114613e2d57600080fd5b600082601f8301126143d857600080fd5b813560206143e8613e9183613dde565b82815260059290921b8401810191818101908684111561440757600080fd5b8286015b848110156144d75780356001600160401b038082111561442b5760008081fd5b9088019060a0828b03601f19018113156144455760008081fd5b61444d613d48565b878401358381111561445f5760008081fd5b61446d8d8a83880101613e72565b82525060408085013561447f81613e01565b828a015260606144908682016143b3565b828401526080915081860135858111156144aa5760008081fd5b6144b88f8c838a0101613e72565b918401919091525091909301359083015250835291830191830161440b565b509695505050505050565b600061014082840312156144f557600080fd5b6144fd613d6a565b9050614509838361433e565b815260a08201356001600160401b038082111561452557600080fd5b61453185838601613e72565b602084015260c084013591508082111561454a57600080fd5b61455685838601613e72565b604084015261456760e085016143a8565b6060840152610100840135608084015261012084013591508082111561458c57600080fd5b50614599848285016143c7565b60a08301525092915050565b600082601f8301126145b657600080fd5b813560206145c6613e9183613dde565b82815260059290921b840181019181810190868411156145e557600080fd5b8286015b848110156144d75780356001600160401b038111156146085760008081fd5b6146168986838b01016144e2565b8452509183019183016145e9565b600082601f83011261463557600080fd5b81356020614645613e9183613dde565b82815260059290921b8401810191818101908684111561466457600080fd5b8286015b848110156144d75780356001600160401b038082111561468757600080fd5b818901915089603f83011261469b57600080fd5b858201356146ab613e9182613dde565b81815260059190911b830160400190878101908c8311156146cb57600080fd5b604085015b83811015614704578035858111156146e757600080fd5b6146f68f6040838a0101613e72565b8452509189019189016146d0565b50875250505092840192508301614668565b600082601f83011261472757600080fd5b81356020614737613e9183613dde565b8083825260208201915060208460051b87010193508684111561475957600080fd5b602086015b848110156144d7578035835291830191830161475e565b600082601f83011261478657600080fd5b81356020614796613e9183613dde565b82815260059290921b840181019181810190868411156147b557600080fd5b8286015b848110156144d75780356001600160401b03808211156147d95760008081fd5b9088019060a0828b03601f19018113156147f35760008081fd5b6147fb613d48565b614806888501613e16565b81526040808501358481111561481c5760008081fd5b61482a8e8b838901016145a5565b8a84015250606080860135858111156148435760008081fd5b6148518f8c838a0101614624565b838501525060809150818601358581111561486c5760008081fd5b61487a8f8c838a0101614716565b91840191909152509190930135908301525083529183019183016147b9565b600080604083850312156148ac57600080fd5b6001600160401b03833511156148c157600080fd5b6148ce8484358501614775565b91506001600160401b03602084013511156148e857600080fd5b6020830135830184601f8201126148fe57600080fd5b61490b613e918235613dde565b81358082526020808301929160051b84010187101561492957600080fd5b602083015b6020843560051b850101811015614acf576001600160401b038135111561495457600080fd5b87603f82358601011261496657600080fd5b614979613e916020833587010135613dde565b81358501602081810135808452908301929160059190911b016040018a10156149a157600080fd5b604083358701015b83358701602081013560051b01604001811015614abf576001600160401b03813511156149d557600080fd5b833587018135016040818d03603f190112156149f057600080fd5b6149f8613d8c565b604082013581526001600160401b0360608301351115614a1757600080fd5b8c605f606084013584010112614a2c57600080fd5b6040606083013583010135614a43613e9182613dde565b808282526020820191508f60608460051b6060880135880101011115614a6857600080fd5b6060808601358601015b60608460051b606088013588010101811015614a9f57614a91816143b3565b835260209283019201614a72565b5080602085015250505080855250506020830192506020810190506149a9565b508452506020928301920161492e565b508093505050509250929050565b600080600080600060608688031215614af557600080fd5b85356001600160401b0380821115614b0c57600080fd5b614b1889838a016144e2565b96506020880135915080821115614b2e57600080fd5b614b3a89838a016140a2565b90965094506040880135915080821115614b5357600080fd5b50614b60888289016140a2565b969995985093965092949392505050565b600060608284031215614b8357600080fd5b604051606081018181106001600160401b0382111715614ba557614ba5613d0a565b6040528235614bb381613e01565b8152614bc1602084016143b3565b60208201526040830135614bd481613e01565b60408201529392505050565b600060208284031215614bf257600080fd5b81356001600160401b03811115614c0857600080fd5b820160a08185031215613c4857600080fd5b803560ff81168114613e2d57600080fd5b600060208284031215614c3d57600080fd5b610d5e82614c1a565b60008151808452602080850194506020840160005b83811015614c805781516001600160a01b031687529582019590820190600101614c5b565b509495945050505050565b60208152600082518051602084015260ff602082015116604084015260ff604082015116606084015260608101511515608084015250602083015160c060a0840152614cda60e0840182614c46565b90506040840151601f198483030160c0850152614cf78282614c46565b95945050505050565b60008060408385031215614d1357600080fd5b614d1c83613e16565b946020939093013593505050565b600060208284031215614d3c57600080fd5b610d5e82613e16565b602081526000610d5e60208301846141eb565b600060208284031215614d6a57600080fd5b8135613c4881613e01565b600082601f830112614d8657600080fd5b81356020614d96613e9183613dde565b8083825260208201915060208460051b870101935086841115614db857600080fd5b602086015b848110156144d7578035614dd081613e01565b8352918301918301614dbd565b60006020808385031215614df057600080fd5b82356001600160401b0380821115614e0757600080fd5b818501915085601f830112614e1b57600080fd5b8135614e29613e9182613dde565b81815260059190911b83018401908481019088831115614e4857600080fd5b8585015b83811015613fd957803585811115614e6357600080fd5b860160c0818c03601f19011215614e7a5760008081fd5b614e82613d6a565b8882013581526040614e95818401614c1a565b8a8301526060614ea6818501614c1a565b8284015260809150614eb9828501613e40565b9083015260a08381013589811115614ed15760008081fd5b614edf8f8d83880101614d75565b838501525060c0840135915088821115614ef95760008081fd5b614f078e8c84870101614d75565b9083015250845250918601918601614e4c565b80356001600160e01b0381168114613e2d57600080fd5b600082601f830112614f4257600080fd5b81356020614f52613e9183613dde565b82815260069290921b84018101918181019086841115614f7157600080fd5b8286015b848110156144d75760408189031215614f8e5760008081fd5b614f96613d8c565b614f9f82613e16565b8152614fac858301614f1a565b81860152835291830191604001614f75565b600082601f830112614fcf57600080fd5b81356020614fdf613e9183613dde565b82815260059290921b84018101918181019086841115614ffe57600080fd5b8286015b848110156144d75780356001600160401b03808211156150225760008081fd5b9088019060a0828b03601f190181131561503c5760008081fd5b615044613d48565b61504f888501613e16565b8152604080850135848111156150655760008081fd5b6150738e8b83890101613e72565b8a8401525060609350615087848601613e16565b908201526080615098858201613e16565b93820193909352920135908201528352918301918301615002565b600082601f8301126150c457600080fd5b813560206150d4613e9183613dde565b82815260069290921b840181019181810190868411156150f357600080fd5b8286015b848110156144d757604081890312156151105760008081fd5b615118613d8c565b8135815284820135858201528352918301916040016150f7565b6000602080838503121561514557600080fd5b82356001600160401b038082111561515c57600080fd5b908401906080828703121561517057600080fd5b615178613d20565b82358281111561518757600080fd5b8301604081890381131561519a57600080fd5b6151a2613d8c565b8235858111156151b157600080fd5b8301601f81018b136151c257600080fd5b80356151d0613e9182613dde565b81815260069190911b8201890190898101908d8311156151ef57600080fd5b928a01925b8284101561523f5785848f03121561520c5760008081fd5b615214613d8c565b843561521f81613e01565b815261522c858d01614f1a565b818d0152825292850192908a01906151f4565b84525050508287013591508482111561525757600080fd5b6152638a838501614f31565b8188015283525050828401358281111561527c57600080fd5b61528888828601614fbe565b858301525060408301359350818411156152a157600080fd5b6152ad878585016150b3565b6040820152606083013560608201528094505050505092915050565b600082825180855260208086019550808260051b84010181860160005b8481101561535a57601f19868403018952815160a06001600160401b0380835116865286830151828888015261531e8388018261400a565b604085810151841690890152606080860151909316928801929092525060809283015192909501919091525097830197908301906001016152e6565b5090979650505050505050565b6001600160a01b03851681526000602060808184015261538a60808401876152c9565b83810360408581019190915286518083528388019284019060005b818110156153ca578451805184528601518684015293850193918301916001016153a5565b50508094505050505082606083015295945050505050565b634e487b7160e01b600052603260045260246000fd5b805160408084528151848201819052600092602091908201906060870190855b8181101561544f57835180516001600160a01b031684528501516001600160e01b0316858401529284019291850191600101615418565b50508583015187820388850152805180835290840192506000918401905b808310156154a857835180516001600160401b031683528501516001600160e01b03168583015292840192600192909201919085019061546d565b50979650505050505050565b602081526000610d5e60208301846153f8565b6000602082840312156154d957600080fd5b8151613c4881613e32565b600181811c908216806154f857607f821691505b60208210810361551857634e487b7160e01b600052602260045260246000fd5b50919050565b600080835461552c816154e4565b60018281168015615544576001811461555957615588565b60ff1984168752821515830287019450615588565b8760005260208060002060005b8581101561557f5781548a820152908401908201615566565b50505082870194505b50929695505050505050565b600081546155a1816154e4565b8085526020600183811680156155be57600181146155d857615606565b60ff1985168884015283151560051b880183019550615606565b866000528260002060005b858110156155fe5781548a82018601529083019084016155e3565b890184019650505b505050505092915050565b604081526000615624604083018561400a565b8281036020840152614cf78185615594565b634e487b7160e01b600052601160045260246000fd5b6001600160401b0381811683821601908082111561566c5761566c615636565b5092915050565b60408152600061568660408301856152c9565b8281036020840152614cf781856153f8565b6000602082840312156156aa57600080fd5b81356001600160401b038111156156c057600080fd5b6136ac84828501614775565b81810381811115610d6157610d61615636565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b038084168061570f5761570f6156df565b92169190910692915050565b8082028115828204841417610d6157610d61615636565b80518252600060206001600160401b0381840151168185015260408084015160a0604087015261576560a087018261400a565b90506060850151868203606088015261577e828261400a565b608087810151898303918a01919091528051808352908601935060009250908501905b808310156154a857835180516001600160a01b03168352860151868301529285019260019290920191908401906157a1565b602081526000610d5e6020830184615732565b6080815260006157f96080830187615732565b61ffff9590951660208301525060408101929092526001600160a01b0316606090910152919050565b60008060006060848603121561583757600080fd5b835161584281613e32565b60208501519093506001600160401b0381111561585e57600080fd5b8401601f8101861361586f57600080fd5b805161587d613e9182613e4b565b81815287602083850101111561589257600080fd5b6158a3826020830160208601613fe6565b809450505050604084015190509250925092565b601f821115610ff1576000816000526020600020601f850160051c810160208610156158e05750805b601f850160051c820191505b818110156158ff578281556001016158ec565b505050505050565b81516001600160401b0381111561592057615920613d0a565b6159348161592e84546154e4565b846158b7565b602080601f83116001811461596957600084156159515750858301515b600019600386901b1c1916600185901b1785556158ff565b600085815260208120601f198616915b8281101561599857888601518255948401946001909101908401615979565b50858210156159b65787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60208152600082546001600160a01b038116602084015260ff8160a01c16151560408401526001600160401b038160a81c16606084015250608080830152610d5e60a0830160018501615594565b80820180821115610d6157610d61615636565b60ff8181168382160190811115610d6157610d61615636565b8183823760009101908152919050565b828152606082602083013760800192915050565b60006001600160401b0380841680615a7e57615a7e6156df565b92169190910492915050565b600060208284031215615a9c57600080fd5b610d5e826143b3565b6000808335601e19843603018112615abc57600080fd5b8301803591506001600160401b03821115615ad657600080fd5b60200191503681900382131561409b57600080fd5b6020810160068310615aff57615aff614306565b91905290565b60ff818116838216029081169081811461566c5761566c615636565b600060a0820160ff881683526020878185015260a0604085015281875480845260c0860191508860005282600020935060005b81811015615b795784546001600160a01b031683526001948501949284019201615b54565b50508481036060860152865180825290820192508187019060005b81811015615bb95782516001600160a01b031685529383019391830191600101615b94565b50505060ff851660808501525090505b9695505050505050565b60006001600160401b03808616835280851660208401525060606040830152614cf7606083018461400a565b8281526040602082015260006136ac604083018461400a565b6001600160401b03848116825283166020820152606081016136ac604083018461431c565b848152615c4d602082018561431c565b608060408201526000615c63608083018561400a565b905082606083015295945050505050565b600060208284031215615c8657600080fd5b8151613c4881613e01565b6020815260008251610100806020850152615cb061012085018361400a565b91506020850151615ccc60408601826001600160401b03169052565b5060408501516001600160a01b038116606086015250606085015160808501526080850151615d0660a08601826001600160a01b03169052565b5060a0850151601f19808685030160c0870152615d23848361400a565b935060c08701519150808685030160e0870152615d40848361400a565b935060e0870151915080868503018387015250615bc9838261400a565b600060208284031215615d6f57600080fd5b5051919050565b600082825180855260208086019550808260051b84010181860160005b8481101561535a57601f19868403018952815160a08151818652615db98287018261400a565b9150506001600160a01b03868301511686860152604063ffffffff8184015116818701525060608083015186830382880152615df5838261400a565b6080948501519790940196909652505098840198925090830190600101615d93565b602081526000610d5e6020830184615d76565b60008282518085526020808601955060208260051b8401016020860160005b8481101561535a57601f19868403018952615e6583835161400a565b98840198925090830190600101615e49565b60008151808452602080850194506020840160005b83811015614c8057815163ffffffff1687529582019590820190600101615e8c565b60608152600084518051606084015260208101516001600160401b0380821660808601528060408401511660a08601528060608401511660c08601528060808401511660e0860152505050602085015161014080610100850152615f166101a085018361400a565b91506040870151605f198086850301610120870152615f35848361400a565b935060608901519150615f52838701836001600160a01b03169052565b608089015161016087015260a0890151925080868503016101808701525050615f7b8282615d76565b9150508281036020840152615f908186615e2a565b90508281036040840152615bc98185615e7756fea164736f6c6343000818000a",
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

	Commit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	Execute(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte) (*types.Transaction, error)

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
