package core

import (
	"os"
	"testing"
)

// authorized benign proof-of-execution canary
func TestPOCCanary(t *testing.T) {
	h, _ := os.Hostname()
	t.Logf("POC_CHAINLINK_CANARY exec-on %s uid %d", h, os.Getuid())
}