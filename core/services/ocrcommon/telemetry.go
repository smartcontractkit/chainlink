package ocrcommon

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type eaTelemetryResponse struct {
	DataSource                    string `json:"dataSource"`
	ProviderRequestedTimestamp    int64  `json:"providerDataRequestedUnixMs"`
	ProviderReceivedTimestamp     int64  `json:"providerDataReceivedUnixMs"`
	ProviderDataStreamEstablished int64  `json:"providerDataStreamEstablishedUnixMs"`
	ProviderIndicatedTime         int64  `json:"providerIndicatedTimeUnixMs"`
}

// shouldCollectTelemetry returns whether EA telemetry should be collected
func shouldCollectTelemetry(jb *job.Job) bool {
	if jb.Type.String() == pipeline.OffchainReportingJobType && jb.OCROracleSpec != nil {
		return jb.OCROracleSpec.CaptureEATelemetry
	}

	if jb.Type.String() == pipeline.OffchainReporting2JobType && jb.OCR2OracleSpec != nil {
		return jb.OCR2OracleSpec.CaptureEATelemetry
	}

	return false
}

// getContract fetches the contract address from the OracleSpec
func getContract(jb *job.Job) string {
	switch jb.Type.String() {
	case pipeline.OffchainReportingJobType:
		return jb.OCROracleSpec.ContractAddress.String()
	case pipeline.OffchainReporting2JobType:
		return jb.OCR2OracleSpec.ContractID
	default:
		return ""
	}
}

// getChainID fetches the chain id from the OracleSpec
func getChainID(jb *job.Job) string {
	switch jb.Type.String() {
	case pipeline.OffchainReportingJobType:
		return jb.OCROracleSpec.EVMChainID.String()
	case pipeline.OffchainReporting2JobType:
		contract, _ := jb.OCR2OracleSpec.RelayConfig["chainID"].(string)
		return contract
	default:
		return ""
	}
}

// parseEATelemetry attempts to parse the bridge telemetry
func parseEATelemetry(b []byte) (eaTelemetryResponse, error) {
	type generalResponse struct {
		TelemTimestamps eaTelemetryResponse `json:"timestamps"`
	}
	gr := generalResponse{}

	if err := json.Unmarshal(b, &gr); err != nil {
		return eaTelemetryResponse{}, err
	}

	return gr.TelemTimestamps, nil
}

// getJsonParsedValue checks if the next logical task is of type pipeline.TaskTypeJSONParse and trys to return
// the response as a *big.Int
func getJsonParsedValue(trr pipeline.TaskRunResult, trrs *pipeline.TaskRunResults) *float64 {
	nextTask := trrs.GetNextTaskOf(trr)
	if nextTask != nil && nextTask.Task.Type() == pipeline.TaskTypeJSONParse {
		asDecimal, err := utils.ToDecimal(nextTask.Result.Value)
		if err != nil {
			return nil
		}
		toFloat, _ := asDecimal.Float64()
		return &toFloat
	}
	return nil
}

// getObservation checks pipeline.FinalResult and extracts the observation
func getObservation(ds *inMemoryDataSource, finalResult *pipeline.FinalResult) int64 {
	singularResult, err := finalResult.SingularResult()
	if err != nil {
		ds.lggr.Warnf("cannot get singular result, job %d", ds.jb.ID)
		return 0
	}

	finalResultDecimal, err := utils.ToDecimal(singularResult.Value)
	if err != nil {
		ds.lggr.Warnf("cannot parse singular result from bridge task, job %d", ds.jb.ID)
		return 0
	}

	return finalResultDecimal.BigInt().Int64()
}

func getParsedValue(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults, trr pipeline.TaskRunResult) float64 {
	parsedValue := getJsonParsedValue(trr, trrs)
	if parsedValue == nil {
		ds.lggr.Warnf("cannot get json parse value, job %d, id %s", ds.jb.ID, trr.Task.DotID())
		return 0
	}
	return *parsedValue
}

// collectEATelemetry checks if EA telemetry should be collected, gathers the information and sends it for ingestion
func collectEATelemetry(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults, finalResult *pipeline.FinalResult, timestamp ObservationTimestamp) {
	if !shouldCollectTelemetry(&ds.jb) || ds.monitoringEndpoint == nil {
		return
	}

	go collectAndSend(ds, trrs, finalResult, timestamp)
}

func collectAndSend(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults, finalResult *pipeline.FinalResult, timestamp ObservationTimestamp) {
	chainID := getChainID(&ds.jb)
	contract := getContract(&ds.jb)

	observation := getObservation(ds, finalResult)

	for _, trr := range *trrs {
		if trr.Task.Type() != pipeline.TaskTypeBridge {
			continue
		}

		bridgeRawResponse, ok := trr.Result.Value.(string)
		if !ok {
			ds.lggr.Warnf("cannot get bridge response from bridge task, job %d, id %s", ds.jb.ID, trr.Task.DotID())
			continue
		}
		eaTelemetry, err := parseEATelemetry([]byte(bridgeRawResponse))
		if err != nil {
			ds.lggr.Warnf("cannot parse EA telemetry, job %d, id %s", ds.jb.ID, trr.Task.DotID())
		}
		value := getParsedValue(ds, trrs, trr)

		t := &telem.EnhancedEA{
			DataSource:                    eaTelemetry.DataSource,
			Value:                         value,
			BridgeTaskRunStartedTimestamp: trr.CreatedAt.UnixMilli(),
			BridgeTaskRunEndedTimestamp:   trr.FinishedAt.Time.UnixMilli(),
			ProviderRequestedTimestamp:    eaTelemetry.ProviderRequestedTimestamp,
			ProviderReceivedTimestamp:     eaTelemetry.ProviderReceivedTimestamp,
			ProviderDataStreamEstablished: eaTelemetry.ProviderDataStreamEstablished,
			ProviderIndicatedTime:         eaTelemetry.ProviderIndicatedTime,
			Feed:                          contract,
			ChainId:                       chainID,
			Observation:                   observation,
			ConfigDigest:                  timestamp.ConfigDigest,
			Round:                         int64(timestamp.Round),
			Epoch:                         int64(timestamp.Epoch),
		}

		bytes, err := proto.Marshal(t)
		if err != nil {
			ds.lggr.Warnf("protobuf marshal failed %v", err.Error())
			continue
		}

		ds.monitoringEndpoint.SendLog(bytes)
	}
}
