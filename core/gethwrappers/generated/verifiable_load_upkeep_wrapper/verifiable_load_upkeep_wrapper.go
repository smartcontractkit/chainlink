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

var VerifiableLoadUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c060405260426101408181526101009182919062005c1961016039815260200160405180608001604052806042815260200162005c5b604291399052620000be906016906002620003c7565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000ee908262000543565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b602082015260189062000120908262000543565b503480156200012e57600080fd5b5060405162005c9d38038062005c9d833981016040819052620001519162000625565b81813380600081620001aa5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001dd57620001dd816200031c565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa1580156200023a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000260919062000668565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002c6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ec919062000699565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c05250620006c0915050565b336001600160a01b03821603620003765760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620001a1565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821562000412579160200282015b8281111562000412578251829062000401908262000543565b5091602001919060010190620003e8565b506200042092915062000424565b5090565b80821115620004205760006200043b828262000445565b5060010162000424565b5080546200045390620004b4565b6000825580601f1062000464575050565b601f01602090049060005260206000209081019062000484919062000487565b50565b5b8082111562000420576000815560010162000488565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004c957607f821691505b602082108103620004ea57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200053e57600081815260208120601f850160051c81016020861015620005195750805b601f850160051c820191505b818110156200053a5782815560010162000525565b5050505b505050565b81516001600160401b038111156200055f576200055f6200049e565b6200057781620005708454620004b4565b84620004f0565b602080601f831160018114620005af5760008415620005965750858301515b600019600386901b1c1916600185901b1785556200053a565b600085815260208120601f198616915b82811015620005e057888601518255948401946001909101908401620005bf565b5085821015620005ff5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200048457600080fd5b600080604083850312156200063957600080fd5b825162000646816200060f565b602084015190925080151581146200065d57600080fd5b809150509250929050565b600080604083850312156200067c57600080fd5b825162000689816200060f565b6020939093015192949293505050565b600060208284031215620006ac57600080fd5b8151620006b9816200060f565b9392505050565b60805160a05160c05160e051615503620007166000396000818161054d0152611f250152600081816109730152613a2301526000818161080c0152613471015260008181610d3e015261344601526155036000f3fe6080604052600436106104ba5760003560e01c806379ea994311610279578063a6b594751161015e578063d6051a72116100d6578063e45530831161008a578063fa333dfb1161006f578063fa333dfb14610f4e578063fba7ffa314611001578063fcdc1f631461102e57600080fd5b8063e455308314610f18578063f2fde38b14610f2e57600080fd5b8063daee1aeb116100bb578063daee1aeb14610eab578063dbef701e14610ecb578063e0114adb14610eeb57600080fd5b8063d6051a7214610e6b578063da6cba4714610e8b57600080fd5b8063b657bc9c1161012d578063c041982211610112578063c041982214610e16578063c98f10b014610e36578063d4c2490014610e4b57600080fd5b8063b657bc9c14610dd5578063becde0e114610df657600080fd5b8063a6b5947514610d60578063a72aa27e14610d80578063af953a4a14610da0578063afb28d1f14610dc057600080fd5b8063924ca578116101f15780639b429354116101c05780639d385eaa116101a55780639d385eaa14610cec5780639d6f1cc714610d0c578063a654824814610d2c57600080fd5b80639b42935414610c8e5780639b51fb0d14610cbb57600080fd5b8063924ca57814610c04578063948108f714610c2457806396cebc7c14610c445780639ac542eb14610c6457600080fd5b80638340507c11610248578063873c75861161022d578063873c758614610b8c5780638da5cb5b14610bac5780638fcb3fba14610bd757600080fd5b80638340507c14610b4c57806386e330af14610b6c57600080fd5b806379ea994314610abf5780637b10399914610adf5780637e7a46dc14610b0c5780638243444a14610b2c57600080fd5b806345d2ec171161039f578063636092e8116103175780637145f11b116102e657806376721303116102cb5780637672130314610a5d578063776898c814610a8a57806379ba509714610aaa57600080fd5b80637145f11b14610a0057806373644cce14610a3057600080fd5b8063636092e81461093c578063642f6cef1461096157806369cdbadb146109a55780636e04ff0d146109d257600080fd5b806351c98be31161036e5780635d4ee7f3116103535780635d4ee7f3146108da5780635f17e616146108ef57806360457ff51461090f57600080fd5b806351c98be31461088d57806357970e93146108ad57600080fd5b806345d2ec17146107cd57806346982093146107fa57806346e7a63e1461082e5780635147cd591461085b57600080fd5b806320e3dbd4116104325780632a9032d311610401578063328ffd11116103e6578063328ffd11146107605780633ebe8d6c1461078d5780634585e33b146107ad57600080fd5b80632a9032d3146106ee5780632b20e3971461070e57600080fd5b806320e3dbd4146106615780632636aecf1461068157806328c4b57b146106a157806329e0a841146106c157600080fd5b806319d97a94116104895780631e0104391161046e5780631e010439146105cf578063206c32e81461060c578063207b65161461064157600080fd5b806319d97a94146105825780631cdde251146105af57600080fd5b806306c1cc00146104c6578063077ac621146104e85780630b7d33e61461051b57806312c550271461053b57600080fd5b366104c157005b600080fd5b3480156104d257600080fd5b506104e66104e13660046141a5565b61105b565b005b3480156104f457600080fd5b50610508610503366004614258565b6112aa565b6040519081526020015b60405180910390f35b34801561052757600080fd5b506104e661053636600461428d565b6112e8565b34801561054757600080fd5b5061056f7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610512565b34801561058e57600080fd5b506105a261059d3660046142d4565b611376565b604051610512919061435b565b3480156105bb57600080fd5b506104e66105ca366004614390565b611433565b3480156105db57600080fd5b506105ef6105ea3660046142d4565b611570565b6040516bffffffffffffffffffffffff9091168152602001610512565b34801561061857600080fd5b5061062c6106273660046143f5565b611604565b60408051928352602083019190915201610512565b34801561064d57600080fd5b506105a261065c3660046142d4565b611687565b34801561066d57600080fd5b506104e661067c366004614421565b6116df565b34801561068d57600080fd5b506104e661069c366004614483565b6118a9565b3480156106ad57600080fd5b506105086106bc3660046144fd565b611b72565b3480156106cd57600080fd5b506106e16106dc3660046142d4565b611bdd565b6040516105129190614529565b3480156106fa57600080fd5b506104e661070936600461466a565b611ce2565b34801561071a57600080fd5b5060115461073b9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610512565b34801561076c57600080fd5b5061050861077b3660046142d4565b60036020526000908152604090205481565b34801561079957600080fd5b506105086107a83660046142d4565b611dc3565b3480156107b957600080fd5b506104e66107c83660046146ee565b611e2c565b3480156107d957600080fd5b506107ed6107e83660046143f5565b61203b565b6040516105129190614724565b34801561080657600080fd5b506105087f000000000000000000000000000000000000000000000000000000000000000081565b34801561083a57600080fd5b506105086108493660046142d4565b600a6020526000908152604090205481565b34801561086757600080fd5b5061087b6108763660046142d4565b6120aa565b60405160ff9091168152602001610512565b34801561089957600080fd5b506104e66108a8366004614768565b61213e565b3480156108b957600080fd5b5060125461073b9073ffffffffffffffffffffffffffffffffffffffff1681565b3480156108e657600080fd5b506104e66121e2565b3480156108fb57600080fd5b506104e661090a3660046147bf565b61231d565b34801561091b57600080fd5b5061050861092a3660046142d4565b60076020526000908152604090205481565b34801561094857600080fd5b506015546105ef906bffffffffffffffffffffffff1681565b34801561096d57600080fd5b506109957f000000000000000000000000000000000000000000000000000000000000000081565b6040519015158152602001610512565b3480156109b157600080fd5b506105086109c03660046142d4565b60086020526000908152604090205481565b3480156109de57600080fd5b506109f26109ed3660046146ee565b6123ea565b6040516105129291906147e1565b348015610a0c57600080fd5b50610995610a1b3660046142d4565b600b6020526000908152604090205460ff1681565b348015610a3c57600080fd5b50610508610a4b3660046142d4565b6000908152600c602052604090205490565b348015610a6957600080fd5b50610508610a783660046142d4565b60046020526000908152604090205481565b348015610a9657600080fd5b50610995610aa53660046142d4565b612517565b348015610ab657600080fd5b506104e6612569565b348015610acb57600080fd5b5061073b610ada3660046142d4565b61266b565b348015610aeb57600080fd5b5060135461073b9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b1857600080fd5b506104e6610b273660046147fc565b6126ff565b348015610b3857600080fd5b506104e6610b473660046147fc565b612790565b348015610b5857600080fd5b506104e6610b67366004614848565b6127ea565b348015610b7857600080fd5b506104e6610b873660046148c6565b612808565b348015610b9857600080fd5b506107ed610ba73660046147bf565b61281b565b348015610bb857600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661073b565b348015610be357600080fd5b50610508610bf23660046142d4565b60056020526000908152604090205481565b348015610c1057600080fd5b506104e6610c1f3660046147bf565b6128d8565b348015610c3057600080fd5b506104e6610c3f366004614977565b612a87565b348015610c5057600080fd5b506104e6610c5f3660046149a7565b612b9f565b348015610c7057600080fd5b5060155461087b906c01000000000000000000000000900460ff1681565b348015610c9a57600080fd5b506104e6610ca93660046147bf565b60009182526009602052604090912055565b348015610cc757600080fd5b5061056f610cd63660046142d4565b600e6020526000908152604090205461ffff1681565b348015610cf857600080fd5b506107ed610d073660046142d4565b612da9565b348015610d1857600080fd5b506105a2610d273660046142d4565b612e0b565b348015610d3857600080fd5b506105087f000000000000000000000000000000000000000000000000000000000000000081565b348015610d6c57600080fd5b506104e6610d7b3660046144fd565b612eb7565b348015610d8c57600080fd5b506104e6610d9b3660046149c4565b612f20565b348015610dac57600080fd5b506104e6610dbb3660046142d4565b612fcb565b348015610dcc57600080fd5b506105a2613051565b348015610de157600080fd5b506105ef610df03660046142d4565b50600090565b348015610e0257600080fd5b506104e6610e1136600461466a565b61305e565b348015610e2257600080fd5b506107ed610e313660046147bf565b6130f8565b348015610e4257600080fd5b506105a26131f5565b348015610e5757600080fd5b506104e6610e663660046149e9565b613202565b348015610e7757600080fd5b5061062c610e863660046147bf565b613281565b348015610e9757600080fd5b506104e6610ea6366004614a0e565b6132ea565b348015610eb757600080fd5b506104e6610ec636600461466a565b613651565b348015610ed757600080fd5b50610508610ee63660046147bf565b61371c565b348015610ef757600080fd5b50610508610f063660046142d4565b60096020526000908152604090205481565b348015610f2457600080fd5b5061050860145481565b348015610f3a57600080fd5b506104e6610f49366004614421565b61374d565b348015610f5a57600080fd5b506105a2610f69366004614a76565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b34801561100d57600080fd5b5061050861101c3660046142d4565b60066020526000908152604090205481565b34801561103a57600080fd5b506105086110493660046142d4565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b3921690611141908c1688614afe565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af11580156111bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111e39190614b42565b5060008860ff1667ffffffffffffffff81111561120257611202614047565b60405190808252806020026020018201604052801561122b578160200160208202803683370190505b50905060005b8960ff168160ff16101561129e57600061124a84613761565b905080838360ff168151811061126257611262614b5d565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061129681614b8c565b915050611231565b50505050505050505050565b600d60205282600052604060002060205281600052604060002081815481106112d257600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e6906113409085908590600401614bab565b600060405180830381600087803b15801561135a57600080fd5b505af115801561136e573d6000803e3d6000fd5b505050505050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa1580156113e7573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261142d9190810190614c11565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa1580156114d2573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115189190810190614c11565b6040518363ffffffff1660e01b8152600401611535929190614bab565b600060405180830381600087803b15801561154f57600080fd5b505af1158015611563573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e01043990602401602060405180830381865afa1580156115e0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061142d9190614c51565b6000828152600d6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493849392919083018282801561166857602002820191906000526020600020905b815481526020019060010190808311611654575b5050505050905061167a81825161382f565b92509250505b9250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b6516906024016113ca565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611775573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117999190614c79565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801561183c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118609190614ca7565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611b675760008989838181106118c9576118c9614b5d565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161190291815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b815260040161192e929190614bab565b600060405180830381600087803b15801561194857600080fd5b505af115801561195c573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa1580156119d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119f69190614cc4565b90508060ff16600103611b52576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611a7f573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611ac59190810190614c11565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611b1e9086908590600401614bab565b600060405180830381600087803b158015611b3857600080fd5b505af1158015611b4c573d6000803e3d6000fd5b50505050505b50508080611b5f90614ce1565b9150506118ad565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611bd393830182828015611bc757602002820191906000526020600020905b815481526020019060010190808311611bb3575b505050505084846138b4565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611c9c573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261142d9190810190614d3c565b8060005b818160ff161015611dbd5760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611d2457611d24614b5d565b905060200201356040518263ffffffff1660e01b8152600401611d4991815260200190565b600060405180830381600087803b158015611d6357600080fd5b505af1158015611d77573d6000803e3d6000fd5b50505050611daa84848360ff16818110611d9357611d93614b5d565b90506020020135600f613a1390919063ffffffff16565b5080611db581614b8c565b915050611ce6565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611e24576000858152600d6020908152604080832061ffff85168452909152902054611e109083614e5b565b915080611e1c81614e6e565b915050611dd9565b509392505050565b60005a90506000611e3f8385018561428d565b5060008181526005602090815260408083205460049092528220549293509190611e67613a1f565b905082600003611e87576000848152600560205260409020819055611fe2565b600084815260036020526040812054611ea08484614e8f565b611eaa9190614e8f565b6000868152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611f1c57602002820191906000526020600020905b815481526020019060010190808311611f08575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff16815103611f975781611f5981614e6e565b6000898152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000868152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055888352600c8252928220805493840181558252902001555b600084815260066020526040812054611ffc906001614e5b565b600086815260066020908152604080832084905560049091529020839055905061202685836128d8565b612031858784612eb7565b5050505050505050565b6000828152600d6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561209d57602002820191906000526020600020905b815481526020019060010190808311612089575b5050505050905092915050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa15801561211a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061142d9190614cc4565b8160005b818110156121db5730635f17e61686868481811061216257612162614b5d565b90506020020135856040518363ffffffff1660e01b815260040161219692919091825263ffffffff16602082015260400190565b600060405180830381600087803b1580156121b057600080fd5b505af11580156121c4573d6000803e3d6000fd5b5050505080806121d390614ce1565b915050612142565b5050505050565b6121ea613ac1565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015612259573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061227d9190614ea2565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156122f5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123199190614b42565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c909152812061235591613f46565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff16116123b1576000848152600d6020908152604080832061ffff85168452909152812061239f91613f46565b806123a981614e6e565b91505061236a565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000606060005a90506000612401858701876142d4565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff81111561243a5761243a614047565b6040519080825280601f01601f191660200182016040528015612464576020820181803683370190505b50604051602001612476929190614bab565b60405160208183030381529060405290506000612491613a1f565b9050600061249e86612517565b90505b835a6124ad9089614e8f565b6124b990612710614e5b565b10156125075781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055816124ff81614ebb565b9250506124a1565b9a91995090975050505050505050565b600081815260056020526040812054810361253457506001919050565b600082815260036020908152604080832054600490925290912054612557613a1f565b6125619190614e8f565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146125ef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa1580156126db573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061142d9190614ca7565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b59061275990869086908690600401614ef0565b600060405180830381600087803b15801561277357600080fd5b505af1158015612787573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d359061275990869086908690600401614ef0565b60176127f68382614fdd565b5060186128038282614fdd565b505050565b8051612319906016906020840190613f64565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612892573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611bd691908101906150f7565b6014546000838152600260205260409020546128f49083614e8f565b1115612319576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa15801561296a573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526129b09190810190614d3c565b6015549091506000906129d79082906c01000000000000000000000000900460ff16614afe565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611dbd57601554612a1a9085906bffffffffffffffffffffffff16612a87565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612b0f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b339190614b42565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401611340565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015612bfe573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612c4491908101906150f7565b80519091506000612c53613a1f565b905060005b828110156121db576000848281518110612c7457612c74614b5d565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612cf4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d189190614cc4565b90508060ff16600103612d94578660ff16600003612d64576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4612d94565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b50508080612da190614ce1565b915050612c58565b6000818152600c6020908152604091829020805483518184028101840190945280845260609392830182828015612dff57602002820191906000526020600020905b815481526020019060010190808311612deb575b50505050509050919050565b60168181548110612e1b57600080fd5b906000526020600020016000915090508054612e3690614f44565b80601f0160208091040260200160405190810160405280929190818152602001828054612e6290614f44565b8015612eaf5780601f10612e8457610100808354040283529160200191612eaf565b820191906000526020600020905b815481529060010190602001808311612e9257829003601f168201915b505050505081565b6000838152600760205260409020545b805a612ed39085614e8f565b612edf90612710614e5b565b1015611dbd5781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612ec7565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b158015612f9857600080fd5b505af1158015612fac573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561303d57600080fd5b505af11580156121db573d6000803e3d6000fd5b60178054612e3690614f44565b8060005b818163ffffffff161015611dbd573063af953a4a858563ffffffff851681811061308e5761308e614b5d565b905060200201356040518263ffffffff1660e01b81526004016130b391815260200190565b600060405180830381600087803b1580156130cd57600080fd5b505af11580156130e1573d6000803e3d6000fd5b5050505080806130f090615188565b915050613062565b60606000613106600f613b44565b9050808410613141576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82600003613156576131538482614e8f565b92505b60008367ffffffffffffffff81111561317157613171614047565b60405190808252806020026020018201604052801561319a578160200160208202803683370190505b50905060005b848110156131ec576131bd6131b58288614e5b565b600f90613b4e565b8282815181106131cf576131cf614b5d565b6020908102919091010152806131e481614ce1565b9150506131a0565b50949350505050565b60188054612e3690614f44565b600061320c613a1f565b90508160ff1660000361324d576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a4505050565b6000828152600c602090815260408083208054825181850281018501909352808352849384939291908301828280156132d957602002820191906000526020600020905b8154815260200190600101908083116132c5575b5050505050905061167a818561382f565b8260005b8181101561136e57600086868381811061330a5761330a614b5d565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161334391815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b815260040161336f929190614bab565b600060405180830381600087803b15801561338957600080fd5b505af115801561339d573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015613413573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134379190614cc4565b90508060ff1660010361363c577f000000000000000000000000000000000000000000000000000000000000000060ff87161561349157507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb308985886040516020016134c591815260200190565b6040516020818303038152906040526134dd906151a1565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa158015613568573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526135ae9190810190614c11565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d35906136079087908590600401614bab565b600060405180830381600087803b15801561362157600080fd5b505af1158015613635573d6000803e3d6000fd5b5050505050505b5050808061364990614ce1565b9150506132ee565b8060005b81811015611dbd57600084848381811061367157613671614b5d565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016136aa91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016136d6929190614bab565b600060405180830381600087803b1580156136f057600080fd5b505af1158015613704573d6000803e3d6000fd5b5050505050808061371490614ce1565b915050613655565b600c602052816000526040600020818154811061373857600080fd5b90600052602060002001600091509150505481565b613755613ac1565b61375e81613b5a565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e11906137bc9086906004016151e3565b6020604051808303816000875af11580156137db573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137ff9190614ea2565b905061380c600f82613c4f565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b8151600090819081908415806138455750808510155b1561384e578094505b60008092505b858310156138aa578660016138698585614e8f565b6138739190614e8f565b8151811061388357613883614b5d565b6020026020010151816138969190614e5b565b9050826138a281614ce1565b935050613854565b9694955050505050565b825160009081908315806138c85750808410155b156138d1578093505b60008467ffffffffffffffff8111156138ec576138ec614047565b604051908082528060200260200182016040528015613915578160200160208202803683370190505b509050600092505b84831015613983578660016139328585614e8f565b61393c9190614e8f565b8151811061394c5761394c614b5d565b602002602001015181848151811061396657613966614b5d565b60209081029190910101528261397b81614ce1565b93505061391d565b61399c816000600184516139979190614e8f565b613c5b565b856064036139d55780600182516139b39190614e8f565b815181106139c3576139c3614b5d565b60200260200101519350505050611bd6565b8060648251886139e59190615335565b6139ef91906153a1565b815181106139ff576139ff614b5d565b602002602001015193505050509392505050565b6000611bd68383613dd3565b60007f000000000000000000000000000000000000000000000000000000000000000015613abc57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613a93573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ab79190614ea2565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613b42576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016125e6565b565b600061142d825490565b6000611bd68383613ecd565b3373ffffffffffffffffffffffffffffffffffffffff821603613bd9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016125e6565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611bd68383613ef7565b8181808203613c6b575050505050565b6000856002613c7a87876153b5565b613c8491906153d5565b613c8e908761543d565b81518110613c9e57613c9e614b5d565b602002602001015190505b818313613dad575b80868481518110613cc457613cc4614b5d565b60200260200101511015613ce45782613cdc81615465565b935050613cb1565b858281518110613cf657613cf6614b5d565b6020026020010151811015613d175781613d0f81615496565b925050613ce4565b818313613da857858281518110613d3057613d30614b5d565b6020026020010151868481518110613d4a57613d4a614b5d565b6020026020010151878581518110613d6457613d64614b5d565b60200260200101888581518110613d7d57613d7d614b5d565b60209081029190910101919091525282613d9681615465565b9350508180613da490615496565b9250505b613ca9565b81851215613dc057613dc0868684613c5b565b8383121561136e5761136e868486613c5b565b60008181526001830160205260408120548015613ebc576000613df7600183614e8f565b8554909150600090613e0b90600190614e8f565b9050818114613e70576000866000018281548110613e2b57613e2b614b5d565b9060005260206000200154905080876000018481548110613e4e57613e4e614b5d565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613e8157613e816154c7565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061142d565b600091505061142d565b5092915050565b6000826000018281548110613ee457613ee4614b5d565b9060005260206000200154905092915050565b6000818152600183016020526040812054613f3e5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561142d565b50600061142d565b508054600082559060005260206000209081019061375e9190613fba565b828054828255906000526020600020908101928215613faa579160200282015b82811115613faa5782518290613f9a9082614fdd565b5091602001919060010190613f84565b50613fb6929150613fcf565b5090565b5b80821115613fb65760008155600101613fbb565b80821115613fb6576000613fe38282613fec565b50600101613fcf565b508054613ff890614f44565b6000825580601f10614008575050565b601f01602090049060005260206000209081019061375e9190613fba565b60ff8116811461375e57600080fd5b63ffffffff8116811461375e57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff8111828210171561409a5761409a614047565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156140e7576140e7614047565b604052919050565b600067ffffffffffffffff82111561410957614109614047565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261414657600080fd5b8135614159614154826140ef565b6140a0565b81815284602083860101111561416e57600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff8116811461375e57600080fd5b600080600080600080600060e0888a0312156141c057600080fd5b87356141cb81614026565b965060208801356141db81614035565b955060408801356141eb81614026565b9450606088013567ffffffffffffffff81111561420757600080fd5b6142138a828b01614135565b94505060808801356142248161418b565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff8116811461425357600080fd5b919050565b60008060006060848603121561426d57600080fd5b8335925061427d60208501614241565b9150604084013590509250925092565b600080604083850312156142a057600080fd5b82359150602083013567ffffffffffffffff8111156142be57600080fd5b6142ca85828601614135565b9150509250929050565b6000602082840312156142e657600080fd5b5035919050565b60005b838110156143085781810151838201526020016142f0565b50506000910152565b600081518084526143298160208601602086016142ed565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611bd66020830184614311565b73ffffffffffffffffffffffffffffffffffffffff8116811461375e57600080fd5b600080600080600080600060e0888a0312156143ab57600080fd5b8735965060208801356143bd8161436e565b955060408801356143cd81614026565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b6000806040838503121561440857600080fd5b8235915061441860208401614241565b90509250929050565b60006020828403121561443357600080fd5b8135611bd68161436e565b60008083601f84011261445057600080fd5b50813567ffffffffffffffff81111561446857600080fd5b6020830191508360208260051b850101111561168057600080fd5b600080600080600080600060c0888a03121561449e57600080fd5b873567ffffffffffffffff8111156144b557600080fd5b6144c18a828b0161443e565b90985096505060208801356144d581614026565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b60008060006060848603121561451257600080fd5b505081359360208301359350604090920135919050565b6020815261455060208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151614569604084018263ffffffff169052565b506040830151610140806060850152614586610160850183614311565b915060608501516145a760808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614613818701836bffffffffffffffffffffffff169052565b86015190506101206146288682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018387015290506146608382614311565b9695505050505050565b6000806020838503121561467d57600080fd5b823567ffffffffffffffff81111561469457600080fd5b6146a08582860161443e565b90969095509350505050565b60008083601f8401126146be57600080fd5b50813567ffffffffffffffff8111156146d657600080fd5b60208301915083602082850101111561168057600080fd5b6000806020838503121561470157600080fd5b823567ffffffffffffffff81111561471857600080fd5b6146a0858286016146ac565b6020808252825182820181905260009190848201906040850190845b8181101561475c57835183529284019291840191600101614740565b50909695505050505050565b60008060006040848603121561477d57600080fd5b833567ffffffffffffffff81111561479457600080fd5b6147a08682870161443e565b90945092505060208401356147b481614035565b809150509250925092565b600080604083850312156147d257600080fd5b50508035926020909101359150565b8215158152604060208201526000611bd36040830184614311565b60008060006040848603121561481157600080fd5b83359250602084013567ffffffffffffffff81111561482f57600080fd5b61483b868287016146ac565b9497909650939450505050565b6000806040838503121561485b57600080fd5b823567ffffffffffffffff8082111561487357600080fd5b61487f86838701614135565b9350602085013591508082111561489557600080fd5b506142ca85828601614135565b600067ffffffffffffffff8211156148bc576148bc614047565b5060051b60200190565b600060208083850312156148d957600080fd5b823567ffffffffffffffff808211156148f157600080fd5b818501915085601f83011261490557600080fd5b8135614913614154826148a2565b81815260059190911b8301840190848101908883111561493257600080fd5b8585015b8381101561496a5780358581111561494e5760008081fd5b61495c8b89838a0101614135565b845250918601918601614936565b5098975050505050505050565b6000806040838503121561498a57600080fd5b82359150602083013561499c8161418b565b809150509250929050565b6000602082840312156149b957600080fd5b8135611bd681614026565b600080604083850312156149d757600080fd5b82359150602083013561499c81614035565b600080604083850312156149fc57600080fd5b82359150602083013561499c81614026565b60008060008060608587031215614a2457600080fd5b843567ffffffffffffffff811115614a3b57600080fd5b614a478782880161443e565b9095509350506020850135614a5b81614026565b91506040850135614a6b81614026565b939692955090935050565b60008060008060008060c08789031215614a8f57600080fd5b8635614a9a8161436e565b95506020870135614aaa81614026565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614b2957614b29614acf565b02949350505050565b8051801515811461425357600080fd5b600060208284031215614b5457600080fd5b611bd682614b32565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614ba257614ba2614acf565b60010192915050565b828152604060208201526000611bd36040830184614311565b600082601f830112614bd557600080fd5b8151614be3614154826140ef565b818152846020838601011115614bf857600080fd5b614c098260208301602087016142ed565b949350505050565b600060208284031215614c2357600080fd5b815167ffffffffffffffff811115614c3a57600080fd5b614c0984828501614bc4565b80516142538161418b565b600060208284031215614c6357600080fd5b8151611bd68161418b565b80516142538161436e565b60008060408385031215614c8c57600080fd5b8251614c978161436e565b6020939093015192949293505050565b600060208284031215614cb957600080fd5b8151611bd68161436e565b600060208284031215614cd657600080fd5b8151611bd681614026565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614d1257614d12614acf565b5060010190565b805161425381614035565b805167ffffffffffffffff8116811461425357600080fd5b600060208284031215614d4e57600080fd5b815167ffffffffffffffff80821115614d6657600080fd5b908301906101408286031215614d7b57600080fd5b614d83614076565b614d8c83614c6e565b8152614d9a60208401614d19565b6020820152604083015182811115614db157600080fd5b614dbd87828601614bc4565b604083015250614dcf60608401614c46565b6060820152614de060808401614c6e565b6080820152614df160a08401614d24565b60a0820152614e0260c08401614d19565b60c0820152614e1360e08401614c46565b60e0820152610100614e26818501614b32565b908201526101208381015183811115614e3e57600080fd5b614e4a88828701614bc4565b918301919091525095945050505050565b8082018082111561142d5761142d614acf565b600061ffff808316818103614e8557614e85614acf565b6001019392505050565b8181038181111561142d5761142d614acf565b600060208284031215614eb457600080fd5b5051919050565b600081614eca57614eca614acf565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b600181811c90821680614f5857607f821691505b602082108103614f91577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561280357600081815260208120601f850160051c81016020861015614fbe5750805b601f850160051c820191505b8181101561136e57828155600101614fca565b815167ffffffffffffffff811115614ff757614ff7614047565b61500b816150058454614f44565b84614f97565b602080601f83116001811461505e57600084156150285750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561136e565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156150ab5788860151825594840194600190910190840161508c565b50858210156150e757878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b6000602080838503121561510a57600080fd5b825167ffffffffffffffff81111561512157600080fd5b8301601f8101851361513257600080fd5b8051615140614154826148a2565b81815260059190911b8201830190838101908783111561515f57600080fd5b928401925b8284101561517d57835182529284019290840190615164565b979650505050505050565b600063ffffffff808316818103614e8557614e85614acf565b80516020808301519190811015614f91577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b6020815260008251610140806020850152615202610160850183614311565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08086850301604087015261523e8483614311565b935060408701519150615269606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526152ca8483614311565b935060e087015191506101008187860301818801526152e98584614311565b9450808801519250506101208187860301818801526153088584614311565b9450808801519250505061532b828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561536d5761536d614acf565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826153b0576153b0615372565b500490565b8181036000831280158383131683831282161715613ec657613ec6614acf565b6000826153e4576153e4615372565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561543857615438614acf565b500590565b808201828112600083128015821682158216171561545d5761545d614acf565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614d1257614d12614acf565b60007f80000000000000000000000000000000000000000000000000000000000000008203614eca57614eca614acf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var VerifiableLoadUpkeepABI = VerifiableLoadUpkeepMetaData.ABI

var VerifiableLoadUpkeepBin = VerifiableLoadUpkeepMetaData.Bin

func DeployVerifiableLoadUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _registrar common.Address, _useArb bool) (common.Address, *types.Transaction, *VerifiableLoadUpkeep, error) {
	parsed, err := VerifiableLoadUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadUpkeepBin), backend, _registrar, _useArb)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifiableLoadUpkeep{address: address, abi: *parsed, VerifiableLoadUpkeepCaller: VerifiableLoadUpkeepCaller{contract: contract}, VerifiableLoadUpkeepTransactor: VerifiableLoadUpkeepTransactor{contract: contract}, VerifiableLoadUpkeepFilterer: VerifiableLoadUpkeepFilterer{contract: contract}}, nil
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) EmittedAgainSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "emittedAgainSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) EmittedSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "emittedSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadUpkeep.Contract.EmittedSig(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadUpkeep.Contract.EmittedSig(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) FeedParamKey() (string, error) {
	return _VerifiableLoadUpkeep.Contract.FeedParamKey(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) FeedParamKey() (string, error) {
	return _VerifiableLoadUpkeep.Contract.FeedParamKey(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadUpkeep.Contract.FeedsHex(&_VerifiableLoadUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadUpkeep.Contract.FeedsHex(&_VerifiableLoadUpkeep.CallOpts, arg0)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetActiveUpkeepIDsDeployedByThisContract(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getActiveUpkeepIDsDeployedByThisContract", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetAllActiveUpkeepIDsOnRegistry(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getAllActiveUpkeepIDsOnRegistry", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetBalance(&_VerifiableLoadUpkeep.CallOpts, id)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetBalance(&_VerifiableLoadUpkeep.CallOpts, id)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.GetForwarder(&_VerifiableLoadUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadUpkeep.Contract.GetForwarder(&_VerifiableLoadUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetLogTriggerConfig(opts *bind.CallOpts, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getLogTriggerConfig", addr, selector, topic0, topic1, topic2, topic3)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getMinBalanceForUpkeep", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadUpkeep.Contract.GetTriggerType(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadUpkeep.Contract.GetTriggerType(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetUpkeepInfo(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getUpkeepInfo", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21UpkeepInfo)).(*KeeperRegistryBase21UpkeepInfo)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadUpkeep.CallOpts, upkeepId)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadUpkeep.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) TimeParamKey() (string, error) {
	return _VerifiableLoadUpkeep.Contract.TimeParamKey(&_VerifiableLoadUpkeep.CallOpts)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepCallerSession) TimeParamKey() (string, error) {
	return _VerifiableLoadUpkeep.Contract.TimeParamKey(&_VerifiableLoadUpkeep.CallOpts)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchPreparingUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchPreparingUpkeeps", upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchPreparingUpkeepsSimple(opts *bind.TransactOpts, upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchPreparingUpkeepsSimple", upkeepIds, log, selector)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds, log, selector)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadUpkeep.TransactOpts, upkeepIds, log, selector)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BatchSendLogs(opts *bind.TransactOpts, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "batchSendLogs", log)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchSendLogs(&_VerifiableLoadUpkeep.TransactOpts, log)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BatchSendLogs(&_VerifiableLoadUpkeep.TransactOpts, log)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) BurnPerformGas(opts *bind.TransactOpts, upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "burnPerformGas", upkeepId, startGas, blockNum)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BurnPerformGas(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.BurnPerformGas(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, startGas, blockNum)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SendLog(opts *bind.TransactOpts, upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "sendLog", upkeepId, log)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SendLog(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, log)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SendLog(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, log)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setFeeds", _feeds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetFeeds(&_VerifiableLoadUpkeep.TransactOpts, _feeds)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetFeeds(&_VerifiableLoadUpkeep.TransactOpts, _feeds)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setParamKeys", _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetParamKeys(&_VerifiableLoadUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetParamKeys(&_VerifiableLoadUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setUpkeepGasLimit", upkeepId, gasLimit)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, cfg)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) TopUpFund(opts *bind.TransactOpts, upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "topUpFund", upkeepId, blockNum)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.TopUpFund(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.TopUpFund(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, blockNum)
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) UpdateLogTriggerConfig1(opts *bind.TransactOpts, upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "updateLogTriggerConfig1", upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) UpdateLogTriggerConfig2(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "updateLogTriggerConfig2", upkeepId, cfg)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, cfg)
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
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadUpkeepLogEmittedIterator, error) {

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

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepLogEmittedIterator{contract: _VerifiableLoadUpkeep.contract, event: "LogEmitted", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
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

type VerifiableLoadUpkeepLogEmittedAgainIterator struct {
	Event *VerifiableLoadUpkeepLogEmittedAgain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadUpkeepLogEmittedAgainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadUpkeepLogEmittedAgain)
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
		it.Event = new(VerifiableLoadUpkeepLogEmittedAgain)
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

func (it *VerifiableLoadUpkeepLogEmittedAgainIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadUpkeepLogEmittedAgainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadUpkeepLogEmittedAgain struct {
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadUpkeepLogEmittedAgainIterator, error) {

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

	logs, sub, err := _VerifiableLoadUpkeep.contract.FilterLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadUpkeepLogEmittedAgainIterator{contract: _VerifiableLoadUpkeep.contract, event: "LogEmittedAgain", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VerifiableLoadUpkeep.contract.WatchLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadUpkeepLogEmittedAgain)
				if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepFilterer) ParseLogEmittedAgain(log types.Log) (*VerifiableLoadUpkeepLogEmittedAgain, error) {
	event := new(VerifiableLoadUpkeepLogEmittedAgain)
	if err := _VerifiableLoadUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifiableLoadUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadUpkeep.abi.Events["LogEmittedAgain"].ID:
		return _VerifiableLoadUpkeep.ParseLogEmittedAgain(log)
	case _VerifiableLoadUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadUpkeep.ParseUpkeepTopUp(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadUpkeepLogEmitted) Topic() common.Hash {
	return common.HexToHash("0x97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf08")
}

func (VerifiableLoadUpkeepLogEmittedAgain) Topic() common.Hash {
	return common.HexToHash("0xc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d")
}

func (VerifiableLoadUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifiableLoadUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifiableLoadUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeep) Address() common.Address {
	return _VerifiableLoadUpkeep.address
}

type VerifiableLoadUpkeepInterface interface {
	BUCKETSIZE(opts *bind.CallOpts) (uint16, error)

	AddLinkAmount(opts *bind.CallOpts) (*big.Int, error)

	BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error)

	Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error)

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

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadUpkeepLogEmitted, error)

	FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadUpkeepLogEmittedAgainIterator, error)

	WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmittedAgain(log types.Log) (*VerifiableLoadUpkeepLogEmittedAgain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadUpkeepOwnershipTransferred, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadUpkeepUpkeepTopUp, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
