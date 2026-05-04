package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

const (
	diagnoseResultsNamePrefix  = "diagnose-"
	maxDiagnoseResultsBasename = 220
	defaultSlowThreshold       = 30 * time.Second
)

// diagnoseResultsDirName returns a repo-root-relative directory basename for
// diagnose output: diagnose-<targetSlug>-<config>-<YYYYMMDDHHMMSS>.
func diagnoseResultsDirName(conf *config.App, goTestArgs []string, now time.Time) string {
	tsPart := now.Format("20060102150405")
	target := guessPackagePatternForSlug(goTestArgs)
	for phase := range 8 {
		cfg := diagnoseConfigDirPartPhase(conf, goTestArgs, phase)
		tail := "-" + cfg + "-" + tsPart
		avail := max(maxDiagnoseResultsBasename-len(diagnoseResultsNamePrefix)-len(tail), 8)
		slug := truncateUTF8MaxBytes(diagnoseTargetSlug(target), avail)
		base := diagnoseResultsNamePrefix + slug + tail
		if len(base) <= maxDiagnoseResultsBasename {
			return base
		}
	}
	return diagnoseResultsNamePrefix + "x" + "-" + tsPart
}

func diagnoseTargetSlug(target string) string {
	t := strings.TrimPrefix(target, "./")
	switch {
	case t == "...":
		return "allpkgs"
	case strings.HasSuffix(t, "/..."):
		t = strings.TrimSuffix(t, "/...") + "_allpkgs"
	}
	t = strings.ReplaceAll(t, "/", "_")
	return sanitizeDirToken(t)
}

func sanitizeDirToken(s string) string {
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
	return b.String()
}

// guessPackagePatternForSlug picks a human-readable slug from go test arguments
// (trailing package patterns). Falls back to "pkgs" if none found.
func guessPackagePatternForSlug(goTestArgs []string) string {
	pkgs := packagePatternsFromEnd(goTestArgs)
	switch len(pkgs) {
	case 0:
		return "pkgs"
	case 1:
		return pkgs[0]
	default:
		return strings.Join(pkgs, "__")
	}
}

func durationDirToken(d time.Duration) string {
	return strings.ReplaceAll(d.String(), ":", "_")
}

func diagnoseConfigDirPartPhase(conf *config.App, goTestArgs []string, phase int) string {
	h := sha256.Sum256([]byte(strings.Join(goTestArgs, "\x00")))
	hash8 := hex.EncodeToString(h[:4])

	dropSlow := phase >= 1
	dropShuffle := phase >= 2
	dropFF := phase >= 3
	shortHash := phase >= 4

	var parts []string
	if conf.Iterations > 0 {
		parts = append(parts, fmt.Sprintf("it%d", conf.Iterations))
	}
	hStr := hash8
	if shortHash {
		hStr = hStr[:4]
	}
	parts = append(parts, "h"+hStr)
	if !dropFF && conf.FailFast {
		parts = append(parts, "ff")
	}
	if !dropShuffle && conf.Shuffle {
		parts = append(parts, "shuffle")
	}
	if !dropSlow {
		slow := conf.SlowThreshold
		if slow == 0 {
			slow = defaultSlowThreshold
		}
		if slow != defaultSlowThreshold {
			parts = append(parts, "slow"+durationDirToken(conf.SlowThreshold))
		}
	}
	return strings.Join(parts, "-")
}

func truncateUTF8MaxBytes(s string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(s) <= maxBytes {
		return s
	}
	s = s[:maxBytes]
	for len(s) > 0 {
		r, size := utf8.DecodeLastRuneInString(s)
		// RuneError is also the rune value U+FFFD; only strip when decoding hit invalid UTF-8 (size 1).
		if r != utf8.RuneError || size != 1 {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}
