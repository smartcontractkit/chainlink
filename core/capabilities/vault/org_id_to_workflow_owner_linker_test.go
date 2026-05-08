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
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
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
		Owner:     "org999",
		Namespace: "ns",
		OrgId:     "org999",
	})

	require.NotNil(t, payload)
	assert.Equal(t, "org999", payload.Owner)
	assert.Equal(t, "org999", payload.OrgId)
	assert.Empty(t, payload.WorkflowOwner)
	assert.Empty(t, resolver.calledWith)
}

func TestCapability_ListSecretIdentifiers_VerifiesWorkflowOwnerAgainstOrgID(t *testing.T) {
	t.Parallel()

	resolver := &testOrgResolver{orgID: "org-999"}
	payload := captureListRequest(t, "request-verify", resolver, true, &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "request-verify",
		Owner:         "trustedowner",
		Namespace:     "ns",
		OrgId:         "org-999",
		WorkflowOwner: "trustedowner",
	})

	require.NotNil(t, payload)
	assert.Equal(t, "trustedowner", payload.Owner)
	assert.Equal(t, "org-999", payload.OrgId)
	assert.Equal(t, "trustedowner", payload.WorkflowOwner)
	assert.Equal(t, []string{"trustedowner"}, resolver.calledWith)
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

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, newVaultOrgIDAsSecretOwnerLimitsFactory(t, true), newTestRequestLifecycleTracker(t))
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

func TestCapability_CreateSecrets_ResolvesOrgIDBeforeLabelValidation(t *testing.T) {
	t.Parallel()

	orgID := "org-123"
	workflowOwner := "0x0001020304050607080900010203040506070809"
	encryptedSecret, capability, store := newCapabilityWithOrgIDEncryptedSecret(t, orgID)

	request := &vaultcommon.CreateSecretsRequest{
		RequestId:     "request-create",
		WorkflowOwner: workflowOwner,
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "secret",
					Namespace: "main",
					Owner:     workflowOwner,
				},
				EncryptedValue: encryptedSecret,
			},
		},
	}

	captured := respondWithCapturedPayload[*vaultcommon.CreateSecretsRequest](t, store, request.RequestId)
	_, err := capability.CreateSecrets(t.Context(), request)
	require.NoError(t, err)
	result := <-captured
	require.NoError(t, result.err)
	payload := result.payload

	assert.Equal(t, orgID, payload.OrgId)
	assert.Equal(t, workflowOwner, payload.WorkflowOwner)
}

func TestCapability_UpdateSecrets_ResolvesOrgIDBeforeLabelValidation(t *testing.T) {
	t.Parallel()

	orgID := "org-123"
	workflowOwner := "0x0001020304050607080900010203040506070809"
	encryptedSecret, capability, store := newCapabilityWithOrgIDEncryptedSecret(t, orgID)

	request := &vaultcommon.UpdateSecretsRequest{
		RequestId:     "request-update",
		WorkflowOwner: workflowOwner,
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "secret",
					Namespace: "main",
					Owner:     workflowOwner,
				},
				EncryptedValue: encryptedSecret,
			},
		},
	}

	captured := respondWithCapturedPayload[*vaultcommon.UpdateSecretsRequest](t, store, request.RequestId)
	_, err := capability.UpdateSecrets(t.Context(), request)
	require.NoError(t, err)
	result := <-captured
	require.NoError(t, result.err)
	payload := result.payload

	assert.Equal(t, orgID, payload.OrgId)
	assert.Equal(t, workflowOwner, payload.WorkflowOwner)
}

func TestCapability_CreateSecrets_AllowsResolvedOrgIDOwner(t *testing.T) {
	t.Parallel()

	orgID := "org123"
	workflowOwner := "0x0001020304050607080900010203040506070809"
	encryptedSecret, capability, store := newCapabilityWithOrgIDEncryptedSecret(t, orgID)

	request := &vaultcommon.CreateSecretsRequest{
		RequestId:     "request-create-org-owner",
		WorkflowOwner: workflowOwner,
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "secret",
					Namespace: "main",
					Owner:     orgID,
				},
				EncryptedValue: encryptedSecret,
			},
		},
	}

	captured := respondWithCapturedPayload[*vaultcommon.CreateSecretsRequest](t, store, request.RequestId)
	_, err := capability.CreateSecrets(t.Context(), request)
	require.NoError(t, err)
	result := <-captured
	require.NoError(t, result.err)
	payload := result.payload

	assert.Equal(t, orgID, payload.OrgId)
	assert.Equal(t, workflowOwner, payload.WorkflowOwner)
	assert.Equal(t, orgID, payload.EncryptedSecrets[0].Id.Owner)
}

func TestCapability_RejectsOwnersOutsideResolvedIdentity(t *testing.T) {
	t.Parallel()

	orgID := "org-123"
	workflowOwner := "0x0001020304050607080900010203040506070809"
	encryptedSecret, capability, store := newCapabilityWithOrgIDEncryptedSecret(t, orgID)

	tests := []struct {
		name      string
		requestID string
		call      func() (*vaulttypes.Response, error)
	}{
		{
			name:      "create",
			requestID: "request-create-owner-mismatch",
			call: func() (*vaulttypes.Response, error) {
				return capability.CreateSecrets(t.Context(), &vaultcommon.CreateSecretsRequest{
					RequestId:     "request-create-owner-mismatch",
					WorkflowOwner: workflowOwner,
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{
							Id: &vaultcommon.SecretIdentifier{
								Key:       "secret",
								Namespace: "main",
								Owner:     "otherowner",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				})
			},
		},
		{
			name:      "update",
			requestID: "request-update-owner-mismatch",
			call: func() (*vaulttypes.Response, error) {
				return capability.UpdateSecrets(t.Context(), &vaultcommon.UpdateSecretsRequest{
					RequestId:     "request-update-owner-mismatch",
					WorkflowOwner: workflowOwner,
					EncryptedSecrets: []*vaultcommon.EncryptedSecret{
						{
							Id: &vaultcommon.SecretIdentifier{
								Key:       "secret",
								Namespace: "main",
								Owner:     "otherowner",
							},
							EncryptedValue: encryptedSecret,
						},
					},
				})
			},
		},
		{
			name:      "delete",
			requestID: "request-delete-owner-mismatch",
			call: func() (*vaulttypes.Response, error) {
				return capability.DeleteSecrets(t.Context(), &vaultcommon.DeleteSecretsRequest{
					RequestId:     "request-delete-owner-mismatch",
					WorkflowOwner: workflowOwner,
					Ids: []*vaultcommon.SecretIdentifier{
						{
							Key:       "secret",
							Namespace: "main",
							Owner:     "otherowner",
						},
					},
				})
			},
		},
		{
			name:      "list",
			requestID: "request-list-owner-mismatch",
			call: func() (*vaulttypes.Response, error) {
				return capability.ListSecretIdentifiers(t.Context(), &vaultcommon.ListSecretIdentifiersRequest{
					RequestId:     "request-list-owner-mismatch",
					Owner:         "otherowner",
					Namespace:     "main",
					WorkflowOwner: workflowOwner,
				})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.call()
			require.ErrorContains(t, err, "must match resolved workflow owner")
			assert.Empty(t, store.GetByIDs([]string{tc.requestID}))
		})
	}
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

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, newVaultOrgIDAsSecretOwnerLimitsFactory(t, true), newTestRequestLifecycleTracker(t))
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

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, resolver, newVaultOrgIDAsSecretOwnerLimitsFactory(t, gateEnabled), newTestRequestLifecycleTracker(t))
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

func newCapabilityWithOrgIDEncryptedSecret(t *testing.T, orgID string) (string, *Capability, *requests.Store[*vaulttypes.Request]) {
	t.Helper()

	_, pk, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	encryptedSecret, err := vaultutils.EncryptSecretWithOrgID("secret-value", pk, orgID)
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clock, expiry)
	reg := coreCapabilities.NewRegistry(lggr)
	lpk := NewLazyPublicKey()
	lpk.Set(pk)

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, lpk, &testOrgResolver{orgID: orgID}, newVaultOrgIDAsSecretOwnerLimitsFactory(t, true), newTestRequestLifecycleTracker(t))
	require.NoError(t, err)
	servicetest.Run(t, capability)

	return encryptedSecret, capability, store
}

type capturedPayload[T proto.Message] struct {
	payload T
	err     error
}

func respondWithCapturedPayload[T proto.Message](t *testing.T, store *requests.Store[*vaulttypes.Request], requestID string) <-chan capturedPayload[T] {
	t.Helper()

	captured := make(chan capturedPayload[T], 1)
	go func() {
		for {
			select {
			case <-t.Context().Done():
				return
			default:
				reqs := store.GetByIDs([]string{requestID})
				if len(reqs) != 1 {
					continue
				}

				payload, ok := reqs[0].Payload.(T)
				if !ok {
					captured <- capturedPayload[T]{err: fmt.Errorf("unexpected payload type %T", reqs[0].Payload)}
					return
				}

				clonedMessage := proto.Clone(payload)
				cloned, ok := clonedMessage.(T)
				if !ok {
					captured <- capturedPayload[T]{err: fmt.Errorf("unexpected cloned payload type %T", clonedMessage)}
					return
				}

				captured <- capturedPayload[T]{payload: cloned}
				reqs[0].SendResponse(t.Context(), &vaulttypes.Response{ID: requestID, Payload: []byte("ok")})
				return
			}
		}
	}()

	return captured
}

func newVaultOrgIDAsSecretOwnerLimitsFactory(t *testing.T, enabled bool) limits.Factory {
	t.Helper()

	getter, err := settings.NewJSONGetter([]byte(fmt.Sprintf(`{"global":{"VaultOrgIdAsSecretOwnerEnabled":%t}}`, enabled)))
	require.NoError(t, err)

	return limits.Factory{Settings: getter}
}
