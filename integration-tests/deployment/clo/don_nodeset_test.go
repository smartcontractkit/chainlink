package clo_test

import (
	"encoding/json"
	"os"
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
	keystoneNops := filterNopNodes(nops.NodeOperators, keystoneFilter)
	ocr3Nops := filterNopNodes(nops.NodeOperators, ocr3Filter)
	workflowNops := filterNopNodes(nops.NodeOperators, workflowFilter)
	assert.Len(t, keystoneNops, 4)
	assert.Len(t, ocr3Nops, 14)
	assert.Len(t, workflowNops, 40)
	b, err := json.MarshalIndent(&keystoneNops, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile("testdata/keystone_nops.json", b, 0644))
	clo.NewDonEnv(clo.DonEnvConfig{
		DonName: keystone.WFDonName,
		Chains:  nil,
		Logger:  logger.TestLogger(t),
		Nops:    keystoneNops,
	})

}

func filterNopNodes(nops []models.NodeOperator, f filterFunc) []*models.NodeOperator {
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
			filterNop := nop
			filterNop.Nodes = res
			out = append(out, &filterNop)
		}
	}
	return out
}

func flattenToNodes(nops []*models.NodeOperator, f filterFunc) []*models.Node {
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

type filterFunc func(n *models.Node) bool

func keystoneFilter(n *models.Node) bool {
	for _, cat := range n.Categories {
		if cat.Name == "Keystone" {
			return true
		}
	}
	return false
}

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

func ptr[T any](a T) *T {
	return &a
}
