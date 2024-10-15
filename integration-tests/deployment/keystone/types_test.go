package keystone

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	v1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
)

func Test_newOcr2Node(t *testing.T) {
	type args struct {
		id        string
		ccfgs     map[chaintype.ChainType]*v1.ChainConfig
		csaPubKey string
	}
	tests := []struct {
		name      string
		args      args
		wantAptos bool
		wantErr   bool
	}{
		{
			name: "no aptos",
			args: args{
				id: "1",
				ccfgs: map[chaintype.ChainType]*v1.ChainConfig{
					chaintype.EVM: {

						Ocr2Config: &v1.OCR2Config{
							P2PKeyBundle: &v1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWMWUKdoAc2ruZf9f55p7NVFj7AFiPm67xjQ8BZBwkqyYv",
								PublicKey: "pubKey",
							},
							OcrKeyBundle: &v1.OCR2Config_OCRKeyBundle{
								BundleId:              "bundleId",
								ConfigPublicKey:       "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OffchainPublicKey:     "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OnchainSigningAddress: "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442",
							},
						},
					},
				},
				csaPubKey: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			},
		},
		{
			name: "with aptos",
			args: args{
				id: "1",
				ccfgs: map[chaintype.ChainType]*v1.ChainConfig{
					chaintype.EVM: {

						Ocr2Config: &v1.OCR2Config{
							P2PKeyBundle: &v1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWMWUKdoAc2ruZf9f55p7NVFj7AFiPm67xjQ8BZBwkqyYv",
								PublicKey: "pubKey",
							},
							OcrKeyBundle: &v1.OCR2Config_OCRKeyBundle{
								BundleId:              "bundleId",
								ConfigPublicKey:       "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OffchainPublicKey:     "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OnchainSigningAddress: "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442",
							},
						},
					},
					chaintype.Aptos: {

						Ocr2Config: &v1.OCR2Config{
							P2PKeyBundle: &v1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWMWUKdoAc2ruZf9f55p7NVFj7AFiPm67xjQ8BZB11111",
								PublicKey: "pubKey",
							},
							OcrKeyBundle: &v1.OCR2Config_OCRKeyBundle{
								BundleId:              "bundleId2",
								ConfigPublicKey:       "0000015fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OffchainPublicKey:     "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OnchainSigningAddress: "111409a8d4f9a18da55c5b2bb08a3f5f68d44777",
							},
						},
					},
				},
				csaPubKey: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			},
			wantAptos: true,
		},
		{
			name: "bad csa key",
			args: args{
				id: "1",
				ccfgs: map[chaintype.ChainType]*v1.ChainConfig{
					chaintype.EVM: {

						Ocr2Config: &v1.OCR2Config{
							P2PKeyBundle: &v1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWMWUKdoAc2ruZf9f55p7NVFj7AFiPm67xjQ8BZBwkqyYv",
								PublicKey: "pubKey",
							},
							OcrKeyBundle: &v1.OCR2Config_OCRKeyBundle{
								BundleId:              "bundleId",
								ConfigPublicKey:       "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OffchainPublicKey:     "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1",
								OnchainSigningAddress: "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442",
							},
						},
					},
				},
				csaPubKey: "not hex",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newOcr2Node(tt.args.id, tt.args.ccfgs, tt.args.csaPubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("newOcr2Node() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			assert.NotNil(t, got.ethOcr2KeyBundle)
			assert.NotNil(t, got.p2pKeyBundle)
			assert.NotNil(t, got.Signer)
			assert.NotNil(t, got.EncryptionPublicKey)
			assert.NotEmpty(t, got.csaKey)
			assert.NotEmpty(t, got.P2PKey)
			assert.Equal(t, tt.wantAptos, got.aptosOcr2KeyBundle != nil)
		})
	}
}

func Test_mapDonsToNodes(t *testing.T) {
	var (
		pubKey   = "03dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1"
		evmSig   = "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442"
		aptosSig = "b35409a8d4f9a18da55c5b2bb08a3f5f68d44442b35409a8d4f9a18da55c5b2bb08a3f5f68d44442"
		peerID   = "p2p_12D3KooWMWUKdoAc2ruZf9f55p7NVFj7AFiPm67xjQ8BZBwkqyYv"
		// todo: these should be defined in common
		writerCap = 3
		ocr3Cap   = 2
	)
	type args struct {
		dons              []DonCapabilities
		excludeBootstraps bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "writer evm only",
			args: args{
				dons: []DonCapabilities{
					{
						Name: "ok writer",
						Nops: []*models.NodeOperator{
							{
								Nodes: []*models.Node{
									{
										PublicKey: &pubKey,
										ChainConfigs: []*models.NodeChainConfig{
											{
												ID: "1",
												Network: &models.Network{
													ChainType: models.ChainTypeEvm,
												},
												Ocr2Config: &models.NodeOCR2Config{
													P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
														PeerID: peerID,
													},
													OcrKeyBundle: &models.NodeOCR2ConfigOCRKeyBundle{
														ConfigPublicKey:       pubKey,
														OffchainPublicKey:     pubKey,
														OnchainSigningAddress: evmSig,
													},
												},
											},
										},
									},
								},
							},
						},
						Capabilities: []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "writer",
								Version:        "1",
								CapabilityType: uint8(writerCap),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "err if no evm chain",
			args: args{
				dons: []DonCapabilities{
					{
						Name: "bad chain",
						Nops: []*models.NodeOperator{
							{
								Nodes: []*models.Node{
									{
										PublicKey: &pubKey,
										ChainConfigs: []*models.NodeChainConfig{
											{
												ID: "1",
												Network: &models.Network{
													ChainType: models.ChainTypeSolana,
												},
												Ocr2Config: &models.NodeOCR2Config{
													P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
														PeerID: peerID,
													},
													OcrKeyBundle: &models.NodeOCR2ConfigOCRKeyBundle{
														ConfigPublicKey:       pubKey,
														OffchainPublicKey:     pubKey,
														OnchainSigningAddress: evmSig,
													},
												},
											},
										},
									},
								},
							},
						},
						Capabilities: []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "writer",
								Version:        "1",
								CapabilityType: uint8(writerCap),
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ocr3 cap evm only",
			args: args{
				dons: []DonCapabilities{
					{
						Name: "bad chain",
						Nops: []*models.NodeOperator{
							{
								Nodes: []*models.Node{
									{
										PublicKey: &pubKey,
										ChainConfigs: []*models.NodeChainConfig{
											{
												ID: "1",
												Network: &models.Network{
													ChainType: models.ChainTypeEvm,
												},
												Ocr2Config: &models.NodeOCR2Config{
													P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
														PeerID: peerID,
													},
													OcrKeyBundle: &models.NodeOCR2ConfigOCRKeyBundle{
														ConfigPublicKey:       pubKey,
														OffchainPublicKey:     pubKey,
														OnchainSigningAddress: evmSig,
													},
												},
											},
										},
									},
								},
							},
						},
						Capabilities: []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "ocr3",
								Version:        "1",
								CapabilityType: uint8(ocr3Cap),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ocr3 cap evm & aptos",
			args: args{
				dons: []DonCapabilities{
					{
						Name: "bad chain",
						Nops: []*models.NodeOperator{
							{
								Nodes: []*models.Node{
									{
										PublicKey: &pubKey,
										ChainConfigs: []*models.NodeChainConfig{
											{
												ID: "1",
												Network: &models.Network{
													ChainType: models.ChainTypeEvm,
												},
												Ocr2Config: &models.NodeOCR2Config{
													P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
														PeerID: peerID,
													},
													OcrKeyBundle: &models.NodeOCR2ConfigOCRKeyBundle{
														ConfigPublicKey:       pubKey,
														OffchainPublicKey:     pubKey,
														OnchainSigningAddress: evmSig,
													},
												},
											},
											{
												ID: "2",
												Network: &models.Network{
													ChainType: models.ChainTypeAptos,
												},
												Ocr2Config: &models.NodeOCR2Config{
													P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
														PeerID: peerID,
													},
													OcrKeyBundle: &models.NodeOCR2ConfigOCRKeyBundle{
														ConfigPublicKey:       pubKey,
														OffchainPublicKey:     pubKey,
														OnchainSigningAddress: aptosSig,
													},
												},
											},
										},
									},
								},
							},
						},
						Capabilities: []kcr.CapabilitiesRegistryCapability{
							{
								LabelledName:   "ocr3",
								Version:        "1",
								CapabilityType: uint8(ocr3Cap),
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mapDonsToNodes(tt.args.dons, tt.args.excludeBootstraps)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapDonsToNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// make sure the clo test data is correct
	wfNops := loadTestNops(t, "testdata/workflow_nodes.json")
	cwNops := loadTestNops(t, "testdata/chain_writer_nodes.json")
	assetNops := loadTestNops(t, "testdata/asset_nodes.json")
	require.Len(t, wfNops, 10)
	require.Len(t, cwNops, 10)
	require.Len(t, assetNops, 16)

	wfDon := DonCapabilities{
		Name:         WFDonName,
		Nops:         wfNops,
		Capabilities: []kcr.CapabilitiesRegistryCapability{OCR3Cap},
	}
	cwDon := DonCapabilities{
		Name:         TargetDonName,
		Nops:         cwNops,
		Capabilities: []kcr.CapabilitiesRegistryCapability{WriteChainCap},
	}
	assetDon := DonCapabilities{
		Name:         StreamDonName,
		Nops:         assetNops,
		Capabilities: []kcr.CapabilitiesRegistryCapability{StreamTriggerCap},
	}
	_, err := mapDonsToNodes([]DonCapabilities{wfDon}, false)
	require.NoError(t, err, "failed to map wf don")
	_, err = mapDonsToNodes([]DonCapabilities{cwDon}, false)
	require.NoError(t, err, "failed to map cw don")
	_, err = mapDonsToNodes([]DonCapabilities{assetDon}, false)
	require.NoError(t, err, "failed to map asset don")
}

func loadTestNops(t *testing.T, pth string) []*models.NodeOperator {
	f, err := os.ReadFile(pth)
	require.NoError(t, err)
	var nops []*models.NodeOperator
	require.NoError(t, json.Unmarshal(f, &nops))
	return nops
}
