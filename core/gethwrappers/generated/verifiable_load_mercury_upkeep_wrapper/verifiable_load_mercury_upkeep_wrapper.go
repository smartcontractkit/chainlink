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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"registrarAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedLabel\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feedList\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"queryLabel\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"query\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"MercuryLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"origin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"firstPerformBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"RegistrarSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkDatas\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxBucketedDelaysForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxDelayForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumTimestampBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTimestampBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getTimestampDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"mercuryCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queryLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractKeeperRegistry2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampBuckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040526005601855601980546001600160681b0319166c140000000002c68af0bb140000179055601960f21b60a05260e160f41b60c0523480156200004557600080fd5b506040516200582538038062005825833981016040819052620000689162000320565b81813380600081620000c15760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000f457620000f4816200025a565b5050601580546001600160a01b0319166001600160a01b0385169081179091556040805163850af0cb60e01b815290516000935063850af0cb9160048082019260a092909190829003018186803b1580156200014f57600080fd5b505afa15801562000164573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200018a91906200038a565b50601780546001600160a01b0319166001600160a01b038381169190911790915560155460408051631b6b6d2360e01b8152905193975091169450631b6b6d2393506004808201935060209291829003018186803b158015620001ec57600080fd5b505afa15801562000201573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000227919062000363565b601680546001600160a01b0319166001600160a01b039290921691909117905550151560f81b6080525062000413915050565b6001600160a01b038116331415620002b55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000b8565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b805163ffffffff811681146200031b57600080fd5b919050565b600080604083850312156200033457600080fd5b82516200034181620003fa565b602084015190925080151581146200035857600080fd5b809150509250929050565b6000602082840312156200037657600080fd5b81516200038381620003fa565b9392505050565b600080600080600060a08688031215620003a357600080fd5b855160038110620003b357600080fd5b9450620003c36020870162000306565b9350620003d36040870162000306565b92506060860151620003e581620003fa565b80925050608086015190509295509295909350565b6001600160a01b03811681146200041057600080fd5b50565b60805160f81c60a05160f01c60c05160f01c6153c4620004616000396000818161069e0152611a530152600081816105740152611b670152600081816109b001526139cf01526153c46000f3fe6080604052600436106104ba5760003560e01c8063776898c8116102795780639fab43861161015e578063d3558528116100d6578063f2fde38b1161008a578063fba7ffa31161006f578063fba7ffa314611056578063fbfb4f7614611083578063fcdc1f63146110a357600080fd5b8063f2fde38b14611016578063fb0ceb041461103657600080fd5b8063dbef701e116100bb578063dbef701e14610fb3578063e0114adb14610fd3578063e45530831461100057600080fd5b8063d355852814610f34578063d6051a7214610f9357600080fd5b8063a79c40431161012d578063b0971e1a11610112578063b0971e1a14610e7c578063c357f1f314610eba578063c804802214610f1457600080fd5b8063a79c404314610e2f578063af953a4a14610e5c57600080fd5b80639fab438614610daf578063a5f5893414610dcf578063a6c60d8914610def578063a72aa27e14610e0f57600080fd5b80638fcb3fba116101f15780639ac542eb116101c05780639b51fb0d116101a55780639b51fb0d14610d3e5780639bb8651114610d6f5780639d385eaa14610d8f57600080fd5b80639ac542eb14610cd55780639b42935414610d1157600080fd5b80638fcb3fba14610c485780639095aa3514610c75578063948108f714610c9557806399cc6b0b14610cb557600080fd5b806380f4df1b1161024857806387dfa9001161022d57806387dfa90014610bdd5780638bc7b77214610bfd5780638da5cb5b14610c1d57600080fd5b806380f4df1b14610ba85780638237831714610bbd57600080fd5b8063776898c814610b2657806379ba509714610b465780637b10399914610b5b5780637e4087b814610b8857600080fd5b80634ad8c9a61161039f578063642f6cef116103175780636e04ff0d116102e65780637145f11b116102cb5780637145f11b14610a9c57806373644cce14610acc5780637672130314610af957600080fd5b80636e04ff0d14610a5c5780637137a70214610a7c57600080fd5b8063642f6cef1461099e578063643b34e9146109e257806369cdbadb14610a0257806369e9b77314610a2f57600080fd5b806358c52c041161036e5780635f17e616116103535780635f17e6161461090f57806360457ff51461092f578063636092e81461095c57600080fd5b806358c52c04146108da5780635d4ee7f3146108fa57600080fd5b80634ad8c9a61461084a5780634d6954451461087857806351c98be31461088d57806357970e93146108ad57600080fd5b80632a9032d3116104325780633ebe8d6c1161040157806345d2ec17116103e657806345d2ec17146107d057806346e7a63e146107f05780634a5479f31461081d57600080fd5b80633ebe8d6c146107905780634585e33b146107b057600080fd5b80632a9032d3146106c05780632b20e397146106e0578063328ffd111461073257806333774d1c1461075f57600080fd5b80631bee00801161048957806320e3dbd41161046e57806320e3dbd41461064a57806328c4b57b1461066c57806329f0e4961461068c57600080fd5b80631bee0080146105e7578063206c32e81461061557600080fd5b806306e3b632146104fe578063077ac6211461053457806312c5502714610562578063177b0eb9146105a957600080fd5b366104f957604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561050a57600080fd5b5061051e61051936600461475c565b6110d0565b60405161052b9190614b1f565b60405180910390f35b34801561054057600080fd5b5061055461054f366004614727565b6111cc565b60405190815260200161052b565b34801561056e57600080fd5b506105967f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff909116815260200161052b565b3480156105b557600080fd5b506105546105c43660046146fb565b6000918252600f6020908152604080842061ffff93909316845291905290205490565b3480156105f357600080fd5b5061060761060236600461449c565b61120a565b60405161052b929190614b32565b34801561062157600080fd5b506106356106303660046146fb565b611513565b6040805192835260208301919091520161052b565b34801561065657600080fd5b5061066a6106653660046142e0565b611596565b005b34801561067857600080fd5b506105546106873660046147b3565b6117b7565b34801561069857600080fd5b506105967f000000000000000000000000000000000000000000000000000000000000000081565b3480156106cc57600080fd5b5061066a6106db3660046143e8565b611822565b3480156106ec57600080fd5b5060155461070d9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161052b565b34801561073e57600080fd5b5061055461074d36600461449c565b60036020526000908152604090205481565b34801561076b57600080fd5b5061059661077a36600461449c565b60116020526000908152604090205461ffff1681565b34801561079c57600080fd5b506105546107ab36600461449c565b6118f5565b3480156107bc57600080fd5b5061066a6107cb3660046144b5565b61195e565b3480156107dc57600080fd5b5061051e6107eb3660046146fb565b61207b565b3480156107fc57600080fd5b5061055461080b36600461449c565b600a6020526000908152604090205481565b34801561082957600080fd5b5061083d61083836600461449c565b6120ea565b60405161052b9190614b72565b34801561085657600080fd5b5061086a6108653660046142fd565b612196565b60405161052b929190614b57565b34801561088457600080fd5b5061083d6121ea565b34801561089957600080fd5b5061066a6108a836600461442a565b6121f7565b3480156108b957600080fd5b5060165461070d9073ffffffffffffffffffffffffffffffffffffffff1681565b3480156108e657600080fd5b5061083d6108f536600461449c565b61229b565b34801561090657600080fd5b5061066a6122b4565b34801561091b57600080fd5b5061066a61092a36600461475c565b61240d565b34801561093b57600080fd5b5061055461094a36600461449c565b60076020526000908152604090205481565b34801561096857600080fd5b50601954610981906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff909116815260200161052b565b3480156109aa57600080fd5b506109d27f000000000000000000000000000000000000000000000000000000000000000081565b604051901515815260200161052b565b3480156109ee57600080fd5b506106356109fd36600461475c565b61257f565b348015610a0e57600080fd5b50610554610a1d36600461449c565b60086020526000908152604090205481565b348015610a3b57600080fd5b5061066a610a4a36600461475c565b60009182526008602052604090912055565b348015610a6857600080fd5b5061086a610a773660046144b5565b612704565b348015610a8857600080fd5b50610554610a97366004614727565b6128b1565b348015610aa857600080fd5b506109d2610ab736600461449c565b600c6020526000908152604090205460ff1681565b348015610ad857600080fd5b50610554610ae736600461449c565b6000908152600d602052604090205490565b348015610b0557600080fd5b50610554610b1436600461449c565b60046020526000908152604090205481565b348015610b3257600080fd5b506109d2610b4136600461449c565b6128d9565b348015610b5257600080fd5b5061066a612929565b348015610b6757600080fd5b5060175461070d9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b9457600080fd5b50610635610ba336600461475c565b612a26565b348015610bb457600080fd5b5061083d612b9e565b348015610bc957600080fd5b50610554610bd836600461477e565b612bab565b348015610be957600080fd5b50610554610bf836600461477e565b612c26565b348015610c0957600080fd5b50610607610c1836600461449c565b612c96565b348015610c2957600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661070d565b348015610c5457600080fd5b50610554610c6336600461449c565b60056020526000908152604090205481565b348015610c8157600080fd5b5061066a610c9036600461484f565b612e0b565b348015610ca157600080fd5b5061066a610cb036600461480f565b61308b565b348015610cc157600080fd5b5061051e610cd03660046146fb565b613223565b348015610ce157600080fd5b50601954610cff906c01000000000000000000000000900460ff1681565b60405160ff909116815260200161052b565b348015610d1d57600080fd5b5061066a610d2c36600461475c565b60009182526009602052604090912055565b348015610d4a57600080fd5b50610596610d5936600461449c565b60126020526000908152604090205461ffff1681565b348015610d7b57600080fd5b5061066a610d8a3660046143e8565b613290565b348015610d9b57600080fd5b5061051e610daa36600461449c565b613361565b348015610dbb57600080fd5b5061066a610dca3660046146af565b6133c3565b348015610ddb57600080fd5b50610554610dea36600461449c565b613468565b348015610dfb57600080fd5b5061066a610e0a36600461449c565b601855565b348015610e1b57600080fd5b5061066a610e2a3660046147df565b6134c9565b348015610e3b57600080fd5b5061066a610e4a36600461475c565b60009182526007602052604090912055565b348015610e6857600080fd5b5061066a610e7736600461449c565b613574565b348015610e8857600080fd5b50610554610e973660046146fb565b6000918252600e6020908152604080842061ffff93909316845291905290205490565b348015610ec657600080fd5b5061066a610ed53660046148a8565b601980547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610f2057600080fd5b5061066a610f2f36600461449c565b6135fa565b348015610f4057600080fd5b5061066a610f4f366004614834565b6019805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b348015610f9f57600080fd5b50610635610fae36600461475c565b613692565b348015610fbf57600080fd5b50610554610fce36600461475c565b6136fb565b348015610fdf57600080fd5b50610554610fee36600461449c565b60096020526000908152604090205481565b34801561100c57600080fd5b5061055460185481565b34801561102257600080fd5b5061066a6110313660046142e0565b61372c565b34801561104257600080fd5b5061055461105136600461475c565b613740565b34801561106257600080fd5b5061055461107136600461449c565b60066020526000908152604090205481565b34801561108f57600080fd5b5061063561109e3660046146fb565b61375c565b3480156110af57600080fd5b506105546110be36600461449c565b60026020526000908152604090205481565b606060006110de60136137d0565b9050808410611119576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8261112b5761112884826150a6565b92505b60008367ffffffffffffffff8111156111465761114661533a565b60405190808252806020026020018201604052801561116f578160200160208202803683370190505b50905060005b848110156111c15761119261118a8288614f2d565b6013906137da565b8282815181106111a4576111a461530b565b6020908102919091010152806111b98161522c565b915050611175565b509150505b92915050565b600e60205282600052604060002060205281600052604060002081815481106111f457600080fd5b9060005260206000200160009250925050505481565b606080600061121960136137d0565b905060008167ffffffffffffffff8111156112365761123661533a565b60405190808252806020026020018201604052801561125f578160200160208202803683370190505b50905060008267ffffffffffffffff81111561127d5761127d61533a565b6040519080825280602002602001820160405280156112a6578160200160208202803683370190505b50905060005b838110156115075760006112c16013836137da565b9050808483815181106112d6576112d661530b565b6020908102919091018101919091526000828152601290915260408082205490517f3ebe8d6c0000000000000000000000000000000000000000000000000000000081526004810184905261ffff90911691903090633ebe8d6c9060240160206040518083038186803b15801561134c57600080fd5b505afa158015611360573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113849190614696565b905060008167ffffffffffffffff8111156113a1576113a161533a565b6040519080825280602002602001820160405280156113ca578160200160208202803683370190505b506000858152600e6020526040812091925090815b8561ffff168161ffff16116114c45761ffff81166000908152602083815260408083208054825181850281018501909352808352919290919083018282801561144757602002820191906000526020600020905b815481526020019060010190808311611433575b5050505050905060005b81518110156114af5781818151811061146c5761146c61530b565b60200260200101518686806114809061522c565b9750815181106114925761149261530b565b6020908102919091010152806114a78161522c565b915050611451565b505080806114bc9061520a565b9150506113df565b506114d0838e866137e6565b8888815181106114e2576114e261530b565b60200260200101818152505050505050505080806114ff9061522c565b9150506112ac565b50909590945092505050565b6000828152600e6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493849392919083018282801561157757602002820191906000526020600020905b815481526020019060010190808311611563575b50505050509050611589818251613946565b92509250505b9250929050565b601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517f850af0cb00000000000000000000000000000000000000000000000000000000815290516000929163850af0cb9160048083019260a0929190829003018186803b15801561162b57600080fd5b505afa15801561163f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116639190614508565b50601780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601554604080517f1b6b6d23000000000000000000000000000000000000000000000000000000008152905193975091169450631b6b6d2393506004808201935060209291829003018186803b15801561170257600080fd5b505afa158015611716573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061173a91906144eb565b601680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601554604051911681527f6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014906020015b60405180910390a15050565b6000838152600d6020908152604080832080548251818502810185019093528083526118189383018282801561180c57602002820191906000526020600020905b8154815260200190600101908083116117f8575b505050505084846137e6565b90505b9392505050565b8060005b818160ff1610156118b6573063c8048022858560ff851681811061184c5761184c61530b565b905060200201356040518263ffffffff1660e01b815260040161187191815260200190565b600060405180830381600087803b15801561188b57600080fd5b505af115801561189f573d6000803e3d6000fd5b5050505080806118ae9061525e565b915050611826565b507fbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b83836040516118e8929190614aca565b60405180910390a1505050565b60008181526012602052604081205461ffff1681805b8261ffff168161ffff1611611956576000858152600e6020908152604080832061ffff851684529091529020546119429083614f2d565b91508061194e8161520a565b91505061190b565b509392505050565b60005a9050600080611972848601866142fd565b9150915060008180602001905181019061198c9190614696565b600081815260056020908152604080832054600490925282205492935091906119b36139cb565b9050826119ed5760008481526005602090815260408083208490556010825282208054600181018255908352912042910155915081611c4d565b600084815260036020526040812054611a0684846150a6565b611a1091906150a6565b6000868152601160209081526040808320546010909252909120805492935061ffff9091169182908110611a4657611a4661530b565b90600052602060002001547f000000000000000000000000000000000000000000000000000000000000000061ffff1642611a8191906150a6565b1115611af05760008681526010602090815260408220805460018101825590835291204291015580611ab28161520a565b600088815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559150505b600086815260126020908152604080832054600e835281842061ffff9091168085529083528184208054835181860281018601909452808452919493909190830182828015611b5e57602002820191906000526020600020905b815481526020019060010190808311611b4a575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff1681511415611bda5781611b9c8161520a565b60008a815260126020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600e6020908152604080832061ffff94851684528252808320805460018181018355918552838520018790558a8452600f83528184209590941683529381528382208054808501825590835281832001859055888252600d81529281208054928301815581529190912001555b600084815260066020526040812054611c67906001614f2d565b6000868152600660209081526040918290208390558151878152908101859052908101859052606081018290529091507f6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede39060800160405180910390a16000858152600460209081526040808320859055601854600290925290912054611cee90846150a6565b1115611f99576017546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810187905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a9060240160006040518083038186803b158015611d5f57600080fd5b505afa158015611d73573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611db99190810190614577565b6017546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810189905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c9060240160206040518083038186803b158015611e2957600080fd5b505afa158015611e3d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e6191906148c5565b601954909150611e859082906c01000000000000000000000000900460ff16614ffe565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611f96576019546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018990526bffffffffffffffffffffffff9091166024820152309063948108f790604401600060405180830381600087803b158015611f1657600080fd5b505af1158015611f2a573d6000803e3d6000fd5b50505060008881526002602090815260409182902087905560195482518b81526bffffffffffffffffffffffff909116918101919091529081018690527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0915060600160405180910390a15b50505b6000858152600760205260409020545b805a611fb5908b6150a6565b611fc190612710614f2d565b10156120025782406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055611fa9565b82863273ffffffffffffffffffffffffffffffffffffffff167fcad583be2d908a590c81c7e332cf11c7a4ea41ecf1e059efac3ea7e83e34f1a58b60008151811061204f5761204f61530b565b60200260200101518b604051612066929190614b85565b60405180910390a45050505050505050505050565b6000828152600e6020908152604080832061ffff851684528252918290208054835181840281018401909452808452606093928301828280156120dd57602002820191906000526020600020905b8154815260200190600101908083116120c9575b5050505050905092915050565b601a81815481106120fa57600080fd5b9060005260206000200160009150905080546121159061517d565b80601f01602080910402602001604051908101604052809291908181526020018280546121419061517d565b801561218e5780601f106121635761010080835404028352916020019161218e565b820191906000526020600020905b81548152906001019060200180831161217157829003601f168201915b505050505081565b60006060600084846040516020016121af929190614a3f565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b601b80546121159061517d565b8160005b818110156122945730635f17e61686868481811061221b5761221b61530b565b90506020020135856040518363ffffffff1660e01b815260040161224f92919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561226957600080fd5b505af115801561227d573d6000803e3d6000fd5b50505050808061228c9061522c565b9150506121fb565b5050505050565b600b60205260009081526040902080546121159061517d565b6122bc613a7c565b6016546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561232657600080fd5b505afa15801561233a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061235e9190614696565b6016546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb90604401602060405180830381600087803b1580156123d157600080fd5b505af11580156123e5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124099190614481565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600d909152812061244591613ffc565b60008281526012602052604081205461ffff16905b8161ffff168161ffff16116124a1576000848152600e6020908152604080832061ffff85168452909152812061248f91613ffc565b806124998161520a565b91505061245a565b5050600082815260126020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055601190915281205461ffff16905b8161ffff168161ffff161161252f576000848152600f6020908152604080832061ffff85168452909152812061251d91613ffc565b806125278161520a565b9150506124e8565b50600083815260106020526040812061254791613ffc565b5050600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6040517f3ebe8d6c00000000000000000000000000000000000000000000000000000000815260048101839052600090819081903090633ebe8d6c9060240160206040518083038186803b1580156125d657600080fd5b505afa1580156125ea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061260e9190614696565b905083158061261d5750808410155b15612626578093505b60008581526012602052604081205485919061ffff16805b6000898152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835291929091908301828280156126a057602002820191906000526020600020905b81548152602001906001019080831161268c575b505050505090506000806126b48388613946565b90925090506126c38287614f2d565b95506126cf81886150a6565b9650600087116126e1575050506126f7565b50505080806126ef90615141565b91505061263e565b5090979596505050505050565b6000606060005a9050600061271b8587018761449c565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff8111156127545761275461533a565b6040519080825280601f01601f19166020018201604052801561277e576020820181803683370190505b50604051602001612790929190614de1565b604051602081830303815290604052905060006127ab6139cb565b905060006127b8866128d9565b90505b835a6127c790896150a6565b6127d390612710614f2d565b10156128145781406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556127bb565b8061282c57600083985098505050505050505061158f565b601b601a601c848960405160200161284691815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f62e8a50d0000000000000000000000000000000000000000000000000000000082526128a89594939291600401614baa565b60405180910390fd5b600f60205282600052604060002060205281600052604060002081815481106111f457600080fd5b6000818152600560205260408120546128f457506001919050565b6000828152600360209081526040808320546004909252909120546129176139cb565b61292191906150a6565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146129aa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016128a8565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040517fa5f589340000000000000000000000000000000000000000000000000000000081526004810183905260009081908190309063a5f589349060240160206040518083038186803b158015612a7d57600080fd5b505afa158015612a91573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ab59190614696565b9050831580612ac45750808410155b15612acd578093505b60008581526011602052604081205485919061ffff16805b6000898152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083529192909190830182828015612b4757602002820191906000526020600020905b815481526020019060010190808311612b33575b50505050509050600080612b5b8388613946565b9092509050612b6a8287614f2d565b9550612b7681886150a6565b965060008711612b88575050506126f7565b5050508080612b9690615141565b915050612ae5565b601c80546121159061517d565b6000838152600e6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493830182828015612c0a57602002820191906000526020600020905b815481526020019060010190808311612bf6575b50505050509050612c1d818583516137e6565b95945050505050565b6000838152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493830182828015612c0a5760200282019190600052602060002090815481526020019060010190808311612bf65750505050509050612c1d818583516137e6565b6060806000612ca560136137d0565b905060008167ffffffffffffffff811115612cc257612cc261533a565b604051908082528060200260200182016040528015612ceb578160200160208202803683370190505b50905060008267ffffffffffffffff811115612d0957612d0961533a565b604051908082528060200260200182016040528015612d32578160200160208202803683370190505b50905060005b83811015611507576000612d4d6013836137da565b6000818152600d6020908152604080832080548251818502810185019093528083529495509293909291830182828015612da657602002820191906000526020600020905b815481526020019060010190808311612d92575b5050505050905081858481518110612dc057612dc061530b565b602002602001018181525050612dd8818a83516137e6565b848481518110612dea57612dea61530b565b60200260200101818152505050508080612e039061522c565b915050612d38565b6040805161014081018252600461010082019081527f746573740000000000000000000000000000000000000000000000000000000061012083015281528151602081810184526000808352818401929092523083850181905263ffffffff8916606085015260808401528351808201855282815260a08401528351908101909352825260c08101919091526bffffffffffffffffffffffff841660e082015260165460155473ffffffffffffffffffffffffffffffffffffffff9182169163095ea7b39116612ede60ff8a1688614ffe565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff166024820152604401602060405180830381600087803b158015612f5757600080fd5b505af1158015612f6b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f8f9190614481565b5060008660ff1667ffffffffffffffff811115612fae57612fae61533a565b604051908082528060200260200182016040528015612fd7578160200160208202803683370190505b50905060005b8760ff168160ff16101561304a576000612ff684613aff565b905080838360ff168151811061300e5761300e61530b565b602090810291909101810191909152600091825260088152604080832088905560079091529020849055806130428161525e565b915050612fdd565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c7118160405161307a9190614b1f565b60405180910390a150505050505050565b6016546017546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b390604401602060405180830381600087803b15801561310e57600080fd5b505af1158015613122573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131469190614481565b506017546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b1580156131c757600080fd5b505af11580156131db573d6000803e3d6000fd5b5050604080518581526bffffffffffffffffffffffff851660208201527f8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e93500190506117ab565b6000828152600f6020908152604080832061ffff851684528252918290208054835181840281018401909452808452606093928301828280156120dd57602002820191906000526020600020908154815260200190600101908083116120c9575050505050905092915050565b8060005b8181101561335b5760008484838181106132b0576132b061530b565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16639fab438682836040516020016132e991815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613315929190614de1565b600060405180830381600087803b15801561332f57600080fd5b505af1158015613343573d6000803e3d6000fd5b505050505080806133539061522c565b915050613294565b50505050565b6000818152600d60209081526040918290208054835181840281018401909452808452606093928301828280156133b757602002820191906000526020600020905b8154815260200190600101908083116133a3575b50505050509050919050565b6017546040517f9fab438600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690639fab43869061341d90869086908690600401614d8d565b600060405180830381600087803b15801561343757600080fd5b505af115801561344b573d6000803e3d6000fd5b5050506000848152600b6020526040902061335b9150838361401a565b60008181526011602052604081205461ffff1681805b8261ffff168161ffff1611611956576000858152600f6020908152604080832061ffff851684529091529020546134b59083614f2d565b9150806134c18161520a565b91505061347e565b6017546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b15801561354157600080fd5b505af1158015613555573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6017546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b1580156135e657600080fd5b505af1158015612294573d6000803e3d6000fd5b6017546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b15801561366657600080fd5b505af115801561367a573d6000803e3d6000fd5b50505050612409816013613c0190919063ffffffff16565b6000828152600d602090815260408083208054825181850281018501909352808352849384939291908301828280156136ea57602002820191906000526020600020905b8154815260200190600101908083116136d6575b505050505090506115898185613946565b600d602052816000526040600020818154811061371757600080fd5b90600052602060002001600091509150505481565b613734613a7c565b61373d81613c0d565b50565b6010602052816000526040600020818154811061371757600080fd5b6000828152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493849392919083018282801561157757602002820191906000526020600020908154815260200190600101908083116115635750505050509050611589818251613946565b60006111c6825490565b600061181b8383613d03565b825160009081908315806137fa5750808410155b15613803578093505b60008467ffffffffffffffff81111561381e5761381e61533a565b604051908082528060200260200182016040528015613847578160200160208202803683370190505b509050600092505b848310156138b55786600161386485856150a6565b61386e91906150a6565b8151811061387e5761387e61530b565b60200260200101518184815181106138985761389861530b565b6020908102919091010152826138ad8161522c565b93505061384f565b6138ce816000600184516138c991906150a6565b613d2d565b85606414156139085780600182516138e691906150a6565b815181106138f6576138f661530b565b6020026020010151935050505061181b565b8060648251886139189190614fc1565b6139229190614fad565b815181106139325761393261530b565b602002602001015193505050509392505050565b81516000908190819084158061395c5750808510155b15613965578094505b60008092505b858310156139c15786600161398085856150a6565b61398a91906150a6565b8151811061399a5761399a61530b565b6020026020010151816139ad9190614f2d565b9050826139b98161522c565b93505061396b565b9694955050505050565b60007f000000000000000000000000000000000000000000000000000000000000000015613a7757606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613a3a57600080fd5b505afa158015613a4e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a729190614696565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613afd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016128a8565b565b6015546040517f08b79da4000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff909116906308b79da490613b5a908690600401614c6d565b602060405180830381600087803b158015613b7457600080fd5b505af1158015613b88573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bac9190614696565b9050613bb9601382613eae565b5060608301516000828152600a6020908152604080832063ffffffff90941690935560a0860151600b8252929091208251613bfa93919291909101906140bc565b5092915050565b600061181b8383613eba565b73ffffffffffffffffffffffffffffffffffffffff8116331415613c8d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016128a8565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613d1a57613d1a61530b565b9060005260206000200154905092915050565b818180821415613d3e575050505050565b6000856002613d4d8787615032565b613d579190614f45565b613d619087614eb9565b81518110613d7157613d7161530b565b602002602001015190505b818313613e80575b80868481518110613d9757613d9761530b565b60200260200101511015613db75782613daf816151d1565b935050613d84565b858281518110613dc957613dc961530b565b6020026020010151811015613dea5781613de2816150e9565b925050613db7565b818313613e7b57858281518110613e0357613e0361530b565b6020026020010151868481518110613e1d57613e1d61530b565b6020026020010151878581518110613e3757613e3761530b565b60200260200101888581518110613e5057613e5061530b565b60209081029190910101919091525282613e69816151d1565b9350508180613e77906150e9565b9250505b613d7c565b81851215613e9357613e93868684613d2d565b83831215613ea657613ea6868486613d2d565b505050505050565b600061181b8383613fad565b60008181526001830160205260408120548015613fa3576000613ede6001836150a6565b8554909150600090613ef2906001906150a6565b9050818114613f57576000866000018281548110613f1257613f1261530b565b9060005260206000200154905080876000018481548110613f3557613f3561530b565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613f6857613f686152dc565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506111c6565b60009150506111c6565b6000818152600183016020526040812054613ff4575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556111c6565b5060006111c6565b508054600082559060005260206000209081019061373d9190614130565b8280546140269061517d565b90600052602060002090601f01602090048101928261404857600085556140ac565b82601f1061407f578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008235161785556140ac565b828001600101855582156140ac579182015b828111156140ac578235825591602001919060010190614091565b506140b8929150614130565b5090565b8280546140c89061517d565b90600052602060002090601f0160209004810192826140ea57600085556140ac565b82601f1061410357805160ff19168380011785556140ac565b828001600101855582156140ac579182015b828111156140ac578251825591602001919060010190614115565b5b808211156140b85760008155600101614131565b805161415081615369565b919050565b60008083601f84011261416757600080fd5b50813567ffffffffffffffff81111561417f57600080fd5b6020830191508360208260051b850101111561158f57600080fd5b8051801515811461415057600080fd5b60008083601f8401126141bc57600080fd5b50813567ffffffffffffffff8111156141d457600080fd5b60208301915083602082850101111561158f57600080fd5b600082601f8301126141fd57600080fd5b813561421061420b82614e73565b614e24565b81815284602083860101111561422557600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261425357600080fd5b815161426161420b82614e73565b81815284602083860101111561427657600080fd5b6142878260208301602087016150bd565b949350505050565b803561ffff8116811461415057600080fd5b80516141508161538b565b805167ffffffffffffffff8116811461415057600080fd5b803560ff8116811461415057600080fd5b80516141508161539d565b6000602082840312156142f257600080fd5b813561181b81615369565b6000806040838503121561431057600080fd5b823567ffffffffffffffff8082111561432857600080fd5b818501915085601f83011261433c57600080fd5b81356020828211156143505761435061533a565b8160051b61435f828201614e24565b8381528281019086840183880185018c101561437a57600080fd5b600093505b858410156143b95780358781111561439657600080fd5b6143a48d87838c01016141ec565b8452506001939093019291840191840161437f565b5097505050860135925050808211156143d157600080fd5b506143de858286016141ec565b9150509250929050565b600080602083850312156143fb57600080fd5b823567ffffffffffffffff81111561441257600080fd5b61441e85828601614155565b90969095509350505050565b60008060006040848603121561443f57600080fd5b833567ffffffffffffffff81111561445657600080fd5b61446286828701614155565b90945092505060208401356144768161538b565b809150509250925092565b60006020828403121561449357600080fd5b61181b8261419a565b6000602082840312156144ae57600080fd5b5035919050565b600080602083850312156144c857600080fd5b823567ffffffffffffffff8111156144df57600080fd5b61441e858286016141aa565b6000602082840312156144fd57600080fd5b815161181b81615369565b600080600080600060a0868803121561452057600080fd5b85516003811061452f57600080fd5b60208701519095506145408161538b565b60408701519094506145518161538b565b606087015190935061456281615369565b80925050608086015190509295509295909350565b60006020828403121561458957600080fd5b815167ffffffffffffffff808211156145a157600080fd5b9083019061014082860312156145b657600080fd5b6145be614dfa565b6145c783614145565b81526145d5602084016142a1565b60208201526040830151828111156145ec57600080fd5b6145f887828601614242565b60408301525061460a606084016142d5565b606082015261461b60808401614145565b608082015261462c60a084016142ac565b60a082015261463d60c084016142a1565b60c082015261464e60e084016142d5565b60e082015261010061466181850161419a565b90820152610120838101518381111561467957600080fd5b61468588828701614242565b918301919091525095945050505050565b6000602082840312156146a857600080fd5b5051919050565b6000806000604084860312156146c457600080fd5b83359250602084013567ffffffffffffffff8111156146e257600080fd5b6146ee868287016141aa565b9497909650939450505050565b6000806040838503121561470e57600080fd5b8235915061471e6020840161428f565b90509250929050565b60008060006060848603121561473c57600080fd5b8335925061474c6020850161428f565b9150604084013590509250925092565b6000806040838503121561476f57600080fd5b50508035926020909101359150565b60008060006060848603121561479357600080fd5b83359250602084013591506147aa6040850161428f565b90509250925092565b6000806000606084860312156147c857600080fd5b505081359360208301359350604090920135919050565b600080604083850312156147f257600080fd5b8235915060208301356148048161538b565b809150509250929050565b6000806040838503121561482257600080fd5b8235915060208301356148048161539d565b60006020828403121561484657600080fd5b61181b826142c4565b600080600080600060a0868803121561486757600080fd5b614870866142c4565b945060208601356148808161538b565b935060408601356148908161539d565b94979396509394606081013594506080013592915050565b6000602082840312156148ba57600080fd5b813561181b8161539d565b6000602082840312156148d757600080fd5b815161181b8161539d565b600081518084526020808501945080840160005b83811015614912578151875295820195908201906001016148f6565b509495945050505050565b600081518084526149358160208601602086016150bd565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8054600090600181811c908083168061498157607f831692505b60208084108214156149bc577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b8388528180156149d35760018114614a0557614a33565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008616828a0152604089019650614a33565b876000528160002060005b86811015614a2b5781548b8201850152908501908301614a10565b8a0183019750505b50505050505092915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015614ab4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552614aa286835161491d565b95509382019390820190600101614a68565b505085840381870152505050612c1d818561491d565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115614b0357600080fd5b8260051b80856040850137600092016040019182525092915050565b60208152600061181b60208301846148e2565b604081526000614b4560408301856148e2565b8281036020840152612c1d81856148e2565b8215158152604060208201526000611818604083018461491d565b60208152600061181b602083018461491d565b604081526000614b98604083018561491d565b8281036020840152612c1d818561491d565b60a081526000614bbd60a0830188614967565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015614c2f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552614c1d8383614967565b94860194925060019182019101614be4565b50508681036040880152614c43818b614967565b9450505050508460608401528281036080840152614c61818561491d565b98975050505050505050565b6020815260008251610100806020850152614c8c61012085018361491d565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe080868503016040870152614cc8848361491d565b935060408701519150614cf3606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a08701519150808685030160c0870152614d44848361491d565b935060c08701519150808685030160e087015250614d62838261491d565b92505060e0850151614d83828601826bffffffffffffffffffffffff169052565b5090949350505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b828152604060208201526000611818604083018461491d565b604051610140810167ffffffffffffffff81118282101715614e1e57614e1e61533a565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614e6b57614e6b61533a565b604052919050565b600067ffffffffffffffff821115614e8d57614e8d61533a565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000808212827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03841381151615614ef357614ef361527e565b827f8000000000000000000000000000000000000000000000000000000000000000038412811615614f2757614f2761527e565b50500190565b60008219821115614f4057614f4061527e565b500190565b600082614f5457614f546152ad565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615614fa857614fa861527e565b500590565b600082614fbc57614fbc6152ad565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614ff957614ff961527e565b500290565b60006bffffffffffffffffffffffff808316818516818304811182151516156150295761502961527e565b02949350505050565b6000808312837f80000000000000000000000000000000000000000000000000000000000000000183128115161561506c5761506c61527e565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0183138116156150a0576150a061527e565b50500390565b6000828210156150b8576150b861527e565b500390565b60005b838110156150d85781810151838201526020016150c0565b8381111561335b5750506000910152565b60007f800000000000000000000000000000000000000000000000000000000000000082141561511b5761511b61527e565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600061ffff8216806151555761515561527e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b600181811c9082168061519157607f821691505b602082108114156151cb577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156152035761520361527e565b5060010190565b600061ffff808316818114156152225761522261527e565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156152035761520361527e565b600060ff821660ff8114156152755761527561527e565b60010192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461373d57600080fd5b63ffffffff8116811461373d57600080fd5b6bffffffffffffffffffffffff8116811461373d57600080fdfea164736f6c6343000806000a",
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) FeedLabel(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "feedLabel")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) FeedLabel() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FeedLabel(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) FeedLabel() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.FeedLabel(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "feeds", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) Feeds(arg0 *big.Int) (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Feeds(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) Feeds(arg0 *big.Int) (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.Feeds(&_VerifiableLoadMercuryUpkeep.CallOpts, arg0)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) MercuryCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "mercuryCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) MercuryCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.MercuryCallback(&_VerifiableLoadMercuryUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) MercuryCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.MercuryCallback(&_VerifiableLoadMercuryUpkeep.CallOpts, values, extraData)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) QueryLabel(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "queryLabel")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) QueryLabel() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.QueryLabel(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) QueryLabel() (string, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.QueryLabel(&_VerifiableLoadMercuryUpkeep.CallOpts)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchRegisterUpkeeps", number, gasLimit, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, number, gasLimit, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, number, gasLimit, amount, checkGasToBurn, performGasToBurn)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchUpdateCheckData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchUpdateCheckData", upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchUpdateCheckData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchUpdateCheckData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchUpdateCheckData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchUpdateCheckData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) UpdateCheckData(opts *bind.TransactOpts, upkeepId *big.Int, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "updateCheckData", upkeepId, checkData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) UpdateCheckData(upkeepId *big.Int, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateCheckData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, checkData)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) UpdateCheckData(upkeepId *big.Int, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateCheckData(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, checkData)
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

	CheckDatas(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error)

	CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error)

	FeedLabel(opts *bind.CallOpts) (string, error)

	Feeds(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error)

	GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error)

	GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetDelaysLengthAtBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, error)

	GetDelaysLengthAtTimestampBucket(opts *bind.CallOpts, upkeepId *big.Int, timestampBucket uint16) (*big.Int, error)

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

	MercuryCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	MinBalanceThresholdMultiplier(opts *bind.CallOpts) (uint8, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	PerformDataSizes(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	PerformGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	PreviousPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	QueryLabel(opts *bind.CallOpts) (string, error)

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

	BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error)

	BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error)

	BatchUpdateCheckData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

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

	UpdateCheckData(opts *bind.TransactOpts, upkeepId *big.Int, checkData []byte) (*types.Transaction, error)

	WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepFundsAdded) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*VerifiableLoadMercuryUpkeepFundsAdded, error)

	FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepInsufficientFundsIterator, error)

	WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepInsufficientFunds) (event.Subscription, error)

	ParseInsufficientFunds(log types.Log) (*VerifiableLoadMercuryUpkeepInsufficientFunds, error)

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
