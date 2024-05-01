package graph_test

import (
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestGraph(t *testing.T) {
	const numNetworks = 100

	g := graph.NewGraph()
	assert.True(t, g.IsEmpty(), "should be empty")
	for i := 0; i < numNetworks; i++ {
		assert.True(t, g.AddNetwork(models.NetworkSelector(i), graph.Data{
			Liquidity: big.NewInt(int64(i * 100)),
		}))
	}
	assert.False(t, g.IsEmpty(), "should not be empty")

	allNetworks := g.GetNetworks()
	assert.Len(t, allNetworks, numNetworks, "networks length should match")

	for i, net := range allNetworks {
		liq, err := g.GetLiquidity(net)
		assert.NoError(t, err)
		assert.Equal(t, int64(i*100), liq.Int64(), "liquidity should match the initial liquidity")
	}

	// add network that already exists returns false
	netSel := models.NetworkSelector(numNetworks - 1)
	assert.False(t, g.AddNetwork(netSel, graph.Data{
		Liquidity: big.NewInt(123),
	}))

	// network does not exist
	liq, err := g.GetLiquidity(models.NetworkSelector(numNetworks + 1))
	assert.Empty(t, liq)
	assert.Error(t, err, "the provided network does not exist should get an error")

	// add some connections between networks and overwrite liquidity
	netSel1 := rand.Intn(numNetworks)
	netSel2 := rand.Intn(numNetworks)
	for netSel2 == netSel1 {
		netSel2 = rand.Intn(numNetworks)
	}
	assert.NoError(t, g.AddConnection(models.NetworkSelector(netSel1), models.NetworkSelector(netSel2)))
	assert.Error(t, g.AddConnection(models.NetworkSelector(numNetworks+1), models.NetworkSelector(numNetworks+2)))

	assert.True(t, g.SetLiquidity(models.NetworkSelector(netSel2), big.NewInt(999)))
	liq, err = g.GetLiquidity(models.NetworkSelector(netSel2))
	assert.NoError(t, err)
	assert.Equal(t, int64(999), liq.Int64())
	assert.False(t, g.SetLiquidity(models.NetworkSelector(999), big.NewInt(999)), "non-existent network")

	g.Reset()
	assert.True(t, g.IsEmpty())
}

func TestNewGraphFromEdges(t *testing.T) {
	var edges []models.Edge
	g, err := graph.NewGraphFromEdges(edges)
	assert.NoError(t, err)
	assert.True(t, g.IsEmpty())

	edges = append(edges, models.NewEdge(models.NetworkSelector(1), models.NetworkSelector(2)))
	g, err = graph.NewGraphFromEdges(edges)
	assert.NoError(t, err)
	assert.False(t, g.IsEmpty())
	neibs, ok := g.GetNeighbors(models.NetworkSelector(1))
	assert.True(t, ok)
	assert.Len(t, neibs, 1)
	assert.Equal(t, models.NetworkSelector(2), neibs[0])

	edges = append(edges, models.NewEdge(models.NetworkSelector(1), models.NetworkSelector(3)))
	g, err = graph.NewGraphFromEdges(edges)
	assert.NoError(t, err)
	neibs, ok = g.GetNeighbors(models.NetworkSelector(1))
	assert.True(t, ok)
	assert.Len(t, neibs, 2)
	assert.Equal(t, models.NetworkSelector(2), neibs[0])
	assert.Equal(t, models.NetworkSelector(3), neibs[1])
}

func TestGraphThreadSafety(t *testing.T) {
	const numWorkers = 50
	const numNetworks = 30

	g := graph.NewGraph()
	for i := 0; i < numNetworks; i++ {
		g.AddNetwork(models.NetworkSelector(i), graph.Data{
			Liquidity: big.NewInt(int64(i * 100)),
		})
	}

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			runGraphOperations(t, numNetworks, g)
		}()
	}
	wg.Wait()
}

// runGraphOperations runs some operations on the provided graph.
// Those operations are intended to be ran concurrently by multiple goroutines to test
// asynchronous behaviour and thready safety.
func runGraphOperations(t *testing.T, numNetworks int, g graph.Graph) {
	g.GetNetworks()
	assert.True(t, g.HasNetwork(models.NetworkSelector(numNetworks-3)))
	assert.False(t, g.HasNetwork(models.NetworkSelector(numNetworks+1234)))
	assert.False(t, g.IsEmpty())
	newNetID := models.NetworkSelector(rand.Intn(numNetworks * 3))
	g.AddNetwork(newNetID, graph.Data{
		Liquidity: big.NewInt(9999),
	})
	_, err := g.GetLiquidity(newNetID)
	assert.NoError(t, err)
	g.SetLiquidity(newNetID, big.NewInt(1234))
	_, err = g.GetLiquidity(newNetID)
	assert.NoError(t, err)
	_ = g.AddConnection(models.NetworkSelector(1), models.NetworkSelector(2))
}

func TestXChainRebalancerData_Equals(t *testing.T) {
	type fields struct {
		RemoteRebalancerAddress   models.Address
		LocalBridgeAdapterAddress models.Address
		RemoteTokenAddress        models.Address
	}
	type args struct {
		other graph.XChainRebalancerData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"equal",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: graph.XChainRebalancerData{
					RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
					LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
					RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
				},
			},
			true,
		},
		{
			"not equal remote rebalancer",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: graph.XChainRebalancerData{
					RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x4")),
					LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
					RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
				},
			},
			false,
		},
		{
			"not equal local bridge",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: graph.XChainRebalancerData{
					RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
					LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
					RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
				},
			},
			false,
		},
		{
			"not equal remote token",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: graph.XChainRebalancerData{
					RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
					LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
					RemoteTokenAddress:        models.Address(common.HexToAddress("0x4")),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := graph.XChainRebalancerData{
				RemoteRebalancerAddress:   tt.fields.RemoteRebalancerAddress,
				LocalBridgeAdapterAddress: tt.fields.LocalBridgeAdapterAddress,
				RemoteTokenAddress:        tt.fields.RemoteTokenAddress,
			}
			got := d.Equals(tt.args.other)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestData_Equals(t *testing.T) {
	type fields struct {
		Liquidity         *big.Int
		TokenAddress      models.Address
		RebalancerAddress models.Address
		XChainRebalancers map[models.NetworkSelector]graph.XChainRebalancerData
		ConfigDigest      models.ConfigDigest
		NetworkSelector   models.NetworkSelector
	}
	type args struct {
		other graph.Data
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"equal",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
					models.NetworkSelector(1): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: graph.Data{
					Liquidity:         big.NewInt(100),
					TokenAddress:      models.Address(common.HexToAddress("0x1")),
					RebalancerAddress: models.Address(common.HexToAddress("0x2")),
					XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
						models.NetworkSelector(1): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			true,
		},
		{
			"not equal liquidity",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
					models.NetworkSelector(1): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: graph.Data{
					Liquidity:         big.NewInt(200),
					TokenAddress:      models.Address(common.HexToAddress("0x1")),
					RebalancerAddress: models.Address(common.HexToAddress("0x2")),
					XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
						models.NetworkSelector(1): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
		{
			"not equal token address",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
					models.NetworkSelector(1): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: graph.Data{
					Liquidity:         big.NewInt(100),
					TokenAddress:      models.Address(common.HexToAddress("0x22")),
					RebalancerAddress: models.Address(common.HexToAddress("0x2")),
					XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
						models.NetworkSelector(1): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
		{
			"not equal rebalancer address",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
					models.NetworkSelector(1): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: graph.Data{
					Liquidity:         big.NewInt(100),
					TokenAddress:      models.Address(common.HexToAddress("0x1")),
					RebalancerAddress: models.Address(common.HexToAddress("0x222")),
					XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
						models.NetworkSelector(1): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
		{
			"not equal xchain rebalancers",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
					models.NetworkSelector(1): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: graph.Data{
					Liquidity:         big.NewInt(100),
					TokenAddress:      models.Address(common.HexToAddress("0x1")),
					RebalancerAddress: models.Address(common.HexToAddress("0x222")),
					XChainRebalancers: map[models.NetworkSelector]graph.XChainRebalancerData{
						models.NetworkSelector(1): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x33")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:        models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := graph.Data{
				Liquidity:         tt.fields.Liquidity,
				TokenAddress:      tt.fields.TokenAddress,
				RebalancerAddress: tt.fields.RebalancerAddress,
				XChainRebalancers: tt.fields.XChainRebalancers,
				ConfigDigest:      tt.fields.ConfigDigest,
				NetworkSelector:   tt.fields.NetworkSelector,
			}
			if got := d.Equals(tt.args.other); got != tt.want {
				t.Errorf("Data.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraph_Equals(t *testing.T) {
	type fields struct {
		genGraph func() graph.Graph
	}
	type args struct {
		other graph.Graph
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"not equal, diff lengths",
			fields{
				genGraph: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(100),
					})
					return g
				},
			},
			args{
				other: graph.NewGraph(),
			},
			false,
		},
		{
			"not equal, diff networks",
			fields{
				genGraph: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(100),
					})
					return g
				},
			},
			args{
				other: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(2), graph.Data{
						Liquidity: big.NewInt(100),
					})
					return g
				}(),
			},
			false,
		},
		{
			"not equal, diff datas",
			fields{
				genGraph: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(100),
					})
					return g
				},
			},
			args{
				other: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(200),
					})
					return g
				}(),
			},
			false,
		},
		{
			"not equal, diff neighbors",
			fields{
				genGraph: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(100),
					})
					g.AddNetwork(models.NetworkSelector(2), graph.Data{
						Liquidity: big.NewInt(200),
					})
					require.NoError(t, g.AddConnection(models.NetworkSelector(1), models.NetworkSelector(2)))
					return g
				},
			},
			args{
				other: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(100),
					})
					g.AddNetwork(models.NetworkSelector(2), graph.Data{
						Liquidity: big.NewInt(200),
					})
					require.NoError(t, g.AddConnection(models.NetworkSelector(2), models.NetworkSelector(1))) // reverse connection
					return g
				}(),
			},
			false,
		},
		{
			"equal",
			fields{
				genGraph: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(100),
					})
					g.AddNetwork(models.NetworkSelector(2), graph.Data{
						Liquidity: big.NewInt(200),
					})
					g.AddNetwork(models.NetworkSelector(3), graph.Data{
						Liquidity: big.NewInt(300),
					})
					require.NoError(t, g.AddConnection(models.NetworkSelector(1), models.NetworkSelector(2)))
					require.NoError(t, g.AddConnection(models.NetworkSelector(1), models.NetworkSelector(3)))
					require.NoError(t, g.AddConnection(models.NetworkSelector(2), models.NetworkSelector(3)))
					return g
				},
			},
			args{
				other: func() graph.Graph {
					g := graph.NewGraph()
					g.AddNetwork(models.NetworkSelector(1), graph.Data{
						Liquidity: big.NewInt(100),
					})
					g.AddNetwork(models.NetworkSelector(2), graph.Data{
						Liquidity: big.NewInt(200),
					})
					g.AddNetwork(models.NetworkSelector(3), graph.Data{
						Liquidity: big.NewInt(300),
					})
					require.NoError(t, g.AddConnection(models.NetworkSelector(1), models.NetworkSelector(2)))
					require.NoError(t, g.AddConnection(models.NetworkSelector(1), models.NetworkSelector(3)))
					require.NoError(t, g.AddConnection(models.NetworkSelector(2), models.NetworkSelector(3)))
					return g
				}(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.fields.genGraph()
			got := g.Equals(tt.args.other)
			require.Equal(t, tt.want, got)
		})
	}
}
