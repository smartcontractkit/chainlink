package offchain

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/node"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

const (
	FilterKeyDONName      = "don_name"
	FilterKeyCSAPublicKey = "csa_public_key"
)

type NodeLabelFilter struct {
	Key   string
	Value string
}

func (f NodeLabelFilter) AddToFilter(filter *nodev1.ListNodesRequest_Filter) *nodev1.ListNodesRequest_Filter {
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

func (f NodeLabelFilter) AddToFilterIfNotPresent(filter *nodev1.ListNodesRequest_Filter) *nodev1.ListNodesRequest_Filter {
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

func (f NodeLabelFilter) ToListFilter() *nodev1.ListNodesRequest_Filter {
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

type NodesSpecifier struct {
	// LabelFilters is a list of filters to apply when selecting the target DON nodes.
	// DEPRECATED: use NodeIDs and CSAKeys instead.
	LabelFilters []NodeLabelFilter
	// NodeIDs is a list of node IDs to target for the job spec.
	NodeIDs []string
	// CSAKeys is a list of CSA keys to target for the job spec.
	CSAKeys []string
}

var ErrNoSpecifiedMethods = fmt.Errorf("no method specified to select nodes")
var ErrMultipleSpecifiedMethods = fmt.Errorf("multiple methods specified to select nodes")

func (s NodesSpecifier) Validate() error {
	// one of don_filters, node_ids, or csa_keys
	if len(s.LabelFilters) == 0 && len(s.NodeIDs) == 0 && len(s.CSAKeys) == 0 {
		return fmt.Errorf("don_filters, node_ids, or csa_keys is required: %w", ErrNoSpecifiedMethods)
	}
	if len(s.NodeIDs) != 0 && (len(s.CSAKeys) != 0 || len(s.LabelFilters) != 0) {
		return fmt.Errorf("only one of node_ids, csa_keys, or don_filters can be provided: %w", ErrMultipleSpecifiedMethods)
	}
	return nil
}

func (s NodesSpecifier) Filter(donName, envName, domain string) *nodev1.ListNodesRequest_Filter {
	if err := s.Validate(); err != nil {
		return nil
	}
	if len(s.NodeIDs) != 0 {
		return &nodev1.ListNodesRequest_Filter{
			Ids: s.NodeIDs,
		}
	}
	if len(s.CSAKeys) != 0 {
		return &nodev1.ListNodesRequest_Filter{
			PublicKeys: s.CSAKeys,
		}
	}
	return labelFilter(domain, envName, s.LabelFilters)

}

func labelFilter(domain, envName string, filters []NodeLabelFilter) *nodev1.ListNodesRequest_Filter {
	filter := &nodev1.ListNodesRequest_Filter{
		Selectors: []*ptypes.Selector{
			{
				Key:   "product",
				Op:    ptypes.SelectorOp_EQ,
				Value: &domain,
			},
			{
				Key:   "environment",
				Op:    ptypes.SelectorOp_EQ,
				Value: &envName,
			},
		},
	}

	for _, f := range filters {
		filter = f.AddToFilter(filter)
	}
	return filter
}
