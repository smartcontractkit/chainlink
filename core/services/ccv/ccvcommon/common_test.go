package ccvcommon

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	evmmocks "github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	ethereumMainnetChainID = 1
	polygonMainnetChainID  = 137
)

// plainChainService is a ChainService that is not a legacyevm.Chain, as in LOOPP mode.
type plainChainService struct {
	commontypes.ChainService
	name string
}

func (p *plainChainService) Name() string { return p.name }

func (p *plainChainService) GetChainInfo(context.Context) (commontypes.ChainInfo, error) {
	return commontypes.ChainInfo{}, nil
}

// newMockChain makes a legacyevm.Chain mock for the given EVM chain ID.
func newMockChain(t *testing.T, chainID int64) *evmmocks.Chain {
	c := evmmocks.NewChain(t)
	c.EXPECT().GetChainInfo(mockAnyCtx()).Return(commontypes.ChainInfo{}, nil).Maybe()
	c.EXPECT().ID().Return(big.NewInt(chainID)).Maybe()
	c.EXPECT().Name().Return("chain-" + big.NewInt(chainID).String()).Maybe()
	return c
}

func mockAnyCtx() context.Context { return context.Background() }

func selectorFor(t *testing.T, chainID uint64) protocol.ChainSelector {
	t.Helper()
	c, ok := chainselectors.ChainByEvmChainID(chainID)
	require.True(t, ok)
	return protocol.ChainSelector(c.Selector)
}

func TestGetLegacyChains(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)

	ethSel := selectorFor(t, ethereumMainnetChainID)
	polySel := selectorFor(t, polygonMainnetChainID)

	t.Run("skips a chain whose chain info fails and keeps the others", func(t *testing.T) {
		bad := evmmocks.NewChain(t)
		bad.EXPECT().GetChainInfo(mockAnyCtx()).Return(commontypes.ChainInfo{}, errors.New("boom"))
		bad.EXPECT().Name().Return("bad").Maybe()

		good := newMockChain(t, ethereumMainnetChainID)

		chains, missing, err := GetLegacyChains(ctx, lggr,
			[]commontypes.ChainService{bad, good}, []protocol.ChainSelector{ethSel, polySel})

		require.NoError(t, err)
		require.Len(t, chains, 1)
		require.Contains(t, chains, ethSel)
		require.Equal(t, []protocol.ChainSelector{polySel}, missing)
	})

	t.Run("skips a chain service that is not a legacyevm.Chain", func(t *testing.T) {
		chains, missing, err := GetLegacyChains(ctx, lggr,
			[]commontypes.ChainService{
				&plainChainService{name: "loopp"},
				newMockChain(t, ethereumMainnetChainID),
			},
			[]protocol.ChainSelector{ethSel})

		require.NoError(t, err)
		require.Len(t, chains, 1)
		require.Empty(t, missing)
	})

	t.Run("reports a config chain that has no chain service", func(t *testing.T) {
		chains, missing, err := GetLegacyChains(ctx, lggr,
			[]commontypes.ChainService{newMockChain(t, ethereumMainnetChainID)},
			[]protocol.ChainSelector{ethSel, polySel})

		require.NoError(t, err)
		require.Len(t, chains, 1)
		require.Equal(t, []protocol.ChainSelector{polySel}, missing)
	})

	t.Run("errors when no config chain resolves", func(t *testing.T) {
		_, _, err := GetLegacyChains(ctx, lggr,
			[]commontypes.ChainService{newMockChain(t, ethereumMainnetChainID)},
			[]protocol.ChainSelector{polySel})

		require.ErrorContains(t, err, "no chain object is available")
	})

	t.Run("no config chains is not an error", func(t *testing.T) {
		chains, missing, err := GetLegacyChains(ctx, lggr, nil, nil)

		require.NoError(t, err)
		require.Empty(t, chains)
		require.Empty(t, missing)
	})
}

func TestMissingChains(t *testing.T) {
	chains := map[protocol.ChainSelector]struct{}{1: {}}

	require.Equal(t, []protocol.ChainSelector{2, 3},
		MissingChains([]protocol.ChainSelector{1, 3, 2, 3}, chains))
	require.Empty(t, MissingChains([]protocol.ChainSelector{1}, chains))
}
