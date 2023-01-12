// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

var OCR2AggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"int192\",\"name\":\"minAnswer_\",\"type\":\"int192\"},{\"internalType\":\"int192\",\"name\":\"maxAnswer_\",\"type\":\"int192\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"billingAccessController\",\"type\":\"address\"},{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"requesterAccessController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"oldLinkToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"newLinkToken\",\"type\":\"address\"}],\"name\":\"LinkTokenSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"observationsTimestamp\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractAccessControllerInterface\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"name\":\"Transmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"ValidatorConfigSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"acceptPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBilling\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLinkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequesterAccessController\",\"outputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId_\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransmitters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorConfig\",\"outputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDetails\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"configCount\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestConfigDigestAndEpoch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"scanLogs\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer_\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp_\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkAvailableForPayment\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"availableBalance\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAnswer\",\"outputs\":[{\"internalType\":\"int192\",\"name\":\"\",\"type\":\"int192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"oracleObservationCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitterAddress\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestNewRound\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"\",\"type\":\"uint80\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceGwei\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"observationPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transmissionPaymentGjuels\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"accountingGas\",\"type\":\"uint24\"}],\"name\":\"setBilling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"_billingAccessController\",\"type\":\"address\"}],\"name\":\"setBillingAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"setLinkToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"payees\",\"type\":\"address[]\"}],\"name\":\"setPayees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAccessControllerInterface\",\"name\":\"requesterAccessController\",\"type\":\"address\"}],\"name\":\"setRequesterAccessController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAggregatorValidatorInterface\",\"name\":\"newValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newGasLimit\",\"type\":\"uint32\"}],\"name\":\"setValidatorConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"transferPayeeship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[3]\",\"name\":\"reportContext\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b50604051620060c1380380620060c1833981016040819052620000349162000537565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be8162000194565b5050601180546001600160a01b0319166001600160a01b038a169081179091556040519091506000907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a908290a3620001178462000240565b7fff0000000000000000000000000000000000000000000000000000000000000060f883901b1660c05280516200015690601090602084019062000479565b506200016283620002b9565b6200016f60008062000334565b50505050601791820b820b604090811b60805290820b90910b901b60a0525062000701565b6001600160a01b038116331415620001ef5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6012546001600160a01b039081169082168114620002b557601280546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291015b60405180910390a15b5050565b620002c36200041b565b600f546001600160a01b039081169082168114620002b557600f80546001600160a01b0319166001600160a01b0384811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae6349101620002ac565b6200033e6200041b565b60408051808201909152600e546001600160a01b03808216808452600160a01b90920463ffffffff16602084015284161415806200038c57508163ffffffff16816020015163ffffffff1614155b1562000416576040805180820182526001600160a01b0385811680835263ffffffff8681166020948501819052600e80546001600160c01b0319168417600160a01b830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b6000546001600160a01b03163314620004775760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640162000082565b565b828054620004879062000695565b90600052602060002090601f016020900481019282620004ab5760008555620004f6565b82601f10620004c657805160ff1916838001178555620004f6565b82800160010185558215620004f6579182015b82811115620004f6578251825591602001919060010190620004d9565b506200050492915062000508565b5090565b5b8082111562000504576000815560010162000509565b8051601781900b81146200053257600080fd5b919050565b600080600080600080600060e0888a0312156200055357600080fd5b87516200056081620006e8565b96506020620005718982016200051f565b96506200058160408a016200051f565b955060608901516200059381620006e8565b60808a0151909550620005a681620006e8565b60a08a015190945060ff81168114620005be57600080fd5b60c08a01519093506001600160401b0380821115620005dc57600080fd5b818b0191508b601f830112620005f157600080fd5b815181811115620006065762000606620006d2565b604051601f8201601f19908116603f01168101908382118183101715620006315762000631620006d2565b816040528281528e868487010111156200064a57600080fd5b600093505b828410156200066e57848401860151818501870152928501926200064f565b82841115620006805760008684830101525b80965050505050505092959891949750929550565b600181811c90821680620006aa57607f821691505b60208210811415620006cc57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114620006fe57600080fd5b50565b60805160401c60a05160401c60c05160f81c61596b62000756600039600061044901526000818161050d0152818161278d0152613d2c015260008181610345015281816127600152613cff015261596b6000f3fe608060405234801561001057600080fd5b50600436106102d35760003560e01c80639bd2c0b111610186578063d09dc339116100e3578063e76d516811610097578063f2fde38b11610071578063f2fde38b1461089a578063fbffd2c1146108ad578063feaf968c146108c057600080fd5b8063e76d516814610856578063eb45716314610874578063eb5dcd6c1461088757600080fd5b8063e3d0e712116100c8578063e3d0e712146107d1578063e4902f82146107e4578063e5fe45771461080c57600080fd5b8063d09dc339146107ab578063daffc4b5146107b357600080fd5b8063b1dc65a41161013a578063b633620c1161011f578063b633620c14610767578063c10753291461077a578063c4c92b371461078d57600080fd5b8063b1dc65a414610741578063b5ab58dc1461075457600080fd5b80639e3ceeab1161016b5780639e3ceeab146106ea578063afcb95d7146106fd578063b121e1471461072e57600080fd5b80639bd2c0b1146106785780639c849b30146106d757600080fd5b8063668a0f02116102345780638205bf6a116101e85780638da5cb5b116101cd5780638da5cb5b146105cc57806398e5b12a1461060b5780639a6fc8f51461062e57600080fd5b80638205bf6a1461056f5780638ac28d5a146105b957600080fd5b80637284e416116102195780637284e4161461052f57806379ba50971461053757806381ff70481461053f57600080fd5b8063668a0f02146104f057806370da2f671461050857600080fd5b80634fb174701161028b57806354fd4d501161027057806354fd4d50146104c0578063643dc105146104c8578063666cab8d146104db57600080fd5b80634fb174701461047d57806350d25bcd1461049257600080fd5b806322adbc78116102bc57806322adbc7814610340578063299372681461037a578063313ce5671461044457600080fd5b80630eafb25b146102d8578063181f5a77146102fe575b600080fd5b6102eb6102e6366004614dd6565b610959565b6040519081526020015b60405180910390f35b60408051808201909152601a81527f4f43523241676772656761746f7220312e302e302d616c70686100000000000060208201525b6040516102f59190615485565b6103677f000000000000000000000000000000000000000000000000000000000000000081565b60405160179190910b81526020016102f5565b610408600b546a0100000000000000000000810463ffffffff908116926e010000000000000000000000000000830482169272010000000000000000000000000000000000008104831692760100000000000000000000000000000000000000000000820416917a01000000000000000000000000000000000000000000000000000090910462ffffff1690565b6040805163ffffffff9687168152948616602086015292851692840192909252909216606082015262ffffff909116608082015260a0016102f5565b61046b7f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff90911681526020016102f5565b61049061048b366004614df3565b610a87565b005b600b546601000000000000900463ffffffff166000908152600c6020526040902054601790810b900b6102eb565b6102eb600681565b6104906104d63660046151c8565b610d8f565b6104e36110b4565b6040516102f591906153a3565b600b546601000000000000900463ffffffff166102eb565b6103677f000000000000000000000000000000000000000000000000000000000000000081565b610333611123565b6104906111ac565b600d54600a546040805163ffffffff808516825264010000000090940490931660208401528201526060016102f5565b600b546601000000000000900463ffffffff9081166000908152600c60205260409020547c01000000000000000000000000000000000000000000000000000000009004166102eb565b6104906105c7366004614dd6565b6112a9565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016102f5565b610613611345565b60405169ffffffffffffffffffff90911681526020016102f5565b61064161063c366004615241565b6114ff565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a0016102f5565b604080518082018252600e5473ffffffffffffffffffffffffffffffffffffffff81168083527401000000000000000000000000000000000000000090910463ffffffff166020928301819052835191825291810191909152016102f5565b6104906106e5366004614e58565b6115c7565b6104906106f8366004614dd6565b61184c565b600a54600b546040805160008152602081019390935261010090910460081c63ffffffff16908201526060016102f5565b61049061073c366004614dd6565b6118fd565b61049061074f366004614f91565b611a25565b6102eb6107623660046150c6565b61206f565b6102eb6107753660046150c6565b6120a5565b610490610788366004614e2c565b6120f7565b60125473ffffffffffffffffffffffffffffffffffffffff166105e6565b6102eb612487565b600f5473ffffffffffffffffffffffffffffffffffffffff166105e6565b6104906107df366004614ec4565b61254c565b6107f76107f2366004614dd6565b612eaa565b60405163ffffffff90911681526020016102f5565b610814612f75565b6040805195865263ffffffff909416602086015260ff9092169284019290925260179190910b606083015267ffffffffffffffff16608082015260a0016102f5565b60115473ffffffffffffffffffffffffffffffffffffffff166105e6565b610490610882366004615098565b613054565b610490610895366004614df3565b61318b565b6104906108a8366004614dd6565b613344565b6104906108bb366004614dd6565b613355565b600b5463ffffffff660100000000000090910481166000818152600c602090815260409182902082516060810184529054601781810b810b810b8084527801000000000000000000000000000000000000000000000000830488169484018590527c01000000000000000000000000000000000000000000000000000000009092049096169190930181905292939190910b9183610641565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff1691810191909152906109cd5750600092915050565b600b5460208201516000917201000000000000000000000000000000000000900463ffffffff169060069060ff16601f8110610a0b57610a0b6158cc565b600881049190910154600b54610a41926007166004026101000a90910463ffffffff9081169166010000000000009004166157c7565b63ffffffff16610a5191906156d6565b610a5f90633b9aca006156d6565b905081604001516bffffffffffffffffffffffff1681610a7f9190615636565b949350505050565b610a8f613366565b60115473ffffffffffffffffffffffffffffffffffffffff908116908316811415610ab957505050565b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff8416906370a082319060240160206040518083038186803b158015610b1e57600080fd5b505afa158015610b32573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b5691906150df565b50610b5f6133e9565b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009073ffffffffffffffffffffffffffffffffffffffff8316906370a082319060240160206040518083038186803b158015610bc757600080fd5b505afa158015610bdb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bff91906150df565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152602482018390529192509083169063a9059cbb90604401602060405180830381600087803b158015610c7357600080fd5b505af1158015610c87573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cab9190615076565b610d16576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f7472616e736665722072656d61696e696e672066756e6473206661696c65640060448201526064015b60405180910390fd5b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff86811691821790925560405190918416907f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a90600090a350505b5050565b60125473ffffffffffffffffffffffffffffffffffffffff16610dc760005473ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610ea257506040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821690636b14daf890610e529033906000903690600401615339565b60206040518083038186803b158015610e6a57600080fd5b505afa158015610e7e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ea29190615076565b610f08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c6044820152606401610d0d565b610f106133e9565b600b80547fffffffffffffffffffffffffffff0000000000000000ffffffffffffffffffff166a010000000000000000000063ffffffff8981169182027fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff16929092176e010000000000000000000000000000898416908102919091177fffffffffffff0000000000000000ffffffffffffffffffffffffffffffffffff1672010000000000000000000000000000000000008985169081027fffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffff1691909117760100000000000000000000000000000000000000000000948916948502177fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff167a01000000000000000000000000000000000000000000000000000062ffffff89169081029190911790955560408051938452602084019290925290820152606081019190915260808101919091527f0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f9060a00160405180910390a1505050505050565b6060600580548060200260200160405190810160405280929190818152602001828054801561111957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116110ee575b5050505050905090565b606060108054611132906157ec565b80601f016020809104026020016040519081016040528092919081815260200182805461115e906157ec565b80156111195780601f1061118057610100808354040283529160200191611119565b820191906000526020600020905b81548152906001019060200180831161118e57509395945050505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461122d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610d0d565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b73ffffffffffffffffffffffffffffffffffffffff818116600090815260136020526040902054163314611339576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4f6e6c792070617965652063616e2077697468647261770000000000000000006044820152606401610d0d565b61134281613878565b50565b6000805473ffffffffffffffffffffffffffffffffffffffff163314806114125750600f546040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690636b14daf8906113c29033906000903690600401615339565b60206040518083038186803b1580156113da57600080fd5b505afa1580156113ee573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114129190615076565b611478576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f4f6e6c79206f776e6572267265717565737465722063616e2063616c6c0000006044820152606401610d0d565b600b54600a546040805191825263ffffffff6101008404600881901c8216602085015260ff811684840152915164ffffffffff9092169366010000000000009004169133917f41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c9181900360600190a26114f281600161564e565b63ffffffff169250505090565b60008080808063ffffffff69ffffffffffffffffffff87161115611531575060009350839250829150819050806115be565b50505063ffffffff8084166000908152600c602090815260409182902082516060810184529054601781810b810b810b8084527801000000000000000000000000000000000000000000000000830487169484018590527c0100000000000000000000000000000000000000000000000000000000909204909516919093018190528695509190920b9250835b91939590929450565b6115cf613366565b828114611638576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f7472616e736d6974746572732e73697a6520213d207061796565732e73697a656044820152606401610d0d565b60005b83811015611845576000858583818110611657576116576158cc565b905060200201602081019061166c9190614dd6565b90506000848484818110611682576116826158cc565b90506020020160208101906116979190614dd6565b73ffffffffffffffffffffffffffffffffffffffff80841660009081526013602052604090205491925016801580806116fb57508273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b611761576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f706179656520616c7265616479207365740000000000000000000000000000006044820152606401610d0d565b73ffffffffffffffffffffffffffffffffffffffff848116600090815260136020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168583169081179091559083161461182e578273ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a45b50505050808061183d90615840565b91505061163b565b5050505050565b611854613366565b600f5473ffffffffffffffffffffffffffffffffffffffff9081169082168114610d8b57600f80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84811691821790925560408051928416835260208301919091527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491015b60405180910390a15050565b73ffffffffffffffffffffffffffffffffffffffff81811660009081526014602052604090205416331461198d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6f6e6c792070726f706f736564207061796565732063616e20616363657074006044820152606401610d0d565b73ffffffffffffffffffffffffffffffffffffffff81811660008181526013602090815260408083208054337fffffffffffffffffffffffff000000000000000000000000000000000000000080831682179093556014909452828520805490921690915590519416939092849290917f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b39190a45050565b60005a604080516101008082018352600b5460ff8116835290810464ffffffffff90811660208085018290526601000000000000840463ffffffff908116968601969096526a01000000000000000000008404861660608601526e01000000000000000000000000000084048616608086015272010000000000000000000000000000000000008404861660a0860152760100000000000000000000000000000000000000000000840490951660c08501527a01000000000000000000000000000000000000000000000000000090920462ffffff1660e08401529394509092918c013591821611611b73576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f7374616c65207265706f727400000000000000000000000000000000000000006044820152606401610d0d565b3360009081526002602052604090205460ff16611bec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f756e617574686f72697a6564207472616e736d697474657200000000000000006044820152606401610d0d565b600a548b3514611c58576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f636f6e666967446967657374206d69736d6174636800000000000000000000006044820152606401610d0d565b611c668a8a8a8a8a8a613b14565b8151611c73906001615676565b60ff168714611cde576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f77726f6e67206e756d626572206f66207369676e6174757265730000000000006044820152606401610d0d565b868514611d47576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f7369676e617475726573206f7574206f6620726567697374726174696f6e00006044820152606401610d0d565b60008a8a604051611d59929190615329565b604051908190038120611d70918e906020016153b6565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828252805160209182012083830190925260008084529083018190529092509060005b8a811015611f795760006001858a8460208110611ddd57611ddd6158cc565b611dea91901a601b615676565b8f8f86818110611dfc57611dfc6158cc565b905060200201358e8e87818110611e1557611e156158cc565b9050602002013560405160008152602001604052604051611e52949392919093845260ff9290921660208401526040830152606082015260800190565b6020604051602081039080840390855afa158015611e74573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081015173ffffffffffffffffffffffffffffffffffffffff811660009081526003602090815290849020838501909452925460ff8082161515808552610100909204169383019390935290955092509050611f52576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f7369676e6174757265206572726f7200000000000000000000000000000000006044820152606401610d0d565b826020015160080260ff166001901b84019350508080611f7190615840565b915050611dbe565b5081827e010101010101010101010101010101010101010101010101010101010101011614612004576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6475706c6963617465207369676e6572000000000000000000000000000000006044820152606401610d0d565b50600091506120539050838d836020020135848e8e8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250613bcb92505050565b905061206183828633614116565b505050505050505050505050565b600063ffffffff82111561208557506000919050565b5063ffffffff166000908152600c6020526040902054601790810b900b90565b600063ffffffff8211156120bb57506000919050565b5063ffffffff9081166000908152600c60205260409020547c010000000000000000000000000000000000000000000000000000000090041690565b60005473ffffffffffffffffffffffffffffffffffffffff163314806121c357506012546040517f6b14daf800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690636b14daf8906121739033906000903690600401615339565b60206040518083038186803b15801561218b57600080fd5b505afa15801561219f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121c39190615076565b612229576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f6e6c79206f776e65722662696c6c696e6741646d696e2063616e2063616c6c6044820152606401610d0d565b6000612233614281565b6011546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015291925060009173ffffffffffffffffffffffffffffffffffffffff909116906370a082319060240160206040518083038186803b1580156122a257600080fd5b505afa1580156122b6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122da91906150df565b905081811015612346576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f696e73756666696369656e742062616c616e63650000000000000000000000006044820152606401610d0d565b60115473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8561237861237286866157b0565b8761447c565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526024820152604401602060405180830381600087803b1580156123e357600080fd5b505af11580156123f7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061241b9190615076565b612481576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610d0d565b50505050565b6011546040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152600091829173ffffffffffffffffffffffffffffffffffffffff909116906370a082319060240160206040518083038186803b1580156124f557600080fd5b505afa158015612509573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061252d91906150df565b90506000612539614281565b9050612545818361573c565b9250505090565b612554613366565b601f865111156125c0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f746f6f206d616e79206f7261636c6573000000000000000000000000000000006044820152606401610d0d565b845186511461262b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f6f7261636c65206c656e677468206d69736d61746368000000000000000000006044820152606401610d0d565b8551612638856003615713565b60ff16106126a2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6661756c74792d6f7261636c65206620746f6f206869676800000000000000006044820152606401610d0d565b6126ae8460ff16614496565b825115612717576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f6f6e636861696e436f6e666967206d75737420626520656d70747900000000006044820152606401610d0d565b6040805160c081018252878152602080820188905260ff87168284015282517f0100000000000000000000000000000000000000000000000000000000000000918101919091527f0000000000000000000000000000000000000000000000000000000000000000601790810b841b60218301527f0000000000000000000000000000000000000000000000000000000000000000900b831b6039820152825160318183030181526051909101909252606081019190915267ffffffffffffffff8316608082015260a08101829052600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000ff1690556128166133e9565b60045460005b8181101561290f57600060048281548110612839576128396158cc565b60009182526020822001546005805473ffffffffffffffffffffffffffffffffffffffff90921693509084908110612873576128736158cc565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff948516835260038252604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016905594168252600290529190912080547fffffffffffffffffffffffffffffffffffff0000000000000000000000000000169055508061290781615840565b91505061281c565b5061291c60046000614af1565b61292860056000614af1565b60005b825151811015612c9d57600360008460000151838151811061294f5761294f6158cc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff16156129ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265706561746564207369676e657220616464726573730000000000000000006044820152606401610d0d565b604080518082019091526001815260ff821660208201528351805160039160009185908110612a1b57612a1b6158cc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281810192909252604001600090812083518154948401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009095169015157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff161761010060ff90951694909402939093179092558401518051600292919084908110612acd57612acd6158cc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff1682528101919091526040016000205460ff1615612b68576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564207472616e736d69747465722061646472657373000000006044820152606401610d0d565b60405180606001604052806001151581526020018260ff16815260200160006bffffffffffffffffffffffff168152506002600085602001518481518110612bb257612bb26158cc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600020835181549385015194909201516bffffffffffffffffffffffff1662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff60ff95909516610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff931515939093167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090941693909317919091179290921617905580612c9581615840565b91505061292b565b5081518051612cb491600491602090910190614b0f565b506020808301518051612ccb926005920190614b0f565b506040820151600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660ff909216919091179055600d80547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff811664010000000063ffffffff438116820292831785559083048116936001939092600092612d5d92869290821691161761564e565b92506101000a81548163ffffffff021916908363ffffffff160217905550612dbc4630600d60009054906101000a900463ffffffff1663ffffffff1686600001518760200151886040015189606001518a608001518b60a00151614500565b600a819055600d5484516020860151604080880151606089015160808a015160a08b015193517f1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e0598612e25988b98919763ffffffff90911696919590949093909290919061553d565b60405180910390a1600b546601000000000000900463ffffffff1660005b845151811015612e9d5781600682601f8110612e6157612e616158cc565b600891828204019190066004026101000a81548163ffffffff021916908363ffffffff1602179055508080612e9590615840565b915050612e43565b5050505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602090815260408083208151606081018352905460ff80821615158084526101008304909116948301949094526201000090046bffffffffffffffffffffffff169181019190915290612f1e5750600092915050565b6006816020015160ff16601f8110612f3857612f386158cc565b600881049190910154600b54612f6e926007166004026101000a90910463ffffffff9081169166010000000000009004166157c7565b9392505050565b600080808080333214612fe4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4f6e6c792063616c6c61626c6520627920454f410000000000000000000000006044820152606401610d0d565b5050600a54600b5463ffffffff6601000000000000820481166000908152600c60205260409020549296610100909204600881901c8216965064ffffffffff169450601783900b93507c010000000000000000000000000000000000000000000000000000000090920490911690565b61305c613366565b60408051808201909152600e5473ffffffffffffffffffffffffffffffffffffffff8082168084527401000000000000000000000000000000000000000090920463ffffffff16602084015284161415806130c757508163ffffffff16816020015163ffffffff1614155b156131865760408051808201825273ffffffffffffffffffffffffffffffffffffffff85811680835263ffffffff8681166020948501819052600e80547fffffffffffffffff00000000000000000000000000000000000000000000000016841774010000000000000000000000000000000000000000830217905586518786015187519316835294820152909392909116917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a35b505050565b73ffffffffffffffffffffffffffffffffffffffff82811660009081526013602052604090205416331461321b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6f6e6c792063757272656e742070617965652063616e207570646174650000006044820152606401610d0d565b3373ffffffffffffffffffffffffffffffffffffffff8216141561329b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f63616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610d0d565b73ffffffffffffffffffffffffffffffffffffffff808316600090815260146020526040902080548383167fffffffffffffffffffffffff0000000000000000000000000000000000000000821681179092559091169081146131865760405173ffffffffffffffffffffffffffffffffffffffff8084169133918616907f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836790600090a4505050565b61334c613366565b611342816145ac565b61335d613366565b611342816146a2565b60005473ffffffffffffffffffffffffffffffffffffffff1633146133e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610d0d565b565b601154600b54604080516103e081019182905273ffffffffffffffffffffffffffffffffffffffff90931692660100000000000090920463ffffffff1691600091600690601f908285855b82829054906101000a900463ffffffff1663ffffffff168152602001906004019060208260030104928301926001038202915080841161343457905050505050509050600060058054806020026020016040519081016040528092919081815260200182805480156134dc57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116134b1575b5050505050905060005b815181101561386a57600060026000848481518110613507576135076158cc565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160029054906101000a90046bffffffffffffffffffffffff166bffffffffffffffffffffffff16905060006002600085858151811061358d5761358d6158cc565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160026101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555060008483601f8110613614576136146158cc565b6020020151600b5490870363ffffffff90811692507201000000000000000000000000000000000000909104168102633b9aca00028201801561385f57600060136000878781518110613669576136696158cc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff90811683529082019290925260409081016000205490517fa9059cbb00000000000000000000000000000000000000000000000000000000815290821660048201819052602482018590529250908a169063a9059cbb90604401602060405180830381600087803b1580156136fe57600080fd5b505af1158015613712573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137369190615076565b61379c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610d0d565b878786601f81106137af576137af6158cc565b602002019063ffffffff16908163ffffffff16815250508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16878781518110613806576138066158cc565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c8560405161385591815260200190565b60405180910390a4505b5050506001016134e6565b50611845600683601f614b99565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600260209081526040918290208251606081018452905460ff80821615158084526101008304909116938301939093526201000090046bffffffffffffffffffffffff16928101929092526138e7575050565b60006138f283610959565b905080156131865773ffffffffffffffffffffffffffffffffffffffff838116600090815260136020526040908190205460115491517fa9059cbb000000000000000000000000000000000000000000000000000000008152908316600482018190526024820185905292919091169063a9059cbb90604401602060405180830381600087803b15801561398557600080fd5b505af1158015613999573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139bd9190615076565b613a23576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610d0d565b600b60000160069054906101000a900463ffffffff166006846020015160ff16601f8110613a5357613a536158cc565b6008810491909101805460079092166004026101000a63ffffffff81810219909316939092169190910291909117905573ffffffffffffffffffffffffffffffffffffffff84811660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff169055601154915186815291841693851692917fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c910160405180910390a450505050565b6000613b218260206156d6565b613b2c8560206156d6565b613b3888610144615636565b613b429190615636565b613b4c9190615636565b613b57906000615636565b9050368114613bc2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f63616c6c64617461206c656e677468206d69736d6174636800000000000000006044820152606401610d0d565b50505050505050565b600080613bd783614743565b9050601f8160400151511115613c49576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6e756d206f62736572766174696f6e73206f7574206f6620626f756e647300006044820152606401610d0d565b604081015151865160ff1610613cbb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f746f6f206665772076616c75657320746f207472757374206d656469616e00006044820152606401610d0d565b64ffffffffff841660208701526040810151805160009190613cdf9060029061569b565b81518110613cef57613cef6158cc565b602002602001015190508060170b7f000000000000000000000000000000000000000000000000000000000000000060170b13158015613d5557507f000000000000000000000000000000000000000000000000000000000000000060170b8160170b13155b613dbb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f6d656469616e206973206f7574206f66206d696e2d6d61782072616e676500006044820152606401610d0d565b60408701805190613dcb82615879565b63ffffffff1663ffffffff168152505060405180606001604052808260170b8152602001836000015163ffffffff1681526020014263ffffffff16815250600c6000896040015163ffffffff1663ffffffff16815260200190815260200160002060008201518160000160006101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908360170b77ffffffffffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160186101000a81548163ffffffff021916908363ffffffff160217905550604082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555090505086600b60008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548164ffffffffff021916908364ffffffffff16021790555060408201518160000160066101000a81548163ffffffff021916908363ffffffff160217905550606082015181600001600a6101000a81548163ffffffff021916908363ffffffff160217905550608082015181600001600e6101000a81548163ffffffff021916908363ffffffff16021790555060a08201518160000160126101000a81548163ffffffff021916908363ffffffff16021790555060c08201518160000160166101000a81548163ffffffff021916908363ffffffff16021790555060e082015181600001601a6101000a81548162ffffff021916908362ffffff160217905550905050866040015163ffffffff167fc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a823385600001518660400151876020015188606001518d8d60405161405f9897969594939291906153d0565b60405180910390a26040808801518351915163ffffffff9283168152600092909116907f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac602719060200160405180910390a3866040015163ffffffff168160170b7f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040516140ef91815260200190565b60405180910390a361410887604001518260170b6147e8565b506060015195945050505050565b60008360170b121561412757612481565b600061414e633b9aca003a04866080015163ffffffff16876060015163ffffffff1661497f565b90506010360260005a905060006141778663ffffffff1685858b60e0015162ffffff16866149a5565b90506000670de0b6b3a764000077ffffffffffffffffffffffffffffffffffffffffffffffff8916830273ffffffffffffffffffffffffffffffffffffffff881660009081526002602052604090205460c08c01519290910492506201000090046bffffffffffffffffffffffff9081169163ffffffff16633b9aca00028284010190811682111561420f5750505050505050612481565b73ffffffffffffffffffffffffffffffffffffffff8816600090815260026020526040902080546bffffffffffffffffffffffff90921662010000027fffffffffffffffffffffffffffffffffffff000000000000000000000000ffff90921691909117905550505050505050505050565b60008060058054806020026020016040519081016040528092919081815260200182805480156142e757602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116142bc575b50508351600b54604080516103e08101918290529697509195660100000000000090910463ffffffff169450600093509150600690601f908285855b82829054906101000a900463ffffffff1663ffffffff16815260200190600401906020826003010492830192600103820291508084116143235790505050505050905060005b838110156143b6578181601f8110614383576143836158cc565b602002015161439290846157c7565b6143a29063ffffffff1687615636565b9550806143ae81615840565b915050614369565b50600b546143e4907201000000000000000000000000000000000000900463ffffffff16633b9aca006156d6565b6143ee90866156d6565b945060005b838110156144745760026000868381518110614411576144116158cc565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff16825281019190915260400160002054614460906201000090046bffffffffffffffffffffffff1687615636565b95508061446c81615840565b9150506143f3565b505050505090565b60008183101561448d575081614490565b50805b92915050565b80600010611342576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f66206d75737420626520706f73697469766500000000000000000000000000006044820152606401610d0d565b6000808a8a8a8a8a8a8a8a8a60405160200161452499989796959493929190615498565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815291905280516020909101207dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167e01000000000000000000000000000000000000000000000000000000000000179150505b9998505050505050505050565b73ffffffffffffffffffffffffffffffffffffffff811633141561462c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610d0d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60125473ffffffffffffffffffffffffffffffffffffffff9081169082168114610d8b57601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff84811691821790925560408051928416835260208301919091527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291016118f1565b6147776040518060800160405280600063ffffffff1681526020016060815260200160608152602001600060170b81525090565b600080606060008580602001905181019061479291906150f8565b929650909450925090506147a68683614a23565b81516040805160208082019690965281519082018252918252805160808101825263ffffffff969096168652938501529183015260170b606082015292915050565b60408051808201909152600e5473ffffffffffffffffffffffffffffffffffffffff81168083527401000000000000000000000000000000000000000090910463ffffffff16602083015261483c57505050565b60006148496001856157c7565b63ffffffff8181166000818152600c60209081526040918290205490870151875192516024810194909452601791820b90910b6044840181905289851660648501526084840189905294955061491993169160a401604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fbeed9b5100000000000000000000000000000000000000000000000000000000179052614ab5565b611845576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f696e73756666696369656e7420676173000000000000000000000000000000006044820152606401610d0d565b6000838381101561499257600285850304015b61499c818461447c565b95945050505050565b600081861015614a11576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6c6566744761732063616e6e6f742065786365656420696e697469616c4761736044820152606401610d0d565b50633b9aca0094039190910101020290565b600081516020614a3391906156d6565b614a3e9060a0615636565b614a49906000615636565b905080835114613186576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f7265706f7274206c656e677468206d69736d61746368000000000000000000006044820152606401610d0d565b60005a6113888110614ae95761138881039050846040820482031115614ae9576000808451602086016000888af150600191505b509392505050565b50805460008255906000526020600020908101906113429190614c2c565b828054828255906000526020600020908101928215614b89579160200282015b82811115614b8957825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614b2f565b50614b95929150614c2c565b5090565b600483019183908215614b895791602002820160005b83821115614bf357835183826101000a81548163ffffffff021916908363ffffffff1602179055509260200192600401602081600301049283019260010302614baf565b8015614c235782816101000a81549063ffffffff0219169055600401602081600301049283019260010302614bf3565b5050614b959291505b5b80821115614b955760008155600101614c2d565b60008083601f840112614c5357600080fd5b50813567ffffffffffffffff811115614c6b57600080fd5b6020830191508360208260051b8501011115614c8657600080fd5b9250929050565b600082601f830112614c9e57600080fd5b81356020614cb3614cae83615612565b6155c3565b80838252828201915082860187848660051b8901011115614cd357600080fd5b60005b85811015614cfb578135614ce98161592a565b84529284019290840190600101614cd6565b5090979650505050505050565b600082601f830112614d1957600080fd5b813567ffffffffffffffff811115614d3357614d336158fb565b614d6460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016155c3565b818152846020838601011115614d7957600080fd5b816020850160208301376000918101602001919091529392505050565b8051601781900b8114614da857600080fd5b919050565b803567ffffffffffffffff81168114614da857600080fd5b803560ff81168114614da857600080fd5b600060208284031215614de857600080fd5b8135612f6e8161592a565b60008060408385031215614e0657600080fd5b8235614e118161592a565b91506020830135614e218161592a565b809150509250929050565b60008060408385031215614e3f57600080fd5b8235614e4a8161592a565b946020939093013593505050565b60008060008060408587031215614e6e57600080fd5b843567ffffffffffffffff80821115614e8657600080fd5b614e9288838901614c41565b90965094506020870135915080821115614eab57600080fd5b50614eb887828801614c41565b95989497509550505050565b60008060008060008060c08789031215614edd57600080fd5b863567ffffffffffffffff80821115614ef557600080fd5b614f018a838b01614c8d565b97506020890135915080821115614f1757600080fd5b614f238a838b01614c8d565b9650614f3160408a01614dc5565b95506060890135915080821115614f4757600080fd5b614f538a838b01614d08565b9450614f6160808a01614dad565b935060a0890135915080821115614f7757600080fd5b50614f8489828a01614d08565b9150509295509295509295565b60008060008060008060008060e0898b031215614fad57600080fd5b606089018a811115614fbe57600080fd5b8998503567ffffffffffffffff80821115614fd857600080fd5b818b0191508b601f830112614fec57600080fd5b813581811115614ffb57600080fd5b8c602082850101111561500d57600080fd5b6020830199508098505060808b013591508082111561502b57600080fd5b6150378c838d01614c41565b909750955060a08b013591508082111561505057600080fd5b5061505d8b828c01614c41565b999c989b50969995989497949560c00135949350505050565b60006020828403121561508857600080fd5b81518015158114612f6e57600080fd5b600080604083850312156150ab57600080fd5b82356150b68161592a565b91506020830135614e218161594c565b6000602082840312156150d857600080fd5b5035919050565b6000602082840312156150f157600080fd5b5051919050565b6000806000806080858703121561510e57600080fd5b84516151198161594c565b809450506020808601519350604086015167ffffffffffffffff81111561513f57600080fd5b8601601f8101881361515057600080fd5b805161515e614cae82615612565b8082825284820191508484018b868560051b870101111561517e57600080fd5b600094505b838510156151a85761519481614d96565b835260019490940193918501918501615183565b5080965050505050506151bd60608601614d96565b905092959194509250565b600080600080600060a086880312156151e057600080fd5b85356151eb8161594c565b945060208601356151fb8161594c565b9350604086013561520b8161594c565b9250606086013561521b8161594c565b9150608086013562ffffff8116811461523357600080fd5b809150509295509295909350565b60006020828403121561525357600080fd5b813569ffffffffffffffffffff81168114612f6e57600080fd5b600081518084526020808501945080840160005b838110156152b357815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101615281565b509495945050505050565b6000815180845260005b818110156152e4576020818501810151868301820152016152c8565b818111156152f6576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8183823760009101908152919050565b73ffffffffffffffffffffffffffffffffffffffff8416815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b602081526000612f6e602083018461526d565b828152608081016060836020840137600081529392505050565b600061010080830160178c810b8552602073ffffffffffffffffffffffffffffffffffffffff8d168187015263ffffffff8c1660408701528360608701528293508a5180845261012087019450818c01935060005b81811015615443578451840b86529482019493820193600101615425565b5050505050828103608084015261545a81886152be565b91505061546c60a083018660170b9052565b8360c083015261459f60e083018464ffffffffff169052565b602081526000612f6e60208301846152be565b60006101208b835273ffffffffffffffffffffffffffffffffffffffff8b16602084015267ffffffffffffffff808b1660408501528160608501526154df8285018b61526d565b915083820360808501526154f3828a61526d565b915060ff881660a085015283820360c085015261551082886152be565b90861660e0850152838103610100850152905061552d81856152be565b9c9b505050505050505050505050565b600061012063ffffffff808d1684528b6020850152808b1660408501525080606084015261556d8184018a61526d565b90508281036080840152615581818961526d565b905060ff871660a084015282810360c084015261559e81876152be565b905067ffffffffffffffff851660e084015282810361010084015261552d81856152be565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561560a5761560a6158fb565b604052919050565b600067ffffffffffffffff82111561562c5761562c6158fb565b5060051b60200190565b600082198211156156495761564961589d565b500190565b600063ffffffff80831681851680830382111561566d5761566d61589d565b01949350505050565b600060ff821660ff84168060ff038211156156935761569361589d565b019392505050565b6000826156d1577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561570e5761570e61589d565b500290565b600060ff821660ff84168160ff04811182151516156157345761573461589d565b029392505050565b6000808312837f8000000000000000000000000000000000000000000000000000000000000000018312811516156157765761577661589d565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0183138116156157aa576157aa61589d565b50500390565b6000828210156157c2576157c261589d565b500390565b600063ffffffff838116908316818110156157e4576157e461589d565b039392505050565b600181811c9082168061580057607f821691505b6020821081141561583a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156158725761587261589d565b5060010190565b600063ffffffff808316818114156158935761589361589d565b6001019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461134257600080fd5b63ffffffff8116811461134257600080fdfea164736f6c6343000806000a",
}

var OCR2AggregatorABI = OCR2AggregatorMetaData.ABI

var OCR2AggregatorBin = OCR2AggregatorMetaData.Bin

func DeployOCR2Aggregator(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, minAnswer_ *big.Int, maxAnswer_ *big.Int, billingAccessController common.Address, requesterAccessController common.Address, decimals_ uint8, description_ string) (common.Address, *types.Transaction, *OCR2Aggregator, error) {
	parsed, err := OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OCR2AggregatorBin), backend, link, minAnswer_, maxAnswer_, billingAccessController, requesterAccessController, decimals_, description_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OCR2Aggregator{OCR2AggregatorCaller: OCR2AggregatorCaller{contract: contract}, OCR2AggregatorTransactor: OCR2AggregatorTransactor{contract: contract}, OCR2AggregatorFilterer: OCR2AggregatorFilterer{contract: contract}}, nil
}

type OCR2Aggregator struct {
	address common.Address
	abi     abi.ABI
	OCR2AggregatorCaller
	OCR2AggregatorTransactor
	OCR2AggregatorFilterer
}

type OCR2AggregatorCaller struct {
	contract *bind.BoundContract
}

type OCR2AggregatorTransactor struct {
	contract *bind.BoundContract
}

type OCR2AggregatorFilterer struct {
	contract *bind.BoundContract
}

type OCR2AggregatorSession struct {
	Contract     *OCR2Aggregator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OCR2AggregatorCallerSession struct {
	Contract *OCR2AggregatorCaller
	CallOpts bind.CallOpts
}

type OCR2AggregatorTransactorSession struct {
	Contract     *OCR2AggregatorTransactor
	TransactOpts bind.TransactOpts
}

type OCR2AggregatorRaw struct {
	Contract *OCR2Aggregator
}

type OCR2AggregatorCallerRaw struct {
	Contract *OCR2AggregatorCaller
}

type OCR2AggregatorTransactorRaw struct {
	Contract *OCR2AggregatorTransactor
}

func NewOCR2Aggregator(address common.Address, backend bind.ContractBackend) (*OCR2Aggregator, error) {
	abi, err := abi.JSON(strings.NewReader(OCR2AggregatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOCR2Aggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OCR2Aggregator{address: address, abi: abi, OCR2AggregatorCaller: OCR2AggregatorCaller{contract: contract}, OCR2AggregatorTransactor: OCR2AggregatorTransactor{contract: contract}, OCR2AggregatorFilterer: OCR2AggregatorFilterer{contract: contract}}, nil
}

func NewOCR2AggregatorCaller(address common.Address, caller bind.ContractCaller) (*OCR2AggregatorCaller, error) {
	contract, err := bindOCR2Aggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorCaller{contract: contract}, nil
}

func NewOCR2AggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OCR2AggregatorTransactor, error) {
	contract, err := bindOCR2Aggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorTransactor{contract: contract}, nil
}

func NewOCR2AggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OCR2AggregatorFilterer, error) {
	contract, err := bindOCR2Aggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorFilterer{contract: contract}, nil
}

func bindOCR2Aggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OCR2AggregatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OCR2Aggregator *OCR2AggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Aggregator.Contract.OCR2AggregatorCaller.contract.Call(opts, result, method, params...)
}

func (_OCR2Aggregator *OCR2AggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.OCR2AggregatorTransactor.contract.Transfer(opts)
}

func (_OCR2Aggregator *OCR2AggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.OCR2AggregatorTransactor.contract.Transact(opts, method, params...)
}

func (_OCR2Aggregator *OCR2AggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OCR2Aggregator.Contract.contract.Call(opts, result, method, params...)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.contract.Transfer(opts)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.contract.Transact(opts, method, params...)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) Decimals() (uint8, error) {
	return _OCR2Aggregator.Contract.Decimals(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) Decimals() (uint8, error) {
	return _OCR2Aggregator.Contract.Decimals(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) Description() (string, error) {
	return _OCR2Aggregator.Contract.Description(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) Description() (string, error) {
	return _OCR2Aggregator.Contract.Description(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetAnswer(&_OCR2Aggregator.CallOpts, roundId)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetAnswer(&_OCR2Aggregator.CallOpts, roundId)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetBilling(opts *bind.CallOpts) (GetBilling,

	error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getBilling")

	outstruct := new(GetBilling)
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

func (_OCR2Aggregator *OCR2AggregatorSession) GetBilling() (GetBilling,

	error) {
	return _OCR2Aggregator.Contract.GetBilling(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetBilling() (GetBilling,

	error) {
	return _OCR2Aggregator.Contract.GetBilling(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetBillingAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getBillingAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) GetBillingAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetBillingAccessController(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetBillingAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetBillingAccessController(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetLinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getLinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) GetLinkToken() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetLinkToken(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetLinkToken() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetLinkToken(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetRequesterAccessController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getRequesterAccessController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) GetRequesterAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetRequesterAccessController(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetRequesterAccessController() (common.Address, error) {
	return _OCR2Aggregator.Contract.GetRequesterAccessController(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getRoundData", roundId)

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

func (_OCR2Aggregator *OCR2AggregatorSession) GetRoundData(roundId *big.Int) (GetRoundData,

	error) {
	return _OCR2Aggregator.Contract.GetRoundData(&_OCR2Aggregator.CallOpts, roundId)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetRoundData(roundId *big.Int) (GetRoundData,

	error) {
	return _OCR2Aggregator.Contract.GetRoundData(&_OCR2Aggregator.CallOpts, roundId)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetTimestamp(&_OCR2Aggregator.CallOpts, roundId)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _OCR2Aggregator.Contract.GetTimestamp(&_OCR2Aggregator.CallOpts, roundId)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetTransmitters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getTransmitters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) GetTransmitters() ([]common.Address, error) {
	return _OCR2Aggregator.Contract.GetTransmitters(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetTransmitters() ([]common.Address, error) {
	return _OCR2Aggregator.Contract.GetTransmitters(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) GetValidatorConfig(opts *bind.CallOpts) (GetValidatorConfig,

	error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "getValidatorConfig")

	outstruct := new(GetValidatorConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.GasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) GetValidatorConfig() (GetValidatorConfig,

	error) {
	return _OCR2Aggregator.Contract.GetValidatorConfig(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) GetValidatorConfig() (GetValidatorConfig,

	error) {
	return _OCR2Aggregator.Contract.GetValidatorConfig(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) LatestAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestAnswer(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestAnswer(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

	error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestConfigDetails")

	outstruct := new(LatestConfigDetails)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfigCount = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ConfigDigest = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCR2Aggregator.Contract.LatestConfigDetails(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestConfigDetails() (LatestConfigDetails,

	error) {
	return _OCR2Aggregator.Contract.LatestConfigDetails(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

	error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestConfigDigestAndEpoch")

	outstruct := new(LatestConfigDigestAndEpoch)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ScanLogs = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ConfigDigest = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Epoch = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _OCR2Aggregator.Contract.LatestConfigDigestAndEpoch(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestConfigDigestAndEpoch() (LatestConfigDigestAndEpoch,

	error) {
	return _OCR2Aggregator.Contract.LatestConfigDigestAndEpoch(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) LatestRound() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestRound(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestRound() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestRound(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestRoundData")

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

func (_OCR2Aggregator *OCR2AggregatorSession) LatestRoundData() (LatestRoundData,

	error) {
	return _OCR2Aggregator.Contract.LatestRoundData(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _OCR2Aggregator.Contract.LatestRoundData(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) LatestTimestamp() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestTimestamp(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestTimestamp() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LatestTimestamp(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LatestTransmissionDetails(opts *bind.CallOpts) (LatestTransmissionDetails,

	error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "latestTransmissionDetails")

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

func (_OCR2Aggregator *OCR2AggregatorSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _OCR2Aggregator.Contract.LatestTransmissionDetails(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LatestTransmissionDetails() (LatestTransmissionDetails,

	error) {
	return _OCR2Aggregator.Contract.LatestTransmissionDetails(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "linkAvailableForPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LinkAvailableForPayment(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) LinkAvailableForPayment() (*big.Int, error) {
	return _OCR2Aggregator.Contract.LinkAvailableForPayment(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) MaxAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "maxAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) MaxAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MaxAnswer(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) MaxAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MaxAnswer(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) MinAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "minAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) MinAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MinAnswer(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) MinAnswer() (*big.Int, error) {
	return _OCR2Aggregator.Contract.MinAnswer(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "oracleObservationCount", transmitterAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _OCR2Aggregator.Contract.OracleObservationCount(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) OracleObservationCount(transmitterAddress common.Address) (uint32, error) {
	return _OCR2Aggregator.Contract.OracleObservationCount(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "owedPayment", transmitterAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _OCR2Aggregator.Contract.OwedPayment(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) OwedPayment(transmitterAddress common.Address) (*big.Int, error) {
	return _OCR2Aggregator.Contract.OwedPayment(&_OCR2Aggregator.CallOpts, transmitterAddress)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) Owner() (common.Address, error) {
	return _OCR2Aggregator.Contract.Owner(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) Owner() (common.Address, error) {
	return _OCR2Aggregator.Contract.Owner(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) TypeAndVersion() (string, error) {
	return _OCR2Aggregator.Contract.TypeAndVersion(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) TypeAndVersion() (string, error) {
	return _OCR2Aggregator.Contract.TypeAndVersion(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OCR2Aggregator.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OCR2Aggregator *OCR2AggregatorSession) Version() (*big.Int, error) {
	return _OCR2Aggregator.Contract.Version(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorCallerSession) Version() (*big.Int, error) {
	return _OCR2Aggregator.Contract.Version(&_OCR2Aggregator.CallOpts)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "acceptOwnership")
}

func (_OCR2Aggregator *OCR2AggregatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptOwnership(&_OCR2Aggregator.TransactOpts)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptOwnership(&_OCR2Aggregator.TransactOpts)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "acceptPayeeship", transmitter)
}

func (_OCR2Aggregator *OCR2AggregatorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptPayeeship(&_OCR2Aggregator.TransactOpts, transmitter)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) AcceptPayeeship(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.AcceptPayeeship(&_OCR2Aggregator.TransactOpts, transmitter)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "requestNewRound")
}

func (_OCR2Aggregator *OCR2AggregatorSession) RequestNewRound() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.RequestNewRound(&_OCR2Aggregator.TransactOpts)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) RequestNewRound() (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.RequestNewRound(&_OCR2Aggregator.TransactOpts)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setBilling", maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_OCR2Aggregator *OCR2AggregatorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBilling(&_OCR2Aggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetBilling(maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBilling(&_OCR2Aggregator.TransactOpts, maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setBillingAccessController", _billingAccessController)
}

func (_OCR2Aggregator *OCR2AggregatorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBillingAccessController(&_OCR2Aggregator.TransactOpts, _billingAccessController)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetBillingAccessController(_billingAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetBillingAccessController(&_OCR2Aggregator.TransactOpts, _billingAccessController)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setConfig", signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_OCR2Aggregator *OCR2AggregatorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetConfig(&_OCR2Aggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetConfig(signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetConfig(&_OCR2Aggregator.TransactOpts, signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) SetLinkToken(opts *bind.TransactOpts, linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setLinkToken", linkToken, recipient)
}

func (_OCR2Aggregator *OCR2AggregatorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetLinkToken(&_OCR2Aggregator.TransactOpts, linkToken, recipient)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetLinkToken(linkToken common.Address, recipient common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetLinkToken(&_OCR2Aggregator.TransactOpts, linkToken, recipient)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setPayees", transmitters, payees)
}

func (_OCR2Aggregator *OCR2AggregatorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetPayees(&_OCR2Aggregator.TransactOpts, transmitters, payees)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetPayees(transmitters []common.Address, payees []common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetPayees(&_OCR2Aggregator.TransactOpts, transmitters, payees)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) SetRequesterAccessController(opts *bind.TransactOpts, requesterAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setRequesterAccessController", requesterAccessController)
}

func (_OCR2Aggregator *OCR2AggregatorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetRequesterAccessController(&_OCR2Aggregator.TransactOpts, requesterAccessController)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetRequesterAccessController(requesterAccessController common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetRequesterAccessController(&_OCR2Aggregator.TransactOpts, requesterAccessController)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) SetValidatorConfig(opts *bind.TransactOpts, newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "setValidatorConfig", newValidator, newGasLimit)
}

func (_OCR2Aggregator *OCR2AggregatorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetValidatorConfig(&_OCR2Aggregator.TransactOpts, newValidator, newGasLimit)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) SetValidatorConfig(newValidator common.Address, newGasLimit uint32) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.SetValidatorConfig(&_OCR2Aggregator.TransactOpts, newValidator, newGasLimit)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "transferOwnership", to)
}

func (_OCR2Aggregator *OCR2AggregatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferOwnership(&_OCR2Aggregator.TransactOpts, to)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferOwnership(&_OCR2Aggregator.TransactOpts, to)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "transferPayeeship", transmitter, proposed)
}

func (_OCR2Aggregator *OCR2AggregatorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferPayeeship(&_OCR2Aggregator.TransactOpts, transmitter, proposed)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) TransferPayeeship(transmitter common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.TransferPayeeship(&_OCR2Aggregator.TransactOpts, transmitter, proposed)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "transmit", reportContext, report, rs, ss, rawVs)
}

func (_OCR2Aggregator *OCR2AggregatorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.Transmit(&_OCR2Aggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) Transmit(reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.Transmit(&_OCR2Aggregator.TransactOpts, reportContext, report, rs, ss, rawVs)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "withdrawFunds", recipient, amount)
}

func (_OCR2Aggregator *OCR2AggregatorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawFunds(&_OCR2Aggregator.TransactOpts, recipient, amount)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) WithdrawFunds(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawFunds(&_OCR2Aggregator.TransactOpts, recipient, amount)
}

func (_OCR2Aggregator *OCR2AggregatorTransactor) WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.contract.Transact(opts, "withdrawPayment", transmitter)
}

func (_OCR2Aggregator *OCR2AggregatorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawPayment(&_OCR2Aggregator.TransactOpts, transmitter)
}

func (_OCR2Aggregator *OCR2AggregatorTransactorSession) WithdrawPayment(transmitter common.Address) (*types.Transaction, error) {
	return _OCR2Aggregator.Contract.WithdrawPayment(&_OCR2Aggregator.TransactOpts, transmitter)
}

type OCR2AggregatorAnswerUpdatedIterator struct {
	Event *OCR2AggregatorAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorAnswerUpdated)
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
		it.Event = new(OCR2AggregatorAnswerUpdated)
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

func (it *OCR2AggregatorAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*OCR2AggregatorAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorAnswerUpdatedIterator{contract: _OCR2Aggregator.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorAnswerUpdated)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseAnswerUpdated(log types.Log) (*OCR2AggregatorAnswerUpdated, error) {
	event := new(OCR2AggregatorAnswerUpdated)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorBillingAccessControllerSetIterator struct {
	Event *OCR2AggregatorBillingAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorBillingAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorBillingAccessControllerSet)
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
		it.Event = new(OCR2AggregatorBillingAccessControllerSet)
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

func (it *OCR2AggregatorBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OCR2AggregatorBillingAccessControllerSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorBillingAccessControllerSetIterator{contract: _OCR2Aggregator.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorBillingAccessControllerSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseBillingAccessControllerSet(log types.Log) (*OCR2AggregatorBillingAccessControllerSet, error) {
	event := new(OCR2AggregatorBillingAccessControllerSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorBillingSetIterator struct {
	Event *OCR2AggregatorBillingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorBillingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorBillingSet)
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
		it.Event = new(OCR2AggregatorBillingSet)
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

func (it *OCR2AggregatorBillingSetIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorBillingSet struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
	Raw                       types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterBillingSet(opts *bind.FilterOpts) (*OCR2AggregatorBillingSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorBillingSetIterator{contract: _OCR2Aggregator.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorBillingSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorBillingSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseBillingSet(log types.Log) (*OCR2AggregatorBillingSet, error) {
	event := new(OCR2AggregatorBillingSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorConfigSetIterator struct {
	Event *OCR2AggregatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorConfigSet)
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
		it.Event = new(OCR2AggregatorConfigSet)
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

func (it *OCR2AggregatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorConfigSet struct {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OCR2AggregatorConfigSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorConfigSetIterator{contract: _OCR2Aggregator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorConfigSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseConfigSet(log types.Log) (*OCR2AggregatorConfigSet, error) {
	event := new(OCR2AggregatorConfigSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorLinkTokenSetIterator struct {
	Event *OCR2AggregatorLinkTokenSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorLinkTokenSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorLinkTokenSet)
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
		it.Event = new(OCR2AggregatorLinkTokenSet)
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

func (it *OCR2AggregatorLinkTokenSetIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, oldLinkToken []common.Address, newLinkToken []common.Address) (*OCR2AggregatorLinkTokenSetIterator, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorLinkTokenSetIterator{contract: _OCR2Aggregator.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorLinkTokenSet, oldLinkToken []common.Address, newLinkToken []common.Address) (event.Subscription, error) {

	var oldLinkTokenRule []interface{}
	for _, oldLinkTokenItem := range oldLinkToken {
		oldLinkTokenRule = append(oldLinkTokenRule, oldLinkTokenItem)
	}
	var newLinkTokenRule []interface{}
	for _, newLinkTokenItem := range newLinkToken {
		newLinkTokenRule = append(newLinkTokenRule, newLinkTokenItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "LinkTokenSet", oldLinkTokenRule, newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorLinkTokenSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseLinkTokenSet(log types.Log) (*OCR2AggregatorLinkTokenSet, error) {
	event := new(OCR2AggregatorLinkTokenSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorNewRoundIterator struct {
	Event *OCR2AggregatorNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorNewRound)
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
		it.Event = new(OCR2AggregatorNewRound)
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

func (it *OCR2AggregatorNewRoundIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*OCR2AggregatorNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorNewRoundIterator{contract: _OCR2Aggregator.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorNewRound)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseNewRound(log types.Log) (*OCR2AggregatorNewRound, error) {
	event := new(OCR2AggregatorNewRound)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorNewTransmissionIterator struct {
	Event *OCR2AggregatorNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorNewTransmission)
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
		it.Event = new(OCR2AggregatorNewTransmission)
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

func (it *OCR2AggregatorNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorNewTransmission struct {
	AggregatorRoundId     uint32
	Answer                *big.Int
	Transmitter           common.Address
	ObservationsTimestamp uint32
	Observations          []*big.Int
	Observers             []byte
	JuelsPerFeeCoin       *big.Int
	ConfigDigest          [32]byte
	EpochAndRound         *big.Int
	Raw                   types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*OCR2AggregatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorNewTransmissionIterator{contract: _OCR2Aggregator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorNewTransmission)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseNewTransmission(log types.Log) (*OCR2AggregatorNewTransmission, error) {
	event := new(OCR2AggregatorNewTransmission)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorOraclePaidIterator struct {
	Event *OCR2AggregatorOraclePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorOraclePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorOraclePaid)
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
		it.Event = new(OCR2AggregatorOraclePaid)
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

func (it *OCR2AggregatorOraclePaidIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*OCR2AggregatorOraclePaidIterator, error) {

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

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorOraclePaidIterator{contract: _OCR2Aggregator.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorOraclePaid)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseOraclePaid(log types.Log) (*OCR2AggregatorOraclePaid, error) {
	event := new(OCR2AggregatorOraclePaid)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorOwnershipTransferRequestedIterator struct {
	Event *OCR2AggregatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorOwnershipTransferRequested)
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
		it.Event = new(OCR2AggregatorOwnershipTransferRequested)
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

func (it *OCR2AggregatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2AggregatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorOwnershipTransferRequestedIterator{contract: _OCR2Aggregator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorOwnershipTransferRequested)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*OCR2AggregatorOwnershipTransferRequested, error) {
	event := new(OCR2AggregatorOwnershipTransferRequested)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorOwnershipTransferredIterator struct {
	Event *OCR2AggregatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorOwnershipTransferred)
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
		it.Event = new(OCR2AggregatorOwnershipTransferred)
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

func (it *OCR2AggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2AggregatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorOwnershipTransferredIterator{contract: _OCR2Aggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorOwnershipTransferred)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*OCR2AggregatorOwnershipTransferred, error) {
	event := new(OCR2AggregatorOwnershipTransferred)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorPayeeshipTransferRequestedIterator struct {
	Event *OCR2AggregatorPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorPayeeshipTransferRequested)
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
		it.Event = new(OCR2AggregatorPayeeshipTransferRequested)
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

func (it *OCR2AggregatorPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*OCR2AggregatorPayeeshipTransferRequestedIterator, error) {

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

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorPayeeshipTransferRequestedIterator{contract: _OCR2Aggregator.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorPayeeshipTransferRequested)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParsePayeeshipTransferRequested(log types.Log) (*OCR2AggregatorPayeeshipTransferRequested, error) {
	event := new(OCR2AggregatorPayeeshipTransferRequested)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorPayeeshipTransferredIterator struct {
	Event *OCR2AggregatorPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorPayeeshipTransferred)
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
		it.Event = new(OCR2AggregatorPayeeshipTransferred)
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

func (it *OCR2AggregatorPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*OCR2AggregatorPayeeshipTransferredIterator, error) {

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

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorPayeeshipTransferredIterator{contract: _OCR2Aggregator.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorPayeeshipTransferred)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParsePayeeshipTransferred(log types.Log) (*OCR2AggregatorPayeeshipTransferred, error) {
	event := new(OCR2AggregatorPayeeshipTransferred)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorRequesterAccessControllerSetIterator struct {
	Event *OCR2AggregatorRequesterAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorRequesterAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorRequesterAccessControllerSet)
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
		it.Event = new(OCR2AggregatorRequesterAccessControllerSet)
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

func (it *OCR2AggregatorRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*OCR2AggregatorRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorRequesterAccessControllerSetIterator{contract: _OCR2Aggregator.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorRequesterAccessControllerSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseRequesterAccessControllerSet(log types.Log) (*OCR2AggregatorRequesterAccessControllerSet, error) {
	event := new(OCR2AggregatorRequesterAccessControllerSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorRoundRequestedIterator struct {
	Event *OCR2AggregatorRoundRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorRoundRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorRoundRequested)
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
		it.Event = new(OCR2AggregatorRoundRequested)
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

func (it *OCR2AggregatorRoundRequestedIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorRoundRequested struct {
	Requester    common.Address
	ConfigDigest [32]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*OCR2AggregatorRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorRoundRequestedIterator{contract: _OCR2Aggregator.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorRoundRequested)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseRoundRequested(log types.Log) (*OCR2AggregatorRoundRequested, error) {
	event := new(OCR2AggregatorRoundRequested)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorTransmittedIterator struct {
	Event *OCR2AggregatorTransmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorTransmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorTransmitted)
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
		it.Event = new(OCR2AggregatorTransmitted)
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

func (it *OCR2AggregatorTransmittedIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorTransmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorTransmitted struct {
	ConfigDigest [32]byte
	Epoch        uint32
	Raw          types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterTransmitted(opts *bind.FilterOpts) (*OCR2AggregatorTransmittedIterator, error) {

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorTransmittedIterator{contract: _OCR2Aggregator.contract, event: "Transmitted", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorTransmitted) (event.Subscription, error) {

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "Transmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorTransmitted)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseTransmitted(log types.Log) (*OCR2AggregatorTransmitted, error) {
	event := new(OCR2AggregatorTransmitted)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "Transmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OCR2AggregatorValidatorConfigSetIterator struct {
	Event *OCR2AggregatorValidatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OCR2AggregatorValidatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OCR2AggregatorValidatorConfigSet)
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
		it.Event = new(OCR2AggregatorValidatorConfigSet)
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

func (it *OCR2AggregatorValidatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *OCR2AggregatorValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OCR2AggregatorValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*OCR2AggregatorValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &OCR2AggregatorValidatorConfigSetIterator{contract: _OCR2Aggregator.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

func (_OCR2Aggregator *OCR2AggregatorFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _OCR2Aggregator.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OCR2AggregatorValidatorConfigSet)
				if err := _OCR2Aggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
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

func (_OCR2Aggregator *OCR2AggregatorFilterer) ParseValidatorConfigSet(log types.Log) (*OCR2AggregatorValidatorConfigSet, error) {
	event := new(OCR2AggregatorValidatorConfigSet)
	if err := _OCR2Aggregator.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetBilling struct {
	MaximumGasPriceGwei       uint32
	ReasonableGasPriceGwei    uint32
	ObservationPaymentGjuels  uint32
	TransmissionPaymentGjuels uint32
	AccountingGas             *big.Int
}
type GetRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type GetValidatorConfig struct {
	Validator common.Address
	GasLimit  uint32
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

func (_OCR2Aggregator *OCR2Aggregator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OCR2Aggregator.abi.Events["AnswerUpdated"].ID:
		return _OCR2Aggregator.ParseAnswerUpdated(log)
	case _OCR2Aggregator.abi.Events["BillingAccessControllerSet"].ID:
		return _OCR2Aggregator.ParseBillingAccessControllerSet(log)
	case _OCR2Aggregator.abi.Events["BillingSet"].ID:
		return _OCR2Aggregator.ParseBillingSet(log)
	case _OCR2Aggregator.abi.Events["ConfigSet"].ID:
		return _OCR2Aggregator.ParseConfigSet(log)
	case _OCR2Aggregator.abi.Events["LinkTokenSet"].ID:
		return _OCR2Aggregator.ParseLinkTokenSet(log)
	case _OCR2Aggregator.abi.Events["NewRound"].ID:
		return _OCR2Aggregator.ParseNewRound(log)
	case _OCR2Aggregator.abi.Events["NewTransmission"].ID:
		return _OCR2Aggregator.ParseNewTransmission(log)
	case _OCR2Aggregator.abi.Events["OraclePaid"].ID:
		return _OCR2Aggregator.ParseOraclePaid(log)
	case _OCR2Aggregator.abi.Events["OwnershipTransferRequested"].ID:
		return _OCR2Aggregator.ParseOwnershipTransferRequested(log)
	case _OCR2Aggregator.abi.Events["OwnershipTransferred"].ID:
		return _OCR2Aggregator.ParseOwnershipTransferred(log)
	case _OCR2Aggregator.abi.Events["PayeeshipTransferRequested"].ID:
		return _OCR2Aggregator.ParsePayeeshipTransferRequested(log)
	case _OCR2Aggregator.abi.Events["PayeeshipTransferred"].ID:
		return _OCR2Aggregator.ParsePayeeshipTransferred(log)
	case _OCR2Aggregator.abi.Events["RequesterAccessControllerSet"].ID:
		return _OCR2Aggregator.ParseRequesterAccessControllerSet(log)
	case _OCR2Aggregator.abi.Events["RoundRequested"].ID:
		return _OCR2Aggregator.ParseRoundRequested(log)
	case _OCR2Aggregator.abi.Events["Transmitted"].ID:
		return _OCR2Aggregator.ParseTransmitted(log)
	case _OCR2Aggregator.abi.Events["ValidatorConfigSet"].ID:
		return _OCR2Aggregator.ParseValidatorConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OCR2AggregatorAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (OCR2AggregatorBillingAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912")
}

func (OCR2AggregatorBillingSet) Topic() common.Hash {
	return common.HexToHash("0x0bf184bf1bba9699114bdceddaf338a1b364252c5e497cc01918dde92031713f")
}

func (OCR2AggregatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (OCR2AggregatorLinkTokenSet) Topic() common.Hash {
	return common.HexToHash("0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a")
}

func (OCR2AggregatorNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (OCR2AggregatorNewTransmission) Topic() common.Hash {
	return common.HexToHash("0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a")
}

func (OCR2AggregatorOraclePaid) Topic() common.Hash {
	return common.HexToHash("0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c")
}

func (OCR2AggregatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OCR2AggregatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OCR2AggregatorPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (OCR2AggregatorPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (OCR2AggregatorRequesterAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634")
}

func (OCR2AggregatorRoundRequested) Topic() common.Hash {
	return common.HexToHash("0x41e3990591fd372502daa15842da15bc7f41c75309ab3ff4f56f1848c178825c")
}

func (OCR2AggregatorTransmitted) Topic() common.Hash {
	return common.HexToHash("0xb04e63db38c49950639fa09d29872f21f5d49d614f3a969d8adf3d4b52e41a62")
}

func (OCR2AggregatorValidatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541")
}

func (_OCR2Aggregator *OCR2Aggregator) Address() common.Address {
	return _OCR2Aggregator.address
}

type OCR2AggregatorInterface interface {
	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	GetBilling(opts *bind.CallOpts) (GetBilling,

		error)

	GetBillingAccessController(opts *bind.CallOpts) (common.Address, error)

	GetLinkToken(opts *bind.CallOpts) (common.Address, error)

	GetRequesterAccessController(opts *bind.CallOpts) (common.Address, error)

	GetRoundData(opts *bind.CallOpts, roundId *big.Int) (GetRoundData,

		error)

	GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	GetTransmitters(opts *bind.CallOpts) ([]common.Address, error)

	GetValidatorConfig(opts *bind.CallOpts) (GetValidatorConfig,

		error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestConfigDetails(opts *bind.CallOpts) (LatestConfigDetails,

		error)

	LatestConfigDigestAndEpoch(opts *bind.CallOpts) (LatestConfigDigestAndEpoch,

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

	OracleObservationCount(opts *bind.CallOpts, transmitterAddress common.Address) (uint32, error)

	OwedPayment(opts *bind.CallOpts, transmitterAddress common.Address) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptPayeeship(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	RequestNewRound(opts *bind.TransactOpts) (*types.Transaction, error)

	SetBilling(opts *bind.TransactOpts, maximumGasPriceGwei uint32, reasonableGasPriceGwei uint32, observationPaymentGjuels uint32, transmissionPaymentGjuels uint32, accountingGas *big.Int) (*types.Transaction, error)

	SetBillingAccessController(opts *bind.TransactOpts, _billingAccessController common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, signers []common.Address, transmitters []common.Address, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte) (*types.Transaction, error)

	SetLinkToken(opts *bind.TransactOpts, linkToken common.Address, recipient common.Address) (*types.Transaction, error)

	SetPayees(opts *bind.TransactOpts, transmitters []common.Address, payees []common.Address) (*types.Transaction, error)

	SetRequesterAccessController(opts *bind.TransactOpts, requesterAccessController common.Address) (*types.Transaction, error)

	SetValidatorConfig(opts *bind.TransactOpts, newValidator common.Address, newGasLimit uint32) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	TransferPayeeship(opts *bind.TransactOpts, transmitter common.Address, proposed common.Address) (*types.Transaction, error)

	Transmit(opts *bind.TransactOpts, reportContext [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error)

	WithdrawFunds(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	WithdrawPayment(opts *bind.TransactOpts, transmitter common.Address) (*types.Transaction, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*OCR2AggregatorAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*OCR2AggregatorAnswerUpdated, error)

	FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OCR2AggregatorBillingAccessControllerSetIterator, error)

	WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorBillingAccessControllerSet) (event.Subscription, error)

	ParseBillingAccessControllerSet(log types.Log) (*OCR2AggregatorBillingAccessControllerSet, error)

	FilterBillingSet(opts *bind.FilterOpts) (*OCR2AggregatorBillingSetIterator, error)

	WatchBillingSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorBillingSet) (event.Subscription, error)

	ParseBillingSet(log types.Log) (*OCR2AggregatorBillingSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OCR2AggregatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OCR2AggregatorConfigSet, error)

	FilterLinkTokenSet(opts *bind.FilterOpts, oldLinkToken []common.Address, newLinkToken []common.Address) (*OCR2AggregatorLinkTokenSetIterator, error)

	WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorLinkTokenSet, oldLinkToken []common.Address, newLinkToken []common.Address) (event.Subscription, error)

	ParseLinkTokenSet(log types.Log) (*OCR2AggregatorLinkTokenSet, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*OCR2AggregatorNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*OCR2AggregatorNewRound, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*OCR2AggregatorNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*OCR2AggregatorNewTransmission, error)

	FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*OCR2AggregatorOraclePaidIterator, error)

	WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error)

	ParseOraclePaid(log types.Log) (*OCR2AggregatorOraclePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2AggregatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OCR2AggregatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OCR2AggregatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OCR2AggregatorOwnershipTransferred, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*OCR2AggregatorPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*OCR2AggregatorPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*OCR2AggregatorPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*OCR2AggregatorPayeeshipTransferred, error)

	FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*OCR2AggregatorRequesterAccessControllerSetIterator, error)

	WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorRequesterAccessControllerSet) (event.Subscription, error)

	ParseRequesterAccessControllerSet(log types.Log) (*OCR2AggregatorRequesterAccessControllerSet, error)

	FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*OCR2AggregatorRoundRequestedIterator, error)

	WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorRoundRequested, requester []common.Address) (event.Subscription, error)

	ParseRoundRequested(log types.Log) (*OCR2AggregatorRoundRequested, error)

	FilterTransmitted(opts *bind.FilterOpts) (*OCR2AggregatorTransmittedIterator, error)

	WatchTransmitted(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorTransmitted) (event.Subscription, error)

	ParseTransmitted(log types.Log) (*OCR2AggregatorTransmitted, error)

	FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*OCR2AggregatorValidatorConfigSetIterator, error)

	WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *OCR2AggregatorValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error)

	ParseValidatorConfigSet(log types.Log) (*OCR2AggregatorValidatorConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
