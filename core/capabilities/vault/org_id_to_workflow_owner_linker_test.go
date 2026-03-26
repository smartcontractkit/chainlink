package vault

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ orgresolver.OrgResolver = (*testOrgResolver)(nil)

type testOrgResolver struct {
	orgID      string
	err        error
	calledWith []string
}

func (t *testOrgResolver) Get(_ context.Context, owner string) (string, error) {
	t.calledWith = append(t.calledWith, owner)
	return t.orgID, t.err
}

func (t *testOrgResolver) Start(context.Context) error { return nil }
func (t *testOrgResolver) Close() error                { return nil }
func (t *testOrgResolver) HealthReport() map[string]error {
	return map[string]error{t.Name(): nil}
}
func (t *testOrgResolver) Name() string { return "test-org-resolver" }
func (t *testOrgResolver) Ready() error { return nil }

func TestCapability_ListSecretIdentifiers_LinksOrgIDFromWorkflowOwner(t *testing.T) {
	t.Parallel()

	requestID := "0xabc123::request-1"
	resolver := &testOrgResolver{orgID: "org-123"}
	payload := captureListRequest(t, requestID, resolver, true, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: requestID,
		Owner:     "0xabc123",
		Namespace: "ns",
	})

	require.NotNil(t, payload)
	assert.Equal(t, "org-123", payload.OrgId)
	assert.Equal(t, "0xabc123", payload.WorkflowOwner)
	assert.Equal(t, []string{"0xabc123"}, resolver.calledWith)
}

func TestCapability_ListSecretIdentifiers_TrustedOrgIDSkipsResolver(t *testing.T) {
	t.Parallel()

	resolver := &testOrgResolver{orgID: "unexpected"}
	payload := captureListRequest(t, "request-2", resolver, true, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "request-2",
		Owner:         "0xabc123",
		Namespace:     "ns",
		OrgId:         "org-999",
		WorkflowOwner: "untrusted-owner",
	})

	require.NotNil(t, payload)
	assert.Equal(t, "org-999", payload.OrgId)
	assert.Equal(t, "untrusted-owner", payload.WorkflowOwner)
	assert.Empty(t, resolver.calledWith)
}

func TestCapability_ListSecretIdentifiers_GateClosedLeavesFieldsUntouched(t *testing.T) {
	t.Parallel()

	resolver := &testOrgResolver{orgID: "unexpected"}
	payload := captureListRequest(t, "0xabc123::request-3", resolver, false, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "0xabc123::request-3",
		Owner:     "0xabc123",
		Namespace: "ns",
	})

	require.NotNil(t, payload)
	assert.Empty(t, payload.OrgId)
	assert.Empty(t, payload.WorkflowOwner)
	assert.Empty(t, resolver.calledWith)
}

func captureListRequest(t *testing.T, requestID string, resolver orgresolver.OrgResolver, gateEnabled bool, req *vaultcommon.ListSecretIdentifiersRequest) *vaultcommon.ListSecretIdentifiersRequest {
	t.Helper()

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, newVaultJWTAuthLimitsFactory(t, gateEnabled))
	require.NoError(t, err)
	servicetest.Run(t, capability)

	var (
		wg              sync.WaitGroup
		capturedPayload *vaultcommon.ListSecretIdentifiersRequest
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) != 1 {
					continue
				}

				payload, ok := reqs[0].Payload.(*vaultcommon.ListSecretIdentifiersRequest)
				require.True(t, ok)
				copied, ok := payload.ProtoReflect().Interface().(*vaultcommon.ListSecretIdentifiersRequest)
				require.True(t, ok)
				capturedPayload = copied
				reqs[0].SendResponse(t.Context(), &vaulttypes.Response{ID: requestID, Payload: []byte("ok")})
				return
			}
		}
	}()

	_, err = capability.ListSecretIdentifiers(t.Context(), req)
	require.NoError(t, err)
	wg.Wait()

	return capturedPayload
}

func newVaultJWTAuthLimitsFactory(t *testing.T, enabled bool) limits.Factory {
	t.Helper()

	getter, err := settings.NewJSONGetter([]byte(fmt.Sprintf(`{"global":{"VaultJWTAuthEnabled":%t}}`, enabled)))
	require.NoError(t, err)

	return limits.Factory{Settings: getter}
}
