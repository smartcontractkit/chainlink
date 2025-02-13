package changeset

import (
	"errors"
	"fmt"
	"strconv"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
)

type RegisterNodesInput struct {
	EnvLabel    string
	ProductName string
	DONs        DONConfigMap
}

type DONConfigMap map[string]DONConfig

type DONConfig struct {
	Name               string    `json:"name"`
	ChannelConfigStore string    `json:"channelConfigStore"`
	Verifier           string    `json:"verifier"`
	Configurator       string    `json:"configurator"`
	Nodes              []NodeCfg `json:"nodes"`
}

type NodeCfg struct {
	Name        string `json:"name"`
	CSAKey      string `json:"csa_key"`
	IsBootstrap bool   `json:"isBootstrap"`
}

// RegisterNodesWithJD registers each node from the config with the Job Distributor.
// It logs errors but continues to register remaining nodes even if some fail (we may revisit this in the future).
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

	for _, don := range cfg.DONs {
		for _, node := range don.Nodes {
			labels := append([]*ptypes.Label(nil), baseLabels...)
			isBootstrapStr := strconv.FormatBool(node.IsBootstrap)

			labels = append(labels, &ptypes.Label{
				Key:   "isBootstrap",
				Value: &isBootstrapStr,
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

	return deployment.ChangesetOutput{}, nil
}

func (cfg RegisterNodesInput) Validate() error {
	if cfg.EnvLabel == "" {
		return errors.New("EnvLabel must not be empty")
	}
	if cfg.ProductName == "" {
		return errors.New("ProductName must not be empty")
	}

	for donName, don := range cfg.DONs {
		if don.Name == "" {
			return fmt.Errorf("DON[%s] has empty Name", donName)
		}
		for _, node := range don.Nodes {
			if node.Name == "" {
				return fmt.Errorf("DON[%s] has node with empty Name", donName)
			}
			if node.CSAKey == "" {
				return fmt.Errorf("DON[%s] node %s has empty CSAKey", donName, node.Name)
			}
		}
	}
	return nil
}
