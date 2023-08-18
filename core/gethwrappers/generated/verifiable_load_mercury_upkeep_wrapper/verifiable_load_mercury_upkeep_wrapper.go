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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"FeedLookup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"newFeeds\",\"type\":\"string[]\"}],\"name\":\"setFeedsHex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c06040526042610140818152610100918291906200607c610160398152602001604051806080016040528060428152602001620060be604291399052620000be90601690600262000365565b50348015620000cc57600080fd5b506040516200610038038062006100833981016040819052620000ef9162000452565b81813380600081620001485760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200017b576200017b81620002ba565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa158015620001d8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001fe919062000495565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801562000264573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200028a9190620004c6565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c052506200065e915050565b336001600160a01b03821603620003145760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200013f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215620003b0579160200282015b82811115620003b057825182906200039f908262000592565b509160200191906001019062000386565b50620003be929150620003c2565b5090565b80821115620003be576000620003d98282620003e3565b50600101620003c2565b508054620003f19062000503565b6000825580601f1062000402575050565b601f01602090049060005260206000209081019062000422919062000425565b50565b5b80821115620003be576000815560010162000426565b6001600160a01b03811681146200042257600080fd5b600080604083850312156200046657600080fd5b825162000473816200043c565b602084015190925080151581146200048a57600080fd5b809150509250929050565b60008060408385031215620004a957600080fd5b8251620004b6816200043c565b6020939093015192949293505050565b600060208284031215620004d957600080fd5b8151620004e6816200043c565b9392505050565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200051857607f821691505b6020821081036200053957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200058d57600081815260208120601f850160051c81016020861015620005685750805b601f850160051c820191505b81811015620005895782815560010162000574565b5050505b505050565b81516001600160401b03811115620005ae57620005ae620004ed565b620005c681620005bf845462000503565b846200053f565b602080601f831160018114620005fe5760008415620005e55750858301515b600019600386901b1c1916600185901b17855562000589565b600085815260208120601f198616915b828110156200062f578886015182559484019460019091019084016200060e565b50858210156200064e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c05160e0516159c8620006b46000396000818161056d0152611fbe0152600081816109c10152613c8901526000818161082c01526136d7015260008181610d3e01526136ac01526159c86000f3fe6080604052600436106104ba5760003560e01c8063776898c811610279578063a6b594751161015e578063d6051a72116100d6578063e45530831161008a578063fa333dfb1161006f578063fa333dfb14610fb5578063fba7ffa314611068578063fcdc1f631461109557600080fd5b8063e455308314610f7f578063f2fde38b14610f9557600080fd5b8063daee1aeb116100bb578063daee1aeb14610f12578063dbef701e14610f32578063e0114adb14610f5257600080fd5b8063d6051a7214610ed2578063da6cba4714610ef257600080fd5b8063b657bc9c1161012d578063c041982211610112578063c041982214610e49578063c98f10b014610e69578063d4c2490014610eb257600080fd5b8063b657bc9c14610e09578063becde0e114610e2957600080fd5b8063a6b5947514610d60578063a72aa27e14610d80578063af953a4a14610da0578063afb28d1f14610dc057600080fd5b8063924ca578116101f15780639b429354116101c05780639d385eaa116101a55780639d385eaa14610cec5780639d6f1cc714610d0c578063a654824814610d2c57600080fd5b80639b42935414610c8e5780639b51fb0d14610cbb57600080fd5b8063924ca57814610c04578063948108f714610c2457806396cebc7c14610c445780639ac542eb14610c6457600080fd5b80637e7a46dc11610248578063873c75861161022d578063873c758614610b8c5780638da5cb5b14610bac5780638fcb3fba14610bd757600080fd5b80637e7a46dc14610b4c5780638243444a14610b6c57600080fd5b8063776898c814610aca57806379ba509714610aea57806379ea994314610aff5780637b10399914610b1f57600080fd5b80634585e33b1161039f5780635f17e6161161031757806369cdbadb116102e65780637145f11b116102cb5780637145f11b14610a4057806373644cce14610a705780637672130314610a9d57600080fd5b806369cdbadb146109f35780636e04ff0d14610a2057600080fd5b80635f17e6161461093d57806360457ff51461095d578063636092e81461098a578063642f6cef146109af57600080fd5b80634b56a42e1161036e57806351c98be31161035357806351c98be3146108db57806357970e93146108fb5780635d4ee7f31461092857600080fd5b80634b56a42e1461087b5780635147cd59146108a957600080fd5b80634585e33b146107cd57806345d2ec17146107ed578063469820931461081a57806346e7a63e1461084e57600080fd5b8063207b65161161043257806329e0a841116104015780632b20e397116103e65780632b20e3971461072e578063328ffd11146107805780633ebe8d6c146107ad57600080fd5b806329e0a841146106e15780632a9032d31461070e57600080fd5b8063207b65161461066157806320e3dbd4146106815780632636aecf146106a157806328c4b57b146106c157600080fd5b806312c55027116104895780631cdde2511161046e5780631cdde251146105cf5780631e010439146105ef578063206c32e81461062c57600080fd5b806312c550271461055b57806319d97a94146105a257600080fd5b806305e25131146104c657806306c1cc00146104e8578063077ac621146105085780630b7d33e61461053b57600080fd5b366104c157005b600080fd5b3480156104d257600080fd5b506104e66104e13660046143dc565b6110c2565b005b3480156104f457600080fd5b506104e66105033660046144fd565b6110d9565b34801561051457600080fd5b506105286105233660046145b0565b611328565b6040519081526020015b60405180910390f35b34801561054757600080fd5b506104e66105563660046145e5565b611366565b34801561056757600080fd5b5061058f7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff9091168152602001610532565b3480156105ae57600080fd5b506105c26105bd36600461462c565b6113f4565b60405161053291906146b3565b3480156105db57600080fd5b506104e66105ea3660046146e8565b6114b1565b3480156105fb57600080fd5b5061060f61060a36600461462c565b6115ee565b6040516bffffffffffffffffffffffff9091168152602001610532565b34801561063857600080fd5b5061064c61064736600461474d565b611683565b60408051928352602083019190915201610532565b34801561066d57600080fd5b506105c261067c36600461462c565b611706565b34801561068d57600080fd5b506104e661069c366004614779565b61175e565b3480156106ad57600080fd5b506104e66106bc3660046147db565b611928565b3480156106cd57600080fd5b506105286106dc366004614855565b611bf1565b3480156106ed57600080fd5b506107016106fc36600461462c565b611c5c565b6040516105329190614881565b34801561071a57600080fd5b506104e66107293660046149c2565b611d61565b34801561073a57600080fd5b5060115461075b9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610532565b34801561078c57600080fd5b5061052861079b36600461462c565b60036020526000908152604090205481565b3480156107b957600080fd5b506105286107c836600461462c565b611e42565b3480156107d957600080fd5b506104e66107e8366004614a46565b611eab565b3480156107f957600080fd5b5061080d61080836600461474d565b6120ca565b6040516105329190614a7c565b34801561082657600080fd5b506105287f000000000000000000000000000000000000000000000000000000000000000081565b34801561085a57600080fd5b5061052861086936600461462c565b600a6020526000908152604090205481565b34801561088757600080fd5b5061089b610896366004614ac0565b612139565b604051610532929190614b8a565b3480156108b557600080fd5b506108c96108c436600461462c565b61218d565b60405160ff9091168152602001610532565b3480156108e757600080fd5b506104e66108f6366004614ba5565b612221565b34801561090757600080fd5b5060125461075b9073ffffffffffffffffffffffffffffffffffffffff1681565b34801561093457600080fd5b506104e66122c5565b34801561094957600080fd5b506104e6610958366004614bfc565b6123fc565b34801561096957600080fd5b5061052861097836600461462c565b60076020526000908152604090205481565b34801561099657600080fd5b5060155461060f906bffffffffffffffffffffffff1681565b3480156109bb57600080fd5b506109e37f000000000000000000000000000000000000000000000000000000000000000081565b6040519015158152602001610532565b3480156109ff57600080fd5b50610528610a0e36600461462c565b60086020526000908152604090205481565b348015610a2c57600080fd5b5061089b610a3b366004614a46565b6124c9565b348015610a4c57600080fd5b506109e3610a5b36600461462c565b600b6020526000908152604090205460ff1681565b348015610a7c57600080fd5b50610528610a8b36600461462c565b6000908152600c602052604090205490565b348015610aa957600080fd5b50610528610ab836600461462c565b60046020526000908152604090205481565b348015610ad657600080fd5b506109e3610ae536600461462c565b6126de565b348015610af657600080fd5b506104e6612730565b348015610b0b57600080fd5b5061075b610b1a36600461462c565b61282d565b348015610b2b57600080fd5b5060135461075b9073ffffffffffffffffffffffffffffffffffffffff1681565b348015610b5857600080fd5b506104e6610b67366004614c1e565b6128c1565b348015610b7857600080fd5b506104e6610b87366004614c1e565b612952565b348015610b9857600080fd5b5061080d610ba7366004614bfc565b6129ac565b348015610bb857600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661075b565b348015610be357600080fd5b50610528610bf236600461462c565b60056020526000908152604090205481565b348015610c1057600080fd5b506104e6610c1f366004614bfc565b612a69565b348015610c3057600080fd5b506104e6610c3f366004614c6a565b612cae565b348015610c5057600080fd5b506104e6610c5f366004614c9a565b612dc6565b348015610c7057600080fd5b506015546108c9906c01000000000000000000000000900460ff1681565b348015610c9a57600080fd5b506104e6610ca9366004614bfc565b60009182526009602052604090912055565b348015610cc757600080fd5b5061058f610cd636600461462c565b600e6020526000908152604090205461ffff1681565b348015610cf857600080fd5b5061080d610d0736600461462c565b612fd0565b348015610d1857600080fd5b506105c2610d2736600461462c565b613032565b348015610d3857600080fd5b506105287f000000000000000000000000000000000000000000000000000000000000000081565b348015610d6c57600080fd5b506104e6610d7b366004614855565b6130de565b348015610d8c57600080fd5b506104e6610d9b366004614cb7565b613147565b348015610dac57600080fd5b506104e6610dbb36600461462c565b6131f2565b348015610dcc57600080fd5b506105c26040518060400160405280600981526020017f666565644964486578000000000000000000000000000000000000000000000081525081565b348015610e1557600080fd5b5061060f610e2436600461462c565b613278565b348015610e3557600080fd5b506104e6610e443660046149c2565b6132d0565b348015610e5557600080fd5b5061080d610e64366004614bfc565b61336a565b348015610e7557600080fd5b506105c26040518060400160405280600b81526020017f626c6f636b4e756d62657200000000000000000000000000000000000000000081525081565b348015610ebe57600080fd5b506104e6610ecd366004614cdc565b613467565b348015610ede57600080fd5b5061064c610eed366004614bfc565b6134e7565b348015610efe57600080fd5b506104e6610f0d366004614d01565b613550565b348015610f1e57600080fd5b506104e6610f2d3660046149c2565b6138b7565b348015610f3e57600080fd5b50610528610f4d366004614bfc565b613982565b348015610f5e57600080fd5b50610528610f6d36600461462c565b60096020526000908152604090205481565b348015610f8b57600080fd5b5061052860145481565b348015610fa157600080fd5b506104e6610fb0366004614779565b6139b3565b348015610fc157600080fd5b506105c2610fd0366004614d69565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b34801561107457600080fd5b5061052861108336600461462c565b60066020526000908152604090205481565b3480156110a157600080fd5b506105286110b036600461462c565b60026020526000908152604090205481565b80516110d59060169060208401906141ac565b5050565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b39216906111bf908c1688614df1565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af115801561123d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112619190614e35565b5060008860ff1667ffffffffffffffff8111156112805761128061428c565b6040519080825280602002602001820160405280156112a9578160200160208202803683370190505b50905060005b8960ff168160ff16101561131c5760006112c8846139c7565b905080838360ff16815181106112e0576112e0614e50565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061131481614e7f565b9150506112af565b50505050505050505050565b600d602052826000526040600020602052816000526040600020818154811061135057600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e6906113be9085908590600401614e9e565b600060405180830381600087803b1580156113d857600080fd5b505af11580156113ec573d6000803e3d6000fd5b505050505050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa158015611465573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114ab9190810190614f04565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa158015611550573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115969190810190614f04565b6040518363ffffffff1660e01b81526004016115b3929190614e9e565b600060405180830381600087803b1580156115cd57600080fd5b505af11580156115e1573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e010439906024015b602060405180830381865afa15801561165f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ab9190614f44565b6000828152600d6020908152604080832061ffff8516845282528083208054825181850281018501909352808352849384939291908301828280156116e757602002820191906000526020600020905b8154815260200190600101908083116116d3575b505050505090506116f9818251613a95565b92509250505b9250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b651690602401611448565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa1580156117f4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118189190614f6c565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa1580156118bb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118df9190614f9a565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611be657600089898381811061194857611948614e50565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161198191815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016119ad929190614e9e565b600060405180830381600087803b1580156119c757600080fd5b505af11580156119db573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015611a51573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a759190614fb7565b90508060ff16600103611bd1576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611afe573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611b449190810190614f04565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611b9d9086908590600401614e9e565b600060405180830381600087803b158015611bb757600080fd5b505af1158015611bcb573d6000803e3d6000fd5b50505050505b50508080611bde90614fd4565b91505061192c565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611c5293830182828015611c4657602002820191906000526020600020905b815481526020019060010190808311611c32575b50505050508484613b1a565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611d1b573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526114ab919081019061502f565b8060005b818160ff161015611e3c5760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611da357611da3614e50565b905060200201356040518263ffffffff1660e01b8152600401611dc891815260200190565b600060405180830381600087803b158015611de257600080fd5b505af1158015611df6573d6000803e3d6000fd5b50505050611e2984848360ff16818110611e1257611e12614e50565b90506020020135600f613c7990919063ffffffff16565b5080611e3481614e7f565b915050611d65565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611ea3576000858152600d6020908152604080832061ffff85168452909152902054611e8f908361514e565b915080611e9b81615161565b915050611e58565b509392505050565b60005a9050600080611ebf84860186614ac0565b91509150600081806020019051810190611ed99190615182565b60008181526005602090815260408083205460049092528220549293509190611f00613c85565b905082600003611f2057600084815260056020526040902081905561207b565b600084815260036020526040812054611f39848461519b565b611f43919061519b565b6000868152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611fb557602002820191906000526020600020905b815481526020019060010190808311611fa1575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff168151036120305781611ff281615161565b6000898152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000868152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055888352600c8252928220805493840181558252902001555b60008481526006602052604081205461209590600161514e565b60008681526006602090815260408083208490556004909152902083905590506120bf8583612a69565b61131c8589846130de565b6000828152600d6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561212c57602002820191906000526020600020905b815481526020019060010190808311612118575b5050505050905092915050565b60006060600084846040516020016121529291906151ae565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa1580156121fd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ab9190614fb7565b8160005b818110156122be5730635f17e61686868481811061224557612245614e50565b90506020020135856040518363ffffffff1660e01b815260040161227992919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561229357600080fd5b505af11580156122a7573d6000803e3d6000fd5b5050505080806122b690614fd4565b915050612225565b5050505050565b6122cd613d27565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561233c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123609190615182565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156123d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110d59190614e35565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c909152812061243491614202565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff1611612490576000848152600d6020908152604080832061ffff85168452909152812061247e91614202565b8061248881615161565b915050612449565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000606060005a905060006124e08587018761462c565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff8111156125195761251961428c565b6040519080825280601f01601f191660200182016040528015612543576020820181803683370190505b50604051602001612555929190614e9e565b60405160208183030381529060405290506000612570613c85565b9050600061257d866126de565b90505b835a61258c908961519b565b6125989061271061514e565b10156125d95781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612580565b806125f15760008398509850505050505050506116ff565b6040518060400160405280600981526020017f666565644964486578000000000000000000000000000000000000000000000081525060166040518060400160405280600b81526020017f626c6f636b4e756d626572000000000000000000000000000000000000000000815250848960405160200161267391815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f7ddd933e0000000000000000000000000000000000000000000000000000000082526126d59594939291600401615295565b60405180910390fd5b60008181526005602052604081205481036126fb57506001919050565b60008281526003602090815260408083205460049092529091205461271e613c85565b612728919061519b565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146127b1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016126d5565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa15801561289d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ab9190614f9a565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b59061291b908690869086906004016153e2565b600060405180830381600087803b15801561293557600080fd5b505af1158015612949573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d359061291b908690869086906004016153e2565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612a23573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611c559190810190615436565b601454600083815260026020526040902054612a85908361519b565b11156110d5576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612afb573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612b41919081019061502f565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612bb6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612bda9190614f44565b601554909150612bfe9082906c01000000000000000000000000900460ff16614df1565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611e3c57601554612c419085906bffffffffffffffffffffffff16612cae565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612d36573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d5a9190614e35565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f7906044016113be565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015612e25573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612e6b9190810190615436565b80519091506000612e7a613c85565b905060005b828110156122be576000848281518110612e9b57612e9b614e50565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612f1b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f3f9190614fb7565b90508060ff16600103612fbb578660ff16600003612f8b576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4612fbb565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b50508080612fc890614fd4565b915050612e7f565b6000818152600c602090815260409182902080548351818402810184019094528084526060939283018282801561302657602002820191906000526020600020905b815481526020019060010190808311613012575b50505050509050919050565b6016818154811061304257600080fd5b90600052602060002001600091509050805461305d90615242565b80601f016020809104026020016040519081016040528092919081815260200182805461308990615242565b80156130d65780601f106130ab576101008083540402835291602001916130d6565b820191906000526020600020905b8154815290600101906020018083116130b957829003601f168201915b505050505081565b6000838152600760205260409020545b805a6130fa908561519b565b6131069061271061514e565b1015611e3c5781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556130ee565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156131bf57600080fd5b505af11580156131d3573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561326457600080fd5b505af11580156122be573d6000803e3d6000fd5b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff169063b657bc9c90602401611642565b8060005b818163ffffffff161015611e3c573063af953a4a858563ffffffff851681811061330057613300614e50565b905060200201356040518263ffffffff1660e01b815260040161332591815260200190565b600060405180830381600087803b15801561333f57600080fd5b505af1158015613353573d6000803e3d6000fd5b505050508080613362906154c7565b9150506132d4565b60606000613378600f613daa565b90508084106133b3576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826000036133c8576133c5848261519b565b92505b60008367ffffffffffffffff8111156133e3576133e361428c565b60405190808252806020026020018201604052801561340c578160200160208202803683370190505b50905060005b8481101561345e5761342f613427828861514e565b600f90613db4565b82828151811061344157613441614e50565b60209081029190910101528061345681614fd4565b915050613412565b50949350505050565b6000613471613c85565b90508160ff166000036134b2576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b505050565b6000828152600c6020908152604080832080548251818502810185019093528083528493849392919083018282801561353f57602002820191906000526020600020905b81548152602001906001019080831161352b575b505050505090506116f98185613a95565b8260005b818110156113ec57600086868381811061357057613570614e50565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016135a991815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016135d5929190614e9e565b600060405180830381600087803b1580156135ef57600080fd5b505af1158015613603573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015613679573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061369d9190614fb7565b90508060ff166001036138a2577f000000000000000000000000000000000000000000000000000000000000000060ff8716156136f757507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb3089858860405160200161372b91815260200190565b604051602081830303815290604052613743906154e0565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa1580156137ce573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526138149190810190614f04565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d359061386d9087908590600401614e9e565b600060405180830381600087803b15801561388757600080fd5b505af115801561389b573d6000803e3d6000fd5b5050505050505b505080806138af90614fd4565b915050613554565b8060005b81811015611e3c5760008484838181106138d7576138d7614e50565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161391091815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b815260040161393c929190614e9e565b600060405180830381600087803b15801561395657600080fd5b505af115801561396a573d6000803e3d6000fd5b5050505050808061397a90614fd4565b9150506138bb565b600c602052816000526040600020818154811061399e57600080fd5b90600052602060002001600091509150505481565b6139bb613d27565b6139c481613dc0565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613a22908690600401615522565b6020604051808303816000875af1158015613a41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a659190615182565b9050613a72600f82613eb5565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b815160009081908190841580613aab5750808510155b15613ab4578094505b60008092505b85831015613b1057866001613acf858561519b565b613ad9919061519b565b81518110613ae957613ae9614e50565b602002602001015181613afc919061514e565b905082613b0881614fd4565b935050613aba565b9694955050505050565b82516000908190831580613b2e5750808410155b15613b37578093505b60008467ffffffffffffffff811115613b5257613b5261428c565b604051908082528060200260200182016040528015613b7b578160200160208202803683370190505b509050600092505b84831015613be957866001613b98858561519b565b613ba2919061519b565b81518110613bb257613bb2614e50565b6020026020010151818481518110613bcc57613bcc614e50565b602090810291909101015282613be181614fd4565b935050613b83565b613c0281600060018451613bfd919061519b565b613ec1565b85606403613c3b578060018251613c19919061519b565b81518110613c2957613c29614e50565b60200260200101519350505050611c55565b806064825188613c4b9190615674565b613c5591906156e0565b81518110613c6557613c65614e50565b602002602001015193505050509392505050565b6000611c558383614039565b60007f000000000000000000000000000000000000000000000000000000000000000015613d2257606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613cf9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d1d9190615182565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613da8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016126d5565b565b60006114ab825490565b6000611c558383614133565b3373ffffffffffffffffffffffffffffffffffffffff821603613e3f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016126d5565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611c55838361415d565b8181808203613ed1575050505050565b6000856002613ee087876156f4565b613eea9190615714565b613ef4908761577c565b81518110613f0457613f04614e50565b602002602001015190505b818313614013575b80868481518110613f2a57613f2a614e50565b60200260200101511015613f4a5782613f42816157a4565b935050613f17565b858281518110613f5c57613f5c614e50565b6020026020010151811015613f7d5781613f75816157d5565b925050613f4a565b81831361400e57858281518110613f9657613f96614e50565b6020026020010151868481518110613fb057613fb0614e50565b6020026020010151878581518110613fca57613fca614e50565b60200260200101888581518110613fe357613fe3614e50565b60209081029190910101919091525282613ffc816157a4565b935050818061400a906157d5565b9250505b613f0f565b8185121561402657614026868684613ec1565b838312156113ec576113ec868486613ec1565b6000818152600183016020526040812054801561412257600061405d60018361519b565b85549091506000906140719060019061519b565b90508181146140d657600086600001828154811061409157614091614e50565b90600052602060002001549050808760000184815481106140b4576140b4614e50565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806140e7576140e761582c565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506114ab565b60009150506114ab565b5092915050565b600082600001828154811061414a5761414a614e50565b9060005260206000200154905092915050565b60008181526001830160205260408120546141a4575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556114ab565b5060006114ab565b8280548282559060005260206000209081019282156141f2579160200282015b828111156141f257825182906141e290826158a1565b50916020019190600101906141cc565b506141fe929150614220565b5090565b50805460008255906000526020600020908101906139c4919061423d565b808211156141fe5760006142348282614252565b50600101614220565b5b808211156141fe576000815560010161423e565b50805461425e90615242565b6000825580601f1061426e575050565b601f0160209004906000526020600020908101906139c4919061423d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff811182821017156142df576142df61428c565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561432c5761432c61428c565b604052919050565b600067ffffffffffffffff82111561434e5761434e61428c565b5060051b60200190565b600067ffffffffffffffff8211156143725761437261428c565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60006143b16143ac84614358565b6142e5565b90508281528383830111156143c557600080fd5b828260208301376000602084830101529392505050565b600060208083850312156143ef57600080fd5b823567ffffffffffffffff8082111561440757600080fd5b818501915085601f83011261441b57600080fd5b81356144296143ac82614334565b81815260059190911b8301840190848101908883111561444857600080fd5b8585015b83811015614495578035858111156144645760008081fd5b8601603f81018b136144765760008081fd5b6144878b898301356040840161439e565b84525091860191860161444c565b5098975050505050505050565b60ff811681146139c457600080fd5b63ffffffff811681146139c457600080fd5b600082601f8301126144d457600080fd5b611c558383356020850161439e565b6bffffffffffffffffffffffff811681146139c457600080fd5b600080600080600080600060e0888a03121561451857600080fd5b8735614523816144a2565b96506020880135614533816144b1565b95506040880135614543816144a2565b9450606088013567ffffffffffffffff81111561455f57600080fd5b61456b8a828b016144c3565b945050608088013561457c816144e3565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff811681146145ab57600080fd5b919050565b6000806000606084860312156145c557600080fd5b833592506145d560208501614599565b9150604084013590509250925092565b600080604083850312156145f857600080fd5b82359150602083013567ffffffffffffffff81111561461657600080fd5b614622858286016144c3565b9150509250929050565b60006020828403121561463e57600080fd5b5035919050565b60005b83811015614660578181015183820152602001614648565b50506000910152565b60008151808452614681816020860160208601614645565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611c556020830184614669565b73ffffffffffffffffffffffffffffffffffffffff811681146139c457600080fd5b600080600080600080600060e0888a03121561470357600080fd5b873596506020880135614715816146c6565b95506040880135614725816144a2565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b6000806040838503121561476057600080fd5b8235915061477060208401614599565b90509250929050565b60006020828403121561478b57600080fd5b8135611c55816146c6565b60008083601f8401126147a857600080fd5b50813567ffffffffffffffff8111156147c057600080fd5b6020830191508360208260051b85010111156116ff57600080fd5b600080600080600080600060c0888a0312156147f657600080fd5b873567ffffffffffffffff81111561480d57600080fd5b6148198a828b01614796565b909850965050602088013561482d816144a2565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b60008060006060848603121561486a57600080fd5b505081359360208301359350604090920135919050565b602081526148a860208201835173ffffffffffffffffffffffffffffffffffffffff169052565b600060208301516148c1604084018263ffffffff169052565b5060408301516101408060608501526148de610160850183614669565b915060608501516148ff60808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e085015161010061496b818701836bffffffffffffffffffffffff169052565b86015190506101206149808682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018387015290506149b88382614669565b9695505050505050565b600080602083850312156149d557600080fd5b823567ffffffffffffffff8111156149ec57600080fd5b6149f885828601614796565b90969095509350505050565b60008083601f840112614a1657600080fd5b50813567ffffffffffffffff811115614a2e57600080fd5b6020830191508360208285010111156116ff57600080fd5b60008060208385031215614a5957600080fd5b823567ffffffffffffffff811115614a7057600080fd5b6149f885828601614a04565b6020808252825182820181905260009190848201906040850190845b81811015614ab457835183529284019291840191600101614a98565b50909695505050505050565b60008060408385031215614ad357600080fd5b823567ffffffffffffffff80821115614aeb57600080fd5b818501915085601f830112614aff57600080fd5b81356020614b0f6143ac83614334565b82815260059290921b84018101918181019089841115614b2e57600080fd5b8286015b84811015614b6657803586811115614b4a5760008081fd5b614b588c86838b01016144c3565b845250918301918301614b32565b5096505086013592505080821115614b7d57600080fd5b50614622858286016144c3565b8215158152604060208201526000611c526040830184614669565b600080600060408486031215614bba57600080fd5b833567ffffffffffffffff811115614bd157600080fd5b614bdd86828701614796565b9094509250506020840135614bf1816144b1565b809150509250925092565b60008060408385031215614c0f57600080fd5b50508035926020909101359150565b600080600060408486031215614c3357600080fd5b83359250602084013567ffffffffffffffff811115614c5157600080fd5b614c5d86828701614a04565b9497909650939450505050565b60008060408385031215614c7d57600080fd5b823591506020830135614c8f816144e3565b809150509250929050565b600060208284031215614cac57600080fd5b8135611c55816144a2565b60008060408385031215614cca57600080fd5b823591506020830135614c8f816144b1565b60008060408385031215614cef57600080fd5b823591506020830135614c8f816144a2565b60008060008060608587031215614d1757600080fd5b843567ffffffffffffffff811115614d2e57600080fd5b614d3a87828801614796565b9095509350506020850135614d4e816144a2565b91506040850135614d5e816144a2565b939692955090935050565b60008060008060008060c08789031215614d8257600080fd5b8635614d8d816146c6565b95506020870135614d9d816144a2565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff80831681851681830481118215151615614e1c57614e1c614dc2565b02949350505050565b805180151581146145ab57600080fd5b600060208284031215614e4757600080fd5b611c5582614e25565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614e9557614e95614dc2565b60010192915050565b828152604060208201526000611c526040830184614669565b600082601f830112614ec857600080fd5b8151614ed66143ac82614358565b818152846020838601011115614eeb57600080fd5b614efc826020830160208701614645565b949350505050565b600060208284031215614f1657600080fd5b815167ffffffffffffffff811115614f2d57600080fd5b614efc84828501614eb7565b80516145ab816144e3565b600060208284031215614f5657600080fd5b8151611c55816144e3565b80516145ab816146c6565b60008060408385031215614f7f57600080fd5b8251614f8a816146c6565b6020939093015192949293505050565b600060208284031215614fac57600080fd5b8151611c55816146c6565b600060208284031215614fc957600080fd5b8151611c55816144a2565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361500557615005614dc2565b5060010190565b80516145ab816144b1565b805167ffffffffffffffff811681146145ab57600080fd5b60006020828403121561504157600080fd5b815167ffffffffffffffff8082111561505957600080fd5b90830190610140828603121561506e57600080fd5b6150766142bb565b61507f83614f61565b815261508d6020840161500c565b60208201526040830151828111156150a457600080fd5b6150b087828601614eb7565b6040830152506150c260608401614f39565b60608201526150d360808401614f61565b60808201526150e460a08401615017565b60a08201526150f560c0840161500c565b60c082015261510660e08401614f39565b60e0820152610100615119818501614e25565b90820152610120838101518381111561513157600080fd5b61513d88828701614eb7565b918301919091525095945050505050565b808201808211156114ab576114ab614dc2565b600061ffff80831681810361517857615178614dc2565b6001019392505050565b60006020828403121561519457600080fd5b5051919050565b818103818111156114ab576114ab614dc2565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b83811015615223577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0888703018552615211868351614669565b955093820193908201906001016151d7565b5050858403818701525050506152398185614669565b95945050505050565b600181811c9082168061525657607f821691505b60208210810361528f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60a0815260006152a860a0830188614669565b602083820381850152818854808452828401915060058382821b86010160008c8152858120815b858110156153a2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe089850301875282825461530a81615242565b80875260018281168015615325576001811461535c5761538b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168d8a01528c8315158b1b8a0101945061538b565b8688528c8820885b848110156153835781548f828d01015283820191508e81019050615364565b8a018e019550505b50998b0199929650505091909101906001016152cf565b50505087810360408901526153b7818c614669565b9550505050505084606084015282810360808401526153d68185614669565b98975050505050505050565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b6000602080838503121561544957600080fd5b825167ffffffffffffffff81111561546057600080fd5b8301601f8101851361547157600080fd5b805161547f6143ac82614334565b81815260059190911b8201830190838101908783111561549e57600080fd5b928401925b828410156154bc578351825292840192908401906154a3565b979650505050505050565b600063ffffffff80831681810361517857615178614dc2565b8051602080830151919081101561528f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b6020815260008251610140806020850152615541610160850183614669565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08086850301604087015261557d8483614669565b9350604087015191506155a8606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e08701526156098483614669565b935060e087015191506101008187860301818801526156288584614669565b9450808801519250506101208187860301818801526156478584614669565b9450808801519250505061566a828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156156ac576156ac614dc2565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826156ef576156ef6156b1565b500490565b818103600083128015838313168383128216171561412c5761412c614dc2565b600082615723576157236156b1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561577757615777614dc2565b500590565b808201828112600083128015821682158216171561579c5761579c614dc2565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361500557615005614dc2565b60007f8000000000000000000000000000000000000000000000000000000000000000820361580657615806614dc2565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b601f8211156134e257600081815260208120601f850160051c810160208610156158825750805b601f850160051c820191505b818110156113ec5782815560010161588e565b815167ffffffffffffffff8111156158bb576158bb61428c565b6158cf816158c98454615242565b8461585b565b602080601f83116001811461592257600084156158ec5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556113ec565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561596f57888601518255948401946001909101908401615950565b50858210156159ab57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactor) SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.contract.Transact(opts, "setFeedsHex", newFeeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeedsHex(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeeds)
}

func (_VerifiableLoadMercuryUpkeep *VerifiableLoadMercuryUpkeepTransactorSession) SetFeedsHex(newFeeds []string) (*types.Transaction, error) {
	return _VerifiableLoadMercuryUpkeep.Contract.SetFeedsHex(&_VerifiableLoadMercuryUpkeep.TransactOpts, newFeeds)
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

	SetFeedsHex(opts *bind.TransactOpts, newFeeds []string) (*types.Transaction, error)

	SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error)

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
