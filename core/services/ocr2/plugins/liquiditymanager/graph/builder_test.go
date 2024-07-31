package graph

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestNewGraphFromEdges(t *testing.T) {
	var edges []models.Edge
	g, err := NewGraphFromEdges(edges)
	assert.NoError(t, err)
	assert.True(t, g.IsEmpty())

	edges = append(edges, models.NewEdge(models.NetworkSelector(1), models.NetworkSelector(2)))
	g, err = NewGraphFromEdges(edges)
	assert.NoError(t, err)
	assert.False(t, g.IsEmpty())
	neibs, ok := g.GetNeighbors(models.NetworkSelector(1), false)
	assert.True(t, ok)
	assert.Len(t, neibs, 1)
	assert.Equal(t, models.NetworkSelector(2), neibs[0])

	edges = append(edges, models.NewEdge(models.NetworkSelector(1), models.NetworkSelector(3)))
	g, err = NewGraphFromEdges(edges)
	assert.NoError(t, err)
	neibs, ok = g.GetNeighbors(models.NetworkSelector(1), false)
	assert.True(t, ok)
	assert.Len(t, neibs, 2)
	assert.Equal(t, models.NetworkSelector(2), neibs[0])
	assert.Equal(t, models.NetworkSelector(3), neibs[1])
}

func TestNewGraphWithData(t *testing.T) {
	type args struct {
		ctx           context.Context //nolint:containedctx
		startNetwork  models.NetworkSelector
		startAddress  models.Address
		getVertexInfo DataGetter
	}
	var (
		rebalNet1 = models.Address(common.HexToAddress("0x1"))
		rebalNet2 = models.Address(common.HexToAddress("0x2"))
		rebalNet3 = models.Address(common.HexToAddress("0x3"))
		rebalNet4 = models.Address(common.HexToAddress("0x4"))
	)
	tests := []struct {
		name    string
		args    args
		want    func() Graph
		wantErr bool
	}{
		{
			"1",
			// 1 is connected to 2 and 3
			// 2 is connected to 1
			// 3 is connected to 1
			args{
				ctx:          testutils.Context(t),
				startNetwork: models.NetworkSelector(1),
				startAddress: rebalNet1,
				getVertexInfo: func(ctx context.Context, v Vertex) (Data, []Vertex, error) {
					switch v.NetworkSelector {
					case 1:
						return Data{
								Liquidity:       big.NewInt(100),
								NetworkSelector: 1,
							}, []Vertex{
								{
									NetworkSelector:  2,
									LiquidityManager: rebalNet2,
								},
								{
									NetworkSelector:  3,
									LiquidityManager: rebalNet3,
								},
							}, nil
					case 2:
						return Data{
								Liquidity:       big.NewInt(200),
								NetworkSelector: 2,
							}, []Vertex{
								{
									NetworkSelector:  1,
									LiquidityManager: rebalNet1,
								},
							}, nil
					case 3:
						return Data{
								Liquidity:       big.NewInt(300),
								NetworkSelector: 3,
							}, []Vertex{
								{
									NetworkSelector:  1,
									LiquidityManager: rebalNet1,
								},
							}, nil
					default:
						return Data{}, nil, nil
					}
				},
			},
			func() Graph {
				g := NewGraph()
				d1 := Data{Liquidity: big.NewInt(100), NetworkSelector: 1}
				d2 := Data{Liquidity: big.NewInt(200), NetworkSelector: 2}
				d3 := Data{Liquidity: big.NewInt(300), NetworkSelector: 3}
				require.NoError(t, g.Add(d1, d2))
				require.NoError(t, g.Add(d1, d3))
				require.NoError(t, g.Add(d2, d1))
				require.NoError(t, g.Add(d3, d1))
				return g
			},
			false,
		},
		{
			"2",
			// 1 is connected to 2, 3 and 4
			// 2 is connected to 1 and 4
			// 3 is connected to 1, 2, and 4
			// 4 is connected to 1, 2, and 3
			args{
				ctx:          testutils.Context(t),
				startNetwork: models.NetworkSelector(1),
				startAddress: rebalNet1,
				getVertexInfo: func(ctx context.Context, v Vertex) (Data, []Vertex, error) {
					switch v.NetworkSelector {
					case 1:
						return Data{
								Liquidity:       big.NewInt(100),
								NetworkSelector: 1,
							}, []Vertex{
								{
									NetworkSelector:  2,
									LiquidityManager: rebalNet2,
								},
								{
									NetworkSelector:  3,
									LiquidityManager: rebalNet3,
								},
								{
									NetworkSelector:  4,
									LiquidityManager: rebalNet4,
								},
							}, nil
					case 2:
						return Data{
								Liquidity:       big.NewInt(200),
								NetworkSelector: 2,
							}, []Vertex{
								{
									NetworkSelector:  1,
									LiquidityManager: rebalNet1,
								},
								{
									NetworkSelector:  4,
									LiquidityManager: rebalNet4,
								},
							}, nil
					case 3:
						return Data{
								Liquidity:       big.NewInt(300),
								NetworkSelector: 3,
							}, []Vertex{
								{
									NetworkSelector:  1,
									LiquidityManager: rebalNet1,
								},
								{
									NetworkSelector:  2,
									LiquidityManager: rebalNet2,
								},
								{
									NetworkSelector:  4,
									LiquidityManager: rebalNet4,
								},
							}, nil
					case 4:
						return Data{
								Liquidity:       big.NewInt(400),
								NetworkSelector: 4,
							}, []Vertex{
								{
									NetworkSelector:  1,
									LiquidityManager: rebalNet1,
								},
								{
									NetworkSelector:  2,
									LiquidityManager: rebalNet2,
								},
								{
									NetworkSelector:  3,
									LiquidityManager: rebalNet3,
								},
							}, nil
					default:
						return Data{}, nil, nil
					}
				},
			},
			func() Graph {
				g := NewGraph()
				d1 := Data{Liquidity: big.NewInt(100), NetworkSelector: 1}
				d2 := Data{Liquidity: big.NewInt(200), NetworkSelector: 2}
				d3 := Data{Liquidity: big.NewInt(300), NetworkSelector: 3}
				d4 := Data{Liquidity: big.NewInt(400), NetworkSelector: 4}
				require.NoError(t, g.Add(d1, d2))
				require.NoError(t, g.Add(d1, d3))
				require.NoError(t, g.Add(d1, d4))
				require.NoError(t, g.Add(d2, d1))
				require.NoError(t, g.Add(d2, d4))
				require.NoError(t, g.Add(d3, d1))
				require.NoError(t, g.Add(d3, d2))
				require.NoError(t, g.Add(d3, d4))
				require.NoError(t, g.Add(d4, d1))
				require.NoError(t, g.Add(d4, d2))
				require.NoError(t, g.Add(d4, d3))
				return g
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGraphWithData(tt.args.ctx, Vertex{
				NetworkSelector:  tt.args.startNetwork,
				LiquidityManager: tt.args.startAddress,
			}, tt.args.getVertexInfo)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, tt.want().Equals(got))
			}
		})
	}
}
