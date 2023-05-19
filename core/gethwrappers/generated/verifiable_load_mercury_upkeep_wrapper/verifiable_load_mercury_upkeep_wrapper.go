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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"registrarAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedLabel\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feedList\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"queryLabel\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"query\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"MercuryLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"origin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"v0\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ed\",\"type\":\"bytes\"}],\"name\":\"MercuryPerformEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"firstPerformBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"RegistrarSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkDatas\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxBucketedDelaysForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxDelayForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumTimestampBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTimestampBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getTimestampDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"mercuryCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queryLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractKeeperRegistry2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newFeedLabel\",\"type\":\"string\"}],\"name\":\"setFeedLabel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newQueryLabel\",\"type\":\"string\"}],\"name\":\"setQueryLabel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampBuckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60056018908155601980546001600160681b0319166c140000000002c68af0bb140000178155601960f21b60a05260e160f41b60c0526101408281527f4554482d5553442d415242495452554d2d544553544e455400000000000000006101605260e09081526101809283527f4254432d5553442d415242495452554d2d544553544e455400000000000000006101a052610100929092526102006040526101c09081527f555344432d5553442d415242495452554d2d544553544e4554000000000000006101e05261012052620000dc90601a90600362000483565b506040805160e08101909152604260608201818152829162005ef66080840139815260200160405180608001604052806042815260200162005f3860429139815260200160405180608001604052806042815260200162005f7a6042913990526200014c90601b90600362000483565b50604080518082019091526009808252683332b2b224a229ba3960b91b60209092019182526200017f91601c91620004e7565b5060408051808201909152600b8082526a313637b1b5a73ab6b132b960a91b6020909201918252620001b491601d91620004e7565b50348015620001c257600080fd5b5060405162005fbc38038062005fbc833981016040819052620001e59162000606565b818133806000816200023e5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b038481169190911790915581161562000271576200027181620003d7565b5050601580546001600160a01b0319166001600160a01b0385169081179091556040805163850af0cb60e01b815290516000935063850af0cb9160048082019260a092909190829003018186803b158015620002cc57600080fd5b505afa158015620002e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000307919062000670565b50601780546001600160a01b0319166001600160a01b038381169190911790915560155460408051631b6b6d2360e01b8152905193975091169450631b6b6d2393506004808201935060209291829003018186803b1580156200036957600080fd5b505afa1580156200037e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620003a4919062000649565b601680546001600160a01b0319166001600160a01b039290921691909117905550151560f81b6080525062000733915050565b6001600160a01b038116331415620004325760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000235565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215620004d5579160200282015b82811115620004d55782518051620004c4918491602090910190620004e7565b5091602001919060010190620004a4565b50620004e392915062000572565b5090565b828054620004f590620006e0565b90600052602060002090601f01602090048101928262000519576000855562000564565b82601f106200053457805160ff191683800117855562000564565b8280016001018555821562000564579182015b828111156200056457825182559160200191906001019062000547565b50620004e392915062000593565b80821115620004e3576000620005898282620005aa565b5060010162000572565b5b80821115620004e3576000815560010162000594565b508054620005b890620006e0565b6000825580601f10620005c9575050565b601f016020900490600052602060002090810190620005e9919062000593565b50565b805163ffffffff811681146200060157600080fd5b919050565b600080604083850312156200061a57600080fd5b825162000627816200071d565b602084015190925080151581146200063e57600080fd5b809150509250929050565b6000602082840312156200065c57600080fd5b815162000669816200071d565b9392505050565b600080600080600060a086880312156200068957600080fd5b8551600381106200069957600080fd5b9450620006a960208701620005ec565b9350620006b960408701620005ec565b92506060860151620006cb816200071d565b80925050608086015190509295509295909350565b600181811c90821680620006f557607f821691505b602082108114156200071757634e487b7160e01b600052602260045260246000fd5b50919050565b6001600160a01b0381168114620005e957600080fd5b60805160f81c60a05160f01c60c05160f01c61577562000781600039600081816107250152611b7a0152600081816105fb0152611c8e015260008181610a570152613bd001526157756000f3fe6080604052600436106105415760003560e01c806379ba5097116102af5780639fab438611610179578063d3558528116100d6578063f2fde38b1161008a578063fba7ffa31161006f578063fba7ffa31461117d578063fbfb4f76146111aa578063fcdc1f63146111ca57600080fd5b8063f2fde38b1461113d578063fb0ceb041461115d57600080fd5b8063dbef701e116100bb578063dbef701e146110da578063e0114adb146110fa578063e45530831461112757600080fd5b8063d35585281461105b578063d6051a72146110ba57600080fd5b8063af953a4a1161012d578063becde0e111610112578063becde0e114610fc1578063c357f1f314610fe1578063c80480221461103b57600080fd5b8063af953a4a14610f63578063b0971e1a14610f8357600080fd5b8063a6c60d891161015e578063a6c60d8914610ef6578063a72aa27e14610f16578063a79c404314610f3657600080fd5b80639fab438614610eb6578063a5f5893414610ed657600080fd5b80638fcb3fba116102275780639b429354116101db5780639bb86511116101c05780639bb8651114610e565780639d385eaa14610e765780639d6f1cc714610e9657600080fd5b80639b42935414610df85780639b51fb0d14610e2557600080fd5b8063948108f71161020c578063948108f714610d7c57806399cc6b0b14610d9c5780639ac542eb14610dbc57600080fd5b80638fcb3fba14610d2f5780639095aa3514610d5c57600080fd5b8063823783171161027e57806387dfa9001161026357806387dfa90014610cc45780638bc7b77214610ce45780638da5cb5b14610d0457600080fd5b80638237831714610c8457806386e330af14610ca457600080fd5b806379ba509714610c0d5780637b10399914610c225780637e4087b814610c4f57806380f4df1b14610c6f57600080fd5b80634ad8c9a61161040b578063643b34e9116103685780637145f11b1161031c57806373644cce1161030157806373644cce14610b935780637672130314610bc0578063776898c814610bed57600080fd5b80637145f11b14610b4357806371934a5214610b7357600080fd5b806369e9b7731161034d57806369e9b77314610ad65780636e04ff0d14610b035780637137a70214610b2357600080fd5b8063643b34e914610a8957806369cdbadb14610aa957600080fd5b80635d4ee7f3116103bf57806360457ff5116103a457806360457ff5146109d6578063636092e814610a03578063642f6cef14610a4557600080fd5b80635d4ee7f3146109a15780635f17e616146109b657600080fd5b806351c98be3116103f057806351c98be31461093457806357970e931461095457806358c52c041461098157600080fd5b80634ad8c9a6146108f15780634d6954451461091f57600080fd5b80632a9032d3116104b95780634585e33b1161046d57806345d2ec171161045257806345d2ec171461087757806346e7a63e146108975780634a5479f3146108c457600080fd5b80634585e33b1461083757806345cdd9d41461085757600080fd5b8063328ffd111161049e578063328ffd11146107b957806333774d1c146107e65780633ebe8d6c1461081757600080fd5b80632a9032d3146107475780632b20e3971461076757600080fd5b80631bee00801161051057806320e3dbd4116104f557806320e3dbd4146106d157806328c4b57b146106f357806329f0e4961461071357600080fd5b80631bee00801461066e578063206c32e81461069c57600080fd5b806306e3b63214610585578063077ac621146105bb57806312c55027146105e9578063177b0eb91461063057600080fd5b3661058057604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561059157600080fd5b506105a56105a0366004614acf565b6111f7565b6040516105b29190614e92565b60405180910390f35b3480156105c757600080fd5b506105db6105d6366004614a9a565b6112f3565b6040519081526020016105b2565b3480156105f557600080fd5b5061061d7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff90911681526020016105b2565b34801561063c57600080fd5b506105db61064b366004614a6e565b6000918252600f6020908152604080842061ffff93909316845291905290205490565b34801561067a57600080fd5b5061068e61068936600461480f565b611331565b6040516105b2929190614ea5565b3480156106a857600080fd5b506106bc6106b7366004614a6e565b61163a565b604080519283526020830191909152016105b2565b3480156106dd57600080fd5b506106f16106ec366004614599565b6116bd565b005b3480156106ff57600080fd5b506105db61070e366004614b26565b6118de565b34801561071f57600080fd5b5061061d7f000000000000000000000000000000000000000000000000000000000000000081565b34801561075357600080fd5b506106f161076236600461475b565b611949565b34801561077357600080fd5b506015546107949073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016105b2565b3480156107c557600080fd5b506105db6107d436600461480f565b60036020526000908152604090205481565b3480156107f257600080fd5b5061061d61080136600461480f565b60116020526000908152604090205461ffff1681565b34801561082357600080fd5b506105db61083236600461480f565b611a1c565b34801561084357600080fd5b506106f1610852366004614828565b611a85565b34801561086357600080fd5b506106f1610872366004614828565b6121a2565b34801561088357600080fd5b506105a5610892366004614a6e565b6121b3565b3480156108a357600080fd5b506105db6108b236600461480f565b600a6020526000908152604090205481565b3480156108d057600080fd5b506108e46108df36600461480f565b612222565b6040516105b29190614ee5565b3480156108fd57600080fd5b5061091161090c3660046145b6565b6122ce565b6040516105b2929190614eca565b34801561092b57600080fd5b506108e4612322565b34801561094057600080fd5b506106f161094f36600461479d565b61232f565b34801561096057600080fd5b506016546107949073ffffffffffffffffffffffffffffffffffffffff1681565b34801561098d57600080fd5b506108e461099c36600461480f565b6123d3565b3480156109ad57600080fd5b506106f16123ec565b3480156109c257600080fd5b506106f16109d1366004614acf565b612545565b3480156109e257600080fd5b506105db6109f136600461480f565b60076020526000908152604090205481565b348015610a0f57600080fd5b50601954610a28906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff90911681526020016105b2565b348015610a5157600080fd5b50610a797f000000000000000000000000000000000000000000000000000000000000000081565b60405190151581526020016105b2565b348015610a9557600080fd5b506106bc610aa4366004614acf565b6126b7565b348015610ab557600080fd5b506105db610ac436600461480f565b60086020526000908152604090205481565b348015610ae257600080fd5b506106f1610af1366004614acf565b60009182526008602052604090912055565b348015610b0f57600080fd5b50610911610b1e366004614828565b61283c565b348015610b2f57600080fd5b506105db610b3e366004614a9a565b6129e9565b348015610b4f57600080fd5b50610a79610b5e36600461480f565b600c6020526000908152604090205460ff1681565b348015610b7f57600080fd5b506106f1610b8e366004614828565b612a11565b348015610b9f57600080fd5b506105db610bae36600461480f565b6000908152600d602052604090205490565b348015610bcc57600080fd5b506105db610bdb36600461480f565b60046020526000908152604090205481565b348015610bf957600080fd5b50610a79610c0836600461480f565b612a1d565b348015610c1957600080fd5b506106f1612a6d565b348015610c2e57600080fd5b506017546107949073ffffffffffffffffffffffffffffffffffffffff1681565b348015610c5b57600080fd5b506106bc610c6a366004614acf565b612b6a565b348015610c7b57600080fd5b506108e4612ce2565b348015610c9057600080fd5b506105db610c9f366004614af1565b612cef565b348015610cb057600080fd5b506106f1610cbf366004614690565b612d6a565b348015610cd057600080fd5b506105db610cdf366004614af1565b612d7d565b348015610cf057600080fd5b5061068e610cff36600461480f565b612ded565b348015610d1057600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610794565b348015610d3b57600080fd5b506105db610d4a36600461480f565b60056020526000908152604090205481565b348015610d6857600080fd5b506106f1610d77366004614bc2565b612f62565b348015610d8857600080fd5b506106f1610d97366004614b82565b6131e2565b348015610da857600080fd5b506105a5610db7366004614a6e565b61337a565b348015610dc857600080fd5b50601954610de6906c01000000000000000000000000900460ff1681565b60405160ff90911681526020016105b2565b348015610e0457600080fd5b506106f1610e13366004614acf565b60009182526009602052604090912055565b348015610e3157600080fd5b5061061d610e4036600461480f565b60126020526000908152604090205461ffff1681565b348015610e6257600080fd5b506106f1610e7136600461475b565b6133e7565b348015610e8257600080fd5b506105a5610e9136600461480f565b6134b8565b348015610ea257600080fd5b506108e4610eb136600461480f565b61351a565b348015610ec257600080fd5b506106f1610ed1366004614a22565b61352a565b348015610ee257600080fd5b506105db610ef136600461480f565b6135cf565b348015610f0257600080fd5b506106f1610f1136600461480f565b601855565b348015610f2257600080fd5b506106f1610f31366004614b52565b613630565b348015610f4257600080fd5b506106f1610f51366004614acf565b60009182526007602052604090912055565b348015610f6f57600080fd5b506106f1610f7e36600461480f565b6136db565b348015610f8f57600080fd5b506105db610f9e366004614a6e565b6000918252600e6020908152604080842061ffff93909316845291905290205490565b348015610fcd57600080fd5b506106f1610fdc36600461475b565b613761565b348015610fed57600080fd5b506106f1610ffc366004614c1b565b601980547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b34801561104757600080fd5b506106f161105636600461480f565b6137fb565b34801561106757600080fd5b506106f1611076366004614ba7565b6019805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b3480156110c657600080fd5b506106bc6110d5366004614acf565b613893565b3480156110e657600080fd5b506105db6110f5366004614acf565b6138fc565b34801561110657600080fd5b506105db61111536600461480f565b60096020526000908152604090205481565b34801561113357600080fd5b506105db60185481565b34801561114957600080fd5b506106f1611158366004614599565b61392d565b34801561116957600080fd5b506105db611178366004614acf565b613941565b34801561118957600080fd5b506105db61119836600461480f565b60066020526000908152604090205481565b3480156111b657600080fd5b506106bc6111c5366004614a6e565b61395d565b3480156111d657600080fd5b506105db6111e536600461480f565b60026020526000908152604090205481565b6060600061120560136139d1565b9050808410611240576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826112525761124f848261543d565b92505b60008367ffffffffffffffff81111561126d5761126d6156eb565b604051908082528060200260200182016040528015611296578160200160208202803683370190505b50905060005b848110156112e8576112b96112b182886152c4565b6013906139db565b8282815181106112cb576112cb6156bc565b6020908102919091010152806112e0816155c3565b91505061129c565b509150505b92915050565b600e602052826000526040600020602052816000526040600020818154811061131b57600080fd5b9060005260206000200160009250925050505481565b606080600061134060136139d1565b905060008167ffffffffffffffff81111561135d5761135d6156eb565b604051908082528060200260200182016040528015611386578160200160208202803683370190505b50905060008267ffffffffffffffff8111156113a4576113a46156eb565b6040519080825280602002602001820160405280156113cd578160200160208202803683370190505b50905060005b8381101561162e5760006113e86013836139db565b9050808483815181106113fd576113fd6156bc565b6020908102919091018101919091526000828152601290915260408082205490517f3ebe8d6c0000000000000000000000000000000000000000000000000000000081526004810184905261ffff90911691903090633ebe8d6c9060240160206040518083038186803b15801561147357600080fd5b505afa158015611487573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ab9190614a09565b905060008167ffffffffffffffff8111156114c8576114c86156eb565b6040519080825280602002602001820160405280156114f1578160200160208202803683370190505b506000858152600e6020526040812091925090815b8561ffff168161ffff16116115eb5761ffff81166000908152602083815260408083208054825181850281018501909352808352919290919083018282801561156e57602002820191906000526020600020905b81548152602001906001019080831161155a575b5050505050905060005b81518110156115d657818181518110611593576115936156bc565b60200260200101518686806115a7906155c3565b9750815181106115b9576115b96156bc565b6020908102919091010152806115ce816155c3565b915050611578565b505080806115e3906155a1565b915050611506565b506115f7838e866139e7565b888881518110611609576116096156bc565b6020026020010181815250505050505050508080611626906155c3565b9150506113d3565b50909590945092505050565b6000828152600e6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493849392919083018282801561169e57602002820191906000526020600020905b81548152602001906001019080831161168a575b505050505090506116b0818251613b47565b92509250505b9250929050565b601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517f850af0cb00000000000000000000000000000000000000000000000000000000815290516000929163850af0cb9160048083019260a0929190829003018186803b15801561175257600080fd5b505afa158015611766573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061178a919061487b565b50601780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601554604080517f1b6b6d23000000000000000000000000000000000000000000000000000000008152905193975091169450631b6b6d2393506004808201935060209291829003018186803b15801561182957600080fd5b505afa15801561183d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611861919061485e565b601680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601554604051911681527f6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014906020015b60405180910390a15050565b6000838152600d60209081526040808320805482518185028101850190935280835261193f9383018282801561193357602002820191906000526020600020905b81548152602001906001019080831161191f575b505050505084846139e7565b90505b9392505050565b8060005b818160ff1610156119dd573063c8048022858560ff8516818110611973576119736156bc565b905060200201356040518263ffffffff1660e01b815260040161199891815260200190565b600060405180830381600087803b1580156119b257600080fd5b505af11580156119c6573d6000803e3d6000fd5b5050505080806119d59061560f565b91505061194d565b507fbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b8383604051611a0f929190614e3d565b60405180910390a1505050565b60008181526012602052604081205461ffff1681805b8261ffff168161ffff1611611a7d576000858152600e6020908152604080832061ffff85168452909152902054611a6990836152c4565b915080611a75816155a1565b915050611a32565b509392505050565b60005a9050600080611a99848601866145b6565b91509150600081806020019051810190611ab39190614a09565b60008181526005602090815260408083205460049092528220549293509190611ada613bcc565b905082611b145760008481526005602090815260408083208490556010825282208054600181018255908352912042910155915081611d74565b600084815260036020526040812054611b2d848461543d565b611b37919061543d565b6000868152601160209081526040808320546010909252909120805492935061ffff9091169182908110611b6d57611b6d6156bc565b90600052602060002001547f000000000000000000000000000000000000000000000000000000000000000061ffff1642611ba8919061543d565b1115611c175760008681526010602090815260408220805460018101825590835291204291015580611bd9816155a1565b600088815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559150505b600086815260126020908152604080832054600e835281842061ffff9091168085529083528184208054835181860281018601909452808452919493909190830182828015611c8557602002820191906000526020600020905b815481526020019060010190808311611c71575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff1681511415611d015781611cc3816155a1565b60008a815260126020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600e6020908152604080832061ffff94851684528252808320805460018181018355918552838520018790558a8452600f83528184209590941683529381528382208054808501825590835281832001859055888252600d81529281208054928301815581529190912001555b600084815260066020526040812054611d8e9060016152c4565b6000868152600660209081526040918290208390558151878152908101859052908101859052606081018290529091507f6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede39060800160405180910390a16000858152600460209081526040808320859055601854600290925290912054611e15908461543d565b11156120c0576017546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810187905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a9060240160006040518083038186803b158015611e8657600080fd5b505afa158015611e9a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611ee091908101906148ea565b6017546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810189905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c9060240160206040518083038186803b158015611f5057600080fd5b505afa158015611f64573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f889190614c38565b601954909150611fac9082906c01000000000000000000000000900460ff16615395565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff1610156120bd576019546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018990526bffffffffffffffffffffffff9091166024820152309063948108f790604401600060405180830381600087803b15801561203d57600080fd5b505af1158015612051573d6000803e3d6000fd5b50505060008881526002602090815260409182902087905560195482518b81526bffffffffffffffffffffffff909116918101919091529081018690527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0915060600160405180910390a15b50505b6000858152600760205260409020545b805a6120dc908b61543d565b6120e8906127106152c4565b10156121295782406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556120d0565b82863273ffffffffffffffffffffffffffffffffffffffff167fcad583be2d908a590c81c7e332cf11c7a4ea41ecf1e059efac3ea7e83e34f1a58b600081518110612176576121766156bc565b60200260200101518b60405161218d929190614ef8565b60405180910390a45050505050505050505050565b6121ae601c83836141fd565b505050565b6000828152600e6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561221557602002820191906000526020600020905b815481526020019060010190808311612201575b5050505050905092915050565b601a818154811061223257600080fd5b90600052602060002001600091509050805461224d90615514565b80601f016020809104026020016040519081016040528092919081815260200182805461227990615514565b80156122c65780601f1061229b576101008083540402835291602001916122c6565b820191906000526020600020905b8154815290600101906020018083116122a957829003601f168201915b505050505081565b60006060600084846040516020016122e7929190614db2565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b601c805461224d90615514565b8160005b818110156123cc5730635f17e616868684818110612353576123536156bc565b90506020020135856040518363ffffffff1660e01b815260040161238792919091825263ffffffff16602082015260400190565b600060405180830381600087803b1580156123a157600080fd5b505af11580156123b5573d6000803e3d6000fd5b5050505080806123c4906155c3565b915050612333565b5050505050565b600b602052600090815260409020805461224d90615514565b6123f4613c7d565b6016546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561245e57600080fd5b505afa158015612472573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124969190614a09565b6016546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb90604401602060405180830381600087803b15801561250957600080fd5b505af115801561251d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061254191906147f4565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600d909152812061257d9161429f565b60008281526012602052604081205461ffff16905b8161ffff168161ffff16116125d9576000848152600e6020908152604080832061ffff8516845290915281206125c79161429f565b806125d1816155a1565b915050612592565b5050600082815260126020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055601190915281205461ffff16905b8161ffff168161ffff1611612667576000848152600f6020908152604080832061ffff8516845290915281206126559161429f565b8061265f816155a1565b915050612620565b50600083815260106020526040812061267f9161429f565b5050600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6040517f3ebe8d6c00000000000000000000000000000000000000000000000000000000815260048101839052600090819081903090633ebe8d6c9060240160206040518083038186803b15801561270e57600080fd5b505afa158015612722573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127469190614a09565b90508315806127555750808410155b1561275e578093505b60008581526012602052604081205485919061ffff16805b6000898152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835291929091908301828280156127d857602002820191906000526020600020905b8154815260200190600101908083116127c4575b505050505090506000806127ec8388613b47565b90925090506127fb82876152c4565b9550612807818861543d565b9650600087116128195750505061282f565b5050508080612827906154d8565b915050612776565b5090979596505050505050565b6000606060005a905060006128538587018761480f565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff81111561288c5761288c6156eb565b6040519080825280601f01601f1916602001820160405280156128b6576020820181803683370190505b506040516020016128c8929190615154565b604051602081830303815290604052905060006128e3613bcc565b905060006128f086612a1d565b90505b835a6128ff908961543d565b61290b906127106152c4565b101561294c5781406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556128f3565b806129645760008398509850505050505050506116b6565b601c601a601d848960405160200161297e91815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f62e8a50d0000000000000000000000000000000000000000000000000000000082526129e09594939291600401614f1d565b60405180910390fd5b600f602052826000526040600020602052816000526040600020818154811061131b57600080fd5b6121ae601d83836141fd565b600081815260056020526040812054612a3857506001919050565b600082815260036020908152604080832054600490925290912054612a5b613bcc565b612a65919061543d565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612aee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016129e0565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040517fa5f589340000000000000000000000000000000000000000000000000000000081526004810183905260009081908190309063a5f589349060240160206040518083038186803b158015612bc157600080fd5b505afa158015612bd5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612bf99190614a09565b9050831580612c085750808410155b15612c11578093505b60008581526011602052604081205485919061ffff16805b6000898152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083529192909190830182828015612c8b57602002820191906000526020600020905b815481526020019060010190808311612c77575b50505050509050600080612c9f8388613b47565b9092509050612cae82876152c4565b9550612cba818861543d565b965060008711612ccc5750505061282f565b5050508080612cda906154d8565b915050612c29565b601d805461224d90615514565b6000838152600e6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493830182828015612d4e57602002820191906000526020600020905b815481526020019060010190808311612d3a575b50505050509050612d61818583516139e7565b95945050505050565b805161254190601a9060208401906142bd565b6000838152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493830182828015612d4e5760200282019190600052602060002090815481526020019060010190808311612d3a5750505050509050612d61818583516139e7565b6060806000612dfc60136139d1565b905060008167ffffffffffffffff811115612e1957612e196156eb565b604051908082528060200260200182016040528015612e42578160200160208202803683370190505b50905060008267ffffffffffffffff811115612e6057612e606156eb565b604051908082528060200260200182016040528015612e89578160200160208202803683370190505b50905060005b8381101561162e576000612ea46013836139db565b6000818152600d6020908152604080832080548251818502810185019093528083529495509293909291830182828015612efd57602002820191906000526020600020905b815481526020019060010190808311612ee9575b5050505050905081858481518110612f1757612f176156bc565b602002602001018181525050612f2f818a83516139e7565b848481518110612f4157612f416156bc565b60200260200101818152505050508080612f5a906155c3565b915050612e8f565b6040805161014081018252600461010082019081527f746573740000000000000000000000000000000000000000000000000000000061012083015281528151602081810184526000808352818401929092523083850181905263ffffffff8916606085015260808401528351808201855282815260a08401528351908101909352825260c08101919091526bffffffffffffffffffffffff841660e082015260165460155473ffffffffffffffffffffffffffffffffffffffff9182169163095ea7b3911661303560ff8a1688615395565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff166024820152604401602060405180830381600087803b1580156130ae57600080fd5b505af11580156130c2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130e691906147f4565b5060008660ff1667ffffffffffffffff811115613105576131056156eb565b60405190808252806020026020018201604052801561312e578160200160208202803683370190505b50905060005b8760ff168160ff1610156131a157600061314d84613d00565b905080838360ff1681518110613165576131656156bc565b602090810291909101810191909152600091825260088152604080832088905560079091529020849055806131998161560f565b915050613134565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c711816040516131d19190614e92565b60405180910390a150505050505050565b6016546017546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b390604401602060405180830381600087803b15801561326557600080fd5b505af1158015613279573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061329d91906147f4565b506017546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b15801561331e57600080fd5b505af1158015613332573d6000803e3d6000fd5b5050604080518581526bffffffffffffffffffffffff851660208201527f8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e93500190506118d2565b6000828152600f6020908152604080832061ffff851684528252918290208054835181840281018401909452808452606093928301828280156122155760200282019190600052602060002090815481526020019060010190808311612201575050505050905092915050565b8060005b818110156134b2576000848483818110613407576134076156bc565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16639fab4386828360405160200161344091815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b815260040161346c929190615154565b600060405180830381600087803b15801561348657600080fd5b505af115801561349a573d6000803e3d6000fd5b505050505080806134aa906155c3565b9150506133eb565b50505050565b6000818152600d602090815260409182902080548351818402810184019094528084526060939283018282801561350e57602002820191906000526020600020905b8154815260200190600101908083116134fa575b50505050509050919050565b601b818154811061223257600080fd5b6017546040517f9fab438600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690639fab43869061358490869086908690600401615100565b600060405180830381600087803b15801561359e57600080fd5b505af11580156135b2573d6000803e3d6000fd5b5050506000848152600b602052604090206134b2915083836141fd565b60008181526011602052604081205461ffff1681805b8261ffff168161ffff1611611a7d576000858152600f6020908152604080832061ffff8516845290915290205461361c90836152c4565b915080613628816155a1565b9150506135e5565b6017546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156136a857600080fd5b505af11580156136bc573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6017546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561374d57600080fd5b505af11580156123cc573d6000803e3d6000fd5b8060005b818163ffffffff1610156134b2573063af953a4a858563ffffffff8516818110613791576137916156bc565b905060200201356040518263ffffffff1660e01b81526004016137b691815260200190565b600060405180830381600087803b1580156137d057600080fd5b505af11580156137e4573d6000803e3d6000fd5b5050505080806137f3906155f5565b915050613765565b6017546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b15801561386757600080fd5b505af115801561387b573d6000803e3d6000fd5b50505050612541816013613e0290919063ffffffff16565b6000828152600d602090815260408083208054825181850281018501909352808352849384939291908301828280156138eb57602002820191906000526020600020905b8154815260200190600101908083116138d7575b505050505090506116b08185613b47565b600d602052816000526040600020818154811061391857600080fd5b90600052602060002001600091509150505481565b613935613c7d565b61393e81613e0e565b50565b6010602052816000526040600020818154811061391857600080fd5b6000828152600f6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493849392919083018282801561169e576020028201919060005260206000209081548152602001906001019080831161168a57505050505090506116b0818251613b47565b60006112ed825490565b60006119428383613f04565b825160009081908315806139fb5750808410155b15613a04578093505b60008467ffffffffffffffff811115613a1f57613a1f6156eb565b604051908082528060200260200182016040528015613a48578160200160208202803683370190505b509050600092505b84831015613ab657866001613a65858561543d565b613a6f919061543d565b81518110613a7f57613a7f6156bc565b6020026020010151818481518110613a9957613a996156bc565b602090810291909101015282613aae816155c3565b935050613a50565b613acf81600060018451613aca919061543d565b613f2e565b8560641415613b09578060018251613ae7919061543d565b81518110613af757613af76156bc565b60200260200101519350505050611942565b806064825188613b199190615358565b613b239190615344565b81518110613b3357613b336156bc565b602002602001015193505050509392505050565b815160009081908190841580613b5d5750808510155b15613b66578094505b60008092505b85831015613bc257866001613b81858561543d565b613b8b919061543d565b81518110613b9b57613b9b6156bc565b602002602001015181613bae91906152c4565b905082613bba816155c3565b935050613b6c565b9694955050505050565b60007f000000000000000000000000000000000000000000000000000000000000000015613c7857606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613c3b57600080fd5b505afa158015613c4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c739190614a09565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613cfe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016129e0565b565b6015546040517f08b79da4000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff909116906308b79da490613d5b908690600401614fe0565b602060405180830381600087803b158015613d7557600080fd5b505af1158015613d89573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613dad9190614a09565b9050613dba6013826140af565b5060608301516000828152600a6020908152604080832063ffffffff90941690935560a0860151600b8252929091208251613dfb9391929190910190614316565b5092915050565b600061194283836140bb565b73ffffffffffffffffffffffffffffffffffffffff8116331415613e8e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016129e0565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613f1b57613f1b6156bc565b9060005260206000200154905092915050565b818180821415613f3f575050505050565b6000856002613f4e87876153c9565b613f5891906152dc565b613f629087615250565b81518110613f7257613f726156bc565b602002602001015190505b818313614081575b80868481518110613f9857613f986156bc565b60200260200101511015613fb85782613fb081615568565b935050613f85565b858281518110613fca57613fca6156bc565b6020026020010151811015613feb5781613fe381615480565b925050613fb8565b81831361407c57858281518110614004576140046156bc565b602002602001015186848151811061401e5761401e6156bc565b6020026020010151878581518110614038576140386156bc565b60200260200101888581518110614051576140516156bc565b6020908102919091010191909152528261406a81615568565b935050818061407890615480565b9250505b613f7d565b8185121561409457614094868684613f2e565b838312156140a7576140a7868486613f2e565b505050505050565b600061194283836141ae565b600081815260018301602052604081205480156141a45760006140df60018361543d565b85549091506000906140f39060019061543d565b9050818114614158576000866000018281548110614113576141136156bc565b9060005260206000200154905080876000018481548110614136576141366156bc565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806141695761416961568d565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506112ed565b60009150506112ed565b60008181526001830160205260408120546141f5575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556112ed565b5060006112ed565b82805461420990615514565b90600052602060002090601f01602090048101928261422b576000855561428f565b82601f10614262578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082351617855561428f565b8280016001018555821561428f579182015b8281111561428f578235825591602001919060010190614274565b5061429b92915061438a565b5090565b508054600082559060005260206000209081019061393e919061438a565b82805482825590600052602060002090810192821561430a579160200282015b8281111561430a57825180516142fa918491602090910190614316565b50916020019190600101906142dd565b5061429b92915061439f565b82805461432290615514565b90600052602060002090601f016020900481019282614344576000855561428f565b82601f1061435d57805160ff191683800117855561428f565b8280016001018555821561428f579182015b8281111561428f57825182559160200191906001019061436f565b5b8082111561429b576000815560010161438b565b8082111561429b5760006143b382826143bc565b5060010161439f565b5080546143c890615514565b6000825580601f106143d8575050565b601f01602090049060005260206000209081019061393e919061438a565b60006144096144048461520a565b615197565b905082815283838301111561441d57600080fd5b828260208301376000602084830101529392505050565b805161443f8161571a565b919050565b60008083601f84011261445657600080fd5b50813567ffffffffffffffff81111561446e57600080fd5b6020830191508360208260051b85010111156116b657600080fd5b8051801515811461443f57600080fd5b60008083601f8401126144ab57600080fd5b50813567ffffffffffffffff8111156144c357600080fd5b6020830191508360208285010111156116b657600080fd5b600082601f8301126144ec57600080fd5b611942838335602085016143f6565b600082601f83011261450c57600080fd5b815161451a6144048261520a565b81815284602083860101111561452f57600080fd5b614540826020830160208701615454565b949350505050565b803561ffff8116811461443f57600080fd5b805161443f8161573c565b805167ffffffffffffffff8116811461443f57600080fd5b803560ff8116811461443f57600080fd5b805161443f8161574e565b6000602082840312156145ab57600080fd5b81356119428161571a565b600080604083850312156145c957600080fd5b823567ffffffffffffffff808211156145e157600080fd5b818501915085601f8301126145f557600080fd5b81356020614605614404836151e6565b8083825282820191508286018a848660051b890101111561462557600080fd5b60005b858110156146605781358781111561463f57600080fd5b61464d8d87838c01016144db565b8552509284019290840190600101614628565b5090975050508601359250508082111561467957600080fd5b50614686858286016144db565b9150509250929050565b600060208083850312156146a357600080fd5b823567ffffffffffffffff808211156146bb57600080fd5b818501915085601f8301126146cf57600080fd5b81356146dd614404826151e6565b80828252858201915085850189878560051b88010111156146fd57600080fd5b60005b8481101561474c5781358681111561471757600080fd5b8701603f81018c1361472857600080fd5b6147398c8a830135604084016143f6565b8552509287019290870190600101614700565b50909998505050505050505050565b6000806020838503121561476e57600080fd5b823567ffffffffffffffff81111561478557600080fd5b61479185828601614444565b90969095509350505050565b6000806000604084860312156147b257600080fd5b833567ffffffffffffffff8111156147c957600080fd5b6147d586828701614444565b90945092505060208401356147e98161573c565b809150509250925092565b60006020828403121561480657600080fd5b61194282614489565b60006020828403121561482157600080fd5b5035919050565b6000806020838503121561483b57600080fd5b823567ffffffffffffffff81111561485257600080fd5b61479185828601614499565b60006020828403121561487057600080fd5b81516119428161571a565b600080600080600060a0868803121561489357600080fd5b8551600381106148a257600080fd5b60208701519095506148b38161573c565b60408701519094506148c48161573c565b60608701519093506148d58161571a565b80925050608086015190509295509295909350565b6000602082840312156148fc57600080fd5b815167ffffffffffffffff8082111561491457600080fd5b90830190610140828603121561492957600080fd5b61493161516d565b61493a83614434565b81526149486020840161455a565b602082015260408301518281111561495f57600080fd5b61496b878286016144fb565b60408301525061497d6060840161458e565b606082015261498e60808401614434565b608082015261499f60a08401614565565b60a08201526149b060c0840161455a565b60c08201526149c160e0840161458e565b60e08201526101006149d4818501614489565b9082015261012083810151838111156149ec57600080fd5b6149f8888287016144fb565b918301919091525095945050505050565b600060208284031215614a1b57600080fd5b5051919050565b600080600060408486031215614a3757600080fd5b83359250602084013567ffffffffffffffff811115614a5557600080fd5b614a6186828701614499565b9497909650939450505050565b60008060408385031215614a8157600080fd5b82359150614a9160208401614548565b90509250929050565b600080600060608486031215614aaf57600080fd5b83359250614abf60208501614548565b9150604084013590509250925092565b60008060408385031215614ae257600080fd5b50508035926020909101359150565b600080600060608486031215614b0657600080fd5b8335925060208401359150614b1d60408501614548565b90509250925092565b600080600060608486031215614b3b57600080fd5b505081359360208301359350604090920135919050565b60008060408385031215614b6557600080fd5b823591506020830135614b778161573c565b809150509250929050565b60008060408385031215614b9557600080fd5b823591506020830135614b778161574e565b600060208284031215614bb957600080fd5b6119428261457d565b600080600080600060a08688031215614bda57600080fd5b614be38661457d565b94506020860135614bf38161573c565b93506040860135614c038161574e565b94979396509394606081013594506080013592915050565b600060208284031215614c2d57600080fd5b81356119428161574e565b600060208284031215614c4a57600080fd5b81516119428161574e565b600081518084526020808501945080840160005b83811015614c8557815187529582019590820190600101614c69565b509495945050505050565b60008151808452614ca8816020860160208601615454565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8054600090600181811c9080831680614cf457607f831692505b6020808410821415614d2f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b838852818015614d465760018114614d7857614da6565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008616828a0152604089019650614da6565b876000528160002060005b86811015614d9e5781548b8201850152908501908301614d83565b8a0183019750505b50505050505092915050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015614e27577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552614e15868351614c90565b95509382019390820190600101614ddb565b505085840381870152505050612d618185614c90565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115614e7657600080fd5b8260051b80856040850137600092016040019182525092915050565b6020815260006119426020830184614c55565b604081526000614eb86040830185614c55565b8281036020840152612d618185614c55565b821515815260406020820152600061193f6040830184614c90565b6020815260006119426020830184614c90565b604081526000614f0b6040830185614c90565b8281036020840152612d618185614c90565b60a081526000614f3060a0830188614cda565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015614fa2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0878403018552614f908383614cda565b94860194925060019182019101614f57565b50508681036040880152614fb6818b614cda565b9450505050508460608401528281036080840152614fd48185614c90565b98975050505050505050565b6020815260008251610100806020850152614fff610120850183614c90565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08086850301604087015261503b8483614c90565b935060408701519150615066606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a08701519150808685030160c08701526150b78483614c90565b935060c08701519150808685030160e0870152506150d58382614c90565b92505060e08501516150f6828601826bffffffffffffffffffffffff169052565b5090949350505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b82815260406020820152600061193f6040830184614c90565b604051610140810167ffffffffffffffff81118282101715615191576151916156eb565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156151de576151de6156eb565b604052919050565b600067ffffffffffffffff821115615200576152006156eb565b5060051b60200190565b600067ffffffffffffffff821115615224576152246156eb565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000808212827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0384138115161561528a5761528a61562f565b827f80000000000000000000000000000000000000000000000000000000000000000384128116156152be576152be61562f565b50500190565b600082198211156152d7576152d761562f565b500190565b6000826152eb576152eb61565e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561533f5761533f61562f565b500590565b6000826153535761535361565e565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156153905761539061562f565b500290565b60006bffffffffffffffffffffffff808316818516818304811182151516156153c0576153c061562f565b02949350505050565b6000808312837f8000000000000000000000000000000000000000000000000000000000000000018312811516156154035761540361562f565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0183138116156154375761543761562f565b50500390565b60008282101561544f5761544f61562f565b500390565b60005b8381101561546f578181015183820152602001615457565b838111156134b25750506000910152565b60007f80000000000000000000000000000000000000000000000000000000000000008214156154b2576154b261562f565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600061ffff8216806154ec576154ec61562f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b600181811c9082168061552857607f821691505b60208210811415615562577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561559a5761559a61562f565b5060010190565b600061ffff808316818114156155b9576155b961562f565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561559a5761559a61562f565b600063ffffffff808316818114156155b9576155b961562f565b600060ff821660ff8114156156265761562661562f565b60010192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461393e57600080fd5b63ffffffff8116811461393e57600080fd5b6bffffffffffffffffffffffff8116811461393e57600080fdfea164736f6c6343000806000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307835353533343434333264353535333434326434313532343234393534353235353464326435343435353335343465343535343030303030303030303030303030",
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetFeedLabel(opts *bind.TransactOpts, newFeedLabel string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setFeedLabel", newFeedLabel)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetFeedLabel(newFeedLabel string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeedLabel(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeedLabel)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetFeedLabel(newFeedLabel string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeedLabel(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeedLabel)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setFeeds", newFeeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetFeeds(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeeds(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetFeeds(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeeds(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeeds)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetQueryLabel(opts *bind.TransactOpts, newQueryLabel string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setQueryLabel", newQueryLabel)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetQueryLabel(newQueryLabel string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetQueryLabel(&_VerifiableLoadMercuryUpkeep.TransactOpts, newQueryLabel)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetQueryLabel(newQueryLabel string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetQueryLabel(&_VerifiableLoadMercuryUpkeep.TransactOpts, newQueryLabel)
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

	BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	CancelUpkeep(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, checkData []byte) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SetAddLinkAmount(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	SetCheckGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error)

	SetFeedLabel(opts *bind.TransactOpts, newFeedLabel string) (*types.Transaction, error)

	SetFeeds(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetMinBalanceThresholdMultiplier(opts *bind.TransactOpts, newMinBalanceThresholdMultiplier uint8) (*types.Transaction, error)

	SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetPerformGasToBurn(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetQueryLabel(opts *bind.TransactOpts, newQueryLabel string) (*types.Transaction, error)

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
