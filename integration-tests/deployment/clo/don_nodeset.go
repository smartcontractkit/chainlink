package clo

import (
	"strings"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
)

// CapabilityNodeSets groups nodes by a given filter function, resulting in a map of don name to nodes.
func CapabilityNodeSets(nops []*models.NodeOperator, donFilters map[string]FilterFuncT[*models.Node]) map[string][]*models.NodeOperator {
	// first drop bootstraps if they exist because they do not serve capabilities
	nonBootstrapNops := FilterNopNodes(nops, func(n *models.Node) bool {
		for _, chain := range n.ChainConfigs {
			if chain.Ocr2Config.IsBootstrap {
				return false
			}
		}
		return true
	})
	// apply given filters to non-bootstrap nodes
	out := make(map[string][]*models.NodeOperator)
	for name, f := range donFilters {
		out[name] = FilterNopNodes(nonBootstrapNops, f)
	}
	return out
}

// FilterNopNodes filters the nodes of each nop by the provided filter function.
// if a nop has no nodes after filtering, it is not included in the output.
func FilterNopNodes(nops []*models.NodeOperator, f FilterFuncT[*models.Node]) []*models.NodeOperator {
	var out []*models.NodeOperator
	for _, nop := range nops {
		var res []*models.Node
		for _, n := range nop.Nodes {
			node := n
			if f(n) {
				res = append(res, node)
			}
		}
		if len(res) > 0 {
			filterNop := *nop
			filterNop.Nodes = res
			out = append(out, &filterNop)
		}
	}
	return out
}

type FilterFuncT[T any] func(n T) bool

func ProductFilterGenerator(p models.ProductType) FilterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		for _, prod := range n.SupportedProducts {
			if prod == p {
				return true
			}
		}
		return false
	}
}

// this could be generalized to a regex filter
func NodeNameFilterGenerator(contains string) FilterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		return strings.Contains(n.Name, contains)
	}
}
