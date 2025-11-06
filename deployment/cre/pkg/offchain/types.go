package offchain

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/node"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

const (
	FilterKeyDONName      = "don_name"
	FilterKeyCSAPublicKey = "csa_public_key"
)

type TargetDONFilter struct {
	Key   string
	Value string
}

func (f TargetDONFilter) AddToFilter(filter *nodev1.ListNodesRequest_Filter) *nodev1.ListNodesRequest_Filter {
	switch f.Key {
	case FilterKeyDONName:
		filter.Selectors = append(filter.Selectors, &ptypes.Selector{
			Op:  ptypes.SelectorOp_EXIST,
			Key: "don-" + f.Value,
		})
	case FilterKeyCSAPublicKey:
		filter.PublicKeys = append(filter.PublicKeys, f.Value)
	default:
		filter.Selectors = append(filter.Selectors, &ptypes.Selector{
			Op:    ptypes.SelectorOp_EQ,
			Key:   f.Key,
			Value: &f.Value,
		})
	}
	return filter
}

func (f TargetDONFilter) AddToFilterIfNotPresent(filter *nodev1.ListNodesRequest_Filter) *nodev1.ListNodesRequest_Filter {
	switch f.Key {
	case FilterKeyDONName:
		for _, s := range filter.Selectors {
			if s.Key == "don-"+f.Value {
				return filter
			}
		}
	case FilterKeyCSAPublicKey:
		for _, pk := range filter.PublicKeys {
			if pk == f.Value {
				return filter
			}
		}
	default:
		for _, s := range filter.Selectors {
			if s.Key == f.Key {
				return filter
			}
		}
	}
	return f.AddToFilter(filter)
}

func (f TargetDONFilter) ToListFilter() *nodev1.ListNodesRequest_Filter {
	filter := &nodev1.ListNodesRequest_Filter{}
	return f.AddToFilter(filter)
}

type NodeCfg struct {
	node.MinimalNodeCfg `yaml:",inline"`
	P2PID               string `json:"p2p_id" yaml:"p2p_id"`
	Zone                string `json:"zone" yaml:"zone"`
}

type DONConfig struct {
	ID             int                          `json:"don_id" yaml:"don_id"`
	Name           string                       `json:"don_name" yaml:"don_name"`
	F              uint8                        `json:"f" yaml:"f"`
	Nodes          []NodeCfg                    `json:"nodes" yaml:"nodes"`
	BootstrapNodes []string                     `json:"bootstrap_nodes,omitempty" yaml:"bootstrap_nodes,omitempty"`
	Capabilities   []contracts.CapabilityConfig `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
}
