package clo_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/test-go/testify/require"
	"google.golang.org/grpc"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var testNops = `
[
  {
    "id": "67",
    "name": "Chainlink Keystone Node Operator 9",
    "nodes": [
      {
        "id": "780",
        "name": "Chainlink Sepolia Prod Keystone One 9",
        "publicKey": "412dc6fe48ea4e34baaa77da2e3b032d39b938597b6f3d61fe7ed183a827a431",
        "connected": true,
        "supportedProducts": [
          "WORKFLOW",
          "OCR3_CAPABILITY"
        ]
      }
    ],
    "createdAt": "2024-08-14T19:00:07.113658Z"
  },
  {
    "id": "68",
    "name": "Chainlink Keystone Node Operator 8",
    "nodes": [
      {
        "id": "781",
        "name": "Chainlink Sepolia Prod Keystone One 8",
        "publicKey": "1141dd1e46797ced9b0fbad49115f18507f6f6e6e3cc86e7e5ba169e58645adc",
        "connected": true,
        "supportedProducts": [
          "WORKFLOW",
          "OCR3_CAPABILITY"
        ]
      }
    ],
    "createdAt": "2024-08-14T20:26:37.622463Z"
  },
  {
    "id": "999",
    "name": "Chainlink Keystone Node Operator 100",
    "nodes": [
      {
        "id": "999",
        "name": "Chainlink Sepolia Prod Keystone One 999",
        "publicKey": "9991dd1e46797ced9b0fbad49115f18507f6f6e6e3cc86e7e5ba169e58999999",
        "connected": true,
        "supportedProducts": [
          "WORKFLOW",
          "OCR3_CAPABILITY"
        ]
      },
      {
        "id": "1000",
        "name": "Chainlink Sepolia Prod Keystone One 1000",
        "publicKey": "1000101e46797ced9b0fbad49115f18507f6f6e6e3cc86e7e5ba169e58641000",
        "connected": true,
        "supportedProducts": [
          "WORKFLOW",
          "OCR3_CAPABILITY"
        ]
      }
    ],
    "createdAt": "2024-08-14T20:26:37.622463Z"
  }
]	
`

func parseTestNops(t *testing.T) []*models.NodeOperator {
	t.Helper()
	var out []*models.NodeOperator
	err := json.Unmarshal([]byte(testNops), &out)
	require.NoError(t, err)
	require.Len(t, out, 3, "wrong number of nops")
	return out
}
func TestJobClient_ListNodes(t *testing.T) {
	lggr := logger.TestLogger(t)
	nops := parseTestNops(t)

	type fields struct {
		NodeOperators []*models.NodeOperator
	}
	type args struct {
		ctx  context.Context
		in   *nodev1.ListNodesRequest
		opts []grpc.CallOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *nodev1.ListNodesResponse
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				NodeOperators: make([]*models.NodeOperator, 0),
			},
			args: args{
				ctx: context.Background(),
				in:  &nodev1.ListNodesRequest{},
			},
			want: &nodev1.ListNodesResponse{},
		},
		{
			name: "one node from one nop",
			fields: fields{
				NodeOperators: nops[0:1],
			},
			args: args{
				ctx: context.Background(),
				in:  &nodev1.ListNodesRequest{},
			},
			want: &nodev1.ListNodesResponse{
				Nodes: []*nodev1.Node{
					{
						Id:          "780",
						Name:        "Chainlink Sepolia Prod Keystone One 9",
						PublicKey:   "412dc6fe48ea4e34baaa77da2e3b032d39b938597b6f3d61fe7ed183a827a431",
						IsConnected: true,
					},
				},
			},
		},
		{
			name: "two nops each with one node",
			fields: fields{
				NodeOperators: nops[0:2],
			},
			args: args{
				ctx: context.Background(),
				in:  &nodev1.ListNodesRequest{},
			},
			want: &nodev1.ListNodesResponse{
				Nodes: []*nodev1.Node{
					{
						Id:          "780",
						Name:        "Chainlink Sepolia Prod Keystone One 9",
						PublicKey:   "412dc6fe48ea4e34baaa77da2e3b032d39b938597b6f3d61fe7ed183a827a431",
						IsConnected: true,
					},
					{
						Id:          "781",
						Name:        "Chainlink Sepolia Prod Keystone One 8",
						PublicKey:   "1141dd1e46797ced9b0fbad49115f18507f6f6e6e3cc86e7e5ba169e58645adc",
						IsConnected: true,
					},
				},
			},
		},
		{
			name: "two nodes from one nop",
			fields: fields{
				NodeOperators: nops[2:3],
			},
			args: args{
				ctx: context.Background(),
				in:  &nodev1.ListNodesRequest{},
			},
			want: &nodev1.ListNodesResponse{
				Nodes: []*nodev1.Node{
					{
						Id:          "999",
						Name:        "Chainlink Sepolia Prod Keystone One 999",
						PublicKey:   "9991dd1e46797ced9b0fbad49115f18507f6f6e6e3cc86e7e5ba169e58999999",
						IsConnected: true,
					},
					{
						Id:          "1000",
						Name:        "Chainlink Sepolia Prod Keystone One 1000",
						PublicKey:   "1000101e46797ced9b0fbad49115f18507f6f6e6e3cc86e7e5ba169e58641000",
						IsConnected: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := clo.NewJobClient(lggr, tt.fields.NodeOperators)

			got, err := j.ListNodes(tt.args.ctx, tt.args.in, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobClient.ListNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobClient.ListNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testNopsWithChainConfigs = `
[
  {
    "id": "67",
    "keys": [
      "keystone-09"
    ],
    "name": "Chainlink Keystone Node Operator 9",
    "metadata": {
      "nodeCount": 1,
      "jobCount": 4
    },
    "nodes": [
      {
        "id": "780",
        "name": "Chainlink Sepolia Prod Keystone One 9",
        "publicKey": "412dc6fe48ea4e34baaa77da2e3b032d39b938597b6f3d61fe7ed183a827a431",
        "chainConfigs": [
          {
            "network": {
              "id": "140",
              "chainID": "421614",
              "chainType": "EVM",
              "name": "Arbitrum Testnet (Sepolia)"
            },
            "accountAddress": "0xbA8E21dFaa0501fCD43146d0b5F21c2B8E0eEdfB",
            "adminAddress": "0x0000000000000000000000000000000000000000",
            "ocr1Config": {
              "p2pKeyBundle": {},
              "ocrKeyBundle": {}
            },
            "ocr2Config": {
              "enabled": true,
              "p2pKeyBundle": {
                "peerID": "p2p_12D3KooWBCMCCZZ8x57AXvJvpCujqhZzTjWXbReaRE8TxNr5dM4U",
                "publicKey": "147d5cc651819b093cd2fdff9760f0f0f77b7ef7798d9e24fc6a350b7300e5d9"
              },
              "ocrKeyBundle": {
                "bundleID": "1c28e76d180d1ed1524e61845fa58a384415de7e51017edf1f8c553e28357772",
                "configPublicKey": "09fced0207611ed618bf0759ab128d9797e15b18e46436be1a56a91e4043ec0e",
                "offchainPublicKey": "c805572b813a072067eab2087ddbee8aa719090e12890b15c01094f0d3f74a5f",
                "onchainSigningAddress": "679296b7c1eb4948efcc87efc550940a182e610c"
              },
              "plugins": {}
            }
          },
          {
            "network": {
              "id": "129",
              "chainID": "11155111",
              "chainType": "EVM",
              "name": "Ethereum Testnet (Sepolia)"
            },
            "accountAddress": "0x0b04cE574E80Da73191Ec141c0016a54A6404056",
            "adminAddress": "0x0000000000000000000000000000000000000000",
            "ocr1Config": {
              "p2pKeyBundle": {},
              "ocrKeyBundle": {}
            },
            "ocr2Config": {
              "enabled": true,
              "p2pKeyBundle": {
                "peerID": "p2p_12D3KooWBCMCCZZ8x57AXvJvpCujqhZzTjWXbReaRE8TxNr5dM4U",
                "publicKey": "147d5cc651819b093cd2fdff9760f0f0f77b7ef7798d9e24fc6a350b7300e5d9"
              },
              "ocrKeyBundle": {
                "bundleID": "1c28e76d180d1ed1524e61845fa58a384415de7e51017edf1f8c553e28357772",
                "configPublicKey": "09fced0207611ed618bf0759ab128d9797e15b18e46436be1a56a91e4043ec0e",
                "offchainPublicKey": "c805572b813a072067eab2087ddbee8aa719090e12890b15c01094f0d3f74a5f",
                "onchainSigningAddress": "679296b7c1eb4948efcc87efc550940a182e610c"
              },
              "plugins": {}
            }
          }
        ],
        "connected": true,
        "supportedProducts": [
          "WORKFLOW",
          "OCR3_CAPABILITY"
        ],
        "categories": [
          {
            "id": "11",
            "name": "Keystone"
          }
        ]
      }
    ],
    "createdAt": "2024-08-14T19:00:07.113658Z"
  },
  {
    "id": "68",
    "keys": [
      "keystone-08"
    ],
    "name": "Chainlink Keystone Node Operator 8",
    "metadata": {
      "nodeCount": 1,
      "jobCount": 4
    },
    "nodes": [
      {
        "id": "781",
        "name": "Chainlink Sepolia Prod Keystone One 8",
        "publicKey": "1141dd1e46797ced9b0fbad49115f18507f6f6e6e3cc86e7e5ba169e58645adc",
        "chainConfigs": [
          {
            "network": {
              "id": "140",
              "chainID": "421614",
              "chainType": "EVM",
              "name": "Arbitrum Testnet (Sepolia)"
            },
            "accountAddress": "0xEa4bC3638660D78Da56f39f6680dCDD0cEAaD2c6",
            "adminAddress": "0x0000000000000000000000000000000000000000",
            "ocr1Config": {
              "p2pKeyBundle": {},
              "ocrKeyBundle": {}
            },
            "ocr2Config": {
              "enabled": true,
              "p2pKeyBundle": {
                "peerID": "p2p_12D3KooWAUagqMycsro27kFznSQRHbhfCBLx8nKD4ptTiUGDe38c",
                "publicKey": "09ca39cd924653c72fbb0e458b629c3efebdad3e29e7cd0b5760754d919ed829"
              },
              "ocrKeyBundle": {
                "bundleID": "be0d639de3ae3cbeaa31ca369514f748ba1d271145cba6796bcc12aace2f64c3",
                "configPublicKey": "e3d4d7a7372a3b1110db0290ab3649eb5fbb0daf6cf3ae02cfe5f367700d9264",
                "offchainPublicKey": "ad08c2a5878cada53521f4e2bb449f191ccca7899246721a0deeea19f7b83f70",
                "onchainSigningAddress": "8c2aa1e6fad88a6006dfb116eb866cbad2910314"
              },
              "plugins": {}
            }
          },
          {
            "network": {
              "id": "129",
              "chainID": "11155111",
              "chainType": "EVM",
              "name": "Ethereum Testnet (Sepolia)"
            },
            "accountAddress": "0x31B179dcF8f9036C30f04bE578793e51bF14A39E",
            "adminAddress": "0x0000000000000000000000000000000000000000",
            "ocr1Config": {
              "p2pKeyBundle": {},
              "ocrKeyBundle": {}
            },
            "ocr2Config": {
              "enabled": true,
              "p2pKeyBundle": {
                "peerID": "p2p_12D3KooWAUagqMycsro27kFznSQRHbhfCBLx8nKD4ptTiUGDe38c",
                "publicKey": "09ca39cd924653c72fbb0e458b629c3efebdad3e29e7cd0b5760754d919ed829"
              },
              "ocrKeyBundle": {
                "bundleID": "be0d639de3ae3cbeaa31ca369514f748ba1d271145cba6796bcc12aace2f64c3",
                "configPublicKey": "e3d4d7a7372a3b1110db0290ab3649eb5fbb0daf6cf3ae02cfe5f367700d9264",
                "offchainPublicKey": "ad08c2a5878cada53521f4e2bb449f191ccca7899246721a0deeea19f7b83f70",
                "onchainSigningAddress": "8c2aa1e6fad88a6006dfb116eb866cbad2910314"
              },
              "plugins": {}
            }
          }
        ],
        "connected": true,
        "supportedProducts": [
          "WORKFLOW",
          "OCR3_CAPABILITY"
        ],
        "categories": [
          {
            "id": "11",
            "name": "Keystone"
          }
        ]
      }
    ],
    "createdAt": "2024-08-14T20:26:37.622463Z"
  }
]`

func TestJobClient_ListNodeChainConfigs(t *testing.T) {
	nops := parseTestNopsWithChainConfigs(t)
	lggr := logger.TestLogger(t)
	type fields struct {
		NodeOperators []*models.NodeOperator
	}
	type args struct {
		ctx  context.Context
		in   *nodev1.ListNodeChainConfigsRequest
		opts []grpc.CallOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *nodev1.ListNodeChainConfigsResponse
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				NodeOperators: make([]*models.NodeOperator, 0),
			},
			args: args{
				ctx: context.Background(),
				in:  &nodev1.ListNodeChainConfigsRequest{},
			},
			want: &nodev1.ListNodeChainConfigsResponse{
				ChainConfigs: make([]*nodev1.ChainConfig, 0),
			},
		},

		{
			name: "no matching nodes",
			fields: fields{
				NodeOperators: nops,
			},
			args: args{
				ctx: context.Background(),
				in: &nodev1.ListNodeChainConfigsRequest{
					Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
						NodeIds: []string{"not-a-node-id"},
					},
				},
			},
			want: &nodev1.ListNodeChainConfigsResponse{
				ChainConfigs: make([]*nodev1.ChainConfig, 0),
			},
		},

		{
			name: "one nop with one node that has two chain configs",
			fields: fields{
				NodeOperators: nops[0:1],
			},
			args: args{
				ctx: context.Background(),
				in:  &nodev1.ListNodeChainConfigsRequest{},
			},
			want: &nodev1.ListNodeChainConfigsResponse{
				ChainConfigs: []*nodev1.ChainConfig{
					{
						Chain: &nodev1.Chain{
							Id:   "421614",
							Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						},
						AccountAddress: "0xbA8E21dFaa0501fCD43146d0b5F21c2B8E0eEdfB",
						AdminAddress:   "0x0000000000000000000000000000000000000000",
						Ocr2Config: &nodev1.OCR2Config{
							Enabled: true,
							P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWBCMCCZZ8x57AXvJvpCujqhZzTjWXbReaRE8TxNr5dM4U",
								PublicKey: "147d5cc651819b093cd2fdff9760f0f0f77b7ef7798d9e24fc6a350b7300e5d9",
							},
							OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
								BundleId:              "1c28e76d180d1ed1524e61845fa58a384415de7e51017edf1f8c553e28357772",
								ConfigPublicKey:       "09fced0207611ed618bf0759ab128d9797e15b18e46436be1a56a91e4043ec0e",
								OffchainPublicKey:     "c805572b813a072067eab2087ddbee8aa719090e12890b15c01094f0d3f74a5f",
								OnchainSigningAddress: "679296b7c1eb4948efcc87efc550940a182e610c",
							},
						},
					},
					{
						Chain: &nodev1.Chain{
							Id:   "11155111",
							Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						},
						AccountAddress: "0x0b04cE574E80Da73191Ec141c0016a54A6404056",
						AdminAddress:   "0x0000000000000000000000000000000000000000",
						Ocr2Config: &nodev1.OCR2Config{
							Enabled: true,
							P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWBCMCCZZ8x57AXvJvpCujqhZzTjWXbReaRE8TxNr5dM4U",
								PublicKey: "147d5cc651819b093cd2fdff9760f0f0f77b7ef7798d9e24fc6a350b7300e5d9",
							},
							OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
								BundleId:              "1c28e76d180d1ed1524e61845fa58a384415de7e51017edf1f8c553e28357772",
								ConfigPublicKey:       "09fced0207611ed618bf0759ab128d9797e15b18e46436be1a56a91e4043ec0e",
								OffchainPublicKey:     "c805572b813a072067eab2087ddbee8aa719090e12890b15c01094f0d3f74a5f",
								OnchainSigningAddress: "679296b7c1eb4948efcc87efc550940a182e610c",
							},
						},
					},
				},
			},
		},

		{
			name: "one nop with one node that has two chain configs matching the filter",
			fields: fields{
				NodeOperators: nops,
			},
			args: args{
				ctx: context.Background(),
				in: &nodev1.ListNodeChainConfigsRequest{
					Filter: &nodev1.ListNodeChainConfigsRequest_Filter{
						NodeIds: []string{"780"},
					},
				},
			},
			want: &nodev1.ListNodeChainConfigsResponse{
				ChainConfigs: []*nodev1.ChainConfig{
					{
						Chain: &nodev1.Chain{
							Id:   "421614",
							Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						},
						AccountAddress: "0xbA8E21dFaa0501fCD43146d0b5F21c2B8E0eEdfB",
						AdminAddress:   "0x0000000000000000000000000000000000000000",
						Ocr2Config: &nodev1.OCR2Config{
							Enabled: true,
							P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWBCMCCZZ8x57AXvJvpCujqhZzTjWXbReaRE8TxNr5dM4U",
								PublicKey: "147d5cc651819b093cd2fdff9760f0f0f77b7ef7798d9e24fc6a350b7300e5d9",
							},
							OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
								BundleId:              "1c28e76d180d1ed1524e61845fa58a384415de7e51017edf1f8c553e28357772",
								ConfigPublicKey:       "09fced0207611ed618bf0759ab128d9797e15b18e46436be1a56a91e4043ec0e",
								OffchainPublicKey:     "c805572b813a072067eab2087ddbee8aa719090e12890b15c01094f0d3f74a5f",
								OnchainSigningAddress: "679296b7c1eb4948efcc87efc550940a182e610c",
							},
						},
					},
					{
						Chain: &nodev1.Chain{
							Id:   "11155111",
							Type: nodev1.ChainType_CHAIN_TYPE_EVM,
						},
						AccountAddress: "0x0b04cE574E80Da73191Ec141c0016a54A6404056",
						AdminAddress:   "0x0000000000000000000000000000000000000000",
						Ocr2Config: &nodev1.OCR2Config{
							Enabled: true,
							P2PKeyBundle: &nodev1.OCR2Config_P2PKeyBundle{
								PeerId:    "p2p_12D3KooWBCMCCZZ8x57AXvJvpCujqhZzTjWXbReaRE8TxNr5dM4U",
								PublicKey: "147d5cc651819b093cd2fdff9760f0f0f77b7ef7798d9e24fc6a350b7300e5d9",
							},
							OcrKeyBundle: &nodev1.OCR2Config_OCRKeyBundle{
								BundleId:              "1c28e76d180d1ed1524e61845fa58a384415de7e51017edf1f8c553e28357772",
								ConfigPublicKey:       "09fced0207611ed618bf0759ab128d9797e15b18e46436be1a56a91e4043ec0e",
								OffchainPublicKey:     "c805572b813a072067eab2087ddbee8aa719090e12890b15c01094f0d3f74a5f",
								OnchainSigningAddress: "679296b7c1eb4948efcc87efc550940a182e610c",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := clo.NewJobClient(lggr, tt.fields.NodeOperators)

			got, err := j.ListNodeChainConfigs(tt.args.ctx, tt.args.in, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobClient.ListNodeChainConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobClient.ListNodeChainConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseTestNopsWithChainConfigs(t *testing.T) []*models.NodeOperator {
	t.Helper()
	var out []*models.NodeOperator
	err := json.Unmarshal([]byte(testNopsWithChainConfigs), &out)
	require.NoError(t, err)
	require.Len(t, out, 2, "wrong number of nops")
	return out
}
