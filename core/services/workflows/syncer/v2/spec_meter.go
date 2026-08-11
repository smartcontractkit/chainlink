package v2

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/resourcemanager"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	meteringpb "github.com/smartcontractkit/chainlink-protos/metering/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// Service-level metering identity constants for the workflow syncer. Service is
// the stable service constant; ResourcePool identifies the workflow_specs_v2
// pool. The billing unit (ResourceType) and per-resource id (ResourceID) are
// carried on each Utilization. The coarse deployment/node/DON dimensions
// (product, tenant, environment, zone, don_id, node_id) are supplied at
// construction via NewSpecMeter.
const (
	meterService      = "workflow-syncer-v2"
	meterResourcePool = "workflow_specs_v2"
	meterResourceType = "operations"
)

// SpecLister is the read-only view of durable workflow-spec storage the
// metering snapshot path needs: the identity columns of every persisted spec.
type SpecLister interface {
	ListWorkflowSpecs(ctx context.Context) ([]*job.WorkflowSpec, error)
}

// specMeter wraps ResourceManager lifecycle and metering identity concerns
// to stop leaking complexity into the event handler
type SpecMeter struct {
	services.Service
	eng *services.Engine

	rm       *resourcemanager.ResourceManager
	identity resourcemanager.ResourceIdentity

	// resolvedDonID holds the workflow DON id once set via SetWorkflowDon.
	resolvedDonID atomic.Pointer[string]

	// unregister removes this meter from the ResourceManager's snapshot
	// registry; set in start, called in close. Nil until started.
	unregister func()

	orgResolver orgresolver.OrgResolver
	specs       SpecLister
}

// NewSpecMeter returns a SpecMeter emitting through rm with the given coarse
// identity (product/tenant/environment/zone/node_id from node metering
// config). The syncer Service/ResourcePool constants are stamped here and
// always overwrite whatever the caller passes; the workflow DON id is supplied
// later via SetWorkflowDon. specs backs the snapshot enumeration. orgResolver
// may be nil (org_id is then always empty).
func NewSpecMeter(
	lggr logger.Logger,
	rm *resourcemanager.ResourceManager,
	identity resourcemanager.ResourceIdentity,
	specs SpecLister,
	orgResolver orgresolver.OrgResolver,
) (*SpecMeter, error) {
	if rm == nil {
		return nil, errors.New("resource manager must be provided")
	}
	if specs == nil {
		return nil, errors.New("spec lister must be provided")
	}
	identity.Service = meterService
	identity.ResourcePool = meterResourcePool

	sm := &SpecMeter{
		rm:          rm,
		identity:    identity,
		orgResolver: orgResolver,
		specs:       specs,
	}
	sm.Service, sm.eng = services.Config{
		Name:  "SpecMeter",
		Start: sm.start,
		Close: sm.close,
	}.NewServiceEngine(lggr)
	return sm, nil
}

func (sm *SpecMeter) start(ctx context.Context) error {
	if err := sm.rm.Start(ctx); err != nil {
		return err
	}
	sm.unregister = sm.rm.Register(sm)
	return nil
}

func (sm *SpecMeter) close() error {
	// Stop snapshotting this meter, then close the ResourceManager service.
	if sm.unregister != nil {
		sm.unregister()
		sm.unregister = nil
	}
	return sm.rm.Close()
}

// SetWorkflowDon supplies the launcher-resolved workflow DON identity. Called
// by the registry after WaitForDon, before any event is dispatched; the value
// is static for the life of the node.
func (sm *SpecMeter) SetWorkflowDon(don commoncap.DON) {
	if sm == nil {
		return
	}
	donID := strconv.FormatUint(uint64(don.ID), 10)
	sm.resolvedDonID.Store(&donID)
}

// EmitSpecDelta emits one metering.v1.MeterRecord (METER_ACTION_UPDATE)
// capturing a signed ±delta to the durable workflow_specs_v2 level for
// workflowID
func (sm *SpecMeter) EmitSpecDelta(ctx context.Context, delta int64, workflowID, owner, eventID string) {
	if sm == nil {
		return
	}
	orgID := contexts.CREValue(ctx).Org
	if orgID == "" && sm.orgResolver != nil && owner != "" {
		if resolved, err := sm.orgResolver.Get(ctx, owner); err != nil {
			sm.eng.Warnw("failed to resolve org ID for metering", "owner", owner, "err", err)
		} else {
			orgID = resolved
		}
	}
	// resource_id = workflow_id (the syncer meters one durable spec per workflow).
	sm.rm.EmitDelta(ctx, sm.baseIdentity(), eventID, delta, resourcemanager.UtilizationFields{
		ResourceType: meterResourceType,
		ResourceID:   workflowID,
		OrgID:        orgID,
	})
}

// baseIdentity returns the meter's identity with the workflow DON id folded in
// once it has been set via SetWorkflowDon.
func (sm *SpecMeter) baseIdentity() resourcemanager.ResourceIdentity {
	id := sm.identity
	if donID := sm.resolvedDonID.Load(); donID != nil {
		nodeID := ""
		if id.Don != nil {
			nodeID = id.Don.NodeID
		}
		id.Don = &resourcemanager.DonIdentity{
			DonID:  *donID,
			NodeID: nodeID,
		}
	}
	return id
}

// GetUtilization implements resourcemanager.Meterable: it returns one
// SnapshotEntry per persisted workflow_specs_v2 spec. The durable resource is
// the stored spec, NOT the running engine.
//
// resource_id is the workflow_id; the organization is resolved fail-open from
// the spec's stored owner through the org resolver. Must fail open.
func (sm *SpecMeter) GetUtilization(ctx context.Context) []resourcemanager.SnapshotEntry {
	specs, err := sm.specs.ListWorkflowSpecs(ctx)
	if err != nil {
		sm.eng.Warnw("failed to list persisted workflow specs for metering snapshot; skipping tick", "err", err)
		return nil
	}
	base := sm.baseIdentity()
	entries := make([]resourcemanager.SnapshotEntry, 0, len(specs))
	for _, spec := range specs {
		if spec == nil {
			continue
		}
		var orgID string
		if sm.orgResolver != nil && spec.WorkflowOwner != "" {
			if resolved, err := sm.orgResolver.Get(ctx, spec.WorkflowOwner); err != nil {
				sm.eng.Warnw("failed to resolve org ID for metering snapshot", "owner", spec.WorkflowOwner, "err", err)
			} else {
				orgID = resolved
			}
		}
		entries = append(entries, resourcemanager.SnapshotEntry{
			Identity: base,
			Utilizations: []*meteringpb.Utilization{
				resourcemanager.NewUtilizationInt(1, resourcemanager.UtilizationFields{
					ResourceType: meterResourceType,
					ResourceID:   spec.WorkflowID,
					OrgID:        orgID,
				}),
			},
		})
	}
	return entries
}
