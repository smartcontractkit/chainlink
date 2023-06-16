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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"registrarAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"FundsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"InsufficientFunds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"firstPerformBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lastBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"PerformingUpkeep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"RegistrarSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"UpkeepsRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_INTERVAL\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"cancelUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkDatas\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getDelaysLengthAtTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxBucketedDelaysForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"getPxDelayForAllUpkeeps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getPxDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInTimestampBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumTimestampBucketedDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTimestampBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"timestampBucket\",\"type\":\"uint16\"}],\"name\":\"getTimestampDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractKeeperRegistry2_0\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"setAddLinkAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setCheckGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractKeeperRegistrar2_0\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"newMinBalanceThresholdMultiplier\",\"type\":\"uint8\"}],\"name\":\"setMinBalanceThresholdMultiplier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformGasToBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newInterval\",\"type\":\"uint256\"}],\"name\":\"setUpkeepTopUpCheckInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampBuckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestampDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"timestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"updateCheckData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040526005601855601980546001600160681b0319166c140000000002c68af0bb140000179055601960f21b60a05260e160f41b60c0523480156200004557600080fd5b50604051620053bc380380620053bc833981016040819052620000689162000320565b81813380600081620000c15760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000f457620000f4816200025a565b5050601580546001600160a01b0319166001600160a01b0385169081179091556040805163850af0cb60e01b815290516000935063850af0cb9160048082019260a092909190829003018186803b1580156200014f57600080fd5b505afa15801562000164573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200018a91906200038a565b50601780546001600160a01b0319166001600160a01b038381169190911790915560155460408051631b6b6d2360e01b8152905193975091169450631b6b6d2393506004808201935060209291829003018186803b158015620001ec57600080fd5b505afa15801562000201573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000227919062000363565b601680546001600160a01b0319166001600160a01b039290921691909117905550151560f81b6080525062000413915050565b6001600160a01b038116331415620002b55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000b8565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b805163ffffffff811681146200031b57600080fd5b919050565b600080604083850312156200033457600080fd5b82516200034181620003fa565b602084015190925080151581146200035857600080fd5b809150509250929050565b6000602082840312156200037657600080fd5b81516200038381620003fa565b9392505050565b600080600080600060a08688031215620003a357600080fd5b855160038110620003b357600080fd5b9450620003c36020870162000306565b9350620003d36040870162000306565b92506060860151620003e581620003fa565b80925050608086015190509295509295909350565b6001600160a01b03811681146200041057600080fd5b50565b60805160f81c60a05160f01c60c05160f01c614f5b620004616000396000818161064d015261199e0152600081816105230152611ab20152600081816108fc015261383f0152614f5b6000f3fe6080604052600436106104695760003560e01c80637b10399911610243578063a6c60d8911610143578063d6051a72116100bb578063f2fde38b1161008a578063fba7ffa31161006f578063fba7ffa314610fbb578063fbfb4f7614610fe8578063fcdc1f631461100857600080fd5b8063f2fde38b14610f7b578063fb0ceb0414610f9b57600080fd5b8063d6051a7214610ef8578063dbef701e14610f18578063e0114adb14610f38578063e455308314610f6557600080fd5b8063b0971e1a11610112578063c357f1f3116100f7578063c357f1f314610e1f578063c804802214610e79578063d355852814610e9957600080fd5b8063b0971e1a14610dc1578063becde0e114610dff57600080fd5b8063a6c60d8914610d34578063a72aa27e14610d54578063a79c404314610d74578063af953a4a14610da157600080fd5b8063948108f7116101d65780639b51fb0d116101a55780639d385eaa1161018a5780639d385eaa14610cd45780639fab438614610cf4578063a5f5893414610d1457600080fd5b80639b51fb0d14610c835780639bb8651114610cb457600080fd5b8063948108f714610bda57806399cc6b0b14610bfa5780639ac542eb14610c1a5780639b42935414610c5657600080fd5b80638bc7b772116102125780638bc7b77214610b425780638da5cb5b14610b625780638fcb3fba14610b8d5780639095aa3514610bba57600080fd5b80637b10399914610ab55780637e4087b814610ae25780638237831714610b0257806387dfa90014610b2257600080fd5b806351c98be31161036957806369cdbadb116102e15780637145f11b116102b057806376721303116102955780637672130314610a53578063776898c814610a8057806379ba509714610aa057600080fd5b80637145f11b146109f657806373644cce14610a2657600080fd5b806369cdbadb1461094e57806369e9b7731461097b5780636e04ff0d146109a85780637137a702146109d657600080fd5b80635f17e61611610338578063636092e81161031d578063636092e8146108a8578063642f6cef146108ea578063643b34e91461092e57600080fd5b80635f17e6161461085b57806360457ff51461087b57600080fd5b806351c98be3146107cc57806357970e93146107ec57806358c52c04146108195780635d4ee7f31461084657600080fd5b806329f0e496116103fc57806333774d1c116103cb5780634585e33b116103b05780634585e33b1461075f57806345d2ec171461077f57806346e7a63e1461079f57600080fd5b806333774d1c1461070e5780633ebe8d6c1461073f57600080fd5b806329f0e4961461063b5780632a9032d31461066f5780632b20e3971461068f578063328ffd11146106e157600080fd5b80631bee0080116104385780631bee008014610596578063206c32e8146105c457806320e3dbd4146105f957806328c4b57b1461061b57600080fd5b806306e3b632146104ad578063077ac621146104e357806312c5502714610511578063177b0eb91461055857600080fd5b366104a857604080513381523460208201527f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874910160405180910390a1005b600080fd5b3480156104b957600080fd5b506104cd6104c8366004614515565b611035565b6040516104da9190614775565b60405180910390f35b3480156104ef57600080fd5b506105036104fe3660046144e0565b611131565b6040519081526020016104da565b34801561051d57600080fd5b506105457f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff90911681526020016104da565b34801561056457600080fd5b506105036105733660046144b4565b6000918252600f6020908152604080842061ffff93909316845291905290205490565b3480156105a257600080fd5b506105b66105b13660046141d0565b61116f565b6040516104da929190614788565b3480156105d057600080fd5b506105e46105df3660046144b4565b611478565b604080519283526020830191909152016104da565b34801561060557600080fd5b506106196106143660046140ff565b6114fb565b005b34801561062757600080fd5b5061050361063636600461456c565b61171c565b34801561064757600080fd5b506105457f000000000000000000000000000000000000000000000000000000000000000081565b34801561067b57600080fd5b5061061961068a36600461411c565b611787565b34801561069b57600080fd5b506015546106bc9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016104da565b3480156106ed57600080fd5b506105036106fc3660046141d0565b60036020526000908152604090205481565b34801561071a57600080fd5b506105456107293660046141d0565b60116020526000908152604090205461ffff1681565b34801561074b57600080fd5b5061050361075a3660046141d0565b61185a565b34801561076b57600080fd5b5061061961077a3660046141e9565b6118c3565b34801561078b57600080fd5b506104cd61079a3660046144b4565b611f65565b3480156107ab57600080fd5b506105036107ba3660046141d0565b600a6020526000908152604090205481565b3480156107d857600080fd5b506106196107e736600461415e565b611fd4565b3480156107f857600080fd5b506016546106bc9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561082557600080fd5b506108396108343660046141d0565b612078565b6040516104da91906147c8565b34801561085257600080fd5b50610619612112565b34801561086757600080fd5b50610619610876366004614515565b61226b565b34801561088757600080fd5b506105036108963660046141d0565b60076020526000908152604090205481565b3480156108b457600080fd5b506019546108cd906bffffffffffffffffffffffff1681565b6040516bffffffffffffffffffffffff90911681526020016104da565b3480156108f657600080fd5b5061091e7f000000000000000000000000000000000000000000000000000000000000000081565b60405190151581526020016104da565b34801561093a57600080fd5b506105e4610949366004614515565b6123dd565b34801561095a57600080fd5b506105036109693660046141d0565b60086020526000908152604090205481565b34801561098757600080fd5b50610619610996366004614515565b60009182526008602052604090912055565b3480156109b457600080fd5b506109c86109c33660046141e9565b612562565b6040516104da9291906147ad565b3480156109e257600080fd5b506105036109f13660046144e0565b61268f565b348015610a0257600080fd5b5061091e610a113660046141d0565b600c6020526000908152604090205460ff1681565b348015610a3257600080fd5b50610503610a413660046141d0565b6000908152600d602052604090205490565b348015610a5f57600080fd5b50610503610a6e3660046141d0565b60046020526000908152604090205481565b348015610a8c57600080fd5b5061091e610a9b3660046141d0565b6126b7565b348015610aac57600080fd5b50610619612707565b348015610ac157600080fd5b506017546106bc9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610aee57600080fd5b506105e4610afd366004614515565b612809565b348015610b0e57600080fd5b50610503610b1d366004614537565b612981565b348015610b2e57600080fd5b50610503610b3d366004614537565b6129fc565b348015610b4e57600080fd5b506105b6610b5d3660046141d0565b612a6c565b348015610b6e57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166106bc565b348015610b9957600080fd5b50610503610ba83660046141d0565b60056020526000908152604090205481565b348015610bc657600080fd5b50610619610bd5366004614608565b612be1565b348015610be657600080fd5b50610619610bf53660046145c8565b612e61565b348015610c0657600080fd5b506104cd610c153660046144b4565b612ff9565b348015610c2657600080fd5b50601954610c44906c01000000000000000000000000900460ff1681565b60405160ff90911681526020016104da565b348015610c6257600080fd5b50610619610c71366004614515565b60009182526009602052604090912055565b348015610c8f57600080fd5b50610545610c9e3660046141d0565b60126020526000908152604090205461ffff1681565b348015610cc057600080fd5b50610619610ccf36600461411c565b613066565b348015610ce057600080fd5b506104cd610cef3660046141d0565b613137565b348015610d0057600080fd5b50610619610d0f3660046143e3565b613199565b348015610d2057600080fd5b50610503610d2f3660046141d0565b61323e565b348015610d4057600080fd5b50610619610d4f3660046141d0565b601855565b348015610d6057600080fd5b50610619610d6f366004614598565b61329f565b348015610d8057600080fd5b50610619610d8f366004614515565b60009182526007602052604090912055565b348015610dad57600080fd5b50610619610dbc3660046141d0565b61334a565b348015610dcd57600080fd5b50610503610ddc3660046144b4565b6000918252600e6020908152604080842061ffff93909316845291905290205490565b348015610e0b57600080fd5b50610619610e1a36600461411c565b6133d0565b348015610e2b57600080fd5b50610619610e3a366004614661565b601980547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff92909216919091179055565b348015610e8557600080fd5b50610619610e943660046141d0565b61346a565b348015610ea557600080fd5b50610619610eb43660046145ed565b6019805460ff9092166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffff909216919091179055565b348015610f0457600080fd5b506105e4610f13366004614515565b613502565b348015610f2457600080fd5b50610503610f33366004614515565b61356b565b348015610f4457600080fd5b50610503610f533660046141d0565b60096020526000908152604090205481565b348015610f7157600080fd5b5061050360185481565b348015610f8757600080fd5b50610619610f963660046140ff565b61359c565b348015610fa757600080fd5b50610503610fb6366004614515565b6135b0565b348015610fc757600080fd5b50610503610fd63660046141d0565b60066020526000908152604090205481565b348015610ff457600080fd5b506105e46110033660046144b4565b6135cc565b34801561101457600080fd5b506105036110233660046141d0565b60026020526000908152604090205481565b606060006110436013613640565b905080841061107e576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826110905761108d8482614c14565b92505b60008367ffffffffffffffff8111156110ab576110ab614ed1565b6040519080825280602002602001820160405280156110d4578160200160208202803683370190505b50905060005b84811015611126576110f76110ef8288614a9b565b60139061364a565b82828151811061110957611109614ea2565b60209081029190910101528061111e81614da9565b9150506110da565b509150505b92915050565b600e602052826000526040600020602052816000526040600020818154811061115957600080fd5b9060005260206000200160009250925050505481565b606080600061117e6013613640565b905060008167ffffffffffffffff81111561119b5761119b614ed1565b6040519080825280602002602001820160405280156111c4578160200160208202803683370190505b50905060008267ffffffffffffffff8111156111e2576111e2614ed1565b60405190808252806020026020018201604052801561120b578160200160208202803683370190505b50905060005b8381101561146c57600061122660138361364a565b90508084838151811061123b5761123b614ea2565b6020908102919091018101919091526000828152601290915260408082205490517f3ebe8d6c0000000000000000000000000000000000000000000000000000000081526004810184905261ffff90911691903090633ebe8d6c9060240160206040518083038186803b1580156112b157600080fd5b505afa1580156112c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112e991906143ca565b905060008167ffffffffffffffff81111561130657611306614ed1565b60405190808252806020026020018201604052801561132f578160200160208202803683370190505b506000858152600e6020526040812091925090815b8561ffff168161ffff16116114295761ffff8116600090815260208381526040808320805482518185028101850190935280835291929091908301828280156113ac57602002820191906000526020600020905b815481526020019060010190808311611398575b5050505050905060005b8151811015611414578181815181106113d1576113d1614ea2565b60200260200101518686806113e590614da9565b9750815181106113f7576113f7614ea2565b60209081029190910101528061140c81614da9565b9150506113b6565b5050808061142190614d87565b915050611344565b50611435838e86613656565b88888151811061144757611447614ea2565b602002602001018181525050505050505050808061146490614da9565b915050611211565b50909590945092505050565b6000828152600e6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156114dc57602002820191906000526020600020905b8154815260200190600101908083116114c8575b505050505090506114ee8182516137b6565b92509250505b9250929050565b601580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517f850af0cb00000000000000000000000000000000000000000000000000000000815290516000929163850af0cb9160048083019260a0929190829003018186803b15801561159057600080fd5b505afa1580156115a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115c8919061423c565b50601780547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601554604080517f1b6b6d23000000000000000000000000000000000000000000000000000000008152905193975091169450631b6b6d2393506004808201935060209291829003018186803b15801561166757600080fd5b505afa15801561167b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061169f919061421f565b601680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff928316179055601554604051911681527f6263309d5d4d1cfececd45a387cda7f14dccde21cf7a1bee1be6561075e61014906020015b60405180910390a15050565b6000838152600d60209081526040808320805482518185028101850190935280835261177d9383018282801561177157602002820191906000526020600020905b81548152602001906001019080831161175d575b50505050508484613656565b90505b9392505050565b8060005b818160ff16101561181b573063c8048022858560ff85168181106117b1576117b1614ea2565b905060200201356040518263ffffffff1660e01b81526004016117d691815260200190565b600060405180830381600087803b1580156117f057600080fd5b505af1158015611804573d6000803e3d6000fd5b50505050808061181390614df5565b91505061178b565b507fbeac20a03a6674e40498fac4356bc86e356c0d761a8d35d436712dc93bc7c74b838360405161184d929190614720565b60405180910390a1505050565b60008181526012602052604081205461ffff1681805b8261ffff168161ffff16116118bb576000858152600e6020908152604080832061ffff851684529091529020546118a79083614a9b565b9150806118b381614d87565b915050611870565b509392505050565b60005a905060006118d68385018561442f565b50600081815260056020908152604080832054600490925282205492935091906118fe61383b565b9050826119385760008481526005602090815260408083208490556010825282208054600181018255908352912042910155915081611b98565b6000848152600360205260408120546119518484614c14565b61195b9190614c14565b6000868152601160209081526040808320546010909252909120805492935061ffff909116918290811061199157611991614ea2565b90600052602060002001547f000000000000000000000000000000000000000000000000000000000000000061ffff16426119cc9190614c14565b1115611a3b57600086815260106020908152604082208054600181018255908352912042910155806119fd81614d87565b600088815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559150505b600086815260126020908152604080832054600e835281842061ffff9091168085529083528184208054835181860281018601909452808452919493909190830182828015611aa957602002820191906000526020600020905b815481526020019060010190808311611a95575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff1681511415611b255781611ae781614d87565b60008a815260126020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000878152600e6020908152604080832061ffff94851684528252808320805460018181018355918552838520018790558a8452600f83528184209590941683529381528382208054808501825590835281832001859055888252600d81529281208054928301815581529190912001555b600084815260066020526040812054611bb2906001614a9b565b6000868152600660209081526040918290208390558151878152908101859052908101859052606081018290529091507f6b6b3eeaaf107627513e76a81662118e7b1d8c78866f70760262115ddcfeede39060800160405180910390a16000858152600460209081526040808320859055601854600290925290912054611c399084614c14565b1115611ee4576017546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810187905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a9060240160006040518083038186803b158015611caa57600080fd5b505afa158015611cbe573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611d0491908101906142ab565b6017546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810189905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c9060240160206040518083038186803b158015611d7457600080fd5b505afa158015611d88573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611dac919061467e565b601954909150611dd09082906c01000000000000000000000000900460ff16614b6c565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611ee1576019546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018990526bffffffffffffffffffffffff9091166024820152309063948108f790604401600060405180830381600087803b158015611e6157600080fd5b505af1158015611e75573d6000803e3d6000fd5b50505060008881526002602090815260409182902087905560195482518b81526bffffffffffffffffffffffff909116918101919091529081018690527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0915060600160405180910390a15b50505b6000858152600760205260409020545b805a611f009089614c14565b611f0c90612710614a9b565b1015611f5a5782406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905582611f5281614ceb565b935050611ef4565b505050505050505050565b6000828152600e6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015611fc757602002820191906000526020600020905b815481526020019060010190808311611fb3575b5050505050905092915050565b8160005b818110156120715730635f17e616868684818110611ff857611ff8614ea2565b90506020020135856040518363ffffffff1660e01b815260040161202c92919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561204657600080fd5b505af115801561205a573d6000803e3d6000fd5b50505050808061206990614da9565b915050611fd8565b5050505050565b600b602052600090815260409020805461209190614cfa565b80601f01602080910402602001604051908101604052809291908181526020018280546120bd90614cfa565b801561210a5780601f106120df5761010080835404028352916020019161210a565b820191906000526020600020905b8154815290600101906020018083116120ed57829003601f168201915b505050505081565b61211a6138ec565b6016546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561218457600080fd5b505afa158015612198573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121bc91906143ca565b6016546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb90604401602060405180830381600087803b15801561222f57600080fd5b505af1158015612243573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061226791906141b5565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600d90915281206122a391613e6c565b60008281526012602052604081205461ffff16905b8161ffff168161ffff16116122ff576000848152600e6020908152604080832061ffff8516845290915281206122ed91613e6c565b806122f781614d87565b9150506122b8565b5050600082815260126020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055601190915281205461ffff16905b8161ffff168161ffff161161238d576000848152600f6020908152604080832061ffff85168452909152812061237b91613e6c565b8061238581614d87565b915050612346565b5060008381526010602052604081206123a591613e6c565b5050600090815260116020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6040517f3ebe8d6c00000000000000000000000000000000000000000000000000000000815260048101839052600090819081903090633ebe8d6c9060240160206040518083038186803b15801561243457600080fd5b505afa158015612448573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061246c91906143ca565b905083158061247b5750808410155b15612484578093505b60008581526012602052604081205485919061ffff16805b6000898152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835291929091908301828280156124fe57602002820191906000526020600020905b8154815260200190600101908083116124ea575b5050505050905060008061251283886137b6565b90925090506125218287614a9b565b955061252d8188614c14565b96506000871161253f57505050612555565b505050808061254d90614caf565b91505061249c565b5090979596505050505050565b6000606060005a90506000612579858701876141d0565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff8111156125b2576125b2614ed1565b6040519080825280601f01601f1916602001820160405280156125dc576020820181803683370190505b506040516020016125ee92919061494f565b6040516020818303038152906040529050600061260961383b565b90506000612616866126b7565b90505b835a6126259089614c14565b61263190612710614a9b565b101561267f5781406000908152600c6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558161267781614ceb565b925050612619565b9a91995090975050505050505050565b600f602052826000526040600020602052816000526040600020818154811061115957600080fd5b6000818152600560205260408120546126d257506001919050565b6000828152600360209081526040808320546004909252909120546126f561383b565b6126ff9190614c14565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461278d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6040517fa5f589340000000000000000000000000000000000000000000000000000000081526004810183905260009081908190309063a5f589349060240160206040518083038186803b15801561286057600080fd5b505afa158015612874573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061289891906143ca565b90508315806128a75750808410155b156128b0578093505b60008581526011602052604081205485919061ffff16805b6000898152600f6020908152604080832061ffff8516845282528083208054825181850281018501909352808352919290919083018282801561292a57602002820191906000526020600020905b815481526020019060010190808311612916575b5050505050905060008061293e83886137b6565b909250905061294d8287614a9b565b95506129598188614c14565b96506000871161296b57505050612555565b505050808061297990614caf565b9150506128c8565b6000838152600e6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938301828280156129e057602002820191906000526020600020905b8154815260200190600101908083116129cc575b505050505090506129f381858351613656565b95945050505050565b6000838152600f6020908152604080832061ffff851684528252808320805482518185028101850190935280835284938301828280156129e057602002820191906000526020600020908154815260200190600101908083116129cc57505050505090506129f381858351613656565b6060806000612a7b6013613640565b905060008167ffffffffffffffff811115612a9857612a98614ed1565b604051908082528060200260200182016040528015612ac1578160200160208202803683370190505b50905060008267ffffffffffffffff811115612adf57612adf614ed1565b604051908082528060200260200182016040528015612b08578160200160208202803683370190505b50905060005b8381101561146c576000612b2360138361364a565b6000818152600d6020908152604080832080548251818502810185019093528083529495509293909291830182828015612b7c57602002820191906000526020600020905b815481526020019060010190808311612b68575b5050505050905081858481518110612b9657612b96614ea2565b602002602001018181525050612bae818a8351613656565b848481518110612bc057612bc0614ea2565b60200260200101818152505050508080612bd990614da9565b915050612b0e565b6040805161014081018252600461010082019081527f746573740000000000000000000000000000000000000000000000000000000061012083015281528151602081810184526000808352818401929092523083850181905263ffffffff8916606085015260808401528351808201855282815260a08401528351908101909352825260c08101919091526bffffffffffffffffffffffff841660e082015260165460155473ffffffffffffffffffffffffffffffffffffffff9182169163095ea7b39116612cb460ff8a1688614b6c565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff166024820152604401602060405180830381600087803b158015612d2d57600080fd5b505af1158015612d41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d6591906141b5565b5060008660ff1667ffffffffffffffff811115612d8457612d84614ed1565b604051908082528060200260200182016040528015612dad578160200160208202803683370190505b50905060005b8760ff168160ff161015612e20576000612dcc8461396f565b905080838360ff1681518110612de457612de4614ea2565b60209081029190910181019190915260009182526008815260408083208890556007909152902084905580612e1881614df5565b915050612db3565b507f2ee10f7eb180441fb9fbba75b10c0162b5390b557712c93426243ca8f383c71181604051612e509190614775565b60405180910390a150505050505050565b6016546017546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b390604401602060405180830381600087803b158015612ee457600080fd5b505af1158015612ef8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f1c91906141b5565b506017546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401600060405180830381600087803b158015612f9d57600080fd5b505af1158015612fb1573d6000803e3d6000fd5b5050604080518581526bffffffffffffffffffffffff851660208201527f8137dc366612bf502338bd8951f835ad8ceba421c4eb3d79c7f9b3ce0ac4762e9350019050611710565b6000828152600f6020908152604080832061ffff85168452825291829020805483518184028101840190945280845260609392830182828015611fc75760200282019190600052602060002090815481526020019060010190808311611fb3575050505050905092915050565b8060005b8181101561313157600084848381811061308657613086614ea2565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16639fab438682836040516020016130bf91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016130eb92919061494f565b600060405180830381600087803b15801561310557600080fd5b505af1158015613119573d6000803e3d6000fd5b5050505050808061312990614da9565b91505061306a565b50505050565b6000818152600d602090815260409182902080548351818402810184019094528084526060939283018282801561318d57602002820191906000526020600020905b815481526020019060010190808311613179575b50505050509050919050565b6017546040517f9fab438600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690639fab4386906131f3908690869086906004016148fb565b600060405180830381600087803b15801561320d57600080fd5b505af1158015613221573d6000803e3d6000fd5b5050506000848152600b6020526040902061313191508383613e8a565b60008181526011602052604081205461ffff1681805b8261ffff168161ffff16116118bb576000858152600f6020908152604080832061ffff8516845290915290205461328b9083614a9b565b91508061329781614d87565b915050613254565b6017546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b15801561331757600080fd5b505af115801561332b573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6017546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b1580156133bc57600080fd5b505af1158015612071573d6000803e3d6000fd5b8060005b818163ffffffff161015613131573063af953a4a858563ffffffff851681811061340057613400614ea2565b905060200201356040518263ffffffff1660e01b815260040161342591815260200190565b600060405180830381600087803b15801561343f57600080fd5b505af1158015613453573d6000803e3d6000fd5b50505050808061346290614ddb565b9150506133d4565b6017546040517fc80480220000000000000000000000000000000000000000000000000000000081526004810183905273ffffffffffffffffffffffffffffffffffffffff9091169063c804802290602401600060405180830381600087803b1580156134d657600080fd5b505af11580156134ea573d6000803e3d6000fd5b50505050612267816013613a7190919063ffffffff16565b6000828152600d6020908152604080832080548251818502810185019093528083528493849392919083018282801561355a57602002820191906000526020600020905b815481526020019060010190808311613546575b505050505090506114ee81856137b6565b600d602052816000526040600020818154811061358757600080fd5b90600052602060002001600091509150505481565b6135a46138ec565b6135ad81613a7d565b50565b6010602052816000526040600020818154811061358757600080fd5b6000828152600f6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156114dc57602002820191906000526020600020908154815260200190600101908083116114c857505050505090506114ee8182516137b6565b600061112b825490565b60006117808383613b73565b8251600090819083158061366a5750808410155b15613673578093505b60008467ffffffffffffffff81111561368e5761368e614ed1565b6040519080825280602002602001820160405280156136b7578160200160208202803683370190505b509050600092505b84831015613725578660016136d48585614c14565b6136de9190614c14565b815181106136ee576136ee614ea2565b602002602001015181848151811061370857613708614ea2565b60209081029190910101528261371d81614da9565b9350506136bf565b61373e816000600184516137399190614c14565b613b9d565b85606414156137785780600182516137569190614c14565b8151811061376657613766614ea2565b60200260200101519350505050611780565b8060648251886137889190614b2f565b6137929190614b1b565b815181106137a2576137a2614ea2565b602002602001015193505050509392505050565b8151600090819081908415806137cc5750808510155b156137d5578094505b60008092505b85831015613831578660016137f08585614c14565b6137fa9190614c14565b8151811061380a5761380a614ea2565b60200260200101518161381d9190614a9b565b90508261382981614da9565b9350506137db565b9694955050505050565b60007f0000000000000000000000000000000000000000000000000000000000000000156138e757606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156138aa57600080fd5b505afa1580156138be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138e291906143ca565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff16331461396d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401612784565b565b6015546040517f08b79da4000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff909116906308b79da4906139ca9086906004016147db565b602060405180830381600087803b1580156139e457600080fd5b505af11580156139f8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a1c91906143ca565b9050613a29601382613d1e565b5060608301516000828152600a6020908152604080832063ffffffff90941690935560a0860151600b8252929091208251613a6a9391929190910190613f2c565b5092915050565b60006117808383613d2a565b73ffffffffffffffffffffffffffffffffffffffff8116331415613afd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401612784565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000826000018281548110613b8a57613b8a614ea2565b9060005260206000200154905092915050565b818180821415613bae575050505050565b6000856002613bbd8787614ba0565b613bc79190614ab3565b613bd19087614a27565b81518110613be157613be1614ea2565b602002602001015190505b818313613cf0575b80868481518110613c0757613c07614ea2565b60200260200101511015613c275782613c1f81614d4e565b935050613bf4565b858281518110613c3957613c39614ea2565b6020026020010151811015613c5a5781613c5281614c57565b925050613c27565b818313613ceb57858281518110613c7357613c73614ea2565b6020026020010151868481518110613c8d57613c8d614ea2565b6020026020010151878581518110613ca757613ca7614ea2565b60200260200101888581518110613cc057613cc0614ea2565b60209081029190910101919091525282613cd981614d4e565b9350508180613ce790614c57565b9250505b613bec565b81851215613d0357613d03868684613b9d565b83831215613d1657613d16868486613b9d565b505050505050565b60006117808383613e1d565b60008181526001830160205260408120548015613e13576000613d4e600183614c14565b8554909150600090613d6290600190614c14565b9050818114613dc7576000866000018281548110613d8257613d82614ea2565b9060005260206000200154905080876000018481548110613da557613da5614ea2565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613dd857613dd8614e73565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061112b565b600091505061112b565b6000818152600183016020526040812054613e645750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561112b565b50600061112b565b50805460008255906000526020600020908101906135ad9190613fa0565b828054613e9690614cfa565b90600052602060002090601f016020900481019282613eb85760008555613f1c565b82601f10613eef578280017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00823516178555613f1c565b82800160010185558215613f1c579182015b82811115613f1c578235825591602001919060010190613f01565b50613f28929150613fa0565b5090565b828054613f3890614cfa565b90600052602060002090601f016020900481019282613f5a5760008555613f1c565b82601f10613f7357805160ff1916838001178555613f1c565b82800160010185558215613f1c579182015b82811115613f1c578251825591602001919060010190613f85565b5b80821115613f285760008155600101613fa1565b8051613fc081614f00565b919050565b60008083601f840112613fd757600080fd5b50813567ffffffffffffffff811115613fef57600080fd5b6020830191508360208260051b85010111156114f457600080fd5b80518015158114613fc057600080fd5b60008083601f84011261402c57600080fd5b50813567ffffffffffffffff81111561404457600080fd5b6020830191508360208285010111156114f457600080fd5b600082601f83011261406d57600080fd5b815161408061407b826149e1565b614992565b81815284602083860101111561409557600080fd5b6140a6826020830160208701614c2b565b949350505050565b803561ffff81168114613fc057600080fd5b8051613fc081614f22565b805167ffffffffffffffff81168114613fc057600080fd5b803560ff81168114613fc057600080fd5b8051613fc081614f34565b60006020828403121561411157600080fd5b813561178081614f00565b6000806020838503121561412f57600080fd5b823567ffffffffffffffff81111561414657600080fd5b61415285828601613fc5565b90969095509350505050565b60008060006040848603121561417357600080fd5b833567ffffffffffffffff81111561418a57600080fd5b61419686828701613fc5565b90945092505060208401356141aa81614f22565b809150509250925092565b6000602082840312156141c757600080fd5b6117808261400a565b6000602082840312156141e257600080fd5b5035919050565b600080602083850312156141fc57600080fd5b823567ffffffffffffffff81111561421357600080fd5b6141528582860161401a565b60006020828403121561423157600080fd5b815161178081614f00565b600080600080600060a0868803121561425457600080fd5b85516003811061426357600080fd5b602087015190955061427481614f22565b604087015190945061428581614f22565b606087015190935061429681614f00565b80925050608086015190509295509295909350565b6000602082840312156142bd57600080fd5b815167ffffffffffffffff808211156142d557600080fd5b9083019061014082860312156142ea57600080fd5b6142f2614968565b6142fb83613fb5565b8152614309602084016140c0565b602082015260408301518281111561432057600080fd5b61432c8782860161405c565b60408301525061433e606084016140f4565b606082015261434f60808401613fb5565b608082015261436060a084016140cb565b60a082015261437160c084016140c0565b60c082015261438260e084016140f4565b60e082015261010061439581850161400a565b9082015261012083810151838111156143ad57600080fd5b6143b98882870161405c565b918301919091525095945050505050565b6000602082840312156143dc57600080fd5b5051919050565b6000806000604084860312156143f857600080fd5b83359250602084013567ffffffffffffffff81111561441657600080fd5b6144228682870161401a565b9497909650939450505050565b6000806040838503121561444257600080fd5b82359150602083013567ffffffffffffffff81111561446057600080fd5b8301601f8101851361447157600080fd5b803561447f61407b826149e1565b81815286602083850101111561449457600080fd5b816020840160208301376000602083830101528093505050509250929050565b600080604083850312156144c757600080fd5b823591506144d7602084016140ae565b90509250929050565b6000806000606084860312156144f557600080fd5b83359250614505602085016140ae565b9150604084013590509250925092565b6000806040838503121561452857600080fd5b50508035926020909101359150565b60008060006060848603121561454c57600080fd5b8335925060208401359150614563604085016140ae565b90509250925092565b60008060006060848603121561458157600080fd5b505081359360208301359350604090920135919050565b600080604083850312156145ab57600080fd5b8235915060208301356145bd81614f22565b809150509250929050565b600080604083850312156145db57600080fd5b8235915060208301356145bd81614f34565b6000602082840312156145ff57600080fd5b611780826140e3565b600080600080600060a0868803121561462057600080fd5b614629866140e3565b9450602086013561463981614f22565b9350604086013561464981614f34565b94979396509394606081013594506080013592915050565b60006020828403121561467357600080fd5b813561178081614f34565b60006020828403121561469057600080fd5b815161178081614f34565b600081518084526020808501945080840160005b838110156146cb578151875295820195908201906001016146af565b509495945050505050565b600081518084526146ee816020860160208601614c2b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815281602082015260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83111561475957600080fd5b8260051b80856040850137600092016040019182525092915050565b602081526000611780602083018461469b565b60408152600061479b604083018561469b565b82810360208401526129f3818561469b565b821515815260406020820152600061177d60408301846146d6565b60208152600061178060208301846146d6565b60208152600082516101008060208501526147fa6101208501836146d6565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08086850301604087015261483684836146d6565b935060408701519150614861606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a08701519150808685030160c08701526148b284836146d6565b935060c08701519150808685030160e0870152506148d083826146d6565b92505060e08501516148f1828601826bffffffffffffffffffffffff169052565b5090949350505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b82815260406020820152600061177d60408301846146d6565b604051610140810167ffffffffffffffff8111828210171561498c5761498c614ed1565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156149d9576149d9614ed1565b604052919050565b600067ffffffffffffffff8211156149fb576149fb614ed1565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b6000808212827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03841381151615614a6157614a61614e15565b827f8000000000000000000000000000000000000000000000000000000000000000038412811615614a9557614a95614e15565b50500190565b60008219821115614aae57614aae614e15565b500190565b600082614ac257614ac2614e44565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615614b1657614b16614e15565b500590565b600082614b2a57614b2a614e44565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615614b6757614b67614e15565b500290565b60006bffffffffffffffffffffffff80831681851681830481118215151615614b9757614b97614e15565b02949350505050565b6000808312837f800000000000000000000000000000000000000000000000000000000000000001831281151615614bda57614bda614e15565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff018313811615614c0e57614c0e614e15565b50500390565b600082821015614c2657614c26614e15565b500390565b60005b83811015614c46578181015183820152602001614c2e565b838111156131315750506000910152565b60007f8000000000000000000000000000000000000000000000000000000000000000821415614c8957614c89614e15565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600061ffff821680614cc357614cc3614e15565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b600081614c8957614c89614e15565b600181811c90821680614d0e57607f821691505b60208210811415614d48577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415614d8057614d80614e15565b5060010190565b600061ffff80831681811415614d9f57614d9f614e15565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415614d8057614d80614e15565b600063ffffffff80831681811415614d9f57614d9f614e15565b600060ff821660ff811415614e0c57614e0c614e15565b60010192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff811681146135ad57600080fd5b63ffffffff811681146135ad57600080fd5b6bffffffffffffffffffffffff811681146135ad57600080fdfea164736f6c6343000806000a",
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
