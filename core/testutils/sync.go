package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils/tests"
)

// TestInterval is just a sensible poll interval that gives fast tests without
// risk of spamming
const TestInterval = 100 * time.Millisecond

// AssertEventually waits for f to return true
func AssertEventually(t *testing.T, f func() bool) {
	assert.Eventually(t, f, tests.WaitTimeout(t), TestInterval/2)
}
