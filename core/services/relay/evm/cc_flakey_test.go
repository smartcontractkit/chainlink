package evm_test

import (
	"bytes"
	"os/exec"
	"testing"
	"time"
)

func TestRunChainComponentsMultipleTimes(t *testing.T) {
	// Configuration for the test run
	testName := "TestChainComponents"
	numRuns := 100
	failCount := 0
	totalTime := time.Duration(0)

	for i := 1; i <= numRuns; i++ {
		t.Logf("Running %s: iteration %d/%d...\n", testName, i, numRuns)

		// Clear the Go test cache
		cleanCmd := exec.Command("go", "clean", "-testcache")
		cleanErr := cleanCmd.Run()
		if cleanErr != nil {
			t.Fatalf("Failed to clear test cache: %v", cleanErr)
		}

		// Start the timer
		start := time.Now()

		// Run the test via `go test`
		cmd := exec.Command("go", "test", "-run", testName, "./...")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()

		// Record the elapsed time
		elapsed := time.Since(start)
		totalTime += elapsed

		if err != nil {
			failCount++
			t.Errorf("Iteration %d failed with error: %s\n", i, err)
			t.Logf("Test output:\n%s\n", out.String())
			t.Logf("Test error output:\n%s\n", stderr.String())
		} else {
			t.Logf("Iteration %d passed.\n", i)
		}

		t.Logf("Test duration: %s\n", elapsed)
	}

	// Calculate the average test time
	averageTime := totalTime / time.Duration(numRuns)

	// Print results
	t.Logf("\nResults for %s:\n", testName)
	t.Logf("Total runs: %d\n", numRuns)
	t.Logf("Failures: %d\n", failCount)
	t.Logf("Pass rate: %.2f%%\n", float64(numRuns-failCount)/float64(numRuns)*100)
	t.Logf("Total time: %s\n", totalTime)
	t.Logf("Average test time: %s\n", averageTime)

	// Optional: Fail the entire test if there were any failures
	if failCount > 0 {
		t.Fatalf("Test failed %d/%d times.\n", failCount, numRuns)
	}
}
