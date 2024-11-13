package clo

import (
	"fmt"

	jd "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment/environment/clo/models"
)

// NewChainConfig creates a new JobDistributor ChainConfig from a clo model NodeChainConfig
func NewChainConfig(chain *models.NodeChainConfig) *jd.ChainConfig {
	return &jd.ChainConfig{
		Chain: &jd.Chain{
			Id:   chain.Network.ChainID,
			Type: jd.ChainType_CHAIN_TYPE_EVM, // TODO: support other chain types
		},

		AccountAddress: chain.AccountAddress,
		AdminAddress:   chain.AdminAddress,
		Ocr2Config: &jd.OCR2Config{
			Enabled: chain.Ocr2Config.Enabled,
			P2PKeyBundle: &jd.OCR2Config_P2PKeyBundle{
				PeerId:    chain.Ocr2Config.P2pKeyBundle.PeerID,
				PublicKey: chain.Ocr2Config.P2pKeyBundle.PublicKey,
			},
			OcrKeyBundle: &jd.OCR2Config_OCRKeyBundle{
				BundleId:              chain.Ocr2Config.OcrKeyBundle.BundleID,
				OnchainSigningAddress: chain.Ocr2Config.OcrKeyBundle.OnchainSigningAddress,
				OffchainPublicKey:     chain.Ocr2Config.OcrKeyBundle.OffchainPublicKey,
				ConfigPublicKey:       chain.Ocr2Config.OcrKeyBundle.ConfigPublicKey,
			},
		},
	}
}

func NodeP2PId(n *models.Node) (string, error) {
	p2pIds := make(map[string]struct{})
	for _, cc := range n.ChainConfigs {
		if cc.Ocr2Config != nil && cc.Ocr2Config.P2pKeyBundle != nil {
			p2pIds[cc.Ocr2Config.P2pKeyBundle.PeerID] = struct{}{}
		}
	}
	if len(p2pIds) == 0 {
		return "", fmt.Errorf("no p2p id found for node %s", n.ID)
	}
	if len(p2pIds) > 1 {
		return "", fmt.Errorf("multiple p2p ids found for node %s", n.ID)
	}
	var p2pId string
	for k := range p2pIds {
		p2pId = k
		break
	}
	return p2pId, nil
}

func NodesToPeerIDs(nodes []*models.Node) ([]string, error) {
	var p2pIds []string
	for _, node := range nodes {
		p2pId, err := NodeP2PId(node)
		if err != nil {
			return nil, err
		}
		p2pIds = append(p2pIds, p2pId)
	}
	return p2pIds, nil
}

func NopsToNodes(nops []*models.NodeOperator) []*models.Node {
	var nodes []*models.Node
	for _, nop := range nops {
		nodes = append(nodes, nop.Nodes...)
	}
	return nodes
}

func NopsToPeerIds(nops []*models.NodeOperator) ([]string, error) {
	return NodesToPeerIDs(NopsToNodes(nops))
}

func SetIdToPeerId(n *models.Node) error {
	p2pId, err := NodeP2PId(n)
	if err != nil {
		return err
	}
	n.ID = p2pId
	return nil
}

// SetNodeIdsToPeerIds sets the ID of each node in the NOPs to the P2P ID of the node
// It mutates the input NOPs
func SetNodeIdsToPeerIds(nops []*models.NodeOperator) error {
	for _, nop := range nops {
		for _, n := range nop.Nodes {
			if err := SetIdToPeerId(n); err != nil {
				return err
			}
		}
	}
	return nil
}
