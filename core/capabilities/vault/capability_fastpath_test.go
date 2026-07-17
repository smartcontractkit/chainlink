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
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"google.golang.org/protobuf/proto"

	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestCapability_Execute_UsesFastPathWhenGateEnabled(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := time.Minute
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	buf := NewFastPathBuffer(lggr, clock, expiry)

	getter, err := settings.NewJSONGetter([]byte(`{"global":{"VaultFastPathGetSecretsEnabled":"true"}}`))
	require.NoError(t, err)
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
		// The OCR plugin drains pending requests before serving them; simulate that here so
		// Complete finds the request in the in-flight map.
		_ = buf.Drain()
		vaultResp := &vault.GetSecretsResponse{
			Responses: []*vault.SecretResponse{{
				Id: &vault.SecretIdentifier{Owner: owner, Namespace: "main", Key: "k"},
				Result: &vault.SecretResponse_Data{
					Data: &vault.SecretData{EncryptedValue: "secret"},
				},
			}},
		}
		payload, err := proto.Marshal(vaultResp)
		require.NoError(t, err)
		buf.Complete(requestID, &vaulttypes.Response{
			ID:      requestID,
			Payload: payload,
			Format:  FastPathResponseFormat,
		})
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
	require.NoError(t, err)

	all, err := store.All()
	require.NoError(t, err)
	require.Empty(t, all, "fast-path request should not enter OCR store")
}
