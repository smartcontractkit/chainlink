package handler

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	evm "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

const (
	ConditionTrigger uint8 = iota
	LogTrigger

	blockNumber            = "blockNumber"
	expectedTypeAndVersion = "KeeperRegistry 2.1.0"
	feedIdHex              = "feedIdHex"
	feedIDs                = "feedIDs"
	timestamp              = "timestamp"
)

var packer = encoding.NewAbiPacker()

var links []string

func (k *Keeper) Debug(ctx context.Context, args []string) {
	if len(args) < 1 {
		failCheckArgs("no upkeepID supplied", nil)
	}
	// test that we are connected to an archive node
	_, err := k.client.BalanceAt(ctx, gethcommon.Address{}, big.NewInt(1))
	if err != nil {
		failCheckConfig("you are not connected to an archive node; try using infura or alchemy", err)
	}
	chainIDBig, err := k.client.ChainID(ctx)
	if err != nil {
		failUnknown("unable to retrieve chainID from rpc client", err)
	}
	chainID := chainIDBig.Int64()
	// connect to registry contract
	latestCallOpts := &bind.CallOpts{Context: ctx}  // always use latest block
	triggerCallOpts := &bind.CallOpts{Context: ctx} // use latest block for conditionals, but use block from tx for log triggers
	registryAddress := gethcommon.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry21, err := iregistry21.NewIKeeperRegistryMaster(registryAddress, k.client)
	if err != nil {
		failUnknown("failed to connect to registry contract", err)
	}
	// verify contract is correct
	typeAndVersion, err := keeperRegistry21.TypeAndVersion(latestCallOpts)
	if err != nil {
		failCheckConfig("failed to get typeAndVersion: are you sure you have the correct contract address?", err)
	}
	if typeAndVersion != expectedTypeAndVersion {
		failCheckConfig(fmt.Sprintf("invalid registry contract: this command can only debug %s, got: %s", expectedTypeAndVersion, typeAndVersion), nil)
	}
	// get upkeepID from command args
	upkeepID := big.NewInt(0)
	upkeepIDNoPrefix := utils.RemoveHexPrefix(args[0])
	_, wasBase10 := upkeepID.SetString(upkeepIDNoPrefix, 10)
	if !wasBase10 {
		_, wasBase16 := upkeepID.SetString(upkeepIDNoPrefix, 16)
		if !wasBase16 {
			failCheckArgs("invalid upkeep ID", nil)
		}
	}
	// get upkeep info
	triggerType, err := keeperRegistry21.GetTriggerType(latestCallOpts, upkeepID)
	if err != nil {
		failUnknown("failed to get trigger type: ", err)
	}
	upkeepInfo, err := keeperRegistry21.GetUpkeep(latestCallOpts, upkeepID)
	if err != nil {
		failUnknown("failed to get trigger type: ", err)
	}
	minBalance, err := keeperRegistry21.GetMinBalance(latestCallOpts, upkeepID)
	if err != nil {
		failUnknown("failed to get min balance: ", err)
	}
	// do basic sanity checks
	if (upkeepInfo.Target == gethcommon.Address{}) {
		failCheckArgs("this upkeep does not exist on this registry", nil)
	}
	addLink("upkeep", common.UpkeepLink(chainID, upkeepID))
	addLink("target", common.ContractExplorerLink(chainID, upkeepInfo.Target))
	if upkeepInfo.Paused {
		resolveIneligible("upkeep is paused")
	}
	if upkeepInfo.MaxValidBlocknumber != math.MaxUint32 {
		resolveIneligible("upkeep is cancelled")
	}
	message("upkeep is active (not paused or cancelled)")
	if upkeepInfo.Balance.Cmp(minBalance) == -1 {
		resolveIneligible("minBalance is < upkeep balance")
	}
	message("upkeep is funded above the min balance")
	if bigmath.Div(bigmath.Mul(bigmath.Sub(upkeepInfo.Balance, minBalance), big.NewInt(100)), minBalance).Cmp(big.NewInt(5)) == -1 {
		warning("upkeep balance is < 5% larger than minBalance")
	}
	// local state for pipeline results
	var checkResult iregistry21.CheckUpkeep
	var blockNum uint64
	var performData []byte
	upkeepNeeded := false
	// check upkeep
	if triggerType == ConditionTrigger {
		message("upkeep identified as conditional trigger")
		tmpCheckResult, err := keeperRegistry21.CheckUpkeep0(latestCallOpts, upkeepID)
		if err != nil {
			failUnknown("failed to check upkeep: ", err)
		}
		checkResult = iregistry21.CheckUpkeep(tmpCheckResult)
		// do tenderly simulation
		rawCall, err := core.RegistryABI.Pack("checkUpkeep", upkeepID, []byte{})
		if err != nil {
			failUnknown("failed to pack raw checkUpkeep call", err)
		}
		addLink("checkUpkeep simulation", tenderlySimLink(k.cfg, chainID, 0, rawCall, registryAddress))
	} else if triggerType == LogTrigger {
		// validate inputs
		message("upkeep identified as log trigger")
		if len(args) != 3 {
			failCheckArgs("txHash and log index must be supplied to command in order to debug log triggered upkeeps", nil)
		}
		txHash := gethcommon.HexToHash(args[1])
		logIndex, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			failCheckArgs("unable to parse log index", err)
		}
		// find transaciton receipt
		_, isPending, err := k.client.TransactionByHash(ctx, txHash)
		if err != nil {
			log.Fatal("failed to fetch tx receipt", err)
		}
		if isPending {
			resolveIneligible(fmt.Sprintf("tx %s is still pending confirmation", txHash))
		}
		receipt, err := k.client.TransactionReceipt(ctx, txHash)
		if err != nil {
			failCheckArgs("failed to fetch tx receipt", err)
		}
		addLink("trigger", common.ExplorerLink(chainID, txHash))
		blockNum = receipt.BlockNumber.Uint64()
		// find matching log event in tx
		var triggeringEvent *types.Log
		for i, log := range receipt.Logs {
			if log.Index == uint(logIndex) {
				triggeringEvent = receipt.Logs[i]
			}
		}
		if triggeringEvent == nil {
			failCheckArgs(fmt.Sprintf("unable to find log with index %d in transaction", logIndex), nil)
		}
		// check that tx for this upkeep / tx was not already performed
		message(fmt.Sprintf("LogTrigger{blockNum: %d, blockHash: %s, txHash: %s, logIndex: %d}", blockNum, receipt.BlockHash.Hex(), txHash, logIndex))
		workID := mustUpkeepWorkID(upkeepID, blockNum, receipt.BlockHash, txHash, logIndex)
		message(fmt.Sprintf("workID computed: %s", hex.EncodeToString(workID[:])))
		hasKey, err := keeperRegistry21.HasDedupKey(latestCallOpts, workID)
		if err != nil {
			failUnknown("failed to check if upkeep was already performed: ", err)
		}
		if hasKey {
			resolveIneligible("upkeep was already performed")
		}
		triggerCallOpts = &bind.CallOpts{Context: ctx, BlockNumber: big.NewInt(receipt.BlockNumber.Int64())}
		rawTriggerConfig, err := keeperRegistry21.GetUpkeepTriggerConfig(triggerCallOpts, upkeepID)
		if err != nil {
			failUnknown("failed to fetch trigger config for upkeep", err)
		}
		triggerConfig, err := packer.UnpackLogTriggerConfig(rawTriggerConfig)
		if err != nil {
			failUnknown("failed to unpack trigger config", err)
		}
		if triggerConfig.FilterSelector > 7 {
			resolveIneligible(fmt.Sprintf("invalid filter selector %d", triggerConfig.FilterSelector))
		}
		if !logMatchesTriggerConfig(triggeringEvent, triggerConfig) {
			resolveIneligible("log does not match trigger config")
		}
		header, err := k.client.HeaderByHash(ctx, receipt.BlockHash)
		if err != nil {
			failUnknown("failed to find block", err)
		}
		triggerData, err := packTriggerData(triggeringEvent, header.Time)
		if err != nil {
			failUnknown("failed to pack trigger data", err)
		}
		checkResult, err = keeperRegistry21.CheckUpkeep(triggerCallOpts, upkeepID, triggerData)
		if err != nil {
			failUnknown("failed to check upkeep", err)
		}
		// do tenderly simulations
		rawCall, err := core.RegistryABI.Pack("checkUpkeep", upkeepID, triggerData)
		if err != nil {
			failUnknown("failed to pack raw checkUpkeep call", err)
		}
		addLink("checkUpkeep simulation", tenderlySimLink(k.cfg, chainID, blockNum, rawCall, registryAddress))
		rawCall = append(core.ILogAutomationABI.Methods["checkLog"].ID, triggerData...)
		addLink("checkLog (direct) simulation", tenderlySimLink(k.cfg, chainID, blockNum, rawCall, upkeepInfo.Target))
	} else {
		resolveIneligible(fmt.Sprintf("invalid trigger type: %d", triggerType))
	}
	upkeepNeeded, performData = checkResult.UpkeepNeeded, checkResult.PerformData
	// handle streams lookup
	if checkResult.UpkeepFailureReason != 0 {
		message(fmt.Sprintf("checkUpkeep failed with UpkeepFailureReason %d", checkResult.UpkeepFailureReason))
	}
	if checkResult.UpkeepFailureReason == uint8(encoding.UpkeepFailureReasonTargetCheckReverted) {
		streamsLookupErr, err := packer.DecodeStreamsLookupRequest(checkResult.PerformData)
		if err == nil {
			message("upkeep reverted with StreamsLookup")
			message(fmt.Sprintf("StreamsLookup data: {FeedParamKey: %s, Feeds: %v, TimeParamKey: %s, Time: %d, ExtraData: %s}", streamsLookupErr.FeedParamKey, streamsLookupErr.Feeds, streamsLookupErr.TimeParamKey, streamsLookupErr.Time.Uint64(), hexutil.Encode(streamsLookupErr.ExtraData)))
			if streamsLookupErr.FeedParamKey == feedIdHex && streamsLookupErr.TimeParamKey == blockNumber {
				message("using mercury lookup v0.2")
				// handle v0.2
				cfg, err := keeperRegistry21.GetUpkeepPrivilegeConfig(triggerCallOpts, upkeepID)
				if err != nil {
					failUnknown("failed to get upkeep privilege config ", err)
				}
				allowed := false
				if len(cfg) > 0 {
					var privilegeConfig evm.UpkeepPrivilegeConfig
					if err := json.Unmarshal(cfg, &privilegeConfig); err != nil {
						failUnknown("failed to unmarshal privilege config ", err)
					}
					allowed = privilegeConfig.MercuryEnabled
				}
				if !allowed {
					resolveIneligible("upkeep reverted with StreamsLookup but is not allowed to access streams")
				}
			} else if streamsLookupErr.FeedParamKey != feedIDs || streamsLookupErr.TimeParamKey != timestamp {
				// handle v0.3
				resolveIneligible("upkeep reverted with StreamsLookup but the configuration is invalid")
			} else {
				message("using mercury lookup v0.3")
			}
			streamsLookup := &StreamsLookup{streamsLookupErr.FeedParamKey, streamsLookupErr.Feeds, streamsLookupErr.TimeParamKey, streamsLookupErr.Time, streamsLookupErr.ExtraData, upkeepID, blockNum}
			handler := NewMercuryLookupHandler(&MercuryCredentials{k.cfg.MercuryLegacyURL, k.cfg.MercuryURL, k.cfg.MercuryID, k.cfg.MercuryKey}, k.rpcClient)
			state, failureReason, values, _, err := handler.doMercuryRequest(ctx, streamsLookup)
			if failureReason == UpkeepFailureReasonInvalidRevertDataInput {
				resolveIneligible("upkeep used invalid revert data")
			}
			if state == InvalidMercuryRequest {
				resolveIneligible("the mercury request data is invalid")
			}
			if err != nil {
				failCheckConfig("failed to do mercury request ", err)
			}
			callbackResult, err := keeperRegistry21.CheckCallback(triggerCallOpts, upkeepID, values, streamsLookup.extraData)
			if err != nil {
				failUnknown("failed to execute mercury callback ", err)
			}
			upkeepNeeded, performData = callbackResult.UpkeepNeeded, callbackResult.PerformData
			// do tenderly simulation
			rawCall, err := core.RegistryABI.Pack("checkCallback", upkeepID, values, streamsLookup.extraData)
			if err != nil {
				failUnknown("failed to pack raw checkUpkeep call", err)
			}
			addLink("checkCallback simulation", tenderlySimLink(k.cfg, chainID, blockNum, rawCall, registryAddress))
		} else {
			message("did not revert with StreamsLookup error")
		}
	}
	if !upkeepNeeded {
		resolveIneligible("upkeep is not needed")
	}
	// simulate perform ukeep
	simulateResult, err := keeperRegistry21.SimulatePerformUpkeep(latestCallOpts, upkeepID, performData)
	if err != nil {
		failUnknown("failed to simulate perform upkeep: ", err)
	}
	if simulateResult.Success {
		resolveEligible()
	}
}

func logMatchesTriggerConfig(log *types.Log, config automation_utils_2_1.LogTriggerConfig) bool {
	if log.Topics[0] != config.Topic0 {
		return false
	}
	if config.FilterSelector&1 > 0 && (len(log.Topics) < 1 || log.Topics[1] != config.Topic1) {
		return false
	}
	if config.FilterSelector&2 > 0 && (len(log.Topics) < 2 || log.Topics[2] != config.Topic2) {
		return false
	}
	if config.FilterSelector&4 > 0 && (len(log.Topics) < 3 || log.Topics[3] != config.Topic3) {
		return false
	}
	return true
}

func packTriggerData(log *types.Log, blockTime uint64) ([]byte, error) {
	var topics [][32]byte
	for _, topic := range log.Topics {
		topics = append(topics, topic)
	}
	b, err := core.UtilsABI.Methods["_log"].Inputs.Pack(&automation_utils_2_1.Log{
		Index:       big.NewInt(int64(log.Index)),
		Timestamp:   big.NewInt(int64(blockTime)),
		TxHash:      log.TxHash,
		BlockNumber: big.NewInt(int64(log.BlockNumber)),
		BlockHash:   log.BlockHash,
		Source:      log.Address,
		Topics:      topics,
		Data:        log.Data,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func mustUpkeepWorkID(upkeepID *big.Int, blockNum uint64, blockHash [32]byte, txHash [32]byte, logIndex int64) [32]byte {
	// TODO - this is a copy of the code in core.UpkeepWorkID
	// We should refactor that code to be more easily exported ex not rely on Trigger structs
	trigger := ocr2keepers.Trigger{
		LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
			TxHash:      txHash,
			Index:       uint32(logIndex),
			BlockNumber: ocr2keepers.BlockNumber(blockNum),
			BlockHash:   blockHash,
		},
	}
	upkeepIdentifier := &ocr2keepers.UpkeepIdentifier{}
	upkeepIdentifier.FromBigInt(upkeepID)
	workID := core.UpkeepWorkID(*upkeepIdentifier, trigger)
	workIDBytes, err := hex.DecodeString(workID)
	if err != nil {
		failUnknown("failed to decode workID", err)
	}
	var result [32]byte
	copy(result[:], workIDBytes[:])
	return result
}

func message(msg string) {
	log.Printf("☑️  %s", msg)
}

func warning(msg string) {
	log.Printf("⚠️  %s", msg)
}

func resolveIneligible(msg string) {
	exit(fmt.Sprintf("✅ %s: this upkeep is not currently elligible", msg), nil, 0)
}

func resolveEligible() {
	exit("❌ this upkeep is currently elligible", nil, 0)
}

func rerun(msg string, err error) {
	exit(fmt.Sprintf("🔁 %s: rerun this command", msg), err, 1)
}

func failUnknown(msg string, err error) {
	exit(fmt.Sprintf("🤷 %s: this should not happen - this script may be broken or your RPC may be experiencing issues", msg), err, 1)
}

func failCheckConfig(msg string, err error) {
	rerun(fmt.Sprintf("%s: check your config", msg), err)
}

func failCheckArgs(msg string, err error) {
	rerun(fmt.Sprintf("%s: check your command arguments", msg), err)
}

func addLink(identifier string, link string) {
	links = append(links, fmt.Sprintf("🔗 %s: %s", identifier, link))
}

func printLinks() {
	for i := 0; i < len(links); i++ {
		log.Println(links[i])
	}
}

func exit(msg string, err error, code int) {
	if err != nil {
		log.Printf("⚠️  %v", err)
	}
	log.Println(msg)
	printLinks()
	os.Exit(code)
}

type TenderlyAPIResponse struct {
	Simulation struct {
		Id string
	}
}

func tenderlySimLink(cfg *config.Config, chainID int64, blockNumber uint64, input []byte, contractAddress gethcommon.Address) string {
	errResult := "<NONE>"
	if cfg.TenderlyAccountName == "" || cfg.TenderlyKey == "" || cfg.TenderlyProjectName == "" {
		warning("tenderly credentials not properly configured - this is optional but helpful")
		return errResult
	}
	values := map[string]interface{}{
		"network_id": fmt.Sprintf("%d", chainID),
		"from":       "0x0000000000000000000000000000000000000000",
		"input":      hexutil.Encode(input),
		"to":         contractAddress.Hex(),
		"gas":        50_000_000,
		"save":       true,
	}
	if blockNumber > 0 {
		values["block_number"] = blockNumber
	}
	jsonData, err := json.Marshal(values)
	if err != nil {
		warning(fmt.Sprintf("unable to marshal tenderly request data: %v", err))
		return errResult
	}
	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.tenderly.co/api/v1/account/%s/project/%s/simulate", cfg.TenderlyAccountName, cfg.TenderlyProjectName),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		warning(fmt.Sprintf("unable to create tenderly request: %v", err))
		return errResult
	}
	request.Header.Set("X-Access-Key", cfg.TenderlyKey)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		warning(fmt.Sprintf("could not run tenderly simulation: %v", err))
		return errResult
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		warning(fmt.Sprintf("unable to read response body from tenderly response: %v", err))
		return errResult
	}
	var responseJSON = &TenderlyAPIResponse{}
	err = json.Unmarshal(body, responseJSON)
	if err != nil {
		warning(fmt.Sprintf("unable to unmarshal tenderly response: %v", err))
		return errResult
	}
	if responseJSON.Simulation.Id == "" {
		warning("unable to simulate tenderly tx")
		return errResult
	}
	return common.TenderlySimLink(responseJSON.Simulation.Id)
}

// TODO - link to performUpkeep tx if exists
