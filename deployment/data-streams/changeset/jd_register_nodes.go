package changeset

import (
	"errors"
	"fmt"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

var RegisterNodesWithJDChangeset = deployment.CreateChangeSet(registerNodesWithJDLogic, registerNodesWithJDPrecondition)

type RegisterNodesInput struct {
	BaseLabels map[string]string
	DONsList   []DONConfig `json:"dons_list,omitempty"`
}

type DONConfig struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	Nodes          []NodeCfg `json:"nodes"`
	BootstrapNodes []NodeCfg `json:"bootstrapNodes"`
}

type NodeCfg struct {
	Name        string `json:"name"`
	CSAKey      string `json:"csa_key"`
	IsBootstrap bool   `json:"isBootstrap,omitempty"`
}

func validateNodeSlice(nodes []NodeCfg, nodeType string, donIndex int) error {
	for _, node := range nodes {
		if node.Name == "" {
			return fmt.Errorf("DON[%d] has %s node with empty Name", donIndex, nodeType)
		}
		if node.CSAKey == "" {
			return fmt.Errorf("DON[%d] %s node %s has empty CSAKey", donIndex, nodeType, node.Name)
		}
	}
	return nil
}

func registerNodesForDON(e deployment.Environment, donName string, donID uint64, nodes []NodeCfg, baseLabels []*ptypes.Label, nodeType jd.NodeType) {
	ntStr := nodeType.String()
	for _, node := range nodes {
		labels := append([]*ptypes.Label(nil), baseLabels...)

		labels = append(labels, &ptypes.Label{
			Key:   "nodeType",
			Value: &ntStr,
		})

		labels = append(labels, &ptypes.Label{
			Key: utils.DonIdentifier(donID, donName),
		})

		nodeID, err := e.Offchain.RegisterNode(e.GetContext(), &nodev1.RegisterNodeRequest{
			Name:      node.Name,
			PublicKey: node.CSAKey,
			Labels:    labels,
		})
		if err != nil {
			e.Logger.Errorw("failed to register node", "node", node.Name, "error", err)
		} else {
			e.Logger.Infow("registered node", "name", node.Name, "id", nodeID)
		}
	}
}

func registerNodesWithJDLogic(e deployment.Environment, cfg RegisterNodesInput) (deployment.ChangesetOutput, error) {
	baseLabels := make([]*ptypes.Label, 0, len(cfg.BaseLabels)+1)
	for key, value := range cfg.BaseLabels {
		baseLabels = append(baseLabels, &ptypes.Label{
			Key:   key,
			Value: &value,
		})
	}

	for _, don := range cfg.DONsList {
		registerNodesForDON(e, don.Name, don.ID, don.Nodes, baseLabels, jd.NodeTypeOracle)
		registerNodesForDON(e, don.Name, don.ID, don.BootstrapNodes, baseLabels, jd.NodeTypeBootstrap)
	}

	return deployment.ChangesetOutput{}, nil
}

func (cfg RegisterNodesInput) Validate() error {
	for key, value := range cfg.BaseLabels {
		if key == "" || value == "" {
			return errors.New("common node labels have empty key or value")
		}
	}
	for i, don := range cfg.DONsList {
		if don.Name == "" {
			return fmt.Errorf("DON[%d] has empty Name", i)
		}
		if err := validateNodeSlice(don.Nodes, "node", i); err != nil {
			return err
		}
		if len(don.BootstrapNodes) == 0 {
			return fmt.Errorf("DON[%d] has no bootstrap nodes", i)
		}
		if err := validateNodeSlice(don.BootstrapNodes, jd.NodeTypeBootstrap.String(), i); err != nil {
			return err
		}
	}
	return nil
}

func registerNodesWithJDPrecondition(_ deployment.Environment, cfg RegisterNodesInput) error {
	return cfg.Validate()
}
