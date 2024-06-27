package ocrcommon

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	mercuryv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	contractAddress := ethkey.EIP55Address(utils.RandomAddress().String())

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
	j.OCROracleSpec.EVMChainID = (*utils.Big)(big.NewInt(1234567890))
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
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.EnhancedEA, "test-network", "test-chainID")

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
			ContractAddress:    ethkey.EIP55AddressFromAddress(feedAddress),
			CaptureEATelemetry: true,
			EVMChainID:         (*utils.Big)(big.NewInt(9)),
		},
	}

	lggr, _ := logger.TestLoggerObserved(t, zap.WarnLevel)
	doneCh := make(chan struct{})
	enhancedTelemService := NewEnhancedTelemetryService(&jb, enhancedTelemChan, doneCh, monitoringEndpoint, lggr.Named("Enhanced Telemetry Mercury"))
	require.NoError(t, enhancedTelemService.Start(testutils.Context(t)))
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
	//enhancedTelemService.StopOnce("EnhancedTelemetryService", func() error { return nil })
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
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.EnhancedEA, "test-network", "test-chainID")
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
	require.NoError(t, enhancedTelemService.Start(testutils.Context(t)))
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
	assert.Contains(t, logs.All()[0].Message, "cannot get bridge response from bridge task")

	badTrrs = &pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: "[]",
			},
		}}
	wg.Add(1)
	enhancedTelemChan <- EnhancedTelemetryData{
		TaskRunResults: *badTrrs,
		FinalResults:   *finalResult,
		RepTimestamp:   observationTimestamp,
	}
	wg.Wait()
	assert.Equal(t, logs.Len(), 4)
	assert.Contains(t, logs.All()[2].Message, "cannot parse EA telemetry")
	assert.Contains(t, logs.All()[3].Message, "cannot get json parse value")
	doneCh <- struct{}{}
}

var trrsMercury = pipeline.TaskRunResults{
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
			Value: int64(321123),
		},
	},
}

func TestGetFinalValues(t *testing.T) {
	e := EnhancedTelemetryService[EnhancedTelemetryMercuryData]{}
	o := mercuryv1.Observation{
		BenchmarkPrice:        mercury.ObsResult[*big.Int]{Val: big.NewInt(111111)},
		Bid:                   mercury.ObsResult[*big.Int]{Val: big.NewInt(222222)},
		Ask:                   mercury.ObsResult[*big.Int]{Val: big.NewInt(333333)},
		CurrentBlockNum:       mercury.ObsResult[int64]{Val: 123456789},
		CurrentBlockHash:      mercury.ObsResult[[]byte]{Val: common.HexToHash("0x123321").Bytes()},
		CurrentBlockTimestamp: mercury.ObsResult[uint64]{Val: 987654321},
	}

	benchmarkPrice, bid, ask, blockNr, blockHash, blockTimestamp := e.getFinalValues(o)
	require.Equal(t, benchmarkPrice, int64(111111))
	require.Equal(t, bid, int64(222222))
	require.Equal(t, ask, int64(333333))
	require.Equal(t, blockNr, int64(123456789))
	require.Equal(t, blockHash, common.HexToHash("0x123321").Bytes())
	require.Equal(t, blockTimestamp, uint64(987654321))

	benchmarkPrice, bid, ask, blockNr, blockHash, blockTimestamp = e.getFinalValues(mercuryv1.Observation{})
	require.Equal(t, benchmarkPrice, int64(0))
	require.Equal(t, bid, int64(0))
	require.Equal(t, ask, int64(0))
	require.Equal(t, blockNr, int64(0))
	require.Nil(t, blockHash)
	require.Equal(t, blockTimestamp, uint64(0))
}

func TestGetPricesFromResults(t *testing.T) {
	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	e := EnhancedTelemetryService[EnhancedTelemetryMercuryData]{
		lggr: lggr,
		job: &job.Job{
			ID: 0,
		},
	}

	benchmarkPrice, bid, ask := e.getPricesFromResults(trrsMercury[0], &trrsMercury)
	require.Equal(t, 123456.123456, benchmarkPrice)
	require.Equal(t, 1234567.1234567, bid)
	require.Equal(t, float64(321123), ask)

	benchmarkPrice, bid, ask = e.getPricesFromResults(trrsMercury[0], &pipeline.TaskRunResults{})
	require.Equal(t, float64(0), benchmarkPrice)
	require.Equal(t, float64(0), bid)
	require.Equal(t, float64(0), ask)
	require.Equal(t, 1, logs.Len())
	require.Contains(t, logs.All()[0].Message, "cannot parse enhanced EA telemetry")

	tt := trrsMercury[:2]
	e.getPricesFromResults(trrsMercury[0], &tt)
	require.Equal(t, 2, logs.Len())
	require.Contains(t, logs.All()[1].Message, "cannot parse enhanced EA telemetry bid price, task is nil")

	tt = trrsMercury[:3]
	e.getPricesFromResults(trrsMercury[0], &tt)
	require.Equal(t, 3, logs.Len())
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
	benchmarkPrice, bid, ask = e.getPricesFromResults(trrsMercury[0], &trrs2)
	require.Equal(t, benchmarkPrice, float64(0))
	require.Equal(t, bid, float64(0))
	require.Equal(t, ask, float64(0))
	require.Equal(t, logs.Len(), 6)
	require.Contains(t, logs.All()[3].Message, "cannot parse enhanced EA telemetry benchmark price")
	require.Contains(t, logs.All()[4].Message, "cannot parse enhanced EA telemetry bid price")
	require.Contains(t, logs.All()[5].Message, "cannot parse enhanced EA telemetry ask price")
}

func TestShouldCollectEnhancedTelemetryMercury(t *testing.T) {

	j := &job.Job{
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

func TestGetAssetSymbolFromRequestData(t *testing.T) {
	e := EnhancedTelemetryService[EnhancedTelemetryMercuryData]{}
	require.Equal(t, e.getAssetSymbolFromRequestData(""), "")
	reqData := `{"data":{"to":"LINK","from":"USD"}}`
	require.Equal(t, e.getAssetSymbolFromRequestData(reqData), "USD/LINK")
}

func TestCollectMercuryEnhancedTelemetry(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryService(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.EnhancedEAMercury, "test-network", "test-chainID")

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
	require.NoError(t, e.Start(testutils.Context(t)))

	wg.Add(1)

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: trrsMercury,
		Observation: mercuryv1.Observation{
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
		DataSource:                    "data-source-name",
		DpBenchmarkPrice:              123456.123456,
		DpBid:                         1234567.1234567,
		DpAsk:                         321123,
		CurrentBlockNumber:            123456789,
		CurrentBlockHash:              common.HexToHash("0x123321").String(),
		CurrentBlockTimestamp:         987654321,
		BridgeTaskRunStartedTimestamp: trrsMercury[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:   trrsMercury[0].FinishedAt.Time.UnixMilli(),
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

	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: pipeline.TaskRunResults{
			pipeline.TaskRunResult{Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
				Result: pipeline.Result{
					Value: nil,
				}},
		},
		Observation: mercuryv1.Observation{},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}
	wg.Add(1)
	trrsMercury[0].Result.Value = ""
	chTelem <- EnhancedTelemetryMercuryData{
		TaskRunResults: trrsMercury,
		Observation:    mercuryv1.Observation{},
		RepTimestamp: types.ReportTimestamp{
			ConfigDigest: types.ConfigDigest{2},
			Epoch:        11,
			Round:        22,
		},
	}

	wg.Wait()
	require.Equal(t, 2, logs.Len())
	require.Contains(t, logs.All()[0].Message, "cannot get bridge response from bridge task")
	require.Contains(t, logs.All()[1].Message, "cannot parse EA telemetry")
	chDone <- struct{}{}
}
