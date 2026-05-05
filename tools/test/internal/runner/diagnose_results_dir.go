package runner

import (
	"strings"
	"time"
	"unicode/utf8"
)

const (
	diagnoseResultsNamePrefix  = "diagnose-"
	maxDiagnoseResultsBasename = 220
)

// diagnoseResultsDirName returns a repo-root-relative directory basename for
// diagnose output: diagnose-<targetSlug>-<YYYYMMDDHHMMSS>. Full argv and harness
// flags live in report.json under the run key (see RunMeta).
func diagnoseResultsDirName(goTestArgs []string, now time.Time) string {
	tsPart := now.Format("20060102150405")
	target := guessPackagePatternForSlug(goTestArgs)
	slug := diagnoseTargetSlug(target)
	tail := "-" + tsPart
	avail := max(maxDiagnoseResultsBasename-len(diagnoseResultsNamePrefix)-len(tail), 1)
	slug = truncateUTF8MaxBytes(slug, avail)
	if slug == "" {
		slug = "x"
	}
	base := diagnoseResultsNamePrefix + slug + tail
	if len(base) <= maxDiagnoseResultsBasename {
		return base
	}
	return diagnoseResultsNamePrefix + "x" + tail
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
