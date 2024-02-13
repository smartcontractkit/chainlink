// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifiable_load_streams_lookup_upkeep_wrapper

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

type KeeperRegistryBase21UpkeepInfo struct {
	Target                   common.Address
	PerformGas               uint32
	CheckData                []byte
	Balance                  *big.Int
	Admin                    common.Address
	MaxValidBlocknumber      uint64
	LastPerformedBlockNumber uint32
	AmountSpent              *big.Int
	Paused                   bool
	OffchainConfig           []byte
}

var VerifiableLoadStreamsLookupUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"errCode\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c0604052604261014081815261010091829190620062496101603981526020016040518060800160405280604281526020016200628b604291399052620000be906016906002620003c7565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000ee908262000543565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b602082015260189062000120908262000543565b503480156200012e57600080fd5b50604051620062cd380380620062cd833981016040819052620001519162000625565b81813380600081620001aa5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001dd57620001dd816200031c565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa1580156200023a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000260919062000668565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002c6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ec919062000699565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c05250620006c0915050565b336001600160a01b03821603620003765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620001a1565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821562000412579160200282015b8281111562000412578251829062000401908262000543565b5091602001919060010190620003e8565b506200042092915062000424565b5090565b80821115620004205760006200043b828262000445565b5060010162000424565b5080546200045390620004b4565b6000825580601f1062000464575050565b601f01602090049060005260206000209081019062000484919062000487565b50565b5b8082111562000420576000815560010162000488565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004c957607f821691505b602082108103620004ea57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200053e57600081815260208120601f850160051c81016020861015620005195750805b601f850160051c820191505b818110156200053a5782815560010162000525565b5050505b505050565b81516001600160401b038111156200055f576200055f6200049e565b6200057781620005708454620004b4565b84620004f0565b602080601f831160018114620005af5760008415620005965750858301515b600019600386901b1c1916600185901b1785556200053a565b600085815260208120601f198616915b82811015620005e057888601518255948401946001909101908401620005bf565b5085821015620005ff5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200048457600080fd5b600080604083850312156200063957600080fd5b825162000646816200060f565b602084015190925080151581146200065d57600080fd5b809150509250929050565b600080604083850312156200067c57600080fd5b825162000689816200060f565b6020939093015192949293505050565b600060208284031215620006ac57600080fd5b8151620006b9816200060f565b9392505050565b60805160a05160c05160e051615b3362000716600039600081816105830152611fb80152600081816109d70152613ce50152600081816108420152613733015260008181610db701526137080152615b336000f3fe6080604052600436106104f05760003560e01c806379ea994311610294578063a6b594751161015e578063d6051a72116100d6578063e45530831161008a578063fa333dfb1161006f578063fa333dfb14610fc6578063fba7ffa314611079578063fcdc1f63146110a657600080fd5b8063e455308314610f90578063f2fde38b14610fa657600080fd5b8063daee1aeb116100bb578063daee1aeb14610f23578063dbef701e14610f43578063e0114adb14610f6357600080fd5b8063d6051a7214610ee3578063da6cba4714610f0357600080fd5b8063b657bc9c1161012d578063c041982211610112578063c041982214610e8e578063c98f10b014610eae578063d4c2490014610ec357600080fd5b8063b657bc9c14610e4e578063becde0e114610e6e57600080fd5b8063a6b5947514610dd9578063a72aa27e14610df9578063af953a4a14610e19578063afb28d1f14610e3957600080fd5b8063924ca5781161020c5780639b51fb0d116101c05780639d6f1cc7116101a55780639d6f1cc714610d62578063a4dc5e3314610d82578063a654824814610da557600080fd5b80639b51fb0d14610d115780639d385eaa14610d4257600080fd5b806396cebc7c116101f157806396cebc7c14610c9a5780639ac542eb14610cba5780639b42935414610ce457600080fd5b8063924ca57814610c5a578063948108f714610c7a57600080fd5b80638340507c11610263578063873c758611610248578063873c758614610be25780638da5cb5b14610c025780638fcb3fba14610c2d57600080fd5b80638340507c14610ba257806386e330af14610bc257600080fd5b806379ea994314610b155780637b10399914610b355780637e7a46dc14610b625780638243444a14610b8257600080fd5b806345d2ec17116103d557806360457ff51161034d5780637145f11b1161030157806376721303116102e65780637672130314610ab3578063776898c814610ae057806379ba509714610b0057600080fd5b80637145f11b14610a5657806373644cce14610a8657600080fd5b8063642f6cef11610332578063642f6cef146109c557806369cdbadb14610a095780636e04ff0d14610a3657600080fd5b806360457ff514610973578063636092e8146109a057600080fd5b80635147cd59116103a457806357970e931161038957806357970e93146109115780635d4ee7f31461093e5780635f17e6161461095357600080fd5b80635147cd59146108bf57806351c98be3146108f157600080fd5b806345d2ec1714610803578063469820931461083057806346e7a63e146108645780634b56a42e1461089157600080fd5b806320e3dbd4116104685780632a9032d311610437578063328ffd111161041c578063328ffd11146107965780633ebe8d6c146107c35780634585e33b146107e357600080fd5b80632a9032d3146107245780632b20e3971461074457600080fd5b806320e3dbd4146106975780632636aecf146106b757806328c4b57b146106d757806329e0a841146106f757600080fd5b806319d97a94116104bf5780631e010439116104a45780631e01043914610605578063206c32e814610642578063207b65161461067757600080fd5b806319d97a94146105b85780631cdde251146105e557600080fd5b806306c1cc00146104fc578063077ac6211461051e5780630b7d33e61461055157806312c550271461057157600080fd5b366104f757005b600080fd5b34801561050857600080fd5b5061051c610517366004614467565b6110d3565b005b34801561052a57600080fd5b5061053e61053936600461451a565b611322565b6040519081526020015b60405180910390f35b34801561055d57600080fd5b5061051c61056c36600461454f565b611360565b34801561057d57600080fd5b506105a57f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610548565b3480156105c457600080fd5b506105d86105d3366004614596565b6113ee565b604051610548919061461d565b3480156105f157600080fd5b5061051c610600366004614652565b6114ab565b34801561061157600080fd5b50610625610620366004614596565b6115e8565b6040516bffffffffffffffffffffffff9091168152602001610548565b34801561064e57600080fd5b5061066261065d3660046146b7565b61167d565b60408051928352602083019190915201610548565b34801561068357600080fd5b506105d8610692366004614596565b611700565b3480156106a357600080fd5b5061051c6106b23660046146e3565b611758565b3480156106c357600080fd5b5061051c6106d2366004614745565b611922565b3480156106e357600080fd5b5061053e6106f23660046147bf565b611beb565b34801561070357600080fd5b50610717610712366004614596565b611c56565b60405161054891906147eb565b34801561073057600080fd5b5061051c61073f36600461492c565b611d5b565b34801561075057600080fd5b506011546107719073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610548565b3480156107a257600080fd5b5061053e6107b1366004614596565b60036020526000908152604090205481565b3480156107cf57600080fd5b5061053e6107de366004614596565b611e3c565b3480156107ef57600080fd5b5061051c6107fe3660046149b0565b611ea5565b34801561080f57600080fd5b5061082361081e3660046146b7565b6120c4565b60405161054891906149e6565b34801561083c57600080fd5b5061053e7f000000000000000000000000000000000000000000000000000000000000000081565b34801561087057600080fd5b5061053e61087f366004614596565b600a6020526000908152604090205481565b34801561089d57600080fd5b506108b16108ac366004614a4e565b612133565b604051610548929190614b18565b3480156108cb57600080fd5b506108df6108da366004614596565b612187565b60405160ff9091168152602001610548565b3480156108fd57600080fd5b5061051c61090c366004614b33565b61221b565b34801561091d57600080fd5b506012546107719073ffffffffffffffffffffffffffffffffffffffff1681565b34801561094a57600080fd5b5061051c6122bf565b34801561095f57600080fd5b5061051c61096e366004614b8a565b6123fa565b34801561097f57600080fd5b5061053e61098e366004614596565b60076020526000908152604090205481565b3480156109ac57600080fd5b50601554610625906bffffffffffffffffffffffff1681565b3480156109d157600080fd5b506109f97f000000000000000000000000000000000000000000000000000000000000000081565b6040519015158152602001610548565b348015610a1557600080fd5b5061053e610a24366004614596565b60086020526000908152604090205481565b348015610a4257600080fd5b506108b1610a513660046149b0565b6124c7565b348015610a6257600080fd5b506109f9610a71366004614596565b600b6020526000908152604090205460ff1681565b348015610a9257600080fd5b5061053e610aa1366004614596565b6000908152600c602052604090205490565b348015610abf57600080fd5b5061053e610ace366004614596565b60046020526000908152604090205481565b348015610aec57600080fd5b506109f9610afb366004614596565b6126f0565b348015610b0c57600080fd5b5061051c612742565b348015610b2157600080fd5b50610771610b30366004614596565b61283f565b348015610b4157600080fd5b506013546107719073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b6e57600080fd5b5061051c610b7d366004614bac565b6128d3565b348015610b8e57600080fd5b5061051c610b9d366004614bac565b612964565b348015610bae57600080fd5b5061051c610bbd366004614bf8565b6129be565b348015610bce57600080fd5b5061051c610bdd366004614c45565b6129dc565b348015610bee57600080fd5b50610823610bfd366004614b8a565b6129ef565b348015610c0e57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610771565b348015610c3957600080fd5b5061053e610c48366004614596565b60056020526000908152604090205481565b348015610c6657600080fd5b5061051c610c75366004614b8a565b612aac565b348015610c8657600080fd5b5061051c610c95366004614cf6565b612cf1565b348015610ca657600080fd5b5061051c610cb5366004614d26565b612e09565b348015610cc657600080fd5b506015546108df906c01000000000000000000000000900460ff1681565b348015610cf057600080fd5b5061051c610cff366004614b8a565b60009182526009602052604090912055565b348015610d1d57600080fd5b506105a5610d2c366004614596565b600e6020526000908152604090205461ffff1681565b348015610d4e57600080fd5b50610823610d5d366004614596565b613013565b348015610d6e57600080fd5b506105d8610d7d366004614596565b613075565b348015610d8e57600080fd5b506108b1610d9d366004614d43565b600192909150565b348015610db157600080fd5b5061053e7f000000000000000000000000000000000000000000000000000000000000000081565b348015610de557600080fd5b5061051c610df43660046147bf565b613121565b348015610e0557600080fd5b5061051c610e14366004614d7d565b61318a565b348015610e2557600080fd5b5061051c610e34366004614596565b613235565b348015610e4557600080fd5b506105d86132bb565b348015610e5a57600080fd5b50610625610e69366004614596565b6132c8565b348015610e7a57600080fd5b5061051c610e8936600461492c565b613320565b348015610e9a57600080fd5b50610823610ea9366004614b8a565b6133ba565b348015610eba57600080fd5b506105d86134b7565b348015610ecf57600080fd5b5061051c610ede366004614da2565b6134c4565b348015610eef57600080fd5b50610662610efe366004614b8a565b613543565b348015610f0f57600080fd5b5061051c610f1e366004614dc7565b6135ac565b348015610f2f57600080fd5b5061051c610f3e36600461492c565b613913565b348015610f4f57600080fd5b5061053e610f5e366004614b8a565b6139de565b348015610f6f57600080fd5b5061053e610f7e366004614596565b60096020526000908152604090205481565b348015610f9c57600080fd5b5061053e60145481565b348015610fb257600080fd5b5061051c610fc13660046146e3565b613a0f565b348015610fd257600080fd5b506105d8610fe1366004614e2f565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b34801561108557600080fd5b5061053e611094366004614596565b60066020526000908152604090205481565b3480156110b257600080fd5b5061053e6110c1366004614596565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b39216906111b9908c1688614eb7565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af1158015611237573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061125b9190614efb565b5060008860ff1667ffffffffffffffff81111561127a5761127a614309565b6040519080825280602002602001820160405280156112a3578160200160208202803683370190505b50905060005b8960ff168160ff1610156113165760006112c284613a23565b905080838360ff16815181106112da576112da614f16565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061130e81614f45565b9150506112a9565b50505050505050505050565b600d602052826000526040600020602052816000526040600020818154811061134a57600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e6906113b89085908590600401614f64565b600060405180830381600087803b1580156113d257600080fd5b505af11580156113e6573d6000803e3d6000fd5b505050505050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa15801561145f573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114a59190810190614fca565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa15801561154a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115909190810190614fca565b6040518363ffffffff1660e01b81526004016115ad929190614f64565b600060405180830381600087803b1580156115c757600080fd5b505af11580156115db573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e010439906024015b602060405180830381865afa158015611659573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114a5919061500a565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156116e157602002820191906000526020600020905b8154815260200190600101908083116116cd575b505050505090506116f3818251613af1565b92509250505b9250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b651690602401611442565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa1580156117ee573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118129190615032565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa1580156118b5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118d99190615060565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611be057600089898381811061194257611942614f16565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161197b91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016119a7929190614f64565b600060405180830381600087803b1580156119c157600080fd5b505af11580156119d5573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015611a4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a6f919061507d565b90508060ff16600103611bcb576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611af8573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b3e9190810190614fca565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611b979086908590600401614f64565b600060405180830381600087803b158015611bb157600080fd5b505af1158015611bc5573d6000803e3d6000fd5b50505050505b50508080611bd89061509a565b915050611926565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611c4c93830182828015611c4057602002820191906000526020600020905b815481526020019060010190808311611c2c575b50505050508484613b76565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611d15573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114a591908101906150f5565b8060005b818160ff161015611e365760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611d9d57611d9d614f16565b905060200201356040518263ffffffff1660e01b8152600401611dc291815260200190565b600060405180830381600087803b158015611ddc57600080fd5b505af1158015611df0573d6000803e3d6000fd5b50505050611e2384848360ff16818110611e0c57611e0c614f16565b90506020020135600f613cd590919063ffffffff16565b5080611e2e81614f45565b915050611d5f565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611e9d576000858152600d6020908152604080832061ffff85168452909152902054611e899083615214565b915080611e9581615227565b915050611e52565b509392505050565b60005a9050600080611eb984860186614a4e565b91509150600081806020019051810190611ed39190615248565b60008181526005602090815260408083205460049092528220549293509190611efa613ce1565b905082600003611f1a576000848152600560205260409020819055612075565b600084815260036020526040812054611f338484615261565b611f3d9190615261565b6000868152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611faf57602002820191906000526020600020905b815481526020019060010190808311611f9b575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff1681510361202a5781611fec81615227565b6000898152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000868152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055888352600c8252928220805493840181558252902001555b60008481526006602052604081205461208f906001615214565b60008681526006602090815260408083208490556004909152902083905590506120b98583612aac565b611316858984613121565b6000828152600d6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561212657602002820191906000526020600020905b815481526020019060010190808311612112575b5050505050905092915050565b600060606000848460405160200161214c929190615274565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa1580156121f7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114a5919061507d565b8160005b818110156122b85730635f17e61686868481811061223f5761223f614f16565b90506020020135856040518363ffffffff1660e01b815260040161227392919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561228d57600080fd5b505af11580156122a1573d6000803e3d6000fd5b5050505080806122b09061509a565b91505061221f565b5050505050565b6122c7613d83565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612336573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061235a9190615248565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156123d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123f69190614efb565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c909152812061243291614208565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff161161248e576000848152600d6020908152604080832061ffff85168452909152812061247c91614208565b8061248681615227565b915050612447565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000606060005a905060006124de85870187614596565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff81111561251757612517614309565b6040519080825280601f01601f191660200182016040528015612541576020820181803683370190505b50604051602001612553929190614f64565b6040516020818303038152906040529050600061256e613ce1565b9050600061257b866126f0565b90505b835a61258a9089615261565b61259690612710615214565b10156125d75781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905561257e565b806125ef5760008398509850505050505050506116f9565b6040517f66656564496448657800000000000000000000000000000000000000000000006020820152600090602901604051602081830303815290604052805190602001206017604051602001612646919061535b565b604051602081830303815290604052805190602001200361266857508161266b565b50425b601760166018838a60405160200161268591815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527ff055e4a20000000000000000000000000000000000000000000000000000000082526126e7959493929160040161548a565b60405180910390fd5b600081815260056020526040812054810361270d57506001919050565b600082815260036020908152604080832054600490925290912054612730613ce1565b61273a9190615261565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146127c3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016126e7565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa1580156128af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114a59190615060565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b59061292d9086908690869060040161554d565b600060405180830381600087803b15801561294757600080fd5b505af115801561295b573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d359061292d9086908690869060040161554d565b60176129ca83826155e7565b5060186129d782826155e7565b505050565b80516123f6906016906020840190614226565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612a66573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611c4f9190810190615701565b601454600083815260026020526040902054612ac89083615261565b11156123f6576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612b3e573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612b8491908101906150f5565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612bf9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c1d919061500a565b601554909150612c419082906c01000000000000000000000000900460ff16614eb7565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611e3657601554612c849085906bffffffffffffffffffffffff16612cf1565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612d79573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d9d9190614efb565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f7906044016113b8565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015612e68573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612eae9190810190615701565b80519091506000612ebd613ce1565b905060005b828110156122b8576000848281518110612ede57612ede614f16565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612f5e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f82919061507d565b90508060ff16600103612ffe578660ff16600003612fce576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4612ffe565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b5050808061300b9061509a565b915050612ec2565b6000818152600c602090815260409182902080548351818402810184019094528084526060939283018282801561306957602002820191906000526020600020905b815481526020019060010190808311613055575b50505050509050919050565b6016818154811061308557600080fd5b9060005260206000200160009150905080546130a090615308565b80601f01602080910402602001604051908101604052809291908181526020018280546130cc90615308565b80156131195780601f106130ee57610100808354040283529160200191613119565b820191906000526020600020905b8154815290600101906020018083116130fc57829003601f168201915b505050505081565b6000838152600760205260409020545b805a61313d9085615261565b61314990612710615214565b1015611e365781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055613131565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b15801561320257600080fd5b505af1158015613216573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b1580156132a757600080fd5b505af11580156122b8573d6000803e3d6000fd5b601780546130a090615308565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff169063b657bc9c9060240161163c565b8060005b818163ffffffff161015611e36573063af953a4a858563ffffffff851681811061335057613350614f16565b905060200201356040518263ffffffff1660e01b815260040161337591815260200190565b600060405180830381600087803b15801561338f57600080fd5b505af11580156133a3573d6000803e3d6000fd5b5050505080806133b290615792565b915050613324565b606060006133c8600f613e06565b9050808410613403576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82600003613418576134158482615261565b92505b60008367ffffffffffffffff81111561343357613433614309565b60405190808252806020026020018201604052801561345c578160200160208202803683370190505b50905060005b848110156134ae5761347f6134778288615214565b600f90613e10565b82828151811061349157613491614f16565b6020908102919091010152806134a68161509a565b915050613462565b50949350505050565b601880546130a090615308565b60006134ce613ce1565b90508160ff1660000361350f576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a4505050565b6000828152600c6020908152604080832080548251818502810185019093528083528493849392919083018282801561359b57602002820191906000526020600020905b815481526020019060010190808311613587575b505050505090506116f38185613af1565b8260005b818110156113e65760008686838181106135cc576135cc614f16565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161360591815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613631929190614f64565b600060405180830381600087803b15801561364b57600080fd5b505af115801561365f573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa1580156136d5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136f9919061507d565b90508060ff166001036138fe577f000000000000000000000000000000000000000000000000000000000000000060ff87161561375357507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb3089858860405160200161378791815260200190565b60405160208183030381529060405261379f906157ab565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa15801561382a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526138709190810190614fca565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906138c99087908590600401614f64565b600060405180830381600087803b1580156138e357600080fd5b505af11580156138f7573d6000803e3d6000fd5b5050505050505b5050808061390b9061509a565b9150506135b0565b8060005b81811015611e3657600084848381811061393357613933614f16565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161396c91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613998929190614f64565b600060405180830381600087803b1580156139b257600080fd5b505af11580156139c6573d6000803e3d6000fd5b505050505080806139d69061509a565b915050613917565b600c60205281600052604060002081815481106139fa57600080fd5b90600052602060002001600091509150505481565b613a17613d83565b613a2081613e1c565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613a7e9086906004016157ed565b6020604051808303816000875af1158015613a9d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ac19190615248565b9050613ace600f82613f11565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b815160009081908190841580613b075750808510155b15613b10578094505b60008092505b85831015613b6c57866001613b2b8585615261565b613b359190615261565b81518110613b4557613b45614f16565b602002602001015181613b589190615214565b905082613b648161509a565b935050613b16565b9694955050505050565b82516000908190831580613b8a5750808410155b15613b93578093505b60008467ffffffffffffffff811115613bae57613bae614309565b604051908082528060200260200182016040528015613bd7578160200160208202803683370190505b509050600092505b84831015613c4557866001613bf48585615261565b613bfe9190615261565b81518110613c0e57613c0e614f16565b6020026020010151818481518110613c2857613c28614f16565b602090810291909101015282613c3d8161509a565b935050613bdf565b613c5e81600060018451613c599190615261565b613f1d565b85606403613c97578060018251613c759190615261565b81518110613c8557613c85614f16565b60200260200101519350505050611c4f565b806064825188613ca7919061593f565b613cb191906159ab565b81518110613cc157613cc1614f16565b602002602001015193505050509392505050565b6000611c4f8383614095565b60007f000000000000000000000000000000000000000000000000000000000000000015613d7e57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613d55573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d799190615248565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613e04576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016126e7565b565b60006114a5825490565b6000611c4f838361418f565b3373ffffffffffffffffffffffffffffffffffffffff821603613e9b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016126e7565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611c4f83836141b9565b8181808203613f2d575050505050565b6000856002613f3c87876159bf565b613f4691906159df565b613f509087615a47565b81518110613f6057613f60614f16565b602002602001015190505b81831361406f575b80868481518110613f8657613f86614f16565b60200260200101511015613fa65782613f9e81615a6f565b935050613f73565b858281518110613fb857613fb8614f16565b6020026020010151811015613fd95781613fd181615aa0565b925050613fa6565b81831361406a57858281518110613ff257613ff2614f16565b602002602001015186848151811061400c5761400c614f16565b602002602001015187858151811061402657614026614f16565b6020026020010188858151811061403f5761403f614f16565b6020908102919091010191909152528261405881615a6f565b935050818061406690615aa0565b9250505b613f6b565b8185121561408257614082868684613f1d565b838312156113e6576113e6868486613f1d565b6000818152600183016020526040812054801561417e5760006140b9600183615261565b85549091506000906140cd90600190615261565b90508181146141325760008660000182815481106140ed576140ed614f16565b906000526020600020015490508087600001848154811061411057614110614f16565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061414357614143615af7565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506114a5565b60009150506114a5565b5092915050565b60008260000182815481106141a6576141a6614f16565b9060005260206000200154905092915050565b6000818152600183016020526040812054614200575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556114a5565b5060006114a5565b5080546000825590600052602060002090810190613a20919061427c565b82805482825590600052602060002090810192821561426c579160200282015b8281111561426c578251829061425c90826155e7565b5091602001919060010190614246565b50614278929150614291565b5090565b5b80821115614278576000815560010161427d565b808211156142785760006142a582826142ae565b50600101614291565b5080546142ba90615308565b6000825580601f106142ca575050565b601f016020900490600052602060002090810190613a20919061427c565b60ff81168114613a2057600080fd5b63ffffffff81168114613a2057600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff8111828210171561435c5761435c614309565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156143a9576143a9614309565b604052919050565b600067ffffffffffffffff8211156143cb576143cb614309565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261440857600080fd5b813561441b614416826143b1565b614362565b81815284602083860101111561443057600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff81168114613a2057600080fd5b600080600080600080600060e0888a03121561448257600080fd5b873561448d816142e8565b9650602088013561449d816142f7565b955060408801356144ad816142e8565b9450606088013567ffffffffffffffff8111156144c957600080fd5b6144d58a828b016143f7565b94505060808801356144e68161444d565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff8116811461451557600080fd5b919050565b60008060006060848603121561452f57600080fd5b8335925061453f60208501614503565b9150604084013590509250925092565b6000806040838503121561456257600080fd5b82359150602083013567ffffffffffffffff81111561458057600080fd5b61458c858286016143f7565b9150509250929050565b6000602082840312156145a857600080fd5b5035919050565b60005b838110156145ca5781810151838201526020016145b2565b50506000910152565b600081518084526145eb8160208601602086016145af565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611c4f60208301846145d3565b73ffffffffffffffffffffffffffffffffffffffff81168114613a2057600080fd5b600080600080600080600060e0888a03121561466d57600080fd5b87359650602088013561467f81614630565b9550604088013561468f816142e8565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b600080604083850312156146ca57600080fd5b823591506146da60208401614503565b90509250929050565b6000602082840312156146f557600080fd5b8135611c4f81614630565b60008083601f84011261471257600080fd5b50813567ffffffffffffffff81111561472a57600080fd5b6020830191508360208260051b85010111156116f957600080fd5b600080600080600080600060c0888a03121561476057600080fd5b873567ffffffffffffffff81111561477757600080fd5b6147838a828b01614700565b9098509650506020880135614797816142e8565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b6000806000606084860312156147d457600080fd5b505081359360208301359350604090920135919050565b6020815261481260208201835173ffffffffffffffffffffffffffffffffffffffff169052565b6000602083015161482b604084018263ffffffff169052565b5060408301516101408060608501526148486101608501836145d3565b9150606085015161486960808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e08501516101006148d5818701836bffffffffffffffffffffffff169052565b86015190506101206148ea8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00183870152905061492283826145d3565b9695505050505050565b6000806020838503121561493f57600080fd5b823567ffffffffffffffff81111561495657600080fd5b61496285828601614700565b90969095509350505050565b60008083601f84011261498057600080fd5b50813567ffffffffffffffff81111561499857600080fd5b6020830191508360208285010111156116f957600080fd5b600080602083850312156149c357600080fd5b823567ffffffffffffffff8111156149da57600080fd5b6149628582860161496e565b6020808252825182820181905260009190848201906040850190845b81811015614a1e57835183529284019291840191600101614a02565b50909695505050505050565b600067ffffffffffffffff821115614a4457614a44614309565b5060051b60200190565b60008060408385031215614a6157600080fd5b823567ffffffffffffffff80821115614a7957600080fd5b818501915085601f830112614a8d57600080fd5b81356020614a9d61441683614a2a565b82815260059290921b84018101918181019089841115614abc57600080fd5b8286015b84811015614af457803586811115614ad85760008081fd5b614ae68c86838b01016143f7565b845250918301918301614ac0565b5096505086013592505080821115614b0b57600080fd5b5061458c858286016143f7565b8215158152604060208201526000611c4c60408301846145d3565b600080600060408486031215614b4857600080fd5b833567ffffffffffffffff811115614b5f57600080fd5b614b6b86828701614700565b9094509250506020840135614b7f816142f7565b809150509250925092565b60008060408385031215614b9d57600080fd5b50508035926020909101359150565b600080600060408486031215614bc157600080fd5b83359250602084013567ffffffffffffffff811115614bdf57600080fd5b614beb8682870161496e565b9497909650939450505050565b60008060408385031215614c0b57600080fd5b823567ffffffffffffffff80821115614c2357600080fd5b614c2f868387016143f7565b93506020850135915080821115614b0b57600080fd5b60006020808385031215614c5857600080fd5b823567ffffffffffffffff80821115614c7057600080fd5b818501915085601f830112614c8457600080fd5b8135614c9261441682614a2a565b81815260059190911b83018401908481019088831115614cb157600080fd5b8585015b83811015614ce957803585811115614ccd5760008081fd5b614cdb8b89838a01016143f7565b845250918601918601614cb5565b5098975050505050505050565b60008060408385031215614d0957600080fd5b823591506020830135614d1b8161444d565b809150509250929050565b600060208284031215614d3857600080fd5b8135611c4f816142e8565b60008060408385031215614d5657600080fd5b8235614d61816142f7565b9150602083013567ffffffffffffffff81111561458057600080fd5b60008060408385031215614d9057600080fd5b823591506020830135614d1b816142f7565b60008060408385031215614db557600080fd5b823591506020830135614d1b816142e8565b60008060008060608587031215614ddd57600080fd5b843567ffffffffffffffff811115614df457600080fd5b614e0087828801614700565b9095509350506020850135614e14816142e8565b91506040850135614e24816142e8565b939692955090935050565b60008060008060008060c08789031215614e4857600080fd5b8635614e5381614630565b95506020870135614e63816142e8565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614ee257614ee2614e88565b02949350505050565b8051801515811461451557600080fd5b600060208284031215614f0d57600080fd5b611c4f82614eeb565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614f5b57614f5b614e88565b60010192915050565b828152604060208201526000611c4c60408301846145d3565b600082601f830112614f8e57600080fd5b8151614f9c614416826143b1565b818152846020838601011115614fb157600080fd5b614fc28260208301602087016145af565b949350505050565b600060208284031215614fdc57600080fd5b815167ffffffffffffffff811115614ff357600080fd5b614fc284828501614f7d565b80516145158161444d565b60006020828403121561501c57600080fd5b8151611c4f8161444d565b805161451581614630565b6000806040838503121561504557600080fd5b825161505081614630565b6020939093015192949293505050565b60006020828403121561507257600080fd5b8151611c4f81614630565b60006020828403121561508f57600080fd5b8151611c4f816142e8565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036150cb576150cb614e88565b5060010190565b8051614515816142f7565b805167ffffffffffffffff8116811461451557600080fd5b60006020828403121561510757600080fd5b815167ffffffffffffffff8082111561511f57600080fd5b90830190610140828603121561513457600080fd5b61513c614338565b61514583615027565b8152615153602084016150d2565b602082015260408301518281111561516a57600080fd5b61517687828601614f7d565b60408301525061518860608401614fff565b606082015261519960808401615027565b60808201526151aa60a084016150dd565b60a08201526151bb60c084016150d2565b60c08201526151cc60e08401614fff565b60e08201526101006151df818501614eeb565b9082015261012083810151838111156151f757600080fd5b61520388828701614f7d565b918301919091525095945050505050565b808201808211156114a5576114a5614e88565b600061ffff80831681810361523e5761523e614e88565b6001019392505050565b60006020828403121561525a57600080fd5b5051919050565b818103818111156114a5576114a5614e88565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156152e9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526152d78683516145d3565b9550938201939082019060010161529d565b5050858403818701525050506152ff81856145d3565b95945050505050565b600181811c9082168061531c57607f821691505b602082108103615355577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600080835461536981615308565b6001828116801561538157600181146153b4576153e3565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841687528215158302870194506153e3565b8760005260208060002060005b858110156153da5781548a8201529084019082016153c1565b50505082870194505b50929695505050505050565b600081546153fc81615308565b80855260206001838116801561541957600181146154515761547f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955061547f565b866000528260002060005b858110156154775781548a820186015290830190840161545c565b890184019650505b505050505092915050565b60a08152600061549d60a08301886153ef565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b8381101561550f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526154fd83836153ef565b948601949250600191820191016154c4565b50508681036040880152615523818b6153ef565b945050505050846060840152828103608084015261554181856145d3565b98975050505050505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f8211156129d757600081815260208120601f850160051c810160208610156155c85750805b601f850160051c820191505b818110156113e6578281556001016155d4565b815167ffffffffffffffff81111561560157615601614309565b6156158161560f8454615308565b846155a1565b602080601f83116001811461566857600084156156325750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556113e6565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156156b557888601518255948401946001909101908401615696565b50858210156156f157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000602080838503121561571457600080fd5b825167ffffffffffffffff81111561572b57600080fd5b8301601f8101851361573c57600080fd5b805161574a61441682614a2a565b81815260059190911b8201830190838101908783111561576957600080fd5b928401925b828410156157875783518252928401929084019061576e565b979650505050505050565b600063ffffffff80831681810361523e5761523e614e88565b80516020808301519190811015615355577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b602081526000825161014080602085015261580c6101608501836145d3565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08086850301604087015261584884836145d3565b935060408701519150615873606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526158d484836145d3565b935060e087015191506101008187860301818801526158f385846145d3565b94508088015192505061012081878603018188015261591285846145d3565b94508088015192505050615935828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561597757615977614e88565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826159ba576159ba61597c565b500490565b818103600083128015838313168383128216171561418857614188614e88565b6000826159ee576159ee61597c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615615a4257615a42614e88565b500590565b8082018281126000831280158216821582161715615a6757615a67614e88565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036150cb576150cb614e88565b60007f80000000000000000000000000000000000000000000000000000000000000008203615ad157615ad1614e88565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var VerifiableLoadStreamsLookupUpkeepABI = VerifiableLoadStreamsLookupUpkeepMetaData.ABI

var VerifiableLoadStreamsLookupUpkeepBin = VerifiableLoadStreamsLookupUpkeepMetaData.Bin

func DeployVerifiableLoadStreamsLookupUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _registrar common.Address, _useArb bool) (common.Address, *types.Transaction, *VerifiableLoadStreamsLookupUpkeep, error) {
	parsed, err := VerifiableLoadStreamsLookupUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadStreamsLookupUpkeepBin), backend, _registrar, _useArb)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifiableLoadStreamsLookupUpkeep{address: address, abi: *parsed, VerifiableLoadStreamsLookupUpkeepCaller: VerifiableLoadStreamsLookupUpkeepCaller{contract: contract}, VerifiableLoadStreamsLookupUpkeepTransactor: VerifiableLoadStreamsLookupUpkeepTransactor{contract: contract}, VerifiableLoadStreamsLookupUpkeepFilterer: VerifiableLoadStreamsLookupUpkeepFilterer{contract: contract}}, nil
}

type VerifiableLoadStreamsLookupUpkeep struct {
	address common.Address
	abi     abi.ABI
	VerifiableLoadStreamsLookupUpkeepCaller
	VerifiableLoadStreamsLookupUpkeepTransactor
	VerifiableLoadStreamsLookupUpkeepFilterer
}

type VerifiableLoadStreamsLookupUpkeepCaller struct {
	contract *bind.BoundContract
}

type VerifiableLoadStreamsLookupUpkeepTransactor struct {
	contract *bind.BoundContract
}

type VerifiableLoadStreamsLookupUpkeepFilterer struct {
	contract *bind.BoundContract
}

type VerifiableLoadStreamsLookupUpkeepSession struct {
	Contract     *VerifiableLoadStreamsLookupUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VerifiableLoadStreamsLookupUpkeepCallerSession struct {
	Contract *VerifiableLoadStreamsLookupUpkeepCaller
	CallOpts bind.CallOpts
}

type VerifiableLoadStreamsLookupUpkeepTransactorSession struct {
	Contract     *VerifiableLoadStreamsLookupUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type VerifiableLoadStreamsLookupUpkeepRaw struct {
	Contract *VerifiableLoadStreamsLookupUpkeep
}

type VerifiableLoadStreamsLookupUpkeepCallerRaw struct {
	Contract *VerifiableLoadStreamsLookupUpkeepCaller
}

type VerifiableLoadStreamsLookupUpkeepTransactorRaw struct {
	Contract *VerifiableLoadStreamsLookupUpkeepTransactor
}

func NewVerifiableLoadStreamsLookupUpkeep(address common.Address, backend bind.ContractBackend) (*VerifiableLoadStreamsLookupUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(VerifiableLoadStreamsLookupUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVerifiableLoadStreamsLookupUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeep{address: address, abi: abi, VerifiableLoadStreamsLookupUpkeepCaller: VerifiableLoadStreamsLookupUpkeepCaller{contract: contract}, VerifiableLoadStreamsLookupUpkeepTransactor: VerifiableLoadStreamsLookupUpkeepTransactor{contract: contract}, VerifiableLoadStreamsLookupUpkeepFilterer: VerifiableLoadStreamsLookupUpkeepFilterer{contract: contract}}, nil
}

func NewVerifiableLoadStreamsLookupUpkeepCaller(address common.Address, caller bind.ContractCaller) (*VerifiableLoadStreamsLookupUpkeepCaller, error) {
	contract, err := bindVerifiableLoadStreamsLookupUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepCaller{contract: contract}, nil
}

func NewVerifiableLoadStreamsLookupUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifiableLoadStreamsLookupUpkeepTransactor, error) {
	contract, err := bindVerifiableLoadStreamsLookupUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepTransactor{contract: contract}, nil
}

func NewVerifiableLoadStreamsLookupUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifiableLoadStreamsLookupUpkeepFilterer, error) {
	contract, err := bindVerifiableLoadStreamsLookupUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepFilterer{contract: contract}, nil
}

func bindVerifiableLoadStreamsLookupUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifiableLoadStreamsLookupUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.VerifiableLoadStreamsLookupUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.VerifiableLoadStreamsLookupUpkeepTransactor.contract.Transfer(opts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.VerifiableLoadStreamsLookupUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.contract.Transfer(opts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) BUCKETSIZE(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "BUCKET_SIZE")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) AddLinkAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "addLinkAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.AddLinkAmount(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.AddLinkAmount(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "bucketedDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BucketedDelays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BucketedDelays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "buckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Buckets(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Buckets(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckCallback(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckCallback(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) CheckErrorHandler(opts *bind.CallOpts, errCode uint32, extraData []byte) (CheckErrorHandler,

	error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "checkErrorHandler", errCode, extraData)

	outstruct := new(CheckErrorHandler)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) CheckErrorHandler(errCode uint32, extraData []byte) (CheckErrorHandler,

	error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckErrorHandler(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, errCode, extraData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) CheckErrorHandler(errCode uint32, extraData []byte) (CheckErrorHandler,

	error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckErrorHandler(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, errCode, extraData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "checkGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "counters", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Counters(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Counters(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "delays", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Delays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Delays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.DummyMap(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.DummyMap(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "eligible", upkeepId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Eligible(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Eligible(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) EmittedAgainSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "emittedAgainSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) EmittedSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "emittedSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.EmittedSig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.EmittedSig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) FeedParamKey() (string, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.FeedParamKey(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) FeedParamKey() (string, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.FeedParamKey(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.FeedsHex(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.FeedsHex(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "firstPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "gasLimits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GasLimits(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GasLimits(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetActiveUpkeepIDsDeployedByThisContract(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getActiveUpkeepIDsDeployedByThisContract", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetAllActiveUpkeepIDsOnRegistry(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getAllActiveUpkeepIDsOnRegistry", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetBalance(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, id)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetBalance(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, id)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getBucketedDelays", upkeepId, bucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getDelays", upkeepId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetDelays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetDelays(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetDelaysLength(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetDelaysLength(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetForwarder(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetForwarder(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetLogTriggerConfig(opts *bind.CallOpts, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getLogTriggerConfig", addr, selector, topic0, topic1, topic2, topic3)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getMinBalanceForUpkeep", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getPxDelayLastNPerforms", upkeepId, p, n)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getSumDelayInBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getSumDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetTriggerType(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetTriggerType(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetUpkeepInfo(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getUpkeepInfo", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21UpkeepInfo)).(*KeeperRegistryBase21UpkeepInfo)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Intervals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "intervals", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Intervals(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Intervals(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) LastTopUpBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "lastTopUpBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.LinkToken(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.LinkToken(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) MinBalanceThresholdMultiplier(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "minBalanceThresholdMultiplier")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Owner() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Owner(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Owner() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Owner(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) PerformDataSizes(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "performDataSizes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PerformDataSizes(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PerformDataSizes(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) PerformGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "performGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) PreviousPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "previousPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Registrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "registrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Registrar() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Registrar(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Registrar() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Registrar(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Registry() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Registry(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) Registry() (common.Address, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Registry(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) TimeParamKey() (string, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.TimeParamKey(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) TimeParamKey() (string, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.TimeParamKey(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) UpkeepTopUpCheckInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "upkeepTopUpCheckInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "useArbitrumBlockNum")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadStreamsLookupUpkeep.CallOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "acceptOwnership")
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.AcceptOwnership(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.AcceptOwnership(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) AddFunds(opts *bind.TransactOpts, upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "addFunds", upkeepId, amount)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.AddFunds(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.AddFunds(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchCancelUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchCancelUpkeeps", upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchPreparingUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchPreparingUpkeeps", upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchPreparingUpkeepsSimple(opts *bind.TransactOpts, upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchPreparingUpkeepsSimple", upkeepIds, log, selector)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds, log, selector)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds, log, selector)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchRegisterUpkeeps", number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchSendLogs(opts *bind.TransactOpts, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchSendLogs", log)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchSendLogs(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, log)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchSendLogs(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, log)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchSetIntervals", upkeepIds, interval)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchUpdatePipelineData", upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "batchWithdrawLinks", upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) BurnPerformGas(opts *bind.TransactOpts, upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "burnPerformGas", upkeepId, startGas, blockNum)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BurnPerformGas(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.BurnPerformGas(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) CheckUpkeep(opts *bind.TransactOpts, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "checkUpkeep", checkData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckUpkeep(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, checkData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) CheckUpkeep(checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckUpkeep(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, checkData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PerformUpkeep(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.PerformUpkeep(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SendLog(opts *bind.TransactOpts, upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "sendLog", upkeepId, log)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SendLog(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, log)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SendLog(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, log)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "setConfig", newRegistrar)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetConfig(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetConfig(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "setFeeds", _feeds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetFeeds(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, _feeds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetFeeds(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, _feeds)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "setInterval", upkeepId, _interval)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetInterval(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetInterval(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "setParamKeys", _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetParamKeys(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetParamKeys(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "setPerformDataSize", upkeepId, value)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "setUpkeepGasLimit", upkeepId, gasLimit)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, cfg)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) TopUpFund(opts *bind.TransactOpts, upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "topUpFund", upkeepId, blockNum)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.TopUpFund(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.TopUpFund(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "transferOwnership", to)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.TransferOwnership(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, to)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.TransferOwnership(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, to)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) UpdateLogTriggerConfig1(opts *bind.TransactOpts, upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "updateLogTriggerConfig1", upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) UpdateLogTriggerConfig2(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "updateLogTriggerConfig2", upkeepId, cfg)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "updateUpkeepPipelineData", upkeepId, pipelineData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "withdrawLinks")
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.WithdrawLinks(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.WithdrawLinks(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.Transact(opts, "withdrawLinks0", upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.contract.RawTransact(opts, nil)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Receive(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepTransactorSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.Receive(&_VerifiableLoadStreamsLookupUpkeep.TransactOpts)
}

type VerifiableLoadStreamsLookupUpkeepLogEmittedIterator struct {
	Event *VerifiableLoadStreamsLookupUpkeepLogEmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadStreamsLookupUpkeepLogEmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadStreamsLookupUpkeepLogEmitted)
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
		it.Event = new(VerifiableLoadStreamsLookupUpkeepLogEmitted)
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

func (it *VerifiableLoadStreamsLookupUpkeepLogEmittedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadStreamsLookupUpkeepLogEmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadStreamsLookupUpkeepLogEmitted struct {
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadStreamsLookupUpkeepLogEmittedIterator, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumRule []interface{}
	for _, blockNumItem := range blockNum {
		blockNumRule = append(blockNumRule, blockNumItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.FilterLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepLogEmittedIterator{contract: _VerifiableLoadStreamsLookupUpkeep.contract, event: "LogEmitted", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumRule []interface{}
	for _, blockNumItem := range blockNum {
		blockNumRule = append(blockNumRule, blockNumItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.WatchLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadStreamsLookupUpkeepLogEmitted)
				if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) ParseLogEmitted(log types.Log) (*VerifiableLoadStreamsLookupUpkeepLogEmitted, error) {
	event := new(VerifiableLoadStreamsLookupUpkeepLogEmitted)
	if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadStreamsLookupUpkeepLogEmittedAgainIterator struct {
	Event *VerifiableLoadStreamsLookupUpkeepLogEmittedAgain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadStreamsLookupUpkeepLogEmittedAgainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadStreamsLookupUpkeepLogEmittedAgain)
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
		it.Event = new(VerifiableLoadStreamsLookupUpkeepLogEmittedAgain)
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

func (it *VerifiableLoadStreamsLookupUpkeepLogEmittedAgainIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadStreamsLookupUpkeepLogEmittedAgainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadStreamsLookupUpkeepLogEmittedAgain struct {
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadStreamsLookupUpkeepLogEmittedAgainIterator, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumRule []interface{}
	for _, blockNumItem := range blockNum {
		blockNumRule = append(blockNumRule, blockNumItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.FilterLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepLogEmittedAgainIterator{contract: _VerifiableLoadStreamsLookupUpkeep.contract, event: "LogEmittedAgain", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

	var upkeepIdRule []interface{}
	for _, upkeepIdItem := range upkeepId {
		upkeepIdRule = append(upkeepIdRule, upkeepIdItem)
	}
	var blockNumRule []interface{}
	for _, blockNumItem := range blockNum {
		blockNumRule = append(blockNumRule, blockNumItem)
	}
	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.WatchLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadStreamsLookupUpkeepLogEmittedAgain)
				if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) ParseLogEmittedAgain(log types.Log) (*VerifiableLoadStreamsLookupUpkeepLogEmittedAgain, error) {
	event := new(VerifiableLoadStreamsLookupUpkeepLogEmittedAgain)
	if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequestedIterator struct {
	Event *VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested)
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
		it.Event = new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested)
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

func (it *VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequestedIterator{contract: _VerifiableLoadStreamsLookupUpkeep.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested)
				if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested, error) {
	event := new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested)
	if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadStreamsLookupUpkeepOwnershipTransferredIterator struct {
	Event *VerifiableLoadStreamsLookupUpkeepOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadStreamsLookupUpkeepOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferred)
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
		it.Event = new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferred)
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

func (it *VerifiableLoadStreamsLookupUpkeepOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadStreamsLookupUpkeepOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadStreamsLookupUpkeepOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepOwnershipTransferredIterator{contract: _VerifiableLoadStreamsLookupUpkeep.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferred)
				if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) ParseOwnershipTransferred(log types.Log) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferred, error) {
	event := new(VerifiableLoadStreamsLookupUpkeepOwnershipTransferred)
	if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadStreamsLookupUpkeepUpkeepTopUpIterator struct {
	Event *VerifiableLoadStreamsLookupUpkeepUpkeepTopUp

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadStreamsLookupUpkeepUpkeepTopUpIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadStreamsLookupUpkeepUpkeepTopUp)
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
		it.Event = new(VerifiableLoadStreamsLookupUpkeepUpkeepTopUp)
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

func (it *VerifiableLoadStreamsLookupUpkeepUpkeepTopUpIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadStreamsLookupUpkeepUpkeepTopUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadStreamsLookupUpkeepUpkeepTopUp struct {
	UpkeepId *big.Int
	Amount   *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadStreamsLookupUpkeepUpkeepTopUpIterator, error) {

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.FilterLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadStreamsLookupUpkeepUpkeepTopUpIterator{contract: _VerifiableLoadStreamsLookupUpkeep.contract, event: "UpkeepTopUp", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepUpkeepTopUp) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadStreamsLookupUpkeep.contract.WatchLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadStreamsLookupUpkeepUpkeepTopUp)
				if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepFilterer) ParseUpkeepTopUp(log types.Log) (*VerifiableLoadStreamsLookupUpkeepUpkeepTopUp, error) {
	event := new(VerifiableLoadStreamsLookupUpkeepUpkeepTopUp)
	if err := _VerifiableLoadStreamsLookupUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type CheckErrorHandler struct {
	UpkeepNeeded bool
	PerformData  []byte
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifiableLoadStreamsLookupUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadStreamsLookupUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadStreamsLookupUpkeep.abi.Events["LogEmittedAgain"].ID:
		return _VerifiableLoadStreamsLookupUpkeep.ParseLogEmittedAgain(log)
	case _VerifiableLoadStreamsLookupUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadStreamsLookupUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadStreamsLookupUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadStreamsLookupUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadStreamsLookupUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadStreamsLookupUpkeep.ParseUpkeepTopUp(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadStreamsLookupUpkeepLogEmitted) Topic() common.Hash {
	return common.HexToHash("0x97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf08")
}

func (VerifiableLoadStreamsLookupUpkeepLogEmittedAgain) Topic() common.Hash {
	return common.HexToHash("0xc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d")
}

func (VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifiableLoadStreamsLookupUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifiableLoadStreamsLookupUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeep) Address() common.Address {
	return _VerifiableLoadStreamsLookupUpkeep.address
}

type VerifiableLoadStreamsLookupUpkeepInterface interface {
	BUCKETSIZE(opts *bind.CallOpts) (uint16, error)

	AddLinkAmount(opts *bind.CallOpts) (*big.Int, error)

	BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error)

	Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error)

	CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error)

	CheckErrorHandler(opts *bind.CallOpts, errCode uint32, extraData []byte) (CheckErrorHandler,

		error)

	CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)

	DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error)

	EmittedAgainSig(opts *bind.CallOpts) ([32]byte, error)

	EmittedSig(opts *bind.CallOpts) ([32]byte, error)

	FeedParamKey(opts *bind.CallOpts) (string, error)

	FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error)

	FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error)

	GetActiveUpkeepIDsDeployedByThisContract(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetAllActiveUpkeepIDsOnRegistry(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error)

	GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error)

	GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error)

	GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error)

	GetLogTriggerConfig(opts *bind.CallOpts, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error)

	GetMinBalanceForUpkeep(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)

	GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error)

	GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error)

	GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error)

	GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error)

	GetUpkeepInfo(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error)

	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)

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

	BatchPreparingUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error)

	BatchPreparingUpkeepsSimple(opts *bind.TransactOpts, upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error)

	BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error)

	BatchSendLogs(opts *bind.TransactOpts, log uint8) (*types.Transaction, error)

	BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error)

	BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error)

	BurnPerformGas(opts *bind.TransactOpts, upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error)

	CheckUpkeep(opts *bind.TransactOpts, checkData []byte) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SendLog(opts *bind.TransactOpts, upkeepId *big.Int, log uint8) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error)

	SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error)

	SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error)

	SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error)

	SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error)

	SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error)

	TopUpFund(opts *bind.TransactOpts, upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateLogTriggerConfig1(opts *bind.TransactOpts, upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error)

	UpdateLogTriggerConfig2(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error)

	UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error)

	WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error)

	WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadStreamsLookupUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadStreamsLookupUpkeepLogEmitted, error)

	FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadStreamsLookupUpkeepLogEmittedAgainIterator, error)

	WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmittedAgain(log types.Log) (*VerifiableLoadStreamsLookupUpkeepLogEmittedAgain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadStreamsLookupUpkeepOwnershipTransferred, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadStreamsLookupUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadStreamsLookupUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadStreamsLookupUpkeepUpkeepTopUp, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
