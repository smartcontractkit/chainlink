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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"firstPerformBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"RegistrarSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkDatas\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxBucketedDelaysForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxDelayForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumTimestampBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTimestampBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getTimestampDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampBuckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6005601855601980546001600160681b0319166c140000000002c68af0bb140000179055606460a052610e1060c0526101a0604052604261012081815260e091829190620064386101403981526020016040518060800160405280604281526020016200647a6042913990526200007b90601a90600262000384565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601b90620000ab908262000500565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b6020820152601c90620000dd908262000500565b50348015620000eb57600080fd5b50604051620064bc380380620064bc8339810160408190526200010e91620005e2565b81813380600081620001675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200019a576200019a81620002d9565b5050601580546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa158015620001f7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200021d919062000625565b50601780546001600160a01b0319166001600160a01b038381169190911790915560155460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801562000283573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002a9919062000656565b601680546001600160a01b0319166001600160a01b0392909216919091179055501515608052506200067d915050565b336001600160a01b03821603620003335760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200015e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215620003cf579160200282015b82811115620003cf5782518290620003be908262000500565b5091602001919060010190620003a5565b50620003dd929150620003e1565b5090565b80821115620003dd576000620003f8828262000402565b50600101620003e1565b508054620004109062000471565b6000825580601f1062000421575050565b601f01602090049060005260206000209081019062000441919062000444565b50565b5b80821115620003dd576000815560010162000445565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200048657607f821691505b602082108103620004a757634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620004fb57600081815260208120601f850160051c81016020861015620004d65750805b601f850160051c820191505b81811015620004f757828155600101620004e2565b5050505b505050565b81516001600160401b038111156200051c576200051c6200045b565b62000534816200052d845462000471565b84620004ad565b602080601f8311600181146200056c5760008415620005535750858301515b600019600386901b1c1916600185901b178555620004f7565b600085815260208120601f198616915b828110156200059d578886015182559484019460019091019084016200057c565b5085821015620005bc5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200044157600080fd5b60008060408385031215620005f657600080fd5b82516200060381620005cc565b602084015190925080151581146200061a57600080fd5b809150509250929050565b600080604083850312156200063957600080fd5b82516200064681620005cc565b6020939093015192949293505050565b6000602082840312156200066957600080fd5b81516200067681620005cc565b9392505050565b60805160a05160c051615d76620006c260003960008181610792015261203801526000818161066a015261214c015260008181610a7701526140880152615d766000f3fe6080604052600436106105415760003560e01c80637b103999116102af578063a79c404311610179578063d6051a72116100d6578063f2fde38b1161008a578063fba7ffa31161006f578063fba7ffa314611186578063fbfb4f76146111b3578063fcdc1f63146111d357600080fd5b8063f2fde38b14611146578063fb0ceb041461116657600080fd5b8063dbef701e116100bb578063dbef701e146110e3578063e0114adb14611103578063e45530831461113057600080fd5b8063d6051a72146110a3578063daee1aeb146110c357600080fd5b8063becde0e11161012d578063c804802211610112578063c80480221461100f578063c98f10b01461102f578063d35585281461104457600080fd5b8063becde0e114610f95578063c357f1f314610fb557600080fd5b8063afb28d1f1161015e578063afb28d1f14610f22578063b0971e1a14610f37578063be61b77514610f7557600080fd5b8063a79c404314610ed5578063af953a4a14610f0257600080fd5b806399cc6b0b116102275780639d6f1cc7116101db578063a6548248116101c0578063a654824814610e61578063a6c60d8914610e95578063a72aa27e14610eb557600080fd5b80639d6f1cc714610e21578063a5f5893414610e4157600080fd5b80639b4293541161020c5780639b42935414610da35780639b51fb0d14610dd05780639d385eaa14610e0157600080fd5b806399cc6b0b14610d475780639ac542eb14610d6757600080fd5b806387dfa9001161027e5780638da5cb5b116102635780638da5cb5b14610ccf5780638fcb3fba14610cfa578063948108f714610d2757600080fd5b806387dfa90014610c8f5780638bc7b77214610caf57600080fd5b80637b10399914610c025780637e4087b814610c2f5780637e7a46dc14610c4f5780638237831714610c6f57600080fd5b806346e7a63e1161040b578063642f6cef116103685780637145f11b1161031c57806376721303116103015780637672130314610ba0578063776898c814610bcd57806379ba509714610bed57600080fd5b80637145f11b14610b4357806373644cce14610b7357600080fd5b806369cdbadb1161034d57806369cdbadb14610ac957806369e9b77314610af65780637137a70214610b2357600080fd5b8063642f6cef14610a65578063643b34e914610aa957600080fd5b806359710992116103bf5780635f17e616116103a45780635f17e616146109d657806360457ff5146109f6578063636092e814610a2357600080fd5b806359710992146109ac5780635d4ee7f3146109c157600080fd5b806351c98be3116103f057806351c98be31461093f57806357970e931461095f57806358c52c041461098c57600080fd5b806346e7a63e146108e45780634b56a42e1461091157600080fd5b806320e3dbd4116104b9578063328ffd111161046d5780633ebe8d6c116104525780633ebe8d6c146108845780634585e33b146108a457806345d2ec17146108c457600080fd5b8063328ffd111461082657806333774d1c1461085357600080fd5b806329f0e4961161049e57806329f0e496146107805780632a9032d3146107b45780632b20e397146107d457600080fd5b806320e3dbd41461074057806328c4b57b1461076057600080fd5b80630d4a4fb111610510578063177b0eb9116104f5578063177b0eb91461069f5780631bee0080146106dd578063206c32e81461070b57600080fd5b80630d4a4fb11461062b57806312c550271461065857600080fd5b806305e251311461058557806306c1cc00146105a757806306e3b632146105c7578063077ac621146105fd57600080fd5b3661058057604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561059157600080fd5b506105a56105a03660046147d2565b611200565b005b3480156105b357600080fd5b506105a56105c23660046148fa565b611217565b3480156105d357600080fd5b506105e76105e2366004614992565b6115d3565b6040516105f491906149ef565b60405180910390f35b34801561060957600080fd5b5061061d610618366004614a14565b6116d2565b6040519081526020016105f4565b34801561063757600080fd5b5061064b610646366004614a49565b611710565b6040516105f49190614ad0565b34801561066457600080fd5b5061068c7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff90911681526020016105f4565b3480156106ab57600080fd5b5061061d6106ba366004614ae3565b6000918252600f6020908152604080842061ffff93909316845291905290205490565b3480156106e957600080fd5b506106fd6106f8366004614a49565b611830565b6040516105f4929190614b0f565b34801561071757600080fd5b5061072b610726366004614ae3565b611b2a565b604080519283526020830191909152016105f4565b34801561074c57600080fd5b506105a561075b366004614b56565b611bad565b34801561076c57600080fd5b5061061d61077b366004614b73565b611dab565b34801561078c57600080fd5b5061068c7f000000000000000000000000000000000000000000000000000000000000000081565b3480156107c057600080fd5b506105a56107cf366004614be4565b611e16565b3480156107e057600080fd5b506015546108019073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016105f4565b34801561083257600080fd5b5061061d610841366004614a49565b60036020526000908152604090205481565b34801561085f57600080fd5b5061068c61086e366004614a49565b60116020526000908152604090205461ffff1681565b34801561089057600080fd5b5061061d61089f366004614a49565b611ee9565b3480156108b057600080fd5b506105a56108bf366004614c68565b611f52565b3480156108d057600080fd5b506105e76108df366004614ae3565b612614565b3480156108f057600080fd5b5061061d6108ff366004614a49565b600a6020526000908152604090205481565b34801561091d57600080fd5b5061093161092c366004614c9e565b612683565b6040516105f4929190614d72565b34801561094b57600080fd5b506105a561095a366004614d8d565b6126d7565b34801561096b57600080fd5b506016546108019073ffffffffffffffffffffffffffffffffffffffff1681565b34801561099857600080fd5b5061064b6109a7366004614a49565b61277b565b3480156109b857600080fd5b506105a5612815565b3480156109cd57600080fd5b506105a5612982565b3480156109e257600080fd5b506105a56109f1366004614992565b612ab9565b348015610a0257600080fd5b5061061d610a11366004614a49565b60076020526000908152604090205481565b348015610a2f57600080fd5b50601954610a48906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff90911681526020016105f4565b348015610a7157600080fd5b50610a997f000000000000000000000000000000000000000000000000000000000000000081565b60405190151581526020016105f4565b348015610ab557600080fd5b5061072b610ac4366004614992565b612c2b565b348015610ad557600080fd5b5061061d610ae4366004614a49565b60086020526000908152604090205481565b348015610b0257600080fd5b506105a5610b11366004614992565b60009182526008602052604090912055565b348015610b2f57600080fd5b5061061d610b3e366004614a14565b612da1565b348015610b4f57600080fd5b50610a99610b5e366004614a49565b600c6020526000908152604090205460ff1681565b348015610b7f57600080fd5b5061061d610b8e366004614a49565b6000908152600d602052604090205490565b348015610bac57600080fd5b5061061d610bbb366004614a49565b60046020526000908152604090205481565b348015610bd957600080fd5b50610a99610be8366004614a49565b612dc9565b348015610bf957600080fd5b506105a5612e1b565b348015610c0e57600080fd5b506017546108019073ffffffffffffffffffffffffffffffffffffffff1681565b348015610c3b57600080fd5b5061072b610c4a366004614992565b612f1d565b348015610c5b57600080fd5b506105a5610c6a366004614de4565b613086565b348015610c7b57600080fd5b5061061d610c8a366004614e30565b61312c565b348015610c9b57600080fd5b5061061d610caa366004614e30565b6131a7565b348015610cbb57600080fd5b506106fd610cca366004614a49565b613217565b348015610cdb57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610801565b348015610d0657600080fd5b5061061d610d15366004614a49565b60056020526000908152604090205481565b348015610d3357600080fd5b506105a5610d42366004614e65565b61338c565b348015610d5357600080fd5b506105e7610d62366004614ae3565b613515565b348015610d7357600080fd5b50601954610d91906c01000000000000000000000000900460ff1681565b60405160ff90911681526020016105f4565b348015610daf57600080fd5b506105a5610dbe366004614992565b60009182526009602052604090912055565b348015610ddc57600080fd5b5061068c610deb366004614a49565b60126020526000908152604090205461ffff1681565b348015610e0d57600080fd5b506105e7610e1c366004614a49565b613582565b348015610e2d57600080fd5b5061064b610e3c366004614a49565b6135e4565b348015610e4d57600080fd5b5061061d610e5c366004614a49565b61360f565b348015610e6d57600080fd5b5061061d7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0881565b348015610ea157600080fd5b506105a5610eb0366004614a49565b601855565b348015610ec157600080fd5b506105a5610ed0366004614e95565b613670565b348015610ee157600080fd5b506105a5610ef0366004614992565b60009182526007602052604090912055565b348015610f0e57600080fd5b506105a5610f1d366004614a49565b61371b565b348015610f2e57600080fd5b5061064b6137a1565b348015610f4357600080fd5b5061061d610f52366004614ae3565b6000918252600e6020908152604080842061ffff93909316845291905290205490565b348015610f8157600080fd5b50610931610f90366004614eba565b6137ae565b348015610fa157600080fd5b506105a5610fb0366004614be4565b613a56565b348015610fc157600080fd5b506105a5610fd0366004614ef6565b601980547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b34801561101b57600080fd5b506105a561102a366004614a49565b613af0565b34801561103b57600080fd5b5061064b613b88565b34801561105057600080fd5b506105a561105f366004614f13565b6019805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b3480156110af57600080fd5b5061072b6110be366004614992565b613b95565b3480156110cf57600080fd5b506105a56110de366004614be4565b613bfe565b3480156110ef57600080fd5b5061061d6110fe366004614992565b613cc9565b34801561110f57600080fd5b5061061d61111e366004614a49565b60096020526000908152604090205481565b34801561113c57600080fd5b5061061d60185481565b34801561115257600080fd5b506105a5611161366004614b56565b613cfa565b34801561117257600080fd5b5061061d611181366004614992565b613d0e565b34801561119257600080fd5b5061061d6111a1366004614a49565b60066020526000908152604090205481565b3480156111bf57600080fd5b5061072b6111ce366004614ae3565b613d2a565b3480156111df57600080fd5b5061061d6111ee366004614a49565b60026020526000908152604090205481565b805161121390601a9060208401906145a2565b5050565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601654601554919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b39216906112fd908c1688614f5d565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af115801561137b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061139f9190614fa1565b5060008860ff1667ffffffffffffffff8111156113be576113be614682565b6040519080825280602002602001820160405280156113e7578160200160208202803683370190505b50905060005b8960ff168160ff16101561159057600061140684613d9e565b90508860ff1660010361153e576040517f0d4a4fb1000000000000000000000000000000000000000000000000000000008152600481018290526000903090630d4a4fb190602401600060405180830381865afa15801561146b573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114b19190810190615009565b6017546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d359061150a908590859060040161503e565b600060405180830381600087803b15801561152457600080fd5b505af1158015611538573d6000803e3d6000fd5b50505050505b80838360ff168151811061155457611554615057565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061158881615086565b9150506113ed565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711816040516115c091906149ef565b60405180910390a1505050505050505050565b606060006115e16013613e8a565b905080841061161c576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826000036116315761162e84826150a5565b92505b60008367ffffffffffffffff81111561164c5761164c614682565b604051908082528060200260200182016040528015611675578160200160208202803683370190505b50905060005b848110156116c75761169861169082886150b8565b601390613e94565b8282815181106116aa576116aa615057565b6020908102919091010152806116bf816150cb565b91505061167b565b509150505b92915050565b600e60205282600052604060002060205281600052604060002081815481106116fa57600080fd5b9060005260206000200160009250925050505481565b606060006040518060c001604052803073ffffffffffffffffffffffffffffffffffffffff168152602001600160ff1681526020017f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0860001b81526020018460405160200161178191815260200190565b60405160208183030381529060405261179990615103565b81526020016000801b81526020016000801b8152509050806040516020016118199190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b604051602081830303815290604052915050919050565b606080600061183f6013613e8a565b905060008167ffffffffffffffff81111561185c5761185c614682565b604051908082528060200260200182016040528015611885578160200160208202803683370190505b50905060008267ffffffffffffffff8111156118a3576118a3614682565b6040519080825280602002602001820160405280156118cc578160200160208202803683370190505b50905060005b83811015611b1e5760006118e7601383613e94565b9050808483815181106118fc576118fc615057565b6020908102919091018101919091526000828152601290915260408082205490517f3ebe8d6c0000000000000000000000000000000000000000000000000000000081526004810184905261ffff90911691903090633ebe8d6c90602401602060405180830381865afa158015611977573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061199b9190615148565b905060008167ffffffffffffffff8111156119b8576119b8614682565b6040519080825280602002602001820160405280156119e1578160200160208202803683370190505b506000858152600e6020526040812091925090815b8561ffff168161ffff1611611adb5761ffff811660009081526020838152604080832080548251818502810185019093528083529192909190830182828015611a5e57602002820191906000526020600020905b815481526020019060010190808311611a4a575b5050505050905060005b8151811015611ac657818181518110611a8357611a83615057565b6020026020010151868680611a97906150cb565b975081518110611aa957611aa9615057565b602090810291909101015280611abe816150cb565b915050611a68565b50508080611ad390615161565b9150506119f6565b50611ae7838e86613ea0565b888881518110611af957611af9615057565b6020026020010181815250505050505050508080611b16906150cb565b9150506118d2565b50909590945092505050565b6000828152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938493929190830182828015611b8e57602002820191906000526020600020905b815481526020019060010190808311611b7a575b50505050509050611ba0818251613fff565b92509250505b9250929050565b601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611c43573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c67919061518d565b50601780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601554604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015611d0a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d2e91906151bb565b601680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601554604051911681527f6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014906020015b60405180910390a15050565b6000838152600d602090815260408083208054825181850281018501909352808352611e0c93830182828015611e0057602002820191906000526020600020905b815481526020019060010190808311611dec575b50505050508484613ea0565b90505b9392505050565b8060005b818160ff161015611eaa573063c8048022858560ff8516818110611e4057611e40615057565b905060200201356040518263ffffffff1660e01b8152600401611e6591815260200190565b600060405180830381600087803b158015611e7f57600080fd5b505af1158015611e93573d6000803e3d6000fd5b505050508080611ea290615086565b915050611e1a565b507fbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b8383604051611edc9291906151d8565b60405180910390a1505050565b60008181526012602052604081205461ffff1681805b8261ffff168161ffff1611611f4a576000858152600e6020908152604080832061ffff85168452909152902054611f3690836150b8565b915080611f4281615161565b915050611eff565b509392505050565b60005a9050600080611f6684860186614c9e565b9150915060008082806020019051810190611f81919061522a565b6000828152600560209081526040808320546004909252822054939550919350909190611fac614084565b905082600003611fe95760008581526005602090815260408083208490556010825282208054600181018255908352912042910155915081612231565b6000611ff585836150a5565b6000878152601160209081526040808320546010909252909120805492935061ffff909116918290811061202b5761202b615057565b90600052602060002001547f000000000000000000000000000000000000000000000000000000000000000061ffff164261206691906150a5565b11156120d5576000878152601060209081526040822080546001810182559083529120429101558061209781615161565b600089815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559150505b600087815260126020908152604080832054600e835281842061ffff909116808552908352818420805483518186028101860190945280845291949390919083018282801561214357602002820191906000526020600020905b81548152602001906001019080831161212f575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff168151036121be578161218081615161565b60008b815260126020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000888152600e6020908152604080832061ffff94851684528252808320805460018181018355918552838520018790558b8452600f83528184209590941683529381528382208054808501825590835281832001859055898252600d81529281208054928301815581529190912001555b60008581526006602052604081205461224b9060016150b8565b600087815260066020908152604091829020839055815189815290810187905290810184905260608101859052608081018290529091507fe1a58b2118f7a6020491ff3fea3e628421dc7392e78ba803adcec9320117af249060a00160405180910390a160008681526004602090815260408083208590556018546002909252909120546122d990846150a5565b1115612566576017546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810188905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa15801561234f573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612395919081019061527c565b6017546040517fb657bc9c000000000000000000000000000000000000000000000000000000008152600481018a905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa15801561240a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061242e91906153ad565b6019549091506124529082906c01000000000000000000000000900460ff16614f5d565b6bffffffffffffffffffffffff1682608001516bffffffffffffffffffffffff161015612563576019546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018a90526bffffffffffffffffffffffff9091166024820152309063948108f790604401600060405180830381600087803b1580156124e357600080fd5b505af11580156124f7573d6000803e3d6000fd5b50505060008981526002602090815260409182902087905560195482518c81526bffffffffffffffffffffffff909116918101919091529081018690527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0915060600160405180910390a15b50505b604051308152829087907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf089060200160405180910390a36000868152600760205260409020545b805a6125b9908c6150a5565b6125c5906127106150b8565b10156126065782406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556125ad565b505050505050505050505050565b6000828152600e6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561267657602002820191906000526020600020905b815481526020019060010190808311612662575b5050505050905092915050565b600060606000848460405160200161269c9291906153ca565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b8160005b818110156127745730635f17e6168686848181106126fb576126fb615057565b90506020020135856040518363ffffffff1660e01b815260040161272f92919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561274957600080fd5b505af115801561275d573d6000803e3d6000fd5b50505050808061276c906150cb565b9150506126db565b5050505050565b600b602052600090815260409020805461279490615455565b80601f01602080910402602001604051908101604052809291908181526020018280546127c090615455565b801561280d5780601f106127e25761010080835404028352916020019161280d565b820191906000526020600020905b8154815290600101906020018083116127f057829003601f168201915b505050505081565b6017546040517f4184e12c0000000000000000000000000000000000000000000000000000000081526000600482018190526103e86024830152600160448301529173ffffffffffffffffffffffffffffffffffffffff1690634184e12c90606401600060405180830381865afa158015612894573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526128da91908101906154a2565b805190915060006128e9614084565b905060005b8281101561297c57600084828151811061290a5761290a615057565b6020026020010151905082817f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0830604051612961919073ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b60405180910390a35080612974816150cb565b9150506128ee565b50505050565b61298a614126565b6016546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156129f9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a1d9190615148565b6016546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af1158015612a95573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112139190614fa1565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600d9091528120612af1916145f8565b60008281526012602052604081205461ffff16905b8161ffff168161ffff1611612b4d576000848152600e6020908152604080832061ffff851684529091528120612b3b916145f8565b80612b4581615161565b915050612b06565b5050600082815260126020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055601190915281205461ffff16905b8161ffff168161ffff1611612bdb576000848152600f6020908152604080832061ffff851684529091528120612bc9916145f8565b80612bd381615161565b915050612b94565b506000838152601060205260408120612bf3916145f8565b5050600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6040517f3ebe8d6c00000000000000000000000000000000000000000000000000000000815260048101839052600090819081903090633ebe8d6c90602401602060405180830381865afa158015612c87573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612cab9190615148565b9050831580612cba5750808410155b15612cc3578093505b60008581526012602052604081205485919061ffff16805b6000898152600e6020908152604080832061ffff85168452825280832080548251818502810185019093528083529192909190830182828015612d3d57602002820191906000526020600020905b815481526020019060010190808311612d29575b50505050509050600080612d518388613fff565b9092509050612d6082876150b8565b9550612d6c81886150a5565b965060008711612d7e57505050612d94565b5050508080612d8c90615533565b915050612cdb565b5090979596505050505050565b600f60205282600052604060002060205281600052604060002081815481106116fa57600080fd5b6000818152600560205260408120548103612de657506001919050565b600082815260036020908152604080832054600490925290912054612e09614084565b612e1391906150a5565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612ea1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040517fa5f589340000000000000000000000000000000000000000000000000000000081526004810183905260009081908190309063a5f5893490602401602060405180830381865afa158015612f79573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f9d9190615148565b9050831580612fac5750808410155b15612fb5578093505b60008581526011602052604081205485919061ffff16805b6000898152600f6020908152604080832061ffff8516845282528083208054825181850281018501909352808352919290919083018282801561302f57602002820191906000526020600020905b81548152602001906001019080831161301b575b505050505090506000806130438388613fff565b909250905061305282876150b8565b955061305e81886150a5565b96506000871161307057505050612d94565b505050808061307e90615533565b915050612fcd565b6017546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b5906130e09086908690869060040161556f565b600060405180830381600087803b1580156130fa57600080fd5b505af115801561310e573d6000803e3d6000fd5b5050506000848152600b60205260409020905061297c82848361560e565b6000838152600e6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849383018282801561318b57602002820191906000526020600020905b815481526020019060010190808311613177575b5050505050905061319e81858351613ea0565b95945050505050565b6000838152600f6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849383018282801561318b5760200282019190600052602060002090815481526020019060010190808311613177575050505050905061319e81858351613ea0565b60608060006132266013613e8a565b905060008167ffffffffffffffff81111561324357613243614682565b60405190808252806020026020018201604052801561326c578160200160208202803683370190505b50905060008267ffffffffffffffff81111561328a5761328a614682565b6040519080825280602002602001820160405280156132b3578160200160208202803683370190505b50905060005b83811015611b1e5760006132ce601383613e94565b6000818152600d602090815260408083208054825181850281018501909352808352949550929390929183018282801561332757602002820191906000526020600020905b815481526020019060010190808311613313575b505050505090508185848151811061334157613341615057565b602002602001018181525050613359818a8351613ea0565b84848151811061336b5761336b615057565b60200260200101818152505050508080613384906150cb565b9150506132b9565b6016546017546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015613414573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134389190614fa1565b506017546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b1580156134b957600080fd5b505af11580156134cd573d6000803e3d6000fd5b5050604080518581526bffffffffffffffffffffffff851660208201527f8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e9350019050611d9f565b6000828152600f6020908152604080832061ffff851684528252918290208054835181840281018401909452808452606093928301828280156126765760200282019190600052602060002090815481526020019060010190808311612662575050505050905092915050565b6000818152600d60209081526040918290208054835181840281018401909452808452606093928301828280156135d857602002820191906000526020600020905b8154815260200190600101908083116135c4575b50505050509050919050565b601a81815481106135f457600080fd5b90600052602060002001600091509050805461279490615455565b60008181526011602052604081205461ffff1681805b8261ffff168161ffff1611611f4a576000858152600f6020908152604080832061ffff8516845290915290205461365c90836150b8565b91508061366881615161565b915050613625565b6017546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156136e857600080fd5b505af11580156136fc573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6017546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561378d57600080fd5b505af1158015612774573d6000803e3d6000fd5b601b805461279490615455565b6000606060005a905060006137c1614084565b90507f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086137f160c0870187615728565b600081811061380257613802615057565b90506020020135036139ce57600061381d60c0870187615728565b600181811061382e5761382e615057565b9050602002013560405160200161384791815260200190565b604051602081830303815290604052905060008180602001905181019061386e9190615148565b9050600061387f60c0890189615728565b600281811061389057613890615057565b905060200201356040516020016138a991815260200190565b60405160208183030381529060405290506000818060200190518101906138d09190615148565b6000848152600860205260409020549091505b805a6138ef90896150a5565b6138fb90613a986150b8565b10156139495781406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558161394181615790565b9250506138e3565b601b601a601c84878660405160200161396c929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e000000000000000000000000000000000000000000000000000000008252612e989594939291600401615860565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f636f756c64206e6f742066696e64206d61746368696e67206576656e7420736960448201527f67000000000000000000000000000000000000000000000000000000000000006064820152608401612e98565b8060005b818163ffffffff16101561297c573063af953a4a858563ffffffff8516818110613a8657613a86615057565b905060200201356040518263ffffffff1660e01b8152600401613aab91815260200190565b600060405180830381600087803b158015613ac557600080fd5b505af1158015613ad9573d6000803e3d6000fd5b505050508080613ae890615923565b915050613a5a565b6017546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b158015613b5c57600080fd5b505af1158015613b70573d6000803e3d6000fd5b505050506112138160136141a990919063ffffffff16565b601c805461279490615455565b6000828152600d60209081526040808320805482518185028101850190935280835284938493929190830182828015613bed57602002820191906000526020600020905b815481526020019060010190808311613bd9575b50505050509050611ba08185613fff565b8060005b8181101561297c576000848483818110613c1e57613c1e615057565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc8283604051602001613c5791815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613c8392919061503e565b600060405180830381600087803b158015613c9d57600080fd5b505af1158015613cb1573d6000803e3d6000fd5b50505050508080613cc1906150cb565b915050613c02565b600d6020528160005260406000208181548110613ce557600080fd5b90600052602060002001600091509150505481565b613d02614126565b613d0b816141b5565b50565b60106020528160005260406000208181548110613ce557600080fd5b6000828152600f6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938493929190830182828015611b8e5760200282019190600052602060002090815481526020019060010190808311611b7a5750505050509050611ba0818251613fff565b6015546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613df990869060040161593c565b6020604051808303816000875af1158015613e18573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e3c9190615148565b9050613e496013826142aa565b5060608301516000828152600a6020908152604080832063ffffffff90941690935560c0860151600b90915291902090613e839082615a8e565b5092915050565b60006116cc825490565b6000611e0f83836142b6565b82516000908190831580613eb45750808410155b15613ebd578093505b60008467ffffffffffffffff811115613ed857613ed8614682565b604051908082528060200260200182016040528015613f01578160200160208202803683370190505b509050600092505b84831015613f6f57866001613f1e85856150a5565b613f2891906150a5565b81518110613f3857613f38615057565b6020026020010151818481518110613f5257613f52615057565b602090810291909101015282613f67816150cb565b935050613f09565b613f8881600060018451613f8391906150a5565b6142e0565b85606403613fc1578060018251613f9f91906150a5565b81518110613faf57613faf615057565b60200260200101519350505050611e0f565b806064825188613fd19190615ba8565b613fdb9190615c14565b81518110613feb57613feb615057565b602002602001015193505050509392505050565b8151600090819081908415806140155750808510155b1561401e578094505b60008092505b8583101561407a5786600161403985856150a5565b61404391906150a5565b8151811061405357614053615057565b60200260200101518161406691906150b8565b905082614072816150cb565b935050614024565b9694955050505050565b60007f00000000000000000000000000000000000000000000000000000000000000001561412157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156140f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061411c9190615148565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff1633146141a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401612e98565b565b6000611e0f8383614460565b3373ffffffffffffffffffffffffffffffffffffffff821603614234576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401612e98565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611e0f8383614553565b60008260000182815481106142cd576142cd615057565b9060005260206000200154905092915050565b81818082036142f0575050505050565b60008560026142ff8787615c28565b6143099190615c48565b6143139087615cb0565b8151811061432357614323615057565b602002602001015190505b818313614432575b8086848151811061434957614349615057565b60200260200101511015614369578261436181615cd8565b935050614336565b85828151811061437b5761437b615057565b602002602001015181101561439c578161439481615d09565b925050614369565b81831361442d578582815181106143b5576143b5615057565b60200260200101518684815181106143cf576143cf615057565b60200260200101518785815181106143e9576143e9615057565b6020026020010188858151811061440257614402615057565b6020908102919091010191909152528261441b81615cd8565b935050818061442990615d09565b9250505b61432e565b81851215614445576144458686846142e0565b83831215614458576144588684866142e0565b505050505050565b600081815260018301602052604081205480156145495760006144846001836150a5565b8554909150600090614498906001906150a5565b90508181146144fd5760008660000182815481106144b8576144b8615057565b90600052602060002001549050808760000184815481106144db576144db615057565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061450e5761450e615d3a565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506116cc565b60009150506116cc565b600081815260018301602052604081205461459a575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556116cc565b5060006116cc565b8280548282559060005260206000209081019282156145e8579160200282015b828111156145e857825182906145d89082615a8e565b50916020019190600101906145c2565b506145f4929150614616565b5090565b5080546000825590600052602060002090810190613d0b9190614633565b808211156145f457600061462a8282614648565b50600101614616565b5b808211156145f45760008155600101614634565b50805461465490615455565b6000825580601f10614664575050565b601f016020900490600052602060002090810190613d0b9190614633565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610160810167ffffffffffffffff811182821017156146d5576146d5614682565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561472257614722614682565b604052919050565b600067ffffffffffffffff82111561474457614744614682565b5060051b60200190565b600067ffffffffffffffff82111561476857614768614682565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006147a76147a28461474e565b6146db565b90508281528383830111156147bb57600080fd5b828260208301376000602084830101529392505050565b600060208083850312156147e557600080fd5b823567ffffffffffffffff808211156147fd57600080fd5b818501915085601f83011261481157600080fd5b813561481f6147a28261472a565b81815260059190911b8301840190848101908883111561483e57600080fd5b8585015b8381101561488b5780358581111561485a5760008081fd5b8601603f81018b1361486c5760008081fd5b61487d8b8983013560408401614794565b845250918601918601614842565b5098975050505050505050565b803560ff811681146148a957600080fd5b919050565b63ffffffff81168114613d0b57600080fd5b600082601f8301126148d157600080fd5b611e0f83833560208501614794565b6bffffffffffffffffffffffff81168114613d0b57600080fd5b600080600080600080600060e0888a03121561491557600080fd5b61491e88614898565b9650602088013561492e816148ae565b955061493c60408901614898565b9450606088013567ffffffffffffffff81111561495857600080fd5b6149648a828b016148c0565b9450506080880135614975816148e0565b9699959850939692959460a0840135945060c09093013592915050565b600080604083850312156149a557600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156149e4578151875295820195908201906001016149c8565b509495945050505050565b602081526000611e0f60208301846149b4565b803561ffff811681146148a957600080fd5b600080600060608486031215614a2957600080fd5b83359250614a3960208501614a02565b9150604084013590509250925092565b600060208284031215614a5b57600080fd5b5035919050565b60005b83811015614a7d578181015183820152602001614a65565b50506000910152565b60008151808452614a9e816020860160208601614a62565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611e0f6020830184614a86565b60008060408385031215614af657600080fd5b82359150614b0660208401614a02565b90509250929050565b604081526000614b2260408301856149b4565b828103602084015261319e81856149b4565b73ffffffffffffffffffffffffffffffffffffffff81168114613d0b57600080fd5b600060208284031215614b6857600080fd5b8135611e0f81614b34565b600080600060608486031215614b8857600080fd5b505081359360208301359350604090920135919050565b60008083601f840112614bb157600080fd5b50813567ffffffffffffffff811115614bc957600080fd5b6020830191508360208260051b8501011115611ba657600080fd5b60008060208385031215614bf757600080fd5b823567ffffffffffffffff811115614c0e57600080fd5b614c1a85828601614b9f565b90969095509350505050565b60008083601f840112614c3857600080fd5b50813567ffffffffffffffff811115614c5057600080fd5b602083019150836020828501011115611ba657600080fd5b60008060208385031215614c7b57600080fd5b823567ffffffffffffffff811115614c9257600080fd5b614c1a85828601614c26565b60008060408385031215614cb157600080fd5b823567ffffffffffffffff80821115614cc957600080fd5b818501915085601f830112614cdd57600080fd5b81356020614ced6147a28361472a565b82815260059290921b84018101918181019089841115614d0c57600080fd5b8286015b84811015614d4457803586811115614d285760008081fd5b614d368c86838b01016148c0565b845250918301918301614d10565b5096505086013592505080821115614d5b57600080fd5b50614d68858286016148c0565b9150509250929050565b8215158152604060208201526000611e0c6040830184614a86565b600080600060408486031215614da257600080fd5b833567ffffffffffffffff811115614db957600080fd5b614dc586828701614b9f565b9094509250506020840135614dd9816148ae565b809150509250925092565b600080600060408486031215614df957600080fd5b83359250602084013567ffffffffffffffff811115614e1757600080fd5b614e2386828701614c26565b9497909650939450505050565b600080600060608486031215614e4557600080fd5b8335925060208401359150614e5c60408501614a02565b90509250925092565b60008060408385031215614e7857600080fd5b823591506020830135614e8a816148e0565b809150509250929050565b60008060408385031215614ea857600080fd5b823591506020830135614e8a816148ae565b600060208284031215614ecc57600080fd5b813567ffffffffffffffff811115614ee357600080fd5b82016101008185031215611e0f57600080fd5b600060208284031215614f0857600080fd5b8135611e0f816148e0565b600060208284031215614f2557600080fd5b611e0f82614898565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614f8857614f88614f2e565b02949350505050565b805180151581146148a957600080fd5b600060208284031215614fb357600080fd5b611e0f82614f91565b600082601f830112614fcd57600080fd5b8151614fdb6147a28261474e565b818152846020838601011115614ff057600080fd5b615001826020830160208701614a62565b949350505050565b60006020828403121561501b57600080fd5b815167ffffffffffffffff81111561503257600080fd5b61500184828501614fbc565b828152604060208201526000611e0c6040830184614a86565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff810361509c5761509c614f2e565b60010192915050565b818103818111156116cc576116cc614f2e565b808201808211156116cc576116cc614f2e565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036150fc576150fc614f2e565b5060010190565b80516020808301519190811015615142577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b60006020828403121561515a57600080fd5b5051919050565b600061ffff80831681810361517857615178614f2e565b6001019392505050565b80516148a981614b34565b600080604083850312156151a057600080fd5b82516151ab81614b34565b6020939093015192949293505050565b6000602082840312156151cd57600080fd5b8151611e0f81614b34565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83111561521157600080fd5b8260051b80856040850137919091016040019392505050565b6000806040838503121561523d57600080fd5b505080516020909101519092909150565b80516148a9816148ae565b80516148a9816148e0565b805167ffffffffffffffff811681146148a957600080fd5b60006020828403121561528e57600080fd5b815167ffffffffffffffff808211156152a657600080fd5b9083019061016082860312156152bb57600080fd5b6152c36146b1565b6152cc83615182565b81526152da60208401615182565b60208201526152eb6040840161524e565b604082015260608301518281111561530257600080fd5b61530e87828601614fbc565b60608301525061532060808401615259565b608082015261533160a08401615182565b60a082015261534260c08401615264565b60c082015261535360e0840161524e565b60e0820152610100615366818501615259565b90820152610120615378848201614f91565b90820152610140838101518381111561539057600080fd5b61539c88828701614fbc565b918301919091525095945050505050565b6000602082840312156153bf57600080fd5b8151611e0f816148e0565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b8381101561543f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261542d868351614a86565b955093820193908201906001016153f3565b50508584038187015250505061319e8185614a86565b600181811c9082168061546957607f821691505b602082108103615142577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060208083850312156154b557600080fd5b825167ffffffffffffffff8111156154cc57600080fd5b8301601f810185136154dd57600080fd5b80516154eb6147a28261472a565b81815260059190911b8201830190838101908783111561550a57600080fd5b928401925b828410156155285783518252928401929084019061550f565b979650505050505050565b600061ffff82168061554757615547614f2e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f82111561560957600081815260208120601f850160051c810160208610156155ea5750805b601f850160051c820191505b81811015614458578281556001016155f6565b505050565b67ffffffffffffffff83111561562657615626614682565b61563a836156348354615455565b836155c3565b6000601f84116001811461568c57600085156156565750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355612774565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156156db57868501358255602094850194600190920191016156bb565b5086821015615716577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261575d57600080fd5b83018035915067ffffffffffffffff82111561577857600080fd5b6020019150600581901b3603821315611ba657600080fd5b60008161579f5761579f614f2e565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600081546157d281615455565b8085526020600183811680156157ef576001811461582757615855565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550615855565b866000528260002060005b8581101561584d5781548a8201860152908301908401615832565b890184019650505b505050505092915050565b60a08152600061587360a08301886157c5565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156158e5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526158d383836157c5565b9486019492506001918201910161589a565b505086810360408801526158f9818b6157c5565b94505050505084606084015282810360808401526159178185614a86565b98975050505050505050565b600063ffffffff80831681810361517857615178614f2e565b602081526000825161014080602085015261595b610160850183614a86565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808685030160408701526159978483614a86565b9350604087015191506159c2606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e0870152615a238483614a86565b935060e08701519150610100818786030181880152615a428584614a86565b945080880151925050610120818786030181880152615a618584614a86565b94508088015192505050615a84828601826bffffffffffffffffffffffff169052565b5090949350505050565b815167ffffffffffffffff811115615aa857615aa8614682565b615abc81615ab68454615455565b846155c3565b602080601f831160018114615b0f5760008415615ad95750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555614458565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015615b5c57888601518255948401946001909101908401615b3d565b5085821015615b9857878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615be057615be0614f2e565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082615c2357615c23615be5565b500490565b8181036000831280158383131683831282161715613e8357613e83614f2e565b600082615c5757615c57615be5565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615615cab57615cab614f2e565b500590565b8082018281126000831280158216821582161715615cd057615cd0614f2e565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036150fc576150fc614f2e565b60007f8000000000000000000000000000000000000000000000000000000000000000820361579f5761579f614f2e565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) TIMESTAMPINTERVAL(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "TIMESTAMP_INTERVAL")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TIMESTAMPINTERVAL() (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TIMESTAMPINTERVAL(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) TIMESTAMPINTERVAL() (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TIMESTAMPINTERVAL(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) CheckDatas(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "checkDatas", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckDatas(arg0 *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckDatas(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) CheckDatas(arg0 *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckDatas(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetDelaysLengthAtBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getDelaysLengthAtBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetDelaysLengthAtBucket(upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLengthAtBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetDelaysLengthAtBucket(upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLengthAtBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetDelaysLengthAtTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getDelaysLengthAtTimestampBucket", upkeepId, timestampBucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetDelaysLengthAtTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLengthAtTimestampBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetDelaysLengthAtTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLengthAtTimestampBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, timestampBucket)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetPxBucketedDelaysForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getPxBucketedDelaysForAllUpkeeps", p)

	if err != nil {
		return *new([]*big.Int), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetPxBucketedDelaysForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxBucketedDelaysForAllUpkeeps(&_VerifiableLoadLogTriggerUpkeep.CallOpts, p)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetPxBucketedDelaysForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxBucketedDelaysForAllUpkeeps(&_VerifiableLoadLogTriggerUpkeep.CallOpts, p)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetPxDelayForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getPxDelayForAllUpkeeps", p)

	if err != nil {
		return *new([]*big.Int), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetPxDelayForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayForAllUpkeeps(&_VerifiableLoadLogTriggerUpkeep.CallOpts, p)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetPxDelayForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayForAllUpkeeps(&_VerifiableLoadLogTriggerUpkeep.CallOpts, p)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetPxDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getPxDelayInBucket", upkeepId, p, bucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetPxDelayInBucket(upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayInBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetPxDelayInBucket(upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayInBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetPxDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getPxDelayInTimestampBucket", upkeepId, p, timestampBucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetPxDelayInTimestampBucket(upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayInTimestampBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, timestampBucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetPxDelayInTimestampBucket(upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayInTimestampBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, timestampBucket)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetSumBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getSumBucketedDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetSumBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumBucketedDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetSumBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumBucketedDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetSumDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getSumDelayInTimestampBucket", upkeepId, timestampBucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetSumDelayInTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayInTimestampBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetSumDelayInTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayInTimestampBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, timestampBucket)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetSumTimestampBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getSumTimestampBucketedDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetSumTimestampBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumTimestampBucketedDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetSumTimestampBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumTimestampBucketedDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetTimestampBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getTimestampBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetTimestampBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetTimestampBucketedDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetTimestampBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetTimestampBucketedDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetTimestampDelays(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getTimestampDelays", upkeepId, timestampBucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetTimestampDelays(upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetTimestampDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetTimestampDelays(upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetTimestampDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, timestampBucket)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) TimestampBuckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "timestampBuckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TimestampBuckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimestampBuckets(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) TimestampBuckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimestampBuckets(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) TimestampDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "timestampDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TimestampDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimestampDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) TimestampDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimestampDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Timestamps(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "timestamps", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Timestamps(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Timestamps(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Timestamps(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Timestamps(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1)
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "cancelUpkeep", upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CancelUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CancelUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) CheckLog(opts *bind.TransactOpts, log Log) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "checkLog", log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckLog(log Log) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) CheckLog(log Log) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log)
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

type VerifiableLoadLogTriggerUpkeepFundsAddedIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepFundsAdded)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepFundsAdded)
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

func (it *VerifiableLoadLogTriggerUpkeepFundsAddedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepFundsAdded struct {
	UpkeepId *big.Int
	Amount   *big.Int
	Raw      types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepFundsAddedIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepFundsAddedIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepFundsAdded) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepFundsAdded)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseFundsAdded(log types.Log) (*VerifiableLoadLogTriggerUpkeepFundsAdded, error) {
	event := new(VerifiableLoadLogTriggerUpkeepFundsAdded)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepInsufficientFundsIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepInsufficientFunds

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepInsufficientFundsIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepInsufficientFunds)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepInsufficientFunds)
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

func (it *VerifiableLoadLogTriggerUpkeepInsufficientFundsIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepInsufficientFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepInsufficientFunds struct {
	Balance  *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepInsufficientFundsIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "InsufficientFunds")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepInsufficientFundsIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "InsufficientFunds", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepInsufficientFunds) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "InsufficientFunds")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepInsufficientFunds)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "InsufficientFunds", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseInsufficientFunds(log types.Log) (*VerifiableLoadLogTriggerUpkeepInsufficientFunds, error) {
	event := new(VerifiableLoadLogTriggerUpkeepInsufficientFunds)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "InsufficientFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

type VerifiableLoadLogTriggerUpkeepPerformingUpkeepIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepPerformingUpkeep)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepPerformingUpkeep)
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

func (it *VerifiableLoadLogTriggerUpkeepPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepPerformingUpkeep struct {
	UpkeepId          *big.Int
	FirstPerformBlock *big.Int
	LastBlock         *big.Int
	PreviousBlock     *big.Int
	Counter           *big.Int
	Raw               types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepPerformingUpkeepIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepPerformingUpkeepIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepPerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepPerformingUpkeep)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParsePerformingUpkeep(log types.Log) (*VerifiableLoadLogTriggerUpkeepPerformingUpkeep, error) {
	event := new(VerifiableLoadLogTriggerUpkeepPerformingUpkeep)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

type VerifiableLoadLogTriggerUpkeepRegistrarSetIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepRegistrarSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepRegistrarSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepRegistrarSet)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepRegistrarSet)
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

func (it *VerifiableLoadLogTriggerUpkeepRegistrarSetIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepRegistrarSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepRegistrarSet struct {
	NewRegistrar common.Address
	Raw          types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterRegistrarSet(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepRegistrarSetIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "RegistrarSet")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepRegistrarSetIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "RegistrarSet", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchRegistrarSet(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepRegistrarSet) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "RegistrarSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepRegistrarSet)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "RegistrarSet", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseRegistrarSet(log types.Log) (*VerifiableLoadLogTriggerUpkeepRegistrarSet, error) {
	event := new(VerifiableLoadLogTriggerUpkeepRegistrarSet)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "RegistrarSet", log); err != nil {
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

type VerifiableLoadLogTriggerUpkeepUpkeepsCancelledIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepUpkeepsCancelled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepsCancelledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepsCancelled)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepsCancelled)
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

func (it *VerifiableLoadLogTriggerUpkeepUpkeepsCancelledIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepsCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepUpkeepsCancelled struct {
	UpkeepIds []*big.Int
	Raw       types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterUpkeepsCancelled(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepsCancelledIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "UpkeepsCancelled")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepUpkeepsCancelledIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "UpkeepsCancelled", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchUpkeepsCancelled(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepsCancelled) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "UpkeepsCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepUpkeepsCancelled)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepsCancelled", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseUpkeepsCancelled(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepsCancelled, error) {
	event := new(VerifiableLoadLogTriggerUpkeepUpkeepsCancelled)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepsCancelled", log); err != nil {
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
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["FundsAdded"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseFundsAdded(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["InsufficientFunds"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseInsufficientFunds(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["PerformingUpkeep"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParsePerformingUpkeep(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["Received"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseReceived(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["RegistrarSet"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseRegistrarSet(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseUpkeepTopUp(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["UpkeepsCancelled"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseUpkeepsCancelled(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["UpkeepsRegistered"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseUpkeepsRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadLogTriggerUpkeepFundsAdded) Topic() common.Hash {
	return common.HexToHash("0x8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e")
}

func (VerifiableLoadLogTriggerUpkeepInsufficientFunds) Topic() common.Hash {
	return common.HexToHash("0x03eb8b54a949acec2cd08fdb6d6bd4647a1f2c907d75d6900648effa92eb147f")
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

func (VerifiableLoadLogTriggerUpkeepPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0xe1a58b2118f7a6020491ff3fea3e628421dc7392e78ba803adcec9320117af24")
}

func (VerifiableLoadLogTriggerUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874")
}

func (VerifiableLoadLogTriggerUpkeepRegistrarSet) Topic() common.Hash {
	return common.HexToHash("0x6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014")
}

func (VerifiableLoadLogTriggerUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (VerifiableLoadLogTriggerUpkeepUpkeepsCancelled) Topic() common.Hash {
	return common.HexToHash("0xbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b")
}

func (VerifiableLoadLogTriggerUpkeepUpkeepsRegistered) Topic() common.Hash {
	return common.HexToHash("0x2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeep) Address() common.Address {
	return _VerifiableLoadLogTriggerUpkeep.address
}

type VerifiableLoadLogTriggerUpkeepInterface interface {
	BUCKETSIZE(opts *bind.CallOpts) (uint16, error)

	TIMESTAMPINTERVAL(opts *bind.CallOpts) (uint16, error)

	AddLinkAmount(opts *bind.CallOpts) (*big.Int, error)

	BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error)

	Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error)

	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	CheckDatas(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error)

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

	GetDelaysLengthAtBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, error)

	GetDelaysLengthAtTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, error)

	GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetPxBucketedDelaysForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error)

	GetPxDelayForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error)

	GetPxDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error)

	GetPxDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error)

	GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error)

	GetSumBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error)

	GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error)

	GetSumDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error)

	GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error)

	GetSumTimestampBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error)

	GetTimestampBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetTimestampDelays(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error)

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

	TimestampBuckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error)

	TimestampDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error)

	Timestamps(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

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

	CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	CheckLog(opts *bind.TransactOpts, log Log) (*types.Transaction, error)

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

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error)

	WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepFundsAdded) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*VerifiableLoadLogTriggerUpkeepFundsAdded, error)

	FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepInsufficientFundsIterator, error)

	WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepInsufficientFunds) (event.Subscription, error)

	ParseInsufficientFunds(log types.Log) (*VerifiableLoadLogTriggerUpkeepInsufficientFunds, error)

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int) (*VerifiableLoadLogTriggerUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadLogTriggerUpkeepLogEmitted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferred, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepPerformingUpkeep) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*VerifiableLoadLogTriggerUpkeepPerformingUpkeep, error)

	FilterReceived(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepReceivedIterator, error)

	WatchReceived(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepReceived) (event.Subscription, error)

	ParseReceived(log types.Log) (*VerifiableLoadLogTriggerUpkeepReceived, error)

	FilterRegistrarSet(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepRegistrarSetIterator, error)

	WatchRegistrarSet(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepRegistrarSet) (event.Subscription, error)

	ParseRegistrarSet(log types.Log) (*VerifiableLoadLogTriggerUpkeepRegistrarSet, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUp, error)

	FilterUpkeepsCancelled(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepsCancelledIterator, error)

	WatchUpkeepsCancelled(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepsCancelled) (event.Subscription, error)

	ParseUpkeepsCancelled(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepsCancelled, error)

	FilterUpkeepsRegistered(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepsRegisteredIterator, error)

	WatchUpkeepsRegistered(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepsRegistered) (event.Subscription, error)

	ParseUpkeepsRegistered(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepsRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
