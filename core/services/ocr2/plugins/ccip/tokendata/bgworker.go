package tokendata

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
)

type msgResult struct {
	TokenAmountIndex int
	Err              error
	Data             []byte
}

type Worker interface {
	// AddJobsFromMsgs will include the provided msgs for background processing.
	AddJobsFromMsgs(ctx context.Context, msgs []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta)

	// GetMsgTokenData returns the token data for the provided msg. If data are not ready it keeps waiting
	// until they get ready. Important: Make sure to pass a proper context with timeout.
	GetMsgTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([][]byte, error)

	GetReaders() map[cciptypes.Address]Reader
}

type BackgroundWorker struct {
	tokenDataReaders map[cciptypes.Address]Reader
	numWorkers       int
	jobsChan         chan cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
	resultsCache     *resultsCache
	timeoutDur       time.Duration
}

func (w *BackgroundWorker) AddJobsFromMsgs(ctx context.Context, msgs []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) {
	go func() {
		for _, msg := range msgs {
			select {
			case <-ctx.Done():
				return
			default:
				if len(msg.TokenAmounts) > 0 {
					w.jobsChan <- msg
				}
			}
		}
	}()
}

func (w *BackgroundWorker) GetReaders() map[cciptypes.Address]Reader {
	return w.tokenDataReaders
}

func (w *BackgroundWorker) GetMsgTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([][]byte, error) {
	res, err := w.getMsgTokenData(ctx, msg.SequenceNumber)
	if err != nil {
		return nil, err
	}

	tokenDatas := make([][]byte, len(msg.TokenAmounts))
	for _, r := range res {
		if r.Err != nil {
			return nil, r.Err
		}
		if r.TokenAmountIndex < 0 || r.TokenAmountIndex >= len(msg.TokenAmounts) {
			return nil, fmt.Errorf("token data index inconsistency")
		}
		tokenDatas[r.TokenAmountIndex] = r.Data
	}

	return tokenDatas, nil
}

func NewBackgroundWorker(
	ctx context.Context,
	tokenDataReaders map[cciptypes.Address]Reader,
	numWorkers int,
	timeoutDur time.Duration,
	expirationDur time.Duration,
) *BackgroundWorker {
	if expirationDur == 0 {
		expirationDur = 24 * time.Hour
	}

	w := &BackgroundWorker{
		tokenDataReaders: tokenDataReaders,
		numWorkers:       numWorkers,
		jobsChan:         make(chan cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, numWorkers*100),
		resultsCache:     newResultsCache(ctx, expirationDur, expirationDur/2),
		timeoutDur:       timeoutDur,
	}

	w.spawnWorkers(ctx)
	return w
}

func (w *BackgroundWorker) spawnWorkers(ctx context.Context) {
	for i := 0; i < w.numWorkers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-w.jobsChan:
					w.workOnMsg(ctx, msg)
				}
			}
		}()
	}
}

func (w *BackgroundWorker) workOnMsg(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) {
	results := make([]msgResult, 0, len(msg.TokenAmounts))

	cachedTokenData := make(map[int]msgResult) // tokenAmount index -> token data
	if cachedData, exists := w.resultsCache.get(msg.SequenceNumber); exists {
		for _, r := range cachedData {
			cachedTokenData[r.TokenAmountIndex] = r
		}
	}

	for i, token := range msg.TokenAmounts {
		offchainTokenDataProvider, exists := w.tokenDataReaders[token.Token]
		if !exists {
			// No token data required
			continue
		}

		// if the result exists in the cache and there wasn't any error keep the existing result
		if cachedResult, exists := cachedTokenData[i]; exists && cachedResult.Err == nil {
			results = append(results, cachedResult)
			continue
		}

		// if there was any error or if the data do not exist in the cache make a call to the provider
		timeoutCtx, cf := context.WithTimeout(ctx, w.timeoutDur)
		tknData, err := offchainTokenDataProvider.ReadTokenData(timeoutCtx, msg, i)
		cf()
		results = append(results, msgResult{
			TokenAmountIndex: i,
			Err:              err,
			Data:             tknData,
		})
	}

	w.resultsCache.add(msg.SequenceNumber, results)
}

func (w *BackgroundWorker) getMsgTokenData(ctx context.Context, seqNum uint64) ([]msgResult, error) {
	if msgTokenData, exists := w.resultsCache.get(seqNum); exists {
		return msgTokenData, nil
	}

	ctx, cf := context.WithTimeout(ctx, w.timeoutDur)
	defer cf()

	// wait until the results are ready or until context timeout is reached
	tick := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		case <-tick.C:
			if msgTokenData, exists := w.resultsCache.get(seqNum); exists {
				return msgTokenData, nil
			}
		}
	}
}

type resultsCache struct {
	expirationDuration time.Duration
	expiresAt          map[uint64]time.Time
	results            map[uint64][]msgResult
	resultsMu          *sync.RWMutex
}

func newResultsCache(ctx context.Context, expirationDuration, cleanupInterval time.Duration) *resultsCache {
	c := &resultsCache{
		expirationDuration: expirationDuration,
		expiresAt:          make(map[uint64]time.Time),
		results:            make(map[uint64][]msgResult),
		resultsMu:          &sync.RWMutex{},
	}

	ticker := time.NewTicker(cleanupInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.cleanExpiredItems()
			}
		}
	}()

	return c
}

func (c *resultsCache) add(msgSeqNum uint64, results []msgResult) {
	c.resultsMu.Lock()
	defer c.resultsMu.Unlock()
	c.results[msgSeqNum] = results
	c.expiresAt[msgSeqNum] = time.Now().Add(c.expirationDuration)
}

func (c *resultsCache) get(msgSeqNum uint64) ([]msgResult, bool) {
	c.resultsMu.RLock()
	defer c.resultsMu.RUnlock()
	v, exists := c.results[msgSeqNum]
	return v, exists
}

func (c *resultsCache) cleanExpiredItems() {
	c.resultsMu.RLock()
	expiredKeys := make([]uint64, 0, len(c.expiresAt))
	for seqNum, expiresAt := range c.expiresAt {
		if expiresAt.Before(time.Now()) {
			expiredKeys = append(expiredKeys, seqNum)
		}
	}
	c.resultsMu.RUnlock()

	if len(expiredKeys) == 0 {
		return
	}

	c.resultsMu.Lock()
	for _, seqNum := range expiredKeys {
		delete(c.results, seqNum)
		delete(c.expiresAt, seqNum)
	}
	c.resultsMu.Unlock()
}
