package changelog

import (
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
	lines := strings.Split(fullText, "\n")

	var (
		inCurrentVersion = false
		currentEntries   []string
		currentEntry     strings.Builder
		pastChangelog    strings.Builder
		pastStarted      bool
	)

	versionHeader := fmt.Sprintf("## %s", version)

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if strings.TrimSpace(line) == versionHeader {
				inCurrentVersion = true
				continue
			} else if inCurrentVersion {
				// We hit the previous version header
				if currentEntry.Len() > 0 {
					currentEntries = append(currentEntries, strings.TrimRight(currentEntry.String(), "\n"))
					currentEntry.Reset()
				}
				inCurrentVersion = false
				pastStarted = true
				pastChangelog.WriteString(line + "\n")
				continue
			}
		}

		if pastStarted {
			pastChangelog.WriteString(line + "\n")
			continue
		}

		if inCurrentVersion {
			// Skip subheaders like "### Minor Changes", "### Patch Changes"
			if strings.HasPrefix(line, "### ") {
				continue
			}

			if strings.HasPrefix(line, "- ") {
				if currentEntry.Len() > 0 {
					currentEntries = append(currentEntries, strings.TrimRight(currentEntry.String(), "\n"))
					currentEntry.Reset()
				}
				currentEntry.WriteString(line + "\n")
			} else if currentEntry.Len() > 0 {
				currentEntry.WriteString(line + "\n")
			}
		}
	}

	if currentEntry.Len() > 0 {
		currentEntries = append(currentEntries, strings.TrimRight(currentEntry.String(), "\n"))
	}

	// Group entries by tags (an entry matching multiple tags appears under all matching tags)
	tagMap := make(map[string][]string)
	for _, entry := range currentEntries {
		matchedAny := false
		for _, tag := range tagsList {
			if tag == "#untagged" {
				continue
			}
			if strings.Contains(entry, tag) {
				tagMap[tag] = append(tagMap[tag], entry)
				matchedAny = true
			}
		}
		if !matchedAny {
			tagMap["#untagged"] = append(tagMap["#untagged"], entry)
		}
	}

	// Build grouped changelog section
	var changelogSection strings.Builder
	var prBodySection strings.Builder

	changelogSection.WriteString(fmt.Sprintf("## %s - PREVIEW\n", version))

	for _, tag := range tagsList {
		entries, exists := tagMap[tag]
		if !exists || len(entries) == 0 {
			continue
		}
		cleanTagName := strings.TrimPrefix(tag, "#")
		heading := fmt.Sprintf("\n## %s\n\n", cleanTagName)

		changelogSection.WriteString(heading)
		prBodySection.WriteString(heading)

		for _, e := range entries {
			changelogSection.WriteString(e + "\n")
			prBodySection.WriteString(e + "\n")
		}
	}

	// Compose full new CHANGELOG.md content
	var newChangelog strings.Builder
	newChangelog.WriteString("# Changelog Chainlink Core\n\n")
	newChangelog.WriteString(changelogSection.String())
	newChangelog.WriteString("\n")
	newChangelog.WriteString(pastChangelog.String())

	// Compose PR body
	var prBody string
	prBodyCandidate := prHeader + prBodySection.String()
	if len(prBodyCandidate) > maxPRDescLength {
		prBody = prHeader + prTruncatedMsg
	} else {
		prBody = prBodyCandidate
	}

	result := &Result{
		Version:      version,
		PRBody:       prBody,
		NewChangelog: newChangelog.String(),
	}

	// Write updated changelog back to disk
	if err := os.WriteFile(resolvedChangelog, []byte(newChangelog.String()), 0o600); err != nil {
		return nil, fmt.Errorf("failed to write updated changelog %s: %w", resolvedChangelog, err)
	}

	// Write to GITHUB_OUTPUT if requested
	if writeGithubOutput {
		if err := githuboutput.AppendVar("version", version); err != nil {
			return nil, fmt.Errorf("failed to write version to GITHUB_OUTPUT: %w", err)
		}
		if err := githuboutput.AppendMultilineVar("pr_body", prBody); err != nil {
			return nil, fmt.Errorf("failed to write pr_body to GITHUB_OUTPUT: %w", err)
		}
	}

	return result, nil
}
