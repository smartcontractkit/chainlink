package standardcapabilities

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type CapDonIDResolver interface {
	awaitCapDonID(ctx context.Context) (uint32, error)
}

type staticCapDonIDResolver struct {
	donID uint32
}

func NewStaticCapDonIDResolver(donID uint32) CapDonIDResolver {
	return &staticCapDonIDResolver{donID: donID}
}

func (s *staticCapDonIDResolver) awaitCapDonID(ctx context.Context) (uint32, error) {
	return s.donID, nil
}

type dynamicCapDonIDResolver struct {
	capabilityID string
	registry     registry.CapabilitiesRegistry
	getPeerID    func() (p2ptypes.PeerID, error)
	lggr         logger.Logger

	timeout       time.Duration
	retryInterval time.Duration
}

func (d *dynamicCapDonIDResolver) awaitCapDonID(ctx context.Context) (uint32, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	for attempt := 1; ; attempt++ {
		donID, err := d.resolveCapDonID(ctx)
		if err == nil {
			d.lggr.Infow("Resolved capability DON ID from registry", "donID", donID, "attempts", attempt)
			return donID, nil
		}
		select {
		case <-ctx.Done():
			d.lggr.Errorw("Could not resolve capability DON ID; falling back to workflow DON ID for event labeling",
				"attempts", attempt, "err", err)
			return 0, nil
		case <-time.After(d.retryInterval):
		}
	}
}

func newDynamicCapDonIDResolver(capabilityID string, registry registry.CapabilitiesRegistry, getPeerID func() (p2ptypes.PeerID, error), lggr logger.Logger) *dynamicCapDonIDResolver {
	return &dynamicCapDonIDResolver{
		capabilityID:  capabilityID,
		registry:      registry,
		getPeerID:     getPeerID,
		lggr:          lggr,
		timeout:       defaultDonIDResolveTimeout,
		retryInterval: defaultDonIDRetryInterval,
	}
}

// resolveCapabilityDonID best-effort resolves the on-chain DON ID this node is
// running the given capability for, by looking up DONsForCapability and filtering
// by the local node's PeerID.
//
// It returns an error only for failures that may clear on retry, such as the
// peer ID or the registry not being available yet right after node start. A
// definitive answer that cannot be honored (no DON, or multiple DONs) returns
// 0 with a warning logged; the plugin then falls back to labeling events with
// the consumer workflow's DON ID. See CRE-4409.
func (d dynamicCapDonIDResolver) resolveCapDonID(ctx context.Context) (uint32, error) {
	peerID, err := d.getPeerID()
	if err != nil {
		return 0, fmt.Errorf("failed to get local peer ID: %w", err)
	}
	dwns, err := d.registry.DONsForCapability(ctx, d.capabilityID)
	if err != nil {
		return 0, fmt.Errorf("DONsForCapability failed: %w", err)
	}
	var matched []uint32
	for _, d := range dwns {
		for _, n := range d.Nodes {
			if n.PeerID != nil && *n.PeerID == peerID {
				matched = append(matched, d.DON.ID)
				break
			}
		}
	}
	switch len(matched) {
	case 1:
		return matched[0], nil
	case 0:
		d.lggr.Warnw("No DON found for local peer on capability; falling back to workflow DON ID for event labeling",
			"peerID", peerID, "capabilityID", d.capabilityID)
		return 0, nil
	default:
		d.lggr.Warnw("Local peer belongs to multiple DONs for capability; cannot disambiguate, falling back to workflow DON ID for event labeling",
			"peerID", peerID, "capabilityID", d.capabilityID, "matched", matched)
		return 0, nil
	}
}
