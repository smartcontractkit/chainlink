package graph

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestGrpah_NodeReaderGetters(t *testing.T) {
	g := NewGraph()

	data1 := Data{
		Liquidity:               big.NewInt(1),
		TokenAddress:            models.Address(common.HexToAddress("0x11")),
		LiquidityManagerAddress: models.Address(common.HexToAddress("0x12")),
		XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{},
		ConfigDigest: models.ConfigDigest{
			ConfigDigest: [32]byte{1},
		},
		NetworkSelector: models.NetworkSelector(1),
	}
	require.True(t, g.(GraphTest).AddNetwork(models.NetworkSelector(1), data1))

	tests := []struct {
		name string
		net  models.NetworkSelector
		data *Data
	}{
		{
			name: "happy path",
			net:  models.NetworkSelector(1),
			data: &data1,
		},
		{
			name: "not exist",
			net:  models.NetworkSelector(333),
			data: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			liq, err := g.GetLiquidity(tc.net)
			if tc.data == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.data.Liquidity, liq)
			}

			tokenAddr, err := g.GetTokenAddress(tc.net)
			if tc.data == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.data.TokenAddress, tokenAddr)
			}

			liqManagerAddr, err := g.GetLiquidityManagerAddress(tc.net)
			if tc.data == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.data.LiquidityManagerAddress, liqManagerAddr)
			}

			xChainData, err := g.GetXChainLiquidityManagerData(tc.net)
			if tc.data == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.data.XChainLiquidityManagers, xChainData)
			}

			data, err := g.GetData(tc.net)
			if tc.data == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, tc.data.Equals(data))
			}
		})
	}
}

func TestGraph_FindPath(t *testing.T) {
	g := NewGraph()

	data1 := Data{
		Liquidity:               big.NewInt(1),
		TokenAddress:            models.Address(common.HexToAddress("0x11")),
		LiquidityManagerAddress: models.Address(common.HexToAddress("0x12")),
		XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{},
		ConfigDigest: models.ConfigDigest{
			ConfigDigest: [32]byte{1},
		},
		NetworkSelector: models.NetworkSelector(1),
	}
	require.True(t, g.(GraphTest).AddNetwork(models.NetworkSelector(1), data1))

	data2 := Data{
		Liquidity:               big.NewInt(2),
		TokenAddress:            models.Address(common.HexToAddress("0x21")),
		LiquidityManagerAddress: models.Address(common.HexToAddress("0x22")),
		XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{},
		ConfigDigest: models.ConfigDigest{
			ConfigDigest: [32]byte{2},
		},
		NetworkSelector: models.NetworkSelector(2),
	}
	require.True(t, g.(GraphTest).AddNetwork(models.NetworkSelector(2), data2))

	data3 := Data{
		Liquidity:               big.NewInt(3),
		TokenAddress:            models.Address(common.HexToAddress("0x31")),
		LiquidityManagerAddress: models.Address(common.HexToAddress("0x32")),
		XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{},
		ConfigDigest: models.ConfigDigest{
			ConfigDigest: [32]byte{3},
		},
		NetworkSelector: models.NetworkSelector(3),
	}
	require.True(t, g.(GraphTest).AddNetwork(models.NetworkSelector(3), data3))

	require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(1), models.NetworkSelector(2)))
	require.NoError(t, g.(GraphTest).AddConnection(models.NetworkSelector(2), models.NetworkSelector(3)))

	tests := []struct {
		name     string
		from     models.NetworkSelector
		to       models.NetworkSelector
		maxEdges int
		want     []models.NetworkSelector
	}{
		{
			name:     "happy path 2 edges",
			from:     models.NetworkSelector(1),
			to:       models.NetworkSelector(3),
			maxEdges: 2,
			want:     []models.NetworkSelector{models.NetworkSelector(2), models.NetworkSelector(3)},
		},
		{
			name:     "happy path 1 edge",
			from:     models.NetworkSelector(1),
			to:       models.NetworkSelector(2),
			maxEdges: 1,
			want:     []models.NetworkSelector{models.NetworkSelector(2)},
		},
		{
			name:     "not enough edges",
			from:     models.NetworkSelector(1),
			to:       models.NetworkSelector(3),
			maxEdges: 1,
			want:     []models.NetworkSelector{},
		},
		{
			name: "no path",
			from: models.NetworkSelector(2),
			to:   models.NetworkSelector(10),
			want: []models.NetworkSelector{},
		},
		{
			name: "same node",
			from: models.NetworkSelector(1),
			to:   models.NetworkSelector(1),
			want: []models.NetworkSelector{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := g.FindPath(tc.from, tc.to, tc.maxEdges, func(nodes ...Data) bool {
				return true
			})
			require.Equal(t, tc.want, path)
		})
	}
}
