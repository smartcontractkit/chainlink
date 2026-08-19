package vault

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
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
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// fakeMetadataRegistry resolves DONByID from an in-memory map, standing in for
// the node's local capabilities registry view.
type fakeMetadataRegistry struct {
	core.UnimplementedCapabilitiesRegistryMetadata
	dons map[uint32]capabilities.DON
}

func (f *fakeMetadataRegistry) DONByID(_ context.Context, donID uint32) (capabilities.DON, error) {
	d, ok := f.dons[donID]
	if !ok {
		return capabilities.DON{}, fmt.Errorf("could not find don %d", donID)
	}
	return d, nil
}

const (
	zoneBDonID = uint32(10)
	zoneADonID = uint32(20)
	// zoneBMixedCaseDonID is a zone-b DON whose registry family casing differs
	// from the zoneBFamily constant, guarding the case-insensitive match.
	zoneBMixedCaseDonID = uint32(30)
)

func newZoneBTestCapability(t *testing.T, settingsJSON string) *Capability {
	t.Helper()
	lggr := logger.TestLogger(t)
	clock := clockwork.NewFakeClock()
	expiry := 10 * time.Second
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler(lggr, store, clock, expiry)

	reg := coreCapabilities.NewRegistry(lggr)
	reg.SetLocalRegistry(&fakeMetadataRegistry{dons: map[uint32]capabilities.DON{
		zoneBDonID:          {ID: zoneBDonID, Name: "workflow_1_zone-b", Families: []string{"zone-b"}},
		zoneADonID:          {ID: zoneADonID, Name: "workflow_1_zone-a", Families: []string{"zone-a"}},
		zoneBMixedCaseDonID: {ID: zoneBMixedCaseDonID, Name: "workflow_1_zone-b_mixed", Families: []string{"Zone-B"}},
	}})

	getter, err := settings.NewJSONGetter([]byte(settingsJSON))
	require.NoError(t, err)
	lf := limits.Factory{Settings: getter}

	capability, err := NewCapability(lggr, clock, expiry, handler, reg, nil, lf, newTestRequestLifecycleTracker(t))
	require.NoError(t, err)
	servicetest.Run(t, capability)
	return capability
}

// allowlistedOwner is the normalized (no 0x, lowercase) form of the workflow
// owner used across these tests.
const allowlistedOwner = "1111111111111111111111111111111111111111"

func zoneBGetSecretsRequest(t *testing.T, owner string) capabilities.CapabilityRequest {
	t.Helper()
	gsr := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id:             &vault.SecretIdentifier{Key: "Foo", Namespace: "Bar", Owner: owner},
				EncryptionKeys: []string{"k"},
			},
		},
	}
	anyproto, err := anypb.New(gsr)
	require.NoError(t, err)
	return capabilities.CapabilityRequest{
		Payload: anyproto,
		Method:  vault.MethodGetSecrets,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner:       owner,
			WorkflowID:          "wf-id",
			WorkflowExecutionID: "exec-id",
			ReferenceID:         "ref-id",
			WorkflowDonID:       zoneBDonID,
		},
	}
}

func TestCapability_Execute_ZoneBRestriction_DeniesNonAllowlistedOwner(t *testing.T) {
	t.Parallel()
	// Master gate enabled, owner NOT allowlisted (default deny).
	capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"}}`)

	_, err := capability.Execute(t.Context(), zoneBGetSecretsRequest(t, "0x"+allowlistedOwner))
	require.Error(t, err)
	require.ErrorContains(t, err, "zone-b workflow DON may only read vault secrets for allowlisted workflow owners")
}

// MUST-2: exercises the allow path through Execute() rather than calling
// enforce() directly, covering the CRE-context population wiring. An allowlisted
// zone-b owner must pass the gate; with no OCR oracle in this unit test the
// request then blocks in handleRequest until the context deadline, so we assert
// it got past the gate (no zone-b denial) and reached the handler (timeout).
func TestCapability_Execute_ZoneBRestriction_AllowsAllowlistedOwner(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"},"owner":{"`+allowlistedOwner+`":{"PerOwner":{"VaultZoneBGetSecretsAllowed":"true"}}}}`)

	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	defer cancel()
	_, err := capability.Execute(ctx, zoneBGetSecretsRequest(t, "0x"+allowlistedOwner))
	require.Error(t, err)
	require.NotContains(t, err.Error(), "zone-b workflow DON may only read vault secrets for allowlisted workflow owners")
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestCapability_enforceZoneBWorkflowRestriction(t *testing.T) {
	t.Parallel()
	ctxWithOwner := func(donID uint32, owner string) (context.Context, uint32) {
		md := capabilities.RequestMetadata{WorkflowOwner: owner, WorkflowDonID: donID}
		return md.ContextWithCRE(t.Context()), donID
	}

	t.Run("gate disabled: allows zone-b caller regardless of owner", func(t *testing.T) {
		t.Parallel()
		capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"false"}}`)
		ctx, donID := ctxWithOwner(zoneBDonID, "0x"+allowlistedOwner)
		require.NoError(t, capability.zoneBRestrictor.enforce(ctx, donID))
	})

	t.Run("gate enabled: allows non-zone-b (zone-a) caller", func(t *testing.T) {
		t.Parallel()
		capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"}}`)
		ctx, donID := ctxWithOwner(zoneADonID, "0xdeadbeef")
		require.NoError(t, capability.zoneBRestrictor.enforce(ctx, donID))
	})

	t.Run("gate enabled: denies zone-b caller with non-allowlisted owner", func(t *testing.T) {
		t.Parallel()
		capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"}}`)
		ctx, donID := ctxWithOwner(zoneBDonID, "0x"+allowlistedOwner)
		err := capability.zoneBRestrictor.enforce(ctx, donID)
		require.ErrorContains(t, err, "allowlisted workflow owners")
	})

	t.Run("gate enabled: allows zone-b caller with allowlisted owner", func(t *testing.T) {
		t.Parallel()
		capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"},"owner":{"`+allowlistedOwner+`":{"PerOwner":{"VaultZoneBGetSecretsAllowed":"true"}}}}`)
		ctx, donID := ctxWithOwner(zoneBDonID, "0x"+allowlistedOwner)
		require.NoError(t, capability.zoneBRestrictor.enforce(ctx, donID))
	})

	t.Run("gate enabled: fails closed for unknown caller DON", func(t *testing.T) {
		t.Parallel()
		capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"}}`)
		md := capabilities.RequestMetadata{WorkflowOwner: "0x" + allowlistedOwner, WorkflowDonID: 999}
		err := capability.zoneBRestrictor.enforce(md.ContextWithCRE(t.Context()), 999)
		require.ErrorContains(t, err, "could not resolve caller workflow DON 999")
	})

	// MUST-3: guards the case-insensitive family match. If isZoneBWorkflowDON
	// regressed to an exact (==) comparison, a "Zone-B" family would be treated
	// as non-zone-b and this caller would be allowed, failing this test.
	t.Run("gate enabled: denies zone-b caller identified by case-insensitive family match", func(t *testing.T) {
		t.Parallel()
		capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"}}`)
		ctx, donID := ctxWithOwner(zoneBMixedCaseDonID, "0x"+allowlistedOwner)
		err := capability.zoneBRestrictor.enforce(ctx, donID)
		require.ErrorContains(t, err, "allowlisted workflow owners")
	})

	// MUST-4: the allowlist is keyed by the normalized owner (0x stripped,
	// lowercased), while callers supply the standard "0x"-prefixed, possibly
	// mixed-case address. If owner matching regressed to a bare string compare,
	// this allowlisted caller would be denied. Uses a hex owner with letters in
	// mixed case so both the 0x-strip and the lowercasing are exercised.
	t.Run("gate enabled: allowlist match normalizes 0x prefix and case of caller owner", func(t *testing.T) {
		t.Parallel()
		const mixedCaseOwner = "AbCdEf0000000000000000000000000000000001"
		const normalizedOwner = "abcdef0000000000000000000000000000000001"
		capability := newZoneBTestCapability(t, `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"},"owner":{"`+normalizedOwner+`":{"PerOwner":{"VaultZoneBGetSecretsAllowed":"true"}}}}`)
		ctx, donID := ctxWithOwner(zoneBDonID, "0x"+mixedCaseOwner)
		require.NoError(t, capability.zoneBRestrictor.enforce(ctx, donID))
	})
}

// toggleableMetadataRegistry resolves DONByID from an in-memory map but can be
// switched to return the same "information not available" error the real
// registry surfaces while its local view is nil/mid-sync.
type toggleableMetadataRegistry struct {
	core.UnimplementedCapabilitiesRegistryMetadata
	dons map[uint32]capabilities.DON
	fail atomic.Bool
}

func (f *toggleableMetadataRegistry) DONByID(_ context.Context, donID uint32) (capabilities.DON, error) {
	if f.fail.Load() {
		return capabilities.DON{}, errors.New("metadataRegistry information not available")
	}
	d, ok := f.dons[donID]
	if !ok {
		return capabilities.DON{}, fmt.Errorf("could not find don %d", donID)
	}
	return d, nil
}

// newOutageTestRestrictor builds a zone-b restrictor backed by a toggleable
// registry so tests can warm the cache and then simulate a registry outage.
func newOutageTestRestrictor(t *testing.T, settingsJSON string) (*zoneBRestrictor, *toggleableMetadataRegistry) {
	t.Helper()
	lggr := logger.TestLogger(t)
	reg := coreCapabilities.NewRegistry(lggr)
	fake := &toggleableMetadataRegistry{dons: map[uint32]capabilities.DON{
		zoneADonID: {ID: zoneADonID, Name: "workflow_1_zone-a", Families: []string{"zone-a"}},
		zoneBDonID: {ID: zoneBDonID, Name: "workflow_1_zone-b", Families: []string{"zone-b"}},
	}}
	reg.SetLocalRegistry(fake)

	getter, err := settings.NewJSONGetter([]byte(settingsJSON))
	require.NoError(t, err)
	z, err := newZoneBRestrictor(lggr, limits.Factory{Settings: getter}, reg)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, z.close()) })
	return z, fake
}

// MUST-2: a transient capabilities registry outage must not fail every vault
// read DON-wide. Once a DON's zone membership has been resolved, a subsequent
// registry lookup failure falls back to the cached value instead of erroring;
// a DON never resolved before still fails closed. Crucially, the cached value
// preserves the security decision in both directions: a cached zone-b DON stays
// denied during the outage rather than slipping through.
func TestZoneBRestrictor_RegistryOutageUsesCachedMembership(t *testing.T) {
	t.Parallel()
	const restrictEnabled = `{"global":{"VaultZoneBWorkflowGetSecretsRestrictEnabled":"true"}}`

	t.Run("cached zone-a membership survives outage (still allowed)", func(t *testing.T) {
		t.Parallel()
		z, fake := newOutageTestRestrictor(t, restrictEnabled)
		md := capabilities.RequestMetadata{WorkflowOwner: "0xdeadbeef", WorkflowDonID: zoneADonID}
		ctx := md.ContextWithCRE(t.Context())

		// Warm the cache while healthy: zone-a caller is allowed.
		require.NoError(t, z.enforce(ctx, zoneADonID))
		// Outage: cached (non-zone-b) membership keeps the read allowed.
		fake.fail.Store(true)
		require.NoError(t, z.enforce(ctx, zoneADonID))
	})

	t.Run("cached zone-b membership survives outage (still denied)", func(t *testing.T) {
		t.Parallel()
		z, fake := newOutageTestRestrictor(t, restrictEnabled)
		md := capabilities.RequestMetadata{WorkflowOwner: "0x" + allowlistedOwner, WorkflowDonID: zoneBDonID}
		ctx := md.ContextWithCRE(t.Context())

		// Warm the cache while healthy: zone-b + non-allowlisted owner is denied,
		// and the zone-b membership is cached.
		require.ErrorContains(t, z.enforce(ctx, zoneBDonID), "allowlisted workflow owners")

		// Outage: the cached zone-b membership must still be enforced. The caller
		// is denied via the allowlist, not allowed through by a resolve failure.
		fake.fail.Store(true)
		err := z.enforce(ctx, zoneBDonID)
		require.ErrorContains(t, err, "allowlisted workflow owners")
		require.NotContains(t, err.Error(), "could not resolve caller workflow DON")
	})

	t.Run("DON never resolved before outage fails closed", func(t *testing.T) {
		t.Parallel()
		z, fake := newOutageTestRestrictor(t, restrictEnabled)
		fake.fail.Store(true)
		md := capabilities.RequestMetadata{WorkflowOwner: "0x" + allowlistedOwner, WorkflowDonID: zoneBDonID}
		err := z.enforce(md.ContextWithCRE(t.Context()), zoneBDonID)
		require.ErrorContains(t, err, "could not resolve caller workflow DON")
	})
}
