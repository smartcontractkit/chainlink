package vault

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPlugin_serveFastPathRequest_CompletesBufferedRequest(t *testing.T) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	clock := clockwork.NewFakeClock()
	buf := vaultcap.NewFastPathBuffer(logger.TestLogger(t), clock, time.Minute)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withVaultOptimizationsEnabled(), withFastPathSource(buf))

	owner := "0x0001020304050607080900010203040506070809"
	id := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "fastpathsecret"}
	rdr := &kv{m: make(map[string]response)}

	plaintext := []byte("fastpath-value")
	var label [32]byte
	copy(label[12:], common.HexToAddress(owner).Bytes())
	ciphertext, err := tdh2easy.EncryptWithLabel(pk, plaintext, label)
	require.NoError(t, err)
	ciphertextBytes, err := ciphertext.Marshal()
	require.NoError(t, err)
	require.NoError(t, newTestWriteStore(t, rdr).WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{
		EncryptedSecret: ciphertextBytes,
	}))

	pubK, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pks := hex.EncodeToString(pubK[:])
	getReq := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{
			Id:             id,
			EncryptionKeys: []string{pks},
		}},
	}
	respCh := buf.Submit("fastpath-req-1", getReq)
	fr := vaultcap.FastPathRequest{ID: "fastpath-req-1", Request: getReq, Expiry: clock.Now().Add(time.Minute)}
	buf.Drain()

	r.serveFastPathRequest(t.Context(), NewReadStore(rdr, r.metrics), fr)

	select {
	case resp := <-respCh:
		require.Empty(t, resp.Error)
		require.Equal(t, vaultcap.FastPathResponseFormat, resp.Format)
		vaultResp := &vaultcommon.GetSecretsResponse{}
		require.NoError(t, proto.Unmarshal(resp.Payload, vaultResp))
		require.Len(t, vaultResp.Responses, 1)
		require.Empty(t, vaultResp.Responses[0].GetError())
		require.NotEmpty(t, vaultResp.Responses[0].GetData().EncryptedDecryptionKeyShares)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for fast-path completion")
	}
}

func TestPlugin_Observation_FastPathDrainsBuffer(t *testing.T) {
	clock := clockwork.NewFakeClock()
	buf := vaultcap.NewFastPathBuffer(logger.TestLogger(t), clock, time.Minute)
	r := newTestReportingPlugin(t, withFastPathSource(buf))
	rdr := &kv{m: make(map[string]response)}

	getReq := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{
			Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "k"},
		}},
	}
	buf.Submit("fastpath-req-obs", getReq)

	_, err := r.Observation(t.Context(), 1, types.AttributedQuery{}, rdr, &blobber{})
	require.NoError(t, err)
}
