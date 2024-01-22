package ocrcommon

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	v1types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v2types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	v3types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

type eaTelemetry struct {
	DataSource                    string
	ProviderRequestedTimestamp    int64
	ProviderReceivedTimestamp     int64
	ProviderDataStreamEstablished int64
	ProviderIndicatedTime         int64
}

type EnhancedTelemetryData struct {
	TaskRunResults pipeline.TaskRunResults
	FinalResults   pipeline.FinalResult
	RepTimestamp   ObservationTimestamp
}

type EnhancedTelemetryMercuryData struct {
	V1Observation              *v1types.Observation
	V2Observation              *v2types.Observation
	V3Observation              *v3types.Observation
	TaskRunResults             pipeline.TaskRunResults
	RepTimestamp               ocrtypes.ReportTimestamp
	FeedVersion                mercuryutils.FeedVersion
	FetchMaxFinalizedTimestamp bool
	IsLinkFeed                 bool
	IsNativeFeed               bool
}

type EnhancedTelemetryService[T EnhancedTelemetryData | EnhancedTelemetryMercuryData] struct {
	services.StateMachine

	chTelem            <-chan T
	chDone             chan struct{}
	monitoringEndpoint commontypes.MonitoringEndpoint
	job                *job.Job
	lggr               logger.Logger
}

func NewEnhancedTelemetryService[T EnhancedTelemetryData | EnhancedTelemetryMercuryData](job *job.Job, chTelem <-chan T, done chan struct{}, me commontypes.MonitoringEndpoint, lggr logger.Logger) *EnhancedTelemetryService[T] {
	return &EnhancedTelemetryService[T]{
		chTelem:            chTelem,
		chDone:             done,
		monitoringEndpoint: me,
		lggr:               lggr,
		job:                job,
	}
}

// Start starts
func (e *EnhancedTelemetryService[T]) Start(context.Context) error {
	return e.StartOnce("EnhancedTelemetryService", func() error {
		go func() {
			e.lggr.Infof("Started enhanced telemetry service for job %d", e.job.ID)
			for {
				select {
				case t := <-e.chTelem:
					switch v := any(t).(type) {
					case EnhancedTelemetryData:
						e.collectEATelemetry(v.TaskRunResults, v.FinalResults, v.RepTimestamp)
					case EnhancedTelemetryMercuryData:
						e.collectMercuryEnhancedTelemetry(v)
					default:
						e.lggr.Errorf("unrecognised telemetry data type: %T", t)
					}
				case <-e.chDone:
					return
				}
			}
		}()
		return nil
	})
}

func (e *EnhancedTelemetryService[T]) Close() error {
	return e.StopOnce("EnhancedTelemetryService", func() error {
		close(e.chDone)
		e.lggr.Infof("Stopping enhanced telemetry service for job %d", e.job.ID)
		return nil
	})
}

// ShouldCollectEnhancedTelemetry returns whether EA telemetry should be collected
func ShouldCollectEnhancedTelemetry(job *job.Job) bool {
	if job.Type.String() == pipeline.OffchainReportingJobType && job.OCROracleSpec != nil {
		return job.OCROracleSpec.CaptureEATelemetry
	}

	if job.Type.String() == pipeline.OffchainReporting2JobType && job.OCR2OracleSpec != nil {
		return job.OCR2OracleSpec.CaptureEATelemetry
	}

	return false
}

// getContract fetches the contract address from the OracleSpec
func (e *EnhancedTelemetryService[T]) getContract() string {
	switch e.job.Type.String() {
	case pipeline.OffchainReportingJobType:
		return e.job.OCROracleSpec.ContractAddress.String()
	case pipeline.OffchainReporting2JobType:
		return e.job.OCR2OracleSpec.ContractID
	default:
		return ""
	}
}

// getChainID fetches the chain id from the OracleSpec
func (e *EnhancedTelemetryService[T]) getChainID() string {
	switch e.job.Type.String() {
	case pipeline.OffchainReportingJobType:
		return e.job.OCROracleSpec.EVMChainID.String()
	case pipeline.OffchainReporting2JobType:
		contract, _ := e.job.OCR2OracleSpec.RelayConfig["chainID"].(string)
		return contract
	default:
		return ""
	}
}

// parseEATelemetry attempts to parse the bridge telemetry
func parseEATelemetry(b []byte) (eaTelemetry, error) {
	type eaTimestamps struct {
		ProviderRequestedTimestamp    int64 `json:"providerDataRequestedUnixMs"`
		ProviderReceivedTimestamp     int64 `json:"providerDataReceivedUnixMs"`
		ProviderDataStreamEstablished int64 `json:"providerDataStreamEstablishedUnixMs"`
		ProviderIndicatedTime         int64 `json:"providerIndicatedTimeUnixMs"`
	}
	type eaMeta struct {
		AdapterName string `json:"adapterName"`
	}

	type eaTelem struct {
		TelemTimestamps eaTimestamps `json:"timestamps"`
		TelemMeta       eaMeta       `json:"meta"`
	}
	t := eaTelem{}

	if err := json.Unmarshal(b, &t); err != nil {
		return eaTelemetry{}, err
	}

	return eaTelemetry{
		DataSource:                    t.TelemMeta.AdapterName,
		ProviderRequestedTimestamp:    t.TelemTimestamps.ProviderRequestedTimestamp,
		ProviderReceivedTimestamp:     t.TelemTimestamps.ProviderReceivedTimestamp,
		ProviderDataStreamEstablished: t.TelemTimestamps.ProviderDataStreamEstablished,
		ProviderIndicatedTime:         t.TelemTimestamps.ProviderIndicatedTime,
	}, nil
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
func (e *EnhancedTelemetryService[T]) getObservation(finalResult *pipeline.FinalResult) int64 {
	singularResult, err := finalResult.SingularResult()
	if err != nil {
		e.lggr.Warnf("cannot get singular result, job %d", e.job.ID)
		return 0
	}

	finalResultDecimal, err := utils.ToDecimal(singularResult.Value)
	if err != nil {
		e.lggr.Warnf("cannot parse singular result from bridge task, job %d", e.job.ID)
		return 0
	}

	return finalResultDecimal.BigInt().Int64()
}

func (e *EnhancedTelemetryService[T]) getParsedValue(trrs *pipeline.TaskRunResults, trr pipeline.TaskRunResult) float64 {
	parsedValue := getJsonParsedValue(trr, trrs)
	if parsedValue == nil {
		e.lggr.Warnf("cannot get json parse value, job %d, id %s", e.job.ID, trr.Task.DotID())
		return 0
	}
	return *parsedValue
}

// collectEATelemetry checks if EA telemetry should be collected, gathers the information and sends it for ingestion
func (e *EnhancedTelemetryService[T]) collectEATelemetry(trrs pipeline.TaskRunResults, finalResult pipeline.FinalResult, timestamp ObservationTimestamp) {
	if e.monitoringEndpoint == nil {
		return
	}

	e.collectAndSend(&trrs, &finalResult, timestamp)
}

func (e *EnhancedTelemetryService[T]) collectAndSend(trrs *pipeline.TaskRunResults, finalResult *pipeline.FinalResult, timestamp ObservationTimestamp) {
	chainID := e.getChainID()
	contract := e.getContract()

	observation := e.getObservation(finalResult)

	for _, trr := range *trrs {
		if trr.Task.Type() != pipeline.TaskTypeBridge {
			continue
		}
		var bridgeName string
		if b, is := trr.Task.(*pipeline.BridgeTask); is {
			bridgeName = b.Name
		}

		if trr.Result.Error != nil {
			e.lggr.Warnw(fmt.Sprintf("cannot get bridge response from bridge task, job=%d, id=%s, name=%q", e.job.ID, trr.Task.DotID(), bridgeName), "err", trr.Result.Error, "jobID", e.job.ID, "dotID", trr.Task.DotID(), "bridgeName", bridgeName)
			continue
		}
		bridgeRawResponse, ok := trr.Result.Value.(string)
		if !ok {
			e.lggr.Warnw(fmt.Sprintf("cannot parse bridge response from bridge task, job=%d, id=%s, name=%q: expected string, got: %v (type %T)", e.job.ID, trr.Task.DotID(), bridgeName, trr.Result.Value, trr.Result.Value), "jobID", e.job.ID, "dotID", trr.Task.DotID(), "bridgeName", bridgeName)
			continue
		}
		eaTelem, err := parseEATelemetry([]byte(bridgeRawResponse))
		if err != nil {
			e.lggr.Warnw(fmt.Sprintf("cannot parse EA telemetry, job=%d, id=%s, name=%q", e.job.ID, trr.Task.DotID(), bridgeName), "err", err, "jobID", e.job.ID, "dotID", trr.Task.DotID(), "bridgeName", bridgeName)
			continue
		}
		value := e.getParsedValue(trrs, trr)

		t := &telem.EnhancedEA{
			DataSource:                    eaTelem.DataSource,
			Value:                         value,
			BridgeTaskRunStartedTimestamp: trr.CreatedAt.UnixMilli(),
			BridgeTaskRunEndedTimestamp:   trr.FinishedAt.Time.UnixMilli(),
			ProviderRequestedTimestamp:    eaTelem.ProviderRequestedTimestamp,
			ProviderReceivedTimestamp:     eaTelem.ProviderReceivedTimestamp,
			ProviderDataStreamEstablished: eaTelem.ProviderDataStreamEstablished,
			ProviderIndicatedTime:         eaTelem.ProviderIndicatedTime,
			Feed:                          contract,
			ChainId:                       chainID,
			Observation:                   observation,
			ConfigDigest:                  timestamp.ConfigDigest,
			Round:                         int64(timestamp.Round),
			Epoch:                         int64(timestamp.Epoch),
		}

		bytes, err := proto.Marshal(t)
		if err != nil {
			e.lggr.Warnw("protobuf marshal failed", "err", err)
			continue
		}

		e.monitoringEndpoint.SendLog(bytes)
	}
}

// collectMercuryEnhancedTelemetry checks if enhanced telemetry should be collected, fetches the information needed and
// sends the telemetry
func (e *EnhancedTelemetryService[T]) collectMercuryEnhancedTelemetry(d EnhancedTelemetryMercuryData) {
	if e.monitoringEndpoint == nil {
		return
	}

	// v1 fields
	var bn int64
	var bh string
	var bt uint64
	// v1+v2+v3 fields
	bp := big.NewInt(0)
	//v1+v3 fields
	bid := big.NewInt(0)
	ask := big.NewInt(0)
	// v2+v3 fields
	var mfts, lp, np int64

	switch {
	case d.V1Observation != nil:
		obs := *d.V1Observation
		if obs.CurrentBlockNum.Err == nil {
			bn = obs.CurrentBlockNum.Val
		}
		if obs.CurrentBlockHash.Err == nil {
			bh = common.BytesToHash(obs.CurrentBlockHash.Val).Hex()
		}
		if obs.CurrentBlockTimestamp.Err == nil {
			bt = obs.CurrentBlockTimestamp.Val
		}
		if obs.BenchmarkPrice.Err == nil && obs.BenchmarkPrice.Val != nil {
			bp = obs.BenchmarkPrice.Val
		}
		if obs.Bid.Err == nil && obs.Bid.Val != nil {
			bid = obs.Bid.Val
		}
		if obs.Ask.Err == nil && obs.Ask.Val != nil {
			ask = obs.Ask.Val
		}
	case d.V2Observation != nil:
		obs := *d.V2Observation
		if obs.MaxFinalizedTimestamp.Err == nil {
			mfts = obs.MaxFinalizedTimestamp.Val
		}
		if obs.LinkPrice.Err == nil && obs.LinkPrice.Val != nil {
			lp = obs.LinkPrice.Val.Int64()
		}
		if obs.NativePrice.Err == nil && obs.NativePrice.Val != nil {
			np = obs.NativePrice.Val.Int64()
		}
		if obs.BenchmarkPrice.Err == nil && obs.BenchmarkPrice.Val != nil {
			bp = obs.BenchmarkPrice.Val
		}
	case d.V3Observation != nil:
		obs := *d.V3Observation
		if obs.MaxFinalizedTimestamp.Err == nil {
			mfts = obs.MaxFinalizedTimestamp.Val
		}
		if obs.LinkPrice.Err == nil && obs.LinkPrice.Val != nil {
			lp = obs.LinkPrice.Val.Int64()
		}
		if obs.NativePrice.Err == nil && obs.NativePrice.Val != nil {
			np = obs.NativePrice.Val.Int64()
		}
		if obs.BenchmarkPrice.Err == nil && obs.BenchmarkPrice.Val != nil {
			bp = obs.BenchmarkPrice.Val
		}
		if obs.Bid.Err == nil && obs.Bid.Val != nil {
			bid = obs.Bid.Val
		}
		if obs.Ask.Err == nil && obs.Ask.Val != nil {
			ask = obs.Ask.Val
		}
	}

	for _, trr := range d.TaskRunResults {
		if trr.Task.Type() != pipeline.TaskTypeBridge {
			continue
		}
		bridgeTask := trr.Task.(*pipeline.BridgeTask)
		bridgeName := bridgeTask.Name

		bridgeRawResponse, ok := trr.Result.Value.(string)
		if !ok {
			e.lggr.Warnw(fmt.Sprintf("cannot get bridge response from bridge task, job=%d, id=%s, name=%q, expected string got %T", e.job.ID, trr.Task.DotID(), bridgeName, trr.Result.Value), "jobID", e.job.ID, "dotID", trr.Task.DotID(), "bridgeName", bridgeName)
			continue
		}
		eaTelem, err := parseEATelemetry([]byte(bridgeRawResponse))
		if err != nil {
			e.lggr.Warnw(fmt.Sprintf("cannot parse EA telemetry, job=%d, id=%s, name=%q", e.job.ID, trr.Task.DotID(), bridgeName), "err", err, "jobID", e.job.ID, "dotID", trr.Task.DotID(), "bridgeName", bridgeName)
		}

		assetSymbol := e.getAssetSymbolFromRequestData(bridgeTask.RequestData)

		benchmarkPrice, bidPrice, askPrice := e.getPricesFromResults(trr, d.TaskRunResults, d.FeedVersion)

		t := &telem.EnhancedEAMercury{
			DataSource:                      eaTelem.DataSource,
			DpBenchmarkPrice:                benchmarkPrice,
			DpBid:                           bidPrice,
			DpAsk:                           askPrice,
			CurrentBlockNumber:              bn,
			CurrentBlockHash:                bh,
			CurrentBlockTimestamp:           bt,
			FetchMaxFinalizedTimestamp:      d.FetchMaxFinalizedTimestamp,
			MaxFinalizedTimestamp:           mfts,
			BridgeTaskRunStartedTimestamp:   trr.CreatedAt.UnixMilli(),
			BridgeTaskRunEndedTimestamp:     trr.FinishedAt.Time.UnixMilli(),
			ProviderRequestedTimestamp:      eaTelem.ProviderRequestedTimestamp,
			ProviderReceivedTimestamp:       eaTelem.ProviderReceivedTimestamp,
			ProviderDataStreamEstablished:   eaTelem.ProviderDataStreamEstablished,
			ProviderIndicatedTime:           eaTelem.ProviderIndicatedTime,
			Feed:                            e.job.OCR2OracleSpec.FeedID.Hex(),
			ObservationBenchmarkPrice:       bp.Int64(),
			ObservationBid:                  bid.Int64(),
			ObservationAsk:                  ask.Int64(),
			ObservationBenchmarkPriceString: stringOrEmpty(bp),
			ObservationBidString:            stringOrEmpty(bid),
			ObservationAskString:            stringOrEmpty(ask),
			IsLinkFeed:                      d.IsLinkFeed,
			LinkPrice:                       lp,
			IsNativeFeed:                    d.IsNativeFeed,
			NativePrice:                     np,
			ConfigDigest:                    d.RepTimestamp.ConfigDigest.Hex(),
			Round:                           int64(d.RepTimestamp.Round),
			Epoch:                           int64(d.RepTimestamp.Epoch),
			AssetSymbol:                     assetSymbol,
			Version:                         uint32(d.FeedVersion),
		}

		bytes, err := proto.Marshal(t)
		if err != nil {
			e.lggr.Warnf("protobuf marshal failed %v", err.Error())
			continue
		}

		e.monitoringEndpoint.SendLog(bytes)
	}
}

// getAssetSymbolFromRequestData parses the requestData of the bridge to generate an asset symbol pair
func (e *EnhancedTelemetryService[T]) getAssetSymbolFromRequestData(requestData string) string {
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

// ShouldCollectEnhancedTelemetryMercury checks if enhanced telemetry should be collected and sent
func ShouldCollectEnhancedTelemetryMercury(jb job.Job) bool {
	if jb.Type.String() == pipeline.OffchainReporting2JobType && jb.OCR2OracleSpec != nil {
		return jb.OCR2OracleSpec.CaptureEATelemetry
	}
	return false
}

// getPricesFromResults parses the pipeline.TaskRunResults for pipeline.TaskTypeJSONParse and gets the benchmarkPrice,
// bid and ask. This functions expects the pipeline.TaskRunResults to be correctly ordered
func (e *EnhancedTelemetryService[T]) getPricesFromResults(startTask pipeline.TaskRunResult, allTasks pipeline.TaskRunResults, mercuryVersion mercuryutils.FeedVersion) (float64, float64, float64) {
	var benchmarkPrice, askPrice, bidPrice float64
	var err error
	// We rely on task results to be sorted in the correct order
	benchmarkPriceTask := allTasks.GetNextTaskOf(startTask)
	if benchmarkPriceTask == nil {
		e.lggr.Warnf("cannot parse enhanced EA telemetry benchmark price, task is nil, job %d, id %s", e.job.ID)
		return 0, 0, 0
	}
	if benchmarkPriceTask.Task.Type() == pipeline.TaskTypeJSONParse {
		if benchmarkPriceTask.Result.Error != nil {
			e.lggr.Warnw(fmt.Sprintf("got error for enhanced EA telemetry benchmark price, job %d, id %s: %s", e.job.ID, benchmarkPriceTask.Task.DotID(), benchmarkPriceTask.Result.Error), "err", benchmarkPriceTask.Result.Error)
		} else {
			benchmarkPrice, err = getResultFloat64(benchmarkPriceTask)
			if err != nil {
				e.lggr.Warnw(fmt.Sprintf("cannot parse enhanced EA telemetry benchmark price, job %d, id %s", e.job.ID, benchmarkPriceTask.Task.DotID()), "err", err)
			}
		}
	}

	// mercury version 2 only supports benchmarkPrice
	if mercuryVersion == 2 {
		return benchmarkPrice, 0, 0
	}

	bidTask := allTasks.GetNextTaskOf(*benchmarkPriceTask)
	if bidTask == nil {
		e.lggr.Warnf("cannot parse enhanced EA telemetry bid price, task is nil, job %d, id %s", e.job.ID)
		return benchmarkPrice, 0, 0
	}

	if bidTask != nil && bidTask.Task.Type() == pipeline.TaskTypeJSONParse {
		if bidTask.Result.Error != nil {
			e.lggr.Warnw(fmt.Sprintf("got error for enhanced EA telemetry bid price, job %d, id %s: %s", e.job.ID, bidTask.Task.DotID(), bidTask.Result.Error), "err", bidTask.Result.Error)
		} else {
			bidPrice, err = getResultFloat64(bidTask)
			if err != nil {
				e.lggr.Warnw(fmt.Sprintf("cannot parse enhanced EA telemetry bid price, job %d, id %s", e.job.ID, bidTask.Task.DotID()), "err", err)
			}
		}
	}

	askTask := allTasks.GetNextTaskOf(*bidTask)
	if askTask == nil {
		e.lggr.Warnf("cannot parse enhanced EA telemetry ask price, task is nil, job %d, id %s", e.job.ID)
		return benchmarkPrice, bidPrice, 0
	}
	if askTask != nil && askTask.Task.Type() == pipeline.TaskTypeJSONParse {
		if bidTask.Result.Error != nil {
			e.lggr.Warnw(fmt.Sprintf("got error for enhanced EA telemetry ask price, job %d, id %s: %s", e.job.ID, askTask.Task.DotID(), askTask.Result.Error), "err", askTask.Result.Error)
		} else {
			askPrice, err = getResultFloat64(askTask)
			if err != nil {
				e.lggr.Warnw(fmt.Sprintf("cannot parse enhanced EA telemetry ask price, job %d, id %s", e.job.ID, askTask.Task.DotID()), "err", err)
			}
		}
	}

	return benchmarkPrice, bidPrice, askPrice
}

// MaybeEnqueueEnhancedTelem sends data to the telemetry channel for processing
func MaybeEnqueueEnhancedTelem(jb job.Job, ch chan<- EnhancedTelemetryMercuryData, data EnhancedTelemetryMercuryData) {
	if ShouldCollectEnhancedTelemetryMercury(jb) {
		EnqueueEnhancedTelem[EnhancedTelemetryMercuryData](ch, data)
	}
}

// EnqueueEnhancedTelem sends data to the telemetry channel for processing
func EnqueueEnhancedTelem[T EnhancedTelemetryData | EnhancedTelemetryMercuryData](ch chan<- T, data T) {
	select {
	case ch <- data:
	default:
	}
}

// getResultFloat64 will check the result type and force it to float64 or returns an error if the conversion cannot be made
func getResultFloat64(task *pipeline.TaskRunResult) (float64, error) {
	result, err := utils.ToDecimal(task.Result.Value)
	if err != nil {
		return 0, err
	}
	resultFloat64, _ := result.Float64()
	return resultFloat64, nil
}

func stringOrEmpty(n *big.Int) string {
	if n.Cmp(big.NewInt(0)) == 0 {
		return ""
	}
	return n.String()
}
