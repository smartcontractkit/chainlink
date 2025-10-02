package changeset

import (
	"errors"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// UpdatesNodesJDChangeset is a changeset that reads node info from a JSON file and updates node name and labels in Job Distributor
var UpdatesNodesJDChangeset = cldf.CreateChangeSet(updatesNodesJDLogic, updatesNodesJDLogicPrecondition)

func updatesNodesJDLogic(env cldf.Environment, c types.UpdateNodeConfig) (cldf.ChangesetOutput, error) {
	nodes := c.Nodes

	for _, node := range nodes {
		n, err := env.Offchain.GetNode(env.GetContext(), &nodev1.GetNodeRequest{
			Id: node.ID,
		})
		if err != nil {
			env.Logger.Errorw("failed to get node", "id", node.ID, "error", err)
			continue
		}
		nodeInfo := n.GetNode()

		var labels []*ptypes.Label
		if node.AppendLabels {
			currentLabels := nodeInfo.GetLabels()
			labels = append(labels, currentLabels...)
		}
		labels = append(labels, node.Labels...)

		_, err = env.Offchain.UpdateNode(env.GetContext(), &nodev1.UpdateNodeRequest{
			Id:        nodeInfo.GetId(),
			Name:      node.Name,
			PublicKey: nodeInfo.GetPublicKey(),
			Labels:    labels,
		})
		if err != nil {
			env.Logger.Errorw("failed to update node", "nodeName", nodeInfo.Name, "err", err)
		} else {
			env.Logger.Infof("node %s updated", nodeInfo.Name)
		}
	}

	return cldf.ChangesetOutput{}, nil
}

func updatesNodesJDLogicPrecondition(env cldf.Environment, c types.UpdateNodeConfig) error {
	if len(c.Nodes) == 0 {
		return errors.New("no nodes provided in the config")
	}

	return nil
}
