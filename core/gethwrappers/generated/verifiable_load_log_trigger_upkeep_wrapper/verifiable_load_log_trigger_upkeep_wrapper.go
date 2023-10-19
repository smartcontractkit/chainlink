// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifiable_load_log_trigger_upkeep_wrapper

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

type Log struct {
	Index       *big.Int
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var VerifiableLoadLogTriggerUpkeepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_useMercury\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"feeds\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"timeParamKey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"StreamsLookup\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"values\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"checkCallback\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"logNum\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_feeds\",\"type\":\"string[]\"}],\"name\":\"setFeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_log\",\"type\":\"uint8\"}],\"name\":\"setLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_feedParamKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_timeParamKey\",\"type\":\"string\"}],\"name\":\"setParamKeys\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useMercury\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c0604052604261014081815261010091829190620066ec6101603981526020016040518060800160405280604281526020016200672e604291399052620000be906016906002620003de565b506040805180820190915260098152680cccacac892c890caf60bb1b6020820152601790620000ee90826200055a565b5060408051808201909152600b81526a313637b1b5a73ab6b132b960a91b60208201526018906200012090826200055a565b503480156200012e57600080fd5b506040516200677038038062006770833981016040819052620001519162000652565b82823380600081620001aa5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620001dd57620001dd8162000333565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa1580156200023a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200026091906200069e565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015620002c6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ec9190620006cf565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c052506019805461ffff191691151561ff00191691909117905550620006f69050565b336001600160a01b038216036200038d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620001a1565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b82805482825590600052602060002090810192821562000429579160200282015b828111156200042957825182906200041890826200055a565b5091602001919060010190620003ff565b50620004379291506200043b565b5090565b80821115620004375760006200045282826200045c565b506001016200043b565b5080546200046a90620004cb565b6000825580601f106200047b575050565b601f0160209004906000526020600020908101906200049b91906200049e565b50565b5b808211156200043757600081556001016200049f565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620004e057607f821691505b6020821081036200050157634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200055557600081815260208120601f850160051c81016020861015620005305750805b601f850160051c820191505b8181101562000551578281556001016200053c565b5050505b505050565b81516001600160401b03811115620005765762000576620004b5565b6200058e81620005878454620004cb565b8462000507565b602080601f831160018114620005c65760008415620005ad5750858301515b600019600386901b1c1916600185901b17855562000551565b600085815260208120601f198616915b82811015620005f757888601518255948401946001909101908401620005d6565b5085821015620006165787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b03811681146200049b57600080fd5b805180151581146200064d57600080fd5b919050565b6000806000606084860312156200066857600080fd5b8351620006758162000626565b925062000685602085016200063c565b915062000695604085016200063c565b90509250925092565b60008060408385031215620006b257600080fd5b8251620006bf8162000626565b6020939093015192949293505050565b600060208284031215620006e257600080fd5b8151620006ef8162000626565b9392505050565b60805160a05160c05160e051615f926200075a600039600081816105b90152612551015260008181610a2d01526140920152600081816108a601528181611fa80152613ae0015260008181610dca01528181611f780152613ab50152615f926000f3fe6080604052600436106105265760003560e01c80637b103999116102af578063af953a4a11610179578063daee1aeb116100d6578063e83ce5581161008a578063fa333dfb1161006f578063fa333dfb14611066578063fba7ffa314611119578063fcdc1f631461114657600080fd5b8063e83ce55814611027578063f2fde38b1461104657600080fd5b8063de818253116100bb578063de81825314610f90578063e0114adb14610fe4578063e45530831461101157600080fd5b8063daee1aeb14610f50578063dbef701e14610f7057600080fd5b8063c41c815b1161012d578063d4c2490011610112578063d4c2490014610ef0578063d6051a7214610f10578063da6cba4714610f3057600080fd5b8063c41c815b14610ec1578063c98f10b014610edb57600080fd5b8063b657bc9c1161015e578063b657bc9c14610e61578063becde0e114610e81578063c041982214610ea157600080fd5b8063af953a4a14610e2c578063afb28d1f14610e4c57600080fd5b8063948108f7116102275780639d385eaa116101db578063a6548248116101c0578063a654824814610db8578063a6b5947514610dec578063a72aa27e14610e0c57600080fd5b80639d385eaa14610d785780639d6f1cc714610d9857600080fd5b80639ac542eb1161020c5780639ac542eb14610cf05780639b42935414610d1a5780639b51fb0d14610d4757600080fd5b8063948108f714610cb057806396cebc7c14610cd057600080fd5b806386e330af1161027e5780638da5cb5b116102635780638da5cb5b14610c385780638fcb3fba14610c63578063924ca57814610c9057600080fd5b806386e330af14610bf8578063873c758614610c1857600080fd5b80637b10399914610b6b5780637e7a46dc14610b985780638243444a14610bb85780638340507c14610bd857600080fd5b806345d2ec17116103f057806360457ff51161036857806373644cce1161031c578063776898c811610301578063776898c814610b1657806379ba509714610b3657806379ea994314610b4b57600080fd5b806373644cce14610abc5780637672130314610ae957600080fd5b8063642f6cef1161034d578063642f6cef14610a1b57806369cdbadb14610a5f5780637145f11b14610a8c57600080fd5b806360457ff5146109c9578063636092e8146109f657600080fd5b80635147cd59116103bf57806357970e93116103a457806357970e93146109675780635d4ee7f3146109945780635f17e616146109a957600080fd5b80635147cd591461091557806351c98be31461094757600080fd5b806345d2ec1714610867578063469820931461089457806346e7a63e146108c85780634b56a42e146108f557600080fd5b806320e3dbd41161049e5780632b20e397116104525780633ebe8d6c116104375780633ebe8d6c146107f957806340691db4146108195780634585e33b1461084757600080fd5b80632b20e3971461077a578063328ffd11146107cc57600080fd5b806328c4b57b1161048357806328c4b57b1461070d57806329e0a8411461072d5780632a9032d31461075a57600080fd5b806320e3dbd4146106cd5780632636aecf146106ed57600080fd5b806319d97a94116104f55780631e010439116104da5780631e0104391461063b578063206c32e814610678578063207b6516146106ad57600080fd5b806319d97a94146105ee5780631cdde2511461061b57600080fd5b806306c1cc0014610532578063077ac621146105545780630b7d33e61461058757806312c55027146105a757600080fd5b3661052d57005b600080fd5b34801561053e57600080fd5b5061055261054d366004614814565b611173565b005b34801561056057600080fd5b5061057461056f3660046148c7565b6113c2565b6040519081526020015b60405180910390f35b34801561059357600080fd5b506105526105a23660046148fc565b611400565b3480156105b357600080fd5b506105db7f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff909116815260200161057e565b3480156105fa57600080fd5b5061060e610609366004614943565b61148e565b60405161057e91906149ca565b34801561062757600080fd5b506105526106363660046149ff565b61154b565b34801561064757600080fd5b5061065b610656366004614943565b611688565b6040516bffffffffffffffffffffffff909116815260200161057e565b34801561068457600080fd5b50610698610693366004614a64565b61171d565b6040805192835260208301919091520161057e565b3480156106b957600080fd5b5061060e6106c8366004614943565b6117a0565b3480156106d957600080fd5b506105526106e8366004614a90565b6117f8565b3480156106f957600080fd5b50610552610708366004614af2565b6119c2565b34801561071957600080fd5b50610574610728366004614b6c565b611c8b565b34801561073957600080fd5b5061074d610748366004614943565b611cf6565b60405161057e9190614b98565b34801561076657600080fd5b50610552610775366004614cd9565b611dfb565b34801561078657600080fd5b506011546107a79073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161057e565b3480156107d857600080fd5b506105746107e7366004614943565b60036020526000908152604090205481565b34801561080557600080fd5b50610574610814366004614943565b611edc565b34801561082557600080fd5b50610839610834366004614d1b565b611f45565b60405161057e929190614d7e565b34801561085357600080fd5b50610552610862366004614ddb565b61244b565b34801561087357600080fd5b50610887610882366004614a64565b61269a565b60405161057e9190614e11565b3480156108a057600080fd5b506105747f000000000000000000000000000000000000000000000000000000000000000081565b3480156108d457600080fd5b506105746108e3366004614943565b600a6020526000908152604090205481565b34801561090157600080fd5b50610839610910366004614e79565b612709565b34801561092157600080fd5b50610935610930366004614943565b61275d565b60405160ff909116815260200161057e565b34801561095357600080fd5b50610552610962366004614f36565b6127f1565b34801561097357600080fd5b506012546107a79073ffffffffffffffffffffffffffffffffffffffff1681565b3480156109a057600080fd5b50610552612895565b3480156109b557600080fd5b506105526109c4366004614f8d565b6129d0565b3480156109d557600080fd5b506105746109e4366004614943565b60076020526000908152604090205481565b348015610a0257600080fd5b5060155461065b906bffffffffffffffffffffffff1681565b348015610a2757600080fd5b50610a4f7f000000000000000000000000000000000000000000000000000000000000000081565b604051901515815260200161057e565b348015610a6b57600080fd5b50610574610a7a366004614943565b60086020526000908152604090205481565b348015610a9857600080fd5b50610a4f610aa7366004614943565b600b6020526000908152604090205460ff1681565b348015610ac857600080fd5b50610574610ad7366004614943565b6000908152600c602052604090205490565b348015610af557600080fd5b50610574610b04366004614943565b60046020526000908152604090205481565b348015610b2257600080fd5b50610a4f610b31366004614943565b612a9d565b348015610b4257600080fd5b50610552612aef565b348015610b5757600080fd5b506107a7610b66366004614943565b612bec565b348015610b7757600080fd5b506013546107a79073ffffffffffffffffffffffffffffffffffffffff1681565b348015610ba457600080fd5b50610552610bb3366004614faf565b612c80565b348015610bc457600080fd5b50610552610bd3366004614faf565b612d11565b348015610be457600080fd5b50610552610bf3366004614ffb565b612d6b565b348015610c0457600080fd5b50610552610c13366004615048565b612d89565b348015610c2457600080fd5b50610887610c33366004614f8d565b612d9c565b348015610c4457600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff166107a7565b348015610c6f57600080fd5b50610574610c7e366004614943565b60056020526000908152604090205481565b348015610c9c57600080fd5b50610552610cab366004614f8d565b612e59565b348015610cbc57600080fd5b50610552610ccb3660046150f9565b61309e565b348015610cdc57600080fd5b50610552610ceb366004615129565b6131b6565b348015610cfc57600080fd5b50601554610935906c01000000000000000000000000900460ff1681565b348015610d2657600080fd5b50610552610d35366004614f8d565b60009182526009602052604090912055565b348015610d5357600080fd5b506105db610d62366004614943565b600e6020526000908152604090205461ffff1681565b348015610d8457600080fd5b50610887610d93366004614943565b6133c0565b348015610da457600080fd5b5061060e610db3366004614943565b613422565b348015610dc457600080fd5b506105747f000000000000000000000000000000000000000000000000000000000000000081565b348015610df857600080fd5b50610552610e07366004614b6c565b6134ce565b348015610e1857600080fd5b50610552610e27366004615146565b613537565b348015610e3857600080fd5b50610552610e47366004614943565b6135e2565b348015610e5857600080fd5b5061060e613668565b348015610e6d57600080fd5b5061065b610e7c366004614943565b613675565b348015610e8d57600080fd5b50610552610e9c366004614cd9565b6136cd565b348015610ead57600080fd5b50610887610ebc366004614f8d565b613767565b348015610ecd57600080fd5b50601954610a4f9060ff1681565b348015610ee757600080fd5b5061060e613864565b348015610efc57600080fd5b50610552610f0b36600461516b565b613871565b348015610f1c57600080fd5b50610698610f2b366004614f8d565b6138f0565b348015610f3c57600080fd5b50610552610f4b366004615190565b613959565b348015610f5c57600080fd5b50610552610f6b366004614cd9565b613cc0565b348015610f7c57600080fd5b50610574610f8b366004614f8d565b613d8b565b348015610f9c57600080fd5b50610552610fab366004615129565b6019805460ff909216610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff909216919091179055565b348015610ff057600080fd5b50610574610fff366004614943565b60096020526000908152604090205481565b34801561101d57600080fd5b5061057460145481565b34801561103357600080fd5b5060195461093590610100900460ff1681565b34801561105257600080fd5b50610552611061366004614a90565b613dbc565b34801561107257600080fd5b5061060e6110813660046151f8565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b34801561112557600080fd5b50610574611134366004614943565b60066020526000908152604090205481565b34801561115257600080fd5b50610574611161366004614943565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b3921690611259908c1688615280565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af11580156112d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112fb91906152c4565b5060008860ff1667ffffffffffffffff81111561131a5761131a6146b6565b604051908082528060200260200182016040528015611343578160200160208202803683370190505b50905060005b8960ff168160ff1610156113b657600061136284613dd0565b905080838360ff168151811061137a5761137a6152df565b602090810291909101810191909152600091825260088152604080832088905560079091529020849055806113ae8161530e565b915050611349565b50505050505050505050565b600d60205282600052604060002060205281600052604060002081815481106113ea57600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e690611458908590859060040161532d565b600060405180830381600087803b15801561147257600080fd5b505af1158015611486573d6000803e3d6000fd5b505050505050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa1580156114ff573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115459190810190615393565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa1580156115ea573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526116309190810190615393565b6040518363ffffffff1660e01b815260040161164d92919061532d565b600060405180830381600087803b15801561166757600080fd5b505af115801561167b573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e010439906024015b602060405180830381865afa1580156116f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061154591906153d3565b6000828152600d6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493849392919083018282801561178157602002820191906000526020600020905b81548152602001906001019080831161176d575b50505050509050611793818251613e9e565b92509250505b9250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b6516906024016114e2565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa15801561188e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118b291906153fb565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa158015611955573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119799190615429565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611c805760008989838181106119e2576119e26152df565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc8283604051602001611a1b91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401611a4792919061532d565b600060405180830381600087803b158015611a6157600080fd5b505af1158015611a75573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015611aeb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b0f9190615446565b90508060ff16600103611c6b576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611b98573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611bde9190810190615393565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611c37908690859060040161532d565b600060405180830381600087803b158015611c5157600080fd5b505af1158015611c65573d6000803e3d6000fd5b50505050505b50508080611c7890615463565b9150506119c6565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611cec93830182828015611ce057602002820191906000526020600020905b815481526020019060010190808311611ccc575b50505050508484613f23565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611db5573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261154591908101906154be565b8060005b818160ff161015611ed65760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611e3d57611e3d6152df565b905060200201356040518263ffffffff1660e01b8152600401611e6291815260200190565b600060405180830381600087803b158015611e7c57600080fd5b505af1158015611e90573d6000803e3d6000fd5b50505050611ec384848360ff16818110611eac57611eac6152df565b90506020020135600f61408290919063ffffffff16565b5080611ece8161530e565b915050611dff565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611f3d576000858152600d6020908152604080832061ffff85168452909152902054611f2990836155dd565b915080611f35816155f0565b915050611ef2565b509392505050565b6000606060005a90506000611f5861408e565b9050600085806020019051810190611f709190615611565b6019549091507f000000000000000000000000000000000000000000000000000000000000000090610100900460ff1615611fc857507f00000000000000000000000000000000000000000000000000000000000000005b80611fd660c08a018a61562a565b6000818110611fe757611fe76152df565b90506020020135036123e957600061200260c08a018a61562a565b6001818110612013576120136152df565b9050602002013560405160200161202c91815260200190565b60405160208183030381529060405290506000818060200190518101906120539190615611565b90508381146120c3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f75706b6565702069647320646f6e2774206d617463680000000000000000000060448201526064015b60405180910390fd5b60006120d260c08c018c61562a565b60028181106120e3576120e36152df565b905060200201356040516020016120fc91815260200190565b60405160208183030381529060405290506000818060200190518101906121239190615611565b9050600061213460c08e018e61562a565b6003818110612145576121456152df565b9050602002013560405160200161215e91815260200190565b60405160208183030381529060405290506000818060200190518101906121859190615429565b6000868152600860205260409020549091505b805a6121a4908d615692565b6121b090613a986155dd565b10156121f15783406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612198565b6040517f6665656449644865780000000000000000000000000000000000000000000000602082015260009060290160405160208183030381529060405280519060200120601760405160200161224891906156f8565b604051602081830303815290604052805190602001200361226a57508361226d565b50425b60195460ff161561231557604080516020810189905290810186905273ffffffffffffffffffffffffffffffffffffffff841660608201526017906016906018908490608001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527ff055e4a20000000000000000000000000000000000000000000000000000000082526120ba9594939291600401615827565b60165460009067ffffffffffffffff811115612333576123336146b6565b60405190808252806020026020018201604052801561236657816020015b60608152602001906001900390816123515790505b5060408051602081018b905290810188905273ffffffffffffffffffffffffffffffffffffffff861660608201529091506000906080016040516020818303038152906040529050600182826040516020016123c39291906158ea565b6040516020818303038152906040529f509f505050505050505050505050505050611799565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f756e6578706563746564206576656e742073696700000000000000000000000060448201526064016120ba565b60005a905060008061245f84860186614e79565b9150915060008060008380602001905181019061247c919061597e565b60008381526005602090815260408083205460049092528220549497509295509093509091906124aa61408e565b9050826000036124ca57600086815260056020526040902081905561260e565b60006124d68683615692565b6000888152600e6020908152604080832054600d835281842061ffff90911680855290835281842080548351818602810186019094528084529596509094919290919083018282801561254857602002820191906000526020600020905b815481526020019060010190808311612534575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff168151036125c35781612585816155f0565b60008b8152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000888152600d6020908152604080832061ffff9094168352928152828220805460018181018355918452828420018590558a8352600c8252928220805493840181558252902001555b6000868152600660205260408120546126289060016155dd565b60008881526006602090815260408083208490556004909152902083905590506126528783612e59565b6040513090839089907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a461268c878b846134ce565b505050505050505050505050565b6000828152600d6020908152604080832061ffff851684528252918290208054835181840281018401909452808452606093928301828280156126fc57602002820191906000526020600020905b8154815260200190600101908083116126e8575b5050505050905092915050565b60006060600084846040516020016127229291906158ea565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00181529190526001969095509350505050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa1580156127cd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115459190615446565b8160005b8181101561288e5730635f17e616868684818110612815576128156152df565b90506020020135856040518363ffffffff1660e01b815260040161284992919091825263ffffffff16602082015260400190565b600060405180830381600087803b15801561286357600080fd5b505af1158015612877573d6000803e3d6000fd5b50505050808061288690615463565b9150506127f5565b5050505050565b61289d614130565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561290c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129309190615611565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156129a8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129cc91906152c4565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c9091528120612a08916145b5565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff1611612a64576000848152600d6020908152604080832061ffff851684529091528120612a52916145b5565b80612a5c816155f0565b915050612a1d565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000818152600560205260408120548103612aba57506001919050565b600082815260036020908152604080832054600490925290912054612add61408e565b612ae79190615692565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612b70576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016120ba565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa158015612c5c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115459190615429565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b590612cda908690869086906004016159ac565b600060405180830381600087803b158015612cf457600080fd5b505af1158015612d08573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d3590612cda908690869086906004016159ac565b6017612d778382615a46565b506018612d848282615a46565b505050565b80516129cc9060169060208401906145d3565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612e13573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611cef9190810190615b60565b601454600083815260026020526040902054612e759083615692565b11156129cc576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa158015612eeb573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612f3191908101906154be565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa158015612fa6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612fca91906153d3565b601554909150612fee9082906c01000000000000000000000000900460ff16615280565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611ed6576015546130319085906bffffffffffffffffffffffff1661309e565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015613126573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061314a91906152c4565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401611458565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015613215573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261325b9190810190615b60565b8051909150600061326a61408e565b905060005b8281101561288e57600084828151811061328b5761328b6152df565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa15801561330b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061332f9190615446565b90508060ff166001036133ab578660ff1660000361337b576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a46133ab565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b505080806133b890615463565b91505061326f565b6000818152600c602090815260409182902080548351818402810184019094528084526060939283018282801561341657602002820191906000526020600020905b815481526020019060010190808311613402575b50505050509050919050565b6016818154811061343257600080fd5b90600052602060002001600091509050805461344d906156a5565b80601f0160208091040260200160405190810160405280929190818152602001828054613479906156a5565b80156134c65780601f1061349b576101008083540402835291602001916134c6565b820191906000526020600020905b8154815290600101906020018083116134a957829003601f168201915b505050505081565b6000838152600760205260409020545b805a6134ea9085615692565b6134f6906127106155dd565b1015611ed65781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556134de565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b1580156135af57600080fd5b505af11580156135c3573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561365457600080fd5b505af115801561288e573d6000803e3d6000fd5b6017805461344d906156a5565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff169063b657bc9c906024016116dc565b8060005b818163ffffffff161015611ed6573063af953a4a858563ffffffff85168181106136fd576136fd6152df565b905060200201356040518263ffffffff1660e01b815260040161372291815260200190565b600060405180830381600087803b15801561373c57600080fd5b505af1158015613750573d6000803e3d6000fd5b50505050808061375f90615bf1565b9150506136d1565b60606000613775600f6141b3565b90508084106137b0576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826000036137c5576137c28482615692565b92505b60008367ffffffffffffffff8111156137e0576137e06146b6565b604051908082528060200260200182016040528015613809578160200160208202803683370190505b50905060005b8481101561385b5761382c61382482886155dd565b600f906141bd565b82828151811061383e5761383e6152df565b60209081029190910101528061385381615463565b91505061380f565b50949350505050565b6018805461344d906156a5565b600061387b61408e565b90508160ff166000036138bc576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a4505050565b6000828152600c6020908152604080832080548251818502810185019093528083528493849392919083018282801561394857602002820191906000526020600020905b815481526020019060010190808311613934575b505050505090506117938185613e9e565b8260005b81811015611486576000868683818110613979576139796152df565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016139b291815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b81526004016139de92919061532d565b600060405180830381600087803b1580156139f857600080fd5b505af1158015613a0c573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa158015613a82573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613aa69190615446565b90508060ff16600103613cab577f000000000000000000000000000000000000000000000000000000000000000060ff871615613b0057507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb30898588604051602001613b3491815260200190565b604051602081830303815290604052613b4c90615c0a565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa158015613bd7573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052613c1d9190810190615393565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590613c76908790859060040161532d565b600060405180830381600087803b158015613c9057600080fd5b505af1158015613ca4573d6000803e3d6000fd5b5050505050505b50508080613cb890615463565b91505061395d565b8060005b81811015611ed6576000848483818110613ce057613ce06152df565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc8283604051602001613d1991815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613d4592919061532d565b600060405180830381600087803b158015613d5f57600080fd5b505af1158015613d73573d6000803e3d6000fd5b50505050508080613d8390615463565b915050613cc4565b600c6020528160005260406000208181548110613da757600080fd5b90600052602060002001600091509150505481565b613dc4614130565b613dcd816141c9565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613e2b908690600401615c4c565b6020604051808303816000875af1158015613e4a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e6e9190615611565b9050613e7b600f826142be565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b815160009081908190841580613eb45750808510155b15613ebd578094505b60008092505b85831015613f1957866001613ed88585615692565b613ee29190615692565b81518110613ef257613ef26152df565b602002602001015181613f0591906155dd565b905082613f1181615463565b935050613ec3565b9694955050505050565b82516000908190831580613f375750808410155b15613f40578093505b60008467ffffffffffffffff811115613f5b57613f5b6146b6565b604051908082528060200260200182016040528015613f84578160200160208202803683370190505b509050600092505b84831015613ff257866001613fa18585615692565b613fab9190615692565b81518110613fbb57613fbb6152df565b6020026020010151818481518110613fd557613fd56152df565b602090810291909101015282613fea81615463565b935050613f8c565b61400b816000600184516140069190615692565b6142ca565b856064036140445780600182516140229190615692565b81518110614032576140326152df565b60200260200101519350505050611cef565b8060648251886140549190615d9e565b61405e9190615e0a565b8151811061406e5761406e6152df565b602002602001015193505050509392505050565b6000611cef8383614442565b60007f00000000000000000000000000000000000000000000000000000000000000001561412b57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015614102573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906141269190615611565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff1633146141b1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016120ba565b565b6000611545825490565b6000611cef838361453c565b3373ffffffffffffffffffffffffffffffffffffffff821603614248576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016120ba565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611cef8383614566565b81818082036142da575050505050565b60008560026142e98787615e1e565b6142f39190615e3e565b6142fd9087615ea6565b8151811061430d5761430d6152df565b602002602001015190505b81831361441c575b80868481518110614333576143336152df565b60200260200101511015614353578261434b81615ece565b935050614320565b858281518110614365576143656152df565b6020026020010151811015614386578161437e81615eff565b925050614353565b8183136144175785828151811061439f5761439f6152df565b60200260200101518684815181106143b9576143b96152df565b60200260200101518785815181106143d3576143d36152df565b602002602001018885815181106143ec576143ec6152df565b6020908102919091010191909152528261440581615ece565b935050818061441390615eff565b9250505b614318565b8185121561442f5761442f8686846142ca565b83831215611486576114868684866142ca565b6000818152600183016020526040812054801561452b576000614466600183615692565b855490915060009061447a90600190615692565b90508181146144df57600086600001828154811061449a5761449a6152df565b90600052602060002001549050808760000184815481106144bd576144bd6152df565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806144f0576144f0615f56565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050611545565b6000915050611545565b5092915050565b6000826000018281548110614553576145536152df565b9060005260206000200154905092915050565b60008181526001830160205260408120546145ad57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155611545565b506000611545565b5080546000825590600052602060002090810190613dcd9190614629565b828054828255906000526020600020908101928215614619579160200282015b8281111561461957825182906146099082615a46565b50916020019190600101906145f3565b5061462592915061463e565b5090565b5b80821115614625576000815560010161462a565b80821115614625576000614652828261465b565b5060010161463e565b508054614667906156a5565b6000825580601f10614677575050565b601f016020900490600052602060002090810190613dcd9190614629565b60ff81168114613dcd57600080fd5b63ffffffff81168114613dcd57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff81118282101715614709576147096146b6565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715614756576147566146b6565b604052919050565b600067ffffffffffffffff821115614778576147786146b6565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f8301126147b557600080fd5b81356147c86147c38261475e565b61470f565b8181528460208386010111156147dd57600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff81168114613dcd57600080fd5b600080600080600080600060e0888a03121561482f57600080fd5b873561483a81614695565b9650602088013561484a816146a4565b9550604088013561485a81614695565b9450606088013567ffffffffffffffff81111561487657600080fd5b6148828a828b016147a4565b9450506080880135614893816147fa565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff811681146148c257600080fd5b919050565b6000806000606084860312156148dc57600080fd5b833592506148ec602085016148b0565b9150604084013590509250925092565b6000806040838503121561490f57600080fd5b82359150602083013567ffffffffffffffff81111561492d57600080fd5b614939858286016147a4565b9150509250929050565b60006020828403121561495557600080fd5b5035919050565b60005b8381101561497757818101518382015260200161495f565b50506000910152565b6000815180845261499881602086016020860161495c565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611cef6020830184614980565b73ffffffffffffffffffffffffffffffffffffffff81168114613dcd57600080fd5b600080600080600080600060e0888a031215614a1a57600080fd5b873596506020880135614a2c816149dd565b95506040880135614a3c81614695565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b60008060408385031215614a7757600080fd5b82359150614a87602084016148b0565b90509250929050565b600060208284031215614aa257600080fd5b8135611cef816149dd565b60008083601f840112614abf57600080fd5b50813567ffffffffffffffff811115614ad757600080fd5b6020830191508360208260051b850101111561179957600080fd5b600080600080600080600060c0888a031215614b0d57600080fd5b873567ffffffffffffffff811115614b2457600080fd5b614b308a828b01614aad565b9098509650506020880135614b4481614695565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b600080600060608486031215614b8157600080fd5b505081359360208301359350604090920135919050565b60208152614bbf60208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151614bd8604084018263ffffffff169052565b506040830151610140806060850152614bf5610160850183614980565b91506060850151614c1660808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e0850151610100614c82818701836bffffffffffffffffffffffff169052565b8601519050610120614c978682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001838701529050614ccf8382614980565b9695505050505050565b60008060208385031215614cec57600080fd5b823567ffffffffffffffff811115614d0357600080fd5b614d0f85828601614aad565b90969095509350505050565b60008060408385031215614d2e57600080fd5b823567ffffffffffffffff80821115614d4657600080fd5b908401906101008287031215614d5b57600080fd5b90925060208401359080821115614d7157600080fd5b50614939858286016147a4565b8215158152604060208201526000611cec6040830184614980565b60008083601f840112614dab57600080fd5b50813567ffffffffffffffff811115614dc357600080fd5b60208301915083602082850101111561179957600080fd5b60008060208385031215614dee57600080fd5b823567ffffffffffffffff811115614e0557600080fd5b614d0f85828601614d99565b6020808252825182820181905260009190848201906040850190845b81811015614e4957835183529284019291840191600101614e2d565b50909695505050505050565b600067ffffffffffffffff821115614e6f57614e6f6146b6565b5060051b60200190565b60008060408385031215614e8c57600080fd5b823567ffffffffffffffff80821115614ea457600080fd5b818501915085601f830112614eb857600080fd5b81356020614ec86147c383614e55565b82815260059290921b84018101918181019089841115614ee757600080fd5b8286015b84811015614f1f57803586811115614f035760008081fd5b614f118c86838b01016147a4565b845250918301918301614eeb565b5096505086013592505080821115614d7157600080fd5b600080600060408486031215614f4b57600080fd5b833567ffffffffffffffff811115614f6257600080fd5b614f6e86828701614aad565b9094509250506020840135614f82816146a4565b809150509250925092565b60008060408385031215614fa057600080fd5b50508035926020909101359150565b600080600060408486031215614fc457600080fd5b83359250602084013567ffffffffffffffff811115614fe257600080fd5b614fee86828701614d99565b9497909650939450505050565b6000806040838503121561500e57600080fd5b823567ffffffffffffffff8082111561502657600080fd5b615032868387016147a4565b93506020850135915080821115614d7157600080fd5b6000602080838503121561505b57600080fd5b823567ffffffffffffffff8082111561507357600080fd5b818501915085601f83011261508757600080fd5b81356150956147c382614e55565b81815260059190911b830184019084810190888311156150b457600080fd5b8585015b838110156150ec578035858111156150d05760008081fd5b6150de8b89838a01016147a4565b8452509186019186016150b8565b5098975050505050505050565b6000806040838503121561510c57600080fd5b82359150602083013561511e816147fa565b809150509250929050565b60006020828403121561513b57600080fd5b8135611cef81614695565b6000806040838503121561515957600080fd5b82359150602083013561511e816146a4565b6000806040838503121561517e57600080fd5b82359150602083013561511e81614695565b600080600080606085870312156151a657600080fd5b843567ffffffffffffffff8111156151bd57600080fd5b6151c987828801614aad565b90955093505060208501356151dd81614695565b915060408501356151ed81614695565b939692955090935050565b60008060008060008060c0878903121561521157600080fd5b863561521c816149dd565b9550602087013561522c81614695565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff808316818516818304811182151516156152ab576152ab615251565b02949350505050565b805180151581146148c257600080fd5b6000602082840312156152d657600080fd5b611cef826152b4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff810361532457615324615251565b60010192915050565b828152604060208201526000611cec6040830184614980565b600082601f83011261535757600080fd5b81516153656147c38261475e565b81815284602083860101111561537a57600080fd5b61538b82602083016020870161495c565b949350505050565b6000602082840312156153a557600080fd5b815167ffffffffffffffff8111156153bc57600080fd5b61538b84828501615346565b80516148c2816147fa565b6000602082840312156153e557600080fd5b8151611cef816147fa565b80516148c2816149dd565b6000806040838503121561540e57600080fd5b8251615419816149dd565b6020939093015192949293505050565b60006020828403121561543b57600080fd5b8151611cef816149dd565b60006020828403121561545857600080fd5b8151611cef81614695565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361549457615494615251565b5060010190565b80516148c2816146a4565b805167ffffffffffffffff811681146148c257600080fd5b6000602082840312156154d057600080fd5b815167ffffffffffffffff808211156154e857600080fd5b9083019061014082860312156154fd57600080fd5b6155056146e5565b61550e836153f0565b815261551c6020840161549b565b602082015260408301518281111561553357600080fd5b61553f87828601615346565b604083015250615551606084016153c8565b6060820152615562608084016153f0565b608082015261557360a084016154a6565b60a082015261558460c0840161549b565b60c082015261559560e084016153c8565b60e08201526101006155a88185016152b4565b9082015261012083810151838111156155c057600080fd5b6155cc88828701615346565b918301919091525095945050505050565b8082018082111561154557611545615251565b600061ffff80831681810361560757615607615251565b6001019392505050565b60006020828403121561562357600080fd5b5051919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261565f57600080fd5b83018035915067ffffffffffffffff82111561567a57600080fd5b6020019150600581901b360382131561179957600080fd5b8181038181111561154557611545615251565b600181811c908216806156b957607f821691505b6020821081036156f2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000808354615706816156a5565b6001828116801561571e576001811461575157615780565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168752821515830287019450615780565b8760005260208060002060005b858110156157775781548a82015290840190820161575e565b50505082870194505b50929695505050505050565b60008154615799816156a5565b8085526020600183811680156157b657600181146157ee5761581c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008516838901528284151560051b890101955061581c565b866000528260002060005b858110156158145781548a82018601529083019084016157f9565b890184019650505b505050505092915050565b60a08152600061583a60a083018861578c565b6020838203818501528188548084528284019150828160051b8501018a6000528360002060005b838110156158ac577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe087840301855261589a838361578c565b94860194925060019182019101615861565b505086810360408801526158c0818b61578c565b94505050505084606084015282810360808401526158de8185614980565b98975050505050505050565b6000604082016040835280855180835260608501915060608160051b8601019250602080880160005b8381101561595f577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa088870301855261594d868351614980565b95509382019390820190600101615913565b5050858403818701525050506159758185614980565b95945050505050565b60008060006060848603121561599357600080fd5b83519250602084015191506040840151614f82816149dd565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b601f821115612d8457600081815260208120601f850160051c81016020861015615a275750805b601f850160051c820191505b8181101561148657828155600101615a33565b815167ffffffffffffffff811115615a6057615a606146b6565b615a7481615a6e84546156a5565b84615a00565b602080601f831160018114615ac75760008415615a915750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555611486565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015615b1457888601518255948401946001909101908401615af5565b5085821015615b5057878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60006020808385031215615b7357600080fd5b825167ffffffffffffffff811115615b8a57600080fd5b8301601f81018513615b9b57600080fd5b8051615ba96147c382614e55565b81815260059190911b82018301908381019087831115615bc857600080fd5b928401925b82841015615be657835182529284019290840190615bcd565b979650505050505050565b600063ffffffff80831681810361560757615607615251565b805160208083015191908110156156f2577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b6020815260008251610140806020850152615c6b610160850183614980565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe080868503016040870152615ca78483614980565b935060408701519150615cd2606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e0870152615d338483614980565b935060e08701519150610100818786030181880152615d528584614980565b945080880151925050610120818786030181880152615d718584614980565b94508088015192505050615d94828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615dd657615dd6615251565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082615e1957615e19615ddb565b500490565b818103600083128015838313168383128216171561453557614535615251565b600082615e4d57615e4d615ddb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f800000000000000000000000000000000000000000000000000000000000000083141615615ea157615ea1615251565b500590565b8082018281126000831280158216821582161715615ec657615ec6615251565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361549457615494615251565b60007f80000000000000000000000000000000000000000000000000000000000000008203615f3057615f30615251565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
}

var VerifiableLoadLogTriggerUpkeepABI = VerifiableLoadLogTriggerUpkeepMetaData.ABI

var VerifiableLoadLogTriggerUpkeepBin = VerifiableLoadLogTriggerUpkeepMetaData.Bin

func DeployVerifiableLoadLogTriggerUpkeep(auth *bind.TransactOpts, backend bind.ContractBackend, _registrar common.Address, _useArb bool, _useMercury bool) (common.Address, *types.Transaction, *VerifiableLoadLogTriggerUpkeep, error) {
	parsed, err := VerifiableLoadLogTriggerUpkeepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifiableLoadLogTriggerUpkeepBin), backend, _registrar, _useArb, _useMercury)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifiableLoadLogTriggerUpkeep{VerifiableLoadLogTriggerUpkeepCaller: VerifiableLoadLogTriggerUpkeepCaller{contract: contract}, VerifiableLoadLogTriggerUpkeepTransactor: VerifiableLoadLogTriggerUpkeepTransactor{contract: contract}, VerifiableLoadLogTriggerUpkeepFilterer: VerifiableLoadLogTriggerUpkeepFilterer{contract: contract}}, nil
}

type VerifiableLoadLogTriggerUpkeep struct {
	address common.Address
	abi     abi.ABI
	VerifiableLoadLogTriggerUpkeepCaller
	VerifiableLoadLogTriggerUpkeepTransactor
	VerifiableLoadLogTriggerUpkeepFilterer
}

type VerifiableLoadLogTriggerUpkeepCaller struct {
	contract *bind.BoundContract
}

type VerifiableLoadLogTriggerUpkeepTransactor struct {
	contract *bind.BoundContract
}

type VerifiableLoadLogTriggerUpkeepFilterer struct {
	contract *bind.BoundContract
}

type VerifiableLoadLogTriggerUpkeepSession struct {
	Contract     *VerifiableLoadLogTriggerUpkeep
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VerifiableLoadLogTriggerUpkeepCallerSession struct {
	Contract *VerifiableLoadLogTriggerUpkeepCaller
	CallOpts bind.CallOpts
}

type VerifiableLoadLogTriggerUpkeepTransactorSession struct {
	Contract     *VerifiableLoadLogTriggerUpkeepTransactor
	TransactOpts bind.TransactOpts
}

type VerifiableLoadLogTriggerUpkeepRaw struct {
	Contract *VerifiableLoadLogTriggerUpkeep
}

type VerifiableLoadLogTriggerUpkeepCallerRaw struct {
	Contract *VerifiableLoadLogTriggerUpkeepCaller
}

type VerifiableLoadLogTriggerUpkeepTransactorRaw struct {
	Contract *VerifiableLoadLogTriggerUpkeepTransactor
}

func NewVerifiableLoadLogTriggerUpkeep(address common.Address, backend bind.ContractBackend) (*VerifiableLoadLogTriggerUpkeep, error) {
	abi, err := abi.JSON(strings.NewReader(VerifiableLoadLogTriggerUpkeepABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeep{address: address, abi: abi, VerifiableLoadLogTriggerUpkeepCaller: VerifiableLoadLogTriggerUpkeepCaller{contract: contract}, VerifiableLoadLogTriggerUpkeepTransactor: VerifiableLoadLogTriggerUpkeepTransactor{contract: contract}, VerifiableLoadLogTriggerUpkeepFilterer: VerifiableLoadLogTriggerUpkeepFilterer{contract: contract}}, nil
}

func NewVerifiableLoadLogTriggerUpkeepCaller(address common.Address, caller bind.ContractCaller) (*VerifiableLoadLogTriggerUpkeepCaller, error) {
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepCaller{contract: contract}, nil
}

func NewVerifiableLoadLogTriggerUpkeepTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifiableLoadLogTriggerUpkeepTransactor, error) {
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepTransactor{contract: contract}, nil
}

func NewVerifiableLoadLogTriggerUpkeepFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifiableLoadLogTriggerUpkeepFilterer, error) {
	contract, err := bindVerifiableLoadLogTriggerUpkeep(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepFilterer{contract: contract}, nil
}

func bindVerifiableLoadLogTriggerUpkeep(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifiableLoadLogTriggerUpkeepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadLogTriggerUpkeep.Contract.VerifiableLoadLogTriggerUpkeepCaller.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.VerifiableLoadLogTriggerUpkeepTransactor.contract.Transfer(opts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.VerifiableLoadLogTriggerUpkeepTransactor.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifiableLoadLogTriggerUpkeep.Contract.contract.Call(opts, result, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.contract.Transfer(opts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.contract.Transact(opts, method, params...)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) BUCKETSIZE(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "BUCKET_SIZE")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) BUCKETSIZE() (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BUCKETSIZE(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) AddLinkAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "addLinkAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddLinkAmount(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) AddLinkAmount() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddLinkAmount(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) BucketedDelays(opts *bind.CallOpts, arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "bucketedDelays", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) BucketedDelays(arg0 *big.Int, arg1 uint16, arg2 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1, arg2)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "buckets", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Buckets(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Buckets(arg0 *big.Int) (uint16, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Buckets(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) CheckCallback(opts *bind.CallOpts, values [][]byte, extraData []byte) (bool, []byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "checkCallback", values, extraData)

	if err != nil {
		return *new(bool), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckCallback(&_VerifiableLoadLogTriggerUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) CheckCallback(values [][]byte, extraData []byte) (bool, []byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckCallback(&_VerifiableLoadLogTriggerUpkeep.CallOpts, values, extraData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) CheckGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "checkGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) CheckGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Counters(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "counters", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Counters(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Counters(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Counters(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Delays(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "delays", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Delays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Delays(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Delays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0, arg1)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) DummyMap(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "dummyMap", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.DummyMap(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) DummyMap(arg0 [32]byte) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.DummyMap(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Eligible(opts *bind.CallOpts, upkeepId *big.Int) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "eligible", upkeepId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Eligible(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Eligible(upkeepId *big.Int) (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Eligible(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) EmittedAgainSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "emittedAgainSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) EmittedAgainSig() ([32]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.EmittedAgainSig(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) EmittedSig(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "emittedSig")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.EmittedSig(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) EmittedSig() ([32]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.EmittedSig(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) FeedParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "feedParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) FeedParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) FeedParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) FeedsHex(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "feedsHex", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedsHex(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) FeedsHex(arg0 *big.Int) (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FeedsHex(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) FirstPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "firstPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) FirstPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.FirstPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GasLimits(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "gasLimits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GasLimits(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GasLimits(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GasLimits(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetActiveUpkeepIDsDeployedByThisContract(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getActiveUpkeepIDsDeployedByThisContract", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadLogTriggerUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetActiveUpkeepIDsDeployedByThisContract(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetActiveUpkeepIDsDeployedByThisContract(&_VerifiableLoadLogTriggerUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetAllActiveUpkeepIDsOnRegistry(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getAllActiveUpkeepIDsOnRegistry", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadLogTriggerUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetAllActiveUpkeepIDsOnRegistry(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetAllActiveUpkeepIDsOnRegistry(&_VerifiableLoadLogTriggerUpkeep.CallOpts, startIndex, maxCount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetBalance(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getBalance", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBalance(&_VerifiableLoadLogTriggerUpkeep.CallOpts, id)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetBalance(id *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBalance(&_VerifiableLoadLogTriggerUpkeep.CallOpts, id)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getBucketedDelays", upkeepId, bucket)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetBucketedDelays(upkeepId *big.Int, bucket uint16) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetBucketedDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getBucketedDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetBucketedDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetBucketedDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetDelays(opts *bind.CallOpts, upkeepId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getDelays", upkeepId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetDelays(upkeepId *big.Int) ([]*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelays(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetDelaysLength(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getDelaysLength", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetDelaysLength(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetDelaysLength(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetForwarder(opts *bind.CallOpts, upkeepID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getForwarder", upkeepID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetForwarder(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetForwarder(upkeepID *big.Int) (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetForwarder(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepID)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetLogTriggerConfig(opts *bind.CallOpts, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getLogTriggerConfig", addr, selector, topic0, topic1, topic2, topic3)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetLogTriggerConfig(addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetLogTriggerConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetMinBalanceForUpkeep(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getMinBalanceForUpkeep", upkeepId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetMinBalanceForUpkeep(upkeepId *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetMinBalanceForUpkeep(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetPxDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getPxDelayLastNPerforms", upkeepId, p, n)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetPxDelayLastNPerforms(upkeepId *big.Int, p *big.Int, n *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetPxDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, p, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetSumDelayInBucket(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getSumDelayInBucket", upkeepId, bucket)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetSumDelayInBucket(upkeepId *big.Int, bucket uint16) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayInBucket(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, bucket)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetSumDelayLastNPerforms(opts *bind.CallOpts, upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getSumDelayLastNPerforms", upkeepId, n)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetSumDelayLastNPerforms(upkeepId *big.Int, n *big.Int) (*big.Int, *big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetSumDelayLastNPerforms(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId, n)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetTriggerType(opts *bind.CallOpts, upkeepId *big.Int) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getTriggerType", upkeepId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetTriggerType(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetTriggerType(upkeepId *big.Int) (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetTriggerType(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetUpkeepInfo(opts *bind.CallOpts, upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getUpkeepInfo", upkeepId)

	if err != nil {
		return *new(KeeperRegistryBase21UpkeepInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(KeeperRegistryBase21UpkeepInfo)).(*KeeperRegistryBase21UpkeepInfo)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetUpkeepInfo(upkeepId *big.Int) (KeeperRegistryBase21UpkeepInfo, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetUpkeepInfo(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getUpkeepPrivilegeConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetUpkeepPrivilegeConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "getUpkeepTriggerConfig", upkeepId)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) GetUpkeepTriggerConfig(upkeepId *big.Int) ([]byte, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.GetUpkeepTriggerConfig(&_VerifiableLoadLogTriggerUpkeep.CallOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Intervals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "intervals", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Intervals(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Intervals(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Intervals(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) LastTopUpBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "lastTopUpBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) LastTopUpBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LastTopUpBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LinkToken(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) LinkToken() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LinkToken(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) LogNum(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "logNum")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) LogNum() (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LogNum(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) LogNum() (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.LogNum(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) MinBalanceThresholdMultiplier(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "minBalanceThresholdMultiplier")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) MinBalanceThresholdMultiplier() (uint8, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.MinBalanceThresholdMultiplier(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Owner() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Owner(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Owner() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Owner(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) PerformDataSizes(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "performDataSizes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformDataSizes(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) PerformDataSizes(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformDataSizes(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) PerformGasToBurns(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "performGasToBurns", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) PerformGasToBurns(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformGasToBurns(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) PreviousPerformBlocks(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "previousPerformBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) PreviousPerformBlocks(arg0 *big.Int) (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PreviousPerformBlocks(&_VerifiableLoadLogTriggerUpkeep.CallOpts, arg0)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Registrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "registrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Registrar() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registrar(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Registrar() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registrar(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Registry() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registry(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) Registry() (common.Address, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Registry(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) TimeParamKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "timeParamKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TimeParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimeParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) TimeParamKey() (string, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TimeParamKey(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) UpkeepTopUpCheckInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "upkeepTopUpCheckInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) UpkeepTopUpCheckInterval() (*big.Int, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpkeepTopUpCheckInterval(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) UseArbitrumBlockNum(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "useArbitrumBlockNum")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) UseArbitrumBlockNum() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseArbitrumBlockNum(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCaller) UseMercury(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VerifiableLoadLogTriggerUpkeep.contract.Call(opts, &out, "useMercury")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UseMercury() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseMercury(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepCallerSession) UseMercury() (bool, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UseMercury(&_VerifiableLoadLogTriggerUpkeep.CallOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "acceptOwnership")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AcceptOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AcceptOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) AddFunds(opts *bind.TransactOpts, upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "addFunds", upkeepId, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddFunds(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) AddFunds(upkeepId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.AddFunds(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, amount)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchCancelUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchCancelUpkeeps", upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchCancelUpkeeps(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchCancelUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchPreparingUpkeeps(opts *bind.TransactOpts, upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchPreparingUpkeeps", upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchPreparingUpkeeps(upkeepIds []*big.Int, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchPreparingUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchPreparingUpkeepsSimple(opts *bind.TransactOpts, upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchPreparingUpkeepsSimple", upkeepIds, log, selector)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, log, selector)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchPreparingUpkeepsSimple(upkeepIds []*big.Int, log uint8, selector uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchPreparingUpkeepsSimple(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, log, selector)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchRegisterUpkeeps(opts *bind.TransactOpts, number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchRegisterUpkeeps", number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchRegisterUpkeeps(number uint8, gasLimit uint32, triggerType uint8, triggerConfig []byte, amount *big.Int, checkGasToBurn *big.Int, performGasToBurn *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchRegisterUpkeeps(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, number, gasLimit, triggerType, triggerConfig, amount, checkGasToBurn, performGasToBurn)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchSendLogs(opts *bind.TransactOpts, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchSendLogs", log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSendLogs(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchSendLogs(log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSendLogs(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchSetIntervals(opts *bind.TransactOpts, upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchSetIntervals", upkeepIds, interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchSetIntervals(upkeepIds []*big.Int, interval uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchSetIntervals(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds, interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchUpdatePipelineData(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchUpdatePipelineData", upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchUpdatePipelineData(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchUpdatePipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BatchWithdrawLinks(opts *bind.TransactOpts, upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "batchWithdrawLinks", upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BatchWithdrawLinks(upkeepIds []*big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BatchWithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepIds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) BurnPerformGas(opts *bind.TransactOpts, upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "burnPerformGas", upkeepId, startGas, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BurnPerformGas(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) BurnPerformGas(upkeepId *big.Int, startGas *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.BurnPerformGas(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, startGas, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) CheckLog(opts *bind.TransactOpts, log Log, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "checkLog", log, checkData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) CheckLog(log Log, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log, checkData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) CheckLog(log Log, checkData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.CheckLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, log, checkData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "performUpkeep", performData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.PerformUpkeep(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, performData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SendLog(opts *bind.TransactOpts, upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "sendLog", upkeepId, log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SendLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SendLog(upkeepId *big.Int, log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SendLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setConfig", newRegistrar)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetConfig(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetConfig(newRegistrar common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetConfig(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, newRegistrar)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setFeeds", _feeds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetFeeds(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _feeds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetFeeds(_feeds []string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetFeeds(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _feeds)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setInterval", upkeepId, _interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetInterval(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetInterval(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetLog(opts *bind.TransactOpts, _log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setLog", _log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetLog(_log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetLog(_log uint8) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetLog(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _log)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetParamKeys(opts *bind.TransactOpts, _feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setParamKeys", _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetParamKeys(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetParamKeys(_feedParamKey string, _timeParamKey string) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetParamKeys(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, _feedParamKey, _timeParamKey)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetPerformDataSize(opts *bind.TransactOpts, upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setPerformDataSize", upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetPerformDataSize(upkeepId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetPerformDataSize(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, value)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetUpkeepGasLimit(opts *bind.TransactOpts, upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setUpkeepGasLimit", upkeepId, gasLimit)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetUpkeepGasLimit(upkeepId *big.Int, gasLimit uint32) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepGasLimit(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, gasLimit)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) SetUpkeepPrivilegeConfig(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "setUpkeepPrivilegeConfig", upkeepId, cfg)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) SetUpkeepPrivilegeConfig(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.SetUpkeepPrivilegeConfig(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) TopUpFund(opts *bind.TransactOpts, upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "topUpFund", upkeepId, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TopUpFund(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) TopUpFund(upkeepId *big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TopUpFund(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, blockNum)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "transferOwnership", to)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TransferOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, to)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.TransferOwnership(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, to)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) UpdateLogTriggerConfig1(opts *bind.TransactOpts, upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "updateLogTriggerConfig1", upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) UpdateLogTriggerConfig1(upkeepId *big.Int, addr common.Address, selector uint8, topic0 [32]byte, topic1 [32]byte, topic2 [32]byte, topic3 [32]byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateLogTriggerConfig1(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, addr, selector, topic0, topic1, topic2, topic3)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) UpdateLogTriggerConfig2(opts *bind.TransactOpts, upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "updateLogTriggerConfig2", upkeepId, cfg)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) UpdateLogTriggerConfig2(upkeepId *big.Int, cfg []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateLogTriggerConfig2(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, cfg)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) UpdateUpkeepPipelineData(opts *bind.TransactOpts, upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "updateUpkeepPipelineData", upkeepId, pipelineData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) UpdateUpkeepPipelineData(upkeepId *big.Int, pipelineData []byte) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.UpdateUpkeepPipelineData(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId, pipelineData)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) WithdrawLinks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "withdrawLinks")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) WithdrawLinks() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) WithdrawLinks0(opts *bind.TransactOpts, upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.Transact(opts, "withdrawLinks0", upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) WithdrawLinks0(upkeepId *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.WithdrawLinks0(&_VerifiableLoadLogTriggerUpkeep.TransactOpts, upkeepId)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.contract.RawTransact(opts, nil)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Receive(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepTransactorSession) Receive() (*types.Transaction, error) {
	return _VerifiableLoadLogTriggerUpkeep.Contract.Receive(&_VerifiableLoadLogTriggerUpkeep.TransactOpts)
}

type VerifiableLoadLogTriggerUpkeepLogEmittedIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepLogEmitted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepLogEmitted)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepLogEmitted)
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

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepLogEmitted struct {
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadLogTriggerUpkeepLogEmittedIterator, error) {

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

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepLogEmittedIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "LogEmitted", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "LogEmitted", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepLogEmitted)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseLogEmitted(log types.Log) (*VerifiableLoadLogTriggerUpkeepLogEmitted, error) {
	event := new(VerifiableLoadLogTriggerUpkeepLogEmitted)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "LogEmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepLogEmittedAgainIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepLogEmittedAgain

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedAgainIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepLogEmittedAgain)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepLogEmittedAgain)
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

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedAgainIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepLogEmittedAgainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepLogEmittedAgain struct {
	UpkeepId *big.Int
	BlockNum *big.Int
	Addr     common.Address
	Raw      types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadLogTriggerUpkeepLogEmittedAgainIterator, error) {

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

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepLogEmittedAgainIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "LogEmittedAgain", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "LogEmittedAgain", upkeepIdRule, blockNumRule, addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepLogEmittedAgain)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseLogEmittedAgain(log types.Log) (*VerifiableLoadLogTriggerUpkeepLogEmittedAgain, error) {
	event := new(VerifiableLoadLogTriggerUpkeepLogEmittedAgain)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "LogEmittedAgain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
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

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, error) {
	event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
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

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseOwnershipTransferred(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferred, error) {
	event := new(VerifiableLoadLogTriggerUpkeepOwnershipTransferred)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator struct {
	Event *VerifiableLoadLogTriggerUpkeepUpkeepTopUp

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
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
		it.Event = new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
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

func (it *VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator) Error() error {
	return it.fail
}

func (it *VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VerifiableLoadLogTriggerUpkeepUpkeepTopUp struct {
	UpkeepId *big.Int
	Amount   *big.Int
	BlockNum *big.Int
	Raw      types.Log
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.FilterLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return &VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator{contract: _VerifiableLoadLogTriggerUpkeep.contract, event: "UpkeepTopUp", logs: logs, sub: sub}, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepTopUp) (event.Subscription, error) {

	logs, sub, err := _VerifiableLoadLogTriggerUpkeep.contract.WatchLogs(opts, "UpkeepTopUp")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
				if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
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

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeepFilterer) ParseUpkeepTopUp(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUp, error) {
	event := new(VerifiableLoadLogTriggerUpkeepUpkeepTopUp)
	if err := _VerifiableLoadLogTriggerUpkeep.contract.UnpackLog(event, "UpkeepTopUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeep) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["LogEmitted"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseLogEmitted(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["LogEmittedAgain"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseLogEmittedAgain(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["OwnershipTransferRequested"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseOwnershipTransferRequested(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["OwnershipTransferred"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseOwnershipTransferred(log)
	case _VerifiableLoadLogTriggerUpkeep.abi.Events["UpkeepTopUp"].ID:
		return _VerifiableLoadLogTriggerUpkeep.ParseUpkeepTopUp(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VerifiableLoadLogTriggerUpkeepLogEmitted) Topic() common.Hash {
	return common.HexToHash("0x97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf08")
}

func (VerifiableLoadLogTriggerUpkeepLogEmittedAgain) Topic() common.Hash {
	return common.HexToHash("0xc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d")
}

func (VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VerifiableLoadLogTriggerUpkeepOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VerifiableLoadLogTriggerUpkeepUpkeepTopUp) Topic() common.Hash {
	return common.HexToHash("0x49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c0")
}

func (_VerifiableLoadLogTriggerUpkeep *VerifiableLoadLogTriggerUpkeep) Address() common.Address {
	return _VerifiableLoadLogTriggerUpkeep.address
}

type VerifiableLoadLogTriggerUpkeepInterface interface {
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

	LogNum(opts *bind.CallOpts) (uint8, error)

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

	UseMercury(opts *bind.CallOpts) (bool, error)

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

	CheckLog(opts *bind.TransactOpts, log Log, checkData []byte) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	SendLog(opts *bind.TransactOpts, upkeepId *big.Int, log uint8) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, newRegistrar common.Address) (*types.Transaction, error)

	SetFeeds(opts *bind.TransactOpts, _feeds []string) (*types.Transaction, error)

	SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error)

	SetLog(opts *bind.TransactOpts, _log uint8) (*types.Transaction, error)

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

	FilterLogEmitted(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadLogTriggerUpkeepLogEmittedIterator, error)

	WatchLogEmitted(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepLogEmitted, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmitted(log types.Log) (*VerifiableLoadLogTriggerUpkeepLogEmitted, error)

	FilterLogEmittedAgain(opts *bind.FilterOpts, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (*VerifiableLoadLogTriggerUpkeepLogEmittedAgainIterator, error)

	WatchLogEmittedAgain(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepLogEmittedAgain, upkeepId []*big.Int, blockNum []*big.Int, addr []common.Address) (event.Subscription, error)

	ParseLogEmittedAgain(log types.Log) (*VerifiableLoadLogTriggerUpkeepLogEmittedAgain, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VerifiableLoadLogTriggerUpkeepOwnershipTransferred, error)

	FilterUpkeepTopUp(opts *bind.FilterOpts) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUpIterator, error)

	WatchUpkeepTopUp(opts *bind.WatchOpts, sink chan<- *VerifiableLoadLogTriggerUpkeepUpkeepTopUp) (event.Subscription, error)

	ParseUpkeepTopUp(log types.Log) (*VerifiableLoadLogTriggerUpkeepUpkeepTopUp, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
