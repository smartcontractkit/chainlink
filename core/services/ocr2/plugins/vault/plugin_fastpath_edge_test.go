package vault

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestPlugin_serveFastPathRequest_SecretNotFoundError(t *testing.T) {
	clock := clockwork.NewFakeClock()
	buf := vaultcap.NewFastPathBuffer(logger.TestLogger(t), clock, time.Minute)
	r := newTestReportingPlugin(t, withFastPathSource(buf))

	getReq := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{
			Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "missing"},
		}},
	}
	respCh := buf.Submit("fastpath-req-not-found", getReq)
	fr := vaultcap.FastPathRequest{ID: "fastpath-req-not-found", Request: getReq, Expiry: clock.Now().Add(time.Minute)}
	buf.Drain()

	r.serveFastPathRequest(t.Context(), NewReadStore(&kv{m: make(map[string]response)}, r.metrics), fr)

	select {
	case resp := <-respCh:
		require.Equal(t, vaultcap.FastPathResponseFormat, resp.Format)
		vaultResp := &vaultcommon.GetSecretsResponse{}
		require.NoError(t, proto.Unmarshal(resp.Payload, vaultResp))
		require.Len(t, vaultResp.Responses, 1)
		require.NotEmpty(t, vaultResp.Responses[0].GetError(), "expected a per-secret error for a missing secret")
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for fast-path not-found response")
	}
}

func TestPlugin_Observation_FastPathDrainsMultipleRequests(t *testing.T) {
	clock := clockwork.NewFakeClock()
	buf := vaultcap.NewFastPathBuffer(logger.TestLogger(t), clock, time.Minute)
	r := newTestReportingPlugin(t, withFastPathSource(buf))

	ids := []string{"fastpath-a", "fastpath-b", "fastpath-c"}
	respChans := make([]<-chan *vaulttypes.Response, len(ids))
	for i, id := range ids {
		getReq := &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{{
				Id: &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: id},
			}},
		}
		respChans[i] = buf.Submit(id, getReq)
	}

	_, err := r.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, &kv{m: make(map[string]response)}, &blobber{})
	require.NoError(t, err)

	for i, id := range ids {
		select {
		case resp := <-respChans[i]:
			require.Equal(t, vaultcap.FastPathResponseFormat, resp.Format, "id=%s", id)
			require.Empty(t, resp.Error, "id=%s", id)
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for fast-path response for id=%s", id)
		}
	}
}
