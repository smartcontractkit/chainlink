package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCapability_Execute_FastPath_ContextCancelled(t *testing.T) {
	// Not parallel: mutates the global cresettings.DefaultGetter.

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := time.Minute
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	buf := NewFastPathBuffer(lggr, clock, expiry)

	getter, err := settings.NewJSONGetter([]byte(`{"global":{"VaultFastPathGetSecretsEnabled":"true"}}`))
	require.NoError(t, err)
	oldGetter := cresettings.DefaultGetter
	cresettings.DefaultGetter = getter
	t.Cleanup(func() { cresettings.DefaultGetter = oldGetter })

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, limits.Factory{Settings: getter}, newTestRequestLifecycleTracker(t), buf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	owner := "testowner"
	gsr := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{{
			Id:             &vault.SecretIdentifier{Owner: owner, Namespace: "main", Key: "k"},
			EncryptionKeys: []string{"enc"},
		}},
	}
	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Cancel the context before the buffer is completed.
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err = capability.Execute(ctx, capabilities.CapabilityRequest{
		Method:  vaulttypes.MethodSecretsGet,
		Payload: anyproto,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       owner,
			WorkflowID:          "wf",
			WorkflowExecutionID: "exec",
			ReferenceID:         "ref",
		},
	})
	require.ErrorIs(t, err, context.Canceled)

	all, err := store.All()
	require.NoError(t, err)
	require.Empty(t, all, "fast-path request should not enter OCR store")
}

func TestCapability_Execute_FastPath_ErrorResponse(t *testing.T) {
	// Not parallel: mutates the global cresettings.DefaultGetter.

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := time.Minute
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	buf := NewFastPathBuffer(lggr, clock, expiry)

	getter, err := settings.NewJSONGetter([]byte(`{"global":{"VaultFastPathGetSecretsEnabled":"true"}}`))
	require.NoError(t, err)
	oldGetter := cresettings.DefaultGetter
	cresettings.DefaultGetter = getter
	t.Cleanup(func() { cresettings.DefaultGetter = oldGetter })

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, limits.Factory{Settings: getter}, newTestRequestLifecycleTracker(t), buf)
	require.NoError(t, err)
	servicetest.Run(t, capability)

	owner := "testowner"
	gsr := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{{
			Id:             &vault.SecretIdentifier{Owner: owner, Namespace: "main", Key: "k"},
			EncryptionKeys: []string{"enc"},
		}},
	}
	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)

	requestID := fmt.Sprintf("%s::%s::%s", "wf", "exec", "ref")
	go func() {
		time.Sleep(10 * time.Millisecond)
		// The OCR plugin drains pending requests before serving them; simulate that here.
		_ = buf.Drain()
		buf.Complete(requestID, &vaulttypes.Response{ID: requestID, Error: "key does not exist"})
	}()

	_, err = capability.Execute(context.Background(), capabilities.CapabilityRequest{
		Method:  vaulttypes.MethodSecretsGet,
		Payload: anyproto,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       owner,
			WorkflowID:          "wf",
			WorkflowExecutionID: "exec",
			ReferenceID:         "ref",
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "key does not exist")

	all, err := store.All()
	require.NoError(t, err)
	require.Empty(t, all, "fast-path request should not enter OCR store")
}
