package wasmtest

import (
	"sync"
	"testing"
)

func TestCreateTestBinary_CancelledContext(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	t.Run("subtest", func(subT *testing.T) {
		go func() {
			defer wg.Done()
			// Wait for the subtest to finish so its context is cancelled
			<-subT.Context().Done()

			// Call CreateTestBinary with the cancelled subtest T.
			// This will invoke the build using subT's cancelled context.
			// Under the old implementation, this fails/panics.
			// Under the new implementation, it should succeed because it doesn't use the test-scoped context.
			_ = CreateTestBinary("core/capabilities/compute/test/simple/cmd", false, subT)
		}()
	})

	wg.Wait()
}
