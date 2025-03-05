// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package data_feeds_cache

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

type DataFeedsCacheWorkflowMetadata struct {
	AllowedSender        common.Address
	AllowedWorkflowOwner common.Address
	AllowedWorkflowName  [10]byte
}

var DataFeedsCacheMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bundleDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"bundleFeedDecimals\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkFeedPermission\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"outputs\":[{\"name\":\"hasPermission\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"feedDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"feedDescription\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAnswer\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBundleDecimals\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"bundleFeedDecimals\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDataIdForProxy\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDecimals\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"feedDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getDescription\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"feedDescription\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeedMetadata\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestAnswer\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestBundle\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"bundle\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestBundleTimestamp\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestRoundData\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestTimestamp\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimestamp\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isFeedAdmin\",\"inputs\":[{\"name\":\"feedAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestAnswer\",\"inputs\":[],\"outputs\":[{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestBundle\",\"inputs\":[],\"outputs\":[{\"name\":\"bundle\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestBundleTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRound\",\"inputs\":[],\"outputs\":[{\"name\":\"round\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"id\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverTokens\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeDataIdMappingsForProxies\",\"inputs\":[{\"name\":\"proxies\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeFeedConfigs\",\"inputs\":[{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBundleFeedConfigs\",\"inputs\":[{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"},{\"name\":\"descriptions\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"decimalsMatrix\",\"type\":\"uint8[][]\",\"internalType\":\"uint8[][]\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDecimalFeedConfigs\",\"inputs\":[{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"},{\"name\":\"descriptions\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFeedAdmin\",\"inputs\":[{\"name\":\"feedAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isAdmin\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateDataIdMappingsForProxies\",\"inputs\":[{\"name\":\"proxies\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AnswerUpdated\",\"inputs\":[{\"name\":\"current\",\"type\":\"int256\",\"indexed\":true,\"internalType\":\"int256\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BundleFeedConfigSet\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"decimals\",\"type\":\"uint8[]\",\"indexed\":false,\"internalType\":\"uint8[]\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BundleReportUpdated\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"bundle\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DecimalFeedConfigSet\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DecimalReportUpdated\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"answer\",\"type\":\"uint224\",\"indexed\":false,\"internalType\":\"uint224\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeedAdminSet\",\"inputs\":[{\"name\":\"feedAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"isAdmin\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeedConfigRemoved\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidUpdatePermission\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"workflowOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"workflowName\",\"type\":\"bytes10\",\"indexed\":false,\"internalType\":\"bytes10\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewRound\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProxyDataIdRemoved\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProxyDataIdUpdated\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleBundleReport\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"reportTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"latestTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleDecimalReport\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"reportTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"latestTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRecovered\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressInsufficientBalance\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ErrorSendingNative\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeedNotConfigured\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requiredBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidDataId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidWorkflowName\",\"inputs\":[{\"name\":\"workflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]},{\"type\":\"error\",\"name\":\"NoMappingForSender\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b614a5b806101576000396000f3fe608060405234801561001057600080fd5b50600436106102925760003560e01c806379ba509711610160578063b5ab58dc116100d8578063ec52b1f51161008c578063feaf968c11610071578063feaf968c146105d3578063feb5d172146105db578063ff25dbc81461067157600080fd5b8063ec52b1f5146105ad578063f2fde38b146105c057600080fd5b8063be4f0a9f116100bd578063be4f0a9f14610567578063cdd251001461057a578063d143dcd91461058d57600080fd5b8063b5ab58dc14610541578063b633620c1461055457600080fd5b80639198274f1161012f5780639a6fc8f5116101145780639a6fc8f51461051e5780639d91348d14610531578063a3d610cc1461053957600080fd5b80639198274f146104d15780639608e18f146104d957600080fd5b806379ba509714610493578063805f21321461049b5780638205bf6a146104ae5780638da5cb5b146104b657600080fd5b80634533dc981161020e5780635f25452b116101c2578063668a0f02116101a7578063668a0f02146104705780636a36e494146104785780637284e4161461048b57600080fd5b80635f25452b146104135780635f3e849f1461045d57600080fd5b806350d25bcd116101f357806350d25bcd146103f057806354fd4d50146103f8578063557a33c21461040057600080fd5b80634533dc98146103bd57806347381b08146103d057600080fd5b8063297dbf561161026557806335f611221161024a57806335f611221461036b5780633a0449741461037e57806343d5ba50146103aa57600080fd5b8063297dbf561461033c578063313ce5671461035157600080fd5b806301ffc9a71461029757806302ccb3ae146102bf578063181f5a77146102df5780631bb1610c1461031b575b600080fd5b6102aa6102a5366004613a2a565b610684565b60405190151581526020015b60405180910390f35b6102d26102cd366004613a89565b610801565b6040516102b69190613af4565b6102d26040518060400160405280601481526020017f446174614665656473436163686520312e302e3000000000000000000000000081525081565b61032e610329366004613a89565b6108f0565b6040519081526020016102b6565b61034f61034a366004613b53565b610976565b005b610359610b2b565b60405160ff90911681526020016102b6565b61034f610379366004613c09565b610b90565b6102aa61038c366004613cf5565b6001600160a01b031660009081526007602052604090205460ff1690565b6103596103b8366004613a89565b611209565b61034f6103cb366004613d12565b611255565b6103e36103de366004613a89565b611873565b6040516102b69190613db8565b61032e61193d565b61032e600781565b61032e61040e366004613a89565b6119cf565b610426610421366004613a89565b611a4d565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a0016102b6565b61034f61046b366004613dfe565b611b11565b61032e611db2565b6102d2610486366004613a89565b611e26565b6102d2611e8d565b61034f611f91565b61034f6104a9366004613e81565b612074565b61032e612782565b6000546040516001600160a01b0390911681526020016102b6565b6102d261281c565b6105056104e7366004613cf5565b6001600160a01b031660009081526002602052604090205460801b90565b6040516001600160801b031990911681526020016102b6565b61042661052c366004613ee6565b612899565b6103e361297e565b61032e612a5e565b61032e61054f366004613f12565b612adb565b61032e610562366004613f12565b612b7a565b61034f610575366004613f2b565b612c22565b61034f610588366004613f2b565b612d10565b6105a061059b366004613f6d565b612fa6565b6040516102b69190613fa0565b61034f6105bb36600461402d565b6131ce565b61034f6105ce366004613cf5565b613275565b610426613289565b6102aa6105e9366004614169565b805160208083015160409384015184516001600160801b031996909616868401526001600160a01b03938416868601529216606085015275ffffffffffffffffffffffffffffffffffffffffffff199091166080808501919091528251808503909101815260a090930182528251928101929092206000908152600990925290205460ff1690565b61032e61067f366004613a89565b613362565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fcce8054600000000000000000000000000000000000000000000000000000000148061071757507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b8061076357507fffffffff0000000000000000000000000000000000000000000000000000000082167f805f213200000000000000000000000000000000000000000000000000000000145b806107af57507fffffffff0000000000000000000000000000000000000000000000000000000082167f5f3e849f00000000000000000000000000000000000000000000000000000000145b806107fb57507fffffffff0000000000000000000000000000000000000000000000000000000082167f181f5a7700000000000000000000000000000000000000000000000000000000145b92915050565b60606001600160801b03198216610844576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031982166000908152600860205260409020600101805461086b9061419d565b80601f01602080910402602001604051908101604052809291908181526020018280546108979061419d565b80156108e45780601f106108b9576101008083540402835291602001916108e4565b820191906000526020600020905b8154815290600101906020018083116108c757829003601f168201915b50505050509050919050565b60006001600160801b03198216610933576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506001600160801b0319166000908152600360205260409020547c0100000000000000000000000000000000000000000000000000000000900463ffffffff1690565b3360009081526007602052604090205460ff166109c6576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b82818114610a00576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b81811015610b2357838382818110610a1d57610a1d6141f0565b9050602002016020810190610a329190613a89565b60026000888885818110610a4857610a486141f0565b9050602002016020810190610a5d9190613cf5565b6001600160a01b03168152602081019190915260400160002080546001600160801b03191660809290921c919091179055838382818110610aa057610aa06141f0565b9050602002016020810190610ab59190613a89565b6001600160801b031916868683818110610ad157610ad16141f0565b9050602002016020810190610ae69190613cf5565b6001600160a01b03167ff31b9e58190970ef07c23d0ba78c358eb3b416e829ef484b29b9993a6b1b285a60405160405180910390a3600101610a03565b505050505050565b3360009081526002602052604081205460801b6001600160801b03198116610b81576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b610b8a816133cb565b91505090565b3360009081526007602052604090205460ff16610bdb576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b801580610be6575086155b15610c1d576040517f60e8b63a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8685141580610c2c5750868314155b15610c63576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b878110156111fe576000898983818110610c8257610c826141f0565b9050602002016020810190610c979190613a89565b90506001600160801b03198116610cda576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031981166000908152600860205260409020600281015415610e755760005b6002820154811015610dfc576000826002018281548110610d2357610d236141f0565b60009182526020808320604080516060808201835260029590950290920180546001600160a01b039081168085526001909201549081168486018190527401000000000000000000000000000000000000000090910460b01b75ffffffffffffffffffffffffffffffffffffffffffff191684840181905283516001600160801b03198d168188015280850193909352958201526080808201959095528151808203909501855260a0019052825192909101919091209092506000908152600960205260409020805460ff191690555050600101610d00565b506001600160801b03198216600090815260086020526040812090610e21828261388b565b610e2f6001830160006138b0565b610e3d6002830160006138ea565b50506040516001600160801b03198316907f871bcdef10dee59b87f17bab788b72faa8dfe1a9cc5bdc45c3baf4c18fa3391090600090a25b60005b848110156110fe576000868683818110610e9457610e946141f0565b905060600201803603810190610eaa919061421f565b905084600003610fcd5780516001600160a01b0316610f035780516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109bd565b60208101516001600160a01b0316610f585760208101516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109bd565b604081015175ffffffffffffffffffffffffffffffffffffffffffff1916610fcd5760408082015190517f114988d500000000000000000000000000000000000000000000000000000000815275ffffffffffffffffffffffffffffffffffffffffffff1990911660048201526024016109bd565b8051602080830180516040808601805182516001600160801b03198c16818801526001600160a01b0397881681850152938716606085015275ffffffffffffffffffffffffffffffffffffffffffff19166080808501919091528251808503909101815260a09093018252825192850192909220600090815260098552908120805460ff191660019081179091556002898101805480840182559084529590922096519490910290950180549385167fffffffffffffffffffffffff000000000000000000000000000000000000000090941693909317835590519184018054915160b01c74010000000000000000000000000000000000000000027fffff000000000000000000000000000000000000000000000000000000000000909216929093169190911717905501610e78565b50868684818110611111576111116141f0565b9050602002810190611123919061423b565b61112e91839161390b565b50888884818110611141576111416141f0565b905060200281019061115391906142a3565b6001830191611163919083614356565b506001600160801b031982167fdfebe0878c5611549f54908260ca12271c7ff3f0ebae0c1de47732612403869e8888868181106111a2576111a26141f0565b90506020028101906111b4919061423b565b8c8c888181106111c6576111c66141f0565b90506020028101906111d891906142a3565b8a8a6040516111ec969594939291906144d1565b60405180910390a25050600101610c66565b505050505050505050565b60006001600160801b0319821661124c576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107fb826133cb565b3360009081526007602052604090205460ff166112a0576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b8015806112ab575084155b156112e2576040517f60e8b63a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84831461131b576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8581101561186a57600087878381811061133a5761133a6141f0565b905060200201602081019061134f9190613a89565b90506001600160801b03198116611392576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b03198116600090815260086020526040902060028101541561152d5760005b60028201548110156114b45760008260020182815481106113db576113db6141f0565b60009182526020808320604080516060808201835260029590950290920180546001600160a01b039081168085526001909201549081168486018190527401000000000000000000000000000000000000000090910460b01b75ffffffffffffffffffffffffffffffffffffffffffff191684840181905283516001600160801b03198d168188015280850193909352958201526080808201959095528151808203909501855260a0019052825192909101919091209092506000908152600960205260409020805460ff1916905550506001016113b8565b506001600160801b031982166000908152600860205260408120906114d9828261388b565b6114e76001830160006138b0565b6114f56002830160006138ea565b50506040516001600160801b03198316907f871bcdef10dee59b87f17bab788b72faa8dfe1a9cc5bdc45c3baf4c18fa3391090600090a25b60005b848110156117b657600086868381811061154c5761154c6141f0565b905060600201803603810190611562919061421f565b9050846000036116855780516001600160a01b03166115bb5780516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109bd565b60208101516001600160a01b03166116105760208101516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109bd565b604081015175ffffffffffffffffffffffffffffffffffffffffffff19166116855760408082015190517f114988d500000000000000000000000000000000000000000000000000000000815275ffffffffffffffffffffffffffffffffffffffffffff1990911660048201526024016109bd565b8051602080830180516040808601805182516001600160801b03198c16818801526001600160a01b0397881681850152938716606085015275ffffffffffffffffffffffffffffffffffffffffffff19166080808501919091528251808503909101815260a09093018252825192850192909220600090815260098552908120805460ff191660019081179091556002898101805480840182559084529590922096519490910290950180549385167fffffffffffffffffffffffff000000000000000000000000000000000000000090941693909317835590519184018054915160b01c74010000000000000000000000000000000000000000027fffff000000000000000000000000000000000000000000000000000000000000909216929093169190911717905501611530565b508686848181106117c9576117c96141f0565b90506020028101906117db91906142a3565b60018301916117eb919083614356565b506001600160801b031982167f2dec0e9ffbb18c6499fc8bee8b9c35f765e76d9dbd436f25dd00a80de267ac0d611821846133cb565b898987818110611833576118336141f0565b905060200281019061184591906142a3565b898960405161185895949392919061454a565b60405180910390a2505060010161131e565b50505050505050565b60606001600160801b031982166118b6576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b03198216600090815260086020908152604091829020805483518184028101840190945280845290918301828280156108e457602002820191906000526020600020906000905b825461010083900a900460ff16815260206001928301818104948501949093039092029101808411611904575094979650505050505050565b3360009081526002602052604081205460801b6001600160801b03198116611993576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b0319166000908152600360205260409020547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b60006001600160801b03198216611a12576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506001600160801b0319166000908152600360205260409020547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b6000808080806001600160801b03198616611a94576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050506001600160801b03199190911660009081526006602090815260408083205460039092529091205490927bffffffffffffffffffffffffffffffffffffffffffffffffffffffff821692507c010000000000000000000000000000000000000000000000000000000090910463ffffffff169081908490565b611b1961348c565b6001600160a01b038316611c065747811115611b6a576040517fcf479181000000000000000000000000000000000000000000000000000000008152476004820152602481018290526044016109bd565b600080836001600160a01b03168360405160006040518083038185875af1925050503d8060008114611bb8576040519150601f19603f3d011682016040523d82523d6000602084013e611bbd565b606091505b509150915081611bff578383826040517fc50febed0000000000000000000000000000000000000000000000000000000081526004016109bd93929190614586565b5050611d60565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b038416906370a0823190602401602060405180830381865afa158015611c63573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c8791906145b7565b811115611d4c576040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b038416906370a0823190602401602060405180830381865afa158015611ceb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d0f91906145b7565b6040517fcf4791810000000000000000000000000000000000000000000000000000000081526004810191909152602481018290526044016109bd565b611d606001600160a01b0384168383613502565b816001600160a01b0316836001600160a01b03167f879f92dded0f26b83c3e00b12e0395dc72cfc3077343d1854ed6988edd1f909683604051611da591815260200190565b60405180910390a3505050565b3360009081526002602052604081205460801b6001600160801b03198116611e08576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b031916600090815260066020526040902054919050565b60606001600160801b03198216611e69576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031982166000908152600560205260409020805461086b9061419d565b3360009081526002602052604090205460609060801b6001600160801b03198116611ee6576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b0319811660009081526008602052604090206001018054611f0d9061419d565b80601f0160208091040260200160405190810160405280929190818152602001828054611f399061419d565b8015611f865780601f10611f5b57610100808354040283529160200191611f86565b820191906000526020600020905b815481529060010190602001808311611f6957829003601f168201915b505050505091505090565b6001546001600160a01b03163314612005576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016109bd565b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000806120b686868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061358292505050565b909250905060006120cb6040602086886145d0565b6120d4916145fa565b90506120e1816060614647565b6120ec90604061465e565b84036124dd576000612100858701876146a9565b905060005b828110156124d6576000828281518110612121576121216141f0565b6020908102919091018101518051604080516001600160801b031983168186015233818301526001600160a01b038b16606082015275ffffffffffffffffffffffffffffffffffffffffffff198a166080808301919091528251808303909101815260a0909101825280519085012060008181526009909552932054919350919060ff1661222257604080513381526001600160a01b038a16602082015275ffffffffffffffffffffffffffffffffffffffffffff198916918101919091526001600160801b03198316907feeeaa8bf618ff6d960c6cf5935e68384f066abcc8b95d0de91bd773c16ae3ae3906060015b60405180910390a25050506124ce565b6001600160801b031982166000908152600360209081526040909120549084015163ffffffff7c010000000000000000000000000000000000000000000000000000000090920482169116116122f3576020838101516001600160801b0319841660008181526003845260409081902054815163ffffffff94851681527c010000000000000000000000000000000000000000000000000000000090910490931693830193909352917fcf16f5f704f981fa2279afa1877dd1fdaa462a03a71ec51b9d3b2416a59a013e9101612212565b604080518082018252848201517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260208086015163ffffffff16818301526001600160801b0319851660009081526006909152918220805491929182906123589061479c565b91829055506001600160801b0319851660008181526003602090815260408083208751888401805163ffffffff9081167c01000000000000000000000000000000000000000000000000000000009081027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff94851617909455888752600486528487208888528652958490208a519151909616928302911690811790945590519283529394508492917f82584589cd7284d4503ed582275e22b2e8f459f9cf4170a7235844e367f966d5910160405180910390a460208086015160405163ffffffff909116815260009183917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271910160405180910390a38085604001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f426040516124c091815260200190565b60405180910390a350505050505b600101612105565b505061186a565b60006124eb858701876147b6565b905060005b81518110156111fe57600082828151811061250d5761250d6141f0565b6020908102919091018101518051604080516001600160801b031983168186015233818301526001600160a01b038b16606082015275ffffffffffffffffffffffffffffffffffffffffffff198a166080808301919091528251808303909101815260a0909101825280519085012060008181526009909552932054919350919060ff1661260e57604080513381526001600160a01b038a16602082015275ffffffffffffffffffffffffffffffffffffffffffff198916918101919091526001600160801b03198316907feeeaa8bf618ff6d960c6cf5935e68384f066abcc8b95d0de91bd773c16ae3ae3906060015b60405180910390a250505061277a565b6001600160801b031982166000908152600560209081526040909120600101549084015163ffffffff9182169116116126a3576020838101516001600160801b0319841660008181526005845260409081902060010154815163ffffffff9485168152931693830193909352917f51001b67094834cc084a0c1feb791cf84a481357aa66b924ba205d4cb56fd98191016125fe565b60408051808201825284820151815260208086015163ffffffff16818301526001600160801b031985166000908152600590915291909120815182919081906126ec908261492a565b5060209182015160019190910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff92831617905590820151825160405191909216916001600160801b03198616917f1dc1bef0b59d624eab3f0ec044781bb5b8594cd64f0ba09d789f5b51acab16149161276d91613af4565b60405180910390a3505050505b6001016124f0565b3360009081526002602052604081205460801b6001600160801b031981166127d8576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b0319166000908152600360205260409020547c0100000000000000000000000000000000000000000000000000000000900463ffffffff16919050565b3360009081526002602052604090205460609060801b6001600160801b03198116612875576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b0319811660009081526005602052604090208054611f0d9061419d565b33600090815260026020526040812054819081908190819060801b6001600160801b031981166128f7576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b69ffffffffffffffffffff871660009081526004602090815260408083206001600160801b0319949094168352929052205495967bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8716967c0100000000000000000000000000000000000000000000000000000000900463ffffffff169550859450879350915050565b3360009081526002602052604090205460609060801b6001600160801b031981166129d7576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b0319811660009081526008602090815260409182902080548351818402810184019094528084529091830182828015611f8657602002820191906000526020600020906000905b825461010083900a900460ff16815260206001928301818104948501949093039092029101808411612a25579050505050505091505090565b3360009081526002602052604081205460801b6001600160801b03198116612ab4576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b03191660009081526005602052604090206001015463ffffffff16919050565b3360009081526002602052604081205460801b6001600160801b03198116612b31576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b60009283526004602090815260408085206001600160801b03199093168552919052909120547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b3360009081526002602052604081205460801b6001600160801b03198116612bd0576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b60009283526004602090815260408085206001600160801b031990931685529190529091205463ffffffff7c010000000000000000000000000000000000000000000000000000000090910416919050565b3360009081526007602052604090205460ff16612c6d576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b8060005b81811015612d0a576000848483818110612c8d57612c8d6141f0565b9050602002016020810190612ca29190613cf5565b6001600160a01b03811660008181526002602052604080822080546001600160801b0319808216909255915194955060809190911b9390841692917f4200186b7bc2d4f13f7888c5bbe9461d57da88705be86521f3d78be691ad1d2a91a35050600101612c71565b50505050565b3360009081526007602052604090205460ff16612d5b576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b60005b81811015612fa1576000838383818110612d7a57612d7a6141f0565b9050602002016020810190612d8f9190613a89565b6001600160801b0319811660009081526008602052604081206002015491925003612df2576040517f8606a85b0000000000000000000000000000000000000000000000000000000081526001600160801b0319821660048201526024016109bd565b60005b6001600160801b03198216600090815260086020526040902060020154811015612f20576001600160801b031982166000908152600860205260408120600201805483908110612e4757612e476141f0565b60009182526020808320604080516060808201835260029590950290920180546001600160a01b039081168085526001909201549081168486018190527401000000000000000000000000000000000000000090910460b01b75ffffffffffffffffffffffffffffffffffffffffffff191684840181905283516001600160801b03198c168188015280850193909352958201526080808201959095528151808203909501855260a0019052825192909101919091209092506000908152600960205260409020805460ff191690555050600101612df5565b506001600160801b03198116600090815260086020526040812090612f45828261388b565b612f536001830160006138b0565b612f616002830160006138ea565b50506040516001600160801b03198216907f871bcdef10dee59b87f17bab788b72faa8dfe1a9cc5bdc45c3baf4c18fa3391090600090a250600101612d5e565b505050565b6001600160801b031983166000908152600860205260408120600281015460609281900361300c576040517f8606a85b0000000000000000000000000000000000000000000000000000000081526001600160801b0319871660048201526024016109bd565b808510613060576040805160008082526020820190925290613056565b60408051606081018252600080825260208083018290529282015282526000199092019101816130295790505b50925050506131c7565b600061306c858761465e565b90508181118061307a575084155b6130845780613086565b815b905061309286826149e9565b67ffffffffffffffff8111156130aa576130aa614066565b6040519080825280602002602001820160405280156130f557816020015b60408051606081018252600080825260208083018290529282015282526000199092019101816130c85790505b50935060005b84518110156131c25760028401613112888361465e565b81548110613122576131226141f0565b60009182526020918290206040805160608101825260029390930290910180546001600160a01b039081168452600190910154908116938301939093527401000000000000000000000000000000000000000090920460b01b75ffffffffffffffffffffffffffffffffffffffffffff19169181019190915285518690839081106131af576131af6141f0565b60209081029190910101526001016130fb565b505050505b9392505050565b6131d661348c565b6001600160a01b038216613221576040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b03831660048201526024016109bd565b6001600160a01b038216600081815260076020526040808220805460ff191685151590811790915590519092917f93a3fa5993d2a54de369386625330cc6d73caee7fece4b3983cf299b264473fd91a35050565b61327d61348c565b61328681613593565b50565b33600090815260026020526040812054819081908190819060801b6001600160801b031981166132e7576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109bd565b6001600160801b03191660009081526006602090815260408083205460039092529091205490967bffffffffffffffffffffffffffffffffffffffffffffffffffffffff821696507c010000000000000000000000000000000000000000000000000000000090910463ffffffff1694508493508692509050565b60006001600160801b031982166133a5576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506001600160801b03191660009081526005602052604090206001015463ffffffff1690565b6000806133d983600761366e565b90507f20000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000082161080159061346f57507f60000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821611155b15613483576131c7602060f883901c6149fc565b50600092915050565b6000546001600160a01b03163314613500576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016109bd565b565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052612fa19084906136d6565b6040810151604a9091015160601c91565b336001600160a01b03821603613605576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016109bd565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040516001600160801b03198316602082015260009060300160405160208183030381529060405282815181106136a7576136a76141f0565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016905092915050565b60006136eb6001600160a01b03841683613752565b9050805160001415801561371057508080602001905181019061370e9190614a15565b155b15612fa1576040517f5274afe70000000000000000000000000000000000000000000000000000000081526001600160a01b03841660048201526024016109bd565b60606131c78383600084600080856001600160a01b031684866040516137789190614a32565b60006040518083038185875af1925050503d80600081146137b5576040519150601f19603f3d011682016040523d82523d6000602084013e6137ba565b606091505b50915091506137ca8683836137d4565b9695505050505050565b6060826137e9576137e482613849565b6131c7565b815115801561380057506001600160a01b0384163b155b15613842576040517f9996b3150000000000000000000000000000000000000000000000000000000081526001600160a01b03851660048201526024016109bd565b50806131c7565b8051156138595780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50805460008255601f01602090049060005260206000209081019061328691906139b4565b5080546138bc9061419d565b6000825580601f106138cc575050565b601f01602090049060005260206000209081019061328691906139b4565b508054600082556002029060005260206000209081019061328691906139c9565b82805482825590600052602060002090601f016020900481019282156139a45791602002820160005b8382111561397557833560ff1683826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302613934565b80156139a25782816101000a81549060ff0219169055600101602081600001049283019260010302613975565b505b506139b09291506139b4565b5090565b5b808211156139b057600081556001016139b5565b5b808211156139b05780547fffffffffffffffffffffffff00000000000000000000000000000000000000001681556001810180547fffff0000000000000000000000000000000000000000000000000000000000001690556002016139ca565b600060208284031215613a3c57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146131c757600080fd5b80356001600160801b031981168114613a8457600080fd5b919050565b600060208284031215613a9b57600080fd5b6131c782613a6c565b60005b83811015613abf578181015183820152602001613aa7565b50506000910152565b60008151808452613ae0816020860160208601613aa4565b601f01601f19169290920160200192915050565b6020815260006131c76020830184613ac8565b60008083601f840112613b1957600080fd5b50813567ffffffffffffffff811115613b3157600080fd5b6020830191508360208260051b8501011115613b4c57600080fd5b9250929050565b60008060008060408587031215613b6957600080fd5b843567ffffffffffffffff811115613b8057600080fd5b613b8c87828801613b07565b909550935050602085013567ffffffffffffffff811115613bac57600080fd5b613bb887828801613b07565b95989497509550505050565b60008083601f840112613bd657600080fd5b50813567ffffffffffffffff811115613bee57600080fd5b602083019150836020606083028501011115613b4c57600080fd5b6000806000806000806000806080898b031215613c2557600080fd5b883567ffffffffffffffff811115613c3c57600080fd5b613c488b828c01613b07565b909950975050602089013567ffffffffffffffff811115613c6857600080fd5b613c748b828c01613b07565b909750955050604089013567ffffffffffffffff811115613c9457600080fd5b613ca08b828c01613b07565b909550935050606089013567ffffffffffffffff811115613cc057600080fd5b613ccc8b828c01613bc4565b999c989b5096995094979396929594505050565b6001600160a01b038116811461328657600080fd5b600060208284031215613d0757600080fd5b81356131c781613ce0565b60008060008060008060608789031215613d2b57600080fd5b863567ffffffffffffffff811115613d4257600080fd5b613d4e89828a01613b07565b909750955050602087013567ffffffffffffffff811115613d6e57600080fd5b613d7a89828a01613b07565b909550935050604087013567ffffffffffffffff811115613d9a57600080fd5b613da689828a01613bc4565b979a9699509497509295939492505050565b602080825282518282018190526000918401906040840190835b81811015613df357835160ff16835260209384019390920191600101613dd2565b509095945050505050565b600080600060608486031215613e1357600080fd5b8335613e1e81613ce0565b92506020840135613e2e81613ce0565b929592945050506040919091013590565b60008083601f840112613e5157600080fd5b50813567ffffffffffffffff811115613e6957600080fd5b602083019150836020828501011115613b4c57600080fd5b60008060008060408587031215613e9757600080fd5b843567ffffffffffffffff811115613eae57600080fd5b613eba87828801613e3f565b909550935050602085013567ffffffffffffffff811115613eda57600080fd5b613bb887828801613e3f565b600060208284031215613ef857600080fd5b813569ffffffffffffffffffff811681146131c757600080fd5b600060208284031215613f2457600080fd5b5035919050565b60008060208385031215613f3e57600080fd5b823567ffffffffffffffff811115613f5557600080fd5b613f6185828601613b07565b90969095509350505050565b600080600060608486031215613f8257600080fd5b613f8b84613a6c565b95602085013595506040909401359392505050565b602080825282518282018190526000918401906040840190835b81811015613df35783516001600160a01b0381511684526001600160a01b03602082015116602085015275ffffffffffffffffffffffffffffffffffffffffffff19604082015116604085015250606083019250602084019350600181019050613fba565b801515811461328657600080fd5b6000806040838503121561404057600080fd5b823561404b81613ce0565b9150602083013561405b8161401f565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156140b8576140b8614066565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156140e7576140e7614066565b604052919050565b803575ffffffffffffffffffffffffffffffffffffffffffff1981168114613a8457600080fd5b60006060828403121561412857600080fd5b614130614095565b9050813561413d81613ce0565b8152602082013561414d81613ce0565b602082015261415e604083016140ef565b604082015292915050565b6000806080838503121561417c57600080fd5b61418583613a6c565b91506141948460208501614116565b90509250929050565b600181811c908216806141b157607f821691505b6020821081036141ea577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006060828403121561423157600080fd5b6131c78383614116565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261427057600080fd5b83018035915067ffffffffffffffff82111561428b57600080fd5b6020019150600581901b3603821315613b4c57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126142d857600080fd5b83018035915067ffffffffffffffff8211156142f357600080fd5b602001915036819003821315613b4c57600080fd5b601f821115612fa157806000526020600020601f840160051c8101602085101561432f5750805b601f840160051c820191505b8181101561434f576000815560010161433b565b5050505050565b67ffffffffffffffff83111561436e5761436e614066565b6143828361437c835461419d565b83614308565b6000601f8411600181146143b6576000851561439e5750838201355b600019600387901b1c1916600186901b17835561434f565b600083815260209020601f19861690835b828110156143e757868501358255602094850194600190920191016143c7565b50868210156144045760001960f88860031b161c19848701351681555b505060018560011b0183555050505050565b818352818160208501375060006020828401015260006020601f19601f840116840101905092915050565b81835260208301925060008160005b848110156144c757813561446381613ce0565b6001600160a01b03168652602082013561447c81613ce0565b6001600160a01b0316602087015275ffffffffffffffffffffffffffffffffffffffffffff196144ae604084016140ef565b1660408701526060958601959190910190600101614450565b5093949350505050565b6060808252810186905260008760808301825b8981101561451357823560ff81168082146144fe57600080fd5b835250602092830192909101906001016144e4565b50838103602085015261452781888a614416565b915050828103604084015261453d818587614441565b9998505050505050505050565b60ff86168152606060208201526000614567606083018688614416565b828103604084015261457a818587614441565b98975050505050505050565b6001600160a01b03841681528260208201526060604082015260006145ae6060830184613ac8565b95945050505050565b6000602082840312156145c957600080fd5b5051919050565b600080858511156145e057600080fd5b838611156145ed57600080fd5b5050820193919092039150565b803560208310156107fb57600019602084900360031b1b1692915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176107fb576107fb614618565b808201808211156107fb576107fb614618565b600067ffffffffffffffff82111561468b5761468b614066565b5060051b60200190565b803563ffffffff81168114613a8457600080fd5b6000602082840312156146bb57600080fd5b813567ffffffffffffffff8111156146d257600080fd5b8201601f810184136146e357600080fd5b80356146f66146f182614671565b6140be565b8082825260208201915060206060840285010192508683111561471857600080fd5b6020840193505b828410156137ca576060848803121561473757600080fd5b61473f614095565b8435815261474f60208601614695565b602082015260408501357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461478357600080fd5b604082015282526060939093019260209091019061471f565b600060001982036147af576147af614618565b5060010190565b6000602082840312156147c857600080fd5b813567ffffffffffffffff8111156147df57600080fd5b8201601f810184136147f057600080fd5b80356147fe6146f182614671565b8082825260208201915060208360051b85010192508683111561482057600080fd5b602084015b8381101561491f57803567ffffffffffffffff81111561484457600080fd5b85016060818a03601f1901121561485a57600080fd5b614862614095565b6020820135815261487560408301614695565b6020820152606082013567ffffffffffffffff81111561489457600080fd5b60208184010192505089601f8301126148ac57600080fd5b813567ffffffffffffffff8111156148c6576148c6614066565b6148d96020601f19601f840116016140be565b8181528b60208386010111156148ee57600080fd5b8160208501602083013760006020838301015280604084015250508085525050602083019250602081019050614825565b509695505050505050565b815167ffffffffffffffff81111561494457614944614066565b61495881614952845461419d565b84614308565b6020601f82116001811461498c57600083156149745750848201515b600019600385901b1c1916600184901b17845561434f565b600084815260208120601f198516915b828110156149bc578785015182556020948501946001909201910161499c565b50848210156149da5786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b818103818111156107fb576107fb614618565b60ff82811682821603908111156107fb576107fb614618565b600060208284031215614a2757600080fd5b81516131c78161401f565b60008251614a44818460208701613aa4565b919091019291505056fea164736f6c634300081a000a",
}

var DataFeedsCacheABI = DataFeedsCacheMetaData.ABI

var DataFeedsCacheBin = DataFeedsCacheMetaData.Bin

func DeployDataFeedsCache(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DataFeedsCache, error) {
	parsed, err := DataFeedsCacheMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DataFeedsCacheBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DataFeedsCache{address: address, abi: *parsed, DataFeedsCacheCaller: DataFeedsCacheCaller{contract: contract}, DataFeedsCacheTransactor: DataFeedsCacheTransactor{contract: contract}, DataFeedsCacheFilterer: DataFeedsCacheFilterer{contract: contract}}, nil
}

type DataFeedsCache struct {
	address common.Address
	abi     abi.ABI
	DataFeedsCacheCaller
	DataFeedsCacheTransactor
	DataFeedsCacheFilterer
}

type DataFeedsCacheCaller struct {
	contract *bind.BoundContract
}

type DataFeedsCacheTransactor struct {
	contract *bind.BoundContract
}

type DataFeedsCacheFilterer struct {
	contract *bind.BoundContract
}

type DataFeedsCacheSession struct {
	Contract     *DataFeedsCache
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type DataFeedsCacheCallerSession struct {
	Contract *DataFeedsCacheCaller
	CallOpts bind.CallOpts
}

type DataFeedsCacheTransactorSession struct {
	Contract     *DataFeedsCacheTransactor
	TransactOpts bind.TransactOpts
}

type DataFeedsCacheRaw struct {
	Contract *DataFeedsCache
}

type DataFeedsCacheCallerRaw struct {
	Contract *DataFeedsCacheCaller
}

type DataFeedsCacheTransactorRaw struct {
	Contract *DataFeedsCacheTransactor
}

func NewDataFeedsCache(address common.Address, backend bind.ContractBackend) (*DataFeedsCache, error) {
	abi, err := abi.JSON(strings.NewReader(DataFeedsCacheABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindDataFeedsCache(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCache{address: address, abi: abi, DataFeedsCacheCaller: DataFeedsCacheCaller{contract: contract}, DataFeedsCacheTransactor: DataFeedsCacheTransactor{contract: contract}, DataFeedsCacheFilterer: DataFeedsCacheFilterer{contract: contract}}, nil
}

func NewDataFeedsCacheCaller(address common.Address, caller bind.ContractCaller) (*DataFeedsCacheCaller, error) {
	contract, err := bindDataFeedsCache(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheCaller{contract: contract}, nil
}

func NewDataFeedsCacheTransactor(address common.Address, transactor bind.ContractTransactor) (*DataFeedsCacheTransactor, error) {
	contract, err := bindDataFeedsCache(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheTransactor{contract: contract}, nil
}

func NewDataFeedsCacheFilterer(address common.Address, filterer bind.ContractFilterer) (*DataFeedsCacheFilterer, error) {
	contract, err := bindDataFeedsCache(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheFilterer{contract: contract}, nil
}

func bindDataFeedsCache(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DataFeedsCacheMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_DataFeedsCache *DataFeedsCacheRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataFeedsCache.Contract.DataFeedsCacheCaller.contract.Call(opts, result, method, params...)
}

func (_DataFeedsCache *DataFeedsCacheRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.DataFeedsCacheTransactor.contract.Transfer(opts)
}

func (_DataFeedsCache *DataFeedsCacheRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.DataFeedsCacheTransactor.contract.Transact(opts, method, params...)
}

func (_DataFeedsCache *DataFeedsCacheCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataFeedsCache.Contract.contract.Call(opts, result, method, params...)
}

func (_DataFeedsCache *DataFeedsCacheTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.contract.Transfer(opts)
}

func (_DataFeedsCache *DataFeedsCacheTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.contract.Transact(opts, method, params...)
}

func (_DataFeedsCache *DataFeedsCacheCaller) BundleDecimals(opts *bind.CallOpts) ([]uint8, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "bundleDecimals")

	if err != nil {
		return *new([]uint8), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) BundleDecimals() ([]uint8, error) {
	return _DataFeedsCache.Contract.BundleDecimals(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) BundleDecimals() ([]uint8, error) {
	return _DataFeedsCache.Contract.BundleDecimals(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) CheckFeedPermission(opts *bind.CallOpts, dataId [16]byte, workflowMetadata DataFeedsCacheWorkflowMetadata) (bool, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "checkFeedPermission", dataId, workflowMetadata)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) CheckFeedPermission(dataId [16]byte, workflowMetadata DataFeedsCacheWorkflowMetadata) (bool, error) {
	return _DataFeedsCache.Contract.CheckFeedPermission(&_DataFeedsCache.CallOpts, dataId, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) CheckFeedPermission(dataId [16]byte, workflowMetadata DataFeedsCacheWorkflowMetadata) (bool, error) {
	return _DataFeedsCache.Contract.CheckFeedPermission(&_DataFeedsCache.CallOpts, dataId, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) Decimals() (uint8, error) {
	return _DataFeedsCache.Contract.Decimals(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) Decimals() (uint8, error) {
	return _DataFeedsCache.Contract.Decimals(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) Description() (string, error) {
	return _DataFeedsCache.Contract.Description(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) Description() (string, error) {
	return _DataFeedsCache.Contract.Description(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getAnswer", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetAnswer(&_DataFeedsCache.CallOpts, roundId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetAnswer(roundId *big.Int) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetAnswer(&_DataFeedsCache.CallOpts, roundId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetBundleDecimals(opts *bind.CallOpts, dataId [16]byte) ([]uint8, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getBundleDecimals", dataId)

	if err != nil {
		return *new([]uint8), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetBundleDecimals(dataId [16]byte) ([]uint8, error) {
	return _DataFeedsCache.Contract.GetBundleDecimals(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetBundleDecimals(dataId [16]byte) ([]uint8, error) {
	return _DataFeedsCache.Contract.GetBundleDecimals(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetDataIdForProxy(opts *bind.CallOpts, proxy common.Address) ([16]byte, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getDataIdForProxy", proxy)

	if err != nil {
		return *new([16]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([16]byte)).(*[16]byte)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetDataIdForProxy(proxy common.Address) ([16]byte, error) {
	return _DataFeedsCache.Contract.GetDataIdForProxy(&_DataFeedsCache.CallOpts, proxy)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetDataIdForProxy(proxy common.Address) ([16]byte, error) {
	return _DataFeedsCache.Contract.GetDataIdForProxy(&_DataFeedsCache.CallOpts, proxy)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetDecimals(opts *bind.CallOpts, dataId [16]byte) (uint8, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getDecimals", dataId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetDecimals(dataId [16]byte) (uint8, error) {
	return _DataFeedsCache.Contract.GetDecimals(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetDecimals(dataId [16]byte) (uint8, error) {
	return _DataFeedsCache.Contract.GetDecimals(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetDescription(opts *bind.CallOpts, dataId [16]byte) (string, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getDescription", dataId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetDescription(dataId [16]byte) (string, error) {
	return _DataFeedsCache.Contract.GetDescription(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetDescription(dataId [16]byte) (string, error) {
	return _DataFeedsCache.Contract.GetDescription(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetFeedMetadata(opts *bind.CallOpts, dataId [16]byte, startIndex *big.Int, maxCount *big.Int) ([]DataFeedsCacheWorkflowMetadata, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getFeedMetadata", dataId, startIndex, maxCount)

	if err != nil {
		return *new([]DataFeedsCacheWorkflowMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataFeedsCacheWorkflowMetadata)).(*[]DataFeedsCacheWorkflowMetadata)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetFeedMetadata(dataId [16]byte, startIndex *big.Int, maxCount *big.Int) ([]DataFeedsCacheWorkflowMetadata, error) {
	return _DataFeedsCache.Contract.GetFeedMetadata(&_DataFeedsCache.CallOpts, dataId, startIndex, maxCount)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetFeedMetadata(dataId [16]byte, startIndex *big.Int, maxCount *big.Int) ([]DataFeedsCacheWorkflowMetadata, error) {
	return _DataFeedsCache.Contract.GetFeedMetadata(&_DataFeedsCache.CallOpts, dataId, startIndex, maxCount)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetLatestAnswer(opts *bind.CallOpts, dataId [16]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getLatestAnswer", dataId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetLatestAnswer(dataId [16]byte) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetLatestAnswer(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetLatestAnswer(dataId [16]byte) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetLatestAnswer(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetLatestBundle(opts *bind.CallOpts, dataId [16]byte) ([]byte, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getLatestBundle", dataId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetLatestBundle(dataId [16]byte) ([]byte, error) {
	return _DataFeedsCache.Contract.GetLatestBundle(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetLatestBundle(dataId [16]byte) ([]byte, error) {
	return _DataFeedsCache.Contract.GetLatestBundle(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetLatestBundleTimestamp(opts *bind.CallOpts, dataId [16]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getLatestBundleTimestamp", dataId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetLatestBundleTimestamp(dataId [16]byte) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetLatestBundleTimestamp(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetLatestBundleTimestamp(dataId [16]byte) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetLatestBundleTimestamp(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetLatestRoundData(opts *bind.CallOpts, dataId [16]byte) (GetLatestRoundData,

	error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getLatestRoundData", dataId)

	outstruct := new(GetLatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetLatestRoundData(dataId [16]byte) (GetLatestRoundData,

	error) {
	return _DataFeedsCache.Contract.GetLatestRoundData(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetLatestRoundData(dataId [16]byte) (GetLatestRoundData,

	error) {
	return _DataFeedsCache.Contract.GetLatestRoundData(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetLatestTimestamp(opts *bind.CallOpts, dataId [16]byte) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getLatestTimestamp", dataId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetLatestTimestamp(dataId [16]byte) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetLatestTimestamp(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetLatestTimestamp(dataId [16]byte) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetLatestTimestamp(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (GetRoundData,

	error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getRoundData", roundId)

	outstruct := new(GetRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetRoundData(roundId *big.Int) (GetRoundData,

	error) {
	return _DataFeedsCache.Contract.GetRoundData(&_DataFeedsCache.CallOpts, roundId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetRoundData(roundId *big.Int) (GetRoundData,

	error) {
	return _DataFeedsCache.Contract.GetRoundData(&_DataFeedsCache.CallOpts, roundId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getTimestamp", roundId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetTimestamp(&_DataFeedsCache.CallOpts, roundId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetTimestamp(roundId *big.Int) (*big.Int, error) {
	return _DataFeedsCache.Contract.GetTimestamp(&_DataFeedsCache.CallOpts, roundId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) IsFeedAdmin(opts *bind.CallOpts, feedAdmin common.Address) (bool, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "isFeedAdmin", feedAdmin)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) IsFeedAdmin(feedAdmin common.Address) (bool, error) {
	return _DataFeedsCache.Contract.IsFeedAdmin(&_DataFeedsCache.CallOpts, feedAdmin)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) IsFeedAdmin(feedAdmin common.Address) (bool, error) {
	return _DataFeedsCache.Contract.IsFeedAdmin(&_DataFeedsCache.CallOpts, feedAdmin)
}

func (_DataFeedsCache *DataFeedsCacheCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) LatestAnswer() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestAnswer(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) LatestAnswer() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestAnswer(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) LatestBundle(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "latestBundle")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) LatestBundle() ([]byte, error) {
	return _DataFeedsCache.Contract.LatestBundle(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) LatestBundle() ([]byte, error) {
	return _DataFeedsCache.Contract.LatestBundle(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) LatestBundleTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "latestBundleTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) LatestBundleTimestamp() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestBundleTimestamp(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) LatestBundleTimestamp() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestBundleTimestamp(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) LatestRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "latestRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) LatestRound() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestRound(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) LatestRound() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestRound(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

	error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "latestRoundData")

	outstruct := new(LatestRoundData)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Answer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.AnsweredInRound = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_DataFeedsCache *DataFeedsCacheSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DataFeedsCache.Contract.LatestRoundData(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) LatestRoundData() (LatestRoundData,

	error) {
	return _DataFeedsCache.Contract.LatestRoundData(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) LatestTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "latestTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) LatestTimestamp() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestTimestamp(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) LatestTimestamp() (*big.Int, error) {
	return _DataFeedsCache.Contract.LatestTimestamp(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) Owner() (common.Address, error) {
	return _DataFeedsCache.Contract.Owner(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) Owner() (common.Address, error) {
	return _DataFeedsCache.Contract.Owner(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DataFeedsCache.Contract.SupportsInterface(&_DataFeedsCache.CallOpts, interfaceId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DataFeedsCache.Contract.SupportsInterface(&_DataFeedsCache.CallOpts, interfaceId)
}

func (_DataFeedsCache *DataFeedsCacheCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) TypeAndVersion() (string, error) {
	return _DataFeedsCache.Contract.TypeAndVersion(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) TypeAndVersion() (string, error) {
	return _DataFeedsCache.Contract.TypeAndVersion(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_DataFeedsCache *DataFeedsCacheSession) Version() (*big.Int, error) {
	return _DataFeedsCache.Contract.Version(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) Version() (*big.Int, error) {
	return _DataFeedsCache.Contract.Version(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "acceptOwnership")
}

func (_DataFeedsCache *DataFeedsCacheSession) AcceptOwnership() (*types.Transaction, error) {
	return _DataFeedsCache.Contract.AcceptOwnership(&_DataFeedsCache.TransactOpts)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DataFeedsCache.Contract.AcceptOwnership(&_DataFeedsCache.TransactOpts)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) OnReport(opts *bind.TransactOpts, metadata []byte, report []byte) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "onReport", metadata, report)
}

func (_DataFeedsCache *DataFeedsCacheSession) OnReport(metadata []byte, report []byte) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.OnReport(&_DataFeedsCache.TransactOpts, metadata, report)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) OnReport(metadata []byte, report []byte) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.OnReport(&_DataFeedsCache.TransactOpts, metadata, report)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) RecoverTokens(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "recoverTokens", token, to, amount)
}

func (_DataFeedsCache *DataFeedsCacheSession) RecoverTokens(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.RecoverTokens(&_DataFeedsCache.TransactOpts, token, to, amount)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) RecoverTokens(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.RecoverTokens(&_DataFeedsCache.TransactOpts, token, to, amount)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) RemoveDataIdMappingsForProxies(opts *bind.TransactOpts, proxies []common.Address) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "removeDataIdMappingsForProxies", proxies)
}

func (_DataFeedsCache *DataFeedsCacheSession) RemoveDataIdMappingsForProxies(proxies []common.Address) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.RemoveDataIdMappingsForProxies(&_DataFeedsCache.TransactOpts, proxies)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) RemoveDataIdMappingsForProxies(proxies []common.Address) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.RemoveDataIdMappingsForProxies(&_DataFeedsCache.TransactOpts, proxies)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) RemoveFeedConfigs(opts *bind.TransactOpts, dataIds [][16]byte) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "removeFeedConfigs", dataIds)
}

func (_DataFeedsCache *DataFeedsCacheSession) RemoveFeedConfigs(dataIds [][16]byte) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.RemoveFeedConfigs(&_DataFeedsCache.TransactOpts, dataIds)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) RemoveFeedConfigs(dataIds [][16]byte) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.RemoveFeedConfigs(&_DataFeedsCache.TransactOpts, dataIds)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) SetBundleFeedConfigs(opts *bind.TransactOpts, dataIds [][16]byte, descriptions []string, decimalsMatrix [][]uint8, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "setBundleFeedConfigs", dataIds, descriptions, decimalsMatrix, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheSession) SetBundleFeedConfigs(dataIds [][16]byte, descriptions []string, decimalsMatrix [][]uint8, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.SetBundleFeedConfigs(&_DataFeedsCache.TransactOpts, dataIds, descriptions, decimalsMatrix, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) SetBundleFeedConfigs(dataIds [][16]byte, descriptions []string, decimalsMatrix [][]uint8, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.SetBundleFeedConfigs(&_DataFeedsCache.TransactOpts, dataIds, descriptions, decimalsMatrix, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) SetDecimalFeedConfigs(opts *bind.TransactOpts, dataIds [][16]byte, descriptions []string, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "setDecimalFeedConfigs", dataIds, descriptions, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheSession) SetDecimalFeedConfigs(dataIds [][16]byte, descriptions []string, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.SetDecimalFeedConfigs(&_DataFeedsCache.TransactOpts, dataIds, descriptions, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) SetDecimalFeedConfigs(dataIds [][16]byte, descriptions []string, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.SetDecimalFeedConfigs(&_DataFeedsCache.TransactOpts, dataIds, descriptions, workflowMetadata)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) SetFeedAdmin(opts *bind.TransactOpts, feedAdmin common.Address, isAdmin bool) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "setFeedAdmin", feedAdmin, isAdmin)
}

func (_DataFeedsCache *DataFeedsCacheSession) SetFeedAdmin(feedAdmin common.Address, isAdmin bool) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.SetFeedAdmin(&_DataFeedsCache.TransactOpts, feedAdmin, isAdmin)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) SetFeedAdmin(feedAdmin common.Address, isAdmin bool) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.SetFeedAdmin(&_DataFeedsCache.TransactOpts, feedAdmin, isAdmin)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "transferOwnership", to)
}

func (_DataFeedsCache *DataFeedsCacheSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.TransferOwnership(&_DataFeedsCache.TransactOpts, to)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.TransferOwnership(&_DataFeedsCache.TransactOpts, to)
}

func (_DataFeedsCache *DataFeedsCacheTransactor) UpdateDataIdMappingsForProxies(opts *bind.TransactOpts, proxies []common.Address, dataIds [][16]byte) (*types.Transaction, error) {
	return _DataFeedsCache.contract.Transact(opts, "updateDataIdMappingsForProxies", proxies, dataIds)
}

func (_DataFeedsCache *DataFeedsCacheSession) UpdateDataIdMappingsForProxies(proxies []common.Address, dataIds [][16]byte) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.UpdateDataIdMappingsForProxies(&_DataFeedsCache.TransactOpts, proxies, dataIds)
}

func (_DataFeedsCache *DataFeedsCacheTransactorSession) UpdateDataIdMappingsForProxies(proxies []common.Address, dataIds [][16]byte) (*types.Transaction, error) {
	return _DataFeedsCache.Contract.UpdateDataIdMappingsForProxies(&_DataFeedsCache.TransactOpts, proxies, dataIds)
}

type DataFeedsCacheAnswerUpdatedIterator struct {
	Event *DataFeedsCacheAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheAnswerUpdated)
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
		it.Event = new(DataFeedsCacheAnswerUpdated)
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

func (it *DataFeedsCacheAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*DataFeedsCacheAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheAnswerUpdatedIterator{contract: _DataFeedsCache.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheAnswerUpdated)
				if err := _DataFeedsCache.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseAnswerUpdated(log types.Log) (*DataFeedsCacheAnswerUpdated, error) {
	event := new(DataFeedsCacheAnswerUpdated)
	if err := _DataFeedsCache.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheBundleFeedConfigSetIterator struct {
	Event *DataFeedsCacheBundleFeedConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheBundleFeedConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheBundleFeedConfigSet)
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
		it.Event = new(DataFeedsCacheBundleFeedConfigSet)
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

func (it *DataFeedsCacheBundleFeedConfigSetIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheBundleFeedConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheBundleFeedConfigSet struct {
	DataId           [16]byte
	Decimals         []uint8
	Description      string
	WorkflowMetadata []DataFeedsCacheWorkflowMetadata
	Raw              types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterBundleFeedConfigSet(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheBundleFeedConfigSetIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "BundleFeedConfigSet", dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheBundleFeedConfigSetIterator{contract: _DataFeedsCache.contract, event: "BundleFeedConfigSet", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchBundleFeedConfigSet(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheBundleFeedConfigSet, dataId [][16]byte) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "BundleFeedConfigSet", dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheBundleFeedConfigSet)
				if err := _DataFeedsCache.contract.UnpackLog(event, "BundleFeedConfigSet", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseBundleFeedConfigSet(log types.Log) (*DataFeedsCacheBundleFeedConfigSet, error) {
	event := new(DataFeedsCacheBundleFeedConfigSet)
	if err := _DataFeedsCache.contract.UnpackLog(event, "BundleFeedConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheBundleReportUpdatedIterator struct {
	Event *DataFeedsCacheBundleReportUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheBundleReportUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheBundleReportUpdated)
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
		it.Event = new(DataFeedsCacheBundleReportUpdated)
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

func (it *DataFeedsCacheBundleReportUpdatedIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheBundleReportUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheBundleReportUpdated struct {
	DataId    [16]byte
	Timestamp *big.Int
	Bundle    []byte
	Raw       types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterBundleReportUpdated(opts *bind.FilterOpts, dataId [][16]byte, timestamp []*big.Int) (*DataFeedsCacheBundleReportUpdatedIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "BundleReportUpdated", dataIdRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheBundleReportUpdatedIterator{contract: _DataFeedsCache.contract, event: "BundleReportUpdated", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchBundleReportUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheBundleReportUpdated, dataId [][16]byte, timestamp []*big.Int) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "BundleReportUpdated", dataIdRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheBundleReportUpdated)
				if err := _DataFeedsCache.contract.UnpackLog(event, "BundleReportUpdated", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseBundleReportUpdated(log types.Log) (*DataFeedsCacheBundleReportUpdated, error) {
	event := new(DataFeedsCacheBundleReportUpdated)
	if err := _DataFeedsCache.contract.UnpackLog(event, "BundleReportUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheDecimalFeedConfigSetIterator struct {
	Event *DataFeedsCacheDecimalFeedConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheDecimalFeedConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheDecimalFeedConfigSet)
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
		it.Event = new(DataFeedsCacheDecimalFeedConfigSet)
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

func (it *DataFeedsCacheDecimalFeedConfigSetIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheDecimalFeedConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheDecimalFeedConfigSet struct {
	DataId           [16]byte
	Decimals         uint8
	Description      string
	WorkflowMetadata []DataFeedsCacheWorkflowMetadata
	Raw              types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterDecimalFeedConfigSet(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheDecimalFeedConfigSetIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "DecimalFeedConfigSet", dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheDecimalFeedConfigSetIterator{contract: _DataFeedsCache.contract, event: "DecimalFeedConfigSet", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchDecimalFeedConfigSet(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheDecimalFeedConfigSet, dataId [][16]byte) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "DecimalFeedConfigSet", dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheDecimalFeedConfigSet)
				if err := _DataFeedsCache.contract.UnpackLog(event, "DecimalFeedConfigSet", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseDecimalFeedConfigSet(log types.Log) (*DataFeedsCacheDecimalFeedConfigSet, error) {
	event := new(DataFeedsCacheDecimalFeedConfigSet)
	if err := _DataFeedsCache.contract.UnpackLog(event, "DecimalFeedConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheDecimalReportUpdatedIterator struct {
	Event *DataFeedsCacheDecimalReportUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheDecimalReportUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheDecimalReportUpdated)
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
		it.Event = new(DataFeedsCacheDecimalReportUpdated)
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

func (it *DataFeedsCacheDecimalReportUpdatedIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheDecimalReportUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheDecimalReportUpdated struct {
	DataId    [16]byte
	RoundId   *big.Int
	Timestamp *big.Int
	Answer    *big.Int
	Raw       types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterDecimalReportUpdated(opts *bind.FilterOpts, dataId [][16]byte, roundId []*big.Int, timestamp []*big.Int) (*DataFeedsCacheDecimalReportUpdatedIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "DecimalReportUpdated", dataIdRule, roundIdRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheDecimalReportUpdatedIterator{contract: _DataFeedsCache.contract, event: "DecimalReportUpdated", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchDecimalReportUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheDecimalReportUpdated, dataId [][16]byte, roundId []*big.Int, timestamp []*big.Int) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "DecimalReportUpdated", dataIdRule, roundIdRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheDecimalReportUpdated)
				if err := _DataFeedsCache.contract.UnpackLog(event, "DecimalReportUpdated", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseDecimalReportUpdated(log types.Log) (*DataFeedsCacheDecimalReportUpdated, error) {
	event := new(DataFeedsCacheDecimalReportUpdated)
	if err := _DataFeedsCache.contract.UnpackLog(event, "DecimalReportUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheFeedAdminSetIterator struct {
	Event *DataFeedsCacheFeedAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheFeedAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheFeedAdminSet)
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
		it.Event = new(DataFeedsCacheFeedAdminSet)
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

func (it *DataFeedsCacheFeedAdminSetIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheFeedAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheFeedAdminSet struct {
	FeedAdmin common.Address
	IsAdmin   bool
	Raw       types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterFeedAdminSet(opts *bind.FilterOpts, feedAdmin []common.Address, isAdmin []bool) (*DataFeedsCacheFeedAdminSetIterator, error) {

	var feedAdminRule []interface{}
	for _, feedAdminItem := range feedAdmin {
		feedAdminRule = append(feedAdminRule, feedAdminItem)
	}
	var isAdminRule []interface{}
	for _, isAdminItem := range isAdmin {
		isAdminRule = append(isAdminRule, isAdminItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "FeedAdminSet", feedAdminRule, isAdminRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheFeedAdminSetIterator{contract: _DataFeedsCache.contract, event: "FeedAdminSet", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchFeedAdminSet(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheFeedAdminSet, feedAdmin []common.Address, isAdmin []bool) (event.Subscription, error) {

	var feedAdminRule []interface{}
	for _, feedAdminItem := range feedAdmin {
		feedAdminRule = append(feedAdminRule, feedAdminItem)
	}
	var isAdminRule []interface{}
	for _, isAdminItem := range isAdmin {
		isAdminRule = append(isAdminRule, isAdminItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "FeedAdminSet", feedAdminRule, isAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheFeedAdminSet)
				if err := _DataFeedsCache.contract.UnpackLog(event, "FeedAdminSet", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseFeedAdminSet(log types.Log) (*DataFeedsCacheFeedAdminSet, error) {
	event := new(DataFeedsCacheFeedAdminSet)
	if err := _DataFeedsCache.contract.UnpackLog(event, "FeedAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheFeedConfigRemovedIterator struct {
	Event *DataFeedsCacheFeedConfigRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheFeedConfigRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheFeedConfigRemoved)
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
		it.Event = new(DataFeedsCacheFeedConfigRemoved)
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

func (it *DataFeedsCacheFeedConfigRemovedIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheFeedConfigRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheFeedConfigRemoved struct {
	DataId [16]byte
	Raw    types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterFeedConfigRemoved(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheFeedConfigRemovedIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "FeedConfigRemoved", dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheFeedConfigRemovedIterator{contract: _DataFeedsCache.contract, event: "FeedConfigRemoved", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchFeedConfigRemoved(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheFeedConfigRemoved, dataId [][16]byte) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "FeedConfigRemoved", dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheFeedConfigRemoved)
				if err := _DataFeedsCache.contract.UnpackLog(event, "FeedConfigRemoved", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseFeedConfigRemoved(log types.Log) (*DataFeedsCacheFeedConfigRemoved, error) {
	event := new(DataFeedsCacheFeedConfigRemoved)
	if err := _DataFeedsCache.contract.UnpackLog(event, "FeedConfigRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheInvalidUpdatePermissionIterator struct {
	Event *DataFeedsCacheInvalidUpdatePermission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheInvalidUpdatePermissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheInvalidUpdatePermission)
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
		it.Event = new(DataFeedsCacheInvalidUpdatePermission)
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

func (it *DataFeedsCacheInvalidUpdatePermissionIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheInvalidUpdatePermissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheInvalidUpdatePermission struct {
	DataId        [16]byte
	Sender        common.Address
	WorkflowOwner common.Address
	WorkflowName  [10]byte
	Raw           types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterInvalidUpdatePermission(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheInvalidUpdatePermissionIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "InvalidUpdatePermission", dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheInvalidUpdatePermissionIterator{contract: _DataFeedsCache.contract, event: "InvalidUpdatePermission", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchInvalidUpdatePermission(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheInvalidUpdatePermission, dataId [][16]byte) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "InvalidUpdatePermission", dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheInvalidUpdatePermission)
				if err := _DataFeedsCache.contract.UnpackLog(event, "InvalidUpdatePermission", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseInvalidUpdatePermission(log types.Log) (*DataFeedsCacheInvalidUpdatePermission, error) {
	event := new(DataFeedsCacheInvalidUpdatePermission)
	if err := _DataFeedsCache.contract.UnpackLog(event, "InvalidUpdatePermission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheNewRoundIterator struct {
	Event *DataFeedsCacheNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheNewRound)
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
		it.Event = new(DataFeedsCacheNewRound)
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

func (it *DataFeedsCacheNewRoundIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*DataFeedsCacheNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheNewRoundIterator{contract: _DataFeedsCache.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheNewRound)
				if err := _DataFeedsCache.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseNewRound(log types.Log) (*DataFeedsCacheNewRound, error) {
	event := new(DataFeedsCacheNewRound)
	if err := _DataFeedsCache.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheOwnershipTransferRequestedIterator struct {
	Event *DataFeedsCacheOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheOwnershipTransferRequested)
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
		it.Event = new(DataFeedsCacheOwnershipTransferRequested)
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

func (it *DataFeedsCacheOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DataFeedsCacheOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheOwnershipTransferRequestedIterator{contract: _DataFeedsCache.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheOwnershipTransferRequested)
				if err := _DataFeedsCache.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseOwnershipTransferRequested(log types.Log) (*DataFeedsCacheOwnershipTransferRequested, error) {
	event := new(DataFeedsCacheOwnershipTransferRequested)
	if err := _DataFeedsCache.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheOwnershipTransferredIterator struct {
	Event *DataFeedsCacheOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheOwnershipTransferred)
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
		it.Event = new(DataFeedsCacheOwnershipTransferred)
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

func (it *DataFeedsCacheOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DataFeedsCacheOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheOwnershipTransferredIterator{contract: _DataFeedsCache.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheOwnershipTransferred)
				if err := _DataFeedsCache.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseOwnershipTransferred(log types.Log) (*DataFeedsCacheOwnershipTransferred, error) {
	event := new(DataFeedsCacheOwnershipTransferred)
	if err := _DataFeedsCache.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheProxyDataIdRemovedIterator struct {
	Event *DataFeedsCacheProxyDataIdRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheProxyDataIdRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheProxyDataIdRemoved)
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
		it.Event = new(DataFeedsCacheProxyDataIdRemoved)
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

func (it *DataFeedsCacheProxyDataIdRemovedIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheProxyDataIdRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheProxyDataIdRemoved struct {
	Proxy  common.Address
	DataId [16]byte
	Raw    types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterProxyDataIdRemoved(opts *bind.FilterOpts, proxy []common.Address, dataId [][16]byte) (*DataFeedsCacheProxyDataIdRemovedIterator, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "ProxyDataIdRemoved", proxyRule, dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheProxyDataIdRemovedIterator{contract: _DataFeedsCache.contract, event: "ProxyDataIdRemoved", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchProxyDataIdRemoved(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheProxyDataIdRemoved, proxy []common.Address, dataId [][16]byte) (event.Subscription, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "ProxyDataIdRemoved", proxyRule, dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheProxyDataIdRemoved)
				if err := _DataFeedsCache.contract.UnpackLog(event, "ProxyDataIdRemoved", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseProxyDataIdRemoved(log types.Log) (*DataFeedsCacheProxyDataIdRemoved, error) {
	event := new(DataFeedsCacheProxyDataIdRemoved)
	if err := _DataFeedsCache.contract.UnpackLog(event, "ProxyDataIdRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheProxyDataIdUpdatedIterator struct {
	Event *DataFeedsCacheProxyDataIdUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheProxyDataIdUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheProxyDataIdUpdated)
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
		it.Event = new(DataFeedsCacheProxyDataIdUpdated)
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

func (it *DataFeedsCacheProxyDataIdUpdatedIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheProxyDataIdUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheProxyDataIdUpdated struct {
	Proxy  common.Address
	DataId [16]byte
	Raw    types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterProxyDataIdUpdated(opts *bind.FilterOpts, proxy []common.Address, dataId [][16]byte) (*DataFeedsCacheProxyDataIdUpdatedIterator, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "ProxyDataIdUpdated", proxyRule, dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheProxyDataIdUpdatedIterator{contract: _DataFeedsCache.contract, event: "ProxyDataIdUpdated", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchProxyDataIdUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheProxyDataIdUpdated, proxy []common.Address, dataId [][16]byte) (event.Subscription, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "ProxyDataIdUpdated", proxyRule, dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheProxyDataIdUpdated)
				if err := _DataFeedsCache.contract.UnpackLog(event, "ProxyDataIdUpdated", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseProxyDataIdUpdated(log types.Log) (*DataFeedsCacheProxyDataIdUpdated, error) {
	event := new(DataFeedsCacheProxyDataIdUpdated)
	if err := _DataFeedsCache.contract.UnpackLog(event, "ProxyDataIdUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheStaleBundleReportIterator struct {
	Event *DataFeedsCacheStaleBundleReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheStaleBundleReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheStaleBundleReport)
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
		it.Event = new(DataFeedsCacheStaleBundleReport)
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

func (it *DataFeedsCacheStaleBundleReportIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheStaleBundleReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheStaleBundleReport struct {
	DataId          [16]byte
	ReportTimestamp *big.Int
	LatestTimestamp *big.Int
	Raw             types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterStaleBundleReport(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheStaleBundleReportIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "StaleBundleReport", dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheStaleBundleReportIterator{contract: _DataFeedsCache.contract, event: "StaleBundleReport", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchStaleBundleReport(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheStaleBundleReport, dataId [][16]byte) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "StaleBundleReport", dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheStaleBundleReport)
				if err := _DataFeedsCache.contract.UnpackLog(event, "StaleBundleReport", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseStaleBundleReport(log types.Log) (*DataFeedsCacheStaleBundleReport, error) {
	event := new(DataFeedsCacheStaleBundleReport)
	if err := _DataFeedsCache.contract.UnpackLog(event, "StaleBundleReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheStaleDecimalReportIterator struct {
	Event *DataFeedsCacheStaleDecimalReport

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheStaleDecimalReportIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheStaleDecimalReport)
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
		it.Event = new(DataFeedsCacheStaleDecimalReport)
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

func (it *DataFeedsCacheStaleDecimalReportIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheStaleDecimalReportIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheStaleDecimalReport struct {
	DataId          [16]byte
	ReportTimestamp *big.Int
	LatestTimestamp *big.Int
	Raw             types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterStaleDecimalReport(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheStaleDecimalReportIterator, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "StaleDecimalReport", dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheStaleDecimalReportIterator{contract: _DataFeedsCache.contract, event: "StaleDecimalReport", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchStaleDecimalReport(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheStaleDecimalReport, dataId [][16]byte) (event.Subscription, error) {

	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "StaleDecimalReport", dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheStaleDecimalReport)
				if err := _DataFeedsCache.contract.UnpackLog(event, "StaleDecimalReport", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseStaleDecimalReport(log types.Log) (*DataFeedsCacheStaleDecimalReport, error) {
	event := new(DataFeedsCacheStaleDecimalReport)
	if err := _DataFeedsCache.contract.UnpackLog(event, "StaleDecimalReport", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type DataFeedsCacheTokenRecoveredIterator struct {
	Event *DataFeedsCacheTokenRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *DataFeedsCacheTokenRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataFeedsCacheTokenRecovered)
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
		it.Event = new(DataFeedsCacheTokenRecovered)
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

func (it *DataFeedsCacheTokenRecoveredIterator) Error() error {
	return it.fail
}

func (it *DataFeedsCacheTokenRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type DataFeedsCacheTokenRecovered struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_DataFeedsCache *DataFeedsCacheFilterer) FilterTokenRecovered(opts *bind.FilterOpts, token []common.Address, to []common.Address) (*DataFeedsCacheTokenRecoveredIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DataFeedsCache.contract.FilterLogs(opts, "TokenRecovered", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DataFeedsCacheTokenRecoveredIterator{contract: _DataFeedsCache.contract, event: "TokenRecovered", logs: logs, sub: sub}, nil
}

func (_DataFeedsCache *DataFeedsCacheFilterer) WatchTokenRecovered(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheTokenRecovered, token []common.Address, to []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _DataFeedsCache.contract.WatchLogs(opts, "TokenRecovered", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(DataFeedsCacheTokenRecovered)
				if err := _DataFeedsCache.contract.UnpackLog(event, "TokenRecovered", log); err != nil {
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

func (_DataFeedsCache *DataFeedsCacheFilterer) ParseTokenRecovered(log types.Log) (*DataFeedsCacheTokenRecovered, error) {
	event := new(DataFeedsCacheTokenRecovered)
	if err := _DataFeedsCache.contract.UnpackLog(event, "TokenRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetLatestRoundData struct {
	Id              *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type GetRoundData struct {
	Id              *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
type LatestRoundData struct {
	Id              *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func (_DataFeedsCache *DataFeedsCache) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _DataFeedsCache.abi.Events["AnswerUpdated"].ID:
		return _DataFeedsCache.ParseAnswerUpdated(log)
	case _DataFeedsCache.abi.Events["BundleFeedConfigSet"].ID:
		return _DataFeedsCache.ParseBundleFeedConfigSet(log)
	case _DataFeedsCache.abi.Events["BundleReportUpdated"].ID:
		return _DataFeedsCache.ParseBundleReportUpdated(log)
	case _DataFeedsCache.abi.Events["DecimalFeedConfigSet"].ID:
		return _DataFeedsCache.ParseDecimalFeedConfigSet(log)
	case _DataFeedsCache.abi.Events["DecimalReportUpdated"].ID:
		return _DataFeedsCache.ParseDecimalReportUpdated(log)
	case _DataFeedsCache.abi.Events["FeedAdminSet"].ID:
		return _DataFeedsCache.ParseFeedAdminSet(log)
	case _DataFeedsCache.abi.Events["FeedConfigRemoved"].ID:
		return _DataFeedsCache.ParseFeedConfigRemoved(log)
	case _DataFeedsCache.abi.Events["InvalidUpdatePermission"].ID:
		return _DataFeedsCache.ParseInvalidUpdatePermission(log)
	case _DataFeedsCache.abi.Events["NewRound"].ID:
		return _DataFeedsCache.ParseNewRound(log)
	case _DataFeedsCache.abi.Events["OwnershipTransferRequested"].ID:
		return _DataFeedsCache.ParseOwnershipTransferRequested(log)
	case _DataFeedsCache.abi.Events["OwnershipTransferred"].ID:
		return _DataFeedsCache.ParseOwnershipTransferred(log)
	case _DataFeedsCache.abi.Events["ProxyDataIdRemoved"].ID:
		return _DataFeedsCache.ParseProxyDataIdRemoved(log)
	case _DataFeedsCache.abi.Events["ProxyDataIdUpdated"].ID:
		return _DataFeedsCache.ParseProxyDataIdUpdated(log)
	case _DataFeedsCache.abi.Events["StaleBundleReport"].ID:
		return _DataFeedsCache.ParseStaleBundleReport(log)
	case _DataFeedsCache.abi.Events["StaleDecimalReport"].ID:
		return _DataFeedsCache.ParseStaleDecimalReport(log)
	case _DataFeedsCache.abi.Events["TokenRecovered"].ID:
		return _DataFeedsCache.ParseTokenRecovered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (DataFeedsCacheAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (DataFeedsCacheBundleFeedConfigSet) Topic() common.Hash {
	return common.HexToHash("0xdfebe0878c5611549f54908260ca12271c7ff3f0ebae0c1de47732612403869e")
}

func (DataFeedsCacheBundleReportUpdated) Topic() common.Hash {
	return common.HexToHash("0x1dc1bef0b59d624eab3f0ec044781bb5b8594cd64f0ba09d789f5b51acab1614")
}

func (DataFeedsCacheDecimalFeedConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2dec0e9ffbb18c6499fc8bee8b9c35f765e76d9dbd436f25dd00a80de267ac0d")
}

func (DataFeedsCacheDecimalReportUpdated) Topic() common.Hash {
	return common.HexToHash("0x82584589cd7284d4503ed582275e22b2e8f459f9cf4170a7235844e367f966d5")
}

func (DataFeedsCacheFeedAdminSet) Topic() common.Hash {
	return common.HexToHash("0x93a3fa5993d2a54de369386625330cc6d73caee7fece4b3983cf299b264473fd")
}

func (DataFeedsCacheFeedConfigRemoved) Topic() common.Hash {
	return common.HexToHash("0x871bcdef10dee59b87f17bab788b72faa8dfe1a9cc5bdc45c3baf4c18fa33910")
}

func (DataFeedsCacheInvalidUpdatePermission) Topic() common.Hash {
	return common.HexToHash("0xeeeaa8bf618ff6d960c6cf5935e68384f066abcc8b95d0de91bd773c16ae3ae3")
}

func (DataFeedsCacheNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (DataFeedsCacheOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (DataFeedsCacheOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (DataFeedsCacheProxyDataIdRemoved) Topic() common.Hash {
	return common.HexToHash("0x4200186b7bc2d4f13f7888c5bbe9461d57da88705be86521f3d78be691ad1d2a")
}

func (DataFeedsCacheProxyDataIdUpdated) Topic() common.Hash {
	return common.HexToHash("0xf31b9e58190970ef07c23d0ba78c358eb3b416e829ef484b29b9993a6b1b285a")
}

func (DataFeedsCacheStaleBundleReport) Topic() common.Hash {
	return common.HexToHash("0x51001b67094834cc084a0c1feb791cf84a481357aa66b924ba205d4cb56fd981")
}

func (DataFeedsCacheStaleDecimalReport) Topic() common.Hash {
	return common.HexToHash("0xcf16f5f704f981fa2279afa1877dd1fdaa462a03a71ec51b9d3b2416a59a013e")
}

func (DataFeedsCacheTokenRecovered) Topic() common.Hash {
	return common.HexToHash("0x879f92dded0f26b83c3e00b12e0395dc72cfc3077343d1854ed6988edd1f9096")
}

func (_DataFeedsCache *DataFeedsCache) Address() common.Address {
	return _DataFeedsCache.address
}

type DataFeedsCacheInterface interface {
	BundleDecimals(opts *bind.CallOpts) ([]uint8, error)

	CheckFeedPermission(opts *bind.CallOpts, dataId [16]byte, workflowMetadata DataFeedsCacheWorkflowMetadata) (bool, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	GetAnswer(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	GetBundleDecimals(opts *bind.CallOpts, dataId [16]byte) ([]uint8, error)

	GetDataIdForProxy(opts *bind.CallOpts, proxy common.Address) ([16]byte, error)

	GetDecimals(opts *bind.CallOpts, dataId [16]byte) (uint8, error)

	GetDescription(opts *bind.CallOpts, dataId [16]byte) (string, error)

	GetFeedMetadata(opts *bind.CallOpts, dataId [16]byte, startIndex *big.Int, maxCount *big.Int) ([]DataFeedsCacheWorkflowMetadata, error)

	GetLatestAnswer(opts *bind.CallOpts, dataId [16]byte) (*big.Int, error)

	GetLatestBundle(opts *bind.CallOpts, dataId [16]byte) ([]byte, error)

	GetLatestBundleTimestamp(opts *bind.CallOpts, dataId [16]byte) (*big.Int, error)

	GetLatestRoundData(opts *bind.CallOpts, dataId [16]byte) (GetLatestRoundData,

		error)

	GetLatestTimestamp(opts *bind.CallOpts, dataId [16]byte) (*big.Int, error)

	GetRoundData(opts *bind.CallOpts, roundId *big.Int) (GetRoundData,

		error)

	GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	IsFeedAdmin(opts *bind.CallOpts, feedAdmin common.Address) (bool, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestBundle(opts *bind.CallOpts) ([]byte, error)

	LatestBundleTimestamp(opts *bind.CallOpts) (*big.Int, error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (LatestRoundData,

		error)

	LatestTimestamp(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	OnReport(opts *bind.TransactOpts, metadata []byte, report []byte) (*types.Transaction, error)

	RecoverTokens(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)

	RemoveDataIdMappingsForProxies(opts *bind.TransactOpts, proxies []common.Address) (*types.Transaction, error)

	RemoveFeedConfigs(opts *bind.TransactOpts, dataIds [][16]byte) (*types.Transaction, error)

	SetBundleFeedConfigs(opts *bind.TransactOpts, dataIds [][16]byte, descriptions []string, decimalsMatrix [][]uint8, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error)

	SetDecimalFeedConfigs(opts *bind.TransactOpts, dataIds [][16]byte, descriptions []string, workflowMetadata []DataFeedsCacheWorkflowMetadata) (*types.Transaction, error)

	SetFeedAdmin(opts *bind.TransactOpts, feedAdmin common.Address, isAdmin bool) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateDataIdMappingsForProxies(opts *bind.TransactOpts, proxies []common.Address, dataIds [][16]byte) (*types.Transaction, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*DataFeedsCacheAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*DataFeedsCacheAnswerUpdated, error)

	FilterBundleFeedConfigSet(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheBundleFeedConfigSetIterator, error)

	WatchBundleFeedConfigSet(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheBundleFeedConfigSet, dataId [][16]byte) (event.Subscription, error)

	ParseBundleFeedConfigSet(log types.Log) (*DataFeedsCacheBundleFeedConfigSet, error)

	FilterBundleReportUpdated(opts *bind.FilterOpts, dataId [][16]byte, timestamp []*big.Int) (*DataFeedsCacheBundleReportUpdatedIterator, error)

	WatchBundleReportUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheBundleReportUpdated, dataId [][16]byte, timestamp []*big.Int) (event.Subscription, error)

	ParseBundleReportUpdated(log types.Log) (*DataFeedsCacheBundleReportUpdated, error)

	FilterDecimalFeedConfigSet(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheDecimalFeedConfigSetIterator, error)

	WatchDecimalFeedConfigSet(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheDecimalFeedConfigSet, dataId [][16]byte) (event.Subscription, error)

	ParseDecimalFeedConfigSet(log types.Log) (*DataFeedsCacheDecimalFeedConfigSet, error)

	FilterDecimalReportUpdated(opts *bind.FilterOpts, dataId [][16]byte, roundId []*big.Int, timestamp []*big.Int) (*DataFeedsCacheDecimalReportUpdatedIterator, error)

	WatchDecimalReportUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheDecimalReportUpdated, dataId [][16]byte, roundId []*big.Int, timestamp []*big.Int) (event.Subscription, error)

	ParseDecimalReportUpdated(log types.Log) (*DataFeedsCacheDecimalReportUpdated, error)

	FilterFeedAdminSet(opts *bind.FilterOpts, feedAdmin []common.Address, isAdmin []bool) (*DataFeedsCacheFeedAdminSetIterator, error)

	WatchFeedAdminSet(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheFeedAdminSet, feedAdmin []common.Address, isAdmin []bool) (event.Subscription, error)

	ParseFeedAdminSet(log types.Log) (*DataFeedsCacheFeedAdminSet, error)

	FilterFeedConfigRemoved(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheFeedConfigRemovedIterator, error)

	WatchFeedConfigRemoved(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheFeedConfigRemoved, dataId [][16]byte) (event.Subscription, error)

	ParseFeedConfigRemoved(log types.Log) (*DataFeedsCacheFeedConfigRemoved, error)

	FilterInvalidUpdatePermission(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheInvalidUpdatePermissionIterator, error)

	WatchInvalidUpdatePermission(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheInvalidUpdatePermission, dataId [][16]byte) (event.Subscription, error)

	ParseInvalidUpdatePermission(log types.Log) (*DataFeedsCacheInvalidUpdatePermission, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*DataFeedsCacheNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*DataFeedsCacheNewRound, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DataFeedsCacheOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*DataFeedsCacheOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DataFeedsCacheOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*DataFeedsCacheOwnershipTransferred, error)

	FilterProxyDataIdRemoved(opts *bind.FilterOpts, proxy []common.Address, dataId [][16]byte) (*DataFeedsCacheProxyDataIdRemovedIterator, error)

	WatchProxyDataIdRemoved(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheProxyDataIdRemoved, proxy []common.Address, dataId [][16]byte) (event.Subscription, error)

	ParseProxyDataIdRemoved(log types.Log) (*DataFeedsCacheProxyDataIdRemoved, error)

	FilterProxyDataIdUpdated(opts *bind.FilterOpts, proxy []common.Address, dataId [][16]byte) (*DataFeedsCacheProxyDataIdUpdatedIterator, error)

	WatchProxyDataIdUpdated(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheProxyDataIdUpdated, proxy []common.Address, dataId [][16]byte) (event.Subscription, error)

	ParseProxyDataIdUpdated(log types.Log) (*DataFeedsCacheProxyDataIdUpdated, error)

	FilterStaleBundleReport(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheStaleBundleReportIterator, error)

	WatchStaleBundleReport(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheStaleBundleReport, dataId [][16]byte) (event.Subscription, error)

	ParseStaleBundleReport(log types.Log) (*DataFeedsCacheStaleBundleReport, error)

	FilterStaleDecimalReport(opts *bind.FilterOpts, dataId [][16]byte) (*DataFeedsCacheStaleDecimalReportIterator, error)

	WatchStaleDecimalReport(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheStaleDecimalReport, dataId [][16]byte) (event.Subscription, error)

	ParseStaleDecimalReport(log types.Log) (*DataFeedsCacheStaleDecimalReport, error)

	FilterTokenRecovered(opts *bind.FilterOpts, token []common.Address, to []common.Address) (*DataFeedsCacheTokenRecoveredIterator, error)

	WatchTokenRecovered(opts *bind.WatchOpts, sink chan<- *DataFeedsCacheTokenRecovered, token []common.Address, to []common.Address) (event.Subscription, error)

	ParseTokenRecovered(log types.Log) (*DataFeedsCacheTokenRecovered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
