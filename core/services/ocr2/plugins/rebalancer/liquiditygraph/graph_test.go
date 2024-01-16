package liquiditygraph_test

import (
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

func TestGraph(t *testing.T) {
	const numNetworks = 100

	g := liquiditygraph.NewGraph()
	assert.True(t, g.IsEmpty(), "should be empty")
	for i := 0; i < numNetworks; i++ {
		assert.True(t, g.AddNetwork(models.NetworkSelector(i), big.NewInt(int64(i*100))))
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
	assert.False(t, g.AddNetwork(netSel, big.NewInt(123)))

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
	assert.True(t, g.AddConnection(models.NetworkSelector(netSel1), models.NetworkSelector(netSel2)))
	assert.False(t, g.AddConnection(models.NetworkSelector(numNetworks+1), models.NetworkSelector(numNetworks+2)))

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

	g := liquiditygraph.NewGraph()
	for i := 0; i < numNetworks; i++ {
		g.AddNetwork(models.NetworkSelector(i), big.NewInt(int64(i*100)))
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
func runGraphOperations(t *testing.T, numNetworks int, g liquiditygraph.LiquidityGraph) {
	g.GetNetworks()
	assert.True(t, g.HasNetwork(models.NetworkSelector(numNetworks-3)))
	assert.False(t, g.HasNetwork(models.NetworkSelector(numNetworks+1234)))
	assert.False(t, g.IsEmpty())
	newNetID := models.NetworkSelector(rand.Intn(numNetworks * 3))
	g.AddNetwork(newNetID, big.NewInt(9999))
	_, err := g.GetLiquidity(newNetID)
	assert.NoError(t, err)
	g.SetLiquidity(newNetID, big.NewInt(1234))
	_, err = g.GetLiquidity(newNetID)
	assert.NoError(t, err)
	g.AddConnection(models.NetworkSelector(1), models.NetworkSelector(2))
}
