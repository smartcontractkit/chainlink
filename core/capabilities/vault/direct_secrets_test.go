package vault

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type fakeDirectSecretsReader struct {
	resp *vault.GetSecretsResponse
	got  *vault.GetSecretsRequest
}

func (f *fakeDirectSecretsReader) GetSecretsDirect(_ context.Context, req *vault.GetSecretsRequest) (*vault.GetSecretsResponse, error) {
	f.got = req
	return f.resp, nil
}

func TestLazyDirectSecretsReader(t *testing.T) {
	t.Parallel()
	l := NewLazyDirectSecretsReader()
	assert.Nil(t, l.Get())

	a, b := &fakeDirectSecretsReader{}, &fakeDirectSecretsReader{}
	l.Set(a)
	l.Set(b)
	l.Clear(a)
	assert.Same(t, b, l.Get(), "clearing a replaced reader is a no-op")
	l.Clear(b)
	assert.Nil(t, l.Get())
}

const directTestOwner = "0x1111111111111111111111111111111111111111"

var directTestID = &vault.SecretIdentifier{Key: "Foo", Namespace: "Bar", Owner: directTestOwner}

func newDirectTestCapability(t *testing.T) (*Capability, *LazyDirectSecretsReader, *requests.Store[*vaulttypes.Request], capabilities.CapabilityRequest, *clockwork.FakeClock) {
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, 10*time.Second)
	directReader := NewLazyDirectSecretsReader()
	capability, err := NewCapability(lggr, clock, 10*time.Second, handler, registry.NewRegistry(lggr), nil, directReader, limits.Factory{Settings: cresettings.DefaultGetter}, newTestRequestLifecycleTracker(t))
	require.NoError(t, err)
	servicetest.Run(t, capability)

	payload, err := anypb.New(&vault.GetSecretsRequest{
		Requests:           []*vault.SecretRequest{{Id: directTestID, EncryptionKeys: []string{"key"}}},
		GetSecretsDirectly: true,
	})
	require.NoError(t, err)
	req := capabilities.CapabilityRequest{
		Payload: payload,
		Method:  vaulttypes.MethodSecretsGet,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       directTestOwner,
			WorkflowID:          "wf",
			WorkflowExecutionID: "exec",
			ReferenceID:         "ref",
		},
	}
	return capability, directReader, store, req, clock
}

func TestCapability_GetSecretsDirect(t *testing.T) {
	t.Parallel()

	t.Run("not ready", func(t *testing.T) {
		t.Parallel()
		capability, _, _, req, clock := newDirectTestCapability(t)
		errCh := make(chan error, 1)
		go func() {
			_, err := capability.Execute(t.Context(), req)
			errCh <- err
		}()
		timeout := time.After(10 * time.Second)
		for {
			select {
			case err := <-errCh:
				require.ErrorIs(t, err, ErrDirectReadNotReady)
				return
			case <-time.After(10 * time.Millisecond):
				clock.Advance(directReadWaitPollInterval)
			case <-timeout:
				t.Fatal("timed out waiting for the not-ready error")
			}
		}
	})

	t.Run("reader appears while waiting", func(t *testing.T) {
		t.Parallel()
		capability, directReader, store, req, clock := newDirectTestCapability(t)
		want := &vault.GetSecretsResponse{Responses: []*vault.SecretResponse{{
			Id:     directTestID,
			Result: &vault.SecretResponse_Data{Data: &vault.SecretData{EncryptedValue: "ct"}},
		}}}
		reader := &fakeDirectSecretsReader{resp: want}
		type result struct {
			resp capabilities.CapabilityResponse
			err  error
		}
		resCh := make(chan result, 1)
		go func() {
			resp, err := capability.Execute(t.Context(), req)
			resCh <- result{resp, err}
		}()
		clock.Advance(directReadWaitPollInterval)
		directReader.Set(reader)
		clock.Advance(directReadWaitPollInterval)
		select {
		case r := <-resCh:
			require.NoError(t, r.err)
			got := &vault.GetSecretsResponse{}
			require.NoError(t, r.resp.Payload.UnmarshalTo(got))
			assert.True(t, proto.Equal(want, got))
			assert.True(t, reader.got.GetSecretsDirectly)
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for the direct read to be served")
		}
		assert.Empty(t, store.GetByIDs([]string{vault.BuildWorkflowGetSecretsRequestID(req.Metadata)}), "nothing is queued for OCR")
	})

	t.Run("context cancelled while waiting", func(t *testing.T) {
		t.Parallel()
		capability, _, _, req, _ := newDirectTestCapability(t)
		ctx, cancel := context.WithCancel(t.Context())
		errCh := make(chan error, 1)
		go func() {
			_, err := capability.Execute(ctx, req)
			errCh <- err
		}()
		cancel()
		select {
		case err := <-errCh:
			require.ErrorIs(t, err, context.Canceled)
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for the cancelled wait to return")
		}
	})

	t.Run("served without OCR", func(t *testing.T) {
		t.Parallel()
		capability, directReader, store, req, _ := newDirectTestCapability(t)
		want := &vault.GetSecretsResponse{Responses: []*vault.SecretResponse{{
			Id:     directTestID,
			Result: &vault.SecretResponse_Data{Data: &vault.SecretData{EncryptedValue: "ct"}},
		}}}
		reader := &fakeDirectSecretsReader{resp: want}
		directReader.Set(reader)

		resp, err := capability.Execute(t.Context(), req)
		require.NoError(t, err)

		got := &vault.GetSecretsResponse{}
		require.NoError(t, resp.Payload.UnmarshalTo(got))
		assert.True(t, proto.Equal(want, got))
		assert.True(t, reader.got.GetSecretsDirectly)
		assert.Empty(t, store.GetByIDs([]string{vault.BuildWorkflowGetSecretsRequestID(req.Metadata)}), "nothing is queued for OCR")
	})

	t.Run("still validated", func(t *testing.T) {
		t.Parallel()
		capability, directReader, _, req, _ := newDirectTestCapability(t)
		directReader.Set(&fakeDirectSecretsReader{resp: &vault.GetSecretsResponse{}})
		req.Metadata.WorkflowOwner = "0x2222222222222222222222222222222222222222"
		_, err := capability.Execute(t.Context(), req)
		require.ErrorContains(t, err, "does not match workflow owner")
	})
}
