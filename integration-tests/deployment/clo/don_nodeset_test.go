package clo_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
)

// this is hacky, but there is no first class concept of a chain writer node in CLO
// in prod, probably better to make an explicit list of pubkeys if we can't add a category or product type
// sufficient for testing
var (
	writerFilter = func(n *models.Node) bool {
		return strings.Contains(n.Name, "Cap One") && !strings.Contains(n.Name, "Boot")
	}

	assetFilter = func(n *models.Node) bool {
		return strings.Contains(n.Name, "DF 2 stage testnet") && !strings.Contains(n.Name, "Bootstrap")
	}

	wfFilter = func(n *models.Node) bool {
		return strings.Contains(n.Name, "Keystone One") && !strings.Contains(n.Name, "Boot")
	}

	aptosFilter = func(n *models.Node) bool {
		return strings.Contains(n.Name, "Keystone Aptos")
	}
)

func TestGenerateNopNodesData(t *testing.T) {
	//t.Skipf("this test is for generating test data only")
	// use for generating keystone deployment test data
	// `./bin/fmscli --config ~/.fmsclient/prod.yaml login`
	// `./bin/fmscli --config ~/.fmsclient/prod.yaml get nodeOperators > /tmp/all-clo-nops.json`

	regenerateFromCLO := true
	d := "/tmp"
	target := filepath.Join(d, "keystone_nops.json")
	if regenerateFromCLO {
		path := "/tmp/clo-staging-nops.json"

		f, err := os.ReadFile(path)
		require.NoError(t, err)
		type cloData struct {
			Nops []*models.NodeOperator `json:"nodeOperators"`
		}
		var d cloData
		require.NoError(t, json.Unmarshal(f, &d))
		require.NotEmpty(t, d.Nops)
		allNops := d.Nops
		sort.Slice(allNops, func(i, j int) bool {
			return allNops[i].ID < allNops[j].ID
		})

		ksFilter := func(n *models.Node) bool {
			return writerFilter(n) || assetFilter(n) || wfFilter(n) || aptosFilter(n)
		}
		ksNops := clo.FilterNopNodes(allNops, ksFilter)
		require.NotEmpty(t, ksNops)
		b, err := json.MarshalIndent(ksNops, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(target, b, 0644)) // nolint: gosec
	}
	keystoneNops := loadTestNops(t, target)

	m := clo.CapabilityNodeSets(keystoneNops, map[string]clo.FilterFuncT[*models.Node]{
		"workflow":    wfFilter,
		"chainWriter": writerFilter,
		"asset":       assetFilter,
		"aptos":       aptosFilter,
	})
	assert.Len(t, m, 4)
	assert.Len(t, m["workflow"], 7)
	assert.Len(t, m["chainWriter"], 4)
	assert.Len(t, m["asset"], 7)
	assert.Len(t, m["aptos"], 3)

	// can be used to derive the test data for the keystone deployment
	updateTestData := true
	if updateTestData {
		d := "/tmp" // change this to the path where you want to write the test, "../deployment/keystone/testdata"
		b, err := json.MarshalIndent(m["workflow"], "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(d, "workflow_nodes.json"), b, 0600))

		b, err = json.MarshalIndent(m["chainWriter"], "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(d, "chain_writer_nodes.json"), b, 0600))
		b, err = json.MarshalIndent(m["asset"], "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(d, "asset_nodes.json"), b, 0600))
		b, err = json.MarshalIndent(m["aptos"], "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(d, "aptos_nodes.json"), b, 0600))
	}
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
