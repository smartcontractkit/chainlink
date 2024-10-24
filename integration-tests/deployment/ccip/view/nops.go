package view

import (
	"context"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

type NopsView struct {
	Nops map[string]NopView `json:"nops,omitempty"`
}

type NopView struct {
	// NodeID is the unique identifier of the node
	NodeID       string                `json:"nodeID"`
	IsBootstrap  bool                  `json:"isBootstrap"`
	OCRKeys      map[string]OCRKeyView `json:"ocrKeys"`
	PayeeAddress string                `json:"payeeAddress"`
	CSAKey       string                `json:"csaKey"`
	IsConnected  bool                  `json:"isConnected"`
	IsEnabled    bool                  `json:"isEnabled"`
}

type OCRKeyView struct {
	OffchainPublicKey         string `json:"offchainPublicKey"`
	OnchainPublicKey          string `json:"onchainPublicKey"`
	PeerID                    string `json:"peerID"`
	TransmitAccount           string `json:"transmitAccount"`
	ConfigEncryptionPublicKey string `json:"configEncryptionPublicKey"`
	KeyBundleID               string `json:"keyBundleID"`
}

func NewNopsView() NopsView {
	return NopsView{
		Nops: make(map[string]NopView),
	}
}

func GenerateNopsView(nodeIds []string, oc deployment.OffchainClient) (NopsView, error) {
	nops := NewNopsView()
	nodes, err := deployment.NodeInfo(nodeIds, oc)
	if err != nil {
		return nops, err
	}
	for _, node := range nodes {
		// get node info
		nodeDetails, err := oc.GetNode(context.Background(), &nodev1.GetNodeRequest{Id: node.NodeID})
		if err != nil {
			return NopsView{}, err
		}
		if nodeDetails == nil || nodeDetails.Node == nil {
			return NopsView{}, fmt.Errorf("failed to get node details from offchain client for node %s", node.NodeID)
		}
		nodeName := nodeDetails.Node.Name
		if nodeName == "" {
			nodeName = node.NodeID
		}
		nop := NopView{
			NodeID:       node.NodeID,
			IsBootstrap:  node.IsBootstrap,
			OCRKeys:      make(map[string]OCRKeyView),
			PayeeAddress: node.AdminAddr,
			CSAKey:       nodeDetails.Node.PublicKey,
			IsConnected:  nodeDetails.Node.IsConnected,
			IsEnabled:    nodeDetails.Node.IsEnabled,
		}
		for sel, ocrConfig := range node.SelToOCRConfig {
			chainid, err := chainsel.ChainIdFromSelector(sel)
			if err != nil {
				return nops, err
			}
			chainName, err := chainsel.NameFromChainId(chainid)
			if err != nil {
				return nops, err
			}
			nop.OCRKeys[chainName] = OCRKeyView{
				OffchainPublicKey:         fmt.Sprintf("%x", ocrConfig.OffchainPublicKey[:]),
				OnchainPublicKey:          fmt.Sprintf("%x", ocrConfig.OnchainPublicKey[:]),
				PeerID:                    ocrConfig.PeerID.String(),
				TransmitAccount:           string(ocrConfig.TransmitAccount),
				ConfigEncryptionPublicKey: fmt.Sprintf("%x", ocrConfig.ConfigEncryptionPublicKey[:]),
				KeyBundleID:               ocrConfig.KeyBundleID,
			}
		}
		nops.Nops[nodeName] = nop
	}
	return nops, nil
}
