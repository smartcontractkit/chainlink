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

	types2 "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	evm21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21"

	commonhex "github.com/smartcontractkit/chainlink-common/pkg/utils/hex"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/streams"
	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

const (
	ConditionTrigger uint8 = iota
	LogTrigger
	expectedTypeAndVersion = "KeeperRegistry 2.1.0"
)

var mercuryPacker = mercury.NewAbiPacker()
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

	// Log triggers: always use block from tx
	// Conditional: use latest block if no block number is provided, otherwise use block from user input
	var triggerCallOpts *bind.CallOpts             // use a certain block
	latestCallOpts := &bind.CallOpts{Context: ctx} // use the latest block

	// connect to registry contract
	registryAddress := gethcommon.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry21, err := iregistry21.NewIKeeperRegistryMaster(registryAddress, k.client)
	if err != nil {
		failUnknown("failed to connect to the registry contract", err)
	}

	// verify contract is correct
	typeAndVersion, err := keeperRegistry21.TypeAndVersion(latestCallOpts)
	if err != nil {
		failCheckConfig("failed to get typeAndVersion: make sure your registry contract address and archive node are valid", err)
	}
	if typeAndVersion != expectedTypeAndVersion {
		failCheckConfig(fmt.Sprintf("invalid registry contract: this command can only debug %s, got: %s", expectedTypeAndVersion, typeAndVersion), nil)
	}
	// get upkeepID from command args
	upkeepID := big.NewInt(0)
	upkeepIDNoPrefix := commonhex.TrimPrefix(args[0])
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
	addLink("upkeep link", common.UpkeepLink(chainID, upkeepID))
	addLink("upkeep contract address", common.ContractExplorerLink(chainID, upkeepInfo.Target))
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
	var workID [32]byte
	var trigger ocr2keepers.Trigger
	upkeepNeeded := false
	// check upkeep
	if triggerType == ConditionTrigger {
		message("upkeep identified as conditional trigger")

		if len(args) > 1 {
			// if a block number is provided, use that block for both checkUpkeep and simulatePerformUpkeep
			blockNum, err = strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				failCheckArgs("unable to parse block number", err)
			}
			triggerCallOpts = &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(blockNum)}
		} else {
			// if no block number is provided, use latest block for both checkUpkeep and simulatePerformUpkeep
			triggerCallOpts = latestCallOpts
		}

		var tmpCheckResult iregistry21.CheckUpkeep0
		tmpCheckResult, err = keeperRegistry21.CheckUpkeep0(triggerCallOpts, upkeepID)
		if err != nil {
			failUnknown("failed to check upkeep: ", err)
		}
		checkResult = iregistry21.CheckUpkeep(tmpCheckResult)
		// do tenderly simulation
		var rawCall []byte
		rawCall, err = core.RegistryABI.Pack("checkUpkeep", upkeepID, []byte{})
		if err != nil {
			failUnknown("failed to pack raw checkUpkeep call", err)
		}
		addLink("checkUpkeep simulation", tenderlySimLink(ctx, k.cfg, chainID, 0, rawCall, registryAddress))
	} else if triggerType == LogTrigger {
		// validate inputs
		message("upkeep identified as log trigger")
		if len(args) != 3 {
			failCheckArgs("txHash and log index must be supplied to command in order to debug log triggered upkeeps", nil)
		}
		txHash := gethcommon.HexToHash(args[1])

		var logIndex int64
		logIndex, err = strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			failCheckArgs("unable to parse log index", err)
		}

		// check that tx is confirmed
		var isPending bool
		_, isPending, err = k.client.TransactionByHash(ctx, txHash)
		if err != nil {
			log.Fatal("failed to get tx by hash", err)
		}
		if isPending {
			resolveIneligible(fmt.Sprintf("tx %s is still pending confirmation", txHash))
		}

		// find transaction receipt
		var receipt *types.Receipt
		receipt, err = k.client.TransactionReceipt(ctx, txHash)
		if err != nil {
			failCheckArgs("failed to fetch tx receipt", err)
		}
		addLink("trigger transaction", common.ExplorerLink(chainID, txHash))
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
		trigger = mustAutomationTrigger(txHash, logIndex, blockNum, receipt.BlockHash)
		workID = mustUpkeepWorkID(upkeepID, trigger)
		message(fmt.Sprintf("workID computed: %s", hex.EncodeToString(workID[:])))
		var hasKey bool
		hasKey, err = keeperRegistry21.HasDedupKey(latestCallOpts, workID)
		if err != nil {
			failUnknown("failed to check if upkeep was already performed: ", err)
		}
		if hasKey {
			resolveIneligible("upkeep was already performed")
		}
		triggerCallOpts = &bind.CallOpts{Context: ctx, BlockNumber: big.NewInt(receipt.BlockNumber.Int64())}
		var rawTriggerConfig []byte
		rawTriggerConfig, err = keeperRegistry21.GetUpkeepTriggerConfig(triggerCallOpts, upkeepID)
		if err != nil {
			failUnknown("failed to fetch trigger config for upkeep", err)
		}
		var triggerConfig automation_utils_2_1.LogTriggerConfig
		triggerConfig, err = packer.UnpackLogTriggerConfig(rawTriggerConfig)
		if err != nil {
			failUnknown("failed to unpack trigger config", err)
		}
		if triggerConfig.FilterSelector > 7 {
			resolveIneligible(fmt.Sprintf("invalid filter selector %d", triggerConfig.FilterSelector))
		}
		if !logMatchesTriggerConfig(triggeringEvent, triggerConfig) {
			resolveIneligible("log does not match trigger config")
		}
		var header *types.Header
		header, err = k.client.HeaderByHash(ctx, receipt.BlockHash)
		if err != nil {
			failUnknown("failed to find block", err)
		}
		var triggerData []byte
		triggerData, err = packTriggerData(triggeringEvent, header.Time)
		if err != nil {
			failUnknown("failed to pack trigger data", err)
		}
		checkResult, err = keeperRegistry21.CheckUpkeep(triggerCallOpts, upkeepID, triggerData)
		if err != nil {
			failUnknown("failed to check upkeep", err)
		}
		// do tenderly simulations
		var rawCall []byte
		rawCall, err = core.RegistryABI.Pack("checkUpkeep", upkeepID, triggerData)
		if err != nil {
			failUnknown("failed to pack raw checkUpkeep call", err)
		}
		addLink("checkUpkeep simulation", tenderlySimLink(ctx, k.cfg, chainID, blockNum, rawCall, registryAddress))
		rawCall = append(core.ILogAutomationABI.Methods["checkLog"].ID, triggerData...)
		addLink("checkLog (direct) simulation", tenderlySimLink(ctx, k.cfg, chainID, blockNum, rawCall, upkeepInfo.Target))
	} else {
		resolveIneligible(fmt.Sprintf("invalid trigger type: %d", triggerType))
	}
	upkeepNeeded, performData = checkResult.UpkeepNeeded, checkResult.PerformData

	if checkResult.UpkeepFailureReason != 0 {
		message(fmt.Sprintf("checkUpkeep failed with UpkeepFailureReason %s", getCheckUpkeepFailureReason(checkResult.UpkeepFailureReason)))
	}

	// handle data streams lookup
	if checkResult.UpkeepFailureReason == uint8(encoding.UpkeepFailureReasonTargetCheckReverted) {
		mc := &types2.MercuryCredentials{LegacyURL: k.cfg.DataStreamsLegacyURL, URL: k.cfg.DataStreamsURL, Username: k.cfg.DataStreamsID, Password: k.cfg.DataStreamsKey}
		mercuryConfig := evm21.NewMercuryConfig(mc, core.StreamsCompatibleABI)
		lggr, _ := logger.NewLogger()
		blockSub := &blockSubscriber{k.client}
		streams := streams.NewStreamsLookup(mercuryConfig, blockSub, k.rpcClient, keeperRegistry21, lggr)

		var streamsLookupErr *mercury.StreamsLookupError
		streamsLookupErr, err = mercuryPacker.DecodeStreamsLookupRequest(checkResult.PerformData)
		if err == nil {
			message("upkeep reverted with StreamsLookup")
			message(fmt.Sprintf("StreamsLookup data: {FeedParamKey: %s, Feeds: %v, TimeParamKey: %s, Time: %d, ExtraData: %s}", streamsLookupErr.FeedParamKey, streamsLookupErr.Feeds, streamsLookupErr.TimeParamKey, streamsLookupErr.Time.Uint64(), hexutil.Encode(streamsLookupErr.ExtraData)))

			streamsLookup := &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: streamsLookupErr.FeedParamKey,
					Feeds:        streamsLookupErr.Feeds,
					TimeParamKey: streamsLookupErr.TimeParamKey,
					Time:         streamsLookupErr.Time,
					ExtraData:    streamsLookupErr.ExtraData,
				},
				UpkeepId: upkeepID,
				Block:    blockNum,
			}

			if streamsLookup.IsMercuryV02() {
				message("using data streams lookup v0.2")
				// check if upkeep is allowed to use mercury v0.2
				var allowed bool
				_, _, _, allowed, err = streams.AllowedToUseMercury(triggerCallOpts, upkeepID)
				if err != nil {
					failUnknown("failed to check if upkeep is allowed to use data streams", err)
				}
				if !allowed {
					resolveIneligible("upkeep reverted with StreamsLookup but is not allowed to access streams")
				}
			} else if streamsLookup.IsMercuryV03() {
				// handle v0.3
				message("using data streams lookup v0.3")
			} else {
				resolveIneligible("upkeep reverted with StreamsLookup but the configuration is invalid")
			}

			if k.cfg.DataStreamsLegacyURL == "" || k.cfg.DataStreamsURL == "" || k.cfg.DataStreamsID == "" || k.cfg.DataStreamsKey == "" {
				failCheckConfig("Data streams configs not set properly, check your DATA_STREAMS_LEGACY_URL, DATA_STREAMS_URL, DATA_STREAMS_ID and DATA_STREAMS_KEY", nil)
			}

			// do mercury request
			automationCheckResult := mustAutomationCheckResult(upkeepID, checkResult, trigger)
			checkResults := []ocr2keepers.CheckResult{automationCheckResult}

			var values [][]byte
			values, err = streams.DoMercuryRequest(ctx, streamsLookup, checkResults, 0)

			if checkResults[0].IneligibilityReason == uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput) {
				resolveIneligible("upkeep used invalid revert data")
			}
			if checkResults[0].PipelineExecutionState == uint8(encoding.InvalidMercuryRequest) {
				resolveIneligible("the data streams request data is invalid")
			}
			if err != nil {
				failCheckConfig("failed to do data streams request ", err)
			}

			// do checkCallback
			err = streams.CheckCallback(ctx, values, streamsLookup, checkResults, 0)
			if err != nil {
				failUnknown("failed to execute data streams callback ", err)
			}
			if checkResults[0].IneligibilityReason != 0 {
				message(fmt.Sprintf("checkCallback failed with UpkeepFailureReason %d", checkResults[0].IneligibilityReason))
			}
			upkeepNeeded, performData = checkResults[0].Eligible, checkResults[0].PerformData
			// do tenderly simulations for checkCallback
			var rawCall []byte
			rawCall, err = core.RegistryABI.Pack("checkCallback", upkeepID, values, streamsLookup.ExtraData)
			if err != nil {
				failUnknown("failed to pack raw checkCallback call", err)
			}
			addLink("checkCallback simulation", tenderlySimLink(ctx, k.cfg, chainID, blockNum, rawCall, registryAddress))
			rawCall, err = core.StreamsCompatibleABI.Pack("checkCallback", values, streamsLookup.ExtraData)
			if err != nil {
				failUnknown("failed to pack raw checkCallback (direct) call", err)
			}
			addLink("checkCallback (direct) simulation", tenderlySimLink(ctx, k.cfg, chainID, blockNum, rawCall, upkeepInfo.Target))
		} else {
			message("did not revert with StreamsLookup error")
		}
	}
	if !upkeepNeeded {
		resolveIneligible("upkeep is not needed")
	}
	// simulate perform upkeep
	simulateResult, err := keeperRegistry21.SimulatePerformUpkeep(triggerCallOpts, upkeepID, performData)
	if err != nil {
		failUnknown("failed to simulate perform upkeep: ", err)
	}

	// do tenderly simulation
	rawCall, err := core.RegistryABI.Pack("simulatePerformUpkeep", upkeepID, performData)
	if err != nil {
		failUnknown("failed to pack raw simulatePerformUpkeep call", err)
	}
	addLink("simulatePerformUpkeep simulation", tenderlySimLink(ctx, k.cfg, chainID, blockNum, rawCall, registryAddress))

	if simulateResult.Success {
		resolveEligible()
	} else {
		// Convert performGas to *big.Int for comparison
		performGasBigInt := new(big.Int).SetUint64(uint64(upkeepInfo.PerformGas))
		// Compare PerformGas and GasUsed
		result := performGasBigInt.Cmp(simulateResult.GasUsed)

		if result < 0 {
			// PerformGas is smaller than GasUsed
			resolveIneligible(fmt.Sprintf("simulate perform upkeep unsuccessful, PerformGas (%d) is lower than GasUsed (%s)", upkeepInfo.PerformGas, simulateResult.GasUsed.String()))
		} else {
			resolveIneligible("simulate perform upkeep unsuccessful")
		}
	}
}

func getCheckUpkeepFailureReason(reasonIndex uint8) string {
	// Copied from KeeperRegistryBase2_1.sol
	reasonStrings := []string{
		"NONE",
		"UPKEEP_CANCELLED",
		"UPKEEP_PAUSED",
		"TARGET_CHECK_REVERTED",
		"UPKEEP_NOT_NEEDED",
		"PERFORM_DATA_EXCEEDS_LIMIT",
		"INSUFFICIENT_BALANCE",
		"CALLBACK_REVERTED",
		"REVERT_DATA_EXCEEDS_LIMIT",
		"REGISTRY_PAUSED",
	}

	if int(reasonIndex) < len(reasonStrings) {
		return reasonStrings[reasonIndex]
	}

	return fmt.Sprintf("Unknown : %d", reasonIndex)
}

func mustAutomationCheckResult(upkeepID *big.Int, checkResult iregistry21.CheckUpkeep, trigger ocr2keepers.Trigger) ocr2keepers.CheckResult {
	upkeepIdentifier := mustUpkeepIdentifier(upkeepID)
	checkResult2 := ocr2keepers.CheckResult{
		Eligible:            checkResult.UpkeepNeeded,
		IneligibilityReason: checkResult.UpkeepFailureReason,
		UpkeepID:            upkeepIdentifier,
		Trigger:             trigger,
		WorkID:              core.UpkeepWorkID(upkeepIdentifier, trigger),
		GasAllocated:        0,
		PerformData:         checkResult.PerformData,
		FastGasWei:          checkResult.FastGasWei,
		LinkNative:          checkResult.LinkNative,
	}

	return checkResult2
}

type blockSubscriber struct {
	ethClient *ethclient.Client
}

func (bs *blockSubscriber) LatestBlock() *ocr2keepers.BlockKey {
	header, err := bs.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil
	}

	return &ocr2keepers.BlockKey{
		Number: ocr2keepers.BlockNumber(header.Number.Uint64()),
		Hash:   header.Hash(),
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

func mustUpkeepWorkID(upkeepID *big.Int, trigger ocr2keepers.Trigger) [32]byte {
	upkeepIdentifier := mustUpkeepIdentifier(upkeepID)

	workID := core.UpkeepWorkID(upkeepIdentifier, trigger)
	workIDBytes, err := hex.DecodeString(workID)
	if err != nil {
		failUnknown("failed to decode workID", err)
	}

	var result [32]byte
	copy(result[:], workIDBytes[:])
	return result
}

func mustUpkeepIdentifier(upkeepID *big.Int) ocr2keepers.UpkeepIdentifier {
	upkeepIdentifier := &ocr2keepers.UpkeepIdentifier{}
	upkeepIdentifier.FromBigInt(upkeepID)
	return *upkeepIdentifier
}

func mustAutomationTrigger(txHash [32]byte, logIndex int64, blockNum uint64, blockHash [32]byte) ocr2keepers.Trigger {
	trigger := ocr2keepers.Trigger{
		LogTriggerExtension: &ocr2keepers.LogTriggerExtension{
			TxHash:      txHash,
			Index:       uint32(logIndex),
			BlockNumber: ocr2keepers.BlockNumber(blockNum),
			BlockHash:   blockHash,
		},
	}
	return trigger
}

func message(msg string) {
	log.Printf("☑️  %s", msg)
}

func warning(msg string) {
	log.Printf("⚠️  %s", msg)
}

func resolveIneligible(msg string) {
	exit(fmt.Sprintf("❌ this upkeep is not eligible: %s", msg), nil, 0)
}

func resolveEligible() {
	exit("✅ this upkeep is eligible", nil, 0)
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

func tenderlySimLink(ctx context.Context, cfg *config.Config, chainID int64, blockNumber uint64, input []byte, contractAddress gethcommon.Address) string {
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
	request, err := http.NewRequestWithContext(
		ctx,
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
