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

type IAutomationV21PlusCommonUpkeepInfoLegacy struct {
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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkErrorHandler\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structIAutomationV21PlusCommon.UpkeepInfoLegacy\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c06040526042610140818152610100918291906200622361016039815260200160405180608001604052806042815260200162006265604291399052620000be906016906002620003c7565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000ee908262000543565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b602082015260189062000120908262000543565b503480156200012e57600080fd5b50604051620062a7380380620062a7833981016040819052620001519162000625565b81813380600081620001aa5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001dd57620001dd816200031c565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa1580156200023a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000260919062000668565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002c6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ec919062000699565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c05250620006c0915050565b336001600160a01b03821603620003765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620001a1565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821562000412579160200282015b8281111562000412578251829062000401908262000543565b5091602001919060010190620003e8565b506200042092915062000424565b5090565b80821115620004205760006200043b828262000445565b5060010162000424565b5080546200045390620004b4565b6000825580601f1062000464575050565b601f01602090049060005260206000209081019062000484919062000487565b50565b5b8082111562000420576000815560010162000488565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004c957607f821691505b602082108103620004ea57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200053e57600081815260208120601f850160051c81016020861015620005195750805b601f850160051c820191505b818110156200053a5782815560010162000525565b5050505b505050565b81516001600160401b038111156200055f576200055f6200049e565b6200057781620005708454620004b4565b84620004f0565b602080601f831160018114620005af5760008415620005965750858301515b600019600386901b1c1916600185901b1785556200053a565b600085815260208120601f198616915b82811015620005e057888601518255948401946001909101908401620005bf565b5085821015620005ff5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200048457600080fd5b600080604083850312156200063957600080fd5b825162000646816200060f565b602084015190925080151581146200065d57600080fd5b809150509250929050565b600080604083850312156200067c57600080fd5b825162000689816200060f565b6020939093015192949293505050565b600060208284031215620006ac57600080fd5b8151620006b9816200060f565b9392505050565b60805160a05160c05160e051615b0d62000716600039600081816105b10152611fcc0152600081816109f70152613cf90152600081816108700152613747015260008181610db4015261371c0152615b0d6000f3fe6080604052600436106104f05760003560e01c806379ba509711610294578063a6b594751161015e578063d6051a72116100d6578063e45530831161008a578063fa333dfb1161006f578063fa333dfb14610fc3578063fba7ffa314611076578063fcdc1f63146110a357600080fd5b8063e455308314610f8d578063f2fde38b14610fa357600080fd5b8063daee1aeb116100bb578063daee1aeb14610f20578063dbef701e14610f40578063e0114adb14610f6057600080fd5b8063d6051a7214610ee0578063da6cba4714610f0057600080fd5b8063b657bc9c1161012d578063c041982211610112578063c041982214610e8b578063c98f10b014610eab578063d4c2490014610ec057600080fd5b8063b657bc9c14610e4b578063becde0e114610e6b57600080fd5b8063a6b5947514610dd6578063a72aa27e14610df6578063af953a4a14610e16578063afb28d1f14610e3657600080fd5b80638fcb3fba1161020c5780639b429354116101c05780639d385eaa116101a55780639d385eaa14610d625780639d6f1cc714610d82578063a654824814610da257600080fd5b80639b42935414610d045780639b51fb0d14610d3157600080fd5b8063948108f7116101f1578063948108f714610c9a57806396cebc7c14610cba5780639ac542eb14610cda57600080fd5b80638fcb3fba14610c4d578063924ca57814610c7a57600080fd5b80638243444a1161026357806386e330af1161024857806386e330af14610be2578063873c758614610c025780638da5cb5b14610c2257600080fd5b80638243444a14610ba25780638340507c14610bc257600080fd5b806379ba509714610b2057806379ea994314610b355780637b10399914610b555780637e7a46dc14610b8257600080fd5b80634585e33b116103d55780635f17e6161161034d5780636e04ff0d1161030157806373644cce116102e657806373644cce14610aa65780637672130314610ad3578063776898c814610b0057600080fd5b80636e04ff0d14610a565780637145f11b14610a7657600080fd5b8063636092e811610332578063636092e8146109c0578063642f6cef146109e557806369cdbadb14610a2957600080fd5b80635f17e6161461097357806360457ff51461099357600080fd5b80634b56a42e116103a457806351c98be31161038957806351c98be31461091157806357970e93146109315780635d4ee7f31461095e57600080fd5b80634b56a42e146108bf5780635147cd59146108df57600080fd5b80634585e33b1461081157806345d2ec1714610831578063469820931461085e57806346e7a63e1461089257600080fd5b8063207b65161161046857806329e0a841116104375780632b20e3971161041c5780632b20e39714610772578063328ffd11146107c45780633ebe8d6c146107f157600080fd5b806329e0a841146107255780632a9032d31461075257600080fd5b8063207b6516146106a557806320e3dbd4146106c55780632636aecf146106e557806328c4b57b1461070557600080fd5b806312c55027116104bf5780631cdde251116104a45780631cdde251146106135780631e01043914610633578063206c32e81461067057600080fd5b806312c550271461059f57806319d97a94146105e657600080fd5b806306c1cc00146104fc578063077ac6211461051e5780630b7d33e6146105515780630fb172fb1461057157600080fd5b366104f757005b600080fd5b34801561050857600080fd5b5061051c61051736600461447b565b6110d0565b005b34801561052a57600080fd5b5061053e61053936600461452e565b61131f565b6040519081526020015b60405180910390f35b34801561055d57600080fd5b5061051c61056c366004614563565b61135d565b34801561057d57600080fd5b5061059161058c366004614563565b6113eb565b604051610548929190614618565b3480156105ab57600080fd5b506105d37f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610548565b3480156105f257600080fd5b50610606610601366004614633565b611403565b604051610548919061464c565b34801561061f57600080fd5b5061051c61062e366004614681565b6114c0565b34801561063f57600080fd5b5061065361064e366004614633565b6115fd565b6040516bffffffffffffffffffffffff9091168152602001610548565b34801561067c57600080fd5b5061069061068b3660046146e6565b611692565b60408051928352602083019190915201610548565b3480156106b157600080fd5b506106066106c0366004614633565b611714565b3480156106d157600080fd5b5061051c6106e0366004614712565b61176c565b3480156106f157600080fd5b5061051c610700366004614774565b611936565b34801561071157600080fd5b5061053e6107203660046147ee565b611bff565b34801561073157600080fd5b50610745610740366004614633565b611c6a565b604051610548919061481a565b34801561075e57600080fd5b5061051c61076d36600461495b565b611d6f565b34801561077e57600080fd5b5060115461079f9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610548565b3480156107d057600080fd5b5061053e6107df366004614633565b60036020526000908152604090205481565b3480156107fd57600080fd5b5061053e61080c366004614633565b611e50565b34801561081d57600080fd5b5061051c61082c3660046149df565b611eb9565b34801561083d57600080fd5b5061085161084c3660046146e6565b6120d8565b6040516105489190614a15565b34801561086a57600080fd5b5061053e7f000000000000000000000000000000000000000000000000000000000000000081565b34801561089e57600080fd5b5061053e6108ad366004614633565b600a6020526000908152604090205481565b3480156108cb57600080fd5b506105916108da366004614a7d565b612147565b3480156108eb57600080fd5b506108ff6108fa366004614633565b61219b565b60405160ff9091168152602001610548565b34801561091d57600080fd5b5061051c61092c366004614b47565b61222f565b34801561093d57600080fd5b5060125461079f9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561096a57600080fd5b5061051c6122d3565b34801561097f57600080fd5b5061051c61098e366004614b9e565b61240e565b34801561099f57600080fd5b5061053e6109ae366004614633565b60076020526000908152604090205481565b3480156109cc57600080fd5b50601554610653906bffffffffffffffffffffffff1681565b3480156109f157600080fd5b50610a197f000000000000000000000000000000000000000000000000000000000000000081565b6040519015158152602001610548565b348015610a3557600080fd5b5061053e610a44366004614633565b60086020526000908152604090205481565b348015610a6257600080fd5b50610591610a713660046149df565b6124db565b348015610a8257600080fd5b50610a19610a91366004614633565b600b6020526000908152604090205460ff1681565b348015610ab257600080fd5b5061053e610ac1366004614633565b6000908152600c602052604090205490565b348015610adf57600080fd5b5061053e610aee366004614633565b60046020526000908152604090205481565b348015610b0c57600080fd5b50610a19610b1b366004614633565b612704565b348015610b2c57600080fd5b5061051c612756565b348015610b4157600080fd5b5061079f610b50366004614633565b612853565b348015610b6157600080fd5b5060135461079f9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b8e57600080fd5b5061051c610b9d366004614bc0565b6128e7565b348015610bae57600080fd5b5061051c610bbd366004614bc0565b612978565b348015610bce57600080fd5b5061051c610bdd366004614c0c565b6129d2565b348015610bee57600080fd5b5061051c610bfd366004614c59565b6129f0565b348015610c0e57600080fd5b50610851610c1d366004614b9e565b612a03565b348015610c2e57600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661079f565b348015610c5957600080fd5b5061053e610c68366004614633565b60056020526000908152604090205481565b348015610c8657600080fd5b5061051c610c95366004614b9e565b612ac0565b348015610ca657600080fd5b5061051c610cb5366004614d0a565b612d05565b348015610cc657600080fd5b5061051c610cd5366004614d3a565b612e1d565b348015610ce657600080fd5b506015546108ff906c01000000000000000000000000900460ff1681565b348015610d1057600080fd5b5061051c610d1f366004614b9e565b60009182526009602052604090912055565b348015610d3d57600080fd5b506105d3610d4c366004614633565b600e6020526000908152604090205461ffff1681565b348015610d6e57600080fd5b50610851610d7d366004614633565b613027565b348015610d8e57600080fd5b50610606610d9d366004614633565b613089565b348015610dae57600080fd5b5061053e7f000000000000000000000000000000000000000000000000000000000000000081565b348015610de257600080fd5b5061051c610df13660046147ee565b613135565b348015610e0257600080fd5b5061051c610e11366004614d57565b61319e565b348015610e2257600080fd5b5061051c610e31366004614633565b613249565b348015610e4257600080fd5b506106066132cf565b348015610e5757600080fd5b50610653610e66366004614633565b6132dc565b348015610e7757600080fd5b5061051c610e8636600461495b565b613334565b348015610e9757600080fd5b50610851610ea6366004614b9e565b6133ce565b348015610eb757600080fd5b506106066134cb565b348015610ecc57600080fd5b5061051c610edb366004614d7c565b6134d8565b348015610eec57600080fd5b50610690610efb366004614b9e565b613557565b348015610f0c57600080fd5b5061051c610f1b366004614da1565b6135c0565b348015610f2c57600080fd5b5061051c610f3b36600461495b565b613927565b348015610f4c57600080fd5b5061053e610f5b366004614b9e565b6139f2565b348015610f6c57600080fd5b5061053e610f7b366004614633565b60096020526000908152604090205481565b348015610f9957600080fd5b5061053e60145481565b348015610faf57600080fd5b5061051c610fbe366004614712565b613a23565b348015610fcf57600080fd5b50610606610fde366004614e09565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b34801561108257600080fd5b5061053e611091366004614633565b60066020526000908152604090205481565b3480156110af57600080fd5b5061053e6110be366004614633565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b39216906111b6908c1688614e91565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af1158015611234573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112589190614ed5565b5060008860ff1667ffffffffffffffff8111156112775761127761431d565b6040519080825280602002602001820160405280156112a0578160200160208202803683370190505b50905060005b8960ff168160ff1610156113135760006112bf84613a37565b905080838360ff16815181106112d7576112d7614ef0565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061130b81614f1f565b9150506112a6565b50505050505050505050565b600d602052826000526040600020602052816000526040600020818154811061134757600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e6906113b59085908590600401614f3e565b600060405180830381600087803b1580156113cf57600080fd5b505af11580156113e3573d6000803e3d6000fd5b505050505050565b604080516000808252602082019092525b9250929050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa158015611474573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114ba9190810190614fa4565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa15801561155f573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115a59190810190614fa4565b6040518363ffffffff1660e01b81526004016115c2929190614f3e565b600060405180830381600087803b1580156115dc57600080fd5b505af11580156115f0573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e010439906024015b602060405180830381865afa15801561166e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ba9190614fe4565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156116f657602002820191906000526020600020905b8154815260200190600101908083116116e2575b50505050509050611708818251613b05565b92509250509250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b651690602401611457565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611802573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611826919061500c565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa1580156118c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118ed919061503a565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611bf457600089898381811061195657611956614ef0565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161198f91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016119bb929190614f3e565b600060405180830381600087803b1580156119d557600080fd5b505af11580156119e9573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015611a5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a839190615057565b90508060ff16600103611bdf576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611b0c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b529190810190614fa4565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611bab9086908590600401614f3e565b600060405180830381600087803b158015611bc557600080fd5b505af1158015611bd9573d6000803e3d6000fd5b50505050505b50508080611bec90615074565b91505061193a565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611c6093830182828015611c5457602002820191906000526020600020905b815481526020019060010190808311611c40575b50505050508484613b8a565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611d29573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114ba91908101906150cf565b8060005b818160ff161015611e4a5760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611db157611db1614ef0565b905060200201356040518263ffffffff1660e01b8152600401611dd691815260200190565b600060405180830381600087803b158015611df057600080fd5b505af1158015611e04573d6000803e3d6000fd5b50505050611e3784848360ff16818110611e2057611e20614ef0565b90506020020135600f613ce990919063ffffffff16565b5080611e4281614f1f565b915050611d73565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611eb1576000858152600d6020908152604080832061ffff85168452909152902054611e9d90836151ee565b915080611ea981615201565b915050611e66565b509392505050565b60005a9050600080611ecd84860186614a7d565b91509150600081806020019051810190611ee79190615222565b60008181526005602090815260408083205460049092528220549293509190611f0e613cf5565b905082600003611f2e576000848152600560205260409020819055612089565b600084815260036020526040812054611f47848461523b565b611f51919061523b565b6000868152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611fc357602002820191906000526020600020905b815481526020019060010190808311611faf575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff1681510361203e578161200081615201565b6000898152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000868152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055888352600c8252928220805493840181558252902001555b6000848152600660205260408120546120a39060016151ee565b60008681526006602090815260408083208490556004909152902083905590506120cd8583612ac0565b611313858984613135565b6000828152600d6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561213a57602002820191906000526020600020905b815481526020019060010190808311612126575b5050505050905092915050565b600060606000848460405160200161216092919061524e565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa15801561220b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ba9190615057565b8160005b818110156122cc5730635f17e61686868481811061225357612253614ef0565b90506020020135856040518363ffffffff1660e01b815260040161228792919091825263ffffffff16602082015260400190565b600060405180830381600087803b1580156122a157600080fd5b505af11580156122b5573d6000803e3d6000fd5b5050505080806122c490615074565b915050612233565b5050505050565b6122db613d97565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561234a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061236e9190615222565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156123e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061240a9190614ed5565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c90915281206124469161421c565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff16116124a2576000848152600d6020908152604080832061ffff8516845290915281206124909161421c565b8061249a81615201565b91505061245b565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000606060005a905060006124f285870187614633565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff81111561252b5761252b61431d565b6040519080825280601f01601f191660200182016040528015612555576020820181803683370190505b50604051602001612567929190614f3e565b60405160208183030381529060405290506000612582613cf5565b9050600061258f86612704565b90505b835a61259e908961523b565b6125aa906127106151ee565b10156125eb5781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612592565b806126035760008398509850505050505050506113fc565b6040517f6665656449644865780000000000000000000000000000000000000000000000602082015260009060290160405160208183030381529060405280519060200120601760405160200161265a9190615335565b604051602081830303815290604052805190602001200361267c57508161267f565b50425b601760166018838a60405160200161269991815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527ff055e4a20000000000000000000000000000000000000000000000000000000082526126fb9594939291600401615464565b60405180910390fd5b600081815260056020526040812054810361272157506001919050565b600082815260036020908152604080832054600490925290912054612744613cf5565b61274e919061523b565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146127d7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016126fb565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa1580156128c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ba919061503a565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b59061294190869086908690600401615527565b600060405180830381600087803b15801561295b57600080fd5b505af115801561296f573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d359061294190869086908690600401615527565b60176129de83826155c1565b5060186129eb82826155c1565b505050565b805161240a90601690602084019061423a565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612a7a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611c6391908101906156db565b601454600083815260026020526040902054612adc908361523b565b111561240a576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612b52573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612b9891908101906150cf565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612c0d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c319190614fe4565b601554909150612c559082906c01000000000000000000000000900460ff16614e91565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611e4a57601554612c989085906bffffffffffffffffffffffff16612d05565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612d8d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612db19190614ed5565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f7906044016113b5565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015612e7c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612ec291908101906156db565b80519091506000612ed1613cf5565b905060005b828110156122cc576000848281518110612ef257612ef2614ef0565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612f72573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f969190615057565b90508060ff16600103613012578660ff16600003612fe2576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4613012565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b5050808061301f90615074565b915050612ed6565b6000818152600c602090815260409182902080548351818402810184019094528084526060939283018282801561307d57602002820191906000526020600020905b815481526020019060010190808311613069575b50505050509050919050565b6016818154811061309957600080fd5b9060005260206000200160009150905080546130b4906152e2565b80601f01602080910402602001604051908101604052809291908181526020018280546130e0906152e2565b801561312d5780601f106131025761010080835404028352916020019161312d565b820191906000526020600020905b81548152906001019060200180831161311057829003601f168201915b505050505081565b6000838152600760205260409020545b805a613151908561523b565b61315d906127106151ee565b1015611e4a5781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055613145565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b15801561321657600080fd5b505af115801561322a573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b1580156132bb57600080fd5b505af11580156122cc573d6000803e3d6000fd5b601780546130b4906152e2565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff169063b657bc9c90602401611651565b8060005b818163ffffffff161015611e4a573063af953a4a858563ffffffff851681811061336457613364614ef0565b905060200201356040518263ffffffff1660e01b815260040161338991815260200190565b600060405180830381600087803b1580156133a357600080fd5b505af11580156133b7573d6000803e3d6000fd5b5050505080806133c69061576c565b915050613338565b606060006133dc600f613e1a565b9050808410613417576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8260000361342c57613429848261523b565b92505b60008367ffffffffffffffff8111156134475761344761431d565b604051908082528060200260200182016040528015613470578160200160208202803683370190505b50905060005b848110156134c25761349361348b82886151ee565b600f90613e24565b8282815181106134a5576134a5614ef0565b6020908102919091010152806134ba81615074565b915050613476565b50949350505050565b601880546130b4906152e2565b60006134e2613cf5565b90508160ff16600003613523576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a4505050565b6000828152600c602090815260408083208054825181850281018501909352808352849384939291908301828280156135af57602002820191906000526020600020905b81548152602001906001019080831161359b575b505050505090506117088185613b05565b8260005b818110156113e35760008686838181106135e0576135e0614ef0565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161361991815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613645929190614f3e565b600060405180830381600087803b15801561365f57600080fd5b505af1158015613673573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa1580156136e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061370d9190615057565b90508060ff16600103613912577f000000000000000000000000000000000000000000000000000000000000000060ff87161561376757507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb3089858860405160200161379b91815260200190565b6040516020818303038152906040526137b390615785565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa15801561383e573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526138849190810190614fa4565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906138dd9087908590600401614f3e565b600060405180830381600087803b1580156138f757600080fd5b505af115801561390b573d6000803e3d6000fd5b5050505050505b5050808061391f90615074565b9150506135c4565b8060005b81811015611e4a57600084848381811061394757613947614ef0565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161398091815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016139ac929190614f3e565b600060405180830381600087803b1580156139c657600080fd5b505af11580156139da573d6000803e3d6000fd5b505050505080806139ea90615074565b91505061392b565b600c6020528160005260406000208181548110613a0e57600080fd5b90600052602060002001600091509150505481565b613a2b613d97565b613a3481613e30565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613a929086906004016157c7565b6020604051808303816000875af1158015613ab1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ad59190615222565b9050613ae2600f82613f25565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b815160009081908190841580613b1b5750808510155b15613b24578094505b60008092505b85831015613b8057866001613b3f858561523b565b613b49919061523b565b81518110613b5957613b59614ef0565b602002602001015181613b6c91906151ee565b905082613b7881615074565b935050613b2a565b9694955050505050565b82516000908190831580613b9e5750808410155b15613ba7578093505b60008467ffffffffffffffff811115613bc257613bc261431d565b604051908082528060200260200182016040528015613beb578160200160208202803683370190505b509050600092505b84831015613c5957866001613c08858561523b565b613c12919061523b565b81518110613c2257613c22614ef0565b6020026020010151818481518110613c3c57613c3c614ef0565b602090810291909101015282613c5181615074565b935050613bf3565b613c7281600060018451613c6d919061523b565b613f31565b85606403613cab578060018251613c89919061523b565b81518110613c9957613c99614ef0565b60200260200101519350505050611c63565b806064825188613cbb9190615919565b613cc59190615985565b81518110613cd557613cd5614ef0565b602002602001015193505050509392505050565b6000611c6383836140a9565b60007f000000000000000000000000000000000000000000000000000000000000000015613d9257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613d69573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d8d9190615222565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613e18576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016126fb565b565b60006114ba825490565b6000611c6383836141a3565b3373ffffffffffffffffffffffffffffffffffffffff821603613eaf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016126fb565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611c6383836141cd565b8181808203613f41575050505050565b6000856002613f508787615999565b613f5a91906159b9565b613f649087615a21565b81518110613f7457613f74614ef0565b602002602001015190505b818313614083575b80868481518110613f9a57613f9a614ef0565b60200260200101511015613fba5782613fb281615a49565b935050613f87565b858281518110613fcc57613fcc614ef0565b6020026020010151811015613fed5781613fe581615a7a565b925050613fba565b81831361407e5785828151811061400657614006614ef0565b602002602001015186848151811061402057614020614ef0565b602002602001015187858151811061403a5761403a614ef0565b6020026020010188858151811061405357614053614ef0565b6020908102919091010191909152528261406c81615a49565b935050818061407a90615a7a565b9250505b613f7f565b8185121561409657614096868684613f31565b838312156113e3576113e3868486613f31565b600081815260018301602052604081205480156141925760006140cd60018361523b565b85549091506000906140e19060019061523b565b905081811461414657600086600001828154811061410157614101614ef0565b906000526020600020015490508087600001848154811061412457614124614ef0565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061415757614157615ad1565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506114ba565b60009150506114ba565b5092915050565b60008260000182815481106141ba576141ba614ef0565b9060005260206000200154905092915050565b6000818152600183016020526040812054614214575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556114ba565b5060006114ba565b5080546000825590600052602060002090810190613a349190614290565b828054828255906000526020600020908101928215614280579160200282015b82811115614280578251829061427090826155c1565b509160200191906001019061425a565b5061428c9291506142a5565b5090565b5b8082111561428c5760008155600101614291565b8082111561428c5760006142b982826142c2565b506001016142a5565b5080546142ce906152e2565b6000825580601f106142de575050565b601f016020900490600052602060002090810190613a349190614290565b60ff81168114613a3457600080fd5b63ffffffff81168114613a3457600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff811182821017156143705761437061431d565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156143bd576143bd61431d565b604052919050565b600067ffffffffffffffff8211156143df576143df61431d565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261441c57600080fd5b813561442f61442a826143c5565b614376565b81815284602083860101111561444457600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff81168114613a3457600080fd5b600080600080600080600060e0888a03121561449657600080fd5b87356144a1816142fc565b965060208801356144b18161430b565b955060408801356144c1816142fc565b9450606088013567ffffffffffffffff8111156144dd57600080fd5b6144e98a828b0161440b565b94505060808801356144fa81614461565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff8116811461452957600080fd5b919050565b60008060006060848603121561454357600080fd5b8335925061455360208501614517565b9150604084013590509250925092565b6000806040838503121561457657600080fd5b82359150602083013567ffffffffffffffff81111561459457600080fd5b6145a08582860161440b565b9150509250929050565b60005b838110156145c55781810151838201526020016145ad565b50506000910152565b600081518084526145e68160208601602086016145aa565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b8215158152604060208201526000611c6060408301846145ce565b60006020828403121561464557600080fd5b5035919050565b602081526000611c6360208301846145ce565b73ffffffffffffffffffffffffffffffffffffffff81168114613a3457600080fd5b600080600080600080600060e0888a03121561469c57600080fd5b8735965060208801356146ae8161465f565b955060408801356146be816142fc565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b600080604083850312156146f957600080fd5b8235915061470960208401614517565b90509250929050565b60006020828403121561472457600080fd5b8135611c638161465f565b60008083601f84011261474157600080fd5b50813567ffffffffffffffff81111561475957600080fd5b6020830191508360208260051b85010111156113fc57600080fd5b600080600080600080600060c0888a03121561478f57600080fd5b873567ffffffffffffffff8111156147a657600080fd5b6147b28a828b0161472f565b90985096505060208801356147c6816142fc565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b60008060006060848603121561480357600080fd5b505081359360208301359350604090920135919050565b6020815261484160208201835173ffffffffffffffffffffffffffffffffffffffff169052565b6000602083015161485a604084018263ffffffff169052565b5060408301516101408060608501526148776101608501836145ce565b9150606085015161489860808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614904818701836bffffffffffffffffffffffff169052565b86015190506101206149198682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00183870152905061495183826145ce565b9695505050505050565b6000806020838503121561496e57600080fd5b823567ffffffffffffffff81111561498557600080fd5b6149918582860161472f565b90969095509350505050565b60008083601f8401126149af57600080fd5b50813567ffffffffffffffff8111156149c757600080fd5b6020830191508360208285010111156113fc57600080fd5b600080602083850312156149f257600080fd5b823567ffffffffffffffff811115614a0957600080fd5b6149918582860161499d565b6020808252825182820181905260009190848201906040850190845b81811015614a4d57835183529284019291840191600101614a31565b50909695505050505050565b600067ffffffffffffffff821115614a7357614a7361431d565b5060051b60200190565b60008060408385031215614a9057600080fd5b823567ffffffffffffffff80821115614aa857600080fd5b818501915085601f830112614abc57600080fd5b81356020614acc61442a83614a59565b82815260059290921b84018101918181019089841115614aeb57600080fd5b8286015b84811015614b2357803586811115614b075760008081fd5b614b158c86838b010161440b565b845250918301918301614aef565b5096505086013592505080821115614b3a57600080fd5b506145a08582860161440b565b600080600060408486031215614b5c57600080fd5b833567ffffffffffffffff811115614b7357600080fd5b614b7f8682870161472f565b9094509250506020840135614b938161430b565b809150509250925092565b60008060408385031215614bb157600080fd5b50508035926020909101359150565b600080600060408486031215614bd557600080fd5b83359250602084013567ffffffffffffffff811115614bf357600080fd5b614bff8682870161499d565b9497909650939450505050565b60008060408385031215614c1f57600080fd5b823567ffffffffffffffff80821115614c3757600080fd5b614c438683870161440b565b93506020850135915080821115614b3a57600080fd5b60006020808385031215614c6c57600080fd5b823567ffffffffffffffff80821115614c8457600080fd5b818501915085601f830112614c9857600080fd5b8135614ca661442a82614a59565b81815260059190911b83018401908481019088831115614cc557600080fd5b8585015b83811015614cfd57803585811115614ce15760008081fd5b614cef8b89838a010161440b565b845250918601918601614cc9565b5098975050505050505050565b60008060408385031215614d1d57600080fd5b823591506020830135614d2f81614461565b809150509250929050565b600060208284031215614d4c57600080fd5b8135611c63816142fc565b60008060408385031215614d6a57600080fd5b823591506020830135614d2f8161430b565b60008060408385031215614d8f57600080fd5b823591506020830135614d2f816142fc565b60008060008060608587031215614db757600080fd5b843567ffffffffffffffff811115614dce57600080fd5b614dda8782880161472f565b9095509350506020850135614dee816142fc565b91506040850135614dfe816142fc565b939692955090935050565b60008060008060008060c08789031215614e2257600080fd5b8635614e2d8161465f565b95506020870135614e3d816142fc565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614ebc57614ebc614e62565b02949350505050565b8051801515811461452957600080fd5b600060208284031215614ee757600080fd5b611c6382614ec5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614f3557614f35614e62565b60010192915050565b828152604060208201526000611c6060408301846145ce565b600082601f830112614f6857600080fd5b8151614f7661442a826143c5565b818152846020838601011115614f8b57600080fd5b614f9c8260208301602087016145aa565b949350505050565b600060208284031215614fb657600080fd5b815167ffffffffffffffff811115614fcd57600080fd5b614f9c84828501614f57565b805161452981614461565b600060208284031215614ff657600080fd5b8151611c6381614461565b80516145298161465f565b6000806040838503121561501f57600080fd5b825161502a8161465f565b6020939093015192949293505050565b60006020828403121561504c57600080fd5b8151611c638161465f565b60006020828403121561506957600080fd5b8151611c63816142fc565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036150a5576150a5614e62565b5060010190565b80516145298161430b565b805167ffffffffffffffff8116811461452957600080fd5b6000602082840312156150e157600080fd5b815167ffffffffffffffff808211156150f957600080fd5b90830190610140828603121561510e57600080fd5b61511661434c565b61511f83615001565b815261512d602084016150ac565b602082015260408301518281111561514457600080fd5b61515087828601614f57565b60408301525061516260608401614fd9565b606082015261517360808401615001565b608082015261518460a084016150b7565b60a082015261519560c084016150ac565b60c08201526151a660e08401614fd9565b60e08201526101006151b9818501614ec5565b9082015261012083810151838111156151d157600080fd5b6151dd88828701614f57565b918301919091525095945050505050565b808201808211156114ba576114ba614e62565b600061ffff80831681810361521857615218614e62565b6001019392505050565b60006020828403121561523457600080fd5b5051919050565b818103818111156114ba576114ba614e62565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156152c3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526152b18683516145ce565b95509382019390820190600101615277565b5050858403818701525050506152d981856145ce565b95945050505050565b600181811c908216806152f657607f821691505b60208210810361532f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000808354615343816152e2565b6001828116801561535b576001811461538e576153bd565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841687528215158302870194506153bd565b8760005260208060002060005b858110156153b45781548a82015290840190820161539b565b50505082870194505b50929695505050505050565b600081546153d6816152e2565b8085526020600183811680156153f3576001811461542b57615459565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550615459565b866000528260002060005b858110156154515781548a8201860152908301908401615436565b890184019650505b505050505092915050565b60a08152600061547760a08301886153c9565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156154e9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526154d783836153c9565b9486019492506001918201910161549e565b505086810360408801526154fd818b6153c9565b945050505050846060840152828103608084015261551b81856145ce565b98975050505050505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f8211156129eb57600081815260208120601f850160051c810160208610156155a25750805b601f850160051c820191505b818110156113e3578281556001016155ae565b815167ffffffffffffffff8111156155db576155db61431d565b6155ef816155e984546152e2565b8461557b565b602080601f831160018114615642576000841561560c5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556113e3565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561568f57888601518255948401946001909101908401615670565b50858210156156cb57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600060208083850312156156ee57600080fd5b825167ffffffffffffffff81111561570557600080fd5b8301601f8101851361571657600080fd5b805161572461442a82614a59565b81815260059190911b8201830190838101908783111561574357600080fd5b928401925b8284101561576157835182529284019290840190615748565b979650505050505050565b600063ffffffff80831681810361521857615218614e62565b8051602080830151919081101561532f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b60208152600082516101408060208501526157e66101608501836145ce565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08086850301604087015261582284836145ce565b93506040870151915061584d606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526158ae84836145ce565b935060e087015191506101008187860301818801526158cd85846145ce565b9450808801519250506101208187860301818801526158ec85846145ce565b9450808801519250505061590f828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561595157615951614e62565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261599457615994615956565b500490565b818103600083128015838313168383128216171561419c5761419c614e62565b6000826159c8576159c8615956565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615615a1c57615a1c614e62565b500590565b8082018281126000831280158216821582161715615a4157615a41614e62565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036150a5576150a5614e62565b60007f80000000000000000000000000000000000000000000000000000000000000008203615aab57615aab614e62565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

	error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.CheckErrorHandler(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, errCode, extraData)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) CheckErrorHandler(errCode *big.Int, extraData []byte) (CheckErrorHandler,

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

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCaller) GetUpkeepInfo(opts *bind.CallOpts, upkeepId *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	var out []interface{}
	err := _VerifiableLoadStreamsLookupUpkeep.contract.Call(opts, &out, "getUpkeepInfo", upkeepId)

	if err != nil {
		return *new(IAutomationV21PlusCommonUpkeepInfoLegacy), err
	}

	out0 := *abi.ConvertType(out[0], new(IAutomationV21PlusCommonUpkeepInfoLegacy)).(*IAutomationV21PlusCommonUpkeepInfoLegacy)

	return out0, err

}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepSession) GetUpkeepInfo(upkeepId *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
	return _VerifiableLoadStreamsLookupUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadStreamsLookupUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadStreamsLookupUpkeep *VerifiableLoadStreamsLookupUpkeepCallerSession) GetUpkeepInfo(upkeepId *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error) {
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

	CheckErrorHandler(opts *bind.CallOpts, errCode *big.Int, extraData []byte) (CheckErrorHandler,

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

	GetUpkeepInfo(opts *bind.CallOpts, upkeepId *big.Int) (IAutomationV21PlusCommonUpkeepInfoLegacy, error)

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
