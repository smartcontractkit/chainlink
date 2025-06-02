package test

import (
	"crypto/rand"
	"math/big"
	"slices"
	"strings"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
)

type NodeConfig struct {
	ChainSelectors []uint64
	Name           string
	Labels         map[string]string
}

func NewNode(t *testing.T, c NodeConfig) *deployment.Node {
	t.Helper()
	k := randSeed(t)
	p2p := p2pkey.MustNewV2XXXTestingOnly(k)
	ocrConfigs := make(map[chain_selectors.ChainDetails]deployment.OCRConfig)
	for _, cs := range c.ChainSelectors {
		ocrConfigs[chain_selectors.ChainDetails{ChainSelector: cs}] = testOCRConfig(t, cs, p2p)
	}
	if c.Labels == nil {
		c.Labels = map[string]string{}
	}
	// make sure to add the p2p_id label b/c downstream systems expect it
	c.Labels["p2p_id"] = p2p.PeerID().String()
	return &deployment.Node{
		NodeID:         c.Name,
		Name:           c.Name,
		PeerID:         p2p.PeerID(),
		CSAKey:         csakey.MustNewV2XXXTestingOnly(k).ID(),
		WorkflowKey:    workflowkey.MustNewXXXTestingOnly(k).ID(),
		AdminAddr:      gethcommon.BigToAddress(k).Hex(),
		Labels:         labelsConversion(c.Labels),
		SelToOCRConfig: ocrConfigs,
	}
}

func NewNodes(t *testing.T, configs []NodeConfig) []*deployment.Node {
	nodes := make([]*deployment.Node, len(configs))
	for i, c := range configs {
		nodes[i] = NewNode(t, c)
	}
	return nodes
}

func randSeed(t *testing.T) *big.Int {
	maxVal := new(big.Int)
	maxVal.Exp(big.NewInt(2), big.NewInt(256), nil)
	randomInt, err := rand.Int(rand.Reader, maxVal)
	require.NoError(t, err)
	return randomInt
}

func labelsConversion(m map[string]string) []*ptypes.Label {
	out := make([]*ptypes.Label, len(m))
	i := 0
	for k, v := range m {
		v := v
		out[i] = &ptypes.Label{Key: k, Value: &v}
		i++
	}
	return out
}

func testOCRConfig(t *testing.T, sel uint64, p2p p2pkey.KeyV2) deployment.OCRConfig {
	t.Helper()
	f, err := chain_selectors.GetSelectorFamily(sel)
	require.NoError(t, err, "selector %d not found", sel)
	seed := p2p.PeerID()
	copy(seed[:], []byte(f))
	require.NoError(t, err)
	return deployment.OCRConfig{
		PeerID:                    p2p.PeerID(),
		OffchainPublicKey:         types2.OffchainPublicKey(seed),
		OnchainPublicKey:          types2.OnchainPublicKey(seed[:32]),
		TransmitAccount:           types2.Account(gethcommon.BytesToAddress(seed[:]).Hex()),
		ConfigEncryptionPublicKey: types2.ConfigEncryptionPublicKey(seed[:32]),
		KeyBundleID:               "fake_orc_bundle_" + f,
	}
}

func ApplyNodeFilter(filter *nodev1.ListNodesRequest_Filter, node *nodev1.Node) bool {
	if filter == nil {
		return true
	}
	if len(filter.Ids) > 0 {
		idx := slices.IndexFunc(filter.Ids, func(id string) bool {
			return node.Id == id
		})
		if idx < 0 {
			return false
		}
	}
	for _, selector := range filter.Selectors {
		idx := slices.IndexFunc(node.Labels, func(label *ptypes.Label) bool {
			return label.Key == selector.Key
		})
		if idx < 0 {
			return false
		}
		label := node.Labels[idx]

		switch selector.Op {
		case ptypes.SelectorOp_IN:
			values := strings.Split(*selector.Value, ",")
			found := slices.Contains(values, *label.Value)
			if !found {
				return false
			}
		case ptypes.SelectorOp_EQ:
			if *label.Value != *selector.Value {
				return false
			}
		case ptypes.SelectorOp_EXIST:
			// do nothing
		default:
			panic("unimplemented selector")
		}
	}
	return true
}
