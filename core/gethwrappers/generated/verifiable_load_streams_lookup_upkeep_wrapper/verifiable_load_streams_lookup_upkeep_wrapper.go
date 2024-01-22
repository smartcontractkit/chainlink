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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c0604052604261014081815261010091829190620061d161016039815260200160405180608001604052806042815260200162006213604291399052620000be906016906002620003c7565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000ee908262000543565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b602082015260189062000120908262000543565b503480156200012e57600080fd5b506040516200625538038062006255833981016040819052620001519162000625565b81813380600081620001aa5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001dd57620001dd816200031c565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa1580156200023a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000260919062000668565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002c6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ec919062000699565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c05250620006c0915050565b336001600160a01b03821603620003765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620001a1565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821562000412579160200282015b8281111562000412578251829062000401908262000543565b5091602001919060010190620003e8565b506200042092915062000424565b5090565b80821115620004205760006200043b828262000445565b5060010162000424565b5080546200045390620004b4565b6000825580601f1062000464575050565b601f01602090049060005260206000209081019062000484919062000487565b50565b5b8082111562000420576000815560010162000488565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004c957607f821691505b602082108103620004ea57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200053e57600081815260208120601f850160051c81016020861015620005195750805b601f850160051c820191505b818110156200053a5782815560010162000525565b5050505b505050565b81516001600160401b038111156200055f576200055f6200049e565b6200057781620005708454620004b4565b84620004f0565b602080601f831160018114620005af5760008415620005965750858301515b600019600386901b1c1916600185901b1785556200053a565b600085815260208120601f198616915b82811015620005e057888601518255948401946001909101908401620005bf565b5085821015620005ff5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200048457600080fd5b600080604083850312156200063957600080fd5b825162000646816200060f565b602084015190925080151581146200065d57600080fd5b809150509250929050565b600080604083850312156200067c57600080fd5b825162000689816200060f565b6020939093015192949293505050565b600060208284031215620006ac57600080fd5b8151620006b9816200060f565b9392505050565b60805160a05160c05160e051615abb62000716600039600081816105680152611f7a0152600081816109bc0152613ca701526000818161082701526136f5015260008181610d7901526136ca0152615abb6000f3fe6080604052600436106104d55760003560e01c806379ea994311610279578063a6b594751161015e578063d6051a72116100d6578063e45530831161008a578063fa333dfb1161006f578063fa333dfb14610f88578063fba7ffa31461103b578063fcdc1f631461106857600080fd5b8063e455308314610f52578063f2fde38b14610f6857600080fd5b8063daee1aeb116100bb578063daee1aeb14610ee5578063dbef701e14610f05578063e0114adb14610f2557600080fd5b8063d6051a7214610ea5578063da6cba4714610ec557600080fd5b8063b657bc9c1161012d578063c041982211610112578063c041982214610e50578063c98f10b014610e70578063d4c2490014610e8557600080fd5b8063b657bc9c14610e10578063becde0e114610e3057600080fd5b8063a6b5947514610d9b578063a72aa27e14610dbb578063af953a4a14610ddb578063afb28d1f14610dfb57600080fd5b8063924ca578116101f15780639b429354116101c05780639d385eaa116101a55780639d385eaa14610d275780639d6f1cc714610d47578063a654824814610d6757600080fd5b80639b42935414610cc95780639b51fb0d14610cf657600080fd5b8063924ca57814610c3f578063948108f714610c5f57806396cebc7c14610c7f5780639ac542eb14610c9f57600080fd5b80638340507c11610248578063873c75861161022d578063873c758614610bc75780638da5cb5b14610be75780638fcb3fba14610c1257600080fd5b80638340507c14610b8757806386e330af14610ba757600080fd5b806379ea994314610afa5780637b10399914610b1a5780637e7a46dc14610b475780638243444a14610b6757600080fd5b806345d2ec17116103ba57806360457ff5116103325780637145f11b116102e657806376721303116102cb5780637672130314610a98578063776898c814610ac557806379ba509714610ae557600080fd5b80637145f11b14610a3b57806373644cce14610a6b57600080fd5b8063642f6cef11610317578063642f6cef146109aa57806369cdbadb146109ee5780636e04ff0d14610a1b57600080fd5b806360457ff514610958578063636092e81461098557600080fd5b80635147cd591161038957806357970e931161036e57806357970e93146108f65780635d4ee7f3146109235780635f17e6161461093857600080fd5b80635147cd59146108a457806351c98be3146108d657600080fd5b806345d2ec17146107e8578063469820931461081557806346e7a63e146108495780634b56a42e1461087657600080fd5b806320e3dbd41161044d5780632a9032d31161041c578063328ffd1111610401578063328ffd111461077b5780633ebe8d6c146107a85780634585e33b146107c857600080fd5b80632a9032d3146107095780632b20e3971461072957600080fd5b806320e3dbd41461067c5780632636aecf1461069c57806328c4b57b146106bc57806329e0a841146106dc57600080fd5b806319d97a94116104a45780631e010439116104895780631e010439146105ea578063206c32e814610627578063207b65161461065c57600080fd5b806319d97a941461059d5780631cdde251146105ca57600080fd5b806306c1cc00146104e1578063077ac621146105035780630b7d33e61461053657806312c550271461055657600080fd5b366104dc57005b600080fd5b3480156104ed57600080fd5b506105016104fc366004614429565b611095565b005b34801561050f57600080fd5b5061052361051e3660046144dc565b6112e4565b6040519081526020015b60405180910390f35b34801561054257600080fd5b50610501610551366004614511565b611322565b34801561056257600080fd5b5061058a7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff909116815260200161052d565b3480156105a957600080fd5b506105bd6105b8366004614558565b6113b0565b60405161052d91906145df565b3480156105d657600080fd5b506105016105e5366004614614565b61146d565b3480156105f657600080fd5b5061060a610605366004614558565b6115aa565b6040516bffffffffffffffffffffffff909116815260200161052d565b34801561063357600080fd5b50610647610642366004614679565b61163f565b6040805192835260208301919091520161052d565b34801561066857600080fd5b506105bd610677366004614558565b6116c2565b34801561068857600080fd5b506105016106973660046146a5565b61171a565b3480156106a857600080fd5b506105016106b7366004614707565b6118e4565b3480156106c857600080fd5b506105236106d7366004614781565b611bad565b3480156106e857600080fd5b506106fc6106f7366004614558565b611c18565b60405161052d91906147ad565b34801561071557600080fd5b506105016107243660046148ee565b611d1d565b34801561073557600080fd5b506011546107569073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161052d565b34801561078757600080fd5b50610523610796366004614558565b60036020526000908152604090205481565b3480156107b457600080fd5b506105236107c3366004614558565b611dfe565b3480156107d457600080fd5b506105016107e3366004614972565b611e67565b3480156107f457600080fd5b50610808610803366004614679565b612086565b60405161052d91906149a8565b34801561082157600080fd5b506105237f000000000000000000000000000000000000000000000000000000000000000081565b34801561085557600080fd5b50610523610864366004614558565b600a6020526000908152604090205481565b34801561088257600080fd5b50610896610891366004614a10565b6120f5565b60405161052d929190614ada565b3480156108b057600080fd5b506108c46108bf366004614558565b612149565b60405160ff909116815260200161052d565b3480156108e257600080fd5b506105016108f1366004614af5565b6121dd565b34801561090257600080fd5b506012546107569073ffffffffffffffffffffffffffffffffffffffff1681565b34801561092f57600080fd5b50610501612281565b34801561094457600080fd5b50610501610953366004614b4c565b6123bc565b34801561096457600080fd5b50610523610973366004614558565b60076020526000908152604090205481565b34801561099157600080fd5b5060155461060a906bffffffffffffffffffffffff1681565b3480156109b657600080fd5b506109de7f000000000000000000000000000000000000000000000000000000000000000081565b604051901515815260200161052d565b3480156109fa57600080fd5b50610523610a09366004614558565b60086020526000908152604090205481565b348015610a2757600080fd5b50610896610a36366004614972565b612489565b348015610a4757600080fd5b506109de610a56366004614558565b600b6020526000908152604090205460ff1681565b348015610a7757600080fd5b50610523610a86366004614558565b6000908152600c602052604090205490565b348015610aa457600080fd5b50610523610ab3366004614558565b60046020526000908152604090205481565b348015610ad157600080fd5b506109de610ae0366004614558565b6126b2565b348015610af157600080fd5b50610501612704565b348015610b0657600080fd5b50610756610b15366004614558565b612801565b348015610b2657600080fd5b506013546107569073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b5357600080fd5b50610501610b62366004614b6e565b612895565b348015610b7357600080fd5b50610501610b82366004614b6e565b612926565b348015610b9357600080fd5b50610501610ba2366004614bba565b612980565b348015610bb357600080fd5b50610501610bc2366004614c07565b61299e565b348015610bd357600080fd5b50610808610be2366004614b4c565b6129b1565b348015610bf357600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610756565b348015610c1e57600080fd5b50610523610c2d366004614558565b60056020526000908152604090205481565b348015610c4b57600080fd5b50610501610c5a366004614b4c565b612a6e565b348015610c6b57600080fd5b50610501610c7a366004614cb8565b612cb3565b348015610c8b57600080fd5b50610501610c9a366004614ce8565b612dcb565b348015610cab57600080fd5b506015546108c4906c01000000000000000000000000900460ff1681565b348015610cd557600080fd5b50610501610ce4366004614b4c565b60009182526009602052604090912055565b348015610d0257600080fd5b5061058a610d11366004614558565b600e6020526000908152604090205461ffff1681565b348015610d3357600080fd5b50610808610d42366004614558565b612fd5565b348015610d5357600080fd5b506105bd610d62366004614558565b613037565b348015610d7357600080fd5b506105237f000000000000000000000000000000000000000000000000000000000000000081565b348015610da757600080fd5b50610501610db6366004614781565b6130e3565b348015610dc757600080fd5b50610501610dd6366004614d05565b61314c565b348015610de757600080fd5b50610501610df6366004614558565b6131f7565b348015610e0757600080fd5b506105bd61327d565b348015610e1c57600080fd5b5061060a610e2b366004614558565b61328a565b348015610e3c57600080fd5b50610501610e4b3660046148ee565b6132e2565b348015610e5c57600080fd5b50610808610e6b366004614b4c565b61337c565b348015610e7c57600080fd5b506105bd613479565b348015610e9157600080fd5b50610501610ea0366004614d2a565b613486565b348015610eb157600080fd5b50610647610ec0366004614b4c565b613505565b348015610ed157600080fd5b50610501610ee0366004614d4f565b61356e565b348015610ef157600080fd5b50610501610f003660046148ee565b6138d5565b348015610f1157600080fd5b50610523610f20366004614b4c565b6139a0565b348015610f3157600080fd5b50610523610f40366004614558565b60096020526000908152604090205481565b348015610f5e57600080fd5b5061052360145481565b348015610f7457600080fd5b50610501610f833660046146a5565b6139d1565b348015610f9457600080fd5b506105bd610fa3366004614db7565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b34801561104757600080fd5b50610523611056366004614558565b60066020526000908152604090205481565b34801561107457600080fd5b50610523611083366004614558565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b392169061117b908c1688614e3f565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af11580156111f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061121d9190614e83565b5060008860ff1667ffffffffffffffff81111561123c5761123c6142cb565b604051908082528060200260200182016040528015611265578160200160208202803683370190505b50905060005b8960ff168160ff1610156112d8576000611284846139e5565b905080838360ff168151811061129c5761129c614e9e565b602090810291909101810191909152600091825260088152604080832088905560079091529020849055806112d081614ecd565b91505061126b565b50505050505050505050565b600d602052826000526040600020602052816000526040600020818154811061130c57600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e69061137a9085908590600401614eec565b600060405180830381600087803b15801561139457600080fd5b505af11580156113a8573d6000803e3d6000fd5b505050505050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa158015611421573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114679190810190614f52565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa15801561150c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115529190810190614f52565b6040518363ffffffff1660e01b815260040161156f929190614eec565b600060405180830381600087803b15801561158957600080fd5b505af115801561159d573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e010439906024015b602060405180830381865afa15801561161b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114679190614f92565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156116a357602002820191906000526020600020905b81548152602001906001019080831161168f575b505050505090506116b5818251613ab3565b92509250505b9250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b651690602401611404565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa1580156117b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117d49190614fba565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015611877573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061189b9190614fe8565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611ba257600089898381811061190457611904614e9e565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161193d91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401611969929190614eec565b600060405180830381600087803b15801561198357600080fd5b505af1158015611997573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015611a0d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a319190615005565b90508060ff16600103611b8d576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611aba573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b009190810190614f52565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611b599086908590600401614eec565b600060405180830381600087803b158015611b7357600080fd5b505af1158015611b87573d6000803e3d6000fd5b50505050505b50508080611b9a90615022565b9150506118e8565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611c0e93830182828015611c0257602002820191906000526020600020905b815481526020019060010190808311611bee575b50505050508484613b38565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611cd7573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611467919081019061507d565b8060005b818160ff161015611df85760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611d5f57611d5f614e9e565b905060200201356040518263ffffffff1660e01b8152600401611d8491815260200190565b600060405180830381600087803b158015611d9e57600080fd5b505af1158015611db2573d6000803e3d6000fd5b50505050611de584848360ff16818110611dce57611dce614e9e565b90506020020135600f613c9790919063ffffffff16565b5080611df081614ecd565b915050611d21565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611e5f576000858152600d6020908152604080832061ffff85168452909152902054611e4b908361519c565b915080611e57816151af565b915050611e14565b509392505050565b60005a9050600080611e7b84860186614a10565b91509150600081806020019051810190611e9591906151d0565b60008181526005602090815260408083205460049092528220549293509190611ebc613ca3565b905082600003611edc576000848152600560205260409020819055612037565b600084815260036020526040812054611ef584846151e9565b611eff91906151e9565b6000868152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611f7157602002820191906000526020600020905b815481526020019060010190808311611f5d575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff16815103611fec5781611fae816151af565b6000898152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000868152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055888352600c8252928220805493840181558252902001555b60008481526006602052604081205461205190600161519c565b600086815260066020908152604080832084905560049091529020839055905061207b8583612a6e565b6112d88589846130e3565b6000828152600d6020908152604080832061ffff851684528252918290208054835181840281018401909452808452606093928301828280156120e857602002820191906000526020600020905b8154815260200190600101908083116120d4575b5050505050905092915050565b600060606000848460405160200161210e9291906151fc565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa1580156121b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114679190615005565b8160005b8181101561227a5730635f17e61686868481811061220157612201614e9e565b90506020020135856040518363ffffffff1660e01b815260040161223592919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561224f57600080fd5b505af1158015612263573d6000803e3d6000fd5b50505050808061227290615022565b9150506121e1565b5050505050565b612289613d45565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156122f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061231c91906151d0565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af1158015612394573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123b89190614e83565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c90915281206123f4916141ca565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff1611612450576000848152600d6020908152604080832061ffff85168452909152812061243e916141ca565b80612448816151af565b915050612409565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000606060005a905060006124a085870187614558565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff8111156124d9576124d96142cb565b6040519080825280601f01601f191660200182016040528015612503576020820181803683370190505b50604051602001612515929190614eec565b60405160208183030381529060405290506000612530613ca3565b9050600061253d866126b2565b90505b835a61254c90896151e9565b6125589061271061519c565b10156125995781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612540565b806125b15760008398509850505050505050506116bb565b6040517f6665656449644865780000000000000000000000000000000000000000000000602082015260009060290160405160208183030381529060405280519060200120601760405160200161260891906152e3565b604051602081830303815290604052805190602001200361262a57508161262d565b50425b601760166018838a60405160200161264791815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527ff055e4a20000000000000000000000000000000000000000000000000000000082526126a99594939291600401615412565b60405180910390fd5b60008181526005602052604081205481036126cf57506001919050565b6000828152600360209081526040808320546004909252909120546126f2613ca3565b6126fc91906151e9565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612785576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016126a9565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa158015612871573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114679190614fe8565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b5906128ef908690869086906004016154d5565b600060405180830381600087803b15801561290957600080fd5b505af115801561291d573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d35906128ef908690869086906004016154d5565b601761298c838261556f565b506018612999828261556f565b505050565b80516123b89060169060208401906141e8565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612a28573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611c119190810190615689565b601454600083815260026020526040902054612a8a90836151e9565b11156123b8576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612b00573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612b46919081019061507d565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612bbb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612bdf9190614f92565b601554909150612c039082906c01000000000000000000000000900460ff16614e3f565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611df857601554612c469085906bffffffffffffffffffffffff16612cb3565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612d3b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d5f9190614e83565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f79060440161137a565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015612e2a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612e709190810190615689565b80519091506000612e7f613ca3565b905060005b8281101561227a576000848281518110612ea057612ea0614e9e565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612f20573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f449190615005565b90508060ff16600103612fc0578660ff16600003612f90576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4612fc0565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b50508080612fcd90615022565b915050612e84565b6000818152600c602090815260409182902080548351818402810184019094528084526060939283018282801561302b57602002820191906000526020600020905b815481526020019060010190808311613017575b50505050509050919050565b6016818154811061304757600080fd5b90600052602060002001600091509050805461306290615290565b80601f016020809104026020016040519081016040528092919081815260200182805461308e90615290565b80156130db5780601f106130b0576101008083540402835291602001916130db565b820191906000526020600020905b8154815290600101906020018083116130be57829003601f168201915b505050505081565b6000838152600760205260409020545b805a6130ff90856151e9565b61310b9061271061519c565b1015611df85781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556130f3565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156131c457600080fd5b505af11580156131d8573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561326957600080fd5b505af115801561227a573d6000803e3d6000fd5b6017805461306290615290565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff169063b657bc9c906024016115fe565b8060005b818163ffffffff161015611df8573063af953a4a858563ffffffff851681811061331257613312614e9e565b905060200201356040518263ffffffff1660e01b815260040161333791815260200190565b600060405180830381600087803b15801561335157600080fd5b505af1158015613365573d6000803e3d6000fd5b5050505080806133749061571a565b9150506132e6565b6060600061338a600f613dc8565b90508084106133c5576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826000036133da576133d784826151e9565b92505b60008367ffffffffffffffff8111156133f5576133f56142cb565b60405190808252806020026020018201604052801561341e578160200160208202803683370190505b50905060005b8481101561347057613441613439828861519c565b600f90613dd2565b82828151811061345357613453614e9e565b60209081029190910101528061346881615022565b915050613424565b50949350505050565b6018805461306290615290565b6000613490613ca3565b90508160ff166000036134d1576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a4505050565b6000828152600c6020908152604080832080548251818502810185019093528083528493849392919083018282801561355d57602002820191906000526020600020905b815481526020019060010190808311613549575b505050505090506116b58185613ab3565b8260005b818110156113a857600086868381811061358e5761358e614e9e565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016135c791815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016135f3929190614eec565b600060405180830381600087803b15801561360d57600080fd5b505af1158015613621573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015613697573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136bb9190615005565b90508060ff166001036138c0577f000000000000000000000000000000000000000000000000000000000000000060ff87161561371557507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb3089858860405160200161374991815260200190565b60405160208183030381529060405261376190615733565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa1580156137ec573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526138329190810190614f52565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d359061388b9087908590600401614eec565b600060405180830381600087803b1580156138a557600080fd5b505af11580156138b9573d6000803e3d6000fd5b5050505050505b505080806138cd90615022565b915050613572565b8060005b81811015611df85760008484838181106138f5576138f5614e9e565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161392e91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b815260040161395a929190614eec565b600060405180830381600087803b15801561397457600080fd5b505af1158015613988573d6000803e3d6000fd5b5050505050808061399890615022565b9150506138d9565b600c60205281600052604060002081815481106139bc57600080fd5b90600052602060002001600091509150505481565b6139d9613d45565b6139e281613dde565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613a40908690600401615775565b6020604051808303816000875af1158015613a5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a8391906151d0565b9050613a90600f82613ed3565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b815160009081908190841580613ac95750808510155b15613ad2578094505b60008092505b85831015613b2e57866001613aed85856151e9565b613af791906151e9565b81518110613b0757613b07614e9e565b602002602001015181613b1a919061519c565b905082613b2681615022565b935050613ad8565b9694955050505050565b82516000908190831580613b4c5750808410155b15613b55578093505b60008467ffffffffffffffff811115613b7057613b706142cb565b604051908082528060200260200182016040528015613b99578160200160208202803683370190505b509050600092505b84831015613c0757866001613bb685856151e9565b613bc091906151e9565b81518110613bd057613bd0614e9e565b6020026020010151818481518110613bea57613bea614e9e565b602090810291909101015282613bff81615022565b935050613ba1565b613c2081600060018451613c1b91906151e9565b613edf565b85606403613c59578060018251613c3791906151e9565b81518110613c4757613c47614e9e565b60200260200101519350505050611c11565b806064825188613c6991906158c7565b613c739190615933565b81518110613c8357613c83614e9e565b602002602001015193505050509392505050565b6000611c118383614057565b60007f000000000000000000000000000000000000000000000000000000000000000015613d4057606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613d17573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d3b91906151d0565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613dc6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016126a9565b565b6000611467825490565b6000611c118383614151565b3373ffffffffffffffffffffffffffffffffffffffff821603613e5d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016126a9565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611c11838361417b565b8181808203613eef575050505050565b6000856002613efe8787615947565b613f089190615967565b613f1290876159cf565b81518110613f2257613f22614e9e565b602002602001015190505b818313614031575b80868481518110613f4857613f48614e9e565b60200260200101511015613f685782613f60816159f7565b935050613f35565b858281518110613f7a57613f7a614e9e565b6020026020010151811015613f9b5781613f9381615a28565b925050613f68565b81831361402c57858281518110613fb457613fb4614e9e565b6020026020010151868481518110613fce57613fce614e9e565b6020026020010151878581518110613fe857613fe8614e9e565b6020026020010188858151811061400157614001614e9e565b6020908102919091010191909152528261401a816159f7565b935050818061402890615a28565b9250505b613f2d565b8185121561404457614044868684613edf565b838312156113a8576113a8868486613edf565b6000818152600183016020526040812054801561414057600061407b6001836151e9565b855490915060009061408f906001906151e9565b90508181146140f45760008660000182815481106140af576140af614e9e565b90600052602060002001549050808760000184815481106140d2576140d2614e9e565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061410557614105615a7f565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611467565b6000915050611467565b5092915050565b600082600001828154811061416857614168614e9e565b9060005260206000200154905092915050565b60008181526001830160205260408120546141c257508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611467565b506000611467565b50805460008255906000526020600020908101906139e2919061423e565b82805482825590600052602060002090810192821561422e579160200282015b8281111561422e578251829061421e908261556f565b5091602001919060010190614208565b5061423a929150614253565b5090565b5b8082111561423a576000815560010161423f565b8082111561423a5760006142678282614270565b50600101614253565b50805461427c90615290565b6000825580601f1061428c575050565b601f0160209004906000526020600020908101906139e2919061423e565b60ff811681146139e257600080fd5b63ffffffff811681146139e257600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff8111828210171561431e5761431e6142cb565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561436b5761436b6142cb565b604052919050565b600067ffffffffffffffff82111561438d5761438d6142cb565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f8301126143ca57600080fd5b81356143dd6143d882614373565b614324565b8181528460208386010111156143f257600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff811681146139e257600080fd5b600080600080600080600060e0888a03121561444457600080fd5b873561444f816142aa565b9650602088013561445f816142b9565b9550604088013561446f816142aa565b9450606088013567ffffffffffffffff81111561448b57600080fd5b6144978a828b016143b9565b94505060808801356144a88161440f565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff811681146144d757600080fd5b919050565b6000806000606084860312156144f157600080fd5b83359250614501602085016144c5565b9150604084013590509250925092565b6000806040838503121561452457600080fd5b82359150602083013567ffffffffffffffff81111561454257600080fd5b61454e858286016143b9565b9150509250929050565b60006020828403121561456a57600080fd5b5035919050565b60005b8381101561458c578181015183820152602001614574565b50506000910152565b600081518084526145ad816020860160208601614571565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611c116020830184614595565b73ffffffffffffffffffffffffffffffffffffffff811681146139e257600080fd5b600080600080600080600060e0888a03121561462f57600080fd5b873596506020880135614641816145f2565b95506040880135614651816142aa565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b6000806040838503121561468c57600080fd5b8235915061469c602084016144c5565b90509250929050565b6000602082840312156146b757600080fd5b8135611c11816145f2565b60008083601f8401126146d457600080fd5b50813567ffffffffffffffff8111156146ec57600080fd5b6020830191508360208260051b85010111156116bb57600080fd5b600080600080600080600060c0888a03121561472257600080fd5b873567ffffffffffffffff81111561473957600080fd5b6147458a828b016146c2565b9098509650506020880135614759816142aa565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b60008060006060848603121561479657600080fd5b505081359360208301359350604090920135919050565b602081526147d460208201835173ffffffffffffffffffffffffffffffffffffffff169052565b600060208301516147ed604084018263ffffffff169052565b50604083015161014080606085015261480a610160850183614595565b9150606085015161482b60808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614897818701836bffffffffffffffffffffffff169052565b86015190506101206148ac8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018387015290506148e48382614595565b9695505050505050565b6000806020838503121561490157600080fd5b823567ffffffffffffffff81111561491857600080fd5b614924858286016146c2565b90969095509350505050565b60008083601f84011261494257600080fd5b50813567ffffffffffffffff81111561495a57600080fd5b6020830191508360208285010111156116bb57600080fd5b6000806020838503121561498557600080fd5b823567ffffffffffffffff81111561499c57600080fd5b61492485828601614930565b6020808252825182820181905260009190848201906040850190845b818110156149e0578351835292840192918401916001016149c4565b50909695505050505050565b600067ffffffffffffffff821115614a0657614a066142cb565b5060051b60200190565b60008060408385031215614a2357600080fd5b823567ffffffffffffffff80821115614a3b57600080fd5b818501915085601f830112614a4f57600080fd5b81356020614a5f6143d8836149ec565b82815260059290921b84018101918181019089841115614a7e57600080fd5b8286015b84811015614ab657803586811115614a9a5760008081fd5b614aa88c86838b01016143b9565b845250918301918301614a82565b5096505086013592505080821115614acd57600080fd5b5061454e858286016143b9565b8215158152604060208201526000611c0e6040830184614595565b600080600060408486031215614b0a57600080fd5b833567ffffffffffffffff811115614b2157600080fd5b614b2d868287016146c2565b9094509250506020840135614b41816142b9565b809150509250925092565b60008060408385031215614b5f57600080fd5b50508035926020909101359150565b600080600060408486031215614b8357600080fd5b83359250602084013567ffffffffffffffff811115614ba157600080fd5b614bad86828701614930565b9497909650939450505050565b60008060408385031215614bcd57600080fd5b823567ffffffffffffffff80821115614be557600080fd5b614bf1868387016143b9565b93506020850135915080821115614acd57600080fd5b60006020808385031215614c1a57600080fd5b823567ffffffffffffffff80821115614c3257600080fd5b818501915085601f830112614c4657600080fd5b8135614c546143d8826149ec565b81815260059190911b83018401908481019088831115614c7357600080fd5b8585015b83811015614cab57803585811115614c8f5760008081fd5b614c9d8b89838a01016143b9565b845250918601918601614c77565b5098975050505050505050565b60008060408385031215614ccb57600080fd5b823591506020830135614cdd8161440f565b809150509250929050565b600060208284031215614cfa57600080fd5b8135611c11816142aa565b60008060408385031215614d1857600080fd5b823591506020830135614cdd816142b9565b60008060408385031215614d3d57600080fd5b823591506020830135614cdd816142aa565b60008060008060608587031215614d6557600080fd5b843567ffffffffffffffff811115614d7c57600080fd5b614d88878288016146c2565b9095509350506020850135614d9c816142aa565b91506040850135614dac816142aa565b939692955090935050565b60008060008060008060c08789031215614dd057600080fd5b8635614ddb816145f2565b95506020870135614deb816142aa565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614e6a57614e6a614e10565b02949350505050565b805180151581146144d757600080fd5b600060208284031215614e9557600080fd5b611c1182614e73565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614ee357614ee3614e10565b60010192915050565b828152604060208201526000611c0e6040830184614595565b600082601f830112614f1657600080fd5b8151614f246143d882614373565b818152846020838601011115614f3957600080fd5b614f4a826020830160208701614571565b949350505050565b600060208284031215614f6457600080fd5b815167ffffffffffffffff811115614f7b57600080fd5b614f4a84828501614f05565b80516144d78161440f565b600060208284031215614fa457600080fd5b8151611c118161440f565b80516144d7816145f2565b60008060408385031215614fcd57600080fd5b8251614fd8816145f2565b6020939093015192949293505050565b600060208284031215614ffa57600080fd5b8151611c11816145f2565b60006020828403121561501757600080fd5b8151611c11816142aa565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361505357615053614e10565b5060010190565b80516144d7816142b9565b805167ffffffffffffffff811681146144d757600080fd5b60006020828403121561508f57600080fd5b815167ffffffffffffffff808211156150a757600080fd5b9083019061014082860312156150bc57600080fd5b6150c46142fa565b6150cd83614faf565b81526150db6020840161505a565b60208201526040830151828111156150f257600080fd5b6150fe87828601614f05565b60408301525061511060608401614f87565b606082015261512160808401614faf565b608082015261513260a08401615065565b60a082015261514360c0840161505a565b60c082015261515460e08401614f87565b60e0820152610100615167818501614e73565b90820152610120838101518381111561517f57600080fd5b61518b88828701614f05565b918301919091525095945050505050565b8082018082111561146757611467614e10565b600061ffff8083168181036151c6576151c6614e10565b6001019392505050565b6000602082840312156151e257600080fd5b5051919050565b8181038181111561146757611467614e10565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015615271577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261525f868351614595565b95509382019390820190600101615225565b5050858403818701525050506152878185614595565b95945050505050565b600181811c908216806152a457607f821691505b6020821081036152dd577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60008083546152f181615290565b60018281168015615309576001811461533c5761536b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008416875282151583028701945061536b565b8760005260208060002060005b858110156153625781548a820152908401908201615349565b50505082870194505b50929695505050505050565b6000815461538481615290565b8085526020600183811680156153a157600181146153d957615407565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b8901019550615407565b866000528260002060005b858110156153ff5781548a82018601529083019084016153e4565b890184019650505b505050505092915050565b60a08152600061542560a0830188615377565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015615497577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526154858383615377565b9486019492506001918201910161544c565b505086810360408801526154ab818b615377565b94505050505084606084015282810360808401526154c98185614595565b98975050505050505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f82111561299957600081815260208120601f850160051c810160208610156155505750805b601f850160051c820191505b818110156113a85782815560010161555c565b815167ffffffffffffffff811115615589576155896142cb565b61559d816155978454615290565b84615529565b602080601f8311600181146155f057600084156155ba5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556113a8565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561563d5788860151825594840194600190910190840161561e565b508582101561567957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000602080838503121561569c57600080fd5b825167ffffffffffffffff8111156156b357600080fd5b8301601f810185136156c457600080fd5b80516156d26143d8826149ec565b81815260059190911b820183019083810190878311156156f157600080fd5b928401925b8284101561570f578351825292840192908401906156f6565b979650505050505050565b600063ffffffff8083168181036151c6576151c6614e10565b805160208083015191908110156152dd577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b6020815260008251610140806020850152615794610160850183614595565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808685030160408701526157d08483614595565b9350604087015191506157fb606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e087015261585c8483614595565b935060e0870151915061010081878603018188015261587b8584614595565b94508088015192505061012081878603018188015261589a8584614595565b945080880151925050506158bd828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156158ff576158ff614e10565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261594257615942615904565b500490565b818103600083128015838313168383128216171561414a5761414a614e10565b60008261597657615976615904565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f8000000000000000000000000000000000000000000000000000000000000000831416156159ca576159ca614e10565b500590565b80820182811260008312801582168215821617156159ef576159ef614e10565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361505357615053614e10565b60007f80000000000000000000000000000000000000000000000000000000000000008203615a5957615a59614e10565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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
