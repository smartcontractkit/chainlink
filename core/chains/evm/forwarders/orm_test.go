package forwarders

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
)

// Tests the atomicity of cleanup function passed to DeleteForwarder, during DELETE operation
func Test_DeleteForwarder(t *testing.T) {
	t.Parallel()
	orm := NewORM(testutils.NewSqlxDB(t))
	addr := testutils.NewAddress()
	chainID := testutils.FixtureChainID
	ctx := testutils.Context(t)

	fwd, err := orm.CreateForwarder(ctx, addr, *big.New(chainID))
	require.NoError(t, err)
	assert.Equal(t, addr, fwd.Address)

	ErrCleaningUp := errors.New("error during cleanup")

	cleanupCalled := 0

	// Cleanup should fail the first time, causing delete to abort.  When cleanup succeeds the second time,
	//  delete should succeed.  Should fail the 3rd and 4th time since the forwarder has already been deleted.
	//  cleanup should only be called the first two times (when DELETE can succeed).
	rets := []error{ErrCleaningUp, nil, nil, ErrCleaningUp}
	expected := []error{ErrCleaningUp, nil, sql.ErrNoRows, sql.ErrNoRows}

	testCleanupFn := func(q sqlutil.DataSource, evmChainID int64, addr common.Address) error {
		require.Less(t, cleanupCalled, len(rets))
		cleanupCalled++
		return rets[cleanupCalled-1]
	}

	for _, expect := range expected {
		err = orm.DeleteForwarder(ctx, fwd.ID, testCleanupFn)
		assert.ErrorIs(t, err, expect)
	}
	assert.Equal(t, 2, cleanupCalled)
}
