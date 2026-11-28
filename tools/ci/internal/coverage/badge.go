package coverage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BadgeOptions describes the report and output files used to publish a badge.
type BadgeOptions struct {
	ReportPath, SVGPath, StatePath string
	SHA, Artifact, Repo, RunURL    string
	Now                            time.Time
}

// BadgeState is the persistent state stored beside the coverage badge.
type BadgeState struct {
	SHA              string   `json:"sha"`
	Repo             string   `json:"repo"`
	GlobalRaw        float64  `json:"global_raw"`
	GlobalSourceOnly float64  `json:"global_source_only"`
	Patch            *float64 `json:"patch"`
	ReportArtifact   string   `json:"report_artifact"`
	RunURL           string   `json:"run_url"`
	UpdatedAt        string   `json:"updated_at"`
}

type badgePercent struct {
	Percent *float64 `json:"Percent"`
}

type badgeReport struct {
	Global           *badgePercent `json:"Global"`
	GlobalSourceOnly *badgePercent `json:"GlobalSourceOnly"`
	Patch            *badgePercent `json:"Patch"`
}

// GenerateBadge writes the SVG and JSON state for a merged coverage report.
func GenerateBadge(opts BadgeOptions) (BadgeState, error) {
	var state BadgeState
	if opts.ReportPath == "" {
		return state, errors.New("report path is required")
	}
	if opts.SVGPath == "" {
		return state, errors.New("SVG path is required")
	}
	if opts.StatePath == "" {
		return state, errors.New("state path is required")
	}

	data, err := os.ReadFile(opts.ReportPath)
	if err != nil {
		return state, fmt.Errorf("read coverage report: %w", err)
	}
	var report badgeReport
	if err := json.Unmarshal(data, &report); err != nil {
		return state, fmt.Errorf("decode coverage report: %w", err)
	}

	source := report.GlobalSourceOnly
	if source == nil {
		source = report.Global
	}
	percent := badgePercentValue(source)
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	state = BadgeState{
		SHA:              opts.SHA,
		Repo:             opts.Repo,
		GlobalRaw:        badgePercentValue(report.Global),
		GlobalSourceOnly: percent,
		ReportArtifact:   opts.Artifact,
		RunURL:           opts.RunURL,
		UpdatedAt:        now.UTC().Truncate(time.Second).Format(time.RFC3339),
	}
	if report.Patch != nil {
		state.Patch = report.Patch.Percent
	}

	value := fmt.Sprintf("%.1f%%", percent)
	color := badgeColor(percent)
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="160" height="20" role="img" aria-label="coverage: %s"><rect width="160" height="20" rx="3" fill="#555"/><rect x="80" width="80" height="20" rx="3" fill="%s"/><rect x="80" width="4" height="20" fill="%s"/><g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11"><text x="40" y="14">coverage</text><text x="120" y="14">%s</text></g></svg>`, value, color, color, value) + "\n"
	if err := os.WriteFile(opts.SVGPath, []byte(svg), 0o644); err != nil {
		return BadgeState{}, fmt.Errorf("write coverage badge: %w", err)
	}
	stateData, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return BadgeState{}, fmt.Errorf("encode coverage state: %w", err)
	}
	stateData = append(stateData, '\n')
	if err := os.WriteFile(opts.StatePath, stateData, 0o644); err != nil {
		return BadgeState{}, fmt.Errorf("write coverage state: %w", err)
	}
	return state, nil
}

func badgePercentValue(stats *badgePercent) float64 {
	if stats == nil || stats.Percent == nil {
		return 0
	}
	return *stats.Percent
}

func badgeColor(percent float64) string {
	switch {
	case percent >= 70:
		return "#2e7d32"
	case percent >= 55:
		return "#f9a825"
	case percent >= 45:
		return "#ef6c00"
	default:
		return "#c62828"
	}
}

// FullReport reports whether all profiles required for a badge are present.
func FullReport(dir string) (bool, error) {
	for _, name := range []string{"unit.txt", "integration.txt", "e2e.txt", "any.json"} {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, fmt.Errorf("stat %s: %w", name, err)
		}
		if !info.Mode().IsRegular() {
			return false, nil
		}
	}
	return true, nil
}
