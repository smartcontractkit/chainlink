// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm_2_evm_offramp

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

type EVM2EVMOffRampDynamicConfig struct {
	PermissionLessExecutionThresholdSeconds uint32
	MaxDataBytes                            uint32
	MaxNumberOfTokensPerMsg                 uint16
	Router                                  common.Address
	PriceRegistry                           common.Address
}

type EVM2EVMOffRampGasLimitOverride struct {
	ReceiverExecutionGasLimit *big.Int
	TokenGasOverrides         []uint32
}

type EVM2EVMOffRampRateLimitToken struct {
	SourceToken common.Address
	DestToken   common.Address
}

type EVM2EVMOffRampStaticConfig struct {
	CommitStore         common.Address
	ChainSelector       uint64
	SourceChainSelector uint64
	OnRamp              common.Address
	PrevOffRamp         common.Address
	RmnProxy            common.Address
	TokenAdminRegistry  common.Address
}

type InternalEVM2EVMMessage struct {
	SourceChainSelector uint64
	Sender              common.Address
	Receiver            common.Address
	SequenceNumber      uint64
	GasLimit            *big.Int
	Strict              bool
	Nonce               uint64
	FeeToken            common.Address
	FeeTokenAmount      *big.Int
	Data                []byte
	TokenAmounts        []ClientEVMTokenAmount
	SourceTokenData     [][]byte
	MessageId           [32]byte
}

type InternalExecutionReport struct {
	Messages          []InternalEVM2EVMMessage
	OffchainTokenData [][][]byte
	Proofs            [][32]byte
	ProofFlagBits     *big.Int
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  *big.Int
	Rate      *big.Int
}

type RateLimiterTokenBucket struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

var EVM2EVMOffRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CommitStoreAlreadyInUse\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"DestinationGasAmountCountMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumOCR2BaseNoChecks.InvalidConfigErrorType\",\"name\":\"errorType\",\"type\":\"uint8\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMessageId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidSourceChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenGasOverride\",\"type\":\"uint256\"}],\"name\":\"InvalidTokenGasOverride\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PriceNotFoundForToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReleased\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePre\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balancePost\",\"type\":\"uint256\"}],\"name\":\"ReleaseOrMintBalanceMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"SkippedAlreadyExecutedMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SkippedIncorrectNonce\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SkippedSenderWithPreviousRampMessageInflight\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"name\":\"TokenAggregateRateLimitAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"name\":\"TokenAggregateRateLimitRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"strict\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"sourceTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.EVM2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllRateLimitTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"sourceTokens\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"destTokens\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getSenderNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAdminRegistry\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"strict\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"sourceTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.EVM2EVMMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReport\",\"name\":\"report\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"receiverExecutionGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"tokenGasOverrides\",\"type\":\"uint32[]\"}],\"internalType\":\"structEVM2EVMOffRamp.GasLimitOverride[]\",\"name\":\"gasLimitOverrides\",\"type\":\"tuple[]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR2Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOffRamp.RateLimitToken[]\",\"name\":\"removes\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMOffRamp.RateLimitToken[]\",\"name\":\"adds\",\"type\":\"tuple[]\"}],\"name\":\"updateRateLimitTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101a06040523480156200001257600080fd5b5060405162006884380380620068848339810160408190526200003591620004ec565b8033806000816200008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c057620000c081620002ca565b50506040805160a081018252602084810180516001600160801b039081168085524263ffffffff169385018490528751151585870181905292518216606080870182905298909601519091166080948501819052600380546001600160a01b031916909217600160801b9485021760ff60a01b1916600160a01b90930292909217905502909117600455469052508201516001600160a01b031615806200016f575081516001600160a01b0316155b8062000186575060c08201516001600160a01b0316155b15620001a5576040516342bcdf7f60e11b815260040160405180910390fd5b81600001516001600160a01b0316634120fccd6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001e8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200020e9190620005b5565b6001600160401b03166001146200023857604051636fc2a20760e11b815260040160405180910390fd5b81516001600160a01b0390811660a090815260408401516001600160401b0390811660c0908152602086015190911660e05260608501518316610100526080850151831661014052908401518216610160528301511661018052620002bd7f8acd72527118c8324937b1a42e02cd246697c3b633f1742f3cae11de233722b362000375565b6101205250620005da9050565b336001600160a01b03821603620003245760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008160c05160e05161010051604051602001620003bf94939291909384526001600160401b039283166020850152911660408301526001600160a01b0316606082015260800190565b604051602081830303815290604052805190602001209050919050565b60405160e081016001600160401b03811182821017156200040d57634e487b7160e01b600052604160045260246000fd5b60405290565b80516001600160a01b03811681146200042b57600080fd5b919050565b80516001600160401b03811681146200042b57600080fd5b80516001600160801b03811681146200042b57600080fd5b6000606082840312156200047357600080fd5b604051606081016001600160401b0381118282101715620004a457634e487b7160e01b600052604160045260246000fd5b806040525080915082518015158114620004bd57600080fd5b8152620004cd6020840162000448565b6020820152620004e06040840162000448565b60408201525092915050565b6000808284036101408112156200050257600080fd5b60e08112156200051157600080fd5b506200051c620003dc565b620005278462000413565b8152620005376020850162000430565b60208201526200054a6040850162000430565b60408201526200055d6060850162000413565b6060820152620005706080850162000413565b60808201526200058360a0850162000413565b60a08201526200059660c0850162000413565b60c08201529150620005ac8460e0850162000460565b90509250929050565b600060208284031215620005c857600080fd5b620005d38262000430565b9392505050565b60805160a05160c05160e05161010051610120516101405161016051610180516161bb620006c9600039600081816102ec01528181611c6a01526135bc0152600081816102bd01528181611c420152611f2701526000818161028e01528181610d8d01528181610df201528181611c18015281816124a4015261250e015260006120c601526000818161025f0152611bee0152600081816101ff0152611b9201526000818161022f01528181611bc601528181611ee40152818161303b01526136e90152600081816101d001528181611b6d01526121ac015260008181611e3e0152611e8a01526161bb6000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c806385572ffb116100d8578063afcb95d71161008c578063c92b283211610066578063c92b2832146105f3578063f077b59214610606578063f2fde38b1461061c57600080fd5b8063afcb95d7146105ad578063b1dc65a4146105cd578063b6113fce146105e057600080fd5b8063873504d7116100bd578063873504d7146105765780638926c4ee146105895780638da5cb5b1461059c57600080fd5b806385572ffb1461053c578063856c82471461054a57600080fd5b8063599f64311161013a5780637437ff9f116101145780637437ff9f1461046157806379ba50971461050457806381ff70481461050c57600080fd5b8063599f643114610414578063666cab8d14610439578063704b6c021461044e57600080fd5b8063181f5a771161016b578063181f5a77146103525780631ef381741461039b578063546719cd146103b057600080fd5b806306285c6914610187578063142a98fc14610332575b600080fd5b61031c6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c08101919091526040518060e001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815250905090565b604051610329919061453c565b60405180910390f35b6103456103403660046145d2565b61062f565b6040516103299190614632565b61038e6040518060400160405280601481526020017f45564d3245564d4f666652616d7020312e352e3000000000000000000000000081525081565b6040516103299190614690565b6103ae6103a93660046148b9565b6106aa565b005b6103b8610a9e565b604051610329919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b6002546001600160a01b03165b6040516001600160a01b039091168152602001610329565b610441610b53565b60405161032991906149cb565b6103ae61045c3660046149de565b610bb5565b6104f76040805160a081018252600080825260208201819052918101829052606081018290526080810191909152506040805160a081018252600a5463ffffffff8082168352640100000000820416602083015268010000000000000000810461ffff16928201929092526a01000000000000000000009091046001600160a01b039081166060830152600b5416608082015290565b60405161032991906149fb565b6103ae610c7e565b6007546005546040805163ffffffff80851682526401000000009094049093166020840152820152606001610329565b6103ae610182366004614a51565b61055d6105583660046149de565b610d61565b60405167ffffffffffffffff9091168152602001610329565b6103ae610584366004614b1d565b610e64565b6103ae610597366004615048565b611037565b6000546001600160a01b0316610421565b604080516001815260006020820181905291810191909152606001610329565b6103ae6105db36600461519f565b61129c565b6103ae6105ee366004615284565b6114a7565b6103ae61060136600461533b565b6117e7565b61060e611852565b6040516103299291906153a9565b6103ae61062a3660046149de565b611978565b600061063d600160046153fd565b600261064a60808561543f565b67ffffffffffffffff1661065e9190615466565b6010600061066d60808761547d565b67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002054901c1660038111156106a4576106a46145ef565b92915050565b84518460ff16601f8211156106f75760016040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ee91906154a4565b60405180910390fd5b806000036107345760006040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ee91906154a4565b61073c611989565b610745856119ff565b60095460005b818110156107bc57600860006009838154811061076a5761076a6154be565b60009182526020808320909101546001600160a01b03168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905560010161074b565b5050865160005b8181101561095f5760008982815181106107df576107df6154be565b60200260200101519050600060028111156107fc576107fc6145ef565b6001600160a01b038216600090815260086020526040902054610100900460ff16600281111561082e5761082e6145ef565b146108685760026040517f367f56a20000000000000000000000000000000000000000000000000000000081526004016106ee91906154a4565b6001600160a01b0381166108a8576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff8316815260208101600290526001600160a01b03821660009081526008602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000161761010083600281111561094b5761094b6145ef565b0217905550905050508060010190506107c3565b5087516109739060099060208b01906144aa565b506006805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908a1617179055600780546109f99146913091906000906109cb9063ffffffff166154ed565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168c8c8c8c8c8c611cc9565b6005819055600780544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff841681179094556040519083048216947f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0594610a8a9487949293918316921691909117908f908f908f908f908f908f90615510565b60405180910390a150505050505050505050565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526040805160a0810182526003546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff1660208501527401000000000000000000000000000000000000000090920460ff161515938301939093526004548084166060840152049091166080820152610b4e90611d56565b905090565b60606009805480602002602001604051908101604052809291908181526020018280548015610bab57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610b8d575b5050505050905090565b6000546001600160a01b03163314801590610bdb57506002546001600160a01b03163314155b15610c12576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383169081179091556040519081527f8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c9060200160405180910390a150565b6001546001600160a01b03163314610cf2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016106ee565b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6001600160a01b0381166000908152600f602052604081205467ffffffffffffffff168082036106a4577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316156106a4576040517f856c82470000000000000000000000000000000000000000000000000000000081526001600160a01b0384811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063856c824790602401602060405180830381865afa158015610e39573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e5d91906155a6565b9392505050565b610e6c611989565b60005b8251811015610f3f57610ea9838281518110610e8d57610e8d6154be565b602002602001015160200151600c611e0890919063ffffffff16565b15610f37577fcbf3cbeaed4ac1d605ed30f4af06c35acaeff2379db7f6146c9cceee83d58782838281518110610ee157610ee16154be565b602002602001015160000151848381518110610eff57610eff6154be565b602002602001015160200151604051610f2e9291906001600160a01b0392831681529116602082015260400190565b60405180910390a15b600101610e6f565b5060005b815181101561103257610f9c828281518110610f6157610f616154be565b602002602001015160200151838381518110610f7f57610f7f6154be565b602002602001015160000151600c611e1d9092919063ffffffff16565b1561102a577ffc23abf7ddbd3c02b1420dafa2355c56c1a06fbb8723862ac14d6bd74177361a828281518110610fd457610fd46154be565b602002602001015160000151838381518110610ff257610ff26154be565b6020026020010151602001516040516110219291906001600160a01b0392831681529116602082015260400190565b60405180910390a15b600101610f43565b505050565b61103f611e3b565b8151518151811461107c576040517f83e3f56400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b818110156112915760008460000151828151811061109f5761109f6154be565b6020026020010151905060008483815181106110bd576110bd6154be565b60209081029190910101518051909150801561112c57826080015181101561112c5761018083015160808401516040517f9c6db58d00000000000000000000000000000000000000000000000000000000815260048101929092526024820152604481018290526064016106ee565b816020015151836101400151511461118e5761018083015160608401516040517f85d2e5bf000000000000000000000000000000000000000000000000000000008152600481019290925267ffffffffffffffff1660248201526044016106ee565b61016083015160005b846101400151518110156112815760008287815181106111b9576111b96154be565b60200260200101518060200190518101906111d49190615608565b90506000856020015183815181106111ee576111ee6154be565b602002602001015163ffffffff169050806000141580156112185750816060015163ffffffff1681105b156112775761018087015160608301516040517fef0c635200000000000000000000000000000000000000000000000000000000815260048101929092526024820185905263ffffffff166044820152606481018290526084016106ee565b5050600101611197565b505050505080600101905061107f565b506110328383611ebc565b6112a6878761290f565b6005548835908082146112ef576040517f93df584c00000000000000000000000000000000000000000000000000000000815260048101829052602481018390526044016106ee565b6112f7611e3b565b6040805183815260208c81013560081c63ffffffff16908201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a13360009081526008602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561137f5761137f6145ef565b6002811115611390576113906145ef565b90525090506002816020015160028111156113ad576113ad6145ef565b1480156113e757506009816000015160ff16815481106113cf576113cf6154be565b6000918252602090912001546001600160a01b031633145b61141d576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50600061142b856020615466565b611436886020615466565b6114428b6101446156d4565b61144c91906156d4565b61145691906156d4565b905036811461149a576040517f8e1192e1000000000000000000000000000000000000000000000000000000008152600481018290523660248201526044016106ee565b5050505050505050505050565b3330146114e0576040517f371a732800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516000808252602082019092528161151d565b60408051808201909152600080825260208201528152602001906001900390816114f65790505b50905060006115306101408701876156e7565b905011156115a7576115a46115496101408701876156e7565b6115596040890160208a016149de565b604080516001600160a01b0390921660208301520160408051601f1981840301815291815261158e9060608b01908b016149de565b61159c6101608b018b61574f565b8a8a8a612966565b90505b6115b56101208601866157b7565b15905080156115c657506080850135155b806115e857506115dc60608601604087016149de565b6001600160a01b03163b155b8061163357506116317f85572ffb0000000000000000000000000000000000000000000000000000000061162260608801604089016149de565b6001600160a01b031690612b8f565b155b1561163e57506117e1565b600a546040805160a08101909152610180870135815260009182916a01000000000000000000009091046001600160a01b031690633cf979839060208082019061168a908c018c6145d2565b67ffffffffffffffff1681526020018a60200160208101906116ac91906149de565b604080516001600160a01b0390921660208301520160408051601f1981840301815291905281526020016116e46101208c018c6157b7565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060200186905261138860808b013561173960608d0160408e016149de565b6040518563ffffffff1660e01b81526004016117589493929190615861565b6000604051808303816000875af1158015611777573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261179f9190810190615926565b5091509150816117dd57806040517f0a8d6e8c0000000000000000000000000000000000000000000000000000000081526004016106ee9190614690565b5050505b50505050565b6000546001600160a01b0316331480159061180d57506002546001600160a01b03163314155b15611844576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61184f600382612bab565b50565b6060806000611861600c612d90565b90508067ffffffffffffffff81111561187c5761187c6146a3565b6040519080825280602002602001820160405280156118a5578160200160208202803683370190505b5092508067ffffffffffffffff8111156118c1576118c16146a3565b6040519080825280602002602001820160405280156118ea578160200160208202803683370190505b50915060005b8181101561197257600080611906600c84612d9b565b915091508086848151811061191d5761191d6154be565b60200260200101906001600160a01b031690816001600160a01b03168152505081858481518110611950576119506154be565b6001600160a01b039092166020928302919091019091015250506001016118f0565b50509091565b611980611989565b61184f81612db9565b6000546001600160a01b031633146119fd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016106ee565b565b600081806020019051810190611a159190615980565b60608101519091506001600160a01b0316611a5c576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600a805460208085015160408087015160608089015163ffffffff9889167fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000909716969096176401000000009890941697909702929092177fffff00000000000000000000000000000000000000000000ffffffffffffffff166801000000000000000061ffff909316929092027fffff0000000000000000000000000000000000000000ffffffffffffffffffff16919091176a01000000000000000000006001600160a01b039485160217909355608080860151600b80547fffffffffffffffffffffffff000000000000000000000000000000000000000016918516919091179055835160e0810185527f0000000000000000000000000000000000000000000000000000000000000000841681527f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff908116938201939093527f0000000000000000000000000000000000000000000000000000000000000000909216828501527f00000000000000000000000000000000000000000000000000000000000000008316948201949094527f00000000000000000000000000000000000000000000000000000000000000008216938101939093527f0000000000000000000000000000000000000000000000000000000000000000811660a08401527f00000000000000000000000000000000000000000000000000000000000000001660c0830152517f7879e20bb60a503429de4a2c912b5904f08a39f2af054c10fb46434b5d61126091611cbd918490615a1f565b60405180910390a15050565b6000808a8a8a8a8a8a8a8a8a604051602001611ced99989796959493929190615ae1565b60408051601f1981840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152611de482606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642611dc891906153fd565b85608001516fffffffffffffffffffffffffffffffff16612e94565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b6000610e5d836001600160a01b038416612ebc565b6000611e33846001600160a01b03851684612ec8565b949350505050565b467f0000000000000000000000000000000000000000000000000000000000000000146119fd576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000060048201524660248201526044016106ee565b6040517f2cbc26bb0000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000060801b77ffffffffffffffff000000000000000000000000000000001660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690632cbc26bb90602401602060405180830381865afa158015611f76573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f9a9190615b69565b15611fd1576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815151600081900361200e576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015151811461204c576040517f57e0e08300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008167ffffffffffffffff811115612067576120676146a3565b604051908082528060200260200182016040528015612090578160200160208202803683370190505b50905060005b82811015612168576000856000015182815181106120b6576120b66154be565b602002602001015190506120ea817f0000000000000000000000000000000000000000000000000000000000000000612ede565b8383815181106120fc576120fc6154be565b602002602001018181525050806101800151838381518110612120576121206154be565b60200260200101511461215f576040517f7185cf6b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50600101612096565b508251604080860151606087015191517f32048875000000000000000000000000000000000000000000000000000000008152921515926000926001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016926332048875926121e292889291600401615bb7565b602060405180830381865afa1580156121ff573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122239190615bed565b90508060000361225f576040517fea75680100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b848110156117dd57600087600001518281518110612282576122826154be565b60200260200101519050600061229b826060015161062f565b905060008160038111156122b1576122b16145ef565b14806122ce575060038160038111156122cc576122cc6145ef565b145b61231457816060015167ffffffffffffffff167fe3dd0bec917c965a133ddb2c84874725ee1e2fd8d763c19efa36d6a11cd82b1f60405160405180910390a25050612907565b606085156123fa5788848151811061232e5761232e6154be565b6020908102919091018101510151600a5490915060009063ffffffff1661235587426153fd565b119050808061237557506003836003811115612373576123736145ef565b145b6123ab576040517f6358b0d000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8985815181106123bd576123bd6154be565b6020026020010151600001516000146123f4578985815181106123e2576123e26154be565b60209081029190910101515160808501525b5061245e565b600082600381111561240e5761240e6145ef565b1461245e57606083015160405167ffffffffffffffff90911681527f67d9ba0f63d427c482c2736300e6d5a34c6691dbcdea8ad35828a1f1ba47e8729060200160405180910390a1505050612907565b60c083015167ffffffffffffffff16156126df576020808401516001600160a01b03166000908152600f909152604081205467ffffffffffffffff1690819003612649577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316156126495760208401516040517f856c82470000000000000000000000000000000000000000000000000000000081526001600160a01b0391821660048201527f00000000000000000000000000000000000000000000000000000000000000009091169063856c824790602401602060405180830381865afa158015612557573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061257b91906155a6565b60c085015190915067ffffffffffffffff16612598826001615c06565b67ffffffffffffffff16146125f95783602001516001600160a01b03168460c0015167ffffffffffffffff167fe44a20935573a783dd0d5991c92d7b6a0eb3173566530364db3ec10e9a990b5d60405160405180910390a350505050612907565b6020848101516001600160a01b03166000908152600f9091526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff83161790555b600083600381111561265d5761265d6145ef565b036126dd5760c084015167ffffffffffffffff1661267c826001615c06565b67ffffffffffffffff16146126dd5783602001516001600160a01b03168460c0015167ffffffffffffffff167fd32ddb11d71e3d63411d37b09f9a8b28664f1cb1338bfd1413c173b0ebf4123760405160405180910390a350505050612907565b505b60008a6020015185815181106126f7576126f76154be565b602002602001015190506127238460600151856000015186610140015151876101200151518551613039565b612732846060015160016131ba565b600080612740868486613264565b915091506127528660600151836131ba565b88156127be57600382600381111561276c5761276c6145ef565b036127be576000856003811115612785576127856145ef565b146127be57806040517fcf19edfd0000000000000000000000000000000000000000000000000000000081526004016106ee9190614690565b60028260038111156127d2576127d26145ef565b1461282a5760038260038111156127eb576127eb6145ef565b1461282a578560600151826040517f9e2616030000000000000000000000000000000000000000000000000000000081526004016106ee929190615c27565b60c086015167ffffffffffffffff16156128b2576000856003811115612852576128526145ef565b036128b2576020808701516001600160a01b03166000908152600f90915260408120805467ffffffffffffffff169161288a83615c45565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505b856101800151866060015167ffffffffffffffff167fd4f851956a5d67c3997d1c9205045fef79bae2947fdee7e9e2641abc7391ef6584846040516128f8929190615c62565b60405180910390a35050505050505b600101612262565b61296261291e82840184615c82565b604080516000808252602082019092529061295c565b6040805180820190915260008152606060208201528152602001906001900390816129345790505b50611ebc565b5050565b60608989808060200260200160405190810160405280939291908181526020016000905b828210156129b6576129a760408302860136819003810190615cb7565b8152602001906001019061298a565b505050505090506000805b8a811015612b715760008888838181106129dd576129dd6154be565b90506020028101906129ef91906157b7565b8101906129fc9190615cd3565b90508451600014612a5757848281518110612a1957612a196154be565b602002602001015163ffffffff16600014612a5757848281518110612a4057612a406154be565b602090810291909101015163ffffffff1660608201525b612ad78d8d84818110612a6c57612a6c6154be565b905060400201602001358c8c848b8b88818110612a8b57612a8b6154be565b9050602002810190612a9d91906157b7565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061355b92505050565b848381518110612ae957612ae96154be565b6020026020010181905250612b25848381518110612b0957612b096154be565b602002602001015160000151600c61392790919063ffffffff16565b15612b6857612b5b848381518110612b3f57612b3f6154be565b6020908102919091010151600b546001600160a01b031661393c565b612b6590846156d4565b92505b506001016129c1565b508015612b8157612b8181613a5d565b509998505050505050505050565b6000612b9a83613a6a565b8015610e5d5750610e5d8383613ace565b8154600090612bd490700100000000000000000000000000000000900463ffffffff16426153fd565b90508015612c765760018301548354612c1c916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612e94565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612c9c916fffffffffffffffffffffffffffffffff9081169116613b9e565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990612d839084908151151581526020808301516fffffffffffffffffffffffffffffffff90811691830191909152604092830151169181019190915260600190565b60405180910390a1505050565b60006106a482613bb4565b6000808080612daa8686613bbf565b909450925050505b9250929050565b336001600160a01b03821603612e2b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016106ee565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612eb385612ea48486615466565b612eae90876156d4565b613b9e565b95945050505050565b6000610e5d8383613bce565b6000611e3384846001600160a01b038516613beb565b60008060001b8284602001518560400151866060015187608001518860a001518960c001518a60e001518b6101000151604051602001612f749897969594939291906001600160a01b039889168152968816602088015267ffffffffffffffff95861660408801526060870194909452911515608086015290921660a0840152921660c082015260e08101919091526101000190565b6040516020818303038152906040528051906020012085610120015180519060200120866101400151604051602001612fad9190615d8f565b60405160208183030381529060405280519060200120876101600151604051602001612fd99190615dfc565b60408051601f198184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e00160405160208183030381529060405280519060200120905092915050565b7f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff168467ffffffffffffffff16146130b2576040517f1279ec8a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff851660048201526024016106ee565b600a5468010000000000000000900461ffff1683111561310a576040517f099d3f7200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff861660048201526024016106ee565b80831461314f576040517f8808f8e700000000000000000000000000000000000000000000000000000000815267ffffffffffffffff861660048201526024016106ee565b600a54640100000000900463ffffffff168211156131b357600a546040517f8693378900000000000000000000000000000000000000000000000000000000815264010000000090910463ffffffff166004820152602481018390526044016106ee565b5050505050565b600060026131c960808561543f565b67ffffffffffffffff166131dd9190615466565b905060006010816131ef60808761547d565b67ffffffffffffffff168152602081019190915260400160002054905081613219600160046153fd565b901b191681836003811115613230576132306145ef565b901b17806010600061324360808861547d565b67ffffffffffffffff16815260208101919091526040016000205550505050565b6040517fb6113fce000000000000000000000000000000000000000000000000000000008152600090606090309063b6113fce906132aa90889088908890600401615e46565b600060405180830381600087803b1580156132c457600080fd5b505af19250505080156132d5575060015b61353e573d808015613303576040519150601f19603f3d011682016040523d82523d6000602084013e613308565b606091505b5061331281615fcf565b7fffffffff00000000000000000000000000000000000000000000000000000000167f0a8d6e8c0000000000000000000000000000000000000000000000000000000014806133aa575061336581615fcf565b7fffffffff00000000000000000000000000000000000000000000000000000000167fe1cd550900000000000000000000000000000000000000000000000000000000145b806133fe57506133b981615fcf565b7fffffffff00000000000000000000000000000000000000000000000000000000167f8d666f6000000000000000000000000000000000000000000000000000000000145b80613452575061340d81615fcf565b7fffffffff00000000000000000000000000000000000000000000000000000000167f78ef802400000000000000000000000000000000000000000000000000000000145b806134a6575061346181615fcf565b7fffffffff00000000000000000000000000000000000000000000000000000000167f0c3b563c00000000000000000000000000000000000000000000000000000000145b806134fa57506134b581615fcf565b7fffffffff00000000000000000000000000000000000000000000000000000000167fae9b4ce900000000000000000000000000000000000000000000000000000000145b1561350a57600392509050613553565b806040517fcf19edfd0000000000000000000000000000000000000000000000000000000081526004016106ee9190614690565b50506040805160208101909152600081526002905b935093915050565b6040805180820190915260008082526020820152600061357e8460200151613c08565b6040517fbbe4f6db0000000000000000000000000000000000000000000000000000000081526001600160a01b0380831660048301529192506000917f0000000000000000000000000000000000000000000000000000000000000000169063bbe4f6db90602401602060405180830381865afa158015613603573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613627919061601f565b90506001600160a01b038116158061366f575061366d6001600160a01b0382167faff2afbf00000000000000000000000000000000000000000000000000000000612b8f565b155b156136b1576040517fae9b4ce90000000000000000000000000000000000000000000000000000000081526001600160a01b03821660048201526024016106ee565b6000806136c98885896060015163ffffffff16613cae565b9150915060008060006137c76040518061010001604052808e81526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020018d6001600160a01b031681526020018f8152602001896001600160a01b031681526020018c6000015181526020018c6040015181526020018b815250604051602401613763919061603c565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f390775370000000000000000000000000000000000000000000000000000000017905287866113886084613df1565b9250925092508261380657816040517fe1cd55090000000000000000000000000000000000000000000000000000000081526004016106ee9190614690565b815160201461384e5781516040517f78ef80240000000000000000000000000000000000000000000000000000000081526020600482015260248101919091526044016106ee565b6000828060200190518101906138649190615bed565b9050866001600160a01b03168c6001600160a01b0316146138f95760006138958d8a613890868a6153fd565b613cae565b509050868110806138af5750816138ac88836153fd565b14155b156138f7576040517fa966e21f0000000000000000000000000000000000000000000000000000000081526004810183905260248101889052604481018290526064016106ee565b505b604080518082019091526001600160a01b039098168852602088015250949550505050505095945050505050565b6000610e5d836001600160a01b038416613f17565b81516040517fd02641a00000000000000000000000000000000000000000000000000000000081526001600160a01b03918216600482015260009182919084169063d02641a0906024016040805180830381865afa1580156139a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139c69190616109565b5190507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116600003613a2f5783516040517f9a655f7b0000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016106ee565b6020840151611e33907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff831690613f23565b61184f6003826000613f60565b6000613a96827f01ffc9a700000000000000000000000000000000000000000000000000000000613ace565b80156106a45750613ac7827fffffffff00000000000000000000000000000000000000000000000000000000613ace565b1592915050565b6040517fffffffff0000000000000000000000000000000000000000000000000000000082166024820152600090819060440160408051601f19818403018152919052602080820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a700000000000000000000000000000000000000000000000000000000178152825192935060009283928392909183918a617530fa92503d91506000519050828015613b87575060208210155b8015613b935750600081115b979650505050505050565b6000818310613bad5781610e5d565b5090919050565b60006106a4826142af565b6000808080612daa86866142ba565b60008181526002830160205260408120819055610e5d83836142e5565b60008281526002840160205260408120829055611e3384846142f1565b60008151602014613c4757816040517f8d666f600000000000000000000000000000000000000000000000000000000081526004016106ee9190614690565b600082806020019051810190613c5d9190615bed565b90506001600160a01b03811180613c75575061040081105b156106a457826040517f8d666f600000000000000000000000000000000000000000000000000000000081526004016106ee9190614690565b6000806000806000613d3d88604051602401613cd991906001600160a01b0391909116815260200190565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f70a082310000000000000000000000000000000000000000000000000000000017905288886113886084613df1565b92509250925082613d7c57816040517fe1cd55090000000000000000000000000000000000000000000000000000000081526004016106ee9190614690565b6020825114613dc45781516040517f78ef80240000000000000000000000000000000000000000000000000000000081526020600482015260248101919091526044016106ee565b81806020019051810190613dd89190615bed565b613de282886153fd565b94509450505050935093915050565b6000606060008361ffff1667ffffffffffffffff811115613e1457613e146146a3565b6040519080825280601f01601f191660200182016040528015613e3e576020820181803683370190505b509150863b613e71577f0c3b563c0000000000000000000000000000000000000000000000000000000060005260046000fd5b5a85811015613ea4577fafa32a2c0000000000000000000000000000000000000000000000000000000060005260046000fd5b8590036040810481038710613edd577f37c3be290000000000000000000000000000000000000000000000000000000060005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d84811115613f005750835b808352806000602085013e50955095509592505050565b6000610e5d83836142fd565b6000670de0b6b3a7640000613f56837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8616615466565b610e5d919061616b565b825474010000000000000000000000000000000000000000900460ff161580613f87575081155b15613f9157505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090613fd790700100000000000000000000000000000000900463ffffffff16426153fd565b905080156140975781831115614019576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546140539083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612e94565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015614134576001600160a01b0384166140e9576040517ff94ebcd100000000000000000000000000000000000000000000000000000000815260048101839052602481018690526044016106ee565b6040517f1a76572a00000000000000000000000000000000000000000000000000000000815260048101839052602481018690526001600160a01b03851660448201526064016106ee565b8483101561422d5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1690600090829061417890826153fd565b614182878a6153fd565b61418c91906156d4565b614196919061616b565b90506001600160a01b0386166141e2576040517f15279c0800000000000000000000000000000000000000000000000000000000815260048101829052602481018690526044016106ee565b6040517fd0c8d23a00000000000000000000000000000000000000000000000000000000815260048101829052602481018690526001600160a01b03871660448201526064016106ee565b61423785846153fd565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b60006106a482614309565b600080806142c88585614313565b600081815260029690960160205260409095205494959350505050565b6000610e5d838361431f565b6000610e5d8383614419565b6000610e5d8383614468565b60006106a4825490565b6000610e5d8383614480565b600081815260018301602052604081205480156144085760006143436001836153fd565b8554909150600090614357906001906153fd565b90508181146143bc576000866000018281548110614377576143776154be565b906000526020600020015490508087600001848154811061439a5761439a6154be565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806143cd576143cd61617f565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506106a4565b60009150506106a4565b5092915050565b6000818152600183016020526040812054614460575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556106a4565b5060006106a4565b60008181526001830160205260408120541515610e5d565b6000826000018281548110614497576144976154be565b9060005260206000200154905092915050565b828054828255906000526020600020908101928215614517579160200282015b8281111561451757825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b039091161782556020909201916001909101906144ca565b50614523929150614527565b5090565b5b808211156145235760008155600101614528565b60e081016106a482846001600160a01b03808251168352602082015167ffffffffffffffff808216602086015280604085015116604086015250508060608301511660608401528060808301511660808401528060a08301511660a08401528060c08301511660c0840152505050565b67ffffffffffffffff8116811461184f57600080fd5b80356145cd816145ac565b919050565b6000602082840312156145e457600080fd5b8135610e5d816145ac565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6004811061462e5761462e6145ef565b9052565b602081016106a4828461461e565b60005b8381101561465b578181015183820152602001614643565b50506000910152565b6000815180845261467c816020860160208601614640565b601f01601f19169290920160200192915050565b602081526000610e5d6020830184614664565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156146f5576146f56146a3565b60405290565b6040516080810167ffffffffffffffff811182821017156146f5576146f56146a3565b6040516101a0810167ffffffffffffffff811182821017156146f5576146f56146a3565b604051601f8201601f1916810167ffffffffffffffff8111828210171561476b5761476b6146a3565b604052919050565b600067ffffffffffffffff82111561478d5761478d6146a3565b5060051b60200190565b6001600160a01b038116811461184f57600080fd5b80356145cd81614797565b600082601f8301126147c857600080fd5b813560206147dd6147d883614773565b614742565b8083825260208201915060208460051b8701019350868411156147ff57600080fd5b602086015b8481101561482457803561481781614797565b8352918301918301614804565b509695505050505050565b803560ff811681146145cd57600080fd5b600067ffffffffffffffff82111561485a5761485a6146a3565b50601f01601f191660200190565b600082601f83011261487957600080fd5b81356148876147d882614840565b81815284602083860101111561489c57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c087890312156148d257600080fd5b863567ffffffffffffffff808211156148ea57600080fd5b6148f68a838b016147b7565b9750602089013591508082111561490c57600080fd5b6149188a838b016147b7565b965061492660408a0161482f565b9550606089013591508082111561493c57600080fd5b6149488a838b01614868565b945061495660808a016145c2565b935060a089013591508082111561496c57600080fd5b5061497989828a01614868565b9150509295509295509295565b60008151808452602080850194506020840160005b838110156149c05781516001600160a01b03168752958201959082019060010161499b565b509495945050505050565b602081526000610e5d6020830184614986565b6000602082840312156149f057600080fd5b8135610e5d81614797565b60a081016106a4828463ffffffff8082511683528060208301511660208401525061ffff604082015116604083015260608101516001600160a01b03808216606085015280608084015116608085015250505050565b600060208284031215614a6357600080fd5b813567ffffffffffffffff811115614a7a57600080fd5b820160a08185031215610e5d57600080fd5b600082601f830112614a9d57600080fd5b81356020614aad6147d883614773565b82815260069290921b84018101918181019086841115614acc57600080fd5b8286015b848110156148245760408189031215614ae95760008081fd5b614af16146d2565b8135614afc81614797565b815281850135614b0b81614797565b81860152835291830191604001614ad0565b60008060408385031215614b3057600080fd5b823567ffffffffffffffff80821115614b4857600080fd5b614b5486838701614a8c565b93506020850135915080821115614b6a57600080fd5b50614b7785828601614a8c565b9150509250929050565b801515811461184f57600080fd5b80356145cd81614b81565b600060408284031215614bac57600080fd5b614bb46146d2565b90508135614bc181614797565b808252506020820135602082015292915050565b600082601f830112614be657600080fd5b81356020614bf66147d883614773565b8083825260208201915060208460061b870101935086841115614c1857600080fd5b602086015b8481101561482457614c2f8882614b9a565b835291830191604001614c1d565b600082601f830112614c4e57600080fd5b81356020614c5e6147d883614773565b82815260059290921b84018101918181019086841115614c7d57600080fd5b8286015b8481101561482457803567ffffffffffffffff811115614ca15760008081fd5b614caf8986838b0101614868565b845250918301918301614c81565b600082601f830112614cce57600080fd5b81356020614cde6147d883614773565b82815260059290921b84018101918181019086841115614cfd57600080fd5b8286015b8481101561482457803567ffffffffffffffff811115614d215760008081fd5b614d2f8986838b0101614c3d565b845250918301918301614d01565b600082601f830112614d4e57600080fd5b81356020614d5e6147d883614773565b8083825260208201915060208460051b870101935086841115614d8057600080fd5b602086015b848110156148245780358352918301918301614d85565b600060808284031215614dae57600080fd5b614db66146fb565b9050813567ffffffffffffffff80821115614dd057600080fd5b818401915084601f830112614de457600080fd5b81356020614df46147d883614773565b82815260059290921b84018101918181019088841115614e1357600080fd5b8286015b84811015614f6f57803586811115614e2e57600080fd5b87016101a0818c03601f19011215614e4557600080fd5b614e4d61471e565b614e588683016145c2565b8152614e66604083016147ac565b86820152614e76606083016147ac565b6040820152614e87608083016145c2565b606082015260a08201356080820152614ea260c08301614b8f565b60a0820152614eb360e083016145c2565b60c0820152610100614ec68184016147ac565b60e083015261012080840135828401526101409150818401358a811115614eec57600080fd5b614efa8f8a83880101614868565b828501525050610160808401358a811115614f1457600080fd5b614f228f8a83880101614bd5565b83850152506101809150818401358a811115614f3d57600080fd5b614f4b8f8a83880101614c3d565b91840191909152506101a09290920135918101919091528352918301918301614e17565b5086525085810135935082841115614f8657600080fd5b614f9287858801614cbd565b90850152506040840135915080821115614fab57600080fd5b50614fb884828501614d3d565b6040830152506060820135606082015292915050565b63ffffffff8116811461184f57600080fd5b600082601f830112614ff157600080fd5b813560206150016147d883614773565b8083825260208201915060208460051b87010193508684111561502357600080fd5b602086015b8481101561482457803561503b81614fce565b8352918301918301615028565b600080604080848603121561505c57600080fd5b833567ffffffffffffffff8082111561507457600080fd5b61508087838801614d9c565b945060209150818601358181111561509757600080fd5b8601601f810188136150a857600080fd5b80356150b66147d882614773565b81815260059190911b8201840190848101908a8311156150d557600080fd5b8584015b83811015615148578035868111156150f15760008081fd5b8501808d03601f19018913156151075760008081fd5b61510f6146d2565b88820135815289820135888111156151275760008081fd5b6151358f8b83860101614fe0565b828b0152508452509186019186016150d9565b50809750505050505050509250929050565b60008083601f84011261516c57600080fd5b50813567ffffffffffffffff81111561518457600080fd5b6020830191508360208260051b8501011115612db257600080fd5b60008060008060008060008060e0898b0312156151bb57600080fd5b606089018a8111156151cc57600080fd5b8998503567ffffffffffffffff808211156151e657600080fd5b818b0191508b601f8301126151fa57600080fd5b81358181111561520957600080fd5b8c602082850101111561521b57600080fd5b6020830199508098505060808b013591508082111561523957600080fd5b6152458c838d0161515a565b909750955060a08b013591508082111561525e57600080fd5b5061526b8b828c0161515a565b999c989b50969995989497949560c00135949350505050565b6000806000806060858703121561529a57600080fd5b843567ffffffffffffffff808211156152b257600080fd5b908601906101a082890312156152c757600080fd5b909450602086013590808211156152dd57600080fd5b6152e98883890161515a565b9095509350604087013591508082111561530257600080fd5b5061530f87828801614fe0565b91505092959194509250565b80356fffffffffffffffffffffffffffffffff811681146145cd57600080fd5b60006060828403121561534d57600080fd5b6040516060810181811067ffffffffffffffff82111715615370576153706146a3565b604052823561537e81614b81565b815261538c6020840161531b565b602082015261539d6040840161531b565b60408201529392505050565b6040815260006153bc6040830185614986565b8281036020840152612eb38185614986565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156106a4576106a46153ce565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600067ffffffffffffffff8084168061545a5761545a615410565b92169190910692915050565b80820281158282048414176106a4576106a46153ce565b600067ffffffffffffffff8084168061549857615498615410565b92169190910492915050565b60208101600383106154b8576154b86145ef565b91905290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600063ffffffff808316818103615506576155066153ce565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b166040850152508060608401526155408184018a614986565b905082810360808401526155548189614986565b905060ff871660a084015282810360c08401526155718187614664565b905067ffffffffffffffff851660e08401528281036101008401526155968185614664565b9c9b505050505050505050505050565b6000602082840312156155b857600080fd5b8151610e5d816145ac565b600082601f8301126155d457600080fd5b81516155e26147d882614840565b8181528460208386010111156155f757600080fd5b611e33826020830160208701614640565b60006020828403121561561a57600080fd5b815167ffffffffffffffff8082111561563257600080fd5b908301906080828603121561564657600080fd5b61564e6146fb565b82518281111561565d57600080fd5b615669878286016155c3565b82525060208301518281111561567e57600080fd5b61568a878286016155c3565b6020830152506040830151828111156156a257600080fd5b6156ae878286016155c3565b604083015250606083015192506156c483614fce565b6060810192909252509392505050565b808201808211156106a4576106a46153ce565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261571c57600080fd5b83018035915067ffffffffffffffff82111561573757600080fd5b6020019150600681901b3603821315612db257600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261578457600080fd5b83018035915067ffffffffffffffff82111561579f57600080fd5b6020019150600581901b3603821315612db257600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126157ec57600080fd5b83018035915067ffffffffffffffff82111561580757600080fd5b602001915036819003821315612db257600080fd5b60008151808452602080850194506020840160005b838110156149c057815180516001600160a01b031688528301518388015260409096019590820190600101615831565b608081528451608082015267ffffffffffffffff60208601511660a08201526000604086015160a060c084015261589c610120840182614664565b905060608701517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80808584030160e08601526158d88383614664565b9250608089015191508085840301610100860152506158f7828261581c565b9250505061590b602083018661ffff169052565b836040830152612eb360608301846001600160a01b03169052565b60008060006060848603121561593b57600080fd5b835161594681614b81565b602085015190935067ffffffffffffffff81111561596357600080fd5b61596f868287016155c3565b925050604084015190509250925092565b600060a0828403121561599257600080fd5b60405160a0810181811067ffffffffffffffff821117156159b5576159b56146a3565b60405282516159c381614fce565b815260208301516159d381614fce565b6020820152604083015161ffff811681146159ed57600080fd5b60408201526060830151615a0081614797565b60608201526080830151615a1381614797565b60808201529392505050565b6101808101615a9082856001600160a01b03808251168352602082015167ffffffffffffffff808216602086015280604085015116604086015250508060608301511660608401528060808301511660808401528060a08301511660a08401528060c08301511660c0840152505050565b825163ffffffff90811660e0840152602084015116610100830152604083015161ffff1661012083015260608301516001600160a01b03908116610140840152608084015116610160830152610e5d565b60006101208b83526001600160a01b038b16602084015267ffffffffffffffff808b166040850152816060850152615b1b8285018b614986565b91508382036080850152615b2f828a614986565b915060ff881660a085015283820360c0850152615b4c8288614664565b90861660e085015283810361010085015290506155968185614664565b600060208284031215615b7b57600080fd5b8151610e5d81614b81565b60008151808452602080850194506020840160005b838110156149c057815187529582019590820190600101615b9b565b606081526000615bca6060830186615b86565b8281036020840152615bdc8186615b86565b915050826040830152949350505050565b600060208284031215615bff57600080fd5b5051919050565b67ffffffffffffffff818116838216019080821115614412576144126153ce565b67ffffffffffffffff8316815260408101610e5d602083018461461e565b600067ffffffffffffffff808316818103615506576155066153ce565b615c6c818461461e565b604060208201526000611e336040830184614664565b600060208284031215615c9457600080fd5b813567ffffffffffffffff811115615cab57600080fd5b611e3384828501614d9c565b600060408284031215615cc957600080fd5b610e5d8383614b9a565b600060208284031215615ce557600080fd5b813567ffffffffffffffff80821115615cfd57600080fd5b9083019060808286031215615d1157600080fd5b615d196146fb565b823582811115615d2857600080fd5b615d3487828601614868565b825250602083013582811115615d4957600080fd5b615d5587828601614868565b602083015250604083013582811115615d6d57600080fd5b615d7987828601614868565b604083015250606083013592506156c483614fce565b602081526000610e5d602083018461581c565b60008282518085526020808601955060208260051b8401016020860160005b84811015615def57601f19868403018952615ddd838351614664565b98840198925090830190600101615dc1565b5090979650505050505050565b602081526000610e5d6020830184615da2565b60008151808452602080850194506020840160005b838110156149c057815163ffffffff1687529582019590820190600101615e24565b60608152615e6160608201855167ffffffffffffffff169052565b60006020850151615e7d60808401826001600160a01b03169052565b5060408501516001600160a01b03811660a084015250606085015167ffffffffffffffff811660c084015250608085015160e083015260a0850151610100615ec88185018315159052565b60c08701519150610120615ee78186018467ffffffffffffffff169052565b60e08801519250610140615f05818701856001600160a01b03169052565b828901519350610160925083838701528189015193506101a091506101808281880152615f36610200880186614664565b9450818a015191507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0808887030184890152615f72868461581c565b948b01518886039091016101c0890152939450615f8f8585615da2565b9450808a01516101e0880152505050508281036020840152615fb18186615da2565b90508281036040840152615fc58185615e0f565b9695505050505050565b6000815160208301517fffffffff00000000000000000000000000000000000000000000000000000000808216935060048310156160175780818460040360031b1b83161693505b505050919050565b60006020828403121561603157600080fd5b8151610e5d81614797565b602081526000825161010080602085015261605b610120850183614664565b91506020850151616078604086018267ffffffffffffffff169052565b5060408501516001600160a01b0381166060860152506060850151608085015260808501516160b260a08601826001600160a01b03169052565b5060a0850151601f19808685030160c08701526160cf8483614664565b935060c08701519150808685030160e08701526160ec8483614664565b935060e0870151915080868503018387015250615fc58382614664565b60006040828403121561611b57600080fd5b6161236146d2565b82517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461614f57600080fd5b8152602083015161615f81614fce565b60208201529392505050565b60008261617a5761617a615410565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000818000a",
}

var EVM2EVMOffRampABI = EVM2EVMOffRampMetaData.ABI

var EVM2EVMOffRampBin = EVM2EVMOffRampMetaData.Bin

func DeployEVM2EVMOffRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig EVM2EVMOffRampStaticConfig, rateLimiterConfig RateLimiterConfig) (common.Address, *types.Transaction, *EVM2EVMOffRamp, error) {
	parsed, err := EVM2EVMOffRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EVM2EVMOffRampBin), backend, staticConfig, rateLimiterConfig)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EVM2EVMOffRamp{address: address, abi: *parsed, EVM2EVMOffRampCaller: EVM2EVMOffRampCaller{contract: contract}, EVM2EVMOffRampTransactor: EVM2EVMOffRampTransactor{contract: contract}, EVM2EVMOffRampFilterer: EVM2EVMOffRampFilterer{contract: contract}}, nil
}

type EVM2EVMOffRamp struct {
	address common.Address
	abi     abi.ABI
	EVM2EVMOffRampCaller
	EVM2EVMOffRampTransactor
	EVM2EVMOffRampFilterer
}

type EVM2EVMOffRampCaller struct {
	contract *bind.BoundContract
}

type EVM2EVMOffRampTransactor struct {
	contract *bind.BoundContract
}

type EVM2EVMOffRampFilterer struct {
	contract *bind.BoundContract
}

type EVM2EVMOffRampSession struct {
	Contract     *EVM2EVMOffRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type EVM2EVMOffRampCallerSession struct {
	Contract *EVM2EVMOffRampCaller
	CallOpts bind.CallOpts
}

type EVM2EVMOffRampTransactorSession struct {
	Contract     *EVM2EVMOffRampTransactor
	TransactOpts bind.TransactOpts
}

type EVM2EVMOffRampRaw struct {
	Contract *EVM2EVMOffRamp
}

type EVM2EVMOffRampCallerRaw struct {
	Contract *EVM2EVMOffRampCaller
}

type EVM2EVMOffRampTransactorRaw struct {
	Contract *EVM2EVMOffRampTransactor
}

func NewEVM2EVMOffRamp(address common.Address, backend bind.ContractBackend) (*EVM2EVMOffRamp, error) {
	abi, err := abi.JSON(strings.NewReader(EVM2EVMOffRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindEVM2EVMOffRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRamp{address: address, abi: abi, EVM2EVMOffRampCaller: EVM2EVMOffRampCaller{contract: contract}, EVM2EVMOffRampTransactor: EVM2EVMOffRampTransactor{contract: contract}, EVM2EVMOffRampFilterer: EVM2EVMOffRampFilterer{contract: contract}}, nil
}

func NewEVM2EVMOffRampCaller(address common.Address, caller bind.ContractCaller) (*EVM2EVMOffRampCaller, error) {
	contract, err := bindEVM2EVMOffRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampCaller{contract: contract}, nil
}

func NewEVM2EVMOffRampTransactor(address common.Address, transactor bind.ContractTransactor) (*EVM2EVMOffRampTransactor, error) {
	contract, err := bindEVM2EVMOffRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampTransactor{contract: contract}, nil
}

func NewEVM2EVMOffRampFilterer(address common.Address, filterer bind.ContractFilterer) (*EVM2EVMOffRampFilterer, error) {
	contract, err := bindEVM2EVMOffRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampFilterer{contract: contract}, nil
}

func bindEVM2EVMOffRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EVM2EVMOffRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMOffRamp.Contract.EVM2EVMOffRampCaller.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.EVM2EVMOffRampTransactor.contract.Transfer(opts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.EVM2EVMOffRampTransactor.contract.Transact(opts, method, params...)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMOffRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.contract.Transfer(opts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.contract.Transact(opts, method, params...)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "ccipReceive", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _EVM2EVMOffRamp.Contract.CcipReceive(&_EVM2EVMOffRamp.CallOpts, arg0)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _EVM2EVMOffRamp.Contract.CcipReceive(&_EVM2EVMOffRamp.CallOpts, arg0)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "currentRateLimiterState")

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMOffRamp.Contract.CurrentRateLimiterState(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMOffRamp.Contract.CurrentRateLimiterState(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) GetAllRateLimitTokens(opts *bind.CallOpts) (GetAllRateLimitTokens,

	error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "getAllRateLimitTokens")

	outstruct := new(GetAllRateLimitTokens)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SourceTokens = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.DestTokens = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) GetAllRateLimitTokens() (GetAllRateLimitTokens,

	error) {
	return _EVM2EVMOffRamp.Contract.GetAllRateLimitTokens(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) GetAllRateLimitTokens() (GetAllRateLimitTokens,

	error) {
	return _EVM2EVMOffRamp.Contract.GetAllRateLimitTokens(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMOffRampDynamicConfig, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(EVM2EVMOffRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMOffRampDynamicConfig)).(*EVM2EVMOffRampDynamicConfig)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) GetDynamicConfig() (EVM2EVMOffRampDynamicConfig, error) {
	return _EVM2EVMOffRamp.Contract.GetDynamicConfig(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) GetDynamicConfig() (EVM2EVMOffRampDynamicConfig, error) {
	return _EVM2EVMOffRamp.Contract.GetDynamicConfig(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) GetExecutionState(opts *bind.CallOpts, sequenceNumber uint64) (uint8, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "getExecutionState", sequenceNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) GetExecutionState(sequenceNumber uint64) (uint8, error) {
	return _EVM2EVMOffRamp.Contract.GetExecutionState(&_EVM2EVMOffRamp.CallOpts, sequenceNumber)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) GetExecutionState(sequenceNumber uint64) (uint8, error) {
	return _EVM2EVMOffRamp.Contract.GetExecutionState(&_EVM2EVMOffRamp.CallOpts, sequenceNumber)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) GetSenderNonce(opts *bind.CallOpts, sender common.Address) (uint64, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "getSenderNonce", sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) GetSenderNonce(sender common.Address) (uint64, error) {
	return _EVM2EVMOffRamp.Contract.GetSenderNonce(&_EVM2EVMOffRamp.CallOpts, sender)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) GetSenderNonce(sender common.Address) (uint64, error) {
	return _EVM2EVMOffRamp.Contract.GetSenderNonce(&_EVM2EVMOffRamp.CallOpts, sender)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) GetStaticConfig(opts *bind.CallOpts) (EVM2EVMOffRampStaticConfig, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(EVM2EVMOffRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMOffRampStaticConfig)).(*EVM2EVMOffRampStaticConfig)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) GetStaticConfig() (EVM2EVMOffRampStaticConfig, error) {
	return _EVM2EVMOffRamp.Contract.GetStaticConfig(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) GetStaticConfig() (EVM2EVMOffRampStaticConfig, error) {
	return _EVM2EVMOffRamp.Contract.GetStaticConfig(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "getTokenLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMOffRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMOffRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) GetTransmitters() ([]common.Address, error) {
	return _EVM2EVMOffRamp.Contract.GetTransmitters(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) GetTransmitters() ([]common.Address, error) {
	return _EVM2EVMOffRamp.Contract.GetTransmitters(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _EVM2EVMOffRamp.Contract.LatestConfigDetails(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _EVM2EVMOffRamp.Contract.LatestConfigDetails(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _EVM2EVMOffRamp.Contract.LatestConfigDigestAndEpoch(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _EVM2EVMOffRamp.Contract.LatestConfigDigestAndEpoch(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) Owner() (common.Address, error) {
	return _EVM2EVMOffRamp.Contract.Owner(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) Owner() (common.Address, error) {
	return _EVM2EVMOffRamp.Contract.Owner(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EVM2EVMOffRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) TypeAndVersion() (string, error) {
	return _EVM2EVMOffRamp.Contract.TypeAndVersion(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampCallerSession) TypeAndVersion() (string, error) {
	return _EVM2EVMOffRamp.Contract.TypeAndVersion(&_EVM2EVMOffRamp.CallOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "acceptOwnership")
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.AcceptOwnership(&_EVM2EVMOffRamp.TransactOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.AcceptOwnership(&_EVM2EVMOffRamp.TransactOpts)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalEVM2EVMMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData, tokenGasOverrides)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) ExecuteSingleMessage(message InternalEVM2EVMMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.ExecuteSingleMessage(&_EVM2EVMOffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) ExecuteSingleMessage(message InternalEVM2EVMMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.ExecuteSingleMessage(&_EVM2EVMOffRamp.TransactOpts, message, offchainTokenData, tokenGasOverrides)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) ManuallyExecute(opts *bind.TransactOpts, report InternalExecutionReport, gasLimitOverrides []EVM2EVMOffRampGasLimitOverride) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "manuallyExecute", report, gasLimitOverrides)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) ManuallyExecute(report InternalExecutionReport, gasLimitOverrides []EVM2EVMOffRampGasLimitOverride) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.ManuallyExecute(&_EVM2EVMOffRamp.TransactOpts, report, gasLimitOverrides)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) ManuallyExecute(report InternalExecutionReport, gasLimitOverrides []EVM2EVMOffRampGasLimitOverride) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.ManuallyExecute(&_EVM2EVMOffRamp.TransactOpts, report, gasLimitOverrides)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "setAdmin", newAdmin)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.SetAdmin(&_EVM2EVMOffRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.SetAdmin(&_EVM2EVMOffRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) SetOCR2Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "setOCR2Config", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) SetOCR2Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.SetOCR2Config(&_EVM2EVMOffRamp.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) SetOCR2Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.SetOCR2Config(&_EVM2EVMOffRamp.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "setRateLimiterConfig", config)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.SetRateLimiterConfig(&_EVM2EVMOffRamp.TransactOpts, config)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.SetRateLimiterConfig(&_EVM2EVMOffRamp.TransactOpts, config)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.TransferOwnership(&_EVM2EVMOffRamp.TransactOpts, to)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.TransferOwnership(&_EVM2EVMOffRamp.TransactOpts, to)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "transmit", reportContext, report, rs, ss, arg4)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.Transmit(&_EVM2EVMOffRamp.TransactOpts, reportContext, report, rs, ss, arg4)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.Transmit(&_EVM2EVMOffRamp.TransactOpts, reportContext, report, rs, ss, arg4)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactor) UpdateRateLimitTokens(opts *bind.TransactOpts, removes []EVM2EVMOffRampRateLimitToken, adds []EVM2EVMOffRampRateLimitToken) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.contract.Transact(opts, "updateRateLimitTokens", removes, adds)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampSession) UpdateRateLimitTokens(removes []EVM2EVMOffRampRateLimitToken, adds []EVM2EVMOffRampRateLimitToken) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.UpdateRateLimitTokens(&_EVM2EVMOffRamp.TransactOpts, removes, adds)
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampTransactorSession) UpdateRateLimitTokens(removes []EVM2EVMOffRampRateLimitToken, adds []EVM2EVMOffRampRateLimitToken) (*types.Transaction, error) {
	return _EVM2EVMOffRamp.Contract.UpdateRateLimitTokens(&_EVM2EVMOffRamp.TransactOpts, removes, adds)
}

type EVM2EVMOffRampAdminSetIterator struct {
	Event *EVM2EVMOffRampAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampAdminSet)
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
		it.Event = new(EVM2EVMOffRampAdminSet)
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

func (it *EVM2EVMOffRampAdminSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampAdminSet struct {
	NewAdmin common.Address
	Raw      types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMOffRampAdminSetIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampAdminSetIterator{contract: _EVM2EVMOffRamp.contract, event: "AdminSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampAdminSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampAdminSet)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseAdminSet(log types.Log) (*EVM2EVMOffRampAdminSet, error) {
	event := new(EVM2EVMOffRampAdminSet)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampAlreadyAttemptedIterator struct {
	Event *EVM2EVMOffRampAlreadyAttempted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampAlreadyAttemptedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampAlreadyAttempted)
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
		it.Event = new(EVM2EVMOffRampAlreadyAttempted)
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

func (it *EVM2EVMOffRampAlreadyAttemptedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampAlreadyAttemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampAlreadyAttempted struct {
	SequenceNumber uint64
	Raw            types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterAlreadyAttempted(opts *bind.FilterOpts) (*EVM2EVMOffRampAlreadyAttemptedIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampAlreadyAttemptedIterator{contract: _EVM2EVMOffRamp.contract, event: "AlreadyAttempted", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampAlreadyAttempted) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "AlreadyAttempted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampAlreadyAttempted)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseAlreadyAttempted(log types.Log) (*EVM2EVMOffRampAlreadyAttempted, error) {
	event := new(EVM2EVMOffRampAlreadyAttempted)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "AlreadyAttempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampConfigChangedIterator struct {
	Event *EVM2EVMOffRampConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampConfigChanged)
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
		it.Event = new(EVM2EVMOffRampConfigChanged)
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

func (it *EVM2EVMOffRampConfigChangedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*EVM2EVMOffRampConfigChangedIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampConfigChangedIterator{contract: _EVM2EVMOffRamp.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampConfigChanged) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampConfigChanged)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseConfigChanged(log types.Log) (*EVM2EVMOffRampConfigChanged, error) {
	event := new(EVM2EVMOffRampConfigChanged)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampConfigSetIterator struct {
	Event *EVM2EVMOffRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampConfigSet)
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
		it.Event = new(EVM2EVMOffRampConfigSet)
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

func (it *EVM2EVMOffRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampConfigSet struct {
	StaticConfig  EVM2EVMOffRampStaticConfig
	DynamicConfig EVM2EVMOffRampDynamicConfig
	Raw           types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMOffRampConfigSetIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampConfigSetIterator{contract: _EVM2EVMOffRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampConfigSet)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseConfigSet(log types.Log) (*EVM2EVMOffRampConfigSet, error) {
	event := new(EVM2EVMOffRampConfigSet)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampConfigSet0Iterator struct {
	Event *EVM2EVMOffRampConfigSet0

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampConfigSet0Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampConfigSet0)
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
		it.Event = new(EVM2EVMOffRampConfigSet0)
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

func (it *EVM2EVMOffRampConfigSet0Iterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampConfigSet0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampConfigSet0 struct {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterConfigSet0(opts *bind.FilterOpts) (*EVM2EVMOffRampConfigSet0Iterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "ConfigSet0")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampConfigSet0Iterator{contract: _EVM2EVMOffRamp.contract, event: "ConfigSet0", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchConfigSet0(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampConfigSet0) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "ConfigSet0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampConfigSet0)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ConfigSet0", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseConfigSet0(log types.Log) (*EVM2EVMOffRampConfigSet0, error) {
	event := new(EVM2EVMOffRampConfigSet0)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ConfigSet0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampExecutionStateChangedIterator struct {
	Event *EVM2EVMOffRampExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampExecutionStateChanged)
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
		it.Event = new(EVM2EVMOffRampExecutionStateChanged)
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

func (it *EVM2EVMOffRampExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampExecutionStateChanged struct {
	SequenceNumber uint64
	MessageId      [32]byte
	State          uint8
	ReturnData     []byte
	Raw            types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sequenceNumber []uint64, messageId [][32]byte) (*EVM2EVMOffRampExecutionStateChangedIterator, error) {

	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "ExecutionStateChanged", sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampExecutionStateChangedIterator{contract: _EVM2EVMOffRamp.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampExecutionStateChanged, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "ExecutionStateChanged", sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampExecutionStateChanged)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseExecutionStateChanged(log types.Log) (*EVM2EVMOffRampExecutionStateChanged, error) {
	event := new(EVM2EVMOffRampExecutionStateChanged)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampOwnershipTransferRequestedIterator struct {
	Event *EVM2EVMOffRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampOwnershipTransferRequested)
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
		it.Event = new(EVM2EVMOffRampOwnershipTransferRequested)
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

func (it *EVM2EVMOffRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOffRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampOwnershipTransferRequestedIterator{contract: _EVM2EVMOffRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampOwnershipTransferRequested)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMOffRampOwnershipTransferRequested, error) {
	event := new(EVM2EVMOffRampOwnershipTransferRequested)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampOwnershipTransferredIterator struct {
	Event *EVM2EVMOffRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampOwnershipTransferred)
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
		it.Event = new(EVM2EVMOffRampOwnershipTransferred)
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

func (it *EVM2EVMOffRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOffRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampOwnershipTransferredIterator{contract: _EVM2EVMOffRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampOwnershipTransferred)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseOwnershipTransferred(log types.Log) (*EVM2EVMOffRampOwnershipTransferred, error) {
	event := new(EVM2EVMOffRampOwnershipTransferred)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampSkippedAlreadyExecutedMessageIterator struct {
	Event *EVM2EVMOffRampSkippedAlreadyExecutedMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampSkippedAlreadyExecutedMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampSkippedAlreadyExecutedMessage)
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
		it.Event = new(EVM2EVMOffRampSkippedAlreadyExecutedMessage)
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

func (it *EVM2EVMOffRampSkippedAlreadyExecutedMessageIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampSkippedAlreadyExecutedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampSkippedAlreadyExecutedMessage struct {
	SequenceNumber uint64
	Raw            types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts, sequenceNumber []uint64) (*EVM2EVMOffRampSkippedAlreadyExecutedMessageIterator, error) {

	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "SkippedAlreadyExecutedMessage", sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampSkippedAlreadyExecutedMessageIterator{contract: _EVM2EVMOffRamp.contract, event: "SkippedAlreadyExecutedMessage", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampSkippedAlreadyExecutedMessage, sequenceNumber []uint64) (event.Subscription, error) {

	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "SkippedAlreadyExecutedMessage", sequenceNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampSkippedAlreadyExecutedMessage)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseSkippedAlreadyExecutedMessage(log types.Log) (*EVM2EVMOffRampSkippedAlreadyExecutedMessage, error) {
	event := new(EVM2EVMOffRampSkippedAlreadyExecutedMessage)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "SkippedAlreadyExecutedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampSkippedIncorrectNonceIterator struct {
	Event *EVM2EVMOffRampSkippedIncorrectNonce

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampSkippedIncorrectNonceIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampSkippedIncorrectNonce)
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
		it.Event = new(EVM2EVMOffRampSkippedIncorrectNonce)
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

func (it *EVM2EVMOffRampSkippedIncorrectNonceIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampSkippedIncorrectNonceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampSkippedIncorrectNonce struct {
	Nonce  uint64
	Sender common.Address
	Raw    types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterSkippedIncorrectNonce(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMOffRampSkippedIncorrectNonceIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "SkippedIncorrectNonce", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampSkippedIncorrectNonceIterator{contract: _EVM2EVMOffRamp.contract, event: "SkippedIncorrectNonce", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchSkippedIncorrectNonce(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampSkippedIncorrectNonce, nonce []uint64, sender []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "SkippedIncorrectNonce", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampSkippedIncorrectNonce)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "SkippedIncorrectNonce", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseSkippedIncorrectNonce(log types.Log) (*EVM2EVMOffRampSkippedIncorrectNonce, error) {
	event := new(EVM2EVMOffRampSkippedIncorrectNonce)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "SkippedIncorrectNonce", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflightIterator struct {
	Event *EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflightIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight)
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
		it.Event = new(EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight)
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

func (it *EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflightIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflightIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight struct {
	Nonce  uint64
	Sender common.Address
	Raw    types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterSkippedSenderWithPreviousRampMessageInflight(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflightIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "SkippedSenderWithPreviousRampMessageInflight", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflightIterator{contract: _EVM2EVMOffRamp.contract, event: "SkippedSenderWithPreviousRampMessageInflight", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchSkippedSenderWithPreviousRampMessageInflight(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight, nonce []uint64, sender []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "SkippedSenderWithPreviousRampMessageInflight", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "SkippedSenderWithPreviousRampMessageInflight", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseSkippedSenderWithPreviousRampMessageInflight(log types.Log) (*EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight, error) {
	event := new(EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "SkippedSenderWithPreviousRampMessageInflight", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampTokenAggregateRateLimitAddedIterator struct {
	Event *EVM2EVMOffRampTokenAggregateRateLimitAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampTokenAggregateRateLimitAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampTokenAggregateRateLimitAdded)
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
		it.Event = new(EVM2EVMOffRampTokenAggregateRateLimitAdded)
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

func (it *EVM2EVMOffRampTokenAggregateRateLimitAddedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampTokenAggregateRateLimitAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampTokenAggregateRateLimitAdded struct {
	SourceToken common.Address
	DestToken   common.Address
	Raw         types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterTokenAggregateRateLimitAdded(opts *bind.FilterOpts) (*EVM2EVMOffRampTokenAggregateRateLimitAddedIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "TokenAggregateRateLimitAdded")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampTokenAggregateRateLimitAddedIterator{contract: _EVM2EVMOffRamp.contract, event: "TokenAggregateRateLimitAdded", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchTokenAggregateRateLimitAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTokenAggregateRateLimitAdded) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "TokenAggregateRateLimitAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampTokenAggregateRateLimitAdded)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitAdded", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseTokenAggregateRateLimitAdded(log types.Log) (*EVM2EVMOffRampTokenAggregateRateLimitAdded, error) {
	event := new(EVM2EVMOffRampTokenAggregateRateLimitAdded)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampTokenAggregateRateLimitRemovedIterator struct {
	Event *EVM2EVMOffRampTokenAggregateRateLimitRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampTokenAggregateRateLimitRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampTokenAggregateRateLimitRemoved)
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
		it.Event = new(EVM2EVMOffRampTokenAggregateRateLimitRemoved)
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

func (it *EVM2EVMOffRampTokenAggregateRateLimitRemovedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampTokenAggregateRateLimitRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampTokenAggregateRateLimitRemoved struct {
	SourceToken common.Address
	DestToken   common.Address
	Raw         types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterTokenAggregateRateLimitRemoved(opts *bind.FilterOpts) (*EVM2EVMOffRampTokenAggregateRateLimitRemovedIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "TokenAggregateRateLimitRemoved")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampTokenAggregateRateLimitRemovedIterator{contract: _EVM2EVMOffRamp.contract, event: "TokenAggregateRateLimitRemoved", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchTokenAggregateRateLimitRemoved(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTokenAggregateRateLimitRemoved) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "TokenAggregateRateLimitRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampTokenAggregateRateLimitRemoved)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitRemoved", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseTokenAggregateRateLimitRemoved(log types.Log) (*EVM2EVMOffRampTokenAggregateRateLimitRemoved, error) {
	event := new(EVM2EVMOffRampTokenAggregateRateLimitRemoved)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampTokensConsumedIterator struct {
	Event *EVM2EVMOffRampTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampTokensConsumed)
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
		it.Event = new(EVM2EVMOffRampTokensConsumed)
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

func (it *EVM2EVMOffRampTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*EVM2EVMOffRampTokensConsumedIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampTokensConsumedIterator{contract: _EVM2EVMOffRamp.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampTokensConsumed)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseTokensConsumed(log types.Log) (*EVM2EVMOffRampTokensConsumed, error) {
	event := new(EVM2EVMOffRampTokensConsumed)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMOffRampTransmittedIterator struct {
	Event *EVM2EVMOffRampTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMOffRampTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMOffRampTransmitted)
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
		it.Event = new(EVM2EVMOffRampTransmitted)
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

func (it *EVM2EVMOffRampTransmittedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMOffRampTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMOffRampTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) FilterTransmitted(opts *bind.FilterOpts) (*EVM2EVMOffRampTransmittedIterator, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMOffRampTransmittedIterator{contract: _EVM2EVMOffRamp.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTransmitted) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMOffRamp.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMOffRampTransmitted)
				if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_EVM2EVMOffRamp *EVM2EVMOffRampFilterer) ParseTransmitted(log types.Log) (*EVM2EVMOffRampTransmitted, error) {
	event := new(EVM2EVMOffRampTransmitted)
	if err := _EVM2EVMOffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetAllRateLimitTokens struct {
	SourceTokens []common.Address
	DestTokens   []common.Address
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

func (_EVM2EVMOffRamp *EVM2EVMOffRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _EVM2EVMOffRamp.abi.Events["AdminSet"].ID:
		return _EVM2EVMOffRamp.ParseAdminSet(log)
	case _EVM2EVMOffRamp.abi.Events["AlreadyAttempted"].ID:
		return _EVM2EVMOffRamp.ParseAlreadyAttempted(log)
	case _EVM2EVMOffRamp.abi.Events["ConfigChanged"].ID:
		return _EVM2EVMOffRamp.ParseConfigChanged(log)
	case _EVM2EVMOffRamp.abi.Events["ConfigSet"].ID:
		return _EVM2EVMOffRamp.ParseConfigSet(log)
	case _EVM2EVMOffRamp.abi.Events["ConfigSet0"].ID:
		return _EVM2EVMOffRamp.ParseConfigSet0(log)
	case _EVM2EVMOffRamp.abi.Events["ExecutionStateChanged"].ID:
		return _EVM2EVMOffRamp.ParseExecutionStateChanged(log)
	case _EVM2EVMOffRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _EVM2EVMOffRamp.ParseOwnershipTransferRequested(log)
	case _EVM2EVMOffRamp.abi.Events["OwnershipTransferred"].ID:
		return _EVM2EVMOffRamp.ParseOwnershipTransferred(log)
	case _EVM2EVMOffRamp.abi.Events["SkippedAlreadyExecutedMessage"].ID:
		return _EVM2EVMOffRamp.ParseSkippedAlreadyExecutedMessage(log)
	case _EVM2EVMOffRamp.abi.Events["SkippedIncorrectNonce"].ID:
		return _EVM2EVMOffRamp.ParseSkippedIncorrectNonce(log)
	case _EVM2EVMOffRamp.abi.Events["SkippedSenderWithPreviousRampMessageInflight"].ID:
		return _EVM2EVMOffRamp.ParseSkippedSenderWithPreviousRampMessageInflight(log)
	case _EVM2EVMOffRamp.abi.Events["TokenAggregateRateLimitAdded"].ID:
		return _EVM2EVMOffRamp.ParseTokenAggregateRateLimitAdded(log)
	case _EVM2EVMOffRamp.abi.Events["TokenAggregateRateLimitRemoved"].ID:
		return _EVM2EVMOffRamp.ParseTokenAggregateRateLimitRemoved(log)
	case _EVM2EVMOffRamp.abi.Events["TokensConsumed"].ID:
		return _EVM2EVMOffRamp.ParseTokensConsumed(log)
	case _EVM2EVMOffRamp.abi.Events["Transmitted"].ID:
		return _EVM2EVMOffRamp.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (EVM2EVMOffRampAdminSet) Topic() common.Hash {
	return common.HexToHash("0x8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c")
}

func (EVM2EVMOffRampAlreadyAttempted) Topic() common.Hash {
	return common.HexToHash("0x67d9ba0f63d427c482c2736300e6d5a34c6691dbcdea8ad35828a1f1ba47e872")
}

func (EVM2EVMOffRampConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (EVM2EVMOffRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0x7879e20bb60a503429de4a2c912b5904f08a39f2af054c10fb46434b5d611260")
}

func (EVM2EVMOffRampConfigSet0) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (EVM2EVMOffRampExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0xd4f851956a5d67c3997d1c9205045fef79bae2947fdee7e9e2641abc7391ef65")
}

func (EVM2EVMOffRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (EVM2EVMOffRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (EVM2EVMOffRampSkippedAlreadyExecutedMessage) Topic() common.Hash {
	return common.HexToHash("0xe3dd0bec917c965a133ddb2c84874725ee1e2fd8d763c19efa36d6a11cd82b1f")
}

func (EVM2EVMOffRampSkippedIncorrectNonce) Topic() common.Hash {
	return common.HexToHash("0xd32ddb11d71e3d63411d37b09f9a8b28664f1cb1338bfd1413c173b0ebf41237")
}

func (EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight) Topic() common.Hash {
	return common.HexToHash("0xe44a20935573a783dd0d5991c92d7b6a0eb3173566530364db3ec10e9a990b5d")
}

func (EVM2EVMOffRampTokenAggregateRateLimitAdded) Topic() common.Hash {
	return common.HexToHash("0xfc23abf7ddbd3c02b1420dafa2355c56c1a06fbb8723862ac14d6bd74177361a")
}

func (EVM2EVMOffRampTokenAggregateRateLimitRemoved) Topic() common.Hash {
	return common.HexToHash("0xcbf3cbeaed4ac1d605ed30f4af06c35acaeff2379db7f6146c9cceee83d58782")
}

func (EVM2EVMOffRampTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (EVM2EVMOffRampTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_EVM2EVMOffRamp *EVM2EVMOffRamp) Address() common.Address {
	return _EVM2EVMOffRamp.address
}

type EVM2EVMOffRampInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error)

	GetAllRateLimitTokens(opts *bind.CallOpts) (GetAllRateLimitTokens,

		error)

	GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMOffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sequenceNumber uint64) (uint8, error)

	GetSenderNonce(opts *bind.CallOpts, sender common.Address) (uint64, error)

	GetStaticConfig(opts *bind.CallOpts) (EVM2EVMOffRampStaticConfig, error)

	GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetTransmitters(opts *bind.CallOpts) ([]common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ExecuteSingleMessage(opts *bind.TransactOpts, message InternalEVM2EVMMessage, offchainTokenData [][]byte, tokenGasOverrides []uint32) (*types.Transaction, error)

	ManuallyExecute(opts *bind.TransactOpts, report InternalExecutionReport, gasLimitOverrides []EVM2EVMOffRampGasLimitOverride) (*types.Transaction, error)

	SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	SetOCR2Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error)

	UpdateRateLimitTokens(opts *bind.TransactOpts, removes []EVM2EVMOffRampRateLimitToken, adds []EVM2EVMOffRampRateLimitToken) (*types.Transaction, error)

	FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMOffRampAdminSetIterator, error)

	WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampAdminSet) (event.Subscription, error)

	ParseAdminSet(log types.Log) (*EVM2EVMOffRampAdminSet, error)

	FilterAlreadyAttempted(opts *bind.FilterOpts) (*EVM2EVMOffRampAlreadyAttemptedIterator, error)

	WatchAlreadyAttempted(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampAlreadyAttempted) (event.Subscription, error)

	ParseAlreadyAttempted(log types.Log) (*EVM2EVMOffRampAlreadyAttempted, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*EVM2EVMOffRampConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*EVM2EVMOffRampConfigChanged, error)

	FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMOffRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*EVM2EVMOffRampConfigSet, error)

	FilterConfigSet0(opts *bind.FilterOpts) (*EVM2EVMOffRampConfigSet0Iterator, error)

	WatchConfigSet0(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampConfigSet0) (event.Subscription, error)

	ParseConfigSet0(log types.Log) (*EVM2EVMOffRampConfigSet0, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sequenceNumber []uint64, messageId [][32]byte) (*EVM2EVMOffRampExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampExecutionStateChanged, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*EVM2EVMOffRampExecutionStateChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOffRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMOffRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMOffRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*EVM2EVMOffRampOwnershipTransferred, error)

	FilterSkippedAlreadyExecutedMessage(opts *bind.FilterOpts, sequenceNumber []uint64) (*EVM2EVMOffRampSkippedAlreadyExecutedMessageIterator, error)

	WatchSkippedAlreadyExecutedMessage(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampSkippedAlreadyExecutedMessage, sequenceNumber []uint64) (event.Subscription, error)

	ParseSkippedAlreadyExecutedMessage(log types.Log) (*EVM2EVMOffRampSkippedAlreadyExecutedMessage, error)

	FilterSkippedIncorrectNonce(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMOffRampSkippedIncorrectNonceIterator, error)

	WatchSkippedIncorrectNonce(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampSkippedIncorrectNonce, nonce []uint64, sender []common.Address) (event.Subscription, error)

	ParseSkippedIncorrectNonce(log types.Log) (*EVM2EVMOffRampSkippedIncorrectNonce, error)

	FilterSkippedSenderWithPreviousRampMessageInflight(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflightIterator, error)

	WatchSkippedSenderWithPreviousRampMessageInflight(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight, nonce []uint64, sender []common.Address) (event.Subscription, error)

	ParseSkippedSenderWithPreviousRampMessageInflight(log types.Log) (*EVM2EVMOffRampSkippedSenderWithPreviousRampMessageInflight, error)

	FilterTokenAggregateRateLimitAdded(opts *bind.FilterOpts) (*EVM2EVMOffRampTokenAggregateRateLimitAddedIterator, error)

	WatchTokenAggregateRateLimitAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTokenAggregateRateLimitAdded) (event.Subscription, error)

	ParseTokenAggregateRateLimitAdded(log types.Log) (*EVM2EVMOffRampTokenAggregateRateLimitAdded, error)

	FilterTokenAggregateRateLimitRemoved(opts *bind.FilterOpts) (*EVM2EVMOffRampTokenAggregateRateLimitRemovedIterator, error)

	WatchTokenAggregateRateLimitRemoved(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTokenAggregateRateLimitRemoved) (event.Subscription, error)

	ParseTokenAggregateRateLimitRemoved(log types.Log) (*EVM2EVMOffRampTokenAggregateRateLimitRemoved, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*EVM2EVMOffRampTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*EVM2EVMOffRampTokensConsumed, error)

	FilterTransmitted(opts *bind.FilterOpts) (*EVM2EVMOffRampTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *EVM2EVMOffRampTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*EVM2EVMOffRampTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
