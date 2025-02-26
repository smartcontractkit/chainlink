package changeset

import (
	"errors"
	"fmt"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
)

type NodeType int

const (
	NodeTypeOracle NodeType = iota
	NodeTypeBootstrap
)

func (nt NodeType) String() string {
	switch nt {
	case NodeTypeOracle:
		return "oracle"
	case NodeTypeBootstrap:
		return "bootstrap"
	default:
		return "unknown"
	}
}

type RegisterNodesInput struct {
	EnvLabel    string
	ProductName string
	// Will be deleted after migration to DONConfigMap
	DONs     DONConfigMap `json:"dons,omitempty"`
	DONsList []DONConfig  `json:"dons_list,omitempty"`
}

type DONConfigMap map[string]DONConfig

type DONConfig struct {
	ID                 int       `json:"id"`
	ChainSelector      string    `json:"chainSelector"`
	Name               string    `json:"name"`
	ChannelConfigStore string    `json:"channelConfigStore"`
	Verifier           string    `json:"verifier"`
	Configurator       string    `json:"configurator"`
	Nodes              []NodeCfg `json:"nodes"`
	BootstrapNodes     []NodeCfg `json:"bootstrapNodes"`
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

func registerNodesForDON(e deployment.Environment, donName string, donID int, nodes []NodeCfg, baseLabels []*ptypes.Label, nodeType NodeType) {
	ntStr := nodeType.String()
	for _, node := range nodes {
		labels := append([]*ptypes.Label(nil), baseLabels...)

		labels = append(labels, &ptypes.Label{
			Key:   "nodeType",
			Value: &ntStr,
		})

		labels = append(labels, &ptypes.Label{
			Key: fmt.Sprintf("don-%d-%s", donID, donName),
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

func RegisterNodesWithJD(e deployment.Environment, cfg RegisterNodesInput) (deployment.ChangesetOutput, error) {
	baseLabels := []*ptypes.Label{
		{
			Key:   "product",
			Value: &cfg.ProductName,
		},
		{
			Key:   "environment",
			Value: &cfg.EnvLabel,
		},
	}

	for _, don := range cfg.DONsList {
		registerNodesForDON(e, don.Name, don.ID, don.Nodes, baseLabels, NodeTypeOracle)
		registerNodesForDON(e, don.Name, don.ID, don.BootstrapNodes, baseLabels, NodeTypeBootstrap)
	}

	return deployment.ChangesetOutput{}, nil
}

func (cfg RegisterNodesInput) Validate() error {
	if cfg.EnvLabel == "" {
		return errors.New("EnvLabel must not be empty")
	}
	if cfg.ProductName == "" {
		return errors.New("ProductName must not be empty")
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
		if err := validateNodeSlice(don.BootstrapNodes, "bootstrap", i); err != nil {
			return err
		}
	}
	return nil
}
