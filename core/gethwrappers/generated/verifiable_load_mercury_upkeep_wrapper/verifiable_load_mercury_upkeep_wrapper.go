// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifiable_load_mercury_upkeep_wrapper

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

var VerifiableLoadMercuryUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"registrarAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"logBlockNum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"origin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"firstPerformBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"RegistrarSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkDatas\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxBucketedDelaysForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxDelayForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumTimestampBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTimestampBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getTimestampDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampBuckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6005601855601980546001600160681b0319166c140000000002c68af0bb140000179055606460a052610e1060c0526101c0604052604261014081815260e091829190620060b0610160398152602001604051806080016040528060428152602001620060f2604291398152602001604051806080016040528060428152602001620061346042913990526200009a90601a90600362000341565b50348015620000a857600080fd5b506040516200617638038062006176833981016040819052620000cb916200042e565b81813380600081620001245760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200015757620001578162000296565b5050601580546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa158015620001b4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001da919062000471565b50601780546001600160a01b0319166001600160a01b038381169190911790915560155460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801562000240573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002669190620004a2565b601680546001600160a01b0319166001600160a01b0392909216919091179055501515608052506200063a915050565b336001600160a01b03821603620002f05760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200011b565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156200038c579160200282015b828111156200038c57825182906200037b90826200056e565b509160200191906001019062000362565b506200039a9291506200039e565b5090565b808211156200039a576000620003b58282620003bf565b506001016200039e565b508054620003cd90620004df565b6000825580601f10620003de575050565b601f016020900490600052602060002090810190620003fe919062000401565b50565b5b808211156200039a576000815560010162000402565b6001600160a01b0381168114620003fe57600080fd5b600080604083850312156200044257600080fd5b82516200044f8162000418565b602084015190925080151581146200046657600080fd5b809150509250929050565b600080604083850312156200048557600080fd5b8251620004928162000418565b6020939093015192949293505050565b600060208284031215620004b557600080fd5b8151620004c28162000418565b9392505050565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004f457607f821691505b6020821081036200051557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200056957600081815260208120601f850160051c81016020861015620005445750805b601f850160051c820191505b81811015620005655782815560010162000550565b5050505b505050565b81516001600160401b038111156200058a576200058a620004c9565b620005a2816200059b8454620004df565b846200051b565b602080601f831160018114620005da5760008415620005c15750858301515b600019600386901b1c1916600185901b17855562000565565b600085815260208120601f198616915b828110156200060b57888601518255948401946001909101908401620005ea565b50858210156200062a5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c051615a316200067f6000396000818161075c01526120330152600081816106340152612147015260008181610a2c0152613e970152615a316000f3fe60806040526004361061050b5760003560e01c806379ba509711610294578063a79c40431161015e578063d6051a72116100d6578063f2fde38b1161008a578063fba7ffa31161006f578063fba7ffa31461116f578063fbfb4f761461119c578063fcdc1f63146111bc57600080fd5b8063f2fde38b1461112f578063fb0ceb041461114f57600080fd5b8063dbef701e116100bb578063dbef701e146110cc578063e0114adb146110ec578063e45530831461111957600080fd5b8063d6051a721461108c578063daee1aeb146110ac57600080fd5b8063becde0e11161012d578063c804802211610112578063c804802214610fc4578063c98f10b014610fe4578063d35585281461102d57600080fd5b8063becde0e114610f4a578063c357f1f314610f6a57600080fd5b8063a79c404314610e76578063af953a4a14610ea3578063afb28d1f14610ec3578063b0971e1a14610f0c57600080fd5b8063948108f71161020c5780639d385eaa116101c0578063a5f58934116101a5578063a5f5893414610e16578063a6c60d8914610e36578063a72aa27e14610e5657600080fd5b80639d385eaa14610dd65780639d6f1cc714610df657600080fd5b80639ac542eb116101f15780639ac542eb14610d3c5780639b42935414610d785780639b51fb0d14610da557600080fd5b8063948108f714610cfc57806399cc6b0b14610d1c57600080fd5b806382378317116102635780638bc7b772116102485780638bc7b77214610c845780638da5cb5b14610ca45780638fcb3fba14610ccf57600080fd5b80638237831714610c4457806387dfa90014610c6457600080fd5b806379ba509714610bc25780637b10399914610bd75780637e4087b814610c045780637e7a46dc14610c2457600080fd5b806346e7a63e116103d5578063642f6cef1161034d5780637137a7021161030157806373644cce116102e657806373644cce14610b485780637672130314610b75578063776898c814610ba257600080fd5b80637137a70214610af85780637145f11b14610b1857600080fd5b806369cdbadb1161033257806369cdbadb14610a7e57806369e9b77314610aab5780636e04ff0d14610ad857600080fd5b8063642f6cef14610a1a578063643b34e914610a5e57600080fd5b806358c52c04116103a45780635f17e616116103895780635f17e6161461098b57806360457ff5146109ab578063636092e8146109d857600080fd5b806358c52c04146109565780635d4ee7f31461097657600080fd5b806346e7a63e146108ae5780634b56a42e146108db57806351c98be31461090957806357970e931461092957600080fd5b806320e3dbd411610483578063328ffd11116104375780633ebe8d6c1161041c5780633ebe8d6c1461084e5780634585e33b1461086e57806345d2ec171461088e57600080fd5b8063328ffd11146107f057806333774d1c1461081d57600080fd5b806329f0e4961161046857806329f0e4961461074a5780632a9032d31461077e5780632b20e3971461079e57600080fd5b806320e3dbd41461070a57806328c4b57b1461072a57600080fd5b80630d4a4fb1116104da578063177b0eb9116104bf578063177b0eb9146106695780631bee0080146106a7578063206c32e8146106d557600080fd5b80630d4a4fb1146105f557806312c550271461062257600080fd5b806305e251311461054f57806306c1cc001461057157806306e3b63214610591578063077ac621146105c757600080fd5b3661054a57604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561055b57600080fd5b5061056f61056a3660046145e1565b6111e9565b005b34801561057d57600080fd5b5061056f61058c366004614709565b611200565b34801561059d57600080fd5b506105b16105ac3660046147a1565b6115bc565b6040516105be91906147fe565b60405180910390f35b3480156105d357600080fd5b506105e76105e2366004614823565b6116bb565b6040519081526020016105be565b34801561060157600080fd5b50610615610610366004614858565b6116f9565b6040516105be91906148df565b34801561062e57600080fd5b506106567f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff90911681526020016105be565b34801561067557600080fd5b506105e76106843660046148f2565b6000918252600f6020908152604080842061ffff93909316845291905290205490565b3480156106b357600080fd5b506106c76106c2366004614858565b611819565b6040516105be92919061491e565b3480156106e157600080fd5b506106f56106f03660046148f2565b611b13565b604080519283526020830191909152016105be565b34801561071657600080fd5b5061056f610725366004614965565b611b96565b34801561073657600080fd5b506105e7610745366004614982565b611d94565b34801561075657600080fd5b506106567f000000000000000000000000000000000000000000000000000000000000000081565b34801561078a57600080fd5b5061056f6107993660046149f3565b611dff565b3480156107aa57600080fd5b506015546107cb9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016105be565b3480156107fc57600080fd5b506105e761080b366004614858565b60036020526000908152604090205481565b34801561082957600080fd5b50610656610838366004614858565b60116020526000908152604090205461ffff1681565b34801561085a57600080fd5b506105e7610869366004614858565b611ed2565b34801561087a57600080fd5b5061056f610889366004614a77565b611f3b565b34801561089a57600080fd5b506105b16108a93660046148f2565b61263c565b3480156108ba57600080fd5b506105e76108c9366004614858565b600a6020526000908152604090205481565b3480156108e757600080fd5b506108fb6108f6366004614aad565b6126ab565b6040516105be929190614b81565b34801561091557600080fd5b5061056f610924366004614b9c565b6126ff565b34801561093557600080fd5b506016546107cb9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561096257600080fd5b50610615610971366004614858565b6127a3565b34801561098257600080fd5b5061056f61283d565b34801561099757600080fd5b5061056f6109a63660046147a1565b612974565b3480156109b757600080fd5b506105e76109c6366004614858565b60076020526000908152604090205481565b3480156109e457600080fd5b506019546109fd906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff90911681526020016105be565b348015610a2657600080fd5b50610a4e7f000000000000000000000000000000000000000000000000000000000000000081565b60405190151581526020016105be565b348015610a6a57600080fd5b506106f5610a793660046147a1565b612ae6565b348015610a8a57600080fd5b506105e7610a99366004614858565b60086020526000908152604090205481565b348015610ab757600080fd5b5061056f610ac63660046147a1565b60009182526008602052604090912055565b348015610ae457600080fd5b506108fb610af3366004614a77565b612c5c565b348015610b0457600080fd5b506105e7610b13366004614823565b612e71565b348015610b2457600080fd5b50610a4e610b33366004614858565b600c6020526000908152604090205460ff1681565b348015610b5457600080fd5b506105e7610b63366004614858565b6000908152600d602052604090205490565b348015610b8157600080fd5b506105e7610b90366004614858565b60046020526000908152604090205481565b348015610bae57600080fd5b50610a4e610bbd366004614858565b612e99565b348015610bce57600080fd5b5061056f612eeb565b348015610be357600080fd5b506017546107cb9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610c1057600080fd5b506106f5610c1f3660046147a1565b612fe8565b348015610c3057600080fd5b5061056f610c3f366004614bf3565b613151565b348015610c5057600080fd5b506105e7610c5f366004614c3f565b6131fd565b348015610c7057600080fd5b506105e7610c7f366004614c3f565b613278565b348015610c9057600080fd5b506106c7610c9f366004614858565b6132e8565b348015610cb057600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166107cb565b348015610cdb57600080fd5b506105e7610cea366004614858565b60056020526000908152604090205481565b348015610d0857600080fd5b5061056f610d17366004614c74565b61345d565b348015610d2857600080fd5b506105b1610d373660046148f2565b6135e6565b348015610d4857600080fd5b50601954610d66906c01000000000000000000000000900460ff1681565b60405160ff90911681526020016105be565b348015610d8457600080fd5b5061056f610d933660046147a1565b60009182526009602052604090912055565b348015610db157600080fd5b50610656610dc0366004614858565b60126020526000908152604090205461ffff1681565b348015610de257600080fd5b506105b1610df1366004614858565b613653565b348015610e0257600080fd5b50610615610e11366004614858565b6136b5565b348015610e2257600080fd5b506105e7610e31366004614858565b6136e0565b348015610e4257600080fd5b5061056f610e51366004614858565b601855565b348015610e6257600080fd5b5061056f610e71366004614ca4565b613741565b348015610e8257600080fd5b5061056f610e913660046147a1565b60009182526007602052604090912055565b348015610eaf57600080fd5b5061056f610ebe366004614858565b6137ec565b348015610ecf57600080fd5b506106156040518060400160405280600981526020017f666565644964486578000000000000000000000000000000000000000000000081525081565b348015610f1857600080fd5b506105e7610f273660046148f2565b6000918252600e6020908152604080842061ffff93909316845291905290205490565b348015610f5657600080fd5b5061056f610f653660046149f3565b613872565b348015610f7657600080fd5b5061056f610f85366004614cc9565b601980547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610fd057600080fd5b5061056f610fdf366004614858565b61390c565b348015610ff057600080fd5b506106156040518060400160405280600b81526020017f626c6f636b4e756d62657200000000000000000000000000000000000000000081525081565b34801561103957600080fd5b5061056f611048366004614ce6565b6019805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b34801561109857600080fd5b506106f56110a73660046147a1565b6139a4565b3480156110b857600080fd5b5061056f6110c73660046149f3565b613a0d565b3480156110d857600080fd5b506105e76110e73660046147a1565b613ad8565b3480156110f857600080fd5b506105e7611107366004614858565b60096020526000908152604090205481565b34801561112557600080fd5b506105e760185481565b34801561113b57600080fd5b5061056f61114a366004614965565b613b09565b34801561115b57600080fd5b506105e761116a3660046147a1565b613b1d565b34801561117b57600080fd5b506105e761118a366004614858565b60066020526000908152604090205481565b3480156111a857600080fd5b506106f56111b73660046148f2565b613b39565b3480156111c857600080fd5b506105e76111d7366004614858565b60026020526000908152604090205481565b80516111fc90601a9060208401906143b1565b5050565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601654601554919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b39216906112e6908c1688614d30565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af1158015611364573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113889190614d74565b5060008860ff1667ffffffffffffffff8111156113a7576113a7614491565b6040519080825280602002602001820160405280156113d0578160200160208202803683370190505b50905060005b8960ff168160ff1610156115795760006113ef84613bad565b90508860ff16600103611527576040517f0d4a4fb1000000000000000000000000000000000000000000000000000000008152600481018290526000903090630d4a4fb190602401600060405180830381865afa158015611454573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261149a9190810190614ddc565b6017546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906114f39085908590600401614e11565b600060405180830381600087803b15801561150d57600080fd5b505af1158015611521573d6000803e3d6000fd5b50505050505b80838360ff168151811061153d5761153d614e2a565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061157181614e59565b9150506113d6565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711816040516115a991906147fe565b60405180910390a1505050505050505050565b606060006115ca6013613c99565b9050808410611605576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8260000361161a576116178482614e78565b92505b60008367ffffffffffffffff81111561163557611635614491565b60405190808252806020026020018201604052801561165e578160200160208202803683370190505b50905060005b848110156116b0576116816116798288614e8b565b601390613ca3565b82828151811061169357611693614e2a565b6020908102919091010152806116a881614e9e565b915050611664565b509150505b92915050565b600e60205282600052604060002060205281600052604060002081815481106116e357600080fd5b9060005260206000200160009250925050505481565b606060006040518060c001604052803073ffffffffffffffffffffffffffffffffffffffff168152602001600160ff1681526020017f8d98eacef480ad8f47c29266a1194f1874fdb68bcc98624964400d6ce72e69ec60001b81526020018460405160200161176a91815260200190565b60405160208183030381529060405261178290614ed6565b81526020016000801b81526020016000801b8152509050806040516020016118029190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b604051602081830303815290604052915050919050565b60608060006118286013613c99565b905060008167ffffffffffffffff81111561184557611845614491565b60405190808252806020026020018201604052801561186e578160200160208202803683370190505b50905060008267ffffffffffffffff81111561188c5761188c614491565b6040519080825280602002602001820160405280156118b5578160200160208202803683370190505b50905060005b83811015611b075760006118d0601383613ca3565b9050808483815181106118e5576118e5614e2a565b6020908102919091018101919091526000828152601290915260408082205490517f3ebe8d6c0000000000000000000000000000000000000000000000000000000081526004810184905261ffff90911691903090633ebe8d6c90602401602060405180830381865afa158015611960573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119849190614f1b565b905060008167ffffffffffffffff8111156119a1576119a1614491565b6040519080825280602002602001820160405280156119ca578160200160208202803683370190505b506000858152600e6020526040812091925090815b8561ffff168161ffff1611611ac45761ffff811660009081526020838152604080832080548251818502810185019093528083529192909190830182828015611a4757602002820191906000526020600020905b815481526020019060010190808311611a33575b5050505050905060005b8151811015611aaf57818181518110611a6c57611a6c614e2a565b6020026020010151868680611a8090614e9e565b975081518110611a9257611a92614e2a565b602090810291909101015280611aa781614e9e565b915050611a51565b50508080611abc90614f34565b9150506119df565b50611ad0838e86613caf565b888881518110611ae257611ae2614e2a565b6020026020010181815250505050505050508080611aff90614e9e565b9150506118bb565b50909590945092505050565b6000828152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938493929190830182828015611b7757602002820191906000526020600020905b815481526020019060010190808311611b63575b50505050509050611b89818251613e0e565b92509250505b9250929050565b601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611c2c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c509190614f60565b50601780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601554604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015611cf3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d179190614f8e565b601680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601554604051911681527f6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014906020015b60405180910390a15050565b6000838152600d602090815260408083208054825181850281018501909352808352611df593830182828015611de957602002820191906000526020600020905b815481526020019060010190808311611dd5575b50505050508484613caf565b90505b9392505050565b8060005b818160ff161015611e93573063c8048022858560ff8516818110611e2957611e29614e2a565b905060200201356040518263ffffffff1660e01b8152600401611e4e91815260200190565b600060405180830381600087803b158015611e6857600080fd5b505af1158015611e7c573d6000803e3d6000fd5b505050508080611e8b90614e59565b915050611e03565b507fbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b8383604051611ec5929190614fab565b60405180910390a1505050565b60008181526012602052604081205461ffff1681805b8261ffff168161ffff1611611f33576000858152600e6020908152604080832061ffff85168452909152902054611f1f9083614e8b565b915080611f2b81614f34565b915050611ee8565b509392505050565b60005a9050600080611f4f84860186614aad565b91509150600081806020019051810190611f699190614f1b565b60008181526005602090815260408083205460049092528220549293509190611f90613e93565b905082600003611fcd576000848152600560209081526040808320849055601082528220805460018101825590835291204291015591508161222c565b600084815260036020526040812054611fe68484614e78565b611ff09190614e78565b6000868152601160209081526040808320546010909252909120805492935061ffff909116918290811061202657612026614e2a565b90600052602060002001547f000000000000000000000000000000000000000000000000000000000000000061ffff16426120619190614e78565b11156120d0576000868152601060209081526040822080546001810182559083529120429101558061209281614f34565b600088815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559150505b600086815260126020908152604080832054600e835281842061ffff909116808552908352818420805483518186028101860190945280845291949390919083018282801561213e57602002820191906000526020600020905b81548152602001906001019080831161212a575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff168151036121b9578161217b81614f34565b60008a815260126020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600e6020908152604080832061ffff94851684528252808320805460018181018355918552838520018790558a8452600f83528184209590941683529381528382208054808501825590835281832001859055888252600d81529281208054928301815581529190912001555b600084815260066020526040812054612246906001614e8b565b6000868152600660209081526040918290208390558151878152908101859052908101859052606081018290529091507f6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede39060800160405180910390a160008581526004602090815260408083208590556018546002909252909120546122cd9084614e78565b111561255a576017546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810187905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612343573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612389919081019061502b565b6017546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810189905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa1580156123fe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612422919061515c565b6019549091506124469082906c01000000000000000000000000900460ff16614d30565b6bffffffffffffffffffffffff1682608001516bffffffffffffffffffffffff161015612557576019546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018990526bffffffffffffffffffffffff9091166024820152309063948108f790604401600060405180830381600087803b1580156124d757600080fd5b505af11580156124eb573d6000803e3d6000fd5b50505060008881526002602090815260409182902087905560195482518b81526bffffffffffffffffffffffff909116918101919091529081018690527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0915060600160405180910390a15b50505b6000858152600760205260409020545b805a612576908b614e78565b61258290612710614e8b565b10156125c35782406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905561256a565b82863273ffffffffffffffffffffffffffffffffffffffff167fcad583be2d908a590c81c7e332cf11c7a4ea41ecf1e059efac3ea7e83e34f1a58b60008151811061261057612610614e2a565b60200260200101518b604051612627929190615179565b60405180910390a45050505050505050505050565b6000828152600e6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561269e57602002820191906000526020600020905b81548152602001906001019080831161268a575b5050505050905092915050565b60006060600084846040516020016126c492919061519e565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b8160005b8181101561279c5730635f17e61686868481811061272357612723614e2a565b90506020020135856040518363ffffffff1660e01b815260040161275792919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561277157600080fd5b505af1158015612785573d6000803e3d6000fd5b50505050808061279490614e9e565b915050612703565b5050505050565b600b60205260009081526040902080546127bc90615229565b80601f01602080910402602001604051908101604052809291908181526020018280546127e890615229565b80156128355780601f1061280a57610100808354040283529160200191612835565b820191906000526020600020905b81548152906001019060200180831161281857829003601f168201915b505050505081565b612845613f35565b6016546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156128b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128d89190614f1b565b6016546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af1158015612950573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111fc9190614d74565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600d90915281206129ac91614407565b60008281526012602052604081205461ffff16905b8161ffff168161ffff1611612a08576000848152600e6020908152604080832061ffff8516845290915281206129f691614407565b80612a0081614f34565b9150506129c1565b5050600082815260126020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055601190915281205461ffff16905b8161ffff168161ffff1611612a96576000848152600f6020908152604080832061ffff851684529091528120612a8491614407565b80612a8e81614f34565b915050612a4f565b506000838152601060205260408120612aae91614407565b5050600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6040517f3ebe8d6c00000000000000000000000000000000000000000000000000000000815260048101839052600090819081903090633ebe8d6c90602401602060405180830381865afa158015612b42573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b669190614f1b565b9050831580612b755750808410155b15612b7e578093505b60008581526012602052604081205485919061ffff16805b6000898152600e6020908152604080832061ffff85168452825280832080548251818502810185019093528083529192909190830182828015612bf857602002820191906000526020600020905b815481526020019060010190808311612be4575b50505050509050600080612c0c8388613e0e565b9092509050612c1b8287614e8b565b9550612c278188614e78565b965060008711612c3957505050612c4f565b5050508080612c4790615276565b915050612b96565b5090979596505050505050565b6000606060005a90506000612c7385870187614858565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff811115612cac57612cac614491565b6040519080825280601f01601f191660200182016040528015612cd6576020820181803683370190505b50604051602001612ce8929190614e11565b60405160208183030381529060405290506000612d03613e93565b90506000612d1086612e99565b90505b835a612d1f9089614e78565b612d2b90612710614e8b565b1015612d6c5781406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612d13565b80612d84576000839850985050505050505050611b8f565b6040518060400160405280600981526020017f6665656449644865780000000000000000000000000000000000000000000000815250601a6040518060400160405280600b81526020017f626c6f636b4e756d6265720000000000000000000000000000000000000000008152508489604051602001612e0691815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e000000000000000000000000000000000000000000000000000000008252612e6895949392916004016152b2565b60405180910390fd5b600f60205282600052604060002060205281600052604060002081815481106116e357600080fd5b6000818152600560205260408120548103612eb657506001919050565b600082815260036020908152604080832054600490925290912054612ed9613e93565b612ee39190614e78565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612f6c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401612e68565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040517fa5f589340000000000000000000000000000000000000000000000000000000081526004810183905260009081908190309063a5f5893490602401602060405180830381865afa158015613044573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130689190614f1b565b90508315806130775750808410155b15613080578093505b60008581526011602052604081205485919061ffff16805b6000898152600f6020908152604080832061ffff851684528252808320805482518185028101850190935280835291929091908301828280156130fa57602002820191906000526020600020905b8154815260200190600101908083116130e6575b5050505050905060008061310e8388613e0e565b909250905061311d8287614e8b565b95506131298188614e78565b96506000871161313b57505050612c4f565b505050808061314990615276565b915050613098565b6017546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b5906131ab908690869086906004016153ff565b600060405180830381600087803b1580156131c557600080fd5b505af11580156131d9573d6000803e3d6000fd5b5050506000848152600b6020526040902090506131f782848361549e565b50505050565b6000838152600e6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849383018282801561325c57602002820191906000526020600020905b815481526020019060010190808311613248575b5050505050905061326f81858351613caf565b95945050505050565b6000838152600f6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849383018282801561325c5760200282019190600052602060002090815481526020019060010190808311613248575050505050905061326f81858351613caf565b60608060006132f76013613c99565b905060008167ffffffffffffffff81111561331457613314614491565b60405190808252806020026020018201604052801561333d578160200160208202803683370190505b50905060008267ffffffffffffffff81111561335b5761335b614491565b604051908082528060200260200182016040528015613384578160200160208202803683370190505b50905060005b83811015611b0757600061339f601383613ca3565b6000818152600d60209081526040808320805482518185028101850190935280835294955092939092918301828280156133f857602002820191906000526020600020905b8154815260200190600101908083116133e4575b505050505090508185848151811061341257613412614e2a565b60200260200101818152505061342a818a8351613caf565b84848151811061343c5761343c614e2a565b6020026020010181815250505050808061345590614e9e565b91505061338a565b6016546017546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af11580156134e5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135099190614d74565b506017546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b15801561358a57600080fd5b505af115801561359e573d6000803e3d6000fd5b5050604080518581526bffffffffffffffffffffffff851660208201527f8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e9350019050611d88565b6000828152600f6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561269e576020028201919060005260206000209081548152602001906001019080831161268a575050505050905092915050565b6000818152600d60209081526040918290208054835181840281018401909452808452606093928301828280156136a957602002820191906000526020600020905b815481526020019060010190808311613695575b50505050509050919050565b601a81815481106136c557600080fd5b9060005260206000200160009150905080546127bc90615229565b60008181526011602052604081205461ffff1681805b8261ffff168161ffff1611611f33576000858152600f6020908152604080832061ffff8516845290915290205461372d9083614e8b565b91508061373981614f34565b9150506136f6565b6017546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156137b957600080fd5b505af11580156137cd573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6017546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561385e57600080fd5b505af115801561279c573d6000803e3d6000fd5b8060005b818163ffffffff1610156131f7573063af953a4a858563ffffffff85168181106138a2576138a2614e2a565b905060200201356040518263ffffffff1660e01b81526004016138c791815260200190565b600060405180830381600087803b1580156138e157600080fd5b505af11580156138f5573d6000803e3d6000fd5b505050508080613904906155b8565b915050613876565b6017546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b15801561397857600080fd5b505af115801561398c573d6000803e3d6000fd5b505050506111fc816013613fb890919063ffffffff16565b6000828152600d602090815260408083208054825181850281018501909352808352849384939291908301828280156139fc57602002820191906000526020600020905b8154815260200190600101908083116139e8575b50505050509050611b898185613e0e565b8060005b818110156131f7576000848483818110613a2d57613a2d614e2a565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc8283604051602001613a6691815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613a92929190614e11565b600060405180830381600087803b158015613aac57600080fd5b505af1158015613ac0573d6000803e3d6000fd5b50505050508080613ad090614e9e565b915050613a11565b600d6020528160005260406000208181548110613af457600080fd5b90600052602060002001600091509150505481565b613b11613f35565b613b1a81613fc4565b50565b60106020528160005260406000208181548110613af457600080fd5b6000828152600f6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938493929190830182828015611b775760200282019190600052602060002090815481526020019060010190808311611b635750505050509050611b89818251613e0e565b6015546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613c089086906004016155d1565b6020604051808303816000875af1158015613c27573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c4b9190614f1b565b9050613c586013826140b9565b5060608301516000828152600a6020908152604080832063ffffffff90941690935560c0860151600b90915291902090613c929082615723565b5092915050565b60006116b5825490565b6000611df883836140c5565b82516000908190831580613cc35750808410155b15613ccc578093505b60008467ffffffffffffffff811115613ce757613ce7614491565b604051908082528060200260200182016040528015613d10578160200160208202803683370190505b509050600092505b84831015613d7e57866001613d2d8585614e78565b613d379190614e78565b81518110613d4757613d47614e2a565b6020026020010151818481518110613d6157613d61614e2a565b602090810291909101015282613d7681614e9e565b935050613d18565b613d9781600060018451613d929190614e78565b6140ef565b85606403613dd0578060018251613dae9190614e78565b81518110613dbe57613dbe614e2a565b60200260200101519350505050611df8565b806064825188613de0919061583d565b613dea91906158a9565b81518110613dfa57613dfa614e2a565b602002602001015193505050509392505050565b815160009081908190841580613e245750808510155b15613e2d578094505b60008092505b85831015613e8957866001613e488585614e78565b613e529190614e78565b81518110613e6257613e62614e2a565b602002602001015181613e759190614e8b565b905082613e8181614e9e565b935050613e33565b9694955050505050565b60007f000000000000000000000000000000000000000000000000000000000000000015613f3057606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613f07573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613f2b9190614f1b565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613fb6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401612e68565b565b6000611df8838361426f565b3373ffffffffffffffffffffffffffffffffffffffff821603614043576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401612e68565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611df88383614362565b60008260000182815481106140dc576140dc614e2a565b9060005260206000200154905092915050565b81818082036140ff575050505050565b600085600261410e87876158bd565b61411891906158dd565b6141229087615945565b8151811061413257614132614e2a565b602002602001015190505b818313614241575b8086848151811061415857614158614e2a565b6020026020010151101561417857826141708161596d565b935050614145565b85828151811061418a5761418a614e2a565b60200260200101518110156141ab57816141a38161599e565b925050614178565b81831361423c578582815181106141c4576141c4614e2a565b60200260200101518684815181106141de576141de614e2a565b60200260200101518785815181106141f8576141f8614e2a565b6020026020010188858151811061421157614211614e2a565b6020908102919091010191909152528261422a8161596d565b93505081806142389061599e565b9250505b61413d565b81851215614254576142548686846140ef565b83831215614267576142678684866140ef565b505050505050565b60008181526001830160205260408120548015614358576000614293600183614e78565b85549091506000906142a790600190614e78565b905081811461430c5760008660000182815481106142c7576142c7614e2a565b90600052602060002001549050808760000184815481106142ea576142ea614e2a565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061431d5761431d6159f5565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506116b5565b60009150506116b5565b60008181526001830160205260408120546143a9575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556116b5565b5060006116b5565b8280548282559060005260206000209081019282156143f7579160200282015b828111156143f757825182906143e79082615723565b50916020019190600101906143d1565b50614403929150614425565b5090565b5080546000825590600052602060002090810190613b1a9190614442565b808211156144035760006144398282614457565b50600101614425565b5b808211156144035760008155600101614443565b50805461446390615229565b6000825580601f10614473575050565b601f016020900490600052602060002090810190613b1a9190614442565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610160810167ffffffffffffffff811182821017156144e4576144e4614491565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561453157614531614491565b604052919050565b600067ffffffffffffffff82111561455357614553614491565b5060051b60200190565b600067ffffffffffffffff82111561457757614577614491565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006145b66145b18461455d565b6144ea565b90508281528383830111156145ca57600080fd5b828260208301376000602084830101529392505050565b600060208083850312156145f457600080fd5b823567ffffffffffffffff8082111561460c57600080fd5b818501915085601f83011261462057600080fd5b813561462e6145b182614539565b81815260059190911b8301840190848101908883111561464d57600080fd5b8585015b8381101561469a578035858111156146695760008081fd5b8601603f81018b1361467b5760008081fd5b61468c8b89830135604084016145a3565b845250918601918601614651565b5098975050505050505050565b803560ff811681146146b857600080fd5b919050565b63ffffffff81168114613b1a57600080fd5b600082601f8301126146e057600080fd5b611df8838335602085016145a3565b6bffffffffffffffffffffffff81168114613b1a57600080fd5b600080600080600080600060e0888a03121561472457600080fd5b61472d886146a7565b9650602088013561473d816146bd565b955061474b604089016146a7565b9450606088013567ffffffffffffffff81111561476757600080fd5b6147738a828b016146cf565b9450506080880135614784816146ef565b9699959850939692959460a0840135945060c09093013592915050565b600080604083850312156147b457600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156147f3578151875295820195908201906001016147d7565b509495945050505050565b602081526000611df860208301846147c3565b803561ffff811681146146b857600080fd5b60008060006060848603121561483857600080fd5b8335925061484860208501614811565b9150604084013590509250925092565b60006020828403121561486a57600080fd5b5035919050565b60005b8381101561488c578181015183820152602001614874565b50506000910152565b600081518084526148ad816020860160208601614871565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611df86020830184614895565b6000806040838503121561490557600080fd5b8235915061491560208401614811565b90509250929050565b60408152600061493160408301856147c3565b828103602084015261326f81856147c3565b73ffffffffffffffffffffffffffffffffffffffff81168114613b1a57600080fd5b60006020828403121561497757600080fd5b8135611df881614943565b60008060006060848603121561499757600080fd5b505081359360208301359350604090920135919050565b60008083601f8401126149c057600080fd5b50813567ffffffffffffffff8111156149d857600080fd5b6020830191508360208260051b8501011115611b8f57600080fd5b60008060208385031215614a0657600080fd5b823567ffffffffffffffff811115614a1d57600080fd5b614a29858286016149ae565b90969095509350505050565b60008083601f840112614a4757600080fd5b50813567ffffffffffffffff811115614a5f57600080fd5b602083019150836020828501011115611b8f57600080fd5b60008060208385031215614a8a57600080fd5b823567ffffffffffffffff811115614aa157600080fd5b614a2985828601614a35565b60008060408385031215614ac057600080fd5b823567ffffffffffffffff80821115614ad857600080fd5b818501915085601f830112614aec57600080fd5b81356020614afc6145b183614539565b82815260059290921b84018101918181019089841115614b1b57600080fd5b8286015b84811015614b5357803586811115614b375760008081fd5b614b458c86838b01016146cf565b845250918301918301614b1f565b5096505086013592505080821115614b6a57600080fd5b50614b77858286016146cf565b9150509250929050565b8215158152604060208201526000611df56040830184614895565b600080600060408486031215614bb157600080fd5b833567ffffffffffffffff811115614bc857600080fd5b614bd4868287016149ae565b9094509250506020840135614be8816146bd565b809150509250925092565b600080600060408486031215614c0857600080fd5b83359250602084013567ffffffffffffffff811115614c2657600080fd5b614c3286828701614a35565b9497909650939450505050565b600080600060608486031215614c5457600080fd5b8335925060208401359150614c6b60408501614811565b90509250925092565b60008060408385031215614c8757600080fd5b823591506020830135614c99816146ef565b809150509250929050565b60008060408385031215614cb757600080fd5b823591506020830135614c99816146bd565b600060208284031215614cdb57600080fd5b8135611df8816146ef565b600060208284031215614cf857600080fd5b611df8826146a7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614d5b57614d5b614d01565b02949350505050565b805180151581146146b857600080fd5b600060208284031215614d8657600080fd5b611df882614d64565b600082601f830112614da057600080fd5b8151614dae6145b18261455d565b818152846020838601011115614dc357600080fd5b614dd4826020830160208701614871565b949350505050565b600060208284031215614dee57600080fd5b815167ffffffffffffffff811115614e0557600080fd5b614dd484828501614d8f565b828152604060208201526000611df56040830184614895565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614e6f57614e6f614d01565b60010192915050565b818103818111156116b5576116b5614d01565b808201808211156116b5576116b5614d01565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614ecf57614ecf614d01565b5060010190565b80516020808301519190811015614f15577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b600060208284031215614f2d57600080fd5b5051919050565b600061ffff808316818103614f4b57614f4b614d01565b6001019392505050565b80516146b881614943565b60008060408385031215614f7357600080fd5b8251614f7e81614943565b6020939093015192949293505050565b600060208284031215614fa057600080fd5b8151611df881614943565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115614fe457600080fd5b8260051b80856040850137919091016040019392505050565b80516146b8816146bd565b80516146b8816146ef565b805167ffffffffffffffff811681146146b857600080fd5b60006020828403121561503d57600080fd5b815167ffffffffffffffff8082111561505557600080fd5b90830190610160828603121561506a57600080fd5b6150726144c0565b61507b83614f55565b815261508960208401614f55565b602082015261509a60408401614ffd565b60408201526060830151828111156150b157600080fd5b6150bd87828601614d8f565b6060830152506150cf60808401615008565b60808201526150e060a08401614f55565b60a08201526150f160c08401615013565b60c082015261510260e08401614ffd565b60e0820152610100615115818501615008565b90820152610120615127848201614d64565b90820152610140838101518381111561513f57600080fd5b61514b88828701614d8f565b918301919091525095945050505050565b60006020828403121561516e57600080fd5b8151611df8816146ef565b60408152600061518c6040830185614895565b828103602084015261326f8185614895565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015615213577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552615201868351614895565b955093820193908201906001016151c7565b50508584038187015250505061326f8185614895565b600181811c9082168061523d57607f821691505b602082108103614f15577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600061ffff82168061528a5761528a614d01565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b60a0815260006152c560a0830188614895565b602083820381850152818854808452828401915060058382821b86010160008c8152858120815b858110156153bf577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe089850301875282825461532781615229565b808752600182811680156153425760018114615379576153a8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168d8a01528c8315158b1b8a010194506153a8565b8688528c8820885b848110156153a05781548f828d01015283820191508e81019050615381565b8a018e019550505b50998b0199929650505091909101906001016152ec565b50505087810360408901526153d4818c614895565b9550505050505084606084015282810360808401526153f38185614895565b98975050505050505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f82111561549957600081815260208120601f850160051c8101602086101561547a5750805b601f850160051c820191505b8181101561426757828155600101615486565b505050565b67ffffffffffffffff8311156154b6576154b6614491565b6154ca836154c48354615229565b83615453565b6000601f84116001811461551c57600085156154e65750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561279c565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b8281101561556b578685013582556020948501946001909201910161554b565b50868210156155a6577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b600063ffffffff808316818103614f4b57614f4b614d01565b60208152600082516101408060208501526155f0610160850183614895565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08086850301604087015261562c8483614895565b935060408701519150615657606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526156b88483614895565b935060e087015191506101008187860301818801526156d78584614895565b9450808801519250506101208187860301818801526156f68584614895565b94508088015192505050615719828601826bffffffffffffffffffffffff169052565b5090949350505050565b815167ffffffffffffffff81111561573d5761573d614491565b6157518161574b8454615229565b84615453565b602080601f8311600181146157a4576000841561576e5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555614267565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156157f1578886015182559484019460019091019084016157d2565b508582101561582d57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561587557615875614d01565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826158b8576158b861587a565b500490565b8181036000831280158383131683831282161715613c9257613c92614d01565b6000826158ec576158ec61587a565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561594057615940614d01565b500590565b808201828112600083128015821682158216171561596557615965614d01565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614ecf57614ecf614d01565b60007f800000000000000000000000000000000000000000000000000000000000000082036159cf576159cf614d01565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307835353533343434333264353535333434326434313532343234393534353235353464326435343435353335343465343535343030303030303030303030303030",
}

var VerifiableLoadMercuryUpkeepABI = VerifiableLoadMercuryUpkeepMetaData.ABI

var VerifiableLoadMercuryUpkeepBin = VerifiableLoadMercuryUpkeepMetaData.Bin

func DeployVerifiableLoadMercuryUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, registrarAddress common.Address, useArb bool) (common.Address, *types.Transaction, *VerifiableLoadMercuryUpkeep, error) {
	parsed, err := VerifiableLoadMercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadMercuryUpkeepBin), backend, registrarAddress, useArb)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifiableLoadMercuryUpkeep{VerifiableLoadMercuryUpkeepCaller: VerifiableLoadMercuryUpkeepCaller{contract: contract}, VerifiableLoadMercuryUpkeepTransactor: VerifiableLoadMercuryUpkeepTransactor{contract: contract}, VerifiableLoadMercuryUpkeepFilterer: VerifiableLoadMercuryUpkeepFilterer{contract: contract}}, nil
}

type VerifiableLoadMercuryUpkeep struct {
	address common.Address
	abi     abi.ABI
	VerifiableLoadMercuryUpkeepCaller
	VerifiableLoadMercuryUpkeepTransactor
	VerifiableLoadMercuryUpkeepFilterer
}

type VerifiableLoadMercuryUpkeepCaller struct {
	contract *bind.BoundContract
}

type VerifiableLoadMercuryUpkeepTransactor struct {
	contract *bind.BoundContract
}

type VerifiableLoadMercuryUpkeepFilterer struct {
	contract *bind.BoundContract
}

type VerifiableLoadMercuryUpkeepSession struct {
	Contract     *VerifiableLoadMercuryUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VerifiableLoadMercuryUpkeepCallerSession struct {
	Contract *VerifiableLoadMercuryUpkeepCaller
	CallOpts bind.CallOpts
}

type VerifiableLoadMercuryUpkeepTransactorSession struct {
	Contract     *VerifiableLoadMercuryUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type VerifiableLoadMercuryUpkeepRaw struct {
	Contract *VerifiableLoadMercuryUpkeep
}

type VerifiableLoadMercuryUpkeepCallerRaw struct {
	Contract *VerifiableLoadMercuryUpkeepCaller
}

type VerifiableLoadMercuryUpkeepTransactorRaw struct {
	Contract *VerifiableLoadMercuryUpkeepTransactor
}

func NewVerifiableLoadMercuryUpkeep(address common.Address, backend bind.ContractBackend) (*VerifiableLoadMercuryUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(VerifiableLoadMercuryUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVerifiableLoadMercuryUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeep{address: address, abi: abi, VerifiableLoadMercuryUpkeepCaller: VerifiableLoadMercuryUpkeepCaller{contract: contract}, VerifiableLoadMercuryUpkeepTransactor: VerifiableLoadMercuryUpkeepTransactor{contract: contract}, VerifiableLoadMercuryUpkeepFilterer: VerifiableLoadMercuryUpkeepFilterer{contract: contract}}, nil
}

func NewVerifiableLoadMercuryUpkeepCaller(address common.Address, caller bind.ContractCaller) (*VerifiableLoadMercuryUpkeepCaller, error) {
	contract, err := bindVerifiableLoadMercuryUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepCaller{contract: contract}, nil
}

func NewVerifiableLoadMercuryUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifiableLoadMercuryUpkeepTransactor, error) {
	contract, err := bindVerifiableLoadMercuryUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepTransactor{contract: contract}, nil
}

func NewVerifiableLoadMercuryUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifiableLoadMercuryUpkeepFilterer, error) {
	contract, err := bindVerifiableLoadMercuryUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepFilterer{contract: contract}, nil
}

func bindVerifiableLoadMercuryUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifiableLoadMercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadMercuryUpkeep.Contract.VerifiableLoadMercuryUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.VerifiableLoadMercuryUpkeepTransactor.contract.Transfer(opts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.VerifiableLoadMercuryUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadMercuryUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.contract.Transfer(opts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) BUCKETSIZE(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "BUCKET_SIZE")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) TIMESTAMPINTERVAL(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "TIMESTAMP_INTERVAL")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) TIMESTAMPINTERVAL() (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TIMESTAMPINTERVAL(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) TIMESTAMPINTERVAL() (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TIMESTAMPINTERVAL(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) AddLinkAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "addLinkAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.AddLinkAmount(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.AddLinkAmount(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "bucketedDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BucketedDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BucketedDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "buckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Buckets(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Buckets(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckCallback(&_VerifiableLoadMercuryUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckCallback(&_VerifiableLoadMercuryUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) CheckDatas(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "checkDatas", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) CheckDatas(arg0 *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckDatas(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) CheckDatas(arg0 *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckDatas(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "checkGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "counters", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Counters(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Counters(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "delays", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Delays(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Delays(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.DummyMap(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.DummyMap(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "eligible", upkeepId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Eligible(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Eligible(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) FeedParamKey() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FeedParamKey(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) FeedParamKey() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FeedParamKey(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FeedsHex(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FeedsHex(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "firstPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "gasLimits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GasLimits(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GasLimits(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetActiveUpkeepIDs(&_VerifiableLoadMercuryUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetActiveUpkeepIDs(&_VerifiableLoadMercuryUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getBucketedDelays", upkeepId, bucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getDelays", upkeepId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelaysLength(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelaysLength(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetDelaysLengthAtBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getDelaysLengthAtBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetDelaysLengthAtBucket(upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelaysLengthAtBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetDelaysLengthAtBucket(upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelaysLengthAtBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetDelaysLengthAtTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getDelaysLengthAtTimestampBucket", upkeepId, timestampBucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetDelaysLengthAtTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelaysLengthAtTimestampBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetDelaysLengthAtTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetDelaysLengthAtTimestampBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getLogTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetLogTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetLogTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetPxBucketedDelaysForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getPxBucketedDelaysForAllUpkeeps", p)

	if err != nil {
		return *new([]*big.Int), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetPxBucketedDelaysForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxBucketedDelaysForAllUpkeeps(&_VerifiableLoadMercuryUpkeep.CallOpts, p)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetPxBucketedDelaysForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxBucketedDelaysForAllUpkeeps(&_VerifiableLoadMercuryUpkeep.CallOpts, p)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetPxDelayForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getPxDelayForAllUpkeeps", p)

	if err != nil {
		return *new([]*big.Int), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetPxDelayForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayForAllUpkeeps(&_VerifiableLoadMercuryUpkeep.CallOpts, p)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetPxDelayForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayForAllUpkeeps(&_VerifiableLoadMercuryUpkeep.CallOpts, p)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetPxDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getPxDelayInBucket", upkeepId, p, bucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetPxDelayInBucket(upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayInBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, p, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetPxDelayInBucket(upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayInBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, p, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetPxDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getPxDelayInTimestampBucket", upkeepId, p, timestampBucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetPxDelayInTimestampBucket(upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayInTimestampBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, p, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetPxDelayInTimestampBucket(upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayInTimestampBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, p, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getPxDelayLastNPerforms", upkeepId, p, n)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetSumBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getSumBucketedDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetSumBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumBucketedDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetSumBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumBucketedDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getSumDelayInBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetSumDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getSumDelayInTimestampBucket", upkeepId, timestampBucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetSumDelayInTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumDelayInTimestampBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetSumDelayInTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumDelayInTimestampBucket(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getSumDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetSumTimestampBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getSumTimestampBucketedDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetSumTimestampBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumTimestampBucketedDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetSumTimestampBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetSumTimestampBucketedDelayLastNPerforms(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetTimestampBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getTimestampBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetTimestampBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetTimestampBucketedDelaysLength(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetTimestampBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetTimestampBucketedDelaysLength(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetTimestampDelays(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getTimestampDelays", upkeepId, timestampBucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetTimestampDelays(upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetTimestampDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetTimestampDelays(upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetTimestampDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Intervals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "intervals", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Intervals(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Intervals(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) LastTopUpBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "lastTopUpBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.LinkToken(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.LinkToken(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) MinBalanceThresholdMultiplier(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "minBalanceThresholdMultiplier")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Owner() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Owner(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Owner() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Owner(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) PerformDataSizes(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "performDataSizes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PerformDataSizes(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PerformDataSizes(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) PerformGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "performGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) PreviousPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "previousPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Registrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "registrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Registrar() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Registrar(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Registrar() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Registrar(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Registry() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Registry(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Registry() (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Registry(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) TimeParamKey() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TimeParamKey(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) TimeParamKey() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TimeParamKey(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) TimestampBuckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "timestampBuckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) TimestampBuckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TimestampBuckets(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) TimestampBuckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TimestampBuckets(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) TimestampDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "timestampDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) TimestampDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TimestampDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) TimestampDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TimestampDelays(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Timestamps(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "timestamps", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Timestamps(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Timestamps(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Timestamps(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Timestamps(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) UpkeepTopUpCheckInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "upkeepTopUpCheckInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "useArbitrumBlockNum")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "acceptOwnership")
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.AcceptOwnership(&_VerifiableLoadMercuryUpkeep.TransactOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.AcceptOwnership(&_VerifiableLoadMercuryUpkeep.TransactOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) AddFunds(opts *bind.TransactOpts, upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "addFunds", upkeepId, amount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.AddFunds(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.AddFunds(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchCancelUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchCancelUpkeeps", upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchRegisterUpkeeps", number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchSetIntervals", upkeepIds, interval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchUpdatePipelineData", upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchWithdrawLinks", upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "cancelUpkeep", upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CancelUpkeep(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CancelUpkeep(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) CheckUpkeep(opts *bind.TransactOpts, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "checkUpkeep", checkData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckUpkeep(&_VerifiableLoadMercuryUpkeep.TransactOpts, checkData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.CheckUpkeep(&_VerifiableLoadMercuryUpkeep.TransactOpts, checkData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PerformUpkeep(&_VerifiableLoadMercuryUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.PerformUpkeep(&_VerifiableLoadMercuryUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetAddLinkAmount(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setAddLinkAmount", amount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetAddLinkAmount(amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetAddLinkAmount(&_VerifiableLoadMercuryUpkeep.TransactOpts, amount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetAddLinkAmount(amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetAddLinkAmount(&_VerifiableLoadMercuryUpkeep.TransactOpts, amount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetCheckGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setCheckGasToBurn", upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetCheckGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetCheckGasToBurn(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetCheckGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetCheckGasToBurn(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setConfig", newRegistrar)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetConfig(&_VerifiableLoadMercuryUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetConfig(&_VerifiableLoadMercuryUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setFeedsHex", newFeeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeedsHex(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeedsHex(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setInterval", upkeepId, _interval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetInterval(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetInterval(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetMinBalanceThresholdMultiplier(opts *bind.TransactOpts, newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setMinBalanceThresholdMultiplier", newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetMinBalanceThresholdMultiplier(newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetMinBalanceThresholdMultiplier(&_VerifiableLoadMercuryUpkeep.TransactOpts, newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetMinBalanceThresholdMultiplier(newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetMinBalanceThresholdMultiplier(&_VerifiableLoadMercuryUpkeep.TransactOpts, newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setPerformDataSize", upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetPerformGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setPerformGasToBurn", upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetPerformGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetPerformGasToBurn(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetPerformGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetPerformGasToBurn(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setUpkeepGasLimit", upkeepId, gasLimit)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetUpkeepTopUpCheckInterval(opts *bind.TransactOpts, newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setUpkeepTopUpCheckInterval", newInterval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetUpkeepTopUpCheckInterval(newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepTopUpCheckInterval(&_VerifiableLoadMercuryUpkeep.TransactOpts, newInterval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetUpkeepTopUpCheckInterval(newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepTopUpCheckInterval(&_VerifiableLoadMercuryUpkeep.TransactOpts, newInterval)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "transferOwnership", to)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TransferOwnership(&_VerifiableLoadMercuryUpkeep.TransactOpts, to)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TransferOwnership(&_VerifiableLoadMercuryUpkeep.TransactOpts, to)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "updateUpkeepPipelineData", upkeepId, pipelineData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "withdrawLinks")
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.WithdrawLinks(&_VerifiableLoadMercuryUpkeep.TransactOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.WithdrawLinks(&_VerifiableLoadMercuryUpkeep.TransactOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "withdrawLinks0", upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.RawTransact(opts, nil)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Receive(&_VerifiableLoadMercuryUpkeep.TransactOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Receive(&_VerifiableLoadMercuryUpkeep.TransactOpts)
}

type VerifiableLoadMercuryUpkeepFundsAddedIterator struct {
	Event *VerifiableLoadMercuryUpkeepFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepFundsAdded)
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
		it.Event = new(VerifiableLoadMercuryUpkeepFundsAdded)
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

func (it *VerifiableLoadMercuryUpkeepFundsAddedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepFundsAdded struct {
	UpkeepId *big.Int
	Amount   *big.Int
	Raw      types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepFundsAddedIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepFundsAddedIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepFundsAdded) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepFundsAdded)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseFundsAdded(log types.Log) (*VerifiableLoadMercuryUpkeepFundsAdded, error) {
	event := new(VerifiableLoadMercuryUpkeepFundsAdded)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepInsufficientFundsIterator struct {
	Event *VerifiableLoadMercuryUpkeepInsufficientFunds

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepInsufficientFundsIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepInsufficientFunds)
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
		it.Event = new(VerifiableLoadMercuryUpkeepInsufficientFunds)
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

func (it *VerifiableLoadMercuryUpkeepInsufficientFundsIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepInsufficientFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepInsufficientFunds struct {
	Balance  *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepInsufficientFundsIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "InsufficientFunds")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepInsufficientFundsIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "InsufficientFunds", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepInsufficientFunds) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "InsufficientFunds")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepInsufficientFunds)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "InsufficientFunds", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseInsufficientFunds(log types.Log) (*VerifiableLoadMercuryUpkeepInsufficientFunds, error) {
	event := new(VerifiableLoadMercuryUpkeepInsufficientFunds)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "InsufficientFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepLogEmittedIterator struct {
	Event *VerifiableLoadMercuryUpkeepLogEmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepLogEmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepLogEmitted)
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
		it.Event = new(VerifiableLoadMercuryUpkeepLogEmitted)
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

func (it *VerifiableLoadMercuryUpkeepLogEmittedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepLogEmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepLogEmitted struct {
	UpkeepId    *big.Int
	LogBlockNum *big.Int
	BlockNum    *big.Int
	Raw         types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int) (*VerifiableLoadMercuryUpkeepLogEmittedIterator, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "LogEmitted", upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepLogEmittedIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "LogEmitted", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepLogEmitted, upkeepId []*big.Int) (event.Subscription, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "LogEmitted", upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepLogEmitted)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseLogEmitted(log types.Log) (*VerifiableLoadMercuryUpkeepLogEmitted, error) {
	event := new(VerifiableLoadMercuryUpkeepLogEmitted)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepMercuryPerformEventIterator struct {
	Event *VerifiableLoadMercuryUpkeepMercuryPerformEvent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepMercuryPerformEventIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepMercuryPerformEvent)
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
		it.Event = new(VerifiableLoadMercuryUpkeepMercuryPerformEvent)
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

func (it *VerifiableLoadMercuryUpkeepMercuryPerformEventIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepMercuryPerformEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepMercuryPerformEvent struct {
	Origin      common.Address
	UpkeepId    *big.Int
	BlockNumber *big.Int
	V0          []byte
	Ed          []byte
	Raw         types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterMercuryPerformEvent(opts *bind.FilterOpts, origin []common.Address, upkeepId []*big.Int, blockNumber []*big.Int) (*VerifiableLoadMercuryUpkeepMercuryPerformEventIterator, error) {

	var originRule []interface{}
	for _, originItem := range origin {
		originRule = append(originRule, originItem)
	}
	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "MercuryPerformEvent", originRule, upkeepIdRule, blockNumberRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepMercuryPerformEventIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "MercuryPerformEvent", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepMercuryPerformEvent, origin []common.Address, upkeepId []*big.Int, blockNumber []*big.Int) (event.Subscription, error) {

	var originRule []interface{}
	for _, originItem := range origin {
		originRule = append(originRule, originItem)
	}
	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "MercuryPerformEvent", originRule, upkeepIdRule, blockNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepMercuryPerformEvent)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "MercuryPerformEvent", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseMercuryPerformEvent(log types.Log) (*VerifiableLoadMercuryUpkeepMercuryPerformEvent, error) {
	event := new(VerifiableLoadMercuryUpkeepMercuryPerformEvent)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "MercuryPerformEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator struct {
	Event *VerifiableLoadMercuryUpkeepOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepOwnershipTransferRequested)
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
		it.Event = new(VerifiableLoadMercuryUpkeepOwnershipTransferRequested)
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

func (it *VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepOwnershipTransferRequested)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadMercuryUpkeepOwnershipTransferRequested, error) {
	event := new(VerifiableLoadMercuryUpkeepOwnershipTransferRequested)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepOwnershipTransferredIterator struct {
	Event *VerifiableLoadMercuryUpkeepOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepOwnershipTransferred)
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
		it.Event = new(VerifiableLoadMercuryUpkeepOwnershipTransferred)
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

func (it *VerifiableLoadMercuryUpkeepOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadMercuryUpkeepOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepOwnershipTransferredIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepOwnershipTransferred)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseOwnershipTransferred(log types.Log) (*VerifiableLoadMercuryUpkeepOwnershipTransferred, error) {
	event := new(VerifiableLoadMercuryUpkeepOwnershipTransferred)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepPerformingUpkeepIterator struct {
	Event *VerifiableLoadMercuryUpkeepPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepPerformingUpkeep)
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
		it.Event = new(VerifiableLoadMercuryUpkeepPerformingUpkeep)
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

func (it *VerifiableLoadMercuryUpkeepPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepPerformingUpkeep struct {
	FirstPerformBlock *big.Int
	LastBlock         *big.Int
	PreviousBlock     *big.Int
	Counter           *big.Int
	Raw               types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepPerformingUpkeepIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepPerformingUpkeepIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepPerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepPerformingUpkeep)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParsePerformingUpkeep(log types.Log) (*VerifiableLoadMercuryUpkeepPerformingUpkeep, error) {
	event := new(VerifiableLoadMercuryUpkeepPerformingUpkeep)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepReceivedIterator struct {
	Event *VerifiableLoadMercuryUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepReceived)
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
		it.Event = new(VerifiableLoadMercuryUpkeepReceived)
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

func (it *VerifiableLoadMercuryUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepReceived struct {
	Sender common.Address
	Value  *big.Int
	Raw    types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterReceived(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepReceivedIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepReceivedIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "Received", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchReceived(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepReceived) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepReceived)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "Received", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseReceived(log types.Log) (*VerifiableLoadMercuryUpkeepReceived, error) {
	event := new(VerifiableLoadMercuryUpkeepReceived)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "Received", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepRegistrarSetIterator struct {
	Event *VerifiableLoadMercuryUpkeepRegistrarSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepRegistrarSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepRegistrarSet)
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
		it.Event = new(VerifiableLoadMercuryUpkeepRegistrarSet)
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

func (it *VerifiableLoadMercuryUpkeepRegistrarSetIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepRegistrarSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepRegistrarSet struct {
	NewRegistrar common.Address
	Raw          types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterRegistrarSet(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepRegistrarSetIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "RegistrarSet")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepRegistrarSetIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "RegistrarSet", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchRegistrarSet(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepRegistrarSet) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "RegistrarSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepRegistrarSet)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "RegistrarSet", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseRegistrarSet(log types.Log) (*VerifiableLoadMercuryUpkeepRegistrarSet, error) {
	event := new(VerifiableLoadMercuryUpkeepRegistrarSet)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "RegistrarSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepUpkeepTopUpIterator struct {
	Event *VerifiableLoadMercuryUpkeepUpkeepTopUp

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepUpkeepTopUpIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepUpkeepTopUp)
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
		it.Event = new(VerifiableLoadMercuryUpkeepUpkeepTopUp)
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

func (it *VerifiableLoadMercuryUpkeepUpkeepTopUpIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepUpkeepTopUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepUpkeepTopUp struct {
	UpkeepId *big.Int
	Amount   *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepUpkeepTopUpIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepUpkeepTopUpIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "UpkeepTopUp", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepUpkeepTopUp) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepUpkeepTopUp)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseUpkeepTopUp(log types.Log) (*VerifiableLoadMercuryUpkeepUpkeepTopUp, error) {
	event := new(VerifiableLoadMercuryUpkeepUpkeepTopUp)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepUpkeepsCancelledIterator struct {
	Event *VerifiableLoadMercuryUpkeepUpkeepsCancelled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepUpkeepsCancelledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepUpkeepsCancelled)
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
		it.Event = new(VerifiableLoadMercuryUpkeepUpkeepsCancelled)
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

func (it *VerifiableLoadMercuryUpkeepUpkeepsCancelledIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepUpkeepsCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepUpkeepsCancelled struct {
	UpkeepIds []*big.Int
	Raw       types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterUpkeepsCancelled(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepUpkeepsCancelledIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "UpkeepsCancelled")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepUpkeepsCancelledIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "UpkeepsCancelled", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchUpkeepsCancelled(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepUpkeepsCancelled) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "UpkeepsCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepUpkeepsCancelled)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "UpkeepsCancelled", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseUpkeepsCancelled(log types.Log) (*VerifiableLoadMercuryUpkeepUpkeepsCancelled, error) {
	event := new(VerifiableLoadMercuryUpkeepUpkeepsCancelled)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "UpkeepsCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadMercuryUpkeepUpkeepsRegisteredIterator struct {
	Event *VerifiableLoadMercuryUpkeepUpkeepsRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepUpkeepsRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepUpkeepsRegistered)
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
		it.Event = new(VerifiableLoadMercuryUpkeepUpkeepsRegistered)
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

func (it *VerifiableLoadMercuryUpkeepUpkeepsRegisteredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepUpkeepsRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepUpkeepsRegistered struct {
	UpkeepIds []*big.Int
	Raw       types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterUpkeepsRegistered(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepUpkeepsRegisteredIterator, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "UpkeepsRegistered")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepUpkeepsRegisteredIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "UpkeepsRegistered", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchUpkeepsRegistered(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepUpkeepsRegistered) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "UpkeepsRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepUpkeepsRegistered)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "UpkeepsRegistered", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseUpkeepsRegistered(log types.Log) (*VerifiableLoadMercuryUpkeepUpkeepsRegistered, error) {
	event := new(VerifiableLoadMercuryUpkeepUpkeepsRegistered)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "UpkeepsRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifiableLoadMercuryUpkeep.abi.Events["FundsAdded"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseFundsAdded(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["InsufficientFunds"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseInsufficientFunds(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["MercuryPerformEvent"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseMercuryPerformEvent(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["PerformingUpkeep"].ID:
		return _VerifiableLoadMercuryUpkeep.ParsePerformingUpkeep(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["Received"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseReceived(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["RegistrarSet"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseRegistrarSet(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseUpkeepTopUp(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["UpkeepsCancelled"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseUpkeepsCancelled(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["UpkeepsRegistered"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseUpkeepsRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadMercuryUpkeepFundsAdded) Topic() common.Hash {
	return common.HexToHash("0x8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e")
}

func (VerifiableLoadMercuryUpkeepInsufficientFunds) Topic() common.Hash {
	return common.HexToHash("0x03eb8b54a949acec2cd08fdb6d6bd4647a1f2c907d75d6900648effa92eb147f")
}

func (VerifiableLoadMercuryUpkeepLogEmitted) Topic() common.Hash {
	return common.HexToHash("0x8d98eacef480ad8f47c29266a1194f1874fdb68bcc98624964400d6ce72e69ec")
}

func (VerifiableLoadMercuryUpkeepMercuryPerformEvent) Topic() common.Hash {
	return common.HexToHash("0xcad583be2d908a590c81c7e332cf11c7a4ea41ecf1e059efac3ea7e83e34f1a5")
}

func (VerifiableLoadMercuryUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifiableLoadMercuryUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifiableLoadMercuryUpkeepPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede3")
}

func (VerifiableLoadMercuryUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874")
}

func (VerifiableLoadMercuryUpkeepRegistrarSet) Topic() common.Hash {
	return common.HexToHash("0x6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014")
}

func (VerifiableLoadMercuryUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (VerifiableLoadMercuryUpkeepUpkeepsCancelled) Topic() common.Hash {
	return common.HexToHash("0xbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b")
}

func (VerifiableLoadMercuryUpkeepUpkeepsRegistered) Topic() common.Hash {
	return common.HexToHash("0x2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711")
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeep) Address() common.Address {
	return _VerifiableLoadMercuryUpkeep.address
}

type VerifiableLoadMercuryUpkeepInterface interface {
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

	BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error)

	BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, checkData []byte) (*types.Transaction, error)

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

	FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepFundsAdded) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*VerifiableLoadMercuryUpkeepFundsAdded, error)

	FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepInsufficientFundsIterator, error)

	WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepInsufficientFunds) (event.Subscription, error)

	ParseInsufficientFunds(log types.Log) (*VerifiableLoadMercuryUpkeepInsufficientFunds, error)

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int) (*VerifiableLoadMercuryUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepLogEmitted, upkeepId []*big.Int) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadMercuryUpkeepLogEmitted, error)

	FilterMercuryPerformEvent(opts *bind.FilterOpts, origin []common.Address, upkeepId []*big.Int, blockNumber []*big.Int) (*VerifiableLoadMercuryUpkeepMercuryPerformEventIterator, error)

	WatchMercuryPerformEvent(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepMercuryPerformEvent, origin []common.Address, upkeepId []*big.Int, blockNumber []*big.Int) (event.Subscription, error)

	ParseMercuryPerformEvent(log types.Log) (*VerifiableLoadMercuryUpkeepMercuryPerformEvent, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadMercuryUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadMercuryUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadMercuryUpkeepOwnershipTransferred, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepPerformingUpkeep) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*VerifiableLoadMercuryUpkeepPerformingUpkeep, error)

	FilterReceived(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepReceivedIterator, error)

	WatchReceived(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepReceived) (event.Subscription, error)

	ParseReceived(log types.Log) (*VerifiableLoadMercuryUpkeepReceived, error)

	FilterRegistrarSet(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepRegistrarSetIterator, error)

	WatchRegistrarSet(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepRegistrarSet) (event.Subscription, error)

	ParseRegistrarSet(log types.Log) (*VerifiableLoadMercuryUpkeepRegistrarSet, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadMercuryUpkeepUpkeepTopUp, error)

	FilterUpkeepsCancelled(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepUpkeepsCancelledIterator, error)

	WatchUpkeepsCancelled(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepUpkeepsCancelled) (event.Subscription, error)

	ParseUpkeepsCancelled(log types.Log) (*VerifiableLoadMercuryUpkeepUpkeepsCancelled, error)

	FilterUpkeepsRegistered(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepUpkeepsRegisteredIterator, error)

	WatchUpkeepsRegistered(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepUpkeepsRegistered) (event.Subscription, error)

	ParseUpkeepsRegistered(log types.Log) (*VerifiableLoadMercuryUpkeepUpkeepsRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
