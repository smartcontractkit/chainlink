package dbdetect

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/smartcontractkit/testrig/modresolve"
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

// IsDiagnoseCommand reports whether argv invokes the testrig diagnose subcommand.
func IsDiagnoseCommand(args []string) bool {
	return slices.Contains(args, "diagnose")
}

// PackageSlug returns a short docker-safe name for the package patterns in argv
// (e.g. ./core/services/... -> core_services).
func PackageSlug(args []string) string {
	patterns := extractPackagePatterns(args)
	switch len(patterns) {
	case 0:
		return "pkgs"
	case 1:
		return patternToSlug(patterns[0])
	default:
		slugs := make([]string, len(patterns))
		for i, p := range patterns {
			slugs[i] = patternToSlug(p)
		}
		return strings.Join(slugs, "__")
	}
}

func patternToSlug(pattern string) string {
	t := strings.TrimPrefix(pattern, "./")
	switch {
	case t == "...":
		return "pkgs"
	case strings.HasSuffix(t, "/..."):
		t = strings.TrimSuffix(t, "/...")
	}
	t = strings.ReplaceAll(t, "/", "_")
	return sanitizeSlugToken(t)
}

func sanitizeSlugToken(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '_', r == '-', r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := b.String()
	if out == "" {
		return "pkgs"
	}
	return out
}

func extractPackagePatterns(args []string) []string {
	var patterns []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			name, _, _ := strings.Cut(strings.TrimLeft(arg, "-"), "=")
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

	moduleDir, adjustedPatterns, err := modresolve.ResolvePatterns(repoRoot, patterns)
	if err != nil {
		return true, err
	}

	goArgs := buildGoListArgs(adjustedPatterns, args)
	if err := validateGoListArgs(goArgs); err != nil {
		return true, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", goArgs...)
	cmd.Dir = moduleDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return true, fmt.Errorf("go list: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	deps := stdout.String()
	for _, targetDep := range postgresTestDeps {
		if strings.Contains(deps, targetDep) {
			return true, nil
		}
	}

	return false, nil
}

// postgresTestDeps lists packages that imply a real Postgres server (CL_DATABASE_URL).
// go list -deps -test must match at least one for the testrig harness to start Postgres.
var postgresTestDeps = []string{
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest",
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight",
	"github.com/smartcontractkit/chainlink/v2/internal/testdb",
}
