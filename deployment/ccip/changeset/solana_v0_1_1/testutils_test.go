package solana_test

import (
	"os"
	"testing"
)

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci && false { // test - disable
		t.Skip("Skipping in CI")
	}
}
