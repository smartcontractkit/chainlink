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

var VerifiableLoadMercuryUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c0604052604261014081815261010091829190620060c161016039815260200160405180608001604052806042815260200162006103604291399052620000be906016906002620003c7565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000ee908262000543565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b602082015260189062000120908262000543565b503480156200012e57600080fd5b506040516200614538038062006145833981016040819052620001519162000625565b81813380600081620001aa5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001dd57620001dd816200031c565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa1580156200023a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000260919062000668565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002c6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ec919062000699565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c05250620006c0915050565b336001600160a01b03821603620003765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620001a1565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821562000412579160200282015b8281111562000412578251829062000401908262000543565b5091602001919060010190620003e8565b506200042092915062000424565b5090565b80821115620004205760006200043b828262000445565b5060010162000424565b5080546200045390620004b4565b6000825580601f1062000464575050565b601f01602090049060005260206000209081019062000484919062000487565b50565b5b8082111562000420576000815560010162000488565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004c957607f821691505b602082108103620004ea57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200053e57600081815260208120601f850160051c81016020861015620005195750805b601f850160051c820191505b818110156200053a5782815560010162000525565b5050505b505050565b81516001600160401b038111156200055f576200055f6200049e565b6200057781620005708454620004b4565b84620004f0565b602080601f831160018114620005af5760008415620005965750858301515b600019600386901b1c1916600185901b1785556200053a565b600085815260208120601f198616915b82811015620005e057888601518255948401946001909101908401620005bf565b5085821015620005ff5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200048457600080fd5b600080604083850312156200063957600080fd5b825162000646816200060f565b602084015190925080151581146200065d57600080fd5b809150509250929050565b600080604083850312156200067c57600080fd5b825162000689816200060f565b6020939093015192949293505050565b600060208284031215620006ac57600080fd5b8151620006b9816200060f565b9392505050565b60805160a05160c05160e0516159ab62000716600039600081816105680152611f7a0152600081816109bc0152613c2b0152600081816108270152613679015260008181610d79015261364e01526159ab6000f3fe6080604052600436106104d55760003560e01c806379ea994311610279578063a6b594751161015e578063d6051a72116100d6578063e45530831161008a578063fa333dfb1161006f578063fa333dfb14610f88578063fba7ffa31461103b578063fcdc1f631461106857600080fd5b8063e455308314610f52578063f2fde38b14610f6857600080fd5b8063daee1aeb116100bb578063daee1aeb14610ee5578063dbef701e14610f05578063e0114adb14610f2557600080fd5b8063d6051a7214610ea5578063da6cba4714610ec557600080fd5b8063b657bc9c1161012d578063c041982211610112578063c041982214610e50578063c98f10b014610e70578063d4c2490014610e8557600080fd5b8063b657bc9c14610e10578063becde0e114610e3057600080fd5b8063a6b5947514610d9b578063a72aa27e14610dbb578063af953a4a14610ddb578063afb28d1f14610dfb57600080fd5b8063924ca578116101f15780639b429354116101c05780639d385eaa116101a55780639d385eaa14610d275780639d6f1cc714610d47578063a654824814610d6757600080fd5b80639b42935414610cc95780639b51fb0d14610cf657600080fd5b8063924ca57814610c3f578063948108f714610c5f57806396cebc7c14610c7f5780639ac542eb14610c9f57600080fd5b80638340507c11610248578063873c75861161022d578063873c758614610bc75780638da5cb5b14610be75780638fcb3fba14610c1257600080fd5b80638340507c14610b8757806386e330af14610ba757600080fd5b806379ea994314610afa5780637b10399914610b1a5780637e7a46dc14610b475780638243444a14610b6757600080fd5b806345d2ec17116103ba57806360457ff5116103325780637145f11b116102e657806376721303116102cb5780637672130314610a98578063776898c814610ac557806379ba509714610ae557600080fd5b80637145f11b14610a3b57806373644cce14610a6b57600080fd5b8063642f6cef11610317578063642f6cef146109aa57806369cdbadb146109ee5780636e04ff0d14610a1b57600080fd5b806360457ff514610958578063636092e81461098557600080fd5b80635147cd591161038957806357970e931161036e57806357970e93146108f65780635d4ee7f3146109235780635f17e6161461093857600080fd5b80635147cd59146108a457806351c98be3146108d657600080fd5b806345d2ec17146107e8578063469820931461081557806346e7a63e146108495780634b56a42e1461087657600080fd5b806320e3dbd41161044d5780632a9032d31161041c578063328ffd1111610401578063328ffd111461077b5780633ebe8d6c146107a85780634585e33b146107c857600080fd5b80632a9032d3146107095780632b20e3971461072957600080fd5b806320e3dbd41461067c5780632636aecf1461069c57806328c4b57b146106bc57806329e0a841146106dc57600080fd5b806319d97a94116104a45780631e010439116104895780631e010439146105ea578063206c32e814610627578063207b65161461065c57600080fd5b806319d97a941461059d5780631cdde251146105ca57600080fd5b806306c1cc00146104e1578063077ac621146105035780630b7d33e61461053657806312c550271461055657600080fd5b366104dc57005b600080fd5b3480156104ed57600080fd5b506105016104fc3660046143ad565b611095565b005b34801561050f57600080fd5b5061052361051e366004614460565b6112e4565b6040519081526020015b60405180910390f35b34801561054257600080fd5b50610501610551366004614495565b611322565b34801561056257600080fd5b5061058a7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff909116815260200161052d565b3480156105a957600080fd5b506105bd6105b83660046144dc565b6113b0565b60405161052d9190614563565b3480156105d657600080fd5b506105016105e5366004614598565b61146d565b3480156105f657600080fd5b5061060a6106053660046144dc565b6115aa565b6040516bffffffffffffffffffffffff909116815260200161052d565b34801561063357600080fd5b506106476106423660046145fd565b61163f565b6040805192835260208301919091520161052d565b34801561066857600080fd5b506105bd6106773660046144dc565b6116c2565b34801561068857600080fd5b50610501610697366004614629565b61171a565b3480156106a857600080fd5b506105016106b736600461468b565b6118e4565b3480156106c857600080fd5b506105236106d7366004614705565b611bad565b3480156106e857600080fd5b506106fc6106f73660046144dc565b611c18565b60405161052d9190614731565b34801561071557600080fd5b50610501610724366004614872565b611d1d565b34801561073557600080fd5b506011546107569073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161052d565b34801561078757600080fd5b506105236107963660046144dc565b60036020526000908152604090205481565b3480156107b457600080fd5b506105236107c33660046144dc565b611dfe565b3480156107d457600080fd5b506105016107e33660046148f6565b611e67565b3480156107f457600080fd5b506108086108033660046145fd565b612086565b60405161052d919061492c565b34801561082157600080fd5b506105237f000000000000000000000000000000000000000000000000000000000000000081565b34801561085557600080fd5b506105236108643660046144dc565b600a6020526000908152604090205481565b34801561088257600080fd5b50610896610891366004614994565b6120f5565b60405161052d929190614a5e565b3480156108b057600080fd5b506108c46108bf3660046144dc565b612149565b60405160ff909116815260200161052d565b3480156108e257600080fd5b506105016108f1366004614a79565b6121dd565b34801561090257600080fd5b506012546107569073ffffffffffffffffffffffffffffffffffffffff1681565b34801561092f57600080fd5b50610501612281565b34801561094457600080fd5b50610501610953366004614ad0565b6123bc565b34801561096457600080fd5b506105236109733660046144dc565b60076020526000908152604090205481565b34801561099157600080fd5b5060155461060a906bffffffffffffffffffffffff1681565b3480156109b657600080fd5b506109de7f000000000000000000000000000000000000000000000000000000000000000081565b604051901515815260200161052d565b3480156109fa57600080fd5b50610523610a093660046144dc565b60086020526000908152604090205481565b348015610a2757600080fd5b50610896610a363660046148f6565b612489565b348015610a4757600080fd5b506109de610a563660046144dc565b600b6020526000908152604090205460ff1681565b348015610a7757600080fd5b50610523610a863660046144dc565b6000908152600c602052604090205490565b348015610aa457600080fd5b50610523610ab33660046144dc565b60046020526000908152604090205481565b348015610ad157600080fd5b506109de610ae03660046144dc565b612636565b348015610af157600080fd5b50610501612688565b348015610b0657600080fd5b50610756610b153660046144dc565b612785565b348015610b2657600080fd5b506013546107569073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b5357600080fd5b50610501610b62366004614af2565b612819565b348015610b7357600080fd5b50610501610b82366004614af2565b6128aa565b348015610b9357600080fd5b50610501610ba2366004614b3e565b612904565b348015610bb357600080fd5b50610501610bc2366004614b8b565b612922565b348015610bd357600080fd5b50610808610be2366004614ad0565b612935565b348015610bf357600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610756565b348015610c1e57600080fd5b50610523610c2d3660046144dc565b60056020526000908152604090205481565b348015610c4b57600080fd5b50610501610c5a366004614ad0565b6129f2565b348015610c6b57600080fd5b50610501610c7a366004614c3c565b612c37565b348015610c8b57600080fd5b50610501610c9a366004614c6c565b612d4f565b348015610cab57600080fd5b506015546108c4906c01000000000000000000000000900460ff1681565b348015610cd557600080fd5b50610501610ce4366004614ad0565b60009182526009602052604090912055565b348015610d0257600080fd5b5061058a610d113660046144dc565b600e6020526000908152604090205461ffff1681565b348015610d3357600080fd5b50610808610d423660046144dc565b612f59565b348015610d5357600080fd5b506105bd610d623660046144dc565b612fbb565b348015610d7357600080fd5b506105237f000000000000000000000000000000000000000000000000000000000000000081565b348015610da757600080fd5b50610501610db6366004614705565b613067565b348015610dc757600080fd5b50610501610dd6366004614c89565b6130d0565b348015610de757600080fd5b50610501610df63660046144dc565b61317b565b348015610e0757600080fd5b506105bd613201565b348015610e1c57600080fd5b5061060a610e2b3660046144dc565b61320e565b348015610e3c57600080fd5b50610501610e4b366004614872565b613266565b348015610e5c57600080fd5b50610808610e6b366004614ad0565b613300565b348015610e7c57600080fd5b506105bd6133fd565b348015610e9157600080fd5b50610501610ea0366004614cae565b61340a565b348015610eb157600080fd5b50610647610ec0366004614ad0565b613489565b348015610ed157600080fd5b50610501610ee0366004614cd3565b6134f2565b348015610ef157600080fd5b50610501610f00366004614872565b613859565b348015610f1157600080fd5b50610523610f20366004614ad0565b613924565b348015610f3157600080fd5b50610523610f403660046144dc565b60096020526000908152604090205481565b348015610f5e57600080fd5b5061052360145481565b348015610f7457600080fd5b50610501610f83366004614629565b613955565b348015610f9457600080fd5b506105bd610fa3366004614d3b565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b34801561104757600080fd5b506105236110563660046144dc565b60066020526000908152604090205481565b34801561107457600080fd5b506105236110833660046144dc565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b392169061117b908c1688614dc3565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af11580156111f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061121d9190614e07565b5060008860ff1667ffffffffffffffff81111561123c5761123c61424f565b604051908082528060200260200182016040528015611265578160200160208202803683370190505b50905060005b8960ff168160ff1610156112d857600061128484613969565b905080838360ff168151811061129c5761129c614e22565b602090810291909101810191909152600091825260088152604080832088905560079091529020849055806112d081614e51565b91505061126b565b50505050505050505050565b600d602052826000526040600020602052816000526040600020818154811061130c57600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e69061137a9085908590600401614e70565b600060405180830381600087803b15801561139457600080fd5b505af11580156113a8573d6000803e3d6000fd5b505050505050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa158015611421573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114679190810190614ed6565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa15801561150c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115529190810190614ed6565b6040518363ffffffff1660e01b815260040161156f929190614e70565b600060405180830381600087803b15801561158957600080fd5b505af115801561159d573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e010439906024015b602060405180830381865afa15801561161b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114679190614f16565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156116a357602002820191906000526020600020905b81548152602001906001019080831161168f575b505050505090506116b5818251613a37565b92509250505b9250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b651690602401611404565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa1580156117b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117d49190614f3e565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015611877573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061189b9190614f6c565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611ba257600089898381811061190457611904614e22565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161193d91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401611969929190614e70565b600060405180830381600087803b15801561198357600080fd5b505af1158015611997573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015611a0d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a319190614f89565b90508060ff16600103611b8d576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611aba573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b009190810190614ed6565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611b599086908590600401614e70565b600060405180830381600087803b158015611b7357600080fd5b505af1158015611b87573d6000803e3d6000fd5b50505050505b50508080611b9a90614fa6565b9150506118e8565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611c0e93830182828015611c0257602002820191906000526020600020905b815481526020019060010190808311611bee575b50505050508484613abc565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611cd7573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114679190810190615001565b8060005b818160ff161015611df85760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611d5f57611d5f614e22565b905060200201356040518263ffffffff1660e01b8152600401611d8491815260200190565b600060405180830381600087803b158015611d9e57600080fd5b505af1158015611db2573d6000803e3d6000fd5b50505050611de584848360ff16818110611dce57611dce614e22565b90506020020135600f613c1b90919063ffffffff16565b5080611df081614e51565b915050611d21565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611e5f576000858152600d6020908152604080832061ffff85168452909152902054611e4b9083615120565b915080611e5781615133565b915050611e14565b509392505050565b60005a9050600080611e7b84860186614994565b91509150600081806020019051810190611e959190615154565b60008181526005602090815260408083205460049092528220549293509190611ebc613c27565b905082600003611edc576000848152600560205260409020819055612037565b600084815260036020526040812054611ef5848461516d565b611eff919061516d565b6000868152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611f7157602002820191906000526020600020905b815481526020019060010190808311611f5d575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff16815103611fec5781611fae81615133565b6000898152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000868152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055888352600c8252928220805493840181558252902001555b600084815260066020526040812054612051906001615120565b600086815260066020908152604080832084905560049091529020839055905061207b85836129f2565b6112d8858984613067565b6000828152600d6020908152604080832061ffff851684528252918290208054835181840281018401909452808452606093928301828280156120e857602002820191906000526020600020905b8154815260200190600101908083116120d4575b5050505050905092915050565b600060606000848460405160200161210e929190615180565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa1580156121b9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114679190614f89565b8160005b8181101561227a5730635f17e61686868481811061220157612201614e22565b90506020020135856040518363ffffffff1660e01b815260040161223592919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561224f57600080fd5b505af1158015612263573d6000803e3d6000fd5b50505050808061227290614fa6565b9150506121e1565b5050505050565b612289613cc9565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156122f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061231c9190615154565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af1158015612394573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123b89190614e07565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c90915281206123f49161414e565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff1611612450576000848152600d6020908152604080832061ffff85168452909152812061243e9161414e565b8061244881615133565b915050612409565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000606060005a905060006124a0858701876144dc565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff8111156124d9576124d961424f565b6040519080825280601f01601f191660200182016040528015612503576020820181803683370190505b50604051602001612515929190614e70565b60405160208183030381529060405290506000612530613c27565b9050600061253d86612636565b90505b835a61254c908961516d565b61255890612710615120565b10156125995781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612540565b806125b15760008398509850505050505050506116bb565b60176016601884896040516020016125cb91815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e00000000000000000000000000000000000000000000000000000000825261262d9594939291600401615302565b60405180910390fd5b600081815260056020526040812054810361265357506001919050565b600082815260036020908152604080832054600490925290912054612676613c27565b612680919061516d565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612709576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161262d565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa1580156127f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114679190614f6c565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b590612873908690869086906004016153c5565b600060405180830381600087803b15801561288d57600080fd5b505af11580156128a1573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d3590612873908690869086906004016153c5565b6017612910838261545f565b50601861291d828261545f565b505050565b80516123b890601690602084019061416c565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa1580156129ac573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611c119190810190615579565b601454600083815260026020526040902054612a0e908361516d565b11156123b8576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612a84573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612aca9190810190615001565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612b3f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b639190614f16565b601554909150612b879082906c01000000000000000000000000900460ff16614dc3565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611df857601554612bca9085906bffffffffffffffffffffffff16612c37565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612cbf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ce39190614e07565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f79060440161137a565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015612dae573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612df49190810190615579565b80519091506000612e03613c27565b905060005b8281101561227a576000848281518110612e2457612e24614e22565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612ea4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ec89190614f89565b90508060ff16600103612f44578660ff16600003612f14576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4612f44565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b50508080612f5190614fa6565b915050612e08565b6000818152600c6020908152604091829020805483518184028101840190945280845260609392830182828015612faf57602002820191906000526020600020905b815481526020019060010190808311612f9b575b50505050509050919050565b60168181548110612fcb57600080fd5b906000526020600020016000915090508054612fe690615214565b80601f016020809104026020016040519081016040528092919081815260200182805461301290615214565b801561305f5780601f106130345761010080835404028352916020019161305f565b820191906000526020600020905b81548152906001019060200180831161304257829003601f168201915b505050505081565b6000838152600760205260409020545b805a613083908561516d565b61308f90612710615120565b1015611df85781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055613077565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b15801561314857600080fd5b505af115801561315c573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b1580156131ed57600080fd5b505af115801561227a573d6000803e3d6000fd5b60178054612fe690615214565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff169063b657bc9c906024016115fe565b8060005b818163ffffffff161015611df8573063af953a4a858563ffffffff851681811061329657613296614e22565b905060200201356040518263ffffffff1660e01b81526004016132bb91815260200190565b600060405180830381600087803b1580156132d557600080fd5b505af11580156132e9573d6000803e3d6000fd5b5050505080806132f89061560a565b91505061326a565b6060600061330e600f613d4c565b9050808410613349576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8260000361335e5761335b848261516d565b92505b60008367ffffffffffffffff8111156133795761337961424f565b6040519080825280602002602001820160405280156133a2578160200160208202803683370190505b50905060005b848110156133f4576133c56133bd8288615120565b600f90613d56565b8282815181106133d7576133d7614e22565b6020908102919091010152806133ec81614fa6565b9150506133a8565b50949350505050565b60188054612fe690615214565b6000613414613c27565b90508160ff16600003613455576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a4505050565b6000828152600c602090815260408083208054825181850281018501909352808352849384939291908301828280156134e157602002820191906000526020600020905b8154815260200190600101908083116134cd575b505050505090506116b58185613a37565b8260005b818110156113a857600086868381811061351257613512614e22565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161354b91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613577929190614e70565b600060405180830381600087803b15801561359157600080fd5b505af11580156135a5573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa15801561361b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061363f9190614f89565b90508060ff16600103613844577f000000000000000000000000000000000000000000000000000000000000000060ff87161561369957507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb308985886040516020016136cd91815260200190565b6040516020818303038152906040526136e590615623565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa158015613770573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526137b69190810190614ed6565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d359061380f9087908590600401614e70565b600060405180830381600087803b15801561382957600080fd5b505af115801561383d573d6000803e3d6000fd5b5050505050505b5050808061385190614fa6565b9150506134f6565b8060005b81811015611df857600084848381811061387957613879614e22565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016138b291815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016138de929190614e70565b600060405180830381600087803b1580156138f857600080fd5b505af115801561390c573d6000803e3d6000fd5b5050505050808061391c90614fa6565b91505061385d565b600c602052816000526040600020818154811061394057600080fd5b90600052602060002001600091509150505481565b61395d613cc9565b61396681613d62565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e11906139c4908690600401615665565b6020604051808303816000875af11580156139e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a079190615154565b9050613a14600f82613e57565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b815160009081908190841580613a4d5750808510155b15613a56578094505b60008092505b85831015613ab257866001613a71858561516d565b613a7b919061516d565b81518110613a8b57613a8b614e22565b602002602001015181613a9e9190615120565b905082613aaa81614fa6565b935050613a5c565b9694955050505050565b82516000908190831580613ad05750808410155b15613ad9578093505b60008467ffffffffffffffff811115613af457613af461424f565b604051908082528060200260200182016040528015613b1d578160200160208202803683370190505b509050600092505b84831015613b8b57866001613b3a858561516d565b613b44919061516d565b81518110613b5457613b54614e22565b6020026020010151818481518110613b6e57613b6e614e22565b602090810291909101015282613b8381614fa6565b935050613b25565b613ba481600060018451613b9f919061516d565b613e63565b85606403613bdd578060018251613bbb919061516d565b81518110613bcb57613bcb614e22565b60200260200101519350505050611c11565b806064825188613bed91906157b7565b613bf79190615823565b81518110613c0757613c07614e22565b602002602001015193505050509392505050565b6000611c118383613fdb565b60007f000000000000000000000000000000000000000000000000000000000000000015613cc457606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613c9b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613cbf9190615154565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613d4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161262d565b565b6000611467825490565b6000611c1183836140d5565b3373ffffffffffffffffffffffffffffffffffffffff821603613de1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161262d565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611c1183836140ff565b8181808203613e73575050505050565b6000856002613e828787615837565b613e8c9190615857565b613e9690876158bf565b81518110613ea657613ea6614e22565b602002602001015190505b818313613fb5575b80868481518110613ecc57613ecc614e22565b60200260200101511015613eec5782613ee4816158e7565b935050613eb9565b858281518110613efe57613efe614e22565b6020026020010151811015613f1f5781613f1781615918565b925050613eec565b818313613fb057858281518110613f3857613f38614e22565b6020026020010151868481518110613f5257613f52614e22565b6020026020010151878581518110613f6c57613f6c614e22565b60200260200101888581518110613f8557613f85614e22565b60209081029190910101919091525282613f9e816158e7565b9350508180613fac90615918565b9250505b613eb1565b81851215613fc857613fc8868684613e63565b838312156113a8576113a8868486613e63565b600081815260018301602052604081205480156140c4576000613fff60018361516d565b85549091506000906140139060019061516d565b905081811461407857600086600001828154811061403357614033614e22565b906000526020600020015490508087600001848154811061405657614056614e22565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806140895761408961596f565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611467565b6000915050611467565b5092915050565b60008260000182815481106140ec576140ec614e22565b9060005260206000200154905092915050565b600081815260018301602052604081205461414657508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611467565b506000611467565b508054600082559060005260206000209081019061396691906141c2565b8280548282559060005260206000209081019282156141b2579160200282015b828111156141b257825182906141a2908261545f565b509160200191906001019061418c565b506141be9291506141d7565b5090565b5b808211156141be57600081556001016141c3565b808211156141be5760006141eb82826141f4565b506001016141d7565b50805461420090615214565b6000825580601f10614210575050565b601f01602090049060005260206000209081019061396691906141c2565b60ff8116811461396657600080fd5b63ffffffff8116811461396657600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff811182821017156142a2576142a261424f565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156142ef576142ef61424f565b604052919050565b600067ffffffffffffffff8211156143115761431161424f565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261434e57600080fd5b813561436161435c826142f7565b6142a8565b81815284602083860101111561437657600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff8116811461396657600080fd5b600080600080600080600060e0888a0312156143c857600080fd5b87356143d38161422e565b965060208801356143e38161423d565b955060408801356143f38161422e565b9450606088013567ffffffffffffffff81111561440f57600080fd5b61441b8a828b0161433d565b945050608088013561442c81614393565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff8116811461445b57600080fd5b919050565b60008060006060848603121561447557600080fd5b8335925061448560208501614449565b9150604084013590509250925092565b600080604083850312156144a857600080fd5b82359150602083013567ffffffffffffffff8111156144c657600080fd5b6144d28582860161433d565b9150509250929050565b6000602082840312156144ee57600080fd5b5035919050565b60005b838110156145105781810151838201526020016144f8565b50506000910152565b600081518084526145318160208601602086016144f5565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611c116020830184614519565b73ffffffffffffffffffffffffffffffffffffffff8116811461396657600080fd5b600080600080600080600060e0888a0312156145b357600080fd5b8735965060208801356145c581614576565b955060408801356145d58161422e565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b6000806040838503121561461057600080fd5b8235915061462060208401614449565b90509250929050565b60006020828403121561463b57600080fd5b8135611c1181614576565b60008083601f84011261465857600080fd5b50813567ffffffffffffffff81111561467057600080fd5b6020830191508360208260051b85010111156116bb57600080fd5b600080600080600080600060c0888a0312156146a657600080fd5b873567ffffffffffffffff8111156146bd57600080fd5b6146c98a828b01614646565b90985096505060208801356146dd8161422e565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b60008060006060848603121561471a57600080fd5b505081359360208301359350604090920135919050565b6020815261475860208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151614771604084018263ffffffff169052565b50604083015161014080606085015261478e610160850183614519565b915060608501516147af60808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e085015161010061481b818701836bffffffffffffffffffffffff169052565b86015190506101206148308682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018387015290506148688382614519565b9695505050505050565b6000806020838503121561488557600080fd5b823567ffffffffffffffff81111561489c57600080fd5b6148a885828601614646565b90969095509350505050565b60008083601f8401126148c657600080fd5b50813567ffffffffffffffff8111156148de57600080fd5b6020830191508360208285010111156116bb57600080fd5b6000806020838503121561490957600080fd5b823567ffffffffffffffff81111561492057600080fd5b6148a8858286016148b4565b6020808252825182820181905260009190848201906040850190845b8181101561496457835183529284019291840191600101614948565b50909695505050505050565b600067ffffffffffffffff82111561498a5761498a61424f565b5060051b60200190565b600080604083850312156149a757600080fd5b823567ffffffffffffffff808211156149bf57600080fd5b818501915085601f8301126149d357600080fd5b813560206149e361435c83614970565b82815260059290921b84018101918181019089841115614a0257600080fd5b8286015b84811015614a3a57803586811115614a1e5760008081fd5b614a2c8c86838b010161433d565b845250918301918301614a06565b5096505086013592505080821115614a5157600080fd5b506144d28582860161433d565b8215158152604060208201526000611c0e6040830184614519565b600080600060408486031215614a8e57600080fd5b833567ffffffffffffffff811115614aa557600080fd5b614ab186828701614646565b9094509250506020840135614ac58161423d565b809150509250925092565b60008060408385031215614ae357600080fd5b50508035926020909101359150565b600080600060408486031215614b0757600080fd5b83359250602084013567ffffffffffffffff811115614b2557600080fd5b614b31868287016148b4565b9497909650939450505050565b60008060408385031215614b5157600080fd5b823567ffffffffffffffff80821115614b6957600080fd5b614b758683870161433d565b93506020850135915080821115614a5157600080fd5b60006020808385031215614b9e57600080fd5b823567ffffffffffffffff80821115614bb657600080fd5b818501915085601f830112614bca57600080fd5b8135614bd861435c82614970565b81815260059190911b83018401908481019088831115614bf757600080fd5b8585015b83811015614c2f57803585811115614c135760008081fd5b614c218b89838a010161433d565b845250918601918601614bfb565b5098975050505050505050565b60008060408385031215614c4f57600080fd5b823591506020830135614c6181614393565b809150509250929050565b600060208284031215614c7e57600080fd5b8135611c118161422e565b60008060408385031215614c9c57600080fd5b823591506020830135614c618161423d565b60008060408385031215614cc157600080fd5b823591506020830135614c618161422e565b60008060008060608587031215614ce957600080fd5b843567ffffffffffffffff811115614d0057600080fd5b614d0c87828801614646565b9095509350506020850135614d208161422e565b91506040850135614d308161422e565b939692955090935050565b60008060008060008060c08789031215614d5457600080fd5b8635614d5f81614576565b95506020870135614d6f8161422e565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614dee57614dee614d94565b02949350505050565b8051801515811461445b57600080fd5b600060208284031215614e1957600080fd5b611c1182614df7565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614e6757614e67614d94565b60010192915050565b828152604060208201526000611c0e6040830184614519565b600082601f830112614e9a57600080fd5b8151614ea861435c826142f7565b818152846020838601011115614ebd57600080fd5b614ece8260208301602087016144f5565b949350505050565b600060208284031215614ee857600080fd5b815167ffffffffffffffff811115614eff57600080fd5b614ece84828501614e89565b805161445b81614393565b600060208284031215614f2857600080fd5b8151611c1181614393565b805161445b81614576565b60008060408385031215614f5157600080fd5b8251614f5c81614576565b6020939093015192949293505050565b600060208284031215614f7e57600080fd5b8151611c1181614576565b600060208284031215614f9b57600080fd5b8151611c118161422e565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614fd757614fd7614d94565b5060010190565b805161445b8161423d565b805167ffffffffffffffff8116811461445b57600080fd5b60006020828403121561501357600080fd5b815167ffffffffffffffff8082111561502b57600080fd5b90830190610140828603121561504057600080fd5b61504861427e565b61505183614f33565b815261505f60208401614fde565b602082015260408301518281111561507657600080fd5b61508287828601614e89565b60408301525061509460608401614f0b565b60608201526150a560808401614f33565b60808201526150b660a08401614fe9565b60a08201526150c760c08401614fde565b60c08201526150d860e08401614f0b565b60e08201526101006150eb818501614df7565b90820152610120838101518381111561510357600080fd5b61510f88828701614e89565b918301919091525095945050505050565b8082018082111561146757611467614d94565b600061ffff80831681810361514a5761514a614d94565b6001019392505050565b60006020828403121561516657600080fd5b5051919050565b8181038181111561146757611467614d94565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b838110156151f5577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08887030185526151e3868351614519565b955093820193908201906001016151a9565b50508584038187015250505061520b8185614519565b95945050505050565b600181811c9082168061522857607f821691505b602082108103615261577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000815461527481615214565b80855260206001838116801561529157600181146152c9576152f7565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b89010195506152f7565b866000528260002060005b858110156152ef5781548a82018601529083019084016152d4565b890184019650505b505050505092915050565b60a08152600061531560a0830188615267565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b83811015615387577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08784030185526153758383615267565b9486019492506001918201910161533c565b5050868103604088015261539b818b615267565b94505050505084606084015282810360808401526153b98185614519565b98975050505050505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f82111561291d57600081815260208120601f850160051c810160208610156154405750805b601f850160051c820191505b818110156113a85782815560010161544c565b815167ffffffffffffffff8111156154795761547961424f565b61548d816154878454615214565b84615419565b602080601f8311600181146154e057600084156154aa5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556113a8565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561552d5788860151825594840194600190910190840161550e565b508582101561556957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000602080838503121561558c57600080fd5b825167ffffffffffffffff8111156155a357600080fd5b8301601f810185136155b457600080fd5b80516155c261435c82614970565b81815260059190911b820183019083810190878311156155e157600080fd5b928401925b828410156155ff578351825292840192908401906155e6565b979650505050505050565b600063ffffffff80831681810361514a5761514a614d94565b80516020808301519190811015615261577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b6020815260008251610140806020850152615684610160850183614519565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0808685030160408701526156c08483614519565b9350604087015191506156eb606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e087015261574c8483614519565b935060e0870151915061010081878603018188015261576b8584614519565b94508088015192505061012081878603018188015261578a8584614519565b945080880151925050506157ad828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156157ef576157ef614d94565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082615832576158326157f4565b500490565b81810360008312801583831316838312821617156140ce576140ce614d94565b600082615866576158666157f4565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f8000000000000000000000000000000000000000000000000000000000000000831416156158ba576158ba614d94565b500590565b80820182811260008312801582168215821617156158df576158df614d94565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614fd757614fd7614d94565b60007f8000000000000000000000000000000000000000000000000000000000000000820361594957615949614d94565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var VerifiableLoadMercuryUpkeepABI = VerifiableLoadMercuryUpkeepMetaData.ABI

var VerifiableLoadMercuryUpkeepBin = VerifiableLoadMercuryUpkeepMetaData.Bin

func DeployVerifiableLoadMercuryUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _registrar common.Address, _useArb bool) (common.Address, *types.Transaction, *VerifiableLoadMercuryUpkeep, error) {
	parsed, err := VerifiableLoadMercuryUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadMercuryUpkeepBin), backend, _registrar, _useArb)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) EmittedAgainSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "emittedAgainSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) EmittedSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "emittedSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.EmittedSig(&_VerifiableLoadMercuryUpkeep.CallOpts)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.EmittedSig(&_VerifiableLoadMercuryUpkeep.CallOpts)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetActiveUpkeepIDsDeployedByThisContract(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getActiveUpkeepIDsDeployedByThisContract", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadMercuryUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadMercuryUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetAllActiveUpkeepIDsOnRegistry(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getAllActiveUpkeepIDsOnRegistry", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadMercuryUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadMercuryUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetBalance(&_VerifiableLoadMercuryUpkeep.CallOpts, id)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetBalance(&_VerifiableLoadMercuryUpkeep.CallOpts, id)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetForwarder(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetForwarder(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetLogTriggerConfig(opts *bind.CallOpts, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getLogTriggerConfig", addr, selector, topic0, topic1, topic2, topic3)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getMinBalanceForUpkeep", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetTriggerType(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetTriggerType(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetUpkeepInfo(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getUpkeepInfo", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21UpkeepInfo)).(*KeeperRegistryBase21UpkeepInfo)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadMercuryUpkeep.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadMercuryUpkeep.CallOpts, upkeepId)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchPreparingUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchPreparingUpkeeps", upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchPreparingUpkeepsSimple(opts *bind.TransactOpts, upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchPreparingUpkeepsSimple", upkeepIds, log, selector)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds, log, selector)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepIds, log, selector)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BatchSendLogs(opts *bind.TransactOpts, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "batchSendLogs", log)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchSendLogs(&_VerifiableLoadMercuryUpkeep.TransactOpts, log)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BatchSendLogs(&_VerifiableLoadMercuryUpkeep.TransactOpts, log)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) BurnPerformGas(opts *bind.TransactOpts, upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "burnPerformGas", upkeepId, startGas, blockNum)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BurnPerformGas(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.BurnPerformGas(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, startGas, blockNum)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SendLog(opts *bind.TransactOpts, upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "sendLog", upkeepId, log)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SendLog(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, log)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SendLog(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, log)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setFeeds", _feeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeeds(&_VerifiableLoadMercuryUpkeep.TransactOpts, _feeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeeds(&_VerifiableLoadMercuryUpkeep.TransactOpts, _feeds)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setParamKeys", _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetParamKeys(&_VerifiableLoadMercuryUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetParamKeys(&_VerifiableLoadMercuryUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setUpkeepGasLimit", upkeepId, gasLimit)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, cfg)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) TopUpFund(opts *bind.TransactOpts, upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "topUpFund", upkeepId, blockNum)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TopUpFund(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.TopUpFund(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, blockNum)
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) UpdateLogTriggerConfig1(opts *bind.TransactOpts, upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "updateLogTriggerConfig1", upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) UpdateLogTriggerConfig2(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "updateLogTriggerConfig2", upkeepId, cfg)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadMercuryUpkeep.TransactOpts, upkeepId, cfg)
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
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadMercuryUpkeepLogEmittedIterator, error) {

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

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepLogEmittedIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "LogEmitted", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
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

type VerifiableLoadMercuryUpkeepLogEmittedAgainIterator struct {
	Event *VerifiableLoadMercuryUpkeepLogEmittedAgain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadMercuryUpkeepLogEmittedAgainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadMercuryUpkeepLogEmittedAgain)
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
		it.Event = new(VerifiableLoadMercuryUpkeepLogEmittedAgain)
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

func (it *VerifiableLoadMercuryUpkeepLogEmittedAgainIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadMercuryUpkeepLogEmittedAgainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadMercuryUpkeepLogEmittedAgain struct {
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadMercuryUpkeepLogEmittedAgainIterator, error) {

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

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.FilterLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadMercuryUpkeepLogEmittedAgainIterator{contract: _VerifiableLoadMercuryUpkeep.contract, event: "LogEmittedAgain", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VerifiableLoadMercuryUpkeep.contract.WatchLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadMercuryUpkeepLogEmittedAgain)
				if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepFilterer) ParseLogEmittedAgain(log types.Log) (*VerifiableLoadMercuryUpkeepLogEmittedAgain, error) {
	event := new(VerifiableLoadMercuryUpkeepLogEmittedAgain)
	if err := _VerifiableLoadMercuryUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifiableLoadMercuryUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["LogEmittedAgain"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseLogEmittedAgain(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadMercuryUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadMercuryUpkeep.ParseUpkeepTopUp(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadMercuryUpkeepLogEmitted) Topic() common.Hash {
	return common.HexToHash("0x97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf08")
}

func (VerifiableLoadMercuryUpkeepLogEmittedAgain) Topic() common.Hash {
	return common.HexToHash("0xc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d")
}

func (VerifiableLoadMercuryUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifiableLoadMercuryUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifiableLoadMercuryUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeep) Address() common.Address {
	return _VerifiableLoadMercuryUpkeep.address
}

type VerifiableLoadMercuryUpkeepInterface interface {
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

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadMercuryUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadMercuryUpkeepLogEmitted, error)

	FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadMercuryUpkeepLogEmittedAgainIterator, error)

	WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmittedAgain(log types.Log) (*VerifiableLoadMercuryUpkeepLogEmittedAgain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadMercuryUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadMercuryUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadMercuryUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadMercuryUpkeepOwnershipTransferred, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadMercuryUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadMercuryUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadMercuryUpkeepUpkeepTopUp, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
