package dbdetect

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// safeGoListArg matches argv tokens built from allowlisted go list flags and package patterns.
var safeGoListArg = regexp.MustCompile(`^[\w./@=,\-*:+$%()\[\]!~]+$`)

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

// goListTwoArgFlags are forwarded to `go list` so build tags and module
// settings match the test invocation.
var goListTwoArgFlags = map[string]bool{
	"tags":    true,
	"mod":     true,
	"modfile": true,
	"overlay": true,
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

func extractGoListFlags(args []string) []string {
	var flags []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		name := strings.TrimLeft(arg, "-")
		if strings.HasPrefix(name, "-") {
			name = strings.TrimLeft(name, "-")
		}
		value, hasEq := "", false
		if eq := strings.Index(name, "="); eq >= 0 {
			value = name[eq+1:]
			name = name[:eq]
			hasEq = true
		}
		if !goListTwoArgFlags[name] {
			continue
		}
		if hasEq {
			flags = append(flags, "-"+name+"="+value)
			continue
		}
		if i+1 >= len(args) {
			continue
		}
		flags = append(flags, "-"+name, args[i+1])
		i++
	}
	return flags
}

func buildGoListArgs(patterns, args []string) []string {
	listFlags := extractGoListFlags(args)
	goArgs := make([]string, 0, 3+len(listFlags)+len(patterns))
	goArgs = append(goArgs, "list", "-deps", "-test")
	goArgs = append(goArgs, listFlags...)
	goArgs = append(goArgs, patterns...)
	return goArgs
}

func validateGoListArgs(goArgs []string) error {
	for _, arg := range goArgs {
		if arg == "" || !safeGoListArg.MatchString(arg) {
			return fmt.Errorf("invalid go list argument %q", arg)
		}
	}
	return nil
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

	goArgs := buildGoListArgs(patterns, args)
	if err := validateGoListArgs(goArgs); err != nil {
		return true, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
