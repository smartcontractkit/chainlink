package mercury

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	htmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	mercurymocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestMercurySetCurrentBlock(t *testing.T) {
	lggr := logger.TestLogger(t)
	ds := datasource{
		lggr: lggr,
	}

	h := evmtypes.Head{
		Number:           testutils.NewRandomPositiveInt64(),
		Hash:             utils.NewHash(),
		ParentHash:       utils.NewHash(),
		Timestamp:        time.Now(),
		BaseFeePerGas:    assets.NewWeiI(testutils.NewRandomPositiveInt64()),
		ReceiptsRoot:     utils.NewHash(),
		TransactionsRoot: utils.NewHash(),
		StateRoot:        utils.NewHash(),
	}

	t.Run("returns head from headtracker if present", func(t *testing.T) {
		headTracker := htmocks.NewHeadTracker(t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("HeadTracker").Return(headTracker)
		headTracker.On("LatestChain").Return(&h, nil)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercury.Observation{}
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, h.Number, obs.CurrentBlockNum.Val)
		assert.Equal(t, h.Hash.Bytes(), obs.CurrentBlockHash.Val)

		chainHeadTracker.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head and eth call succeeds", func(t *testing.T) {
		ethClient := evmclimocks.NewClient(t)
		headTracker := htmocks.NewHeadTracker(t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("Client").Return(ethClient)
		chainHeadTracker.On("HeadTracker").Return(headTracker)
		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return(nil)
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&h, nil)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercury.Observation{}
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, h.Number, obs.CurrentBlockNum.Val)
		assert.Equal(t, h.Hash.Bytes(), obs.CurrentBlockHash.Val)

		chainHeadTracker.AssertExpectations(t)
		ethClient.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head and eth call fails", func(t *testing.T) {
		ethClient := evmclimocks.NewClient(t)
		headTracker := htmocks.NewHeadTracker(t)
		chainHeadTracker := mercurymocks.NewChainHeadTracker(t)

		chainHeadTracker.On("Client").Return(ethClient)
		chainHeadTracker.On("HeadTracker").Return(headTracker)
		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return(nil)
		err := errors.New("foo")
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, err)

		ds.chainHeadTracker = chainHeadTracker

		obs := relaymercury.Observation{}
		ds.setCurrentBlock(context.Background(), &obs)

		assert.Equal(t, err, obs.CurrentBlockNum.Err)
		assert.Equal(t, err, obs.CurrentBlockHash.Err)

		chainHeadTracker.AssertExpectations(t)
		ethClient.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})
}
