package changeset

import (
	"errors"
	"fmt"
	"strconv"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
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

func registerNodesToJDLogic(env deployment.Environment, c types.NodeConfig) (deployment.ChangesetOutput, error) {
	dons, _ := LoadJSON[[]*DONConfigSchema](c.InputFileName, c.InputFS)

	for _, don := range dons {
		for _, node := range don.Nodes {
			n, err := env.Offchain.GetNode(env.GetContext(), &nodev1.GetNodeRequest{
				PublicKey: &node.CSAKey,
			})

			// base labels
			labels := []*ptypes.Label{
				&ptypes.Label{
					Key:   "product",
					Value: pointer.To(productLabel),
				},
				&ptypes.Label{
					Key:   "domain",
					Value: pointer.To(productLabel),
				},
				&ptypes.Label{
					Key:   productLabel,
					Value: pointer.To(""),
				},
				&ptypes.Label{
					Key:   "environment",
					Value: pointer.To(env.Name),
				},
				&ptypes.Label{
					Key:   "don_id",
					Value: pointer.To(strconv.Itoa(don.ID)),
				},
			}
			if node.IsBootstrap {
				labels = append(labels, &ptypes.Label{
					Key:   devenv.LabelNodeTypeKey,
					Value: pointer.To(devenv.LabelNodeTypeValueBootstrap),
				})
			} else {
				labels = append(labels, &ptypes.Label{
					Key:   devenv.LabelNodeTypeKey,
					Value: pointer.To(devenv.LabelNodeTypeValuePlugin),
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

	return deployment.ChangesetOutput{}, nil
}

func registerNodesToJDLogicPrecondition(env deployment.Environment, c types.NodeConfig) error {
	if c.InputFileName == "" {
		return errors.New("input file name is required")
	}

	_, err := LoadJSON[[]*DONConfigSchema](c.InputFileName, c.InputFS)
	if err != nil {
		return fmt.Errorf("failed to load don config input file: %w", err)
	}

	return nil
}
