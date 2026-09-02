package matrix

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var sanitizeRegex = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// SanitizeTestID sanitizes test name to be a legal artifact name.
func SanitizeTestID(name string) string {
	return sanitizeRegex.ReplaceAllString(name, "_")
}

// InMemoryOptions contains parameters for generating in-memory test matrix.
type InMemoryOptions struct {
	ConfigFile string
	RunID      string
	RunAttempt string
	SpotFlag   string
}

// RawInMemoryEntry matches structure in .github/in-memory-tests.json.
type RawInMemoryEntry struct {
	Name     string `json:"name"`
	Test     string `json:"test"`
	Timeout  string `json:"timeout"`
	Parallel int    `json:"parallel"`
	Plugins  bool   `json:"plugins"`
	RunsOn   string `json:"runs_on"`
	FreeDisk bool   `json:"free_disk"`
	Aptos    string `json:"aptos"`
	Sui      string `json:"sui"`
}

// InMemoryEntry is the resolved test matrix entry.
type InMemoryEntry struct {
	Name       string `json:"name"`
	Test       string `json:"test"`
	Timeout    string `json:"timeout"`
	JobTimeout int    `json:"job_timeout"`
	Parallel   int    `json:"parallel"`
	Plugins    bool   `json:"plugins"`
	RunsOn     string `json:"runs_on"`
	TestID     string `json:"test_id"`
	FreeDisk   bool   `json:"free_disk"`
	Aptos      string `json:"aptos"`
	Sui        string `json:"sui"`
}

// BuildInMemoryMatrix parses in-memory test configuration and builds matrix.
func BuildInMemoryMatrix(ctx context.Context, opts InMemoryOptions) ([]InMemoryEntry, error) {
	data, err := os.ReadFile(opts.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read in-memory test config file %s: %w", opts.ConfigFile, err)
	}

	var rawEntries []RawInMemoryEntry
	if err := json.Unmarshal(data, &rawEntries); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", opts.ConfigFile, err)
	}

	spotFlag := opts.SpotFlag
	if spotFlag == "" {
		spotFlag = "spot=co"
	}
	runAttempt := opts.RunAttempt
	if runAttempt == "" {
		runAttempt = "1"
	}

	entries := make([]InMemoryEntry, len(rawEntries))
	for i, raw := range rawEntries {
		timeoutMinutes := 15
		trimmed := strings.TrimSuffix(raw.Timeout, "m")
		if val, err := strconv.Atoi(trimmed); err == nil {
			timeoutMinutes = val
		}
		jobTimeout := timeoutMinutes + 5

		var runsOn string
		if raw.RunsOn == "" {
			runsOn = "ubuntu-latest"
		} else {
			runsOn = fmt.Sprintf("runs-on=%s-%d-%s/%s/family=m7i+m8i/%s/image=ubuntu24-full-x64/extras=s3-cache+tmpfs",
				opts.RunID, i, runAttempt, raw.RunsOn, spotFlag)
		}

		entries[i] = InMemoryEntry{
			Name:       raw.Name,
			Test:       raw.Test,
			Timeout:    raw.Timeout,
			JobTimeout: jobTimeout,
			Parallel:   raw.Parallel,
			Plugins:    raw.Plugins,
			RunsOn:     runsOn,
			TestID:     SanitizeTestID(raw.Name),
			FreeDisk:   raw.FreeDisk,
			Aptos:      raw.Aptos,
			Sui:        raw.Sui,
		}
	}

	return entries, nil
}
