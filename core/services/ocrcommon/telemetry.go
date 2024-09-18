package ocrcommon

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	v1types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v2types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	v3types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	v4types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v4"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type EATelemetry struct {
	DataSource                    string
	ProviderRequestedTimestamp    int64
	ProviderReceivedTimestamp     int64
	ProviderDataStreamEstablished int64
	ProviderIndicatedTime         int64

	DpBenchmarkPrice              float64
	DpBid                         float64
	DpAsk                         float64
	BridgeTaskRunStartedTimestamp int64
	BridgeTaskRunEndedTimestamp   int64
	AssetSymbol                   string
	BridgeRequestData             string
}

type EnhancedTelemetryData struct {
	TaskRunResults pipeline.TaskRunResults
	FinalResults   pipeline.FinalResult
	RepTimestamp   ObservationTimestamp
}

type EnhancedTelemetryMercuryData struct {
	V1Observation                *v1types.Observation
	V2Observation                *v2types.Observation
	V3Observation                *v3types.Observation
	V4Observation                *v4types.Observation
	TaskRunResults               pipeline.TaskRunResults
	RepTimestamp                 ocrtypes.ReportTimestamp
	FeedVersion                  mercuryutils.FeedVersion
	FetchMaxFinalizedTimestamp   bool
	IsLinkFeed                   bool
	IsNativeFeed                 bool
	DpInvariantViolationDetected bool
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

func ParseMercuryEATelemetry(lggr logger.Logger, trrs pipeline.TaskRunResults, feedVersion mercuryutils.FeedVersion) (eaTelemetryValues []EATelemetry) {
	for _, trr := range trrs {
		if trr.Task.Type() != pipeline.TaskTypeBridge {
			continue
		}
		bridgeTask := trr.Task.(*pipeline.BridgeTask)
		bridgeName := bridgeTask.Name

		bridgeRawResponse, ok := trr.Result.Value.(string)
		if !ok {
			lggr.Warnw(fmt.Sprintf("cannot get bridge response from bridge task, id=%s, name=%q, expected string got %T", trr.Task.DotID(), bridgeName, trr.Result.Value), "dotID", trr.Task.DotID(), "bridgeName", bridgeName)
			continue
		}
		eaTelem, err := parseEATelemetry([]byte(bridgeRawResponse))
		if err != nil {
			lggr.Warnw(fmt.Sprintf("cannot parse EA telemetry, id=%s, name=%q", trr.Task.DotID(), bridgeName), "err", err, "dotID", trr.Task.DotID(), "bridgeName", bridgeName)
		}
		eaTelem.BridgeRequestData = bridgeTask.RequestData
		eaTelem.DpBenchmarkPrice, eaTelem.DpBid, eaTelem.DpAsk = getPricesFromBridgeTask(lggr, trr, trrs, feedVersion)

		eaTelem.BridgeTaskRunStartedTimestamp = trr.CreatedAt.UnixMilli()
		eaTelem.BridgeTaskRunEndedTimestamp = trr.FinishedAt.Time.UnixMilli()

		parsedBridgeData := parseBridgeRequestData(bridgeTask.RequestData, feedVersion)
		if parsedBridgeData.IsMarketStatus {
			// Only collect telemetry for pricing bridges.
			continue
		}

		eaTelem.AssetSymbol = parsedBridgeData.AssetSymbol

		eaTelemetryValues = append(eaTelemetryValues, eaTelem)
	}
	return
}

// parseEATelemetry attempts to parse the bridge telemetry
func parseEATelemetry(b []byte) (EATelemetry, error) {
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
		return EATelemetry{}, err
	}

	return EATelemetry{
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
	// v1+v2+v3+v4 fields
	bp := big.NewInt(0)
	// v1+v3 fields
	bid := big.NewInt(0)
	ask := big.NewInt(0)
	// v2+v3 fields
	var mfts, lp, np int64
	// v4 fields
	var marketStatus telem.MarketStatus

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
	case d.V4Observation != nil:
		obs := *d.V4Observation
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
		if obs.MarketStatus.Err == nil {
			marketStatus = telem.MarketStatus(obs.MarketStatus.Val)
		}
	}

	eaTelemetryValues := ParseMercuryEATelemetry(logger.Sugared(e.lggr).With("jobID", e.job.ID), d.TaskRunResults, d.FeedVersion)
	for _, eaTelem := range eaTelemetryValues {
		t := &telem.EnhancedEAMercury{
			DataSource:                      eaTelem.DataSource,
			DpBenchmarkPrice:                eaTelem.DpBenchmarkPrice,
			DpBid:                           eaTelem.DpBid,
			DpAsk:                           eaTelem.DpAsk,
			DpInvariantViolationDetected:    d.DpInvariantViolationDetected,
			CurrentBlockNumber:              bn,
			CurrentBlockHash:                bh,
			CurrentBlockTimestamp:           bt,
			FetchMaxFinalizedTimestamp:      d.FetchMaxFinalizedTimestamp,
			MaxFinalizedTimestamp:           mfts,
			BridgeTaskRunStartedTimestamp:   eaTelem.BridgeTaskRunStartedTimestamp,
			BridgeTaskRunEndedTimestamp:     eaTelem.BridgeTaskRunEndedTimestamp,
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
			ObservationMarketStatus:         marketStatus,
			IsLinkFeed:                      d.IsLinkFeed,
			LinkPrice:                       lp,
			IsNativeFeed:                    d.IsNativeFeed,
			NativePrice:                     np,
			ConfigDigest:                    d.RepTimestamp.ConfigDigest.Hex(),
			Round:                           int64(d.RepTimestamp.Round),
			Epoch:                           int64(d.RepTimestamp.Epoch),
			BridgeRequestData:               eaTelem.BridgeRequestData,
			AssetSymbol:                     eaTelem.AssetSymbol,
			Version:                         uint32(d.FeedVersion),
		}
		e.lggr.Debugw(fmt.Sprintf("EA Telemetry = %+v", t), "feedID", e.job.OCR2OracleSpec.FeedID.Hex(), "jobID", e.job.ID, "datasource", eaTelem.DataSource)
		bytes, err := proto.Marshal(t)
		if err != nil {
			e.lggr.Warnf("protobuf marshal failed %v", err.Error())
			continue
		}

		e.monitoringEndpoint.SendLog(bytes)
	}
}

type telemetryAttributes struct {
	PriceType *string `json:"priceType"`
}

func parseTelemetryAttributes(a string) (telemetryAttributes, error) {
	attrs := &telemetryAttributes{}
	err := json.Unmarshal([]byte(a), attrs)
	if err != nil {
		return telemetryAttributes{}, err
	}
	return *attrs, nil
}

type bridgeRequestData struct {
	AssetSymbol    string
	IsMarketStatus bool
}

// parseRequestData parses the requestData of the bridge.
func parseBridgeRequestData(requestData string, mercuryVersion mercuryutils.FeedVersion) bridgeRequestData {
	type reqDataPayload struct {
		Endpoint *string `json:"endpoint"`
		To       *string `json:"to"`
		From     *string `json:"from"`
		Address  *string `json:"address"` // used for view function ea only
		Market   *string `json:"market"`  // used for market status ea only
	}
	type reqData struct {
		Data reqDataPayload `json:"data"`
	}

	rd := &reqData{}
	err := json.Unmarshal([]byte(requestData), rd)
	if err != nil {
		return bridgeRequestData{}
	}

	if mercuryVersion == 4 && ((rd.Data.Endpoint != nil && *rd.Data.Endpoint == "market-status") || (rd.Data.Market != nil && *rd.Data.Market != "")) {
		return bridgeRequestData{
			AssetSymbol:    *rd.Data.Market,
			IsMarketStatus: true,
		}
	}

	if rd.Data.From != nil && rd.Data.To != nil {
		return bridgeRequestData{AssetSymbol: *rd.Data.From + "/" + *rd.Data.To}
	}

	if rd.Data.Address != nil {
		return bridgeRequestData{AssetSymbol: *rd.Data.Address}
	}

	return bridgeRequestData{}
}

// ShouldCollectEnhancedTelemetryMercury checks if enhanced telemetry should be collected and sent
func ShouldCollectEnhancedTelemetryMercury(jb job.Job) bool {
	if jb.Type.String() == pipeline.OffchainReporting2JobType && jb.OCR2OracleSpec != nil {
		return jb.OCR2OracleSpec.CaptureEATelemetry
	}
	return false
}

const (
	bid          = "bid"
	ask          = "ask"
	benchmark    = "benchmark"
	exchangeRate = "exchangeRate"
)

func getPricesFromBridgeTask(lggr logger.Logger, bridgeTask pipeline.TaskRunResult, allTasks pipeline.TaskRunResults, mercuryVersion mercuryutils.FeedVersion) (float64, float64, float64) {
	var benchmarkPrice, bidPrice, askPrice float64

	// This will assume that all fields we care about are tagged with the correct priceType
	benchmarkPrice, bidPrice, askPrice = getPricesFromBridgeTaskByTelemetryField(lggr, bridgeTask, allTasks)

	// If prices weren't parsed by telemetry fields - attempt to get prices using the legacy method
	// This is for backwards compatibility with job specs that don't have the telemetry attributes set
	if benchmarkPrice == 0 && bidPrice == 0 && askPrice == 0 {
		benchmarkP, bidP, askP := getPricesFromResultsByOrder(lggr, bridgeTask, allTasks, mercuryVersion)
		bidPrice = bidP
		askPrice = askP
		benchmarkPrice = benchmarkP
	}

	return benchmarkPrice, bidPrice, askPrice
}

// CollectTaskRunResultsWithTags collects TaskRunResults for descendent tasks with non-empty TaskTags.
func collectTaskRunResultsWithTags(bridgeTask pipeline.TaskRunResult, allTasks pipeline.TaskRunResults) []pipeline.TaskRunResult {
	startTask := bridgeTask.Task
	descendants := startTask.GetDescendantTasks()
	var taskRunResultsWithTags []pipeline.TaskRunResult
	for _, task := range descendants {
		trr := allTasks.GetTaskRunResultOf(task)
		if trr != nil {
			if trr.Task.TaskTags() != "" {
				taskRunResultsWithTags = append(taskRunResultsWithTags, *trr)
			}
		}
	}
	return taskRunResultsWithTags
}

// getPricesFromBridgeTaskByTelemetryField attempts to parse prices from via telemetry fields in the TaskTags
func getPricesFromBridgeTaskByTelemetryField(lggr logger.Logger, bridgeTask pipeline.TaskRunResult, allTasks pipeline.TaskRunResults) (float64, float64, float64) {
	var benchmarkPrice, bidPrice, askPrice float64

	// Outputs are the mapped tasks from this task.
	var tasksWithTags = collectTaskRunResultsWithTags(bridgeTask, allTasks)

	for _, trr := range tasksWithTags {
		attributes, err := parseTelemetryAttributes(trr.Task.TaskTags())
		if err != nil {
			lggr.Warnw(fmt.Sprintf("found telemetry attributes but cannot them, taskTags=%s", trr.Task.TaskTags()), "err", err)
			continue
		}

		if attributes.PriceType != nil {
			switch *attributes.PriceType {
			case bid:
				bidPrice = parsePriceFromTask(lggr, trr)
			case ask:
				askPrice = parsePriceFromTask(lggr, trr)
			case benchmark:
				benchmarkPrice = parsePriceFromTask(lggr, trr)
			case exchangeRate:
				price := parsePriceFromTask(lggr, trr)
				benchmarkPrice, bidPrice, askPrice = price, price, price
			case "":
				lggr.Warnw(fmt.Sprintf("no priceType found in attributes, parsedAttributes=%+v, id %s", attributes, trr.Task.DotID()))
			}
		}
	}

	return benchmarkPrice, bidPrice, askPrice
}

func parsePriceFromTask(lggr logger.Logger, trr pipeline.TaskRunResult) float64 {
	var val float64
	if trr.Result.Error != nil {
		lggr.Warnw(fmt.Sprintf("got error on EA telemetry price task, id %s: %s", trr.Task.DotID(), trr.Result.Error), "err", trr.Result.Error)
		return 0
	}
	val, err := getResultFloat64(&trr)
	if err != nil {
		lggr.Warnw(fmt.Sprintf("cannot parse EA telemetry price to float64, DOT id %s", trr.Task.DotID()), "task_type", trr.Task.Type(), "task_tags", trr.Task.TaskTags(), "err", err)
	}
	return val
}

// getPricesFromResultsByOrder parses the pipeline.TaskRunResults for pipeline.TaskTypeJSONParse and gets the benchmarkPrice,
// bid and ask. This functions expects the pipeline.TaskRunResults to be correctly ordered
func getPricesFromResultsByOrder(lggr logger.Logger, startTask pipeline.TaskRunResult, allTasks pipeline.TaskRunResults, mercuryVersion mercuryutils.FeedVersion) (float64, float64, float64) {
	var benchmarkPrice, askPrice, bidPrice float64

	// We rely on task results to be sorted in the correct order
	benchmarkPriceTask := allTasks.GetNextTaskOf(startTask)
	if benchmarkPriceTask == nil {
		lggr.Warn("cannot parse enhanced EA telemetry benchmark price, task is nil")
		return 0, 0, 0
	}
	if benchmarkPriceTask.Task.Type() == pipeline.TaskTypeJSONParse {
		benchmarkPrice = parsePriceFromTask(lggr, *benchmarkPriceTask)
	}

	// mercury versions 2 and 4 only supports benchmarkPrice
	if mercuryVersion == 2 || mercuryVersion == 4 {
		return benchmarkPrice, 0, 0
	}

	bidTask := allTasks.GetNextTaskOf(*benchmarkPriceTask)
	if bidTask == nil {
		lggr.Warnf("cannot parse enhanced EA telemetry bid price, task is nil, id %s", benchmarkPriceTask.Task.DotID())
		return benchmarkPrice, 0, 0
	}

	if bidTask.Task.Type() == pipeline.TaskTypeJSONParse {
		bidPrice = parsePriceFromTask(lggr, *bidTask)
	}

	askTask := allTasks.GetNextTaskOf(*bidTask)
	if askTask == nil {
		lggr.Warnf("cannot parse enhanced EA telemetry ask price, task is nil, id %s", benchmarkPriceTask.Task.DotID())
		return benchmarkPrice, bidPrice, 0
	}
	if askTask.Task.Type() == pipeline.TaskTypeJSONParse {
		askPrice = parsePriceFromTask(lggr, *askTask)
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
