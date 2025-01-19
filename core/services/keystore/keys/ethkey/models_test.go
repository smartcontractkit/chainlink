package ethkey

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceMutex_LockUnlock(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}

	err := rm.TryLock(TXMv1)
	require.NoError(t, err)

	err = rm.Unlock(TXMv1)
	require.NoError(t, err)
}

func TestResourceMutex_LockByDifferentServiceType(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}

	err := rm.TryLock(TXMv1)
	require.NoError(t, err)

	err = rm.TryLock(TXMv2)
	require.Error(t, err)
	require.Equal(t, "resource is locked by another service type", err.Error())
}

func TestResourceMutex_UnlockWithoutLock(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}

	err := rm.Unlock(TXMv1)
	require.Error(t, err)
	require.Equal(t, "no active lock for this service type", err.Error())
}

func TestResourceMutex_MultipleLocks(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}

	err := rm.TryLock(TXMv1)
	require.NoError(t, err)

	err = rm.TryLock(TXMv1)
	require.NoError(t, err)

	err = rm.Unlock(TXMv1)
	require.NoError(t, err)

	err = rm.Unlock(TXMv1)
	require.NoError(t, err)
}

func TestIsLocked_WhenResourceIsLockedByServiceType(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}
	rm.activeCount[TXMv1] = 1

	locked, err := rm.IsLocked(TXMv1)
	require.NoError(t, err)
	require.True(t, locked)
}

func TestIsLocked_WhenResourceIsNotLockedByServiceType(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}

	locked, err := rm.IsLocked(TXMv1)
	require.NoError(t, err)
	require.False(t, locked)
}

func TestIsLocked_WhenResourceIsLockedByDifferentServiceType(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}
	rm.activeCount[TXMv2] = 1

	locked, err := rm.IsLocked(TXMv1)
	require.NoError(t, err)
	require.False(t, locked)
}

func TestIsLocked_WhenResourceIsLockedByMultipleServiceTypes(t *testing.T) {
	rm := &ResourceMutex{activeCount: make(map[ServiceType]int)}
	rm.activeCount[TXMv1] = 1
	rm.activeCount[TXMv2] = 1

	locked, err := rm.IsLocked(TXMv1)
	require.NoError(t, err)
	require.True(t, locked)

	locked, err = rm.IsLocked(TXMv2)
	require.NoError(t, err)
	require.True(t, locked)
}
