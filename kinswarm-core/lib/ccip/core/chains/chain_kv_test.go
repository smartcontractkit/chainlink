package chains_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
)

func Test_ChainKV(t *testing.T) {
	var (
		testChainID = "id"
		testChain   = &testChainService{name: "test chain"}
	)
	// test empty case
	empty := make(map[string]*testChainService)
	kv := chains.NewChainsKV[*testChainService](empty)
	c, err := kv.Get(testChainID)
	assert.Nil(t, c)
	assert.ErrorIs(t, err, chains.ErrNoSuchChainID)

	assert.Equal(t, kv.Len(), 0)
	assert.Len(t, kv.Slice(), 0)

	cs, err := kv.List()
	assert.NoError(t, err)
	assert.Len(t, cs, 0)

	// test with one chain
	onechain := map[string]*testChainService{testChainID: testChain}
	kv = chains.NewChainsKV[*testChainService](onechain)
	c, err = kv.Get(testChainID)
	assert.Equal(t, c, testChain)
	assert.NoError(t, err)

	assert.Equal(t, kv.Len(), 1)
	assert.Len(t, kv.Slice(), 1)

	cs, err = kv.List()
	assert.NoError(t, err)
	assert.Len(t, cs, 1)

	//List explicit chain
	cs, err = kv.List(testChainID)
	assert.NoError(t, err)
	assert.Len(t, cs, 1)
	assert.Equal(t, testChain, cs[0])

	//List no such id
	cs, err = kv.List("no such id")
	assert.Error(t, err)
	assert.Len(t, cs, 0)
}

type testChainService struct {
	name string
}

// Start the service. Must quit immediately if the context is cancelled.
// The given context applies to Start function only and must not be retained.
func (s *testChainService) Start(_ context.Context) error {
	return nil
}

// Close stops the Service.
// Invariants: Usually after this call the Service cannot be started
// again, you need to build a new Service to do so.
func (s *testChainService) Close() error {
	return nil
}

// Name returns the fully qualified name of the service
func (s *testChainService) Name() string {
	return s.name
}

// Ready should return nil if ready, or an error message otherwise.
func (s *testChainService) Ready() error {
	return nil
}

// HealthReport returns a full health report of the callee including it's dependencies.
// key is the dep name, value is nil if healthy, or error message otherwise.
func (s *testChainService) HealthReport() map[string]error {
	return map[string]error{}
}

// Implement [types.LatestHead] interface
func (s *testChainService) LatestHead(_ context.Context) (head types.Head, err error) {
	return
}

// Implement [types.ChainService] interface
func (s *testChainService) GetChainStatus(ctx context.Context) (stat types.ChainStatus, err error) {
	return
}
func (s *testChainService) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	return
}
func (s *testChainService) Transact(ctx context.Context, from string, to string, amount *big.Int, balanceCheck bool) error {
	return nil
}
