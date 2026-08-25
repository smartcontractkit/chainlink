package changeset

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// AllowedTags is the list of valid changeset tags recognized by the release process.
var AllowedTags = []string{
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
}

// Result contains the result of changeset tag validation.
type Result struct {
	HasTags   bool     `json:"has_tags"`
	FoundTags []string `json:"found_tags"`
}

// CheckTags parses a changeset markdown file, validates the frontmatter semver value for chainlink,
// and searches for recognized release tags in the changeset content.
func CheckTags(filePath string) (Result, error) {
	if filePath == "" {
		return Result{}, errors.New("no changeset file path provided")
	}

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return Result{}, fmt.Errorf("file '%s' does not exist", filePath)
	}

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return Result{}, fmt.Errorf("failed to read changeset file: %w", err)
	}

	content := string(contentBytes)

	// Extract YAML frontmatter
	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return Result{}, fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	var meta map[string]string
	if err := yaml.Unmarshal([]byte(frontmatter), &meta); err != nil {
		return Result{}, fmt.Errorf("invalid changeset frontmatter YAML: %w", err)
	}

	semverVal, ok := meta["chainlink"]
	if !ok || (semverVal != "major" && semverVal != "minor" && semverVal != "patch") {
		return Result{}, errors.New("invalid changeset semvar value for 'chainlink'. Must be 'major', 'minor', or 'patch'")
	}

	// Scan lines for tags
	var foundTags []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		for _, tag := range AllowedTags {
			if strings.Contains(line, tag) {
				foundTags = append(foundTags, tag)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("failed scanning changeset content: %w", err)
	}

	hasTags := len(foundTags) > 0

	return Result{
		HasTags:   hasTags,
		FoundTags: foundTags,
	}, nil
}

func extractFrontmatter(content string) (string, error) {
	lines := strings.Split(content, "\n")
	firstIndex := -1
	secondIndex := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if firstIndex == -1 {
				firstIndex = i
				continue
			}
			secondIndex = i
			break
		}
	}

	if firstIndex == -1 || secondIndex == -1 || secondIndex <= firstIndex {
		return "", errors.New("frontmatter delimiters '---' not found")
	}

	return strings.Join(lines[firstIndex+1:secondIndex], "\n"), nil
}
