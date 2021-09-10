// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testoffchainaggregator

import (
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
)

var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

const TestOffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAdminAccessController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_oldLinkToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_newLinkToken\",\"type\":\"address\"}],\"name\":\"LinkTokenSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"ValidatorConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingData\",\"outputs\":[{\"internalType\":\"uint16[31]\",\"name\":\"observationsCounts\",\"type\":\"uint16[31]\"},{\"internalType\":\"uint256[31]\",\"name\":\"gasReimbursements\",\"type\":\"uint256[31]\"},{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"_newValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_newGasLimit\",\"type\":\"uint32\"}],\"name\":\"setValidatorConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testAccountingGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"testBurnLINK\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"testDecodeReport\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"int192[]\",\"name\":\"\",\"type\":\"int192[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"txGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reasonableGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maximumGasPrice\",\"type\":\"uint256\"}],\"name\":\"testImpliedGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"testPayee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_x\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"_y\",\"type\":\"uint16\"}],\"name\":\"testSaturatingAddUint16\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitterOrSigner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amountLinkWei\",\"type\":\"uint256\"}],\"name\":\"testSetGasReimbursements\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_amount\",\"type\":\"uint16\"}],\"name\":\"testSetOracleObservationCount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testTotalLinkDue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"linkDue\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callDataCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"}],\"name\":\"testTransmitterGasCostEthWei\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorConfig\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var TestOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b50604051620064153803806200641583398181016040526101408110156200003957600080fd5b508051602080830151604080850151606086015160808088015160a089015160c08a015160e08b01516101008c0151610120909c01518851808a019099526004895263151154d560e21b9a89019a909a52600080546001600160a01b0319163317815594859052999a979995989497929691959094909390918b918b918b918b918b918b918b918b918b918b918b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b89620000e8878787878762000287565b600980546001600160a01b0319166001600160a01b0384169081179091556040516000907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a908290a36200013c8162000379565b6200014662000610565b6200015062000610565b60005b601f8160ff161015620001a0576001838260ff16601f81106200017257fe5b61ffff909216602092909202015260018260ff8316601f81106200019257fe5b602002015260010162000153565b50620001b0600b83601f6200062f565b50620001c0600f82601f620006cc565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e0525050835162000209935060329250602085019150620006fd565b506200021583620003f2565b62000222600080620004ca565b8560170b60a08160170b60401b815250508460170b60c08160170b60401b815250505050505050505050505050506001603360006101000a81548160ff0219169083151502179055505050505050505050505050505050505050505050505062000796565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526008805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b600a546001600160a01b039081169082168114620003ee57600a80546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b0316331462000452576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6031546001600160a01b039081169082168114620003ee57603180546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6000546001600160a01b031633146200052a576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b604080518082019091526030546001600160a01b03808216808452600160a01b90920463ffffffff16602084015284161415806200057857508163ffffffff16816020015163ffffffff1614155b156200060b576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052603080546001600160a01b031916841763ffffffff60a01b1916600160a01b8302179055865187860151875193168352948201528451919493909216927fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541928290030190a35b505050565b604051806103e00160405280601f906020820280368337509192915050565b600283019183908215620006ba5791602002820160005b838211156200068857835183826101000a81548161ffff021916908361ffff160217905550926020019260020160208160010104928301926001030262000646565b8015620006b85782816101000a81549061ffff021916905560020160208160010104928301926001030262000688565b505b50620006c89291506200077f565b5090565b82601f8101928215620006ba579160200282015b82811115620006ba578251825591602001919060010190620006e0565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620007355760008555620006ba565b82601f106200075057805160ff1916838001178555620006ba565b82800160010185558215620006ba5791820182811115620006ba578251825591602001919060010190620006e0565b5b80821115620006c8576000815560010162000780565b60805160f81c60a05160401c60c05160401c60e05160f81c615c33620007e2600039806113eb52508061191c5280614e5c5250806113665280614e2f5250806125da5250615c336000f3fe608060405234801561001057600080fd5b506004361061038e5760003560e01c806398e5b12a116101de578063c10753291161010f578063e5fe4577116100ad578063f2fde38b1161007c578063f2fde38b146110ea578063fa98a1c714611110578063fbffd2c114611139578063feaf968c1461115f5761038e565b8063e5fe45771461103b578063e76d516814611082578063eb4571631461108a578063eb5dcd6c146110bc5761038e565b8063dc7f0124116100e9578063dc7f012414610d97578063e285191114610d9f578063e3d0e71214610dcb578063e4902f82146110155761038e565b8063c107532914610d33578063d09dc33914610d5f578063d18bf87e14610d675761038e565b8063a118f2491161017c578063b1dc65a411610156578063b1dc65a414610ba3578063b5ab58dc14610cb4578063b633620c14610cd1578063bd82470614610cee5761038e565b8063a118f24914610b3a578063acfe7f9c14610b60578063b121e14714610b7d5761038e565b80639b764d97116101b85780639b764d97146109f05780639c849b3014610a305780639e3ceeab14610aee5780639eb6e06014610b145761038e565b806398e5b12a1461094e578063996e8298146109755780639a6fc8f51461097d5761038e565b806366cfeaf1116102c35780638038e4a1116102615780638823da6c116102305780638823da6c146108ca5780638ac28d5a146108f05780638da5cb5b146109165780638e0566de1461091e5761038e565b80638038e4a114610833578063814118341461083b57806381ff7048146108935780638205bf6a146108c25761038e565b806370efdf2d1161029d57806370efdf2d146107ab5780637284e416146107cf57806377096177146107d757806379ba50971461082b5761038e565b806366cfeaf1146106d35780636b14daf8146106db57806370da2f67146107a35761038e565b8063313ce567116103305780634fb174701161030a5780634fb174701461068d57806350d25bcd146106bb57806354fd4d50146106c3578063668a0f02146106cb5761038e565b8063313ce567146104c25780633b5cdfa2146104e05780633c04967b146105df5761038e565b8063102a474b1161036c578063102a474b146103dd578063181f5a77146103e557806322adbc781461046257806329937268146104815761038e565b80630a756983146103935780630b69df861461039d5780630eafb25b146103b7575b600080fd5b61039b611167565b005b6103a5611200565b60408051918252519081900360200190f35b6103a5600480360360208110156103cd57600080fd5b50356001600160a01b0316611206565b6103a5611335565b6103ed611344565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561042757818101518382015260200161040f565b50505050905090810190601f1680156104545780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61046a611364565b6040805160179290920b8252519081900360200190f35b610489611388565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b6104ca6113e9565b6040805160ff9092168252519081900360200190f35b610584600480360360208110156104f657600080fd5b810190602081018135600160201b81111561051057600080fd5b82018360208201111561052257600080fd5b803590602001918460018302840111600160201b8311171561054357600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061140d945050505050565b6040518083815260200180602001828103825283818151815260200191508051906020019060200280838360005b838110156105ca5781810151838201526020016105b2565b50505050905001935050505060405180910390f35b6105e7611423565b60405180886103e080838360005b8381101561060d5781810151838201526020016105f5565b5050505090500187601f60200280838360005b83811015610638578181015183820152602001610620565b505050509050018663ffffffff1681526020018563ffffffff1681526020018463ffffffff1681526020018363ffffffff1681526020018263ffffffff16815260200197505050505050505060405180910390f35b61039b600480360360408110156106a357600080fd5b506001600160a01b0381358116916020013516611544565b6103a56117d7565b6103a561185f565b6103a5611864565b6103a56118ec565b61078f600480360360408110156106f157600080fd5b6001600160a01b038235169190810190604081016020820135600160201b81111561071b57600080fd5b82018360208201111561072d57600080fd5b803590602001918460018302840111600160201b8311171561074e57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506118f2945050505050565b604080519115158252519081900360200190f35b61046a61191a565b6107b361193e565b604080516001600160a01b039092168252519081900360200190f35b6103ed61194d565b610806600480360360808110156107ed57600080fd5b50803590602081013590604081013590606001356119d5565b604080516fffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61039b6119ec565b61039b611aa2565b610843611b3c565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561087f578181015183820152602001610867565b505050509050019250505060405180910390f35b61089b611b9e565b6040805163ffffffff94851681529290931660208301528183015290519081900360600190f35b6103a5611bb9565b61039b600480360360208110156108e057600080fd5b50356001600160a01b0316611c41565b61039b6004803603602081101561090657600080fd5b50356001600160a01b0316611d13565b6107b3611d8a565b610926611d99565b604080516001600160a01b03909316835263ffffffff90911660208301528051918290030190f35b610956611dcc565b6040805169ffffffffffffffffffff9092168252519081900360200190f35b6107b3611f70565b6109a66004803603602081101561099357600080fd5b503569ffffffffffffffffffff16611f7f565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b610a1960048036036040811015610a0657600080fd5b5061ffff81358116916020013516612020565b6040805161ffff9092168252519081900360200190f35b61039b60048036036040811015610a4657600080fd5b810190602081018135600160201b811115610a6057600080fd5b820183602082011115610a7257600080fd5b803590602001918460208302840111600160201b83111715610a9357600080fd5b919390929091602081019035600160201b811115610ab057600080fd5b820183602082011115610ac257600080fd5b803590602001918460208302840111600160201b83111715610ae357600080fd5b50909250905061202c565b61039b60048036036020811015610b0457600080fd5b50356001600160a01b0316612246565b6107b360048036036020811015610b2a57600080fd5b50356001600160a01b0316612315565b61039b60048036036020811015610b5057600080fd5b50356001600160a01b0316612333565b61039b60048036036020811015610b7657600080fd5b5035612394565b61039b60048036036020811015610b9357600080fd5b50356001600160a01b0316612418565b61039b600480360360e0811015610bb957600080fd5b810181608081016060820135600160201b811115610bd657600080fd5b820183602082011115610be857600080fd5b803590602001918460018302840111600160201b83111715610c0957600080fd5b919390929091602081019035600160201b811115610c2657600080fd5b820183602082011115610c3857600080fd5b803590602001918460208302840111600160201b83111715610c5957600080fd5b919390929091602081019035600160201b811115610c7657600080fd5b820183602082011115610c8857600080fd5b803590602001918460208302840111600160201b83111715610ca957600080fd5b9193509150356124f9565b6103a560048036036020811015610cca57600080fd5b5035612a41565b6103a560048036036020811015610ce757600080fd5b5035612aca565b61039b600480360360a0811015610d0457600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516612b53565b61039b60048036036040811015610d4957600080fd5b506001600160a01b038135169060200135612c82565b6103a5612f4f565b61039b60048036036040811015610d7d57600080fd5b5080356001600160a01b0316906020013561ffff16612fe0565b61078f613037565b61039b60048036036040811015610db557600080fd5b506001600160a01b038135169060200135613040565b61039b600480360360c0811015610de157600080fd5b810190602081018135600160201b811115610dfb57600080fd5b820183602082011115610e0d57600080fd5b803590602001918460208302840111600160201b83111715610e2e57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b811115610e7d57600080fd5b820183602082011115610e8f57600080fd5b803590602001918460208302840111600160201b83111715610eb057600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929560ff853516959094909350604081019250602001359050600160201b811115610f0a57600080fd5b820183602082011115610f1c57600080fd5b803590602001918460018302840111600160201b83111715610f3d57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929567ffffffffffffffff853516959094909350604081019250602001359050600160201b811115610fa157600080fd5b820183602082011115610fb357600080fd5b803590602001918460018302840111600160201b83111715610fd457600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506130f5945050505050565b610a196004803603602081101561102b57600080fd5b50356001600160a01b0316613994565b611043613a41565b6040805195865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b6107b3613af5565b61039b600480360360408110156110a057600080fd5b5080356001600160a01b0316906020013563ffffffff16613b04565b61039b600480360360408110156110d257600080fd5b506001600160a01b0381358116916020013516613c57565b61039b6004803603602081101561110057600080fd5b50356001600160a01b0316613d9a565b6103a56004803603606081101561112657600080fd5b5080359060208101359060400135613e43565b61039b6004803603602081101561114f57600080fd5b50356001600160a01b0316613e58565b6109a6613eb9565b6000546001600160a01b031633146111bf576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b60335460ff16156111fe576033805460ff191690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b61179390565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff80821684528594840191610100900416600281111561124857fe5b600281111561125357fe5b905250905060008160200151600281111561126a57fe5b141561127a576000915050611330565b6040805160a08101825260085463ffffffff8082168352600160201b820481166020840152600160401b8204811693830193909352600160601b8104831660608301819052600160801b909104909216608082015282519091600091600190600b9060ff16601f81106112e957fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca000290506001600f846000015160ff16601f811061132757fe5b01540301925050505b919050565b600061133f613f58565b905090565b6060604051806060016040528060288152602001615bdb60289139905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b6040805160a08101825260085463ffffffff808216808452600160201b8304821660208501819052600160401b84048316958501869052600160601b8404831660608601819052600160801b90940490921660809094018490529490939290565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000606061141a8361410c565b91509150915091565b61142b615acb565b611433615acb565b6040805160a08101825260085463ffffffff808216808452600160201b8304821660208501819052600160401b84048316858701819052600160601b8504841660608701819052600160801b9095049093166080860181905286516103e0810190975260009687968796879687969295600b95600f9591939187601f8282826020028201916000905b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116114bc575050604080516103e0810191829052959c508b9450601f93509150839050845b815481526020019060010190808311611512575050505050955097509750975097509750975097505090919293949596565b6000546001600160a01b0316331461159c576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b6009546001600160a01b039081169083168114156115ba57506117d3565b604080516370a0823160e01b815230600482015290516001600160a01b038516916370a08231916024808301926020929190829003018186803b15801561160057600080fd5b505afa158015611614573d6000803e3d6000fd5b505050506040513d602081101561162a57600080fd5b5061163590506141c3565b6000816001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b15801561168457600080fd5b505afa158015611698573d6000803e3d6000fd5b505050506040513d60208110156116ae57600080fd5b50516040805163a9059cbb60e01b81526001600160a01b0386811660048301526024820184905291519293509084169163a9059cbb916044808201926020929091908290030181600087803b15801561170657600080fd5b505af115801561171a573d6000803e3d6000fd5b505050506040513d602081101561173057600080fd5b5051611783576040805162461bcd60e51b815260206004820152601f60248201527f7472616e736665722072656d61696e696e672066756e6473206661696c656400604482015290519081900360640190fd5b600980546001600160a01b0319166001600160a01b0386811691821790925560405190918416907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a90600090a350505b5050565b600061181a336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b611857576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b61133f61454c565b600481565b60006118a7336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b6118e4576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b61133f614577565b60025490565b60006118fe838361458c565b8061191157506001600160a01b03831632145b90505b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6031546001600160a01b031690565b6060611990336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b6119cd576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b61133f6145bc565b60006119e385858585614649565b95945050505050565b6001546001600160a01b03163314611a4b576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b03163314611afa576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b60335460ff166111fe576033805460ff191660011790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b60606007805480602002602001604051908101604052809291908181526020018280548015611b9457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611b76575b5050505050905090565b60045460025463ffffffff80831692600160201b9004169192565b6000611bfc336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b611c39576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b61133f6146d5565b6000546001600160a01b03163314611c99576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b6001600160a01b03811660009081526034602052604090205460ff1615611d10576001600160a01b038116600081815260346020908152604091829020805460ff19169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b6001600160a01b038181166000908152600d6020526040902054163314611d81576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b611d108161470a565b6000546001600160a01b031681565b604080518082019091526030546001600160a01b038116808352600160a01b90910463ffffffff16602090920182905291565b600080546001600160a01b0316331480611e8f575060315460408051630d629b5f60e31b815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b158015611e6257600080fd5b505afa158015611e76573d6000803e3d6000fd5b505050506040513d6020811015611e8c57600080fd5b50515b611ee0576040805162461bcd60e51b815260206004820152601d60248201527f4f6e6c79206f776e6572267265717565737465722063616e2063616c6c000000604482015290519081900360640190fd5b604080518082018252602e5464ffffffffff8116825263ffffffff65010000000000820481166020808501919091526002548551908152600884901c9092169082015260ff909116818401529151909133917f41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c9181900360600190a2806020015160010163ffffffff1691505090565b600a546001600160a01b031690565b6000806000806000611fc8336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b612005576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b61200e866148f2565b939a9299509097509550909350915050565b60006119118383614a28565b6000546001600160a01b03163314612084576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b8281146120d8576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b8381101561223f5760008585838181106120f157fe5b905060200201356001600160a01b03169050600084848481811061211157fe5b6001600160a01b038581166000908152600d6020908152604090912054920293909301358316935090911690508015808061215d5750826001600160a01b0316826001600160a01b0316145b6121ae576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b038481166000908152600d6020526040902080546001600160a01b0319168583169081179091559083161461222f57826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b5050600190920191506120db9050565b5050505050565b6000546001600160a01b0316331461229e576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b6031546001600160a01b0390811690821681146117d357603180546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6001600160a01b039081166000908152600d60205260409020541690565b6000546001600160a01b0316331461238b576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b611d1081614a40565b6009546040805163a9059cbb60e01b8152600160048201526024810184905290516001600160a01b039092169163a9059cbb916044808201926020929091908290030181600087803b1580156123e957600080fd5b505af11580156123fd573d6000803e3d6000fd5b505050506040513d602081101561241357600080fd5b505050565b6001600160a01b038181166000908152600e6020526040902054163314612486576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b6001600160a01b038181166000818152600d602090815260408083208054336001600160a01b03198083168217909355600e909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c0135916125499184918491908e908e9081908401838280828437600092019190915250614abb92505050565b6040805160608101825260025480825260035460ff808216602085015261010090910416928201929092529083146125c8576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b6125d68b8b8b8b8b8b6151c4565b60007f000000000000000000000000000000000000000000000000000000000000000015612623576002826020015183604001510160ff168161261557fe5b0460010160ff169050612631565b816020015160010160ff1690505b888114612685576040805162461bcd60e51b815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e617475726573000000000000604482015290519081900360640190fd5b8887146126d9576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b3360009081526005602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561271657fe5b600281111561272157fe5b905250905060028160200151600281111561273857fe5b14801561276c57506007816000015160ff168154811061275457fe5b6000918252602090912001546001600160a01b031633145b6127bd576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b50505050506000888860405180838380828437808301925050509250505060405180910390208a60405160200180838152602001826003602002808284378083019250505092505050604051602081830303815290604052805190602001209050612826615acb565b61282e615aea565b60005b88811015612a1b57600060018588846020811061284a57fe5b1a601b018d8d8681811061285a57fe5b905060200201358c8c8781811061286d57fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156128c8573d6000803e3d6000fd5b505060408051601f198101516001600160a01b03811660009081526005602090815290849020838501909452835460ff8082168552929650929450840191610100900416600281111561291757fe5b600281111561292257fe5b905250925060018360200151600281111561293957fe5b1461298b576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f811061299c57fe5b6020020151156129f3576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f8110612a0857fe5b9115156020909202015250600101612831565b5050505063ffffffff8110612a2c57fe5b612a368133615230565b505050505050505050565b6000612a84336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b612ac1576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b611914826153a7565b6000612b0d336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b612b4a576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b611914826153dd565b600a546000546001600160a01b039182169116331480612c14575060408051630d629b5f60e31b815233600482018181526024830193845236604484018190526001600160a01b03861694636b14daf8946000939190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b158015612be757600080fd5b505afa158015612bfb573d6000803e3d6000fd5b505050506040513d6020811015612c1157600080fd5b50515b612c65576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b612c6d6141c3565b612c7a868686868661541d565b505050505050565b6000546001600160a01b0316331480612d445750600a5460408051630d629b5f60e31b815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f909101601f1916909201965060209550909350505081840390508186803b158015612d1757600080fd5b505afa158015612d2b573d6000803e3d6000fd5b505050506040513d6020811015612d4157600080fd5b50515b612d95576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6000612d9f613f58565b600954604080516370a0823160e01b815230600482015290519293506000926001600160a01b03909216916370a0823191602480820192602092909190829003018186803b158015612df057600080fd5b505afa158015612e04573d6000803e3d6000fd5b505050506040513d6020811015612e1a57600080fd5b5051905081811015612e73576040805162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b6009546001600160a01b031663a9059cbb85612e9185850387615536565b6040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015612ed757600080fd5b505af1158015612eeb573d6000803e3d6000fd5b505050506040513d6020811015612f0157600080fd5b5051612f49576040805162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b604482015290519081900360640190fd5b50505050565b600954604080516370a0823160e01b8152306004820152905160009283926001600160a01b03909116916370a0823191602480820192602092909190829003018186803b158015612f9f57600080fd5b505afa158015612fb3573d6000803e3d6000fd5b505050506040513d6020811015612fc957600080fd5b505190506000612fd7613f58565b90910391505090565b6001600160a01b0382166000908152600560205260409020546001820190600b9060ff16601f811061300e57fe5b601091828204019190066002026101000a81548161ffff021916908361ffff1602179055505050565b60335460ff1681565b60006001600160a01b038316600090815260056020526040902054610100900460ff16600281111561306e57fe5b14156130c1576040805162461bcd60e51b815260206004820152600f60248201527f6164647265737320756e6b6e6f776e0000000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b0382166000908152600560205260409020546001820190600f9060ff16601f81106130ef57fe5b01555050565b855185518560ff16601f831115613153576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b600081116131a8576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b8183146131e65760405162461bcd60e51b8152600401808060200182810382526024815260200180615c036024913960400191505060405180910390fd5b80600302831161323d576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b6000546001600160a01b03163314613295576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a08101869052906132dc908861554d565b600654156133ca5760068054600019810191600091839081106132fb57fe5b6000918252602082200154600780546001600160a01b039092169350908490811061332257fe5b60009182526020808320909101546001600160a01b03858116845260059092526040808420805461ffff199081169091559290911680845292208054909116905560068054919250908061337257fe5b600082815260209020810160001990810180546001600160a01b0319169055019055600780548061339f57fe5b600082815260209020810160001990810180546001600160a01b0319169055019055506132dc915050565b60005b8151518110156136cb57600060056000846000015184815181106133ed57fe5b6020908102919091018101516001600160a01b0316825281019190915260400160002054610100900460ff16600281111561342457fe5b14613476576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff821681526001602082015282518051600591600091859081106134a157fe5b6020908102919091018101516001600160a01b03168252818101929092526040016000208251815460ff191660ff90911617808255918301519091829061ff0019166101008360028111156134f257fe5b0217905550600091506135029050565b600560008460200151848151811061351657fe5b6020908102919091018101516001600160a01b0316825281019190915260400160002054610100900460ff16600281111561354d57fe5b1461359f576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff8216815260208101600281525060056000846020015184815181106135cc57fe5b6020908102919091018101516001600160a01b03168252818101929092526040016000208251815460ff191660ff90911617808255918301519091829061ff00191661010083600281111561361d57fe5b02179055505082518051600692508390811061363557fe5b602090810291909101810151825460018101845560009384529282902090920180546001600160a01b0319166001600160a01b03909316929092179091558201518051600791908390811061368657fe5b60209081029190910181015182546001808201855560009485529290932090920180546001600160a01b0319166001600160a01b0390931692909217909155016133cd565b5060408101516003805460ff831660ff19909116179055600480544363ffffffff908116600160201b90810267ffffffff0000000019841617808316600101831663ffffffff199091161793849055855160208701516060880151608089015160a08a015194909604851697469761374c9789973097921695949391615567565b60026000018190555050816000015151600260010160016101000a81548160ff021916908360ff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff16856000015186602001518760400151886060015189608001518a60a00151604051808a63ffffffff1681526020018981526020018863ffffffff16815260200180602001806020018760ff168152602001806020018667ffffffffffffffff1681526020018060200185810385528b818151815260200191508051906020019060200280838360005b8381101561385557818101518382015260200161383d565b5050505090500185810384528a818151815260200191508051906020019060200280838360005b8381101561389457818101518382015260200161387c565b50505050905001858103835288818151815260200191508051906020019080838360005b838110156138d05781810151838201526020016138b8565b50505050905090810190601f1680156138fd5780820380516001836020036101000a031916815260200191505b50858103825286518152865160209182019188019080838360005b83811015613930578181015183820152602001613918565b50505050905090810190601f16801561395d5780820380516001836020036101000a031916815260200191505b509d505050505050505050505050505060405180910390a1613987826040015183606001516117d3565b5050505050505050505050565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff8082168452859484019161010090041660028111156139d657fe5b60028111156139e157fe5b90525090506000816020015160028111156139f857fe5b1415613a08576000915050611330565b6001600b826000015160ff16601f8110613a1e57fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b600080808080333214613a9b576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050600254602e5463ffffffff65010000000000820481166000908152602f60205260409020549296600883901c909116955064ffffffffff9091169350601782900b9250600160c01b90910467ffffffffffffffff1690565b6009546001600160a01b031690565b6000546001600160a01b03163314613b5c576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b604080518082019091526030546001600160a01b03808216808452600160a01b90920463ffffffff1660208401528416141580613ba957508163ffffffff16816020015163ffffffff1614155b15612413576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052603080546001600160a01b03191684177fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff16600160a01b8302179055865187860151875193168352948201528451919493909216927fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541928290030190a3505050565b6001600160a01b038281166000908152600d6020526040902054163314613cc5576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b336001600160a01b0382161415613d23576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b6001600160a01b038083166000908152600e6020526040902080548383166001600160a01b031982168117909255909116908114612413576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b6000546001600160a01b03163314613df2576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613e50848484615780565b949350505050565b6000546001600160a01b03163314613eb0576040805162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015290519081900360640190fd5b611d108161579d565b6000806000806000613f02336000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506118f292505050565b613f3f576040805162461bcd60e51b81526020600482015260096024820152684e6f2061636365737360b81b604482015290519081900360640190fd5b613f47615814565b945094509450945094509091929394565b604080516103e0810191829052600091829190600b90601f908285855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411613f755790505050505050905060005b601f811015613fe55760018282601f8110613fce57fe5b60200201510361ffff169290920191600101613fb7565b506040805160a08101825260085463ffffffff8082168352600160201b82048116602080850191909152600160401b8304821684860152600160601b8304821660608501819052600160801b90930490911660808401526007805485518184028101840190965280865296909202633b9aca000295929360009390929183018282801561409b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161407d575b5050604080516103e08101918290529495506000949350600f9250601f915082845b8154815260200190600101908083116140bd575050505050905060005b82518110156141045760018282601f81106140f157fe5b60200201510395909501946001016140da565b505050505090565b6000606082806020019051604081101561412557600080fd5b815160208301805160405192949293830192919084600160201b82111561414b57600080fd5b90830190602082018581111561416057600080fd5b82518660208202830111600160201b8211171561417c57600080fd5b82525081516020918201928201910280838360005b838110156141a9578181015183820152602001614191565b505050509050016040525050508092508193505050915091565b6040805160a08101825260085463ffffffff8082168352600160201b820481166020840152600160401b8204811683850152600160601b820481166060840152600160801b90910416608082015260095482516103e081019384905291926001600160a01b0390911691600091600b90601f908285855b82829054906101000a900461ffff1661ffff168152602001906002019060208260010104928301926001038202915080841161423a575050604080516103e08101918290529596506000959450600f9350601f9250905082845b81548152602001906001019080831161429457505050505090506000600780548060200260200160405190810160405280929190818152602001828054801561430657602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116142e8575b5050505050905060005b815181101561453057600060018483601f811061432957fe5b6020020151039050600060018684601f811061434157fe5b60200201510361ffff169050600082896060015163ffffffff168302633b9aca00020190506000811115614525576000600d600087878151811061438157fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b03169050886001600160a01b031663a9059cbb82846040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b15801561441657600080fd5b505af115801561442a573d6000803e3d6000fd5b505050506040513d602081101561444057600080fd5b5051614488576040805162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b604482015290519081900360640190fd5b60018886601f811061449657fe5b61ffff909216602092909202015260018786601f81106144b257fe5b602002018181525050886001600160a01b0316816001600160a01b03168787815181106144db57fe5b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c856040518082815260200191505060405180910390a4505b505050600101614310565b5061453e600b84601f615b01565b50612c7a600f83601f615b97565b602e5465010000000000900463ffffffff166000908152602f6020526040902054601790810b900b90565b602e5465010000000000900463ffffffff1690565b6001600160a01b03821660009081526034602052604081205460ff168061191157505060335460ff161592915050565b60328054604080516020601f6002600019610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015611b945780601f1061461d57610100808354040283529160200191611b94565b820191906000526020600020905b81548152906001019060200180831161462b57509395945050505050565b6000818510156146a0576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b818503830161179301633b9aca00858202026fffffffffffffffffffffffffffffffff81106146cb57fe5b9695505050505050565b602e5465010000000000900463ffffffff166000908152602f6020526040902054600160c01b900467ffffffffffffffff1690565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561475057fe5b600281111561475b57fe5b9052509050600061476b83611206565b90508015612413576001600160a01b038084166000908152600d6020908152604080832054600954825163a9059cbb60e01b8152918616600483018190526024830188905292519295169363a9059cbb9360448084019491939192918390030190829087803b1580156147dd57600080fd5b505af11580156147f1573d6000803e3d6000fd5b505050506040513d602081101561480757600080fd5b505161484f576040805162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b604482015290519081900360640190fd5b6001600b846000015160ff16601f811061486557fe5b601091828204019190066002026101000a81548161ffff021916908361ffff1602179055506001600f846000015160ff16601f81106148a057fe5b01556009546040805184815290516001600160a01b039283169284811692908816917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c9181900360200190a450505050565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e740000000000000000000000000000000000815250906149cd5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561499257818101518382015260200161497a565b50505050905090810190601f1680156149bf5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5050505063ffffffff83166000908152602f6020908152604091829020825180840190935254601781810b810b810b808552600160c01b90920467ffffffffffffffff1693909201839052949594900b939092508291508490565b60006119118261ffff168461ffff160161ffff615536565b6001600160a01b03811660009081526034602052604090205460ff16611d10576001600160a01b038116600081815260346020908152604091829020805460ff19166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b60408051808201909152602e5464ffffffffff8082168084526501000000000090920463ffffffff166020840152841611614b3d576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b60006060614b4a8461410c565b9092509050614b598482615872565b600354815160ff90911690601f1015614bb9576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b80600202825111614c11576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b6000825167ffffffffffffffff81118015614c2b57600080fd5b506040519080825280601f01601f191660200182016040528015614c56576020820181803683370190505b509050614c61615acb565b60005b8451811015614d50576000868260208110614c7b57fe5b1a90508281601f8110614c8a57fe5b602002015115614ce1576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b60018382601f8110614cef57fe5b91151560209283029190910152879083908110614d0857fe5b1a60f81b848381518110614d1857fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050600101614c64565b50614d5a826158d6565b64ffffffffff8816865260005b6001855103811015614e05576000858260010181518110614d8457fe5b602002602001015160170b868381518110614d9b57fe5b602002602001015160170b1315905080614dfc576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b50600101614d67565b506000846002865181614e1457fe5b0481518110614e1f57fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b13158015614e8557507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b614ed6576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b86602001805180919060010163ffffffff1663ffffffff168152505060405180604001604052808260170b81526020014267ffffffffffffffff16815250602f6000896020015163ffffffff1663ffffffff16815260200190815260200160002060008201518160000160006101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908360170b77ffffffffffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160186101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550905050866020015163ffffffff167f8235efcbf95cfe12e2d5afec1e5e568dc529cb92d6a9b4195da079f1411244f8823388878f8f604051808760170b8152602001866001600160a01b0316815260200180602001806020018581526020018464ffffffffff168152602001838103835287818151815260200191508051906020019060200280838360005b8381101561505c578181015183820152602001615044565b50505050905001838103825286818151815260200191508051906020019080838360005b83811015615098578181015183820152602001615080565b50505050905090810190601f1680156150c55780820380516001836020036101000a031916815260200191505b509850505050505050505060405180910390a260208088015160408051428152905160009363ffffffff909316927f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271928290030190a3866020015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a361517287602001518260170b615945565b50508451602e805460209097015163ffffffff16650100000000000268ffffffff00000000001964ffffffffff90931664ffffffffff1990981697909717919091169590951790945550505050505050565b602083810286019082020161014401368114615227576040805162461bcd60e51b815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d617463680000000000000000604482015290519081900360640190fd5b50505050505050565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561527657fe5b600281111561528157fe5b905250905060028160200151600281111561529857fe5b146152a257600080fd5b6040805160a08101825260085463ffffffff808216808452600160201b8304821660208501819052600160401b8404831695850195909552600160601b830482166060850152600160801b909204166080830152909160009161530c91633b9aca003a0491615780565b90506010360260005a9050600061532b8863ffffffff16858585614649565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff1683028161535957fe5b049050856080015163ffffffff16633b9aca000281600f896000015160ff16601f811061538257fe5b01540101600f886000015160ff16601f811061539a57fe5b0155505050505050505050565b600063ffffffff8211156153bd57506000611330565b5063ffffffff166000908152602f6020526040902054601790810b900b90565b600063ffffffff8211156153f357506000611330565b5063ffffffff166000908152602f6020526040902054600160c01b900467ffffffffffffffff1690565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526008805463ffffffff1916871767ffffffff000000001916600160201b8702176bffffffff00000000000000001916600160401b8502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff16600160601b8402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b600081831015615547575081611914565b50919050565b6155556141c3565b5050602e805464ffffffffff19169055565b6000808a8a8a8a8a8a8a8a8a604051602001808a8152602001896001600160a01b031681526020018867ffffffffffffffff16815260200180602001806020018760ff168152602001806020018667ffffffffffffffff1681526020018060200185810385528b818151815260200191508051906020019060200280838360005b838110156156005781810151838201526020016155e8565b5050505090500185810384528a818151815260200191508051906020019060200280838360005b8381101561563f578181015183820152602001615627565b50505050905001858103835288818151815260200191508051906020019080838360005b8381101561567b578181015183820152602001615663565b50505050905090810190601f1680156156a85780820380516001836020036101000a031916815260200191505b50858103825286518152865160209182019188019080838360005b838110156156db5781810151838201526020016156c3565b50505050905090810190601f1680156157085780820380516001836020036101000a031916815260200191505b5060408051601f1981840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179f505050505050505050505050505050509998505050505050505050565b6000838381101561579357600285850304015b6119e38184615536565b600a546001600160a01b0390811690821681146117d357600a80546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b602e5465010000000000900463ffffffff166000818152602f6020908152604091829020825180840190935254601781810b810b810b808552600160c01b90920467ffffffffffffffff1693909201839052929392900b9181908490565b60008151602002606001600001905080835114612413576040805162461bcd60e51b815260206004820152601660248201527f7265706f7274206c656e677468206d69736d6174636800000000000000000000604482015290519081900360640190fd5b604080516103e081019182905261593791839190600b90601f90826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116158f55790505050505050615a56565b6117d390600b90601f615b01565b604080518082019091526030546001600160a01b038116808352600160a01b90910463ffffffff16602083015261597c57506117d3565b600019830163ffffffff8181166000818152602f602090815260408083205487518884015183517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810197909752601792830b90920b602487018190528b88166044880152606487018b9052925192966001600160a01b039091169563beed9b51959290911693608480830194919391928390030190829088803b158015615a2757600080fd5b5087f193505050508015615a4d57506040513d6020811015615a4857600080fd5b505160015b612c7a5761223f565b615a5e615acb565b60005b8351811015615ac3576000848281518110615a7857fe5b016020015160f81c9050615a9d8482601f8110615a9157fe5b60200201516001614a28565b848260ff16601f8110615aac57fe5b61ffff909216602092909202015250600101615a61565b509092915050565b604051806103e00160405280601f906020820280368337509192915050565b604080518082019091526000808252602082015290565b600283019183908215615b875791602002820160005b83821115615b5757835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302615b17565b8015615b855782816101000a81549061ffff0219169055600201602081600101049283019260010302615b57565b505b50615b93929150615bc5565b5090565b82601f8101928215615b87579160200282015b82811115615b87578251825591602001919060010190615baa565b5b80821115615b935760008155600101615bc656fe416363657373436f6e74726f6c6c65644f6666636861696e41676772656761746f7220332e302e306f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000706000a"

func DeployTestOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAccessController common.Address, _requesterAdminAccessController common.Address) (common.Address, *types.Transaction, *TestOffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(TestOffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TestOffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _minAnswer, _maxAnswer, _billingAccessController, _requesterAdminAccessController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestOffchainAggregator{TestOffchainAggregatorCaller: TestOffchainAggregatorCaller{contract: contract}, TestOffchainAggregatorTransactor: TestOffchainAggregatorTransactor{contract: contract}, TestOffchainAggregatorFilterer: TestOffchainAggregatorFilterer{contract: contract}}, nil
}

type TestOffchainAggregator struct {
	address common.Address
	abi     abi.ABI
	TestOffchainAggregatorCaller
	TestOffchainAggregatorTransactor
	TestOffchainAggregatorFilterer
}

type TestOffchainAggregatorCaller struct {
	contract *bind.BoundContract
}

type TestOffchainAggregatorTransactor struct {
	contract *bind.BoundContract
}

type TestOffchainAggregatorFilterer struct {
	contract *bind.BoundContract
}

type TestOffchainAggregatorSession struct {
	Contract     *TestOffchainAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TestOffchainAggregatorCallerSession struct {
	Contract *TestOffchainAggregatorCaller
	CallOpts bind.CallOpts
}

type TestOffchainAggregatorTransactorSession struct {
	Contract     *TestOffchainAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type TestOffchainAggregatorRaw struct {
	Contract *TestOffchainAggregator
}

type TestOffchainAggregatorCallerRaw struct {
	Contract *TestOffchainAggregatorCaller
}

type TestOffchainAggregatorTransactorRaw struct {
	Contract *TestOffchainAggregatorTransactor
}

func NewTestOffchainAggregator(address common.Address, backend bind.ContractBackend) (*TestOffchainAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(TestOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTestOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregator{address: address, abi: abi, TestOffchainAggregatorCaller: TestOffchainAggregatorCaller{contract: contract}, TestOffchainAggregatorTransactor: TestOffchainAggregatorTransactor{contract: contract}, TestOffchainAggregatorFilterer: TestOffchainAggregatorFilterer{contract: contract}}, nil
}

func NewTestOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*TestOffchainAggregatorCaller, error) {
	contract, err := bindTestOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorCaller{contract: contract}, nil
}

func NewTestOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*TestOffchainAggregatorTransactor, error) {
	contract, err := bindTestOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorTransactor{contract: contract}, nil
}

func NewTestOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*TestOffchainAggregatorFilterer, error) {
	contract, err := bindTestOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorFilterer{contract: contract}, nil
}

func bindTestOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorTransactor.contract.Transfer(opts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.contract.Transfer(opts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) BillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "billingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) BillingAccessController() (common.Address, error) {
	return _TestOffchainAggregator.Contract.BillingAccessController(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) BillingAccessController() (common.Address, error) {
	return _TestOffchainAggregator.Contract.BillingAccessController(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) BillingData(opts *bind.CallOpts) (BillingData,

	error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "billingData")

	outstruct := new(BillingData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ObservationsCounts = *abi.ConvertType(out[0], new([31]uint16)).(*[31]uint16)
	outstruct.GasReimbursements = *abi.ConvertType(out[1], new([31]*big.Int)).(*[31]*big.Int)
	outstruct.MaximumGasPrice = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.ReasonableGasPrice = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.MicroLinkPerEth = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.LinkGweiPerObservation = *abi.ConvertType(out[5], new(uint32)).(*uint32)
	outstruct.LinkGweiPerTransmission = *abi.ConvertType(out[6], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) BillingData() (BillingData,

	error) {
	return _TestOffchainAggregator.Contract.BillingData(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) BillingData() (BillingData,

	error) {
	return _TestOffchainAggregator.Contract.BillingData(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) CheckEnabled() (bool, error) {
	return _TestOffchainAggregator.Contract.CheckEnabled(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _TestOffchainAggregator.Contract.CheckEnabled(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) Decimals() (uint8, error) {
	return _TestOffchainAggregator.Contract.Decimals(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _TestOffchainAggregator.Contract.Decimals(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) Description() (string, error) {
	return _TestOffchainAggregator.Contract.Description(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Description() (string, error) {
	return _TestOffchainAggregator.Contract.Description(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetAnswer(&_TestOffchainAggregator.CallOpts, _roundId)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetAnswer(&_TestOffchainAggregator.CallOpts, _roundId)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetBilling(opts *bind.CallOpts) (GetBilling,

	error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(GetBilling)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaximumGasPrice = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ReasonableGasPrice = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.MicroLinkPerEth = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.LinkGweiPerObservation = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.LinkGweiPerTransmission = *abi.ConvertType(out[4], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetBilling() (GetBilling,

	error) {
	return _TestOffchainAggregator.Contract.GetBilling(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetBilling() (GetBilling,

	error) {
	return _TestOffchainAggregator.Contract.GetBilling(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetConfigDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getConfigDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetConfigDigest() ([32]byte, error) {
	return _TestOffchainAggregator.Contract.GetConfigDigest(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetConfigDigest() ([32]byte, error) {
	return _TestOffchainAggregator.Contract.GetConfigDigest(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetLinkToken() (common.Address, error) {
	return _TestOffchainAggregator.Contract.GetLinkToken(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetLinkToken() (common.Address, error) {
	return _TestOffchainAggregator.Contract.GetLinkToken(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

	outstruct := new(GetRoundData)
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

func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _TestOffchainAggregator.Contract.GetRoundData(&_TestOffchainAggregator.CallOpts, _roundId)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _TestOffchainAggregator.Contract.GetRoundData(&_TestOffchainAggregator.CallOpts, _roundId)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetTimestamp(&_TestOffchainAggregator.CallOpts, _roundId)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.GetTimestamp(&_TestOffchainAggregator.CallOpts, _roundId)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _TestOffchainAggregator.Contract.HasAccess(&_TestOffchainAggregator.CallOpts, _user, _calldata)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _TestOffchainAggregator.Contract.HasAccess(&_TestOffchainAggregator.CallOpts, _user, _calldata)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestAnswer(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestAnswer(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _TestOffchainAggregator.Contract.LatestConfigDetails(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _TestOffchainAggregator.Contract.LatestConfigDetails(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestRound(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestRound(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
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

func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _TestOffchainAggregator.Contract.LatestRoundData(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _TestOffchainAggregator.Contract.LatestRoundData(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestTimestamp(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LatestTimestamp(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (LatestTransmissionDetails,

	error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "latestTransmissionDetails")

	outstruct := new(LatestTransmissionDetails)
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

func (_TestOffchainAggregator *TestOffchainAggregatorSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _TestOffchainAggregator.Contract.LatestTransmissionDetails(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _TestOffchainAggregator.Contract.LatestTransmissionDetails(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LinkAvailableForPayment(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.LinkAvailableForPayment(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MaxAnswer(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MaxAnswer(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MinAnswer(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.MinAnswer(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _TestOffchainAggregator.Contract.OracleObservationCount(&_TestOffchainAggregator.CallOpts, _signerOrTransmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _TestOffchainAggregator.Contract.OracleObservationCount(&_TestOffchainAggregator.CallOpts, _signerOrTransmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.OwedPayment(&_TestOffchainAggregator.CallOpts, _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.OwedPayment(&_TestOffchainAggregator.CallOpts, _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) Owner() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Owner(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _TestOffchainAggregator.Contract.Owner(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) RequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "requesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) RequesterAccessController() (common.Address, error) {
	return _TestOffchainAggregator.Contract.RequesterAccessController(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) RequesterAccessController() (common.Address, error) {
	return _TestOffchainAggregator.Contract.RequesterAccessController(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestAccountingGasCost(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testAccountingGasCost")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestAccountingGasCost() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestAccountingGasCost(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestAccountingGasCost() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestAccountingGasCost(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestDecodeReport(opts *bind.CallOpts, report []byte) ([32]byte, []*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testDecodeReport", report)

	if err != nil {
		return *new([32]byte), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestDecodeReport(report []byte) ([32]byte, []*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestDecodeReport(&_TestOffchainAggregator.CallOpts, report)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestDecodeReport(report []byte) ([32]byte, []*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestDecodeReport(&_TestOffchainAggregator.CallOpts, report)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestImpliedGasPrice(opts *bind.CallOpts, txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testImpliedGasPrice", txGasPrice, reasonableGasPrice, maximumGasPrice)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestImpliedGasPrice(txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestImpliedGasPrice(&_TestOffchainAggregator.CallOpts, txGasPrice, reasonableGasPrice, maximumGasPrice)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestImpliedGasPrice(txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestImpliedGasPrice(&_TestOffchainAggregator.CallOpts, txGasPrice, reasonableGasPrice, maximumGasPrice)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestPayee(opts *bind.CallOpts, _transmitter common.Address) (common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testPayee", _transmitter)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestPayee(_transmitter common.Address) (common.Address, error) {
	return _TestOffchainAggregator.Contract.TestPayee(&_TestOffchainAggregator.CallOpts, _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestPayee(_transmitter common.Address) (common.Address, error) {
	return _TestOffchainAggregator.Contract.TestPayee(&_TestOffchainAggregator.CallOpts, _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestSaturatingAddUint16(opts *bind.CallOpts, _x uint16, _y uint16) (uint16, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testSaturatingAddUint16", _x, _y)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestSaturatingAddUint16(_x uint16, _y uint16) (uint16, error) {
	return _TestOffchainAggregator.Contract.TestSaturatingAddUint16(&_TestOffchainAggregator.CallOpts, _x, _y)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestSaturatingAddUint16(_x uint16, _y uint16) (uint16, error) {
	return _TestOffchainAggregator.Contract.TestSaturatingAddUint16(&_TestOffchainAggregator.CallOpts, _x, _y)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestTotalLinkDue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testTotalLinkDue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestTotalLinkDue() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTotalLinkDue(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestTotalLinkDue() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTotalLinkDue(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TestTransmitterGasCostEthWei(opts *bind.CallOpts, initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "testTransmitterGasCostEthWei", initialGas, gasPrice, callDataCost, gasLeft)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestTransmitterGasCostEthWei(initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTransmitterGasCostEthWei(&_TestOffchainAggregator.CallOpts, initialGas, gasPrice, callDataCost, gasLeft)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TestTransmitterGasCostEthWei(initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error) {
	return _TestOffchainAggregator.Contract.TestTransmitterGasCostEthWei(&_TestOffchainAggregator.CallOpts, initialGas, gasPrice, callDataCost, gasLeft)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _TestOffchainAggregator.Contract.Transmitters(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _TestOffchainAggregator.Contract.Transmitters(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TypeAndVersion() (string, error) {
	return _TestOffchainAggregator.Contract.TypeAndVersion(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) TypeAndVersion() (string, error) {
	return _TestOffchainAggregator.Contract.TypeAndVersion(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) ValidatorConfig(opts *bind.CallOpts) (ValidatorConfig,

	error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "validatorConfig")

	outstruct := new(ValidatorConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.GasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) ValidatorConfig() (ValidatorConfig,

	error) {
	return _TestOffchainAggregator.Contract.ValidatorConfig(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) ValidatorConfig() (ValidatorConfig,

	error) {
	return _TestOffchainAggregator.Contract.ValidatorConfig(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestOffchainAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) Version() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.Version(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _TestOffchainAggregator.Contract.Version(&_TestOffchainAggregator.CallOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "acceptOwnership")
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptOwnership(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptOwnership(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AcceptPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "addAccess", _user)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AddAccess(&_TestOffchainAggregator.TransactOpts, _user)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.AddAccess(&_TestOffchainAggregator.TransactOpts, _user)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "disableAccessCheck")
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.DisableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.DisableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "enableAccessCheck")
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.EnableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.EnableAccessCheck(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "removeAccess", _user)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.RemoveAccess(&_TestOffchainAggregator.TransactOpts, _user)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.RemoveAccess(&_TestOffchainAggregator.TransactOpts, _user)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "requestNewRound")
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.RequestNewRound(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.RequestNewRound(&_TestOffchainAggregator.TransactOpts)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBilling(&_TestOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBilling(&_TestOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBillingAccessController(&_TestOffchainAggregator.TransactOpts, _billingAccessController)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetBillingAccessController(&_TestOffchainAggregator.TransactOpts, _billingAccessController)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetConfig(&_TestOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetConfig(&_TestOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetLinkToken(opts *bind.TransactOpts, _linkToken common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setLinkToken", _linkToken, _recipient)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetLinkToken(_linkToken common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetLinkToken(&_TestOffchainAggregator.TransactOpts, _linkToken, _recipient)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetLinkToken(_linkToken common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetLinkToken(&_TestOffchainAggregator.TransactOpts, _linkToken, _recipient)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetPayees(&_TestOffchainAggregator.TransactOpts, _transmitters, _payees)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetPayees(&_TestOffchainAggregator.TransactOpts, _transmitters, _payees)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, _requesterAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setRequesterAccessController", _requesterAccessController)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetRequesterAccessController(&_TestOffchainAggregator.TransactOpts, _requesterAccessController)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetRequesterAccessController(&_TestOffchainAggregator.TransactOpts, _requesterAccessController)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) SetValidatorConfig(opts *bind.TransactOpts, _newValidator common.Address, _newGasLimit uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "setValidatorConfig", _newValidator, _newGasLimit)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) SetValidatorConfig(_newValidator common.Address, _newGasLimit uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetValidatorConfig(&_TestOffchainAggregator.TransactOpts, _newValidator, _newGasLimit)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) SetValidatorConfig(_newValidator common.Address, _newGasLimit uint32) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.SetValidatorConfig(&_TestOffchainAggregator.TransactOpts, _newValidator, _newGasLimit)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TestBurnLINK(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "testBurnLINK", amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestBurnLINK(amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestBurnLINK(&_TestOffchainAggregator.TransactOpts, amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TestBurnLINK(amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestBurnLINK(&_TestOffchainAggregator.TransactOpts, amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TestSetGasReimbursements(opts *bind.TransactOpts, _transmitterOrSigner common.Address, _amountLinkWei *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "testSetGasReimbursements", _transmitterOrSigner, _amountLinkWei)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestSetGasReimbursements(_transmitterOrSigner common.Address, _amountLinkWei *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetGasReimbursements(&_TestOffchainAggregator.TransactOpts, _transmitterOrSigner, _amountLinkWei)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TestSetGasReimbursements(_transmitterOrSigner common.Address, _amountLinkWei *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetGasReimbursements(&_TestOffchainAggregator.TransactOpts, _transmitterOrSigner, _amountLinkWei)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TestSetOracleObservationCount(opts *bind.TransactOpts, _oracle common.Address, _amount uint16) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "testSetOracleObservationCount", _oracle, _amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TestSetOracleObservationCount(_oracle common.Address, _amount uint16) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetOracleObservationCount(&_TestOffchainAggregator.TransactOpts, _oracle, _amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TestSetOracleObservationCount(_oracle common.Address, _amount uint16) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TestSetOracleObservationCount(&_TestOffchainAggregator.TransactOpts, _oracle, _amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferOwnership(&_TestOffchainAggregator.TransactOpts, _to)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferOwnership(&_TestOffchainAggregator.TransactOpts, _to)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter, _proposed)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.TransferPayeeship(&_TestOffchainAggregator.TransactOpts, _transmitter, _proposed)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.Transmit(&_TestOffchainAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.Transmit(&_TestOffchainAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawFunds(&_TestOffchainAggregator.TransactOpts, _recipient, _amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawFunds(&_TestOffchainAggregator.TransactOpts, _recipient, _amount)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawPayment(&_TestOffchainAggregator.TransactOpts, _transmitter)
}

func (_TestOffchainAggregator *TestOffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _TestOffchainAggregator.Contract.WithdrawPayment(&_TestOffchainAggregator.TransactOpts, _transmitter)
}

type TestOffchainAggregatorAddedAccessIterator struct {
	Event *TestOffchainAggregatorAddedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorAddedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorAddedAccess)
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
		it.Event = new(TestOffchainAggregatorAddedAccess)
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

func (it *TestOffchainAggregatorAddedAccessIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*TestOffchainAggregatorAddedAccessIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorAddedAccessIterator{contract: _TestOffchainAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorAddedAccess) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorAddedAccess)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseAddedAccess(log types.Log) (*TestOffchainAggregatorAddedAccess, error) {
	event := new(TestOffchainAggregatorAddedAccess)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorAnswerUpdatedIterator struct {
	Event *TestOffchainAggregatorAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorAnswerUpdated)
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
		it.Event = new(TestOffchainAggregatorAnswerUpdated)
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

func (it *TestOffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*TestOffchainAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorAnswerUpdatedIterator{contract: _TestOffchainAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorAnswerUpdated)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*TestOffchainAggregatorAnswerUpdated, error) {
	event := new(TestOffchainAggregatorAnswerUpdated)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *TestOffchainAggregatorBillingAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorBillingAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorBillingAccessControllerSet)
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
		it.Event = new(TestOffchainAggregatorBillingAccessControllerSet)
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

func (it *TestOffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*TestOffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorBillingAccessControllerSetIterator{contract: _TestOffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorBillingAccessControllerSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*TestOffchainAggregatorBillingAccessControllerSet, error) {
	event := new(TestOffchainAggregatorBillingAccessControllerSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorBillingSetIterator struct {
	Event *TestOffchainAggregatorBillingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorBillingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorBillingSet)
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
		it.Event = new(TestOffchainAggregatorBillingSet)
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

func (it *TestOffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*TestOffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorBillingSetIterator{contract: _TestOffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorBillingSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*TestOffchainAggregatorBillingSet, error) {
	event := new(TestOffchainAggregatorBillingSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorCheckAccessDisabledIterator struct {
	Event *TestOffchainAggregatorCheckAccessDisabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorCheckAccessDisabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorCheckAccessDisabled)
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
		it.Event = new(TestOffchainAggregatorCheckAccessDisabled)
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

func (it *TestOffchainAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorCheckAccessDisabled struct {
	Raw types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*TestOffchainAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorCheckAccessDisabledIterator{contract: _TestOffchainAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorCheckAccessDisabled)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*TestOffchainAggregatorCheckAccessDisabled, error) {
	event := new(TestOffchainAggregatorCheckAccessDisabled)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorCheckAccessEnabledIterator struct {
	Event *TestOffchainAggregatorCheckAccessEnabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorCheckAccessEnabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorCheckAccessEnabled)
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
		it.Event = new(TestOffchainAggregatorCheckAccessEnabled)
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

func (it *TestOffchainAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorCheckAccessEnabled struct {
	Raw types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*TestOffchainAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorCheckAccessEnabledIterator{contract: _TestOffchainAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorCheckAccessEnabled)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*TestOffchainAggregatorCheckAccessEnabled, error) {
	event := new(TestOffchainAggregatorCheckAccessEnabled)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorConfigSetIterator struct {
	Event *TestOffchainAggregatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorConfigSet)
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
		it.Event = new(TestOffchainAggregatorConfigSet)
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

func (it *TestOffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	OnchainConfig             []byte
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*TestOffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorConfigSetIterator{contract: _TestOffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorConfigSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*TestOffchainAggregatorConfigSet, error) {
	event := new(TestOffchainAggregatorConfigSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorLinkTokenSetIterator struct {
	Event *TestOffchainAggregatorLinkTokenSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorLinkTokenSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorLinkTokenSet)
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
		it.Event = new(TestOffchainAggregatorLinkTokenSet)
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

func (it *TestOffchainAggregatorLinkTokenSetIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, _oldLinkToken []common.Address, _newLinkToken []common.Address) (*TestOffchainAggregatorLinkTokenSetIterator, error) {

	var _oldLinkTokenRule []interface{}
	for _, _oldLinkTokenItem := range _oldLinkToken {
		_oldLinkTokenRule = append(_oldLinkTokenRule, _oldLinkTokenItem)
	}
	var _newLinkTokenRule []interface{}
	for _, _newLinkTokenItem := range _newLinkToken {
		_newLinkTokenRule = append(_newLinkTokenRule, _newLinkTokenItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "LinkTokenSet", _oldLinkTokenRule, _newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorLinkTokenSetIterator{contract: _TestOffchainAggregator.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorLinkTokenSet, _oldLinkToken []common.Address, _newLinkToken []common.Address) (event.Subscription, error) {

	var _oldLinkTokenRule []interface{}
	for _, _oldLinkTokenItem := range _oldLinkToken {
		_oldLinkTokenRule = append(_oldLinkTokenRule, _oldLinkTokenItem)
	}
	var _newLinkTokenRule []interface{}
	for _, _newLinkTokenItem := range _newLinkToken {
		_newLinkTokenRule = append(_newLinkTokenRule, _newLinkTokenItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "LinkTokenSet", _oldLinkTokenRule, _newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorLinkTokenSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseLinkTokenSet(log types.Log) (*TestOffchainAggregatorLinkTokenSet, error) {
	event := new(TestOffchainAggregatorLinkTokenSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorNewRoundIterator struct {
	Event *TestOffchainAggregatorNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorNewRound)
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
		it.Event = new(TestOffchainAggregatorNewRound)
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

func (it *TestOffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*TestOffchainAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorNewRoundIterator{contract: _TestOffchainAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorNewRound)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseNewRound(log types.Log) (*TestOffchainAggregatorNewRound, error) {
	event := new(TestOffchainAggregatorNewRound)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorNewTransmissionIterator struct {
	Event *TestOffchainAggregatorNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorNewTransmission)
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
		it.Event = new(TestOffchainAggregatorNewTransmission)
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

func (it *TestOffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	ConfigDigest      [32]byte
	EpochAndRound     *big.Int
	Raw               types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*TestOffchainAggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorNewTransmissionIterator{contract: _TestOffchainAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorNewTransmission)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*TestOffchainAggregatorNewTransmission, error) {
	event := new(TestOffchainAggregatorNewTransmission)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorOraclePaidIterator struct {
	Event *TestOffchainAggregatorOraclePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorOraclePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorOraclePaid)
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
		it.Event = new(TestOffchainAggregatorOraclePaid)
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

func (it *TestOffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*TestOffchainAggregatorOraclePaidIterator, error) {

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

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorOraclePaidIterator{contract: _TestOffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorOraclePaid)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*TestOffchainAggregatorOraclePaid, error) {
	event := new(TestOffchainAggregatorOraclePaid)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *TestOffchainAggregatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorOwnershipTransferRequested)
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
		it.Event = new(TestOffchainAggregatorOwnershipTransferRequested)
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

func (it *TestOffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestOffchainAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorOwnershipTransferRequestedIterator{contract: _TestOffchainAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorOwnershipTransferRequested)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*TestOffchainAggregatorOwnershipTransferRequested, error) {
	event := new(TestOffchainAggregatorOwnershipTransferRequested)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorOwnershipTransferredIterator struct {
	Event *TestOffchainAggregatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorOwnershipTransferred)
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
		it.Event = new(TestOffchainAggregatorOwnershipTransferred)
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

func (it *TestOffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestOffchainAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorOwnershipTransferredIterator{contract: _TestOffchainAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorOwnershipTransferred)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*TestOffchainAggregatorOwnershipTransferred, error) {
	event := new(TestOffchainAggregatorOwnershipTransferred)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *TestOffchainAggregatorPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorPayeeshipTransferRequested)
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
		it.Event = new(TestOffchainAggregatorPayeeshipTransferRequested)
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

func (it *TestOffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*TestOffchainAggregatorPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorPayeeshipTransferRequestedIterator{contract: _TestOffchainAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorPayeeshipTransferRequested)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*TestOffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(TestOffchainAggregatorPayeeshipTransferRequested)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorPayeeshipTransferredIterator struct {
	Event *TestOffchainAggregatorPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorPayeeshipTransferred)
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
		it.Event = new(TestOffchainAggregatorPayeeshipTransferred)
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

func (it *TestOffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*TestOffchainAggregatorPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorPayeeshipTransferredIterator{contract: _TestOffchainAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorPayeeshipTransferred)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*TestOffchainAggregatorPayeeshipTransferred, error) {
	event := new(TestOffchainAggregatorPayeeshipTransferred)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorRemovedAccessIterator struct {
	Event *TestOffchainAggregatorRemovedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorRemovedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorRemovedAccess)
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
		it.Event = new(TestOffchainAggregatorRemovedAccess)
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

func (it *TestOffchainAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*TestOffchainAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorRemovedAccessIterator{contract: _TestOffchainAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorRemovedAccess)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseRemovedAccess(log types.Log) (*TestOffchainAggregatorRemovedAccess, error) {
	event := new(TestOffchainAggregatorRemovedAccess)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorRequesterAccessControllerSetIterator struct {
	Event *TestOffchainAggregatorRequesterAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorRequesterAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorRequesterAccessControllerSet)
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
		it.Event = new(TestOffchainAggregatorRequesterAccessControllerSet)
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

func (it *TestOffchainAggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*TestOffchainAggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorRequesterAccessControllerSetIterator{contract: _TestOffchainAggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorRequesterAccessControllerSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*TestOffchainAggregatorRequesterAccessControllerSet, error) {
	event := new(TestOffchainAggregatorRequesterAccessControllerSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorRoundRequestedIterator struct {
	Event *TestOffchainAggregatorRoundRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorRoundRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorRoundRequested)
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
		it.Event = new(TestOffchainAggregatorRoundRequested)
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

func (it *TestOffchainAggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [32]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*TestOffchainAggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorRoundRequestedIterator{contract: _TestOffchainAggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorRoundRequested)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseRoundRequested(log types.Log) (*TestOffchainAggregatorRoundRequested, error) {
	event := new(TestOffchainAggregatorRoundRequested)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type TestOffchainAggregatorValidatorConfigSetIterator struct {
	Event *TestOffchainAggregatorValidatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TestOffchainAggregatorValidatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestOffchainAggregatorValidatorConfigSet)
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
		it.Event = new(TestOffchainAggregatorValidatorConfigSet)
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

func (it *TestOffchainAggregatorValidatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *TestOffchainAggregatorValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TestOffchainAggregatorValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*TestOffchainAggregatorValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &TestOffchainAggregatorValidatorConfigSetIterator{contract: _TestOffchainAggregator.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _TestOffchainAggregator.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TestOffchainAggregatorValidatorConfigSet)
				if err := _TestOffchainAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
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

func (_TestOffchainAggregator *TestOffchainAggregatorFilterer) ParseValidatorConfigSet(log types.Log) (*TestOffchainAggregatorValidatorConfigSet, error) {
	event := new(TestOffchainAggregatorValidatorConfigSet)
	if err := _TestOffchainAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BillingData struct {
	ObservationsCounts      [31]uint16
	GasReimbursements       [31]*big.Int
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}
type GetBilling struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
}
type GetRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestConfigDetails struct {
	ConfigCount  uint32
	BlockNumber  uint32
	ConfigDigest [32]byte
}
type LatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestTransmissionDetails struct {
	ConfigDigest    [32]byte
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp uint64
}
type ValidatorConfig struct {
	Validator common.Address
	GasLimit  uint32
}

func (_TestOffchainAggregator *TestOffchainAggregator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TestOffchainAggregator.abi.Events["AddedAccess"].ID:
		return _TestOffchainAggregator.ParseAddedAccess(log)
	case _TestOffchainAggregator.abi.Events["AnswerUpdated"].ID:
		return _TestOffchainAggregator.ParseAnswerUpdated(log)
	case _TestOffchainAggregator.abi.Events["BillingAccessControllerSet"].ID:
		return _TestOffchainAggregator.ParseBillingAccessControllerSet(log)
	case _TestOffchainAggregator.abi.Events["BillingSet"].ID:
		return _TestOffchainAggregator.ParseBillingSet(log)
	case _TestOffchainAggregator.abi.Events["CheckAccessDisabled"].ID:
		return _TestOffchainAggregator.ParseCheckAccessDisabled(log)
	case _TestOffchainAggregator.abi.Events["CheckAccessEnabled"].ID:
		return _TestOffchainAggregator.ParseCheckAccessEnabled(log)
	case _TestOffchainAggregator.abi.Events["ConfigSet"].ID:
		return _TestOffchainAggregator.ParseConfigSet(log)
	case _TestOffchainAggregator.abi.Events["LinkTokenSet"].ID:
		return _TestOffchainAggregator.ParseLinkTokenSet(log)
	case _TestOffchainAggregator.abi.Events["NewRound"].ID:
		return _TestOffchainAggregator.ParseNewRound(log)
	case _TestOffchainAggregator.abi.Events["NewTransmission"].ID:
		return _TestOffchainAggregator.ParseNewTransmission(log)
	case _TestOffchainAggregator.abi.Events["OraclePaid"].ID:
		return _TestOffchainAggregator.ParseOraclePaid(log)
	case _TestOffchainAggregator.abi.Events["OwnershipTransferRequested"].ID:
		return _TestOffchainAggregator.ParseOwnershipTransferRequested(log)
	case _TestOffchainAggregator.abi.Events["OwnershipTransferred"].ID:
		return _TestOffchainAggregator.ParseOwnershipTransferred(log)
	case _TestOffchainAggregator.abi.Events["PayeeshipTransferRequested"].ID:
		return _TestOffchainAggregator.ParsePayeeshipTransferRequested(log)
	case _TestOffchainAggregator.abi.Events["PayeeshipTransferred"].ID:
		return _TestOffchainAggregator.ParsePayeeshipTransferred(log)
	case _TestOffchainAggregator.abi.Events["RemovedAccess"].ID:
		return _TestOffchainAggregator.ParseRemovedAccess(log)
	case _TestOffchainAggregator.abi.Events["RequesterAccessControllerSet"].ID:
		return _TestOffchainAggregator.ParseRequesterAccessControllerSet(log)
	case _TestOffchainAggregator.abi.Events["RoundRequested"].ID:
		return _TestOffchainAggregator.ParseRoundRequested(log)
	case _TestOffchainAggregator.abi.Events["ValidatorConfigSet"].ID:
		return _TestOffchainAggregator.ParseValidatorConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TestOffchainAggregatorAddedAccess) Topic() common.Hash {
	return common.HexToHash("0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4")
}

func (TestOffchainAggregatorAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (TestOffchainAggregatorBillingAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912")
}

func (TestOffchainAggregatorBillingSet) Topic() common.Hash {
	return common.HexToHash("0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6")
}

func (TestOffchainAggregatorCheckAccessDisabled) Topic() common.Hash {
	return common.HexToHash("0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638")
}

func (TestOffchainAggregatorCheckAccessEnabled) Topic() common.Hash {
	return common.HexToHash("0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480")
}

func (TestOffchainAggregatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (TestOffchainAggregatorLinkTokenSet) Topic() common.Hash {
	return common.HexToHash("0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a")
}

func (TestOffchainAggregatorNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (TestOffchainAggregatorNewTransmission) Topic() common.Hash {
	return common.HexToHash("0x8235efcbf95cfe12e2d5afec1e5e568dc529cb92d6a9b4195da079f1411244f8")
}

func (TestOffchainAggregatorOraclePaid) Topic() common.Hash {
	return common.HexToHash("0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c")
}

func (TestOffchainAggregatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (TestOffchainAggregatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (TestOffchainAggregatorPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (TestOffchainAggregatorPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (TestOffchainAggregatorRemovedAccess) Topic() common.Hash {
	return common.HexToHash("0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1")
}

func (TestOffchainAggregatorRequesterAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634")
}

func (TestOffchainAggregatorRoundRequested) Topic() common.Hash {
	return common.HexToHash("0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c")
}

func (TestOffchainAggregatorValidatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541")
}

func (_TestOffchainAggregator *TestOffchainAggregator) Address() common.Address {
	return _TestOffchainAggregator.address
}

type TestOffchainAggregatorInterface interface {
	BillingAccessController(opts *bind.CallOpts) (common.Address, error)

	BillingData(opts *bind.CallOpts) (BillingData,

		error)

	CheckEnabled(opts *bind.CallOpts) (bool, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error)

	GetBilling(opts *bind.CallOpts) (GetBilling,

		error)

	GetConfigDigest(opts *bind.CallOpts) ([32]byte, error)

	GetLinkToken(opts *bind.CallOpts) (common.Address, error)

	GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

		error)

	GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error)

	HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	LatestTimestamp(opts *bind.CallOpts) (*big.Int, error)

	LatestTransmissionDetails(opts *bind.CallOpts) (LatestTransmissionDetails,

		error)

	LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error)

	MaxAnswer(opts *bind.CallOpts) (*big.Int, error)

	MinAnswer(opts *bind.CallOpts) (*big.Int, error)

	OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error)

	OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	RequesterAccessController(opts *bind.CallOpts) (common.Address, error)

	TestAccountingGasCost(opts *bind.CallOpts) (*big.Int, error)

	TestDecodeReport(opts *bind.CallOpts, report []byte) ([32]byte, []*big.Int, error)

	TestImpliedGasPrice(opts *bind.CallOpts, txGasPrice *big.Int, reasonableGasPrice *big.Int, maximumGasPrice *big.Int) (*big.Int, error)

	TestPayee(opts *bind.CallOpts, _transmitter common.Address) (common.Address, error)

	TestSaturatingAddUint16(opts *bind.CallOpts, _x uint16, _y uint16) (uint16, error)

	TestTotalLinkDue(opts *bind.CallOpts) (*big.Int, error)

	TestTransmitterGasCostEthWei(opts *bind.CallOpts, initialGas *big.Int, gasPrice *big.Int, callDataCost *big.Int, gasLeft *big.Int) (*big.Int, error)

	Transmitters(opts *bind.CallOpts) ([]common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	ValidatorConfig(opts *bind.CallOpts) (ValidatorConfig,

		error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error)

	AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error)

	DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error)

	EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error)

	RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error)

	RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error)

	SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error)

	SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error)

	SetLinkToken(opts *bind.TransactOpts, _linkToken common.Address, _recipient common.Address) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error)

	SetRequesterAccessController(opts *bind.TransactOpts, _requesterAccessController common.Address) (*types.Transaction, error)

	SetValidatorConfig(opts *bind.TransactOpts, _newValidator common.Address, _newGasLimit uint32) (*types.Transaction, error)

	TestBurnLINK(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	TestSetGasReimbursements(opts *bind.TransactOpts, _transmitterOrSigner common.Address, _amountLinkWei *big.Int) (*types.Transaction, error)

	TestSetOracleObservationCount(opts *bind.TransactOpts, _oracle common.Address, _amount uint16) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error)

	FilterAddedAccess(opts *bind.FilterOpts) (*TestOffchainAggregatorAddedAccessIterator, error)

	WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorAddedAccess) (event.Subscription, error)

	ParseAddedAccess(log types.Log) (*TestOffchainAggregatorAddedAccess, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*TestOffchainAggregatorAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*TestOffchainAggregatorAnswerUpdated, error)

	FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*TestOffchainAggregatorBillingAccessControllerSetIterator, error)

	WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorBillingAccessControllerSet) (event.Subscription, error)

	ParseBillingAccessControllerSet(log types.Log) (*TestOffchainAggregatorBillingAccessControllerSet, error)

	FilterBillingSet(opts *bind.FilterOpts) (*TestOffchainAggregatorBillingSetIterator, error)

	WatchBillingSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorBillingSet) (event.Subscription, error)

	ParseBillingSet(log types.Log) (*TestOffchainAggregatorBillingSet, error)

	FilterCheckAccessDisabled(opts *bind.FilterOpts) (*TestOffchainAggregatorCheckAccessDisabledIterator, error)

	WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorCheckAccessDisabled) (event.Subscription, error)

	ParseCheckAccessDisabled(log types.Log) (*TestOffchainAggregatorCheckAccessDisabled, error)

	FilterCheckAccessEnabled(opts *bind.FilterOpts) (*TestOffchainAggregatorCheckAccessEnabledIterator, error)

	WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorCheckAccessEnabled) (event.Subscription, error)

	ParseCheckAccessEnabled(log types.Log) (*TestOffchainAggregatorCheckAccessEnabled, error)

	FilterConfigSet(opts *bind.FilterOpts) (*TestOffchainAggregatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*TestOffchainAggregatorConfigSet, error)

	FilterLinkTokenSet(opts *bind.FilterOpts, _oldLinkToken []common.Address, _newLinkToken []common.Address) (*TestOffchainAggregatorLinkTokenSetIterator, error)

	WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorLinkTokenSet, _oldLinkToken []common.Address, _newLinkToken []common.Address) (event.Subscription, error)

	ParseLinkTokenSet(log types.Log) (*TestOffchainAggregatorLinkTokenSet, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*TestOffchainAggregatorNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*TestOffchainAggregatorNewRound, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*TestOffchainAggregatorNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*TestOffchainAggregatorNewTransmission, error)

	FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*TestOffchainAggregatorOraclePaidIterator, error)

	WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error)

	ParseOraclePaid(log types.Log) (*TestOffchainAggregatorOraclePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestOffchainAggregatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*TestOffchainAggregatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestOffchainAggregatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*TestOffchainAggregatorOwnershipTransferred, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*TestOffchainAggregatorPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*TestOffchainAggregatorPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*TestOffchainAggregatorPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*TestOffchainAggregatorPayeeshipTransferred, error)

	FilterRemovedAccess(opts *bind.FilterOpts) (*TestOffchainAggregatorRemovedAccessIterator, error)

	WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorRemovedAccess) (event.Subscription, error)

	ParseRemovedAccess(log types.Log) (*TestOffchainAggregatorRemovedAccess, error)

	FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*TestOffchainAggregatorRequesterAccessControllerSetIterator, error)

	WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorRequesterAccessControllerSet) (event.Subscription, error)

	ParseRequesterAccessControllerSet(log types.Log) (*TestOffchainAggregatorRequesterAccessControllerSet, error)

	FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*TestOffchainAggregatorRoundRequestedIterator, error)

	WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorRoundRequested, requester []common.Address) (event.Subscription, error)

	ParseRoundRequested(log types.Log) (*TestOffchainAggregatorRoundRequested, error)

	FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*TestOffchainAggregatorValidatorConfigSetIterator, error)

	WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *TestOffchainAggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error)

	ParseValidatorConfigSet(log types.Log) (*TestOffchainAggregatorValidatorConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
