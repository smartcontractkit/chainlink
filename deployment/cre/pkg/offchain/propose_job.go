package offchain

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
)

// ProposeJobRequest contains parameters for proposing a job to nodes in JD.
// Jobs are always proposed to nodes filtered at least by domain and environment.
// Additional node labels can be provided to further filter nodes.
type ProposeJobRequest struct {
	Job         string // toml
	Domain      string
	Environment string
	// labels to filter nodes by
	NodeLabels map[string]string
	PublicKeys []string // optional
	// labels to set on the new/updated job object
	JobLabels      map[string]string
	OffchainClient cldf_offchain.Client
	Lggr           logger.Logger

	ExtraSelectors []*ptypes.Selector // optional
}

func (r ProposeJobRequest) Validate() error {
	if r.Job == "" {
		return errors.New("job is empty")
	}
	// TODO validate valid toml
	if r.Domain == "" {
		return errors.New("domain is empty")
	}
	if r.Environment == "" {
		return errors.New("environment is empty")
	}
	if r.OffchainClient == nil {
		return errors.New("offchain client is nil")
	}
	if r.Lggr == nil {
		return errors.New("logger is nil")
	}

	return nil
}

func ProposeJob(ctx context.Context, req ProposeJobRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}
	var merr error
	// always filter at least by product and environment
	domainKey := req.Domain
	selectors := []*ptypes.Selector{
		{
			Key:   "product",
			Op:    ptypes.SelectorOp_EQ,
			Value: &domainKey,
		},
		{
			Key:   "environment",
			Op:    ptypes.SelectorOp_EQ,
			Value: &req.Environment,
		},
	}
	for key, value := range req.NodeLabels {
		selectors = append(selectors, &ptypes.Selector{
			Key:   key,
			Op:    ptypes.SelectorOp_EQ,
			Value: pointer.To(value),
		})
	}
	selectors = append(selectors, req.ExtraSelectors...)

	nodes, err := req.OffchainClient.ListNodes(ctx, &nodev1.ListNodesRequest{Filter: &nodev1.ListNodesRequest_Filter{
		Enabled:    1,
		Selectors:  selectors,
		PublicKeys: req.PublicKeys,
	}})
	if err != nil {
		return err
	}

	for _, node := range nodes.GetNodes() {
		_, err1 := req.OffchainClient.ProposeJob(ctx,
			&jobv1.ProposeJobRequest{
				NodeId: node.Id,
				Spec:   req.Job,
				Labels: convertLabels(req.JobLabels),
			})
		if err1 != nil {
			req.Lggr.Infow("Failed to propose job to node", "nodeId", node.Id, "nodeName", node.Name)
			merr = errors.Join(merr, fmt.Errorf("error proposing job to node %s spec %s : %w", node.Id, req.Job, err1))
		} else {
			req.Lggr.Infow("Successfully proposed job to node", "nodeId", node.Id, "nodeName", node.Name)
		}
	}

	return merr
}

func convertLabels(labels map[string]string) []*ptypes.Label {
	res := make([]*ptypes.Label, 0, len(labels))
	for k, v := range labels {
		newVal := v
		res = append(res, &ptypes.Label{
			Key:   k,
			Value: &newVal,
		})
	}

	return res
}
