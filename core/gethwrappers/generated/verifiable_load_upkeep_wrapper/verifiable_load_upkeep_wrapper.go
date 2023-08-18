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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_useArb\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogEmittedAgain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"UpkeepTopUp\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BUCKET_SIZE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"addFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addLinkAmount\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"batchPreparingUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"}],\"name\":\"batchPreparingUpkeepsSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"number\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"triggerType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"triggerConfig\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"checkGasToBurn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"performGasToBurn\",\"type\":\"uint256\"}],\"name\":\"batchRegisterUpkeeps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"batchSendLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"interval\",\"type\":\"uint32\"}],\"name\":\"batchSetIntervals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchUpdatePipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"upkeepIds\",\"type\":\"uint256[]\"}],\"name\":\"batchWithdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"bucketedDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"buckets\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"burnPerformGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"checkGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"counters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"dummyMap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"eligible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedAgainSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emittedSig\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feedParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"feedsHex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"firstPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasLimits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveUpkeepIDsDeployedByThisContract\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getAllActiveUpkeepIDsOnRegistry\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getBucketedDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getBucketedDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelays\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getDelaysLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepID\",\"type\":\"uint256\"}],\"name\":\"getForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"getLogTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"logTrigger\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getMinBalanceForUpkeep\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getPxDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"bucket\",\"type\":\"uint16\"}],\"name\":\"getSumDelayInBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"getSumDelayLastNPerforms\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getTriggerType\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"performGas\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"maxValidBlocknumber\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lastPerformedBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"amountSpent\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"internalType\":\"structKeeperRegistryBase2_1.UpkeepInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepPrivilegeConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"getUpkeepTriggerConfig\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"intervals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastTopUpBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBalanceThresholdMultiplier\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performDataSizes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"performGasToBurns\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"previousPerformBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIKeeperRegistryMaster\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"log\",\"type\":\"uint8\"}],\"name\":\"sendLog\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAutomationRegistrar2_1\",\"name\":\"newRegistrar\",\"type\":\"address\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interval\",\"type\":\"uint256\"}],\"name\":\"setInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPerformDataSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"}],\"name\":\"setUpkeepGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"setUpkeepPrivilegeConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeParamKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"topUpFund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"selector\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic2\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"topic3\",\"type\":\"bytes32\"}],\"name\":\"updateLogTriggerConfig1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cfg\",\"type\":\"bytes\"}],\"name\":\"updateLogTriggerConfig2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pipelineData\",\"type\":\"bytes\"}],\"name\":\"updateUpkeepPipelineData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upkeepTopUpCheckInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"useArbitrumBlockNum\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"upkeepId\",\"type\":\"uint256\"}],\"name\":\"withdrawLinks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x7f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf086080527fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d60a0526005601455601580546001600160681b0319166c140000000002c68af0bb140000179055606460e0526101c06040526042610140818152610100918291906200592461016039815260200160405180608001604052806042815260200162005966604291399052620000be90601690600262000365565b50348015620000cc57600080fd5b50604051620059a8380380620059a8833981016040819052620000ef9162000452565b81813380600081620001485760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156200017b576200017b81620002ba565b5050601180546001600160a01b0319166001600160a01b038516908117909155604080516330fe427560e21b815281516000945063c3f909d4926004808401939192918290030181865afa158015620001d8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001fe919062000495565b50601380546001600160a01b0319166001600160a01b038381169190911790915560115460408051631b6b6d2360e01b81529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801562000264573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200028a9190620004c6565b601280546001600160a01b0319166001600160a01b039290921691909117905550151560c052506200065e915050565b336001600160a01b03821603620003145760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016200013f565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215620003b0579160200282015b82811115620003b057825182906200039f908262000592565b509160200191906001019062000386565b50620003be929150620003c2565b5090565b80821115620003be576000620003d98282620003e3565b50600101620003c2565b508054620003f19062000503565b6000825580601f1062000402575050565b601f01602090049060005260206000209081019062000422919062000425565b50565b5b80821115620003be576000815560010162000426565b6001600160a01b03811681146200042257600080fd5b600080604083850312156200046657600080fd5b825162000473816200043c565b602084015190925080151581146200048a57600080fd5b809150509250929050565b60008060408385031215620004a957600080fd5b8251620004b6816200043c565b6020939093015192949293505050565b600060208284031215620004d957600080fd5b8151620004e6816200043c565b9392505050565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200051857607f821691505b6020821081036200053957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200058d57600081815260208120601f850160051c81016020861015620005685750805b601f850160051c820191505b81811015620005895782815560010162000574565b5050505b505050565b81516001600160401b03811115620005ae57620005ae620004ed565b620005c681620005bf845462000503565b846200053f565b602080601f831160018114620005fe5760008415620005e55750858301515b600019600386901b1c1916600185901b17855562000589565b600085815260208120601f198616915b828110156200062f578886015182559484019460019091019084016200060e565b50858210156200064e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a05160c05160e051615270620006b4600039600081816105170152611f1701526000818161093d0152613ab80152600081816107d60152613506015260008181610cc801526134db01526152706000f3fe6080604052600436106104845760003560e01c806379ba50971161025e578063a72aa27e11610143578063da6cba47116100bb578063e45530831161008a578063fa333dfb1161006f578063fa333dfb14610f3f578063fba7ffa314610ff2578063fcdc1f631461101f57600080fd5b8063e455308314610f09578063f2fde38b14610f1f57600080fd5b8063da6cba4714610e7c578063daee1aeb14610e9c578063dbef701e14610ebc578063e0114adb14610edc57600080fd5b8063becde0e111610112578063c98f10b0116100f7578063c98f10b014610df3578063d4c2490014610e3c578063d6051a7214610e5c57600080fd5b8063becde0e114610db3578063c041982214610dd357600080fd5b8063a72aa27e14610d0a578063af953a4a14610d2a578063afb28d1f14610d4a578063b657bc9c14610d9357600080fd5b8063948108f7116101d65780639b51fb0d116101a55780639d6f1cc71161018a5780639d6f1cc714610c96578063a654824814610cb6578063a6b5947514610cea57600080fd5b80639b51fb0d14610c455780639d385eaa14610c7657600080fd5b8063948108f714610bae57806396cebc7c14610bce5780639ac542eb14610bee5780639b42935414610c1857600080fd5b80638243444a1161022d5780638da5cb5b116102125780638da5cb5b14610b365780638fcb3fba14610b61578063924ca57814610b8e57600080fd5b80638243444a14610af6578063873c758614610b1657600080fd5b806379ba509714610a7457806379ea994314610a895780637b10399914610aa95780637e7a46dc14610ad657600080fd5b80634585e33b1161038457806360457ff5116102fc5780636e04ff0d116102cb57806373644cce116102b057806373644cce146109fa5780637672130314610a27578063776898c814610a5457600080fd5b80636e04ff0d1461099c5780637145f11b146109ca57600080fd5b806360457ff5146108d9578063636092e814610906578063642f6cef1461092b57806369cdbadb1461096f57600080fd5b80635147cd591161035357806357970e931161033857806357970e93146108775780635d4ee7f3146108a45780635f17e616146108b957600080fd5b80635147cd591461082557806351c98be31461085757600080fd5b80634585e33b1461077757806345d2ec171461079757806346982093146107c457806346e7a63e146107f857600080fd5b8063207b65161161041757806329e0a841116103e65780632b20e397116103cb5780632b20e397146106d8578063328ffd111461072a5780633ebe8d6c1461075757600080fd5b806329e0a8411461068b5780632a9032d3146106b857600080fd5b8063207b65161461060b57806320e3dbd41461062b5780632636aecf1461064b57806328c4b57b1461066b57600080fd5b806319d97a941161045357806319d97a941461054c5780631cdde251146105795780631e01043914610599578063206c32e8146105d657600080fd5b806306c1cc0014610490578063077ac621146104b25780630b7d33e6146104e557806312c550271461050557600080fd5b3661048b57005b600080fd5b34801561049c57600080fd5b506104b06104ab36600461418c565b61104c565b005b3480156104be57600080fd5b506104d26104cd36600461423f565b61129b565b6040519081526020015b60405180910390f35b3480156104f157600080fd5b506104b0610500366004614274565b6112d9565b34801561051157600080fd5b506105397f000000000000000000000000000000000000000000000000000000000000000081565b60405161ffff90911681526020016104dc565b34801561055857600080fd5b5061056c6105673660046142bb565b611367565b6040516104dc9190614342565b34801561058557600080fd5b506104b0610594366004614377565b611424565b3480156105a557600080fd5b506105b96105b43660046142bb565b611561565b6040516bffffffffffffffffffffffff90911681526020016104dc565b3480156105e257600080fd5b506105f66105f13660046143dc565b6115f6565b604080519283526020830191909152016104dc565b34801561061757600080fd5b5061056c6106263660046142bb565b611679565b34801561063757600080fd5b506104b0610646366004614408565b6116d1565b34801561065757600080fd5b506104b061066636600461446a565b61189b565b34801561067757600080fd5b506104d26106863660046144e4565b611b64565b34801561069757600080fd5b506106ab6106a63660046142bb565b611bcf565b6040516104dc9190614510565b3480156106c457600080fd5b506104b06106d3366004614651565b611cd4565b3480156106e457600080fd5b506011546107059073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016104dc565b34801561073657600080fd5b506104d26107453660046142bb565b60036020526000908152604090205481565b34801561076357600080fd5b506104d26107723660046142bb565b611db5565b34801561078357600080fd5b506104b06107923660046146d5565b611e1e565b3480156107a357600080fd5b506107b76107b23660046143dc565b61202d565b6040516104dc919061470b565b3480156107d057600080fd5b506104d27f000000000000000000000000000000000000000000000000000000000000000081565b34801561080457600080fd5b506104d26108133660046142bb565b600a6020526000908152604090205481565b34801561083157600080fd5b506108456108403660046142bb565b61209c565b60405160ff90911681526020016104dc565b34801561086357600080fd5b506104b061087236600461474f565b612130565b34801561088357600080fd5b506012546107059073ffffffffffffffffffffffffffffffffffffffff1681565b3480156108b057600080fd5b506104b06121d4565b3480156108c557600080fd5b506104b06108d43660046147a6565b61230f565b3480156108e557600080fd5b506104d26108f43660046142bb565b60076020526000908152604090205481565b34801561091257600080fd5b506015546105b9906bffffffffffffffffffffffff1681565b34801561093757600080fd5b5061095f7f000000000000000000000000000000000000000000000000000000000000000081565b60405190151581526020016104dc565b34801561097b57600080fd5b506104d261098a3660046142bb565b60086020526000908152604090205481565b3480156109a857600080fd5b506109bc6109b73660046146d5565b6123dc565b6040516104dc9291906147c8565b3480156109d657600080fd5b5061095f6109e53660046142bb565b600b6020526000908152604090205460ff1681565b348015610a0657600080fd5b506104d2610a153660046142bb565b6000908152600c602052604090205490565b348015610a3357600080fd5b506104d2610a423660046142bb565b60046020526000908152604090205481565b348015610a6057600080fd5b5061095f610a6f3660046142bb565b612509565b348015610a8057600080fd5b506104b061255b565b348015610a9557600080fd5b50610705610aa43660046142bb565b61265d565b348015610ab557600080fd5b506013546107059073ffffffffffffffffffffffffffffffffffffffff1681565b348015610ae257600080fd5b506104b0610af13660046147e3565b6126f1565b348015610b0257600080fd5b506104b0610b113660046147e3565b612782565b348015610b2257600080fd5b506107b7610b313660046147a6565b6127dc565b348015610b4257600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610705565b348015610b6d57600080fd5b506104d2610b7c3660046142bb565b60056020526000908152604090205481565b348015610b9a57600080fd5b506104b0610ba93660046147a6565b612899565b348015610bba57600080fd5b506104b0610bc936600461482f565b612ade565b348015610bda57600080fd5b506104b0610be936600461485f565b612bf6565b348015610bfa57600080fd5b50601554610845906c01000000000000000000000000900460ff1681565b348015610c2457600080fd5b506104b0610c333660046147a6565b60009182526009602052604090912055565b348015610c5157600080fd5b50610539610c603660046142bb565b600e6020526000908152604090205461ffff1681565b348015610c8257600080fd5b506107b7610c913660046142bb565b612e00565b348015610ca257600080fd5b5061056c610cb13660046142bb565b612e62565b348015610cc257600080fd5b506104d27f000000000000000000000000000000000000000000000000000000000000000081565b348015610cf657600080fd5b506104b0610d053660046144e4565b612f0e565b348015610d1657600080fd5b506104b0610d2536600461487c565b612f77565b348015610d3657600080fd5b506104b0610d453660046142bb565b613022565b348015610d5657600080fd5b5061056c6040518060400160405280600981526020017f666565644964486578000000000000000000000000000000000000000000000081525081565b348015610d9f57600080fd5b506105b9610dae3660046142bb565b6130a8565b348015610dbf57600080fd5b506104b0610dce366004614651565b613100565b348015610ddf57600080fd5b506107b7610dee3660046147a6565b61319a565b348015610dff57600080fd5b5061056c6040518060400160405280600b81526020017f626c6f636b4e756d62657200000000000000000000000000000000000000000081525081565b348015610e4857600080fd5b506104b0610e573660046148a1565b613297565b348015610e6857600080fd5b506105f6610e773660046147a6565b613316565b348015610e8857600080fd5b506104b0610e973660046148c6565b61337f565b348015610ea857600080fd5b506104b0610eb7366004614651565b6136e6565b348015610ec857600080fd5b506104d2610ed73660046147a6565b6137b1565b348015610ee857600080fd5b506104d2610ef73660046142bb565b60096020526000908152604090205481565b348015610f1557600080fd5b506104d260145481565b348015610f2b57600080fd5b506104b0610f3a366004614408565b6137e2565b348015610f4b57600080fd5b5061056c610f5a36600461492e565b6040805160c0808201835273ffffffffffffffffffffffffffffffffffffffff9890981680825260ff97881660208084019182528385019889526060808501988952608080860198895260a095860197885286519283019490945291519099168985015296519688019690965293519486019490945290519184019190915251828401528051808303909301835260e0909101905290565b348015610ffe57600080fd5b506104d261100d3660046142bb565b60066020526000908152604090205481565b34801561102b57600080fd5b506104d261103a3660046142bb565b60026020526000908152604090205481565b6040805161018081018252600461014082019081527f746573740000000000000000000000000000000000000000000000000000000061016083015281528151602081810184526000808352818401929092523083850181905263ffffffff8b166060850152608084015260ff808a1660a08501528451808301865283815260c085015260e0840189905284519182019094529081526101008201526bffffffffffffffffffffffff8516610120820152601254601154919273ffffffffffffffffffffffffffffffffffffffff9182169263095ea7b3921690611132908c16886149b6565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b16815273ffffffffffffffffffffffffffffffffffffffff90921660048301526bffffffffffffffffffffffff1660248201526044016020604051808303816000875af11580156111b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111d491906149fa565b5060008860ff1667ffffffffffffffff8111156111f3576111f361402e565b60405190808252806020026020018201604052801561121c578160200160208202803683370190505b50905060005b8960ff168160ff16101561128f57600061123b846137f6565b905080838360ff168151811061125357611253614a15565b6020908102919091018101919091526000918252600881526040808320889055600790915290208490558061128781614a44565b915050611222565b50505050505050505050565b600d60205282600052604060002060205281600052604060002081815481106112c357600080fd5b9060005260206000200160009250925050505481565b6013546040517f0b7d33e600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690630b7d33e6906113319085908590600401614a63565b600060405180830381600087803b15801561134b57600080fd5b505af115801561135f573d6000803e3d6000fd5b505050505050565b6013546040517f19d97a940000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff16906319d97a94906024015b600060405180830381865afa1580156113d8573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261141e9190810190614ac9565b92915050565b6013546040517ffa333dfb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff888116600483015260ff8816602483015260448201879052606482018690526084820185905260a4820184905290911690634ee88d35908990309063fa333dfb9060c401600060405180830381865afa1580156114c3573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526115099190810190614ac9565b6040518363ffffffff1660e01b8152600401611526929190614a63565b600060405180830381600087803b15801561154057600080fd5b505af1158015611554573d6000803e3d6000fd5b5050505050505050505050565b6013546040517f1e0104390000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690631e010439906024015b602060405180830381865afa1580156115d2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061141e9190614b09565b6000828152600d6020908152604080832061ffff85168452825280832080548251818502810185019093528083528493849392919083018282801561165a57602002820191906000526020600020905b815481526020019060010190808311611646575b5050505050905061166c8182516138c4565b92509250505b9250929050565b6013546040517f207b65160000000000000000000000000000000000000000000000000000000081526004810183905260609173ffffffffffffffffffffffffffffffffffffffff169063207b6516906024016113bb565b601180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8316908117909155604080517fc3f909d400000000000000000000000000000000000000000000000000000000815281516000939263c3f909d492600480820193918290030181865afa158015611767573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061178b9190614b31565b50601380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691909117909155601154604080517f1b6b6d230000000000000000000000000000000000000000000000000000000081529051939450911691631b6b6d23916004808201926020929091908290030181865afa15801561182e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118529190614b5f565b601280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff929092169190911790555050565b8560005b81811015611b595760008989838181106118bb576118bb614a15565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016118f491815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401611920929190614a63565b600060405180830381600087803b15801561193a57600080fd5b505af115801561194e573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa1580156119c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119e89190614b7c565b90508060ff16600103611b44576040517ffa333dfb000000000000000000000000000000000000000000000000000000008152306004820181905260ff8b166024830152604482018a9052606482018890526084820188905260a4820187905260009163fa333dfb9060c401600060405180830381865afa158015611a71573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611ab79190810190614ac9565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d3590611b109086908590600401614a63565b600060405180830381600087803b158015611b2a57600080fd5b505af1158015611b3e573d6000803e3d6000fd5b50505050505b50508080611b5190614b99565b91505061189f565b505050505050505050565b6000838152600c602090815260408083208054825181850281018501909352808352611bc593830182828015611bb957602002820191906000526020600020905b815481526020019060010190808311611ba5575b50505050508484613949565b90505b9392505050565b604080516101408101825260008082526020820181905260609282018390528282018190526080820181905260a0820181905260c0820181905260e082018190526101008201526101208101919091526013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff9091169063c7c3a19a90602401600060405180830381865afa158015611c8e573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261141e9190810190614bf4565b8060005b818160ff161015611daf5760135473ffffffffffffffffffffffffffffffffffffffff1663c8048022858560ff8516818110611d1657611d16614a15565b905060200201356040518263ffffffff1660e01b8152600401611d3b91815260200190565b600060405180830381600087803b158015611d5557600080fd5b505af1158015611d69573d6000803e3d6000fd5b50505050611d9c84848360ff16818110611d8557611d85614a15565b90506020020135600f613aa890919063ffffffff16565b5080611da781614a44565b915050611cd8565b50505050565b6000818152600e602052604081205461ffff1681805b8261ffff168161ffff1611611e16576000858152600d6020908152604080832061ffff85168452909152902054611e029083614d13565b915080611e0e81614d26565b915050611dcb565b509392505050565b60005a90506000611e3183850185614274565b5060008181526005602090815260408083205460049092528220549293509190611e59613ab4565b905082600003611e79576000848152600560205260409020819055611fd4565b600084815260036020526040812054611e928484614d47565b611e9c9190614d47565b6000868152600e6020908152604080832054600d835281842061ffff909116808552908352818420805483518186028101860190945280845295965090949192909190830182828015611f0e57602002820191906000526020600020905b815481526020019060010190808311611efa575b505050505090507f000000000000000000000000000000000000000000000000000000000000000061ffff16815103611f895781611f4b81614d26565b6000898152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff83161790559250505b506000868152600d6020908152604080832061ffff909416835292815282822080546001818101835591845282842001859055888352600c8252928220805493840181558252902001555b600084815260066020526040812054611fee906001614d13565b60008681526006602090815260408083208490556004909152902083905590506120188583612899565b612023858784612f0e565b5050505050505050565b6000828152600d6020908152604080832061ffff8516845282529182902080548351818402810184019094528084526060939283018282801561208f57602002820191906000526020600020905b81548152602001906001019080831161207b575b5050505050905092915050565b6013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff1690635147cd5990602401602060405180830381865afa15801561210c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061141e9190614b7c565b8160005b818110156121cd5730635f17e61686868481811061215457612154614a15565b90506020020135856040518363ffffffff1660e01b815260040161218892919091825263ffffffff16602082015260400190565b600060405180830381600087803b1580156121a257600080fd5b505af11580156121b6573d6000803e3d6000fd5b5050505080806121c590614b99565b915050612134565b5050505050565b6121dc613b56565b6012546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009173ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa15801561224b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061226f9190614d5a565b6012546040517fa9059cbb0000000000000000000000000000000000000000000000000000000081523360048201526024810183905291925073ffffffffffffffffffffffffffffffffffffffff169063a9059cbb906044016020604051808303816000875af11580156122e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061230b91906149fa565b5050565b60008281526003602090815260408083208490556005825280832083905560068252808320839055600c909152812061234791613fdb565b6000828152600e602052604081205461ffff16905b8161ffff168161ffff16116123a3576000848152600d6020908152604080832061ffff85168452909152812061239191613fdb565b8061239b81614d26565b91505061235c565b5050506000908152600e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000169055565b6000606060005a905060006123f3858701876142bb565b60008181526009602090815260408083205460089092528220549293509190838367ffffffffffffffff81111561242c5761242c61402e565b6040519080825280601f01601f191660200182016040528015612456576020820181803683370190505b50604051602001612468929190614a63565b60405160208183030381529060405290506000612483613ab4565b9050600061249086612509565b90505b835a61249f9089614d47565b6124ab90612710614d13565b10156124f95781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055816124f181614d73565b925050612493565b9a91995090975050505050505050565b600081815260056020526040812054810361252657506001919050565b600082815260036020908152604080832054600490925290912054612549613ab4565b6125539190614d47565b101592915050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146125e1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6013546040517f79ea99430000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff16906379ea994390602401602060405180830381865afa1580156126cd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061141e9190614b5f565b6013546040517fcd7f71b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063cd7f71b59061274b90869086908690600401614da8565b600060405180830381600087803b15801561276557600080fd5b505af1158015612779573d6000803e3d6000fd5b50505050505050565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911690634ee88d359061274b90869086908690600401614da8565b6013546040517f06e3b632000000000000000000000000000000000000000000000000000000008152600481018490526024810183905260609173ffffffffffffffffffffffffffffffffffffffff16906306e3b63290604401600060405180830381865afa158015612853573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611bc89190810190614dfc565b6014546000838152600260205260409020546128b59083614d47565b111561230b576013546040517fc7c3a19a0000000000000000000000000000000000000000000000000000000081526004810184905260009173ffffffffffffffffffffffffffffffffffffffff169063c7c3a19a90602401600060405180830381865afa15801561292b573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526129719190810190614bf4565b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810186905291925060009173ffffffffffffffffffffffffffffffffffffffff9091169063b657bc9c90602401602060405180830381865afa1580156129e6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a0a9190614b09565b601554909150612a2e9082906c01000000000000000000000000900460ff166149b6565b6bffffffffffffffffffffffff1682606001516bffffffffffffffffffffffff161015611daf57601554612a719085906bffffffffffffffffffffffff16612ade565b60008481526002602090815260409182902085905560155482518781526bffffffffffffffffffffffff909116918101919091529081018490527f49d4100ab0124eb4a9a65dc4ea08d6412a43f6f05c49194983f5b322bcc0a5c09060600160405180910390a150505050565b6012546013546040517f095ea7b300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff91821660048201526bffffffffffffffffffffffff8416602482015291169063095ea7b3906044016020604051808303816000875af1158015612b66573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b8a91906149fa565b506013546040517f948108f7000000000000000000000000000000000000000000000000000000008152600481018490526bffffffffffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063948108f790604401611331565b6040517fc04198220000000000000000000000000000000000000000000000000000000081526000600482018190526024820181905290309063c041982290604401600060405180830381865afa158015612c55573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612c9b9190810190614dfc565b80519091506000612caa613ab4565b905060005b828110156121cd576000848281518110612ccb57612ccb614a15565b60209081029190910101516013546040517f5147cd590000000000000000000000000000000000000000000000000000000081526004810183905291925060009173ffffffffffffffffffffffffffffffffffffffff90911690635147cd5990602401602060405180830381865afa158015612d4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d6f9190614b7c565b90508060ff16600103612deb578660ff16600003612dbb576040513090859084907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4612deb565b6040513090859084907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a45b50508080612df890614b99565b915050612caf565b6000818152600c6020908152604091829020805483518184028101840190945280845260609392830182828015612e5657602002820191906000526020600020905b815481526020019060010190808311612e42575b50505050509050919050565b60168181548110612e7257600080fd5b906000526020600020016000915090508054612e8d90614ea2565b80601f0160208091040260200160405190810160405280929190818152602001828054612eb990614ea2565b8015612f065780601f10612edb57610100808354040283529160200191612f06565b820191906000526020600020905b815481529060010190602001808311612ee957829003601f168201915b505050505081565b6000838152600760205260409020545b805a612f2a9085614d47565b612f3690612710614d13565b1015611daf5781406000908152600b6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055612f1e565b6013546040517fa72aa27e0000000000000000000000000000000000000000000000000000000081526004810184905263ffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff9091169063a72aa27e90604401600060405180830381600087803b158015612fef57600080fd5b505af1158015613003573d6000803e3d6000fd5b505050600092835250600a602052604090912063ffffffff9091169055565b6013546040517f744bfe610000000000000000000000000000000000000000000000000000000081526004810183905230602482015273ffffffffffffffffffffffffffffffffffffffff9091169063744bfe6190604401600060405180830381600087803b15801561309457600080fd5b505af11580156121cd573d6000803e3d6000fd5b6013546040517fb657bc9c0000000000000000000000000000000000000000000000000000000081526004810183905260009173ffffffffffffffffffffffffffffffffffffffff169063b657bc9c906024016115b5565b8060005b818163ffffffff161015611daf573063af953a4a858563ffffffff851681811061313057613130614a15565b905060200201356040518263ffffffff1660e01b815260040161315591815260200190565b600060405180830381600087803b15801561316f57600080fd5b505af1158015613183573d6000803e3d6000fd5b50505050808061319290614ef5565b915050613104565b606060006131a8600f613bd9565b90508084106131e3576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826000036131f8576131f58482614d47565b92505b60008367ffffffffffffffff8111156132135761321361402e565b60405190808252806020026020018201604052801561323c578160200160208202803683370190505b50905060005b8481101561328e5761325f6132578288614d13565b600f90613be3565b82828151811061327157613271614a15565b60209081029190910101528061328681614b99565b915050613242565b50949350505050565b60006132a1613ab4565b90508160ff166000036132e2576040513090829085907f97009585a4d2440f981ab6f6eec514343e1e6b2aa9b991a26998e6806f41bf0890600090a4505050565b6040513090829085907fc76416badc8398ce17c93eab7b4f60f263241694cf503e4df24f233a8cc1c50d90600090a4505050565b6000828152600c6020908152604080832080548251818502810185019093528083528493849392919083018282801561336e57602002820191906000526020600020905b81548152602001906001019080831161335a575b5050505050905061166c81856138c4565b8260005b8181101561135f57600086868381811061339f5761339f614a15565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc82836040516020016133d891815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401613404929190614a63565b600060405180830381600087803b15801561341e57600080fd5b505af1158015613432573d6000803e3d6000fd5b50506013546040517f5147cd59000000000000000000000000000000000000000000000000000000008152600481018590526000935073ffffffffffffffffffffffffffffffffffffffff9091169150635147cd5990602401602060405180830381865afa1580156134a8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134cc9190614b7c565b90508060ff166001036136d1577f000000000000000000000000000000000000000000000000000000000000000060ff87161561352657507f00000000000000000000000000000000000000000000000000000000000000005b60003073ffffffffffffffffffffffffffffffffffffffff1663fa333dfb3089858860405160200161355a91815260200190565b60405160208183030381529060405261357290614f0e565b60405160e086901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909416600485015260ff90921660248401526044830152606482015260006084820181905260a482015260c401600060405180830381865afa1580156135fd573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526136439190810190614ac9565b6013546040517f4ee88d3500000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff1690634ee88d359061369c9087908590600401614a63565b600060405180830381600087803b1580156136b657600080fd5b505af11580156136ca573d6000803e3d6000fd5b5050505050505b505080806136de90614b99565b915050613383565b8060005b81811015611daf57600084848381811061370657613706614a15565b9050602002013590503073ffffffffffffffffffffffffffffffffffffffff16637e7a46dc828360405160200161373f91815260200190565b6040516020818303038152906040526040518363ffffffff1660e01b815260040161376b929190614a63565b600060405180830381600087803b15801561378557600080fd5b505af1158015613799573d6000803e3d6000fd5b505050505080806137a990614b99565b9150506136ea565b600c60205281600052604060002081815481106137cd57600080fd5b90600052602060002001600091509150505481565b6137ea613b56565b6137f381613bef565b50565b6011546040517f3f678e11000000000000000000000000000000000000000000000000000000008152600091829173ffffffffffffffffffffffffffffffffffffffff90911690633f678e1190613851908690600401614f50565b6020604051808303816000875af1158015613870573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138949190614d5a565b90506138a1600f82613ce4565b506060909201516000838152600a6020526040902063ffffffff90911690555090565b8151600090819081908415806138da5750808510155b156138e3578094505b60008092505b8583101561393f578660016138fe8585614d47565b6139089190614d47565b8151811061391857613918614a15565b60200260200101518161392b9190614d13565b90508261393781614b99565b9350506138e9565b9694955050505050565b8251600090819083158061395d5750808410155b15613966578093505b60008467ffffffffffffffff8111156139815761398161402e565b6040519080825280602002602001820160405280156139aa578160200160208202803683370190505b509050600092505b84831015613a18578660016139c78585614d47565b6139d19190614d47565b815181106139e1576139e1614a15565b60200260200101518184815181106139fb576139fb614a15565b602090810291909101015282613a1081614b99565b9350506139b2565b613a3181600060018451613a2c9190614d47565b613cf0565b85606403613a6a578060018251613a489190614d47565b81518110613a5857613a58614a15565b60200260200101519350505050611bc8565b806064825188613a7a91906150a2565b613a84919061510e565b81518110613a9457613a94614a15565b602002602001015193505050509392505050565b6000611bc88383613e68565b60007f000000000000000000000000000000000000000000000000000000000000000015613b5157606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613b28573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b4c9190614d5a565b905090565b504390565b60005473ffffffffffffffffffffffffffffffffffffffff163314613bd7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016125d8565b565b600061141e825490565b6000611bc88383613f62565b3373ffffffffffffffffffffffffffffffffffffffff821603613c6e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016125d8565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000611bc88383613f8c565b8181808203613d00575050505050565b6000856002613d0f8787615122565b613d199190615142565b613d2390876151aa565b81518110613d3357613d33614a15565b602002602001015190505b818313613e42575b80868481518110613d5957613d59614a15565b60200260200101511015613d795782613d71816151d2565b935050613d46565b858281518110613d8b57613d8b614a15565b6020026020010151811015613dac5781613da481615203565b925050613d79565b818313613e3d57858281518110613dc557613dc5614a15565b6020026020010151868481518110613ddf57613ddf614a15565b6020026020010151878581518110613df957613df9614a15565b60200260200101888581518110613e1257613e12614a15565b60209081029190910101919091525282613e2b816151d2565b9350508180613e3990615203565b9250505b613d3e565b81851215613e5557613e55868684613cf0565b8383121561135f5761135f868486613cf0565b60008181526001830160205260408120548015613f51576000613e8c600183614d47565b8554909150600090613ea090600190614d47565b9050818114613f05576000866000018281548110613ec057613ec0614a15565b9060005260206000200154905080876000018481548110613ee357613ee3614a15565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613f1657613f16615234565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061141e565b600091505061141e565b5092915050565b6000826000018281548110613f7957613f79614a15565b9060005260206000200154905092915050565b6000818152600183016020526040812054613fd35750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561141e565b50600061141e565b50805460008255906000526020600020908101906137f391905b808211156140095760008155600101613ff5565b5090565b60ff811681146137f357600080fd5b63ffffffff811681146137f357600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610140810167ffffffffffffffff811182821017156140815761408161402e565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156140ce576140ce61402e565b604052919050565b600067ffffffffffffffff8211156140f0576140f061402e565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261412d57600080fd5b813561414061413b826140d6565b614087565b81815284602083860101111561415557600080fd5b816020850160208301376000918101602001919091529392505050565b6bffffffffffffffffffffffff811681146137f357600080fd5b600080600080600080600060e0888a0312156141a757600080fd5b87356141b28161400d565b965060208801356141c28161401c565b955060408801356141d28161400d565b9450606088013567ffffffffffffffff8111156141ee57600080fd5b6141fa8a828b0161411c565b945050608088013561420b81614172565b9699959850939692959460a0840135945060c09093013592915050565b803561ffff8116811461423a57600080fd5b919050565b60008060006060848603121561425457600080fd5b8335925061426460208501614228565b9150604084013590509250925092565b6000806040838503121561428757600080fd5b82359150602083013567ffffffffffffffff8111156142a557600080fd5b6142b18582860161411c565b9150509250929050565b6000602082840312156142cd57600080fd5b5035919050565b60005b838110156142ef5781810151838201526020016142d7565b50506000910152565b600081518084526143108160208601602086016142d4565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611bc860208301846142f8565b73ffffffffffffffffffffffffffffffffffffffff811681146137f357600080fd5b600080600080600080600060e0888a03121561439257600080fd5b8735965060208801356143a481614355565b955060408801356143b48161400d565b969995985095966060810135965060808101359560a0820135955060c0909101359350915050565b600080604083850312156143ef57600080fd5b823591506143ff60208401614228565b90509250929050565b60006020828403121561441a57600080fd5b8135611bc881614355565b60008083601f84011261443757600080fd5b50813567ffffffffffffffff81111561444f57600080fd5b6020830191508360208260051b850101111561167257600080fd5b600080600080600080600060c0888a03121561448557600080fd5b873567ffffffffffffffff81111561449c57600080fd5b6144a88a828b01614425565b90985096505060208801356144bc8161400d565b96999598509596604081013596506060810135956080820135955060a0909101359350915050565b6000806000606084860312156144f957600080fd5b505081359360208301359350604090920135919050565b6020815261453760208201835173ffffffffffffffffffffffffffffffffffffffff169052565b60006020830151614550604084018263ffffffff169052565b50604083015161014080606085015261456d6101608501836142f8565b9150606085015161458e60808601826bffffffffffffffffffffffff169052565b50608085015173ffffffffffffffffffffffffffffffffffffffff811660a08601525060a085015167ffffffffffffffff811660c08601525060c085015163ffffffff811660e08601525060e08501516101006145fa818701836bffffffffffffffffffffffff169052565b860151905061012061460f8682018315159052565b8601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00183870152905061464783826142f8565b9695505050505050565b6000806020838503121561466457600080fd5b823567ffffffffffffffff81111561467b57600080fd5b61468785828601614425565b90969095509350505050565b60008083601f8401126146a557600080fd5b50813567ffffffffffffffff8111156146bd57600080fd5b60208301915083602082850101111561167257600080fd5b600080602083850312156146e857600080fd5b823567ffffffffffffffff8111156146ff57600080fd5b61468785828601614693565b6020808252825182820181905260009190848201906040850190845b8181101561474357835183529284019291840191600101614727565b50909695505050505050565b60008060006040848603121561476457600080fd5b833567ffffffffffffffff81111561477b57600080fd5b61478786828701614425565b909450925050602084013561479b8161401c565b809150509250925092565b600080604083850312156147b957600080fd5b50508035926020909101359150565b8215158152604060208201526000611bc560408301846142f8565b6000806000604084860312156147f857600080fd5b83359250602084013567ffffffffffffffff81111561481657600080fd5b61482286828701614693565b9497909650939450505050565b6000806040838503121561484257600080fd5b82359150602083013561485481614172565b809150509250929050565b60006020828403121561487157600080fd5b8135611bc88161400d565b6000806040838503121561488f57600080fd5b8235915060208301356148548161401c565b600080604083850312156148b457600080fd5b8235915060208301356148548161400d565b600080600080606085870312156148dc57600080fd5b843567ffffffffffffffff8111156148f357600080fd5b6148ff87828801614425565b90955093505060208501356149138161400d565b915060408501356149238161400d565b939692955090935050565b60008060008060008060c0878903121561494757600080fd5b863561495281614355565b955060208701356149628161400d565b95989597505050506040840135936060810135936080820135935060a0909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006bffffffffffffffffffffffff808316818516818304811182151516156149e1576149e1614987565b02949350505050565b8051801515811461423a57600080fd5b600060208284031215614a0c57600080fd5b611bc8826149ea565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff8103614a5a57614a5a614987565b60010192915050565b828152604060208201526000611bc560408301846142f8565b600082601f830112614a8d57600080fd5b8151614a9b61413b826140d6565b818152846020838601011115614ab057600080fd5b614ac18260208301602087016142d4565b949350505050565b600060208284031215614adb57600080fd5b815167ffffffffffffffff811115614af257600080fd5b614ac184828501614a7c565b805161423a81614172565b600060208284031215614b1b57600080fd5b8151611bc881614172565b805161423a81614355565b60008060408385031215614b4457600080fd5b8251614b4f81614355565b6020939093015192949293505050565b600060208284031215614b7157600080fd5b8151611bc881614355565b600060208284031215614b8e57600080fd5b8151611bc88161400d565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614bca57614bca614987565b5060010190565b805161423a8161401c565b805167ffffffffffffffff8116811461423a57600080fd5b600060208284031215614c0657600080fd5b815167ffffffffffffffff80821115614c1e57600080fd5b908301906101408286031215614c3357600080fd5b614c3b61405d565b614c4483614b26565b8152614c5260208401614bd1565b6020820152604083015182811115614c6957600080fd5b614c7587828601614a7c565b604083015250614c8760608401614afe565b6060820152614c9860808401614b26565b6080820152614ca960a08401614bdc565b60a0820152614cba60c08401614bd1565b60c0820152614ccb60e08401614afe565b60e0820152610100614cde8185016149ea565b908201526101208381015183811115614cf657600080fd5b614d0288828701614a7c565b918301919091525095945050505050565b8082018082111561141e5761141e614987565b600061ffff808316818103614d3d57614d3d614987565b6001019392505050565b8181038181111561141e5761141e614987565b600060208284031215614d6c57600080fd5b5051919050565b600081614d8257614d82614987565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b83815260406020820152816040820152818360608301376000818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b60006020808385031215614e0f57600080fd5b825167ffffffffffffffff80821115614e2757600080fd5b818501915085601f830112614e3b57600080fd5b815181811115614e4d57614e4d61402e565b8060051b9150614e5e848301614087565b8181529183018401918481019088841115614e7857600080fd5b938501935b83851015614e9657845182529385019390850190614e7d565b98975050505050505050565b600181811c90821680614eb657607f821691505b602082108103614eef577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600063ffffffff808316818103614d3d57614d3d614987565b80516020808301519190811015614eef577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b6020815260008251610140806020850152614f6f6101608501836142f8565b915060208501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe080868503016040870152614fab84836142f8565b935060408701519150614fd6606087018373ffffffffffffffffffffffffffffffffffffffff169052565b606087015163ffffffff811660808801529150608087015173ffffffffffffffffffffffffffffffffffffffff811660a0880152915060a087015160ff811660c0880152915060c08701519150808685030160e087015261503784836142f8565b935060e0870151915061010081878603018188015261505685846142f8565b94508088015192505061012081878603018188015261507585846142f8565b94508088015192505050615098828601826bffffffffffffffffffffffff169052565b5090949350505050565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156150da576150da614987565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261511d5761511d6150df565b500490565b8181036000831280158383131683831282161715613f5b57613f5b614987565b600082615151576151516150df565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f8000000000000000000000000000000000000000000000000000000000000000831416156151a5576151a5614987565b500590565b80820182811260008312801582168215821617156151ca576151ca614987565b505092915050565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203614bca57614bca614987565b60007f80000000000000000000000000000000000000000000000000000000000000008203614d8257614d82614987565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c6343000810000a307834353534343832643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030",
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

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactor) SetInterval(opts *bind.TransactOpts, upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.contract.Transact(opts, "setInterval", upkeepId, _interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetInterval(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, _interval)
}

func (_VerifiableLoadUpkeep *VerifiableLoadUpkeepTransactorSession) SetInterval(upkeepId *big.Int, _interval *big.Int) (*types.Transaction, error) {
	return _VerifiableLoadUpkeep.Contract.SetInterval(&_VerifiableLoadUpkeep.TransactOpts, upkeepId, _interval)
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
