package clo_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func loadTestData(t *testing.T) *clo.GetNodeOperatorsResponse {
	f, err := os.ReadFile("testdata/nops.json")
	require.NoError(t, err)
	var nops clo.GetNodeOperatorsResponse
	require.NoError(t, json.Unmarshal(f, &nops))
	return &nops
}

func TestDonNodeset(t *testing.T) {
	nops := loadTestData(t)

	keystoneNops := filterNopNodes(nops.NodeOperators, categoryNameFilterGenerator("Keystone")) //filterNops(nops.NodeOperators, keystoneFilter)
	//ocr3Nops := filterNopNodes(nops.NodeOperators, ocr3Filter)
	assert.Len(t, keystoneNops, 10)
	//assert.Len(t, ocr3Nops, 14)
	b, err := json.MarshalIndent(&keystoneNops, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile("testdata/keystone_nops.json", b, 0644))

	workflowNodes := filterNopNodes(keystoneNops, productFilterGenerator(models.ProductTypeWorkflow))

	// this is hacky, but there is no first class concept of a chain writer node in CLO
	// in prod, probably better to make an explicit list of pubkeys if we can't add a category or product type
	assert.Len(t, workflowNodes, 10)

	b, err = json.MarshalIndent(&workflowNodes, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile("testdata/workflow_nodes.json", b, 0644))

	chainWriterNodes := filterNopNodes(keystoneNops, nodeNameFilterGenerator("Keystone Cap One"))
	assert.Len(t, chainWriterNodes, 10)

	b, err = json.MarshalIndent(&chainWriterNodes, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile("testdata/chain_writer_nodes.json", b, 0644))

	clo.NewDonEnv(clo.DonEnvConfig{
		DonName: keystone.WFDonName,
		Chains:  nil,
		Logger:  logger.TestLogger(t),
		Nops:    keystoneNops,
	})

}

// filterNops filters the input nops by the provided filter function.
func filterNops(nops []*models.NodeOperator, f filterFuncT[*models.NodeOperator]) []*models.NodeOperator {
	var out []*models.NodeOperator
	for _, nop := range nops {
		if f(nop) {
			nop := nop
			out = append(out, nop)
		}
	}
	return out
}

// filterNopNodes filters the nodes of each nop by the provided filter function.
// if a nop has no nodes after filtering, it is not included in the output.
func filterNopNodes(nops []*models.NodeOperator, f filterFuncT[*models.Node]) []*models.NodeOperator {
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

// flattenToNodes flattens the nodes of each nop into a single slice of nodes.
func flattenToNodes(nops []*models.NodeOperator, f filterFuncT[*models.Node]) []*models.Node {
	var out []*models.Node
	for _, nop := range nops {
		for _, n := range nop.Nodes {
			if f(n) {
				out = append(out, n)

			}
		}
	}
	return out
}

type filterFuncT[T any] func(n T) bool

// type filterFunc func(n *models.Node) bool
func keystoneFilter(n *models.Node) bool {
	for _, cat := range n.Categories {
		if cat.Name == "Keystone" {
			return true
		}
	}
	return false
}

func productFilterGenerator(p models.ProductType) filterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		for _, prod := range n.SupportedProducts {
			if prod == p {
				return true
			}
		}
		return false
	}
}

func categoryNameFilterGenerator(name string) filterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		for _, cat := range n.Categories {
			if cat.Name == name {
				return true
			}
		}
		return false
	}
}

func publicKeyFilterGenerator(pubKey ...string) filterFuncT[*models.Node] {
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
func nodeNameFilterGenerator(contains string) filterFuncT[*models.Node] {
	return func(n *models.Node) bool {
		return strings.Contains(n.Name, contains)
	}
}

// this is hacky
var chainWriterFilter = nodeNameFilterGenerator("Keystone Cap One")

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

/*
	func ocr3Filter(n *models.Node) bool {
		for _, p := range n.SupportedProducts {
			if p == models.ProductTypeOcr3Capability {
				return true
			}
		}
		return false
	}

	func workflowFilter(n *models.Node) bool {
		for _, p := range n.SupportedProducts {
			if p == models.ProductTypeWorkflow {
				return true
			}
		}
		return false
	}
*/
func ptr[T any](a T) *T {
	return &a
}
