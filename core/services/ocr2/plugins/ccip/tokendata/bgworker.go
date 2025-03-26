package tokendata

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type msgResult struct {
	TokenAmountIndex int
	Err              error
	Data             []byte
}

type Worker interface {
	job.ServiceCtx
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
	resultsCache     *cache.Cache
	timeoutDur       time.Duration

	services.StateMachine
	wg       sync.WaitGroup
	stopChan services.StopChan
}

func NewBackgroundWorker(
	tokenDataReaders map[cciptypes.Address]Reader,
	numWorkers int,
	timeoutDur time.Duration,
	expirationDur time.Duration,
) *BackgroundWorker {
	if expirationDur == 0 {
		expirationDur = 24 * time.Hour
	}

	return &BackgroundWorker{
		tokenDataReaders: tokenDataReaders,
		numWorkers:       numWorkers,
		jobsChan:         make(chan cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, numWorkers*200),
		resultsCache:     cache.New(expirationDur, expirationDur/2),
		timeoutDur:       timeoutDur,
		stopChan:         make(services.StopChan),
	}
}

func (w *BackgroundWorker) Start(context.Context) error {
	return w.StateMachine.StartOnce("Token BackgroundWorker", func() error {
		for i := 0; i < w.numWorkers; i++ {
			w.wg.Add(1)
			w.run()
		}
		return nil
	})
}

func (w *BackgroundWorker) Close() error {
	return w.StateMachine.StopOnce("Token BackgroundWorker", func() error {
		close(w.stopChan)
		w.wg.Wait()
		return nil
	})
}

func (w *BackgroundWorker) AddJobsFromMsgs(ctx context.Context, msgs []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) {
	w.wg.Add(1)
	go func(ctx context.Context) {
		defer w.wg.Done()
		ctx, cancel := w.stopChan.Ctx(ctx)
		defer cancel()

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
	}(ctx)
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
			return nil, errors.New("token data index inconsistency")
		}
		tokenDatas[r.TokenAmountIndex] = r.Data
	}

	return tokenDatas, nil
}

func (w *BackgroundWorker) run() {
	go func() {
		defer w.wg.Done()
		ctx, cancel := w.stopChan.NewCtx()
		defer cancel()

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

func (w *BackgroundWorker) workOnMsg(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) {
	results := make([]msgResult, 0, len(msg.TokenAmounts))

	cachedTokenData := make(map[int]msgResult) // tokenAmount index -> token data
	if cachedData, exists := w.getFromCache(msg.SequenceNumber); exists {
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

	w.resultsCache.Set(strconv.FormatUint(msg.SequenceNumber, 10), results, cache.DefaultExpiration)
}

func (w *BackgroundWorker) getMsgTokenData(ctx context.Context, seqNum uint64) ([]msgResult, error) {
	if msgTokenData, exists := w.getFromCache(seqNum); exists {
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
			if msgTokenData, exists := w.getFromCache(seqNum); exists {
				return msgTokenData, nil
			}
		}
	}
}

func (w *BackgroundWorker) getFromCache(seqNum uint64) ([]msgResult, bool) {
	rawResult, found := w.resultsCache.Get(strconv.FormatUint(seqNum, 10))
	if !found {
		return nil, false
	}
	return rawResult.([]msgResult), true
}
