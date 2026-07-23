package onchain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func TestChartNodeNameForWorker(t *testing.T) {
	t.Parallel()

	p2pKey, err := crypto.NewP2PKey("dev-password")
	require.NoError(t, err)
	peer := strings.TrimPrefix(p2pKey.PeerID.String(), "p2p_")

	worker := &cre.NodeMetadata{Host: "internal-host", Keys: &secrets.NodeKeys{P2PKey: p2pKey}}
	runtime := map[string]domain.NodeRuntimeInfo{
		"node-0": {PeerID: peer},
	}
	require.Equal(t, "node-0", chartNodeNameForWorker(worker, runtime))
}

func TestJdChainConfigNodes_IncludesBootstrap(t *testing.T) {
	t.Parallel()

	topology := testWorkflowBootstrapTopology(t)
	nodes := jdChainConfigNodes(topology, nil)
	require.Len(t, nodes, 5)

	workers, bootstraps := 0, 0
	for _, n := range nodes {
		switch {
		case n.nodeMeta.HasRole(cre.BootstrapNode):
			bootstraps++
			require.Equal(t, "bootstrap", n.donName)
		case n.nodeMeta.HasRole(cre.WorkerNode):
			workers++
			require.Equal(t, "workflow", n.donName)
		}
	}
	require.Equal(t, 4, workers)
	require.Equal(t, 1, bootstraps)
}
