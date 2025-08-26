package test

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/deployment"
)

type node struct {
	ID   string
	Name string
	// Transmitter key/OCR keys for this node
	PeerID      p2pkey.PeerID
	CSA         string
	WorkflowKey string
	IsBoostrap  bool

	OCRConfigs map[chain_selectors.ChainDetails]deployment.OCRConfig
}

// viewOnlyDon represents a DON that is backed by view of nodes, not actual, useable nodes
type viewOnlyDon struct {
	name string
	m    map[string]*deployment.Node
}

func newViewOnlyDon(name string, nodes []*deployment.Node) *viewOnlyDon {
	m := make(map[string]*deployment.Node)
	for _, n := range nodes {
		m[n.PeerID.String()] = n
	}
	return &viewOnlyDon{name: name, m: m}
}

func (d *viewOnlyDon) GetP2PIDs() P2PIDs {
	var out []p2pkey.PeerID
	for _, n := range d.m {
		out = append(out, n.PeerID)
	}
	return out
}

func (d *viewOnlyDon) N() int {
	return len(d.m)
}

func (d *viewOnlyDon) F() int {
	return (d.N() - 1) / 3
}

func (d *viewOnlyDon) Name() string {
	return d.name
}

func (d *viewOnlyDon) AllNodes() (map[string]node, error) {
	out := make(map[string]node)
	for k, v := range d.m {
		out[k] = node{
			ID:          v.NodeID,
			Name:        v.Name,
			PeerID:      v.PeerID,
			CSA:         v.CSAKey,
			WorkflowKey: v.WorkflowKey,
			OCRConfigs:  v.SelToOCRConfig,
			IsBoostrap:  v.MultiAddr != "",
		}
	}
	return out, nil
}
