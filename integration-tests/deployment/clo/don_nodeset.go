package clo

import (
	"strings"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
)

func DonNodesets(nops []*models.NodeOperator, donFilters map[string]FilterFuncT[*models.Node]) map[string][]*models.NodeOperator {
	// first drop bootstraps if they exist
	nonBootstrapNops := filterNopNodes(nops, func(n *models.Node) bool {
		for _, chain := range n.ChainConfigs {
			if chain.Ocr2Config.IsBootstrap {
				return false
			}
		}
		return true
	})
	out := make(map[string][]*models.NodeOperator)
	for name, f := range donFilters {
		out[name] = filterNopNodes(nonBootstrapNops, f)
	}
	return out
}

// filterNopNodes filters the nodes of each nop by the provided filter function.
// if a nop has no nodes after filtering, it is not included in the output.
func filterNopNodes(nops []*models.NodeOperator, f FilterFuncT[*models.Node]) []*models.NodeOperator {
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

func categoryNameFilterGenerator(name string) FilterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		for _, cat := range n.Categories {
			if cat.Name == name {
				return true
			}
		}
		return false
	}
}

func publicKeyFilterGenerator(pubKey ...string) FilterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		if n.PublicKey == nil {
			return false
		}
		found := false
		for _, key := range pubKey {
			if *n.PublicKey == key {
				found = true
				break
			}
		}
		return found
	}
}

// this could be generalized to a regex filter
func NodeNameFilterGenerator(contains string) FilterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		return strings.Contains(n.Name, contains)
	}
}

// this is hacky
var chainWriterFilter = NodeNameFilterGenerator("Keystone Cap One")

func keystoneNopFilter(nop *models.NodeOperator) bool {
	nodeFilter := categoryNameFilterGenerator("Keystone")
	//isKeystoneNop := false
	for _, node := range nop.Nodes {
		if nodeFilter(node) {
			return true
		}
	}
	return false
}
