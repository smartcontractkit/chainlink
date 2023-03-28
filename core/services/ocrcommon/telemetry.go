package ocrcommon

import (
	"encoding/json"
	"math/big"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type eaTelemetryResponse struct {
	DataSource                    string `json:"data_source"`
	ProviderRequestedProtocol     string `json:"provider_requested_protocol"`
	ProviderRequestedTimestamp    int64  `json:"provider_requested_timestamp"`
	ProviderReceivedTimestamp     int64  `json:"provider_received_timestamp"`
	ProviderDataStreamEstablished int64  `json:"provider_data_stream_established"`
	ProviderDataReceived          int64  `json:"provider_data_received"`
	ProviderIndicatedTime         int64  `json:"provider_indicated_time"`
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
		Telemetry eaTelemetryResponse `json:"telemetry"`
	}
	gr := generalResponse{}

	if err := json.Unmarshal(b, &gr); err != nil {
		return eaTelemetryResponse{}, err
	}

	return gr.Telemetry, nil
}

// getJsonParsedValue checks if the next logical task is of type pipeline.TaskTypeJSONParse and trys to return
// the response as a *big.Int
func getJsonParsedValue(trr pipeline.TaskRunResult, trrs *pipeline.TaskRunResults) *big.Int {
	nextTask := trrs.GetNextTaskOf(trr)
	if nextTask != nil && nextTask.Task.Type() == pipeline.TaskTypeJSONParse {
		asDecimal, err := utils.ToDecimal(nextTask.Result.Value)
		if err != nil {
			return nil
		}
		return asDecimal.BigInt()
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

func getParsedValue(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults, trr pipeline.TaskRunResult) int64 {
	parsedValue := getJsonParsedValue(trr, trrs)
	if parsedValue == nil {
		ds.lggr.Warnf("cannot get json parse value, job %d, id %s", ds.jb.ID, trr.Task.DotID())
		return 0
	}
	return parsedValue.Int64()
}

// collectEATelemetry checks if EA telemetry should be collected, gathers the information and sends it for ingestion
func collectEATelemetry(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults, finalResult *pipeline.FinalResult) {
	if !shouldCollectTelemetry(&ds.jb) || ds.monitoringEndpoint == nil {
		return
	}

	go collectAndSend(ds, trrs, finalResult)
}

func collectAndSend(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults, finalResult *pipeline.FinalResult) {
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
			ProviderRequestedProtocol:     eaTelemetry.ProviderRequestedProtocol,
			ProviderRequestedTimestamp:    eaTelemetry.ProviderRequestedTimestamp,
			ProviderReceivedTimestamp:     eaTelemetry.ProviderReceivedTimestamp,
			ProviderDataStreamEstablished: eaTelemetry.ProviderDataStreamEstablished,
			ProviderDataReceived:          eaTelemetry.ProviderDataReceived,
			ProviderIndicatedTime:         eaTelemetry.ProviderIndicatedTime,
			Feed:                          contract,
			ChainId:                       chainID,
			Observation:                   observation,
			ConfigDigest:                  "",
			Round:                         0,
			Epoch:                         0,
		}

		bytes, err := proto.Marshal(t)
		if err != nil {
			ds.lggr.Warnf("protobuf marshal failed %v", err.Error())
			continue
		}
		ds.monitoringEndpoint.SendLog(bytes)
	}
}
