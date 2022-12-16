package directrequestocr_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	drocr_serv "github.com/smartcontractkit/chainlink/core/services/directrequestocr"
	drocr_mocks "github.com/smartcontractkit/chainlink/core/services/directrequestocr/mocks"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr/config"
)

func preparePlugin(t *testing.T, batchSize uint32) (types.ReportingPlugin, *drocr_mocks.ORM) {
	lggr := logger.TestLogger(t)
	ocrLogger := logger.NewOCRWrapper(lggr, true, func(msg string) {})
	orm := drocr_mocks.NewORM(t)
	factory := directrequestocr.DirectRequestReportingPluginFactory{
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
	return plugin, orm
}

func newRequestID() drocr_serv.RequestID {
	return testutils.Random32Byte()
}

func newRequest() drocr_serv.Request {
	return drocr_serv.Request{RequestID: newRequestID(), State: drocr_serv.IN_PROGRESS}
}

func newRequestWithResult(result []byte) drocr_serv.Request {
	return drocr_serv.Request{RequestID: newRequestID(), State: drocr_serv.RESULT_READY, Result: result}
}

func newRequestFinalized() drocr_serv.Request {
	return drocr_serv.Request{RequestID: newRequestID(), State: drocr_serv.FINALIZED}
}

func newRequestTimedOut() drocr_serv.Request {
	return drocr_serv.Request{RequestID: newRequestID(), State: drocr_serv.TIMED_OUT}
}

func newRequestConfirmed() drocr_serv.Request {
	return drocr_serv.Request{RequestID: newRequestID(), State: drocr_serv.CONFIRMED}
}

func buildObservation(t *testing.T, requestId []byte, compResult []byte, compError []byte, observer uint8) types.AttributedObservation {
	observationProto := directrequestocr.Observation{
		ProcessedRequests: []*directrequestocr.ProcessedRequest{{
			RequestID: requestId,
			Result:    compResult,
			Error:     compError,
		}},
	}
	raw, err := proto.Marshal(&observationProto)
	require.NoError(t, err)
	return types.AttributedObservation{
		Observation: raw,
		Observer:    commontypes.OracleID(observer),
	}
}

func TestDRReporting_Query(t *testing.T) {
	t.Parallel()
	const batchSize = 10
	plugin, orm := preparePlugin(t, batchSize)
	reqs := []drocr_serv.Request{newRequest(), newRequest()}
	orm.On("FindOldestEntriesByState", drocr_serv.RESULT_READY, uint32(batchSize)).Return(reqs, nil)

	q, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	require.NoError(t, err)

	queryProto := &directrequestocr.Query{}
	err = proto.Unmarshal(q, queryProto)
	require.NoError(t, err)
	require.Equal(t, 2, len(queryProto.RequestIDs))
	require.Equal(t, reqs[0].RequestID[:], queryProto.RequestIDs[0])
	require.Equal(t, reqs[1].RequestID[:], queryProto.RequestIDs[1])
}

func TestDRReporting_Observation(t *testing.T) {
	t.Parallel()
	plugin, orm := preparePlugin(t, 10)

	req1 := newRequestWithResult([]byte("abc"))
	req2 := newRequest()
	req3 := newRequestWithResult([]byte("def"))
	req4 := newRequestTimedOut()
	nonexistentId := newRequestID()

	orm.On("FindById", req1.RequestID).Return(&req1, nil)
	orm.On("FindById", req2.RequestID).Return(&req2, nil)
	orm.On("FindById", req3.RequestID).Return(&req3, nil)
	orm.On("FindById", req4.RequestID).Return(&req4, nil)
	orm.On("FindById", nonexistentId).Return(nil, errors.New("nonexistent ID"))

	// Query asking for 5 requests (with duplicates), out of which:
	//   - two are ready
	//   - one is still in progress
	//   - one has timed out
	//   - one doesn't exist
	queryProto := directrequestocr.Query{}
	queryProto.RequestIDs = [][]byte{req1.RequestID[:], req1.RequestID[:], req2.RequestID[:], req3.RequestID[:], req4.RequestID[:], nonexistentId[:], req4.RequestID[:]}
	marshalled, err := proto.Marshal(&queryProto)
	require.NoError(t, err)

	obs, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, marshalled)
	require.NoError(t, err)

	observationProto := &directrequestocr.Observation{}
	err = proto.Unmarshal(obs, observationProto)
	require.NoError(t, err)
	require.Equal(t, len(observationProto.ProcessedRequests), 2)
	require.Equal(t, observationProto.ProcessedRequests[0].RequestID, req1.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[0].Result, []byte("abc"))
	require.Equal(t, observationProto.ProcessedRequests[1].RequestID, req3.RequestID[:])
	require.Equal(t, observationProto.ProcessedRequests[1].Result, []byte("def"))
}

func TestDRReporting_Report(t *testing.T) {
	t.Parallel()
	plugin, _ := preparePlugin(t, 10)
	codec, err := directrequestocr.NewReportCodec()
	require.NoError(t, err)
	reqId1, reqId2, reqId3 := newRequestID(), newRequestID(), newRequestID()
	compResult := []byte("aaa")

	queryProto := directrequestocr.Query{}
	queryProto.RequestIDs = [][]byte{reqId1[:], reqId2[:], reqId3[:]}
	marshalledQuery, err := proto.Marshal(&queryProto)
	require.NoError(t, err)

	obs := []types.AttributedObservation{
		buildObservation(t, reqId1[:], compResult, []byte{}, 1),
		buildObservation(t, reqId1[:], compResult, []byte{}, 2),
	}

	// Two observations are not enough to produce a report
	produced, reportBytes, err := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, marshalledQuery, obs)
	require.False(t, produced)
	require.Nil(t, reportBytes)
	require.NoError(t, err)

	// Three observations with the same requestID should produce a report
	obs = append(obs, buildObservation(t, reqId1[:], compResult, []byte{}, 3))
	produced, reportBytes, err = plugin.Report(testutils.Context(t), types.ReportTimestamp{}, marshalledQuery, obs)
	require.True(t, produced)
	require.NoError(t, err)

	decoded, err := codec.DecodeReport(reportBytes)
	require.NoError(t, err)

	require.Equal(t, 1, len(decoded))
	require.Equal(t, reqId1[:], decoded[0].RequestID)
	require.Equal(t, compResult, decoded[0].Result)
	require.Equal(t, []byte{}, decoded[0].Error)
}

func getReportBytes(t *testing.T, codec *directrequestocr.ReportCodec, reqs ...drocr_serv.Request) []byte {
	var report []*directrequestocr.ProcessedRequest
	for _, req := range reqs {
		req := req
		report = append(report, &directrequestocr.ProcessedRequest{RequestID: req.RequestID[:], Result: req.Result})
	}
	reportBytes, err := codec.EncodeReport(report)
	require.NoError(t, err)
	return reportBytes
}

func TestDRReporting_ShouldAcceptFinalizedReport(t *testing.T) {
	t.Parallel()
	plugin, orm := preparePlugin(t, 10)
	codec, err := directrequestocr.NewReportCodec()
	require.NoError(t, err)

	req1 := newRequestWithResult([]byte("xxx")) // nonexistent
	req2 := newRequestWithResult([]byte("abc"))
	req3 := newRequestFinalized()
	req4 := newRequestTimedOut()

	orm.On("FindById", req1.RequestID).Return(nil, errors.New("nonexistent ID"))
	orm.On("FindById", req2.RequestID).Return(&req2, nil)
	orm.On("SetFinalized", req2.RequestID, mock.Anything, mock.Anything).Return(nil)
	orm.On("FindById", req3.RequestID).Return(&req3, nil)
	orm.On("SetFinalized", req3.RequestID, mock.Anything, mock.Anything).Return(errors.New("same state"))
	orm.On("FindById", req4.RequestID).Return(&req4, nil)
	orm.On("SetFinalized", req4.RequestID, mock.Anything, mock.Anything).Return(errors.New("already timed out"))

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
	plugin, orm := preparePlugin(t, 10)
	codec, err := directrequestocr.NewReportCodec()
	require.NoError(t, err)

	req1 := newRequestWithResult([]byte("xxx")) // nonexistent
	req2 := newRequestWithResult([]byte("abc"))
	req3 := newRequestFinalized()
	req4 := newRequestTimedOut()
	req5 := newRequestConfirmed()

	orm.On("FindById", req1.RequestID).Return(nil, errors.New("nonexistent ID"))
	orm.On("FindById", req2.RequestID).Return(&req2, nil)
	orm.On("FindById", req3.RequestID).Return(&req3, nil)
	orm.On("FindById", req4.RequestID).Return(&req4, nil)
	orm.On("FindById", req5.RequestID).Return(&req5, nil)

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
