package cache

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	cachedValue = "cached_value"
)

func TestGet_InitDataForTheFirstTime(t *testing.T) {
	lp := lpMocks.NewLogPoller(t)
	lp.On("LatestBlock", mock.Anything).Maybe().Return(logpoller.LogPollerBlock{BlockNumber: 100}, nil)

	contract := newCachedContract(lp, "", []string{"value1"}, 0)

	value, err := contract.Get(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, "value1", value)
}

func TestGet_ReturnDataFromCacheIfNoNewEvents(t *testing.T) {
	latestBlock := int64(100)
	lp := lpMocks.NewLogPoller(t)
	mockLogPollerQuery(lp, latestBlock)

	contract := newCachedContract(lp, cachedValue, []string{"value1"}, latestBlock)

	value, err := contract.Get(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, cachedValue, value)
}

func TestGet_CallOriginForNewEvents(t *testing.T) {
	latestBlock := int64(100)
	lp := lpMocks.NewLogPoller(t)
	m := mockLogPollerQuery(lp, latestBlock+1)

	contract := newCachedContract(lp, cachedValue, []string{"value1", "value2", "value3"}, latestBlock)

	// First call
	value, err := contract.Get(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, "value1", value)

	currentBlock := contract.lastChangeBlock
	require.Equal(t, latestBlock+1, currentBlock)

	m.Unset()
	mockLogPollerQuery(lp, latestBlock+1)

	// Second call doesn't change anything
	value, err = contract.Get(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, "value1", value)
	require.Equal(t, int64(101), contract.lastChangeBlock)
}

func TestGet_CacheProgressing(t *testing.T) {
	firstBlock := int64(100)
	secondBlock := int64(105)
	thirdBlock := int64(110)

	lp := lpMocks.NewLogPoller(t)
	m := mockLogPollerQuery(lp, secondBlock)

	contract := newCachedContract(lp, cachedValue, []string{"value1", "value2", "value3"}, firstBlock)

	// First call
	value, err := contract.Get(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, "value1", value)
	require.Equal(t, secondBlock, contract.lastChangeBlock)

	m.Unset()
	mockLogPollerQuery(lp, thirdBlock)

	// Second call
	value, err = contract.Get(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, "value2", value)
	require.Equal(t, thirdBlock, contract.lastChangeBlock)
}

func TestGet_ConcurrentAccess(t *testing.T) {
	mockedPoller := lpMocks.NewLogPoller(t)
	progressingPoller := ProgressingLogPoller{
		LogPoller:   mockedPoller,
		latestBlock: 1,
	}

	iterations := 100
	originValues := make([]string, iterations)
	for i := 0; i < iterations; i++ {
		originValues[i] = "value_" + strconv.Itoa(i)
	}
	contract := newCachedContract(&progressingPoller, "empty", originValues, 1)

	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			_, _ = contract.Get(testutils.Context(t))
		}()
	}
	wg.Wait()

	// 1 init block + 100 iterations
	require.Equal(t, int64(101), contract.lastChangeBlock)
}

func newCachedContract(lp logpoller.LogPoller, cacheValue string, originValue []string, lastChangeBlock int64) *CachedChain[string] {
	return &CachedChain[string]{
		observedEvents:          []common.Hash{{}},
		logPoller:               lp,
		address:                 []common.Address{{}},
		optimisticConfirmations: 0,

		lock:            &sync.RWMutex{},
		value:           cacheValue,
		lastChangeBlock: lastChangeBlock,
		origin:          &FakeContractOrigin{values: originValue},
	}
}

func mockLogPollerQuery(lp *lpMocks.LogPoller, latestBlock int64) *mock.Call {
	return lp.On("LatestBlockByEventSigsAddrsWithConfs", mock.Anything, []common.Hash{{}}, []common.Address{{}}, logpoller.Confirmations(0), mock.Anything).
		Maybe().Return(latestBlock, nil)
}

type ProgressingLogPoller struct {
	*lpMocks.LogPoller
	latestBlock int64
	lock        sync.Mutex
}

func (lp *ProgressingLogPoller) LatestBlockByEventSigsAddrsWithConfs(int64, []common.Hash, []common.Address, logpoller.Confirmations, ...pg.QOpt) (int64, error) {
	lp.lock.Lock()
	defer lp.lock.Unlock()
	lp.latestBlock++
	return lp.latestBlock, nil
}

type FakeContractOrigin struct {
	values  []string
	counter int
	lock    sync.Mutex
}

func (f *FakeContractOrigin) CallOrigin(context.Context) (string, error) {
	f.lock.Lock()
	defer func() {
		f.counter++
		f.lock.Unlock()
	}()
	return f.values[f.counter], nil
}

func (f *FakeContractOrigin) Copy(value string) string {
	return value
}
