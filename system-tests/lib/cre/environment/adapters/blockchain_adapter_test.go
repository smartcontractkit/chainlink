package adapters

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
)

func TestBlockchainAdapterDescribeAndRewrite(t *testing.T) {
	adapter := NewBlockchainAdapter()
	out := &blockchain.Output{
		Nodes: []*blockchain.Node{
			{
				ExternalHTTPUrl: "http://10.0.0.10:8545",
				ExternalWSUrl:   "ws://10.0.0.10:8546",
			},
		},
	}

	refs, err := adapter.DescribeEndpoints("blockchain:0:anvil", out)
	if err != nil {
		t.Fatalf("expected describe to succeed: %v", err)
	}
	if len(refs) != 2 {
		t.Fatalf("expected two endpoint refs, got %d", len(refs))
	}

	bindings := []tunnel.TunnelBinding{
		{
			EndpointRef: tunnel.EndpointRef{EndpointName: "node-0-http"},
			LocalURL:    "http://127.0.0.1:18080",
		},
		{
			EndpointRef: tunnel.EndpointRef{EndpointName: "node-0-ws"},
			LocalURL:    "ws://127.0.0.1:18081",
		},
	}

	if err := adapter.RewriteWithBindings(out, bindings); err != nil {
		t.Fatalf("expected rewrite to succeed: %v", err)
	}

	if out.Nodes[0].ExternalHTTPUrl != "http://127.0.0.1:18080" {
		t.Fatalf("unexpected rewritten http url: %s", out.Nodes[0].ExternalHTTPUrl)
	}
	if out.Nodes[0].ExternalWSUrl != "ws://127.0.0.1:18081" {
		t.Fatalf("unexpected rewritten ws url: %s", out.Nodes[0].ExternalWSUrl)
	}
}
