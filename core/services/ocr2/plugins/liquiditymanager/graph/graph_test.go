package graph

import (
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestGraph(t *testing.T) {
	const numNetworks = 100

	g := NewGraph()
	assert.True(t, g.IsEmpty(), "should be empty")
	for i := 0; i < numNetworks; i++ {
		assert.True(t, g.(GraphTest).AddNetwork(models.NetworkSelector(i), Data{
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
	assert.False(t, g.(GraphTest).AddNetwork(netSel, Data{
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
	assert.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(netSel1), models.NetworkSelector(netSel2)))
	assert.Error(t, g.(GraphTest).AddConnection(models.NetworkSelector(numNetworks+1), models.NetworkSelector(numNetworks+2)))

	assert.True(t, g.SetLiquidity(models.NetworkSelector(netSel2), big.NewInt(999)))
	liq, err = g.GetLiquidity(models.NetworkSelector(netSel2))
	assert.NoError(t, err)
	assert.Equal(t, int64(999), liq.Int64())
	assert.False(t, g.SetLiquidity(models.NetworkSelector(999), big.NewInt(999)), "non-existent network")

	g.Reset()
	assert.True(t, g.IsEmpty())
}

func TestGraphThreadSafety(t *testing.T) {
	const numWorkers = 50
	const numNetworks = 30

	g := NewGraph()
	for i := 0; i < numNetworks; i++ {
		g.(GraphTest).AddNetwork(models.NetworkSelector(i), Data{
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

// runGraphOperations runs some operations on the provided
// Those operations are intended to be ran concurrently by multiple goroutines to test
// asynchronous behaviour and thready safety.
func runGraphOperations(t *testing.T, numNetworks int, g Graph) {
	g.GetNetworks()
	assert.True(t, g.(GraphTest).HasNetwork(models.NetworkSelector(numNetworks-3)))
	assert.False(t, g.(GraphTest).HasNetwork(models.NetworkSelector(numNetworks+1234)))
	assert.False(t, g.IsEmpty())
	newNetID := models.NetworkSelector(rand.Intn(numNetworks * 3))
	g.(GraphTest).AddNetwork(newNetID, Data{
		Liquidity: big.NewInt(9999),
	})
	_, err := g.GetLiquidity(newNetID)
	assert.NoError(t, err)
	g.SetLiquidity(newNetID, big.NewInt(1234))
	_, err = g.GetLiquidity(newNetID)
	assert.NoError(t, err)
	_ = g.(GraphTest).AddConnection(models.NetworkSelector(1), models.NetworkSelector(2))
}

func TestGraph_Equals(t *testing.T) {
	type fields struct {
		genGraph func() Graph
	}
	type args struct {
		other Graph
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
				genGraph: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
						Liquidity: big.NewInt(100),
					})
					return g
				},
			},
			args{
				other: NewGraph(),
			},
			false,
		},
		{
			"not equal, diff networks",
			fields{
				genGraph: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
						Liquidity: big.NewInt(100),
					})
					return g
				},
			},
			args{
				other: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(2), Data{
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
				genGraph: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
						Liquidity: big.NewInt(100),
					})
					return g
				},
			},
			args{
				other: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
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
				genGraph: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
						Liquidity: big.NewInt(100),
					})
					g.(GraphTest).AddNetwork(models.NetworkSelector(2), Data{
						Liquidity: big.NewInt(200),
					})
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(1), models.NetworkSelector(2)))
					return g
				},
			},
			args{
				other: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
						Liquidity: big.NewInt(100),
					})
					g.(GraphTest).AddNetwork(models.NetworkSelector(2), Data{
						Liquidity: big.NewInt(200),
					})
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(2), models.NetworkSelector(1))) // reverse connection
					return g
				}(),
			},
			false,
		},
		{
			"equal",
			fields{
				genGraph: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
						Liquidity: big.NewInt(100),
					})
					g.(GraphTest).AddNetwork(models.NetworkSelector(2), Data{
						Liquidity: big.NewInt(200),
					})
					g.(GraphTest).AddNetwork(models.NetworkSelector(3), Data{
						Liquidity: big.NewInt(300),
					})
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(1), models.NetworkSelector(2)))
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(1), models.NetworkSelector(3)))
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(2), models.NetworkSelector(3)))
					return g
				},
			},
			args{
				other: func() Graph {
					g := NewGraph()
					g.(GraphTest).AddNetwork(models.NetworkSelector(1), Data{
						Liquidity: big.NewInt(100),
					})
					g.(GraphTest).AddNetwork(models.NetworkSelector(2), Data{
						Liquidity: big.NewInt(200),
					})
					g.(GraphTest).AddNetwork(models.NetworkSelector(3), Data{
						Liquidity: big.NewInt(300),
					})
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(1), models.NetworkSelector(2)))
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(1), models.NetworkSelector(3)))
					require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(2), models.NetworkSelector(3)))
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
