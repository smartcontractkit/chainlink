package discoverer

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func Test_EvmDiscoverer_DiscoverBalances(t *testing.T) {
	tests := []struct {
		name            string
		initialGraph    map[models.NetworkSelector]*big.Int
		liquidityGetter func(ctx context.Context, network models.NetworkSelector, lmAddress common.Address) (*big.Int, error)
		wantGraph       map[models.NetworkSelector]*big.Int
		wantErr         bool
	}{
		{
			name:         "empty",
			initialGraph: map[models.NetworkSelector]*big.Int{},
			liquidityGetter: func(ctx context.Context, network models.NetworkSelector, lmAddress common.Address) (*big.Int, error) {
				return big.NewInt(100), nil
			},
			wantGraph: map[models.NetworkSelector]*big.Int{},
		},
		{
			name: "happy path",
			initialGraph: map[models.NetworkSelector]*big.Int{
				1: big.NewInt(100),
				2: big.NewInt(100),
				3: big.NewInt(100),
			},
			liquidityGetter: func(ctx context.Context, network models.NetworkSelector, lmAddress common.Address) (*big.Int, error) {
				liq := big.NewInt(0).Mul(big.NewInt(100), big.NewInt(int64(network)))
				return liq, nil
			},
			wantGraph: map[models.NetworkSelector]*big.Int{
				1: big.NewInt(100),
				2: big.NewInt(100 * 2),
				3: big.NewInt(100 * 3),
			},
		},
		{
			name: "error",
			initialGraph: map[models.NetworkSelector]*big.Int{
				1: big.NewInt(100),
				2: big.NewInt(100),
				3: big.NewInt(100),
			},
			liquidityGetter: func(ctx context.Context, network models.NetworkSelector, lmAddress common.Address) (*big.Int, error) {
				if network%2 == 0 {
					return nil, fmt.Errorf("dummy test error")
				}
				liq := big.NewInt(0).Mul(big.NewInt(100), big.NewInt(int64(network)))
				return liq, nil
			},
			wantGraph: map[models.NetworkSelector]*big.Int{
				1: big.NewInt(100),
				2: big.NewInt(100),     // got error
				3: big.NewInt(100 * 3), // 3 is the only one that should be updated
			},
			wantErr: true,
		},
		{
			name: "10 networks",
			initialGraph: map[models.NetworkSelector]*big.Int{
				1:  big.NewInt(100),
				2:  big.NewInt(100),
				3:  big.NewInt(100),
				4:  big.NewInt(100),
				5:  big.NewInt(100),
				6:  big.NewInt(100),
				7:  big.NewInt(100),
				8:  big.NewInt(100),
				9:  big.NewInt(100),
				10: big.NewInt(100),
			},
			liquidityGetter: func(ctx context.Context, network models.NetworkSelector, lmAddress common.Address) (*big.Int, error) {
				liq := big.NewInt(0).Mul(big.NewInt(100), big.NewInt(int64(network)))
				return liq, nil
			},
			wantGraph: map[models.NetworkSelector]*big.Int{
				1:  big.NewInt(100),
				2:  big.NewInt(100 * 2),
				3:  big.NewInt(100 * 3),
				4:  big.NewInt(100 * 4),
				5:  big.NewInt(100 * 5),
				6:  big.NewInt(100 * 6),
				7:  big.NewInt(100 * 7),
				8:  big.NewInt(100 * 8),
				9:  big.NewInt(100 * 9),
				10: big.NewInt(100 * 10),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := graph.NewGraph()
			for network, liq := range tc.initialGraph {
				g.(graph.GraphTest).AddNetwork(network, graph.Data{Liquidity: liq})
			}
			d := &evmDiscoverer{
				lggr:            logger.TestLogger(t),
				liquidityGetter: tc.liquidityGetter,
			}
			err := d.DiscoverBalances(testutils.Context(t), g)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			for network, expectedLiq := range tc.wantGraph {
				liq, err := g.GetLiquidity(network)
				require.NoError(t, err)
				require.Equalf(t, expectedLiq, liq, "wrong liquidity for network %d", network)
			}
		})
	}
}
