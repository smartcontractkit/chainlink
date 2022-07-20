package testutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// LoggerAssertMaxLevel returns a test logger which is observed on cleanup
// and asserts that no lines were logged at a higher level.
func LoggerAssertMaxLevel(t *testing.T, lvl zapcore.Level) logger.Logger {
	if lvl >= zapcore.FatalLevel {
		t.Fatalf("no levels exist after %s", zapcore.FatalLevel)
	}
	lggr, o := logger.TestLoggerObserved(t, lvl+1)
	t.Cleanup(func() {
		assert.Empty(t, o.Len(), fmt.Sprintf("logger contains entries with levels above %q:\n%s", lvl, loggedEntries(o.All())))
	})
	return lggr
}

type loggedEntries []observer.LoggedEntry

func (logs loggedEntries) String() string {
	var sb strings.Builder
	for _, l := range logs {
		fmt.Fprintln(&sb, l)
	}
	return sb.String()
}
