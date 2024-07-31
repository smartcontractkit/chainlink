package rebalalgo

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestPingPong(t *testing.T) {
	netA := models.NetworkSelector(1)
	netB := models.NetworkSelector(2)
	netC := models.NetworkSelector(3)

	testCases := []struct {
		name         string
		balances     map[models.NetworkSelector]int64
		lanes        [][2]models.NetworkSelector
		inflight     []UnexecutedTransfer
		expTransfers []models.Transfer
		expErr       bool
	}{
		{
			name:     "simple - two networks one is dry",
			balances: map[models.NetworkSelector]int64{netA: 100, netB: 0},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, {netB, netA}, // A <--> B
			},
			expTransfers: []models.Transfer{
				{From: netA, To: netB, Amount: ubig.NewI(100)},
			},
		},
		{
			name:     "determinism check - two networks both with non negative balance",
			balances: map[models.NetworkSelector]int64{netA: 100, netB: 100},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, {netB, netA}, // A <--> B
			},
			expTransfers: []models.Transfer{
				{From: netA, To: netB, Amount: ubig.NewI(100)},
			},
		},
		{
			name:     "skip non-bidirectional lane",
			balances: map[models.NetworkSelector]int64{netA: 100, netB: 0},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, // A --> B
			},
		},
		{
			name: "three networks - two are dry",
			balances: map[models.NetworkSelector]int64{
				netA: 0,
				netB: 30,
				netC: 0,
			},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, {netB, netA}, // A <--> B
				{netB, netC}, {netC, netB}, // B <--> C
				{netC, netA}, {netA, netC}, // A <--> C
			},
			expTransfers: []models.Transfer{
				{From: netB, To: netA, Amount: ubig.NewI(15)},
				{From: netB, To: netC, Amount: ubig.NewI(15)},
			},
		},
		{
			name: "three networks - one is dry",
			balances: map[models.NetworkSelector]int64{
				netA: 15,
				netB: 0,
				netC: 15,
			},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, {netB, netA}, // A <--> B
				{netB, netC}, {netC, netB}, // B <--> C
				{netC, netA}, {netA, netC}, // A <--> C
			},
			expTransfers: []models.Transfer{
				{From: netA, To: netB, Amount: ubig.NewI(7)},
				{From: netA, To: netC, Amount: ubig.NewI(7)},
				{From: netC, To: netB, Amount: ubig.NewI(15)},
			},
		},
		{
			name:     "three networks - all are dry",
			balances: map[models.NetworkSelector]int64{netA: 0, netB: 0, netC: 0},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, {netB, netA}, // A <--> B
				{netB, netC}, {netC, netB}, // B <--> C
				{netC, netA}, {netA, netC}, // A <--> C
			},
		},
		{
			name: "three networks - one lane non-bidirectional",
			balances: map[models.NetworkSelector]int64{
				netA: 0,
				netB: 30,
				netC: 0,
			},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, {netB, netA}, // A <--> B
				{netB, netC},               // B --> C
				{netC, netA}, {netA, netC}, // C <--> A
			},
			expTransfers: []models.Transfer{
				{From: netB, To: netA, Amount: ubig.NewI(30)},
			},
		},
		{
			name: "three networks - balance less than networks number",
			balances: map[models.NetworkSelector]int64{
				netA: 1,
				netB: 0,
				netC: 0,
			},
			lanes: [][2]models.NetworkSelector{
				{netA, netB}, {netB, netA}, // A <--> B
				{netC, netA}, {netA, netC}, // A <--> C
			},
			expTransfers: []models.Transfer{
				{From: netA, To: netB, Amount: ubig.NewI(1)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pp := NewPingPong()
			g := genGraph(t, tc.balances, tc.lanes)

			transfers, err := pp.ComputeTransfersToBalance(g, tc.inflight)
			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, len(tc.expTransfers), len(transfers))
			for i := range tc.expTransfers {
				expTransfer := tc.expTransfers[i]
				actualTransfer := transfers[i]

				assert.Equal(t, expTransfer.From, actualTransfer.From)
				assert.Equal(t, expTransfer.To, actualTransfer.To)
				assert.Equal(t, expTransfer.Amount, actualTransfer.Amount)
			}
		})
	}
}

// Verifies that ping pong will never reach balance by simulating multiple tests with
// random networks and connections.
func TestPingPong_infinity(t *testing.T) {
	const (
		// number of tests to run
		numSimulations = 50

		// max number of networks per test, result will be in range [2, maxNetworks]
		maxNetworks = 20

		// max number of lanes per test, result will be in range [1, maxBidirectionalLanes]
		maxBidirectionalLanes = 10

		// number of execution rounds per test
		numRoundsPerSimulation = 500
	)

	for i := 0; i < numSimulations; i++ {
		runPingPongInfinitySimulation(t, numRoundsPerSimulation, maxNetworks, maxBidirectionalLanes)
	}
}

func runPingPongInfinitySimulation(t *testing.T, rounds, maxNets, maxLanes int) {
	assert.Greaterf(t, maxNets, 2, "at least 2 networks are required")
	assert.Greaterf(t, maxLanes, 1, "at least one lane is required")
	numNetworks := 2 + rand.Intn(maxNets-1)
	numBidirectionalLanes := 1 + rand.Intn(maxLanes)

	// Generate graph with random networks and balances.
	networks := make([]models.NetworkSelector, numNetworks)
	balances := make(map[models.NetworkSelector]int64, numNetworks)
	for i := range networks {
		networks[i] = models.NetworkSelector(i + 1)
		balances[networks[i]] = int64(rand.Intn(1000))
	}
	lanes := getRandomBidirectionalLanes(networks, numBidirectionalLanes)
	g := genGraph(t, balances, lanes)

	pp := NewPingPong()
	inflightTransfers := make([]UnexecutedTransfer, 0)

	for round := 0; round < rounds; round++ {
		// Cleanup the executed transfers from memory.
		filteredInflights := make([]UnexecutedTransfer, 0, len(inflightTransfers))
		for _, tr := range inflightTransfers {
			if tr.TransferStatus() != models.TransferStatusExecuted {
				filteredInflights = append(filteredInflights, tr)
			}
		}
		inflightTransfers = filteredInflights

		transfersToBalance, err := pp.ComputeTransfersToBalance(g, inflightTransfers)
		assert.NoError(t, err)

		if len(inflightTransfers) == 0 {
			assert.True(t, len(transfersToBalance) > 0, "balance should not be reached")
		}

		pendingTransfers := make([]UnexecutedTransfer, len(transfersToBalance))
		for i, tr := range transfersToBalance {
			pendingTransfers[i] = models.PendingTransfer{
				Transfer: models.Transfer{
					From:   tr.From,
					To:     tr.To,
					Amount: tr.Amount,
				},
			}
		}

		// Find some random inflight transfers and mark them as done by applying them to the graph.
		inflightTransfers = append(inflightTransfers, pendingTransfers...)
		rand.Shuffle(len(inflightTransfers), func(i, j int) {
			inflightTransfers[i], inflightTransfers[j] = inflightTransfers[j], inflightTransfers[i]
		})

		numInflightToApply := 1 + rand.Intn(len(inflightTransfers))
		numApplied := 0
		for idx, inf := range inflightTransfers {
			sourceLiq, err := g.GetLiquidity(inf.FromNetwork())
			assert.NoError(t, err)
			destLiq, err := g.GetLiquidity(inf.ToNetwork())
			assert.NoError(t, err)

			newSourceLiq := big.NewInt(0).Sub(sourceLiq, inf.TransferAmount())
			newDestLiq := big.NewInt(0).Add(destLiq, inf.TransferAmount())

			g.SetLiquidity(inf.FromNetwork(), newSourceLiq)
			g.SetLiquidity(inf.ToNetwork(), newDestLiq)

			// XXX: this is ugly, but we need to update the transfer status to executed.
			// We can't do it directly because the Transfer interface doesn't have a setter.
			pendingTr, ok := inflightTransfers[idx].(models.PendingTransfer)
			require.True(t, ok)

			pendingTr.Status = models.TransferStatusExecuted
			inflightTransfers[idx] = pendingTr

			numApplied++
			if numApplied >= numInflightToApply {
				break
			}
		}
	}
}

// Generates a graph from the provided lanes and balances.
func genGraph(t testing.TB, balances map[models.NetworkSelector]int64, lanes [][2]models.NetworkSelector) graph.Graph {
	g := graph.NewGraph()
	for netSel, balance := range balances {
		g.(graph.GraphTest).AddNetwork(netSel, graph.Data{
			Liquidity:       big.NewInt(balance),
			NetworkSelector: netSel,
		})
	}

	for _, lane := range lanes {
		assert.NoError(t, g.(graph.GraphTest).AddConnection(lane[0], lane[1]))
	}

	return g
}

// Computes and pseudo-randomly returns bidirectional lanes from the provided networks.
func getRandomBidirectionalLanes(networks []models.NetworkSelector, n int) [][2]models.NetworkSelector {
	type Lane struct {
		From models.NetworkSelector
		To   models.NetworkSelector
	}

	allLanes := make([]Lane, 0)
	for _, netFrom := range networks {
		for _, netTo := range networks {
			if netFrom <= netTo {
				continue
			}

			allLanes = append(allLanes, Lane{From: netFrom, To: netTo})
		}
	}

	rand.Shuffle(len(allLanes), func(i, j int) {
		allLanes[i], allLanes[j] = allLanes[j], allLanes[i]
	})

	if n > len(allLanes) {
		n = len(allLanes)
	}

	lanes := make([][2]models.NetworkSelector, 0)
	for _, lane := range allLanes[:n] {
		lanes = append(lanes,
			[2]models.NetworkSelector{lane.From, lane.To},
			[2]models.NetworkSelector{lane.To, lane.From},
		)
	}
	return lanes
}
