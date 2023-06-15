package txmgr_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func Test_SendEveryStrategy(t *testing.T) {
	t.Parallel()

	s := txmgrcommon.SendEveryStrategy{}

	assert.Equal(t, uuid.NullUUID{}, s.Subject())

	n, err := s.PruneQueue(nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)
}

func Test_DropOldestStrategy_Subject(t *testing.T) {
	t.Parallel()
	cfg := configtest.NewGeneralConfig(t, nil)

	subject := uuid.New()
	s := txmgrcommon.NewDropOldestStrategy(subject, 1, cfg.Database().DefaultQueryTimeout())

	assert.True(t, s.Subject().Valid)
	assert.Equal(t, subject, s.Subject().UUID)
}

func Test_DropOldestStrategy_PruneQueue(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	subject := uuid.New()
	queueSize := uint32(2)
	queryTimeout := cfg.Database().DefaultQueryTimeout()
	mockTxStore := mocks.NewEvmTxStore(t)

	t.Run("calls PrineUnstartedTxQueue for the given subject and queueSize, ignoring fromAddress", func(t *testing.T) {
		strategy1 := txmgrcommon.NewDropOldestStrategy(subject, queueSize, queryTimeout)
		mockTxStore.On("PruneUnstartedTxQueue", queueSize, subject, mock.Anything, mock.Anything).Once().Return(int64(2), nil)
		n, err := strategy1.PruneQueue(mockTxStore, pg.WithQueryer(db))
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)
	})
}
