package syncerlimiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	user1String = "119BFD3D78fbb740c614432975CBE829E26C490e"
	user2String = "219BFD3D78fbb740c614432975CBE829E26C490e"
	user3String = "319BFD3D78fbb740c614432975CBE829E26C490e"
	user4String = "419BFD3D78fbb740c614432975CBE829E26C490e"
	user5String = "519BFD3D78fbb740c614432975CBE829E26C490e"
)

func TestWorkflowLimits(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	config := Config{
		Global:   3,
		PerOwner: 1,
		PerOwnerOverrides: map[string]int32{
			"0x" + user5String: 2,
		},
	}
	wsl, err := NewWorkflowLimits(lggr, config, limits.Factory{Logger: lggr.Named("Limits")})
	require.NoError(t, err)

	ctx1 := contexts.WithCRE(t.Context(), contexts.CRE{Owner: user1String})
	require.NoError(t, wsl.Use(ctx1, 1))
	// Global 1/3, PerOwner 1/1

	ctx2 := contexts.WithCRE(t.Context(), contexts.CRE{Owner: user2String})
	require.NoError(t, wsl.Use(ctx2, 1))
	// Global 2/3, PerOwner 1/1

	err = wsl.Use(ctx1, 1)
	require.Error(t, err)
	var errLimited limits.ErrorResourceLimited[int]
	if assert.ErrorAs(t, err, &errLimited) {
		require.Equal(t, settings.ScopeOwner, errLimited.Scope)
	}
	// Global 2/3, PerOwner 1/1 exceeded

	ctx3 := contexts.WithCRE(t.Context(), contexts.CRE{Owner: user3String})
	require.NoError(t, wsl.Use(ctx3, 1))
	// Global 3/3, PerOwner 1/1 (one each user)

	require.Error(t, wsl.Use(ctx2, 1))
	// Global 3/3, PerOwner 1/1 Global and PerOwner exceeded

	require.NoError(t, wsl.Free(ctx2, 1))
	// Global 2/3, User2 PerOwner 0/1

	require.NoError(t, wsl.Use(ctx2, 1))
	// Global 3/3, PerOwner 1/1 (one each user)

	ctx4 := contexts.WithCRE(t.Context(), contexts.CRE{Owner: user4String})
	err = wsl.Use(ctx4, 1)
	require.Error(t, err)
	if assert.ErrorAs(t, err, &errLimited) {
		require.Equal(t, settings.ScopeGlobal, errLimited.Scope)
	}
	// Global 3/3, PerOwner 0/1 Global exceeded

	ctx5 := contexts.WithCRE(t.Context(), contexts.CRE{Owner: user5String})
	err = wsl.Use(ctx5, 1)
	require.Error(t, err)
	if assert.ErrorAs(t, err, &errLimited) {
		require.Equal(t, settings.ScopeGlobal, errLimited.Scope)
	}
	// Global 3/3, PerOwner 0/1 Global exceeded

	// Drop global limit
	require.NoError(t, wsl.Free(ctx1, 1))
	require.NoError(t, wsl.Free(ctx2, 1))
	require.NoError(t, wsl.Free(ctx3, 1))
	// Global 0/3

	// add external owner
	require.NoError(t, wsl.Use(ctx5, 1))
	// Global 1/3, PerOwner 1/2

	require.NoError(t, wsl.Use(ctx5, 1))
	// Global 2/3, PerOwner 2/2 Override allows 2

	err = wsl.Use(ctx5, 1)
	require.Error(t, err)
	if assert.ErrorAs(t, err, &errLimited) {
		require.Equal(t, settings.ScopeOwner, errLimited.Scope)
	}
	// Global 2/3, PerOwner 2/2 Override exceeded
}
