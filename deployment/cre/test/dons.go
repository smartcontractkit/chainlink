package test

import (
	"sort"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// testDon is an interface for a DON that can be used in tests
type testDon interface {
	GetP2PIDs() P2PIDs
	N() int
	F() int
	Name() string
}

var _ testDon = (*memoryDon)(nil)

// memoryDon is backed by in-memory nodes running the full chainlink Application
type memoryDon struct {
	name string
	m    map[string]memory.Node
}

func newMemoryDon(name string, m map[string]memory.Node) *memoryDon {
	return &memoryDon{name: name, m: m}
}

func (d *memoryDon) GetP2PIDs() P2PIDs {
	var out []p2pkey.PeerID
	for _, n := range d.m {
		out = append(out, n.Keys.PeerID)
	}
	return out
}

func (d *memoryDon) N() int {
	return len(d.m)
}

func (d *memoryDon) F() int {
	return (d.N() - 1) / 3
}

func (d *memoryDon) Name() string {
	return d.name
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

// testDons is a collection of testDon with convenience methods commonly used in tests
type testDons interface {
	Get(name string) testDon
	Put(d testDon)
	List() []testDon
	// Unique list of p2pIDs across all dons
	P2PIDs() P2PIDs
	AllNodes() (map[string]node, error)
}

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

var _ testDons = (*memoryDons)(nil)

// memoryDons implements [testDons] backed by [memoryDon]
type memoryDons struct {
	dons map[string]*memoryDon
}

func newMemoryDons() *memoryDons {
	return &memoryDons{dons: make(map[string]*memoryDon)}
}

func (d *memoryDons) Get(name string) testDon {
	x := d.dons[name]
	return x
}

func (d *memoryDons) Put(d2 testDon) {
	d.dons[d2.Name()] = d2.(*memoryDon)
}

func (d *memoryDons) List() []testDon {
	out := make([]testDon, 0, len(d.dons))
	donNames := make([]string, 0, len(d.dons))
	for k := range d.dons {
		donNames = append(donNames, k)
	}
	sort.Strings(donNames)
	for _, name := range donNames {
		out = append(out, d.dons[name])
	}
	return out
}

func (d *memoryDons) P2PIDs() P2PIDs {
	var out P2PIDs
	for _, d := range d.dons {
		out = append(out, d.GetP2PIDs()...)
	}
	return out.Unique()
}

func (d *memoryDons) AllNodes() (map[string]node, error) {
	out := make(map[string]node)
	for _, don := range d.dons {
		for k, v := range don.m {
			dNode, err := v.DeploymentNode()
			if err != nil {
				return nil, err
			}

			out[k] = node{
				ID:          v.ID,
				Name:        v.Name,
				PeerID:      v.Keys.PeerID,
				CSA:         v.Keys.CSA.PublicKeyString(),
				WorkflowKey: v.Keys.WorkflowKey.PublicKeyString(),
				IsBoostrap:  v.IsBoostrap,
				OCRConfigs:  dNode.SelToOCRConfig,
			}
		}
	}
	return out, nil
}

func (d *memoryDons) AllNodesForJD() map[string]memory.Node {
	out := make(map[string]memory.Node)
	for _, don := range d.dons {
		for k, v := range don.m {
			out[k] = v
		}
	}
	return out
}

// viewOnlyDons implements [testDons] backed by [viewOnlyDon]
type viewOnlyDons struct {
	dons map[string]*viewOnlyDon
}

func newViewOnlyDons() *viewOnlyDons {
	return &viewOnlyDons{dons: make(map[string]*viewOnlyDon)}
}

func (d *viewOnlyDons) Get(name string) testDon {
	x := d.dons[name]
	return x
}

func (d *viewOnlyDons) Put(d2 testDon) {
	d.dons[d2.Name()] = d2.(*viewOnlyDon)
}

func (d *viewOnlyDons) List() []testDon {
	out := make([]testDon, 0, len(d.dons))
	donNames := make([]string, 0, len(d.dons))
	for k := range d.dons {
		donNames = append(donNames, k)
	}
	sort.Strings(donNames)
	for _, name := range donNames {
		out = append(out, d.dons[name])
	}
	return out
}

func (d *viewOnlyDons) P2PIDs() P2PIDs {
	var out P2PIDs
	for _, d := range d.dons {
		out = append(out, d.GetP2PIDs()...)
	}
	return out.Unique()
}

func (d *viewOnlyDons) AllNodes() (map[string]node, error) {
	out := make(map[string]node)
	for _, don := range d.dons {
		for k, v := range don.m {
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
	}
	return out, nil
}

func (d *viewOnlyDons) NodeList() deployment.Nodes {
	tmp := make(map[string]*deployment.Node)
	for _, don := range d.dons {
		for k, v := range don.m {
			tmp[k] = v
		}
	}

	nodes := make([]deployment.Node, 0, len(tmp))
	for _, v := range tmp {
		nodes = append(nodes, *v)
	}
	return nodes
}
