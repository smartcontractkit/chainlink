package discoverer

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func Test_discover(t *testing.T) {
	type args struct {
		ctx           context.Context
		startNetwork  models.NetworkSelector
		startAddress  models.Address
		getVertexInfo func(ctx context.Context, network models.NetworkSelector, rebalancerAddress models.Address) (graph.Data, []dataItem, error)
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
		want    func() graph.Graph
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
				getVertexInfo: func(ctx context.Context, network models.NetworkSelector, rebalancerAddress models.Address) (graph.Data, []dataItem, error) {
					switch network {
					case 1:
						return graph.Data{
								Liquidity: big.NewInt(100),
							}, []dataItem{
								{
									networkSelector:   2,
									rebalancerAddress: rebalNet2,
								},
								{
									networkSelector:   3,
									rebalancerAddress: rebalNet3,
								},
							}, nil
					case 2:
						return graph.Data{
								Liquidity: big.NewInt(200),
							}, []dataItem{
								{
									networkSelector:   1,
									rebalancerAddress: rebalNet1,
								},
							}, nil
					case 3:
						return graph.Data{
								Liquidity: big.NewInt(300),
							}, []dataItem{
								{
									networkSelector:   1,
									rebalancerAddress: rebalNet1,
								},
							}, nil
					default:
						return graph.Data{}, nil, nil
					}
				},
			},
			func() graph.Graph {
				g := graph.NewGraph()
				g.AddNetwork(1, graph.Data{Liquidity: big.NewInt(100)})
				g.AddNetwork(2, graph.Data{Liquidity: big.NewInt(200)})
				g.AddNetwork(3, graph.Data{Liquidity: big.NewInt(300)})
				require.NoError(t, g.AddConnection(1, 2))
				require.NoError(t, g.AddConnection(1, 3))
				require.NoError(t, g.AddConnection(2, 1))
				require.NoError(t, g.AddConnection(3, 1))
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
				getVertexInfo: func(ctx context.Context, network models.NetworkSelector, rebalancerAddress models.Address) (graph.Data, []dataItem, error) {
					switch network {
					case 1:
						return graph.Data{
								Liquidity: big.NewInt(100),
							}, []dataItem{
								{
									networkSelector:   2,
									rebalancerAddress: rebalNet2,
								},
								{
									networkSelector:   3,
									rebalancerAddress: rebalNet3,
								},
								{
									networkSelector:   4,
									rebalancerAddress: rebalNet4,
								},
							}, nil
					case 2:
						return graph.Data{
								Liquidity: big.NewInt(200),
							}, []dataItem{
								{
									networkSelector:   1,
									rebalancerAddress: rebalNet1,
								},
								{
									networkSelector:   4,
									rebalancerAddress: rebalNet4,
								},
							}, nil
					case 3:
						return graph.Data{
								Liquidity: big.NewInt(300),
							}, []dataItem{
								{
									networkSelector:   1,
									rebalancerAddress: rebalNet1,
								},
								{
									networkSelector:   2,
									rebalancerAddress: rebalNet2,
								},
								{
									networkSelector:   4,
									rebalancerAddress: rebalNet4,
								},
							}, nil
					case 4:
						return graph.Data{

								Liquidity: big.NewInt(400),
							}, []dataItem{
								{
									networkSelector:   1,
									rebalancerAddress: rebalNet1,
								},
								{
									networkSelector:   2,
									rebalancerAddress: rebalNet2,
								},
								{
									networkSelector:   3,
									rebalancerAddress: rebalNet3,
								},
							}, nil
					default:
						return graph.Data{}, nil, nil
					}
				},
			},
			func() graph.Graph {
				g := graph.NewGraph()
				g.AddNetwork(1, graph.Data{Liquidity: big.NewInt(100)})
				g.AddNetwork(2, graph.Data{Liquidity: big.NewInt(200)})
				g.AddNetwork(3, graph.Data{Liquidity: big.NewInt(300)})
				g.AddNetwork(4, graph.Data{Liquidity: big.NewInt(400)})
				require.NoError(t, g.AddConnection(1, 2))
				require.NoError(t, g.AddConnection(1, 3))
				require.NoError(t, g.AddConnection(1, 4))
				require.NoError(t, g.AddConnection(2, 1))
				require.NoError(t, g.AddConnection(2, 4))
				require.NoError(t, g.AddConnection(3, 1))
				require.NoError(t, g.AddConnection(3, 2))
				require.NoError(t, g.AddConnection(3, 4))
				require.NoError(t, g.AddConnection(4, 1))
				require.NoError(t, g.AddConnection(4, 2))
				require.NoError(t, g.AddConnection(4, 3))
				return g
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := discover(tt.args.ctx, tt.args.startNetwork, tt.args.startAddress, tt.args.getVertexInfo)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, tt.want().Equals(got))
			}
		})
	}
}
