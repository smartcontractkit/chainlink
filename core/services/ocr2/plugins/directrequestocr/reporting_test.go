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

func PreparePlugin(t *testing.T) (types.ReportingPlugin, drocr_serv.ORM) {
	ocrLogger := logger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})

	orm := drocr_serv.NewInMemoryORM()
	factory := directrequestocr.DirectRequestReportingPluginFactory{
		Logger:    ocrLogger,
		PluginORM: orm,
	}

	pluginConfig := config.ReportingPluginConfigWrapper{
		Config: &config.ReportingPluginConfig{
			MaxRequestBatchSize: 10,
		},
	}
	pluginConfigBytes, err := config.EncodeReportingPluginConfig(&pluginConfig)
	require.NoError(t, err)
	plugin, _, _ := factory.NewReportingPlugin(types.ReportingPluginConfig{
		N:              4,
		F:              1,
		OffchainConfig: pluginConfigBytes,
	})
	return plugin, orm
}

func intToByte32(id int) [32]byte {
	byteArr := (*[32]byte)([]byte(fmt.Sprintf("%032d", id)))
	return *byteArr
}

func TestDRReporting_Query(t *testing.T) {
	t.Parallel()
	plugin, orm := PreparePlugin(t)
	reqId1, reqId2 := intToByte32(13), intToByte32(67)
	txHash := common.HexToHash("0xabc")

	// Two requests but only one ready
	id1, err := orm.CreateRequest(reqId1, time.Now(), &txHash)
	require.NoError(t, err)
	_, err = orm.CreateRequest(reqId2, time.Now(), &txHash)
	require.NoError(t, err)
	err = orm.SetResult(id1, 1, []byte{}, time.Now())
	require.NoError(t, err)

	q, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	require.NoError(t, err)

	queryProto := &directrequestocr.Query{}
	err = proto.Unmarshal(q, queryProto)
	require.NoError(t, err)
	require.Equal(t, 1, len(queryProto.RequestIDs))
	require.Equal(t, reqId1[:], queryProto.RequestIDs[0])
}

func TestDRReporting_Observation(t *testing.T) {
	t.Parallel()
	plugin, orm := PreparePlugin(t)
	reqId1, reqId2, reqId3 := intToByte32(13), intToByte32(14), intToByte32(15)
	txHash := common.HexToHash("0xabc")

	id1, err := orm.CreateRequest(reqId1, time.Now(), &txHash)
	require.NoError(t, err)
	_, err = orm.CreateRequest(reqId2, time.Now(), &txHash)
	require.NoError(t, err)
	id3, err := orm.CreateRequest(reqId3, time.Now(), &txHash)
	require.NoError(t, err)

	// Query asking for 3 requests but we only have 2 of them ready
	queryProto := directrequestocr.Query{}
	queryProto.RequestIDs = [][]byte{reqId1[:], reqId2[:], reqId3[:]}
	marshalled, err := proto.Marshal(&queryProto)
	require.NoError(t, err)
	err = orm.SetResult(id1, 1, []byte("abc"), time.Now())
	require.NoError(t, err)
	err = orm.SetError(id3, 1, drocr_serv.USER_EXCEPTION, "Bug LOL!", time.Now())
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

func buildObservation(t *testing.T, requestId []byte, compResult []byte, compError []byte, observer uint8) types.AttributedObservation {
	observationProto := directrequestocr.Observation{
		ProcessedRequests: []*directrequestocr.ProcessedRequest{&directrequestocr.ProcessedRequest{
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

func TestDRReporting_Report(t *testing.T) {
	t.Parallel()
	plugin, _ := PreparePlugin(t)
	codec, err := directrequestocr.NewReportCodec()
	require.NoError(t, err)
	reqId1, reqId2, reqId3 := intToByte32(13), intToByte32(14), intToByte32(15)
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
