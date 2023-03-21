package smoke

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc/pb"
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
)

func TestMercuryServerAPI(t *testing.T) {
	testEnv, err := mercury.NewEnv(t.Name(), "smoke", mercury.DefaultResources)

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	admin := mercury.User{
		Id:       testEnv.MSInfo.AdminId,
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
	rpcNodeConf, csaKey := genMockedRpcNodeConf()
	err = testEnv.AddMercuryServer(&initUsers, rpcNodeConf)
	require.NoError(t, err)
	msUrl := testEnv.MSInfo.LocalUrl

	// Setup wsrpc client for one of the rpc nodes defined in the conf
	wsrpcLggr, _ := logger.NewLogger()
	wsrpcUrl := testEnv.MSInfo.LocalWsrpcUrl[6:len(testEnv.MSInfo.LocalWsrpcUrl)]
	wsrpcClient := wsrpc.NewClient(wsrpcLggr, csaKey, testEnv.MSInfo.RpcPubKey, wsrpcUrl)
	err = wsrpcClient.Start(context.Background())
	require.NoError(t, err)

	t.Run("GET /admin/user as admin role", func(t *testing.T) {
		c := client.NewMercuryServerClient(msUrl, admin.Id, admin.Key)
		users, resp, err := c.GetUsers()
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t, len(initUsers), len(users))
	})

	t.Run("GET /admin/user as user role", func(t *testing.T) {
		c := client.NewMercuryServerClient(msUrl, user.Id, user.Key)
		users, resp, err := c.GetUsers()
		require.NoError(t, err)
		require.Equal(t, 401, resp.StatusCode)
		require.Equal(t, 0, len(users))
	})

	t.Run("WSRPC LatestReport() empty", func(t *testing.T) {
		req := &pb.LatestReportRequest{
			FeedId: feedID[:],
		}
		resp, err := wsrpcClient.LatestReport(context.Background(), req)
		require.NoError(t, err)
		require.Nil(t, resp.Report)
	})

	// TODO: test validFromBlockNum must be less or equal to currentBlockNumber

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

func genMockedRpcNodeConf() (*[]mercury.RpcNode, csakey.KeyV2) {
	_, privKey, _ := ed25519.GenerateKey(rand.Reader)
	csaKey := csakey.Raw(privKey).Key()

	rpcNodeConf := &[]mercury.RpcNode{
		{
			Id:            "0",
			Status:        "active",
			NodeAddress:   []string{"0x9aF03D0296F21f59aB956e83f9d969F544a021Fa"},
			OracleAddress: "0x0000000000000000000000000000000000000000",
			CsaKeys: []mercury.CsaKey{
				{
					NodeName:    "0",
					NodeAddress: "0x9aF03D0296F21f59aB956e83f9d969F544a021Fa",
					PublicKey:   csaKey.PublicKeyString(),
				},
			},
			Ocr2ConfigPublicKey:   []string{"fdff12ced64d6419b432f5096aa9b3de04531cf923b0142095f3e40014e81305"},
			Ocr2OffchainPublicKey: []string{"93400913aedd411ed6ec5d13c83ca7d666636a43dfd1195d62b3f4c0e1e6ce49"},
			Ocr2OnchainPublicKey:  []string{"01f2b0776f613604149579c8aebcf6ccf091b765"},
		},
	}

	return rpcNodeConf, csaKey
}

func genUuid() string {
	id, _ := uuid.NewV4()
	return id.String()
}
