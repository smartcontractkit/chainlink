// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifiable_load_upkeep_wrapper

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

var VerifiableLoadUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"registrarAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"logBlockNum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"firstPerformBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"RegistrarSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkDatas\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxBucketedDelaysForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxDelayForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumTimestampBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTimestampBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getTimestampDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampBuckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040526005601855601980546001600160681b0319166c140000000002c68af0bb140000179055606460a052610e1060c0523480156200004057600080fd5b506040516200567f3803806200567f8339810160408190526200006391620002f2565b81813380600081620000bc5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000ef57620000ef816200022e565b5050601580546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa1580156200014c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000172919062000335565b50601780546001600160a01b0319166001600160a01b038381169190911790915560155460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620001d8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001fe919062000366565b601680546001600160a01b0319166001600160a01b0392909216919091179055501515608052506200038d915050565b336001600160a01b03821603620002885760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000b3565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6001600160a01b0381168114620002ef57600080fd5b50565b600080604083850312156200030657600080fd5b82516200031381620002d9565b602084015190925080151581146200032a57600080fd5b809150509250929050565b600080604083850312156200034957600080fd5b82516200035681620002d9565b6020939093015192949293505050565b6000602082840312156200037957600080fd5b81516200038681620002d9565b9392505050565b60805160a05160c0516152ad620003d2600039600081816106b50152611e8901526000818161058d0152611f9d0152600081816109570152613b2e01526152ad6000f3fe6080604052600436106104845760003560e01c8063776898c81161025e578063a72aa27e11610143578063daee1aeb116100bb578063f2fde38b1161008a578063fba7ffa31161006f578063fba7ffa314610ff6578063fbfb4f7614611023578063fcdc1f631461104357600080fd5b8063f2fde38b14610fb6578063fb0ceb0414610fd657600080fd5b8063daee1aeb14610f33578063dbef701e14610f53578063e0114adb14610f73578063e455308314610fa057600080fd5b8063becde0e111610112578063c8048022116100f7578063c804802214610e94578063d355852814610eb4578063d6051a7214610f1357600080fd5b8063becde0e114610e1a578063c357f1f314610e3a57600080fd5b8063a72aa27e14610d6f578063a79c404314610d8f578063af953a4a14610dbc578063b0971e1a14610ddc57600080fd5b80638fcb3fba116101d65780639b429354116101a55780639d385eaa1161018a5780639d385eaa14610d0f578063a5f5893414610d2f578063a6c60d8914610d4f57600080fd5b80639b42935414610cb15780639b51fb0d14610cde57600080fd5b80638fcb3fba14610c08578063948108f714610c3557806399cc6b0b14610c555780639ac542eb14610c7557600080fd5b80637e7a46dc1161022d57806387dfa9001161021257806387dfa90014610b9d5780638bc7b77214610bbd5780638da5cb5b14610bdd57600080fd5b80637e7a46dc14610b5d5780638237831714610b7d57600080fd5b8063776898c814610adb57806379ba509714610afb5780637b10399914610b105780637e4087b814610b3d57600080fd5b806345d2ec1711610384578063642f6cef116102fc5780636e04ff0d116102cb5780637145f11b116102b05780637145f11b14610a5157806373644cce14610a815780637672130314610aae57600080fd5b80636e04ff0d14610a035780637137a70214610a3157600080fd5b8063642f6cef14610945578063643b34e91461098957806369cdbadb146109a957806369e9b773146109d657600080fd5b806358c52c04116103535780635f17e616116103385780635f17e616146108b657806360457ff5146108d6578063636092e81461090357600080fd5b806358c52c04146108815780635d4ee7f3146108a157600080fd5b806345d2ec17146107e757806346e7a63e1461080757806351c98be31461083457806357970e931461085457600080fd5b806320e3dbd4116104175780632b20e397116103e657806333774d1c116103cb57806333774d1c146107765780633ebe8d6c146107a75780634585e33b146107c757600080fd5b80632b20e397146106f7578063328ffd111461074957600080fd5b806320e3dbd41461066357806328c4b57b1461068357806329f0e496146106a35780632a9032d3146106d757600080fd5b806312c550271161045357806312c550271461057b578063177b0eb9146105c25780631bee008014610600578063206c32e81461062e57600080fd5b806306c1cc00146104c857806306e3b632146104ea578063077ac621146105205780630d4a4fb11461054e57600080fd5b366104c357604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b3480156104d457600080fd5b506104e86104e3366004614200565b611070565b005b3480156104f657600080fd5b5061050a610505366004614298565b61142c565b60405161051791906142f5565b60405180910390f35b34801561052c57600080fd5b5061054061053b36600461431a565b61152b565b604051908152602001610517565b34801561055a57600080fd5b5061056e61056936600461434f565b611569565b60405161051791906143d6565b34801561058757600080fd5b506105af7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610517565b3480156105ce57600080fd5b506105406105dd3660046143e9565b6000918252600f6020908152604080842061ffff93909316845291905290205490565b34801561060c57600080fd5b5061062061061b36600461434f565b611689565b604051610517929190614415565b34801561063a57600080fd5b5061064e6106493660046143e9565b611983565b60408051928352602083019190915201610517565b34801561066f57600080fd5b506104e861067e36600461445c565b611a06565b34801561068f57600080fd5b5061054061069e366004614479565b611c04565b3480156106af57600080fd5b506105af7f000000000000000000000000000000000000000000000000000000000000000081565b3480156106e357600080fd5b506104e86106f23660046144ea565b611c6f565b34801561070357600080fd5b506015546107249073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610517565b34801561075557600080fd5b5061054061076436600461434f565b60036020526000908152604090205481565b34801561078257600080fd5b506105af61079136600461434f565b60116020526000908152604090205461ffff1681565b3480156107b357600080fd5b506105406107c236600461434f565b611d42565b3480156107d357600080fd5b506104e86107e236600461456e565b611dab565b3480156107f357600080fd5b5061050a6108023660046143e9565b612431565b34801561081357600080fd5b5061054061082236600461434f565b600a6020526000908152604090205481565b34801561084057600080fd5b506104e861084f3660046145a4565b6124a0565b34801561086057600080fd5b506016546107249073ffffffffffffffffffffffffffffffffffffffff1681565b34801561088d57600080fd5b5061056e61089c36600461434f565b612544565b3480156108ad57600080fd5b506104e86125de565b3480156108c257600080fd5b506104e86108d1366004614298565b612719565b3480156108e257600080fd5b506105406108f136600461434f565b60076020526000908152604090205481565b34801561090f57600080fd5b50601954610928906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff9091168152602001610517565b34801561095157600080fd5b506109797f000000000000000000000000000000000000000000000000000000000000000081565b6040519015158152602001610517565b34801561099557600080fd5b5061064e6109a4366004614298565b61288b565b3480156109b557600080fd5b506105406109c436600461434f565b60086020526000908152604090205481565b3480156109e257600080fd5b506104e86109f1366004614298565b60009182526008602052604090912055565b348015610a0f57600080fd5b50610a23610a1e36600461456e565b612a01565b6040516105179291906145fb565b348015610a3d57600080fd5b50610540610a4c36600461431a565b612b2e565b348015610a5d57600080fd5b50610979610a6c36600461434f565b600c6020526000908152604090205460ff1681565b348015610a8d57600080fd5b50610540610a9c36600461434f565b6000908152600d602052604090205490565b348015610aba57600080fd5b50610540610ac936600461434f565b60046020526000908152604090205481565b348015610ae757600080fd5b50610979610af636600461434f565b612b56565b348015610b0757600080fd5b506104e8612ba8565b348015610b1c57600080fd5b506017546107249073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b4957600080fd5b5061064e610b58366004614298565b612caa565b348015610b6957600080fd5b506104e8610b78366004614616565b612e13565b348015610b8957600080fd5b50610540610b98366004614662565b612ebf565b348015610ba957600080fd5b50610540610bb8366004614662565b612f3a565b348015610bc957600080fd5b50610620610bd836600461434f565b612faa565b348015610be957600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610724565b348015610c1457600080fd5b50610540610c2336600461434f565b60056020526000908152604090205481565b348015610c4157600080fd5b506104e8610c50366004614697565b61311f565b348015610c6157600080fd5b5061050a610c703660046143e9565b6132a8565b348015610c8157600080fd5b50601954610c9f906c01000000000000000000000000900460ff1681565b60405160ff9091168152602001610517565b348015610cbd57600080fd5b506104e8610ccc366004614298565b60009182526009602052604090912055565b348015610cea57600080fd5b506105af610cf936600461434f565b60126020526000908152604090205461ffff1681565b348015610d1b57600080fd5b5061050a610d2a36600461434f565b613315565b348015610d3b57600080fd5b50610540610d4a36600461434f565b613377565b348015610d5b57600080fd5b506104e8610d6a36600461434f565b601855565b348015610d7b57600080fd5b506104e8610d8a3660046146c7565b6133d8565b348015610d9b57600080fd5b506104e8610daa366004614298565b60009182526007602052604090912055565b348015610dc857600080fd5b506104e8610dd736600461434f565b613483565b348015610de857600080fd5b50610540610df73660046143e9565b6000918252600e6020908152604080842061ffff93909316845291905290205490565b348015610e2657600080fd5b506104e8610e353660046144ea565b613509565b348015610e4657600080fd5b506104e8610e553660046146ec565b601980547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610ea057600080fd5b506104e8610eaf36600461434f565b6135a3565b348015610ec057600080fd5b506104e8610ecf366004614709565b6019805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b348015610f1f57600080fd5b5061064e610f2e366004614298565b61363b565b348015610f3f57600080fd5b506104e8610f4e3660046144ea565b6136a4565b348015610f5f57600080fd5b50610540610f6e366004614298565b61376f565b348015610f7f57600080fd5b50610540610f8e36600461434f565b60096020526000908152604090205481565b348015610fac57600080fd5b5061054060185481565b348015610fc257600080fd5b506104e8610fd136600461445c565b6137a0565b348015610fe257600080fd5b50610540610ff1366004614298565b6137b4565b34801561100257600080fd5b5061054061101136600461434f565b60066020526000908152604090205481565b34801561102f57600080fd5b5061064e61103e3660046143e9565b6137d0565b34801561104f57600080fd5b5061054061105e36600461434f565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601654601554919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b3921690611156908c1688614753565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af11580156111d4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111f89190614797565b5060008860ff1667ffffffffffffffff811115611217576112176140a2565b604051908082528060200260200182016040528015611240578160200160208202803683370190505b50905060005b8960ff168160ff1610156113e957600061125f84613844565b90508860ff16600103611397576040517f0d4a4fb1000000000000000000000000000000000000000000000000000000008152600481018290526000903090630d4a4fb190602401600060405180830381865afa1580156112c4573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261130a91908101906147ff565b6017546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906113639085908590600401614834565b600060405180830381600087803b15801561137d57600080fd5b505af1158015611391573d6000803e3d6000fd5b50505050505b80838360ff16815181106113ad576113ad61484d565b602090810291909101810191909152600091825260088152604080832088905560079091529020849055806113e18161487c565b915050611246565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c7118160405161141991906142f5565b60405180910390a1505050505050505050565b6060600061143a6013613930565b9050808410611475576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8260000361148a57611487848261489b565b92505b60008367ffffffffffffffff8111156114a5576114a56140a2565b6040519080825280602002602001820160405280156114ce578160200160208202803683370190505b50905060005b84811015611520576114f16114e982886148ae565b60139061393a565b8282815181106115035761150361484d565b602090810291909101015280611518816148c1565b9150506114d4565b509150505b92915050565b600e602052826000526040600020602052816000526040600020818154811061155357600080fd5b9060005260206000200160009250925050505481565b606060006040518060c001604052803073ffffffffffffffffffffffffffffffffffffffff168152602001600160ff1681526020017f8d98eacef480ad8f47c29266a1194f1874fdb68bcc98624964400d6ce72e69ec60001b8152602001846040516020016115da91815260200190565b6040516020818303038152906040526115f2906148f9565b81526020016000801b81526020016000801b8152509050806040516020016116729190600060c08201905073ffffffffffffffffffffffffffffffffffffffff835116825260ff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b604051602081830303815290604052915050919050565b60608060006116986013613930565b905060008167ffffffffffffffff8111156116b5576116b56140a2565b6040519080825280602002602001820160405280156116de578160200160208202803683370190505b50905060008267ffffffffffffffff8111156116fc576116fc6140a2565b604051908082528060200260200182016040528015611725578160200160208202803683370190505b50905060005b8381101561197757600061174060138361393a565b9050808483815181106117555761175561484d565b6020908102919091018101919091526000828152601290915260408082205490517f3ebe8d6c0000000000000000000000000000000000000000000000000000000081526004810184905261ffff90911691903090633ebe8d6c90602401602060405180830381865afa1580156117d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117f4919061493e565b905060008167ffffffffffffffff811115611811576118116140a2565b60405190808252806020026020018201604052801561183a578160200160208202803683370190505b506000858152600e6020526040812091925090815b8561ffff168161ffff16116119345761ffff8116600090815260208381526040808320805482518185028101850190935280835291929091908301828280156118b757602002820191906000526020600020905b8154815260200190600101908083116118a3575b5050505050905060005b815181101561191f578181815181106118dc576118dc61484d565b60200260200101518686806118f0906148c1565b9750815181106119025761190261484d565b602090810291909101015280611917816148c1565b9150506118c1565b5050808061192c90614957565b91505061184f565b50611940838e86613946565b8888815181106119525761195261484d565b602002602001018181525050505050505050808061196f906148c1565b91505061172b565b50909590945092505050565b6000828152600e6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156119e757602002820191906000526020600020905b8154815260200190600101908083116119d3575b505050505090506119f9818251613aa5565b92509250505b9250929050565b601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611a9c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ac09190614983565b50601780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601554604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015611b63573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b8791906149b1565b601680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601554604051911681527f6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014906020015b60405180910390a15050565b6000838152600d602090815260408083208054825181850281018501909352808352611c6593830182828015611c5957602002820191906000526020600020905b815481526020019060010190808311611c45575b50505050508484613946565b90505b9392505050565b8060005b818160ff161015611d03573063c8048022858560ff8516818110611c9957611c9961484d565b905060200201356040518263ffffffff1660e01b8152600401611cbe91815260200190565b600060405180830381600087803b158015611cd857600080fd5b505af1158015611cec573d6000803e3d6000fd5b505050508080611cfb9061487c565b915050611c73565b507fbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b8383604051611d359291906149ce565b60405180910390a1505050565b60008181526012602052604081205461ffff1681805b8261ffff168161ffff1611611da3576000858152600e6020908152604080832061ffff85168452909152902054611d8f90836148ae565b915080611d9b81614957565b915050611d58565b509392505050565b60005a90506000611dbe83850185614a20565b5060008181526005602090815260408083205460049092528220549293509190611de6613b2a565b905082600003611e235760008481526005602090815260408083208490556010825282208054600181018255908352912042910155915081612082565b600084815260036020526040812054611e3c848461489b565b611e46919061489b565b6000868152601160209081526040808320546010909252909120805492935061ffff9091169182908110611e7c57611e7c61484d565b90600052602060002001547f000000000000000000000000000000000000000000000000000000000000000061ffff1642611eb7919061489b565b1115611f265760008681526010602090815260408220805460018101825590835291204291015580611ee881614957565b600088815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559150505b600086815260126020908152604080832054600e835281842061ffff9091168085529083528184208054835181860281018601909452808452919493909190830182828015611f9457602002820191906000526020600020905b815481526020019060010190808311611f80575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff1681510361200f5781611fd181614957565b60008a815260126020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600e6020908152604080832061ffff94851684528252808320805460018181018355918552838520018790558a8452600f83528184209590941683529381528382208054808501825590835281832001859055888252600d81529281208054928301815581529190912001555b60008481526006602052604081205461209c9060016148ae565b6000868152600660209081526040918290208390558151878152908101859052908101859052606081018290529091507f6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede39060800160405180910390a16000858152600460209081526040808320859055601854600290925290912054612123908461489b565b11156123b0576017546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810187905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612199573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526121df9190810190614a95565b6017546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810189905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612254573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122789190614bc6565b60195490915061229c9082906c01000000000000000000000000900460ff16614753565b6bffffffffffffffffffffffff1682608001516bffffffffffffffffffffffff1610156123ad576019546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018990526bffffffffffffffffffffffff9091166024820152309063948108f790604401600060405180830381600087803b15801561232d57600080fd5b505af1158015612341573d6000803e3d6000fd5b50505060008881526002602090815260409182902087905560195482518b81526bffffffffffffffffffffffff909116918101919091529081018690527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0915060600160405180910390a15b50505b6000858152600760205260409020545b805a6123cc908961489b565b6123d8906127106148ae565b10156124265782406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558261241e81614be3565b9350506123c0565b505050505050505050565b6000828152600e6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561249357602002820191906000526020600020905b81548152602001906001019080831161247f575b5050505050905092915050565b8160005b8181101561253d5730635f17e6168686848181106124c4576124c461484d565b90506020020135856040518363ffffffff1660e01b81526004016124f892919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561251257600080fd5b505af1158015612526573d6000803e3d6000fd5b505050508080612535906148c1565b9150506124a4565b5050505050565b600b602052600090815260409020805461255d90614c18565b80601f016020809104026020016040519081016040528092919081815260200182805461258990614c18565b80156125d65780601f106125ab576101008083540402835291602001916125d6565b820191906000526020600020905b8154815290600101906020018083116125b957829003601f168201915b505050505081565b6125e6613bcc565b6016546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612655573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612679919061493e565b6016546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156126f1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127159190614797565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600d909152812061275191614048565b60008281526012602052604081205461ffff16905b8161ffff168161ffff16116127ad576000848152600e6020908152604080832061ffff85168452909152812061279b91614048565b806127a581614957565b915050612766565b5050600082815260126020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055601190915281205461ffff16905b8161ffff168161ffff161161283b576000848152600f6020908152604080832061ffff85168452909152812061282991614048565b8061283381614957565b9150506127f4565b50600083815260106020526040812061285391614048565b5050600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6040517f3ebe8d6c00000000000000000000000000000000000000000000000000000000815260048101839052600090819081903090633ebe8d6c90602401602060405180830381865afa1580156128e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061290b919061493e565b905083158061291a5750808410155b15612923578093505b60008581526012602052604081205485919061ffff16805b6000898152600e6020908152604080832061ffff8516845282528083208054825181850281018501909352808352919290919083018282801561299d57602002820191906000526020600020905b815481526020019060010190808311612989575b505050505090506000806129b18388613aa5565b90925090506129c082876148ae565b95506129cc818861489b565b9650600087116129de575050506129f4565b50505080806129ec90614c65565b91505061293b565b5090979596505050505050565b6000606060005a90506000612a188587018761434f565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff811115612a5157612a516140a2565b6040519080825280601f01601f191660200182016040528015612a7b576020820181803683370190505b50604051602001612a8d929190614834565b60405160208183030381529060405290506000612aa8613b2a565b90506000612ab586612b56565b90505b835a612ac4908961489b565b612ad0906127106148ae565b1015612b1e5781406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905581612b1681614be3565b925050612ab8565b9a91995090975050505050505050565b600f602052826000526040600020602052816000526040600020818154811061155357600080fd5b6000818152600560205260408120548103612b7357506001919050565b600082815260036020908152604080832054600490925290912054612b96613b2a565b612ba0919061489b565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612c2e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040517fa5f589340000000000000000000000000000000000000000000000000000000081526004810183905260009081908190309063a5f5893490602401602060405180830381865afa158015612d06573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d2a919061493e565b9050831580612d395750808410155b15612d42578093505b60008581526011602052604081205485919061ffff16805b6000898152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083529192909190830182828015612dbc57602002820191906000526020600020905b815481526020019060010190808311612da8575b50505050509050600080612dd08388613aa5565b9092509050612ddf82876148ae565b9550612deb818861489b565b965060008711612dfd575050506129f4565b5050508080612e0b90614c65565b915050612d5a565b6017546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b590612e6d90869086908690600401614ca1565b600060405180830381600087803b158015612e8757600080fd5b505af1158015612e9b573d6000803e3d6000fd5b5050506000848152600b602052604090209050612eb9828483614d40565b50505050565b6000838152600e6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493830182828015612f1e57602002820191906000526020600020905b815481526020019060010190808311612f0a575b50505050509050612f3181858351613946565b95945050505050565b6000838152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493830182828015612f1e5760200282019190600052602060002090815481526020019060010190808311612f0a5750505050509050612f3181858351613946565b6060806000612fb96013613930565b905060008167ffffffffffffffff811115612fd657612fd66140a2565b604051908082528060200260200182016040528015612fff578160200160208202803683370190505b50905060008267ffffffffffffffff81111561301d5761301d6140a2565b604051908082528060200260200182016040528015613046578160200160208202803683370190505b50905060005b8381101561197757600061306160138361393a565b6000818152600d60209081526040808320805482518185028101850190935280835294955092939092918301828280156130ba57602002820191906000526020600020905b8154815260200190600101908083116130a6575b50505050509050818584815181106130d4576130d461484d565b6020026020010181815250506130ec818a8351613946565b8484815181106130fe576130fe61484d565b60200260200101818152505050508080613117906148c1565b91505061304c565b6016546017546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af11580156131a7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131cb9190614797565b506017546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b15801561324c57600080fd5b505af1158015613260573d6000803e3d6000fd5b5050604080518581526bffffffffffffffffffffffff851660208201527f8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e9350019050611bf8565b6000828152600f6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015612493576020028201919060005260206000209081548152602001906001019080831161247f575050505050905092915050565b6000818152600d602090815260409182902080548351818402810184019094528084526060939283018282801561336b57602002820191906000526020600020905b815481526020019060010190808311613357575b50505050509050919050565b60008181526011602052604081205461ffff1681805b8261ffff168161ffff1611611da3576000858152600f6020908152604080832061ffff851684529091529020546133c490836148ae565b9150806133d081614957565b91505061338d565b6017546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b15801561345057600080fd5b505af1158015613464573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6017546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b1580156134f557600080fd5b505af115801561253d573d6000803e3d6000fd5b8060005b818163ffffffff161015612eb9573063af953a4a858563ffffffff85168181106135395761353961484d565b905060200201356040518263ffffffff1660e01b815260040161355e91815260200190565b600060405180830381600087803b15801561357857600080fd5b505af115801561358c573d6000803e3d6000fd5b50505050808061359b90614e5a565b91505061350d565b6017546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b15801561360f57600080fd5b505af1158015613623573d6000803e3d6000fd5b50505050612715816013613c4f90919063ffffffff16565b6000828152600d6020908152604080832080548251818502810185019093528083528493849392919083018282801561369357602002820191906000526020600020905b81548152602001906001019080831161367f575b505050505090506119f98185613aa5565b8060005b81811015612eb95760008484838181106136c4576136c461484d565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016136fd91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613729929190614834565b600060405180830381600087803b15801561374357600080fd5b505af1158015613757573d6000803e3d6000fd5b50505050508080613767906148c1565b9150506136a8565b600d602052816000526040600020818154811061378b57600080fd5b90600052602060002001600091509150505481565b6137a8613bcc565b6137b181613c5b565b50565b6010602052816000526040600020818154811061378b57600080fd5b6000828152600f6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156119e757602002820191906000526020600020908154815260200190600101908083116119d357505050505090506119f9818251613aa5565b6015546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e119061389f908690600401614e73565b6020604051808303816000875af11580156138be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138e2919061493e565b90506138ef601382613d50565b5060608301516000828152600a6020908152604080832063ffffffff90941690935560c0860151600b909152919020906139299082614fc5565b5092915050565b6000611525825490565b6000611c688383613d5c565b8251600090819083158061395a5750808410155b15613963578093505b60008467ffffffffffffffff81111561397e5761397e6140a2565b6040519080825280602002602001820160405280156139a7578160200160208202803683370190505b509050600092505b84831015613a15578660016139c4858561489b565b6139ce919061489b565b815181106139de576139de61484d565b60200260200101518184815181106139f8576139f861484d565b602090810291909101015282613a0d816148c1565b9350506139af565b613a2e81600060018451613a29919061489b565b613d86565b85606403613a67578060018251613a45919061489b565b81518110613a5557613a5561484d565b60200260200101519350505050611c68565b806064825188613a7791906150df565b613a81919061514b565b81518110613a9157613a9161484d565b602002602001015193505050509392505050565b815160009081908190841580613abb5750808510155b15613ac4578094505b60008092505b85831015613b2057866001613adf858561489b565b613ae9919061489b565b81518110613af957613af961484d565b602002602001015181613b0c91906148ae565b905082613b18816148c1565b935050613aca565b9694955050505050565b60007f000000000000000000000000000000000000000000000000000000000000000015613bc757606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613b9e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bc2919061493e565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613c4d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401612c25565b565b6000611c688383613f06565b3373ffffffffffffffffffffffffffffffffffffffff821603613cda576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401612c25565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611c688383613ff9565b6000826000018281548110613d7357613d7361484d565b9060005260206000200154905092915050565b8181808203613d96575050505050565b6000856002613da5878761515f565b613daf919061517f565b613db990876151e7565b81518110613dc957613dc961484d565b602002602001015190505b818313613ed8575b80868481518110613def57613def61484d565b60200260200101511015613e0f5782613e078161520f565b935050613ddc565b858281518110613e2157613e2161484d565b6020026020010151811015613e425781613e3a81615240565b925050613e0f565b818313613ed357858281518110613e5b57613e5b61484d565b6020026020010151868481518110613e7557613e7561484d565b6020026020010151878581518110613e8f57613e8f61484d565b60200260200101888581518110613ea857613ea861484d565b60209081029190910101919091525282613ec18161520f565b9350508180613ecf90615240565b9250505b613dd4565b81851215613eeb57613eeb868684613d86565b83831215613efe57613efe868486613d86565b505050505050565b60008181526001830160205260408120548015613fef576000613f2a60018361489b565b8554909150600090613f3e9060019061489b565b9050818114613fa3576000866000018281548110613f5e57613f5e61484d565b9060005260206000200154905080876000018481548110613f8157613f8161484d565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613fb457613fb4615271565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611525565b6000915050611525565b600081815260018301602052604081205461404057508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611525565b506000611525565b50805460008255906000526020600020908101906137b191905b808211156140765760008155600101614062565b5090565b803560ff8116811461408b57600080fd5b919050565b63ffffffff811681146137b157600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610160810167ffffffffffffffff811182821017156140f5576140f56140a2565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614142576141426140a2565b604052919050565b600067ffffffffffffffff821115614164576141646140a2565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f8301126141a157600080fd5b81356141b46141af8261414a565b6140fb565b8181528460208386010111156141c957600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff811681146137b157600080fd5b600080600080600080600060e0888a03121561421b57600080fd5b6142248861407a565b9650602088013561423481614090565b95506142426040890161407a565b9450606088013567ffffffffffffffff81111561425e57600080fd5b61426a8a828b01614190565b945050608088013561427b816141e6565b9699959850939692959460a0840135945060c09093013592915050565b600080604083850312156142ab57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156142ea578151875295820195908201906001016142ce565b509495945050505050565b602081526000611c6860208301846142ba565b803561ffff8116811461408b57600080fd5b60008060006060848603121561432f57600080fd5b8335925061433f60208501614308565b9150604084013590509250925092565b60006020828403121561436157600080fd5b5035919050565b60005b8381101561438357818101518382015260200161436b565b50506000910152565b600081518084526143a4816020860160208601614368565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611c68602083018461438c565b600080604083850312156143fc57600080fd5b8235915061440c60208401614308565b90509250929050565b60408152600061442860408301856142ba565b8281036020840152612f3181856142ba565b73ffffffffffffffffffffffffffffffffffffffff811681146137b157600080fd5b60006020828403121561446e57600080fd5b8135611c688161443a565b60008060006060848603121561448e57600080fd5b505081359360208301359350604090920135919050565b60008083601f8401126144b757600080fd5b50813567ffffffffffffffff8111156144cf57600080fd5b6020830191508360208260051b85010111156119ff57600080fd5b600080602083850312156144fd57600080fd5b823567ffffffffffffffff81111561451457600080fd5b614520858286016144a5565b90969095509350505050565b60008083601f84011261453e57600080fd5b50813567ffffffffffffffff81111561455657600080fd5b6020830191508360208285010111156119ff57600080fd5b6000806020838503121561458157600080fd5b823567ffffffffffffffff81111561459857600080fd5b6145208582860161452c565b6000806000604084860312156145b957600080fd5b833567ffffffffffffffff8111156145d057600080fd5b6145dc868287016144a5565b90945092505060208401356145f081614090565b809150509250925092565b8215158152604060208201526000611c65604083018461438c565b60008060006040848603121561462b57600080fd5b83359250602084013567ffffffffffffffff81111561464957600080fd5b6146558682870161452c565b9497909650939450505050565b60008060006060848603121561467757600080fd5b833592506020840135915061468e60408501614308565b90509250925092565b600080604083850312156146aa57600080fd5b8235915060208301356146bc816141e6565b809150509250929050565b600080604083850312156146da57600080fd5b8235915060208301356146bc81614090565b6000602082840312156146fe57600080fd5b8135611c68816141e6565b60006020828403121561471b57600080fd5b611c688261407a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff8083168185168183048111821515161561477e5761477e614724565b02949350505050565b8051801515811461408b57600080fd5b6000602082840312156147a957600080fd5b611c6882614787565b600082601f8301126147c357600080fd5b81516147d16141af8261414a565b8181528460208386010111156147e657600080fd5b6147f7826020830160208701614368565b949350505050565b60006020828403121561481157600080fd5b815167ffffffffffffffff81111561482857600080fd5b6147f7848285016147b2565b828152604060208201526000611c65604083018461438c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff810361489257614892614724565b60010192915050565b8181038181111561152557611525614724565b8082018082111561152557611525614724565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036148f2576148f2614724565b5060010190565b80516020808301519190811015614938577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8160200360031b1b821691505b50919050565b60006020828403121561495057600080fd5b5051919050565b600061ffff80831681810361496e5761496e614724565b6001019392505050565b805161408b8161443a565b6000806040838503121561499657600080fd5b82516149a18161443a565b6020939093015192949293505050565b6000602082840312156149c357600080fd5b8151611c688161443a565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115614a0757600080fd5b8260051b80856040850137919091016040019392505050565b60008060408385031215614a3357600080fd5b82359150602083013567ffffffffffffffff811115614a5157600080fd5b614a5d85828601614190565b9150509250929050565b805161408b81614090565b805161408b816141e6565b805167ffffffffffffffff8116811461408b57600080fd5b600060208284031215614aa757600080fd5b815167ffffffffffffffff80821115614abf57600080fd5b908301906101608286031215614ad457600080fd5b614adc6140d1565b614ae583614978565b8152614af360208401614978565b6020820152614b0460408401614a67565b6040820152606083015182811115614b1b57600080fd5b614b27878286016147b2565b606083015250614b3960808401614a72565b6080820152614b4a60a08401614978565b60a0820152614b5b60c08401614a7d565b60c0820152614b6c60e08401614a67565b60e0820152610100614b7f818501614a72565b90820152610120614b91848201614787565b908201526101408381015183811115614ba957600080fd5b614bb5888287016147b2565b918301919091525095945050505050565b600060208284031215614bd857600080fd5b8151611c68816141e6565b600081614bf257614bf2614724565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600181811c90821680614c2c57607f821691505b602082108103614938577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600061ffff821680614c7957614c79614724565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f821115614d3b57600081815260208120601f850160051c81016020861015614d1c5750805b601f850160051c820191505b81811015613efe57828155600101614d28565b505050565b67ffffffffffffffff831115614d5857614d586140a2565b614d6c83614d668354614c18565b83614cf5565b6000601f841160018114614dbe5760008515614d885750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b17835561253d565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b82811015614e0d5786850135825560209485019460019092019101614ded565b5086821015614e48577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b600063ffffffff80831681810361496e5761496e614724565b6020815260008251610140806020850152614e9261016085018361438c565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe080868503016040870152614ece848361438c565b935060408701519150614ef9606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e0870152614f5a848361438c565b935060e08701519150610100818786030181880152614f79858461438c565b945080880151925050610120818786030181880152614f98858461438c565b94508088015192505050614fbb828601826bffffffffffffffffffffffff169052565b5090949350505050565b815167ffffffffffffffff811115614fdf57614fdf6140a2565b614ff381614fed8454614c18565b84614cf5565b602080601f83116001811461504657600084156150105750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613efe565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561509357888601518255948401946001909101908401615074565b50858210156150cf57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561511757615117614724565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261515a5761515a61511c565b500490565b818103600083128015838313168383128216171561392957613929614724565b60008261518e5761518e61511c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f8000000000000000000000000000000000000000000000000000000000000000831416156151e2576151e2614724565b500590565b808201828112600083128015821682158216171561520757615207614724565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036148f2576148f2614724565b60007f80000000000000000000000000000000000000000000000000000000000000008203614bf257614bf2614724565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a",
}

var VerifiableLoadUpkeepABI = VerifiableLoadUpkeepMetaData.ABI

var VerifiableLoadUpkeepBin = VerifiableLoadUpkeepMetaData.Bin

func DeployVerifiableLoadUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, registrarAddress common.Address, useArb bool) (common.Address, *types.Transaction, *VerifiableLoadUpkeep, error) {
	parsed, err := VerifiableLoadUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadUpkeepBin), backend, registrarAddress, useArb)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifiableLoadUpkeep{VerifiableLoadUpkeepCaller: VerifiableLoadUpkeepCaller{contract: contract}, VerifiableLoadUpkeepTransactor: VerifiableLoadUpkeepTransactor{contract: contract}, VerifiableLoadUpkeepFilterer: VerifiableLoadUpkeepFilterer{contract: contract}}, nil
}

type VerifiableLoadUpkeep struct {
	address common.Address
	abi     abi.ABI
	VerifiableLoadUpkeepCaller
	VerifiableLoadUpkeepTransactor
	VerifiableLoadUpkeepFilterer
}

type VerifiableLoadUpkeepCaller struct {
	contract *bind.BoundContract
}

type VerifiableLoadUpkeepTransactor struct {
	contract *bind.BoundContract
}

type VerifiableLoadUpkeepFilterer struct {
	contract *bind.BoundContract
}

type VerifiableLoadUpkeepSession struct {
	Contract     *VerifiableLoadUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VerifiableLoadUpkeepCallerSession struct {
	Contract *VerifiableLoadUpkeepCaller
	CallOpts bind.CallOpts
}

type VerifiableLoadUpkeepTransactorSession struct {
	Contract     *VerifiableLoadUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type VerifiableLoadUpkeepRaw struct {
	Contract *VerifiableLoadUpkeep
}

type VerifiableLoadUpkeepCallerRaw struct {
	Contract *VerifiableLoadUpkeepCaller
}

type VerifiableLoadUpkeepTransactorRaw struct {
	Contract *VerifiableLoadUpkeepTransactor
}

func NewVerifiableLoadUpkeep(address common.Address, backend bind.ContractBackend) (*VerifiableLoadUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(VerifiableLoadUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVerifiableLoadUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeep{address: address, abi: abi, VerifiableLoadUpkeepCaller: VerifiableLoadUpkeepCaller{contract: contract}, VerifiableLoadUpkeepTransactor: VerifiableLoadUpkeepTransactor{contract: contract}, VerifiableLoadUpkeepFilterer: VerifiableLoadUpkeepFilterer{contract: contract}}, nil
}

func NewVerifiableLoadUpkeepCaller(address common.Address, caller bind.ContractCaller) (*VerifiableLoadUpkeepCaller, error) {
	contract, err := bindVerifiableLoadUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepCaller{contract: contract}, nil
}

func NewVerifiableLoadUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifiableLoadUpkeepTransactor, error) {
	contract, err := bindVerifiableLoadUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepTransactor{contract: contract}, nil
}

func NewVerifiableLoadUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifiableLoadUpkeepFilterer, error) {
	contract, err := bindVerifiableLoadUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepFilterer{contract: contract}, nil
}

func bindVerifiableLoadUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifiableLoadUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadUpkeep.Contract.VerifiableLoadUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.VerifiableLoadUpkeepTransactor.contract.Transfer(opts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.VerifiableLoadUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.contract.Transfer(opts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) BUCKETSIZE(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "BUCKET_SIZE")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) TIMESTAMPINTERVAL(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "TIMESTAMP_INTERVAL")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) TIMESTAMPINTERVAL() (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.TIMESTAMPINTERVAL(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) TIMESTAMPINTERVAL() (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.TIMESTAMPINTERVAL(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) AddLinkAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "addLinkAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.AddLinkAmount(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.AddLinkAmount(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "bucketedDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.BucketedDelays(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.BucketedDelays(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "buckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.Buckets(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.Buckets(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) CheckDatas(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "checkDatas", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) CheckDatas(arg0 *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.CheckDatas(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) CheckDatas(arg0 *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.CheckDatas(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "checkGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "counters", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Counters(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Counters(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "delays", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Delays(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Delays(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadUpkeep.Contract.DummyMap(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadUpkeep.Contract.DummyMap(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "eligible", upkeepId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadUpkeep.Contract.Eligible(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadUpkeep.Contract.Eligible(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "firstPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "gasLimits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GasLimits(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GasLimits(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getActiveUpkeepIDs", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetActiveUpkeepIDs(&_VerifiableLoadUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetActiveUpkeepIDs(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetActiveUpkeepIDs(&_VerifiableLoadUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getBucketedDelays", upkeepId, bucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getDelays", upkeepId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelays(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelays(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelaysLength(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelaysLength(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetDelaysLengthAtBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getDelaysLengthAtBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetDelaysLengthAtBucket(upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelaysLengthAtBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetDelaysLengthAtBucket(upkeepId *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelaysLengthAtBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetDelaysLengthAtTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getDelaysLengthAtTimestampBucket", upkeepId, timestampBucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetDelaysLengthAtTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelaysLengthAtTimestampBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetDelaysLengthAtTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetDelaysLengthAtTimestampBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetLogTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getLogTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetLogTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetLogTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetPxBucketedDelaysForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getPxBucketedDelaysForAllUpkeeps", p)

	if err != nil {
		return *new([]*big.Int), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetPxBucketedDelaysForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxBucketedDelaysForAllUpkeeps(&_VerifiableLoadUpkeep.CallOpts, p)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetPxBucketedDelaysForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxBucketedDelaysForAllUpkeeps(&_VerifiableLoadUpkeep.CallOpts, p)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetPxDelayForAllUpkeeps(opts *bind.CallOpts, p *big.Int) ([]*big.Int, []*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getPxDelayForAllUpkeeps", p)

	if err != nil {
		return *new([]*big.Int), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetPxDelayForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayForAllUpkeeps(&_VerifiableLoadUpkeep.CallOpts, p)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetPxDelayForAllUpkeeps(p *big.Int) ([]*big.Int, []*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayForAllUpkeeps(&_VerifiableLoadUpkeep.CallOpts, p)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetPxDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getPxDelayInBucket", upkeepId, p, bucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetPxDelayInBucket(upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayInBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, p, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetPxDelayInBucket(upkeepId *big.Int, p *big.Int, bucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayInBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, p, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetPxDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getPxDelayInTimestampBucket", upkeepId, p, timestampBucket)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetPxDelayInTimestampBucket(upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayInTimestampBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, p, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetPxDelayInTimestampBucket(upkeepId *big.Int, p *big.Int, timestampBucket uint16) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayInTimestampBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, p, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getPxDelayLastNPerforms", upkeepId, p, n)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetSumBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getSumBucketedDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetSumBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumBucketedDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetSumBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumBucketedDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getSumDelayInBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetSumDelayInTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getSumDelayInTimestampBucket", upkeepId, timestampBucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetSumDelayInTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumDelayInTimestampBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetSumDelayInTimestampBucket(upkeepId *big.Int, timestampBucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumDelayInTimestampBucket(&_VerifiableLoadUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getSumDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetSumTimestampBucketedDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getSumTimestampBucketedDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetSumTimestampBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumTimestampBucketedDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetSumTimestampBucketedDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetSumTimestampBucketedDelayLastNPerforms(&_VerifiableLoadUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetTimestampBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getTimestampBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetTimestampBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetTimestampBucketedDelaysLength(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetTimestampBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetTimestampBucketedDelaysLength(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetTimestampDelays(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getTimestampDelays", upkeepId, timestampBucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetTimestampDelays(upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetTimestampDelays(&_VerifiableLoadUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetTimestampDelays(upkeepId *big.Int, timestampBucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetTimestampDelays(&_VerifiableLoadUpkeep.CallOpts, upkeepId, timestampBucket)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Intervals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "intervals", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Intervals(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Intervals(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) LastTopUpBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "lastTopUpBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.LinkToken(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.LinkToken(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) MinBalanceThresholdMultiplier(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "minBalanceThresholdMultiplier")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Owner() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.Owner(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Owner() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.Owner(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) PerformDataSizes(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "performDataSizes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.PerformDataSizes(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.PerformDataSizes(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) PerformGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "performGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) PreviousPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "previousPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Registrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "registrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Registrar() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.Registrar(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Registrar() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.Registrar(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Registry() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.Registry(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Registry() (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.Registry(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) TimestampBuckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "timestampBuckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) TimestampBuckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.TimestampBuckets(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) TimestampBuckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadUpkeep.Contract.TimestampBuckets(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) TimestampDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "timestampDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) TimestampDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.TimestampDelays(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) TimestampDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.TimestampDelays(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) Timestamps(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "timestamps", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Timestamps(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Timestamps(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) Timestamps(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.Timestamps(&_VerifiableLoadUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) UpkeepTopUpCheckInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "upkeepTopUpCheckInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "useArbitrumBlockNum")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "acceptOwnership")
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.AcceptOwnership(&_VerifiableLoadUpkeep.TransactOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.AcceptOwnership(&_VerifiableLoadUpkeep.TransactOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) AddFunds(opts *bind.TransactOpts, upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "addFunds", upkeepId, amount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.AddFunds(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.AddFunds(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchCancelUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchCancelUpkeeps", upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchRegisterUpkeeps", number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchSetIntervals", upkeepIds, interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchUpdatePipelineData", upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchWithdrawLinks", upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "cancelUpkeep", upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.CancelUpkeep(&_VerifiableLoadUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) CancelUpkeep(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.CancelUpkeep(&_VerifiableLoadUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) CheckUpkeep(opts *bind.TransactOpts, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "checkUpkeep", checkData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.CheckUpkeep(&_VerifiableLoadUpkeep.TransactOpts, checkData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.CheckUpkeep(&_VerifiableLoadUpkeep.TransactOpts, checkData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.PerformUpkeep(&_VerifiableLoadUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.PerformUpkeep(&_VerifiableLoadUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetAddLinkAmount(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setAddLinkAmount", amount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetAddLinkAmount(amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetAddLinkAmount(&_VerifiableLoadUpkeep.TransactOpts, amount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetAddLinkAmount(amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetAddLinkAmount(&_VerifiableLoadUpkeep.TransactOpts, amount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetCheckGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setCheckGasToBurn", upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetCheckGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetCheckGasToBurn(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetCheckGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetCheckGasToBurn(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setConfig", newRegistrar)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetConfig(&_VerifiableLoadUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetConfig(&_VerifiableLoadUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setInterval", upkeepId, _interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetInterval(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetInterval(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetMinBalanceThresholdMultiplier(opts *bind.TransactOpts, newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setMinBalanceThresholdMultiplier", newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetMinBalanceThresholdMultiplier(newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetMinBalanceThresholdMultiplier(&_VerifiableLoadUpkeep.TransactOpts, newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetMinBalanceThresholdMultiplier(newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetMinBalanceThresholdMultiplier(&_VerifiableLoadUpkeep.TransactOpts, newMinBalanceThresholdMultiplier)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setPerformDataSize", upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetPerformGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setPerformGasToBurn", upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetPerformGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetPerformGasToBurn(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetPerformGasToBurn(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetPerformGasToBurn(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setUpkeepGasLimit", upkeepId, gasLimit)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetUpkeepTopUpCheckInterval(opts *bind.TransactOpts, newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setUpkeepTopUpCheckInterval", newInterval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetUpkeepTopUpCheckInterval(newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepTopUpCheckInterval(&_VerifiableLoadUpkeep.TransactOpts, newInterval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetUpkeepTopUpCheckInterval(newInterval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepTopUpCheckInterval(&_VerifiableLoadUpkeep.TransactOpts, newInterval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "transferOwnership", to)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.TransferOwnership(&_VerifiableLoadUpkeep.TransactOpts, to)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.TransferOwnership(&_VerifiableLoadUpkeep.TransactOpts, to)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "updateUpkeepPipelineData", upkeepId, pipelineData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "withdrawLinks")
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.WithdrawLinks(&_VerifiableLoadUpkeep.TransactOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.WithdrawLinks(&_VerifiableLoadUpkeep.TransactOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "withdrawLinks0", upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.RawTransact(opts, nil)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.Receive(&_VerifiableLoadUpkeep.TransactOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.Receive(&_VerifiableLoadUpkeep.TransactOpts)
}

type VerifiableLoadUpkeepFundsAddedIterator struct {
	Event *VerifiableLoadUpkeepFundsAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepFundsAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepFundsAdded)
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
		it.Event = new(VerifiableLoadUpkeepFundsAdded)
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

func (it *VerifiableLoadUpkeepFundsAddedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepFundsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepFundsAdded struct {
	UpkeepId *big.Int
	Amount   *big.Int
	Raw      types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadUpkeepFundsAddedIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepFundsAddedIterator{contract: _VerifiableLoadUpkeep.contract, event: "FundsAdded", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepFundsAdded) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "FundsAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepFundsAdded)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "FundsAdded", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseFundsAdded(log types.Log) (*VerifiableLoadUpkeepFundsAdded, error) {
	event := new(VerifiableLoadUpkeepFundsAdded)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "FundsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepInsufficientFundsIterator struct {
	Event *VerifiableLoadUpkeepInsufficientFunds

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepInsufficientFundsIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepInsufficientFunds)
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
		it.Event = new(VerifiableLoadUpkeepInsufficientFunds)
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

func (it *VerifiableLoadUpkeepInsufficientFundsIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepInsufficientFundsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepInsufficientFunds struct {
	Balance  *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadUpkeepInsufficientFundsIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "InsufficientFunds")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepInsufficientFundsIterator{contract: _VerifiableLoadUpkeep.contract, event: "InsufficientFunds", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepInsufficientFunds) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "InsufficientFunds")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepInsufficientFunds)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "InsufficientFunds", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseInsufficientFunds(log types.Log) (*VerifiableLoadUpkeepInsufficientFunds, error) {
	event := new(VerifiableLoadUpkeepInsufficientFunds)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "InsufficientFunds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepLogEmittedIterator struct {
	Event *VerifiableLoadUpkeepLogEmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepLogEmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepLogEmitted)
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
		it.Event = new(VerifiableLoadUpkeepLogEmitted)
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

func (it *VerifiableLoadUpkeepLogEmittedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepLogEmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepLogEmitted struct {
	UpkeepId    *big.Int
	LogBlockNum *big.Int
	BlockNum    *big.Int
	Raw         types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int) (*VerifiableLoadUpkeepLogEmittedIterator, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "LogEmitted", upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepLogEmittedIterator{contract: _VerifiableLoadUpkeep.contract, event: "LogEmitted", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepLogEmitted, upkeepId []*big.Int) (event.Subscription, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "LogEmitted", upkeepIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepLogEmitted)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseLogEmitted(log types.Log) (*VerifiableLoadUpkeepLogEmitted, error) {
	event := new(VerifiableLoadUpkeepLogEmitted)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepOwnershipTransferRequestedIterator struct {
	Event *VerifiableLoadUpkeepOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepOwnershipTransferRequested)
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
		it.Event = new(VerifiableLoadUpkeepOwnershipTransferRequested)
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

func (it *VerifiableLoadUpkeepOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadUpkeepOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepOwnershipTransferRequestedIterator{contract: _VerifiableLoadUpkeep.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepOwnershipTransferRequested)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadUpkeepOwnershipTransferRequested, error) {
	event := new(VerifiableLoadUpkeepOwnershipTransferRequested)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepOwnershipTransferredIterator struct {
	Event *VerifiableLoadUpkeepOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepOwnershipTransferred)
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
		it.Event = new(VerifiableLoadUpkeepOwnershipTransferred)
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

func (it *VerifiableLoadUpkeepOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadUpkeepOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepOwnershipTransferredIterator{contract: _VerifiableLoadUpkeep.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepOwnershipTransferred)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseOwnershipTransferred(log types.Log) (*VerifiableLoadUpkeepOwnershipTransferred, error) {
	event := new(VerifiableLoadUpkeepOwnershipTransferred)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepPerformingUpkeepIterator struct {
	Event *VerifiableLoadUpkeepPerformingUpkeep

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepPerformingUpkeepIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepPerformingUpkeep)
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
		it.Event = new(VerifiableLoadUpkeepPerformingUpkeep)
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

func (it *VerifiableLoadUpkeepPerformingUpkeepIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepPerformingUpkeepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepPerformingUpkeep struct {
	FirstPerformBlock *big.Int
	LastBlock         *big.Int
	PreviousBlock     *big.Int
	Counter           *big.Int
	Raw               types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterPerformingUpkeep(opts *bind.FilterOpts) (*VerifiableLoadUpkeepPerformingUpkeepIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepPerformingUpkeepIterator{contract: _VerifiableLoadUpkeep.contract, event: "PerformingUpkeep", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepPerformingUpkeep) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "PerformingUpkeep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepPerformingUpkeep)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParsePerformingUpkeep(log types.Log) (*VerifiableLoadUpkeepPerformingUpkeep, error) {
	event := new(VerifiableLoadUpkeepPerformingUpkeep)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "PerformingUpkeep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepReceivedIterator struct {
	Event *VerifiableLoadUpkeepReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepReceived)
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
		it.Event = new(VerifiableLoadUpkeepReceived)
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

func (it *VerifiableLoadUpkeepReceivedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepReceived struct {
	Sender common.Address
	Value  *big.Int
	Raw    types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterReceived(opts *bind.FilterOpts) (*VerifiableLoadUpkeepReceivedIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepReceivedIterator{contract: _VerifiableLoadUpkeep.contract, event: "Received", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchReceived(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepReceived) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "Received")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepReceived)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "Received", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseReceived(log types.Log) (*VerifiableLoadUpkeepReceived, error) {
	event := new(VerifiableLoadUpkeepReceived)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "Received", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepRegistrarSetIterator struct {
	Event *VerifiableLoadUpkeepRegistrarSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepRegistrarSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepRegistrarSet)
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
		it.Event = new(VerifiableLoadUpkeepRegistrarSet)
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

func (it *VerifiableLoadUpkeepRegistrarSetIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepRegistrarSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepRegistrarSet struct {
	NewRegistrar common.Address
	Raw          types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterRegistrarSet(opts *bind.FilterOpts) (*VerifiableLoadUpkeepRegistrarSetIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "RegistrarSet")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepRegistrarSetIterator{contract: _VerifiableLoadUpkeep.contract, event: "RegistrarSet", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchRegistrarSet(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepRegistrarSet) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "RegistrarSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepRegistrarSet)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "RegistrarSet", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseRegistrarSet(log types.Log) (*VerifiableLoadUpkeepRegistrarSet, error) {
	event := new(VerifiableLoadUpkeepRegistrarSet)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "RegistrarSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepUpkeepTopUpIterator struct {
	Event *VerifiableLoadUpkeepUpkeepTopUp

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepUpkeepTopUpIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepUpkeepTopUp)
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
		it.Event = new(VerifiableLoadUpkeepUpkeepTopUp)
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

func (it *VerifiableLoadUpkeepUpkeepTopUpIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepUpkeepTopUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepUpkeepTopUp struct {
	UpkeepId *big.Int
	Amount   *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadUpkeepUpkeepTopUpIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepUpkeepTopUpIterator{contract: _VerifiableLoadUpkeep.contract, event: "UpkeepTopUp", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepUpkeepTopUp) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepUpkeepTopUp)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseUpkeepTopUp(log types.Log) (*VerifiableLoadUpkeepUpkeepTopUp, error) {
	event := new(VerifiableLoadUpkeepUpkeepTopUp)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepUpkeepsCancelledIterator struct {
	Event *VerifiableLoadUpkeepUpkeepsCancelled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepUpkeepsCancelledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepUpkeepsCancelled)
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
		it.Event = new(VerifiableLoadUpkeepUpkeepsCancelled)
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

func (it *VerifiableLoadUpkeepUpkeepsCancelledIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepUpkeepsCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepUpkeepsCancelled struct {
	UpkeepIds []*big.Int
	Raw       types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterUpkeepsCancelled(opts *bind.FilterOpts) (*VerifiableLoadUpkeepUpkeepsCancelledIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "UpkeepsCancelled")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepUpkeepsCancelledIterator{contract: _VerifiableLoadUpkeep.contract, event: "UpkeepsCancelled", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchUpkeepsCancelled(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepUpkeepsCancelled) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "UpkeepsCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepUpkeepsCancelled)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "UpkeepsCancelled", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseUpkeepsCancelled(log types.Log) (*VerifiableLoadUpkeepUpkeepsCancelled, error) {
	event := new(VerifiableLoadUpkeepUpkeepsCancelled)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "UpkeepsCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadUpkeepUpkeepsRegisteredIterator struct {
	Event *VerifiableLoadUpkeepUpkeepsRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepUpkeepsRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepUpkeepsRegistered)
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
		it.Event = new(VerifiableLoadUpkeepUpkeepsRegistered)
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

func (it *VerifiableLoadUpkeepUpkeepsRegisteredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepUpkeepsRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepUpkeepsRegistered struct {
	UpkeepIds []*big.Int
	Raw       types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterUpkeepsRegistered(opts *bind.FilterOpts) (*VerifiableLoadUpkeepUpkeepsRegisteredIterator, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "UpkeepsRegistered")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepUpkeepsRegisteredIterator{contract: _VerifiableLoadUpkeep.contract, event: "UpkeepsRegistered", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchUpkeepsRegistered(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepUpkeepsRegistered) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "UpkeepsRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepUpkeepsRegistered)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "UpkeepsRegistered", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseUpkeepsRegistered(log types.Log) (*VerifiableLoadUpkeepUpkeepsRegistered, error) {
	event := new(VerifiableLoadUpkeepUpkeepsRegistered)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "UpkeepsRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifiableLoadUpkeep.abi.Events["FundsAdded"].ID:
		return _VerifiableLoadUpkeep.ParseFundsAdded(log)
	case _VerifiableLoadUpkeep.abi.Events["InsufficientFunds"].ID:
		return _VerifiableLoadUpkeep.ParseInsufficientFunds(log)
	case _VerifiableLoadUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadUpkeep.abi.Events["PerformingUpkeep"].ID:
		return _VerifiableLoadUpkeep.ParsePerformingUpkeep(log)
	case _VerifiableLoadUpkeep.abi.Events["Received"].ID:
		return _VerifiableLoadUpkeep.ParseReceived(log)
	case _VerifiableLoadUpkeep.abi.Events["RegistrarSet"].ID:
		return _VerifiableLoadUpkeep.ParseRegistrarSet(log)
	case _VerifiableLoadUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadUpkeep.ParseUpkeepTopUp(log)
	case _VerifiableLoadUpkeep.abi.Events["UpkeepsCancelled"].ID:
		return _VerifiableLoadUpkeep.ParseUpkeepsCancelled(log)
	case _VerifiableLoadUpkeep.abi.Events["UpkeepsRegistered"].ID:
		return _VerifiableLoadUpkeep.ParseUpkeepsRegistered(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadUpkeepFundsAdded) Topic() common.Hash {
	return common.HexToHash("0x8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e")
}

func (VerifiableLoadUpkeepInsufficientFunds) Topic() common.Hash {
	return common.HexToHash("0x03eb8b54a949acec2cd08fdb6d6bd4647a1f2c907d75d6900648effa92eb147f")
}

func (VerifiableLoadUpkeepLogEmitted) Topic() common.Hash {
	return common.HexToHash("0x8d98eacef480ad8f47c29266a1194f1874fdb68bcc98624964400d6ce72e69ec")
}

func (VerifiableLoadUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifiableLoadUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifiableLoadUpkeepPerformingUpkeep) Topic() common.Hash {
	return common.HexToHash("0x6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede3")
}

func (VerifiableLoadUpkeepReceived) Topic() common.Hash {
	return common.HexToHash("0x88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874")
}

func (VerifiableLoadUpkeepRegistrarSet) Topic() common.Hash {
	return common.HexToHash("0x6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014")
}

func (VerifiableLoadUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (VerifiableLoadUpkeepUpkeepsCancelled) Topic() common.Hash {
	return common.HexToHash("0xbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b")
}

func (VerifiableLoadUpkeepUpkeepsRegistered) Topic() common.Hash {
	return common.HexToHash("0x2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711")
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeep) Address() common.Address {
	return _VerifiableLoadUpkeep.address
}

type VerifiableLoadUpkeepInterface interface {
	BUCKETSIZE(opts *bind.CallOpts) (uint16, error)

	TIMESTAMPINTERVAL(opts *bind.CallOpts) (uint16, error)

	AddLinkAmount(opts *bind.CallOpts) (*big.Int, error)

	BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error)

	Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error)

	CheckDatas(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error)

	CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error)

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

	FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadUpkeepFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepFundsAdded) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*VerifiableLoadUpkeepFundsAdded, error)

	FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadUpkeepInsufficientFundsIterator, error)

	WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepInsufficientFunds) (event.Subscription, error)

	ParseInsufficientFunds(log types.Log) (*VerifiableLoadUpkeepInsufficientFunds, error)

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int) (*VerifiableLoadUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepLogEmitted, upkeepId []*big.Int) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadUpkeepLogEmitted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadUpkeepOwnershipTransferred, error)

	FilterPerformingUpkeep(opts *bind.FilterOpts) (*VerifiableLoadUpkeepPerformingUpkeepIterator, error)

	WatchPerformingUpkeep(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepPerformingUpkeep) (event.Subscription, error)

	ParsePerformingUpkeep(log types.Log) (*VerifiableLoadUpkeepPerformingUpkeep, error)

	FilterReceived(opts *bind.FilterOpts) (*VerifiableLoadUpkeepReceivedIterator, error)

	WatchReceived(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepReceived) (event.Subscription, error)

	ParseReceived(log types.Log) (*VerifiableLoadUpkeepReceived, error)

	FilterRegistrarSet(opts *bind.FilterOpts) (*VerifiableLoadUpkeepRegistrarSetIterator, error)

	WatchRegistrarSet(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepRegistrarSet) (event.Subscription, error)

	ParseRegistrarSet(log types.Log) (*VerifiableLoadUpkeepRegistrarSet, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadUpkeepUpkeepTopUp, error)

	FilterUpkeepsCancelled(opts *bind.FilterOpts) (*VerifiableLoadUpkeepUpkeepsCancelledIterator, error)

	WatchUpkeepsCancelled(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepUpkeepsCancelled) (event.Subscription, error)

	ParseUpkeepsCancelled(log types.Log) (*VerifiableLoadUpkeepUpkeepsCancelled, error)

	FilterUpkeepsRegistered(opts *bind.FilterOpts) (*VerifiableLoadUpkeepUpkeepsRegisteredIterator, error)

	WatchUpkeepsRegistered(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepUpkeepsRegistered) (event.Subscription, error)

	ParseUpkeepsRegistered(log types.Log) (*VerifiableLoadUpkeepUpkeepsRegistered, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
