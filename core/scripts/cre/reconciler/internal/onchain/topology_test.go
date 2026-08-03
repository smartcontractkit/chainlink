package onchain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func testSecretsTOML(t *testing.T) string {
	t.Helper()

	p2pKey, err := crypto.NewP2PKey("dev-password")
	require.NoError(t, err)
	dkgKey, err := crypto.NewDKGRecipientKey("dev-password")
	require.NoError(t, err)

	keys := &secrets.NodeKeys{
		P2PKey: p2pKey,
		DKGKey: dkgKey,
		EVM:    map[uint64]*crypto.EVMKey{},
	}
	secretsTOML, err := keys.ToNodeSecretsTOML()
	require.NoError(t, err)
	return secretsTOML
}

func testNodeSpec(t *testing.T, secretsTOML string, role cre.NodeType) *cre.NodeSpecWithRole {
	t.Helper()

	return &cre.NodeSpecWithRole{
		Input: &clnode.Input{
			Node: &clnode.NodeInput{
				TestSecretsOverrides: secretsTOML,
			},
		},
		Roles: []cre.NodeType{role},
	}
}

func testWorkflowBootstrapTopology(t *testing.T) *cre.Topology {
	t.Helper()

	provider := infra.Provider{Type: infra.Kubernetes}

	workerSpecs := make([]*cre.NodeSpecWithRole, 4)
	for i := range workerSpecs {
		workerSpecs[i] = testNodeSpec(t, testSecretsTOML(t), cre.WorkerNode)
	}
	workflowNS := &cre.NodeSet{
		Input: &ns.Input{
			Name:         "workflow",
			Nodes:        4,
			OverrideMode: "all",
		},
		NodeSpecs:          workerSpecs,
		Capabilities:       []string{"cron", "consensus", "don-time"},
		DONTypes:           []string{"workflow"},
		SupportedEVMChains: []uint64{1337},
		DonFamily:          "workflow",
	}

	bootstrapDon := &domain.DON{Name: "bootstrap", DONTypes: []string{"bootstrap"}}
	bootstrapNS := newBootstrapOnlyNodeSet(
		bootstrapDon,
		[]*cre.NodeSpecWithRole{testNodeSpec(t, testSecretsTOML(t), cre.BootstrapNode)},
		[]uint64{1337},
	)

	topology, err := cre.NewTopology([]*cre.NodeSet{workflowNS, bootstrapNS}, provider, cre.CapabilityConfigs{})
	require.NoError(t, err)
	return topology
}

func testSingleWorkerTopology(t *testing.T) *cre.Topology {
	t.Helper()

	secretsTOML := testSecretsTOML(t)
	provider := infra.Provider{Type: infra.Kubernetes}
	workflowNS := &cre.NodeSet{
		Input: &ns.Input{
			Name:         "workflow",
			Nodes:        1,
			OverrideMode: "all",
		},
		NodeSpecs:          []*cre.NodeSpecWithRole{testNodeSpec(t, secretsTOML, cre.WorkerNode)},
		Capabilities:       []string{"cron", "consensus", "don-time"},
		DONTypes:           []string{"workflow"},
		SupportedEVMChains: []uint64{1337},
		DonFamily:          "workflow",
	}
	bootstrapDon := &domain.DON{Name: "bootstrap", DONTypes: []string{"bootstrap"}}
	bootstrapNS := newBootstrapOnlyNodeSet(
		bootstrapDon,
		[]*cre.NodeSpecWithRole{testNodeSpec(t, secretsTOML, cre.BootstrapNode)},
		[]uint64{1337},
	)
	topology, err := cre.NewTopology([]*cre.NodeSet{workflowNS, bootstrapNS}, provider, cre.CapabilityConfigs{})
	require.NoError(t, err)
	return topology
}

func TestRequiredEVMChainIDsForDON(t *testing.T) {
	t.Parallel()

	don := &domain.DON{
		Name:         "workflow",
		Capabilities: []string{"cron", "evm-1337", "http-action"},
	}
	chains := requiredEVMChainIDsForDON(don, []uint64{1337, 31337}, 1337)
	require.Equal(t, []uint64{1337}, chains)
}

func TestRequiredEVMChainIDsForDON_BaseEVMCapability(t *testing.T) {
	t.Parallel()

	don := &domain.DON{
		Name:         "workflow",
		Capabilities: []string{"evm"},
	}
	chains := requiredEVMChainIDsForDON(don, []uint64{1337, 31337}, 1337)
	require.Equal(t, []uint64{1337, 31337}, chains)
}

func TestNewGatewayNodeSet_PerNode(t *testing.T) {
	t.Parallel()

	specs := []*cre.NodeSpecWithRole{{Roles: []cre.NodeType{cre.WorkerNode}}}
	nodeSet := newGatewayNodeSet(
		"gateway-don",
		specs,
		"workflow",
		[]uint64{1337},
	)

	require.Equal(t, "gateway-don", nodeSet.Name)
	require.Equal(t, "workflow", nodeSet.GatewayDonID)
	require.Equal(t, []string{"gateway"}, nodeSet.DONTypes)
	require.Len(t, nodeSet.NodeSpecs, 1)
	require.Equal(t, []cre.NodeType{cre.GatewayNode}, nodeSet.NodeSpecs[0].Roles)
}

func TestBuildGatewayNodeSets_MultipleNodes(t *testing.T) {
	t.Parallel()

	specsA := []*cre.NodeSpecWithRole{{}}
	specsB := []*cre.NodeSpecWithRole{{}}
	nsA := newGatewayNodeSet("gateway-don-a", specsA, "workflow-a", []uint64{1337})
	nsB := newGatewayNodeSet("gateway-don-b", specsB, "workflow-b", []uint64{1337})

	require.NotEqual(t, nsA.GatewayDonID, nsB.GatewayDonID)
	require.NotEqual(t, nsA.Name, nsB.Name)
	require.Equal(t, "gateway-don-a", nsA.Name)
	require.Equal(t, "gateway-don-b", nsB.Name)
}

func TestBuildTopology_IncludesBootstrapOnlyDON(t *testing.T) {
	t.Parallel()

	secretsTOML := testSecretsTOML(t)
	provider := infra.Provider{Type: infra.Kubernetes}

	workerSpecs := make([]*cre.NodeSpecWithRole, 4)
	for i := range workerSpecs {
		workerSpecs[i] = testNodeSpec(t, secretsTOML, cre.WorkerNode)
	}
	workflowNS := &cre.NodeSet{
		Input: &ns.Input{
			Name:         "workflow",
			Nodes:        4,
			OverrideMode: "all",
		},
		NodeSpecs:          workerSpecs,
		Capabilities:       []string{"cron", "consensus", "don-time"},
		DONTypes:           []string{"workflow"},
		SupportedEVMChains: []uint64{1337},
		DonFamily:          "workflow",
	}

	bootstrapDon := &domain.DON{Name: "bootstrap", DONTypes: []string{"bootstrap"}}
	bootstrapNS := newBootstrapOnlyNodeSet(
		bootstrapDon,
		[]*cre.NodeSpecWithRole{testNodeSpec(t, secretsTOML, cre.BootstrapNode)},
		[]uint64{1337},
	)

	topology, err := cre.NewTopology([]*cre.NodeSet{workflowNS, bootstrapNS}, provider, cre.CapabilityConfigs{})
	require.NoError(t, err)

	_, hasBootstrap := topology.DonsMetadata.Bootstrap()
	require.True(t, hasBootstrap)

	foundBootstrapDON := false
	for _, donMeta := range topology.DonsMetadata.List() {
		if donMeta.Name == "bootstrap" {
			foundBootstrapDON = true
		}
	}
	require.True(t, foundBootstrapDON)
}

func TestBuildBootstrapOnlyNodeSet_MissingBootstrap_Fails(t *testing.T) {
	t.Parallel()

	d := &Deployer{}
	cv := &domain.ChartValues{}
	state := &domain.StateFile{}
	don := &domain.DON{Name: "bootstrap", DONTypes: []string{"bootstrap"}}

	_, _, _, err := d.buildBootstrapOnlyNodeSet(context.Background(), don, []uint64{1337}, 1337, nil, cv, state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no resolvable bootstrap node")
}

func TestBuildBootstrapOnlyNodeSet_RejectsGatewayDON(t *testing.T) {
	t.Parallel()

	d := &Deployer{}
	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{{Name: "gateway-0", NodeType: domain.RoleGateway}},
	}
	state := &domain.StateFile{}
	don := &domain.DON{Name: "don-3", DONTypes: []string{"gateway"}}

	_, _, _, err := d.buildBootstrapOnlyNodeSet(context.Background(), don, []uint64{1337}, 1337, nil, cv, state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a bootstrap DON")
}

func TestBuildTopologyRouting_GatewayDONNotBootstrapOnly(t *testing.T) {
	t.Parallel()

	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "gateway-0", NodeType: domain.RoleGateway, Namespace: "zone-a-gateway", DONName: "don-3"},
		},
	}
	gatewayDon := domain.DON{Name: "don-3", DONTypes: []string{"gateway"}}
	require.Empty(t, gatewayDon.WorkerNodes(cv))
	require.False(t, gatewayDon.IsBootstrapOnly(cv))
	require.True(t, gatewayDon.IsGatewayDon())
}

func TestBuildBootstrapOnlyNodeSet_CreatesZeroWorkerNodeSet(t *testing.T) {
	t.Parallel()

	secretsTOML := testSecretsTOML(t)
	don := &domain.DON{Name: "bootstrap", DONTypes: []string{"bootstrap"}}
	nodeSet := newBootstrapOnlyNodeSet(
		don,
		[]*cre.NodeSpecWithRole{testNodeSpec(t, secretsTOML, cre.BootstrapNode)},
		[]uint64{1337},
	)

	require.Equal(t, 0, nodeSet.Nodes)
	require.Len(t, nodeSet.NodeSpecs, 1)
	require.Equal(t, []cre.NodeType{cre.BootstrapNode}, nodeSet.NodeSpecs[0].Roles)
}
