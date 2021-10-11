package headtracker_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/core/services/headtracker/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_HeadSaver_addHeads(t *testing.T) {
	uncleHash := utils.NewHash()
	cfg := new(htmocks.Config)
	cfg.Test(t)
	cfg.On("EvmFinalityDepth").Return(uint32(1))
	hs := headtracker.NewHeadSaver(logger.Default, nil, cfg)

	var heads []*eth.Head
	var parentHash common.Hash
	for i := 0; i < 5; i++ {
		hash := utils.NewHash()
		h := eth.NewHead(big.NewInt(int64(i)), hash, parentHash, uint64(time.Now().Unix()), utils.NewBigI(0))
		heads = append(heads, &h)
		if i == 2 {
			// uncled block
			h := eth.NewHead(big.NewInt(int64(i)), uncleHash, parentHash, uint64(time.Now().Unix()), utils.NewBigI(0))
			heads = append(heads, &h)
		}
		parentHash = hash
	}
	headtracker.AddHeads(hs, heads, 6)
	// Add duplicates (should be ignored)
	headtracker.AddHeads(hs, heads[2:5], 6)

	ch := hs.LatestChain()
	assert.Equal(t, 6, len(headtracker.Heads(hs)))
	require.NotNil(t, ch)
	require.Equal(t, 5, int(ch.ChainLength()))

	ch = hs.Chain(uncleHash)
	assert.Equal(t, 6, len(headtracker.Heads(hs)))
	require.NotNil(t, ch)
	require.Equal(t, 3, int(ch.ChainLength()))

	// Adding beyond the limit truncates
	headtracker.AddHeads(hs, heads, 2)
	assert.Equal(t, 2, len(headtracker.Heads(hs)))
	ch = hs.LatestChain()
	require.NotNil(t, ch)
	require.Equal(t, 2, int(ch.ChainLength()))
}
