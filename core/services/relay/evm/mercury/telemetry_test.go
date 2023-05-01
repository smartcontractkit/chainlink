package mercury

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

const bridgeResponse = `{
			"meta":{
				"adapterName":"data-source-name"
			},
			"timestamps":{
				"providerDataRequestedUnixMs":92233720368547760,
				"providerDataReceivedUnixMs":-92233720368547760,
				"providerDataStreamEstablishedUnixMs":1,
				"providerIndicatedTimeUnixMs":-123456789
			}
		}`

// Order of tasks matter, mercury expects:
// benchmarkPrice in position 0
// bid in position 1
// ask in position 2
// blockNumber in position 3
// blockHash in position 4
// blockTimestamp in position 5
var finalTrrs = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.MultiplyTask{},
		Result: pipeline.Result{
			Value: 111111,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.MultiplyTask{},
		Result: pipeline.Result{
			Value: 222222,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.MultiplyTask{},
		Result: pipeline.Result{
			Value: 333333,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.LookupTask{},
		Result: pipeline.Result{
			Value: int64(123456789),
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.LookupTask{},
		Result: pipeline.Result{
			Value: common.HexToHash("0x123321"),
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.LookupTask{},
		Result: pipeline.Result{
			Value: uint64(123456789),
		},
	}}
var trrs = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			BaseTask:    pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			RequestData: `{"data":{"to":"LINK","from":"USD"}}`,
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds1_benchmark", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: float64(123456.123456),
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(2, "ds2_bid", nil, nil, 2),
		},
		Result: pipeline.Result{
			Value: float64(1234567.1234567),
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(3, "ds3_ask", nil, nil, 3),
		},
		Result: pipeline.Result{
			Value: float64(123456789.1),
		},
	},
}

func TestGetFinalValues(t *testing.T) {
	ds := datasource{}

	benchmarkPrice, bid, ask, blockNr, blockHash, blockTimestamp := getFinalValues(&ds, &finalTrrs)
	require.Equal(t, benchmarkPrice, int64(111111))
	require.Equal(t, bid, int64(222222))
	require.Equal(t, ask, int64(333333))
	require.Equal(t, blockNr, int64(123456789))
	require.Equal(t, blockHash, common.HexToHash("0x123321").Bytes())
	require.Equal(t, blockTimestamp, uint64(123456789))

	benchmarkPrice, bid, ask, blockNr, blockHash, blockTimestamp = getFinalValues(&ds, &pipeline.TaskRunResults{})
	require.Equal(t, benchmarkPrice, int64(0))
	require.Equal(t, bid, int64(0))
	require.Equal(t, ask, int64(0))
	require.Equal(t, blockNr, int64(0))
	require.Nil(t, blockHash)
	require.Equal(t, blockTimestamp, uint64(0))
}

func TestGetPricesFromResults(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zap.WarnLevel)
	ds := datasource{
		lggr: lggr,
	}
	benchmarkPrice, bid, ask := getPricesFromResults(&ds, trrs[0], &trrs)
	require.Equal(t, benchmarkPrice, float64(123456.123456))
	require.Equal(t, bid, float64(1234567.1234567))
	require.Equal(t, ask, float64(123456789.1))

	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	ds = datasource{
		lggr: lggr,
	}
	benchmarkPrice, bid, ask = getPricesFromResults(&ds, trrs[0], &pipeline.TaskRunResults{})
	require.Equal(t, benchmarkPrice, float64(0))
	require.Equal(t, bid, float64(0))
	require.Equal(t, ask, float64(0))
	require.Equal(t, logs.Len(), 1)
	require.Contains(t, logs.All()[0].Message, "cannot parse enhanced EA telemetry")

	tt := trrs[:2]
	getPricesFromResults(&ds, trrs[0], &tt)
	require.Equal(t, logs.Len(), 2)
	require.Contains(t, logs.All()[1].Message, "cannot parse enhanced EA telemetry bid price, task is nil")

	tt = trrs[:3]
	getPricesFromResults(&ds, trrs[0], &tt)
	require.Equal(t, logs.Len(), 3)
	require.Contains(t, logs.All()[2].Message, "cannot parse enhanced EA telemetry ask price, task is nil")

	trrs2 := pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: bridgeResponse,
			},
		},
		pipeline.TaskRunResult{
			Task: &pipeline.JSONParseTask{
				BaseTask: pipeline.NewBaseTask(1, "ds1_benchmark", nil, nil, 1),
			},
			Result: pipeline.Result{
				Value: nil,
			},
		},
		pipeline.TaskRunResult{
			Task: &pipeline.JSONParseTask{
				BaseTask: pipeline.NewBaseTask(2, "ds2_bid", nil, nil, 2),
			},
			Result: pipeline.Result{
				Value: nil,
			},
		},
		pipeline.TaskRunResult{
			Task: &pipeline.JSONParseTask{
				BaseTask: pipeline.NewBaseTask(3, "ds3_ask", nil, nil, 3),
			},
			Result: pipeline.Result{
				Value: nil,
			},
		}}
	benchmarkPrice, bid, ask = getPricesFromResults(&ds, trrs[0], &trrs2)
	require.Equal(t, benchmarkPrice, float64(0))
	require.Equal(t, bid, float64(0))
	require.Equal(t, ask, float64(0))
	require.Equal(t, logs.Len(), 6)
	require.Contains(t, logs.All()[3].Message, "cannot parse enhanced EA telemetry benchmark price")
	require.Contains(t, logs.All()[4].Message, "cannot parse enhanced EA telemetry bid price")
	require.Contains(t, logs.All()[5].Message, "cannot parse enhanced EA telemetry ask price")
}

func TestShouldCollectEnhancedTelemetryMercury(t *testing.T) {
	jb := job.Job{
		Type: job.Type(pipeline.OffchainReporting2JobType),
		OCR2OracleSpec: &job.OCR2OracleSpec{
			CaptureEATelemetry: true,
		},
	}

	require.Equal(t, shouldCollectEnhancedTelemetryMercury(&jb), true)
	jb.OCR2OracleSpec.CaptureEATelemetry = false
	require.Equal(t, shouldCollectEnhancedTelemetryMercury(&jb), false)
	jb.OCR2OracleSpec.CaptureEATelemetry = true
	jb.Type = job.Type(pipeline.CronJobType)
	require.Equal(t, shouldCollectEnhancedTelemetryMercury(&jb), false)
}

func TestGetAssetSymbolFromRequestData(t *testing.T) {
	require.Equal(t, getAssetSymbolFromRequestData(""), "")
	reqData := `{"data":{"to":"LINK","from":"USD"}}`
	require.Equal(t, getAssetSymbolFromRequestData(reqData), "USD/LINK")
}

func TestCollectMercuryEnhancedTelemetry(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryIngressClient(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.EnhancedEAMercury)

	var sentMessage []byte
	ingressClient.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return().Run(func(args mock.Arguments) {
		sentMessage = args[0].(synchronization.TelemPayload).Telemetry
		wg.Done()
	})

	jb := job.Job{
		Type: job.Type(pipeline.OffchainReporting2JobType),
		OCR2OracleSpec: &job.OCR2OracleSpec{
			CaptureEATelemetry: true,
			FeedID:             common.HexToHash("0x111"),
		},
	}

	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	ds := datasource{
		lggr:               lggr,
		jb:                 jb,
		monitoringEndpoint: monitoringEndpoint,
	}

	wg.Add(1)
	collectMercuryEnhancedTelemetry(&ds, finalTrrs, &trrs, ocrtypes.ReportTimestamp{
		ConfigDigest: ocrtypes.ConfigDigest{2},
		Epoch:        11,
		Round:        22,
	})

	expectedTelemetry := telem.EnhancedEAMercury{
		DataSource:                    "data-source-name",
		DpBenchmarkPrice:              123456.123456,
		DpBid:                         1234567.1234567,
		DpAsk:                         123456789.1,
		CurrentBlockNumber:            123456789,
		CurrentBlockHash:              common.HexToHash("0x123321").String(),
		CurrentBlockTimestamp:         123456789,
		BridgeTaskRunStartedTimestamp: trrs[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:   trrs[0].FinishedAt.Time.UnixMilli(),
		ProviderRequestedTimestamp:    92233720368547760,
		ProviderReceivedTimestamp:     -92233720368547760,
		ProviderDataStreamEstablished: 1,
		ProviderIndicatedTime:         -123456789,
		Feed:                          common.HexToHash("0x111").String(),
		ObservationBenchmarkPrice:     111111,
		ObservationBid:                222222,
		ObservationAsk:                333333,
		ConfigDigest:                  "0200000000000000000000000000000000000000000000000000000000000000",
		Round:                         22,
		Epoch:                         11,
		AssetSymbol:                   "USD/LINK",
	}

	expectedMessage, _ := proto.Marshal(&expectedTelemetry)
	wg.Wait()
	require.Equal(t, expectedMessage, sentMessage)

	trrs[0].Result.Value = ""
	wg.Add(1)
	collectMercuryEnhancedTelemetry(&ds, finalTrrs, &trrs, ocrtypes.ReportTimestamp{
		ConfigDigest: ocrtypes.ConfigDigest{2},
		Epoch:        11,
		Round:        22,
	})
	wg.Wait()
	require.Equal(t, logs.Len(), 1)
	require.Contains(t, logs.All()[0].Message, "cannot parse EA telemetry")

	trrs[0].Result.Value = nil
	collectMercuryEnhancedTelemetry(&ds, finalTrrs, &trrs, ocrtypes.ReportTimestamp{
		ConfigDigest: ocrtypes.ConfigDigest{2},
		Epoch:        11,
		Round:        22,
	})
	require.Equal(t, logs.Len(), 2)
	require.Contains(t, logs.All()[1].Message, "cannot get bridge response from bridge task")
}
