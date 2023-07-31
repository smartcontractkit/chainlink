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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_autoLog\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_useMercury\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"autoLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_autoLog\",\"type\":\"bool\"}],\"name\":\"setAutoLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_useMercury\",\"type\":\"bool\"}],\"name\":\"setUseMercury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useMercury\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460c0526101a0604052604261012081815260e091829190620054fe6101403981526020016040518060800160405280604281526020016200554060429139905262000099906016906002620003c4565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000c9908262000540565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152601890620000fb908262000540565b503480156200010957600080fd5b5060405162005582380380620055828339810160408190526200012c9162000638565b83833380600081620001855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001b857620001b88162000319565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa15801562000215573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200023b919062000697565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002a1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002c79190620006c8565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560a052506019805461ffff191692151561ff001916929092176101009115159190910217905550620006ef9050565b336001600160a01b03821603620003735760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200017c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156200040f579160200282015b828111156200040f5782518290620003fe908262000540565b5091602001919060010190620003e5565b506200041d92915062000421565b5090565b808211156200041d57600062000438828262000442565b5060010162000421565b5080546200045090620004b1565b6000825580601f1062000461575050565b601f01602090049060005260206000209081019062000481919062000484565b50565b5b808211156200041d576000815560010162000485565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004c657607f821691505b602082108103620004e757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200053b57600081815260208120601f850160051c81016020861015620005165750805b601f850160051c820191505b81811015620005375782815560010162000522565b5050505b505050565b81516001600160401b038111156200055c576200055c6200049b565b62000574816200056d8454620004b1565b84620004ed565b602080601f831160018114620005ac5760008415620005935750858301515b600019600386901b1c1916600185901b17855562000537565b600085815260208120601f198616915b82811015620005dd57888601518255948401946001909101908401620005bc565b5085821015620005fc5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200048157600080fd5b805180151581146200063357600080fd5b919050565b600080600080608085870312156200064f57600080fd5b84516200065c816200060c565b93506200066c6020860162000622565b92506200067c6040860162000622565b91506200068c6060860162000622565b905092959194509250565b60008060408385031215620006ab57600080fd5b8251620006b8816200060c565b6020939093015192949293505050565b600060208284031215620006db57600080fd5b8151620006e8816200060c565b9392505050565b60805160a05160c051614dc36200073b600039600081816105970152611e3e0152600081816108d301526130b3015260008181610c410152818161151c0152611a1c0152614dc36000f3fe60806040526004361061044e5760003560e01c806379ba509711610243578063af953a4a11610143578063d90c4a68116100bb578063e15e1b341161008a578063f2fde38b1161006f578063f2fde38b14610f5b578063fba7ffa314610f7b578063fcdc1f6314610fa857600080fd5b8063e15e1b3414610ef9578063e455308314610f4557600080fd5b8063d90c4a6814610e72578063daee1aeb14610e8c578063dbef701e14610eac578063e0114adb14610ecc57600080fd5b8063c41c815b11610112578063c98f10b0116100f7578063c98f10b014610dde578063d355852814610df3578063d6051a7214610e5257600080fd5b8063c41c815b14610d9f578063c804802214610dbe57600080fd5b8063af953a4a14610cf0578063afb28d1f14610d10578063becde0e114610d25578063c357f1f314610d4557600080fd5b80639b429354116101d6578063a6548248116101a5578063a6c60d891161018a578063a6c60d8914610c83578063a72aa27e14610ca3578063a79c404314610cc357600080fd5b8063a654824814610c2f578063a6b5947514610c6357600080fd5b80639b42935414610b915780639b51fb0d14610bbe5780639d385eaa14610bef5780639d6f1cc714610c0f57600080fd5b80638fcb3fba116102125780638fcb3fba14610ae8578063924ca57814610b15578063948108f714610b355780639ac542eb14610b5557600080fd5b806379ba509714610a5b5780637b10399914610a705780637e7a46dc14610a9d5780638da5cb5b14610abd57600080fd5b806346e7a63e1161034e578063636092e8116102e15780636b2120c6116102b057806373644cce1161029557806373644cce146109e15780637672130314610a0e578063776898c814610a3b57600080fd5b80636b2120c61461095f5780637145f11b146109b157600080fd5b8063636092e81461087f578063642f6cef146108c157806369cdbadb1461090557806369e9b7731461093257600080fd5b8063597109921161031d57806359710992146108085780635d4ee7f31461081d5780635f17e6161461083257806360457ff51461085257600080fd5b806346e7a63e1461076e5780634b56a42e1461079b57806351c98be3146107bb57806357970e93146107db57600080fd5b806320e3dbd4116103e1578063328ffd11116103b057806340691db41161039557806340691db4146107005780634585e33b1461072e57806345d2ec171461074e57600080fd5b8063328ffd11146106b35780633ebe8d6c146106e057600080fd5b806320e3dbd41461060157806328c4b57b146106215780632a9032d3146106415780632b20e3971461066157600080fd5b80630d4a4fb11161041d5780630d4a4fb1146105385780630e577d421461056557806312c5502714610585578063206c32e8146105cc57600080fd5b806305e251311461049257806306c1cc00146104b457806306e3b632146104d4578063077ac6211461050a57600080fd5b3661048d57604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561049e57600080fd5b506104b26104ad3660046139e0565b610fd5565b005b3480156104c057600080fd5b506104b26104cf366004613b01565b610fec565b3480156104e057600080fd5b506104f46104ef366004613b9d565b6113a8565b6040516105019190613bbf565b60405180910390f35b34801561051657600080fd5b5061052a610525366004613c1a565b6114a7565b604051908152602001610501565b34801561054457600080fd5b50610558610553366004613c4f565b6114e5565b6040516105019190613cd6565b34801561057157600080fd5b506104b2610580366004613c4f565b611602565b34801561059157600080fd5b506105b97f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610501565b3480156105d857600080fd5b506105ec6105e7366004613ce9565b61164a565b60408051928352602083019190915201610501565b34801561060d57600080fd5b506104b261061c366004613d37565b6116cd565b34801561062d57600080fd5b5061052a61063c366004613d54565b611897565b34801561064d57600080fd5b506104b261065c366004613dc5565b611902565b34801561066d57600080fd5b5060115461068e9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610501565b3480156106bf57600080fd5b5061052a6106ce366004613c4f565b60036020526000908152604090205481565b3480156106ec57600080fd5b5061052a6106fb366004613c4f565b61199c565b34801561070c57600080fd5b5061072061071b366004613e07565b611a05565b604051610501929190613e74565b34801561073a57600080fd5b506104b2610749366004613ed1565b611d3d565b34801561075a57600080fd5b506104f4610769366004613ce9565b611f9a565b34801561077a57600080fd5b5061052a610789366004613c4f565b600a6020526000908152604090205481565b3480156107a757600080fd5b506107206107b6366004613f07565b612009565b3480156107c757600080fd5b506104b26107d6366004613fc4565b61205d565b3480156107e757600080fd5b5060125461068e9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561081457600080fd5b506104b2612101565b34801561082957600080fd5b506104b26122ec565b34801561083e57600080fd5b506104b261084d366004613b9d565b612423565b34801561085e57600080fd5b5061052a61086d366004613c4f565b60076020526000908152604090205481565b34801561088b57600080fd5b506015546108a4906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff9091168152602001610501565b3480156108cd57600080fd5b506108f57f000000000000000000000000000000000000000000000000000000000000000081565b6040519015158152602001610501565b34801561091157600080fd5b5061052a610920366004613c4f565b60086020526000908152604090205481565b34801561093e57600080fd5b506104b261094d366004613b9d565b60009182526008602052604090912055565b34801561096b57600080fd5b506104b261097a366004614029565b60198054911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff909216919091179055565b3480156109bd57600080fd5b506108f56109cc366004613c4f565b600b6020526000908152604090205460ff1681565b3480156109ed57600080fd5b5061052a6109fc366004613c4f565b6000908152600c602052604090205490565b348015610a1a57600080fd5b5061052a610a29366004613c4f565b60046020526000908152604090205481565b348015610a4757600080fd5b506108f5610a56366004613c4f565b6124f0565b348015610a6757600080fd5b506104b2612542565b348015610a7c57600080fd5b5060135461068e9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610aa957600080fd5b506104b2610ab8366004614046565b61263f565b348015610ac957600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661068e565b348015610af457600080fd5b5061052a610b03366004613c4f565b60056020526000908152604090205481565b348015610b2157600080fd5b506104b2610b30366004613b9d565b6126d0565b348015610b4157600080fd5b506104b2610b50366004614092565b612915565b348015610b6157600080fd5b50601554610b7f906c01000000000000000000000000900460ff1681565b60405160ff9091168152602001610501565b348015610b9d57600080fd5b506104b2610bac366004613b9d565b60009182526009602052604090912055565b348015610bca57600080fd5b506105b9610bd9366004613c4f565b600e6020526000908152604090205461ffff1681565b348015610bfb57600080fd5b506104f4610c0a366004613c4f565b612a5e565b348015610c1b57600080fd5b50610558610c2a366004613c4f565b612ac0565b348015610c3b57600080fd5b5061052a7f000000000000000000000000000000000000000000000000000000000000000081565b348015610c6f57600080fd5b506104b2610c7e366004613d54565b612b6c565b348015610c8f57600080fd5b506104b2610c9e366004613c4f565b601455565b348015610caf57600080fd5b506104b2610cbe3660046140c2565b612bd5565b348015610ccf57600080fd5b506104b2610cde366004613b9d565b60009182526007602052604090912055565b348015610cfc57600080fd5b506104b2610d0b366004613c4f565b612c80565b348015610d1c57600080fd5b50610558612d06565b348015610d3157600080fd5b506104b2610d40366004613dc5565b612d13565b348015610d5157600080fd5b506104b2610d603660046140e7565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610dab57600080fd5b506019546108f590610100900460ff1681565b348015610dca57600080fd5b506104b2610dd9366004613c4f565b612dad565b348015610dea57600080fd5b50610558612e45565b348015610dff57600080fd5b506104b2610e0e366004614104565b6015805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b348015610e5e57600080fd5b506105ec610e6d366004613b9d565b612e52565b348015610e7e57600080fd5b506019546108f59060ff1681565b348015610e9857600080fd5b506104b2610ea7366004613dc5565b612ebb565b348015610eb857600080fd5b5061052a610ec7366004613b9d565b612f86565b348015610ed857600080fd5b5061052a610ee7366004613c4f565b60096020526000908152604090205481565b348015610f0557600080fd5b506104b2610f14366004614029565b601980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b348015610f5157600080fd5b5061052a60145481565b348015610f6757600080fd5b506104b2610f76366004613d37565b612fb7565b348015610f8757600080fd5b5061052a610f96366004613c4f565b60066020526000908152604090205481565b348015610fb457600080fd5b5061052a610fc3366004613c4f565b60026020526000908152604090205481565b8051610fe89060169060208401906137b0565b5050565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b39216906110d2908c1688614150565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af1158015611150573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611174919061418f565b5060008860ff1667ffffffffffffffff81111561119357611193613890565b6040519080825280602002602001820160405280156111bc578160200160208202803683370190505b50905060005b8960ff168160ff1610156113655760006111db84612fcb565b90508860ff16600103611313576040517f0d4a4fb1000000000000000000000000000000000000000000000000000000008152600481018290526000903090630d4a4fb190602401600060405180830381865afa158015611240573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261128691908101906141f9565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906112df908590859060040161422e565b600060405180830381600087803b1580156112f957600080fd5b505af115801561130d573d6000803e3d6000fd5b50505050505b80838360ff168151811061132957611329614247565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061135d81614276565b9150506111c2565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711816040516113959190613bbf565b60405180910390a1505050505050505050565b606060006113b6600f613099565b90508084106113f1576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82600003611406576114038482614295565b92505b60008367ffffffffffffffff81111561142157611421613890565b60405190808252806020026020018201604052801561144a578160200160208202803683370190505b50905060005b8481101561149c5761146d61146582886142a8565b600f906130a3565b82828151811061147f5761147f614247565b602090810291909101015280611494816142bb565b915050611450565b509150505b92915050565b600d60205282600052604060002060205281600052604060002081815481106114cf57600080fd5b9060005260206000200160009250925050505481565b606060006040518060c001604052803073ffffffffffffffffffffffffffffffffffffffff168152602001600160ff1681526020017f000000000000000000000000000000000000000000000000000000000000000081526020018460405160200161155391815260200190565b60405160208183030381529060405261156b906142f3565b81526020016000801b81526020016000801b8152509050806040516020016115eb9190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b604051602081830303815290604052915050919050565b600061160c6130af565b604051308152909150819083907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a35050565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156116ae57602002820191906000526020600020905b81548152602001906001019080831161169a575b505050505090506116c0818251613151565b92509250505b9250929050565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611763573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117879190614343565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801561182a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061184e9190614371565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b6000838152600c6020908152604080832080548251818502810185019093528083526118f8938301828280156118ec57602002820191906000526020600020905b8154815260200190600101908083116118d8575b505050505084846131d6565b90505b9392505050565b8060005b818160ff161015611996573063c8048022858560ff851681811061192c5761192c614247565b905060200201356040518263ffffffff1660e01b815260040161195191815260200190565b600060405180830381600087803b15801561196b57600080fd5b505af115801561197f573d6000803e3d6000fd5b50505050808061198e90614276565b915050611906565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff16116119fd576000858152600d6020908152604080832061ffff851684529091529020546119e990836142a8565b9150806119f58161438e565b9150506119b2565b509392505050565b6000606060005a90506000611a186130af565b90507f0000000000000000000000000000000000000000000000000000000000000000611a4860c08801886143af565b6000818110611a5957611a59614247565b9050602002013503611cb5576000611a7460c08801886143af565b6001818110611a8557611a85614247565b90506020020135604051602001611a9e91815260200190565b6040516020818303038152906040529050600081806020019051810190611ac59190614417565b90506000611ad660c08a018a6143af565b6002818110611ae757611ae7614247565b90506020020135604051602001611b0091815260200190565b6040516020818303038152906040529050600081806020019051810190611b279190614417565b6000848152600860205260409020549091505b805a611b469089614295565b611b5290613a986142a8565b1015611b935781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055611b3a565b601954610100900460ff1615611c3157601760166018848786604051602001611bc6929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e000000000000000000000000000000000000000000000000000000008252611c289594939291600401614518565b60405180910390fd5b604080516001808252818301909252600091816020015b6060815260200190600190039081611c4857505060408051602081018890528082018690528151808203830181526060820190925291925090600190611c9490849084906080016145db565b6040516020818303038152906040529a509a505050505050505050506116c6565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401611c28565b60005a9050600080611d5184860186613f07565b9150915060008082806020019051810190611d6c919061466f565b6000828152600560209081526040808320546004909252822054939550919350909190611d976130af565b905082600003611db7576000858152600560205260409020819055611efb565b6000611dc38583614295565b6000878152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611e3557602002820191906000526020600020905b815481526020019060010190808311611e21575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff16815103611eb05781611e728161438e565b60008a8152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055898352600c8252928220805493840181558252902001555b600085815260066020526040812054611f159060016142a8565b6000878152600660209081526040808320849055600490915290208390559050611f3f86836126d0565b60195460ff1615611f8257604051308152829087907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a35b611f8d868a84612b6c565b5050505050505050505050565b6000828152600d6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015611ffc57602002820191906000526020600020905b815481526020019060010190808311611fe8575b5050505050905092915050565b60006060600084846040516020016120229291906145db565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b8160005b818110156120fa5730635f17e61686868481811061208157612081614247565b90506020020135856040518363ffffffff1660e01b81526004016120b592919091825263ffffffff16602082015260400190565b600060405180830381600087803b1580156120cf57600080fd5b505af11580156120e3573d6000803e3d6000fd5b5050505080806120f2906142bb565b915050612061565b5050505050565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600060048201819052602482018190529173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612178573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526121be9190810190614693565b805190915060006121cd6130af565b905060005b828110156119965760008482815181106121ee576121ee614247565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa15801561226e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122929190614724565b90508060ff166001036122d757604051308152849083907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a35b505080806122e4906142bb565b9150506121d2565b6122f4613335565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612363573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123879190614417565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156123ff573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fe8919061418f565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c909152812061245b91613806565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff16116124b7576000848152600d6020908152604080832061ffff8516845290915281206124a591613806565b806124af8161438e565b915050612470565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b600081815260056020526040812054810361250d57506001919050565b6000828152600360209081526040808320546004909252909120546125306130af565b61253a9190614295565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146125c3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401611c28565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b59061269990869086908690600401614741565b600060405180830381600087803b1580156126b357600080fd5b505af11580156126c7573d6000803e3d6000fd5b50505050505050565b6014546000838152600260205260409020546126ec9083614295565b1115610fe8576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612762573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526127a891908101906147c3565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa15801561281d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061284191906148e2565b6015549091506128659082906c01000000000000000000000000900460ff16614150565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611996576015546128a89085906bffffffffffffffffffffffff16612915565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af115801561299d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129c1919061418f565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b158015612a4257600080fd5b505af1158015612a56573d6000803e3d6000fd5b505050505050565b6000818152600c6020908152604091829020805483518184028101840190945280845260609392830182828015612ab457602002820191906000526020600020905b815481526020019060010190808311612aa0575b50505050509050919050565b60168181548110612ad057600080fd5b906000526020600020016000915090508054612aeb90614430565b80601f0160208091040260200160405190810160405280929190818152602001828054612b1790614430565b8015612b645780601f10612b3957610100808354040283529160200191612b64565b820191906000526020600020905b815481529060010190602001808311612b4757829003601f168201915b505050505081565b6000838152600760205260409020545b805a612b889085614295565b612b94906127106142a8565b10156119965781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612b7c565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b158015612c4d57600080fd5b505af1158015612c61573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b158015612cf257600080fd5b505af11580156120fa573d6000803e3d6000fd5b60178054612aeb90614430565b8060005b818163ffffffff161015611996573063af953a4a858563ffffffff8516818110612d4357612d43614247565b905060200201356040518263ffffffff1660e01b8152600401612d6891815260200190565b600060405180830381600087803b158015612d8257600080fd5b505af1158015612d96573d6000803e3d6000fd5b505050508080612da5906148ff565b915050612d17565b6013546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b158015612e1957600080fd5b505af1158015612e2d573d6000803e3d6000fd5b50505050610fe881600f6133b890919063ffffffff16565b60188054612aeb90614430565b6000828152600c60209081526040808320805482518185028101850190935280835284938493929190830182828015612eaa57602002820191906000526020600020905b815481526020019060010190808311612e96575b505050505090506116c08185613151565b8060005b81811015611996576000848483818110612edb57612edb614247565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc8283604051602001612f1491815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401612f4092919061422e565b600060405180830381600087803b158015612f5a57600080fd5b505af1158015612f6e573d6000803e3d6000fd5b50505050508080612f7e906142bb565b915050612ebf565b600c6020528160005260406000208181548110612fa257600080fd5b90600052602060002001600091509150505481565b612fbf613335565b612fc8816133c4565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613026908690600401614918565b6020604051808303816000875af1158015613045573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130699190614417565b9050613076600f826134b9565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b60006114a1825490565b60006118fb83836134c5565b60007f00000000000000000000000000000000000000000000000000000000000000001561314c57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613123573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131479190614417565b905090565b504390565b8151600090819081908415806131675750808510155b15613170578094505b60008092505b858310156131cc5786600161318b8585614295565b6131959190614295565b815181106131a5576131a5614247565b6020026020010151816131b891906142a8565b9050826131c4816142bb565b935050613176565b9694955050505050565b825160009081908315806131ea5750808410155b156131f3578093505b60008467ffffffffffffffff81111561320e5761320e613890565b604051908082528060200260200182016040528015613237578160200160208202803683370190505b509050600092505b848310156132a5578660016132548585614295565b61325e9190614295565b8151811061326e5761326e614247565b602002602001015181848151811061328857613288614247565b60209081029190910101528261329d816142bb565b93505061323f565b6132be816000600184516132b99190614295565b6134ef565b856064036132f75780600182516132d59190614295565b815181106132e5576132e5614247565b602002602001015193505050506118fb565b8060648251886133079190614a6a565b6133119190614ad6565b8151811061332157613321614247565b602002602001015193505050509392505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146133b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611c28565b565b60006118fb8383613667565b3373ffffffffffffffffffffffffffffffffffffffff821603613443576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611c28565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006118fb8383613761565b60008260000182815481106134dc576134dc614247565b9060005260206000200154905092915050565b81818082036134ff575050505050565b600085600261350e8787614aea565b6135189190614b0a565b6135229087614b72565b8151811061353257613532614247565b602002602001015190505b818313613641575b8086848151811061355857613558614247565b60200260200101511015613578578261357081614b9a565b935050613545565b85828151811061358a5761358a614247565b60200260200101518110156135ab57816135a381614bcb565b925050613578565b81831361363c578582815181106135c4576135c4614247565b60200260200101518684815181106135de576135de614247565b60200260200101518785815181106135f8576135f8614247565b6020026020010188858151811061361157613611614247565b6020908102919091010191909152528261362a81614b9a565b935050818061363890614bcb565b9250505b61353d565b81851215613654576136548686846134ef565b83831215612a5657612a568684866134ef565b6000818152600183016020526040812054801561375057600061368b600183614295565b855490915060009061369f90600190614295565b90508181146137045760008660000182815481106136bf576136bf614247565b90600052602060002001549050808760000184815481106136e2576136e2614247565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061371557613715614c22565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506114a1565b60009150506114a1565b5092915050565b60008181526001830160205260408120546137a8575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556114a1565b5060006114a1565b8280548282559060005260206000209081019282156137f6579160200282015b828111156137f657825182906137e69082614c9c565b50916020019190600101906137d0565b50613802929150613824565b5090565b5080546000825590600052602060002090810190612fc89190613841565b808211156138025760006138388282613856565b50600101613824565b5b808211156138025760008155600101613842565b50805461386290614430565b6000825580601f10613872575050565b601f016020900490600052602060002090810190612fc89190613841565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff811182821017156138e3576138e3613890565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561393057613930613890565b604052919050565b600067ffffffffffffffff82111561395257613952613890565b5060051b60200190565b600067ffffffffffffffff82111561397657613976613890565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006139b56139b08461395c565b6138e9565b90508281528383830111156139c957600080fd5b828260208301376000602084830101529392505050565b600060208083850312156139f357600080fd5b823567ffffffffffffffff80821115613a0b57600080fd5b818501915085601f830112613a1f57600080fd5b8135613a2d6139b082613938565b81815260059190911b83018401908481019088831115613a4c57600080fd5b8585015b83811015613a9957803585811115613a685760008081fd5b8601603f81018b13613a7a5760008081fd5b613a8b8b89830135604084016139a2565b845250918601918601613a50565b5098975050505050505050565b60ff81168114612fc857600080fd5b63ffffffff81168114612fc857600080fd5b600082601f830112613ad857600080fd5b6118fb838335602085016139a2565b6bffffffffffffffffffffffff81168114612fc857600080fd5b600080600080600080600060e0888a031215613b1c57600080fd5b8735613b2781613aa6565b96506020880135613b3781613ab5565b95506040880135613b4781613aa6565b9450606088013567ffffffffffffffff811115613b6357600080fd5b613b6f8a828b01613ac7565b9450506080880135613b8081613ae7565b9699959850939692959460a0840135945060c09093013592915050565b60008060408385031215613bb057600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015613bf757835183529284019291840191600101613bdb565b50909695505050505050565b803561ffff81168114613c1557600080fd5b919050565b600080600060608486031215613c2f57600080fd5b83359250613c3f60208501613c03565b9150604084013590509250925092565b600060208284031215613c6157600080fd5b5035919050565b60005b83811015613c83578181015183820152602001613c6b565b50506000910152565b60008151808452613ca4816020860160208601613c68565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006118fb6020830184613c8c565b60008060408385031215613cfc57600080fd5b82359150613d0c60208401613c03565b90509250929050565b73ffffffffffffffffffffffffffffffffffffffff81168114612fc857600080fd5b600060208284031215613d4957600080fd5b81356118fb81613d15565b600080600060608486031215613d6957600080fd5b505081359360208301359350604090920135919050565b60008083601f840112613d9257600080fd5b50813567ffffffffffffffff811115613daa57600080fd5b6020830191508360208260051b85010111156116c657600080fd5b60008060208385031215613dd857600080fd5b823567ffffffffffffffff811115613def57600080fd5b613dfb85828601613d80565b90969095509350505050565b60008060408385031215613e1a57600080fd5b823567ffffffffffffffff80821115613e3257600080fd5b908401906101008287031215613e4757600080fd5b90925060208401359080821115613e5d57600080fd5b50613e6a85828601613ac7565b9150509250929050565b82151581526040602082015260006118f86040830184613c8c565b60008083601f840112613ea157600080fd5b50813567ffffffffffffffff811115613eb957600080fd5b6020830191508360208285010111156116c657600080fd5b60008060208385031215613ee457600080fd5b823567ffffffffffffffff811115613efb57600080fd5b613dfb85828601613e8f565b60008060408385031215613f1a57600080fd5b823567ffffffffffffffff80821115613f3257600080fd5b818501915085601f830112613f4657600080fd5b81356020613f566139b083613938565b82815260059290921b84018101918181019089841115613f7557600080fd5b8286015b84811015613fad57803586811115613f915760008081fd5b613f9f8c86838b0101613ac7565b845250918301918301613f79565b5096505086013592505080821115613e5d57600080fd5b600080600060408486031215613fd957600080fd5b833567ffffffffffffffff811115613ff057600080fd5b613ffc86828701613d80565b909450925050602084013561401081613ab5565b809150509250925092565b8015158114612fc857600080fd5b60006020828403121561403b57600080fd5b81356118fb8161401b565b60008060006040848603121561405b57600080fd5b83359250602084013567ffffffffffffffff81111561407957600080fd5b61408586828701613e8f565b9497909650939450505050565b600080604083850312156140a557600080fd5b8235915060208301356140b781613ae7565b809150509250929050565b600080604083850312156140d557600080fd5b8235915060208301356140b781613ab5565b6000602082840312156140f957600080fd5b81356118fb81613ae7565b60006020828403121561411657600080fd5b81356118fb81613aa6565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff8083168185168183048111821515161561417b5761417b614121565b02949350505050565b8051613c158161401b565b6000602082840312156141a157600080fd5b81516118fb8161401b565b600082601f8301126141bd57600080fd5b81516141cb6139b08261395c565b8181528460208386010111156141e057600080fd5b6141f1826020830160208701613c68565b949350505050565b60006020828403121561420b57600080fd5b815167ffffffffffffffff81111561422257600080fd5b6141f1848285016141ac565b8281526040602082015260006118f86040830184613c8c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff810361428c5761428c614121565b60010192915050565b818103818111156114a1576114a1614121565b808201808211156114a1576114a1614121565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036142ec576142ec614121565b5060010190565b80516020808301519190811015614332577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b8051613c1581613d15565b6000806040838503121561435657600080fd5b825161436181613d15565b6020939093015192949293505050565b60006020828403121561438357600080fd5b81516118fb81613d15565b600061ffff8083168181036143a5576143a5614121565b6001019392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126143e457600080fd5b83018035915067ffffffffffffffff8211156143ff57600080fd5b6020019150600581901b36038213156116c657600080fd5b60006020828403121561442957600080fd5b5051919050565b600181811c9082168061444457607f821691505b602082108103614332577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000815461448a81614430565b8085526020600183811680156144a757600181146144df5761450d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955061450d565b866000528260002060005b858110156145055781548a82018601529083019084016144ea565b890184019650505b505050505092915050565b60a08152600061452b60a083018861447d565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b8381101561459d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087840301855261458b838361447d565b94860194925060019182019101614552565b505086810360408801526145b1818b61447d565b94505050505084606084015282810360808401526145cf8185613c8c565b98975050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015614650577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261463e868351613c8c565b95509382019390820190600101614604565b5050858403818701525050506146668185613c8c565b95945050505050565b6000806040838503121561468257600080fd5b505080516020909101519092909150565b600060208083850312156146a657600080fd5b825167ffffffffffffffff8111156146bd57600080fd5b8301601f810185136146ce57600080fd5b80516146dc6139b082613938565b81815260059190911b820183019083810190878311156146fb57600080fd5b928401925b8284101561471957835182529284019290840190614700565b979650505050505050565b60006020828403121561473657600080fd5b81516118fb81613aa6565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b8051613c1581613ab5565b8051613c1581613ae7565b805167ffffffffffffffff81168114613c1557600080fd5b6000602082840312156147d557600080fd5b815167ffffffffffffffff808211156147ed57600080fd5b90830190610140828603121561480257600080fd5b61480a6138bf565b61481383614338565b815261482160208401614795565b602082015260408301518281111561483857600080fd5b614844878286016141ac565b604083015250614856606084016147a0565b606082015261486760808401614338565b608082015261487860a084016147ab565b60a082015261488960c08401614795565b60c082015261489a60e084016147a0565b60e08201526101006148ad818501614184565b9082015261012083810151838111156148c557600080fd5b6148d1888287016141ac565b918301919091525095945050505050565b6000602082840312156148f457600080fd5b81516118fb81613ae7565b600063ffffffff8083168181036143a5576143a5614121565b6020815260008251610140806020850152614937610160850183613c8c565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808685030160408701526149738483613c8c565b93506040870151915061499e606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526149ff8483613c8c565b935060e08701519150610100818786030181880152614a1e8584613c8c565b945080880151925050610120818786030181880152614a3d8584613c8c565b94508088015192505050614a60828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614aa257614aa2614121565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082614ae557614ae5614aa7565b500490565b818103600083128015838313168383128216171561375a5761375a614121565b600082614b1957614b19614aa7565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615614b6d57614b6d614121565b500590565b8082018281126000831280158216821582161715614b9257614b92614121565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036142ec576142ec614121565b60007f80000000000000000000000000000000000000000000000000000000000000008203614bfc57614bfc614121565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b601f821115614c9757600081815260208120601f850160051c81016020861015614c785750805b601f850160051c820191505b81811015612a5657828155600101614c84565b505050565b815167ffffffffffffffff811115614cb657614cb6613890565b614cca81614cc48454614430565b84614c51565b602080601f831160018114614d1d5760008415614ce75750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555612a56565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015614d6a57888601518255948401946001909101908401614d4b565b5085821015614da657878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var VerifiableLoadLogTriggerUpkeepABI = VerifiableLoadLogTriggerUpkeepMetaData.ABI

var VerifiableLoadLogTriggerUpkeepBin = VerifiableLoadLogTriggerUpkeepMetaData.Bin

func DeployVerifiableLoadLogTriggerUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _registrar common.Address, _useArb bool, _autoLog bool, _useMercury bool) (common.Address, *types.Transaction, *VerifiableLoadLogTriggerUpkeep, error) {
	parsed, err := VerifiableLoadLogTriggerUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadLogTriggerUpkeepBin), backend, _registrar, _useArb, _autoLog, _useMercury)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) AutoLog(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "autoLog")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) AutoLog() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AutoLog(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) AutoLog() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AutoLog(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) UseMercury(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "useMercury")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UseMercury() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseMercury(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) UseMercury() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseMercury(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SendLog(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "sendLog", upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SendLog(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SendLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SendLog(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SendLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetAutoLog(opts *bind.TransactOpts, _autoLog bool) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setAutoLog", _autoLog)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetAutoLog(_autoLog bool) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetAutoLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _autoLog)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetAutoLog(_autoLog bool) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetAutoLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _autoLog)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetUseMercury(opts *bind.TransactOpts, _useMercury bool) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setUseMercury", _useMercury)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetUseMercury(_useMercury bool) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUseMercury(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _useMercury)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetUseMercury(_useMercury bool) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUseMercury(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _useMercury)
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

	AutoLog(opts *bind.CallOpts) (bool, error)

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

	UseMercury(opts *bind.CallOpts) (bool, error)

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

	SendLog(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	SetAddLinkAmount(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	SetAutoLog(opts *bind.TransactOpts, _autoLog bool) (*types.Transaction, error)

	SetCheckGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error)

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetMinBalanceThresholdMultiplier(opts *bind.TransactOpts, newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error)

	SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetPerformGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepTopUpCheckInterval(opts *bind.TransactOpts, newInterval *big.Int) (*types.Transaction, error)

	SetUseMercury(opts *bind.TransactOpts, _useMercury bool) (*types.Transaction, error)

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
