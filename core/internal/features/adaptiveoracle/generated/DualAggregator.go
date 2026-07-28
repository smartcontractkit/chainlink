// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package adaptiveoracle

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
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
	_ = time.Tick
	_ = context.Background
)

// DualAggregatorMetaData contains all meta data concerning the DualAggregator contract.
var DualAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"link\",\"type\":\"address\",\"internalType\":\"contractLinkTokenInterface\"},{\"name\":\"minAnswer_\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"maxAnswer_\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"billingAccessController\",\"type\":\"address\",\"internalType\":\"contractAccessControllerInterface\"},{\"name\":\"requesterAccessController\",\"type\":\"address\",\"internalType\":\"contractAccessControllerInterface\"},{\"name\":\"decimals_\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"description_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"secondaryProxy_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"cutoffTime_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxSyncIterations_\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addAccess\",\"inputs\":[{\"name\":\"_user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"disableAccessCheck\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enableAccessCheck\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAdaptiveOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAnswer\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBilling\",\"inputs\":[],\"outputs\":[{\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accountingGas\",\"type\":\"uint24\",\"internalType\":\"uint24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBillingAccessController\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAccessControllerInterface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLinkToken\",\"inputs\":[],\"outputs\":[{\"name\":\"linkToken\",\"type\":\"address\",\"internalType\":\"contractLinkTokenInterface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequesterAccessController\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAccessControllerInterface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"roundId_\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimestamp\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransmitters\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidatorConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"contractAggregatorValidatorInterface\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasAccess\",\"inputs\":[{\"name\":\"_user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_calldata\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestAnswer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDetails\",\"inputs\":[],\"outputs\":[{\"name\":\"configCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestConfigDigestAndEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"scanLogs\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRound\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestTransmissionDetails\",\"inputs\":[],\"outputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"latestAnswer_\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"latestTimestamp_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"linkAvailableForPayment\",\"inputs\":[],\"outputs\":[{\"name\":\"availableBalance\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxAnswer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minAnswer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oracleObservationCount\",\"inputs\":[{\"name\":\"transmitterAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owedPayment\",\"inputs\":[{\"name\":\"transmitterAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAccess\",\"inputs\":[{\"name\":\"_user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestNewRound\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resetToReferenceRate\",\"inputs\":[{\"name\":\"referenceRate\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAdaptiveOracle\",\"inputs\":[{\"name\":\"adaptiveOracle\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBilling\",\"inputs\":[{\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accountingGas\",\"type\":\"uint24\",\"internalType\":\"uint24\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBillingAccessController\",\"inputs\":[{\"name\":\"_billingAccessController\",\"type\":\"address\",\"internalType\":\"contractAccessControllerInterface\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setConfig\",\"inputs\":[{\"name\":\"signers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"onchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCutoffTime\",\"inputs\":[{\"name\":\"_cutoffTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLinkToken\",\"inputs\":[{\"name\":\"linkToken\",\"type\":\"address\",\"internalType\":\"contractLinkTokenInterface\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPayees\",\"inputs\":[{\"name\":\"transmitters\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"payees\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRequesterAccessController\",\"inputs\":[{\"name\":\"requesterAccessController\",\"type\":\"address\",\"internalType\":\"contractAccessControllerInterface\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setValidatorConfig\",\"inputs\":[{\"name\":\"newValidator\",\"type\":\"address\",\"internalType\":\"contractAggregatorValidatorInterface\"},{\"name\":\"newGasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferPayeeship\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transmit\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[3]\",\"internalType\":\"bytes32[3]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transmitSecondary\",\"inputs\":[{\"name\":\"reportContext\",\"type\":\"bytes32[3]\",\"internalType\":\"bytes32[3]\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"ss\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"rawVs\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPayment\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AdaptiveOracleSet\",\"inputs\":[{\"name\":\"old\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AdaptiveRateReset\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"referenceRate\",\"type\":\"int192\",\"indexed\":false,\"internalType\":\"int192\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddedAccess\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AnswerUpdated\",\"inputs\":[{\"name\":\"current\",\"type\":\"int256\",\"indexed\":true,\"internalType\":\"int256\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingAccessControllerSet\",\"inputs\":[{\"name\":\"old\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillingSet\",\"inputs\":[{\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"accountingGas\",\"type\":\"uint24\",\"indexed\":false,\"internalType\":\"uint24\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CheckAccessDisabled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CheckAccessEnabled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigSet\",\"inputs\":[{\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"configCount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signers\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"transmitters\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"f\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"onchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"offchainConfigVersion\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"offchainConfig\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CutoffTimeSet\",\"inputs\":[{\"name\":\"cutoffTime\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LinkTokenSet\",\"inputs\":[{\"name\":\"oldLinkToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\"},{\"name\":\"newLinkToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewRound\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewTransmission\",\"inputs\":[{\"name\":\"aggregatorRoundId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"answer\",\"type\":\"int192\",\"indexed\":false,\"internalType\":\"int192\"},{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"observationsTimestamp\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"observations\",\"type\":\"int192[]\",\"indexed\":false,\"internalType\":\"int192[]\"},{\"name\":\"observers\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"juelsPerFeeCoin\",\"type\":\"int192\",\"indexed\":false,\"internalType\":\"int192\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"epochAndRound\",\"type\":\"uint40\",\"indexed\":false,\"internalType\":\"uint40\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OraclePaid\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payee\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"linkToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferRequested\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PayeeshipTransferred\",\"inputs\":[{\"name\":\"transmitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"previous\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PrimaryFeedUnlocked\",\"inputs\":[{\"name\":\"primaryRoundId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemovedAccess\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequesterAccessControllerSet\",\"inputs\":[{\"name\":\"old\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\"},{\"name\":\"current\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoundRequested\",\"inputs\":[{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SecondaryRoundIdUpdated\",\"inputs\":[{\"name\":\"secondaryRoundId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transmitted\",\"inputs\":[{\"name\":\"configDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorConfigSet\",\"inputs\":[{\"name\":\"previousValidator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\"},{\"name\":\"previousGasLimit\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"currentValidator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\"},{\"name\":\"currentGasLimit\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CalldataLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotTransferPayeeToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ConfigDigestMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FaultyOracleFTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientGas\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOnChainConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LeftGasCannotExceedInitialGas\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxSyncIterationsReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MedianIsOutOfMinMaxRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NumObservationsOutOfBounds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyAdaptiveOracleCanCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCallableByEOA\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyCurrentPayeeCanUpdate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyOwnerAndBillingAdminCanCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyOwnerAndRequesterCanCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPayeeCanWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyProposedPayeesCanAccept\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OracleLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PayeeAlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedSignerAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RepeatedTransmitterAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReportLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureError\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignaturesOutOfRegistration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaleReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooFewValuesToTrustMedian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooManyOracles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferRemainingFundsFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransmittersSizeNotEqualPayeeSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedTransmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongNumberOfSignatures\",\"inputs\":[]}]",
	Bin: "0x61012060405234801561001157600080fd5b506040516153f83803806153f88339810160408190526100309161059e565b33806000816100865760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100b6576100b6816101c0565b50506001805460ff60a01b1916600160a01b1790555060ff851661010052601789810b60805288900b60a0526001600160a01b0380841660c05263ffffffff821660e05260158054918c166001600160a01b0319909216821790556040516000907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a908290a361014587610269565b61014e866102e1565b610159600080610358565b6012805463ffffffff191663ffffffff84169081179091556040519081527fb24a681ce3399a408a89fd0c2b59dfc24bdad592b1c7ec7671cf060596c1c4d19060200160405180910390a160136101b0858261070b565b50505050505050505050506107c9565b336001600160a01b038216036102185760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161007d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6016546001600160a01b0390811690821681146102dd57601680546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291015b60405180910390a15b5050565b6102e961043b565b6010546001600160a01b0390811690821681146102dd57601080546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491016102d4565b61036061043b565b60408051808201909152600f546001600160a01b03808216808452600160a01b90920463ffffffff16602084015284161415806103ad57508163ffffffff16816020015163ffffffff1614155b15610436576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600f80546001600160c01b0319168417600160a01b830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6000546001600160a01b031633146104955760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161007d565b565b6001600160a01b03811681146104ac57600080fd5b50565b8051601781900b81146104c157600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126104ed57600080fd5b81516001600160401b03811115610506576105066104c6565b604051601f8201601f19908116603f011681016001600160401b0381118282101715610534576105346104c6565b60405281815283820160200185101561054c57600080fd5b60005b8281101561056b5760208186018101518383018201520161054f565b506000918101602001919091529392505050565b80516104c181610497565b805163ffffffff811681146104c157600080fd5b6000806000806000806000806000806101408b8d0312156105be57600080fd5b8a516105c981610497565b99506105d760208c016104af565b98506105e560408c016104af565b975060608b01516105f581610497565b60808c015190975061060681610497565b60a08c015190965060ff8116811461061d57600080fd5b60c08c01519095506001600160401b0381111561063957600080fd5b6106458d828e016104dc565b94505061065460e08c0161057f565b92506106636101008c0161058a565b91506106726101208c0161058a565b90509295989b9194979a5092959850565b600181811c9082168061069757607f821691505b6020821081036106b757634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111561043657806000526020600020601f840160051c810160208510156106e45750805b601f840160051c820191505b8181101561070457600081556001016106f0565b5050505050565b81516001600160401b03811115610724576107246104c6565b610738816107328454610683565b846106bd565b6020601f82116001811461076c57600083156107545750848201515b600019600385901b1c1916600184901b178455610704565b600084815260208120601f198516915b8281101561079c578785015182556020948501946001909201910161077c565b50848210156107ba5786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b60805160a05160c05160e05161010051614bc861083060003960006104090152600081816131bd01526133d4015260006128ad0152600081816104cf01528181611bb501526138ca01526000818161038301528181611b8d015261389f0152614bc86000f3fe608060405234801561001057600080fd5b506004361061030c5760003560e01c80639a6fc8f51161019d578063c4c92b37116100e9578063e5fe4577116100a2578063eb5dcd6c1161007c578063eb5dcd6c14610826578063f2fde38b14610839578063fbffd2c11461084c578063feaf968c1461085f57600080fd5b8063e5fe4577146107b9578063e76d516814610802578063eb4571631461081357600080fd5b8063c4c92b3714610740578063d09dc33914610751578063daffc4b514610759578063dc7f01241461076a578063e3d0e7121461077e578063e4902f821461079157600080fd5b8063b121e14711610156578063b5ab58dc11610130578063b5ab58dc146106f4578063b633620c14610707578063ba0cb29e1461071a578063c10753291461072d57600080fd5b8063b121e147146106bb578063b17f2a6b146106ce578063b1dc65a4146106e157600080fd5b80639a6fc8f5146105c65780639bd2c0b1146106105780639c849b30146106515780639e3ceeab14610664578063a118f24914610677578063afcb95d71461068a57600080fd5b80636b14daf81161025c5780638205bf6a116102155780638da5cb5b116101ef5780638da5cb5b1461056c578063945260481461057d57806398c89ea21461059057806398e5b12a146105a357600080fd5b80638205bf6a1461053e5780638823da6c146105465780638ac28d5a1461055957600080fd5b80636b14daf8146104aa57806370da2f67146104cd5780637284e416146104f657806379ba5097146104fe5780638038e4a11461050657806381ff70481461050e57600080fd5b8063364efd83116102c957806354fd4d50116102a357806354fd4d5014610473578063643dc1051461047a578063666cab8d1461048d578063668a0f02146104a257600080fd5b8063364efd83146104335780634fb174701461045857806350d25bcd1461046b57600080fd5b80630a756983146103115780630eafb25b1461031b578063181f5a771461034157806322adbc781461038157806329937268146103aa578063313ce56714610402575b600080fd5b610319610867565b005b61032e610329366004614097565b6108ba565b6040519081526020015b60405180910390f35b6103746040518060400160405280601481526020017304475616c41676772656761746f7220312e302e360641b81525081565b60405161033891906140fa565b7f000000000000000000000000000000000000000000000000000000000000000060170b61032e565b600d54600c546040805163ffffffff600160701b850481168252600160901b850481166020830152600160b01b8504811692820192909252600160d01b90930416606083015262ffffff16608082015260a001610338565b60405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610338565b6014546001600160a01b03165b6040516001600160a01b039091168152602001610338565b61031961046636600461410d565b6109bf565b61032e610ba8565b600661032e565b610319610488366004614158565b610bd8565b610495610d6c565b6040516103389190614216565b61032e610dce565b6104bd6104b83660046142de565b610de3565b6040519015158152602001610338565b7f000000000000000000000000000000000000000000000000000000000000000060170b61032e565b610374610e0b565b610319610e94565b610319610f43565b600e54600b546040805163ffffffff80851682526401000000009094049093166020840152820152606001610338565b61032e610f9a565b610319610554366004614097565b610fd1565b610319610567366004614097565b611053565b6000546001600160a01b0316610440565b61031961058b36600461432d565b611096565b61031961059e366004614097565b6111f2565b6105ab611262565b60405169ffffffffffffffffffff9091168152602001610338565b6105d96105d4366004614346565b611393565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a001610338565b604080518082018252600f546001600160a01b038116808352600160a01b90910463ffffffff16602092830181905283519182529181019190915201610338565b61031961065f3660046143bd565b611440565b610319610672366004614097565b6115bb565b610319610685366004614097565b61163a565b600b54600d546040805160008152602081019390935261010090910460081c63ffffffff1690820152606001610338565b6103196106c9366004614097565b6116b6565b6103196106dc36600461442c565b611763565b6103196106ef366004614449565b6117b3565b61032e61070236600461432d565b6117cf565b61032e61071536600461432d565b61180d565b610319610728366004614449565b611852565b61031961073b36600461453a565b611864565b6016546001600160a01b0316610440565b61032e611a66565b6010546001600160a01b0316610440565b6001546104bd90600160a01b900460ff1681565b61031961078c36600461462c565b611af6565b6107a461079f366004614097565b61211b565b60405163ffffffff9091168152602001610338565b6107c16121d1565b6040805195865263ffffffff909416602086015260ff9092169284019290925260179190910b60608301526001600160401b0316608082015260a001610338565b6015546001600160a01b0316610440565b610319610821366004614708565b612250565b61031961083436600461410d565b612333565b610319610847366004614097565b61240d565b61031961085a366004614097565b61241e565b6105d961242f565b61086f6124af565b600154600160a01b900460ff16156108b8576001805460ff60a01b191690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b6001600160a01b03811660009081526003602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b0316918101919091529061091c5750600092915050565b600d546020820151600091600160b01b900463ffffffff169060079060ff16601f811061094b5761094b614736565b600881049190910154600d5461097e926007166004026101000a90910463ffffffff90811691600160301b900416614762565b63ffffffff1661098e919061477e565b61099c90633b9aca0061477e565b905081604001516001600160601b0316816109b79190614795565b949350505050565b6109c76124af565b6015546001600160a01b039081169083168190036109e457505050565b6040516370a0823160e01b81523060048201526001600160a01b038416906370a0823190602401602060405180830381865afa158015610a28573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a4c91906147a8565b50610a55612502565b6040516370a0823160e01b81523060048201526000906001600160a01b038316906370a0823190602401602060405180830381865afa158015610a9c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ac091906147a8565b60405163a9059cbb60e01b81526001600160a01b038581166004830152602482018390529192509083169063a9059cbb906044016020604051808303816000875af1158015610b13573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b3791906147c1565b610b5457604051633b92843d60e11b815260040160405180910390fd5b601580546001600160a01b0319166001600160a01b0386811691821790925560405190918416907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a90600090a350505b5050565b600060116000610bb6612878565b63ffffffff16815260208101919091526040016000206001015460170b919050565b6000546001600160a01b0316331480610c625750601654604051630d629b5f60e31b81526001600160a01b0390911690636b14daf890610c2190339060009036906004016147e3565b602060405180830381865afa158015610c3e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c6291906147c1565b610c7f576040516391ed77c560e01b815260040160405180910390fd5b610c87612502565b600d805467ffffffffffffffff60701b1916600160701b63ffffffff88811691820263ffffffff60901b191692909217600160901b8884169081029190911767ffffffffffffffff60b01b1916600160b01b88851690810263ffffffff60d01b191691909117600160d01b94881694850217909455600c805462ffffff191662ffffff871690811790915560408051938452602084019290925290820193909352606081019190915260808101919091527f0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f9060a00160405180910390a15050505050565b60606006805480602002602001604051908101604052809291908181526020018280548015610dc457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610da6575b5050505050905090565b6000610dd8612878565b63ffffffff16905090565b6000610def8383612982565b80610e0257506001600160a01b03831632145b90505b92915050565b606060138054610e1a90614823565b80601f0160208091040260200160405190810160405280929190818152602001828054610e4690614823565b8015610dc45780601f10610e6857610100808354040283529160200191610dc4565b820191906000526020600020905b815481529060010190602001808311610e7657509395945050505050565b6001546001600160a01b03163314610eec5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064015b60405180910390fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610f4b6124af565b600154600160a01b900460ff166108b8576001805460ff60a01b1916600160a01b1790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b600060116000610fa8612878565b63ffffffff9081168252602082019290925260400160002060010154600160e01b900416919050565b610fd96124af565b6001600160a01b03811660009081526002602052604090205460ff1615611050576001600160a01b038116600081815260026020908152604091829020805460ff1916905590519182527f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d191015b60405180910390a15b50565b6001600160a01b0381811660009081526017602052604090205416331461108d57604051632ab4a3db60e01b815260040160405180910390fd5b611050816129b9565b6014546001600160a01b031633146110c157604051633c0ce55f60e21b815260040160405180910390fd5b600d546000906110df90600160301b900463ffffffff166001614857565b600d80546dffffffffffffffff0000000000001916600160301b63ffffffff84811691820263ffffffff60501b191692909217600160501b82021790925560408051608081018252601787900b808252602080830182815242861684860181815260608601918252600089815260118552879020955186546001600160c01b0319166001600160c01b039182161787559251600196909601805491519251969093166001600160e01b031990911617600160c01b91881691909102176001600160e01b0316600160e01b949096169390930294909417909155905190815292935084927f6f643b21ecba6d2a464beb4b6fa87271920464cfdea7d6ff1b75dd3c7f143747910160405180910390a2505050565b6111fa6124af565b6014546001600160a01b039081169082168114610ba457601480546001600160a01b0319166001600160a01b0384811691821790925560405190918316907f1969b480ede741c56a9a724f5457349d09b3ee4574a1dbb1b75a2e668b68050490600090a35050565b600080546001600160a01b031633148015906112f15750601054604051630d629b5f60e31b81526001600160a01b0390911690636b14daf8906112ae90339060009036906004016147e3565b602060405180830381865afa1580156112cb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112ef91906147c1565b155b1561130f5760405163099b888b60e31b815260040160405180910390fd5b600d54600b546040805191825263ffffffff6101008404600881901c8216602085015260ff811684840152915164ffffffffff90921693600160301b9004169133917f41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c9181900360600190a2611386816001614857565b63ffffffff169250505090565b60008060008060006113a3612878565b63ffffffff168669ffffffffffffffffffff1611156113d057506000935083925082915081905080611437565b5050505063ffffffff82811660009081526011602090815260409182902082516080810184528154601790810b82526001909201549182900b928101839052600160c01b82048516938101849052600160e01b909104909316606090930183905284935091835b91939590929450565b6114486124af565b82811461146857604051633d2f942960e01b815260040160405180910390fd5b60005b838110156115b457600085858381811061148757611487614736565b905060200201602081019061149c9190614097565b905060008484848181106114b2576114b2614736565b90506020020160208101906114c79190614097565b6001600160a01b038084166000908152601760205260409020549192501680158015816115065750826001600160a01b0316826001600160a01b031614155b15611524576040516315d5c0c560e31b815260040160405180910390fd5b6001600160a01b03848116600090815260176020526040902080546001600160a01b031916858316908117909155908316146115a557826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050505080600101905061146b565b5050505050565b6115c36124af565b6010546001600160a01b039081169082168114610ba457601080546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491015b60405180910390a15050565b6116426124af565b6001600160a01b03811660009081526002602052604090205460ff16611050576001600160a01b038116600081815260026020908152604091829020805460ff1916600117905590519182527f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49101611047565b6001600160a01b038181166000908152601860205260409020541633146116f0576040516332cce5df60e11b815260040160405180910390fd5b6001600160a01b0381811660008181526017602090815260408083208054336001600160a01b031980831682179093556018909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b61176b6124af565b6012805463ffffffff191663ffffffff83169081179091556040519081527fb24a681ce3399a408a89fd0c2b59dfc24bdad592b1c7ec7671cf060596c1c4d190602001611047565b6117c588888888888888886000612ba2565b5050505050505050565b60006117d9612878565b63ffffffff168211156117ee57506000919050565b5063ffffffff1660009081526011602052604090206001015460170b90565b6000611817612878565b63ffffffff1682111561182c57506000919050565b5063ffffffff908116600090815260116020526040902060010154600160e01b90041690565b6117c588888888888888886001612ba2565b6000546001600160a01b031633148015906118f25750601654604051630d629b5f60e31b81526001600160a01b0390911690636b14daf8906118af90339060009036906004016147e3565b602060405180830381865afa1580156118cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118f091906147c1565b155b15611910576040516391ed77c560e01b815260040160405180910390fd5b600061191a612e33565b6015546040516370a0823160e01b81523060048201529192506000916001600160a01b03909116906370a0823190602401602060405180830381865afa158015611968573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061198c91906147a8565b9050818110156119af57604051631e9acf1760e31b815260040160405180910390fd5b6015546001600160a01b031663a9059cbb856119d46119ce8686614873565b87612fe9565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015611a1f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a4391906147c1565b611a605760405163356680b760e01b815260040160405180910390fd5b50505050565b6015546040516370a0823160e01b815230600482015260009182916001600160a01b03909116906370a0823190602401602060405180830381865afa158015611ab3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ad791906147a8565b90506000611ae3612e33565b9050611aef8183614886565b9250505090565b611afe6124af565b601f86511115611b2157604051630974082760e21b815260040160405180910390fd5b8451865114611b43576040516304a14cb760e31b815260040160405180910390fd5b8551611b508560036148a6565b60ff1610611b7157604051631064b94d60e11b815260040160405180910390fd5b611b7d8460ff16613000565b60408051600160f81b60208201527f0000000000000000000000000000000000000000000000000000000000000000821b60218201527f000000000000000000000000000000000000000000000000000000000000000090911b603982015260510160405160208183030381529060405280519060200120838051906020012014611c1b576040516354408ee360e11b815260040160405180910390fd5b6040805160c0810182526001600160401b038416815260ff86166020820152908101849052606081018290526080810187905260a08101869052600d805465ffffffffff0019169055611c6c612502565b60055460005b81811015611d1357600060058281548110611c8f57611c8f614736565b6000918252602082200154600680546001600160a01b0390921693509084908110611cbc57611cbc614736565b60009182526020808320909101546001600160a01b039485168352600482526040808420805461ffff1916905594168252600390529190912080546dffffffffffffffffffffffffffff1916905550600101611c72565b50611d2060056000613f57565b611d2c60066000613f57565b60005b826080015151811015611f48576004600084608001518381518110611d5657611d56614736565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff1615611d9b576040516316c6131560e01b815260040160405180910390fd5b60405180604001604052806001151581526020018260ff168152506004600085608001518481518110611dd057611dd0614736565b6020908102919091018101516001600160a01b0316825281810192909252604001600090812083518154949093015160ff166101000261ff00199315159390931661ffff19909416939093179190911790915560a08401518051600392919084908110611e3f57611e3f614736565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff1615611e845760405163358f4d1d60e21b815260040160405180910390fd5b60405180606001604052806001151581526020018260ff16815260200160006001600160601b0316815250600360008560a001518481518110611ec957611ec9614736565b6020908102919091018101516001600160a01b03168252818101929092526040908101600020835181549385015194909201516001600160601b0316620100000262010000600160701b031960ff959095166101000261ff00199315159390931661ffff19909416939093179190911792909216179055600101611d2f565b5060808201518051611f6291600591602090910190613f75565b5060a08201518051611f7c91600691602090910190613f75565b506020820151600d805460ff191660ff909216919091179055600e805467ffffffff0000000019811664010000000063ffffffff438116820292831785559083048116936001939092600092611fd9928692908216911617614857565b92506101000a81548163ffffffff021916908363ffffffff1602179055506120384630600e60009054906101000a900463ffffffff1663ffffffff1686608001518760a00151886020015189604001518a600001518b60600151613021565b600b819055600e54608085015160a08601516020870151604080890151895160608b015192517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e059861209e988b98919763ffffffff9091169691959094919391926148c2565b60405180910390a1600d54600160301b900463ffffffff1660005b84608001515181101561210e5781600782601f81106120da576120da614736565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff1602179055508060010190506120b9565b5050505050505050505050565b6001600160a01b03811660009081526003602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046001600160601b0316918101919091529061217d5750600092915050565b6007816020015160ff16601f811061219757612197614736565b600881049190910154600d546121ca926007166004026101000a90910463ffffffff90811691600160301b900416614762565b9392505050565b6000808080803332146121f7576040516374e2cd5160e01b815260040160405180910390fd5b5050600b54600d5463ffffffff600160301b82048116600090815260116020526040902080546001909101549397610100909304600881901c8316975064ffffffffff16955060170b9350600160e01b90920490911690565b6122586124af565b60408051808201909152600f546001600160a01b03808216808452600160a01b90920463ffffffff16602084015284161415806122a557508163ffffffff16816020015163ffffffff1614155b1561232e576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600f80546001600160c01b0319168417600160a01b830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6001600160a01b0382811660009081526017602052604090205416331461236d57604051635cbe80b560e11b815260040160405180910390fd5b6001600160a01b038116330361239657604051633cef863360e11b815260040160405180910390fd5b6001600160a01b03808316600090815260186020526040902080548383166001600160a01b03198216811790925590911690811461232e576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b6124156124af565b6110508161307d565b6124266124af565b61105081613126565b600080600080600080612440612878565b63ffffffff90811660008181526011602090815260409182902082516080810184528154601790810b82526001909201549182900b928101839052600160c01b82048616938101849052600160e01b909104909416606090940184905291999198509650909450879350915050565b6000546001600160a01b031633146108b85760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610ee3565b601554600d54604080516103e08101918290526001600160a01b0390931692600160301b90920463ffffffff1691600091600790601f908285855b82829054906101000a900463ffffffff1663ffffffff168152602001906004019060208260030104928301926001038202915080841161253d57905050505050509050600060068054806020026020016040519081016040528092919081815260200182805480156125d857602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116125ba575b5050505050905060005b815181101561286a5760006003600084848151811061260357612603614736565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160029054906101000a90046001600160601b03166001600160601b0316905060006003600085858151811061266557612665614736565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060000160026101000a8154816001600160601b0302191690836001600160601b0316021790555060008483601f81106126c8576126c8614736565b6020020151600d5490870363ffffffff9081169250600160b01b909104168102633b9aca00028201801561285f5760006017600087878151811061270e5761270e614736565b6020908102919091018101516001600160a01b03908116835290820192909252604090810160002054905163a9059cbb60e01b815290821660048201819052602482018590529250908a169063a9059cbb906044016020604051808303816000875af1158015612782573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127a691906147c1565b6127c35760405163356680b760e01b815260040160405180910390fd5b878786601f81106127d6576127d6614736565b602002019063ffffffff16908163ffffffff1681525050886001600160a01b0316816001600160a01b031687878151811061281357612813614736565b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c8560405161285591815260200190565b60405180910390a4505b5050506001016125e2565b506115b4600783601f613fda565b600d5460009063ffffffff600160301b8204811691600160501b81049091169060ff600160f01b909104166001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016330361292a5760125463ffffffff83811660009081526011602052604090206001015442926129079290811691600160e01b900416614857565b63ffffffff1610156129235761291b613195565b935050505090565b5092915050565b8163ffffffff168363ffffffff160361297a5780801561296a575063ffffffff808416600090815260116020526040902060010154600160e01b90041642145b1561297a5761291b600184614762565b509092915050565b6001600160a01b03821660009081526002602052604081205460ff1680610e02575050600154600160a01b900460ff161592915050565b6001600160a01b0381166000908152600360209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046001600160601b031692810192909252612a16575050565b6000612a21836108ba565b9050801561232e576001600160a01b038381166000908152601760205260409081902054601554915163a9059cbb60e01b8152908316600482018190526024820185905292919091169063a9059cbb906044016020604051808303816000875af1158015612a93573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ab791906147c1565b612ad45760405163356680b760e01b815260040160405180910390fd5b600d60000160069054906101000a900463ffffffff166007846020015160ff16601f8110612b0457612b04614736565b6008810491909101805460079092166004026101000a63ffffffff8181021990931693909216919091029190911790556001600160a01b03848116600081815260036020908152604091829020805462010000600160701b0319169055601554915186815291841693851692917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c910160405180910390a450505050565b60005a9050612bb38a89888761325f565b6000612bf48a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061331092505050565b6040805161012081018252600d5460ff808216835261010080830464ffffffffff166020850152600160301b830463ffffffff90811695850195909552600160501b830485166060850152600160701b830485166080850152600160901b8304851660a0850152600160b01b8304851660c0850152600160d01b830490941660e0840152600160f01b909104161515918101919091529091508315612d4257600080612c9f846133aa565b915091508115612d3f578063ffffffff16836060015163ffffffff1610612cd957604051637c01d16560e11b815260040160405180910390fd5b600d805463ffffffff60501b1916600160501b63ffffffff8416908102919091179091556040517f8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b90600090a2612d358385600001518761351a565b5050505050612e28565b50505b602081810151908d01359064ffffffffff8083169116141580612d685750816101000151155b15612dc257816020015164ffffffffff168164ffffffffff1611612d9f57604051637c01d16560e11b815260040160405180910390fd5b612daf8d8d8d8d8d8d8d8d61362b565b612dbd828e35838689613805565b612dfe565b600d54604051600160301b90910463ffffffff16907fda2435684a37fba6f7841e49b59e6ad975e462bbebd28ec9da4ed9746a6992be90600090a25b600d805460ff60f01b1916600160f01b871515021790558251612e239083908661351a565b505050505b505050505050505050565b6000806006805480602002602001604051908101604052809291908181526020018280548015612e8c57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612e6e575b50508351600d54604080516103e08101918290529697509195600160301b90910463ffffffff169450600093509150600790601f908285855b82829054906101000a900463ffffffff1663ffffffff1681526020019060040190602082600301049283019260010382029150808411612ec55790505050505050905060005b83811015612f4e578181601f8110612f2557612f25614736565b6020020151612f349084614762565b612f449063ffffffff1687614795565b9550600101612f0b565b50600d54612f6d90600160b01b900463ffffffff16633b9aca0061477e565b612f77908661477e565b945060005b83811015612fe15760036000868381518110612f9a57612f9a614736565b6020908102919091018101516001600160a01b0316825281019190915260400160002054612fd7906201000090046001600160601b031687614795565b9550600101612f7c565b505050505090565b600081831015612ffa575081610e05565b50919050565b600081116110505760405163039d1a4d60e41b815260040160405180910390fd5b6000808a8a8a8a8a8a8a8a8a60405160200161304599989796959493929190614958565b60408051808303601f1901815291905280516020909101206001600160f01b0316600160f01b179150505b9998505050505050505050565b336001600160a01b038216036130d55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610ee3565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6016546001600160a01b039081169082168114610ba457601680546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912910161162e565b600d54600090600160301b900463ffffffff16805b63ffffffff8116156132485763ffffffff7f0000000000000000000000000000000000000000000000000000000000000000166131e78284614762565b63ffffffff16146132485760125463ffffffff82811660009081526011602052604090206001015442926132269290811691600160e01b900416614857565b63ffffffff1610156132385792915050565b61324181614992565b90506131aa565b5050600d54600160501b900463ffffffff16919050565b3360009081526003602052604090205460ff1661328f57604051631b41e11d60e31b815260040160405180910390fd5b600b548435146132b25760405163dfdcf8e760e01b815260040160405180910390fd5b6132bd838383613d21565b600d546132ce9060ff1660016149b2565b60ff1682146132f0576040516371253a2560e01b815260040160405180910390fd5b808214611a605760405163a75d88af60e01b815260040160405180910390fd5b60408051608081018252600080825260208201526060918101829052818101919091526000806000808580602001905181019061334d91906149dd565b935093509350935061335f8683613d86565b815160408051602081018690526000910160408051918152928152825160808101845260179490940b845263ffffffff90961660208401525081019390935260608301525092915050565b600d546000908190600160301b900463ffffffff16805b63ffffffff81161561350d5763ffffffff7f0000000000000000000000000000000000000000000000000000000000000000166133fe8284614762565b63ffffffff1603613422576040516330f4b89560e21b815260040160405180910390fd5b63ffffffff80821660009081526011602090815260409182902082516080810184528154601790810b82526001909201549182900b81840152600160c01b82048516938101849052600160e01b9091048416606082015290880151909216111561349457506000958695509350505050565b856020015163ffffffff16816040015163ffffffff161480156134ea5750606086015180516134c590600290614aa8565b815181106134d5576134d5614736565b602002602001015160170b816000015160170b145b156134fc575060019590945092505050565b5061350681614992565b90506133c1565b5060009485945092505050565b60008260170b121561352b57505050565b6000613552633b9aca003a048560a0015163ffffffff16866080015163ffffffff16613dcf565b90506010360260005a600c5490915060009061357f9063ffffffff8716908690869062ffffff1686613df5565b90506000670de0b6b3a76400006001600160c01b03881683023360009081526003602052604090205460e08b01519290910492506201000090046001600160601b039081169163ffffffff16633b9aca0002828401019081168211156135eb5750505050505050505050565b33600090815260036020526040902080546001600160601b03909216620100000262010000600160701b0319909216919091179055505050505050505050565b6000878760405161363d929190614aca565b604051908190038120613654918b90602001614ada565b60408051601f19818403018152828252805160209182012083830190925260008084529083018190529092509060005b878110156137c35760006001858784602081106136a3576136a3614736565b6136b091901a601b6149b2565b8c8c868181106136c2576136c2614736565b905060200201358b8b878181106136db576136db614736565b9050602002013560405160008152602001604052604051613718949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa15801561373a573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526004602090815290849020838501909452925460ff80821615158085526101009092041693830193909352909550925090506137a45760405163669233e360e11b815260040160405180910390fd5b826020015160080260ff166001901b8401935050806001019050613684565b5081827e01010101010101010101010101010101010101010101010101010101010101161461210e57604051638044bb3360e01b815260040160405180910390fd5b601f826060015151111561382f5760405160016293ddfb60e01b0319815260040160405180910390fd5b846000015160ff168260600151511161385b57604051635765bdd760e01b815260040160405180910390fd5b64ffffffffff83166020860152606082015180516000919061387f90600290614aa8565b8151811061388f5761388f614736565b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b13806138f257507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b135b156139105760405163650c8d9360e11b815260040160405180910390fd5b6040808701805163ffffffff16600090815260116020529190912054815160179190910b9161393e82614af0565b63ffffffff1690525060145482906001600160a01b031680156139d457604051632861711b60e21b8152601785810b600483015284900b60248201526001600160a01b0382169063a185c46c906044016020604051808303816000875af11580156139ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139d191906147a8565b91505b60408051608081018252601786810b825284900b60208083019182528981015163ffffffff908116848601908152428216606086019081528f87015183166000908152601190945295909220935184546001600160c01b039182166001600160c01b0319909116178555925160019094018054925195518216600160e01b026001600160e01b0396909216600160c01b026001600160e01b0319909316949093169390931717929092161790558415613ac3576040808a015163ffffffff1660608b0181905290517f8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b90600090a25b88600d60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600e6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160000160126101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160000160166101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600001601a6101000a81548163ffffffff021916908363ffffffff16021790555061010082015181600001601e6101000a81548160ff021916908315150217905550905050886040015163ffffffff167fc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a853389602001518a606001518b604001518c600001518f8f604051613c76989796959493929190614b15565b60405180910390a26040808a0151602080890151925163ffffffff9384168152600093909216917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271910160405180910390a3886040015163ffffffff168460170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f42604051613d0891815260200190565b60405180910390a3612e2889604001518560170b613e2a565b6000613d2e82602061477e565b613d3984602061477e565b613d4586610144614795565b613d4f9190614795565b613d599190614795565b613d64906000614795565b9050368114611a605760405163b4d895d560e01b815260040160405180910390fd5b600081516020613d96919061477e565b613da19060a0614795565b613dac906000614795565b90508083511461232e576040516306a70a0b60e51b815260040160405180910390fd5b60008383811015613de257600285850304015b613dec8184612fe9565b95945050505050565b600081861015613e185760405163fbf484ab60e01b815260040160405180910390fd5b50633b9aca0094039190910101020290565b60408051808201909152600f546001600160a01b038116808352600160a01b90910463ffffffff166020830152613e6057505050565b6000613e6d600185614762565b63ffffffff818116600081815260116020526040808220549051602481019390935260170b60448301819052928816606483015260848201879052929350909190613ef49060a40160408051601f19818403018152919052602080820180516001600160e01b031663beed9b5160e01b17905286519087015163ffffffff16611388613f1d565b91505080613f15576040516307099c5360e21b815260040160405180910390fd5b505050505050565b6000805a838110613f4d57839003604081048103851015613f4d57600080885160208a0160008a8af19250600191505b5094509492505050565b5080546000825590600052602060002090810190611050919061406d565b828054828255906000526020600020908101928215613fca579160200282015b82811115613fca57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190613f95565b50613fd692915061406d565b5090565b600483019183908215613fca5791602002820160005b8382111561403457835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302613ff0565b80156140645782816101000a81549063ffffffff0219169055600401602081600301049283019260010302614034565b5050613fd69291505b5b80821115613fd6576000815560010161406e565b6001600160a01b038116811461105057600080fd5b6000602082840312156140a957600080fd5b81356121ca81614082565b6000815180845260005b818110156140da576020818501810151868301820152016140be565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000610e0260208301846140b4565b6000806040838503121561412057600080fd5b823561412b81614082565b9150602083013561413b81614082565b809150509250929050565b63ffffffff8116811461105057600080fd5b600080600080600060a0868803121561417057600080fd5b853561417b81614146565b9450602086013561418b81614146565b9350604086013561419b81614146565b925060608601356141ab81614146565b9150608086013562ffffff811681146141c357600080fd5b809150509295509295909350565b600081518084526020840193506020830160005b8281101561420c5781516001600160a01b03168652602095860195909101906001016141e5565b5093949350505050565b602081526000610e0260208301846141d1565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171561426757614267614229565b604052919050565b600082601f83011261428057600080fd5b81356001600160401b0381111561429957614299614229565b6142ac601f8201601f191660200161423f565b8181528460208386010111156142c157600080fd5b816020850160208301376000918101602001919091529392505050565b600080604083850312156142f157600080fd5b82356142fc81614082565b915060208301356001600160401b0381111561431757600080fd5b6143238582860161426f565b9150509250929050565b60006020828403121561433f57600080fd5b5035919050565b60006020828403121561435857600080fd5b813569ffffffffffffffffffff811681146121ca57600080fd5b60008083601f84011261438457600080fd5b5081356001600160401b0381111561439b57600080fd5b6020830191508360208260051b85010111156143b657600080fd5b9250929050565b600080600080604085870312156143d357600080fd5b84356001600160401b038111156143e957600080fd5b6143f587828801614372565b90955093505060208501356001600160401b0381111561441457600080fd5b61442087828801614372565b95989497509550505050565b60006020828403121561443e57600080fd5b81356121ca81614146565b60008060008060008060008060e0898b03121561446557600080fd5b606089018a81111561447657600080fd5b899850356001600160401b0381111561448e57600080fd5b8901601f81018b1361449f57600080fd5b80356001600160401b038111156144b557600080fd5b8b60208284010111156144c757600080fd5b6020919091019750955060808901356001600160401b038111156144ea57600080fd5b6144f68b828c01614372565b90965094505060a08901356001600160401b0381111561451557600080fd5b6145218b828c01614372565b999c989b50969995989497949560c00135949350505050565b6000806040838503121561454d57600080fd5b823561455881614082565b946020939093013593505050565b60006001600160401b0382111561457f5761457f614229565b5060051b60200190565b600082601f83011261459a57600080fd5b81356145ad6145a882614566565b61423f565b8082825260208201915060208360051b8601019250858311156145cf57600080fd5b602085015b838110156145f55780356145e781614082565b8352602092830192016145d4565b5095945050505050565b803560ff8116811461461057600080fd5b919050565b80356001600160401b038116811461461057600080fd5b60008060008060008060c0878903121561464557600080fd5b86356001600160401b0381111561465b57600080fd5b61466789828a01614589565b96505060208701356001600160401b0381111561468357600080fd5b61468f89828a01614589565b95505061469e604088016145ff565b935060608701356001600160401b038111156146b957600080fd5b6146c589828a0161426f565b9350506146d460808801614615565b915060a08701356001600160401b038111156146ef57600080fd5b6146fb89828a0161426f565b9150509295509295509295565b6000806040838503121561471b57600080fd5b823561472681614082565b9150602083013561413b81614146565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b63ffffffff8281168282160390811115610e0557610e0561474c565b8082028115828204841417610e0557610e0561474c565b80820180821115610e0557610e0561474c565b6000602082840312156147ba57600080fd5b5051919050565b6000602082840312156147d357600080fd5b815180151581146121ca57600080fd5b6001600160a01b03841681526040602082018190528101829052818360608301376000818301606090810191909152601f909201601f1916010192915050565b600181811c9082168061483757607f821691505b602082108103612ffa57634e487b7160e01b600052602260045260246000fd5b63ffffffff8181168382160190811115610e0557610e0561474c565b81810381811115610e0557610e0561474c565b81810360008312801583831316838312821617156129235761292361474c565b60ff81811683821602908116908181146129235761292361474c565b63ffffffff8a16815288602082015263ffffffff88166040820152610120606082015260006148f56101208301896141d1565b828103608084015261490781896141d1565b905060ff871660a084015282810360c084015261492481876140b4565b90506001600160401b03851660e084015282810361010084015261494881856140b4565b9c9b505050505050505050505050565b8981526001600160a01b03891660208201526001600160401b0388166040820152610120606082018190526000906148f5908301896141d1565b600063ffffffff8216806149a8576149a861474c565b6000190192915050565b60ff8181168382160190811115610e0557610e0561474c565b8051601781900b811461461057600080fd5b600080600080608085870312156149f357600080fd5b84516149fe81614146565b6020860151604087015191955093506001600160401b03811115614a2157600080fd5b8501601f81018713614a3257600080fd5b8051614a406145a882614566565b8082825260208201915060208360051b850101925089831115614a6257600080fd5b6020840193505b82841015614a8b57614a7a846149cb565b825260209384019390910190614a69565b9450614a9d92505050606086016149cb565b905092959194509250565b600082614ac557634e487b7160e01b600052601260045260246000fd5b500490565b8183823760009101908152919050565b8281526080810160608360208401379392505050565b600063ffffffff821663ffffffff8103614b0c57614b0c61474c565b60010192915050565b600061010082018a60170b835260018060a01b038a16602084015263ffffffff8916604084015261010060608401528088518083526101208501915060208a01925060005b81811015614b7b57835160170b835260209384019390920191600101614b5a565b50508381036080850152614b8f81896140b4565b92505050614ba260a083018660170b9052565b8360c083015261307060e083018464ffffffffff16905256fea164736f6c634300081a000a",
}

// DualAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use DualAggregatorMetaData.ABI instead.
var DualAggregatorABI = DualAggregatorMetaData.ABI

// DualAggregatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DualAggregatorMetaData.Bin instead.
var DualAggregatorBin = DualAggregatorMetaData.Bin

// DeployDualAggregator deploys a new Ethereum contract, binding an instance of DualAggregator to it.
func DeployDualAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, minAnswer_ *big.Int, maxAnswer_ *big.Int, billingAccessController common.Address, requesterAccessController common.Address, decimals_ uint8, description_ string, secondaryProxy_ common.Address, cutoffTime_ uint32, maxSyncIterations_ uint32) (common.Address, *types.Transaction, *DualAggregator, error) {
	parsed, err := DualAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DualAggregatorBin), backend, link, minAnswer_, maxAnswer_, billingAccessController, requesterAccessController, decimals_, description_, secondaryProxy_, cutoffTime_, maxSyncIterations_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DualAggregator{DualAggregatorCaller: DualAggregatorCaller{contract: contract}, DualAggregatorTransactor: DualAggregatorTransactor{contract: contract}, DualAggregatorFilterer: DualAggregatorFilterer{contract: contract}}, nil
}

// DualAggregator is an auto generated Go binding around an Ethereum contract.
type DualAggregator struct {
	DualAggregatorCaller     // Read-only binding to the contract
	DualAggregatorTransactor // Write-only binding to the contract
	DualAggregatorFilterer   // Log filterer for contract events
}

// DualAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DualAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DualAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DualAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DualAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DualAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DualAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DualAggregatorSession struct {
	Contract     *DualAggregator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DualAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DualAggregatorCallerSession struct {
	Contract *DualAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// DualAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DualAggregatorTransactorSession struct {
	Contract     *DualAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DualAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DualAggregatorRaw struct {
	Contract *DualAggregator // Generic contract binding to access the raw methods on
}

// DualAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DualAggregatorCallerRaw struct {
	Contract *DualAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// DualAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DualAggregatorTransactorRaw struct {
	Contract *DualAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDualAggregator creates a new instance of DualAggregator, bound to a specific deployed contract.
func NewDualAggregator(address common.Address, backend bind.ContractBackend) (*DualAggregator, error) {
	contract, err := bindDualAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DualAggregator{DualAggregatorCaller: DualAggregatorCaller{contract: contract}, DualAggregatorTransactor: DualAggregatorTransactor{contract: contract}, DualAggregatorFilterer: DualAggregatorFilterer{contract: contract}}, nil
}

// NewDualAggregatorCaller creates a new read-only instance of DualAggregator, bound to a specific deployed contract.
func NewDualAggregatorCaller(address common.Address, caller bind.ContractCaller) (*DualAggregatorCaller, error) {
	contract, err := bindDualAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCaller{contract: contract}, nil
}

// NewDualAggregatorTransactor creates a new write-only instance of DualAggregator, bound to a specific deployed contract.
func NewDualAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*DualAggregatorTransactor, error) {
	contract, err := bindDualAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorTransactor{contract: contract}, nil
}

// NewDualAggregatorFilterer creates a new log filterer instance of DualAggregator, bound to a specific deployed contract.
func NewDualAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*DualAggregatorFilterer, error) {
	contract, err := bindDualAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorFilterer{contract: contract}, nil
}

// bindDualAggregator binds a generic wrapper to an already deployed contract.
func bindDualAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DualAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DualAggregator *DualAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DualAggregator.Contract.DualAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DualAggregator *DualAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.Contract.DualAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DualAggregator *DualAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DualAggregator.Contract.DualAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DualAggregator *DualAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DualAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DualAggregator *DualAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DualAggregator *DualAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DualAggregator.Contract.contract.Transact(opts, method, params...)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_DualAggregator *DualAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_DualAggregator *DualAggregatorSession) CheckEnabled() (bool, error) {
	return _DualAggregator.Contract.CheckEnabled(&_DualAggregator.CallOpts)
}

// CheckEnabled is a free data retrieval call binding the contract method 0xdc7f0124.
//
// Solidity: function checkEnabled() view returns(bool)
func (_DualAggregator *DualAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _DualAggregator.Contract.CheckEnabled(&_DualAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DualAggregator *DualAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DualAggregator *DualAggregatorSession) Decimals() (uint8, error) {
	return _DualAggregator.Contract.Decimals(&_DualAggregator.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DualAggregator *DualAggregatorCallerSession) Decimals() (uint8, error) {
	return _DualAggregator.Contract.Decimals(&_DualAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_DualAggregator *DualAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_DualAggregator *DualAggregatorSession) Description() (string, error) {
	return _DualAggregator.Contract.Description(&_DualAggregator.CallOpts)
}

// Description is a free data retrieval call binding the contract method 0x7284e416.
//
// Solidity: function description() view returns(string)
func (_DualAggregator *DualAggregatorCallerSession) Description() (string, error) {
	return _DualAggregator.Contract.Description(&_DualAggregator.CallOpts)
}

// GetAdaptiveOracle is a free data retrieval call binding the contract method 0x364efd83.
//
// Solidity: function getAdaptiveOracle() view returns(address)
func (_DualAggregator *DualAggregatorCaller) GetAdaptiveOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getAdaptiveOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAdaptiveOracle is a free data retrieval call binding the contract method 0x364efd83.
//
// Solidity: function getAdaptiveOracle() view returns(address)
func (_DualAggregator *DualAggregatorSession) GetAdaptiveOracle() (common.Address, error) {
	return _DualAggregator.Contract.GetAdaptiveOracle(&_DualAggregator.CallOpts)
}

// GetAdaptiveOracle is a free data retrieval call binding the contract method 0x364efd83.
//
// Solidity: function getAdaptiveOracle() view returns(address)
func (_DualAggregator *DualAggregatorCallerSession) GetAdaptiveOracle() (common.Address, error) {
	return _DualAggregator.Contract.GetAdaptiveOracle(&_DualAggregator.CallOpts)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_DualAggregator *DualAggregatorCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_DualAggregator *DualAggregatorSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetAnswer(&_DualAggregator.CallOpts, roundId)
}

// GetAnswer is a free data retrieval call binding the contract method 0xb5ab58dc.
//
// Solidity: function getAnswer(uint256 roundId) view returns(int256)
func (_DualAggregator *DualAggregatorCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetAnswer(&_DualAggregator.CallOpts, roundId)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_DualAggregator *DualAggregatorCaller) GetBilling(opts *bind.CallOpts) (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(struct {
		MaximumGasPriceGwei       uint32
		ReasonableGasPriceGwei    uint32
		ObservationPaymentGjuels  uint32
		TransmissionPaymentGjuels uint32
		AccountingGas             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPriceGwei = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ReasonableGasPriceGwei = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ObservationPaymentGjuels = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.TransmissionPaymentGjuels = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.AccountingGas = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_DualAggregator *DualAggregatorSession) GetBilling() (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	return _DualAggregator.Contract.GetBilling(&_DualAggregator.CallOpts)
}

// GetBilling is a free data retrieval call binding the contract method 0x29937268.
//
// Solidity: function getBilling() view returns(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_DualAggregator *DualAggregatorCallerSession) GetBilling() (struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}, error) {
	return _DualAggregator.Contract.GetBilling(&_DualAggregator.CallOpts)
}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_DualAggregator *DualAggregatorCaller) GetBillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getBillingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_DualAggregator *DualAggregatorSession) GetBillingAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetBillingAccessController(&_DualAggregator.CallOpts)
}

// GetBillingAccessController is a free data retrieval call binding the contract method 0xc4c92b37.
//
// Solidity: function getBillingAccessController() view returns(address)
func (_DualAggregator *DualAggregatorCallerSession) GetBillingAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetBillingAccessController(&_DualAggregator.CallOpts)
}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_DualAggregator *DualAggregatorCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_DualAggregator *DualAggregatorSession) GetLinkToken() (common.Address, error) {
	return _DualAggregator.Contract.GetLinkToken(&_DualAggregator.CallOpts)
}

// GetLinkToken is a free data retrieval call binding the contract method 0xe76d5168.
//
// Solidity: function getLinkToken() view returns(address linkToken)
func (_DualAggregator *DualAggregatorCallerSession) GetLinkToken() (common.Address, error) {
	return _DualAggregator.Contract.GetLinkToken(&_DualAggregator.CallOpts)
}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_DualAggregator *DualAggregatorCaller) GetRequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getRequesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_DualAggregator *DualAggregatorSession) GetRequesterAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetRequesterAccessController(&_DualAggregator.CallOpts)
}

// GetRequesterAccessController is a free data retrieval call binding the contract method 0xdaffc4b5.
//
// Solidity: function getRequesterAccessController() view returns(address)
func (_DualAggregator *DualAggregatorCallerSession) GetRequesterAccessController() (common.Address, error) {
	return _DualAggregator.Contract.GetRequesterAccessController(&_DualAggregator.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_DualAggregator *DualAggregatorCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getRoundData", roundId)

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_DualAggregator *DualAggregatorSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _DualAggregator.Contract.GetRoundData(&_DualAggregator.CallOpts, roundId)
}

// GetRoundData is a free data retrieval call binding the contract method 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_DualAggregator *DualAggregatorCallerSession) GetRoundData(roundId *big.Int) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _DualAggregator.Contract.GetRoundData(&_DualAggregator.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_DualAggregator *DualAggregatorCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_DualAggregator *DualAggregatorSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetTimestamp(&_DualAggregator.CallOpts, roundId)
}

// GetTimestamp is a free data retrieval call binding the contract method 0xb633620c.
//
// Solidity: function getTimestamp(uint256 roundId) view returns(uint256)
func (_DualAggregator *DualAggregatorCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _DualAggregator.Contract.GetTimestamp(&_DualAggregator.CallOpts, roundId)
}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_DualAggregator *DualAggregatorCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_DualAggregator *DualAggregatorSession) GetTransmitters() ([]common.Address, error) {
	return _DualAggregator.Contract.GetTransmitters(&_DualAggregator.CallOpts)
}

// GetTransmitters is a free data retrieval call binding the contract method 0x666cab8d.
//
// Solidity: function getTransmitters() view returns(address[])
func (_DualAggregator *DualAggregatorCallerSession) GetTransmitters() ([]common.Address, error) {
	return _DualAggregator.Contract.GetTransmitters(&_DualAggregator.CallOpts)
}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_DualAggregator *DualAggregatorCaller) GetValidatorConfig(opts *bind.CallOpts) (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "getValidatorConfig")

	outstruct := new(struct {
		Validator common.Address
		GasLimit  uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.GasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_DualAggregator *DualAggregatorSession) GetValidatorConfig() (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	return _DualAggregator.Contract.GetValidatorConfig(&_DualAggregator.CallOpts)
}

// GetValidatorConfig is a free data retrieval call binding the contract method 0x9bd2c0b1.
//
// Solidity: function getValidatorConfig() view returns(address validator, uint32 gasLimit)
func (_DualAggregator *DualAggregatorCallerSession) GetValidatorConfig() (struct {
	Validator common.Address
	GasLimit  uint32
}, error) {
	return _DualAggregator.Contract.GetValidatorConfig(&_DualAggregator.CallOpts)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_DualAggregator *DualAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_DualAggregator *DualAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _DualAggregator.Contract.HasAccess(&_DualAggregator.CallOpts, _user, _calldata)
}

// HasAccess is a free data retrieval call binding the contract method 0x6b14daf8.
//
// Solidity: function hasAccess(address _user, bytes _calldata) view returns(bool)
func (_DualAggregator *DualAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _DualAggregator.Contract.HasAccess(&_DualAggregator.CallOpts, _user, _calldata)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.LatestAnswer(&_DualAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.LatestAnswer(&_DualAggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_DualAggregator *DualAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(struct {
		ConfigCount  uint32
		BlockNumber  uint32
		ConfigDigest [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_DualAggregator *DualAggregatorSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _DualAggregator.Contract.LatestConfigDetails(&_DualAggregator.CallOpts)
}

// LatestConfigDetails is a free data retrieval call binding the contract method 0x81ff7048.
//
// Solidity: function latestConfigDetails() view returns(uint32 configCount, uint32 blockNumber, bytes32 configDigest)
func (_DualAggregator *DualAggregatorCallerSession) LatestConfigDetails() (struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}, error) {
	return _DualAggregator.Contract.LatestConfigDetails(&_DualAggregator.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_DualAggregator *DualAggregatorCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(struct {
		ScanLogs     bool
		ConfigDigest [32]byte
		Epoch        uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_DualAggregator *DualAggregatorSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _DualAggregator.Contract.LatestConfigDigestAndEpoch(&_DualAggregator.CallOpts)
}

// LatestConfigDigestAndEpoch is a free data retrieval call binding the contract method 0xafcb95d7.
//
// Solidity: function latestConfigDigestAndEpoch() view returns(bool scanLogs, bytes32 configDigest, uint32 epoch)
func (_DualAggregator *DualAggregatorCallerSession) LatestConfigDigestAndEpoch() (struct {
	ScanLogs     bool
	ConfigDigest [32]byte
	Epoch        uint32
}, error) {
	return _DualAggregator.Contract.LatestConfigDigestAndEpoch(&_DualAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_DualAggregator *DualAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_DualAggregator *DualAggregatorSession) LatestRound() (*big.Int, error) {
	return _DualAggregator.Contract.LatestRound(&_DualAggregator.CallOpts)
}

// LatestRound is a free data retrieval call binding the contract method 0x668a0f02.
//
// Solidity: function latestRound() view returns(uint256)
func (_DualAggregator *DualAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _DualAggregator.Contract.LatestRound(&_DualAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_DualAggregator *DualAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(struct {
		RoundId         *big.Int
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RoundId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_DualAggregator *DualAggregatorSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _DualAggregator.Contract.LatestRoundData(&_DualAggregator.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (_DualAggregator *DualAggregatorCallerSession) LatestRoundData() (struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}, error) {
	return _DualAggregator.Contract.LatestRoundData(&_DualAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_DualAggregator *DualAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_DualAggregator *DualAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _DualAggregator.Contract.LatestTimestamp(&_DualAggregator.CallOpts)
}

// LatestTimestamp is a free data retrieval call binding the contract method 0x8205bf6a.
//
// Solidity: function latestTimestamp() view returns(uint256)
func (_DualAggregator *DualAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _DualAggregator.Contract.LatestTimestamp(&_DualAggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_DualAggregator *DualAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(struct {
		ConfigDigest    [32]byte
		Epoch           uint32
		Round           uint8
		LatestAnswer    *big.Int
		LatestTimestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Round = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.LatestAnswer = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.LatestTimestamp = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_DualAggregator *DualAggregatorSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _DualAggregator.Contract.LatestTransmissionDetails(&_DualAggregator.CallOpts)
}

// LatestTransmissionDetails is a free data retrieval call binding the contract method 0xe5fe4577.
//
// Solidity: function latestTransmissionDetails() view returns(bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
func (_DualAggregator *DualAggregatorCallerSession) LatestTransmissionDetails() (struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}, error) {
	return _DualAggregator.Contract.LatestTransmissionDetails(&_DualAggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_DualAggregator *DualAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_DualAggregator *DualAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _DualAggregator.Contract.LinkAvailableForPayment(&_DualAggregator.CallOpts)
}

// LinkAvailableForPayment is a free data retrieval call binding the contract method 0xd09dc339.
//
// Solidity: function linkAvailableForPayment() view returns(int256 availableBalance)
func (_DualAggregator *DualAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _DualAggregator.Contract.LinkAvailableForPayment(&_DualAggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MaxAnswer(&_DualAggregator.CallOpts)
}

// MaxAnswer is a free data retrieval call binding the contract method 0x70da2f67.
//
// Solidity: function maxAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MaxAnswer(&_DualAggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorSession) MinAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MinAnswer(&_DualAggregator.CallOpts)
}

// MinAnswer is a free data retrieval call binding the contract method 0x22adbc78.
//
// Solidity: function minAnswer() view returns(int256)
func (_DualAggregator *DualAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _DualAggregator.Contract.MinAnswer(&_DualAggregator.CallOpts)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_DualAggregator *DualAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "oracleObservationCount", transmitterAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_DualAggregator *DualAggregatorSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _DualAggregator.Contract.OracleObservationCount(&_DualAggregator.CallOpts, transmitterAddress)
}

// OracleObservationCount is a free data retrieval call binding the contract method 0xe4902f82.
//
// Solidity: function oracleObservationCount(address transmitterAddress) view returns(uint32)
func (_DualAggregator *DualAggregatorCallerSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _DualAggregator.Contract.OracleObservationCount(&_DualAggregator.CallOpts, transmitterAddress)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_DualAggregator *DualAggregatorCaller) OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "owedPayment", transmitterAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_DualAggregator *DualAggregatorSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _DualAggregator.Contract.OwedPayment(&_DualAggregator.CallOpts, transmitterAddress)
}

// OwedPayment is a free data retrieval call binding the contract method 0x0eafb25b.
//
// Solidity: function owedPayment(address transmitterAddress) view returns(uint256)
func (_DualAggregator *DualAggregatorCallerSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _DualAggregator.Contract.OwedPayment(&_DualAggregator.CallOpts, transmitterAddress)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DualAggregator *DualAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DualAggregator *DualAggregatorSession) Owner() (common.Address, error) {
	return _DualAggregator.Contract.Owner(&_DualAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DualAggregator *DualAggregatorCallerSession) Owner() (common.Address, error) {
	return _DualAggregator.Contract.Owner(&_DualAggregator.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_DualAggregator *DualAggregatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_DualAggregator *DualAggregatorSession) TypeAndVersion() (string, error) {
	return _DualAggregator.Contract.TypeAndVersion(&_DualAggregator.CallOpts)
}

// TypeAndVersion is a free data retrieval call binding the contract method 0x181f5a77.
//
// Solidity: function typeAndVersion() view returns(string)
func (_DualAggregator *DualAggregatorCallerSession) TypeAndVersion() (string, error) {
	return _DualAggregator.Contract.TypeAndVersion(&_DualAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_DualAggregator *DualAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DualAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_DualAggregator *DualAggregatorSession) Version() (*big.Int, error) {
	return _DualAggregator.Contract.Version(&_DualAggregator.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_DualAggregator *DualAggregatorCallerSession) Version() (*big.Int, error) {
	return _DualAggregator.Contract.Version(&_DualAggregator.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DualAggregator *DualAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DualAggregator *DualAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptOwnership(&_DualAggregator.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DualAggregator *DualAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptOwnership(&_DualAggregator.TransactOpts)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_DualAggregator *DualAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "acceptPayeeship", transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_DualAggregator *DualAggregatorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptPayeeship(&_DualAggregator.TransactOpts, transmitter)
}

// AcceptPayeeship is a paid mutator transaction binding the contract method 0xb121e147.
//
// Solidity: function acceptPayeeship(address transmitter) returns()
func (_DualAggregator *DualAggregatorTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AcceptPayeeship(&_DualAggregator.TransactOpts, transmitter)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_DualAggregator *DualAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "addAccess", _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_DualAggregator *DualAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AddAccess(&_DualAggregator.TransactOpts, _user)
}

// AddAccess is a paid mutator transaction binding the contract method 0xa118f249.
//
// Solidity: function addAccess(address _user) returns()
func (_DualAggregator *DualAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.AddAccess(&_DualAggregator.TransactOpts, _user)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_DualAggregator *DualAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "disableAccessCheck")
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_DualAggregator *DualAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.DisableAccessCheck(&_DualAggregator.TransactOpts)
}

// DisableAccessCheck is a paid mutator transaction binding the contract method 0x0a756983.
//
// Solidity: function disableAccessCheck() returns()
func (_DualAggregator *DualAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.DisableAccessCheck(&_DualAggregator.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_DualAggregator *DualAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "enableAccessCheck")
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_DualAggregator *DualAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.EnableAccessCheck(&_DualAggregator.TransactOpts)
}

// EnableAccessCheck is a paid mutator transaction binding the contract method 0x8038e4a1.
//
// Solidity: function enableAccessCheck() returns()
func (_DualAggregator *DualAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _DualAggregator.Contract.EnableAccessCheck(&_DualAggregator.TransactOpts)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_DualAggregator *DualAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "removeAccess", _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_DualAggregator *DualAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.RemoveAccess(&_DualAggregator.TransactOpts, _user)
}

// RemoveAccess is a paid mutator transaction binding the contract method 0x8823da6c.
//
// Solidity: function removeAccess(address _user) returns()
func (_DualAggregator *DualAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.RemoveAccess(&_DualAggregator.TransactOpts, _user)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_DualAggregator *DualAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "requestNewRound")
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_DualAggregator *DualAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _DualAggregator.Contract.RequestNewRound(&_DualAggregator.TransactOpts)
}

// RequestNewRound is a paid mutator transaction binding the contract method 0x98e5b12a.
//
// Solidity: function requestNewRound() returns(uint80)
func (_DualAggregator *DualAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _DualAggregator.Contract.RequestNewRound(&_DualAggregator.TransactOpts)
}

// ResetToReferenceRate is a paid mutator transaction binding the contract method 0x94526048.
//
// Solidity: function resetToReferenceRate(int256 referenceRate) returns()
func (_DualAggregator *DualAggregatorTransactor) ResetToReferenceRate(opts *bind.TransactOpts, referenceRate *big.Int) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "resetToReferenceRate", referenceRate)
}

// ResetToReferenceRate is a paid mutator transaction binding the contract method 0x94526048.
//
// Solidity: function resetToReferenceRate(int256 referenceRate) returns()
func (_DualAggregator *DualAggregatorSession) ResetToReferenceRate(referenceRate *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.ResetToReferenceRate(&_DualAggregator.TransactOpts, referenceRate)
}

// ResetToReferenceRate is a paid mutator transaction binding the contract method 0x94526048.
//
// Solidity: function resetToReferenceRate(int256 referenceRate) returns()
func (_DualAggregator *DualAggregatorTransactorSession) ResetToReferenceRate(referenceRate *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.ResetToReferenceRate(&_DualAggregator.TransactOpts, referenceRate)
}

// SetAdaptiveOracle is a paid mutator transaction binding the contract method 0x98c89ea2.
//
// Solidity: function setAdaptiveOracle(address adaptiveOracle) returns()
func (_DualAggregator *DualAggregatorTransactor) SetAdaptiveOracle(opts *bind.TransactOpts, adaptiveOracle common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setAdaptiveOracle", adaptiveOracle)
}

// SetAdaptiveOracle is a paid mutator transaction binding the contract method 0x98c89ea2.
//
// Solidity: function setAdaptiveOracle(address adaptiveOracle) returns()
func (_DualAggregator *DualAggregatorSession) SetAdaptiveOracle(adaptiveOracle common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetAdaptiveOracle(&_DualAggregator.TransactOpts, adaptiveOracle)
}

// SetAdaptiveOracle is a paid mutator transaction binding the contract method 0x98c89ea2.
//
// Solidity: function setAdaptiveOracle(address adaptiveOracle) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetAdaptiveOracle(adaptiveOracle common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetAdaptiveOracle(&_DualAggregator.TransactOpts, adaptiveOracle)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_DualAggregator *DualAggregatorTransactor) SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setBilling", maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_DualAggregator *DualAggregatorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBilling(&_DualAggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBilling is a paid mutator transaction binding the contract method 0x643dc105.
//
// Solidity: function setBilling(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBilling(&_DualAggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_DualAggregator *DualAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_DualAggregator *DualAggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBillingAccessController(&_DualAggregator.TransactOpts, _billingAccessController)
}

// SetBillingAccessController is a paid mutator transaction binding the contract method 0xfbffd2c1.
//
// Solidity: function setBillingAccessController(address _billingAccessController) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetBillingAccessController(&_DualAggregator.TransactOpts, _billingAccessController)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_DualAggregator *DualAggregatorTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_DualAggregator *DualAggregatorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetConfig(&_DualAggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetConfig is a paid mutator transaction binding the contract method 0xe3d0e712.
//
// Solidity: function setConfig(address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetConfig(&_DualAggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

// SetCutoffTime is a paid mutator transaction binding the contract method 0xb17f2a6b.
//
// Solidity: function setCutoffTime(uint32 _cutoffTime) returns()
func (_DualAggregator *DualAggregatorTransactor) SetCutoffTime(opts *bind.TransactOpts, _cutoffTime uint32) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setCutoffTime", _cutoffTime)
}

// SetCutoffTime is a paid mutator transaction binding the contract method 0xb17f2a6b.
//
// Solidity: function setCutoffTime(uint32 _cutoffTime) returns()
func (_DualAggregator *DualAggregatorSession) SetCutoffTime(_cutoffTime uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetCutoffTime(&_DualAggregator.TransactOpts, _cutoffTime)
}

// SetCutoffTime is a paid mutator transaction binding the contract method 0xb17f2a6b.
//
// Solidity: function setCutoffTime(uint32 _cutoffTime) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetCutoffTime(_cutoffTime uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetCutoffTime(&_DualAggregator.TransactOpts, _cutoffTime)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_DualAggregator *DualAggregatorTransactor) SetLinkToken(opts *bind.TransactOpts, linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setLinkToken", linkToken, recipient)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_DualAggregator *DualAggregatorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetLinkToken(&_DualAggregator.TransactOpts, linkToken, recipient)
}

// SetLinkToken is a paid mutator transaction binding the contract method 0x4fb17470.
//
// Solidity: function setLinkToken(address linkToken, address recipient) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetLinkToken(&_DualAggregator.TransactOpts, linkToken, recipient)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_DualAggregator *DualAggregatorTransactor) SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setPayees", transmitters, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_DualAggregator *DualAggregatorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetPayees(&_DualAggregator.TransactOpts, transmitters, payees)
}

// SetPayees is a paid mutator transaction binding the contract method 0x9c849b30.
//
// Solidity: function setPayees(address[] transmitters, address[] payees) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetPayees(&_DualAggregator.TransactOpts, transmitters, payees)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_DualAggregator *DualAggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, requesterAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setRequesterAccessController", requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_DualAggregator *DualAggregatorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetRequesterAccessController(&_DualAggregator.TransactOpts, requesterAccessController)
}

// SetRequesterAccessController is a paid mutator transaction binding the contract method 0x9e3ceeab.
//
// Solidity: function setRequesterAccessController(address requesterAccessController) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetRequesterAccessController(&_DualAggregator.TransactOpts, requesterAccessController)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_DualAggregator *DualAggregatorTransactor) SetValidatorConfig(opts *bind.TransactOpts, newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "setValidatorConfig", newValidator, newGasLimit)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_DualAggregator *DualAggregatorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetValidatorConfig(&_DualAggregator.TransactOpts, newValidator, newGasLimit)
}

// SetValidatorConfig is a paid mutator transaction binding the contract method 0xeb457163.
//
// Solidity: function setValidatorConfig(address newValidator, uint32 newGasLimit) returns()
func (_DualAggregator *DualAggregatorTransactorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _DualAggregator.Contract.SetValidatorConfig(&_DualAggregator.TransactOpts, newValidator, newGasLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_DualAggregator *DualAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transferOwnership", to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_DualAggregator *DualAggregatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferOwnership(&_DualAggregator.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to) returns()
func (_DualAggregator *DualAggregatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferOwnership(&_DualAggregator.TransactOpts, to)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_DualAggregator *DualAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_DualAggregator *DualAggregatorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferPayeeship(&_DualAggregator.TransactOpts, transmitter, proposed)
}

// TransferPayeeship is a paid mutator transaction binding the contract method 0xeb5dcd6c.
//
// Solidity: function transferPayeeship(address transmitter, address proposed) returns()
func (_DualAggregator *DualAggregatorTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransferPayeeship(&_DualAggregator.TransactOpts, transmitter, proposed)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_DualAggregator *DualAggregatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_DualAggregator *DualAggregatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.Transmit(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xb1dc65a4.
//
// Solidity: function transmit(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_DualAggregator *DualAggregatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.Transmit(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// TransmitSecondary is a paid mutator transaction binding the contract method 0xba0cb29e.
//
// Solidity: function transmitSecondary(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_DualAggregator *DualAggregatorTransactor) TransmitSecondary(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "transmitSecondary", reportContext, report, rs, ss, rawVs)
}

// TransmitSecondary is a paid mutator transaction binding the contract method 0xba0cb29e.
//
// Solidity: function transmitSecondary(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_DualAggregator *DualAggregatorSession) TransmitSecondary(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransmitSecondary(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// TransmitSecondary is a paid mutator transaction binding the contract method 0xba0cb29e.
//
// Solidity: function transmitSecondary(bytes32[3] reportContext, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_DualAggregator *DualAggregatorTransactorSession) TransmitSecondary(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _DualAggregator.Contract.TransmitSecondary(&_DualAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_DualAggregator *DualAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "withdrawFunds", recipient, amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_DualAggregator *DualAggregatorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawFunds(&_DualAggregator.TransactOpts, recipient, amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0xc1075329.
//
// Solidity: function withdrawFunds(address recipient, uint256 amount) returns()
func (_DualAggregator *DualAggregatorTransactorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawFunds(&_DualAggregator.TransactOpts, recipient, amount)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_DualAggregator *DualAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.contract.Transact(opts, "withdrawPayment", transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_DualAggregator *DualAggregatorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawPayment(&_DualAggregator.TransactOpts, transmitter)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address transmitter) returns()
func (_DualAggregator *DualAggregatorTransactorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _DualAggregator.Contract.WithdrawPayment(&_DualAggregator.TransactOpts, transmitter)
}

// DualAggregatorAdaptiveOracleSetIterator is returned from FilterAdaptiveOracleSet and is used to iterate over the raw logs and unpacked data for AdaptiveOracleSet events raised by the DualAggregator contract.
type DualAggregatorAdaptiveOracleSetIterator struct {
	Event *DualAggregatorAdaptiveOracleSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorAdaptiveOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorAdaptiveOracleSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorAdaptiveOracleSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorAdaptiveOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorAdaptiveOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorAdaptiveOracleSet represents a AdaptiveOracleSet event raised by the DualAggregator contract.
type DualAggregatorAdaptiveOracleSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAdaptiveOracleSet is a free log retrieval operation binding the contract event 0x1969b480ede741c56a9a724f5457349d09b3ee4574a1dbb1b75a2e668b680504.
//
// Solidity: event AdaptiveOracleSet(address indexed old, address indexed current)
func (_DualAggregator *DualAggregatorFilterer) FilterAdaptiveOracleSet(opts *bind.FilterOpts, old []common.Address, current []common.Address) (*DualAggregatorAdaptiveOracleSetIterator, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "AdaptiveOracleSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorAdaptiveOracleSetIterator{contract: _DualAggregator.contract, event: "AdaptiveOracleSet", logs: logs, sub: sub}, nil
}

// WatchAdaptiveOracleSet is a free log subscription operation binding the contract event 0x1969b480ede741c56a9a724f5457349d09b3ee4574a1dbb1b75a2e668b680504.
//
// Solidity: event AdaptiveOracleSet(address indexed old, address indexed current)
func (_DualAggregator *DualAggregatorFilterer) WatchAdaptiveOracleSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorAdaptiveOracleSet, old []common.Address, current []common.Address) (event.Subscription, error) {

	var oldRule []interface{}
	for _, oldItem := range old {
		oldRule = append(oldRule, oldItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "AdaptiveOracleSet", oldRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorAdaptiveOracleSet)
				if err := _DualAggregator.contract.UnpackLog(event, "AdaptiveOracleSet", log); err != nil {
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

// ParseAdaptiveOracleSet is a log parse operation binding the contract event 0x1969b480ede741c56a9a724f5457349d09b3ee4574a1dbb1b75a2e668b680504.
//
// Solidity: event AdaptiveOracleSet(address indexed old, address indexed current)
func (_DualAggregator *DualAggregatorFilterer) ParseAdaptiveOracleSet(log types.Log) (*DualAggregatorAdaptiveOracleSet, error) {
	event := new(DualAggregatorAdaptiveOracleSet)
	if err := _DualAggregator.contract.UnpackLog(event, "AdaptiveOracleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorAdaptiveRateResetIterator is returned from FilterAdaptiveRateReset and is used to iterate over the raw logs and unpacked data for AdaptiveRateReset events raised by the DualAggregator contract.
type DualAggregatorAdaptiveRateResetIterator struct {
	Event *DualAggregatorAdaptiveRateReset // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorAdaptiveRateResetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorAdaptiveRateReset)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorAdaptiveRateReset)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorAdaptiveRateResetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorAdaptiveRateResetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorAdaptiveRateReset represents a AdaptiveRateReset event raised by the DualAggregator contract.
type DualAggregatorAdaptiveRateReset struct {
	RoundId       uint32
	ReferenceRate *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdaptiveRateReset is a free log retrieval operation binding the contract event 0x6f643b21ecba6d2a464beb4b6fa87271920464cfdea7d6ff1b75dd3c7f143747.
//
// Solidity: event AdaptiveRateReset(uint32 indexed roundId, int192 referenceRate)
func (_DualAggregator *DualAggregatorFilterer) FilterAdaptiveRateReset(opts *bind.FilterOpts, roundId []uint32) (*DualAggregatorAdaptiveRateResetIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "AdaptiveRateReset", roundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorAdaptiveRateResetIterator{contract: _DualAggregator.contract, event: "AdaptiveRateReset", logs: logs, sub: sub}, nil
}

// WatchAdaptiveRateReset is a free log subscription operation binding the contract event 0x6f643b21ecba6d2a464beb4b6fa87271920464cfdea7d6ff1b75dd3c7f143747.
//
// Solidity: event AdaptiveRateReset(uint32 indexed roundId, int192 referenceRate)
func (_DualAggregator *DualAggregatorFilterer) WatchAdaptiveRateReset(opts *bind.WatchOpts, sink chan<- *DualAggregatorAdaptiveRateReset, roundId []uint32) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "AdaptiveRateReset", roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorAdaptiveRateReset)
				if err := _DualAggregator.contract.UnpackLog(event, "AdaptiveRateReset", log); err != nil {
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

// ParseAdaptiveRateReset is a log parse operation binding the contract event 0x6f643b21ecba6d2a464beb4b6fa87271920464cfdea7d6ff1b75dd3c7f143747.
//
// Solidity: event AdaptiveRateReset(uint32 indexed roundId, int192 referenceRate)
func (_DualAggregator *DualAggregatorFilterer) ParseAdaptiveRateReset(log types.Log) (*DualAggregatorAdaptiveRateReset, error) {
	event := new(DualAggregatorAdaptiveRateReset)
	if err := _DualAggregator.contract.UnpackLog(event, "AdaptiveRateReset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorAddedAccessIterator is returned from FilterAddedAccess and is used to iterate over the raw logs and unpacked data for AddedAccess events raised by the DualAggregator contract.
type DualAggregatorAddedAccessIterator struct {
	Event *DualAggregatorAddedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorAddedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorAddedAccess)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorAddedAccess)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorAddedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorAddedAccess represents a AddedAccess event raised by the DualAggregator contract.
type DualAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddedAccess is a free log retrieval operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_DualAggregator *DualAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*DualAggregatorAddedAccessIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorAddedAccessIterator{contract: _DualAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

// WatchAddedAccess is a free log subscription operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_DualAggregator *DualAggregatorFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *DualAggregatorAddedAccess) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorAddedAccess)
				if err := _DualAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

// ParseAddedAccess is a log parse operation binding the contract event 0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4.
//
// Solidity: event AddedAccess(address user)
func (_DualAggregator *DualAggregatorFilterer) ParseAddedAccess(log types.Log) (*DualAggregatorAddedAccess, error) {
	event := new(DualAggregatorAddedAccess)
	if err := _DualAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorAnswerUpdatedIterator is returned from FilterAnswerUpdated and is used to iterate over the raw logs and unpacked data for AnswerUpdated events raised by the DualAggregator contract.
type DualAggregatorAnswerUpdatedIterator struct {
	Event *DualAggregatorAnswerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorAnswerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorAnswerUpdated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorAnswerUpdated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorAnswerUpdated represents a AnswerUpdated event raised by the DualAggregator contract.
type DualAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAnswerUpdated is a free log retrieval operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_DualAggregator *DualAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*DualAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorAnswerUpdatedIterator{contract: _DualAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnswerUpdated is a free log subscription operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_DualAggregator *DualAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *DualAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorAnswerUpdated)
				if err := _DualAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

// ParseAnswerUpdated is a log parse operation binding the contract event 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (_DualAggregator *DualAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*DualAggregatorAnswerUpdated, error) {
	event := new(DualAggregatorAnswerUpdated)
	if err := _DualAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorBillingAccessControllerSetIterator is returned from FilterBillingAccessControllerSet and is used to iterate over the raw logs and unpacked data for BillingAccessControllerSet events raised by the DualAggregator contract.
type DualAggregatorBillingAccessControllerSetIterator struct {
	Event *DualAggregatorBillingAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorBillingAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorBillingAccessControllerSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorBillingAccessControllerSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorBillingAccessControllerSet represents a BillingAccessControllerSet event raised by the DualAggregator contract.
type DualAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBillingAccessControllerSet is a free log retrieval operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_DualAggregator *DualAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*DualAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorBillingAccessControllerSetIterator{contract: _DualAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchBillingAccessControllerSet is a free log subscription operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_DualAggregator *DualAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorBillingAccessControllerSet)
				if err := _DualAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

// ParseBillingAccessControllerSet is a log parse operation binding the contract event 0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912.
//
// Solidity: event BillingAccessControllerSet(address old, address current)
func (_DualAggregator *DualAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*DualAggregatorBillingAccessControllerSet, error) {
	event := new(DualAggregatorBillingAccessControllerSet)
	if err := _DualAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorBillingSetIterator is returned from FilterBillingSet and is used to iterate over the raw logs and unpacked data for BillingSet events raised by the DualAggregator contract.
type DualAggregatorBillingSetIterator struct {
	Event *DualAggregatorBillingSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorBillingSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorBillingSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorBillingSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorBillingSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorBillingSet represents a BillingSet event raised by the DualAggregator contract.
type DualAggregatorBillingSet struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterBillingSet is a free log retrieval operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_DualAggregator *DualAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*DualAggregatorBillingSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorBillingSetIterator{contract: _DualAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

// WatchBillingSet is a free log subscription operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_DualAggregator *DualAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorBillingSet)
				if err := _DualAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

// ParseBillingSet is a log parse operation binding the contract event 0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f.
//
// Solidity: event BillingSet(uint32 maximumGasPriceGwei, uint32 reasonableGasPriceGwei, uint32 observationPaymentGjuels, uint32 transmissionPaymentGjuels, uint24 accountingGas)
func (_DualAggregator *DualAggregatorFilterer) ParseBillingSet(log types.Log) (*DualAggregatorBillingSet, error) {
	event := new(DualAggregatorBillingSet)
	if err := _DualAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorCheckAccessDisabledIterator is returned from FilterCheckAccessDisabled and is used to iterate over the raw logs and unpacked data for CheckAccessDisabled events raised by the DualAggregator contract.
type DualAggregatorCheckAccessDisabledIterator struct {
	Event *DualAggregatorCheckAccessDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorCheckAccessDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorCheckAccessDisabled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorCheckAccessDisabled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorCheckAccessDisabled represents a CheckAccessDisabled event raised by the DualAggregator contract.
type DualAggregatorCheckAccessDisabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessDisabled is a free log retrieval operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_DualAggregator *DualAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*DualAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCheckAccessDisabledIterator{contract: _DualAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessDisabled is a free log subscription operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_DualAggregator *DualAggregatorFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *DualAggregatorCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorCheckAccessDisabled)
				if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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

// ParseCheckAccessDisabled is a log parse operation binding the contract event 0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638.
//
// Solidity: event CheckAccessDisabled()
func (_DualAggregator *DualAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*DualAggregatorCheckAccessDisabled, error) {
	event := new(DualAggregatorCheckAccessDisabled)
	if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorCheckAccessEnabledIterator is returned from FilterCheckAccessEnabled and is used to iterate over the raw logs and unpacked data for CheckAccessEnabled events raised by the DualAggregator contract.
type DualAggregatorCheckAccessEnabledIterator struct {
	Event *DualAggregatorCheckAccessEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorCheckAccessEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorCheckAccessEnabled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorCheckAccessEnabled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorCheckAccessEnabled represents a CheckAccessEnabled event raised by the DualAggregator contract.
type DualAggregatorCheckAccessEnabled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCheckAccessEnabled is a free log retrieval operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_DualAggregator *DualAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*DualAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCheckAccessEnabledIterator{contract: _DualAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

// WatchCheckAccessEnabled is a free log subscription operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_DualAggregator *DualAggregatorFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *DualAggregatorCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorCheckAccessEnabled)
				if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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

// ParseCheckAccessEnabled is a log parse operation binding the contract event 0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480.
//
// Solidity: event CheckAccessEnabled()
func (_DualAggregator *DualAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*DualAggregatorCheckAccessEnabled, error) {
	event := new(DualAggregatorCheckAccessEnabled)
	if err := _DualAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorConfigSetIterator is returned from FilterConfigSet and is used to iterate over the raw logs and unpacked data for ConfigSet events raised by the DualAggregator contract.
type DualAggregatorConfigSetIterator struct {
	Event *DualAggregatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorConfigSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorConfigSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorConfigSet represents a ConfigSet event raised by the DualAggregator contract.
type DualAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterConfigSet is a free log retrieval operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_DualAggregator *DualAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*DualAggregatorConfigSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorConfigSetIterator{contract: _DualAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

// WatchConfigSet is a free log subscription operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_DualAggregator *DualAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorConfigSet)
				if err := _DualAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

// ParseConfigSet is a log parse operation binding the contract event 0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05.
//
// Solidity: event ConfigSet(uint32 previousConfigBlockNumber, bytes32 configDigest, uint64 configCount, address[] signers, address[] transmitters, uint8 f, bytes onchainConfig, uint64 offchainConfigVersion, bytes offchainConfig)
func (_DualAggregator *DualAggregatorFilterer) ParseConfigSet(log types.Log) (*DualAggregatorConfigSet, error) {
	event := new(DualAggregatorConfigSet)
	if err := _DualAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorCutoffTimeSetIterator is returned from FilterCutoffTimeSet and is used to iterate over the raw logs and unpacked data for CutoffTimeSet events raised by the DualAggregator contract.
type DualAggregatorCutoffTimeSetIterator struct {
	Event *DualAggregatorCutoffTimeSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorCutoffTimeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorCutoffTimeSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorCutoffTimeSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorCutoffTimeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorCutoffTimeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorCutoffTimeSet represents a CutoffTimeSet event raised by the DualAggregator contract.
type DualAggregatorCutoffTimeSet struct {
	CutoffTime uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCutoffTimeSet is a free log retrieval operation binding the contract event 0xb24a681ce3399a408a89fd0c2b59dfc24bdad592b1c7ec7671cf060596c1c4d1.
//
// Solidity: event CutoffTimeSet(uint32 cutoffTime)
func (_DualAggregator *DualAggregatorFilterer) FilterCutoffTimeSet(opts *bind.FilterOpts) (*DualAggregatorCutoffTimeSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "CutoffTimeSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorCutoffTimeSetIterator{contract: _DualAggregator.contract, event: "CutoffTimeSet", logs: logs, sub: sub}, nil
}

// WatchCutoffTimeSet is a free log subscription operation binding the contract event 0xb24a681ce3399a408a89fd0c2b59dfc24bdad592b1c7ec7671cf060596c1c4d1.
//
// Solidity: event CutoffTimeSet(uint32 cutoffTime)
func (_DualAggregator *DualAggregatorFilterer) WatchCutoffTimeSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorCutoffTimeSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "CutoffTimeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorCutoffTimeSet)
				if err := _DualAggregator.contract.UnpackLog(event, "CutoffTimeSet", log); err != nil {
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

// ParseCutoffTimeSet is a log parse operation binding the contract event 0xb24a681ce3399a408a89fd0c2b59dfc24bdad592b1c7ec7671cf060596c1c4d1.
//
// Solidity: event CutoffTimeSet(uint32 cutoffTime)
func (_DualAggregator *DualAggregatorFilterer) ParseCutoffTimeSet(log types.Log) (*DualAggregatorCutoffTimeSet, error) {
	event := new(DualAggregatorCutoffTimeSet)
	if err := _DualAggregator.contract.UnpackLog(event, "CutoffTimeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorLinkTokenSetIterator is returned from FilterLinkTokenSet and is used to iterate over the raw logs and unpacked data for LinkTokenSet events raised by the DualAggregator contract.
type DualAggregatorLinkTokenSetIterator struct {
	Event *DualAggregatorLinkTokenSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorLinkTokenSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorLinkTokenSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorLinkTokenSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorLinkTokenSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorLinkTokenSet represents a LinkTokenSet event raised by the DualAggregator contract.
type DualAggregatorLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLinkTokenSet is a free log retrieval operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_DualAggregator *DualAggregatorFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, oldLinkToken []common.Address, newLinkToken []common.Address) (*DualAggregatorLinkTokenSetIterator, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorLinkTokenSetIterator{contract: _DualAggregator.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

// WatchLinkTokenSet is a free log subscription operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_DualAggregator *DualAggregatorFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorLinkTokenSet, oldLinkToken []common.Address, newLinkToken []common.Address) (event.Subscription, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorLinkTokenSet)
				if err := _DualAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
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

// ParseLinkTokenSet is a log parse operation binding the contract event 0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a.
//
// Solidity: event LinkTokenSet(address indexed oldLinkToken, address indexed newLinkToken)
func (_DualAggregator *DualAggregatorFilterer) ParseLinkTokenSet(log types.Log) (*DualAggregatorLinkTokenSet, error) {
	event := new(DualAggregatorLinkTokenSet)
	if err := _DualAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the DualAggregator contract.
type DualAggregatorNewRoundIterator struct {
	Event *DualAggregatorNewRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorNewRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorNewRound)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorNewRound)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorNewRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorNewRound represents a NewRound event raised by the DualAggregator contract.
type DualAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewRound is a free log retrieval operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_DualAggregator *DualAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*DualAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorNewRoundIterator{contract: _DualAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

// WatchNewRound is a free log subscription operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_DualAggregator *DualAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *DualAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorNewRound)
				if err := _DualAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

// ParseNewRound is a log parse operation binding the contract event 0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271.
//
// Solidity: event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)
func (_DualAggregator *DualAggregatorFilterer) ParseNewRound(log types.Log) (*DualAggregatorNewRound, error) {
	event := new(DualAggregatorNewRound)
	if err := _DualAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorNewTransmissionIterator is returned from FilterNewTransmission and is used to iterate over the raw logs and unpacked data for NewTransmission events raised by the DualAggregator contract.
type DualAggregatorNewTransmissionIterator struct {
	Event *DualAggregatorNewTransmission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorNewTransmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorNewTransmission)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorNewTransmission)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorNewTransmission represents a NewTransmission event raised by the DualAggregator contract.
type DualAggregatorNewTransmission struct {
	AggregatorRoundId     uint32
	Answer                *big.Int
	Transmitter           common.Address
	ObservationsTimestamp uint32
	Observations          []*big.Int
	Observers             []byte
	JuelsPerFeeCoin       *big.Int
	ConfigDigest          [32]byte
	EpochAndRound         *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterNewTransmission is a free log retrieval operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_DualAggregator *DualAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*DualAggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorNewTransmissionIterator{contract: _DualAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

// WatchNewTransmission is a free log subscription operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_DualAggregator *DualAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *DualAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorNewTransmission)
				if err := _DualAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

// ParseNewTransmission is a log parse operation binding the contract event 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a.
//
// Solidity: event NewTransmission(uint32 indexed aggregatorRoundId, int192 answer, address transmitter, uint32 observationsTimestamp, int192[] observations, bytes observers, int192 juelsPerFeeCoin, bytes32 configDigest, uint40 epochAndRound)
func (_DualAggregator *DualAggregatorFilterer) ParseNewTransmission(log types.Log) (*DualAggregatorNewTransmission, error) {
	event := new(DualAggregatorNewTransmission)
	if err := _DualAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorOraclePaidIterator is returned from FilterOraclePaid and is used to iterate over the raw logs and unpacked data for OraclePaid events raised by the DualAggregator contract.
type DualAggregatorOraclePaidIterator struct {
	Event *DualAggregatorOraclePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorOraclePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorOraclePaid)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorOraclePaid)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorOraclePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorOraclePaid represents a OraclePaid event raised by the DualAggregator contract.
type DualAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOraclePaid is a free log retrieval operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_DualAggregator *DualAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*DualAggregatorOraclePaidIterator, error) {

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

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorOraclePaidIterator{contract: _DualAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

// WatchOraclePaid is a free log subscription operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_DualAggregator *DualAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *DualAggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorOraclePaid)
				if err := _DualAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

// ParseOraclePaid is a log parse operation binding the contract event 0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c.
//
// Solidity: event OraclePaid(address indexed transmitter, address indexed payee, uint256 amount, address indexed linkToken)
func (_DualAggregator *DualAggregatorFilterer) ParseOraclePaid(log types.Log) (*DualAggregatorOraclePaid, error) {
	event := new(DualAggregatorOraclePaid)
	if err := _DualAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorOwnershipTransferRequestedIterator is returned from FilterOwnershipTransferRequested and is used to iterate over the raw logs and unpacked data for OwnershipTransferRequested events raised by the DualAggregator contract.
type DualAggregatorOwnershipTransferRequestedIterator struct {
	Event *DualAggregatorOwnershipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorOwnershipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorOwnershipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorOwnershipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorOwnershipTransferRequested represents a OwnershipTransferRequested event raised by the DualAggregator contract.
type DualAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferRequested is a free log retrieval operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_DualAggregator *DualAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DualAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorOwnershipTransferRequestedIterator{contract: _DualAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferRequested is a free log subscription operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_DualAggregator *DualAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorOwnershipTransferRequested)
				if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

// ParseOwnershipTransferRequested is a log parse operation binding the contract event 0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278.
//
// Solidity: event OwnershipTransferRequested(address indexed from, address indexed to)
func (_DualAggregator *DualAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*DualAggregatorOwnershipTransferRequested, error) {
	event := new(DualAggregatorOwnershipTransferRequested)
	if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DualAggregator contract.
type DualAggregatorOwnershipTransferredIterator struct {
	Event *DualAggregatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the DualAggregator contract.
type DualAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_DualAggregator *DualAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DualAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorOwnershipTransferredIterator{contract: _DualAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_DualAggregator *DualAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DualAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorOwnershipTransferred)
				if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_DualAggregator *DualAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*DualAggregatorOwnershipTransferred, error) {
	event := new(DualAggregatorOwnershipTransferred)
	if err := _DualAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorPayeeshipTransferRequestedIterator is returned from FilterPayeeshipTransferRequested and is used to iterate over the raw logs and unpacked data for PayeeshipTransferRequested events raised by the DualAggregator contract.
type DualAggregatorPayeeshipTransferRequestedIterator struct {
	Event *DualAggregatorPayeeshipTransferRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorPayeeshipTransferRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorPayeeshipTransferRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorPayeeshipTransferRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorPayeeshipTransferRequested represents a PayeeshipTransferRequested event raised by the DualAggregator contract.
type DualAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferRequested is a free log retrieval operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_DualAggregator *DualAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*DualAggregatorPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorPayeeshipTransferRequestedIterator{contract: _DualAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferRequested is a free log subscription operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_DualAggregator *DualAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorPayeeshipTransferRequested)
				if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

// ParsePayeeshipTransferRequested is a log parse operation binding the contract event 0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367.
//
// Solidity: event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed)
func (_DualAggregator *DualAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*DualAggregatorPayeeshipTransferRequested, error) {
	event := new(DualAggregatorPayeeshipTransferRequested)
	if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorPayeeshipTransferredIterator is returned from FilterPayeeshipTransferred and is used to iterate over the raw logs and unpacked data for PayeeshipTransferred events raised by the DualAggregator contract.
type DualAggregatorPayeeshipTransferredIterator struct {
	Event *DualAggregatorPayeeshipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorPayeeshipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorPayeeshipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorPayeeshipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorPayeeshipTransferred represents a PayeeshipTransferred event raised by the DualAggregator contract.
type DualAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPayeeshipTransferred is a free log retrieval operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_DualAggregator *DualAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*DualAggregatorPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorPayeeshipTransferredIterator{contract: _DualAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

// WatchPayeeshipTransferred is a free log subscription operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_DualAggregator *DualAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *DualAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorPayeeshipTransferred)
				if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

// ParsePayeeshipTransferred is a log parse operation binding the contract event 0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3.
//
// Solidity: event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current)
func (_DualAggregator *DualAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*DualAggregatorPayeeshipTransferred, error) {
	event := new(DualAggregatorPayeeshipTransferred)
	if err := _DualAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorPrimaryFeedUnlockedIterator is returned from FilterPrimaryFeedUnlocked and is used to iterate over the raw logs and unpacked data for PrimaryFeedUnlocked events raised by the DualAggregator contract.
type DualAggregatorPrimaryFeedUnlockedIterator struct {
	Event *DualAggregatorPrimaryFeedUnlocked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorPrimaryFeedUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorPrimaryFeedUnlocked)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorPrimaryFeedUnlocked)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorPrimaryFeedUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorPrimaryFeedUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorPrimaryFeedUnlocked represents a PrimaryFeedUnlocked event raised by the DualAggregator contract.
type DualAggregatorPrimaryFeedUnlocked struct {
	PrimaryRoundId uint32
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPrimaryFeedUnlocked is a free log retrieval operation binding the contract event 0xda2435684a37fba6f7841e49b59e6ad975e462bbebd28ec9da4ed9746a6992be.
//
// Solidity: event PrimaryFeedUnlocked(uint32 indexed primaryRoundId)
func (_DualAggregator *DualAggregatorFilterer) FilterPrimaryFeedUnlocked(opts *bind.FilterOpts, primaryRoundId []uint32) (*DualAggregatorPrimaryFeedUnlockedIterator, error) {

	var primaryRoundIdRule []interface{}
	for _, primaryRoundIdItem := range primaryRoundId {
		primaryRoundIdRule = append(primaryRoundIdRule, primaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "PrimaryFeedUnlocked", primaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorPrimaryFeedUnlockedIterator{contract: _DualAggregator.contract, event: "PrimaryFeedUnlocked", logs: logs, sub: sub}, nil
}

// WatchPrimaryFeedUnlocked is a free log subscription operation binding the contract event 0xda2435684a37fba6f7841e49b59e6ad975e462bbebd28ec9da4ed9746a6992be.
//
// Solidity: event PrimaryFeedUnlocked(uint32 indexed primaryRoundId)
func (_DualAggregator *DualAggregatorFilterer) WatchPrimaryFeedUnlocked(opts *bind.WatchOpts, sink chan<- *DualAggregatorPrimaryFeedUnlocked, primaryRoundId []uint32) (event.Subscription, error) {

	var primaryRoundIdRule []interface{}
	for _, primaryRoundIdItem := range primaryRoundId {
		primaryRoundIdRule = append(primaryRoundIdRule, primaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "PrimaryFeedUnlocked", primaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorPrimaryFeedUnlocked)
				if err := _DualAggregator.contract.UnpackLog(event, "PrimaryFeedUnlocked", log); err != nil {
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

// ParsePrimaryFeedUnlocked is a log parse operation binding the contract event 0xda2435684a37fba6f7841e49b59e6ad975e462bbebd28ec9da4ed9746a6992be.
//
// Solidity: event PrimaryFeedUnlocked(uint32 indexed primaryRoundId)
func (_DualAggregator *DualAggregatorFilterer) ParsePrimaryFeedUnlocked(log types.Log) (*DualAggregatorPrimaryFeedUnlocked, error) {
	event := new(DualAggregatorPrimaryFeedUnlocked)
	if err := _DualAggregator.contract.UnpackLog(event, "PrimaryFeedUnlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorRemovedAccessIterator is returned from FilterRemovedAccess and is used to iterate over the raw logs and unpacked data for RemovedAccess events raised by the DualAggregator contract.
type DualAggregatorRemovedAccessIterator struct {
	Event *DualAggregatorRemovedAccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorRemovedAccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorRemovedAccess)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorRemovedAccess)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorRemovedAccess represents a RemovedAccess event raised by the DualAggregator contract.
type DualAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedAccess is a free log retrieval operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_DualAggregator *DualAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*DualAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorRemovedAccessIterator{contract: _DualAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

// WatchRemovedAccess is a free log subscription operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_DualAggregator *DualAggregatorFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *DualAggregatorRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorRemovedAccess)
				if err := _DualAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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

// ParseRemovedAccess is a log parse operation binding the contract event 0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1.
//
// Solidity: event RemovedAccess(address user)
func (_DualAggregator *DualAggregatorFilterer) ParseRemovedAccess(log types.Log) (*DualAggregatorRemovedAccess, error) {
	event := new(DualAggregatorRemovedAccess)
	if err := _DualAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorRequesterAccessControllerSetIterator is returned from FilterRequesterAccessControllerSet and is used to iterate over the raw logs and unpacked data for RequesterAccessControllerSet events raised by the DualAggregator contract.
type DualAggregatorRequesterAccessControllerSetIterator struct {
	Event *DualAggregatorRequesterAccessControllerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorRequesterAccessControllerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorRequesterAccessControllerSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorRequesterAccessControllerSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorRequesterAccessControllerSet represents a RequesterAccessControllerSet event raised by the DualAggregator contract.
type DualAggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRequesterAccessControllerSet is a free log retrieval operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_DualAggregator *DualAggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*DualAggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorRequesterAccessControllerSetIterator{contract: _DualAggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

// WatchRequesterAccessControllerSet is a free log subscription operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_DualAggregator *DualAggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorRequesterAccessControllerSet)
				if err := _DualAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

// ParseRequesterAccessControllerSet is a log parse operation binding the contract event 0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634.
//
// Solidity: event RequesterAccessControllerSet(address old, address current)
func (_DualAggregator *DualAggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*DualAggregatorRequesterAccessControllerSet, error) {
	event := new(DualAggregatorRequesterAccessControllerSet)
	if err := _DualAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorRoundRequestedIterator is returned from FilterRoundRequested and is used to iterate over the raw logs and unpacked data for RoundRequested events raised by the DualAggregator contract.
type DualAggregatorRoundRequestedIterator struct {
	Event *DualAggregatorRoundRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorRoundRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorRoundRequested)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorRoundRequested)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorRoundRequested represents a RoundRequested event raised by the DualAggregator contract.
type DualAggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [32]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRoundRequested is a free log retrieval operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_DualAggregator *DualAggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*DualAggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorRoundRequestedIterator{contract: _DualAggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

// WatchRoundRequested is a free log subscription operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_DualAggregator *DualAggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *DualAggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorRoundRequested)
				if err := _DualAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

// ParseRoundRequested is a log parse operation binding the contract event 0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c.
//
// Solidity: event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round)
func (_DualAggregator *DualAggregatorFilterer) ParseRoundRequested(log types.Log) (*DualAggregatorRoundRequested, error) {
	event := new(DualAggregatorRoundRequested)
	if err := _DualAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorSecondaryRoundIdUpdatedIterator is returned from FilterSecondaryRoundIdUpdated and is used to iterate over the raw logs and unpacked data for SecondaryRoundIdUpdated events raised by the DualAggregator contract.
type DualAggregatorSecondaryRoundIdUpdatedIterator struct {
	Event *DualAggregatorSecondaryRoundIdUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorSecondaryRoundIdUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorSecondaryRoundIdUpdated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorSecondaryRoundIdUpdated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorSecondaryRoundIdUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorSecondaryRoundIdUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorSecondaryRoundIdUpdated represents a SecondaryRoundIdUpdated event raised by the DualAggregator contract.
type DualAggregatorSecondaryRoundIdUpdated struct {
	SecondaryRoundId uint32
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSecondaryRoundIdUpdated is a free log retrieval operation binding the contract event 0x8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b.
//
// Solidity: event SecondaryRoundIdUpdated(uint32 indexed secondaryRoundId)
func (_DualAggregator *DualAggregatorFilterer) FilterSecondaryRoundIdUpdated(opts *bind.FilterOpts, secondaryRoundId []uint32) (*DualAggregatorSecondaryRoundIdUpdatedIterator, error) {

	var secondaryRoundIdRule []interface{}
	for _, secondaryRoundIdItem := range secondaryRoundId {
		secondaryRoundIdRule = append(secondaryRoundIdRule, secondaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "SecondaryRoundIdUpdated", secondaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorSecondaryRoundIdUpdatedIterator{contract: _DualAggregator.contract, event: "SecondaryRoundIdUpdated", logs: logs, sub: sub}, nil
}

// WatchSecondaryRoundIdUpdated is a free log subscription operation binding the contract event 0x8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b.
//
// Solidity: event SecondaryRoundIdUpdated(uint32 indexed secondaryRoundId)
func (_DualAggregator *DualAggregatorFilterer) WatchSecondaryRoundIdUpdated(opts *bind.WatchOpts, sink chan<- *DualAggregatorSecondaryRoundIdUpdated, secondaryRoundId []uint32) (event.Subscription, error) {

	var secondaryRoundIdRule []interface{}
	for _, secondaryRoundIdItem := range secondaryRoundId {
		secondaryRoundIdRule = append(secondaryRoundIdRule, secondaryRoundIdItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "SecondaryRoundIdUpdated", secondaryRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorSecondaryRoundIdUpdated)
				if err := _DualAggregator.contract.UnpackLog(event, "SecondaryRoundIdUpdated", log); err != nil {
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

// ParseSecondaryRoundIdUpdated is a log parse operation binding the contract event 0x8d530b9ddc4b318d28fdd4c3a21fcfecece54c1a72a824f262985b99afef009b.
//
// Solidity: event SecondaryRoundIdUpdated(uint32 indexed secondaryRoundId)
func (_DualAggregator *DualAggregatorFilterer) ParseSecondaryRoundIdUpdated(log types.Log) (*DualAggregatorSecondaryRoundIdUpdated, error) {
	event := new(DualAggregatorSecondaryRoundIdUpdated)
	if err := _DualAggregator.contract.UnpackLog(event, "SecondaryRoundIdUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorTransmittedIterator is returned from FilterTransmitted and is used to iterate over the raw logs and unpacked data for Transmitted events raised by the DualAggregator contract.
type DualAggregatorTransmittedIterator struct {
	Event *DualAggregatorTransmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorTransmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorTransmitted)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorTransmitted)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorTransmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorTransmitted represents a Transmitted event raised by the DualAggregator contract.
type DualAggregatorTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransmitted is a free log retrieval operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_DualAggregator *DualAggregatorFilterer) FilterTransmitted(opts *bind.FilterOpts) (*DualAggregatorTransmittedIterator, error) {

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &DualAggregatorTransmittedIterator{contract: _DualAggregator.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

// WatchTransmitted is a free log subscription operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_DualAggregator *DualAggregatorFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *DualAggregatorTransmitted) (event.Subscription, error) {

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorTransmitted)
				if err := _DualAggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

// ParseTransmitted is a log parse operation binding the contract event 0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62.
//
// Solidity: event Transmitted(bytes32 configDigest, uint32 epoch)
func (_DualAggregator *DualAggregatorFilterer) ParseTransmitted(log types.Log) (*DualAggregatorTransmitted, error) {
	event := new(DualAggregatorTransmitted)
	if err := _DualAggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DualAggregatorValidatorConfigSetIterator is returned from FilterValidatorConfigSet and is used to iterate over the raw logs and unpacked data for ValidatorConfigSet events raised by the DualAggregator contract.
type DualAggregatorValidatorConfigSetIterator struct {
	Event *DualAggregatorValidatorConfigSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DualAggregatorValidatorConfigSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DualAggregatorValidatorConfigSet)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DualAggregatorValidatorConfigSet)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DualAggregatorValidatorConfigSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DualAggregatorValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DualAggregatorValidatorConfigSet represents a ValidatorConfigSet event raised by the DualAggregator contract.
type DualAggregatorValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorConfigSet is a free log retrieval operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_DualAggregator *DualAggregatorFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*DualAggregatorValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _DualAggregator.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &DualAggregatorValidatorConfigSetIterator{contract: _DualAggregator.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

// WatchValidatorConfigSet is a free log subscription operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_DualAggregator *DualAggregatorFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *DualAggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _DualAggregator.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DualAggregatorValidatorConfigSet)
				if err := _DualAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
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

// ParseValidatorConfigSet is a log parse operation binding the contract event 0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541.
//
// Solidity: event ValidatorConfigSet(address indexed previousValidator, uint32 previousGasLimit, address indexed currentValidator, uint32 currentGasLimit)
func (_DualAggregator *DualAggregatorFilterer) ParseValidatorConfigSet(log types.Log) (*DualAggregatorValidatorConfigSet, error) {
	event := new(DualAggregatorValidatorConfigSet)
	if err := _DualAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
