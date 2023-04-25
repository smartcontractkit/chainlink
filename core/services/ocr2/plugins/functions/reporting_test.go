package functions_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	functions_srv "github.com/smartcontractkit/chainlink/v2/core/services/functions"
	functions_mocks "github.com/smartcontractkit/chainlink/v2/core/services/functions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
)

func preparePlugin(t *testing.T, batchSize uint32) (types.ReportingPlugin, *functions_mocks.ORM, *functions.ReportCodec) {
	lggr := logger.TestLogger(t)
	ocrLogger := logger.NewOCRWrapper(lggr, true, func(msg string) {})
	orm := functions_mocks.NewORM(t)
	factory := functions.FunctionsReportingPluginFactory{
		Logger:    ocrLogger,
		PluginORM: orm,
	}

	pluginConfig := config.ReportingPluginConfigWrapper{
		Config: &config.ReportingPluginConfig{
			MaxRequestBatchSize: batchSize,
		},
	}
	pluginConfigBytes, err := config.EncodeReportingPluginConfig(&pluginConfig)
	require.NoError(t, err)
	plugin, _, err := factory.NewReportingPlugin(types.ReportingPluginConfig{
		N:              4,
		F:              1,
		OffchainConfig: pluginConfigBytes,
	})
	require.NoError(t, err)
	codec, err := functions.NewReportCodec()
	require.NoError(t, err)
	return plugin, orm, codec
}

func newRequestID() functions_srv.RequestID {
	return testutils.Random32Byte()
}

func newRequest() functions_srv.Request {
	return functions_srv.Request{RequestID: newRequestID(), State: functions_srv.IN_PROGRESS}
}

func newRequestWithResult(result []byte) functions_srv.Request {
	return functions_srv.Request{RequestID: newRequestID(), State: functions_srv.RESULT_READY, Result: result}
}

func newRequestFinalized() functions_srv.Request {
	return functions_srv.Request{RequestID: newRequestID(), State: functions_srv.FINALIZED}
}

func newRequestTimedOut() functions_srv.Request {
	return functions_srv.Request{RequestID: newRequestID(), State: functions_srv.TIMED_OUT}
}

func newRequestConfirmed() functions_srv.Request {
	return functions_srv.Request{RequestID: newRequestID(), State: functions_srv.CONFIRMED}
}

func newMarshalledQuery(t *testing.T, reqIDs ...functions_srv.RequestID) []byte {
	queryProto := functions.Query{}
	queryProto.RequestIDs = [][]byte{}
	for _, id := range reqIDs {
		id := id
		queryProto.RequestIDs = append(queryProto.RequestIDs, id[:])
	}
	marshalled, err := proto.Marshal(&queryProto)
	require.NoError(t, err)
	return marshalled
}

func newProcessedRequest(requestId functions_srv.RequestID, compResult []byte, compError []byte) *functions.ProcessedRequest {
	return &functions.ProcessedRequest{
		RequestID: requestId[:],
		Result:    compResult,
		Error:     compError,
	}
}

func newObservation(t *testing.T, observerId uint8, requests ...*functions.ProcessedRequest) types.AttributedObservation {
	observationProto := functions.Observation{ProcessedRequests: requests}
	raw, err := proto.Marshal(&observationProto)
	require.NoError(t, err)
	return types.AttributedObservation{
		Observation: raw,
		Observer:    commontypes.OracleID(observerId),
	}
}

func TestDRReporting_Query(t *testing.T) {
	t.Parallel()
	const batchSize = 10
	plugin, orm, _ := preparePlugin(t, batchSize)
	reqs := []functions_srv.Request{newRequest(), newRequest()}
	orm.On("FindOldestEntriesByState", functions_srv.RESULT_READY, uint32(batchSize), mock.Anything).Return(reqs, nil)

	q, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	require.NoError(t, err)

	queryProto := &functions.Query{}
	err = proto.Unmarshal(q, queryProto)
	require.NoError(t, err)
	require.Equal(t, 2, len(queryProto.RequestIDs))
	require.Equal(t, reqs[0].RequestID[:], queryProto.RequestIDs[0])
	require.Equal(t, reqs[1].RequestID[:], queryProto.RequestIDs[1])
}

func TestDRReporting_Observation(t *testing.T) {
	t.Parallel()
	plugin, orm, _ := preparePlugin(t, 10)

	req1 := newRequestWithResult([]byte("abc"))
	req2 := newRequest()
	req3 := newRequestWithResult([]byte("def"))
	req4 := newRequestTimedOut()
	nonexistentId := newRequestID()

	orm.On("FindById", req1.RequestID, mock.Anything).Return(&req1, nil)
	orm.On("FindById", req2.RequestID, mock.Anything).Return(&req2, nil)
	orm.On("FindById", req3.RequestID, mock.Anything).Return(&req3, nil)
	orm.On("FindById", req4.RequestID, mock.Anything).Return(&req4, nil)
	orm.On("FindById", nonexistentId, mock.Anything).Return(nil, errors.New("nonexistent ID"))

	// Query asking for 5 requests (with duplicates), out of which:
	//   - two are ready
	//   - one is still in progress
	//   - one has timed out
	//   - one doesn't exist
	query := newMarshalledQuery(t, req1.RequestID, req1.RequestID, req2.RequestID, req3.RequestID, req4.RequestID, nonexistentId, req4.RequestID)
	obs, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, query)
	require.NoError(t, err)

	observationProto := &functions.Observation{}
	err = proto.Unmarshal(obs, observationProto)
	require.NoError(t, err)
	require.Equal(t, len(observationProto.ProcessedRequests), 2)
	require.Equal(t, observationProto.ProcessedRequests[0].RequestID, req1.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[0].Result, []byte("abc"))
	require.Equal(t, observationProto.ProcessedRequests[1].RequestID, req3.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[1].Result, []byte("def"))
}

func TestDRReporting_Observation_IncorrectQuery(t *testing.T) {
	t.Parallel()
	plugin, orm, _ := preparePlugin(t, 10)

	req1 := newRequestWithResult([]byte("abc"))
	invalidId := []byte("invalid")

	orm.On("FindById", req1.RequestID, mock.Anything).Return(&req1, nil)

	// Query asking for 3 requests (with duplicates), out of which:
	//   - two are invalid
	//   - one is ready
	queryProto := functions.Query{}
	queryProto.RequestIDs = [][]byte{invalidId, req1.RequestID[:], invalidId}
	marshalled, err := proto.Marshal(&queryProto)
	require.NoError(t, err)

	obs, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, marshalled)
	require.NoError(t, err)
	observationProto := &functions.Observation{}
	err = proto.Unmarshal(obs, observationProto)
	require.NoError(t, err)
	require.Equal(t, len(observationProto.ProcessedRequests), 1)
	require.Equal(t, observationProto.ProcessedRequests[0].RequestID, req1.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[0].Result, []byte("abc"))
}

func TestDRReporting_Report(t *testing.T) {
	t.Parallel()
	plugin, _, codec := preparePlugin(t, 10)
	reqId1, reqId2, reqId3 := newRequestID(), newRequestID(), newRequestID()
	compResult := []byte("aaa")
	procReq1 := newProcessedRequest(reqId1, compResult, []byte{})
	procReq2 := newProcessedRequest(reqId2, compResult, []byte{})

	query := newMarshalledQuery(t, reqId1, reqId2, reqId3, reqId1, reqId2) // duplicates should be ignored
	obs := []types.AttributedObservation{
		newObservation(t, 1, procReq2, procReq1),
		newObservation(t, 2, procReq1, procReq2),
	}

	// Two observations are not enough to produce a report
	produced, reportBytes, err := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	require.False(t, produced)
	require.Nil(t, reportBytes)
	require.NoError(t, err)

	// Three observations with the same requestID should produce a report
	obs = append(obs, newObservation(t, 3, procReq1, procReq2))
	produced, reportBytes, err = plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	require.True(t, produced)
	require.NoError(t, err)

	decoded, err := codec.DecodeReport(reportBytes)
	require.NoError(t, err)
	require.Equal(t, 2, len(decoded))
	require.Equal(t, reqId1[:], decoded[0].RequestID)
	require.Equal(t, compResult, decoded[0].Result)
	require.Equal(t, []byte{}, decoded[0].Error)
	require.Equal(t, reqId2[:], decoded[1].RequestID)
	require.Equal(t, compResult, decoded[1].Result)
	require.Equal(t, []byte{}, decoded[1].Error)
}

func TestDRReporting_Report_DeterministicOrderOfRequests(t *testing.T) {
	t.Parallel()
	plugin, _, codec := preparePlugin(t, 10)
	reqId1, reqId2, reqId3 := newRequestID(), newRequestID(), newRequestID()
	compResult := []byte("aaa")

	query := newMarshalledQuery(t, reqId1, reqId2, reqId3, reqId1, reqId2) // duplicates should be ignored
	procReq1 := newProcessedRequest(reqId1, compResult, []byte{})
	procReq2 := newProcessedRequest(reqId2, compResult, []byte{})
	procReq3 := newProcessedRequest(reqId3, compResult, []byte{})
	obs := []types.AttributedObservation{
		newObservation(t, 1, procReq1, procReq2, procReq3),
		newObservation(t, 2, procReq2, procReq1, procReq3),
		newObservation(t, 3, procReq3, procReq2, procReq1),
	}

	produced1, reportBytes1, err1 := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	produced2, reportBytes2, err2 := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	require.True(t, produced1)
	require.True(t, produced2)
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.Equal(t, reportBytes1, reportBytes2)

	decoded, err := codec.DecodeReport(reportBytes1)
	require.NoError(t, err)
	require.Equal(t, 3, len(decoded))
}

func TestDRReporting_Report_IncorrectObservation(t *testing.T) {
	t.Parallel()
	plugin, _, _ := preparePlugin(t, 10)
	reqId1 := newRequestID()
	compResult := []byte("aaa")

	query := newMarshalledQuery(t, reqId1)
	req := newProcessedRequest(reqId1, compResult, []byte{})

	// There are 4 observations but all are coming from the same node
	obs := []types.AttributedObservation{newObservation(t, 1, req, req, req, req)}
	produced, reportBytes, err := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	require.False(t, produced)
	require.Nil(t, reportBytes)
	require.NoError(t, err)
}

func getReportBytes(t *testing.T, codec *functions.ReportCodec, reqs ...functions_srv.Request) []byte {
	var report []*functions.ProcessedRequest
	for _, req := range reqs {
		req := req
		report = append(report, &functions.ProcessedRequest{RequestID: req.RequestID[:], Result: req.Result})
	}
	reportBytes, err := codec.EncodeReport(report)
	require.NoError(t, err)
	return reportBytes
}

func TestDRReporting_ShouldAcceptFinalizedReport(t *testing.T) {
	t.Parallel()
	plugin, orm, codec := preparePlugin(t, 10)

	req1 := newRequestWithResult([]byte("xxx")) // nonexistent
	req2 := newRequestWithResult([]byte("abc"))
	req3 := newRequestFinalized()
	req4 := newRequestTimedOut()

	orm.On("FindById", req1.RequestID, mock.Anything).Return(nil, errors.New("nonexistent ID"))
	orm.On("FindById", req2.RequestID, mock.Anything).Return(&req2, nil)
	orm.On("SetFinalized", req2.RequestID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	orm.On("FindById", req3.RequestID, mock.Anything).Return(&req3, nil)
	orm.On("SetFinalized", req3.RequestID, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("same state"))
	orm.On("FindById", req4.RequestID, mock.Anything).Return(&req4, nil)
	orm.On("SetFinalized", req4.RequestID, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("already timed out"))

	// Attempting to transmit 2 requests, out of which:
	//   - one was already accepted for transmission earlier
	//   - one has timed out
	should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, getReportBytes(t, codec, req3, req4))
	require.NoError(t, err)
	require.False(t, should)

	// Attempting to transmit 2 requests, out of which:
	//   - one is ready
	//   - one was already accepted for transmission earlier
	should, err = plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, getReportBytes(t, codec, req2, req3))
	require.NoError(t, err)
	require.True(t, should)

	// Attempting to transmit 2 requests, out of which:
	//   - one doesn't exist
	//   - one has timed out
	should, err = plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, getReportBytes(t, codec, req1, req4))
	require.NoError(t, err)
	require.True(t, should)
}

func TestDRReporting_ShouldTransmitAcceptedReport(t *testing.T) {
	t.Parallel()
	plugin, orm, codec := preparePlugin(t, 10)

	req1 := newRequestWithResult([]byte("xxx")) // nonexistent
	req2 := newRequestWithResult([]byte("abc"))
	req3 := newRequestFinalized()
	req4 := newRequestTimedOut()
	req5 := newRequestConfirmed()

	orm.On("FindById", req1.RequestID, mock.Anything).Return(nil, errors.New("nonexistent ID"))
	orm.On("FindById", req2.RequestID, mock.Anything).Return(&req2, nil)
	orm.On("FindById", req3.RequestID, mock.Anything).Return(&req3, nil)
	orm.On("FindById", req4.RequestID, mock.Anything).Return(&req4, nil)
	orm.On("FindById", req5.RequestID, mock.Anything).Return(&req5, nil)

	// Attempting to transmit 2 requests, out of which:
	//   - one was already confirmed on chain
	//   - one has timed out
	should, err := plugin.ShouldTransmitAcceptedReport(testutils.Context(t), types.ReportTimestamp{}, getReportBytes(t, codec, req5, req4))
	require.NoError(t, err)
	require.False(t, should)

	// Attempting to transmit 2 requests, out of which:
	//   - one is ready
	//   - one in finalized
	should, err = plugin.ShouldTransmitAcceptedReport(testutils.Context(t), types.ReportTimestamp{}, getReportBytes(t, codec, req2, req3))
	require.NoError(t, err)
	require.True(t, should)

	// Attempting to transmit 2 requests, out of which:
	//   - one doesn't exist
	//   - one is ready
	should, err = plugin.ShouldTransmitAcceptedReport(testutils.Context(t), types.ReportTimestamp{}, getReportBytes(t, codec, req1, req2))
	require.NoError(t, err)
	require.True(t, should)
}
