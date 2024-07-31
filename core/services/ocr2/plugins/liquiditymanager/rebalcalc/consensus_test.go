package rebalcalc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestMedianLiquidityPerChain(t *testing.T) {
	type args struct {
		observations []models.Observation
		f            int
	}
	tests := []struct {
		name    string
		args    args
		want    []models.NetworkLiquidity
		wantErr bool
	}{
		{
			"no observations",
			args{[]models.Observation{}, 1},
			[]models.NetworkLiquidity{},
			true,
		},
		{
			"single observation",
			args{[]models.Observation{
				{},
			}, 1},
			[]models.NetworkLiquidity{},
			true,
		},
		{
			"multiple observations",
			args{[]models.Observation{
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: 1, Liquidity: ubig.NewI(1)},
						{Network: 2, Liquidity: ubig.NewI(2)},
					},
				},
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: 1, Liquidity: ubig.NewI(2)},
						{Network: 2, Liquidity: ubig.NewI(3)},
					},
				},
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: 1, Liquidity: ubig.NewI(3)},
						{Network: 2, Liquidity: ubig.NewI(4)},
					},
				},
			}, 1},
			[]models.NetworkLiquidity{
				{Network: 1, Liquidity: ubig.NewI(2)},
				{Network: 2, Liquidity: ubig.NewI(3)},
			},
			false,
		},
		{
			"below bft",
			args{[]models.Observation{
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: 1, Liquidity: ubig.NewI(1)},
						{Network: 2, Liquidity: ubig.NewI(2)},
					},
				},
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: 3, Liquidity: ubig.NewI(2)},
						{Network: 3, Liquidity: ubig.NewI(6)},
					},
				},
				{
					LiquidityPerChain: []models.NetworkLiquidity{
						{Network: 3, Liquidity: ubig.NewI(4)},
					},
				},
			}, 1},
			[]models.NetworkLiquidity{
				{Network: 3, Liquidity: ubig.NewI(4)},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MedianLiquidityPerChain(tt.args.observations, tt.args.f)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPendingTransfersConsensus(t *testing.T) {
	type args struct {
		observations []models.Observation
		f            int
	}
	tests := []struct {
		name        string
		args        args
		numExpected int
		wantErr     bool
	}{
		{
			"no observations",
			args{[]models.Observation{}, 1},
			0,
			true,
		},
		{
			"not enough observations",
			args{[]models.Observation{
				{}, {}, {},
			}, 2},
			0,
			true,
		},
		{
			"enough observations",
			args{[]models.Observation{
				{
					PendingTransfers: []models.PendingTransfer{
						{Transfer: models.Transfer{From: 1, To: 2, Amount: ubig.NewI(1)}},
						{Transfer: models.Transfer{From: 2, To: 3, Amount: ubig.NewI(2)}},
					},
				},
				{
					PendingTransfers: []models.PendingTransfer{
						{Transfer: models.Transfer{From: 1, To: 2, Amount: ubig.NewI(1)}},
						{Transfer: models.Transfer{From: 2, To: 3, Amount: ubig.NewI(2)}},
						{Transfer: models.Transfer{From: 3, To: 4, Amount: ubig.NewI(3)}}, // should not be included.
					},
				},
				{
					PendingTransfers: []models.PendingTransfer{
						{Transfer: models.Transfer{From: 1, To: 2, Amount: ubig.NewI(1)}},
						{Transfer: models.Transfer{From: 2, To: 3, Amount: ubig.NewI(2)}},
					},
				},
			}, 1},
			2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PendingTransfersConsensus(tt.args.observations, tt.args.f)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.numExpected, len(got))
			}
		})
	}
}

func TestConfigDigestsConsensus(t *testing.T) {
	type args struct {
		observations []models.Observation
		f            int
	}
	tests := []struct {
		name    string
		args    args
		want    []models.ConfigDigestWithMeta
		wantErr bool
	}{
		{
			"no observations",
			args{[]models.Observation{}, 1},
			[]models.ConfigDigestWithMeta{},
			true,
		},
		{
			"not enough observations",
			args{[]models.Observation{
				{}, {}, {},
			}, 2},
			[]models.ConfigDigestWithMeta{},
			true,
		},
		{
			"enough observations",
			args{[]models.Observation{
				{
					ConfigDigests: []models.ConfigDigestWithMeta{
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x1"))}, NetworkSel: 1},
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x2"))}, NetworkSel: 2},
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x3"))}, NetworkSel: 3},
					},
				},
				{
					ConfigDigests: []models.ConfigDigestWithMeta{
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x1"))}, NetworkSel: 1},
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x2"))}, NetworkSel: 2},
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x3"))}, NetworkSel: 3},
					},
				},
				{
					ConfigDigests: []models.ConfigDigestWithMeta{
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x1"))}, NetworkSel: 1},
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x2"))}, NetworkSel: 2},
						{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x3"))}, NetworkSel: 3},
					},
				},
			}, 1},
			[]models.ConfigDigestWithMeta{
				{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x1"))}, NetworkSel: 1},
				{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x2"))}, NetworkSel: 2},
				{Digest: models.ConfigDigest{ConfigDigest: ocrtypes.ConfigDigest(common.HexToHash("0x3"))}, NetworkSel: 3},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConfigDigestsConsensus(tt.args.observations, tt.args.f)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGraphEdgesConsensus(t *testing.T) {
	type args struct {
		observations []models.Observation
		f            int
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Edge
		wantErr bool
	}{
		{
			"no observations",
			args{[]models.Observation{}, 1},
			[]models.Edge{},
			true,
		},
		{
			"not enough observations",
			args{[]models.Observation{
				{}, {}, {},
			}, 2},
			[]models.Edge{},
			true,
		},
		{
			"enough observations",
			args{[]models.Observation{
				{
					Edges: []models.Edge{
						{Source: 1, Dest: 2},
						{Source: 2, Dest: 3},
					},
				},
				{
					Edges: []models.Edge{
						{Source: 1, Dest: 2},
						{Source: 2, Dest: 3},
						{Source: 3, Dest: 4}, // should not be included.
					},
				},
				{
					Edges: []models.Edge{
						{Source: 1, Dest: 2},
						{Source: 2, Dest: 3},
					},
				},
			}, 1},
			[]models.Edge{
				{Source: 1, Dest: 2},
				{Source: 2, Dest: 3},
			},
			false,
		}, {
			"differently ordered edges",
			args{[]models.Observation{
				{
					Edges: []models.Edge{
						{Source: 1, Dest: 4},
						{Source: 1, Dest: 1},
						{Source: 1, Dest: 3},
						{Source: 2, Dest: 1},
						{Source: 1, Dest: 2},
					},
				},
				{
					Edges: []models.Edge{
						{Source: 2, Dest: 1},
						{Source: 1, Dest: 4},
						{Source: 1, Dest: 3},
						{Source: 1, Dest: 2},
						{Source: 1, Dest: 1},
					},
				},
				{
					Edges: []models.Edge{
						{Source: 1, Dest: 2},
						{Source: 1, Dest: 4},
						{Source: 2, Dest: 1},
						{Source: 1, Dest: 1},
						{Source: 1, Dest: 3},
					},
				},
			}, 1},
			[]models.Edge{
				{Source: 1, Dest: 1},
				{Source: 1, Dest: 2},
				{Source: 1, Dest: 3},
				{Source: 1, Dest: 4},
				{Source: 2, Dest: 1},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GraphEdgesConsensus(tt.args.observations, tt.args.f)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInflightTransfersConsensus(t *testing.T) {
	type args struct {
		observations []models.Observation
		f            int
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Transfer
		wantErr bool
	}{
		{
			"no observations",
			args{[]models.Observation{}, 1},
			[]models.Transfer{},
			true,
		},
		{
			"not enough observations",
			args{[]models.Observation{
				{}, {}, {},
			}, 2},
			[]models.Transfer{},
			true,
		},
		{
			"enough observations",
			args{[]models.Observation{
				{
					InflightTransfers: []models.Transfer{
						{From: 1, To: 2, Amount: ubig.NewI(1)},
						{From: 2, To: 3, Amount: ubig.NewI(2)},
					},
				},
				{
					InflightTransfers: []models.Transfer{
						{From: 1, To: 2, Amount: ubig.NewI(1)},
						{From: 2, To: 3, Amount: ubig.NewI(2)},
						{From: 3, To: 4, Amount: ubig.NewI(3)}, // should not be included.
					},
				},
				{
					InflightTransfers: []models.Transfer{
						{From: 1, To: 2, Amount: ubig.NewI(1)},
						{From: 2, To: 3, Amount: ubig.NewI(2)},
					},
				},
			}, 1},
			[]models.Transfer{
				{From: 1, To: 2, Amount: ubig.NewI(1)},
				{From: 2, To: 3, Amount: ubig.NewI(2)},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InflightTransfersConsensus(tt.args.observations, tt.args.f)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
