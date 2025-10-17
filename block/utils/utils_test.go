package utils

import (
	"testing"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
)

func TestKeyedMutex(t *testing.T) {
	t.Parallel()

	var km KeyedMutex
	unlock1 := km.LockInt64(1)
	unlock2 := km.LockInt64(2)

	awaiter := testutils.NewAwaiter()
	go func() {
		km.LockInt64(1)()
		km.LockInt64(2)()
		awaiter.ItHappened()
	}()

	unlock2()
	unlock1()
	awaiter.AwaitOrFail(t)
}
