package operations

import (
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

// StandardCapabilityNodeFilter is the JD node selection of the standard capability proposal.
// jdDONName is the resolved JD DON name (JDDONName if set, else DONName).
func StandardCapabilityNodeFilter(jdDONName, environment, domain string, donFilters []offchain.TargetDONFilter) *nodev1.ListNodesRequest_Filter {
	nodeType := PluginNodeType
	filter := &nodev1.ListNodesRequest_Filter{
		Selectors: []*ptypes.Selector{
			{Key: "don-" + jdDONName, Op: ptypes.SelectorOp_EXIST},
			{Key: "environment", Op: ptypes.SelectorOp_EQ, Value: &environment},
			{Key: "product", Op: ptypes.SelectorOp_EQ, Value: &domain},
			{Key: "type", Op: ptypes.SelectorOp_EQ, Value: &nodeType},
		},
	}
	for _, f := range donFilters {
		filter = f.AddToFilter(filter)
	}
	for _, f := range donFilters {
		filter = f.AddToFilterIfNotPresent(filter)
	}

	return filter
}

// PluginDONNodeFilter is the JD node selection of the OCR3 and Ring proposals: the DON filters
// restricted to plugin nodes.
func PluginDONNodeFilter(donFilters []offchain.TargetDONFilter) *nodev1.ListNodesRequest_Filter {
	filter := &nodev1.ListNodesRequest_Filter{}
	for _, f := range donFilters {
		filter = f.AddToFilter(filter)
	}

	return offchain.TargetDONFilter{Key: "type", Value: PluginNodeType}.AddToFilter(filter)
}
