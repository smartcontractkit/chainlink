package test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"maps"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
	changeset2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	envtest "github.com/smartcontractkit/chainlink/deployment/environment/test"
	changeset3 "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

const (
	DONName           = "test-don"
	RegistryQualifier = "test-registry"
	Zone              = "test-zone-1"
	TotalNodes        = 4
	// EnvironmentName is the name of the environment used in the test
	//
	// This is set by the runtime loader, but the constant is not exposed so we define it here.
	//
	// This will be fixed in a future release of the chainlink-deployments-framework.
	EnvironmentName = "test_environment"
)

type EnvWrapperV2 struct {
	t *testing.T

	Runtime *runtime.Runtime
	Don     *viewOnlyDon

	TestJD *envtest.JDNodeService

	Env              *cldf.Environment
	RegistrySelector uint64
	RegistryAddress  common.Address
	AptosSelector    uint64
}

type donConfig struct {
	Name             string
	N                int
	F                int
	CapabilityConfig map[string]*pb.CapabilityConfig
	Labels           map[string]string
	RegistryChainSel uint64
}

// TODO CRE-999; aptos can be made optional
func initRuntime(t *testing.T, lggr logger.Logger) *EnvWrapperV2 {
	var (
		registryChainSel = chain_selectors.TEST_90000001.Selector
		// by inspection, the only chain that is needed is evm, but some callers
		// expect aptos keys and therefore an aptos selector to use for generating
		// the keys
		aptosChainSel = chain_selectors.APTOS_LOCALNET.Selector

		donCfg = donConfig{
			Name:             DONName,
			N:                TotalNodes,
			F:                (TotalNodes-1)/3 + 1,
			RegistryChainSel: registryChainSel,
		}
	)

	// Setup the view only DON. Only need one DON
	don := newViewOnlyNodes(t, registryChainSel, aptosChainSel, donCfg)

	// Setup JD service
	nodes := make(deployment.Nodes, 0, don.N())
	for _, v := range don.m {
		nodes = append(nodes, *v)
	}
	jd := envtest.NewJDService(nodes)

	// Setup the runtime environment
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulatedWithConfig(t, []uint64{registryChainSel}, onchain.EVMSimLoaderConfig{
			NumAdditionalAccounts: 1,
		}),
		environment.WithLogger(lggr),
		environment.WithOffchainClient(jd),
		environment.WithNodeIDs(nodes.IDs()),
	))
	require.NoError(t, err)

	h := &EnvWrapperV2{
		t:                t,
		Runtime:          rt,
		TestJD:           jd,
		Don:              don,
		AptosSelector:    aptosChainSel,
		RegistrySelector: registryChainSel,
	}

	// Deploy the capabilities registry
	deployCapabilitiesRegistry(t, h)

	registryAddrs := h.Runtime.Environment().DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(h.RegistrySelector),
		datastore.AddressRefByType("CapabilitiesRegistry"),
	)
	require.Len(t, registryAddrs, 1)
	registryAddress := registryAddrs[0].Address
	h.RegistryAddress = common.HexToAddress(registryAddress)

	return h
}

func NewTestHarness(t *testing.T, useMCMS bool) *EnvWrapperV2 {
	t.Helper()

	return SetupEnvV2(t, useMCMS)
}

// SetupEnvV2 starts an environment with a single DON, 4 nodes and a capabilities registry v2 deployed and configured.
func SetupEnvV2(t *testing.T, useMCMS bool) *EnvWrapperV2 {
	t.Helper()

	lggr := logger.Test(t)
	h := initRuntime(t, lggr)
	t.Log("Initialized runtime", "registryChainSel", h.RegistrySelector)

	configureCapabilitiesRegistry(t, h)

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		h.RegistryAddress,
		h.Runtime.Environment().BlockChains.EVMChains()[h.RegistrySelector].Client,
	)
	require.NoError(t, err)
	require.NotNil(t, capReg)

	gotNodes, err := capReg.GetNodesByP2PIds(nil, h.Don.GetP2PIDs().Bytes32())
	require.NoError(t, err)
	require.Len(t, gotNodes, len(h.Don.GetP2PIDs()))
	require.Len(t, gotNodes, h.Don.N()) // +1 for bootstrap
	for _, n := range gotNodes {
		require.Equal(t, "test-capability@1.0.0", n.CapabilityIds[0])
	}

	gotDON, err := capReg.GetDONByName(nil, h.Don.Name())
	require.NoError(t, err)
	require.Len(t, gotDON.NodeP2PIds, len(h.Don.GetP2PIDs()))

	// Sort both slices before comparison
	sort.Slice(gotDON.NodeP2PIds, func(i, j int) bool {
		return bytes.Compare(gotDON.NodeP2PIds[i][:], gotDON.NodeP2PIds[j][:]) < 0
	})
	sortedNodesP2PIDsBytes := make([][32]byte, len(h.Don.GetP2PIDs()))
	copy(sortedNodesP2PIDsBytes, h.Don.GetP2PIDs().Bytes32())
	sort.Slice(sortedNodesP2PIDsBytes, func(i, j int) bool {
		return bytes.Compare(sortedNodesP2PIDsBytes[i][:], sortedNodesP2PIDsBytes[j][:]) < 0
	})
	for i, id := range gotDON.NodeP2PIds {
		require.Equal(t, sortedNodesP2PIDsBytes[i], id)
	}

	if useMCMS {
		t.Log("Setting up MCMS infrastructure...")
		timelockCfgs := map[uint64]cldfproposalutils.MCMSWithTimelockConfig{
			h.RegistrySelector: cldftesthelpers.SingleGroupTimelockConfig(t),
		}

		err = h.Runtime.Exec(
			runtime.ChangesetTask(cldf.CreateLegacyChangeSet(mcmschangesets.DeployMCMSWithTimelockV2), timelockCfgs),
		)
		require.NoError(t, err, "failed to deploy MCMS infrastructure")
		t.Log("MCMS infrastructure deployed successfully")

		t.Log("Transferring ownership to MCMS...")
		err = h.Runtime.Exec(
			runtime.ChangesetTask(
				cldf.CreateLegacyChangeSet(changeset3.AcceptAllOwnershipsProposal),
				&changeset3.AcceptAllOwnershipRequest{
					ChainSelector: h.RegistrySelector,
					MinDelay:      0,
				},
			),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
		)
		require.NoError(t, err, "failed to transfer ownership to MCMS")
		t.Log("Ownership transferred to MCMS successfully")
	}

	// Set the Environment because some changesets still use the Environment pointer
	// To be removed once all changesets below use the runtime instead
	h.Env = new(h.Runtime.Environment())

	return h
}

func configureCapabilitiesRegistry(t *testing.T, h *EnvWrapperV2) {
	t.Helper()

	chainID, err := chain_selectors.GetChainIDFromSelector(h.RegistrySelector)
	require.NoError(t, err)

	registryChainDetails, err := chain_selectors.GetChainDetailsByChainIDAndFamily(chainID, chain_selectors.FamilyEVM)
	require.NoError(t, err)

	donNodes, err := h.Don.AllNodes()
	require.NoError(t, err)

	var nodes []changeset2.CapabilitiesRegistryNodeParams
	for _, n := range donNodes {
		p2pID := n.PeerID.String()
		ocrConfig, ok := n.OCRConfigs[registryChainDetails]
		require.True(t, ok, "node %s does not have OCR config for registry chain %d", n.Name, h.RegistrySelector)

		nodes = append(nodes, changeset2.CapabilitiesRegistryNodeParams{
			NOP:                 "Operator 1",
			P2pID:               p2pID,
			CsaKey:              n.CSA,
			EncryptionPublicKey: n.WorkflowKey,
			Signer:              hex.EncodeToString(ocrConfig.OnchainPublicKey),
			CapabilityIDs: []string{
				"test-capability@1.0.0",
			},
		})
	}

	err = h.Runtime.Exec(
		runtime.ChangesetTask(changeset2.ConfigureCapabilitiesRegistry{},
			changeset2.ConfigureCapabilitiesRegistryInput{
				ChainSelector:               h.RegistrySelector,
				CapabilitiesRegistryAddress: h.RegistryAddress.Hex(),
				Nops: []changeset2.CapabilitiesRegistryNodeOperator{
					{
						Name:  "Operator 1",
						Admin: common.HexToAddress("0x01"),
					},
				},
				Nodes: nodes,
				Capabilities: []changeset2.CapabilitiesRegistryCapability{
					{
						CapabilityID: "test-capability@1.0.0",
						Metadata:     map[string]any{"capabilityType": 2},
					},
				},
				DONs: []changeset2.CapabilitiesRegistryNewDONParams{
					{
						Name:        h.Don.Name(),
						F:           uint8(h.Don.F()), //nolint:gosec // disable G115
						Nodes:       h.Don.GetP2PIDs().Strings(),
						DonFamilies: []string{"test-family"},
						Config:      map[string]any{"defaultConfig": map[string]any{}},
						CapabilityConfigurations: []changeset2.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityID: "test-capability@1.0.0",
							},
						},
						IsPublic:         true,
						AcceptsWorkflows: true,
					},
				},
			},
		),
	)
	require.NoError(t, err)
}

// deployCapabilitiesRegistry deploys the capabilities registry in the runtime environment
func deployCapabilitiesRegistry(t *testing.T, h *EnvWrapperV2) {
	t.Helper()

	err := h.Runtime.Exec(
		runtime.ChangesetTask(changeset2.DeployCapabilitiesRegistry{},
			changeset2.DeployCapabilitiesRegistryInput{
				ChainSelector: h.RegistrySelector,
				Qualifier:     RegistryQualifier,
			},
		),
	)
	require.NoError(t, err)
}

// newViewOnlyNodes creates a view only DON with the given configuration.
func newViewOnlyNodes(t *testing.T, registryChainSel, aptosChainSel uint64, donCfg donConfig) *viewOnlyDon {
	var nodesCfg []envtest.NodeConfig

	for i := 0; i < donCfg.N; i++ {
		labels := map[string]string{
			"don-" + donCfg.Name: donCfg.Name,
			"environment":        EnvironmentName,
			"product":            "cre",
			"type":               "plugin",
			"zone":               Zone,
		}
		if donCfg.Labels != nil {
			maps.Copy(labels, donCfg.Labels)
		}

		nCfg := envtest.NodeConfig{
			ChainSelectors: []uint64{registryChainSel, aptosChainSel, chain_selectors.SOLANA_DEVNET.Selector},
			Name:           fmt.Sprintf("%s-%d", donCfg.Name, i),
			Labels:         labels,
		}
		nodesCfg = append(nodesCfg, nCfg)
	}

	btLabels := map[string]string{
		"don-" + donCfg.Name: donCfg.Name,
		"environment":        EnvironmentName,
		"product":            "cre",
		"type":               "bootstrap",
		"zone":               Zone,
	}
	if donCfg.Labels != nil {
		maps.Copy(btLabels, donCfg.Labels)
	}
	nodesCfg = append(nodesCfg, envtest.NodeConfig{
		ChainSelectors: []uint64{registryChainSel, aptosChainSel, chain_selectors.SOLANA_DEVNET.Selector},
		Name:           donCfg.Name + "-bootstrap",
		Labels:         btLabels,
	})

	n := envtest.NewNodes(t, nodesCfg)
	require.Len(t, n, donCfg.N+1) // +1 for bootstrap

	return newViewOnlyDon(donCfg.Name, n)
}
