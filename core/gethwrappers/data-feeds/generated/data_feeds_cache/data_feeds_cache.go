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
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bundleDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkFeedPermission\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAnswer\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBundleDecimals\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDataIdForProxy\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDecimals\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getDescription\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeedMetadata\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"},{\"name\":\"startIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestAnswer\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestBundle\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"bundle\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestBundleTimestamp\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestRoundData\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestTimestamp\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimestamp\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isFeedAdmin\",\"inputs\":[{\"name\":\"feedAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestAnswer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestBundle\",\"inputs\":[],\"outputs\":[{\"name\":\"bundle\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestBundleTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRound\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverTokens\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeDataIdMappingsForProxies\",\"inputs\":[{\"name\":\"proxies\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeFeedConfigs\",\"inputs\":[{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBundleFeedConfigs\",\"inputs\":[{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"},{\"name\":\"descriptions\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"decimalsMatrix\",\"type\":\"uint8[][]\",\"internalType\":\"uint8[][]\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDecimalFeedConfigs\",\"inputs\":[{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"},{\"name\":\"descriptions\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFeedAdmin\",\"inputs\":[{\"name\":\"feedAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isAdmin\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateDataIdMappingsForProxies\",\"inputs\":[{\"name\":\"proxies\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"dataIds\",\"type\":\"bytes16[]\",\"internalType\":\"bytes16[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AnswerUpdated\",\"inputs\":[{\"name\":\"current\",\"type\":\"int256\",\"indexed\":true,\"internalType\":\"int256\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BundleFeedConfigSet\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"decimals\",\"type\":\"uint8[]\",\"indexed\":false,\"internalType\":\"uint8[]\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BundleReportUpdated\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"bundle\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DecimalFeedConfigSet\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"workflowMetadata\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"components\":[{\"name\":\"allowedSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DecimalReportUpdated\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"answer\",\"type\":\"uint224\",\"indexed\":false,\"internalType\":\"uint224\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeedAdminSet\",\"inputs\":[{\"name\":\"feedAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"isAdmin\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeedConfigRemoved\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidUpdatePermission\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"workflowOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"workflowName\",\"type\":\"bytes10\",\"indexed\":false,\"internalType\":\"bytes10\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewRound\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProxyDataIdRemoved\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProxyDataIdUpdated\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleBundleReport\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"reportTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"latestTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleDecimalReport\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"indexed\":true,\"internalType\":\"bytes16\"},{\"name\":\"reportTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"latestTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRecovered\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressInsufficientBalance\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ErrorSendingNative\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeedNotConfigured\",\"inputs\":[{\"name\":\"dataId\",\"type\":\"bytes16\",\"internalType\":\"bytes16\"}]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requiredBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidDataId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidWorkflowName\",\"inputs\":[{\"name\":\"workflowName\",\"type\":\"bytes10\",\"internalType\":\"bytes10\"}]},{\"type\":\"error\",\"name\":\"NoMappingForSender\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f80fd5b5033805f816100655760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b5f80546001600160a01b0319166001600160a01b0384811691909117909155811615610094576100948161009c565b505050610144565b336001600160a01b038216036100f45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005c565b600180546001600160a01b0319166001600160a01b038381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6148ee806101515f395ff3fe608060405234801561000f575f80fd5b5060043610610283575f3560e01c806379ba509711610157578063b5ab58dc116100d2578063ec52b1f511610088578063feaf968c1161006e578063feaf968c146105c0578063feb5d172146105c8578063ff25dbc81461065d575f80fd5b8063ec52b1f51461059a578063f2fde38b146105ad575f80fd5b8063be4f0a9f116100b8578063be4f0a9f14610554578063cdd2510014610567578063d143dcd91461057a575f80fd5b8063b5ab58dc1461052e578063b633620c14610541575f80fd5b80639198274f116101275780639a6fc8f51161010d5780639a6fc8f51461050b5780639d91348d1461051e578063a3d610cc14610526575f80fd5b80639198274f146104bf5780639608e18f146104c7575f80fd5b806379ba509714610482578063805f21321461048a5780638205bf6a1461049d5780638da5cb5b146104a5575f80fd5b80634533dc98116102015780635f25452b116101b7578063668a0f021161019d578063668a0f021461045f5780636a36e494146104675780637284e4161461047a575f80fd5b80635f25452b146104025780635f3e849f1461044c575f80fd5b806350d25bcd116101e757806350d25bcd146103df57806354fd4d50146103e7578063557a33c2146103ef575f80fd5b80634533dc98146103ac57806347381b08146103bf575f80fd5b8063297dbf561161025657806335f611221161023c57806335f611221461035b5780633a0449741461036e57806343d5ba5014610399575f80fd5b8063297dbf561461032c578063313ce56714610341575f80fd5b806301ffc9a71461028757806302ccb3ae146102af578063181f5a77146102cf5780631bb1610c1461030b575b5f80fd5b61029a610295366004613972565b610670565b60405190151581526020015b60405180910390f35b6102c26102bd3660046139cd565b6107ec565b6040516102a69190613a14565b6102c26040518060400160405280601481526020017f446174614665656473436163686520312e302e3000000000000000000000000081525081565b61031e6103193660046139cd565b6108d8565b6040519081526020016102a6565b61033f61033a366004613a6e565b61095c565b005b610349610b0d565b60405160ff90911681526020016102a6565b61033f610369366004613b1b565b610b71565b61029a61037c366004613bfd565b6001600160a01b03165f9081526007602052604090205460ff1690565b6103496103a73660046139cd565b6111da565b61033f6103ba366004613c18565b611225565b6103d26103cd3660046139cd565b611833565b6040516102a69190613cb7565b61031e6118f9565b61031e600781565b61031e6103fd3660046139cd565b611989565b6104156104103660046139cd565b611a05565b6040805169ffffffffffffffffffff968716815260208101959095528401929092526060830152909116608082015260a0016102a6565b61033f61045a366004613cfc565b611ac7565b61031e611d60565b6102c26104753660046139cd565b611dd2565b6102c2611e38565b61033f611f38565b61033f610498366004613d78565b61201a565b61031e612715565b5f546040516001600160a01b0390911681526020016102a6565b6102c26127ad565b6104f26104d5366004613bfd565b6001600160a01b03165f9081526002602052604090205460801b90565b6040516001600160801b031990911681526020016102a6565b610415610519366004613dd8565b612828565b6103d261290b565b61031e6129e6565b61031e61053c366004613e01565b612a61565b61031e61054f366004613e01565b612afe565b61033f610562366004613e18565b612ba4565b61033f610575366004613e18565b612c8e565b61058d610588366004613e57565b612f17565b6040516102a69190613e87565b61033f6105a8366004613f12565b613136565b61033f6105bb366004613bfd565b6131dc565b6104156131f0565b61029a6105d6366004614047565b805160208083015160409384015184516001600160801b031996909616868401526001600160a01b03938416868601529216606085015275ffffffffffffffffffffffffffffffffffffffffffff199091166080808501919091528251808503909101815260a090930182528251928101929092205f908152600990925290205460ff1690565b61031e61066b3660046139cd565b6132c7565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167fcce8054600000000000000000000000000000000000000000000000000000000148061070257507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b8061074e57507fffffffff0000000000000000000000000000000000000000000000000000000082167f805f213200000000000000000000000000000000000000000000000000000000145b8061079a57507fffffffff0000000000000000000000000000000000000000000000000000000082167f5f3e849f00000000000000000000000000000000000000000000000000000000145b806107e657507fffffffff0000000000000000000000000000000000000000000000000000000082167f181f5a7700000000000000000000000000000000000000000000000000000000145b92915050565b60606001600160801b0319821661082f576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031982165f908152600860205260409020600101805461085590614079565b80601f016020809104026020016040519081016040528092919081815260200182805461088190614079565b80156108cc5780601f106108a3576101008083540402835291602001916108cc565b820191905f5260205f20905b8154815290600101906020018083116108af57829003601f168201915b50505050509050919050565b5f6001600160801b0319821661091a576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506001600160801b0319165f908152600360205260409020547c0100000000000000000000000000000000000000000000000000000000900463ffffffff1690565b335f9081526007602052604090205460ff166109ab576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b828181146109e5576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5b81811015610b0557838382818110610a0157610a016140ca565b9050602002016020810190610a1691906139cd565b60025f888885818110610a2b57610a2b6140ca565b9050602002016020810190610a409190613bfd565b6001600160a01b0316815260208101919091526040015f2080546001600160801b03191660809290921c919091179055838382818110610a8257610a826140ca565b9050602002016020810190610a9791906139cd565b6001600160801b031916868683818110610ab357610ab36140ca565b9050602002016020810190610ac89190613bfd565b6001600160a01b03167ff31b9e58190970ef07c23d0ba78c358eb3b416e829ef484b29b9993a6b1b285a60405160405180910390a36001016109e7565b505050505050565b335f9081526002602052604081205460801b6001600160801b03198116610b62576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b610b6b8161332e565b91505090565b335f9081526007602052604090205460ff16610bbb576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b801580610bc6575086155b15610bfd576040517f60e8b63a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8685141580610c0c5750868314155b15610c43576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5b878110156111cf575f898983818110610c6057610c606140ca565b9050602002016020810190610c7591906139cd565b90506001600160801b03198116610cb8576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031981165f908152600860205260409020600281015415610e4a575f5b6002820154811015610dd5575f826002018281548110610cfe57610cfe6140ca565b5f9182526020808320604080516060808201835260029590950290920180546001600160a01b039081168085526001909201549081168486018190527401000000000000000000000000000000000000000090910460b01b75ffffffffffffffffffffffffffffffffffffffffffff191684840181905283516001600160801b03198d168188015280850193909352958201526080808201959095528151808203909501855260a0019052825192909101919091209092505f908152600960205260409020805460ff191690555050600101610cdc565b506001600160801b031982165f90815260086020526040812090610df982826137e2565b610e06600183015f613804565b610e13600283015f61383b565b50506040516001600160801b03198316907f871bcdef10dee59b87f17bab788b72faa8dfe1a9cc5bdc45c3baf4c18fa33910905f90a25b5f5b848110156110cf575f868683818110610e6757610e676140ca565b905060600201803603810190610e7d91906140f7565b9050845f03610f9f5780516001600160a01b0316610ed55780516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109a2565b60208101516001600160a01b0316610f2a5760208101516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109a2565b604081015175ffffffffffffffffffffffffffffffffffffffffffff1916610f9f5760408082015190517f114988d500000000000000000000000000000000000000000000000000000000815275ffffffffffffffffffffffffffffffffffffffffffff1990911660048201526024016109a2565b8051602080830180516040808601805182516001600160801b03198c16818801526001600160a01b0397881681850152938716606085015275ffffffffffffffffffffffffffffffffffffffffffff19166080808501919091528251808503909101815260a090930182528251928501929092205f90815260098552908120805460ff191660019081179091556002898101805480840182559084529590922096519490910290950180549385167fffffffffffffffffffffffff000000000000000000000000000000000000000090941693909317835590519184018054915160b01c74010000000000000000000000000000000000000000027fffff000000000000000000000000000000000000000000000000000000000000909216929093169190911717905501610e4c565b508686848181106110e2576110e26140ca565b90506020028101906110f49190614111565b6110ff918391613859565b50888884818110611112576111126140ca565b90506020028101906111249190614175565b6001830191611134919083614221565b506001600160801b031982167fdfebe0878c5611549f54908260ca12271c7ff3f0ebae0c1de47732612403869e888886818110611173576111736140ca565b90506020028101906111859190614111565b8c8c88818110611197576111976140ca565b90506020028101906111a99190614175565b8a8a6040516111bd96959493929190614392565b60405180910390a25050600101610c45565b505050505050505050565b5f6001600160801b0319821661121c576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107e68261332e565b335f9081526007602052604090205460ff1661126f576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b80158061127a575084155b156112b1576040517f60e8b63a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8483146112ea576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5b8581101561182a575f878783818110611307576113076140ca565b905060200201602081019061131c91906139cd565b90506001600160801b0319811661135f576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031981165f9081526008602052604090206002810154156114f1575f5b600282015481101561147c575f8260020182815481106113a5576113a56140ca565b5f9182526020808320604080516060808201835260029590950290920180546001600160a01b039081168085526001909201549081168486018190527401000000000000000000000000000000000000000090910460b01b75ffffffffffffffffffffffffffffffffffffffffffff191684840181905283516001600160801b03198d168188015280850193909352958201526080808201959095528151808203909501855260a0019052825192909101919091209092505f908152600960205260409020805460ff191690555050600101611383565b506001600160801b031982165f908152600860205260408120906114a082826137e2565b6114ad600183015f613804565b6114ba600283015f61383b565b50506040516001600160801b03198316907f871bcdef10dee59b87f17bab788b72faa8dfe1a9cc5bdc45c3baf4c18fa33910905f90a25b5f5b84811015611776575f86868381811061150e5761150e6140ca565b90506060020180360381019061152491906140f7565b9050845f036116465780516001600160a01b031661157c5780516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109a2565b60208101516001600160a01b03166115d15760208101516040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016109a2565b604081015175ffffffffffffffffffffffffffffffffffffffffffff19166116465760408082015190517f114988d500000000000000000000000000000000000000000000000000000000815275ffffffffffffffffffffffffffffffffffffffffffff1990911660048201526024016109a2565b8051602080830180516040808601805182516001600160801b03198c16818801526001600160a01b0397881681850152938716606085015275ffffffffffffffffffffffffffffffffffffffffffff19166080808501919091528251808503909101815260a090930182528251928501929092205f90815260098552908120805460ff191660019081179091556002898101805480840182559084529590922096519490910290950180549385167fffffffffffffffffffffffff000000000000000000000000000000000000000090941693909317835590519184018054915160b01c74010000000000000000000000000000000000000000027fffff0000000000000000000000000000000000000000000000000000000000009092169290931691909117179055016114f3565b50868684818110611789576117896140ca565b905060200281019061179b9190614175565b60018301916117ab919083614221565b506001600160801b031982167f2dec0e9ffbb18c6499fc8bee8b9c35f765e76d9dbd436f25dd00a80de267ac0d6117e18461332e565b8989878181106117f3576117f36140ca565b90506020028101906118059190614175565b8989604051611818959493929190614409565b60405180910390a250506001016112ec565b50505050505050565b60606001600160801b03198216611876576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031982165f90815260086020908152604091829020805483518184028101840190945280845290918301828280156108cc57602002820191905f5260205f20905f905b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116118c0575094979650505050505050565b335f9081526002602052604081205460801b6001600160801b0319811661194e576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b0319165f908152600360205260409020547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b5f6001600160801b031982166119cb576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506001600160801b0319165f908152600360205260409020547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1690565b5f808080806001600160801b03198616611a4b576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050506001600160801b0319919091165f9081526006602090815260408083205460039092529091205490927bffffffffffffffffffffffffffffffffffffffffffffffffffffffff821692507c010000000000000000000000000000000000000000000000000000000090910463ffffffff169081908490565b611acf6133ed565b6001600160a01b038316611bb85747811115611b20576040517fcf479181000000000000000000000000000000000000000000000000000000008152476004820152602481018290526044016109a2565b5f80836001600160a01b0316836040515f6040518083038185875af1925050503d805f8114611b6a576040519150601f19603f3d011682016040523d82523d5f602084013e611b6f565b606091505b509150915081611bb1578383826040517fc50febed0000000000000000000000000000000000000000000000000000000081526004016109a293929190614444565b5050611d0e565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b038416906370a0823190602401602060405180830381865afa158015611c13573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611c379190614474565b811115611cfa576040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b038416906370a0823190602401602060405180830381865afa158015611c99573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611cbd9190614474565b6040517fcf4791810000000000000000000000000000000000000000000000000000000081526004810191909152602481018290526044016109a2565b611d0e6001600160a01b0384168383613462565b816001600160a01b0316836001600160a01b03167f879f92dded0f26b83c3e00b12e0395dc72cfc3077343d1854ed6988edd1f909683604051611d5391815260200190565b60405180910390a3505050565b335f9081526002602052604081205460801b6001600160801b03198116611db5576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b0319165f90815260066020526040902054919050565b60606001600160801b03198216611e15576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160801b031982165f908152600560205260409020805461085590614079565b335f9081526002602052604090205460609060801b6001600160801b03198116611e90576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b031981165f9081526008602052604090206001018054611eb690614079565b80601f0160208091040260200160405190810160405280929190818152602001828054611ee290614079565b8015611f2d5780601f10611f0457610100808354040283529160200191611f2d565b820191905f5260205f20905b815481529060010190602001808311611f1057829003601f168201915b505050505091505090565b6001546001600160a01b03163314611fac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016109a2565b5f8054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b5f8061205a86868080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506134e292505050565b90925090505f61206e60406020868861448b565b612077916144b2565b90506120848160606144fc565b61208f906040614513565b8403612477575f6120a28587018761455c565b90505f5b82811015612470575f8282815181106120c1576120c16140ca565b6020908102919091018101518051604080516001600160801b031983168186015233818301526001600160a01b038b16606082015275ffffffffffffffffffffffffffffffffffffffffffff198a166080808301919091528251808303909101815260a090910182528051908501205f8181526009909552932054919350919060ff166121c157604080513381526001600160a01b038a16602082015275ffffffffffffffffffffffffffffffffffffffffffff198916918101919091526001600160801b03198316907feeeaa8bf618ff6d960c6cf5935e68384f066abcc8b95d0de91bd773c16ae3ae3906060015b60405180910390a2505050612468565b6001600160801b031982165f908152600360209081526040909120549084015163ffffffff7c01000000000000000000000000000000000000000000000000000000009092048216911611612290576020838101516001600160801b031984165f8181526003845260409081902054815163ffffffff94851681527c010000000000000000000000000000000000000000000000000000000090910490931693830193909352917fcf16f5f704f981fa2279afa1877dd1fdaa462a03a71ec51b9d3b2416a59a013e91016121b1565b604080518082018252848201517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16815260208086015163ffffffff16818301526001600160801b031985165f9081526006909152918220805491929182906122f490614648565b91829055506001600160801b031985165f8181526003602090815260408083208751888401805163ffffffff9081167c01000000000000000000000000000000000000000000000000000000009081027bffffffffffffffffffffffffffffffffffffffffffffffffffffffff94851617909455888752600486528487208888528652958490208a519151909616928302911690811790945590519283529394508492917f82584589cd7284d4503ed582275e22b2e8f459f9cf4170a7235844e367f966d5910160405180910390a460208086015160405163ffffffff90911681525f9183917f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271910160405180910390a38085604001517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f4260405161245a91815260200190565b60405180910390a350505050505b6001016120a6565b505061182a565b5f61248485870187614660565b90505f5b81518110156111cf575f8282815181106124a4576124a46140ca565b6020908102919091018101518051604080516001600160801b031983168186015233818301526001600160a01b038b16606082015275ffffffffffffffffffffffffffffffffffffffffffff198a166080808301919091528251808303909101815260a090910182528051908501205f8181526009909552932054919350919060ff166125a457604080513381526001600160a01b038a16602082015275ffffffffffffffffffffffffffffffffffffffffffff198916918101919091526001600160801b03198316907feeeaa8bf618ff6d960c6cf5935e68384f066abcc8b95d0de91bd773c16ae3ae3906060015b60405180910390a250505061270d565b6001600160801b031982165f908152600560209081526040909120600101549084015163ffffffff918216911611612637576020838101516001600160801b031984165f8181526005845260409081902060010154815163ffffffff9485168152931693830193909352917f51001b67094834cc084a0c1feb791cf84a481357aa66b924ba205d4cb56fd9819101612594565b60408051808201825284820151815260208086015163ffffffff16818301526001600160801b031985165f9081526005909152919091208151829190819061267f90826147c9565b5060209182015160019190910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff92831617905590820151825160405191909216916001600160801b03198616917f1dc1bef0b59d624eab3f0ec044781bb5b8594cd64f0ba09d789f5b51acab16149161270091613a14565b60405180910390a3505050505b600101612488565b335f9081526002602052604081205460801b6001600160801b0319811661276a576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b0319165f908152600360205260409020547c0100000000000000000000000000000000000000000000000000000000900463ffffffff16919050565b335f9081526002602052604090205460609060801b6001600160801b03198116612805576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b031981165f9081526005602052604090208054611eb690614079565b335f90815260026020526040812054819081908190819060801b6001600160801b03198116612885576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b69ffffffffffffffffffff87165f9081526004602090815260408083206001600160801b0319949094168352929052205495967bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8716967c0100000000000000000000000000000000000000000000000000000000900463ffffffff169550859450879350915050565b335f9081526002602052604090205460609060801b6001600160801b03198116612963576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b031981165f9081526008602090815260409182902080548351818402810184019094528084529091830182828015611f2d57602002820191905f5260205f20905f905b825461010083900a900460ff168152602060019283018181049485019490930390920291018084116129ad579050505050505091505090565b335f9081526002602052604081205460801b6001600160801b03198116612a3b576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b0319165f9081526005602052604090206001015463ffffffff16919050565b335f9081526002602052604081205460801b6001600160801b03198116612ab6576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b5f9283526004602090815260408085206001600160801b03199093168552919052909120547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff16919050565b335f9081526002602052604081205460801b6001600160801b03198116612b53576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b5f9283526004602090815260408085206001600160801b031990931685529190529091205463ffffffff7c010000000000000000000000000000000000000000000000000000000090910416919050565b335f9081526007602052604090205460ff16612bee576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b805f5b81811015612c88575f848483818110612c0c57612c0c6140ca565b9050602002016020810190612c219190613bfd565b6001600160a01b0381165f8181526002602052604080822080546001600160801b0319808216909255915194955060809190911b9390841692917f4200186b7bc2d4f13f7888c5bbe9461d57da88705be86521f3d78be691ad1d2a91a35050600101612bf1565b50505050565b335f9081526007602052604090205460ff16612cd8576040517fd86ad9cf0000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b5f5b81811015612f12575f838383818110612cf557612cf56140ca565b9050602002016020810190612d0a91906139cd565b6001600160801b031981165f9081526008602052604081206002015491925003612d6c576040517f8606a85b0000000000000000000000000000000000000000000000000000000081526001600160801b0319821660048201526024016109a2565b5f5b6001600160801b031982165f90815260086020526040902060020154811015612e95576001600160801b031982165f908152600860205260408120600201805483908110612dbe57612dbe6140ca565b5f9182526020808320604080516060808201835260029590950290920180546001600160a01b039081168085526001909201549081168486018190527401000000000000000000000000000000000000000090910460b01b75ffffffffffffffffffffffffffffffffffffffffffff191684840181905283516001600160801b03198c168188015280850193909352958201526080808201959095528151808203909501855260a0019052825192909101919091209092505f908152600960205260409020805460ff191690555050600101612d6e565b506001600160801b031981165f90815260086020526040812090612eb982826137e2565b612ec6600183015f613804565b612ed3600283015f61383b565b50506040516001600160801b03198216907f871bcdef10dee59b87f17bab788b72faa8dfe1a9cc5bdc45c3baf4c18fa33910905f90a250600101612cda565b505050565b6001600160801b031983165f9081526008602052604081206002810154606092819003612f7c576040517f8606a85b0000000000000000000000000000000000000000000000000000000081526001600160801b0319871660048201526024016109a2565b808510612fcd57604080515f8082526020820190925290612fc3565b604080516060810182525f80825260208083018290529282015282525f19909201910181612f985790505b509250505061312f565b5f612fd88587614513565b905081811180612fe6575084155b612ff05780612ff2565b815b9050612ffe8682614884565b67ffffffffffffffff81111561301657613016613f49565b60405190808252806020026020018201604052801561305f57816020015b604080516060810182525f80825260208083018290529282015282525f199092019101816130345790505b5093505f5b845181101561312a576002840161307b8883614513565b8154811061308b5761308b6140ca565b5f9182526020918290206040805160608101825260029390930290910180546001600160a01b039081168452600190910154908116938301939093527401000000000000000000000000000000000000000090920460b01b75ffffffffffffffffffffffffffffffffffffffffffff1916918101919091528551869083908110613117576131176140ca565b6020908102919091010152600101613064565b505050505b9392505050565b61313e6133ed565b6001600160a01b038216613189576040517f8e4c8aa60000000000000000000000000000000000000000000000000000000081526001600160a01b03831660048201526024016109a2565b6001600160a01b0382165f81815260076020526040808220805460ff191685151590811790915590519092917f93a3fa5993d2a54de369386625330cc6d73caee7fece4b3983cf299b264473fd91a35050565b6131e46133ed565b6131ed816134f3565b50565b335f90815260026020526040812054819081908190819060801b6001600160801b0319811661324d576040517f718b09d00000000000000000000000000000000000000000000000000000000081523360048201526024016109a2565b6001600160801b0319165f9081526006602090815260408083205460039092529091205490967bffffffffffffffffffffffffffffffffffffffffffffffffffffffff821696507c010000000000000000000000000000000000000000000000000000000090910463ffffffff1694508493508692509050565b5f6001600160801b03198216613309576040517f0760371200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506001600160801b0319165f9081526005602052604090206001015463ffffffff1690565b5f8061333b8360076135cd565b90507f20000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000008216108015906133d157507f60000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821611155b156133e55761312f602060f883901c614897565b505f92915050565b5f546001600160a01b03163314613460576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016109a2565b565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052612f12908490613634565b6040810151604a9091015160601c91565b336001600160a01b03821603613565576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016109a2565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b038381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6040516001600160801b0319831660208201525f906030016040516020818303038152906040528281518110613605576136056140ca565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016905092915050565b5f6136486001600160a01b038416836136ae565b905080515f1415801561366c57508080602001905181019061366a91906148b0565b155b15612f12576040517f5274afe70000000000000000000000000000000000000000000000000000000081526001600160a01b03841660048201526024016109a2565b606061312f83835f845f80856001600160a01b031684866040516136d291906148cb565b5f6040518083038185875af1925050503d805f811461370c576040519150601f19603f3d011682016040523d82523d5f602084013e613711565b606091505b509150915061372186838361372b565b9695505050505050565b6060826137405761373b826137a0565b61312f565b815115801561375757506001600160a01b0384163b155b15613799576040517f9996b3150000000000000000000000000000000000000000000000000000000081526001600160a01b03851660048201526024016109a2565b508061312f565b8051156137b05780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5080545f8255601f0160209004905f5260205f20908101906131ed91906138fd565b50805461381090614079565b5f825580601f1061381f575050565b601f0160209004905f5260205f20908101906131ed91906138fd565b5080545f8255600202905f5260205f20908101906131ed9190613911565b828054828255905f5260205f2090601f016020900481019282156138ed579160200282015f5b838211156138bf57833560ff1683826101000a81548160ff021916908360ff16021790555092602001926001016020815f0104928301926001030261387f565b80156138eb5782816101000a81549060ff02191690556001016020815f010492830192600103026138bf565b505b506138f99291506138fd565b5090565b5b808211156138f9575f81556001016138fe565b5b808211156138f95780547fffffffffffffffffffffffff00000000000000000000000000000000000000001681556001810180547fffff000000000000000000000000000000000000000000000000000000000000169055600201613912565b5f60208284031215613982575f80fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461312f575f80fd5b80356001600160801b0319811681146139c8575f80fd5b919050565b5f602082840312156139dd575f80fd5b61312f826139b1565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f61312f60208301846139e6565b5f8083601f840112613a36575f80fd5b50813567ffffffffffffffff811115613a4d575f80fd5b6020830191508360208260051b8501011115613a67575f80fd5b9250929050565b5f805f8060408587031215613a81575f80fd5b843567ffffffffffffffff811115613a97575f80fd5b613aa387828801613a26565b909550935050602085013567ffffffffffffffff811115613ac2575f80fd5b613ace87828801613a26565b95989497509550505050565b5f8083601f840112613aea575f80fd5b50813567ffffffffffffffff811115613b01575f80fd5b602083019150836020606083028501011115613a67575f80fd5b5f805f805f805f806080898b031215613b32575f80fd5b883567ffffffffffffffff811115613b48575f80fd5b613b548b828c01613a26565b909950975050602089013567ffffffffffffffff811115613b73575f80fd5b613b7f8b828c01613a26565b909750955050604089013567ffffffffffffffff811115613b9e575f80fd5b613baa8b828c01613a26565b909550935050606089013567ffffffffffffffff811115613bc9575f80fd5b613bd58b828c01613ada565b999c989b5096995094979396929594505050565b6001600160a01b03811681146131ed575f80fd5b5f60208284031215613c0d575f80fd5b813561312f81613be9565b5f805f805f8060608789031215613c2d575f80fd5b863567ffffffffffffffff811115613c43575f80fd5b613c4f89828a01613a26565b909750955050602087013567ffffffffffffffff811115613c6e575f80fd5b613c7a89828a01613a26565b909550935050604087013567ffffffffffffffff811115613c99575f80fd5b613ca589828a01613ada565b979a9699509497509295939492505050565b602080825282518282018190525f918401906040840190835b81811015613cf157835160ff16835260209384019390920191600101613cd0565b509095945050505050565b5f805f60608486031215613d0e575f80fd5b8335613d1981613be9565b92506020840135613d2981613be9565b929592945050506040919091013590565b5f8083601f840112613d4a575f80fd5b50813567ffffffffffffffff811115613d61575f80fd5b602083019150836020828501011115613a67575f80fd5b5f805f8060408587031215613d8b575f80fd5b843567ffffffffffffffff811115613da1575f80fd5b613dad87828801613d3a565b909550935050602085013567ffffffffffffffff811115613dcc575f80fd5b613ace87828801613d3a565b5f60208284031215613de8575f80fd5b813569ffffffffffffffffffff8116811461312f575f80fd5b5f60208284031215613e11575f80fd5b5035919050565b5f8060208385031215613e29575f80fd5b823567ffffffffffffffff811115613e3f575f80fd5b613e4b85828601613a26565b90969095509350505050565b5f805f60608486031215613e69575f80fd5b613e72846139b1565b95602085013595506040909401359392505050565b602080825282518282018190525f918401906040840190835b81811015613cf15783516001600160a01b0381511684526001600160a01b03602082015116602085015275ffffffffffffffffffffffffffffffffffffffffffff19604082015116604085015250606083019250602084019350600181019050613ea0565b80151581146131ed575f80fd5b5f8060408385031215613f23575f80fd5b8235613f2e81613be9565b91506020830135613f3e81613f05565b809150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040516060810167ffffffffffffffff81118282101715613f9957613f99613f49565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715613fc857613fc8613f49565b604052919050565b803575ffffffffffffffffffffffffffffffffffffffffffff19811681146139c8575f80fd5b5f60608284031215614006575f80fd5b61400e613f76565b9050813561401b81613be9565b8152602082013561402b81613be9565b602082015261403c60408301613fd0565b604082015292915050565b5f8060808385031215614058575f80fd5b614061836139b1565b91506140708460208501613ff6565b90509250929050565b600181811c9082168061408d57607f821691505b6020821081036140c4577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f60608284031215614107575f80fd5b61312f8383613ff6565b5f8083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112614144575f80fd5b83018035915067ffffffffffffffff82111561415e575f80fd5b6020019150600581901b3603821315613a67575f80fd5b5f8083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126141a8575f80fd5b83018035915067ffffffffffffffff8211156141c2575f80fd5b602001915036819003821315613a67575f80fd5b601f821115612f1257805f5260205f20601f840160051c810160208510156141fb5750805b601f840160051c820191505b8181101561421a575f8155600101614207565b5050505050565b67ffffffffffffffff83111561423957614239613f49565b61424d836142478354614079565b836141d6565b5f601f84116001811461427e575f85156142675750838201355b5f19600387901b1c1916600186901b17835561421a565b5f83815260208120601f198716915b828110156142ad578685013582556020948501946001909201910161428d565b50868210156142c9575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b81835281816020850137505f602082840101525f6020601f19601f840116840101905092915050565b8183526020830192505f815f5b8481101561438857813561432481613be9565b6001600160a01b03168652602082013561433d81613be9565b6001600160a01b0316602087015275ffffffffffffffffffffffffffffffffffffffffffff1961436f60408401613fd0565b1660408701526060958601959190910190600101614311565b5093949350505050565b606080825281018690525f8760808301825b898110156143d257823560ff81168082146143bd575f80fd5b835250602092830192909101906001016143a4565b5083810360208501526143e681888a6142db565b91505082810360408401526143fc818587614304565b9998505050505050505050565b60ff86168152606060208201525f6144256060830186886142db565b8281036040840152614438818587614304565b98975050505050505050565b6001600160a01b0384168152826020820152606060408201525f61446b60608301846139e6565b95945050505050565b5f60208284031215614484575f80fd5b5051919050565b5f8085851115614499575f80fd5b838611156144a5575f80fd5b5050820193919092039150565b803560208310156107e6575f19602084900360031b1b1692915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b80820281158282048414176107e6576107e66144cf565b808201808211156107e6576107e66144cf565b5f67ffffffffffffffff82111561453f5761453f613f49565b5060051b60200190565b803563ffffffff811681146139c8575f80fd5b5f6020828403121561456c575f80fd5b813567ffffffffffffffff811115614582575f80fd5b8201601f81018413614592575f80fd5b80356145a56145a082614526565b613f9f565b808282526020820191506020606084028501019250868311156145c6575f80fd5b6020840193505b8284101561372157606084880312156145e4575f80fd5b6145ec613f76565b843581526145fc60208601614549565b602082015260408501357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8116811461462f575f80fd5b60408201528252606093909301926020909101906145cd565b5f5f198203614659576146596144cf565b5060010190565b5f60208284031215614670575f80fd5b813567ffffffffffffffff811115614686575f80fd5b8201601f81018413614696575f80fd5b80356146a46145a082614526565b8082825260208201915060208360051b8501019250868311156146c5575f80fd5b602084015b838110156147be57803567ffffffffffffffff8111156146e8575f80fd5b85016060818a03601f190112156146fd575f80fd5b614705613f76565b6020820135815261471860408301614549565b6020820152606082013567ffffffffffffffff811115614736575f80fd5b60208184010192505089601f83011261474d575f80fd5b813567ffffffffffffffff81111561476757614767613f49565b61477a6020601f19601f84011601613f9f565b8181528b602083860101111561478e575f80fd5b816020850160208301375f60208383010152806040840152505080855250506020830192506020810190506146ca565b509695505050505050565b815167ffffffffffffffff8111156147e3576147e3613f49565b6147f7816147f18454614079565b846141d6565b6020601f821160018114614829575f83156148125750848201515b5f19600385901b1c1916600184901b17845561421a565b5f84815260208120601f198516915b828110156148585787850151825560209485019460019092019101614838565b508482101561487557868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b818103818111156107e6576107e66144cf565b60ff82811682821603908111156107e6576107e66144cf565b5f602082840312156148c0575f80fd5b815161312f81613f05565b5f82518060208501845e5f92019182525091905056fea164736f6c634300081a000a",
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

func (_DataFeedsCache *DataFeedsCacheCaller) GetLatestRoundData(opts *bind.CallOpts, dataId [16]byte) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getLatestRoundData", dataId)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetLatestRoundData(dataId [16]byte) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _DataFeedsCache.Contract.GetLatestRoundData(&_DataFeedsCache.CallOpts, dataId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetLatestRoundData(dataId [16]byte) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
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

func (_DataFeedsCache *DataFeedsCacheCaller) GetRoundData(opts *bind.CallOpts, roundId *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "getRoundData", roundId)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, err

}

func (_DataFeedsCache *DataFeedsCacheSession) GetRoundData(roundId *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _DataFeedsCache.Contract.GetRoundData(&_DataFeedsCache.CallOpts, roundId)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) GetRoundData(roundId *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
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

func (_DataFeedsCache *DataFeedsCacheCaller) LatestRoundData(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _DataFeedsCache.contract.Call(opts, &out, "latestRoundData")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, err

}

func (_DataFeedsCache *DataFeedsCacheSession) LatestRoundData() (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _DataFeedsCache.Contract.LatestRoundData(&_DataFeedsCache.CallOpts)
}

func (_DataFeedsCache *DataFeedsCacheCallerSession) LatestRoundData() (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
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

	GetLatestRoundData(opts *bind.CallOpts, dataId [16]byte) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error)

	GetLatestTimestamp(opts *bind.CallOpts, dataId [16]byte) (*big.Int, error)

	GetRoundData(opts *bind.CallOpts, roundId *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error)

	GetTimestamp(opts *bind.CallOpts, roundId *big.Int) (*big.Int, error)

	IsFeedAdmin(opts *bind.CallOpts, feedAdmin common.Address) (bool, error)

	LatestAnswer(opts *bind.CallOpts) (*big.Int, error)

	LatestBundle(opts *bind.CallOpts) ([]byte, error)

	LatestBundleTimestamp(opts *bind.CallOpts) (*big.Int, error)

	LatestRound(opts *bind.CallOpts) (*big.Int, error)

	LatestRoundData(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error)

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
