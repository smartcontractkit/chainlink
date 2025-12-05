package solana_test

import (
	"os"
	"testing"
)

func skipInCI(t *testing.T) {
	t.Helper()

	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}
