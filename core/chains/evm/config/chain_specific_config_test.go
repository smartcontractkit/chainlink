package config

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const csvPath = "./defaults/chain_configs.csv"
const exportDefaultsToCSV = true
const profilePerms = 0664

func TestChainSpecificConfigDefaultSets(t *testing.T) {

	chainData := parseChainSpecificDefaults(chainSpecificConfigDefaultSets)
	b := writeChainSpecificDefaults(chainData, t, exportDefaultsToCSV)

	f, err := os.Open(csvPath)
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
			{"balanceMonitorEnabled", strconv.FormatBool(settings.balanceMonitorEnabled)},
			{"balanceMonitorBlockDelay", strconv.FormatUint(uint64(settings.balanceMonitorBlockDelay), 10)},
			{"blockEmissionIdleWarningThreshold", settings.blockEmissionIdleWarningThreshold.String()},
			{"blockHistoryEstimatorBatchSize", strconv.FormatUint(uint64(settings.blockHistoryEstimatorBatchSize), 10)},
			{"blockHistoryEstimatorBlockDelay", strconv.FormatUint(uint64(settings.blockHistoryEstimatorBlockDelay), 10)},
			{"blockHistoryEstimatorBlockHistorySize", strconv.FormatUint(uint64(settings.blockHistoryEstimatorBlockHistorySize), 10)},
			{"blockHistoryEstimatorTransactionPercentile", strconv.FormatUint(uint64(settings.blockHistoryEstimatorTransactionPercentile), 10)},
			{"chainType", string(settings.chainType)},
			{"eip1559DynamicFees", strconv.FormatBool(settings.eip1559DynamicFees)},
			{"ethTxReaperInterval", settings.ethTxReaperInterval.String()},
			{"ethTxReaperThreshold", settings.ethTxReaperThreshold.String()},
			{"ethTxResendAfterThreshold", settings.ethTxResendAfterThreshold.String()},
			{"finalityDepth", strconv.FormatUint(uint64(settings.finalityDepth), 10)},
			{"gasBumpPercent", strconv.FormatUint(uint64(settings.gasBumpPercent), 10)},
			{"gasBumpThreshold", strconv.FormatUint(uint64(settings.gasBumpThreshold), 10)},
			{"gasBumpTxDepth", strconv.FormatUint(uint64(settings.gasBumpTxDepth), 10)},
			{"gasBumpWei", settings.gasBumpWei.String()},
			{"gasEstimatorMode", settings.gasEstimatorMode},
			{"gasFeeCapDefault", settings.gasFeeCapDefault.String()},
			{"gasLimitDefault", strconv.FormatUint(uint64(settings.gasLimitDefault), 10)},
			{"gasLimitMultiplier", strconv.FormatFloat(float64(settings.gasLimitMultiplier), 'f', 2, 64)},
			{"gasLimitTransfer", strconv.FormatUint(uint64(settings.gasLimitTransfer), 10)},
			{"gasPriceDefault", settings.gasPriceDefault.String()},
			{"gasTipCapDefault", settings.gasTipCapDefault.String()},
			{"gasTipCapMinimum", settings.gasTipCapMinimum.String()},
			{"headTrackerHistoryDepth", strconv.FormatUint(uint64(settings.headTrackerHistoryDepth), 10)},
			{"headTrackerMaxBufferSize", strconv.FormatUint(uint64(settings.headTrackerMaxBufferSize), 10)},
			{"headTrackerSamplingInterval", settings.headTrackerSamplingInterval.String()},
			{"linkContractAddress", settings.linkContractAddress},
			{"logBackfillBatchSize", strconv.FormatUint(uint64(settings.logBackfillBatchSize), 10)},
			{"logPollInterval", settings.logPollInterval.String()},
			{"maxGasPriceWei", settings.maxGasPriceWei.String()},
			{"maxInFlightTransactions", strconv.FormatUint(uint64(settings.maxInFlightTransactions), 10)},
			{"maxQueuedTransactions", strconv.FormatUint(uint64(settings.maxQueuedTransactions), 10)},
			{"minGasPriceWei", settings.minGasPriceWei.String()},
			{"minIncomingConfirmations", strconv.FormatUint(uint64(settings.minIncomingConfirmations), 10)},
			{"minimumContractPayment", settings.minimumContractPayment.String()},
			{"nodeDeadAfterNoNewHeadersThreshold", settings.nodeDeadAfterNoNewHeadersThreshold.String()},
			{"nodePollFailureThreshold", strconv.FormatUint(uint64(settings.nodePollFailureThreshold), 10)},
			{"nodePollInterval", settings.nodePollInterval.String()},
			{"nonceAutoSync", strconv.FormatBool(settings.nonceAutoSync)},
			{"useForwarders", strconv.FormatBool(settings.useForwarders)},
			{"rpcDefaultBatchSize", strconv.FormatUint(uint64(settings.rpcDefaultBatchSize), 10)},
			{"complete", strconv.FormatBool(settings.complete)},
			{"ocrContractConfirmations", strconv.FormatUint(uint64(settings.ocrContractConfirmations), 10)},
			{"ocrContractTransmitterTransmitTimeout", settings.ocrContractTransmitterTransmitTimeout.String()},
			{"ocrDatabaseTimeout", settings.ocrDatabaseTimeout.String()},
			{"ocrObservationGracePeriod", settings.ocrObservationGracePeriod.String()},
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

func writeChainSpecificDefaults(chainData [][]string, t *testing.T, exportDefaultsToCSV bool) []byte {

	if exportDefaultsToCSV {
		mode := os.O_RDWR | os.O_CREATE | os.O_TRUNC
		csvfile, err := os.OpenFile(csvPath, mode, profilePerms)
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
