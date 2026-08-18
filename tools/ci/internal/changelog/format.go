package changelog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/githuboutput"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/paths"
)

const (
	maxPRDescLength = 65000
	prHeader        = "This PR is a preview of the changes that will be included in the next release. Please do not merge this PR.\n---\n"
	prTruncatedMsg  = "The changelog content is too long for the PR description. Please view the full changelog in the [CHANGELOG.md](https://github.com/smartcontractkit/chainlink/blob/changesets/release-preview/CHANGELOG.md)"
)

// Ordered list of tags to group changelog entries by.
var tagsList = []string{
	"#nops",
	"#added",
	"#changed",
	"#removed",
	"#updated",
	"#deprecation_notice",
	"#breaking_change",
	"#db_update",
	"#wip",
	"#bugfix",
	"#internal",
	"#untagged",
}

// Result holds the parsed and generated changelog output.
type Result struct {
	Version      string
	PRBody       string
	NewChangelog string
}

// ReadVersionFromPackageJSON parses the "version" field from package.json.
func ReadVersionFromPackageJSON(pkgJSONPath string) (string, error) {
	resolved := paths.ResolveFromRepoRoot(pkgJSONPath)
	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", pkgJSONPath, err)
	}

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", pkgJSONPath, err)
	}
	if pkg.Version == "" {
		return "", fmt.Errorf("version field empty in %s", pkgJSONPath)
	}
	return pkg.Version, nil
}

// Format processes the CHANGELOG.md file, groups new entries by tag, updates the file, and sets GITHUB_OUTPUT.
func Format(changelogPath, packageJSONPath string, writeGithubOutput bool) (*Result, error) {
	version, err := ReadVersionFromPackageJSON(packageJSONPath)
	if err != nil {
		return nil, err
	}

	resolvedChangelog := paths.ResolveFromRepoRoot(changelogPath)
	contentBytes, err := os.ReadFile(resolvedChangelog)
	if err != nil {
		return nil, fmt.Errorf("failed to read changelog %s: %w", changelogPath, err)
	}

	fullText := string(contentBytes)
	reader := bufio.NewReader(strings.NewReader(fullText))

	var (
		inCurrentVersion = false
		currentEntries   []string
		currentEntry     strings.Builder
		previousHistory  strings.Builder
		pastFirstVersion = false
	)

	for {
		lineWithNewline, err := reader.ReadString('\n')
		if err != nil && len(lineWithNewline) == 0 {
			break
		}
		line := strings.TrimRight(lineWithNewline, "\r\n")

		if strings.HasPrefix(line, "## ") {
			if inCurrentVersion {
				// Second version header encountered - start of previous history
				if currentEntry.Len() > 0 {
					currentEntries = append(currentEntries, strings.TrimSpace(currentEntry.String()))
					currentEntry.Reset()
				}
				inCurrentVersion = false
				pastFirstVersion = true
			} else if !pastFirstVersion {
				// First version header encountered
				inCurrentVersion = true
				continue
			}
		}

		if pastFirstVersion {
			previousHistory.WriteString(line)
			previousHistory.WriteString("\n")
			if err != nil {
				break
			}
			continue
		}

		if inCurrentVersion {
			if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
				if currentEntry.Len() > 0 {
					currentEntries = append(currentEntries, strings.TrimSpace(currentEntry.String()))
					currentEntry.Reset()
				}
				currentEntry.WriteString(line)
			} else if currentEntry.Len() > 0 {
				currentEntry.WriteString("\n")
				currentEntry.WriteString(line)
			}
		}

		if err != nil {
			break
		}
	}

	if currentEntry.Len() > 0 {
		currentEntries = append(currentEntries, strings.TrimSpace(currentEntry.String()))
	}

	// Group entries by tag
	taggedMap := make(map[string][]string)
	for _, rawTag := range tagsList {
		tagKey := strings.TrimPrefix(rawTag, "#")
		taggedMap[tagKey] = []string{}
	}

	for _, entry := range currentEntries {
		matched := false
		for _, rawTag := range tagsList {
			if rawTag == "#untagged" {
				continue
			}
			if strings.Contains(entry, rawTag) {
				tagKey := strings.TrimPrefix(rawTag, "#")
				taggedMap[tagKey] = append(taggedMap[tagKey], entry)
				matched = true
			}
		}
		if !matched && len(entry) > 0 {
			taggedMap["untagged"] = append(taggedMap["untagged"], entry)
		}
	}

	// Build grouped changelog section
	var changelogSection strings.Builder
	changelogSection.WriteString("# Changelog Chainlink Core\n\n")
	fmt.Fprintf(&changelogSection, "## %s - PREVIEW\n", version)

	for _, rawTag := range tagsList {
		tagKey := strings.TrimPrefix(rawTag, "#")
		entries := taggedMap[tagKey]
		if len(entries) == 0 {
			continue
		}

		fmt.Fprintf(&changelogSection, "\n## %s\n\n", tagKey)
		changelogSection.WriteString(strings.Join(entries, "\n\n"))
		changelogSection.WriteString("\n")
	}

	previewChangelog := changelogSection.String()

	// Assemble final CHANGELOG.md content
	var finalChangelog strings.Builder
	finalChangelog.WriteString(previewChangelog)
	if previousHistory.Len() > 0 {
		finalChangelog.WriteString("\n")
		finalChangelog.WriteString(strings.TrimSpace(previousHistory.String()))
		finalChangelog.WriteString("\n")
	}

	newChangelogStr := finalChangelog.String()

	// Build PR body with length check
	var prBody string
	if len(previewChangelog) > maxPRDescLength {
		prBody = prHeader + prTruncatedMsg
	} else {
		prBody = prHeader + previewChangelog
	}

	// Write new changelog back to disk
	//nolint:gosec // CHANGELOG.md is committed to the repo and must remain readable
	if err := os.WriteFile(resolvedChangelog, []byte(newChangelogStr), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write updated %s: %w", resolvedChangelog, err)
	}

	// Optionally write to GITHUB_OUTPUT
	if writeGithubOutput {
		if err := githuboutput.AppendVar("version", version); err != nil {
			return nil, fmt.Errorf("failed to write version to GITHUB_OUTPUT: %w", err)
		}
		if err := githuboutput.AppendMultilineVar("pr_body", prBody); err != nil {
			return nil, fmt.Errorf("failed to write pr_body to GITHUB_OUTPUT: %w", err)
		}
	}

	return &Result{
		Version:      version,
		PRBody:       prBody,
		NewChangelog: newChangelogStr,
	}, nil
}
