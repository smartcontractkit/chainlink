package operations

import nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

type JDNode struct {
	ID              string               `json:"id" yaml:"id"`                           // Unique identifier for the node.
	Name            string               `json:"name" yaml:"name"`                       // Human-readable name for the node.
	PublicKey       string               `json:"publicKey" yaml:"publicKey"`             // Public key used for secure communications.
	IsEnabled       bool                 `json:"isEnabled" yaml:"isEnabled"`             // Indicates if the node is currently enabled.
	IsConnected     bool                 `json:"isConnected" yaml:"isConnected"`         // Indicates if the node is currently connected to the network.
	Labels          []map[string]string  `json:"labels" yaml:"labels"`                   // Set of labels associated with the node.
	CreatedAt       string               `json:"createdAt" yaml:"createdAt"`             // Timestamp when the node was created.
	UpdatedAt       string               `json:"updatedAt" yaml:"updatedAt"`             // Timestamp when the node was last updated.
	WorkflowKey     *string              `json:"workflowKey" yaml:"workflowKey"`         // Workflow Public key
	P2PKeyBundles   []JDNodeP2PKeyBundle `json:"p2pKeyBundles" yaml:"p2pKeyBundles"`     // List of P2P key bundles associated with the node.
	NopFriendlyName string               `json:"nopFriendlyName" yaml:"nopFriendlyName"` // Friendly name defined by NOP
	Version         string               `json:"version" yaml:"version"`                 // Node Version
}

type JDNodeP2PKeyBundle struct {
	PeerID    string `json:"peerID" yaml:"peerID"`
	PublicKey string `json:"publicKey" yaml:"publicKey"`
}

func NewJDNodeFromProto(n *nodev1.Node) JDNode {
	if n == nil {
		return JDNode{}
	}

	var labels []map[string]string
	if n.Labels != nil {
		for _, l := range n.Labels {
			value := ""
			if l.Value != nil {
				value = *l.Value
			}
			labels = append(labels, map[string]string{l.Key: value})
		}
	}
	var p2pKeyBundles []JDNodeP2PKeyBundle
	for _, b := range n.P2PKeyBundles {
		if b == nil {
			continue
		}

		p2pKeyBundles = append(p2pKeyBundles, JDNodeP2PKeyBundle{
			PeerID:    b.PeerId,
			PublicKey: b.PublicKey,
		})
	}

	return JDNode{
		ID:              n.Id,
		Name:            n.Name,
		PublicKey:       n.PublicKey,
		IsEnabled:       n.IsEnabled,
		IsConnected:     n.IsConnected,
		Labels:          labels,
		P2PKeyBundles:   p2pKeyBundles,
		CreatedAt:       n.CreatedAt.String(),
		UpdatedAt:       n.UpdatedAt.String(),
		WorkflowKey:     n.WorkflowKey,
		NopFriendlyName: n.NopFriendlyName,
		Version:         n.Version,
	}
}
