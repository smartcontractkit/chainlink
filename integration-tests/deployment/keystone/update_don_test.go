package keystone_test

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	kstest "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/test-go/testify/require"
)

func TestUpdateDon(t *testing.T) {
	type args struct {
		lggr        logger.Logger
		req         *kslib.UpdateDonRequest
		setupDonCfg setupUpdateDonTestConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *kslib.UpdateDonResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupUpdateDonTest(t, tt.args.lggr, tt.args.setupDonCfg)
			got, err := kslib.UpdateDon(tt.args.lggr, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateDon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateDon() = %v, want %v", got, tt.want)
			}
		})
	}
}

type setupUpdateDonTestConfig struct {
	dons []kslib.DonCapabilities
}

type setupUpdateDonTestResult struct {
	registry *kcr.CapabilitiesRegistry
	chain    deployment.Chain
}

func setupUpdateDonTest(t *testing.T, lggr logger.Logger, cfg setupUpdateDonTestConfig) {
	t.Helper()
	req := newSetupTestRegistryRequest(t, cfg.dons)
	kstest.SetupTestRegistry(t, lggr, req)
	return
}

func newSetupTestRegistryRequest(t *testing.T, dons []kslib.DonCapabilities) *kstest.SetupTestRegistryRequest {
	t.Helper()
	allNops := make(map[string]*models.NodeOperator)
	for _, don := range dons {
		for _, nop := range don.Nops {
			nop := nop
			n, exists := allNops[nop.ID]
			if exists {
				nop.Nodes = append(n.Nodes, nop.Nodes...)
			}
			allNops[nop.ID] = nop
		}
	}
	var nops []*models.NodeOperator
	for _, nop := range allNops {
		nops = append(nops, nop)
	}
	nopsToNodes := makeNopToNodes(t, nops)
	testDons := makeTestDon(t, dons)
	p2pToCapabilities := makeP2PToCapabilities(t, dons)
	req := &kstest.SetupTestRegistryRequest{
		NopToNodes:        nopsToNodes,
		Dons:              testDons,
		P2pToCapabilities: p2pToCapabilities,
	}
	return req
}

func makeNopToNodes(t *testing.T, cloNops []*models.NodeOperator) map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc {
	nopToNodes := make(map[kcr.CapabilitiesRegistryNodeOperator][]*kslib.P2PSignerEnc)

	for _, nop := range cloNops {
		// all chain configs are the same wrt admin address & node keys
		// so we can just use the first one
		chainConfig := nop.Nodes[0].ChainConfigs[0]
		crnop := kcr.CapabilitiesRegistryNodeOperator{
			Name:  nop.Name,
			Admin: common.HexToAddress(chainConfig.AdminAddress),
		}
		var nodes []*kslib.P2PSignerEnc
		for _, node := range nop.Nodes {
			require.NotNil(t, node.PublicKey, "public key is nil %s", node.ID)
			p, err := kslib.NewP2PSignerEncFromCLO(chainConfig, *node.PublicKey)
			require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
			nodes = append(nodes, p)
		}
		nopToNodes[crnop] = nodes
	}
	return nopToNodes
}

func makeP2PToCapabilities(t *testing.T, dons []kslib.DonCapabilities) map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability {
	p2pToCapabilities := make(map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability)
	for _, don := range dons {
		for _, nop := range don.Nops {
			for _, node := range nop.Nodes {
				for _, cap := range don.Capabilities {
					p, err := kslib.NewP2PSignerEncFromCLO(node.ChainConfigs[0], *node.PublicKey)
					require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
					p2pToCapabilities[p.P2PKey] = append(p2pToCapabilities[p.P2PKey], cap)
				}
			}
		}
	}
	return p2pToCapabilities
}

func makeTestDon(t *testing.T, dons []kslib.DonCapabilities) []kstest.Don {
	out := make([]kstest.Don, len(dons))
	for i, don := range dons {
		out[i] = testDon(t, don)
	}
	return out
}

func testDon(t *testing.T, don kslib.DonCapabilities) kstest.Don {
	var p2pids []p2pkey.PeerID
	for _, nop := range don.Nops {
		// all chain configs are the same wrt admin address & node keys
		// so we can just use the first one
		chainConfig := nop.Nodes[0].ChainConfigs[0]
		for _, node := range nop.Nodes {
			p, err := kslib.NewP2PSignerEncFromCLO(chainConfig, *node.PublicKey)
			require.NoError(t, err, "failed to make p2p signer enc from clo nod %s", node.ID)
			p2pids = append(p2pids, p.P2PKey)
		}
	}

	var capabilityConfigs []kslib.CapabilityConfig
	for _, cap := range don.Capabilities {
		capabilityConfigs = append(capabilityConfigs, kslib.CapabilityConfig{
			Capability: cap,
		})
	}
	return kstest.Don{
		Name:              don.Name,
		P2PIDs:            p2pids,
		CapabilityConfigs: capabilityConfigs,
	}
}

func newP2PSignerEnc(signer [32]byte, p2pkey p2pkey.PeerID, encryptionPublicKey [32]byte) *kslib.P2PSignerEnc {
	return &kslib.P2PSignerEnc{
		Signer:              signer,
		P2PKey:              p2pkey,
		EncryptionPublicKey: encryptionPublicKey,
	}
}
