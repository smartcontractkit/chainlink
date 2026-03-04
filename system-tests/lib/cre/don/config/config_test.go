package config

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func TestResolveGatewayConnectorURL_PlacementMatrix(t *testing.T) {
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	tests := []struct {
		name            string
		callerPlacement string
		targetPlacement string
		wantURL         string
	}{
		{
			name:            "local caller local target uses internal",
			callerPlacement: "local",
			targetPlacement: "local",
			wantURL:         "ws://bootstrap-gateway-node0:5003/node",
		},
		{
			name:            "local caller remote target uses external ec2",
			callerPlacement: "local",
			targetPlacement: "remote",
			wantURL:         "ws://203.0.113.10:5003/node",
		},
		{
			name:            "remote caller local target uses docker host external",
			callerPlacement: "remote",
			targetPlacement: "local",
			wantURL:         "ws://" + strings.TrimPrefix(framework.HostDockerInternal(), "http://") + ":5003/node",
		},
		{
			name:            "remote caller remote target uses internal",
			callerPlacement: "remote",
			targetPlacement: "remote",
			wantURL:         "ws://bootstrap-gateway-node0:5003/node",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topology, gateway := mustBuildGatewayTopology(t, tt.targetPlacement)

			gotURL, err := resolveGatewayConnectorURL(tt.callerPlacement, topology, gateway, "")
			require.NoError(t, err, "resolveGatewayConnectorURL should not fail")
			require.Equal(t, tt.wantURL, gotURL, "unexpected gateway connector URL")
		})
	}
}

func TestResolveGatewayConnectorURL_RemoteHostOverride(t *testing.T) {
	topology, gateway := mustBuildGatewayTopology(t, "remote")
	gotURL, err := resolveGatewayConnectorURL("local", topology, gateway, "203.0.113.22")
	require.NoError(t, err, "resolveGatewayConnectorURL should use explicit remote host override")
	require.Equal(t, "ws://203.0.113.22:5003/node", gotURL, "unexpected gateway connector URL")
}

func TestResolveNodeFacingBootstrapAddress_PlacementMatrix(t *testing.T) {
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	tests := []struct {
		name               string
		callerPlacement    string
		bootstrapPlacement string
		bootstrapHost      string
		internalPort       int
		externalPort       int
		remoteHostIP       string
		want               string
	}{
		{
			name:               "local caller local bootstrap uses internal host",
			callerPlacement:    "local",
			bootstrapPlacement: "local",
			bootstrapHost:      "bootstrap-node",
			internalPort:       5001,
			externalPort:       15001,
			remoteHostIP:       "203.0.113.10",
			want:               "bootstrap-node:5001",
		},
		{
			name:               "local caller remote bootstrap uses external host override",
			callerPlacement:    "local",
			bootstrapPlacement: "remote",
			bootstrapHost:      "bootstrap-node",
			internalPort:       5001,
			externalPort:       15001,
			remoteHostIP:       "203.0.113.10",
			want:               "203.0.113.10:15001",
		},
		{
			name:               "remote caller local bootstrap uses docker host external",
			callerPlacement:    "remote",
			bootstrapPlacement: "local",
			bootstrapHost:      "bootstrap-node",
			internalPort:       5001,
			externalPort:       15001,
			remoteHostIP:       "203.0.113.10",
			want:               strings.TrimPrefix(framework.HostDockerInternal(), "http://") + ":5001",
		},
		{
			name:               "remote caller remote bootstrap uses internal host",
			callerPlacement:    "remote",
			bootstrapPlacement: "remote",
			bootstrapHost:      "bootstrap-node",
			internalPort:       5001,
			externalPort:       15001,
			remoteHostIP:       "203.0.113.10",
			want:               "bootstrap-node:5001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveNodeFacingBootstrapAddress(
				tt.callerPlacement,
				tt.bootstrapPlacement,
				tt.bootstrapHost,
				tt.internalPort,
				tt.externalPort,
				tt.remoteHostIP,
			)
			require.NoError(t, err, "resolveNodeFacingBootstrapAddress should not fail")
			require.Equal(t, tt.want, got, "unexpected resolved bootstrap address")
		})
	}
}

type fakeBlockchain struct {
	selector uint64
	id       uint64
	out      *blockchain.Output
}

func (f *fakeBlockchain) ChainSelector() uint64 { return f.selector }
func (f *fakeBlockchain) ChainID() uint64       { return f.id }
func (f *fakeBlockchain) ChainFamily() string   { return f.out.Family }
func (f *fakeBlockchain) IsFamily(chainFamily string) bool {
	return strings.EqualFold(f.out.Family, chainFamily)
}
func (f *fakeBlockchain) Fund(_ context.Context, _ string, _ uint64) error { return nil }
func (f *fakeBlockchain) CtfOutput() *blockchain.Output                    { return f.out }
func (f *fakeBlockchain) ToCldfChain() (cldf_chain.BlockChain, error)      { return nil, nil }

var _ creblockchains.Blockchain = (*fakeBlockchain)(nil)

func TestFindEVMChains_AllowsMissingWSForTron(t *testing.T) {
	nodeSet := &cre.NodeSet{
		Input: &ns.Input{
			Name: "workflow",
		},
		Placement:          "local",
		SupportedEVMChains: []uint64{TronEVMChainID},
	}
	donMetadata, err := cre.NewDonMetadata(nodeSet, 1, infra.Provider{Type: infra.Docker}, nil)
	require.NoError(t, err)

	input := cre.GenerateConfigsInput{
		DonMetadata: donMetadata,
		Blockchains: map[uint64]creblockchains.Blockchain{
			TronEVMChainID: &fakeBlockchain{
				selector: TronEVMChainID,
				id:       TronEVMChainID,
				out: &blockchain.Output{
					Type:   blockchain.TypeTron,
					Family: blockchain.FamilyEVM,
					Nodes: []*blockchain.Node{
						{
							InternalHTTPUrl: "http://tron:9090/jsonrpc",
							ExternalHTTPUrl: "http://localhost:9090/jsonrpc",
							InternalWSUrl:   "",
							ExternalWSUrl:   "",
						},
					},
				},
			},
		},
		BlockchainPlacementBySelector: map[uint64]string{
			TronEVMChainID: "local",
		},
	}

	evmChains, err := findEVMChains(input)
	require.NoError(t, err, "tron should not require WS endpoint resolution")
	require.Len(t, evmChains, 1)
	require.Equal(t, TronEVMChainID, evmChains[0].ChainID)
	require.NotEmpty(t, evmChains[0].HTTPRPC)
	require.Empty(t, evmChains[0].WSRPC, "tron WSRPC should remain empty when source has no ws endpoint")
}

func mustBuildGatewayTopology(t *testing.T, targetPlacement string) (*cre.Topology, *cre.DonGatewayConfiguration) {
	t.Helper()

	provider := infra.Provider{Type: infra.Docker}
	nodeSet := &cre.NodeSet{
		Input: &ns.Input{Name: "workflow"},
		NodeSpecs: []*cre.NodeSpecWithRole{
			{
				Input: &clnode.Input{Node: &clnode.NodeInput{}},
				Roles: []cre.NodeType{cre.BootstrapNode},
			},
		},
		Placement: targetPlacement,
	}
	donMetadata, err := cre.NewDonMetadata(nodeSet, 1, provider, nil)
	require.NoError(t, err, "failed to build DonMetadata")
	donsMetadata, err := cre.NewDonsMetadata([]*cre.DonMetadata{donMetadata}, provider)
	require.NoError(t, err, "failed to build DonsMetadata")

	gateway := &cre.DonGatewayConfiguration{
		GatewayConfiguration: &cre.GatewayConfiguration{
			NodeUUID: donMetadata.NodesMetadata[0].UUID,
			Outgoing: cre.Outgoing{
				Host: "bootstrap-gateway-node0",
				Port: 5003,
				Path: "/node",
			},
			AuthGatewayID: "gateway-node-0",
		},
	}

	return &cre.Topology{DonsMetadata: donsMetadata}, gateway
}
