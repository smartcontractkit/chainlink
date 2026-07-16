package changeset

import (
	"errors"
	"strconv"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

// RegisterNodesToJDChangeset is a changeset that reads node info from a JSON file and registers them in Job Distributor
// Register a node with a set of base labels and optionally with additional extra labels
var RegisterNodesToJDChangeset = cldf.CreateChangeSet(registerNodesToJDLogic, registerNodesToJDLogicPrecondition)

type MinimalNodeCfg struct {
	Name        string          `json:"name"`
	CSAKey      string          `json:"csa_key"`
	IsBootstrap bool            `json:"is_bootstrap"`
	Labels      []*ptypes.Label `json:"labels"`
}

type DONConfigSchema struct {
	ID    int              `json:"id"`
	Name  string           `json:"name"`
	Nodes []MinimalNodeCfg `json:"nodes"`
}

const productLabel = "data-feeds"

func registerNodesToJDLogic(env cldf.Environment, c types.RegisterNodeConfig) (cldf.ChangesetOutput, error) {
	dons := c.DONs

	for _, don := range dons {
		for _, node := range don.Nodes {
			n, err := env.Offchain.GetNode(env.GetContext(), &nodev1.GetNodeRequest{
				PublicKey: &node.CSAKey,
			})

			// base labels
			labels := []*ptypes.Label{
				{
					Key:   "product",
					Value: new(productLabel),
				},
				{
					Key:   "domain",
					Value: new(productLabel),
				},
				{
					Key:   productLabel,
					Value: new(""),
				},
				{
					Key:   "environment",
					Value: new(env.Name),
				},
				{
					Key:   "don_id",
					Value: new(strconv.Itoa(don.ID)),
				},
			}
			if node.IsBootstrap {
				labels = append(labels, &ptypes.Label{
					Key:   devenv.LabelNodeTypeKey,
					Value: new(devenv.LabelNodeTypeValueBootstrap),
				})
			} else {
				labels = append(labels, &ptypes.Label{
					Key:   devenv.LabelNodeTypeKey,
					Value: new(devenv.LabelNodeTypeValuePlugin),
				})
			}
			// extra labels
			labels = append(labels, node.Labels...)

			if err != nil {
				env.Logger.Infow("Node not found, attempting to register", "name", node.Name)
				newNode, err := env.Offchain.RegisterNode(env.GetContext(), &nodev1.RegisterNodeRequest{
					Name:      node.Name,
					PublicKey: node.CSAKey,
					Labels:    labels,
				})
				if err != nil {
					env.Logger.Errorw("failed to register node", "don", don.Name, "node", node.Name, "error", err)
				} else {
					env.Logger.Infow("registered node", "name", node.Name, "id", newNode.Node.Id)
				}
				continue
			}
			env.Logger.Infow("Node already registered, use UpdatesNodesJDChangeset to update node labels or name", "name", node.Name, "id", n.Node.Id)
		}
	}

	return cldf.ChangesetOutput{}, nil
}

func registerNodesToJDLogicPrecondition(env cldf.Environment, c types.RegisterNodeConfig) error {
	if len(c.DONs) == 0 {
		return errors.New("no DONs provided in the configuration")
	}

	return nil
}
