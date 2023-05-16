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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"registrarAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"firstPerformBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"RegistrarSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkDatas\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxBucketedDelaysForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxDelayForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumTimestampBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTimestampBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getTimestampDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractKeeperRegistry2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampBuckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040526005601855601980546001600160681b0319166c140000000002c68af0bb140000179055601960f21b60a05260e160f41b60c0523480156200004557600080fd5b50604051620052cd380380620052cd833981016040819052620000689162000320565b81813380600081620000c15760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000f457620000f4816200025a565b5050601580546001600160a01b0319166001600160a01b0385169081179091556040805163850af0cb60e01b815290516000935063850af0cb9160048082019260a092909190829003018186803b1580156200014f57600080fd5b505afa15801562000164573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200018a91906200038a565b50601780546001600160a01b0319166001600160a01b038381169190911790915560155460408051631b6b6d2360e01b8152905193975091169450631b6b6d2393506004808201935060209291829003018186803b158015620001ec57600080fd5b505afa15801562000201573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000227919062000363565b601680546001600160a01b0319166001600160a01b039290921691909117905550151560f81b6080525062000413915050565b6001600160a01b038116331415620002b55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000b8565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b805163ffffffff811681146200031b57600080fd5b919050565b600080604083850312156200033457600080fd5b82516200034181620003fa565b602084015190925080151581146200035857600080fd5b809150509250929050565b6000602082840312156200037657600080fd5b81516200038381620003fa565b9392505050565b600080600080600060a08688031215620003a357600080fd5b855160038110620003b357600080fd5b9450620003c36020870162000306565b9350620003d36040870162000306565b92506060860151620003e581620003fa565b80925050608086015190509295509295909350565b6001600160a01b03811681146200041057600080fd5b50565b60805160f81c60a05160f01c60c05160f01c614e6c620004616000396000818161063201526119630152600081816105080152611a770152600081816108e1015261376a0152614e6c6000f3fe60806040526004361061044e5760003560e01c806379ba509711610243578063a5f5893411610143578063d6051a72116100bb578063f2fde38b1161008a578063fba7ffa31161006f578063fba7ffa314610f80578063fbfb4f7614610fad578063fcdc1f6314610fcd57600080fd5b8063f2fde38b14610f40578063fb0ceb0414610f6057600080fd5b8063d6051a7214610ebd578063dbef701e14610edd578063e0114adb14610efd578063e455308314610f2a57600080fd5b8063af953a4a11610112578063c357f1f3116100f7578063c357f1f314610de4578063c804802214610e3e578063d355852814610e5e57600080fd5b8063af953a4a14610d86578063b0971e1a14610da657600080fd5b8063a5f5893414610cf9578063a6c60d8914610d19578063a72aa27e14610d39578063a79c404314610d5957600080fd5b80639095aa35116101d65780639b429354116101a55780639bb865111161018a5780639bb8651114610c995780639d385eaa14610cb95780639fab438614610cd957600080fd5b80639b42935414610c3b5780639b51fb0d14610c6857600080fd5b80639095aa3514610b9f578063948108f714610bbf57806399cc6b0b14610bdf5780639ac542eb14610bff57600080fd5b806387dfa9001161021257806387dfa90014610b075780638bc7b77214610b275780638da5cb5b14610b475780638fcb3fba14610b7257600080fd5b806379ba509714610a855780637b10399914610a9a5780637e4087b814610ac75780638237831714610ae757600080fd5b806351c98be31161034e578063643b34e9116102e15780637137a702116102b057806373644cce1161029557806373644cce14610a0b5780637672130314610a38578063776898c814610a6557600080fd5b80637137a702146109bb5780637145f11b146109db57600080fd5b8063643b34e91461091357806369cdbadb1461093357806369e9b773146109605780636e04ff0d1461098d57600080fd5b80635f17e6161161031d5780635f17e6161461084057806360457ff514610860578063636092e81461088d578063642f6cef146108cf57600080fd5b806351c98be3146107b157806357970e93146107d157806358c52c04146107fe5780635d4ee7f31461082b57600080fd5b806329f0e496116103e157806333774d1c116103b05780634585e33b116103955780634585e33b1461074457806345d2ec171461076457806346e7a63e1461078457600080fd5b806333774d1c146106f35780633ebe8d6c1461072457600080fd5b806329f0e496146106205780632a9032d3146106545780632b20e39714610674578063328ffd11146106c657600080fd5b80631bee00801161041d5780631bee00801461057b578063206c32e8146105a957806320e3dbd4146105de57806328c4b57b1461060057600080fd5b806306e3b63214610492578063077ac621146104c857806312c55027146104f6578063177b0eb91461053d57600080fd5b3661048d57604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b34801561049e57600080fd5b506104b26104ad366004614440565b610ffa565b6040516104bf91906146a0565b60405180910390f35b3480156104d457600080fd5b506104e86104e336600461440b565b6110f6565b6040519081526020016104bf565b34801561050257600080fd5b5061052a7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff90911681526020016104bf565b34801561054957600080fd5b506104e86105583660046143df565b6000918252600f6020908152604080842061ffff93909316845291905290205490565b34801561058757600080fd5b5061059b6105963660046140fb565b611134565b6040516104bf9291906146b3565b3480156105b557600080fd5b506105c96105c43660046143df565b61143d565b604080519283526020830191909152016104bf565b3480156105ea57600080fd5b506105fe6105f936600461402a565b6114c0565b005b34801561060c57600080fd5b506104e861061b366004614497565b6116e1565b34801561062c57600080fd5b5061052a7f000000000000000000000000000000000000000000000000000000000000000081565b34801561066057600080fd5b506105fe61066f366004614047565b61174c565b34801561068057600080fd5b506015546106a19073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016104bf565b3480156106d257600080fd5b506104e86106e13660046140fb565b60036020526000908152604090205481565b3480156106ff57600080fd5b5061052a61070e3660046140fb565b60116020526000908152604090205461ffff1681565b34801561073057600080fd5b506104e861073f3660046140fb565b61181f565b34801561075057600080fd5b506105fe61075f366004614114565b611888565b34801561077057600080fd5b506104b261077f3660046143df565b611f2a565b34801561079057600080fd5b506104e861079f3660046140fb565b600a6020526000908152604090205481565b3480156107bd57600080fd5b506105fe6107cc366004614089565b611f99565b3480156107dd57600080fd5b506016546106a19073ffffffffffffffffffffffffffffffffffffffff1681565b34801561080a57600080fd5b5061081e6108193660046140fb565b61203d565b6040516104bf91906146f3565b34801561083757600080fd5b506105fe6120d7565b34801561084c57600080fd5b506105fe61085b366004614440565b612230565b34801561086c57600080fd5b506104e861087b3660046140fb565b60076020526000908152604090205481565b34801561089957600080fd5b506019546108b2906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff90911681526020016104bf565b3480156108db57600080fd5b506109037f000000000000000000000000000000000000000000000000000000000000000081565b60405190151581526020016104bf565b34801561091f57600080fd5b506105c961092e366004614440565b6123a2565b34801561093f57600080fd5b506104e861094e3660046140fb565b60086020526000908152604090205481565b34801561096c57600080fd5b506105fe61097b366004614440565b60009182526008602052604090912055565b34801561099957600080fd5b506109ad6109a8366004614114565b612527565b6040516104bf9291906146d8565b3480156109c757600080fd5b506104e86109d636600461440b565b612654565b3480156109e757600080fd5b506109036109f63660046140fb565b600c6020526000908152604090205460ff1681565b348015610a1757600080fd5b506104e8610a263660046140fb565b6000908152600d602052604090205490565b348015610a4457600080fd5b506104e8610a533660046140fb565b60046020526000908152604090205481565b348015610a7157600080fd5b50610903610a803660046140fb565b61267c565b348015610a9157600080fd5b506105fe6126cc565b348015610aa657600080fd5b506017546106a19073ffffffffffffffffffffffffffffffffffffffff1681565b348015610ad357600080fd5b506105c9610ae2366004614440565b6127ce565b348015610af357600080fd5b506104e8610b02366004614462565b612946565b348015610b1357600080fd5b506104e8610b22366004614462565b6129c1565b348015610b3357600080fd5b5061059b610b423660046140fb565b612a31565b348015610b5357600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166106a1565b348015610b7e57600080fd5b506104e8610b8d3660046140fb565b60056020526000908152604090205481565b348015610bab57600080fd5b506105fe610bba366004614533565b612ba6565b348015610bcb57600080fd5b506105fe610bda3660046144f3565b612e26565b348015610beb57600080fd5b506104b2610bfa3660046143df565b612fbe565b348015610c0b57600080fd5b50601954610c29906c01000000000000000000000000900460ff1681565b60405160ff90911681526020016104bf565b348015610c4757600080fd5b506105fe610c56366004614440565b60009182526009602052604090912055565b348015610c7457600080fd5b5061052a610c833660046140fb565b60126020526000908152604090205461ffff1681565b348015610ca557600080fd5b506105fe610cb4366004614047565b61302b565b348015610cc557600080fd5b506104b2610cd43660046140fb565b6130fc565b348015610ce557600080fd5b506105fe610cf436600461430e565b61315e565b348015610d0557600080fd5b506104e8610d143660046140fb565b613203565b348015610d2557600080fd5b506105fe610d343660046140fb565b601855565b348015610d4557600080fd5b506105fe610d543660046144c3565b613264565b348015610d6557600080fd5b506105fe610d74366004614440565b60009182526007602052604090912055565b348015610d9257600080fd5b506105fe610da13660046140fb565b61330f565b348015610db257600080fd5b506104e8610dc13660046143df565b6000918252600e6020908152604080842061ffff93909316845291905290205490565b348015610df057600080fd5b506105fe610dff36600461458c565b601980547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610e4a57600080fd5b506105fe610e593660046140fb565b613395565b348015610e6a57600080fd5b506105fe610e79366004614518565b6019805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b348015610ec957600080fd5b506105c9610ed8366004614440565b61342d565b348015610ee957600080fd5b506104e8610ef8366004614440565b613496565b348015610f0957600080fd5b506104e8610f183660046140fb565b60096020526000908152604090205481565b348015610f3657600080fd5b506104e860185481565b348015610f4c57600080fd5b506105fe610f5b36600461402a565b6134c7565b348015610f6c57600080fd5b506104e8610f7b366004614440565b6134db565b348015610f8c57600080fd5b506104e8610f9b3660046140fb565b60066020526000908152604090205481565b348015610fb957600080fd5b506105c9610fc83660046143df565b6134f7565b348015610fd957600080fd5b506104e8610fe83660046140fb565b60026020526000908152604090205481565b60606000611008601361356b565b9050808410611043576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82611055576110528482614b3f565b92505b60008367ffffffffffffffff81111561107057611070614de2565b604051908082528060200260200182016040528015611099578160200160208202803683370190505b50905060005b848110156110eb576110bc6110b482886149c6565b601390613575565b8282815181106110ce576110ce614db3565b6020908102919091010152806110e381614cd4565b91505061109f565b509150505b92915050565b600e602052826000526040600020602052816000526040600020818154811061111e57600080fd5b9060005260206000200160009250925050505481565b6060806000611143601361356b565b905060008167ffffffffffffffff81111561116057611160614de2565b604051908082528060200260200182016040528015611189578160200160208202803683370190505b50905060008267ffffffffffffffff8111156111a7576111a7614de2565b6040519080825280602002602001820160405280156111d0578160200160208202803683370190505b50905060005b838110156114315760006111eb601383613575565b90508084838151811061120057611200614db3565b6020908102919091018101919091526000828152601290915260408082205490517f3ebe8d6c0000000000000000000000000000000000000000000000000000000081526004810184905261ffff90911691903090633ebe8d6c9060240160206040518083038186803b15801561127657600080fd5b505afa15801561128a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112ae91906142f5565b905060008167ffffffffffffffff8111156112cb576112cb614de2565b6040519080825280602002602001820160405280156112f4578160200160208202803683370190505b506000858152600e6020526040812091925090815b8561ffff168161ffff16116113ee5761ffff81166000908152602083815260408083208054825181850281018501909352808352919290919083018282801561137157602002820191906000526020600020905b81548152602001906001019080831161135d575b5050505050905060005b81518110156113d95781818151811061139657611396614db3565b60200260200101518686806113aa90614cd4565b9750815181106113bc576113bc614db3565b6020908102919091010152806113d181614cd4565b91505061137b565b505080806113e690614cb2565b915050611309565b506113fa838e86613581565b88888151811061140c5761140c614db3565b602002602001018181525050505050505050808061142990614cd4565b9150506111d6565b50909590945092505050565b6000828152600e6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156114a157602002820191906000526020600020905b81548152602001906001019080831161148d575b505050505090506114b38182516136e1565b92509250505b9250929050565b601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517f850af0cb00000000000000000000000000000000000000000000000000000000815290516000929163850af0cb9160048083019260a0929190829003018186803b15801561155557600080fd5b505afa158015611569573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061158d9190614167565b50601780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601554604080517f1b6b6d23000000000000000000000000000000000000000000000000000000008152905193975091169450631b6b6d2393506004808201935060209291829003018186803b15801561162c57600080fd5b505afa158015611640573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611664919061414a565b601680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601554604051911681527f6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014906020015b60405180910390a15050565b6000838152600d6020908152604080832080548251818502810185019093528083526117429383018282801561173657602002820191906000526020600020905b815481526020019060010190808311611722575b50505050508484613581565b90505b9392505050565b8060005b818160ff1610156117e0573063c8048022858560ff851681811061177657611776614db3565b905060200201356040518263ffffffff1660e01b815260040161179b91815260200190565b600060405180830381600087803b1580156117b557600080fd5b505af11580156117c9573d6000803e3d6000fd5b5050505080806117d890614d06565b915050611750565b507fbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b838360405161181292919061464b565b60405180910390a1505050565b60008181526012602052604081205461ffff1681805b8261ffff168161ffff1611611880576000858152600e6020908152604080832061ffff8516845290915290205461186c90836149c6565b91508061187881614cb2565b915050611835565b509392505050565b60005a9050600061189b8385018561435a565b50600081815260056020908152604080832054600490925282205492935091906118c3613766565b9050826118fd5760008481526005602090815260408083208490556010825282208054600181018255908352912042910155915081611b5d565b6000848152600360205260408120546119168484614b3f565b6119209190614b3f565b6000868152601160209081526040808320546010909252909120805492935061ffff909116918290811061195657611956614db3565b90600052602060002001547f000000000000000000000000000000000000000000000000000000000000000061ffff16426119919190614b3f565b1115611a0057600086815260106020908152604082208054600181018255908352912042910155806119c281614cb2565b600088815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559150505b600086815260126020908152604080832054600e835281842061ffff9091168085529083528184208054835181860281018601909452808452919493909190830182828015611a6e57602002820191906000526020600020905b815481526020019060010190808311611a5a575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff1681511415611aea5781611aac81614cb2565b60008a815260126020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600e6020908152604080832061ffff94851684528252808320805460018181018355918552838520018790558a8452600f83528184209590941683529381528382208054808501825590835281832001859055888252600d81529281208054928301815581529190912001555b600084815260066020526040812054611b779060016149c6565b6000868152600660209081526040918290208390558151878152908101859052908101859052606081018290529091507f6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede39060800160405180910390a16000858152600460209081526040808320859055601854600290925290912054611bfe9084614b3f565b1115611ea9576017546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810187905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a9060240160006040518083038186803b158015611c6f57600080fd5b505afa158015611c83573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611cc991908101906141d6565b6017546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810189905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c9060240160206040518083038186803b158015611d3957600080fd5b505afa158015611d4d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d7191906145a9565b601954909150611d959082906c01000000000000000000000000900460ff16614a97565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611ea6576019546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018990526bffffffffffffffffffffffff9091166024820152309063948108f790604401600060405180830381600087803b158015611e2657600080fd5b505af1158015611e3a573d6000803e3d6000fd5b50505060008881526002602090815260409182902087905560195482518b81526bffffffffffffffffffffffff909116918101919091529081018690527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0915060600160405180910390a15b50505b6000858152600760205260409020545b805a611ec59089614b3f565b611ed1906127106149c6565b1015611f1f5782406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905582611f1781614c16565b935050611eb9565b505050505050505050565b6000828152600e6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015611f8c57602002820191906000526020600020905b815481526020019060010190808311611f78575b5050505050905092915050565b8160005b818110156120365730635f17e616868684818110611fbd57611fbd614db3565b90506020020135856040518363ffffffff1660e01b8152600401611ff192919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561200b57600080fd5b505af115801561201f573d6000803e3d6000fd5b50505050808061202e90614cd4565b915050611f9d565b5050505050565b600b602052600090815260409020805461205690614c25565b80601f016020809104026020016040519081016040528092919081815260200182805461208290614c25565b80156120cf5780601f106120a4576101008083540402835291602001916120cf565b820191906000526020600020905b8154815290600101906020018083116120b257829003601f168201915b505050505081565b6120df613817565b6016546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561214957600080fd5b505afa15801561215d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061218191906142f5565b6016546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb90604401602060405180830381600087803b1580156121f457600080fd5b505af1158015612208573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061222c91906140e0565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600d909152812061226891613d97565b60008281526012602052604081205461ffff16905b8161ffff168161ffff16116122c4576000848152600e6020908152604080832061ffff8516845290915281206122b291613d97565b806122bc81614cb2565b91505061227d565b5050600082815260126020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055601190915281205461ffff16905b8161ffff168161ffff1611612352576000848152600f6020908152604080832061ffff85168452909152812061234091613d97565b8061234a81614cb2565b91505061230b565b50600083815260106020526040812061236a91613d97565b5050600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6040517f3ebe8d6c00000000000000000000000000000000000000000000000000000000815260048101839052600090819081903090633ebe8d6c9060240160206040518083038186803b1580156123f957600080fd5b505afa15801561240d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061243191906142f5565b90508315806124405750808410155b15612449578093505b60008581526012602052604081205485919061ffff16805b6000898152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835291929091908301828280156124c357602002820191906000526020600020905b8154815260200190600101908083116124af575b505050505090506000806124d783886136e1565b90925090506124e682876149c6565b95506124f28188614b3f565b9650600087116125045750505061251a565b505050808061251290614bda565b915050612461565b5090979596505050505050565b6000606060005a9050600061253e858701876140fb565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff81111561257757612577614de2565b6040519080825280601f01601f1916602001820160405280156125a1576020820181803683370190505b506040516020016125b392919061487a565b604051602081830303815290604052905060006125ce613766565b905060006125db8661267c565b90505b835a6125ea9089614b3f565b6125f6906127106149c6565b10156126445781406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558161263c81614c16565b9250506125de565b9a91995090975050505050505050565b600f602052826000526040600020602052816000526040600020818154811061111e57600080fd5b60008181526005602052604081205461269757506001919050565b6000828152600360209081526040808320546004909252909120546126ba613766565b6126c49190614b3f565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612752576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040517fa5f589340000000000000000000000000000000000000000000000000000000081526004810183905260009081908190309063a5f589349060240160206040518083038186803b15801561282557600080fd5b505afa158015612839573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061285d91906142f5565b905083158061286c5750808410155b15612875578093505b60008581526011602052604081205485919061ffff16805b6000898152600f6020908152604080832061ffff851684528252808320805482518185028101850190935280835291929091908301828280156128ef57602002820191906000526020600020905b8154815260200190600101908083116128db575b5050505050905060008061290383886136e1565b909250905061291282876149c6565b955061291e8188614b3f565b9650600087116129305750505061251a565b505050808061293e90614bda565b91505061288d565b6000838152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938301828280156129a557602002820191906000526020600020905b815481526020019060010190808311612991575b505050505090506129b881858351613581565b95945050505050565b6000838152600f6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938301828280156129a5576020028201919060005260206000209081548152602001906001019080831161299157505050505090506129b881858351613581565b6060806000612a40601361356b565b905060008167ffffffffffffffff811115612a5d57612a5d614de2565b604051908082528060200260200182016040528015612a86578160200160208202803683370190505b50905060008267ffffffffffffffff811115612aa457612aa4614de2565b604051908082528060200260200182016040528015612acd578160200160208202803683370190505b50905060005b83811015611431576000612ae8601383613575565b6000818152600d6020908152604080832080548251818502810185019093528083529495509293909291830182828015612b4157602002820191906000526020600020905b815481526020019060010190808311612b2d575b5050505050905081858481518110612b5b57612b5b614db3565b602002602001018181525050612b73818a8351613581565b848481518110612b8557612b85614db3565b60200260200101818152505050508080612b9e90614cd4565b915050612ad3565b6040805161014081018252600461010082019081527f746573740000000000000000000000000000000000000000000000000000000061012083015281528151602081810184526000808352818401929092523083850181905263ffffffff8916606085015260808401528351808201855282815260a08401528351908101909352825260c08101919091526bffffffffffffffffffffffff841660e082015260165460155473ffffffffffffffffffffffffffffffffffffffff9182169163095ea7b39116612c7960ff8a1688614a97565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff166024820152604401602060405180830381600087803b158015612cf257600080fd5b505af1158015612d06573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d2a91906140e0565b5060008660ff1667ffffffffffffffff811115612d4957612d49614de2565b604051908082528060200260200182016040528015612d72578160200160208202803683370190505b50905060005b8760ff168160ff161015612de5576000612d918461389a565b905080838360ff1681518110612da957612da9614db3565b60209081029190910181019190915260009182526008815260408083208890556007909152902084905580612ddd81614d06565b915050612d78565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c71181604051612e1591906146a0565b60405180910390a150505050505050565b6016546017546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b390604401602060405180830381600087803b158015612ea957600080fd5b505af1158015612ebd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ee191906140e0565b506017546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b158015612f6257600080fd5b505af1158015612f76573d6000803e3d6000fd5b5050604080518581526bffffffffffffffffffffffff851660208201527f8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e93500190506116d5565b6000828152600f6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015611f8c5760200282019190600052602060002090815481526020019060010190808311611f78575050505050905092915050565b8060005b818110156130f657600084848381811061304b5761304b614db3565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16639fab4386828360405160200161308491815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016130b092919061487a565b600060405180830381600087803b1580156130ca57600080fd5b505af11580156130de573d6000803e3d6000fd5b505050505080806130ee90614cd4565b91505061302f565b50505050565b6000818152600d602090815260409182902080548351818402810184019094528084526060939283018282801561315257602002820191906000526020600020905b81548152602001906001019080831161313e575b50505050509050919050565b6017546040517f9fab438600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690639fab4386906131b890869086908690600401614826565b600060405180830381600087803b1580156131d257600080fd5b505af11580156131e6573d6000803e3d6000fd5b5050506000848152600b602052604090206130f691508383613db5565b60008181526011602052604081205461ffff1681805b8261ffff168161ffff1611611880576000858152600f6020908152604080832061ffff8516845290915290205461325090836149c6565b91508061325c81614cb2565b915050613219565b6017546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156132dc57600080fd5b505af11580156132f0573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6017546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561338157600080fd5b505af1158015612036573d6000803e3d6000fd5b6017546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b15801561340157600080fd5b505af1158015613415573d6000803e3d6000fd5b5050505061222c81601361399c90919063ffffffff16565b6000828152600d6020908152604080832080548251818502810185019093528083528493849392919083018282801561348557602002820191906000526020600020905b815481526020019060010190808311613471575b505050505090506114b381856136e1565b600d60205281600052604060002081815481106134b257600080fd5b90600052602060002001600091509150505481565b6134cf613817565b6134d8816139a8565b50565b601060205281600052604060002081815481106134b257600080fd5b6000828152600f6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156114a1576020028201919060005260206000209081548152602001906001019080831161148d57505050505090506114b38182516136e1565b60006110f0825490565b60006117458383613a9e565b825160009081908315806135955750808410155b1561359e578093505b60008467ffffffffffffffff8111156135b9576135b9614de2565b6040519080825280602002602001820160405280156135e2578160200160208202803683370190505b509050600092505b84831015613650578660016135ff8585614b3f565b6136099190614b3f565b8151811061361957613619614db3565b602002602001015181848151811061363357613633614db3565b60209081029190910101528261364881614cd4565b9350506135ea565b613669816000600184516136649190614b3f565b613ac8565b85606414156136a35780600182516136819190614b3f565b8151811061369157613691614db3565b60200260200101519350505050611745565b8060648251886136b39190614a5a565b6136bd9190614a46565b815181106136cd576136cd614db3565b602002602001015193505050509392505050565b8151600090819081908415806136f75750808510155b15613700578094505b60008092505b8583101561375c5786600161371b8585614b3f565b6137259190614b3f565b8151811061373557613735614db3565b60200260200101518161374891906149c6565b90508261375481614cd4565b935050613706565b9694955050505050565b60007f00000000000000000000000000000000000000000000000000000000000000001561381257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156137d557600080fd5b505afa1580156137e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061380d91906142f5565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613898576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401612749565b565b6015546040517f08b79da4000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff909116906308b79da4906138f5908690600401614706565b602060405180830381600087803b15801561390f57600080fd5b505af1158015613923573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061394791906142f5565b9050613954601382613c49565b5060608301516000828152600a6020908152604080832063ffffffff90941690935560a0860151600b82529290912082516139959391929190910190613e57565b5092915050565b60006117458383613c55565b73ffffffffffffffffffffffffffffffffffffffff8116331415613a28576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401612749565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613ab557613ab5614db3565b9060005260206000200154905092915050565b818180821415613ad9575050505050565b6000856002613ae88787614acb565b613af291906149de565b613afc9087614952565b81518110613b0c57613b0c614db3565b602002602001015190505b818313613c1b575b80868481518110613b3257613b32614db3565b60200260200101511015613b525782613b4a81614c79565b935050613b1f565b858281518110613b6457613b64614db3565b6020026020010151811015613b855781613b7d81614b82565b925050613b52565b818313613c1657858281518110613b9e57613b9e614db3565b6020026020010151868481518110613bb857613bb8614db3565b6020026020010151878581518110613bd257613bd2614db3565b60200260200101888581518110613beb57613beb614db3565b60209081029190910101919091525282613c0481614c79565b9350508180613c1290614b82565b9250505b613b17565b81851215613c2e57613c2e868684613ac8565b83831215613c4157613c41868486613ac8565b505050505050565b60006117458383613d48565b60008181526001830160205260408120548015613d3e576000613c79600183614b3f565b8554909150600090613c8d90600190614b3f565b9050818114613cf2576000866000018281548110613cad57613cad614db3565b9060005260206000200154905080876000018481548110613cd057613cd0614db3565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613d0357613d03614d84565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506110f0565b60009150506110f0565b6000818152600183016020526040812054613d8f575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556110f0565b5060006110f0565b50805460008255906000526020600020908101906134d89190613ecb565b828054613dc190614c25565b90600052602060002090601f016020900481019282613de35760008555613e47565b82601f10613e1a578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555613e47565b82800160010185558215613e47579182015b82811115613e47578235825591602001919060010190613e2c565b50613e53929150613ecb565b5090565b828054613e6390614c25565b90600052602060002090601f016020900481019282613e855760008555613e47565b82601f10613e9e57805160ff1916838001178555613e47565b82800160010185558215613e47579182015b82811115613e47578251825591602001919060010190613eb0565b5b80821115613e535760008155600101613ecc565b8051613eeb81614e11565b919050565b60008083601f840112613f0257600080fd5b50813567ffffffffffffffff811115613f1a57600080fd5b6020830191508360208260051b85010111156114b957600080fd5b80518015158114613eeb57600080fd5b60008083601f840112613f5757600080fd5b50813567ffffffffffffffff811115613f6f57600080fd5b6020830191508360208285010111156114b957600080fd5b600082601f830112613f9857600080fd5b8151613fab613fa68261490c565b6148bd565b818152846020838601011115613fc057600080fd5b613fd1826020830160208701614b56565b949350505050565b803561ffff81168114613eeb57600080fd5b8051613eeb81614e33565b805167ffffffffffffffff81168114613eeb57600080fd5b803560ff81168114613eeb57600080fd5b8051613eeb81614e45565b60006020828403121561403c57600080fd5b813561174581614e11565b6000806020838503121561405a57600080fd5b823567ffffffffffffffff81111561407157600080fd5b61407d85828601613ef0565b90969095509350505050565b60008060006040848603121561409e57600080fd5b833567ffffffffffffffff8111156140b557600080fd5b6140c186828701613ef0565b90945092505060208401356140d581614e33565b809150509250925092565b6000602082840312156140f257600080fd5b61174582613f35565b60006020828403121561410d57600080fd5b5035919050565b6000806020838503121561412757600080fd5b823567ffffffffffffffff81111561413e57600080fd5b61407d85828601613f45565b60006020828403121561415c57600080fd5b815161174581614e11565b600080600080600060a0868803121561417f57600080fd5b85516003811061418e57600080fd5b602087015190955061419f81614e33565b60408701519094506141b081614e33565b60608701519093506141c181614e11565b80925050608086015190509295509295909350565b6000602082840312156141e857600080fd5b815167ffffffffffffffff8082111561420057600080fd5b90830190610140828603121561421557600080fd5b61421d614893565b61422683613ee0565b815261423460208401613feb565b602082015260408301518281111561424b57600080fd5b61425787828601613f87565b6040830152506142696060840161401f565b606082015261427a60808401613ee0565b608082015261428b60a08401613ff6565b60a082015261429c60c08401613feb565b60c08201526142ad60e0840161401f565b60e08201526101006142c0818501613f35565b9082015261012083810151838111156142d857600080fd5b6142e488828701613f87565b918301919091525095945050505050565b60006020828403121561430757600080fd5b5051919050565b60008060006040848603121561432357600080fd5b83359250602084013567ffffffffffffffff81111561434157600080fd5b61434d86828701613f45565b9497909650939450505050565b6000806040838503121561436d57600080fd5b82359150602083013567ffffffffffffffff81111561438b57600080fd5b8301601f8101851361439c57600080fd5b80356143aa613fa68261490c565b8181528660208385010111156143bf57600080fd5b816020840160208301376000602083830101528093505050509250929050565b600080604083850312156143f257600080fd5b8235915061440260208401613fd9565b90509250929050565b60008060006060848603121561442057600080fd5b8335925061443060208501613fd9565b9150604084013590509250925092565b6000806040838503121561445357600080fd5b50508035926020909101359150565b60008060006060848603121561447757600080fd5b833592506020840135915061448e60408501613fd9565b90509250925092565b6000806000606084860312156144ac57600080fd5b505081359360208301359350604090920135919050565b600080604083850312156144d657600080fd5b8235915060208301356144e881614e33565b809150509250929050565b6000806040838503121561450657600080fd5b8235915060208301356144e881614e45565b60006020828403121561452a57600080fd5b6117458261400e565b600080600080600060a0868803121561454b57600080fd5b6145548661400e565b9450602086013561456481614e33565b9350604086013561457481614e45565b94979396509394606081013594506080013592915050565b60006020828403121561459e57600080fd5b813561174581614e45565b6000602082840312156145bb57600080fd5b815161174581614e45565b600081518084526020808501945080840160005b838110156145f6578151875295820195908201906001016145da565b509495945050505050565b60008151808452614619816020860160208601614b56565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83111561468457600080fd5b8260051b80856040850137600092016040019182525092915050565b60208152600061174560208301846145c6565b6040815260006146c660408301856145c6565b82810360208401526129b881856145c6565b82151581526040602082015260006117426040830184614601565b6020815260006117456020830184614601565b6020815260008251610100806020850152614725610120850183614601565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808685030160408701526147618483614601565b93506040870151915061478c606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a08701519150808685030160c08701526147dd8483614601565b935060c08701519150808685030160e0870152506147fb8382614601565b92505060e085015161481c828601826bffffffffffffffffffffffff169052565b5090949350505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b8281526040602082015260006117426040830184614601565b604051610140810167ffffffffffffffff811182821017156148b7576148b7614de2565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561490457614904614de2565b604052919050565b600067ffffffffffffffff82111561492657614926614de2565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000808212827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0384138115161561498c5761498c614d26565b827f80000000000000000000000000000000000000000000000000000000000000000384128116156149c0576149c0614d26565b50500190565b600082198211156149d9576149d9614d26565b500190565b6000826149ed576149ed614d55565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615614a4157614a41614d26565b500590565b600082614a5557614a55614d55565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614a9257614a92614d26565b500290565b60006bffffffffffffffffffffffff80831681851681830481118215151615614ac257614ac2614d26565b02949350505050565b6000808312837f800000000000000000000000000000000000000000000000000000000000000001831281151615614b0557614b05614d26565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff018313811615614b3957614b39614d26565b50500390565b600082821015614b5157614b51614d26565b500390565b60005b83811015614b71578181015183820152602001614b59565b838111156130f65750506000910152565b60007f8000000000000000000000000000000000000000000000000000000000000000821415614bb457614bb4614d26565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600061ffff821680614bee57614bee614d26565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b600081614bb457614bb4614d26565b600181811c90821680614c3957607f821691505b60208210811415614c73577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415614cab57614cab614d26565b5060010190565b600061ffff80831681811415614cca57614cca614d26565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415614cab57614cab614d26565b600060ff821660ff811415614d1d57614d1d614d26565b60010192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff811681146134d857600080fd5b63ffffffff811681146134d857600080fd5b6bffffffffffffffffffffffff811681146134d857600080fdfea164736f6c6343000806000a",
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchRegisterUpkeeps", number, gasLimit, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, number, gasLimit, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, number, gasLimit, amount, checkGasToBurn, performGasToBurn)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchUpdateCheckData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchUpdateCheckData", upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchUpdateCheckData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchUpdateCheckData(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchUpdateCheckData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchUpdateCheckData(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) UpdateCheckData(opts *bind.TransactOpts, upkeepId *big.Int, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "updateCheckData", upkeepId, checkData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) UpdateCheckData(upkeepId *big.Int, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateCheckData(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, checkData)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) UpdateCheckData(upkeepId *big.Int, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateCheckData(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, checkData)
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

	FilterFundsAdded(opts *bind.FilterOpts) (*VerifiableLoadUpkeepFundsAddedIterator, error)

	WatchFundsAdded(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepFundsAdded) (event.Subscription, error)

	ParseFundsAdded(log types.Log) (*VerifiableLoadUpkeepFundsAdded, error)

	FilterInsufficientFunds(opts *bind.FilterOpts) (*VerifiableLoadUpkeepInsufficientFundsIterator, error)

	WatchInsufficientFunds(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepInsufficientFunds) (event.Subscription, error)

	ParseInsufficientFunds(log types.Log) (*VerifiableLoadUpkeepInsufficientFunds, error)

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
