package config

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/stretchr/testify/assert"
)

const chainSpecificConfigsPath = "./defaults/chain_specific_configs.csv"
const defaultConfigsPath = "./defaults/default_configs.csv"
const exportDefaultsToCSV = false
const profilePerms = 0664

func TestConfigDefaultSets(t *testing.T) {
	chainAgnosticDefaultConfigs := []string{
		"OCRObservationTimeout",
		"OCRBlockchainTimeout",
		"OCRTraceLogging",
		"FeatureOffchainReporting",
		"KeeperCheckUpkeepGasPriceFeatureEnabled",
		"KeeperDefaultTransactionQueueDepth",
		"KeeperGasPriceBufferPercent",
		"KeeperGasTipCapBufferPercent",
		"KeeperBaseFeeBufferPercent",
		"KeeperMaximumGracePeriod",
		"KeeperRegistryCheckGasOverhead",
		"KeeperRegistryPerformGasOverhead",
		"KeeperRegistrySyncInterval",
		"KeeperRegistrySyncUpkeepQueueSize",
		"KeeperTurnLookBack",
		"KeeperTurnFlagEnabled",
		"Dev",
		"BlockBackfillDepth",
		"DatabaseLockingMode",
	}
	var data [][]string
	for _, fieldName := range chainAgnosticDefaultConfigs {
		envName := envvar.Name(fieldName)
		defaultValue, err := envvar.DefaultValue(fieldName)
		if !err {
			t.Fatal("Failed to retrieve default value: ", err)
		}
		data = append(data, []string{envName, defaultValue})
	}
	b := writeChainSpecificDefaults(data, defaultConfigsPath, t, exportDefaultsToCSV)
	f, err := os.Open(defaultConfigsPath)
	if err != nil {
		t.Fatal("Failed opening file: ", err)
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(b), string(fileBytes))
}

func TestChainSpecificConfigDefaultSets(t *testing.T) {
	chainData := parseChainSpecificDefaults(chainSpecificConfigDefaultSets)
	b := writeChainSpecificDefaults(chainData, chainSpecificConfigsPath, t, exportDefaultsToCSV)

	f, err := os.Open(chainSpecificConfigsPath)
	if err != nil {
		t.Fatal("Failed opening file: ", err)
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(b), string(fileBytes))

}

func parseChainSpecificDefaults(chainConfigs map[int64]chainSpecificConfigDefaultSet) [][]string {
	orderedChainIDs := make([]int, 0)
	for k := range chainSpecificConfigDefaultSets {
		orderedChainIDs = append(orderedChainIDs, int(k))
	}
	sort.Ints(orderedChainIDs)

	var matrix [][]string
	for index, id := range orderedChainIDs {
		settings := chainConfigs[int64(id)]
		data := [][]string{
			{"EVM_EIP1559_DYNAMIC_FEES", strconv.FormatBool(settings.eip1559DynamicFees)},
			{"ETH_TX_REAPER_INTERVAL", settings.ethTxReaperInterval.String()},
			{"ETH_TX_REAPER_THRESHOLD", settings.ethTxReaperThreshold.String()},
			{"ETH_TX_RESEND_AFTER_THRESHOLD", settings.ethTxResendAfterThreshold.String()},
			{"ETH_FINALITY_DEPTH", strconv.FormatUint(uint64(settings.finalityDepth), 10)},
			{"ETH_GAS_BUMP_PERCENT", strconv.FormatUint(uint64(settings.gasBumpPercent), 10)},
			{"ETH_GAS_BUMP_THRESHOLD", strconv.FormatUint(uint64(settings.gasBumpThreshold), 10)},
			{"ETH_GAS_BUMP_TX_DEPTH", strconv.FormatUint(uint64(settings.gasBumpTxDepth), 10)},
			{"ETH_GAS_BUMP_WEI", settings.gasBumpWei.String()},
			{"EVM_GAS_FEE_CAP_DEFAULT", settings.gasFeeCapDefault.String()},
			{"ETH_GAS_LIMIT_DEFAULT", strconv.FormatUint(uint64(settings.gasLimitDefault), 10)},
			{"ETH_GAS_LIMIT_MULTIPLIER", strconv.FormatFloat(float64(settings.gasLimitMultiplier), 'f', 2, 64)},
			{"ETH_GAS_LIMIT_TRANSFER", strconv.FormatUint(uint64(settings.gasLimitTransfer), 10)},
			{"ETH_GAS_PRICE_DEFAULT", settings.gasPriceDefault.String()},
			{"EVM_GAS_TIP_CAP_DEFAULT", settings.gasTipCapDefault.String()},
			{"EVM_GAS_TIP_CAP_MINIMUM", settings.gasTipCapMinimum.String()},
			{"ETH_HEAD_TRACKER_HISTORY_DEPTH", strconv.FormatUint(uint64(settings.headTrackerHistoryDepth), 10)},
			{"ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", strconv.FormatUint(uint64(settings.headTrackerMaxBufferSize), 10)},
			{"ETH_HEAD_TRACKER_SAMPLING_INTERVAL", settings.headTrackerSamplingInterval.String()},
			{"LINK_CONTRACT_ADDRESS", settings.linkContractAddress},
			{"ETH_LOG_BACKFILL_BATCH_SIZE", strconv.FormatUint(uint64(settings.logBackfillBatchSize), 10)},
			{"ETH_LOG_POLL_INTERVAL", settings.logPollInterval.String()},
			{"ETH_MAX_GAS_PRICE_WEI", settings.maxGasPriceWei.String()},
			{"ETH_MAX_IN_FLIGHT_TRANSACTIONS", strconv.FormatUint(uint64(settings.maxInFlightTransactions), 10)},
			{"ETH_MAX_QUEUED_TRANSACTIONS", strconv.FormatUint(uint64(settings.maxQueuedTransactions), 10)},
			{"ETH_MIN_GAS_PRICE_WEI", settings.minGasPriceWei.String()},
			{"ETH_USE_FORWARDERS", strconv.FormatBool(settings.useForwarders)},
			{"ETH_RPC_DEFAULT_BATCH_SIZE", strconv.FormatUint(uint64(settings.rpcDefaultBatchSize), 10)},
			{"OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT", settings.ocrContractTransmitterTransmitTimeout.String()},
			{"OCR_DATABASE_TIMEOUT", settings.ocrDatabaseTimeout.String()},
			{"OCR_OBSERVATION_GRACE_PERIOD", settings.ocrObservationGracePeriod.String()},

			{"balanceMonitorEnabled", strconv.FormatBool(settings.balanceMonitorEnabled)},
			{"balanceMonitorBlockDelay", strconv.FormatUint(uint64(settings.balanceMonitorBlockDelay), 10)},
			{"blockEmissionIdleWarningThreshold", settings.blockEmissionIdleWarningThreshold.String()},
			{"blockHistoryEstimatorBatchSize", strconv.FormatUint(uint64(settings.blockHistoryEstimatorBatchSize), 10)},
			{"blockHistoryEstimatorBlockDelay", strconv.FormatUint(uint64(settings.blockHistoryEstimatorBlockDelay), 10)},
			{"blockHistoryEstimatorBlockHistorySize", strconv.FormatUint(uint64(settings.blockHistoryEstimatorBlockHistorySize), 10)},
			{"blockHistoryEstimatorTransactionPercentile", strconv.FormatUint(uint64(settings.blockHistoryEstimatorTransactionPercentile), 10)},
			{"chainType", string(settings.chainType)},
			{"gasEstimatorMode", settings.gasEstimatorMode},
			{"minIncomingConfirmations", strconv.FormatUint(uint64(settings.minIncomingConfirmations), 10)},
			{"minimumContractPayment", settings.minimumContractPayment.String()},
			{"nodeDeadAfterNoNewHeadersThreshold", settings.nodeDeadAfterNoNewHeadersThreshold.String()},
			{"nodePollFailureThreshold", strconv.FormatUint(uint64(settings.nodePollFailureThreshold), 10)},
			{"nodePollInterval", settings.nodePollInterval.String()},
			{"nonceAutoSync", strconv.FormatBool(settings.nonceAutoSync)},
			{"complete", strconv.FormatBool(settings.complete)},
			{"ocrContractConfirmations", strconv.FormatUint(uint64(settings.ocrContractConfirmations), 10)},
		}
		// Add config names col.
		if index == 0 {
			matrix = append(matrix, data...)
		} else {
			for i, row := range data {
				matrix[i] = append(matrix[i], row[1])
			}
		}
	}
	var matrixWithHeader [][]string

	// Append header.
	var header []string
	header = append(header, "configValueName")
	for _, id := range orderedChainIDs {
		header = append(header, strconv.FormatInt(int64(id), 10))
	}
	matrixWithHeader = append(matrixWithHeader, header)
	matrixWithHeader = append(matrixWithHeader, matrix...)

	return matrixWithHeader
}

func writeChainSpecificDefaults(chainData [][]string, filePath string, t *testing.T, exportDefaultsToCSV bool) []byte {

	if exportDefaultsToCSV {
		mode := os.O_RDWR | os.O_CREATE | os.O_TRUNC
		csvfile, err := os.OpenFile(filePath, mode, profilePerms)
		if err != nil {
			t.Fatal("Failed creating file: ", err)
		}
		csvwriter := csv.NewWriter(csvfile)
		for _, row := range chainData {
			_ = csvwriter.Write(row)
		}
		csvwriter.Flush()
	}

	var b bytes.Buffer
	byteWriter := csv.NewWriter(&b)
	for _, row := range chainData {
		_ = byteWriter.Write(row)
	}
	byteWriter.Flush()

	return b.Bytes()
}
