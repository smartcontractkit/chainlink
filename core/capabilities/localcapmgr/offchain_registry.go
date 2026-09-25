package localcapmgr

import (
	"context"
	"sort"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
)

// Divergence kinds recorded when the offchain capabilities registry disagrees with the
// on-chain registry. They are emitted as telemetry only: Phase 2 observes and cross-validates
// but does not yet let the offchain payload affect running capabilities (that gated cutover is
// Phase 4). "missing" means the on-chain registry has something the offchain payload lacks;
// "extra" means the offchain payload describes something not present on-chain for this node.
const (
	divergenceMissingDON        = "missing_don"
	divergenceMissingCapability = "missing_capability"
	divergenceExtraDON          = "extra_don"
	divergenceExtraCapability   = "extra_capability"
)

// offchainCrossCheck is the outcome of comparing the parsed offchain registry against the
// on-chain DON set this node belongs to.
type offchainCrossCheck struct {
	version       uint64
	divergences   map[string]int64 // kind -> count
	matchedCaps   int64
	comparedDONs  int64
	offchainEmpty bool // true when no offchain payload has been applied yet
}

// crossValidateOffchain compares the applied offchain capabilities registry against the
// on-chain DON set and emits telemetry. It never mutates desired state or blocks a capability;
// during the parallel-run (Phase 2) the on-chain registry remains authoritative and the
// offchain payload is observed only. It fails closed in the sense that a payload that cannot
// be parsed, or that diverges from on-chain, is never adopted — it is recorded and ignored.
//
// Only allowlisted capabilities are compared, matching buildDesiredState: those are the only
// capabilities this node would run, so they are the only ones whose config matters here.
func (m *localCapabilityManager) crossValidateOffchain(ctx context.Context, allMyDONs []registry.DON) {
	if m.offchainRegistry == nil {
		return // offchain registry not wired (e.g. feature off, or tests)
	}
	reg, version := m.offchainRegistry.LoadParsed()
	m.recordOffchainCheck(ctx, m.computeOffchainCrossCheck(reg, version, allMyDONs))
}

// computeOffchainCrossCheck compares a parsed offchain registry against the on-chain DON set
// and returns the divergence summary. It is pure (no metrics, no mutation) so it can be unit
// tested; crossValidateOffchain wires it to telemetry.
func (m *localCapabilityManager) computeOffchainCrossCheck(reg *capabilitiespb.OffchainCapabilitiesRegistry, version uint64, allMyDONs []registry.DON) offchainCrossCheck {
	check := offchainCrossCheck{version: version, divergences: map[string]int64{}}

	if reg == nil {
		check.offchainEmpty = true
		return check
	}

	offchainDONs := reg.GetDons()

	// on-chain -> offchain: every allowlisted (DON, capability) on-chain should be present
	// offchain.
	seenOffchainCap := map[uint32]map[string]bool{}
	for _, don := range allMyDONs {
		allowlisted := m.allowlistedCapIDs(don)
		if len(allowlisted) == 0 {
			continue
		}
		check.comparedDONs++

		offDON, ok := offchainDONs[don.ID]
		if !ok {
			check.divergences[divergenceMissingDON]++
			m.lggr.Warnw("Offchain registry missing on-chain DON", "donID", don.ID, "offchainVersion", version)
			continue
		}

		offCaps := offDON.GetCapabilityConfigs()
		seenOffchainCap[don.ID] = map[string]bool{}
		for _, capID := range allowlisted {
			if _, ok := offCaps[capID]; ok {
				check.matchedCaps++
				seenOffchainCap[don.ID][capID] = true
			} else {
				check.divergences[divergenceMissingCapability]++
				m.lggr.Warnw("Offchain registry missing on-chain capability",
					"donID", don.ID, "capID", capID, "offchainVersion", version)
			}
		}
	}

	// offchain -> on-chain: offchain DONs/capabilities not present on-chain for this node.
	onchainDONs := map[uint32]map[string]struct{}{}
	for _, don := range allMyDONs {
		caps := map[string]struct{}{}
		for capID := range don.CapabilityConfigurations {
			caps[capID] = struct{}{}
		}
		onchainDONs[don.ID] = caps
	}
	for donID, offDON := range offchainDONs {
		onCaps, ok := onchainDONs[donID]
		if !ok {
			check.divergences[divergenceExtraDON]++
			continue
		}
		for capID := range offDON.GetCapabilityConfigs() {
			if _, ok := onCaps[capID]; !ok {
				check.divergences[divergenceExtraCapability]++
			}
		}
	}

	return check
}

// allowlistedCapIDs returns the capability IDs configured on a DON that this node is
// allowlisted to launch, mirroring buildDesiredState's filter.
func (m *localCapabilityManager) allowlistedCapIDs(don registry.DON) []string {
	var out []string
	for capID := range don.CapabilityConfigurations {
		if m.localCfg != nil && m.localCfg.IsAllowlisted(capID) {
			out = append(out, capID)
		}
	}
	sort.Strings(out)
	return out
}
