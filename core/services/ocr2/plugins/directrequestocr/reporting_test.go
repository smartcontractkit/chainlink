package directrequestocr_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	drocr_serv "github.com/smartcontractkit/chainlink/core/services/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr/config"
)

func reqID(id int) drocr_serv.RequestID {
	byteArr := (*[32]byte)([]byte(fmt.Sprintf("%032d\n", id)))
	return *byteArr
}

func sliceToByte32(slice []byte) [32]byte {
	var res [32]byte
	copy(res[:], slice[:32])
	return res
}

func preparePlugin(t *testing.T, batchSize uint32) (types.ReportingPlugin, drocr_serv.ORM) {
	ocrLogger := logger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})

	orm := drocr_serv.NewInMemoryORM()
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

func createRequest(t *testing.T, orm drocr_serv.ORM, id [32]byte) {
	testTxHash := common.HexToHash("0xabc")
	err := orm.CreateRequest(id, time.Now(), &testTxHash)
	require.NoError(t, err)
}

func createRequestWithResult(t *testing.T, orm drocr_serv.ORM, id [32]byte, result []byte) {
	createRequest(t, orm, id)
	err := orm.SetResult(id, 1, result, time.Now())
	require.NoError(t, err)
}

func createRequestWithError(t *testing.T, orm drocr_serv.ORM, id [32]byte, errVal []byte) {
	createRequest(t, orm, id)
	err := orm.SetError(id, 1, drocr_serv.USER_EXCEPTION, errVal, time.Now())
	require.NoError(t, err)
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

func TestDRReporting_Query_PickOnlyReadyRequests(t *testing.T) {
	t.Parallel()
	plugin, orm := preparePlugin(t, 10)
	reqId1, reqId2 := reqID(13), reqID(67)

	// Two requests but only one ready
	createRequestWithResult(t, orm, reqId1, []byte{})
	createRequest(t, orm, reqId2)

	q, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	require.NoError(t, err)

	queryProto := &directrequestocr.Query{}
	err = proto.Unmarshal(q, queryProto)
	require.NoError(t, err)
	require.Equal(t, 1, len(queryProto.RequestIDs))
	require.Equal(t, reqId1[:], queryProto.RequestIDs[0])
}

func TestDRReporting_Query_LimitToBatchSize(t *testing.T) {
	t.Parallel()
	plugin, orm := preparePlugin(t, 5)

	for i := 0; i < 20; i++ {
		createRequestWithResult(t, orm, reqID(10+i), []byte{})
	}

	// 20 results are ready but batch size is only 5
	q, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	require.NoError(t, err)

	queryProto := &directrequestocr.Query{}
	err = proto.Unmarshal(q, queryProto)
	require.NoError(t, err)
	require.Equal(t, 5, len(queryProto.RequestIDs))
	uniqueCounter := make(map[[32]byte]bool)
	for _, r := range queryProto.RequestIDs {
		uniqueCounter[sliceToByte32(r)] = true
	}
	require.Equal(t, 5, len(uniqueCounter), "wrong number of unique IDs")
}

func TestDRReporting_Observation(t *testing.T) {
	t.Parallel()
	plugin, orm := preparePlugin(t, 10)
	reqId1, reqId2, reqId3, reqId4 := reqID(13), reqID(14), reqID(15), reqID(16)

	createRequestWithResult(t, orm, reqId1, []byte("abc"))
	createRequest(t, orm, reqId2)
	createRequestWithError(t, orm, reqId3, []byte("Bug LOL!"))

	// Query asking for 4 requests (+ one duplicate) but we've only seen 3 of them, 2 of which are ready
	queryProto := directrequestocr.Query{}
	queryProto.RequestIDs = [][]byte{reqId1[:], reqId1[:], reqId2[:], reqId3[:], reqId4[:]}
	marshalled, err := proto.Marshal(&queryProto)
	require.NoError(t, err)

	obs, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, marshalled)
	require.NoError(t, err)

	observationProto := &directrequestocr.Observation{}
	err = proto.Unmarshal(obs, observationProto)
	require.NoError(t, err)
	require.Equal(t, len(observationProto.ProcessedRequests), 2)
	require.Equal(t, observationProto.ProcessedRequests[0].RequestID, reqId1[:])
	require.Equal(t, observationProto.ProcessedRequests[0].Result, []byte("abc"))
	require.Equal(t, len(observationProto.ProcessedRequests[0].Error), 0)
	require.Equal(t, observationProto.ProcessedRequests[1].RequestID, reqId3[:])
	require.Equal(t, len(observationProto.ProcessedRequests[1].Result), 0)
	require.Equal(t, observationProto.ProcessedRequests[1].Error, []byte("Bug LOL!"))
}

func TestDRReporting_Report(t *testing.T) {
	t.Parallel()
	plugin, _ := preparePlugin(t, 10)
	codec, err := directrequestocr.NewReportCodec()
	require.NoError(t, err)
	reqId1, reqId2, reqId3 := reqID(13), reqID(14), reqID(15)
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
