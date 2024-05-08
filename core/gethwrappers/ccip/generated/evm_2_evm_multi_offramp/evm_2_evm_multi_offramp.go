// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm_2_evm_multi_offramp

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

type EVM2EVMMultiOffRampDynamicConfig struct {
	PermissionLessExecutionThresholdSeconds uint32
	Router                                  common.Address
	PriceRegistry                           common.Address
	MaxNumberOfTokensPerMsg                 uint16
	MaxDataBytes                            uint32
	MaxPoolReleaseOrMintGas                 uint32
}

type EVM2EVMMultiOffRampRateLimitToken struct {
	SourceToken common.Address
	DestToken   common.Address
}

type EVM2EVMMultiOffRampSourceChainConfig struct {
	IsEnabled    bool
	PrevOffRamp  common.Address
	OnRamp       common.Address
	MetadataHash [32]byte
}

type EVM2EVMMultiOffRampSourceChainConfigArgs struct {
	SourceChainSelector uint64
	IsEnabled           bool
	PrevOffRamp         common.Address
	OnRamp              common.Address
}

type EVM2EVMMultiOffRampStaticConfig struct {
	CommitStore   common.Address
	ChainSelector uint64
	ArmProxy      common.Address
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

var EVM2EVMMultiOffRampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyAttempted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"AlreadyExecuted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadARMSignal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CanOnlySelfCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CommitStoreAlreadyInUse\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"ConfigDigestMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyReport\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"error\",\"type\":\"bytes\"}],\"name\":\"ExecutionError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"ForkedChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"got\",\"type\":\"uint256\"}],\"name\":\"InvalidDataLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidEVMAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"InvalidManualExecutionGasLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMessageId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"InvalidNewState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"InvalidSourceChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionGasLimitMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManualExecutionNotYetEnabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualSize\",\"type\":\"uint256\"}],\"name\":\"MessageTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"notPool\",\"type\":\"address\"}],\"name\":\"NotACompatiblePool\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByAdminOrOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OracleCannotBeZeroAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PriceNotFoundForToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"error\",\"type\":\"bytes\"}],\"name\":\"ReceiverError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RootNotCommitted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"TokenDataMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"error\",\"type\":\"bytes\"}],\"name\":\"TokenHandlingError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedTransmitter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedTokenData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"UnsupportedNumberOfTokens\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"WrongMessageLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.StaticConfig\",\"name\":\"staticConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.DynamicConfig\",\"name\":\"dynamicConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"state\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"ExecutionStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SkippedIncorrectNonce\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SkippedSenderWithPreviousRampMessageInflight\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"metadataHash\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfig\",\"name\":\"sourceConfig\",\"type\":\"tuple\"}],\"name\":\"SourceChainConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"SourceChainSelectorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"name\":\"TokenAggregateRateLimitAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"name\":\"TokenAggregateRateLimitRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfigArgs[]\",\"name\":\"sourceChainConfigUpdates\",\"type\":\"tuple[]\"}],\"name\":\"applySourceChainConfigUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"strict\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"sourceTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.EVM2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[]\"}],\"name\":\"executeSingleMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllRateLimitTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"sourceTokens\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"destTokens\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynamicConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"permissionLessExecutionThresholdSeconds\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"priceRegistry\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"maxNumberOfTokensPerMsg\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxDataBytes\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxPoolReleaseOrMintGas\",\"type\":\"uint32\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.DynamicConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"}],\"name\":\"getExecutionState\",\"outputs\":[{\"internalType\":\"enumInternal.MessageExecutionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getSenderNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"}],\"name\":\"getSourceChainConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"prevOffRamp\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"metadataHash\",\"type\":\"bytes32\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.SourceChainConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStaticConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"commitStore\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"armProxy\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.StaticConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"strict\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"tokenAmounts\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"sourceTokenData\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"internalType\":\"structInternal.EVM2EVMMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[][]\",\"name\":\"offchainTokenData\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proofs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"proofFlagBits\",\"type\":\"uint256\"}],\"internalType\":\"structInternal.ExecutionReport\",\"name\":\"report\",\"type\":\"tuple\"},{\"internalType\":\"uint256[]\",\"name\":\"gasLimitOverrides\",\"type\":\"uint256[]\"}],\"name\":\"manuallyExecute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setOCR2Config\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.RateLimitToken[]\",\"name\":\"removes\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"sourceToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destToken\",\"type\":\"address\"}],\"internalType\":\"structEVM2EVMMultiOffRamp.RateLimitToken[]\",\"name\":\"adds\",\"type\":\"tuple[]\"}],\"name\":\"updateRateLimitTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162006a6238038062006a62833981016040819052620000359162000639565b8033806000816200008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c057620000c081620001b5565b50506040805160a081018252602084810180516001600160801b039081168085524263ffffffff169385018490528751151585870181905292518216606086018190529790950151166080938401819052600380546001600160a01b031916909517600160801b9384021760ff60a01b1916600160a01b909202919091179093559091029092176004555046905282516001600160a01b031662000177576040516342bcdf7f60e11b815260040160405180910390fd5b82516001600160a01b0390811660a05260208401516001600160401b031660c05260408401511660e052620001ac8262000260565b505050620007f4565b336001600160a01b038216036200020f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005b81518110156200048b576000828281518110620002845762000284620007b6565b602090810291909101015180516060820151919250906001600160a01b0316620002c1576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160401b0381166000908152600f60205260409020600101546001600160a01b031662000328576040516001600160401b03821681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a15b6000604051806080016040528084602001511515815260200184604001516001600160a01b0316815260200184606001516001600160a01b031681526020016200039e8486606001517f8acd72527118c8324937b1a42e02cd246697c3b633f1742f3cae11de233722b36200048f60201b60201c565b90526001600160401b0383166000818152600f60209081526040918290208451815486840180516001600160a81b0319909216921515610100600160a81b03198116939093176101006001600160a01b039384160217845587860180516001860180546001600160a01b0319169185169190911790556060808a01805160029097019690965587519485529151831695840195909552935116938101939093525190820152919250907fdba8597411dc0624375cfff476f6173674609571f4d98d294dd3a47af07927849060800160405180910390a2505050806200048390620007cc565b905062000263565b5050565b60c05160408051602081018490526001600160401b0380871692820192909252911660608201526001600160a01b038316608082015260009060a0016040516020818303038152906040528051906020012090509392505050565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b0381118282101715620005255762000525620004ea565b60405290565b604051608081016001600160401b0381118282101715620005255762000525620004ea565b604051601f8201601f191681016001600160401b03811182821017156200057b576200057b620004ea565b604052919050565b80516001600160a01b03811681146200059b57600080fd5b919050565b80516001600160401b03811681146200059b57600080fd5b805180151581146200059b57600080fd5b80516001600160801b03811681146200059b57600080fd5b600060608284031215620005f457600080fd5b620005fe62000500565b90506200060b82620005b8565b81526200061b60208301620005c9565b60208201526200062e60408301620005c9565b604082015292915050565b600080600083850360e08112156200065057600080fd5b6060808212156200066057600080fd5b6200066a62000500565b9150620006778662000583565b8252602062000688818801620005a0565b8184015260406200069b81890162000583565b84820152878301519396506001600160401b0380851115620006bc57600080fd5b848901945089601f860112620006d157600080fd5b845181811115620006e657620006e6620004ea565b620006f6848260051b0162000550565b818152848101925060079190911b86018401908b8211156200071757600080fd5b958401955b8187101562000794576080878d031215620007375760008081fd5b620007416200052b565b6200074c88620005a0565b81526200075b868901620005b8565b868201526200076c85890162000583565b858201526200077d87890162000583565b81880152835260809690960195918401916200071c565b80985050505050505050620007ad8560808601620005e1565b90509250925092565b634e487b7160e01b600052603260045260246000fd5b600060018201620007ed57634e487b7160e01b600052601160045260246000fd5b5060010190565b60805160a05160c05160e0516161f56200086d6000396000818161023601528181611d120152611ecf01526000818161020601528181611ce90152613aed0152600081816101ca01528181611cbb015261218b015260008181610e8401528181610ed00152818161136201526113ae01526161f56000f3fe608060405234801561001057600080fd5b50600436106101985760003560e01c806381ff7048116100e3578063b1dc65a41161008c578063f077b59211610066578063f077b592146106f7578063f2fde38b1461070d578063f52121a51461072057600080fd5b8063b1dc65a4146105df578063c92b2832146105f2578063e9d68a8e1461060557600080fd5b8063873504d7116100bd578063873504d71461058e5780638da5cb5b146105a1578063afcb95d7146105bf57600080fd5b806381ff7048146104f757806385572ffb14610527578063856c82471461053557600080fd5b8063666cab8d116101455780637437ff9f1161011f5780637437ff9f146103e557806379ba5097146104dc5780637f63b711146104e457600080fd5b8063666cab8d146103aa578063704b6c02146103bf578063740f4150146103d257600080fd5b80631ef38174116101765780631ef38174146102f2578063546719cd14610307578063599f64311461036b57600080fd5b806306285c691461019d578063142a98fc14610289578063181f5a77146102a9575b600080fd5b610273604080516060810182526000808252602082018190529181019190915260405180606001604052807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000067ffffffffffffffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16815250905090565b604051610280919061471c565b60405180910390f35b61029c610297366004614787565b610733565b604051610280919061480e565b6102e56040518060400160405280601d81526020017f45564d3245564d4d756c74694f666652616d7020312e362e302d64657600000081525081565b604051610280919061488a565b610305610300366004614b1b565b6107ae565b005b61030f610c6d565b604051610280919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610280565b6103b2610d22565b6040516102809190614c39565b6103056103cd366004614c4c565b610d91565b6103056103e03660046150a8565b610e81565b6104cf6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a0810191909152506040805160c081018252600a5463ffffffff808216835273ffffffffffffffffffffffffffffffffffffffff64010000000090920482166020840152600b549182169383019390935261ffff7401000000000000000000000000000000000000000082041660608301527601000000000000000000000000000000000000000000008104831660808301527a010000000000000000000000000000000000000000000000000000900490911660a082015290565b6040516102809190615163565b610305611003565b6103056104f23660046151d2565b611100565b6007546005546040805163ffffffff80851682526401000000009094049093166020840152820152606001610280565b6103056101983660046152b6565b610575610543366004614c4c565b73ffffffffffffffffffffffffffffffffffffffff1660009081526010602052604090205467ffffffffffffffff1690565b60405167ffffffffffffffff9091168152602001610280565b61030561059c366004615382565b611114565b60005473ffffffffffffffffffffffffffffffffffffffff16610385565b604080516001815260006020820181905291810191909152606001610280565b6103056105ed36600461542b565b61130c565b610305610600366004615530565b61159d565b6106a7610613366004614787565b6040805160808101825260008082526020820181905291810182905260608101919091525067ffffffffffffffff166000908152600f60209081526040918290208251608081018452815460ff81161515825273ffffffffffffffffffffffffffffffffffffffff610100909104811693820193909352600182015490921692820192909252600290910154606082015290565b6040805182511515815260208084015173ffffffffffffffffffffffffffffffffffffffff9081169183019190915283830151169181019190915260609182015191810191909152608001610280565b6106ff61161f565b604051610280929190615580565b61030561071b366004614c4c565b611796565b61030561072e3660046155a5565b6117a7565b60006107416001600461562e565b600261074e608085615670565b67ffffffffffffffff166107629190615697565b601160006107716080876156ae565b67ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002054901c1660038111156107a8576107a86147a4565b92915050565b84518460ff16601f821115610824576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f746f6f206d616e79207472616e736d697474657273000000000000000000000060448201526064015b60405180910390fd5b8060000361088e576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f7369746976650000000000000000000000000000604482015260640161081b565b610896611a72565b61089f85611af5565b60095460005b8181101561092b5760086000600983815481106108c4576108c46156d5565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905561092481615704565b90506108a5565b50875160005b81811015610b275760008a828151811061094d5761094d6156d5565b602002602001015190506000600281111561096a5761096a6147a4565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260086020526040902054610100900460ff1660028111156109a9576109a96147a4565b14610a10576040517f89a6198900000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015260640161081b565b73ffffffffffffffffffffffffffffffffffffffff8116610a5d576040517fd6c62c9b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805180820190915260ff83168152602081016002905273ffffffffffffffffffffffffffffffffffffffff821660009081526008602090815260409091208251815460ff9091167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082168117835592840151919283917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001617610100836002811115610b0d57610b0d6147a4565b02179055509050505080610b2090615704565b9050610931565b508851610b3b9060099060208c0190614686565b506006805460ff838116610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000909216908b161717905560078054610bc1914691309190600090610b939063ffffffff1661573c565b91906101000a81548163ffffffff021916908363ffffffff160217905563ffffffff168d8d8d8d8d8d611d70565b6005600001819055506000600760049054906101000a900463ffffffff16905043600760046101000a81548163ffffffff021916908363ffffffff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600560000154600760009054906101000a900463ffffffff168e8e8e8e8e8e604051610c589998979695949392919061575f565b60405180910390a15050505050505050505050565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526040805160a0810182526003546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff1660208501527401000000000000000000000000000000000000000090920460ff161515938301939093526004548084166060840152049091166080820152610d1d90611e1b565b905090565b60606009805480602002602001604051908101604052809291908181526020018280548015610d8757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311610d5c575b5050505050905090565b60005473ffffffffffffffffffffffffffffffffffffffff163314801590610dd1575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610e08576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c9060200160405180910390a150565b467f000000000000000000000000000000000000000000000000000000000000000014610f0c576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f0000000000000000000000000000000000000000000000000000000000000000600482015267ffffffffffffffff4616602482015260440161081b565b81515181518114610f49576040517f83e3f56400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610ff3576000838281518110610f6857610f686156d5565b6020026020010151905080600014158015610fa157508451805183908110610f9257610f926156d5565b60200260200101516080015181105b15610fe2576040517f085e39cf000000000000000000000000000000000000000000000000000000008152600481018390526024810182905260440161081b565b50610fec81615704565b9050610f4c565b50610ffe8383611ecd565b505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611084576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161081b565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611108611a72565b611111816128d4565b50565b61111c611a72565b60005b82518110156112045761115983828151811061113d5761113d6156d5565b602002602001015160200151600c612b9690919063ffffffff16565b156111f4577fcbf3cbeaed4ac1d605ed30f4af06c35acaeff2379db7f6146c9cceee83d58782838281518110611191576111916156d5565b6020026020010151600001518483815181106111af576111af6156d5565b6020026020010151602001516040516111eb92919073ffffffffffffffffffffffffffffffffffffffff92831681529116602082015260400190565b60405180910390a15b6111fd81615704565b905061111f565b5060005b8151811015610ffe57611261828281518110611226576112266156d5565b602002602001015160200151838381518110611244576112446156d5565b602002602001015160000151600c612bbf9092919063ffffffff16565b156112fc577ffc23abf7ddbd3c02b1420dafa2355c56c1a06fbb8723862ac14d6bd74177361a828281518110611299576112996156d5565b6020026020010151600001518383815181106112b7576112b76156d5565b6020026020010151602001516040516112f392919073ffffffffffffffffffffffffffffffffffffffff92831681529116602082015260400190565b60405180910390a15b61130581615704565b9050611208565b6113168787612bea565b60055488359080821461135f576040517f93df584c000000000000000000000000000000000000000000000000000000008152600481018290526024810183905260440161081b565b467f0000000000000000000000000000000000000000000000000000000000000000146113e0576040517f0f01ce850000000000000000000000000000000000000000000000000000000081527f0000000000000000000000000000000000000000000000000000000000000000600482015246602482015260440161081b565b6040805183815260208c81013560081c63ffffffff16908201527fb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62910160405180910390a13360009081526008602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115611468576114686147a4565b6002811115611479576114796147a4565b9052509050600281602001516002811115611496576114966147a4565b1480156114dd57506009816000015160ff16815481106114b8576114b86156d5565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1633145b611513576040517fda0f08e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506000611521856020615697565b61152c886020615697565b6115388b6101446157f5565b61154291906157f5565b61154c91906157f5565b9050368114611590576040517f8e1192e10000000000000000000000000000000000000000000000000000000081526004810182905236602482015260440161081b565b5050505050505050505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633148015906115dd575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15611614576040517ff6cd562000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611111600382612c0d565b60608061162c600c612df2565b67ffffffffffffffff8111156116445761164461489d565b60405190808252806020026020018201604052801561166d578160200160208202803683370190505b50915061167a600c612df2565b67ffffffffffffffff8111156116925761169261489d565b6040519080825280602002602001820160405280156116bb578160200160208202803683370190505b50905060005b6116cb600c612df2565b811015611791576000806116e0600c84612dfd565b91509150808584815181106116f7576116f76156d5565b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505081848481518110611744576117446156d5565b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505050508061178a90615704565b90506116c1565b509091565b61179e611a72565b61111181612e1b565b3330146117e0576040517f371a732800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516000808252602082019092528161181d565b60408051808201909152600080825260208201528152602001906001900390816117f65790505b5061014084015151909150156118d6576101408301516040805160608101909152602085015173ffffffffffffffffffffffffffffffffffffffff1660808201526118d391908060a08101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152908252875167ffffffffffffffff1660208301528781015173ffffffffffffffffffffffffffffffffffffffff1691015261016086015185612f10565b90505b610120830151511580156118ec57506080830151155b806119105750604083015173ffffffffffffffffffffffffffffffffffffffff163b155b8061195d5750604083015161195b9073ffffffffffffffffffffffffffffffffffffffff167f85572ffb0000000000000000000000000000000000000000000000000000000061332a565b155b1561196757505050565b600a546000908190640100000000900473ffffffffffffffffffffffffffffffffffffffff16633cf9798361199c8786613346565b611388886080015189604001516040518563ffffffff1660e01b81526004016119c89493929190615859565b6000604051808303816000875af11580156119e7573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611a2d9190810190615970565b509150915081611a6b57806040517f0a8d6e8c00000000000000000000000000000000000000000000000000000000815260040161081b919061488a565b5050505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611af3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161081b565b565b600081806020019051810190611b0b91906159de565b602081015190915073ffffffffffffffffffffffffffffffffffffffff16611b5f576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051600a805460208085015173ffffffffffffffffffffffffffffffffffffffff908116640100000000027fffffffffffffffff00000000000000000000000000000000000000000000000090931663ffffffff9586161792909217909255604080850151600b805460608089015160808a015160a08b01518a167a010000000000000000000000000000000000000000000000000000027fffff00000000ffffffffffffffffffffffffffffffffffffffffffffffffffff91909a1676010000000000000000000000000000000000000000000002167fffff0000000000000000ffffffffffffffffffffffffffffffffffffffffffff61ffff90921674010000000000000000000000000000000000000000027fffffffffffffffffffff000000000000000000000000000000000000000000009094169588169590951792909217919091169290921795909517909455805193840181527f00000000000000000000000000000000000000000000000000000000000000008216845267ffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016928401929092527f00000000000000000000000000000000000000000000000000000000000000001682820152517f1406529d5e62f482fd1cd1a7e2ff03ff94c54a21619be336f55043d7c26fe83e91611d64918490615a8a565b60405180910390a15050565b6000808a8a8a8a8a8a8a8a8a604051602001611d9499989796959493929190615b3e565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150509998505050505050505050565b6040805160a081018252600080825260208201819052918101829052606081018290526080810191909152611ea982606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff1642611e8d919061562e565b85608001516fffffffffffffffffffffffffffffffff166133f6565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663397796f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015611f38573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f5c9190615bd3565b15611f93576040517fc148371500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8151516000819003611fd0576040517ebf199700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826020015151811461200e576040517f57e0e08300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008167ffffffffffffffff8111156120295761202961489d565b604051908082528060200260200182016040528015612052578160200160208202803683370190505b50905060005b8281101561214057600085600001518281518110612078576120786156d5565b602002602001015190506120ba81600f6000846000015167ffffffffffffffff1667ffffffffffffffff1681526020019081526020016000206002015461341e565b8383815181106120cc576120cc6156d5565b6020026020010181815250508061018001518383815181106120f0576120f06156d5565b60200260200101511461212f576040517f7185cf6b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5061213981615704565b9050612058565b50604080850151606086015191517f3204887500000000000000000000000000000000000000000000000000000000815260009273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016926332048875926121c192879291600401615c20565b602060405180830381865afa1580156121de573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122029190615c56565b90508060000361223e576040517fea75680100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8351151560005b848110156128cb57600087600001518281518110612265576122656156d5565b60200260200101519050600061227e8260600151610733565b90506000816003811115612294576122946147a4565b14806122b1575060038160038111156122af576122af6147a4565b145b6122f95760608201516040517f50a6e05200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161081b565b83156123b657600a5460009063ffffffff16612315874261562e565b119050808061233557506003826003811115612333576123336147a4565b145b61236b576040517f6358b0d000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b88848151811061237d5761237d6156d5565b60200260200101516000146123b05788848151811061239e5761239e6156d5565b60200260200101518360800181815250505b50612413565b60008160038111156123ca576123ca6147a4565b146124135760608201516040517f67d9ba0f00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161081b565b60208083015173ffffffffffffffffffffffffffffffffffffffff908116600090815260108352604080822054865167ffffffffffffffff9081168452600f90955291205492169161010090041681158015612484575073ffffffffffffffffffffffffffffffffffffffff811615155b156126075760208401516040517f856c824700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201529082169063856c824790602401602060405180830381865afa1580156124fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061251f9190615c6f565b60c085015190925067ffffffffffffffff1661253c836001615c8c565b67ffffffffffffffff16146125aa57836020015173ffffffffffffffffffffffffffffffffffffffff168460c0015167ffffffffffffffff167fe44a20935573a783dd0d5991c92d7b6a0eb3173566530364db3ec10e9a990b5d60405160405180910390a3505050506128bb565b60208481015173ffffffffffffffffffffffffffffffffffffffff16600090815260109091526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff84161790555b600083600381111561261b5761261b6147a4565b036126a85760c084015167ffffffffffffffff1661263a836001615c8c565b67ffffffffffffffff16146126a857836020015173ffffffffffffffffffffffffffffffffffffffff168460c0015167ffffffffffffffff167fd32ddb11d71e3d63411d37b09f9a8b28664f1cb1338bfd1413c173b0ebf4123760405160405180910390a3505050506128bb565b60008b6020015186815181106126c0576126c06156d5565b602002602001015190506126ec85606001518660000151876101400151518861012001515185516135a4565b6126fb8560600151600161374a565b60008061270887846137f4565b9150915061271a87606001518361374a565b88801561273857506003826003811115612736576127366147a4565b145b1561277157806040517fcf19edfd00000000000000000000000000000000000000000000000000000000815260040161081b919061488a565b6003826003811115612785576127856147a4565b141580156127a5575060028260038111156127a2576127a26147a4565b14155b156127e4578660600151826040517f9e26160300000000000000000000000000000000000000000000000000000000815260040161081b929190615cad565b60008660038111156127f8576127f86147a4565b036128655760208088015173ffffffffffffffffffffffffffffffffffffffff166000908152601090915260408120805467ffffffffffffffff169161283d83615ccb565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505b866101800151876060015167ffffffffffffffff167fd4f851956a5d67c3997d1c9205045fef79bae2947fdee7e9e2641abc7391ef6584846040516128ab929190615ce8565b60405180910390a3505050505050505b6128c481615704565b9050612245565b50505050505050565b60005b8151811015612b925760008282815181106128f4576128f46156d5565b6020908102919091010151805160608201519192509073ffffffffffffffffffffffffffffffffffffffff16612956576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600f602052604090206001015473ffffffffffffffffffffffffffffffffffffffff166129cb5760405167ffffffffffffffff821681527ff4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb99060200160405180910390a15b60006040518060800160405280846020015115158152602001846040015173ffffffffffffffffffffffffffffffffffffffff168152602001846060015173ffffffffffffffffffffffffffffffffffffffff168152602001612a538486606001517f8acd72527118c8324937b1a42e02cd246697c3b633f1742f3cae11de233722b3613ae7565b905267ffffffffffffffff83166000818152600f60209081526040918290208451815486840180517fffffffffffffffffffffff0000000000000000000000000000000000000000009092169215157fffffffffffffffffffffff0000000000000000000000000000000000000000ff81169390931761010073ffffffffffffffffffffffffffffffffffffffff9384160217845587860180516001860180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169185169190911790556060808a01805160029097019690965587519485529151831695840195909552935116938101939093525190820152919250907fdba8597411dc0624375cfff476f6173674609571f4d98d294dd3a47af07927849060800160405180910390a250505080612b8b90615704565b90506128d7565b5050565b6000612bb88373ffffffffffffffffffffffffffffffffffffffff8416613b76565b9392505050565b6000612be28473ffffffffffffffffffffffffffffffffffffffff851684613b82565b949350505050565b612b92612bf982840184615d08565b604080516000815260208101909152611ecd565b8154600090612c3690700100000000000000000000000000000000900463ffffffff164261562e565b90508015612cd85760018301548354612c7e916fffffffffffffffffffffffffffffffff808216928116918591700100000000000000000000000000000000909104166133f6565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612cfe916fffffffffffffffffffffffffffffffff9081169116613ba5565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990612de59084908151151581526020808301516fffffffffffffffffffffffffffffffff90811691830191909152604092830151169181019190915260600190565b60405180910390a1505050565b60006107a882613bbb565b6000808080612e0c8686613bc6565b909450925050505b9250929050565b3373ffffffffffffffffffffffffffffffffffffffff821603612e9a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161081b565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b836000805b8651811015613311576000858281518110612f3257612f326156d5565b6020026020010151806020019051810190612f4d9190615d3d565b90506000612f5e8260200151613bd5565b9050612fa073ffffffffffffffffffffffffffffffffffffffff82167faff2afbf0000000000000000000000000000000000000000000000000000000061332a565b612fee576040517fae9b4ce900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161081b565b60008061315a634059f55b60e01b6040518060e001604052808d6000015181526020018d6020015167ffffffffffffffff1681526020018d6040015173ffffffffffffffffffffffffffffffffffffffff1681526020018e8981518110613057576130576156d5565b602002602001015160200151815260200187600001518152602001876040015181526020018b898151811061308e5761308e6156d5565b60200260200101518152506040516024016130a99190615df2565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152600b54859063ffffffff7a010000000000000000000000000000000000000000000000000000909104166113886084613c30565b50915091508161319857806040517fe1cd550900000000000000000000000000000000000000000000000000000000815260040161081b919061488a565b80516040146131e2578051604080517f78ef80240000000000000000000000000000000000000000000000000000000081526004810191909152602481019190915260440161081b565b600080828060200190518101906131f99190615ecc565b9150915061320682613d56565b898881518110613218576132186156d5565b60200260200101516000019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505080898881518110613269576132696156d5565b602002602001015160200181815250506132aa89888151811061328e5761328e6156d5565b602002602001015160000151600c613ded90919063ffffffff16565b156132fa576132ed8988815181106132c4576132c46156d5565b6020908102919091010151600b5473ffffffffffffffffffffffffffffffffffffffff16613e0f565b6132f790896157f5565b97505b5050505050508061330a90615704565b9050612f15565b5080156133215761332181613f4a565b50949350505050565b600061333583613f57565b8015612bb85750612bb88383613fbb565b6040805160a08101825260008082526020820152606091810182905281810182905260808101919091526040518060a001604052808461018001518152602001846000015167ffffffffffffffff16815260200184602001516040516020016133cb919073ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b6040516020818303038152906040528152602001846101200151815260200183815250905092915050565b6000613415856134068486615697565b61341090876157f5565b613ba5565b95945050505050565b60008060001b8284602001518560400151866060015187608001518860a001518960c001518a60e001518b61010001516040516020016134c198979695949392919073ffffffffffffffffffffffffffffffffffffffff9889168152968816602088015267ffffffffffffffff95861660408801526060870194909452911515608086015290921660a0840152921660c082015260e08101919091526101000190565b60405160208183030381529060405280519060200120856101200151805190602001208661014001516040516020016134fa9190615ef0565b604051602081830303815290604052805190602001208761016001516040516020016135269190615f58565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282528051602091820120908301979097528101949094526060840192909252608083015260a082015260c081019190915260e00160405160208183030381529060405280519060200120905092915050565b67ffffffffffffffff84166000908152600f602052604090206001015473ffffffffffffffffffffffffffffffffffffffff16613619576040517f1279ec8a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015260240161081b565b600b5474010000000000000000000000000000000000000000900461ffff1683111561367d576040517f099d3f7200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015260240161081b565b8083146136c2576040517f8808f8e700000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8616600482015260240161081b565b600b54760100000000000000000000000000000000000000000000900463ffffffff16821115611a6b57600b546040517f8693378900000000000000000000000000000000000000000000000000000000815276010000000000000000000000000000000000000000000090910463ffffffff1660048201526024810183905260440161081b565b60006002613759608085615670565b67ffffffffffffffff1661376d9190615697565b9050600060118161377f6080876156ae565b67ffffffffffffffff1681526020810191909152604001600020549050816137a96001600461562e565b901b1916818360038111156137c0576137c06147a4565b901b1780601160006137d36080886156ae565b67ffffffffffffffff16815260208101919091526040016000205550505050565b6040517ff52121a5000000000000000000000000000000000000000000000000000000008152600090606090309063f52121a5906138389087908790600401615f6b565b600060405180830381600087803b15801561385257600080fd5b505af1925050508015613863575060015b613acc573d808015613891576040519150601f19603f3d011682016040523d82523d6000602084013e613896565b606091505b506138a0816160f5565b7fffffffff00000000000000000000000000000000000000000000000000000000167f0a8d6e8c00000000000000000000000000000000000000000000000000000000148061393857506138f3816160f5565b7fffffffff00000000000000000000000000000000000000000000000000000000167fe1cd550900000000000000000000000000000000000000000000000000000000145b8061398c5750613947816160f5565b7fffffffff00000000000000000000000000000000000000000000000000000000167f8d666f6000000000000000000000000000000000000000000000000000000000145b806139e0575061399b816160f5565b7fffffffff00000000000000000000000000000000000000000000000000000000167f78ef802400000000000000000000000000000000000000000000000000000000145b80613a3457506139ef816160f5565b7fffffffff00000000000000000000000000000000000000000000000000000000167f0c3b563c00000000000000000000000000000000000000000000000000000000145b80613a885750613a43816160f5565b7fffffffff00000000000000000000000000000000000000000000000000000000167fae9b4ce900000000000000000000000000000000000000000000000000000000145b15613a9857600392509050612e14565b806040517fcf19edfd00000000000000000000000000000000000000000000000000000000815260040161081b919061488a565b50506040805160208101909152600081526002909250929050565b600081847f000000000000000000000000000000000000000000000000000000000000000085604051602001613b57949392919093845267ffffffffffffffff92831660208501529116604083015273ffffffffffffffffffffffffffffffffffffffff16606082015260800190565b6040516020818303038152906040528051906020012090509392505050565b6000612bb8838361408a565b6000612be2848473ffffffffffffffffffffffffffffffffffffffff85166140a7565b6000818310613bb45781612bb8565b5090919050565b60006107a8826140c4565b6000808080612e0c86866140cf565b60008151602014613c1457816040517f8d666f6000000000000000000000000000000000000000000000000000000000815260040161081b919061488a565b6107a882806020019051810190613c2b9190615c56565b613d56565b6000606060008361ffff1667ffffffffffffffff811115613c5357613c5361489d565b6040519080825280601f01601f191660200182016040528015613c7d576020820181803683370190505b509150863b613cb0577f0c3b563c0000000000000000000000000000000000000000000000000000000060005260046000fd5b5a85811015613ce3577fafa32a2c0000000000000000000000000000000000000000000000000000000060005260046000fd5b8590036040810481038710613d1c577f37c3be290000000000000000000000000000000000000000000000000000000060005260046000fd5b505a6000808a5160208c0160008c8cf193505a900390503d84811115613d3f5750835b808352806000602085013e50955095509592505050565b600073ffffffffffffffffffffffffffffffffffffffff821180613d7a5750600a82105b15613de957604080516020810184905201604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f8d666f6000000000000000000000000000000000000000000000000000000000825261081b9160040161488a565b5090565b6000612bb88373ffffffffffffffffffffffffffffffffffffffff84166140fa565b81516040517fd02641a000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff918216600482015260009182919084169063d02641a0906024016040805180830381865afa158015613e82573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ea69190616145565b5190507bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116600003613f1c5783516040517f9a655f7b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161081b565b6020840151612be2907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff831690614106565b6111116003826000614143565b6000613f83827f01ffc9a700000000000000000000000000000000000000000000000000000000613fbb565b80156107a85750613fb4827fffffffff00000000000000000000000000000000000000000000000000000000613fbb565b1592915050565b604080517fffffffff000000000000000000000000000000000000000000000000000000008316602480830191909152825180830390910181526044909101909152602080820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f01ffc9a700000000000000000000000000000000000000000000000000000000178152825160009392849283928392918391908a617530fa92503d91506000519050828015614073575060208210155b801561407f5750600081115b979650505050505050565b60008181526002830160205260408120819055612bb883836144c6565b60008281526002840160205260408120829055612be284846144d2565b60006107a8826144de565b600080806140dd85856144e8565b600081815260029690960160205260409095205494959350505050565b6000612bb883836144f4565b6000670de0b6b3a7640000614139837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8616615697565b612bb891906161a5565b825474010000000000000000000000000000000000000000900460ff16158061416a575081155b1561417457505050565b825460018401546fffffffffffffffffffffffffffffffff808316929116906000906141ba90700100000000000000000000000000000000900463ffffffff164261562e565b9050801561427a57818311156141fc576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546142369083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166133f6565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b848210156143315773ffffffffffffffffffffffffffffffffffffffff84166142d9576040517ff94ebcd1000000000000000000000000000000000000000000000000000000008152600481018390526024810186905260440161081b565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff8516604482015260640161081b565b848310156144445760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290614375908261562e565b61437f878a61562e565b61438991906157f5565b61439391906161a5565b905073ffffffffffffffffffffffffffffffffffffffff86166143ec576040517f15279c08000000000000000000000000000000000000000000000000000000008152600481018290526024810186905260440161081b565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff8716604482015260640161081b565b61444e858461562e565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b6000612bb88383614513565b6000612bb8838361460d565b60006107a8825490565b6000612bb8838361465c565b6000612bb8838360008181526001830160205260408120541515612bb8565b600081815260018301602052604081205480156145fc57600061453760018361562e565b855490915060009061454b9060019061562e565b90508181146145b057600086600001828154811061456b5761456b6156d5565b906000526020600020015490508087600001848154811061458e5761458e6156d5565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806145c1576145c16161b9565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506107a8565b60009150506107a8565b5092915050565b6000818152600183016020526040812054614654575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556107a8565b5060006107a8565b6000826000018281548110614673576146736156d5565b9060005260206000200154905092915050565b828054828255906000526020600020908101928215614700579160200282015b8281111561470057825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906146a6565b50613de99291505b80821115613de95760008155600101614708565b606081016107a88284805173ffffffffffffffffffffffffffffffffffffffff908116835260208083015167ffffffffffffffff169084015260409182015116910152565b67ffffffffffffffff8116811461111157600080fd5b803561478281614761565b919050565b60006020828403121561479957600080fd5b8135612bb881614761565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6004811061480a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b602081016107a882846147d3565b60005b8381101561483757818101518382015260200161481f565b50506000910152565b6000815180845261485881602086016020860161481c565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612bb86020830184614840565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156148ef576148ef61489d565b60405290565b6040516101a0810167ffffffffffffffff811182821017156148ef576148ef61489d565b6040516080810167ffffffffffffffff811182821017156148ef576148ef61489d565b6040516060810167ffffffffffffffff811182821017156148ef576148ef61489d565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156149a6576149a661489d565b604052919050565b600067ffffffffffffffff8211156149c8576149c861489d565b5060051b60200190565b73ffffffffffffffffffffffffffffffffffffffff8116811461111157600080fd5b8035614782816149d2565b600082601f830112614a1057600080fd5b81356020614a25614a20836149ae565b61495f565b82815260059290921b84018101918181019086841115614a4457600080fd5b8286015b84811015614a68578035614a5b816149d2565b8352918301918301614a48565b509695505050505050565b803560ff8116811461478257600080fd5b600067ffffffffffffffff821115614a9e57614a9e61489d565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112614adb57600080fd5b8135614ae9614a2082614a84565b818152846020838601011115614afe57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215614b3457600080fd5b863567ffffffffffffffff80821115614b4c57600080fd5b614b588a838b016149ff565b97506020890135915080821115614b6e57600080fd5b614b7a8a838b016149ff565b9650614b8860408a01614a73565b95506060890135915080821115614b9e57600080fd5b614baa8a838b01614aca565b9450614bb860808a01614777565b935060a0890135915080821115614bce57600080fd5b50614bdb89828a01614aca565b9150509295509295509295565b600081518084526020808501945080840160005b83811015614c2e57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614bfc565b509495945050505050565b602081526000612bb86020830184614be8565b600060208284031215614c5e57600080fd5b8135612bb8816149d2565b801515811461111157600080fd5b803561478281614c69565b600082601f830112614c9357600080fd5b81356020614ca3614a20836149ae565b82815260069290921b84018101918181019086841115614cc257600080fd5b8286015b84811015614a685760408189031215614cdf5760008081fd5b614ce76148cc565b8135614cf2816149d2565b81528185013585820152835291830191604001614cc6565b600082601f830112614d1b57600080fd5b81356020614d2b614a20836149ae565b82815260059290921b84018101918181019086841115614d4a57600080fd5b8286015b84811015614a6857803567ffffffffffffffff811115614d6e5760008081fd5b614d7c8986838b0101614aca565b845250918301918301614d4e565b60006101a08284031215614d9d57600080fd5b614da56148f5565b9050614db082614777565b8152614dbe602083016149f4565b6020820152614dcf604083016149f4565b6040820152614de060608301614777565b606082015260808201356080820152614dfb60a08301614c77565b60a0820152614e0c60c08301614777565b60c0820152614e1d60e083016149f4565b60e082015261010082810135908201526101208083013567ffffffffffffffff80821115614e4a57600080fd5b614e5686838701614aca565b83850152610140925082850135915080821115614e7257600080fd5b614e7e86838701614c82565b83850152610160925082850135915080821115614e9a57600080fd5b50614ea785828601614d0a565b82840152505061018080830135818301525092915050565b600082601f830112614ed057600080fd5b81356020614ee0614a20836149ae565b82815260059290921b84018101918181019086841115614eff57600080fd5b8286015b84811015614a6857803567ffffffffffffffff811115614f235760008081fd5b614f318986838b0101614d0a565b845250918301918301614f03565b600082601f830112614f5057600080fd5b81356020614f60614a20836149ae565b82815260059290921b84018101918181019086841115614f7f57600080fd5b8286015b84811015614a685780358352918301918301614f83565b600060808284031215614fac57600080fd5b614fb4614919565b9050813567ffffffffffffffff80821115614fce57600080fd5b818401915084601f830112614fe257600080fd5b81356020614ff2614a20836149ae565b82815260059290921b8401810191818101908884111561501157600080fd5b8286015b848110156150495780358681111561502d5760008081fd5b61503b8b86838b0101614d8a565b845250918301918301615015565b508652508581013593508284111561506057600080fd5b61506c87858801614ebf565b9085015250604084013591508082111561508557600080fd5b5061509284828501614f3f565b6040830152506060820135606082015292915050565b600080604083850312156150bb57600080fd5b823567ffffffffffffffff808211156150d357600080fd5b6150df86838701614f9a565b93506020915081850135818111156150f657600080fd5b85019050601f8101861361510957600080fd5b8035615117614a20826149ae565b81815260059190911b8201830190838101908883111561513657600080fd5b928401925b828410156151545783358252928401929084019061513b565b80955050505050509250929050565b60c081016107a8828463ffffffff808251168352602082015173ffffffffffffffffffffffffffffffffffffffff8082166020860152806040850151166040860152505061ffff60608301511660608401528060808301511660808401528060a08301511660a0840152505050565b600060208083850312156151e557600080fd5b823567ffffffffffffffff8111156151fc57600080fd5b8301601f8101851361520d57600080fd5b803561521b614a20826149ae565b81815260079190911b8201830190838101908783111561523a57600080fd5b928401925b8284101561407f57608084890312156152585760008081fd5b615260614919565b843561526b81614761565b81528486013561527a81614c69565b8187015260408581013561528d816149d2565b908201526060858101356152a0816149d2565b908201528252608093909301929084019061523f565b6000602082840312156152c857600080fd5b813567ffffffffffffffff8111156152df57600080fd5b820160a08185031215612bb857600080fd5b600082601f83011261530257600080fd5b81356020615312614a20836149ae565b82815260069290921b8401810191818101908684111561533157600080fd5b8286015b84811015614a68576040818903121561534e5760008081fd5b6153566148cc565b8135615361816149d2565b815281850135615370816149d2565b81860152835291830191604001615335565b6000806040838503121561539557600080fd5b823567ffffffffffffffff808211156153ad57600080fd5b6153b9868387016152f1565b935060208501359150808211156153cf57600080fd5b506153dc858286016152f1565b9150509250929050565b60008083601f8401126153f857600080fd5b50813567ffffffffffffffff81111561541057600080fd5b6020830191508360208260051b8501011115612e1457600080fd5b60008060008060008060008060e0898b03121561544757600080fd5b606089018a81111561545857600080fd5b8998503567ffffffffffffffff8082111561547257600080fd5b818b0191508b601f83011261548657600080fd5b81358181111561549557600080fd5b8c60208285010111156154a757600080fd5b6020830199508098505060808b01359150808211156154c557600080fd5b6154d18c838d016153e6565b909750955060a08b01359150808211156154ea57600080fd5b506154f78b828c016153e6565b999c989b50969995989497949560c00135949350505050565b80356fffffffffffffffffffffffffffffffff8116811461478257600080fd5b60006060828403121561554257600080fd5b61554a61493c565b823561555581614c69565b815261556360208401615510565b602082015261557460408401615510565b60408201529392505050565b6040815260006155936040830185614be8565b82810360208401526134158185614be8565b600080604083850312156155b857600080fd5b823567ffffffffffffffff808211156155d057600080fd5b6155dc86838701614d8a565b935060208501359150808211156155f257600080fd5b506153dc85828601614d0a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156107a8576107a86155ff565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600067ffffffffffffffff8084168061568b5761568b615641565b92169190910692915050565b80820281158282048414176107a8576107a86155ff565b600067ffffffffffffffff808416806156c9576156c9615641565b92169190910492915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203615735576157356155ff565b5060010190565b600063ffffffff808316818103615755576157556155ff565b6001019392505050565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261578f8184018a614be8565b905082810360808401526157a38189614be8565b905060ff871660a084015282810360c08401526157c08187614840565b905067ffffffffffffffff851660e08401528281036101008401526157e58185614840565b9c9b505050505050505050505050565b808201808211156107a8576107a86155ff565b600081518084526020808501945080840160005b83811015614c2e578151805173ffffffffffffffffffffffffffffffffffffffff168852830151838801526040909601959082019060010161581c565b608081528451608082015267ffffffffffffffff60208601511660a08201526000604086015160a060c0840152615894610120840182614840565b905060608701517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80808584030160e08601526158d08383614840565b9250608089015191508085840301610100860152506158ef8282615808565b92505050615903602083018661ffff169052565b836040830152613415606083018473ffffffffffffffffffffffffffffffffffffffff169052565b600082601f83011261593c57600080fd5b815161594a614a2082614a84565b81815284602083860101111561595f57600080fd5b612be282602083016020870161481c565b60008060006060848603121561598557600080fd5b835161599081614c69565b602085015190935067ffffffffffffffff8111156159ad57600080fd5b6159b98682870161592b565b925050604084015190509250925092565b805163ffffffff8116811461478257600080fd5b600060c082840312156159f057600080fd5b60405160c0810181811067ffffffffffffffff82111715615a1357615a1361489d565b604052615a1f836159ca565b81526020830151615a2f816149d2565b60208201526040830151615a42816149d2565b6040820152606083015161ffff81168114615a5c57600080fd5b6060820152615a6d608084016159ca565b6080820152615a7e60a084016159ca565b60a08201529392505050565b6101208101615ad08285805173ffffffffffffffffffffffffffffffffffffffff908116835260208083015167ffffffffffffffff169084015260409182015116910152565b612bb8606083018463ffffffff808251168352602082015173ffffffffffffffffffffffffffffffffffffffff8082166020860152806040850151166040860152505061ffff60608301511660608401528060808301511660808401528060a08301511660a0840152505050565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b166040850152816060850152615b858285018b614be8565b91508382036080850152615b99828a614be8565b915060ff881660a085015283820360c0850152615bb68288614840565b90861660e085015283810361010085015290506157e58185614840565b600060208284031215615be557600080fd5b8151612bb881614c69565b600081518084526020808501945080840160005b83811015614c2e57815187529582019590820190600101615c04565b606081526000615c336060830186615bf0565b8281036020840152615c458186615bf0565b915050826040830152949350505050565b600060208284031215615c6857600080fd5b5051919050565b600060208284031215615c8157600080fd5b8151612bb881614761565b67ffffffffffffffff818116838216019080821115614606576146066155ff565b67ffffffffffffffff8316815260408101612bb860208301846147d3565b600067ffffffffffffffff808316818103615755576157556155ff565b615cf281846147d3565b604060208201526000612be26040830184614840565b600060208284031215615d1a57600080fd5b813567ffffffffffffffff811115615d3157600080fd5b612be284828501614f9a565b600060208284031215615d4f57600080fd5b815167ffffffffffffffff80821115615d6757600080fd5b9083019060608286031215615d7b57600080fd5b615d8361493c565b825182811115615d9257600080fd5b615d9e8782860161592b565b825250602083015182811115615db357600080fd5b615dbf8782860161592b565b602083015250604083015182811115615dd757600080fd5b615de38782860161592b565b60408301525095945050505050565b602081526000825160e06020840152615e0f610100840182614840565b905067ffffffffffffffff60208501511660408401526040840151615e4c606085018273ffffffffffffffffffffffffffffffffffffffff169052565b506060840151608084015260808401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808584030160a0860152615e918383614840565b925060a08601519150808584030160c0860152615eae8383614840565b925060c08601519150808584030160e0860152506134158282614840565b60008060408385031215615edf57600080fd5b505080516020909101519092909150565b602081526000612bb86020830184615808565b600081518084526020808501808196508360051b8101915082860160005b85811015615f4b578284038952615f39848351614840565b98850198935090840190600101615f21565b5091979650505050505050565b602081526000612bb86020830184615f03565b60408152615f8660408201845167ffffffffffffffff169052565b60006020840151615faf606084018273ffffffffffffffffffffffffffffffffffffffff169052565b50604084015173ffffffffffffffffffffffffffffffffffffffff8116608084015250606084015167ffffffffffffffff811660a084015250608084015160c083015260a084015180151560e08401525060c084015161010061601d8185018367ffffffffffffffff169052565b60e086015191506101206160488186018473ffffffffffffffffffffffffffffffffffffffff169052565b81870151925061014091508282860152808701519250506101a061016081818701526160786101e0870185614840565b93508288015192507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc06101808188870301818901526160b78686615808565b9550828a015194508188870301848901526160d28686615f03565b9550808a01516101c0890152505050505082810360208401526134158185615f03565b6000815160208301517fffffffff000000000000000000000000000000000000000000000000000000008082169350600483101561613d5780818460040360031b1b83161693505b505050919050565b60006040828403121561615757600080fd5b61615f6148cc565b82517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461618b57600080fd5b8152616199602084016159ca565b60208201529392505050565b6000826161b4576161b4615641565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000813000a",
}

var EVM2EVMMultiOffRampABI = EVM2EVMMultiOffRampMetaData.ABI

var EVM2EVMMultiOffRampBin = EVM2EVMMultiOffRampMetaData.Bin

func DeployEVM2EVMMultiOffRamp(auth *bind.TransactOpts, backend bind.ContractBackend, staticConfig EVM2EVMMultiOffRampStaticConfig, sourceChainConfigs []EVM2EVMMultiOffRampSourceChainConfigArgs, rateLimiterConfig RateLimiterConfig) (common.Address, *types.Transaction, *EVM2EVMMultiOffRamp, error) {
	parsed, err := EVM2EVMMultiOffRampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EVM2EVMMultiOffRampBin), backend, staticConfig, sourceChainConfigs, rateLimiterConfig)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EVM2EVMMultiOffRamp{address: address, abi: *parsed, EVM2EVMMultiOffRampCaller: EVM2EVMMultiOffRampCaller{contract: contract}, EVM2EVMMultiOffRampTransactor: EVM2EVMMultiOffRampTransactor{contract: contract}, EVM2EVMMultiOffRampFilterer: EVM2EVMMultiOffRampFilterer{contract: contract}}, nil
}

type EVM2EVMMultiOffRamp struct {
	address common.Address
	abi     abi.ABI
	EVM2EVMMultiOffRampCaller
	EVM2EVMMultiOffRampTransactor
	EVM2EVMMultiOffRampFilterer
}

type EVM2EVMMultiOffRampCaller struct {
	contract *bind.BoundContract
}

type EVM2EVMMultiOffRampTransactor struct {
	contract *bind.BoundContract
}

type EVM2EVMMultiOffRampFilterer struct {
	contract *bind.BoundContract
}

type EVM2EVMMultiOffRampSession struct {
	Contract     *EVM2EVMMultiOffRamp
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type EVM2EVMMultiOffRampCallerSession struct {
	Contract *EVM2EVMMultiOffRampCaller
	CallOpts bind.CallOpts
}

type EVM2EVMMultiOffRampTransactorSession struct {
	Contract     *EVM2EVMMultiOffRampTransactor
	TransactOpts bind.TransactOpts
}

type EVM2EVMMultiOffRampRaw struct {
	Contract *EVM2EVMMultiOffRamp
}

type EVM2EVMMultiOffRampCallerRaw struct {
	Contract *EVM2EVMMultiOffRampCaller
}

type EVM2EVMMultiOffRampTransactorRaw struct {
	Contract *EVM2EVMMultiOffRampTransactor
}

func NewEVM2EVMMultiOffRamp(address common.Address, backend bind.ContractBackend) (*EVM2EVMMultiOffRamp, error) {
	abi, err := abi.JSON(strings.NewReader(EVM2EVMMultiOffRampABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindEVM2EVMMultiOffRamp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRamp{address: address, abi: abi, EVM2EVMMultiOffRampCaller: EVM2EVMMultiOffRampCaller{contract: contract}, EVM2EVMMultiOffRampTransactor: EVM2EVMMultiOffRampTransactor{contract: contract}, EVM2EVMMultiOffRampFilterer: EVM2EVMMultiOffRampFilterer{contract: contract}}, nil
}

func NewEVM2EVMMultiOffRampCaller(address common.Address, caller bind.ContractCaller) (*EVM2EVMMultiOffRampCaller, error) {
	contract, err := bindEVM2EVMMultiOffRamp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampCaller{contract: contract}, nil
}

func NewEVM2EVMMultiOffRampTransactor(address common.Address, transactor bind.ContractTransactor) (*EVM2EVMMultiOffRampTransactor, error) {
	contract, err := bindEVM2EVMMultiOffRamp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampTransactor{contract: contract}, nil
}

func NewEVM2EVMMultiOffRampFilterer(address common.Address, filterer bind.ContractFilterer) (*EVM2EVMMultiOffRampFilterer, error) {
	contract, err := bindEVM2EVMMultiOffRamp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampFilterer{contract: contract}, nil
}

func bindEVM2EVMMultiOffRamp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EVM2EVMMultiOffRampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMMultiOffRamp.Contract.EVM2EVMMultiOffRampCaller.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.EVM2EVMMultiOffRampTransactor.contract.Transfer(opts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.EVM2EVMMultiOffRampTransactor.contract.Transact(opts, method, params...)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EVM2EVMMultiOffRamp.Contract.contract.Call(opts, result, method, params...)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.contract.Transfer(opts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.contract.Transact(opts, method, params...)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "ccipReceive", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _EVM2EVMMultiOffRamp.Contract.CcipReceive(&_EVM2EVMMultiOffRamp.CallOpts, arg0)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) CcipReceive(arg0 ClientAny2EVMMessage) error {
	return _EVM2EVMMultiOffRamp.Contract.CcipReceive(&_EVM2EVMMultiOffRamp.CallOpts, arg0)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "currentRateLimiterState")

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMMultiOffRamp.Contract.CurrentRateLimiterState(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) CurrentRateLimiterState() (RateLimiterTokenBucket, error) {
	return _EVM2EVMMultiOffRamp.Contract.CurrentRateLimiterState(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetAllRateLimitTokens(opts *bind.CallOpts) (GetAllRateLimitTokens,

	error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getAllRateLimitTokens")

	outstruct := new(GetAllRateLimitTokens)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SourceTokens = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.DestTokens = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetAllRateLimitTokens() (GetAllRateLimitTokens,

	error) {
	return _EVM2EVMMultiOffRamp.Contract.GetAllRateLimitTokens(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetAllRateLimitTokens() (GetAllRateLimitTokens,

	error) {
	return _EVM2EVMMultiOffRamp.Contract.GetAllRateLimitTokens(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMMultiOffRampDynamicConfig, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getDynamicConfig")

	if err != nil {
		return *new(EVM2EVMMultiOffRampDynamicConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOffRampDynamicConfig)).(*EVM2EVMMultiOffRampDynamicConfig)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetDynamicConfig() (EVM2EVMMultiOffRampDynamicConfig, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetDynamicConfig(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetDynamicConfig() (EVM2EVMMultiOffRampDynamicConfig, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetDynamicConfig(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetExecutionState(opts *bind.CallOpts, sequenceNumber uint64) (uint8, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getExecutionState", sequenceNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetExecutionState(sequenceNumber uint64) (uint8, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetExecutionState(&_EVM2EVMMultiOffRamp.CallOpts, sequenceNumber)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetExecutionState(sequenceNumber uint64) (uint8, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetExecutionState(&_EVM2EVMMultiOffRamp.CallOpts, sequenceNumber)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetSenderNonce(opts *bind.CallOpts, sender common.Address) (uint64, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getSenderNonce", sender)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetSenderNonce(sender common.Address) (uint64, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetSenderNonce(&_EVM2EVMMultiOffRamp.CallOpts, sender)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetSenderNonce(sender common.Address) (uint64, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetSenderNonce(&_EVM2EVMMultiOffRamp.CallOpts, sender)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getSourceChainConfig", sourceChainSelector)

	if err != nil {
		return *new(EVM2EVMMultiOffRampSourceChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOffRampSourceChainConfig)).(*EVM2EVMMultiOffRampSourceChainConfig)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetSourceChainConfig(sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetSourceChainConfig(&_EVM2EVMMultiOffRamp.CallOpts, sourceChainSelector)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetSourceChainConfig(sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetSourceChainConfig(&_EVM2EVMMultiOffRamp.CallOpts, sourceChainSelector)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetStaticConfig(opts *bind.CallOpts) (EVM2EVMMultiOffRampStaticConfig, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getStaticConfig")

	if err != nil {
		return *new(EVM2EVMMultiOffRampStaticConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(EVM2EVMMultiOffRampStaticConfig)).(*EVM2EVMMultiOffRampStaticConfig)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetStaticConfig() (EVM2EVMMultiOffRampStaticConfig, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetStaticConfig(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetStaticConfig() (EVM2EVMMultiOffRampStaticConfig, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetStaticConfig(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getTokenLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetTokenLimitAdmin() (common.Address, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetTokenLimitAdmin(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) GetTransmitters() ([]common.Address, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetTransmitters(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) GetTransmitters() ([]common.Address, error) {
	return _EVM2EVMMultiOffRamp.Contract.GetTransmitters(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _EVM2EVMMultiOffRamp.Contract.LatestConfigDetails(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _EVM2EVMMultiOffRamp.Contract.LatestConfigDetails(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _EVM2EVMMultiOffRamp.Contract.LatestConfigDigestAndEpoch(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _EVM2EVMMultiOffRamp.Contract.LatestConfigDigestAndEpoch(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) Owner() (common.Address, error) {
	return _EVM2EVMMultiOffRamp.Contract.Owner(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) Owner() (common.Address, error) {
	return _EVM2EVMMultiOffRamp.Contract.Owner(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EVM2EVMMultiOffRamp.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) TypeAndVersion() (string, error) {
	return _EVM2EVMMultiOffRamp.Contract.TypeAndVersion(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampCallerSession) TypeAndVersion() (string, error) {
	return _EVM2EVMMultiOffRamp.Contract.TypeAndVersion(&_EVM2EVMMultiOffRamp.CallOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "acceptOwnership")
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.AcceptOwnership(&_EVM2EVMMultiOffRamp.TransactOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.AcceptOwnership(&_EVM2EVMMultiOffRamp.TransactOpts)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []EVM2EVMMultiOffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "applySourceChainConfigUpdates", sourceChainConfigUpdates)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []EVM2EVMMultiOffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.ApplySourceChainConfigUpdates(&_EVM2EVMMultiOffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) ApplySourceChainConfigUpdates(sourceChainConfigUpdates []EVM2EVMMultiOffRampSourceChainConfigArgs) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.ApplySourceChainConfigUpdates(&_EVM2EVMMultiOffRamp.TransactOpts, sourceChainConfigUpdates)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) ExecuteSingleMessage(opts *bind.TransactOpts, message InternalEVM2EVMMessage, offchainTokenData [][]byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "executeSingleMessage", message, offchainTokenData)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) ExecuteSingleMessage(message InternalEVM2EVMMessage, offchainTokenData [][]byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.ExecuteSingleMessage(&_EVM2EVMMultiOffRamp.TransactOpts, message, offchainTokenData)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) ExecuteSingleMessage(message InternalEVM2EVMMessage, offchainTokenData [][]byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.ExecuteSingleMessage(&_EVM2EVMMultiOffRamp.TransactOpts, message, offchainTokenData)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) ManuallyExecute(opts *bind.TransactOpts, report InternalExecutionReport, gasLimitOverrides []*big.Int) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "manuallyExecute", report, gasLimitOverrides)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) ManuallyExecute(report InternalExecutionReport, gasLimitOverrides []*big.Int) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.ManuallyExecute(&_EVM2EVMMultiOffRamp.TransactOpts, report, gasLimitOverrides)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) ManuallyExecute(report InternalExecutionReport, gasLimitOverrides []*big.Int) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.ManuallyExecute(&_EVM2EVMMultiOffRamp.TransactOpts, report, gasLimitOverrides)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "setAdmin", newAdmin)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.SetAdmin(&_EVM2EVMMultiOffRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.SetAdmin(&_EVM2EVMMultiOffRamp.TransactOpts, newAdmin)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) SetOCR2Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "setOCR2Config", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) SetOCR2Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.SetOCR2Config(&_EVM2EVMMultiOffRamp.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) SetOCR2Config(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.SetOCR2Config(&_EVM2EVMMultiOffRamp.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "setRateLimiterConfig", config)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.SetRateLimiterConfig(&_EVM2EVMMultiOffRamp.TransactOpts, config)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) SetRateLimiterConfig(config RateLimiterConfig) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.SetRateLimiterConfig(&_EVM2EVMMultiOffRamp.TransactOpts, config)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "transferOwnership", to)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.TransferOwnership(&_EVM2EVMMultiOffRamp.TransactOpts, to)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.TransferOwnership(&_EVM2EVMMultiOffRamp.TransactOpts, to)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "transmit", reportContext, report, rs, ss, arg4)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.Transmit(&_EVM2EVMMultiOffRamp.TransactOpts, reportContext, report, rs, ss, arg4)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.Transmit(&_EVM2EVMMultiOffRamp.TransactOpts, reportContext, report, rs, ss, arg4)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactor) UpdateRateLimitTokens(opts *bind.TransactOpts, removes []EVM2EVMMultiOffRampRateLimitToken, adds []EVM2EVMMultiOffRampRateLimitToken) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.contract.Transact(opts, "updateRateLimitTokens", removes, adds)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampSession) UpdateRateLimitTokens(removes []EVM2EVMMultiOffRampRateLimitToken, adds []EVM2EVMMultiOffRampRateLimitToken) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.UpdateRateLimitTokens(&_EVM2EVMMultiOffRamp.TransactOpts, removes, adds)
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampTransactorSession) UpdateRateLimitTokens(removes []EVM2EVMMultiOffRampRateLimitToken, adds []EVM2EVMMultiOffRampRateLimitToken) (*types.Transaction, error) {
	return _EVM2EVMMultiOffRamp.Contract.UpdateRateLimitTokens(&_EVM2EVMMultiOffRamp.TransactOpts, removes, adds)
}

type EVM2EVMMultiOffRampAdminSetIterator struct {
	Event *EVM2EVMMultiOffRampAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampAdminSet)
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
		it.Event = new(EVM2EVMMultiOffRampAdminSet)
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

func (it *EVM2EVMMultiOffRampAdminSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampAdminSet struct {
	NewAdmin common.Address
	Raw      types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampAdminSetIterator, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampAdminSetIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "AdminSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampAdminSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "AdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampAdminSet)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseAdminSet(log types.Log) (*EVM2EVMMultiOffRampAdminSet, error) {
	event := new(EVM2EVMMultiOffRampAdminSet)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "AdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampConfigSetIterator struct {
	Event *EVM2EVMMultiOffRampConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampConfigSet)
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
		it.Event = new(EVM2EVMMultiOffRampConfigSet)
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

func (it *EVM2EVMMultiOffRampConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampConfigSet struct {
	StaticConfig  EVM2EVMMultiOffRampStaticConfig
	DynamicConfig EVM2EVMMultiOffRampDynamicConfig
	Raw           types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampConfigSetIterator, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampConfigSetIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampConfigSet) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampConfigSet)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseConfigSet(log types.Log) (*EVM2EVMMultiOffRampConfigSet, error) {
	event := new(EVM2EVMMultiOffRampConfigSet)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampConfigSet0Iterator struct {
	Event *EVM2EVMMultiOffRampConfigSet0

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampConfigSet0Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampConfigSet0)
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
		it.Event = new(EVM2EVMMultiOffRampConfigSet0)
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

func (it *EVM2EVMMultiOffRampConfigSet0Iterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampConfigSet0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampConfigSet0 struct {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterConfigSet0(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampConfigSet0Iterator, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "ConfigSet0")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampConfigSet0Iterator{contract: _EVM2EVMMultiOffRamp.contract, event: "ConfigSet0", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchConfigSet0(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampConfigSet0) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "ConfigSet0")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampConfigSet0)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "ConfigSet0", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseConfigSet0(log types.Log) (*EVM2EVMMultiOffRampConfigSet0, error) {
	event := new(EVM2EVMMultiOffRampConfigSet0)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "ConfigSet0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampExecutionStateChangedIterator struct {
	Event *EVM2EVMMultiOffRampExecutionStateChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampExecutionStateChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampExecutionStateChanged)
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
		it.Event = new(EVM2EVMMultiOffRampExecutionStateChanged)
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

func (it *EVM2EVMMultiOffRampExecutionStateChangedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampExecutionStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampExecutionStateChanged struct {
	SequenceNumber uint64
	MessageId      [32]byte
	State          uint8
	ReturnData     []byte
	Raw            types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterExecutionStateChanged(opts *bind.FilterOpts, sequenceNumber []uint64, messageId [][32]byte) (*EVM2EVMMultiOffRampExecutionStateChangedIterator, error) {

	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "ExecutionStateChanged", sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampExecutionStateChangedIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "ExecutionStateChanged", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampExecutionStateChanged, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error) {

	var sequenceNumberRule []interface{}
	for _, sequenceNumberItem := range sequenceNumber {
		sequenceNumberRule = append(sequenceNumberRule, sequenceNumberItem)
	}
	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "ExecutionStateChanged", sequenceNumberRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampExecutionStateChanged)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseExecutionStateChanged(log types.Log) (*EVM2EVMMultiOffRampExecutionStateChanged, error) {
	event := new(EVM2EVMMultiOffRampExecutionStateChanged)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "ExecutionStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampOwnershipTransferRequestedIterator struct {
	Event *EVM2EVMMultiOffRampOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampOwnershipTransferRequested)
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
		it.Event = new(EVM2EVMMultiOffRampOwnershipTransferRequested)
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

func (it *EVM2EVMMultiOffRampOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOffRampOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampOwnershipTransferRequestedIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampOwnershipTransferRequested)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMMultiOffRampOwnershipTransferRequested, error) {
	event := new(EVM2EVMMultiOffRampOwnershipTransferRequested)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampOwnershipTransferredIterator struct {
	Event *EVM2EVMMultiOffRampOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampOwnershipTransferred)
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
		it.Event = new(EVM2EVMMultiOffRampOwnershipTransferred)
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

func (it *EVM2EVMMultiOffRampOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOffRampOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampOwnershipTransferredIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampOwnershipTransferred)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseOwnershipTransferred(log types.Log) (*EVM2EVMMultiOffRampOwnershipTransferred, error) {
	event := new(EVM2EVMMultiOffRampOwnershipTransferred)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampSkippedIncorrectNonceIterator struct {
	Event *EVM2EVMMultiOffRampSkippedIncorrectNonce

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampSkippedIncorrectNonceIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampSkippedIncorrectNonce)
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
		it.Event = new(EVM2EVMMultiOffRampSkippedIncorrectNonce)
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

func (it *EVM2EVMMultiOffRampSkippedIncorrectNonceIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampSkippedIncorrectNonceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampSkippedIncorrectNonce struct {
	Nonce  uint64
	Sender common.Address
	Raw    types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterSkippedIncorrectNonce(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMMultiOffRampSkippedIncorrectNonceIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "SkippedIncorrectNonce", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampSkippedIncorrectNonceIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "SkippedIncorrectNonce", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchSkippedIncorrectNonce(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSkippedIncorrectNonce, nonce []uint64, sender []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "SkippedIncorrectNonce", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampSkippedIncorrectNonce)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SkippedIncorrectNonce", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseSkippedIncorrectNonce(log types.Log) (*EVM2EVMMultiOffRampSkippedIncorrectNonce, error) {
	event := new(EVM2EVMMultiOffRampSkippedIncorrectNonce)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SkippedIncorrectNonce", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflightIterator struct {
	Event *EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflightIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight)
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
		it.Event = new(EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight)
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

func (it *EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflightIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflightIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight struct {
	Nonce  uint64
	Sender common.Address
	Raw    types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterSkippedSenderWithPreviousRampMessageInflight(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflightIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "SkippedSenderWithPreviousRampMessageInflight", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflightIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "SkippedSenderWithPreviousRampMessageInflight", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchSkippedSenderWithPreviousRampMessageInflight(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight, nonce []uint64, sender []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "SkippedSenderWithPreviousRampMessageInflight", nonceRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SkippedSenderWithPreviousRampMessageInflight", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseSkippedSenderWithPreviousRampMessageInflight(log types.Log) (*EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight, error) {
	event := new(EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SkippedSenderWithPreviousRampMessageInflight", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampSourceChainConfigSetIterator struct {
	Event *EVM2EVMMultiOffRampSourceChainConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampSourceChainConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampSourceChainConfigSet)
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
		it.Event = new(EVM2EVMMultiOffRampSourceChainConfigSet)
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

func (it *EVM2EVMMultiOffRampSourceChainConfigSetIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampSourceChainConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceConfig        EVM2EVMMultiOffRampSourceChainConfig
	Raw                 types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*EVM2EVMMultiOffRampSourceChainConfigSetIterator, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampSourceChainConfigSetIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "SourceChainConfigSet", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error) {

	var sourceChainSelectorRule []interface{}
	for _, sourceChainSelectorItem := range sourceChainSelector {
		sourceChainSelectorRule = append(sourceChainSelectorRule, sourceChainSelectorItem)
	}

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "SourceChainConfigSet", sourceChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampSourceChainConfigSet)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseSourceChainConfigSet(log types.Log) (*EVM2EVMMultiOffRampSourceChainConfigSet, error) {
	event := new(EVM2EVMMultiOffRampSourceChainConfigSet)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SourceChainConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampSourceChainSelectorAddedIterator struct {
	Event *EVM2EVMMultiOffRampSourceChainSelectorAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampSourceChainSelectorAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampSourceChainSelectorAdded)
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
		it.Event = new(EVM2EVMMultiOffRampSourceChainSelectorAdded)
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

func (it *EVM2EVMMultiOffRampSourceChainSelectorAddedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampSourceChainSelectorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampSourceChainSelectorAdded struct {
	SourceChainSelector uint64
	Raw                 types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampSourceChainSelectorAddedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampSourceChainSelectorAddedIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "SourceChainSelectorAdded", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSourceChainSelectorAdded) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "SourceChainSelectorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampSourceChainSelectorAdded)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseSourceChainSelectorAdded(log types.Log) (*EVM2EVMMultiOffRampSourceChainSelectorAdded, error) {
	event := new(EVM2EVMMultiOffRampSourceChainSelectorAdded)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "SourceChainSelectorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampTokenAggregateRateLimitAddedIterator struct {
	Event *EVM2EVMMultiOffRampTokenAggregateRateLimitAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampTokenAggregateRateLimitAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampTokenAggregateRateLimitAdded)
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
		it.Event = new(EVM2EVMMultiOffRampTokenAggregateRateLimitAdded)
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

func (it *EVM2EVMMultiOffRampTokenAggregateRateLimitAddedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampTokenAggregateRateLimitAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampTokenAggregateRateLimitAdded struct {
	SourceToken common.Address
	DestToken   common.Address
	Raw         types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterTokenAggregateRateLimitAdded(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampTokenAggregateRateLimitAddedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "TokenAggregateRateLimitAdded")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampTokenAggregateRateLimitAddedIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "TokenAggregateRateLimitAdded", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchTokenAggregateRateLimitAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampTokenAggregateRateLimitAdded) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "TokenAggregateRateLimitAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampTokenAggregateRateLimitAdded)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitAdded", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseTokenAggregateRateLimitAdded(log types.Log) (*EVM2EVMMultiOffRampTokenAggregateRateLimitAdded, error) {
	event := new(EVM2EVMMultiOffRampTokenAggregateRateLimitAdded)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampTokenAggregateRateLimitRemovedIterator struct {
	Event *EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampTokenAggregateRateLimitRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved)
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
		it.Event = new(EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved)
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

func (it *EVM2EVMMultiOffRampTokenAggregateRateLimitRemovedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampTokenAggregateRateLimitRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved struct {
	SourceToken common.Address
	DestToken   common.Address
	Raw         types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterTokenAggregateRateLimitRemoved(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampTokenAggregateRateLimitRemovedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "TokenAggregateRateLimitRemoved")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampTokenAggregateRateLimitRemovedIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "TokenAggregateRateLimitRemoved", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchTokenAggregateRateLimitRemoved(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "TokenAggregateRateLimitRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitRemoved", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseTokenAggregateRateLimitRemoved(log types.Log) (*EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved, error) {
	event := new(EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "TokenAggregateRateLimitRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type EVM2EVMMultiOffRampTransmittedIterator struct {
	Event *EVM2EVMMultiOffRampTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *EVM2EVMMultiOffRampTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EVM2EVMMultiOffRampTransmitted)
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
		it.Event = new(EVM2EVMMultiOffRampTransmitted)
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

func (it *EVM2EVMMultiOffRampTransmittedIterator) Error() error {
	return it.fail
}

func (it *EVM2EVMMultiOffRampTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type EVM2EVMMultiOffRampTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) FilterTransmitted(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampTransmittedIterator, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &EVM2EVMMultiOffRampTransmittedIterator{contract: _EVM2EVMMultiOffRamp.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampTransmitted) (event.Subscription, error) {

	logs, sub, err := _EVM2EVMMultiOffRamp.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(EVM2EVMMultiOffRampTransmitted)
				if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRampFilterer) ParseTransmitted(log types.Log) (*EVM2EVMMultiOffRampTransmitted, error) {
	event := new(EVM2EVMMultiOffRampTransmitted)
	if err := _EVM2EVMMultiOffRamp.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRamp) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _EVM2EVMMultiOffRamp.abi.Events["AdminSet"].ID:
		return _EVM2EVMMultiOffRamp.ParseAdminSet(log)
	case _EVM2EVMMultiOffRamp.abi.Events["ConfigSet"].ID:
		return _EVM2EVMMultiOffRamp.ParseConfigSet(log)
	case _EVM2EVMMultiOffRamp.abi.Events["ConfigSet0"].ID:
		return _EVM2EVMMultiOffRamp.ParseConfigSet0(log)
	case _EVM2EVMMultiOffRamp.abi.Events["ExecutionStateChanged"].ID:
		return _EVM2EVMMultiOffRamp.ParseExecutionStateChanged(log)
	case _EVM2EVMMultiOffRamp.abi.Events["OwnershipTransferRequested"].ID:
		return _EVM2EVMMultiOffRamp.ParseOwnershipTransferRequested(log)
	case _EVM2EVMMultiOffRamp.abi.Events["OwnershipTransferred"].ID:
		return _EVM2EVMMultiOffRamp.ParseOwnershipTransferred(log)
	case _EVM2EVMMultiOffRamp.abi.Events["SkippedIncorrectNonce"].ID:
		return _EVM2EVMMultiOffRamp.ParseSkippedIncorrectNonce(log)
	case _EVM2EVMMultiOffRamp.abi.Events["SkippedSenderWithPreviousRampMessageInflight"].ID:
		return _EVM2EVMMultiOffRamp.ParseSkippedSenderWithPreviousRampMessageInflight(log)
	case _EVM2EVMMultiOffRamp.abi.Events["SourceChainConfigSet"].ID:
		return _EVM2EVMMultiOffRamp.ParseSourceChainConfigSet(log)
	case _EVM2EVMMultiOffRamp.abi.Events["SourceChainSelectorAdded"].ID:
		return _EVM2EVMMultiOffRamp.ParseSourceChainSelectorAdded(log)
	case _EVM2EVMMultiOffRamp.abi.Events["TokenAggregateRateLimitAdded"].ID:
		return _EVM2EVMMultiOffRamp.ParseTokenAggregateRateLimitAdded(log)
	case _EVM2EVMMultiOffRamp.abi.Events["TokenAggregateRateLimitRemoved"].ID:
		return _EVM2EVMMultiOffRamp.ParseTokenAggregateRateLimitRemoved(log)
	case _EVM2EVMMultiOffRamp.abi.Events["Transmitted"].ID:
		return _EVM2EVMMultiOffRamp.ParseTransmitted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (EVM2EVMMultiOffRampAdminSet) Topic() common.Hash {
	return common.HexToHash("0x8fe72c3e0020beb3234e76ae6676fa576fbfcae600af1c4fea44784cf0db329c")
}

func (EVM2EVMMultiOffRampConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1406529d5e62f482fd1cd1a7e2ff03ff94c54a21619be336f55043d7c26fe83e")
}

func (EVM2EVMMultiOffRampConfigSet0) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (EVM2EVMMultiOffRampExecutionStateChanged) Topic() common.Hash {
	return common.HexToHash("0xd4f851956a5d67c3997d1c9205045fef79bae2947fdee7e9e2641abc7391ef65")
}

func (EVM2EVMMultiOffRampOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (EVM2EVMMultiOffRampOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (EVM2EVMMultiOffRampSkippedIncorrectNonce) Topic() common.Hash {
	return common.HexToHash("0xd32ddb11d71e3d63411d37b09f9a8b28664f1cb1338bfd1413c173b0ebf41237")
}

func (EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight) Topic() common.Hash {
	return common.HexToHash("0xe44a20935573a783dd0d5991c92d7b6a0eb3173566530364db3ec10e9a990b5d")
}

func (EVM2EVMMultiOffRampSourceChainConfigSet) Topic() common.Hash {
	return common.HexToHash("0xdba8597411dc0624375cfff476f6173674609571f4d98d294dd3a47af0792784")
}

func (EVM2EVMMultiOffRampSourceChainSelectorAdded) Topic() common.Hash {
	return common.HexToHash("0xf4c1390c70e5c0f491ae1ccbc06f9117cbbadf2767b247b3bc203280f24c0fb9")
}

func (EVM2EVMMultiOffRampTokenAggregateRateLimitAdded) Topic() common.Hash {
	return common.HexToHash("0xfc23abf7ddbd3c02b1420dafa2355c56c1a06fbb8723862ac14d6bd74177361a")
}

func (EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved) Topic() common.Hash {
	return common.HexToHash("0xcbf3cbeaed4ac1d605ed30f4af06c35acaeff2379db7f6146c9cceee83d58782")
}

func (EVM2EVMMultiOffRampTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (_EVM2EVMMultiOffRamp *EVM2EVMMultiOffRamp) Address() common.Address {
	return _EVM2EVMMultiOffRamp.address
}

type EVM2EVMMultiOffRampInterface interface {
	CcipReceive(opts *bind.CallOpts, arg0 ClientAny2EVMMessage) error

	CurrentRateLimiterState(opts *bind.CallOpts) (RateLimiterTokenBucket, error)

	GetAllRateLimitTokens(opts *bind.CallOpts) (GetAllRateLimitTokens,

		error)

	GetDynamicConfig(opts *bind.CallOpts) (EVM2EVMMultiOffRampDynamicConfig, error)

	GetExecutionState(opts *bind.CallOpts, sequenceNumber uint64) (uint8, error)

	GetSenderNonce(opts *bind.CallOpts, sender common.Address) (uint64, error)

	GetSourceChainConfig(opts *bind.CallOpts, sourceChainSelector uint64) (EVM2EVMMultiOffRampSourceChainConfig, error)

	GetStaticConfig(opts *bind.CallOpts) (EVM2EVMMultiOffRampStaticConfig, error)

	GetTokenLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetTransmitters(opts *bind.CallOpts) ([]common.Address, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

		error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplySourceChainConfigUpdates(opts *bind.TransactOpts, sourceChainConfigUpdates []EVM2EVMMultiOffRampSourceChainConfigArgs) (*types.Transaction, error)

	ExecuteSingleMessage(opts *bind.TransactOpts, message InternalEVM2EVMMessage, offchainTokenData [][]byte) (*types.Transaction, error)

	ManuallyExecute(opts *bind.TransactOpts, report InternalExecutionReport, gasLimitOverrides []*big.Int) (*types.Transaction, error)

	SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error)

	SetOCR2Config(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetRateLimiterConfig(opts *bind.TransactOpts, config RateLimiterConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, arg4 [32]byte) (*types.Transaction, error)

	UpdateRateLimitTokens(opts *bind.TransactOpts, removes []EVM2EVMMultiOffRampRateLimitToken, adds []EVM2EVMMultiOffRampRateLimitToken) (*types.Transaction, error)

	FilterAdminSet(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampAdminSetIterator, error)

	WatchAdminSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampAdminSet) (event.Subscription, error)

	ParseAdminSet(log types.Log) (*EVM2EVMMultiOffRampAdminSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*EVM2EVMMultiOffRampConfigSet, error)

	FilterConfigSet0(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampConfigSet0Iterator, error)

	WatchConfigSet0(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampConfigSet0) (event.Subscription, error)

	ParseConfigSet0(log types.Log) (*EVM2EVMMultiOffRampConfigSet0, error)

	FilterExecutionStateChanged(opts *bind.FilterOpts, sequenceNumber []uint64, messageId [][32]byte) (*EVM2EVMMultiOffRampExecutionStateChangedIterator, error)

	WatchExecutionStateChanged(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampExecutionStateChanged, sequenceNumber []uint64, messageId [][32]byte) (event.Subscription, error)

	ParseExecutionStateChanged(log types.Log) (*EVM2EVMMultiOffRampExecutionStateChanged, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOffRampOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*EVM2EVMMultiOffRampOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EVM2EVMMultiOffRampOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*EVM2EVMMultiOffRampOwnershipTransferred, error)

	FilterSkippedIncorrectNonce(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMMultiOffRampSkippedIncorrectNonceIterator, error)

	WatchSkippedIncorrectNonce(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSkippedIncorrectNonce, nonce []uint64, sender []common.Address) (event.Subscription, error)

	ParseSkippedIncorrectNonce(log types.Log) (*EVM2EVMMultiOffRampSkippedIncorrectNonce, error)

	FilterSkippedSenderWithPreviousRampMessageInflight(opts *bind.FilterOpts, nonce []uint64, sender []common.Address) (*EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflightIterator, error)

	WatchSkippedSenderWithPreviousRampMessageInflight(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight, nonce []uint64, sender []common.Address) (event.Subscription, error)

	ParseSkippedSenderWithPreviousRampMessageInflight(log types.Log) (*EVM2EVMMultiOffRampSkippedSenderWithPreviousRampMessageInflight, error)

	FilterSourceChainConfigSet(opts *bind.FilterOpts, sourceChainSelector []uint64) (*EVM2EVMMultiOffRampSourceChainConfigSetIterator, error)

	WatchSourceChainConfigSet(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSourceChainConfigSet, sourceChainSelector []uint64) (event.Subscription, error)

	ParseSourceChainConfigSet(log types.Log) (*EVM2EVMMultiOffRampSourceChainConfigSet, error)

	FilterSourceChainSelectorAdded(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampSourceChainSelectorAddedIterator, error)

	WatchSourceChainSelectorAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampSourceChainSelectorAdded) (event.Subscription, error)

	ParseSourceChainSelectorAdded(log types.Log) (*EVM2EVMMultiOffRampSourceChainSelectorAdded, error)

	FilterTokenAggregateRateLimitAdded(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampTokenAggregateRateLimitAddedIterator, error)

	WatchTokenAggregateRateLimitAdded(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampTokenAggregateRateLimitAdded) (event.Subscription, error)

	ParseTokenAggregateRateLimitAdded(log types.Log) (*EVM2EVMMultiOffRampTokenAggregateRateLimitAdded, error)

	FilterTokenAggregateRateLimitRemoved(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampTokenAggregateRateLimitRemovedIterator, error)

	WatchTokenAggregateRateLimitRemoved(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved) (event.Subscription, error)

	ParseTokenAggregateRateLimitRemoved(log types.Log) (*EVM2EVMMultiOffRampTokenAggregateRateLimitRemoved, error)

	FilterTransmitted(opts *bind.FilterOpts) (*EVM2EVMMultiOffRampTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *EVM2EVMMultiOffRampTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*EVM2EVMMultiOffRampTransmitted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
