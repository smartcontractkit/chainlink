package vault

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestFastPathBuffer_SubmitDrainComplete(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	buf := newFastPathBuffer(logger.TestLogger(t), clock, time.Minute)
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{{Id: &vaultcommon.SecretIdentifier{Key: "k"}}},
	}

	respCh := buf.Submit("req-1", req)
	drained := buf.Drain()
	require.Len(t, drained, 1)
	require.Equal(t, "req-1", drained[0].ID)

	buf.Complete("req-1", &vaulttypes.Response{ID: "req-1", Payload: []byte("ok"), Format: FastPathResponseFormat})
	select {
	case resp := <-respCh:
		require.Equal(t, []byte("ok"), resp.Payload)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for fast-path response")
	}
}

func TestFastPathBuffer_ExpireOlderThan(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	buf := newFastPathBuffer(logger.TestLogger(t), clock, time.Second)
	req := &vaultcommon.GetSecretsRequest{}
	respCh := buf.Submit("req-expire", req)

	clock.Advance(2 * time.Second)
	buf.ExpireOlderThan(clock.Now())

	select {
	case resp := <-respCh:
		require.NotEmpty(t, resp.Error)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for expiry response")
	}
}

func TestFastPathBuffer_CloseDrainsWithError(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	buf := newFastPathBuffer(logger.TestLogger(t), clock, time.Minute)
	req := &vaultcommon.GetSecretsRequest{}
	respCh := buf.Submit("req-close", req)
	buf.drainAllWithError("vault capability closed")

	select {
	case resp := <-respCh:
		require.Equal(t, "vault capability closed", resp.Error)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for close response")
	}
}
