package syncerlimiter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkflowLimits(t *testing.T) {
	t.Parallel()

	config := Config{
		Global:   3,
		PerOwner: 1,
		PerOwnerOverrides: map[string]int32{
			"ext-owner": 2,
		},
	}
	wsl, err := NewWorkflowLimits(config)
	require.Equal(t, int32(3), wsl.config.Global)
	require.Equal(t, int32(1), wsl.config.PerOwner)
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

	allowOwner, allowGlobal = wsl.Allow("ext-owner")
	require.True(t, allowOwner)
	require.False(t, allowGlobal)
	// Global 3/3, PerOwner 0/1 Global exceeded

	// Drop global limit
	wsl.Decrement("user1")
	wsl.Decrement("user2")
	wsl.Decrement("user3")
	// Global 0/3

	// add external owner
	allowOwner, allowGlobal = wsl.Allow("ext-owner")
	require.True(t, allowOwner && allowGlobal)
	// Global 1/3, PerOwner 1/2

	allowOwner, allowGlobal = wsl.Allow("ext-owner")
	require.True(t, allowOwner && allowGlobal)
	// Global 2/3, PerOwner 2/2 Override allows 2

	allowOwner, allowGlobal = wsl.Allow("ext-owner")
	require.False(t, allowOwner)
	require.True(t, allowGlobal)
	// Global 2/3, PerOwner 2/2 Override exceeded
}

func TestLimits_getPerOwnerLimit(t *testing.T) {
	config := Config{PerOwner: defaultPerOwner}

	tests := []struct {
		name      string
		limits    *Limits
		owner     string
		wantLimit int32
	}{
		{
			name: "no overrides",
			limits: func() *Limits {
				l, err := NewWorkflowLimits(config)
				require.NoError(t, err)
				return l
			}(),
			owner:     "owner1",
			wantLimit: defaultPerOwner,
		},
		{
			name: "override exists",
			limits: func() *Limits {
				config.PerOwnerOverrides = map[string]int32{
					"owner2": 20,
				}
				l, err := NewWorkflowLimits(config)
				require.NoError(t, err)
				return l
			}(),
			owner:     "owner2",
			wantLimit: 20,
		},
		{
			name: "override does not exist",
			limits: func() *Limits {
				config.PerOwnerOverrides = map[string]int32{
					"owner2": 20,
				}
				l, err := NewWorkflowLimits(config)
				require.NoError(t, err)
				return l
			}(),
			owner:     "owner3",
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
