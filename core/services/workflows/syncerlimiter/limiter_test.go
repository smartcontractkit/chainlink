package syncerlimiter_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
)

func TestWorkflowLimits(t *testing.T) {
	t.Parallel()

	config := syncerlimiter.Config{
		Global:   3,
		PerOwner: 1,
	}
	wsl, err := syncerlimiter.NewWorkflowLimits(config)
	require.NoError(t, err)
	allowOwner, allowGlobal := wsl.Allow("user1")
	require.True(t, allowOwner && allowGlobal)
	// Global 1/3, PerOwner 1/1

	allowOwner, allowGlobal = wsl.Allow("user2")
	require.True(t, allowOwner && allowGlobal)
	// Global 2/3, PerOwner 1/1

	allowOwner, allowGlobal = wsl.Allow("user1")
	require.True(t, allowGlobal)
	require.False(t, allowOwner)
	// Global 2/3, PerOwner 1/1 exceeded

	allowOwner, allowGlobal = wsl.Allow("user3")
	require.True(t, allowOwner && allowGlobal)
	// Global 3/3, PerOwner 1/1 (one each user)

	allowOwner, allowGlobal = wsl.Allow("user2")
	require.False(t, allowOwner)
	require.False(t, allowGlobal)
	// Global 3/3, PerOwner 1/1 Global and PerOwner exceeded

	wsl.Decrement("user2")
	// Global 2/3, User2 PerOwner 0/1

	allowOwner, allowGlobal = wsl.Allow("user2")
	require.True(t, allowOwner && allowGlobal)
	// Global 3/3, PerOwner 1/1 (one each user)

	wsl.Decrement("non-existent-user")
	allowOwner, allowGlobal = wsl.Allow("non-existent-user")
	require.True(t, allowOwner)
	require.False(t, allowGlobal)
	// Global 3/3, PerOwner 0/1 Global exceeded
}
