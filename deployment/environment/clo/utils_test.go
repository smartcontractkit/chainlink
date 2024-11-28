package clo

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/environment/clo/models"
	"github.com/stretchr/testify/assert"
)

func TestSetNodeIdsToPeerIds(t *testing.T) {
	type args struct {
		nops []*models.NodeOperator
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.NodeOperator
		wantErr bool
	}{
		{
			name: "no nodes",
			args: args{
				nops: []*models.NodeOperator{
					{
						ID: "nop1",
					},
				},
			},
			want: []*models.NodeOperator{
				{
					ID: "nop1",
				},
			},
		},
		{
			name: "error no p2p key bundle",
			args: args{
				nops: []*models.NodeOperator{
					{
						ID: "nop1",
						Nodes: []*models.Node{
							{
								ID: "node1",
								ChainConfigs: []*models.NodeChainConfig{
									{
										Ocr2Config: &models.NodeOCR2Config{},
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error multiple p2p key bundle",
			args: args{
				nops: []*models.NodeOperator{
					{
						ID: "nop1",
						Nodes: []*models.Node{
							{
								ID: "node1",
								ChainConfigs: []*models.NodeChainConfig{
									{
										Ocr2Config: &models.NodeOCR2Config{
											P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
												PeerID: "peer1",
											},
										},
									},
									{
										Ocr2Config: &models.NodeOCR2Config{
											P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
												PeerID: "peer2",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "multiple nodes",
			args: args{
				nops: []*models.NodeOperator{
					{
						ID: "nop1",
						Nodes: []*models.Node{
							{
								ID: "node1",
								ChainConfigs: []*models.NodeChainConfig{
									{
										Ocr2Config: &models.NodeOCR2Config{
											P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
												PeerID: "peer1",
											},
										},
									},
								},
							},
							{
								ID: "node2",
								ChainConfigs: []*models.NodeChainConfig{
									{
										Ocr2Config: &models.NodeOCR2Config{
											P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
												PeerID: "another peer id",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []*models.NodeOperator{
				{
					ID: "nop1",
					Nodes: []*models.Node{
						{
							ID: "peer1",
							ChainConfigs: []*models.NodeChainConfig{
								{
									Ocr2Config: &models.NodeOCR2Config{
										P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
											PeerID: "peer1",
										},
									},
								},
							},
						},
						{
							ID: "another peer id",
							ChainConfigs: []*models.NodeChainConfig{
								{
									Ocr2Config: &models.NodeOCR2Config{
										P2pKeyBundle: &models.NodeOCR2ConfigP2PKeyBundle{
											PeerID: "another peer id",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetNodeIdsToPeerIds(tt.args.nops)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetNodeIdsToPeerIds() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			assert.EqualValues(t, tt.args.nops, tt.want)
		})
	}
}
