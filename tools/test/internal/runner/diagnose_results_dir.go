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
func diagnoseResultsDirName(conf *config.App, target string, now time.Time) string {
	tsPart := now.Format("20060102150405")
	for phase := range 8 {
		cfg := diagnoseConfigDirPartPhase(conf, phase)
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

func diagnoseRunDirToken(run string) (full string, hashOnly string) {
	h := sha256.Sum256([]byte(run))
	hash8 := hex.EncodeToString(h[:4])
	hashOnly = "r" + hash8

	const maxRunes = 40
	rs := []rune(run)
	if len(rs) > maxRunes {
		rs = rs[:maxRunes]
	}
	short := sanitizeDirToken(string(rs))
	if short == "" {
		return hashOnly, hashOnly
	}
	return "r" + short + "-" + hash8, hashOnly
}

func durationDirToken(d time.Duration) string {
	return strings.ReplaceAll(d.String(), ":", "_")
}

func diagnoseConfigDirPartPhase(conf *config.App, phase int) string {
	dropSlow := phase >= 1
	runHashOnly := phase >= 2
	dropCPU := phase >= 3
	dropPar := phase >= 4
	dropFF := phase >= 5
	dropShuffle := phase >= 6
	dropRace := phase >= 7

	var parts []string
	if conf.Iterations > 0 {
		parts = append(parts, fmt.Sprintf("it%d", conf.Iterations))
	}
	parts = append(parts, "to"+durationDirToken(conf.Timeout))
	if !dropRace && conf.Race {
		parts = append(parts, "race")
	}
	if !dropShuffle && conf.Shuffle {
		parts = append(parts, "shuffle")
	}
	if !dropFF && conf.FailFast {
		parts = append(parts, "ff")
	}
	if !dropPar && conf.Parallel > 0 {
		parts = append(parts, fmt.Sprintf("p%d", conf.Parallel))
	}
	if !dropCPU && conf.CPU != "" {
		parts = append(parts, "cpu-"+strings.ReplaceAll(conf.CPU, ",", "-"))
	}
	if conf.Run != "" {
		full, hash := diagnoseRunDirToken(conf.Run)
		if runHashOnly {
			parts = append(parts, hash)
		} else {
			parts = append(parts, full)
		}
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
		r, _ := utf8.DecodeLastRuneInString(s)
		if r != utf8.RuneError {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}
