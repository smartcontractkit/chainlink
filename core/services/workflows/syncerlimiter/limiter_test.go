package syncerlimiter

import (
	"testing"

	"github.com/stretchr/testify/require"

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
			"0x519BFD3D78fbb740c614432975CBE829E26C490e": 2,
		},
	}
	wsl, err := NewWorkflowLimits(lggr, config)
	require.Equal(t, int32(3), wsl.config.Global)
	require.Equal(t, int32(1), wsl.config.PerOwner)
	require.NoError(t, err)

	allowOwner, allowGlobal := wsl.Allow(user1String)
	require.True(t, allowOwner && allowGlobal)
	// Global 1/3, PerOwner 1/1

	allowOwner, allowGlobal = wsl.Allow(user2String)
	require.True(t, allowOwner && allowGlobal)
	// Global 2/3, PerOwner 1/1

	allowOwner, allowGlobal = wsl.Allow(user1String)
	require.True(t, allowGlobal)
	require.False(t, allowOwner)
	// Global 2/3, PerOwner 1/1 exceeded

	allowOwner, allowGlobal = wsl.Allow(user3String)
	require.True(t, allowOwner && allowGlobal)
	// Global 3/3, PerOwner 1/1 (one each user)

	allowOwner, allowGlobal = wsl.Allow(user2String)
	require.False(t, allowOwner)
	require.False(t, allowGlobal)
	// Global 3/3, PerOwner 1/1 Global and PerOwner exceeded

	wsl.Decrement(user2String)
	// Global 2/3, User2 PerOwner 0/1

	allowOwner, allowGlobal = wsl.Allow(user2String)
	require.True(t, allowOwner && allowGlobal)
	// Global 3/3, PerOwner 1/1 (one each user)

	wsl.Decrement(user4String)
	allowOwner, allowGlobal = wsl.Allow(user4String)
	require.True(t, allowOwner)
	require.False(t, allowGlobal)
	// Global 3/3, PerOwner 0/1 Global exceeded

	allowOwner, allowGlobal = wsl.Allow(user5String)
	require.True(t, allowOwner)
	require.False(t, allowGlobal)
	// Global 3/3, PerOwner 0/1 Global exceeded

	// Drop global limit
	wsl.Decrement(user1String)
	wsl.Decrement(user2String)
	wsl.Decrement(user3String)
	// Global 0/3

	// add external owner
	allowOwner, allowGlobal = wsl.Allow(user5String)
	require.True(t, allowOwner && allowGlobal)
	// Global 1/3, PerOwner 1/2

	allowOwner, allowGlobal = wsl.Allow(user5String)
	require.True(t, allowOwner && allowGlobal)
	// Global 2/3, PerOwner 2/2 Override allows 2

	allowOwner, allowGlobal = wsl.Allow(user5String)
	require.False(t, allowOwner)
	require.True(t, allowGlobal)
	// Global 2/3, PerOwner 2/2 Override exceeded
}

func TestLimits_getPerOwnerLimit(t *testing.T) {
	config := Config{PerOwner: defaultPerOwner}
	lggr := logger.TestLogger(t)

	tests := []struct {
		name      string
		limits    *Limits
		owner     string
		wantLimit int32
	}{
		{
			name: "no overrides",
			limits: func() *Limits {
				l, err := NewWorkflowLimits(lggr, config)
				require.NoError(t, err)
				return l
			}(),
			owner:     user1String,
			wantLimit: defaultPerOwner,
		},
		{
			name: "override exists",
			limits: func() *Limits {
				config.PerOwnerOverrides = map[string]int32{
					"0x" + user2String: 20,
				}
				l, err := NewWorkflowLimits(lggr, config)
				require.NoError(t, err)
				return l
			}(),
			owner:     user2String,
			wantLimit: 20,
		},
		{
			name: "ignores checksum",
			limits: func() *Limits {
				config.PerOwnerOverrides = map[string]int32{
					"0x1b5B31d3AB780cd192DE677707EA5A3Cb53c6d67": 8,
					"0xcec984E1dA8Ad5C012bEE7C4abc334F07Ee61A6e": 0,
				}
				l, err := NewWorkflowLimits(lggr, config)
				require.NoError(t, err)
				return l
			}(),
			owner:     "1b5b31d3ab780cd192de677707ea5a3cb53c6d67",
			wantLimit: 8,
		},
		{
			name: "override does not exist",
			limits: func() *Limits {
				config.PerOwnerOverrides = map[string]int32{
					"0x" + user2String: 20,
				}
				l, err := NewWorkflowLimits(lggr, config)
				require.NoError(t, err)
				return l
			}(),
			owner:     user3String,
			wantLimit: defaultPerOwner,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLimit := tt.limits.getPerOwnerLimit(tt.owner); gotLimit != tt.wantLimit {
				t.Errorf("getPerOwnerLimit() = %v, want %v", gotLimit, tt.wantLimit)
			}
		})
	}
}
