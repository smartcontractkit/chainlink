package ocrcommon

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercuryv1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercuryv2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercuryv4 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v4"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
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

var trrs = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:     "test-bridge-1",
			BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds1_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "123456.123456789",
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:     "test-bridge-2",
			BaseTask: pipeline.NewBaseTask(0, "ds2", nil, nil, 0),
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds2_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "12345678",
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:     "test-bridge-3",
			BaseTask: pipeline.NewBaseTask(0, "ds3", nil, nil, 0),
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds3_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "1234567890",
		},
	},
}

func TestShouldCollectTelemetry(t *testing.T) {
	j := job.Job{
		OCROracleSpec:  &job.OCROracleSpec{CaptureEATelemetry: true},
		OCR2OracleSpec: &job.OCR2OracleSpec{CaptureEATelemetry: true},
	}

	j.Type = job.Type(pipeline.OffchainReportingJobType)
	assert.True(t, ShouldCollectEnhancedTelemetry(&j))
	j.OCROracleSpec.CaptureEATelemetry = false
	assert.False(t, ShouldCollectEnhancedTelemetry(&j))

	j.Type = job.Type(pipeline.OffchainReporting2JobType)
	assert.True(t, ShouldCollectEnhancedTelemetry(&j))
	j.OCR2OracleSpec.CaptureEATelemetry = false
	assert.False(t, ShouldCollectEnhancedTelemetry(&j))

	j.Type = job.Type(pipeline.VRFJobType)
	assert.False(t, ShouldCollectEnhancedTelemetry(&j))
}

func TestGetContract(t *testing.T) {
	j := job.Job{
		OCROracleSpec:  &job.OCROracleSpec{CaptureEATelemetry: true},
		OCR2OracleSpec: &job.OCR2OracleSpec{CaptureEATelemetry: true},
	}
	e := EnhancedTelemetryService[EnhancedTelemetryData]{
		job:  &j,
		lggr: nil,
	}
	contractAddress := evmtypes.EIP55Address(utils.RandomAddress().String())

	j.Type = job.Type(pipeline.OffchainReportingJobType)
	j.OCROracleSpec.ContractAddress = contractAddress
	assert.Equal(t, contractAddress.String(), e.getContract())

	j.Type = job.Type(pipeline.OffchainReporting2JobType)
	j.OCR2OracleSpec.ContractID = contractAddress.String()
	assert.Equal(t, contractAddress.String(), e.getContract())

	j.Type = job.Type(pipeline.VRFJobType)
	assert.Empty(t, e.getContract())
}

func TestGetChainID(t *testing.T) {
	j := job.Job{
		OCROracleSpec:  &job.OCROracleSpec{CaptureEATelemetry: true},
		OCR2OracleSpec: &job.OCR2OracleSpec{CaptureEATelemetry: true},
	}
	e := EnhancedTelemetryService[EnhancedTelemetryData]{
		job:  &j,
		lggr: nil,
	}

	j.Type = job.Type(pipeline.OffchainReportingJobType)
	j.OCROracleSpec.EVMChainID = (*ubig.Big)(big.NewInt(1234567890))
	assert.Equal(t, "1234567890", e.getChainID())

	j.Type = job.Type(pipeline.OffchainReporting2JobType)
	j.OCR2OracleSpec.RelayConfig = make(map[string]interface{})
	j.OCR2OracleSpec.RelayConfig["chainID"] = "foo"
	assert.Equal(t, "foo", e.getChainID())

	j.Type = job.Type(pipeline.VRFJobType)
	assert.Empty(t, e.getChainID())
}

func TestParseEATelemetry(t *testing.T) {
	ea, err := parseEATelemetry([]byte(bridgeResponse))
	assert.NoError(t, err)
	assert.Equal(t, ea.DataSource, "data-source-name")
	assert.Equal(t, ea.ProviderRequestedTimestamp, int64(92233720368547760))
	assert.Equal(t, ea.ProviderReceivedTimestamp, int64(-92233720368547760))
	assert.Equal(t, ea.ProviderDataStreamEstablished, int64(1))
	assert.Equal(t, ea.ProviderIndicatedTime, int64(-123456789))

	_, err = parseEATelemetry(nil)
	assert.Error(t, err)
}

func TestGetJsonParsedValue(t *testing.T) {
	resp := getJsonParsedValue(trrs[0], &trrs)
	assert.Equal(t, 123456.123456789, *resp)

	trrs[1].Result.Value = nil
	resp = getJsonParsedValue(trrs[0], &trrs)
	assert.Nil(t, resp)

	resp = getJsonParsedValue(trrs[1], &trrs)
	assert.Nil(t, resp)
}

func TestSendEATelemetry(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("test-network", "test-chainID", "0xa", synchronization.EnhancedEA)

	var sentMessage []byte
	ingressClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		sentMessage = args[1].([]byte)
		wg.Done()
	})

	feedAddress := utils.RandomAddress()

	enhancedTelemChan := make(chan EnhancedTelemetryData, 100)
	jb := job.Job{
		Type: job.Type(pipeline.OffchainReportingJobType),
		OCROracleSpec: &job.OCROracleSpec{
			ContractAddress:    evmtypes.EIP55AddressFromAddress(feedAddress),
			CaptureEATelemetry: true,
			EVMChainID:         (*ubig.Big)(big.NewInt(9)),
		},
	}

	lggr, _ := logger.TestLoggerObserved(t, zap.WarnLevel)
	doneCh := make(chan struct{})
	enhancedTelemService := NewEnhancedTelemetryService(&jb, enhancedTelemChan, doneCh, monitoringEndpoint, lggr.Named("Enhanced Telemetry Mercury"))
	servicetest.Run(t, enhancedTelemService)
	trrs := pipeline.TaskRunResults{
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
				BaseTask: pipeline.NewBaseTask(1, "ds1", nil, nil, 1),
			},
			Result: pipeline.Result{
				Value: "123456789.1234567",
			},
		},
	}
	fr := pipeline.FinalResult{
		Values:      []interface{}{"123456"},
		AllErrors:   nil,
		FatalErrors: []error{nil},
	}

	observationTimestamp := ObservationTimestamp{
		Round:        15,
		Epoch:        738,
		ConfigDigest: "config digest hex",
	}

	wg.Add(1)
	enhancedTelemChan <- EnhancedTelemetryData{
		TaskRunResults: trrs,
		FinalResults:   fr,
		RepTimestamp:   observationTimestamp,
	}

	expectedTelemetry := telem.EnhancedEA{
		DataSource:                    "data-source-name",
		Value:                         123456789.1234567,
		BridgeTaskRunStartedTimestamp: trrs[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:   trrs[0].FinishedAt.Time.UnixMilli(),
		ProviderRequestedTimestamp:    92233720368547760,
		ProviderReceivedTimestamp:     -92233720368547760,
		ProviderDataStreamEstablished: 1,
		ProviderIndicatedTime:         -123456789,
		Feed:                          feedAddress.String(),
		ChainId:                       "9",
		Observation:                   123456,
		Round:                         15,
		Epoch:                         738,
		ConfigDigest:                  "config digest hex",
	}

	expectedMessage, _ := proto.Marshal(&expectedTelemetry)
	wg.Wait()
	assert.Equal(t, expectedMessage, sentMessage)
	doneCh <- struct{}{}
}

func TestGetObservation(t *testing.T) {
	j := job.Job{
		OCROracleSpec:  &job.OCROracleSpec{CaptureEATelemetry: true},
		OCR2OracleSpec: &job.OCR2OracleSpec{CaptureEATelemetry: true},
	}

	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	e := EnhancedTelemetryService[EnhancedTelemetryData]{
		job:  &j,
		lggr: lggr,
	}

	obs := e.getObservation(&pipeline.FinalResult{})
	assert.Equal(t, obs, int64(0))
	assert.Equal(t, logs.Len(), 1)
	assert.Contains(t, logs.All()[0].Message, "cannot get singular result")

	finalResult := &pipeline.FinalResult{
		Values:      []interface{}{"123456"},
		AllErrors:   nil,
		FatalErrors: []error{nil},
	}
	obs = e.getObservation(finalResult)
	assert.Equal(t, obs, int64(123456))
}

func TestCollectAndSend(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("test-network", "test-chainID", "0xa", synchronization.EnhancedEA)
	ingressClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		wg.Done()
	})

	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	jb := job.Job{
		ID:   1234567890,
		Type: job.Type(pipeline.OffchainReportingJobType),
		OCROracleSpec: &job.OCROracleSpec{
			CaptureEATelemetry: true,
		},
	}

	enhancedTelemChan := make(chan EnhancedTelemetryData, 100)
	doneCh := make(chan struct{})

	enhancedTelemService := NewEnhancedTelemetryService(&jb, enhancedTelemChan, doneCh, monitoringEndpoint, lggr.Named("Enhanced Telemetry"))
	servicetest.Run(t, enhancedTelemService)
	finalResult := &pipeline.FinalResult{
		Values:      []interface{}{"123456"},
		AllErrors:   nil,
		FatalErrors: []error{nil},
	}

	badTrrs := &pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: nil,
			},
		},
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds2", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: bridgeResponse,
			},
		}}

	observationTimestamp := ObservationTimestamp{
		Round:        0,
		Epoch:        0,
		ConfigDigest: "",
	}
	wg.Add(1)
	enhancedTelemChan <- EnhancedTelemetryData{
		TaskRunResults: *badTrrs,
		FinalResults:   *finalResult,
		RepTimestamp:   observationTimestamp,
	}

	wg.Wait()
	assert.Equal(t, logs.Len(), 2)
	assert.Contains(t, logs.All()[0].Message, "cannot parse bridge response from bridge task")

	badTrrs = &pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: "[]",
			},
		}}
	enhancedTelemChan <- EnhancedTelemetryData{
		TaskRunResults: *badTrrs,
		FinalResults:   *finalResult,
		RepTimestamp:   observationTimestamp,
	}
	wg.Wait()
	assert.Equal(t, 2, logs.Len())
	assert.Contains(t, logs.All()[0].Message, "cannot parse bridge response from bridge task")
	assert.Contains(t, logs.All()[1].Message, "cannot get json parse value")
	doneCh <- struct{}{}
}

var trrsMercuryV1 = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:        "link-usd-test-bridge-v1",
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
			Value: int64(321123),
		},
	},
}

var trrsMercuryV2 = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:        "link-usd-test-bridge-v2",
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
}

var trrsMercuryV4 = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:        "link-usd-test-bridge-v2",
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
			Value: 123456.123456,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:        "market-status-bridge",
			BaseTask:    pipeline.NewBaseTask(2, "ds2", nil, nil, 2),
			RequestData: `{"data":{"endpoint":"market-status","market":"forex"}}`,
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(3, "market_status", nil, nil, 3),
		},
		Result: pipeline.Result{
			Value: 2.0,
		},
	},
}

func TestGetPricesFromBridgeByTelemetryField(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zap.WarnLevel)
	// These are intentionally out of order from the "legacy" method which expects order of `benchmark, bid, ask`
	jsonParseTaskBid := pipeline.JSONParseTask{
		BaseTask: pipeline.NewBaseTask(1, "json_parse_2", nil, nil, 2),
	}
	jsonParseTaskBid.BaseTask.Tags = `{"priceType": "bid"}`
	jsonParseTaskAsk := pipeline.JSONParseTask{
		BaseTask: pipeline.NewBaseTask(2, "json_parse_3", nil, nil, 3),
	}
	jsonParseTaskAsk.BaseTask.Tags = `{"priceType": "ask"}`
	jsonParseTaskBenchmark := pipeline.JSONParseTask{
		BaseTask: pipeline.NewBaseTask(3, "json_parse_1", nil, nil, 1),
	}
	jsonParseTaskBenchmark.BaseTask.Tags = `{"priceType": "benchmark"}`

	bridgeOutputs := []pipeline.Task{&jsonParseTaskAsk, &jsonParseTaskBid, &jsonParseTaskBenchmark}

	bridgeTask := pipeline.BridgeTask{
		Name:     "bridge-task",
		BaseTask: pipeline.NewBaseTask(0, "bridge", nil, bridgeOutputs, 0),
	}

	// Create task run results
	taskRunResults := pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &bridgeTask,
			Result: pipeline.Result{
				Value: bridgeResponse,
			},
		},
		pipeline.TaskRunResult{
			Task: &jsonParseTaskBenchmark,
			Result: pipeline.Result{
				Value: "123456.123456",
			},
		},
		pipeline.TaskRunResult{
			Task: &jsonParseTaskBid,
			Result: pipeline.Result{
				Value: "1234567.1234567",
			},
		},
		pipeline.TaskRunResult{
			Task: &jsonParseTaskAsk,
			Result: pipeline.Result{
				Value: "321123",
			},
		},
	}

	benchmarkPrice, bidPrice, askPrice := getPricesFromBridgeTask(lggr, taskRunResults[0], taskRunResults, 1)

	require.Equal(t, 123456.123456, benchmarkPrice)
	require.Equal(t, 1234567.1234567, bidPrice)
	require.Equal(t, 321123.0, askPrice)

	// now removing the TaskTags will throw off the parsed order - and we'll be parsing the "incorrect" prices
	// according to the legacy ordering approach
	jsonParseTaskAsk.BaseTask.Tags = ""
	jsonParseTaskBid.BaseTask.Tags = ""
	jsonParseTaskBenchmark.BaseTask.Tags = ""

	wrongBenchmarkPrice, wrongBidPrice, wrongAskPrice := getPricesFromBridgeTask(lggr, taskRunResults[0], taskRunResults, 1)
	require.Equal(t, 1234567.1234567, wrongBenchmarkPrice)
	require.Equal(t, 321123.0, wrongBidPrice)
	require.Equal(t, 123456.123456, wrongAskPrice)
}

func TestGetPricesFromBridgeTaskByOrder(t *testing.T) {
	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)

	benchmarkPrice, bid, ask := getPricesFromBridgeTask(lggr, trrsMercuryV1[0], trrsMercuryV1, 1)
	require.Equal(t, 123456.123456, benchmarkPrice)
	require.Equal(t, 1234567.1234567, bid)
	require.Equal(t, float64(321123), ask)

	benchmarkPrice, bid, ask = getPricesFromBridgeTask(lggr, trrsMercuryV1[0], pipeline.TaskRunResults{}, 1)
	require.Equal(t, float64(0), benchmarkPrice)
	require.Equal(t, float64(0), bid)
	require.Equal(t, float64(0), ask)
	require.Equal(t, 1, logs.Len())
	require.Contains(t, logs.All()[0].Message, "cannot parse enhanced EA telemetry")

	tt := trrsMercuryV1[:2]
	getPricesFromBridgeTask(lggr, trrsMercuryV1[0], tt, 1)
	require.Equal(t, 2, logs.Len())
	require.Contains(t, logs.All()[1].Message, "cannot parse enhanced EA telemetry bid price, task is nil")

	tt = trrsMercuryV1[:3]
	getPricesFromBridgeTask(lggr, trrsMercuryV1[0], tt, 1)
	require.Equal(t, 3, logs.Len())
	require.Contains(t, logs.All()[2].Message, "cannot parse enhanced EA telemetry ask price, task is nil")

	trrs2 := pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				Name:     "test-bridge-1",
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
	benchmarkPrice, bid, ask = getPricesFromBridgeTask(lggr, trrsMercuryV1[0], trrs2, 3)
	require.Equal(t, benchmarkPrice, float64(0))
	require.Equal(t, bid, float64(0))
	require.Equal(t, ask, float64(0))
	require.Equal(t, logs.Len(), 6)
	require.Contains(t, logs.All()[3].Message, "cannot parse EA telemetry price to float64, DOT id ds1_benchmark")
	require.Contains(t, logs.All()[4].Message, "cannot parse EA telemetry price to float64, DOT id ds2_bid")
	require.Contains(t, logs.All()[5].Message, "cannot parse EA telemetry price to float64, DOT id ds3_ask")

	benchmarkPrice, bid, ask = getPricesFromBridgeTask(lggr, trrsMercuryV1[0], trrsMercuryV2, 2)
	require.Equal(t, 123456.123456, benchmarkPrice)
	require.Equal(t, float64(0), bid)
	require.Equal(t, float64(0), ask)
}

func TestShouldCollectEnhancedTelemetryMercury(t *testing.T) {
	j := job.Job{
		Type: job.Type(pipeline.OffchainReporting2JobType),
		OCR2OracleSpec: &job.OCR2OracleSpec{
			CaptureEATelemetry: true,
		},
	}

	require.Equal(t, ShouldCollectEnhancedTelemetryMercury(j), true)
	j.OCR2OracleSpec.CaptureEATelemetry = false
	require.Equal(t, ShouldCollectEnhancedTelemetryMercury(j), false)
	j.OCR2OracleSpec.CaptureEATelemetry = true
	j.Type = job.Type(pipeline.CronJobType)
	require.Equal(t, ShouldCollectEnhancedTelemetryMercury(j), false)
}

func TestParseBridgeRequestData(t *testing.T) {
	require.Equal(t, parseBridgeRequestData("", 2), bridgeRequestData{})

	reqData := `{"data":{"to":"LINK","from":"USD"}}`
	require.Equal(t, parseBridgeRequestData(reqData, 2), bridgeRequestData{AssetSymbol: "USD/LINK"})

	reqData = `{"data":{"to":"LINK","from":"USD","market":"forex"}}`
	require.Equal(t, parseBridgeRequestData(reqData, 2), bridgeRequestData{AssetSymbol: "USD/LINK"})

	reqData = `{"data":{"endpoint":"market-status","market":"forex"}}`
	require.Equal(t, parseBridgeRequestData(reqData, 4), bridgeRequestData{AssetSymbol: "forex", IsMarketStatus: true})

	reqData = `{"data":{"market":"metals"}}`
	require.Equal(t, parseBridgeRequestData(reqData, 4), bridgeRequestData{AssetSymbol: "metals", IsMarketStatus: true})

	viewFunctionReqData := `{"data":{"address":"0x12345678", "signature": "function stEthPerToken() view returns (int256)"}}`
	require.Equal(t, parseBridgeRequestData(viewFunctionReqData, 3), bridgeRequestData{AssetSymbol: "0x12345678"})
}

func getViewFunctionTaskRunResults() pipeline.TaskRunResults {
	var taskViewFunctionParseValue = func() pipeline.MultiplyTask {
		task := pipeline.MultiplyTask{
			BaseTask: pipeline.NewBaseTask(3, "ds1_parse", nil, nil, 3),
			Times:    "1",
		}
		task.BaseTask.Tags = `{"priceType": "exchangeRate"}`
		return task
	}()

	var taskViewFunctionDecode = pipeline.ETHABIDecodeTask{
		ABI:      "uint256 data",
		BaseTask: pipeline.NewBaseTask(2, "ds1_decode", nil, []pipeline.Task{&taskViewFunctionParseValue}, 2),
	}

	var taskViewFunctionJSONParse = pipeline.JSONParseTask{
		BaseTask: pipeline.NewBaseTask(1, "ds1_parse", nil, []pipeline.Task{&taskViewFunctionDecode}, 1),
	}

	const viewFunctionBridgeResponse = `{
		  "data": {
			"result": "0x000000000000000000000000000000000000000000000000105ba6a589b23a81"
		  },
		  "statusCode": 200,
		  "result": "0x000000000000000000000000000000000000000000000000105ba6a589b23a81",
		  "timestamps": {
			"providerDataRequestedUnixMs": 1726243598046,
			"providerDataReceivedUnixMs": 1726243598341
		  },
		  "meta": {
			"adapterName": "VIEW_FUNCTION"
		  }
		}`

	var taskViewFunctionBridgeRequest = pipeline.BridgeTask{
		Name:        "bridge-view-function",
		BaseTask:    pipeline.NewBaseTask(0, "ds1", nil, []pipeline.Task{&taskViewFunctionJSONParse}, 0),
		RequestData: `{"data":{"address":"0x1234","signature":"function stEthPerToken() external view returns (uint256)"}}`,
	}

	return pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &taskViewFunctionBridgeRequest,
			Result: pipeline.Result{
				Value: viewFunctionBridgeResponse,
			},
		},
		pipeline.TaskRunResult{
			Task: &taskViewFunctionJSONParse,
			Result: pipeline.Result{
				Value: `0x000000000000000000000000000000000000000000000000105ba6a589b23a81`,
			},
		},
		pipeline.TaskRunResult{
			Task: &taskViewFunctionDecode,
			Result: pipeline.Result{
				Value: map[string]interface{}{
					"data": big.NewInt(1178718957397490305),
				},
			},
		},
		pipeline.TaskRunResult{
			Task: &taskViewFunctionParseValue,
			Result: pipeline.Result{
				Value: decimal.NewFromInt(1178718957397490305),
			},
		},
	}
}

func TestCollectMercuryEnhancedTelemetryV1ViewFunction(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("test-network", "test-chainID", "0xa", synchronization.EnhancedEAMercury)

	var sentMessage []byte
	ingressClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		sentMessage = args[1].([]byte)
		wg.Done()
	})

	lggr, _ := logger.TestLoggerObserved(t, zap.WarnLevel)
	chTelem := make(chan EnhancedTelemetryMercuryData, 100)
	chDone := make(chan struct{})
	feedID := common.HexToHash("0x111")
	e := EnhancedTelemetryService[EnhancedTelemetryMercuryData]{
		chDone:  chDone,
		chTelem: chTelem,
		job: &job.Job{
			Type: job.Type(pipeline.OffchainReporting2JobType),
			OCR2OracleSpec: &job.OCR2OracleSpec{
				CaptureEATelemetry: true,
				FeedID:             &feedID,
			},
		},
		lggr:               lggr,
		monitoringEndpoint: monitoringEndpoint,
	}
	servicetest.Run(t, &e)

	wg.Add(1)

	taskRunResults := getViewFunctionTaskRunResults()

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: taskRunResults,
		V1Observation: &mercuryv1.Observation{
			BenchmarkPrice:        mercury.ObsResult[*big.Int]{Val: big.NewInt(111111)},
			Bid:                   mercury.ObsResult[*big.Int]{Val: big.NewInt(222222)},
			Ask:                   mercury.ObsResult[*big.Int]{Val: big.NewInt(333333)},
			CurrentBlockNum:       mercury.ObsResult[int64]{Val: 123456789},
			CurrentBlockHash:      mercury.ObsResult[[]byte]{Val: common.HexToHash("0x123321").Bytes()},
			CurrentBlockTimestamp: mercury.ObsResult[uint64]{Val: 987654321},
		},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}

	expectedTelemetry := telem.EnhancedEAMercury{
		DataSource:                      "VIEW_FUNCTION",
		DpBenchmarkPrice:                1178718957397490400,
		DpBid:                           1178718957397490400,
		DpAsk:                           1178718957397490400,
		CurrentBlockNumber:              123456789,
		CurrentBlockHash:                common.HexToHash("0x123321").String(),
		CurrentBlockTimestamp:           987654321,
		BridgeTaskRunStartedTimestamp:   taskRunResults[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:     taskRunResults[0].FinishedAt.Time.UnixMilli(),
		ProviderRequestedTimestamp:      1726243598046,
		ProviderReceivedTimestamp:       1726243598341,
		ProviderDataStreamEstablished:   0,
		ProviderIndicatedTime:           0,
		Feed:                            common.HexToHash("0x111").String(),
		ObservationBenchmarkPrice:       111111,
		ObservationBid:                  222222,
		ObservationAsk:                  333333,
		ConfigDigest:                    "0200000000000000000000000000000000000000000000000000000000000000",
		Round:                           22,
		Epoch:                           11,
		BridgeRequestData:               `{"data":{"address":"0x1234","signature":"function stEthPerToken() external view returns (uint256)"}}`,
		AssetSymbol:                     "0x1234",
		ObservationBenchmarkPriceString: "111111",
		ObservationBidString:            "222222",
		ObservationAskString:            "333333",
	}

	expectedMessage, _ := proto.Marshal(&expectedTelemetry)
	wg.Wait()
	require.Equal(t, expectedMessage, sentMessage)

	chDone <- struct{}{}
}

func TestCollectMercuryEnhancedTelemetryV1(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("test-network", "test-chainID", "0xa", synchronization.EnhancedEAMercury)

	var sentMessage []byte
	ingressClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		sentMessage = args[1].([]byte)
		wg.Done()
	})

	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	chTelem := make(chan EnhancedTelemetryMercuryData, 100)
	chDone := make(chan struct{})
	feedID := common.HexToHash("0x111")
	e := EnhancedTelemetryService[EnhancedTelemetryMercuryData]{
		chDone:  chDone,
		chTelem: chTelem,
		job: &job.Job{
			Type: job.Type(pipeline.OffchainReporting2JobType),
			OCR2OracleSpec: &job.OCR2OracleSpec{
				CaptureEATelemetry: true,
				FeedID:             &feedID,
			},
		},
		lggr:               lggr,
		monitoringEndpoint: monitoringEndpoint,
	}
	servicetest.Run(t, &e)

	wg.Add(1)

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: trrsMercuryV1,
		V1Observation: &mercuryv1.Observation{
			BenchmarkPrice:        mercury.ObsResult[*big.Int]{Val: big.NewInt(111111)},
			Bid:                   mercury.ObsResult[*big.Int]{Val: big.NewInt(222222)},
			Ask:                   mercury.ObsResult[*big.Int]{Val: big.NewInt(333333)},
			CurrentBlockNum:       mercury.ObsResult[int64]{Val: 123456789},
			CurrentBlockHash:      mercury.ObsResult[[]byte]{Val: common.HexToHash("0x123321").Bytes()},
			CurrentBlockTimestamp: mercury.ObsResult[uint64]{Val: 987654321},
		},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}

	expectedTelemetry := telem.EnhancedEAMercury{
		DataSource:                      "data-source-name",
		DpBenchmarkPrice:                123456.123456,
		DpBid:                           1234567.1234567,
		DpAsk:                           321123,
		CurrentBlockNumber:              123456789,
		CurrentBlockHash:                common.HexToHash("0x123321").String(),
		CurrentBlockTimestamp:           987654321,
		BridgeTaskRunStartedTimestamp:   trrsMercuryV1[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:     trrsMercuryV1[0].FinishedAt.Time.UnixMilli(),
		ProviderRequestedTimestamp:      92233720368547760,
		ProviderReceivedTimestamp:       -92233720368547760,
		ProviderDataStreamEstablished:   1,
		ProviderIndicatedTime:           -123456789,
		Feed:                            common.HexToHash("0x111").String(),
		ObservationBenchmarkPrice:       111111,
		ObservationBid:                  222222,
		ObservationAsk:                  333333,
		ConfigDigest:                    "0200000000000000000000000000000000000000000000000000000000000000",
		Round:                           22,
		Epoch:                           11,
		BridgeRequestData:               `{"data":{"to":"LINK","from":"USD"}}`,
		AssetSymbol:                     "USD/LINK",
		ObservationBenchmarkPriceString: "111111",
		ObservationBidString:            "222222",
		ObservationAskString:            "333333",
	}

	expectedMessage, _ := proto.Marshal(&expectedTelemetry)
	wg.Wait()
	require.Equal(t, expectedMessage, sentMessage)

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: pipeline.TaskRunResults{
			pipeline.TaskRunResult{Task: &pipeline.BridgeTask{
				Name:     "test-mercury-bridge-1",
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
				Result: pipeline.Result{
					Value: nil,
				}},
		},
		V1Observation: &mercuryv1.Observation{},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}
	wg.Add(1)
	trrsMercuryV1[0].Result.Value = ""
	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: trrsMercuryV1,
		V1Observation:  &mercuryv1.Observation{},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}

	wg.Wait()
	require.Equal(t, 2, logs.Len())
	require.Contains(t, logs.All()[0].Message, `cannot get bridge response from bridge task, id=ds1, name="test-mercury-bridge-1"`)
	require.Contains(t, logs.All()[1].Message, "cannot parse EA telemetry")
	chDone <- struct{}{}
}

func TestCollectMercuryEnhancedTelemetryV2(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("test-network", "test-chainID", "0xa", synchronization.EnhancedEAMercury)

	var sentMessage []byte
	ingressClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		sentMessage = args[1].([]byte)
		wg.Done()
	})

	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	chTelem := make(chan EnhancedTelemetryMercuryData, 100)
	chDone := make(chan struct{})
	feedID := common.HexToHash("0x111")
	e := EnhancedTelemetryService[EnhancedTelemetryMercuryData]{
		chDone:  chDone,
		chTelem: chTelem,
		job: &job.Job{
			Type: job.Type(pipeline.OffchainReporting2JobType),
			OCR2OracleSpec: &job.OCR2OracleSpec{
				CaptureEATelemetry: true,
				FeedID:             &feedID,
			},
		},
		lggr:               lggr,
		monitoringEndpoint: monitoringEndpoint,
	}
	servicetest.Run(t, &e)

	wg.Add(1)

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: trrsMercuryV2,
		V2Observation: &mercuryv2.Observation{
			BenchmarkPrice:        mercury.ObsResult[*big.Int]{Val: big.NewInt(111111)},
			MaxFinalizedTimestamp: mercury.ObsResult[int64]{Val: 321},
			LinkPrice:             mercury.ObsResult[*big.Int]{Val: big.NewInt(4321)},
			NativePrice:           mercury.ObsResult[*big.Int]{Val: big.NewInt(54321)},
		},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}

	expectedTelemetry := telem.EnhancedEAMercury{
		DataSource:                      "data-source-name",
		DpBenchmarkPrice:                123456.123456,
		CurrentBlockNumber:              0,
		CurrentBlockHash:                "",
		CurrentBlockTimestamp:           0,
		BridgeTaskRunStartedTimestamp:   trrsMercuryV1[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:     trrsMercuryV1[0].FinishedAt.Time.UnixMilli(),
		ProviderRequestedTimestamp:      92233720368547760,
		ProviderReceivedTimestamp:       -92233720368547760,
		ProviderDataStreamEstablished:   1,
		ProviderIndicatedTime:           -123456789,
		Feed:                            common.HexToHash("0x111").String(),
		ObservationBenchmarkPrice:       111111,
		ObservationBid:                  0,
		ObservationAsk:                  0,
		ConfigDigest:                    "0200000000000000000000000000000000000000000000000000000000000000",
		Round:                           22,
		Epoch:                           11,
		BridgeRequestData:               `{"data":{"to":"LINK","from":"USD"}}`,
		AssetSymbol:                     "USD/LINK",
		ObservationBenchmarkPriceString: "111111",
		MaxFinalizedTimestamp:           321,
		LinkPrice:                       4321,
		NativePrice:                     54321,
	}

	expectedMessage, _ := proto.Marshal(&expectedTelemetry)
	wg.Wait()

	require.Equal(t, expectedMessage, sentMessage)

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: pipeline.TaskRunResults{
			pipeline.TaskRunResult{Task: &pipeline.BridgeTask{
				Name:     "test-mercury-bridge-2",
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
				Result: pipeline.Result{
					Value: nil,
				}},
		},
		V2Observation: &mercuryv2.Observation{},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}
	wg.Add(1)
	trrsMercuryV2[0].Result.Value = ""
	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: trrsMercuryV2,
		V2Observation:  &mercuryv2.Observation{},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}

	wg.Wait()
	require.Equal(t, 4, logs.Len())
	require.Contains(t, logs.All()[0].Message, "cannot parse enhanced EA telemetry bid price")
	require.Contains(t, logs.All()[1].Message, "cannot get bridge response from bridge task")
	require.Contains(t, logs.All()[2].Message, "cannot parse EA telemetry")
	require.Contains(t, logs.All()[3].Message, "cannot parse enhanced EA telemetry bid price")
	chDone <- struct{}{}
}

func TestCollectMercuryEnhancedTelemetryV4(t *testing.T) {
	ingressClient := mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("test-network", "test-chainID", "0xa", synchronization.EnhancedEAMercury)

	sentMessageCh := make(chan []byte)
	ingressClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		sentMessageCh <- args[1].([]byte)
	})

	lggr, _ := logger.TestLoggerObserved(t, zap.WarnLevel)
	chTelem := make(chan EnhancedTelemetryMercuryData, 100)
	chDone := make(chan struct{})
	feedID := common.HexToHash("0x0004")
	e := EnhancedTelemetryService[EnhancedTelemetryMercuryData]{
		chDone:  chDone,
		chTelem: chTelem,
		job: &job.Job{
			Type: job.Type(pipeline.OffchainReporting2JobType),
			OCR2OracleSpec: &job.OCR2OracleSpec{
				CaptureEATelemetry: true,
				FeedID:             &feedID,
			},
		},
		lggr:               lggr,
		monitoringEndpoint: monitoringEndpoint,
	}
	servicetest.Run(t, &e)

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: trrsMercuryV4,
		FeedVersion:    4,
		V4Observation: &mercuryv4.Observation{
			BenchmarkPrice:        mercury.ObsResult[*big.Int]{Val: big.NewInt(111111)},
			MarketStatus:          mercury.ObsResult[uint32]{Val: 2},
			MaxFinalizedTimestamp: mercury.ObsResult[int64]{Val: 321},
			LinkPrice:             mercury.ObsResult[*big.Int]{Val: big.NewInt(4321)},
			NativePrice:           mercury.ObsResult[*big.Int]{Val: big.NewInt(54321)},
		},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}

	expectedPricingTelemetry := telem.EnhancedEAMercury{
		DataSource:                      "data-source-name",
		DpBenchmarkPrice:                123456.123456,
		BridgeTaskRunStartedTimestamp:   trrsMercuryV4[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:     trrsMercuryV4[0].FinishedAt.Time.UnixMilli(),
		ProviderRequestedTimestamp:      92233720368547760,
		ProviderReceivedTimestamp:       -92233720368547760,
		ProviderDataStreamEstablished:   1,
		ProviderIndicatedTime:           -123456789,
		Feed:                            common.HexToHash("0x0004").String(),
		ObservationBenchmarkPrice:       111111,
		ObservationMarketStatus:         2,
		ConfigDigest:                    "0200000000000000000000000000000000000000000000000000000000000000",
		Round:                           22,
		Epoch:                           11,
		AssetSymbol:                     "USD/LINK",
		ObservationBenchmarkPriceString: "111111",
		MaxFinalizedTimestamp:           321,
		LinkPrice:                       4321,
		NativePrice:                     54321,
		Version:                         4,
		BridgeRequestData:               `{"data":{"to":"LINK","from":"USD"}}`,
	}
	expectedPricingMessage, _ := proto.Marshal(&expectedPricingTelemetry)
	require.Equal(t, expectedPricingMessage, <-sentMessageCh)

	chDone <- struct{}{}

	// Verify that no other telemetry is sent.
	require.Len(t, sentMessageCh, 0)
}
