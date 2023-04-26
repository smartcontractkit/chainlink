package mercury

import (
	"encoding/json"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

func collectMercuryEnhancedTelemetry(ds *datasource, trrs *pipeline.TaskRunResults, repts ocrtypes.ReportTimestamp) {
	if !shouldCollectEnhancedTelemetryMercury(&ds.jb) || ds.monitoringEndpoint == nil {
		return
	}

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

		parse, _ := ds.parse(*trrs)

		benchmarkPrice := float64(0)
		benchmarkParseTask := trrs.GetNextTaskOf(trr)
		if benchmarkParseTask.Task.Type() == pipeline.TaskTypeJSONParse {
			benchmarkPrice, ok = benchmarkParseTask.Result.Value.(float64)
			if !ok {
				ds.lggr.Warnf("cannot parse enhanced EA telemetry benchmark price, job %d, id %s", ds.jb.ID, trr.Task.DotID())
			}
		}

		bidPrice := float64(0)
		bidParseTask := trrs.GetNextTaskOf(*benchmarkParseTask)
		if bidParseTask.Task.Type() == pipeline.TaskTypeJSONParse {
			bidPrice, ok = bidParseTask.Result.Value.(float64)
			if !ok {
				ds.lggr.Warnf("cannot parse enhanced EA telemetry bid price, job %d, id %s", ds.jb.ID, trr.Task.DotID())
			}
		}

		askPrice := float64(0)
		askParseTask := trrs.GetNextTaskOf(*bidParseTask)
		if askParseTask.Task.Type() == pipeline.TaskTypeJSONParse {
			askPrice, ok = askParseTask.Result.Value.(float64)
			if !ok {
				ds.lggr.Warnf("cannot parse enhanced EA telemetry ask price, job %d, id %s", ds.jb.ID, trr.Task.DotID())
			}
		}

		obsBenchmarkPrice := int64(0)
		obsBid := int64(0)
		obsAsk := int64(0)

		if parse.BenchmarkPrice.Val != nil {
			obsBenchmarkPrice = parse.BenchmarkPrice.Val.Int64()
		}

		if parse.Bid.Val != nil {
			obsBid = parse.Bid.Val.Int64()
		}

		if parse.Ask.Val != nil {
			obsAsk = parse.Ask.Val.Int64()
		}

		t := &telem.EnhancedEAMercury{
			DataSource:                    eaTelem.DataSource,
			DpBenchmarkPrice:              benchmarkPrice,
			DpBid:                         bidPrice,
			DpAsk:                         askPrice,
			CurrentBlockNumber:            parse.CurrentBlockNum.Val,
			CurrentBlockHash:              string(parse.CurrentBlockHash.Val),
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

func getAssetSymbolFromRequestData(requestData string) string {
	type ReqToFrom struct {
		To   string `json:"to"`
		From string `json:"from"`
	}
	type reqData struct {
		Data ReqToFrom `json:"data"`
	}

	rd := &reqData{}
	err := json.Unmarshal([]byte(requestData), rd)
	if err != nil {
		return ""
	}

	return rd.Data.From + "/" + rd.Data.To
}

func shouldCollectEnhancedTelemetryMercury(jb *job.Job) bool {
	if jb.Type.String() == pipeline.OffchainReporting2JobType && jb.OCR2OracleSpec != nil {
		return jb.OCR2OracleSpec.CaptureEATelemetry
	}

	return false
}
