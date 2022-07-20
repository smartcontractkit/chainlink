// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accesscontrolledoffchainaggregator

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

const AccessControlledOffchainAggregatorABI = "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"},{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_link\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"_minAnswer\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"_maxAnswer\",\"type\":\"int192\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"AddedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"CheckAccessEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_oldLinkToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_newLinkToken\",\"type\":\"address\"}],\"name\":\"LinkTokenSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RemovedAccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"ValidatorConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"addAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"billingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableAccessCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"_roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerOrTransmitter\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"removeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_threshold\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"_offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"_offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"_linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"_newValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_newGasLimit\",\"type\":\"uint32\"}],\"name\":\"setValidatorConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorConfig\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var AccessControlledOffchainAggregatorBin = "0x6101006040523480156200001257600080fd5b506040516200655b3803806200655b83398181016040526101808110156200003957600080fd5b815160208301516040808501516060860151608087015160a088015160c089015160e08a01516101008b01516101208c01516101408d01516101608e0180519a519c9e9b9d999c989b979a969995989497939692959194939182019284640100000000821115620000a957600080fd5b908301906020820185811115620000bf57600080fd5b8251640100000000811182820188101715620000da57600080fd5b82525081516020918201929091019080838360005b8381101562000109578181015183820152602001620000ef565b50505050905090810190601f168015620001375780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b03191633178155608052508b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b89620001758787878787620002f3565b600980546001600160a01b0319166001600160a01b0384169081179091556040516000907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a908290a3620001c981620003e5565b620001d36200067c565b620001dd6200067c565b60005b601f8160ff1610156200022d576001838260ff16601f8110620001ff57fe5b61ffff909216602092909202015260018260ff8316601f81106200021f57fe5b6020020152600101620001e0565b506200023d600b83601f6200069b565b506200024d600f82601f62000738565b505050505060f887901b7fff000000000000000000000000000000000000000000000000000000000000001660e052505083516200029693506032925060208501915062000769565b50620002a2836200045e565b620002af60008062000536565b50505050601791820b820b604090811b60a05290820b90910b901b60c05250506033805460ff1916600117905550620008029e505050505050505050505050505050565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a1660809889018190526008805463ffffffff1916871763ffffffff60201b191664010000000087021763ffffffff60401b19166801000000000000000085021763ffffffff60601b19166c0100000000000000000000000084021763ffffffff60801b1916600160801b830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b600a546001600160a01b0390811690821681146200045a57600a80546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15b5050565b6000546001600160a01b03163314620004be576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6031546001600160a01b0390811690821681146200045a57603180546001600160a01b0319166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6000546001600160a01b0316331462000596576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b604080518082019091526030546001600160a01b03808216808452600160a01b90920463ffffffff1660208401528416141580620005e457508163ffffffff16816020015163ffffffff1614155b1562000677576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052603080546001600160a01b031916841763ffffffff60a01b1916600160a01b8302179055865187860151875193168352948201528451919493909216927fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541928290030190a35b505050565b604051806103e00160405280601f906020820280368337509192915050565b600283019183908215620007265791602002820160005b83821115620006f457835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302620006b2565b8015620007245782816101000a81549061ffff0219169055600201602081600101049283019260010302620006f4565b505b5062000734929150620007eb565b5090565b82601f810192821562000726579160200282015b82811115620007265782518255916020019190600101906200074c565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620007a1576000855562000726565b82601f10620007bc57805160ff191683800117855562000726565b82800160010185558215620007265791820182811115620007265782518255916020019190600101906200074c565b5b80821115620007345760008155600101620007ec565b60805160f81c60a05160401c60c05160401c60e05160f81c615d0d6200084e600039806110b952508061152b5280614aaa5250806110185280614a7d5250806122685250615d0d6000f3fe608060405234801561001057600080fd5b506004361061030a5760003560e01c806398e5b12a1161019c578063c1075329116100ee578063e76d516811610097578063f2fde38b11610071578063f2fde38b14610d99578063fbffd2c114610dbf578063feaf968c14610de55761030a565b8063e76d516814610d31578063eb45716314610d39578063eb5dcd6c14610d6b5761030a565b8063e3d0e712116100c8578063e3d0e71214610a5b578063e4902f8214610cad578063e5fe457714610cea5761030a565b8063c107532914610a1f578063d09dc33914610a4b578063dc7f012414610a535761030a565b8063a118f24911610150578063b5ab58dc1161012a578063b5ab58dc146109a0578063b633620c146109bd578063bd824706146109da5761030a565b8063a118f2491461083d578063b121e14714610863578063b1dc65a4146108895761030a565b80639a6fc8f5116101815780639a6fc8f5146106e25780639c849b30146107555780639e3ceeab146108175761030a565b806398e5b12a146106b3578063996e8298146106da5761030a565b806370da2f671161026057806381ff7048116102095780638ac28d5a116101e35780638ac28d5a146106555780638da5cb5b1461067b5780638e0566de146106835761030a565b806381ff7048146105f85780638205bf6a146106275780638823da6c1461062f5761030a565b806379ba50971161023a57806379ba5097146105905780638038e4a11461059857806381411834146105a05761030a565b806370da2f671461055c57806370efdf2d146105645780637284e416146105885761030a565b8063313ce567116102c257806354fd4d501161029c57806354fd4d5014610482578063668a0f021461048a5780636b14daf8146104925761030a565b8063313ce5671461042e5780634fb174701461044c57806350d25bcd1461047a5761030a565b8063181f5a77116102f3578063181f5a771461035157806322adbc78146103ce57806329937268146103ed5761030a565b80630a7569831461030f5780630eafb25b14610319575b600080fd5b610317610ded565b005b61033f6004803603602081101561032f57600080fd5b50356001600160a01b0316610eab565b60408051918252519081900360200190f35b610359610ff6565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561039357818101518382015260200161037b565b50505050905090810190601f1680156103c05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103d6611016565b6040805160179290920b8252519081900360200190f35b6103f561103a565b6040805163ffffffff96871681529486166020860152928516848401529084166060840152909216608082015290519081900360a00190f35b6104366110b7565b6040805160ff9092168252519081900360200190f35b6103176004803603604081101561046257600080fd5b506001600160a01b03813581169160200135166110db565b61033f6113bf565b61033f611460565b61033f611465565b610548600480360360408110156104a857600080fd5b6001600160a01b0382351691908101906040810160208201356401000000008111156104d357600080fd5b8201836020820111156104e557600080fd5b8035906020019184600183028401116401000000008311171561050757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611501945050505050565b604080519115158252519081900360200190f35b6103d6611529565b61056c61154d565b604080516001600160a01b039092168252519081900360200190f35b61035961155c565b6103176115f8565b6103176116c6565b6105a8611785565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156105e45781810151838201526020016105cc565b505050509050019250505060405180910390f35b6106006117e7565b6040805163ffffffff94851681529290931660208301528183015290519081900360600190f35b61033f611803565b6103176004803603602081101561064557600080fd5b50356001600160a01b031661189f565b6103176004803603602081101561066b57600080fd5b50356001600160a01b0316611996565b61056c611a0d565b61068b611a1c565b604080516001600160a01b03909316835263ffffffff90911660208301528051918290030190f35b6106bb611a60565b6040805169ffffffffffffffffffff9092168252519081900360200190f35b61056c611c3b565b61070b600480360360208110156106f857600080fd5b503569ffffffffffffffffffff16611c4a565b604051808669ffffffffffffffffffff1681526020018581526020018481526020018381526020018269ffffffffffffffffffff1681526020019550505050505060405180910390f35b6103176004803603604081101561076b57600080fd5b81019060208101813564010000000081111561078657600080fd5b82018360208201111561079857600080fd5b803590602001918460208302840111640100000000831117156107ba57600080fd5b9193909290916020810190356401000000008111156107d857600080fd5b8201836020820111156107ea57600080fd5b8035906020019184602083028401116401000000008311171561080c57600080fd5b509092509050611cff565b6103176004803603602081101561082d57600080fd5b50356001600160a01b0316611f38565b6103176004803603602081101561085357600080fd5b50356001600160a01b0316612026565b6103176004803603602081101561087957600080fd5b50356001600160a01b031661208e565b610317600480360360e081101561089f57600080fd5b8101816080810160608201356401000000008111156108bd57600080fd5b8201836020820111156108cf57600080fd5b803590602001918460018302840111640100000000831117156108f157600080fd5b91939092909160208101903564010000000081111561090f57600080fd5b82018360208201111561092157600080fd5b8035906020019184602083028401116401000000008311171561094357600080fd5b91939092909160208101903564010000000081111561096157600080fd5b82018360208201111561097357600080fd5b8035906020019184602083028401116401000000008311171561099557600080fd5b919350915035612187565b61033f600480360360208110156109b657600080fd5b50356126ed565b61033f600480360360208110156109d357600080fd5b503561278a565b610317600480360360a08110156109f057600080fd5b5063ffffffff813581169160208101358216916040820135811691606081013582169160809091013516612827565b61031760048036036040811015610a3557600080fd5b506001600160a01b03813516906020013561298d565b61033f612cb5565b610548612d5f565b610317600480360360c0811015610a7157600080fd5b810190602081018135640100000000811115610a8c57600080fd5b820183602082011115610a9e57600080fd5b80359060200191846020830284011164010000000083111715610ac057600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050640100000000811115610b1057600080fd5b820183602082011115610b2257600080fd5b80359060200191846020830284011164010000000083111715610b4457600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929560ff853516959094909350604081019250602001359050640100000000811115610b9f57600080fd5b820183602082011115610bb157600080fd5b80359060200191846001830284011164010000000083111715610bd357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929567ffffffffffffffff853516959094909350604081019250602001359050640100000000811115610c3857600080fd5b820183602082011115610c4a57600080fd5b80359060200191846001830284011164010000000083111715610c6c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550612d68945050505050565b610cd360048036036020811015610cc357600080fd5b50356001600160a01b03166137ac565b6040805161ffff9092168252519081900360200190f35b610cf2613859565b6040805195865263ffffffff909416602086015260ff9092168484015260170b606084015267ffffffffffffffff166080830152519081900360a00190f35b61056c613922565b61031760048036036040811015610d4f57600080fd5b5080356001600160a01b0316906020013563ffffffff16613931565b61031760048036036040811015610d8157600080fd5b506001600160a01b0381358116916020013516613ac6565b61031760048036036020811015610daf57600080fd5b50356001600160a01b0316613c21565b61031760048036036020811015610dd557600080fd5b50356001600160a01b0316613ce9565b61070b613d51565b6000546001600160a01b03163314610e4c576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60335460ff1615610ea957603380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f53963890600090a15b565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff808216845285948401916101009004166002811115610eed57fe5b6002811115610ef857fe5b9052509050600081602001516002811115610f0f57fe5b1415610f1f576000915050610ff1565b6040805160a08101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c010000000000000000000000008104831660608301819052700100000000000000000000000000000000909104909216608082015282519091600091600190600b9060ff16601f8110610faa57fe5b601091828204019190066002029054906101000a900461ffff160361ffff1602633b9aca000290506001600f846000015160ff16601f8110610fe857fe5b01540301925050505b919050565b6060604051806060016040528060288152602001615cb560289139905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b6040805160a08101825260085463ffffffff808216808452640100000000830482166020850181905268010000000000000000840483169585018690526c01000000000000000000000000840483166060860181905270010000000000000000000000000000000090940490921660809094018490529490939290565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000546001600160a01b0316331461113a576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6009546001600160a01b0390811690831681141561115857506113bb565b604080517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015290516001600160a01b038516916370a08231916024808301926020929190829003018186803b1580156111b757600080fd5b505afa1580156111cb573d6000803e3d6000fd5b505050506040513d60208110156111e157600080fd5b506111ec9050613e04565b6000816001600160a01b03166370a08231306040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b15801561123b57600080fd5b505afa15801561124f573d6000803e3d6000fd5b505050506040513d602081101561126557600080fd5b5051604080517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b0386811660048301526024820184905291519293509084169163a9059cbb916044808201926020929091908290030181600087803b1580156112d657600080fd5b505af11580156112ea573d6000803e3d6000fd5b505050506040513d602081101561130057600080fd5b5051611353576040805162461bcd60e51b815260206004820152601f60248201527f7472616e736665722072656d61696e696e672066756e6473206661696c656400604482015290519081900360640190fd5b600980547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0386811691821790925560405190918416907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a90600090a350505b5050565b6000611402336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b611453576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61145b6141b4565b905090565b600481565b60006114a8336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b6114f9576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61145b6141df565b600061150d83836141f4565b8061152057506001600160a01b03831632145b90505b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6031546001600160a01b031690565b606061159f336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b6115f0576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61145b614224565b6001546001600160a01b03163314611657576040805162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015290519081900360640190fd5b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b03163314611725576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b60335460ff16610ea957603380547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517faebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c348090600090a1565b606060078054806020026020016040519081016040528092919081815260200182805480156117dd57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116117bf575b5050505050905090565b60045460025463ffffffff808316926401000000009004169192565b6000611846336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b611897576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61145b6142cf565b6000546001600160a01b031633146118fe576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6001600160a01b03811660009081526034602052604090205460ff1615611993576001600160a01b03811660008181526034602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055815192835290517f3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d19281900390910190a15b50565b6001600160a01b038181166000908152600d6020526040902054163314611a04576040805162461bcd60e51b815260206004820152601760248201527f4f6e6c792070617965652063616e207769746864726177000000000000000000604482015290519081900360640190fd5b61199381614319565b6000546001600160a01b031681565b604080518082019091526030546001600160a01b0381168083527401000000000000000000000000000000000000000090910463ffffffff16602090920182905291565b600080546001600160a01b0316331480611b5a5750603154604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015611b2d57600080fd5b505afa158015611b41573d6000803e3d6000fd5b505050506040513d6020811015611b5757600080fd5b50515b611bab576040805162461bcd60e51b815260206004820152601d60248201527f4f6e6c79206f776e6572267265717565737465722063616e2063616c6c000000604482015290519081900360640190fd5b604080518082018252602e5464ffffffffff8116825263ffffffff65010000000000820481166020808501919091526002548551908152600884901c9092169082015260ff909116818401529151909133917f41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c9181900360600190a2806020015160010163ffffffff1691505090565b600a546001600160a01b031690565b6000806000806000611c93336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b611ce4576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b611ced86614525565b939a9299509097509550909350915050565b6000546001600160a01b03163314611d5e576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b828114611db2576040805162461bcd60e51b815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a65604482015290519081900360640190fd5b60005b83811015611f31576000858583818110611dcb57fe5b905060200201356001600160a01b031690506000848484818110611deb57fe5b6001600160a01b038581166000908152600d60209081526040909120549202939093013583169350909116905080158080611e375750826001600160a01b0316826001600160a01b0316145b611e88576040805162461bcd60e51b815260206004820152601160248201527f706179656520616c726561647920736574000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b038481166000908152600d6020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685831690811790915590831614611f2157826001600160a01b0316826001600160a01b0316856001600160a01b03167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b505060019092019150611db59050565b5050505050565b6000546001600160a01b03163314611f97576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6031546001600160a01b0390811690821681146113bb57603180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349281900390910190a15050565b6000546001600160a01b03163314612085576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b61199381614670565b6001600160a01b038181166000908152600e60205260409020541633146120fc576040805162461bcd60e51b815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e2061636365707400604482015290519081900360640190fd5b6001600160a01b038181166000818152600d602090815260408083208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217909355600e909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a604080516020601f8b018190048102820181019092528981529192508a3591818c0135916121d79184918491908e908e908190840183828082843760009201919091525061470992505050565b6040805160608101825260025480825260035460ff80821660208501526101009091041692820192909252908314612256576040805162461bcd60e51b815260206004820152601560248201527f636f6e666967446967657374206d69736d617463680000000000000000000000604482015290519081900360640190fd5b6122648b8b8b8b8b8b614e42565b60007f0000000000000000000000000000000000000000000000000000000000000000156122b1576002826020015183604001510160ff16816122a357fe5b0460010160ff1690506122bf565b816020015160010160ff1690505b888114612313576040805162461bcd60e51b815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e617475726573000000000000604482015290519081900360640190fd5b888714612367576040805162461bcd60e51b815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e0000604482015290519081900360640190fd5b3360009081526005602090815260408083208151808301909252805460ff808216845292939192918401916101009091041660028111156123a457fe5b60028111156123af57fe5b90525090506002816020015160028111156123c657fe5b1480156123fa57506007816000015160ff16815481106123e257fe5b6000918252602090912001546001600160a01b031633145b61244b576040805162461bcd60e51b815260206004820152601860248201527f756e617574686f72697a6564207472616e736d69747465720000000000000000604482015290519081900360640190fd5b50505050506000888860405180838380828437808301925050509250505060405180910390208a604051602001808381526020018260036020028082843780830192505050925050506040516020818303038152906040528051906020012090506124b4615ba5565b6124bc615bc4565b60005b888110156126c75760006001858884602081106124d857fe5b1a601b018d8d868181106124e857fe5b905060200201358c8c878181106124fb57fe5b9050602002013560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015612556573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101516001600160a01b03811660009081526005602090815290849020838501909452835460ff808216855292965092945084019161010090041660028111156125c357fe5b60028111156125ce57fe5b90525092506001836020015160028111156125e557fe5b14612637576040805162461bcd60e51b815260206004820152601e60248201527f61646472657373206e6f7420617574686f72697a656420746f207369676e0000604482015290519081900360640190fd5b8251849060ff16601f811061264857fe5b60200201511561269f576040805162461bcd60e51b815260206004820152601460248201527f6e6f6e2d756e69717565207369676e6174757265000000000000000000000000604482015290519081900360640190fd5b600184846000015160ff16601f81106126b457fe5b91151560209092020152506001016124bf565b5050505063ffffffff81106126d857fe5b6126e28133614eae565b505050505050505050565b6000612730336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b612781576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61152382615041565b60006127cd336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b61281e576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b61152382615077565b600a546000546001600160a01b03918216911633148061291f5750604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b03861694636b14daf8946000939190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b1580156128f257600080fd5b505afa158015612906573d6000803e3d6000fd5b505050506040513d602081101561291c57600080fd5b50515b612970576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b612978613e04565b61298586868686866150cc565b505050505050565b6000546001600160a01b0316331480612a865750600a54604080517f6b14daf800000000000000000000000000000000000000000000000000000000815233600482018181526024830193845236604484018190526001600160a01b0390951694636b14daf894929360009391929190606401848480828437600083820152604051601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016909201965060209550909350505081840390508186803b158015612a5957600080fd5b505afa158015612a6d573d6000803e3d6000fd5b505050506040513d6020811015612a8357600080fd5b50515b612ad7576040805162461bcd60e51b815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c604482015290519081900360640190fd5b6000612ae1615246565b600954604080517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015290519293506000926001600160a01b03909216916370a0823191602480820192602092909190829003018186803b158015612b4b57600080fd5b505afa158015612b5f573d6000803e3d6000fd5b505050506040513d6020811015612b7557600080fd5b5051905081811015612bce576040805162461bcd60e51b815260206004820152601460248201527f696e73756666696369656e742062616c616e6365000000000000000000000000604482015290519081900360640190fd5b6009546001600160a01b031663a9059cbb85612bec85850387615416565b6040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b158015612c3257600080fd5b505af1158015612c46573d6000803e3d6000fd5b505050506040513d6020811015612c5c57600080fd5b5051612caf576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b50505050565b600954604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905160009283926001600160a01b03909116916370a0823191602480820192602092909190829003018186803b158015612d1e57600080fd5b505afa158015612d32573d6000803e3d6000fd5b505050506040513d6020811015612d4857600080fd5b505190506000612d56615246565b90910391505090565b60335460ff1681565b855185518560ff16601f831115612dc6576040805162461bcd60e51b815260206004820152601060248201527f746f6f206d616e79207369676e65727300000000000000000000000000000000604482015290519081900360640190fd5b60008111612e1b576040805162461bcd60e51b815260206004820152601a60248201527f7468726573686f6c64206d75737420626520706f736974697665000000000000604482015290519081900360640190fd5b818314612e595760405162461bcd60e51b8152600401808060200182810382526024815260200180615cdd6024913960400191505060405180910390fd5b806003028311612eb0576040805162461bcd60e51b815260206004820181905260248201527f6661756c74792d6f7261636c65207468726573686f6c6420746f6f2068696768604482015290519081900360640190fd5b6000546001600160a01b03163314612f0f576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b6040805160c0810182528a8152602081018a905260ff89169181018290526060810188905267ffffffffffffffff8716608082015260a0810186905290612f56908861542d565b600654156130eb57600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81019160009183908110612f9357fe5b6000918252602082200154600780546001600160a01b0390921693509084908110612fba57fe5b60009182526020808320909101546001600160a01b0385811684526005909252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009081169091559290911680845292208054909116905560068054919250908061302757fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055600780548061308a57fe5b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550612f56915050565b60005b815151811015613492576000600560008460000151848151811061310e57fe5b6020908102919091018101516001600160a01b0316825281019190915260400160002054610100900460ff16600281111561314557fe5b14613197576040805162461bcd60e51b815260206004820152601760248201527f7265706561746564207369676e65722061646472657373000000000000000000604482015290519081900360640190fd5b6040805180820190915260ff821681526001602082015282518051600591600091859081106131c257fe5b6020908102919091018101516001600160a01b0316825281810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661010083600281111561324e57fe5b02179055506000915061325e9050565b600560008460200151848151811061327257fe5b6020908102919091018101516001600160a01b0316825281019190915260400160002054610100900460ff1660028111156132a957fe5b146132fb576040805162461bcd60e51b815260206004820152601c60248201527f7265706561746564207472616e736d6974746572206164647265737300000000604482015290519081900360640190fd5b6040805180820190915260ff82168152602081016002815250600560008460200151848151811061332857fe5b6020908102919091018101516001600160a01b0316825281810192909252604001600020825181547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff9091161780825591830151909182907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008360028111156133b457fe5b0217905550508251805160069250839081106133cc57fe5b602090810291909101810151825460018101845560009384529282902090920180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03909316929092179091558201518051600791908390811061343557fe5b60209081029190910181015182546001808201855560009485529290932090920180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0390931692909217909155016130ee565b5060408101516003805460ff83167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00909116179055600480544363ffffffff9081166401000000009081027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff84161780831660010183167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000009091161793849055855160208701516060880151608089015160a08a01519490960485169746976135649789973097921695949391615461565b60026000018190555050816000015151600260010160016101000a81548160ff021916908360ff1602179055507f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0581600260000154600460009054906101000a900463ffffffff16856000015186602001518760400151886060015189608001518a60a00151604051808a63ffffffff1681526020018981526020018863ffffffff16815260200180602001806020018760ff168152602001806020018667ffffffffffffffff1681526020018060200185810385528b818151815260200191508051906020019060200280838360005b8381101561366d578181015183820152602001613655565b5050505090500185810384528a818151815260200191508051906020019060200280838360005b838110156136ac578181015183820152602001613694565b50505050905001858103835288818151815260200191508051906020019080838360005b838110156136e85781810151838201526020016136d0565b50505050905090810190601f1680156137155780820380516001836020036101000a031916815260200191505b50858103825286518152865160209182019188019080838360005b83811015613748578181015183820152602001613730565b50505050905090810190601f1680156137755780820380516001836020036101000a031916815260200191505b509d505050505050505050505050505060405180910390a161379f826040015183606001516113bb565b5050505050505050505050565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff8082168452859484019161010090041660028111156137ee57fe5b60028111156137f957fe5b905250905060008160200151600281111561381057fe5b1415613820576000915050610ff1565b6001600b826000015160ff16601f811061383657fe5b601091828204019190066002029054906101000a900461ffff1603915050919050565b6000808080803332146138b3576040805162461bcd60e51b815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f41000000000000000000000000604482015290519081900360640190fd5b5050600254602e5463ffffffff65010000000000820481166000908152602f60205260409020549296600883901c909116955064ffffffffff9091169350601782900b9250780100000000000000000000000000000000000000000000000090910467ffffffffffffffff1690565b6009546001600160a01b031690565b6000546001600160a01b03163314613990576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b604080518082019091526030546001600160a01b038082168084527401000000000000000000000000000000000000000090920463ffffffff16602084015284161415806139ee57508163ffffffff16816020015163ffffffff1614155b15613ac1576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052603080547fffffffffffffffffffffffff00000000000000000000000000000000000000001684177fffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000008302179055865187860151875193168352948201528451919493909216927fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541928290030190a35b505050565b6001600160a01b038281166000908152600d6020526040902054163314613b34576040805162461bcd60e51b815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e20757064617465000000604482015290519081900360640190fd5b336001600160a01b0382161415613b92576040805162461bcd60e51b815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015290519081900360640190fd5b6001600160a01b038083166000908152600e6020526040902080548383167fffffffffffffffffffffffff000000000000000000000000000000000000000082168117909255909116908114613ac1576040516001600160a01b038084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b6000546001600160a01b03163314613c80576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000546001600160a01b03163314613d48576040805162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015290519081900360640190fd5b61199381615698565b6000806000806000613d9a336000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061150192505050565b613deb576040805162461bcd60e51b815260206004820152600960248201527f4e6f206163636573730000000000000000000000000000000000000000000000604482015290519081900360640190fd5b613df3615727565b945094509450945094509091929394565b6040805160a08101825260085463ffffffff80821683526401000000008204811660208401526801000000000000000082048116838501526c0100000000000000000000000082048116606084015270010000000000000000000000000000000090910416608082015260095482516103e081019384905291926001600160a01b0390911691600091600b90601f908285855b82829054906101000a900461ffff1661ffff1681526020019060020190602082600101049283019260010382029150808411613e97575050604080516103e08101918290529596506000959450600f9350601f9250905082845b815481526020019060010190808311613ef1575050505050905060006007805480602002602001604051908101604052809291908181526020018280548015613f6357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613f45575b5050505050905060005b815181101561419857600060018483601f8110613f8657fe5b6020020151039050600060018684601f8110613f9e57fe5b60200201510361ffff169050600082896060015163ffffffff168302633b9aca0002019050600081111561418d576000600d6000878781518110613fde57fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060009054906101000a90046001600160a01b03169050886001600160a01b031663a9059cbb82846040518363ffffffff1660e01b815260040180836001600160a01b0316815260200182815260200192505050602060405180830381600087803b15801561407357600080fd5b505af1158015614087573d6000803e3d6000fd5b505050506040513d602081101561409d57600080fd5b50516140f0576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b60018886601f81106140fe57fe5b61ffff909216602092909202015260018786601f811061411a57fe5b602002018181525050886001600160a01b0316816001600160a01b031687878151811061414357fe5b60200260200101516001600160a01b03167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c856040518082815260200191505060405180910390a4505b505050600101613f6d565b506141a6600b84601f615bdb565b50612985600f83601f615c71565b602e5465010000000000900463ffffffff166000908152602f6020526040902054601790810b900b90565b602e5465010000000000900463ffffffff1690565b6001600160a01b03821660009081526034602052604081205460ff168061152057505060335460ff161592915050565b60328054604080516020601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156117dd5780601f106142a3576101008083540402835291602001916117dd565b820191906000526020600020905b8154815290600101906020018083116142b157509395945050505050565b602e5465010000000000900463ffffffff166000908152602f60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff8082168452929391929184019161010090910416600281111561435f57fe5b600281111561436a57fe5b9052509050600061437a83610eab565b90508015613ac1576001600160a01b038084166000908152600d602090815260408083205460095482517fa9059cbb000000000000000000000000000000000000000000000000000000008152918616600483018190526024830188905292519295169363a9059cbb9360448084019491939192918390030190829087803b15801561440557600080fd5b505af1158015614419573d6000803e3d6000fd5b505050506040513d602081101561442f57600080fd5b5051614482576040805162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015290519081900360640190fd5b6001600b846000015160ff16601f811061449857fe5b601091828204019190066002026101000a81548161ffff021916908361ffff1602179055506001600f846000015160ff16601f81106144d357fe5b01556009546040805184815290516001600160a01b039283169284811692908816917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c9181900360200190a450505050565b600080600080600063ffffffff8669ffffffffffffffffffff1611156040518060400160405280600f81526020017f4e6f20646174612070726573656e740000000000000000000000000000000000815250906146005760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b838110156145c55781810151838201526020016145ad565b50505050905090810190601f1680156145f25780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5050505063ffffffff83166000908152602f6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052949594900b939092508291508490565b6001600160a01b03811660009081526034602052604090205460ff16611993576001600160a01b03811660008181526034602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055815192835290517f87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db49281900390910190a150565b60408051808201909152602e5464ffffffffff8082168084526501000000000090920463ffffffff16602084015284161161478b576040805162461bcd60e51b815260206004820152600c60248201527f7374616c65207265706f72740000000000000000000000000000000000000000604482015290519081900360640190fd5b600060606147988461579a565b90925090506147a78482615853565b600354815160ff90911690601f1015614807576040805162461bcd60e51b815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e64730000604482015290519081900360640190fd5b8060020282511161485f576040805162461bcd60e51b815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e0000604482015290519081900360640190fd5b6000825167ffffffffffffffff8111801561487957600080fd5b506040519080825280601f01601f1916602001820160405280156148a4576020820181803683370190505b5090506148af615ba5565b60005b845181101561499e5760008682602081106148c957fe5b1a90508281601f81106148d857fe5b60200201511561492f576040805162461bcd60e51b815260206004820152601760248201527f6f6273657276657220696e646578207265706561746564000000000000000000604482015290519081900360640190fd5b60018382601f811061493d57fe5b9115156020928302919091015287908390811061495657fe5b1a60f81b84838151811061496657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350506001016148b2565b506149a8826158b7565b64ffffffffff8816865260005b6001855103811015614a535760008582600101815181106149d257fe5b602002602001015160170b8683815181106149e957fe5b602002602001015160170b1315905080614a4a576040805162461bcd60e51b815260206004820152601760248201527f6f62736572766174696f6e73206e6f7420736f72746564000000000000000000604482015290519081900360640190fd5b506001016149b5565b506000846002865181614a6257fe5b0481518110614a6d57fe5b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b13158015614ad357507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b614b24576040805162461bcd60e51b815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e67650000604482015290519081900360640190fd5b86602001805180919060010163ffffffff1663ffffffff168152505060405180604001604052808260170b81526020014267ffffffffffffffff16815250602f6000896020015163ffffffff1663ffffffff16815260200190815260200160002060008201518160000160006101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908360170b77ffffffffffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160186101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550905050866020015163ffffffff167f8235efcbf95cfe12e2d5afec1e5e568dc529cb92d6a9b4195da079f1411244f8823388878f8f604051808760170b8152602001866001600160a01b0316815260200180602001806020018581526020018464ffffffffff168152602001838103835287818151815260200191508051906020019060200280838360005b83811015614caa578181015183820152602001614c92565b50505050905001838103825286818151815260200191508051906020019080838360005b83811015614ce6578181015183820152602001614cce565b50505050905090810190601f168015614d135780820380516001836020036101000a031916815260200191505b509850505050505050505060405180910390a260208088015160408051428152905160009363ffffffff909316927f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271928290030190a3866020015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040518082815260200191505060405180910390a3614dc087602001518260170b615926565b50508451602e805460209097015163ffffffff1665010000000000027fffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffffff64ffffffffff9093167fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000090981697909717919091169590951790945550505050505050565b602083810286019082020161014401368114614ea5576040805162461bcd60e51b815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d617463680000000000000000604482015290519081900360640190fd5b50505050505050565b6001600160a01b03811660009081526005602090815260408083208151808301909252805460ff80821684529293919291840191610100909104166002811115614ef457fe5b6002811115614eff57fe5b9052509050600281602001516002811115614f1657fe5b14614f2057600080fd5b6040805160a08101825260085463ffffffff80821680845264010000000083048216602085018190526801000000000000000084048316958501959095526c010000000000000000000000008304821660608501527001000000000000000000000000000000009092041660808301529091600091614fa691633b9aca003a0491615a66565b90506010360260005a90506000614fc58863ffffffff16858585615a8c565b6fffffffffffffffffffffffffffffffff1690506000620f4240866040015163ffffffff16830281614ff357fe5b049050856080015163ffffffff16633b9aca000281600f896000015160ff16601f811061501c57fe5b01540101600f886000015160ff16601f811061503457fe5b0155505050505050505050565b600063ffffffff82111561505757506000610ff1565b5063ffffffff166000908152602f6020526040902054601790810b900b90565b600063ffffffff82111561508d57506000610ff1565b5063ffffffff166000908152602f60205260409020547801000000000000000000000000000000000000000000000000900467ffffffffffffffff1690565b6040805160a0808201835263ffffffff88811680845288821660208086018290528984168688018190528985166060808901829052958a166080988901819052600880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001687177fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff166401000000008702177fffffffffffffffffffffffffffffffffffffffff00000000ffffffffffffffff16680100000000000000008502177fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff166c010000000000000000000000008402177fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000830217905589519586529285019390935283880152928201529283015291517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6929181900390910190a15050505050565b604080516103e0810191829052600091829190600b90601f908285855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116152635790505050505050905060005b601f8110156152d35760018282601f81106152bc57fe5b60200201510361ffff1692909201916001016152a5565b506040805160a08101825260085463ffffffff8082168352640100000000820481166020808501919091526801000000000000000083048216848601526c01000000000000000000000000830482166060850181905270010000000000000000000000000000000090930490911660808401526007805485518184028101840190965280865296909202633b9aca00029592936000939092918301828280156153a557602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311615387575b5050604080516103e08101918290529495506000949350600f9250601f915082845b8154815260200190600101908083116153c7575050505050905060005b825181101561540e5760018282601f81106153fb57fe5b60200201510395909501946001016153e4565b505050505090565b600081831015615427575081611523565b50919050565b615435613e04565b5050602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000169055565b6000808a8a8a8a8a8a8a8a8a604051602001808a8152602001896001600160a01b031681526020018867ffffffffffffffff16815260200180602001806020018760ff168152602001806020018667ffffffffffffffff1681526020018060200185810385528b818151815260200191508051906020019060200280838360005b838110156154fa5781810151838201526020016154e2565b5050505090500185810384528a818151815260200191508051906020019060200280838360005b83811015615539578181015183820152602001615521565b50505050905001858103835288818151815260200191508051906020019080838360005b8381101561557557818101518382015260200161555d565b50505050905090810190601f1680156155a25780820380516001836020036101000a031916815260200191505b50858103825286518152865160209182019188019080838360005b838110156155d55781810151838201526020016155bd565b50505050905090810190601f1680156156025780820380516001836020036101000a031916815260200191505b50604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179f505050505050505050505050505050509998505050505050505050565b600a546001600160a01b0390811690821681146113bb57600a80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03848116918217909255604080519284168352602083019190915280517f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d489129281900390910190a15050565b602e5465010000000000900463ffffffff166000818152602f6020908152604091829020825180840190935254601781810b810b810b808552780100000000000000000000000000000000000000000000000090920467ffffffffffffffff1693909201839052929392900b9181908490565b600060608280602001905160408110156157b357600080fd5b8151602083018051604051929492938301929190846401000000008211156157da57600080fd5b9083019060208201858111156157ef57600080fd5b825186602082028301116401000000008211171561580c57600080fd5b82525081516020918201928201910280838360005b83811015615839578181015183820152602001615821565b505050509050016040525050508092508193505050915091565b60008151602002606001600001905080835114613ac1576040805162461bcd60e51b815260206004820152601660248201527f7265706f7274206c656e677468206d69736d6174636800000000000000000000604482015290519081900360640190fd5b604080516103e081019182905261591891839190600b90601f90826000855b82829054906101000a900461ffff1661ffff16815260200190600201906020826001010492830192600103820291508084116158d65790505050505050615b18565b6113bb90600b90601f615bdb565b604080518082019091526030546001600160a01b0381168083527401000000000000000000000000000000000000000090910463ffffffff16602083015261596e57506113bb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff830163ffffffff8181166000818152602f602090815260408083205487518884015183517fbeed9b510000000000000000000000000000000000000000000000000000000081526004810197909752601792830b90920b602487018190528b88166044880152606487018b9052925192966001600160a01b039091169563beed9b51959290911693608480830194919391928390030190829088803b158015615a3757600080fd5b5087f193505050508015615a5d57506040513d6020811015615a5857600080fd5b505160015b61298557611f31565b60008383811015615a7957600285850304015b615a838184615416565b95945050505050565b600081851015615ae3576040805162461bcd60e51b815260206004820181905260248201527f6761734c6566742063616e6e6f742065786365656420696e697469616c476173604482015290519081900360640190fd5b818503830161179301633b9aca00858202026fffffffffffffffffffffffffffffffff8110615b0e57fe5b9695505050505050565b615b20615ba5565b60005b8351811015615b85576000848281518110615b3a57fe5b016020015160f81c9050615b5f8482601f8110615b5357fe5b60200201516001615b8d565b848260ff16601f8110615b6e57fe5b61ffff909216602092909202015250600101615b23565b509092915050565b60006115208261ffff168461ffff160161ffff615416565b604051806103e00160405280601f906020820280368337509192915050565b604080518082019091526000808252602082015290565b600283019183908215615c615791602002820160005b83821115615c3157835183826101000a81548161ffff021916908361ffff1602179055509260200192600201602081600101049283019260010302615bf1565b8015615c5f5782816101000a81549061ffff0219169055600201602081600101049283019260010302615c31565b505b50615c6d929150615c9f565b5090565b82601f8101928215615c61579160200282015b82811115615c61578251825591602001919060010190615c84565b5b80821115615c6d5760008155600101615ca056fe416363657373436f6e74726f6c6c65644f6666636861696e41676772656761746f7220332e302e306f7261636c6520616464726573736573206f7574206f6620726567697374726174696f6ea164736f6c6343000706000a"

func DeployAccessControlledOffchainAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32, _link common.Address, _minAnswer *big.Int, _maxAnswer *big.Int, _billingAccessController common.Address, _requesterAccessController common.Address, _decimals uint8, description string) (common.Address, *types.Transaction, *AccessControlledOffchainAggregator, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledOffchainAggregatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccessControlledOffchainAggregatorBin), backend, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission, _link, _minAnswer, _maxAnswer, _billingAccessController, _requesterAccessController, _decimals, description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessControlledOffchainAggregator{AccessControlledOffchainAggregatorCaller: AccessControlledOffchainAggregatorCaller{contract: contract}, AccessControlledOffchainAggregatorTransactor: AccessControlledOffchainAggregatorTransactor{contract: contract}, AccessControlledOffchainAggregatorFilterer: AccessControlledOffchainAggregatorFilterer{contract: contract}}, nil
}

type AccessControlledOffchainAggregator struct {
	address common.Address
	abi     abi.ABI
	AccessControlledOffchainAggregatorCaller
	AccessControlledOffchainAggregatorTransactor
	AccessControlledOffchainAggregatorFilterer
}

type AccessControlledOffchainAggregatorCaller struct {
	contract *bind.BoundContract
}

type AccessControlledOffchainAggregatorTransactor struct {
	contract *bind.BoundContract
}

type AccessControlledOffchainAggregatorFilterer struct {
	contract *bind.BoundContract
}

type AccessControlledOffchainAggregatorSession struct {
	Contract     *AccessControlledOffchainAggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AccessControlledOffchainAggregatorCallerSession struct {
	Contract *AccessControlledOffchainAggregatorCaller
	CallOpts bind.CallOpts
}

type AccessControlledOffchainAggregatorTransactorSession struct {
	Contract     *AccessControlledOffchainAggregatorTransactor
	TransactOpts bind.TransactOpts
}

type AccessControlledOffchainAggregatorRaw struct {
	Contract *AccessControlledOffchainAggregator
}

type AccessControlledOffchainAggregatorCallerRaw struct {
	Contract *AccessControlledOffchainAggregatorCaller
}

type AccessControlledOffchainAggregatorTransactorRaw struct {
	Contract *AccessControlledOffchainAggregatorTransactor
}

func NewAccessControlledOffchainAggregator(address common.Address, backend bind.ContractBackend) (*AccessControlledOffchainAggregator, error) {
	abi, err := abi.JSON(strings.NewReader(AccessControlledOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAccessControlledOffchainAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregator{address: address, abi: abi, AccessControlledOffchainAggregatorCaller: AccessControlledOffchainAggregatorCaller{contract: contract}, AccessControlledOffchainAggregatorTransactor: AccessControlledOffchainAggregatorTransactor{contract: contract}, AccessControlledOffchainAggregatorFilterer: AccessControlledOffchainAggregatorFilterer{contract: contract}}, nil
}

func NewAccessControlledOffchainAggregatorCaller(address common.Address, caller bind.ContractCaller) (*AccessControlledOffchainAggregatorCaller, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCaller{contract: contract}, nil
}

func NewAccessControlledOffchainAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControlledOffchainAggregatorTransactor, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorTransactor{contract: contract}, nil
}

func NewAccessControlledOffchainAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControlledOffchainAggregatorFilterer, error) {
	contract, err := bindAccessControlledOffchainAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorFilterer{contract: contract}, nil
}

func bindAccessControlledOffchainAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlledOffchainAggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transfer(opts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AccessControlledOffchainAggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessControlledOffchainAggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transfer(opts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.contract.Transact(opts, method, params...)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) BillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "billingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) BillingAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.BillingAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) BillingAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.BillingAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) CheckEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "checkEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) CheckEnabled() (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.CheckEnabled(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Decimals() (uint8, error) {
	return _AccessControlledOffchainAggregator.Contract.Decimals(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Description() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.Description(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getAnswer", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetAnswer(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetAnswer(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetAnswer(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetBilling(opts *bind.CallOpts) (GetBilling,

	error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getBilling")

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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetBilling() (GetBilling,

	error) {
	return _AccessControlledOffchainAggregator.Contract.GetBilling(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetBilling() (GetBilling,

	error) {
	return _AccessControlledOffchainAggregator.Contract.GetBilling(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetLinkToken() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.GetLinkToken(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetLinkToken() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.GetLinkToken(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetRoundData(opts *bind.CallOpts, _roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getRoundData", _roundId)

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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _AccessControlledOffchainAggregator.Contract.GetRoundData(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetRoundData(_roundId *big.Int) (GetRoundData,

	error) {
	return _AccessControlledOffchainAggregator.Contract.GetRoundData(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) GetTimestamp(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "getTimestamp", _roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) GetTimestamp(_roundId *big.Int) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.GetTimestamp(&_AccessControlledOffchainAggregator.CallOpts, _roundId)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) HasAccess(opts *bind.CallOpts, _user common.Address, _calldata []byte) (bool, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "hasAccess", _user, _calldata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) HasAccess(_user common.Address, _calldata []byte) (bool, error) {
	return _AccessControlledOffchainAggregator.Contract.HasAccess(&_AccessControlledOffchainAggregator.CallOpts, _user, _calldata)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AccessControlledOffchainAggregator.Contract.LatestConfigDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _AccessControlledOffchainAggregator.Contract.LatestConfigDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestRound() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRound(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRound(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestRoundData")

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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRoundData(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _AccessControlledOffchainAggregator.Contract.LatestRoundData(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTimestamp(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTimestamp(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (LatestTransmissionDetails,

	error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "latestTransmissionDetails")

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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTransmissionDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _AccessControlledOffchainAggregator.Contract.LatestTransmissionDetails(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.LinkAvailableForPayment(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MaxAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.MinAnswer(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OracleObservationCount(opts *bind.CallOpts, _signerOrTransmitter common.Address) (uint16, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "oracleObservationCount", _signerOrTransmitter)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OracleObservationCount(_signerOrTransmitter common.Address) (uint16, error) {
	return _AccessControlledOffchainAggregator.Contract.OracleObservationCount(&_AccessControlledOffchainAggregator.CallOpts, _signerOrTransmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) OwedPayment(opts *bind.CallOpts, _transmitter common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "owedPayment", _transmitter)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) OwedPayment(_transmitter common.Address) (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.OwedPayment(&_AccessControlledOffchainAggregator.CallOpts, _transmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Owner() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Owner(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) RequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "requesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) RequesterAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.RequesterAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) RequesterAccessController() (common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.RequesterAccessController(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Transmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "transmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Transmitters() ([]common.Address, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmitters(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) TypeAndVersion() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.TypeAndVersion(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) TypeAndVersion() (string, error) {
	return _AccessControlledOffchainAggregator.Contract.TypeAndVersion(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) ValidatorConfig(opts *bind.CallOpts) (ValidatorConfig,

	error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "validatorConfig")

	outstruct := new(ValidatorConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.GasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) ValidatorConfig() (ValidatorConfig,

	error) {
	return _AccessControlledOffchainAggregator.Contract.ValidatorConfig(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) ValidatorConfig() (ValidatorConfig,

	error) {
	return _AccessControlledOffchainAggregator.Contract.ValidatorConfig(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessControlledOffchainAggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Version() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.Version(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorCallerSession) Version() (*big.Int, error) {
	return _AccessControlledOffchainAggregator.Contract.Version(&_AccessControlledOffchainAggregator.CallOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "acceptOwnership")
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptOwnership(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptOwnership(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "acceptPayeeship", _transmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AcceptPayeeship(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AcceptPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) AddAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "addAccess", _user)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AddAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) AddAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.AddAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) DisableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "disableAccessCheck")
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.DisableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) DisableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.DisableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) EnableAccessCheck(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "enableAccessCheck")
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.EnableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) EnableAccessCheck() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.EnableAccessCheck(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) RemoveAccess(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "removeAccess", _user)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RemoveAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) RemoveAccess(_user common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RemoveAccess(&_AccessControlledOffchainAggregator.TransactOpts, _user)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "requestNewRound")
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RequestNewRound(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.RequestNewRound(&_AccessControlledOffchainAggregator.TransactOpts)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetBilling(opts *bind.TransactOpts, _maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setBilling", _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBilling(&_AccessControlledOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetBilling(_maximumGasPrice uint32, _reasonableGasPrice uint32, _microLinkPerEth uint32, _linkGweiPerObservation uint32, _linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBilling(&_AccessControlledOffchainAggregator.TransactOpts, _maximumGasPrice, _reasonableGasPrice, _microLinkPerEth, _linkGweiPerObservation, _linkGweiPerTransmission)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBillingAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _billingAccessController)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetBillingAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _billingAccessController)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setConfig", _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetConfig(&_AccessControlledOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetConfig(_signers []common.Address, _transmitters []common.Address, _threshold uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetConfig(&_AccessControlledOffchainAggregator.TransactOpts, _signers, _transmitters, _threshold, _onchainConfig, _offchainConfigVersion, _offchainConfig)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetLinkToken(opts *bind.TransactOpts, _linkToken common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setLinkToken", _linkToken, _recipient)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetLinkToken(_linkToken common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetLinkToken(&_AccessControlledOffchainAggregator.TransactOpts, _linkToken, _recipient)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetLinkToken(_linkToken common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetLinkToken(&_AccessControlledOffchainAggregator.TransactOpts, _linkToken, _recipient)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetPayees(opts *bind.TransactOpts, _transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setPayees", _transmitters, _payees)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetPayees(&_AccessControlledOffchainAggregator.TransactOpts, _transmitters, _payees)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetPayees(_transmitters []common.Address, _payees []common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetPayees(&_AccessControlledOffchainAggregator.TransactOpts, _transmitters, _payees)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, _requesterAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setRequesterAccessController", _requesterAccessController)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetRequesterAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _requesterAccessController)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetRequesterAccessController(_requesterAccessController common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetRequesterAccessController(&_AccessControlledOffchainAggregator.TransactOpts, _requesterAccessController)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) SetValidatorConfig(opts *bind.TransactOpts, _newValidator common.Address, _newGasLimit uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "setValidatorConfig", _newValidator, _newGasLimit)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) SetValidatorConfig(_newValidator common.Address, _newGasLimit uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetValidatorConfig(&_AccessControlledOffchainAggregator.TransactOpts, _newValidator, _newGasLimit)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) SetValidatorConfig(_newValidator common.Address, _newGasLimit uint32) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.SetValidatorConfig(&_AccessControlledOffchainAggregator.TransactOpts, _newValidator, _newGasLimit)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transferOwnership", _to)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferOwnership(&_AccessControlledOffchainAggregator.TransactOpts, _to)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) TransferOwnership(_to common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferOwnership(&_AccessControlledOffchainAggregator.TransactOpts, _to)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transferPayeeship", _transmitter, _proposed)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter, _proposed)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) TransferPayeeship(_transmitter common.Address, _proposed common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.TransferPayeeship(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter, _proposed)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmit(&_AccessControlledOffchainAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.Transmit(&_AccessControlledOffchainAggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "withdrawFunds", _recipient, _amount)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawFunds(&_AccessControlledOffchainAggregator.TransactOpts, _recipient, _amount)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) WithdrawFunds(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawFunds(&_AccessControlledOffchainAggregator.TransactOpts, _recipient, _amount)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.contract.Transact(opts, "withdrawPayment", _transmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawPayment(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorTransactorSession) WithdrawPayment(_transmitter common.Address) (*types.Transaction, error) {
	return _AccessControlledOffchainAggregator.Contract.WithdrawPayment(&_AccessControlledOffchainAggregator.TransactOpts, _transmitter)
}

type AccessControlledOffchainAggregatorAddedAccessIterator struct {
	Event *AccessControlledOffchainAggregatorAddedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorAddedAccess)
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
		it.Event = new(AccessControlledOffchainAggregatorAddedAccess)
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

func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorAddedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorAddedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterAddedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorAddedAccessIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorAddedAccessIterator{contract: _AccessControlledOffchainAggregator.contract, event: "AddedAccess", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorAddedAccess) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "AddedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorAddedAccess)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseAddedAccess(log types.Log) (*AccessControlledOffchainAggregatorAddedAccess, error) {
	event := new(AccessControlledOffchainAggregatorAddedAccess)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AddedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorAnswerUpdatedIterator struct {
	Event *AccessControlledOffchainAggregatorAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorAnswerUpdated)
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
		it.Event = new(AccessControlledOffchainAggregatorAnswerUpdated)
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

func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AccessControlledOffchainAggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorAnswerUpdatedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorAnswerUpdated)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseAnswerUpdated(log types.Log) (*AccessControlledOffchainAggregatorAnswerUpdated, error) {
	event := new(AccessControlledOffchainAggregatorAnswerUpdated)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorBillingAccessControllerSetIterator struct {
	Event *AccessControlledOffchainAggregatorBillingAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
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
		it.Event = new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
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

func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorBillingAccessControllerSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*AccessControlledOffchainAggregatorBillingAccessControllerSet, error) {
	event := new(AccessControlledOffchainAggregatorBillingAccessControllerSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorBillingSetIterator struct {
	Event *AccessControlledOffchainAggregatorBillingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorBillingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorBillingSet)
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
		it.Event = new(AccessControlledOffchainAggregatorBillingSet)
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

func (it *AccessControlledOffchainAggregatorBillingSetIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorBillingSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorBillingSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseBillingSet(log types.Log) (*AccessControlledOffchainAggregatorBillingSet, error) {
	event := new(AccessControlledOffchainAggregatorBillingSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorCheckAccessDisabledIterator struct {
	Event *AccessControlledOffchainAggregatorCheckAccessDisabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorCheckAccessDisabled)
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
		it.Event = new(AccessControlledOffchainAggregatorCheckAccessDisabled)
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

func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorCheckAccessDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorCheckAccessDisabled struct {
	Raw types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterCheckAccessDisabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessDisabledIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCheckAccessDisabledIterator{contract: _AccessControlledOffchainAggregator.contract, event: "CheckAccessDisabled", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorCheckAccessDisabled) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "CheckAccessDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorCheckAccessDisabled)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseCheckAccessDisabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessDisabled, error) {
	event := new(AccessControlledOffchainAggregatorCheckAccessDisabled)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorCheckAccessEnabledIterator struct {
	Event *AccessControlledOffchainAggregatorCheckAccessEnabled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorCheckAccessEnabled)
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
		it.Event = new(AccessControlledOffchainAggregatorCheckAccessEnabled)
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

func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorCheckAccessEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorCheckAccessEnabled struct {
	Raw types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterCheckAccessEnabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessEnabledIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorCheckAccessEnabledIterator{contract: _AccessControlledOffchainAggregator.contract, event: "CheckAccessEnabled", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorCheckAccessEnabled) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "CheckAccessEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorCheckAccessEnabled)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseCheckAccessEnabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessEnabled, error) {
	event := new(AccessControlledOffchainAggregatorCheckAccessEnabled)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "CheckAccessEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorConfigSetIterator struct {
	Event *AccessControlledOffchainAggregatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorConfigSet)
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
		it.Event = new(AccessControlledOffchainAggregatorConfigSet)
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

func (it *AccessControlledOffchainAggregatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorConfigSet struct {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorConfigSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorConfigSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorConfigSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseConfigSet(log types.Log) (*AccessControlledOffchainAggregatorConfigSet, error) {
	event := new(AccessControlledOffchainAggregatorConfigSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorLinkTokenSetIterator struct {
	Event *AccessControlledOffchainAggregatorLinkTokenSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorLinkTokenSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorLinkTokenSet)
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
		it.Event = new(AccessControlledOffchainAggregatorLinkTokenSet)
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

func (it *AccessControlledOffchainAggregatorLinkTokenSetIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, _oldLinkToken []common.Address, _newLinkToken []common.Address) (*AccessControlledOffchainAggregatorLinkTokenSetIterator, error) {

	var _oldLinkTokenRule []interface{}
	for _, _oldLinkTokenItem := range _oldLinkToken {
		_oldLinkTokenRule = append(_oldLinkTokenRule, _oldLinkTokenItem)
	}
	var _newLinkTokenRule []interface{}
	for _, _newLinkTokenItem := range _newLinkToken {
		_newLinkTokenRule = append(_newLinkTokenRule, _newLinkTokenItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "LinkTokenSet", _oldLinkTokenRule, _newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorLinkTokenSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorLinkTokenSet, _oldLinkToken []common.Address, _newLinkToken []common.Address) (event.Subscription, error) {

	var _oldLinkTokenRule []interface{}
	for _, _oldLinkTokenItem := range _oldLinkToken {
		_oldLinkTokenRule = append(_oldLinkTokenRule, _oldLinkTokenItem)
	}
	var _newLinkTokenRule []interface{}
	for _, _newLinkTokenItem := range _newLinkToken {
		_newLinkTokenRule = append(_newLinkTokenRule, _newLinkTokenItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "LinkTokenSet", _oldLinkTokenRule, _newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorLinkTokenSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseLinkTokenSet(log types.Log) (*AccessControlledOffchainAggregatorLinkTokenSet, error) {
	event := new(AccessControlledOffchainAggregatorLinkTokenSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorNewRoundIterator struct {
	Event *AccessControlledOffchainAggregatorNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorNewRound)
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
		it.Event = new(AccessControlledOffchainAggregatorNewRound)
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

func (it *AccessControlledOffchainAggregatorNewRoundIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AccessControlledOffchainAggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorNewRoundIterator{contract: _AccessControlledOffchainAggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorNewRound)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseNewRound(log types.Log) (*AccessControlledOffchainAggregatorNewRound, error) {
	event := new(AccessControlledOffchainAggregatorNewRound)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorNewTransmissionIterator struct {
	Event *AccessControlledOffchainAggregatorNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorNewTransmission)
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
		it.Event = new(AccessControlledOffchainAggregatorNewTransmission)
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

func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	ConfigDigest      [32]byte
	EpochAndRound     *big.Int
	Raw               types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*AccessControlledOffchainAggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorNewTransmissionIterator{contract: _AccessControlledOffchainAggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorNewTransmission)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseNewTransmission(log types.Log) (*AccessControlledOffchainAggregatorNewTransmission, error) {
	event := new(AccessControlledOffchainAggregatorNewTransmission)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorOraclePaidIterator struct {
	Event *AccessControlledOffchainAggregatorOraclePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorOraclePaid)
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
		it.Event = new(AccessControlledOffchainAggregatorOraclePaid)
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

func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*AccessControlledOffchainAggregatorOraclePaidIterator, error) {

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

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorOraclePaidIterator{contract: _AccessControlledOffchainAggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorOraclePaid)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOraclePaid(log types.Log) (*AccessControlledOffchainAggregatorOraclePaid, error) {
	event := new(AccessControlledOffchainAggregatorOraclePaid)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
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
		it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
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

func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferRequested, error) {
	event := new(AccessControlledOffchainAggregatorOwnershipTransferRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorOwnershipTransferredIterator struct {
	Event *AccessControlledOffchainAggregatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferred)
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
		it.Event = new(AccessControlledOffchainAggregatorOwnershipTransferred)
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

func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledOffchainAggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorOwnershipTransferredIterator{contract: _AccessControlledOffchainAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorOwnershipTransferred)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferred, error) {
	event := new(AccessControlledOffchainAggregatorOwnershipTransferred)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
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
		it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
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

func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferRequested, error) {
	event := new(AccessControlledOffchainAggregatorPayeeshipTransferRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorPayeeshipTransferredIterator struct {
	Event *AccessControlledOffchainAggregatorPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferred)
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
		it.Event = new(AccessControlledOffchainAggregatorPayeeshipTransferred)
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

func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*AccessControlledOffchainAggregatorPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorPayeeshipTransferredIterator{contract: _AccessControlledOffchainAggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorPayeeshipTransferred)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferred, error) {
	event := new(AccessControlledOffchainAggregatorPayeeshipTransferred)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorRemovedAccessIterator struct {
	Event *AccessControlledOffchainAggregatorRemovedAccess

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorRemovedAccess)
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
		it.Event = new(AccessControlledOffchainAggregatorRemovedAccess)
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

func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorRemovedAccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorRemovedAccess struct {
	User common.Address
	Raw  types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterRemovedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorRemovedAccessIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorRemovedAccessIterator{contract: _AccessControlledOffchainAggregator.contract, event: "RemovedAccess", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRemovedAccess) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "RemovedAccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorRemovedAccess)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseRemovedAccess(log types.Log) (*AccessControlledOffchainAggregatorRemovedAccess, error) {
	event := new(AccessControlledOffchainAggregatorRemovedAccess)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RemovedAccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator struct {
	Event *AccessControlledOffchainAggregatorRequesterAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
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
		it.Event = new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
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

func (it *AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*AccessControlledOffchainAggregatorRequesterAccessControllerSet, error) {
	event := new(AccessControlledOffchainAggregatorRequesterAccessControllerSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorRoundRequestedIterator struct {
	Event *AccessControlledOffchainAggregatorRoundRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorRoundRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorRoundRequested)
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
		it.Event = new(AccessControlledOffchainAggregatorRoundRequested)
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

func (it *AccessControlledOffchainAggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [32]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*AccessControlledOffchainAggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorRoundRequestedIterator{contract: _AccessControlledOffchainAggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorRoundRequested)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseRoundRequested(log types.Log) (*AccessControlledOffchainAggregatorRoundRequested, error) {
	event := new(AccessControlledOffchainAggregatorRoundRequested)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type AccessControlledOffchainAggregatorValidatorConfigSetIterator struct {
	Event *AccessControlledOffchainAggregatorValidatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AccessControlledOffchainAggregatorValidatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlledOffchainAggregatorValidatorConfigSet)
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
		it.Event = new(AccessControlledOffchainAggregatorValidatorConfigSet)
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

func (it *AccessControlledOffchainAggregatorValidatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *AccessControlledOffchainAggregatorValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AccessControlledOffchainAggregatorValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*AccessControlledOffchainAggregatorValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlledOffchainAggregatorValidatorConfigSetIterator{contract: _AccessControlledOffchainAggregator.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _AccessControlledOffchainAggregator.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AccessControlledOffchainAggregatorValidatorConfigSet)
				if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregatorFilterer) ParseValidatorConfigSet(log types.Log) (*AccessControlledOffchainAggregatorValidatorConfigSet, error) {
	event := new(AccessControlledOffchainAggregatorValidatorConfigSet)
	if err := _AccessControlledOffchainAggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AccessControlledOffchainAggregator.abi.Events["AddedAccess"].ID:
		return _AccessControlledOffchainAggregator.ParseAddedAccess(log)
	case _AccessControlledOffchainAggregator.abi.Events["AnswerUpdated"].ID:
		return _AccessControlledOffchainAggregator.ParseAnswerUpdated(log)
	case _AccessControlledOffchainAggregator.abi.Events["BillingAccessControllerSet"].ID:
		return _AccessControlledOffchainAggregator.ParseBillingAccessControllerSet(log)
	case _AccessControlledOffchainAggregator.abi.Events["BillingSet"].ID:
		return _AccessControlledOffchainAggregator.ParseBillingSet(log)
	case _AccessControlledOffchainAggregator.abi.Events["CheckAccessDisabled"].ID:
		return _AccessControlledOffchainAggregator.ParseCheckAccessDisabled(log)
	case _AccessControlledOffchainAggregator.abi.Events["CheckAccessEnabled"].ID:
		return _AccessControlledOffchainAggregator.ParseCheckAccessEnabled(log)
	case _AccessControlledOffchainAggregator.abi.Events["ConfigSet"].ID:
		return _AccessControlledOffchainAggregator.ParseConfigSet(log)
	case _AccessControlledOffchainAggregator.abi.Events["LinkTokenSet"].ID:
		return _AccessControlledOffchainAggregator.ParseLinkTokenSet(log)
	case _AccessControlledOffchainAggregator.abi.Events["NewRound"].ID:
		return _AccessControlledOffchainAggregator.ParseNewRound(log)
	case _AccessControlledOffchainAggregator.abi.Events["NewTransmission"].ID:
		return _AccessControlledOffchainAggregator.ParseNewTransmission(log)
	case _AccessControlledOffchainAggregator.abi.Events["OraclePaid"].ID:
		return _AccessControlledOffchainAggregator.ParseOraclePaid(log)
	case _AccessControlledOffchainAggregator.abi.Events["OwnershipTransferRequested"].ID:
		return _AccessControlledOffchainAggregator.ParseOwnershipTransferRequested(log)
	case _AccessControlledOffchainAggregator.abi.Events["OwnershipTransferred"].ID:
		return _AccessControlledOffchainAggregator.ParseOwnershipTransferred(log)
	case _AccessControlledOffchainAggregator.abi.Events["PayeeshipTransferRequested"].ID:
		return _AccessControlledOffchainAggregator.ParsePayeeshipTransferRequested(log)
	case _AccessControlledOffchainAggregator.abi.Events["PayeeshipTransferred"].ID:
		return _AccessControlledOffchainAggregator.ParsePayeeshipTransferred(log)
	case _AccessControlledOffchainAggregator.abi.Events["RemovedAccess"].ID:
		return _AccessControlledOffchainAggregator.ParseRemovedAccess(log)
	case _AccessControlledOffchainAggregator.abi.Events["RequesterAccessControllerSet"].ID:
		return _AccessControlledOffchainAggregator.ParseRequesterAccessControllerSet(log)
	case _AccessControlledOffchainAggregator.abi.Events["RoundRequested"].ID:
		return _AccessControlledOffchainAggregator.ParseRoundRequested(log)
	case _AccessControlledOffchainAggregator.abi.Events["ValidatorConfigSet"].ID:
		return _AccessControlledOffchainAggregator.ParseValidatorConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AccessControlledOffchainAggregatorAddedAccess) Topic() common.Hash {
	return common.HexToHash("0x87286ad1f399c8e82bf0c4ef4fcdc570ea2e1e92176e5c848b6413545b885db4")
}

func (AccessControlledOffchainAggregatorAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (AccessControlledOffchainAggregatorBillingAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912")
}

func (AccessControlledOffchainAggregatorBillingSet) Topic() common.Hash {
	return common.HexToHash("0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6")
}

func (AccessControlledOffchainAggregatorCheckAccessDisabled) Topic() common.Hash {
	return common.HexToHash("0x3be8a977a014527b50ae38adda80b56911c267328965c98ddc385d248f539638")
}

func (AccessControlledOffchainAggregatorCheckAccessEnabled) Topic() common.Hash {
	return common.HexToHash("0xaebf329500988c6488a0074e5a0a9ff304561fc5c6fc877aeb1d59c8282c3480")
}

func (AccessControlledOffchainAggregatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (AccessControlledOffchainAggregatorLinkTokenSet) Topic() common.Hash {
	return common.HexToHash("0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a")
}

func (AccessControlledOffchainAggregatorNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (AccessControlledOffchainAggregatorNewTransmission) Topic() common.Hash {
	return common.HexToHash("0x8235efcbf95cfe12e2d5afec1e5e568dc529cb92d6a9b4195da079f1411244f8")
}

func (AccessControlledOffchainAggregatorOraclePaid) Topic() common.Hash {
	return common.HexToHash("0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c")
}

func (AccessControlledOffchainAggregatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (AccessControlledOffchainAggregatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (AccessControlledOffchainAggregatorPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (AccessControlledOffchainAggregatorPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (AccessControlledOffchainAggregatorRemovedAccess) Topic() common.Hash {
	return common.HexToHash("0x3d68a6fce901d20453d1a7aa06bf3950302a735948037deb182a8db66df2a0d1")
}

func (AccessControlledOffchainAggregatorRequesterAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634")
}

func (AccessControlledOffchainAggregatorRoundRequested) Topic() common.Hash {
	return common.HexToHash("0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c")
}

func (AccessControlledOffchainAggregatorValidatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541")
}

func (_AccessControlledOffchainAggregator *AccessControlledOffchainAggregator) Address() common.Address {
	return _AccessControlledOffchainAggregator.address
}

type AccessControlledOffchainAggregatorInterface interface {
	BillingAccessController(opts *bind.CallOpts) (common.Address, error)

	CheckEnabled(opts *bind.CallOpts) (bool, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, _roundId *big.Int) (*big.Int, error)

	GetBilling(opts *bind.CallOpts) (GetBilling,

		error)

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

	TransferOwnership(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, _transmitter common.Address, _proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, _transmitter common.Address) (*types.Transaction, error)

	FilterAddedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorAddedAccessIterator, error)

	WatchAddedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorAddedAccess) (event.Subscription, error)

	ParseAddedAccess(log types.Log) (*AccessControlledOffchainAggregatorAddedAccess, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*AccessControlledOffchainAggregatorAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*AccessControlledOffchainAggregatorAnswerUpdated, error)

	FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingAccessControllerSetIterator, error)

	WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorBillingAccessControllerSet) (event.Subscription, error)

	ParseBillingAccessControllerSet(log types.Log) (*AccessControlledOffchainAggregatorBillingAccessControllerSet, error)

	FilterBillingSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorBillingSetIterator, error)

	WatchBillingSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorBillingSet) (event.Subscription, error)

	ParseBillingSet(log types.Log) (*AccessControlledOffchainAggregatorBillingSet, error)

	FilterCheckAccessDisabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessDisabledIterator, error)

	WatchCheckAccessDisabled(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorCheckAccessDisabled) (event.Subscription, error)

	ParseCheckAccessDisabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessDisabled, error)

	FilterCheckAccessEnabled(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorCheckAccessEnabledIterator, error)

	WatchCheckAccessEnabled(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorCheckAccessEnabled) (event.Subscription, error)

	ParseCheckAccessEnabled(log types.Log) (*AccessControlledOffchainAggregatorCheckAccessEnabled, error)

	FilterConfigSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*AccessControlledOffchainAggregatorConfigSet, error)

	FilterLinkTokenSet(opts *bind.FilterOpts, _oldLinkToken []common.Address, _newLinkToken []common.Address) (*AccessControlledOffchainAggregatorLinkTokenSetIterator, error)

	WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorLinkTokenSet, _oldLinkToken []common.Address, _newLinkToken []common.Address) (event.Subscription, error)

	ParseLinkTokenSet(log types.Log) (*AccessControlledOffchainAggregatorLinkTokenSet, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*AccessControlledOffchainAggregatorNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*AccessControlledOffchainAggregatorNewRound, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*AccessControlledOffchainAggregatorNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*AccessControlledOffchainAggregatorNewTransmission, error)

	FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*AccessControlledOffchainAggregatorOraclePaidIterator, error)

	WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error)

	ParseOraclePaid(log types.Log) (*AccessControlledOffchainAggregatorOraclePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledOffchainAggregatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AccessControlledOffchainAggregatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*AccessControlledOffchainAggregatorOwnershipTransferred, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*AccessControlledOffchainAggregatorPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*AccessControlledOffchainAggregatorPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*AccessControlledOffchainAggregatorPayeeshipTransferred, error)

	FilterRemovedAccess(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorRemovedAccessIterator, error)

	WatchRemovedAccess(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRemovedAccess) (event.Subscription, error)

	ParseRemovedAccess(log types.Log) (*AccessControlledOffchainAggregatorRemovedAccess, error)

	FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*AccessControlledOffchainAggregatorRequesterAccessControllerSetIterator, error)

	WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRequesterAccessControllerSet) (event.Subscription, error)

	ParseRequesterAccessControllerSet(log types.Log) (*AccessControlledOffchainAggregatorRequesterAccessControllerSet, error)

	FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*AccessControlledOffchainAggregatorRoundRequestedIterator, error)

	WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorRoundRequested, requester []common.Address) (event.Subscription, error)

	ParseRoundRequested(log types.Log) (*AccessControlledOffchainAggregatorRoundRequested, error)

	FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*AccessControlledOffchainAggregatorValidatorConfigSetIterator, error)

	WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *AccessControlledOffchainAggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error)

	ParseValidatorConfigSet(log types.Log) (*AccessControlledOffchainAggregatorValidatorConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
