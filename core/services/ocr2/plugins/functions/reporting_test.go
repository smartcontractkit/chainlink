package functions_test

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	functions_srv "github.com/smartcontractkit/chainlink/v2/core/services/functions"
	functions_mocks "github.com/smartcontractkit/chainlink/v2/core/services/functions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
)

func preparePlugin(t *testing.T, batchSize uint32, maxTotalGasLimit uint32) (types.ReportingPlugin, *functions_mocks.ORM, encoding.ReportCodec, *functions_mocks.OffchainTransmitter) {
	lggr := logger.TestLogger(t)
	ocrLogger := commonlogger.NewOCRWrapper(lggr, true, func(msg string) {})
	orm := functions_mocks.NewORM(t)
	offchainTransmitter := functions_mocks.NewOffchainTransmitter(t)
	factory := functions.FunctionsReportingPluginFactory{
		Logger:              ocrLogger,
		PluginORM:           orm,
		ContractVersion:     1,
		OffchainTransmitter: offchainTransmitter,
	}

	pluginConfig := config.ReportingPluginConfigWrapper{
		Config: &config.ReportingPluginConfig{
			MaxRequestBatchSize:       batchSize,
			MaxReportTotalCallbackGas: maxTotalGasLimit,
		},
	}
	pluginConfigBytes, err := config.EncodeReportingPluginConfig(&pluginConfig)
	require.NoError(t, err)
	plugin, _, err := factory.NewReportingPlugin(tests.Context(t), types.ReportingPluginConfig{
		N:              4,
		F:              1,
		OffchainConfig: pluginConfigBytes,
	})
	require.NoError(t, err)
	codec, err := encoding.NewReportCodec(1)
	require.NoError(t, err)
	return plugin, orm, codec, offchainTransmitter
}

func newRequestID() functions_srv.RequestID {
	return testutils.Random32Byte()
}

func newRequest() functions_srv.Request {
	var gasLimit uint32 = 100000
	return functions_srv.Request{RequestID: newRequestID(), State: functions_srv.IN_PROGRESS, CoordinatorContractAddress: &common.Address{1}, CallbackGasLimit: &gasLimit}
}

func newRequestWithResult(result []byte) functions_srv.Request {
	req := newRequest()
	req.State = functions_srv.RESULT_READY
	req.Result = result
	return req
}

func newRequestFinalized() functions_srv.Request {
	req := newRequest()
	req.State = functions_srv.FINALIZED
	return req
}

func newRequestTimedOut() functions_srv.Request {
	req := newRequest()
	req.State = functions_srv.TIMED_OUT
	return req
}

func newRequestConfirmed() functions_srv.Request {
	req := newRequest()
	req.State = functions_srv.CONFIRMED
	return req
}

func newMarshalledQuery(t *testing.T, reqIDs ...functions_srv.RequestID) []byte {
	queryProto := encoding.Query{}
	queryProto.RequestIDs = [][]byte{}
	for _, id := range reqIDs {
		id := id
		queryProto.RequestIDs = append(queryProto.RequestIDs, id[:])
	}
	marshalled, err := proto.Marshal(&queryProto)
	require.NoError(t, err)
	return marshalled
}

func newProcessedRequest(requestId functions_srv.RequestID, compResult []byte, compError []byte) *encoding.ProcessedRequest {
	return &encoding.ProcessedRequest{
		RequestID:           requestId[:],
		Result:              compResult,
		Error:               compError,
		CoordinatorContract: []byte{1},
	}
}

func newProcessedRequestWithMeta(requestId functions_srv.RequestID, compResult []byte, compError []byte, callbackGasLimit uint32, coordinatorContract []byte, onchainMetadata []byte) *encoding.ProcessedRequest {
	return &encoding.ProcessedRequest{
		RequestID:           requestId[:],
		Result:              compResult,
		Error:               compError,
		CallbackGasLimit:    callbackGasLimit,
		CoordinatorContract: coordinatorContract,
		OnchainMetadata:     onchainMetadata,
	}
}

func newObservation(t *testing.T, observerId uint8, requests ...*encoding.ProcessedRequest) types.AttributedObservation {
	observationProto := encoding.Observation{ProcessedRequests: requests}
	raw, err := proto.Marshal(&observationProto)
	require.NoError(t, err)
	return types.AttributedObservation{
		Observation: raw,
		Observer:    commontypes.OracleID(observerId),
	}
}

func TestFunctionsReporting_Query(t *testing.T) {
	t.Parallel()
	const batchSize = 10
	plugin, orm, _, _ := preparePlugin(t, batchSize, 0)
	reqs := []functions_srv.Request{newRequest(), newRequest()}
	orm.On("FindOldestEntriesByState", mock.Anything, functions_srv.RESULT_READY, uint32(batchSize), mock.Anything).Return(reqs, nil)

	q, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	require.NoError(t, err)

	queryProto := &encoding.Query{}
	err = proto.Unmarshal(q, queryProto)
	require.NoError(t, err)
	require.Equal(t, 2, len(queryProto.RequestIDs))
	require.Equal(t, reqs[0].RequestID[:], queryProto.RequestIDs[0])
	require.Equal(t, reqs[1].RequestID[:], queryProto.RequestIDs[1])
}

func TestFunctionsReporting_Query_HandleCoordinatorMismatch(t *testing.T) {
	t.Parallel()
	const batchSize = 10
	plugin, orm, _, _ := preparePlugin(t, batchSize, 1000000)
	reqs := []functions_srv.Request{newRequest(), newRequest()}
	reqs[0].CoordinatorContractAddress = &common.Address{1}
	reqs[1].CoordinatorContractAddress = &common.Address{2}
	orm.On("FindOldestEntriesByState", mock.Anything, functions_srv.RESULT_READY, uint32(batchSize), mock.Anything).Return(reqs, nil)

	q, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	require.NoError(t, err)

	queryProto := &encoding.Query{}
	err = proto.Unmarshal(q, queryProto)
	require.NoError(t, err)
	require.Equal(t, 1, len(queryProto.RequestIDs))
	require.Equal(t, reqs[0].RequestID[:], queryProto.RequestIDs[0])
	// reqs[1] should be excluded from this query because it has a different coordinator address
}

func TestFunctionsReporting_Observation(t *testing.T) {
	t.Parallel()
	plugin, orm, _, _ := preparePlugin(t, 10, 0)

	req1 := newRequestWithResult([]byte("abc"))
	req2 := newRequest()
	req3 := newRequestWithResult([]byte("def"))
	req4 := newRequestTimedOut()
	nonexistentId := newRequestID()

	orm.On("FindById", mock.Anything, req1.RequestID, mock.Anything).Return(&req1, nil)
	orm.On("FindById", mock.Anything, req2.RequestID, mock.Anything).Return(&req2, nil)
	orm.On("FindById", mock.Anything, req3.RequestID, mock.Anything).Return(&req3, nil)
	orm.On("FindById", mock.Anything, req4.RequestID, mock.Anything).Return(&req4, nil)
	orm.On("FindById", mock.Anything, nonexistentId, mock.Anything).Return(nil, errors.New("nonexistent ID"))

	// Query asking for 5 requests (with duplicates), out of which:
	//   - two are ready
	//   - one is still in progress
	//   - one has timed out
	//   - one doesn't exist
	query := newMarshalledQuery(t, req1.RequestID, req1.RequestID, req2.RequestID, req3.RequestID, req4.RequestID, nonexistentId, req4.RequestID)
	obs, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, query)
	require.NoError(t, err)

	observationProto := &encoding.Observation{}
	err = proto.Unmarshal(obs, observationProto)
	require.NoError(t, err)
	require.Equal(t, len(observationProto.ProcessedRequests), 2)
	require.Equal(t, observationProto.ProcessedRequests[0].RequestID, req1.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[0].Result, []byte("abc"))
	require.Equal(t, observationProto.ProcessedRequests[1].RequestID, req3.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[1].Result, []byte("def"))
}

func TestFunctionsReporting_Observation_IncorrectQuery(t *testing.T) {
	t.Parallel()
	plugin, orm, _, _ := preparePlugin(t, 10, 0)

	req1 := newRequestWithResult([]byte("abc"))
	invalidId := []byte("invalid")

	orm.On("FindById", mock.Anything, req1.RequestID, mock.Anything).Return(&req1, nil)

	// Query asking for 3 requests (with duplicates), out of which:
	//   - two are invalid
	//   - one is ready
	queryProto := encoding.Query{}
	queryProto.RequestIDs = [][]byte{invalidId, req1.RequestID[:], invalidId}
	marshalled, err := proto.Marshal(&queryProto)
	require.NoError(t, err)

	obs, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, marshalled)
	require.NoError(t, err)
	observationProto := &encoding.Observation{}
	err = proto.Unmarshal(obs, observationProto)
	require.NoError(t, err)
	require.Equal(t, len(observationProto.ProcessedRequests), 1)
	require.Equal(t, observationProto.ProcessedRequests[0].RequestID, req1.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[0].Result, []byte("abc"))
}

func TestFunctionsReporting_Report(t *testing.T) {
	t.Parallel()
	plugin, _, codec, _ := preparePlugin(t, 10, 1000000)
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

func TestFunctionsReporting_Report_WithGasLimitAndMetadata(t *testing.T) {
	t.Parallel()
	plugin, _, codec, _ := preparePlugin(t, 10, 300000)
	reqId1, reqId2, reqId3 := newRequestID(), newRequestID(), newRequestID()
	compResult := []byte("aaa")
	gasLimit1, gasLimit2 := uint32(100_000), uint32(200_000)
	coordinatorContract := common.Address{1}
	meta1, meta2 := []byte("meta1"), []byte("meta2")
	procReq1 := newProcessedRequestWithMeta(reqId1, compResult, []byte{}, gasLimit1, coordinatorContract[:], meta1)
	procReq2 := newProcessedRequestWithMeta(reqId2, compResult, []byte{}, gasLimit2, coordinatorContract[:], meta2)

	query := newMarshalledQuery(t, reqId1, reqId2, reqId3, reqId1, reqId2) // duplicates should be ignored
	obs := []types.AttributedObservation{
		newObservation(t, 1, procReq2, procReq1),
		newObservation(t, 2, procReq1, procReq2),
		newObservation(t, 3, procReq1, procReq2),
	}

	produced, reportBytes, err := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	require.True(t, produced)
	require.NoError(t, err)

	decoded, err := codec.DecodeReport(reportBytes)
	require.NoError(t, err)
	require.Equal(t, 2, len(decoded))

	require.Equal(t, reqId1[:], decoded[0].RequestID)
	require.Equal(t, compResult, decoded[0].Result)
	require.Equal(t, []byte{}, decoded[0].Error)
	require.Equal(t, coordinatorContract[:], decoded[0].CoordinatorContract)
	require.Equal(t, meta1, decoded[0].OnchainMetadata)
	// CallbackGasLimit is not ABI-encoded

	require.Equal(t, reqId2[:], decoded[1].RequestID)
	require.Equal(t, compResult, decoded[1].Result)
	require.Equal(t, []byte{}, decoded[1].Error)
	require.Equal(t, coordinatorContract[:], decoded[1].CoordinatorContract)
	require.Equal(t, meta2, decoded[1].OnchainMetadata)
	// CallbackGasLimit is not ABI-encoded
}

func TestFunctionsReporting_Report_HandleCoordinatorMismatch(t *testing.T) {
	t.Parallel()
	plugin, _, codec, _ := preparePlugin(t, 10, 300000)
	reqId1, reqId2, reqId3 := newRequestID(), newRequestID(), newRequestID()
	compResult, meta := []byte("aaa"), []byte("meta")
	coordinatorContractA, coordinatorContractB := common.Address{1}, common.Address{2}
	procReq1 := newProcessedRequestWithMeta(reqId1, compResult, []byte{}, 0, coordinatorContractA[:], meta)
	procReq2 := newProcessedRequestWithMeta(reqId2, compResult, []byte{}, 0, coordinatorContractB[:], meta)
	procReq3 := newProcessedRequestWithMeta(reqId3, compResult, []byte{}, 0, coordinatorContractA[:], meta)

	query := newMarshalledQuery(t, reqId1, reqId2, reqId3, reqId1, reqId2) // duplicates should be ignored
	obs := []types.AttributedObservation{
		newObservation(t, 1, procReq2, procReq3, procReq1),
		newObservation(t, 2, procReq1, procReq2, procReq3),
		newObservation(t, 3, procReq3, procReq1, procReq2),
	}

	produced, reportBytes, err := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	require.True(t, produced)
	require.NoError(t, err)

	decoded, err := codec.DecodeReport(reportBytes)
	require.NoError(t, err)
	require.Equal(t, 2, len(decoded))

	require.Equal(t, reqId1[:], decoded[0].RequestID)
	require.Equal(t, reqId3[:], decoded[1].RequestID)
	// reqId2	should be excluded from this report because it has a different coordinator address
}

func TestFunctionsReporting_Report_CallbackGasLimitExceeded(t *testing.T) {
	t.Parallel()
	plugin, _, codec, _ := preparePlugin(t, 10, 200000)
	reqId1, reqId2 := newRequestID(), newRequestID()
	compResult := []byte("aaa")
	gasLimit1, gasLimit2 := uint32(100_000), uint32(200_000)
	coordinatorContract1, coordinatorContract2 := common.Address{1}, common.Address{2}
	procReq1 := newProcessedRequestWithMeta(reqId1, compResult, []byte{}, gasLimit1, coordinatorContract1[:], []byte{})
	procReq2 := newProcessedRequestWithMeta(reqId2, compResult, []byte{}, gasLimit2, coordinatorContract2[:], []byte{})

	query := newMarshalledQuery(t, reqId1, reqId2)
	obs := []types.AttributedObservation{
		newObservation(t, 1, procReq2, procReq1),
		newObservation(t, 2, procReq1, procReq2),
		newObservation(t, 3, procReq1, procReq2),
	}

	produced, reportBytes, err := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, query, obs)
	require.True(t, produced)
	require.NoError(t, err)

	decoded, err := codec.DecodeReport(reportBytes)
	require.NoError(t, err)
	// Gas limit is set to 200k per report so we can only fit the first request
	require.Equal(t, 1, len(decoded))
	require.Equal(t, reqId1[:], decoded[0].RequestID)
	require.Equal(t, compResult, decoded[0].Result)
	require.Equal(t, []byte{}, decoded[0].Error)
	require.Equal(t, coordinatorContract1[:], decoded[0].CoordinatorContract)
}

func TestFunctionsReporting_Report_DeterministicOrderOfRequests(t *testing.T) {
	t.Parallel()
	plugin, _, codec, _ := preparePlugin(t, 10, 0)
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

func TestFunctionsReporting_Report_IncorrectObservation(t *testing.T) {
	t.Parallel()
	plugin, _, _, _ := preparePlugin(t, 10, 0)
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

func getReportBytes(t *testing.T, codec encoding.ReportCodec, reqs ...functions_srv.Request) []byte {
	var report []*encoding.ProcessedRequest
	for _, req := range reqs {
		req := req
		report = append(report, &encoding.ProcessedRequest{
			RequestID:           req.RequestID[:],
			Result:              req.Result,
			Error:               req.Error,
			CallbackGasLimit:    *req.CallbackGasLimit,
			CoordinatorContract: req.CoordinatorContractAddress[:],
			OnchainMetadata:     req.OnchainMetadata,
		})
	}
	reportBytes, err := codec.EncodeReport(report)
	require.NoError(t, err)
	return reportBytes
}

func TestFunctionsReporting_ShouldAcceptFinalizedReport(t *testing.T) {
	t.Parallel()
	plugin, orm, codec, _ := preparePlugin(t, 10, 0)

	req1 := newRequestWithResult([]byte("xxx")) // nonexistent
	req2 := newRequestWithResult([]byte("abc"))
	req3 := newRequestFinalized()
	req4 := newRequestTimedOut()

	orm.On("FindById", mock.Anything, req1.RequestID, mock.Anything).Return(nil, errors.New("nonexistent ID"))
	orm.On("FindById", mock.Anything, req2.RequestID, mock.Anything).Return(&req2, nil)
	orm.On("SetFinalized", mock.Anything, req2.RequestID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	orm.On("FindById", mock.Anything, req3.RequestID, mock.Anything).Return(&req3, nil)
	orm.On("SetFinalized", mock.Anything, req3.RequestID, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("same state"))
	orm.On("FindById", mock.Anything, req4.RequestID, mock.Anything).Return(&req4, nil)
	orm.On("SetFinalized", mock.Anything, req4.RequestID, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("already timed out"))

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

func TestFunctionsReporting_ShouldAcceptFinalizedReport_OffchainTransmission(t *testing.T) {
	t.Parallel()
	plugin, orm, codec, offchainTransmitter := preparePlugin(t, 10, 0)
	req1 := newRequestWithResult([]byte("abc"))
	req1.OnchainMetadata = []byte(functions_srv.OffchainRequestMarker)

	orm.On("FindById", mock.Anything, req1.RequestID, mock.Anything).Return(&req1, nil)
	orm.On("SetFinalized", mock.Anything, req1.RequestID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	offchainTransmitter.On("TransmitReport", mock.Anything, mock.Anything).Return(nil)

	should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, getReportBytes(t, codec, req1))
	require.NoError(t, err)
	require.False(t, should)
}

func TestFunctionsReporting_ShouldTransmitAcceptedReport(t *testing.T) {
	t.Parallel()
	plugin, orm, codec, _ := preparePlugin(t, 10, 0)

	req1 := newRequestWithResult([]byte("xxx")) // nonexistent
	req2 := newRequestWithResult([]byte("abc"))
	req3 := newRequestFinalized()
	req4 := newRequestTimedOut()
	req5 := newRequestConfirmed()

	orm.On("FindById", mock.Anything, req1.RequestID, mock.Anything).Return(nil, errors.New("nonexistent ID"))
	orm.On("FindById", mock.Anything, req2.RequestID, mock.Anything).Return(&req2, nil)
	orm.On("FindById", mock.Anything, req3.RequestID, mock.Anything).Return(&req3, nil)
	orm.On("FindById", mock.Anything, req4.RequestID, mock.Anything).Return(&req4, nil)
	orm.On("FindById", mock.Anything, req5.RequestID, mock.Anything).Return(&req5, nil)

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

func TestFunctionsReporting_ShouldIncludeCoordinator(t *testing.T) {
	t.Parallel()

	zeroAddr, coord1, coord2 := &common.Address{}, &common.Address{1}, &common.Address{2}

	// should never pass nil requestCoordinator
	newCoord, err := functions.ShouldIncludeCoordinator(nil, nil)
	require.Error(t, err)
	require.Nil(t, newCoord)

	// should never pass zero requestCoordinator
	newCoord, err = functions.ShouldIncludeCoordinator(zeroAddr, nil)
	require.Error(t, err)
	require.Nil(t, newCoord)

	// overwrite nil reportCoordinator
	newCoord, err = functions.ShouldIncludeCoordinator(coord1, nil)
	require.NoError(t, err)
	require.Equal(t, coord1, newCoord)

	// same address is fine
	newCoord, err = functions.ShouldIncludeCoordinator(coord1, newCoord)
	require.NoError(t, err)
	require.Equal(t, coord1, newCoord)

	// different address is not accepted
	newCoord, err = functions.ShouldIncludeCoordinator(coord2, newCoord)
	require.Error(t, err)
	require.Equal(t, coord1, newCoord)
}
