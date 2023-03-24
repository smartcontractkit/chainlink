package mercury

import (
	"context"
	"encoding/base64"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/core/utils"
	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var (
	feedID              = [32]byte{'f', 'o', 'o'}
	sampleReport        = mercuryactions.BuildSampleReport(feedID)
	sig2                = ocrtypes.AttributedOnchainSignature{Signature: mustDecodeBase64("kbeuRczizOJCxBzj7MUAFpz3yl2WRM6K/f0ieEBvA+oTFUaKslbQey10krumVjzAvlvKxMfyZo0WkOgNyfF6xwE="), Signer: 2}
	sig3                = ocrtypes.AttributedOnchainSignature{Signature: mustDecodeBase64("9jz4b6Dh2WhXxQ97a6/S9UNjSfrEi9016XKTrfN0mLQFDiNuws23x7Z4n+6g0sqKH/hnxx1VukWUH/ohtw83/wE="), Signer: 3}
	sampleSigs          = []ocrtypes.AttributedOnchainSignature{sig2, sig3}
	sampleReportContext = ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: mustHexToConfigDigest("0x0001fc30092226b37f6924b464e16a54a7978a9a524519a73403af64d487dc45"),
			Epoch:        6,
			Round:        28,
		},
		ExtraHash: [32]uint8{27, 144, 106, 73, 166, 228, 123, 166, 179, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114},
	}
	rs, ss, vs = genRsSsVs()
)

func genRsSsVs() ([][32]byte, [][32]byte, [32]byte) {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range sampleSigs {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	return rs, ss, vs
}

func TestMercuryServerAPI(t *testing.T) {
	testEnv, err := mercury.NewEnv(t.Name(), "lukaszf-smoke", mercury.DefaultResources)

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	admin := mercury.User{
		Id:       testEnv.MSInfo.UserId,
		Key:      "admintestkey",
		Secret:   "mz1I4AgYtvo3Wumrgtlyh9VWkCf/IzZ6JROnuw==",
		Role:     "admin",
		Disabled: false,
	}
	user := mercury.User{
		Id:       genUuid(),
		Key:      "admintestkey",
		Secret:   "mz1I4AgYtvo3Wumrgtlyh9VWkCf/IzZ6JROnuw==",
		Role:     "user",
		Disabled: false,
	}
	initUsers := []mercury.User{admin, user}

	// Setup mercury server with mocked rpc conf
	rpcPubKey, nodesCsaKeys, err := testEnv.AddMercuryServer(&initUsers)
	require.NoError(t, err)
	msUrl := testEnv.MSInfo.LocalUrl
	csaKey := nodesCsaKeys[0]

	// Start wsrpc client for one of the rpc nodes defined in the conf
	wsrpcLggr, _ := logger.NewLogger()
	wsrpcUrl := testEnv.MSInfo.LocalWsrpcUrl[6:len(testEnv.MSInfo.LocalWsrpcUrl)]
	wsrpcClient := wsrpc.NewClient(wsrpcLggr, csaKey.KeyV2, rpcPubKey, wsrpcUrl)
	err = wsrpcClient.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() {
		err := wsrpcClient.Close()
		require.NoError(t, err)
	})

	t.Run("GET /admin/user as admin role", func(t *testing.T) {
		c := client.NewMercuryServerClient(msUrl, admin.Id, admin.Key)
		users, resp, err := c.GetUsers()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t, len(initUsers)+1, len(users)) // include bootstrap user
	})

	t.Run("GET /admin/user as user role", func(t *testing.T) {
		c := client.NewMercuryServerClient(msUrl, user.Id, user.Key)
		users, resp, err := c.GetUsers()
		require.NoError(t, err)
		require.Equal(t, 401, resp.StatusCode)
		require.Equal(t, 0, len(users))
	})

	t.Run("WSRPC LatestReport() returns no report for empty db", func(t *testing.T) {
		err := testEnv.ClearMercuryReportsInDb()
		require.NoError(t, err)

		req := &pb.LatestReportRequest{
			FeedId: feedID[:],
		}
		resp, err := wsrpcClient.LatestReport(context.Background(), req)
		require.NoError(t, err)
		require.Nil(t, resp.Report)
	})

	t.Run("WSRPC LatestReport() returns newly created report", func(t *testing.T) {
		var rs [][32]byte
		var ss [][32]byte
		var vs [32]byte
		for i, as := range sampleSigs {
			r, s, v, err := evmutil.SplitSignature(as.Signature)
			if err != nil {
				panic("eventTransmit(ev): error in SplitSignature")
			}
			rs = append(rs, r)
			ss = append(ss, s)
			vs[i] = v
		}
		rawReportCtx := evmutil.RawReportContext(sampleReportContext)

		payload, err := PayloadTypes.Pack(rawReportCtx, []byte(sampleReport), rs, ss, vs)
		require.NoError(t, err)

		// Transmit report
		req := &pb.TransmitRequest{
			Payload: payload,
		}
		res, err := wsrpcClient.Transmit(context.Background(), req)
		require.NoError(t, err)
		_ = res

		// Get latest report
		req2 := &pb.LatestReportRequest{
			FeedId: feedID[:],
		}
		res2, err := wsrpcClient.LatestReport(context.Background(), req2)
		require.NoError(t, err)
		require.NotNil(t, res2.Report)
		// TODO: validate report fields
	})

	t.Run("WSRPC: validFromBlockNum must be less or equal to currentBlockNumber", func(t *testing.T) {
		report := mercuryactions.Report{
			FeedId:                feedID,
			ObservationsTimestamp: uint32(42),
			BenchmarkPrice:        big.NewInt(242),
			Bid:                   big.NewInt(243),
			Ask:                   big.NewInt(244),
			CurrentBlockNum:       uint64(143),
			CurrentBlockHash:      utils.NewHash(),
			ValidFromBlockNum:     uint64(144),
		}
		reportBytes, err := report.Pack()
		require.NoError(t, err)
		rawReportCtx := evmutil.RawReportContext(sampleReportContext)
		payload, err := PayloadTypes.Pack(rawReportCtx, reportBytes, rs, ss, vs)
		require.NoError(t, err)

		// Transmit report
		req := &pb.TransmitRequest{
			Payload: payload,
		}
		res, _ := wsrpcClient.Transmit(context.Background(), req)
		require.Equal(t, "failed to store report", res.Error)
	})

	t.Run("WSRPC: on block overlap", func(t *testing.T) {
		err := testEnv.ClearMercuryReportsInDb()
		require.NoError(t, err)

		report1 := mercuryactions.Report{
			FeedId:                feedID,
			ObservationsTimestamp: uint32(42),
			BenchmarkPrice:        big.NewInt(242),
			Bid:                   big.NewInt(243),
			Ask:                   big.NewInt(244),
			CurrentBlockNum:       uint64(201),
			CurrentBlockHash:      utils.NewHash(),
			ValidFromBlockNum:     uint64(200),
		}
		reportBytes1, err := report1.Pack()
		require.NoError(t, err)
		rawReportCtx1 := evmutil.RawReportContext(sampleReportContext)
		payload1, err := PayloadTypes.Pack(rawReportCtx1, reportBytes1, rs, ss, vs)
		require.NoError(t, err)

		// Transmit report
		req1 := &pb.TransmitRequest{
			Payload: payload1,
		}
		res1, _ := wsrpcClient.Transmit(context.Background(), req1)
		require.Empty(t, res1.Error)

		report2 := mercuryactions.Report{
			FeedId:                feedID,
			ObservationsTimestamp: uint32(42),
			BenchmarkPrice:        big.NewInt(242),
			Bid:                   big.NewInt(243),
			Ask:                   big.NewInt(244),
			CurrentBlockNum:       uint64(201),
			CurrentBlockHash:      utils.NewHash(),
			ValidFromBlockNum:     uint64(200),
		}
		reportBytes2, err := report2.Pack()
		require.NoError(t, err)
		rawReportCtx2 := evmutil.RawReportContext(sampleReportContext)
		payload2, err := PayloadTypes.Pack(rawReportCtx2, reportBytes2, rs, ss, vs)
		require.NoError(t, err)

		// Transmit report
		req2 := &pb.TransmitRequest{
			Payload: payload2,
		}
		res2, _ := wsrpcClient.Transmit(context.Background(), req2)
		require.Empty(t, res2.Error)

		// Check only 1 report was saved in the db
		reports, err := testEnv.GetAllReportsFromMercuryDb()
		require.Empty(t, err)
		require.Equal(t, 1, len(reports))
	})
}

var PayloadTypes = mercuryactions.GetPayloadTypes()

func mustDecodeBase64(s string) (b []byte) {
	var err error
	b, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustHexToConfigDigest(s string) (cd ocrtypes.ConfigDigest) {
	b := hexutil.MustDecode(s)
	var err error
	cd, err = ocrtypes.BytesToConfigDigest(b)
	if err != nil {
		panic(err)
	}
	return
}

func genUuid() string {
	id, _ := uuid.NewV4()
	return id.String()
}
