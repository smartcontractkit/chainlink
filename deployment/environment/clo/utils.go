package clo

import (
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
