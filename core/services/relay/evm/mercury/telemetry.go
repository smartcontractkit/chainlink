package mercury

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

// collectMercuryEnhancedTelemetry checks if enhanced telemetry should be collected, fetches the information needed and
// sends the telemetry
func collectMercuryEnhancedTelemetry(ds *datasource, finalTrrs pipeline.TaskRunResults, trrs *pipeline.TaskRunResults, repts ocrtypes.ReportTimestamp) {
	if !shouldCollectEnhancedTelemetryMercury(&ds.jb) || ds.monitoringEndpoint == nil {
		return
	}

	obsBenchmarkPrice, obsBid, obsAsk, obsBlockNum, obsBlockHash, obsBlockTimestamp := getFinalValues(ds, &finalTrrs)

	for _, trr := range *trrs {
		if trr.Task.Type() != pipeline.TaskTypeBridge {
			continue
		}
		bridgeTask := trr.Task.(*pipeline.BridgeTask)

		bridgeRawResponse, ok := trr.Result.Value.(string)
		if !ok {
			ds.lggr.Warnf("cannot get bridge response from bridge task, job %d, id %s", ds.jb.ID, trr.Task.DotID())
			continue
		}
		eaTelem, err := ocrcommon.ParseEATelemetry([]byte(bridgeRawResponse))
		if err != nil {
			ds.lggr.Warnf("cannot parse EA telemetry, job %d, id %s", ds.jb.ID, trr.Task.DotID())
		}

		assetSymbol := getAssetSymbolFromRequestData(bridgeTask.RequestData)
		benchmarkPrice, bidPrice, askPrice := getPricesFromResults(ds, trr, trrs)

		t := &telem.EnhancedEAMercury{
			DataSource:                    eaTelem.DataSource,
			DpBenchmarkPrice:              benchmarkPrice,
			DpBid:                         bidPrice,
			DpAsk:                         askPrice,
			CurrentBlockNumber:            obsBlockNum,
			CurrentBlockHash:              common.BytesToHash(obsBlockHash).String(),
			CurrentBlockTimestamp:         obsBlockTimestamp,
			BridgeTaskRunStartedTimestamp: trr.CreatedAt.UnixMilli(),
			BridgeTaskRunEndedTimestamp:   trr.FinishedAt.Time.UnixMilli(),
			ProviderRequestedTimestamp:    eaTelem.ProviderRequestedTimestamp,
			ProviderReceivedTimestamp:     eaTelem.ProviderReceivedTimestamp,
			ProviderDataStreamEstablished: eaTelem.ProviderDataStreamEstablished,
			ProviderIndicatedTime:         eaTelem.ProviderIndicatedTime,
			Feed:                          ds.jb.OCR2OracleSpec.FeedID.Hex(),
			ObservationBenchmarkPrice:     obsBenchmarkPrice,
			ObservationBid:                obsBid,
			ObservationAsk:                obsAsk,
			ConfigDigest:                  repts.ConfigDigest.Hex(),
			Round:                         int64(repts.Round),
			Epoch:                         int64(repts.Epoch),
			AssetSymbol:                   assetSymbol,
		}

		bytes, err := proto.Marshal(t)
		if err != nil {
			ds.lggr.Warnf("protobuf marshal failed %v", err.Error())
			continue
		}

		ds.monitoringEndpoint.SendLog(bytes)
	}
}

// getAssetSymbolFromRequestData parses the requestData of the bridge to generate an asset symbol pair
func getAssetSymbolFromRequestData(requestData string) string {
	type reqDataPayload struct {
		To   string `json:"to"`
		From string `json:"from"`
	}
	type reqData struct {
		Data reqDataPayload `json:"data"`
	}

	rd := &reqData{}
	err := json.Unmarshal([]byte(requestData), rd)
	if err != nil {
		return ""
	}

	return rd.Data.From + "/" + rd.Data.To
}

// shouldCollectEnhancedTelemetryMercury checks if enhanced telemetry should be collected and sent
func shouldCollectEnhancedTelemetryMercury(jb *job.Job) bool {
	if jb.Type.String() == pipeline.OffchainReporting2JobType && jb.OCR2OracleSpec != nil {
		return jb.OCR2OracleSpec.CaptureEATelemetry
	}
	return false
}

// getPricesFromResults parses the pipeline.TaskRunResults for pipeline.TaskTypeJSONParse and gets the benchmarkPrice,
// bid and ask. This functions expects the pipeline.TaskRunResults to be correctly ordered
func getPricesFromResults(ds *datasource, startTask pipeline.TaskRunResult, allTasks *pipeline.TaskRunResults) (float64, float64, float64) {
	var benchmarkPrice, askPrice, bidPrice float64
	var ok bool
	//We rely on task results to be sorted in the correct order
	benchmarkPriceTask := allTasks.GetNextTaskOf(startTask)
	if benchmarkPriceTask == nil {
		ds.lggr.Warnf("cannot parse enhanced EA telemetry benchmark price, task is nil, job %d, id %s", ds.jb.ID)
		return 0, 0, 0
	}
	if benchmarkPriceTask.Task.Type() == pipeline.TaskTypeJSONParse {
		benchmarkPrice, ok = benchmarkPriceTask.Result.Value.(float64)
		if !ok {
			ds.lggr.Warnf("cannot parse enhanced EA telemetry benchmark price, job %d, id %s", ds.jb.ID, benchmarkPriceTask.Task.DotID())
		}
	}

	bidTask := allTasks.GetNextTaskOf(*benchmarkPriceTask)
	if bidTask == nil {
		ds.lggr.Warnf("cannot parse enhanced EA telemetry bid price, task is nil, job %d, id %s", ds.jb.ID)
		return 0, 0, 0
	}
	if bidTask.Task.Type() == pipeline.TaskTypeJSONParse {
		bidPrice, ok = bidTask.Result.Value.(float64)
		if !ok {
			ds.lggr.Warnf("cannot parse enhanced EA telemetry bid price, job %d, id %s", ds.jb.ID, bidTask.Task.DotID())
		}
	}

	askTask := allTasks.GetNextTaskOf(*bidTask)
	if askTask == nil {
		ds.lggr.Warnf("cannot parse enhanced EA telemetry ask price, task is nil, job %d, id %s", ds.jb.ID)
		return 0, 0, 0
	}
	if askTask.Task.Type() == pipeline.TaskTypeJSONParse {
		askPrice, ok = askTask.Result.Value.(float64)
		if !ok {
			ds.lggr.Warnf("cannot parse enhanced EA telemetry ask price, job %d, id %s", ds.jb.ID, askTask.Task.DotID())
		}
	}

	return benchmarkPrice, bidPrice, askPrice
}

// getFinalValues runs a parse on the pipeline.TaskRunResults and returns the values
func getFinalValues(ds *datasource, trrs *pipeline.TaskRunResults) (int64, int64, int64, int64, []byte, uint64) {
	var benchmarkPrice, bid, ask int64
	parse, _ := ds.parse(*trrs)

	if parse.BenchmarkPrice.Val != nil {
		benchmarkPrice = parse.BenchmarkPrice.Val.Int64()
	}
	if parse.Bid.Val != nil {
		bid = parse.Bid.Val.Int64()
	}
	if parse.Ask.Val != nil {
		ask = parse.Ask.Val.Int64()
	}

	return benchmarkPrice, bid, ask, parse.CurrentBlockNum.Val, parse.CurrentBlockHash.Val, parse.CurrentBlockTimestamp.Val
}
