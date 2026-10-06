package vault

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

// zoneBFamily is the DON family (in the capabilities registry) identifying the
// zone-b workflow DON whose vault GetSecrets reads are restricted to an
// allowlist of workflow owners. Any family starting with it (e.g. the
// "zone-b_workflows" family used by sharded workflow DONs) is treated as zone-b.
const zoneBFamily = "zone-b"

// zoneBRestrictor enforces that GetSecrets reads originating from a zone-b
// workflow DON are limited to an allowlist of workflow owners. Zone membership
// is resolved authoritatively from the node's own capabilities registry view,
// never from caller-supplied metadata.
type zoneBRestrictor struct {
	lggr                 logger.Logger
	capabilitiesRegistry registry.CapabilitiesRegistry
	// restrictEnabled is the master gate. When open, GetSecrets reads from a
	// zone-b workflow DON are restricted to allowlisted workflow owners.
	restrictEnabled limits.GateLimiter
	// ownerAllowed is the owner-scoped allowlist gate consulted only for zone-b
	// callers when the master gate is open.
	ownerAllowed limits.GateLimiter

	// zoneCacheMu guards zoneCache.
	zoneCacheMu sync.RWMutex
	// zoneCache holds the last successfully-resolved zone-b membership per
	// WorkflowDonID. It is the fallback when the capabilities registry view is
	// transiently unavailable, so a registry blip does not fail every vault read
	// DON-wide (see isZoneBWorkflowDON).
	zoneCache map[uint32]bool
}

func newZoneBRestrictor(lggr logger.Logger, limitsFactory limits.Factory, capabilitiesRegistry registry.CapabilitiesRegistry) (*zoneBRestrictor, error) {
	restrictEnabled, err := limits.MakeGateLimiter(limitsFactory, cresettings.Default.VaultZoneBWorkflowGetSecretsRestrictEnabled)
	if err != nil {
		return nil, fmt.Errorf("failed to create zone-b restrict gate limiter: %w", err)
	}
	ownerAllowed, err := limits.MakeGateLimiter(limitsFactory, cresettings.Default.PerOwner.VaultZoneBGetSecretsAllowed)
	if err != nil {
		return nil, fmt.Errorf("failed to create zone-b owner allowlist gate limiter: %w", err)
	}
	if restrictEnabled == nil || ownerAllowed == nil {
		return nil, errors.New("zone-b restrictor requires non-nil gate limiters")
	}
	return &zoneBRestrictor{
		lggr:                 logger.Named(lggr, "ZoneBRestrictor"),
		capabilitiesRegistry: capabilitiesRegistry,
		restrictEnabled:      restrictEnabled,
		ownerAllowed:         ownerAllowed,
		zoneCache:            make(map[uint32]bool),
	}, nil
}

// enforce denies GetSecrets reads originating from a zone-b workflow DON unless
// the calling workflow owner is allowlisted. It is a no-op unless the master
// gate (VaultZoneBWorkflowGetSecretsRestrictEnabled) is open and the caller
// resolves to a zone-b DON. The owner is read from ctx, which must already carry
// the (normalized) CRE owner via RequestMetadata.ContextWithCRE.
func (z *zoneBRestrictor) enforce(ctx context.Context, workflowDonID uint32) error {
	enabled, err := z.restrictEnabled.IsOpen(ctx)
	if err != nil {
		return fmt.Errorf("could not evaluate zone-b vault read restriction gate: %w", err)
	}
	if !enabled {
		return nil
	}

	isZoneB, err := z.isZoneBWorkflowDON(ctx, workflowDonID)
	if err != nil {
		// Fail closed: if we cannot authoritatively resolve the caller's zone, do
		// not proceed. The registry is in-process, so this only fires for an
		// unknown/unregistered WorkflowDonID.
		return err
	}
	if !isZoneB {
		return nil
	}

	if err := z.ownerAllowed.AllowErr(ctx); err != nil {
		return fmt.Errorf("zone-b workflow DON may only read vault secrets for allowlisted workflow owners: %w", err)
	}
	return nil
}

// isZoneBWorkflowDON authoritatively resolves the caller's WorkflowDonID against
// the node's own capabilities registry view and reports whether that DON is in
// the zone-b family. This does not trust any caller-supplied zone data.
func (z *zoneBRestrictor) isZoneBWorkflowDON(ctx context.Context, workflowDonID uint32) (bool, error) {
	don, err := z.capabilitiesRegistry.DONByID(ctx, workflowDonID)
	if err != nil {
		// The registry view can be transiently unavailable (e.g. not yet synced
		// after startup, or nil mid-update: DONByID returns "metadataRegistry
		// information not available"). That error is not specific to zone-b
		// callers, so failing closed here would block every vault GetSecrets read
		// DON-wide. Fall back to the last successfully-resolved membership for this
		// DON; only a never-before-resolved DON fails closed.
		if cached, ok := z.cachedZoneMembership(workflowDonID); ok {
			z.lggr.Warnw("capabilities registry lookup failed; using cached zone-b membership",
				"workflowDonID", workflowDonID, "isZoneB", cached, "err", err)
			return cached, nil
		}
		return false, fmt.Errorf("could not resolve caller workflow DON %d for zone-b vault read restriction: %w", workflowDonID, err)
	}
	isZoneB := slices.ContainsFunc(don.Families, isZoneBFamily)
	z.storeZoneMembership(workflowDonID, isZoneB)
	return isZoneB, nil
}

// isZoneBFamily reports whether family identifies a zone-b workflow DON, i.e.
// starts with zoneBFamily. This covers the base "zone-b" family as well as
// sharding-specific variants such as "zone-b_workflows". Each workflow shard is
// its own DON in the registry, so every shard of a zone-b workflow DON is
// resolved and restricted by its own WorkflowDonID. A prefix match fails
// closed: over-matching only restricts more DONs, never fewer. The match is
// case-insensitive since family casing may vary across registry sources.
func isZoneBFamily(family string) bool {
	return strings.HasPrefix(strings.ToLower(family), zoneBFamily)
}

func (z *zoneBRestrictor) cachedZoneMembership(workflowDonID uint32) (bool, bool) {
	z.zoneCacheMu.RLock()
	defer z.zoneCacheMu.RUnlock()
	isZoneB, ok := z.zoneCache[workflowDonID]
	return isZoneB, ok
}

func (z *zoneBRestrictor) storeZoneMembership(workflowDonID uint32, isZoneB bool) {
	z.zoneCacheMu.Lock()
	defer z.zoneCacheMu.Unlock()
	z.zoneCache[workflowDonID] = isZoneB
}

func (z *zoneBRestrictor) close() error {
	var err error
	if cerr := z.restrictEnabled.Close(); cerr != nil {
		err = errors.Join(err, fmt.Errorf("error closing zone-b restrict gate limiter: %w", cerr))
	}
	if cerr := z.ownerAllowed.Close(); cerr != nil {
		err = errors.Join(err, fmt.Errorf("error closing zone-b owner allowlist gate limiter: %w", cerr))
	}
	return err
}
