package pgtest

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/internal/pgtestenv"
)

// EnvDisableAutoContainers mirrors [pgtestenv.EnvDisableAutoContainers].
const EnvDisableAutoContainers = pgtestenv.EnvDisableAutoContainers

// EnsureAutoPostgres delegates to [pgtestenv.EnsureAutoPostgres] (importable from module root tests).
func EnsureAutoPostgres(t testing.TB) {
	pgtestenv.EnsureAutoPostgres(t)
}
