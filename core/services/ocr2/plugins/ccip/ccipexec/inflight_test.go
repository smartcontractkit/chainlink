package ccipexec

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInflightReportsContainer_add(t *testing.T) {
	lggr := logger.TestLogger(t)
	container := newInflightExecReportsContainer(time.Second)

	err := container.add(lggr, []cciptypes.EVM2EVMMessage{
		{SequenceNumber: 1}, {SequenceNumber: 2}, {SequenceNumber: 3},
	})
	require.NoError(t, err)
	err = container.add(lggr, []cciptypes.EVM2EVMMessage{
		{SequenceNumber: 1},
	})
	require.Error(t, err)
	require.Equal(t, "report is already in flight", err.Error())
	require.Equal(t, 1, len(container.getAll()))
}

func TestInflightReportsContainer_expire(t *testing.T) {
	lggr := logger.TestLogger(t)
	container := newInflightExecReportsContainer(time.Second)

	err := container.add(lggr, []cciptypes.EVM2EVMMessage{
		{SequenceNumber: 1}, {SequenceNumber: 2}, {SequenceNumber: 3},
	})
	require.NoError(t, err)
	container.reports[0].createdAt = time.Now().Add(-time.Second * 5)
	require.Equal(t, 1, len(container.getAll()))

	container.expire(lggr)
	require.Equal(t, 0, len(container.getAll()))
}
