// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifiable_load_log_trigger_upkeep_wrapper

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

type Log struct {
	Index       *big.Int
	TxIndex     *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var VerifiableLoadLogTriggerUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460c0526101a0604052604261012081815260e091829190620052386101403981526020016040518060800160405280604281526020016200527a60429139905262000099906016906002620003a2565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000c990826200051e565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152601890620000fb90826200051e565b503480156200010957600080fd5b50604051620052bc380380620052bc8339810160408190526200012c9162000600565b81813380600081620001855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001b857620001b881620002f7565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa15801562000215573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200023b919062000643565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002a1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002c7919062000674565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560a052506200069b915050565b336001600160a01b03821603620003515760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200017c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215620003ed579160200282015b82811115620003ed5782518290620003dc90826200051e565b5091602001919060010190620003c3565b50620003fb929150620003ff565b5090565b80821115620003fb57600062000416828262000420565b50600101620003ff565b5080546200042e906200048f565b6000825580601f106200043f575050565b601f0160209004906000526020600020908101906200045f919062000462565b50565b5b80821115620003fb576000815560010162000463565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004a457607f821691505b602082108103620004c557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200051957600081815260208120601f850160051c81016020861015620004f45750805b601f850160051c820191505b81811015620005155782815560010162000500565b5050505b505050565b81516001600160401b038111156200053a576200053a62000479565b62000552816200054b84546200048f565b84620004cb565b602080601f8311600181146200058a5760008415620005715750858301515b600019600386901b1c1916600185901b17855562000515565b600085815260208120601f198616915b82811015620005bb578886015182559484019460019091019084016200059a565b5085821015620005da5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200045f57600080fd5b600080604083850312156200061457600080fd5b82516200062181620005ea565b602084015190925080151581146200063857600080fd5b809150509250929050565b600080604083850312156200065757600080fd5b82516200066481620005ea565b6020939093015192949293505050565b6000602082840312156200068757600080fd5b81516200069481620005ea565b9392505050565b60805160a05160c051614b51620006e7600039600081816104f00152611bf101526000818161082c015261303e015260008181610b480152818161139e01526118560152614b516000f3fe6080604052600436106103c75760003560e01c806379ba5097116101f2578063a79c40431161010d578063d6051a72116100a0578063e45530831161006f578063e455308314610dc7578063f2fde38b14610ddd578063fba7ffa314610dfd578063fcdc1f6314610e2a57600080fd5b8063d6051a7214610d3a578063daee1aeb14610d5a578063dbef701e14610d7a578063e0114adb14610d9a57600080fd5b8063c357f1f3116100dc578063c357f1f314610c4c578063c804802214610ca6578063c98f10b014610cc6578063d355852814610cdb57600080fd5b8063a79c404314610bca578063af953a4a14610bf7578063afb28d1f14610c17578063becde0e114610c2c57600080fd5b80639b42935411610185578063a654824811610154578063a654824814610b36578063a6b5947514610b6a578063a6c60d8914610b8a578063a72aa27e14610baa57600080fd5b80639b42935414610a985780639b51fb0d14610ac55780639d385eaa14610af65780639d6f1cc714610b1657600080fd5b80638fcb3fba116101c15780638fcb3fba146109ef578063924ca57814610a1c578063948108f714610a3c5780639ac542eb14610a5c57600080fd5b806379ba5097146109625780637b103999146109775780637e7a46dc146109a45780638da5cb5b146109c457600080fd5b806346e7a63e116102e2578063636092e8116102755780637145f11b116102445780637145f11b146108b857806373644cce146108e85780637672130314610915578063776898c81461094257600080fd5b8063636092e8146107d8578063642f6cef1461081a57806369cdbadb1461085e57806369e9b7731461088b57600080fd5b806359710992116102b157806359710992146107615780635d4ee7f3146107765780635f17e6161461078b57806360457ff5146107ab57600080fd5b806346e7a63e146106c75780634b56a42e146106f457806351c98be31461071457806357970e931461073457600080fd5b806328c4b57b1161035a5780633ebe8d6c116103295780633ebe8d6c1461063957806340691db4146106595780634585e33b1461068757806345d2ec17146106a757600080fd5b806328c4b57b1461057a5780632a9032d31461059a5780632b20e397146105ba578063328ffd111461060c57600080fd5b80630d4a4fb1116103965780630d4a4fb1146104b157806312c55027146104de578063206c32e81461052557806320e3dbd41461055a57600080fd5b806305e251311461040b57806306c1cc001461042d57806306e3b6321461044d578063077ac6211461048357600080fd5b3661040657604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561041757600080fd5b5061042b610426366004613787565b610e57565b005b34801561043957600080fd5b5061042b6104483660046138a8565b610e6e565b34801561045957600080fd5b5061046d610468366004613944565b61122a565b60405161047a9190613966565b60405180910390f35b34801561048f57600080fd5b506104a361049e3660046139c1565b611329565b60405190815260200161047a565b3480156104bd57600080fd5b506104d16104cc3660046139f6565b611367565b60405161047a9190613a7d565b3480156104ea57600080fd5b506105127f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff909116815260200161047a565b34801561053157600080fd5b50610545610540366004613a90565b611484565b6040805192835260208301919091520161047a565b34801561056657600080fd5b5061042b610575366004613ade565b611507565b34801561058657600080fd5b506104a3610595366004613afb565b6116d1565b3480156105a657600080fd5b5061042b6105b5366004613b6c565b61173c565b3480156105c657600080fd5b506011546105e79073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161047a565b34801561061857600080fd5b506104a36106273660046139f6565b60036020526000908152604090205481565b34801561064557600080fd5b506104a36106543660046139f6565b6117d6565b34801561066557600080fd5b50610679610674366004613bae565b61183f565b60405161047a929190613c1b565b34801561069357600080fd5b5061042b6106a2366004613c78565b611af0565b3480156106b357600080fd5b5061046d6106c2366004613a90565b611d41565b3480156106d357600080fd5b506104a36106e23660046139f6565b600a6020526000908152604090205481565b34801561070057600080fd5b5061067961070f366004613cae565b611db0565b34801561072057600080fd5b5061042b61072f366004613d6b565b611e04565b34801561074057600080fd5b506012546105e79073ffffffffffffffffffffffffffffffffffffffff1681565b34801561076d57600080fd5b5061042b611ea8565b34801561078257600080fd5b5061042b612093565b34801561079757600080fd5b5061042b6107a6366004613944565b6121ca565b3480156107b757600080fd5b506104a36107c63660046139f6565b60076020526000908152604090205481565b3480156107e457600080fd5b506015546107fd906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff909116815260200161047a565b34801561082657600080fd5b5061084e7f000000000000000000000000000000000000000000000000000000000000000081565b604051901515815260200161047a565b34801561086a57600080fd5b506104a36108793660046139f6565b60086020526000908152604090205481565b34801561089757600080fd5b5061042b6108a6366004613944565b60009182526008602052604090912055565b3480156108c457600080fd5b5061084e6108d33660046139f6565b600b6020526000908152604090205460ff1681565b3480156108f457600080fd5b506104a36109033660046139f6565b6000908152600c602052604090205490565b34801561092157600080fd5b506104a36109303660046139f6565b60046020526000908152604090205481565b34801561094e57600080fd5b5061084e61095d3660046139f6565b612297565b34801561096e57600080fd5b5061042b6122e9565b34801561098357600080fd5b506013546105e79073ffffffffffffffffffffffffffffffffffffffff1681565b3480156109b057600080fd5b5061042b6109bf366004613dc2565b6123e6565b3480156109d057600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166105e7565b3480156109fb57600080fd5b506104a3610a0a3660046139f6565b60056020526000908152604090205481565b348015610a2857600080fd5b5061042b610a37366004613944565b612477565b348015610a4857600080fd5b5061042b610a57366004613e0e565b6126bc565b348015610a6857600080fd5b50601554610a86906c01000000000000000000000000900460ff1681565b60405160ff909116815260200161047a565b348015610aa457600080fd5b5061042b610ab3366004613944565b60009182526009602052604090912055565b348015610ad157600080fd5b50610512610ae03660046139f6565b600e6020526000908152604090205461ffff1681565b348015610b0257600080fd5b5061046d610b113660046139f6565b612805565b348015610b2257600080fd5b506104d1610b313660046139f6565b612867565b348015610b4257600080fd5b506104a37f000000000000000000000000000000000000000000000000000000000000000081565b348015610b7657600080fd5b5061042b610b85366004613afb565b612913565b348015610b9657600080fd5b5061042b610ba53660046139f6565b601455565b348015610bb657600080fd5b5061042b610bc5366004613e3e565b61297c565b348015610bd657600080fd5b5061042b610be5366004613944565b60009182526007602052604090912055565b348015610c0357600080fd5b5061042b610c123660046139f6565b612a27565b348015610c2357600080fd5b506104d1612aad565b348015610c3857600080fd5b5061042b610c47366004613b6c565b612aba565b348015610c5857600080fd5b5061042b610c67366004613e63565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610cb257600080fd5b5061042b610cc13660046139f6565b612b54565b348015610cd257600080fd5b506104d1612bec565b348015610ce757600080fd5b5061042b610cf6366004613e80565b6015805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b348015610d4657600080fd5b50610545610d55366004613944565b612bf9565b348015610d6657600080fd5b5061042b610d75366004613b6c565b612c62565b348015610d8657600080fd5b506104a3610d95366004613944565b612d2d565b348015610da657600080fd5b506104a3610db53660046139f6565b60096020526000908152604090205481565b348015610dd357600080fd5b506104a360145481565b348015610de957600080fd5b5061042b610df8366004613ade565b612d5e565b348015610e0957600080fd5b506104a3610e183660046139f6565b60066020526000908152604090205481565b348015610e3657600080fd5b506104a3610e453660046139f6565b60026020526000908152604090205481565b8051610e6a906016906020840190613557565b5050565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b3921690610f54908c1688613ecc565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af1158015610fd2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ff69190613f10565b5060008860ff1667ffffffffffffffff81111561101557611015613637565b60405190808252806020026020018201604052801561103e578160200160208202803683370190505b50905060005b8960ff168160ff1610156111e757600061105d84612d72565b90508860ff16600103611195576040517f0d4a4fb1000000000000000000000000000000000000000000000000000000008152600481018290526000903090630d4a4fb190602401600060405180830381865afa1580156110c2573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526111089190810190613f78565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906111619085908590600401613fad565b600060405180830381600087803b15801561117b57600080fd5b505af115801561118f573d6000803e3d6000fd5b50505050505b80838360ff16815181106111ab576111ab613fc6565b602090810291909101810191909152600091825260088152604080832088905560079091529020849055806111df81613ff5565b915050611044565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711816040516112179190613966565b60405180910390a1505050505050505050565b60606000611238600f612e40565b9050808410611273576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82600003611288576112858482614014565b92505b60008367ffffffffffffffff8111156112a3576112a3613637565b6040519080825280602002602001820160405280156112cc578160200160208202803683370190505b50905060005b8481101561131e576112ef6112e78288614027565b600f90612e4a565b82828151811061130157611301613fc6565b6020908102919091010152806113168161403a565b9150506112d2565b509150505b92915050565b600d602052826000526040600020602052816000526040600020818154811061135157600080fd5b9060005260206000200160009250925050505481565b606060006040518060c001604052803073ffffffffffffffffffffffffffffffffffffffff168152602001600160ff1681526020017f00000000000000000000000000000000000000000000000000000000000000008152602001846040516020016113d591815260200190565b6040516020818303038152906040526113ed90614072565b81526020016000801b81526020016000801b81525090508060405160200161146d9190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b604051602081830303815290604052915050919050565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156114e857602002820191906000526020600020905b8154815260200190600101908083116114d4575b505050505090506114fa818251612e56565b92509250505b9250929050565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa15801561159d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115c191906140c2565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015611664573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061168891906140f0565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b6000838152600c6020908152604080832080548251818502810185019093528083526117329383018282801561172657602002820191906000526020600020905b815481526020019060010190808311611712575b50505050508484612edb565b90505b9392505050565b8060005b818160ff1610156117d0573063c8048022858560ff851681811061176657611766613fc6565b905060200201356040518263ffffffff1660e01b815260040161178b91815260200190565b600060405180830381600087803b1580156117a557600080fd5b505af11580156117b9573d6000803e3d6000fd5b5050505080806117c890613ff5565b915050611740565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611837576000858152600d6020908152604080832061ffff851684529091529020546118239083614027565b91508061182f8161410d565b9150506117ec565b509392505050565b6000606060005a9050600061185261303a565b90507f000000000000000000000000000000000000000000000000000000000000000061188260c088018861412e565b600081811061189357611893613fc6565b9050602002013503611a685760006118ae60c088018861412e565b60018181106118bf576118bf613fc6565b905060200201356040516020016118d891815260200190565b60405160208183030381529060405290506000818060200190518101906118ff9190614196565b9050600061191060c08a018a61412e565b600281811061192157611921613fc6565b9050602002013560405160200161193a91815260200190565b60405160208183030381529060405290506000818060200190518101906119619190614196565b6000848152600860205260409020549091505b805a6119809089614014565b61198c90613a98614027565b10156119da5781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055816119d2816141af565b925050611974565b6017601660188487866040516020016119fd929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e000000000000000000000000000000000000000000000000000000008252611a5f95949392916004016142cc565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401611a5f565b60005a9050600080611b0484860186613cae565b9150915060008082806020019051810190611b1f919061438f565b6000828152600560209081526040808320546004909252822054939550919350909190611b4a61303a565b905082600003611b6a576000858152600560205260409020819055611cae565b6000611b768583614014565b6000878152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611be857602002820191906000526020600020905b815481526020019060010190808311611bd4575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff16815103611c635781611c258161410d565b60008a8152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055898352600c8252928220805493840181558252902001555b600085815260066020526040812054611cc8906001614027565b6000878152600660209081526040808320849055600490915290208390559050611cf28683612477565b604051308152829087907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a3611d34868a84612913565b5050505050505050505050565b6000828152600d6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015611da357602002820191906000526020600020905b815481526020019060010190808311611d8f575b5050505050905092915050565b6000606060008484604051602001611dc99291906143b3565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b8160005b81811015611ea15730635f17e616868684818110611e2857611e28613fc6565b90506020020135856040518363ffffffff1660e01b8152600401611e5c92919091825263ffffffff16602082015260400190565b600060405180830381600087803b158015611e7657600080fd5b505af1158015611e8a573d6000803e3d6000fd5b505050508080611e999061403a565b915050611e08565b5050505050565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600060048201819052602482018190529173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015611f1f573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611f659190810190614447565b80519091506000611f7461303a565b905060005b828110156117d0576000848281518110611f9557611f95613fc6565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612015573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061203991906144d8565b90508060ff1660010361207e57604051308152849083907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a35b5050808061208b9061403a565b915050611f79565b61209b6130dc565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561210a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061212e9190614196565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156121a6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e6a9190613f10565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c9091528120612202916135ad565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff161161225e576000848152600d6020908152604080832061ffff85168452909152812061224c916135ad565b806122568161410d565b915050612217565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b60008181526005602052604081205481036122b457506001919050565b6000828152600360209081526040808320546004909252909120546122d761303a565b6122e19190614014565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461236a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401611a5f565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b590612440908690869086906004016144f5565b600060405180830381600087803b15801561245a57600080fd5b505af115801561246e573d6000803e3d6000fd5b50505050505050565b6014546000838152600260205260409020546124939083614014565b1115610e6a576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612509573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261254f9190810190614577565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa1580156125c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125e89190614696565b60155490915061260c9082906c01000000000000000000000000900460ff16613ecc565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff1610156117d05760155461264f9085906bffffffffffffffffffffffff166126bc565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612744573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127689190613f10565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b1580156127e957600080fd5b505af11580156127fd573d6000803e3d6000fd5b505050505050565b6000818152600c602090815260409182902080548351818402810184019094528084526060939283018282801561285b57602002820191906000526020600020905b815481526020019060010190808311612847575b50505050509050919050565b6016818154811061287757600080fd5b906000526020600020016000915090508054612892906141e4565b80601f01602080910402602001604051908101604052809291908181526020018280546128be906141e4565b801561290b5780601f106128e05761010080835404028352916020019161290b565b820191906000526020600020905b8154815290600101906020018083116128ee57829003601f168201915b505050505081565b6000838152600760205260409020545b805a61292f9085614014565b61293b90612710614027565b10156117d05781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612923565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156129f457600080fd5b505af1158015612a08573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b158015612a9957600080fd5b505af1158015611ea1573d6000803e3d6000fd5b60178054612892906141e4565b8060005b818163ffffffff1610156117d0573063af953a4a858563ffffffff8516818110612aea57612aea613fc6565b905060200201356040518263ffffffff1660e01b8152600401612b0f91815260200190565b600060405180830381600087803b158015612b2957600080fd5b505af1158015612b3d573d6000803e3d6000fd5b505050508080612b4c906146b3565b915050612abe565b6013546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b158015612bc057600080fd5b505af1158015612bd4573d6000803e3d6000fd5b50505050610e6a81600f61315f90919063ffffffff16565b60188054612892906141e4565b6000828152600c60209081526040808320805482518185028101850190935280835284938493929190830182828015612c5157602002820191906000526020600020905b815481526020019060010190808311612c3d575b505050505090506114fa8185612e56565b8060005b818110156117d0576000848483818110612c8257612c82613fc6565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc8283604051602001612cbb91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401612ce7929190613fad565b600060405180830381600087803b158015612d0157600080fd5b505af1158015612d15573d6000803e3d6000fd5b50505050508080612d259061403a565b915050612c66565b600c6020528160005260406000208181548110612d4957600080fd5b90600052602060002001600091509150505481565b612d666130dc565b612d6f8161316b565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190612dcd9086906004016146cc565b6020604051808303816000875af1158015612dec573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e109190614196565b9050612e1d600f82613260565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b6000611323825490565b6000611735838361326c565b815160009081908190841580612e6c5750808510155b15612e75578094505b60008092505b85831015612ed157866001612e908585614014565b612e9a9190614014565b81518110612eaa57612eaa613fc6565b602002602001015181612ebd9190614027565b905082612ec98161403a565b935050612e7b565b9694955050505050565b82516000908190831580612eef5750808410155b15612ef8578093505b60008467ffffffffffffffff811115612f1357612f13613637565b604051908082528060200260200182016040528015612f3c578160200160208202803683370190505b509050600092505b84831015612faa57866001612f598585614014565b612f639190614014565b81518110612f7357612f73613fc6565b6020026020010151818481518110612f8d57612f8d613fc6565b602090810291909101015282612fa28161403a565b935050612f44565b612fc381600060018451612fbe9190614014565b613296565b85606403612ffc578060018251612fda9190614014565b81518110612fea57612fea613fc6565b60200260200101519350505050611735565b80606482518861300c919061481e565b613016919061488a565b8151811061302657613026613fc6565b602002602001015193505050509392505050565b60007f0000000000000000000000000000000000000000000000000000000000000000156130d757606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156130ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130d29190614196565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331461315d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611a5f565b565b6000611735838361340e565b3373ffffffffffffffffffffffffffffffffffffffff8216036131ea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611a5f565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006117358383613508565b600082600001828154811061328357613283613fc6565b9060005260206000200154905092915050565b81818082036132a6575050505050565b60008560026132b5878761489e565b6132bf91906148be565b6132c99087614926565b815181106132d9576132d9613fc6565b602002602001015190505b8183136133e8575b808684815181106132ff576132ff613fc6565b6020026020010151101561331f57826133178161494e565b9350506132ec565b85828151811061333157613331613fc6565b6020026020010151811015613352578161334a8161497f565b92505061331f565b8183136133e35785828151811061336b5761336b613fc6565b602002602001015186848151811061338557613385613fc6565b602002602001015187858151811061339f5761339f613fc6565b602002602001018885815181106133b8576133b8613fc6565b602090810291909101019190915252826133d18161494e565b93505081806133df9061497f565b9250505b6132e4565b818512156133fb576133fb868684613296565b838312156127fd576127fd868486613296565b600081815260018301602052604081205480156134f7576000613432600183614014565b855490915060009061344690600190614014565b90508181146134ab57600086600001828154811061346657613466613fc6565b906000526020600020015490508087600001848154811061348957613489613fc6565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806134bc576134bc6149b0565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611323565b6000915050611323565b5092915050565b600081815260018301602052604081205461354f57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611323565b506000611323565b82805482825590600052602060002090810192821561359d579160200282015b8281111561359d578251829061358d9082614a2a565b5091602001919060010190613577565b506135a99291506135cb565b5090565b5080546000825590600052602060002090810190612d6f91906135e8565b808211156135a95760006135df82826135fd565b506001016135cb565b5b808211156135a957600081556001016135e9565b508054613609906141e4565b6000825580601f10613619575050565b601f016020900490600052602060002090810190612d6f91906135e8565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff8111828210171561368a5761368a613637565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156136d7576136d7613637565b604052919050565b600067ffffffffffffffff8211156136f9576136f9613637565b5060051b60200190565b600067ffffffffffffffff82111561371d5761371d613637565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600061375c61375784613703565b613690565b905082815283838301111561377057600080fd5b828260208301376000602084830101529392505050565b6000602080838503121561379a57600080fd5b823567ffffffffffffffff808211156137b257600080fd5b818501915085601f8301126137c657600080fd5b81356137d4613757826136df565b81815260059190911b830184019084810190888311156137f357600080fd5b8585015b838110156138405780358581111561380f5760008081fd5b8601603f81018b136138215760008081fd5b6138328b8983013560408401613749565b8452509186019186016137f7565b5098975050505050505050565b60ff81168114612d6f57600080fd5b63ffffffff81168114612d6f57600080fd5b600082601f83011261387f57600080fd5b61173583833560208501613749565b6bffffffffffffffffffffffff81168114612d6f57600080fd5b600080600080600080600060e0888a0312156138c357600080fd5b87356138ce8161384d565b965060208801356138de8161385c565b955060408801356138ee8161384d565b9450606088013567ffffffffffffffff81111561390a57600080fd5b6139168a828b0161386e565b94505060808801356139278161388e565b9699959850939692959460a0840135945060c09093013592915050565b6000806040838503121561395757600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b8181101561399e57835183529284019291840191600101613982565b50909695505050505050565b803561ffff811681146139bc57600080fd5b919050565b6000806000606084860312156139d657600080fd5b833592506139e6602085016139aa565b9150604084013590509250925092565b600060208284031215613a0857600080fd5b5035919050565b60005b83811015613a2a578181015183820152602001613a12565b50506000910152565b60008151808452613a4b816020860160208601613a0f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006117356020830184613a33565b60008060408385031215613aa357600080fd5b82359150613ab3602084016139aa565b90509250929050565b73ffffffffffffffffffffffffffffffffffffffff81168114612d6f57600080fd5b600060208284031215613af057600080fd5b813561173581613abc565b600080600060608486031215613b1057600080fd5b505081359360208301359350604090920135919050565b60008083601f840112613b3957600080fd5b50813567ffffffffffffffff811115613b5157600080fd5b6020830191508360208260051b850101111561150057600080fd5b60008060208385031215613b7f57600080fd5b823567ffffffffffffffff811115613b9657600080fd5b613ba285828601613b27565b90969095509350505050565b60008060408385031215613bc157600080fd5b823567ffffffffffffffff80821115613bd957600080fd5b908401906101008287031215613bee57600080fd5b90925060208401359080821115613c0457600080fd5b50613c118582860161386e565b9150509250929050565b82151581526040602082015260006117326040830184613a33565b60008083601f840112613c4857600080fd5b50813567ffffffffffffffff811115613c6057600080fd5b60208301915083602082850101111561150057600080fd5b60008060208385031215613c8b57600080fd5b823567ffffffffffffffff811115613ca257600080fd5b613ba285828601613c36565b60008060408385031215613cc157600080fd5b823567ffffffffffffffff80821115613cd957600080fd5b818501915085601f830112613ced57600080fd5b81356020613cfd613757836136df565b82815260059290921b84018101918181019089841115613d1c57600080fd5b8286015b84811015613d5457803586811115613d385760008081fd5b613d468c86838b010161386e565b845250918301918301613d20565b5096505086013592505080821115613c0457600080fd5b600080600060408486031215613d8057600080fd5b833567ffffffffffffffff811115613d9757600080fd5b613da386828701613b27565b9094509250506020840135613db78161385c565b809150509250925092565b600080600060408486031215613dd757600080fd5b83359250602084013567ffffffffffffffff811115613df557600080fd5b613e0186828701613c36565b9497909650939450505050565b60008060408385031215613e2157600080fd5b823591506020830135613e338161388e565b809150509250929050565b60008060408385031215613e5157600080fd5b823591506020830135613e338161385c565b600060208284031215613e7557600080fd5b81356117358161388e565b600060208284031215613e9257600080fd5b81356117358161384d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615613ef757613ef7613e9d565b02949350505050565b805180151581146139bc57600080fd5b600060208284031215613f2257600080fd5b61173582613f00565b600082601f830112613f3c57600080fd5b8151613f4a61375782613703565b818152846020838601011115613f5f57600080fd5b613f70826020830160208701613a0f565b949350505050565b600060208284031215613f8a57600080fd5b815167ffffffffffffffff811115613fa157600080fd5b613f7084828501613f2b565b8281526040602082015260006117326040830184613a33565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff810361400b5761400b613e9d565b60010192915050565b8181038181111561132357611323613e9d565b8082018082111561132357611323613e9d565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361406b5761406b613e9d565b5060010190565b805160208083015191908110156140b1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b80516139bc81613abc565b600080604083850312156140d557600080fd5b82516140e081613abc565b6020939093015192949293505050565b60006020828403121561410257600080fd5b815161173581613abc565b600061ffff80831681810361412457614124613e9d565b6001019392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261416357600080fd5b83018035915067ffffffffffffffff82111561417e57600080fd5b6020019150600581901b360382131561150057600080fd5b6000602082840312156141a857600080fd5b5051919050565b6000816141be576141be613e9d565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600181811c908216806141f857607f821691505b6020821081036140b1577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000815461423e816141e4565b80855260206001838116801561425b5760018114614293576142c1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b89010195506142c1565b866000528260002060005b858110156142b95781548a820186015290830190840161429e565b890184019650505b505050505092915050565b60a0815260006142df60a0830188614231565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015614351577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087840301855261433f8383614231565b94860194925060019182019101614306565b50508681036040880152614365818b614231565b94505050505084606084015282810360808401526143838185613a33565b98975050505050505050565b600080604083850312156143a257600080fd5b505080516020909101519092909150565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015614428577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552614416868351613a33565b955093820193908201906001016143dc565b50508584038187015250505061443e8185613a33565b95945050505050565b6000602080838503121561445a57600080fd5b825167ffffffffffffffff81111561447157600080fd5b8301601f8101851361448257600080fd5b8051614490613757826136df565b81815260059190911b820183019083810190878311156144af57600080fd5b928401925b828410156144cd578351825292840192908401906144b4565b979650505050505050565b6000602082840312156144ea57600080fd5b81516117358161384d565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b80516139bc8161385c565b80516139bc8161388e565b805167ffffffffffffffff811681146139bc57600080fd5b60006020828403121561458957600080fd5b815167ffffffffffffffff808211156145a157600080fd5b9083019061014082860312156145b657600080fd5b6145be613666565b6145c7836140b7565b81526145d560208401614549565b60208201526040830151828111156145ec57600080fd5b6145f887828601613f2b565b60408301525061460a60608401614554565b606082015261461b608084016140b7565b608082015261462c60a0840161455f565b60a082015261463d60c08401614549565b60c082015261464e60e08401614554565b60e0820152610100614661818501613f00565b90820152610120838101518381111561467957600080fd5b61468588828701613f2b565b918301919091525095945050505050565b6000602082840312156146a857600080fd5b81516117358161388e565b600063ffffffff80831681810361412457614124613e9d565b60208152600082516101408060208501526146eb610160850183613a33565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808685030160408701526147278483613a33565b935060408701519150614752606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526147b38483613a33565b935060e087015191506101008187860301818801526147d28584613a33565b9450808801519250506101208187860301818801526147f18584613a33565b94508088015192505050614814828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561485657614856613e9d565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826148995761489961485b565b500490565b818103600083128015838313168383128216171561350157613501613e9d565b6000826148cd576148cd61485b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561492157614921613e9d565b500590565b808201828112600083128015821682158216171561494657614946613e9d565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361406b5761406b613e9d565b60007f800000000000000000000000000000000000000000000000000000000000000082036141be576141be613e9d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b601f821115614a2557600081815260208120601f850160051c81016020861015614a065750805b601f850160051c820191505b818110156127fd57828155600101614a12565b505050565b815167ffffffffffffffff811115614a4457614a44613637565b614a5881614a5284546141e4565b846149df565b602080601f831160018114614aab5760008415614a755750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556127fd565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015614af857888601518255948401946001909101908401614ad9565b5085821015614b3457878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var VerifiableLoadLogTriggerUpkeepABI = VerifiableLoadLogTriggerUpkeepMetaData.ABI

var VerifiableLoadLogTriggerUpkeepBin = VerifiableLoadLogTriggerUpkeepMetaData.Bin

func DeployVerifiableLoadLogTriggerUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _registrar common.Address, _useArb bool) (common.Address, *types.Transaction, *VerifiableLoadLogTriggerUpkeep, error) {
	parsed, err := VerifiableLoadLogTriggerUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadLogTriggerUpkeepBin), backend, _registrar, _useArb)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifiableLoadLogTriggerUpkeep{VerifiableLoadLogTriggerUpkeepCaller: VerifiableLoadLogTriggerUpkeepCaller{contract: contract}, VerifiableLoadLogTriggerUpkeepTransactor: VerifiableLoadLogTriggerUpkeepTransactor{contract: contract}, VerifiableLoadLogTriggerUpkeepFilterer: VerifiableLoadLogTriggerUpkeepFilterer{contract: contract}}, nil
}

type VerifiableLoadLogTriggerUpkeep struct {
	address common.Address
	abi     abi.ABI
	VerifiableLoadLogTriggerUpkeepCaller
	VerifiableLoadLogTriggerUpkeepTransactor
	VerifiableLoadLogTriggerUpkeepFilterer
}

type VerifiableLoadLogTriggerUpkeepCaller struct {
	contract *bind.BoundContract
}

type VerifiableLoadLogTriggerUpkeepTransactor struct {
	contract *bind.BoundContract
}

type VerifiableLoadLogTriggerUpkeepFilterer struct {
	contract *bind.BoundContract
}

type VerifiableLoadLogTriggerUpkeepSession struct {
	Contract     *VerifiableLoadLogTriggerUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VerifiableLoadLogTriggerUpkeepCallerSession struct {
	Contract *VerifiableLoadLogTriggerUpkeepCaller
	CallOpts bind.CallOpts
}

type VerifiableLoadLogTriggerUpkeepTransactorSession struct {
	Contract     *VerifiableLoadLogTriggerUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type VerifiableLoadLogTriggerUpkeepRaw struct {
	Contract *VerifiableLoadLogTriggerUpkeep
}

type VerifiableLoadLogTriggerUpkeepCallerRaw struct {
	Contract *VerifiableLoadLogTriggerUpkeepCaller
}

type VerifiableLoadLogTriggerUpkeepTransactorRaw struct {
	Contract *VerifiableLoadLogTriggerUpkeepTransactor
}

func NewVerifiableLoadLogTriggerUpkeep(address common.Address, backend bind.ContractBackend) (*VerifiableLoadLogTriggerUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(VerifiableLoadLogTriggerUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeep{address: address, abi: abi, VerifiableLoadLogTriggerUpkeepCaller: VerifiableLoadLogTriggerUpkeepCaller{contract: contract}, VerifiableLoadLogTriggerUpkeepTransactor: VerifiableLoadLogTriggerUpkeepTransactor{contract: contract}, VerifiableLoadLogTriggerUpkeepFilterer: VerifiableLoadLogTriggerUpkeepFilterer{contract: contract}}, nil
}

func NewVerifiableLoadLogTriggerUpkeepCaller(address common.Address, caller bind.ContractCaller) (*VerifiableLoadLogTriggerUpkeepCaller, error) {
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepCaller{contract: contract}, nil
}

func NewVerifiableLoadLogTriggerUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifiableLoadLogTriggerUpkeepTransactor, error) {
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepTransactor{contract: contract}, nil
}

func NewVerifiableLoadLogTriggerUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifiableLoadLogTriggerUpkeepFilterer, error) {
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepFilterer{contract: contract}, nil
}

func bindVerifiableLoadLogTriggerUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifiableLoadLogTriggerUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadLogTriggerUpkeep.Contract.VerifiableLoadLogTriggerUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.VerifiableLoadLogTriggerUpkeepTransactor.contract.Transfer(opts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.VerifiableLoadLogTriggerUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadLogTriggerUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.contract.Transfer(opts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) BUCKETSIZE(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "BUCKET_SIZE")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) AddLinkAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "addLinkAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddLinkAmount(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddLinkAmount(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "bucketedDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "buckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Buckets(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Buckets(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckCallback(&_VerifiableLoadLogTriggerUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckCallback(&_VerifiableLoadLogTriggerUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "checkGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "counters", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Counters(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Counters(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "delays", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Delays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Delays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.DummyMap(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.DummyMap(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "eligible", upkeepId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Eligible(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Eligible(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) EmittedSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "emittedSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.EmittedSig(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.EmittedSig(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) FeedParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) FeedParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedsHex(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedsHex(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "firstPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "gasLimits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GasLimits(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GasLimits(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetActiveUpkeepIDs(&_VerifiableLoadLogTriggerUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetActiveUpkeepIDs(&_VerifiableLoadLogTriggerUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getBucketedDelays", upkeepId, bucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getDelays", upkeepId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getLogTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetLogTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetLogTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getPxDelayLastNPerforms", upkeepId, p, n)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getSumDelayInBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getSumDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Intervals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "intervals", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Intervals(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Intervals(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) LastTopUpBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "lastTopUpBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LinkToken(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LinkToken(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) MinBalanceThresholdMultiplier(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "minBalanceThresholdMultiplier")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Owner() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Owner(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Owner() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Owner(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) PerformDataSizes(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "performDataSizes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformDataSizes(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformDataSizes(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) PerformGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "performGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) PreviousPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "previousPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Registrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "registrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Registrar() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registrar(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Registrar() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registrar(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Registry() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registry(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Registry() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registry(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TimeParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimeParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) TimeParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimeParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) UpkeepTopUpCheckInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "upkeepTopUpCheckInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "useArbitrumBlockNum")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "acceptOwnership")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AcceptOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AcceptOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) AddFunds(opts *bind.TransactOpts, upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "addFunds", upkeepId, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddFunds(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddFunds(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchCancelUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchCancelUpkeeps", upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchRegisterUpkeeps", number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchSendLogs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchSendLogs")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchSendLogs() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSendLogs(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchSendLogs() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSendLogs(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchSetIntervals", upkeepIds, interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchUpdatePipelineData", upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchWithdrawLinks", upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BurnPerformGas(opts *bind.TransactOpts, upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "burnPerformGas", upkeepId, startGas, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BurnPerformGas(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BurnPerformGas(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "cancelUpkeep", upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CancelUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CancelUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) CheckLog(opts *bind.TransactOpts, log Log, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "checkLog", log, checkData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckLog(log Log, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log, checkData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) CheckLog(log Log, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log, checkData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetAddLinkAmount(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setAddLinkAmount", amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetAddLinkAmount(amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetAddLinkAmount(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetAddLinkAmount(amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetAddLinkAmount(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetCheckGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setCheckGasToBurn", upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetCheckGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetCheckGasToBurn(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetCheckGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetCheckGasToBurn(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setConfig", newRegistrar)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetConfig(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetConfig(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setFeedsHex", newFeeds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetFeedsHex(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newFeeds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetFeedsHex(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newFeeds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setInterval", upkeepId, _interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetInterval(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetInterval(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetMinBalanceThresholdMultiplier(opts *bind.TransactOpts, newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setMinBalanceThresholdMultiplier", newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetMinBalanceThresholdMultiplier(newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetMinBalanceThresholdMultiplier(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetMinBalanceThresholdMultiplier(newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetMinBalanceThresholdMultiplier(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setPerformDataSize", upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetPerformGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setPerformGasToBurn", upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetPerformGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetPerformGasToBurn(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetPerformGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetPerformGasToBurn(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setUpkeepGasLimit", upkeepId, gasLimit)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetUpkeepTopUpCheckInterval(opts *bind.TransactOpts, newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setUpkeepTopUpCheckInterval", newInterval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetUpkeepTopUpCheckInterval(newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepTopUpCheckInterval(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newInterval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetUpkeepTopUpCheckInterval(newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepTopUpCheckInterval(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newInterval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) TopUpFund(opts *bind.TransactOpts, upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "topUpFund", upkeepId, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TopUpFund(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TopUpFund(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "transferOwnership", to)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TransferOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, to)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TransferOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, to)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "updateUpkeepPipelineData", upkeepId, pipelineData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "withdrawLinks")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "withdrawLinks0", upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.RawTransact(opts, nil)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Receive(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Receive(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

type VerifiableLoadLogTriggerUpkeepLogEmittedIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepLogEmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepLogEmitted)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepLogEmitted)
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

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepLogEmitted struct {
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int) (*VerifiableLoadLogTriggerUpkeepLogEmittedIterator, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumRule []interface{}
	for _, blockNumItem := range blockNum {
		blockNumRule = append(blockNumRule, blockNumItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepLogEmittedIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "LogEmitted", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int) (event.Subscription, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumRule []interface{}
	for _, blockNumItem := range blockNum {
		blockNumRule = append(blockNumRule, blockNumItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepLogEmitted)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseLogEmitted(log types.Log) (*VerifiableLoadLogTriggerUpkeepLogEmitted, error) {
	event := new(VerifiableLoadLogTriggerUpkeepLogEmitted)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
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

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, error) {
	event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
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

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseOwnershipTransferred(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferred, error) {
	event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepReceivedIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepReceived)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepReceived)
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

func (it *VerifiableLoadLogTriggerUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepReceived struct {
	Sender common.Address
	Value  *big.Int
	Raw    types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterReceived(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepReceivedIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepReceivedIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "Received", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchReceived(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepReceived) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepReceived)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "Received", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseReceived(log types.Log) (*VerifiableLoadLogTriggerUpkeepReceived, error) {
	event := new(VerifiableLoadLogTriggerUpkeepReceived)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "Received", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepUpkeepTopUp

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
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

func (it *VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepUpkeepTopUp struct {
	UpkeepId *big.Int
	Amount   *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "UpkeepTopUp", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepTopUp) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseUpkeepTopUp(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUp, error) {
	event := new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepUpkeepsRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepsRegistered)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepsRegistered)
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

func (it *VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepUpkeepsRegistered struct {
	UpkeepIds []*big.Int
	Raw       types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterUpkeepsRegistered(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "UpkeepsRegistered")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "UpkeepsRegistered", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchUpkeepsRegistered(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepsRegistered) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "UpkeepsRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepUpkeepsRegistered)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepsRegistered", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseUpkeepsRegistered(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepsRegistered, error) {
	event := new(VerifiableLoadLogTriggerUpkeepUpkeepsRegistered)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepsRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["Received"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseReceived(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseUpkeepTopUp(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["UpkeepsRegistered"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseUpkeepsRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadLogTriggerUpkeepLogEmitted) Topic() common.Hash {
	return common.HexToHash("0x97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf08")
}

func (VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifiableLoadLogTriggerUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifiableLoadLogTriggerUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874")
}

func (VerifiableLoadLogTriggerUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (VerifiableLoadLogTriggerUpkeepUpkeepsRegistered) Topic() common.Hash {
	return common.HexToHash("0x2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeep) Address() common.Address {
	return _VerifiableLoadLogTriggerUpkeep.address
}

type VerifiableLoadLogTriggerUpkeepInterface interface {
	BUCKETSIZE(opts *bind.CallOpts) (uint16, error)

	AddLinkAmount(opts *bind.CallOpts) (*big.Int, error)

	BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error)

	Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error)

	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error)

	EmittedSig(opts *bind.CallOpts) ([32]byte, error)

	FeedParamKey(opts *bind.CallOpts) (string, error)

	FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error)

	GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error)

	GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error)

	GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error)

	GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error)

	Intervals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	LastTopUpBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	LinkToken(opts *bind.CallOpts) (common.Address, error)

	MinBalanceThresholdMultiplier(opts *bind.CallOpts) (uint8, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	PerformDataSizes(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	PerformGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	PreviousPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Registrar(opts *bind.CallOpts) (common.Address, error)

	Registry(opts *bind.CallOpts) (common.Address, error)

	TimeParamKey(opts *bind.CallOpts) (string, error)

	UpkeepTopUpCheckInterval(opts *bind.CallOpts) (*big.Int, error)

	UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddFunds(opts *bind.TransactOpts, upkeepId *big.Int, amount *big.Int) (*types.Transaction, error)

	BatchCancelUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error)

	BatchSendLogs(opts *bind.TransactOpts) (*types.Transaction, error)

	BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error)

	BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	BurnPerformGas(opts *bind.TransactOpts, upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	CheckLog(opts *bind.TransactOpts, log Log, checkData []byte) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetAddLinkAmount(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	SetCheckGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error)

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetMinBalanceThresholdMultiplier(opts *bind.TransactOpts, newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error)

	SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetPerformGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepTopUpCheckInterval(opts *bind.TransactOpts, newInterval *big.Int) (*types.Transaction, error)

	TopUpFund(opts *bind.TransactOpts, upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error)

	WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int) (*VerifiableLoadLogTriggerUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadLogTriggerUpkeepLogEmitted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferred, error)

	FilterReceived(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepReceivedIterator, error)

	WatchReceived(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepReceived) (event.Subscription, error)

	ParseReceived(log types.Log) (*VerifiableLoadLogTriggerUpkeepReceived, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUp, error)

	FilterUpkeepsRegistered(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator, error)

	WatchUpkeepsRegistered(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepsRegistered) (event.Subscription, error)

	ParseUpkeepsRegistered(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepsRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
