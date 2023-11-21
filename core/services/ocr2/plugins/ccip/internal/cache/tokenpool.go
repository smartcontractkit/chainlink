package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

func NewTokenPools(
	lggr logger.Logger,
	lp logpoller.LogPoller,
	offRamp ccipdata.OffRampReader,
	numWorkers int,
) *CachedChain[map[common.Address]common.Address] {
	return &CachedChain[map[common.Address]common.Address]{
		observedEvents:  offRamp.TokenEvents(),
		logPoller:       lp,
		address:         []common.Address{offRamp.Address()},
		lock:            &sync.RWMutex{},
		value:           make(map[common.Address]common.Address),
		lastChangeBlock: 0,
		origin:          newTokenPoolsOrigin(lggr, offRamp, numWorkers),
	}
}

func newTokenPoolsOrigin(
	lggr logger.Logger,
	offRamp ccipdata.OffRampReader,
	numWorkers int) *tokenPools {
	return &tokenPools{
		lggr:       lggr,
		offRamp:    offRamp,
		numWorkers: numWorkers,
	}
}

type tokenPools struct {
	lggr       logger.Logger
	offRamp    ccipdata.OffRampReader
	numWorkers int
}

func (t *tokenPools) Copy(value map[common.Address]common.Address) map[common.Address]common.Address {
	return copyMap(value)
}

func (t *tokenPools) CallOrigin(ctx context.Context) (map[common.Address]common.Address, error) {
	destTokens, err := t.offRamp.GetDestinationTokens(ctx)
	if err != nil {
		return nil, err
	}

	eg := new(errgroup.Group)
	eg.SetLimit(t.numWorkers)
	var mu sync.Mutex

	mapping := make(map[common.Address]common.Address, len(destTokens))
	for _, token := range destTokens {
		token := token
		eg.Go(func() error {
			poolAddress, err := t.offRamp.GetPoolByDestToken(ctx, token)
			if err != nil {
				return fmt.Errorf("get token pool for token '%s': %w", token, err)
			}

			mu.Lock()
			mapping[token] = poolAddress
			mu.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return mapping, nil
}
