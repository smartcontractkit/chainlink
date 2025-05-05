package changeset

import (
	"errors"
	"fmt"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// UpdatesNodesJDChangeset is a changeset that reads node info from a JSON file and updates node name and labels in Job Distributor
var UpdatesNodesJDChangeset = cldf.CreateChangeSet(updatesNodesJDLogic, updatesNodesJDLogicPrecondition)

type NodeConfigSchema struct {
	ID           string          `json:"id"`            // node id
	Name         string          `json:"name"`          // new node name
	Labels       []*ptypes.Label `json:"labels"`        // new labels
	AppendLabels bool            `json:"append_labels"` // if true, append new labels to existing labels, otherwise replace
}

func updatesNodesJDLogic(env deployment.Environment, c types.NodeConfig) (deployment.ChangesetOutput, error) {
	nodes, _ := LoadJSON[[]*NodeConfigSchema](c.InputFileName, c.InputFS)

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

	return deployment.ChangesetOutput{}, nil
}

func updatesNodesJDLogicPrecondition(env deployment.Environment, c types.NodeConfig) error {
	if c.InputFileName == "" {
		return errors.New("input file name is required")
	}

	_, err := LoadJSON[[]*NodeConfigSchema](c.InputFileName, c.InputFS)
	if err != nil {
		return fmt.Errorf("failed to load node config input file: %w", err)
	}

	return nil
}
