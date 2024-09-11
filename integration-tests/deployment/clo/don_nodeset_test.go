package clo_test

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
)

func TestDonNodeset(t *testing.T) {
	keystoneNops := loadTestNops(t, "testdata/keystone_nops.json")

	// this is hacky, but there is no first class concept of a chain writer node in CLO
	// in prod, probably better to make an explicit list of pubkeys if we can't add a category or product type
	// sufficient for testing
	writerFilter := func(n *models.Node) bool {
		return strings.Contains(n.Name, "ks-writer")
	}

	m := clo.CapabilityNodeSets(keystoneNops, map[string]clo.FilterFuncT[*models.Node]{
		"workflow":    clo.ProductFilterGenerator(models.ProductTypeWorkflow),
		"chainWriter": writerFilter,
	})
	assert.Len(t, m, 2)
	assert.Len(t, m["workflow"], 10)
	assert.Len(t, m["chainWriter"], 10)

	// can be used to derive the test data for the keystone deployment
	updateTestData := true
	if updateTestData {
		b, err := json.MarshalIndent(m["workflow"], "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile("testdata/workflow_nodes.json", b, 0644)) // nolint: gosec

		b, err = json.MarshalIndent(m["chainWriter"], "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile("testdata/chain_writer_nodes.json", b, 0644)) // nolint: gosec
	}
	gotWFNops := m["workflow"]
	sort.Slice(gotWFNops, func(i, j int) bool {
		return gotWFNops[i].ID < gotWFNops[j].ID
	})
	expectedWorkflowNops := loadTestNops(t, "testdata/workflow_nodes.json")
	assert.True(t, reflect.DeepEqual(gotWFNops, expectedWorkflowNops), "workflow nodes do not match")

	gotChainWriterNops := m["chainWriter"]
	sort.Slice(gotChainWriterNops, func(i, j int) bool {
		return gotChainWriterNops[i].ID < gotChainWriterNops[j].ID
	})
	expectedChainWriterNops := loadTestNops(t, "testdata/chain_writer_nodes.json")
	assert.True(t, reflect.DeepEqual(gotChainWriterNops, expectedChainWriterNops), "chain writer nodes do not match")
}

func loadTestNops(t *testing.T, path string) []*models.NodeOperator {
	f, err := os.ReadFile(path)
	require.NoError(t, err)
	var nodes []*models.NodeOperator
	require.NoError(t, json.Unmarshal(f, &nodes))
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})
	return nodes
}
