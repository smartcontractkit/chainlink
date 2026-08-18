package githuboutput

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnvFilePath returns the path set in $GITHUB_OUTPUT, or "" when unset.
func EnvFilePath() string {
	return os.Getenv("GITHUB_OUTPUT")
}

// AppendVar appends a single-line k=v pair to $GITHUB_OUTPUT. No-op when $GITHUB_OUTPUT is unset.
func AppendVar(key, value string) error {
	file := EnvFilePath()
	if file == "" {
		return nil
	}
	return AppendToFile(file, fmt.Sprintf("%s=%s\n", key, value))
}

// AppendMultilineVar appends a delimited (heredoc-style) variable to $GITHUB_OUTPUT.
// No-op when $GITHUB_OUTPUT is unset.
func AppendMultilineVar(key, value string) error {
	file := EnvFilePath()
	if file == "" {
		return nil
	}
	return AppendToFile(file, fmt.Sprintf("%s<<EOF\n%s\nEOF\n", key, value))
}

// AppendVars appends each k=v pair to $GITHUB_OUTPUT. No-op when $GITHUB_OUTPUT is unset.
func AppendVars(vars map[string]string) error {
	for key, value := range vars {
		if err := AppendVar(key, value); err != nil {
			return err
		}
	}
	return nil
}

// AppendToFile appends content to path, creating the file if needed.
func AppendToFile(path, content string) error {
	f, err := os.OpenFile(filepath.Clean(path), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}
