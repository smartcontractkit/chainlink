package dbdetect

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func looksLikeGoPackagePattern(arg string) bool {
	return strings.Contains(arg, ".") ||
		strings.Contains(arg, "/") ||
		strings.Contains(arg, "...")
}

var testBinaryTwoArgSuffixFlags = map[string]bool{
	"-run":      true,
	"-bench":    true,
	"-skip":     true,
	"-fuzz":     true,
	"-count":    true,
	"-timeout":  true,
	"-tags":     true,
	"-parallel": true,
}

// harnessRootValueFlags are testrig root flags (see dbflags.Register) that take a
// separate value token and must not be treated as Go package patterns.
var harnessRootValueFlags = map[string]bool{
	"database-url":     true,
	"postgres-version": true,
}

func extractPackagePatterns(args []string) []string {
	var patterns []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			name := strings.Split(strings.TrimLeft(arg, "-"), "=")[0]
			if (testBinaryTwoArgSuffixFlags["-"+name] || harnessRootValueFlags[name]) && !strings.Contains(arg, "=") {
				i++
			}
			continue
		}
		if arg == "diagnose" || arg == "gotestsum" || arg == "init-skill" || arg == "show-skill" {
			continue
		}
		if looksLikeGoPackagePattern(arg) {
			patterns = append(patterns, arg)
		}
	}
	return patterns
}

// NeedsPostgres checks if Postgres is needed for the given arguments and repository root.
func NeedsPostgres(repoRoot string, args []string) (bool, error) {
	// 1. Check for -short flag.
	for _, arg := range args {
		if arg == "-short" || arg == "--short" || strings.HasPrefix(arg, "-short=") || strings.HasPrefix(arg, "--short=") {
			return false, nil
		}
	}

	patterns := extractPackagePatterns(args)
	if len(patterns) == 0 {
		return true, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 2. Run go list -deps -test
	goArgs := append([]string{"list", "-deps", "-test"}, patterns...)
	cmd := exec.CommandContext(ctx, "go", goArgs...)
	cmd.Dir = repoRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return true, fmt.Errorf("go list: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	targetDep := "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	needsDB := strings.Contains(stdout.String(), targetDep)

	return needsDB, nil
}
