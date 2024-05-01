package capreg

import (
	"context"
	"fmt"

	"go.uber.org/multierr"
)

var _ Local = (*capabilityLauncher)(nil)

type capabilityLauncher struct {
	myP2PID        []byte
	capabilitySvcs map[CapabilityID]CapabilityService
	myDONs         map[uint32]DON
}

func NewCapabilityLauncher(myP2PID []byte) *capabilityLauncher {
	return &capabilityLauncher{
		myP2PID:        myP2PID,
		capabilitySvcs: make(map[CapabilityID]CapabilityService),
		myDONs:         make(map[uint32]DON),
	}
}

// Sync implements Local.
// For the launcher, we need to filter out state that is not relevant to this node.
// This is done by checking if this node is a member of any new DONs using the p2p ID.
func (l *capabilityLauncher) Sync(ctx context.Context, s State) error {
	var errs error

	// filter out DONs that this node is not a member of.
	relevantDONs := filterRelevantDONs(l.myP2PID, s)

	// This shouldn't happen - nodes should only be added to a DON if they support the capability
	// of that DON.
	errs = multierr.Append(errs, l.handleDeletedDONs(ctx, relevantDONs, errs))

	errs = multierr.Append(errs, l.handleNewAndUpdatedDONs(ctx, s, errs))

	return errs
}

func (l *capabilityLauncher) handleNewAndUpdatedDONs(ctx context.Context, s State, errs error) error {
	for _, don := range s.DONs {
		myDON, ok := l.myDONs[don.ID]
		if ok {
			// If there has been a change to the Nodes field, we probably want to restart all capability services.
			if nodesChanged(myDON.Nodes, don.Nodes) {
				l.myDONs[don.ID] = don
				for _, cc := range don.CapabilityConfigurations {
					capSvc, ok2 := l.capabilitySvcs[cc.CapabilityID]
					if !ok2 {
						errs = multierr.Append(errs,
							fmt.Errorf("restart: capability service not found for capability ID %s", cc.CapabilityID))
						continue
					}

					if err := capSvc.Update(ctx, don); err != nil {
						errs = multierr.Append(errs,
							fmt.Errorf("failed to update capability service %s with updated DON configuration (new nodes property): %w, DON: %+v", cc.CapabilityID, err, don))
					}
				}
				continue
			}

			// If there has only been a change to a subset of the capability configurations, we only want to update
			// the capability services that are affected.
			removed, newOrUpdated := capabilityDiff(l.myDONs[don.ID].CapabilityConfigurations, don.CapabilityConfigurations)
			for _, capabilityID := range newOrUpdated {
				capSvc, ok2 := l.capabilitySvcs[capabilityID]
				if !ok2 {
					errs = multierr.Append(errs, fmt.Errorf("update: capability service not found for capability ID %s", capabilityID))
					continue
				}

				if err := capSvc.Update(ctx, don); err != nil {
					errs = multierr.Append(errs,
						fmt.Errorf("failed to update capability service %s with new DON configuration: %w, DON: %+v",
							capabilityID, err, don))
				}
			}
			for _, capabilityID := range removed {
				capSvc, ok2 := l.capabilitySvcs[capabilityID]
				if !ok2 {
					errs = multierr.Append(errs, fmt.Errorf("update: capability service not found for capability ID %s", capabilityID))
					continue
				}

				if err := capSvc.Stop(ctx, don); err != nil {
					errs = multierr.Append(errs,
						fmt.Errorf("failed to stop capability service %s with updated DON configuration: %w, DON: %+v",
							capabilityID, err, don))
				}
			}

			l.myDONs[don.ID] = don
		} else {
			// New DON, start capability services specified by it.
			l.myDONs[don.ID] = don
			for _, cc := range don.CapabilityConfigurations {
				capSvc, ok2 := l.capabilitySvcs[cc.CapabilityID]
				if !ok2 {
					// This shouldn't happen - nodes should only be added to a DON if they support the capability
					// of that DON.
					errs = multierr.Append(errs, fmt.Errorf("start: capability service not found for capability ID %s", cc.CapabilityID))
					continue
				}

				if err := capSvc.Start(ctx, don); err != nil {
					errs = multierr.Append(errs,
						fmt.Errorf("failed to start capability service %s with new DON configuration: %w, DON: %+v",
							cc.CapabilityID, err, don))
				}
			}
		}
	}
	return errs
}

func (l *capabilityLauncher) handleDeletedDONs(ctx context.Context, relevantDONs map[uint32]DON, errs error) error {
	for _, don := range l.myDONs {
		_, found := relevantDONs[don.ID]
		if !found {
			// DON present in our state and not onchain state == DON has been deleted from the onchain state.
			// Spin down capability services for this particular DON.
			for _, cc := range don.CapabilityConfigurations {
				capSvc, ok := l.capabilitySvcs[cc.CapabilityID]
				if !ok {
					errs = multierr.Append(errs, fmt.Errorf("delete: capability service not found for capability ID %s", cc.CapabilityID))
					continue
				}

				if err := capSvc.Stop(ctx, don); err != nil {
					errs = multierr.Append(errs,
						fmt.Errorf("failed to stop capability service %s with deleted DON configuration: %w, DON: %+v",
							cc.CapabilityID, err, don))
				}
			}
			delete(l.myDONs, don.ID)
		}
	}
	return errs
}

// RegisterCapabilityService registers the provided capability service with the launcher.
func (l *capabilityLauncher) RegisterCapabilityService(c CapabilityService) error {
	if c == nil {
		return fmt.Errorf("cannot register nil capability")
	}

	if c.CapabilityID() == "" {
		return fmt.Errorf("cannot register capability with empty ID")
	}

	if _, exists := l.capabilitySvcs[c.CapabilityID()]; exists {
		return fmt.Errorf("capability with ID %s already registered", c.CapabilityID())
	}

	l.capabilitySvcs[c.CapabilityID()] = c
	return nil
}

// Close implements Local.
// For the launcher, we need to stop all capability services.
func (l *capabilityLauncher) Close() error {
	var errs error
	for _, svc := range l.capabilitySvcs {
		if err := svc.Close(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}
