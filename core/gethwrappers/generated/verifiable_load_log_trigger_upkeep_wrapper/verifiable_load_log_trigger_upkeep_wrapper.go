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
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460c0526101a0604052604261012081815260e09182919062005486610140398152602001604051806080016040528060428152602001620054c860429139905262000099906016906002620003c4565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000c9908262000540565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152601890620000fb908262000540565b503480156200010957600080fd5b506040516200550a3803806200550a8339810160408190526200012c9162000638565b83833380600081620001855760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001b857620001b88162000319565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa15801562000215573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200023b919062000697565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002a1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002c79190620006c8565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560a052506019805461ffff191692151561ff001916929092176101009115159190910217905550620006ef9050565b336001600160a01b03821603620003735760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200017c565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156200040f579160200282015b828111156200040f5782518290620003fe908262000540565b5091602001919060010190620003e5565b506200041d92915062000421565b5090565b808211156200041d57600062000438828262000442565b5060010162000421565b5080546200045090620004b1565b6000825580601f1062000461575050565b601f01602090049060005260206000209081019062000481919062000484565b50565b5b808211156200041d576000815560010162000485565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004c657607f821691505b602082108103620004e757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200053b57600081815260208120601f850160051c81016020861015620005165750805b601f850160051c820191505b81811015620005375782815560010162000522565b5050505b505050565b81516001600160401b038111156200055c576200055c6200049b565b62000574816200056d8454620004b1565b84620004ed565b602080601f831160018114620005ac5760008415620005935750858301515b600019600386901b1c1916600185901b17855562000537565b600085815260208120601f198616915b82811015620005dd57888601518255948401946001909101908401620005bc565b5085821015620005fc5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200048157600080fd5b805180151581146200063357600080fd5b919050565b600080600080608085870312156200064f57600080fd5b84516200065c816200060c565b93506200066c6020860162000622565b92506200067c6040860162000622565b91506200068c6060860162000622565b905092959194509250565b60008060408385031215620006ab57600080fd5b8251620006b8816200060c565b6020939093015192949293505050565b600060208284031215620006db57600080fd5b8151620006e8816200060c565b9392505050565b60805160a05160c051614d4b6200073b600039600081816105970152611db70152600081816108d3015261302c015260008181610c410152818161151c0152611a1c0152614d4b6000f3fe60806040526004361061044e5760003560e01c806379ba509711610243578063af953a4a11610143578063d90c4a68116100bb578063e15e1b341161008a578063f2fde38b1161006f578063f2fde38b14610f5b578063fba7ffa314610f7b578063fcdc1f6314610fa857600080fd5b8063e15e1b3414610ef9578063e455308314610f4557600080fd5b8063d90c4a6814610e72578063daee1aeb14610e8c578063dbef701e14610eac578063e0114adb14610ecc57600080fd5b8063c41c815b11610112578063c98f10b0116100f7578063c98f10b014610dde578063d355852814610df3578063d6051a7214610e5257600080fd5b8063c41c815b14610d9f578063c804802214610dbe57600080fd5b8063af953a4a14610cf0578063afb28d1f14610d10578063becde0e114610d25578063c357f1f314610d4557600080fd5b80639b429354116101d6578063a6548248116101a5578063a6c60d891161018a578063a6c60d8914610c83578063a72aa27e14610ca3578063a79c404314610cc357600080fd5b8063a654824814610c2f578063a6b5947514610c6357600080fd5b80639b42935414610b915780639b51fb0d14610bbe5780639d385eaa14610bef5780639d6f1cc714610c0f57600080fd5b80638fcb3fba116102125780638fcb3fba14610ae8578063924ca57814610b15578063948108f714610b355780639ac542eb14610b5557600080fd5b806379ba509714610a5b5780637b10399914610a705780637e7a46dc14610a9d5780638da5cb5b14610abd57600080fd5b806346e7a63e1161034e578063636092e8116102e15780636b2120c6116102b057806373644cce1161029557806373644cce146109e15780637672130314610a0e578063776898c814610a3b57600080fd5b80636b2120c61461095f5780637145f11b146109b157600080fd5b8063636092e81461087f578063642f6cef146108c157806369cdbadb1461090557806369e9b7731461093257600080fd5b8063597109921161031d57806359710992146108085780635d4ee7f31461081d5780635f17e6161461083257806360457ff51461085257600080fd5b806346e7a63e1461076e5780634b56a42e1461079b57806351c98be3146107bb57806357970e93146107db57600080fd5b806320e3dbd4116103e1578063328ffd11116103b057806340691db41161039557806340691db4146107005780634585e33b1461072e57806345d2ec171461074e57600080fd5b8063328ffd11146106b35780633ebe8d6c146106e057600080fd5b806320e3dbd41461060157806328c4b57b146106215780632a9032d3146106415780632b20e3971461066157600080fd5b80630d4a4fb11161041d5780630d4a4fb1146105385780630e577d421461056557806312c5502714610585578063206c32e8146105cc57600080fd5b806305e251311461049257806306c1cc00146104b457806306e3b632146104d4578063077ac6211461050a57600080fd5b3661048d57604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561049e57600080fd5b506104b26104ad366004613959565b610fd5565b005b3480156104c057600080fd5b506104b26104cf366004613a7a565b610fec565b3480156104e057600080fd5b506104f46104ef366004613b16565b6113a8565b6040516105019190613b38565b60405180910390f35b34801561051657600080fd5b5061052a610525366004613b93565b6114a7565b604051908152602001610501565b34801561054457600080fd5b50610558610553366004613bc8565b6114e5565b6040516105019190613c4f565b34801561057157600080fd5b506104b2610580366004613bc8565b611602565b34801561059157600080fd5b506105b97f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610501565b3480156105d857600080fd5b506105ec6105e7366004613c62565b61164a565b60408051928352602083019190915201610501565b34801561060d57600080fd5b506104b261061c366004613cb0565b6116cd565b34801561062d57600080fd5b5061052a61063c366004613ccd565b611897565b34801561064d57600080fd5b506104b261065c366004613d3e565b611902565b34801561066d57600080fd5b5060115461068e9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610501565b3480156106bf57600080fd5b5061052a6106ce366004613bc8565b60036020526000908152604090205481565b3480156106ec57600080fd5b5061052a6106fb366004613bc8565b61199c565b34801561070c57600080fd5b5061072061071b366004613d80565b611a05565b604051610501929190613ded565b34801561073a57600080fd5b506104b2610749366004613e4a565b611cb6565b34801561075a57600080fd5b506104f4610769366004613c62565b611f13565b34801561077a57600080fd5b5061052a610789366004613bc8565b600a6020526000908152604090205481565b3480156107a757600080fd5b506107206107b6366004613e80565b611f82565b3480156107c757600080fd5b506104b26107d6366004613f3d565b611fd6565b3480156107e757600080fd5b5060125461068e9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561081457600080fd5b506104b261207a565b34801561082957600080fd5b506104b2612265565b34801561083e57600080fd5b506104b261084d366004613b16565b61239c565b34801561085e57600080fd5b5061052a61086d366004613bc8565b60076020526000908152604090205481565b34801561088b57600080fd5b506015546108a4906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff9091168152602001610501565b3480156108cd57600080fd5b506108f57f000000000000000000000000000000000000000000000000000000000000000081565b6040519015158152602001610501565b34801561091157600080fd5b5061052a610920366004613bc8565b60086020526000908152604090205481565b34801561093e57600080fd5b506104b261094d366004613b16565b60009182526008602052604090912055565b34801561096b57600080fd5b506104b261097a366004613fa2565b60198054911515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff909216919091179055565b3480156109bd57600080fd5b506108f56109cc366004613bc8565b600b6020526000908152604090205460ff1681565b3480156109ed57600080fd5b5061052a6109fc366004613bc8565b6000908152600c602052604090205490565b348015610a1a57600080fd5b5061052a610a29366004613bc8565b60046020526000908152604090205481565b348015610a4757600080fd5b506108f5610a56366004613bc8565b612469565b348015610a6757600080fd5b506104b26124bb565b348015610a7c57600080fd5b5060135461068e9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610aa957600080fd5b506104b2610ab8366004613fbf565b6125b8565b348015610ac957600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661068e565b348015610af457600080fd5b5061052a610b03366004613bc8565b60056020526000908152604090205481565b348015610b2157600080fd5b506104b2610b30366004613b16565b612649565b348015610b4157600080fd5b506104b2610b5036600461400b565b61288e565b348015610b6157600080fd5b50601554610b7f906c01000000000000000000000000900460ff1681565b60405160ff9091168152602001610501565b348015610b9d57600080fd5b506104b2610bac366004613b16565b60009182526009602052604090912055565b348015610bca57600080fd5b506105b9610bd9366004613bc8565b600e6020526000908152604090205461ffff1681565b348015610bfb57600080fd5b506104f4610c0a366004613bc8565b6129d7565b348015610c1b57600080fd5b50610558610c2a366004613bc8565b612a39565b348015610c3b57600080fd5b5061052a7f000000000000000000000000000000000000000000000000000000000000000081565b348015610c6f57600080fd5b506104b2610c7e366004613ccd565b612ae5565b348015610c8f57600080fd5b506104b2610c9e366004613bc8565b601455565b348015610caf57600080fd5b506104b2610cbe36600461403b565b612b4e565b348015610ccf57600080fd5b506104b2610cde366004613b16565b60009182526007602052604090912055565b348015610cfc57600080fd5b506104b2610d0b366004613bc8565b612bf9565b348015610d1c57600080fd5b50610558612c7f565b348015610d3157600080fd5b506104b2610d40366004613d3e565b612c8c565b348015610d5157600080fd5b506104b2610d60366004614060565b601580547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610dab57600080fd5b506019546108f590610100900460ff1681565b348015610dca57600080fd5b506104b2610dd9366004613bc8565b612d26565b348015610dea57600080fd5b50610558612dbe565b348015610dff57600080fd5b506104b2610e0e36600461407d565b6015805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b348015610e5e57600080fd5b506105ec610e6d366004613b16565b612dcb565b348015610e7e57600080fd5b506019546108f59060ff1681565b348015610e9857600080fd5b506104b2610ea7366004613d3e565b612e34565b348015610eb857600080fd5b5061052a610ec7366004613b16565b612eff565b348015610ed857600080fd5b5061052a610ee7366004613bc8565b60096020526000908152604090205481565b348015610f0557600080fd5b506104b2610f14366004613fa2565b601980547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b348015610f5157600080fd5b5061052a60145481565b348015610f6757600080fd5b506104b2610f76366004613cb0565b612f30565b348015610f8757600080fd5b5061052a610f96366004613bc8565b60066020526000908152604090205481565b348015610fb457600080fd5b5061052a610fc3366004613bc8565b60026020526000908152604090205481565b8051610fe8906016906020840190613729565b5050565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b39216906110d2908c16886140c9565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af1158015611150573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111749190614108565b5060008860ff1667ffffffffffffffff81111561119357611193613809565b6040519080825280602002602001820160405280156111bc578160200160208202803683370190505b50905060005b8960ff168160ff1610156113655760006111db84612f44565b90508860ff16600103611313576040517f0d4a4fb1000000000000000000000000000000000000000000000000000000008152600481018290526000903090630d4a4fb190602401600060405180830381865afa158015611240573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526112869190810190614172565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906112df90859085906004016141a7565b600060405180830381600087803b1580156112f957600080fd5b505af115801561130d573d6000803e3d6000fd5b50505050505b80838360ff1681518110611329576113296141c0565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061135d816141ef565b9150506111c2565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711816040516113959190613b38565b60405180910390a1505050505050505050565b606060006113b6600f613012565b90508084106113f1576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8260000361140657611403848261420e565b92505b60008367ffffffffffffffff81111561142157611421613809565b60405190808252806020026020018201604052801561144a578160200160208202803683370190505b50905060005b8481101561149c5761146d6114658288614221565b600f9061301c565b82828151811061147f5761147f6141c0565b60209081029190910101528061149481614234565b915050611450565b509150505b92915050565b600d60205282600052604060002060205281600052604060002081815481106114cf57600080fd5b9060005260206000200160009250925050505481565b606060006040518060c001604052803073ffffffffffffffffffffffffffffffffffffffff168152602001600160ff1681526020017f000000000000000000000000000000000000000000000000000000000000000081526020018460405160200161155391815260200190565b60405160208183030381529060405261156b9061426c565b81526020016000801b81526020016000801b8152509050806040516020016115eb9190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b604051602081830303815290604052915050919050565b600061160c613028565b604051308152909150819083907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a35050565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156116ae57602002820191906000526020600020905b81548152602001906001019080831161169a575b505050505090506116c08182516130ca565b92509250505b9250929050565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611763573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061178791906142bc565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801561182a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061184e91906142ea565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b6000838152600c6020908152604080832080548251818502810185019093528083526118f8938301828280156118ec57602002820191906000526020600020905b8154815260200190600101908083116118d8575b5050505050848461314f565b90505b9392505050565b8060005b818160ff161015611996573063c8048022858560ff851681811061192c5761192c6141c0565b905060200201356040518263ffffffff1660e01b815260040161195191815260200190565b600060405180830381600087803b15801561196b57600080fd5b505af115801561197f573d6000803e3d6000fd5b50505050808061198e906141ef565b915050611906565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff16116119fd576000858152600d6020908152604080832061ffff851684529091529020546119e99083614221565b9150806119f581614307565b9150506119b2565b509392505050565b6000606060005a90506000611a18613028565b90507f0000000000000000000000000000000000000000000000000000000000000000611a4860c0880188614328565b6000818110611a5957611a596141c0565b9050602002013503611c2e576000611a7460c0880188614328565b6001818110611a8557611a856141c0565b90506020020135604051602001611a9e91815260200190565b6040516020818303038152906040529050600081806020019051810190611ac59190614390565b90506000611ad660c08a018a614328565b6002818110611ae757611ae76141c0565b90506020020135604051602001611b0091815260200190565b6040516020818303038152906040529050600081806020019051810190611b279190614390565b6000848152600860205260409020549091505b805a611b46908961420e565b611b5290613a98614221565b1015611ba05781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905581611b98816143a9565b925050611b3a565b601760166018848786604051602001611bc3929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e000000000000000000000000000000000000000000000000000000008252611c2595949392916004016144c6565b60405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401611c25565b60005a9050600080611cca84860186613e80565b9150915060008082806020019051810190611ce59190614589565b6000828152600560209081526040808320546004909252822054939550919350909190611d10613028565b905082600003611d30576000858152600560205260409020819055611e74565b6000611d3c858361420e565b6000878152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611dae57602002820191906000526020600020905b815481526020019060010190808311611d9a575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff16815103611e295781611deb81614307565b60008a8152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055898352600c8252928220805493840181558252902001555b600085815260066020526040812054611e8e906001614221565b6000878152600660209081526040808320849055600490915290208390559050611eb88683612649565b60195460ff1615611efb57604051308152829087907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a35b611f06868a84612ae5565b5050505050505050505050565b6000828152600d6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015611f7557602002820191906000526020600020905b815481526020019060010190808311611f61575b5050505050905092915050565b6000606060008484604051602001611f9b9291906145ad565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b8160005b818110156120735730635f17e616868684818110611ffa57611ffa6141c0565b90506020020135856040518363ffffffff1660e01b815260040161202e92919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561204857600080fd5b505af115801561205c573d6000803e3d6000fd5b50505050808061206b90614234565b915050611fda565b5050505050565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600060048201819052602482018190529173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa1580156120f1573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526121379190810190614641565b80519091506000612146613028565b905060005b82811015611996576000848281518110612167576121676141c0565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa1580156121e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061220b91906146d2565b90508060ff1660010361225057604051308152849083907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a35b5050808061225d90614234565b91505061214b565b61226d6132ae565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156122dc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123009190614390565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af1158015612378573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fe89190614108565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c90915281206123d49161377f565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff1611612430576000848152600d6020908152604080832061ffff85168452909152812061241e9161377f565b8061242881614307565b9150506123e9565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b600081815260056020526040812054810361248657506001919050565b6000828152600360209081526040808320546004909252909120546124a9613028565b6124b3919061420e565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461253c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401611c25565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b590612612908690869086906004016146ef565b600060405180830381600087803b15801561262c57600080fd5b505af1158015612640573d6000803e3d6000fd5b50505050505050565b601454600083815260026020526040902054612665908361420e565b1115610fe8576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa1580156126db573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526127219190810190614771565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612796573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127ba9190614890565b6015549091506127de9082906c01000000000000000000000000900460ff166140c9565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611996576015546128219085906bffffffffffffffffffffffff1661288e565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612916573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061293a9190614108565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b1580156129bb57600080fd5b505af11580156129cf573d6000803e3d6000fd5b505050505050565b6000818152600c6020908152604091829020805483518184028101840190945280845260609392830182828015612a2d57602002820191906000526020600020905b815481526020019060010190808311612a19575b50505050509050919050565b60168181548110612a4957600080fd5b906000526020600020016000915090508054612a64906143de565b80601f0160208091040260200160405190810160405280929190818152602001828054612a90906143de565b8015612add5780601f10612ab257610100808354040283529160200191612add565b820191906000526020600020905b815481529060010190602001808311612ac057829003601f168201915b505050505081565b6000838152600760205260409020545b805a612b01908561420e565b612b0d90612710614221565b10156119965781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612af5565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b158015612bc657600080fd5b505af1158015612bda573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b158015612c6b57600080fd5b505af1158015612073573d6000803e3d6000fd5b60178054612a64906143de565b8060005b818163ffffffff161015611996573063af953a4a858563ffffffff8516818110612cbc57612cbc6141c0565b905060200201356040518263ffffffff1660e01b8152600401612ce191815260200190565b600060405180830381600087803b158015612cfb57600080fd5b505af1158015612d0f573d6000803e3d6000fd5b505050508080612d1e906148ad565b915050612c90565b6013546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b158015612d9257600080fd5b505af1158015612da6573d6000803e3d6000fd5b50505050610fe881600f61333190919063ffffffff16565b60188054612a64906143de565b6000828152600c60209081526040808320805482518185028101850190935280835284938493929190830182828015612e2357602002820191906000526020600020905b815481526020019060010190808311612e0f575b505050505090506116c081856130ca565b8060005b81811015611996576000848483818110612e5457612e546141c0565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc8283604051602001612e8d91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401612eb99291906141a7565b600060405180830381600087803b158015612ed357600080fd5b505af1158015612ee7573d6000803e3d6000fd5b50505050508080612ef790614234565b915050612e38565b600c6020528160005260406000208181548110612f1b57600080fd5b90600052602060002001600091509150505481565b612f386132ae565b612f418161333d565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190612f9f9086906004016148c6565b6020604051808303816000875af1158015612fbe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612fe29190614390565b9050612fef600f82613432565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b60006114a1825490565b60006118fb838361343e565b60007f0000000000000000000000000000000000000000000000000000000000000000156130c557606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561309c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130c09190614390565b905090565b504390565b8151600090819081908415806130e05750808510155b156130e9578094505b60008092505b8583101561314557866001613104858561420e565b61310e919061420e565b8151811061311e5761311e6141c0565b6020026020010151816131319190614221565b90508261313d81614234565b9350506130ef565b9694955050505050565b825160009081908315806131635750808410155b1561316c578093505b60008467ffffffffffffffff81111561318757613187613809565b6040519080825280602002602001820160405280156131b0578160200160208202803683370190505b509050600092505b8483101561321e578660016131cd858561420e565b6131d7919061420e565b815181106131e7576131e76141c0565b6020026020010151818481518110613201576132016141c0565b60209081029190910101528261321681614234565b9350506131b8565b61323781600060018451613232919061420e565b613468565b8560640361327057806001825161324e919061420e565b8151811061325e5761325e6141c0565b602002602001015193505050506118fb565b8060648251886132809190614a18565b61328a9190614a84565b8151811061329a5761329a6141c0565b602002602001015193505050509392505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461332f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401611c25565b565b60006118fb83836135e0565b3373ffffffffffffffffffffffffffffffffffffffff8216036133bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401611c25565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60006118fb83836136da565b6000826000018281548110613455576134556141c0565b9060005260206000200154905092915050565b8181808203613478575050505050565b60008560026134878787614a98565b6134919190614ab8565b61349b9087614b20565b815181106134ab576134ab6141c0565b602002602001015190505b8183136135ba575b808684815181106134d1576134d16141c0565b602002602001015110156134f157826134e981614b48565b9350506134be565b858281518110613503576135036141c0565b6020026020010151811015613524578161351c81614b79565b9250506134f1565b8183136135b55785828151811061353d5761353d6141c0565b6020026020010151868481518110613557576135576141c0565b6020026020010151878581518110613571576135716141c0565b6020026020010188858151811061358a5761358a6141c0565b602090810291909101019190915252826135a381614b48565b93505081806135b190614b79565b9250505b6134b6565b818512156135cd576135cd868684613468565b838312156129cf576129cf868486613468565b600081815260018301602052604081205480156136c957600061360460018361420e565b85549091506000906136189060019061420e565b905081811461367d576000866000018281548110613638576136386141c0565b906000526020600020015490508087600001848154811061365b5761365b6141c0565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061368e5761368e614baa565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506114a1565b60009150506114a1565b5092915050565b6000818152600183016020526040812054613721575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556114a1565b5060006114a1565b82805482825590600052602060002090810192821561376f579160200282015b8281111561376f578251829061375f9082614c24565b5091602001919060010190613749565b5061377b92915061379d565b5090565b5080546000825590600052602060002090810190612f4191906137ba565b8082111561377b5760006137b182826137cf565b5060010161379d565b5b8082111561377b57600081556001016137bb565b5080546137db906143de565b6000825580601f106137eb575050565b601f016020900490600052602060002090810190612f4191906137ba565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff8111828210171561385c5761385c613809565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156138a9576138a9613809565b604052919050565b600067ffffffffffffffff8211156138cb576138cb613809565b5060051b60200190565b600067ffffffffffffffff8211156138ef576138ef613809565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600061392e613929846138d5565b613862565b905082815283838301111561394257600080fd5b828260208301376000602084830101529392505050565b6000602080838503121561396c57600080fd5b823567ffffffffffffffff8082111561398457600080fd5b818501915085601f83011261399857600080fd5b81356139a6613929826138b1565b81815260059190911b830184019084810190888311156139c557600080fd5b8585015b83811015613a12578035858111156139e15760008081fd5b8601603f81018b136139f35760008081fd5b613a048b898301356040840161391b565b8452509186019186016139c9565b5098975050505050505050565b60ff81168114612f4157600080fd5b63ffffffff81168114612f4157600080fd5b600082601f830112613a5157600080fd5b6118fb8383356020850161391b565b6bffffffffffffffffffffffff81168114612f4157600080fd5b600080600080600080600060e0888a031215613a9557600080fd5b8735613aa081613a1f565b96506020880135613ab081613a2e565b95506040880135613ac081613a1f565b9450606088013567ffffffffffffffff811115613adc57600080fd5b613ae88a828b01613a40565b9450506080880135613af981613a60565b9699959850939692959460a0840135945060c09093013592915050565b60008060408385031215613b2957600080fd5b50508035926020909101359150565b6020808252825182820181905260009190848201906040850190845b81811015613b7057835183529284019291840191600101613b54565b50909695505050505050565b803561ffff81168114613b8e57600080fd5b919050565b600080600060608486031215613ba857600080fd5b83359250613bb860208501613b7c565b9150604084013590509250925092565b600060208284031215613bda57600080fd5b5035919050565b60005b83811015613bfc578181015183820152602001613be4565b50506000910152565b60008151808452613c1d816020860160208601613be1565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006118fb6020830184613c05565b60008060408385031215613c7557600080fd5b82359150613c8560208401613b7c565b90509250929050565b73ffffffffffffffffffffffffffffffffffffffff81168114612f4157600080fd5b600060208284031215613cc257600080fd5b81356118fb81613c8e565b600080600060608486031215613ce257600080fd5b505081359360208301359350604090920135919050565b60008083601f840112613d0b57600080fd5b50813567ffffffffffffffff811115613d2357600080fd5b6020830191508360208260051b85010111156116c657600080fd5b60008060208385031215613d5157600080fd5b823567ffffffffffffffff811115613d6857600080fd5b613d7485828601613cf9565b90969095509350505050565b60008060408385031215613d9357600080fd5b823567ffffffffffffffff80821115613dab57600080fd5b908401906101008287031215613dc057600080fd5b90925060208401359080821115613dd657600080fd5b50613de385828601613a40565b9150509250929050565b82151581526040602082015260006118f86040830184613c05565b60008083601f840112613e1a57600080fd5b50813567ffffffffffffffff811115613e3257600080fd5b6020830191508360208285010111156116c657600080fd5b60008060208385031215613e5d57600080fd5b823567ffffffffffffffff811115613e7457600080fd5b613d7485828601613e08565b60008060408385031215613e9357600080fd5b823567ffffffffffffffff80821115613eab57600080fd5b818501915085601f830112613ebf57600080fd5b81356020613ecf613929836138b1565b82815260059290921b84018101918181019089841115613eee57600080fd5b8286015b84811015613f2657803586811115613f0a5760008081fd5b613f188c86838b0101613a40565b845250918301918301613ef2565b5096505086013592505080821115613dd657600080fd5b600080600060408486031215613f5257600080fd5b833567ffffffffffffffff811115613f6957600080fd5b613f7586828701613cf9565b9094509250506020840135613f8981613a2e565b809150509250925092565b8015158114612f4157600080fd5b600060208284031215613fb457600080fd5b81356118fb81613f94565b600080600060408486031215613fd457600080fd5b83359250602084013567ffffffffffffffff811115613ff257600080fd5b613ffe86828701613e08565b9497909650939450505050565b6000806040838503121561401e57600080fd5b82359150602083013561403081613a60565b809150509250929050565b6000806040838503121561404e57600080fd5b82359150602083013561403081613a2e565b60006020828403121561407257600080fd5b81356118fb81613a60565b60006020828403121561408f57600080fd5b81356118fb81613a1f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff808316818516818304811182151516156140f4576140f461409a565b02949350505050565b8051613b8e81613f94565b60006020828403121561411a57600080fd5b81516118fb81613f94565b600082601f83011261413657600080fd5b8151614144613929826138d5565b81815284602083860101111561415957600080fd5b61416a826020830160208701613be1565b949350505050565b60006020828403121561418457600080fd5b815167ffffffffffffffff81111561419b57600080fd5b61416a84828501614125565b8281526040602082015260006118f86040830184613c05565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff81036142055761420561409a565b60010192915050565b818103818111156114a1576114a161409a565b808201808211156114a1576114a161409a565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036142655761426561409a565b5060010190565b805160208083015191908110156142ab577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b8051613b8e81613c8e565b600080604083850312156142cf57600080fd5b82516142da81613c8e565b6020939093015192949293505050565b6000602082840312156142fc57600080fd5b81516118fb81613c8e565b600061ffff80831681810361431e5761431e61409a565b6001019392505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261435d57600080fd5b83018035915067ffffffffffffffff82111561437857600080fd5b6020019150600581901b36038213156116c657600080fd5b6000602082840312156143a257600080fd5b5051919050565b6000816143b8576143b861409a565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600181811c908216806143f257607f821691505b6020821081036142ab577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60008154614438816143de565b808552602060018381168015614455576001811461448d576144bb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b89010195506144bb565b866000528260002060005b858110156144b35781548a8201860152908301908401614498565b890184019650505b505050505092915050565b60a0815260006144d960a083018861442b565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b8381101561454b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552614539838361442b565b94860194925060019182019101614500565b5050868103604088015261455f818b61442b565b945050505050846060840152828103608084015261457d8185613c05565b98975050505050505050565b6000806040838503121561459c57600080fd5b505080516020909101519092909150565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015614622577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552614610868351613c05565b955093820193908201906001016145d6565b5050858403818701525050506146388185613c05565b95945050505050565b6000602080838503121561465457600080fd5b825167ffffffffffffffff81111561466b57600080fd5b8301601f8101851361467c57600080fd5b805161468a613929826138b1565b81815260059190911b820183019083810190878311156146a957600080fd5b928401925b828410156146c7578351825292840192908401906146ae565b979650505050505050565b6000602082840312156146e457600080fd5b81516118fb81613a1f565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b8051613b8e81613a2e565b8051613b8e81613a60565b805167ffffffffffffffff81168114613b8e57600080fd5b60006020828403121561478357600080fd5b815167ffffffffffffffff8082111561479b57600080fd5b9083019061014082860312156147b057600080fd5b6147b8613838565b6147c1836142b1565b81526147cf60208401614743565b60208201526040830151828111156147e657600080fd5b6147f287828601614125565b6040830152506148046060840161474e565b6060820152614815608084016142b1565b608082015261482660a08401614759565b60a082015261483760c08401614743565b60c082015261484860e0840161474e565b60e082015261010061485b8185016140fd565b90820152610120838101518381111561487357600080fd5b61487f88828701614125565b918301919091525095945050505050565b6000602082840312156148a257600080fd5b81516118fb81613a60565b600063ffffffff80831681810361431e5761431e61409a565b60208152600082516101408060208501526148e5610160850183613c05565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808685030160408701526149218483613c05565b93506040870151915061494c606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526149ad8483613c05565b935060e087015191506101008187860301818801526149cc8584613c05565b9450808801519250506101208187860301818801526149eb8584613c05565b94508088015192505050614a0e828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614a5057614a5061409a565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082614a9357614a93614a55565b500490565b81810360008312801583831316838312821617156136d3576136d361409a565b600082614ac757614ac7614a55565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615614b1b57614b1b61409a565b500590565b8082018281126000831280158216821582161715614b4057614b4061409a565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036142655761426561409a565b60007f800000000000000000000000000000000000000000000000000000000000000082036143b8576143b861409a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b601f821115614c1f57600081815260208120601f850160051c81016020861015614c005750805b601f850160051c820191505b818110156129cf57828155600101614c0c565b505050565b815167ffffffffffffffff811115614c3e57614c3e613809565b614c5281614c4c84546143de565b84614bd9565b602080601f831160018114614ca55760008415614c6f5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556129cf565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015614cf257888601518255948401946001909101908401614cd3565b5085821015614d2e57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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
