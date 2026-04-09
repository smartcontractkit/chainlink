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

	resolver := &testOrgResolver{orgID: "org-123"}
	payload := captureListRequest(t, "request-1", resolver, true, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "request-1",
		Owner:         "0xabc123",
		Namespace:     "ns",
		WorkflowOwner: "0xabc123",
	})

	require.NotNil(t, payload)
	assert.Equal(t, "org-123", payload.OrgId)
	assert.Equal(t, "0xabc123", payload.WorkflowOwner)
	assert.Equal(t, []string{"0xabc123"}, resolver.calledWith)
}

func TestCapability_ListSecretIdentifiers_OrgIDOnlySkipsResolver(t *testing.T) {
	t.Parallel()

	resolver := &testOrgResolver{orgID: "unexpected"}
	payload := captureListRequest(t, "request-2", resolver, true, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-2",
		Owner:     "0xabc123",
		Namespace: "ns",
		OrgId:     "org-999",
	})

	require.NotNil(t, payload)
	assert.Equal(t, "org-999", payload.OrgId)
	assert.Empty(t, payload.WorkflowOwner)
	assert.Empty(t, resolver.calledWith)
}

func TestCapability_ListSecretIdentifiers_VerifiesWorkflowOwnerAgainstOrgID(t *testing.T) {
	t.Parallel()

	resolver := &testOrgResolver{orgID: "org-999"}
	payload := captureListRequest(t, "request-verify", resolver, true, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "request-verify",
		Owner:         "0xabc123",
		Namespace:     "ns",
		OrgId:         "org-999",
		WorkflowOwner: "trusted-owner",
	})

	require.NotNil(t, payload)
	assert.Equal(t, "org-999", payload.OrgId)
	assert.Equal(t, "trusted-owner", payload.WorkflowOwner)
	assert.Equal(t, []string{"trusted-owner"}, resolver.calledWith)
}

func TestCapability_ListSecretIdentifiers_RejectsWorkflowOwnerOrgIDMismatch(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	resolver := &testOrgResolver{orgID: "org-actual"}

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, newVaultOrgIDAsSecretOwnerLimitsFactory(t, true))
	require.NoError(t, err)
	servicetest.Run(t, capability)

	_, err = capability.ListSecretIdentifiers(t.Context(), &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "request-mismatch",
		Owner:         "0xabc123",
		Namespace:     "ns",
		OrgId:         "org-request",
		WorkflowOwner: "trusted-owner",
	})
	require.ErrorContains(t, err, `workflow owner "trusted-owner" resolves to org_id "org-actual", does not match request org_id "org-request"`)
	assert.Equal(t, []string{"trusted-owner"}, resolver.calledWith)
	assert.Empty(t, store.GetByIDs([]string{"request-mismatch"}))
}

func TestCapability_ListSecretIdentifiers_GateClosedLeavesFieldsUntouched(t *testing.T) {
	t.Parallel()

	resolver := &testOrgResolver{orgID: "unexpected"}
	payload := captureListRequest(t, "request-3", resolver, false, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-3",
		Owner:     "0xabc123",
		Namespace: "ns",
	})

	require.NotNil(t, payload)
	assert.Empty(t, payload.OrgId)
	assert.Empty(t, payload.WorkflowOwner)
	assert.Empty(t, resolver.calledWith)
}

func TestCapability_ListSecretIdentifiers_RejectsMissingWorkflowOwnerWhenOrgIDMissing(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	resolver := &testOrgResolver{orgID: "org-actual"}

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, newVaultOrgIDAsSecretOwnerLimitsFactory(t, true))
	require.NoError(t, err)
	servicetest.Run(t, capability)

	_, err = capability.ListSecretIdentifiers(t.Context(), &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "request-missing-workflow-owner",
		Owner:     "0xabc123",
		Namespace: "ns",
	})
	require.ErrorContains(t, err, "org_id and workflow owner cannot both be empty")
	assert.Empty(t, resolver.calledWith)
	assert.Empty(t, store.GetByIDs([]string{"request-missing-workflow-owner"}))
}

func captureListRequest(t *testing.T, requestID string, resolver orgresolver.OrgResolver, gateEnabled bool, req *vaultcommon.ListSecretIdentifiersRequest) *vaultcommon.ListSecretIdentifiersRequest {
	t.Helper()

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, newVaultOrgIDAsSecretOwnerLimitsFactory(t, gateEnabled))
	require.NoError(t, err)
	servicetest.Run(t, capability)

	var (
		wg              sync.WaitGroup
		capturedPayload *vaultcommon.ListSecretIdentifiersRequest
		capturedOK      bool
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
				if !ok {
					return
				}
				copied, ok := payload.ProtoReflect().Interface().(*vaultcommon.ListSecretIdentifiersRequest)
				if !ok {
					return
				}
				capturedPayload = copied
				capturedOK = true
				reqs[0].SendResponse(t.Context(), &vaulttypes.Response{ID: requestID, Payload: []byte("ok")})
				return
			}
		}
	}()

	_, err = capability.ListSecretIdentifiers(t.Context(), req)
	require.NoError(t, err)
	wg.Wait()
	require.True(t, capturedOK)

	return capturedPayload
}

func newVaultOrgIDAsSecretOwnerLimitsFactory(t *testing.T, enabled bool) limits.Factory {
	t.Helper()

	getter, err := settings.NewJSONGetter([]byte(fmt.Sprintf(`{"global":{"VaultOrgIdAsSecretOwnerEnabled":%t}}`, enabled)))
	require.NoError(t, err)

	return limits.Factory{Settings: getter}
}
