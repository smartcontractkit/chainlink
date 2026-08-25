package ghaction

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// Action handles writing to GitHub Actions output files and formatting workflow commands.
type Action struct {
	out        io.Writer
	outputPath string
	envPath    string
}

// New creates a new Action context. If outputPath or envPath are empty,
// it falls back to GITHUB_OUTPUT and GITHUB_ENV environment variables.
func New(out io.Writer, outputPath, envPath string) *Action {
	if out == nil {
		out = os.Stdout
	}
	if outputPath == "" {
		outputPath = os.Getenv("GITHUB_OUTPUT")
	}
	if envPath == "" {
		envPath = os.Getenv("GITHUB_ENV")
	}
	return &Action{
		out:        out,
		outputPath: outputPath,
		envPath:    envPath,
	}
}

// SetOutput writes a key-value pair to GITHUB_OUTPUT file, or fallback to out.
func (a *Action) SetOutput(key, value string) error {
	return a.writeFileOrFallback(a.outputPath, key, value)
}

// SetEnv writes a key-value pair to GITHUB_ENV file, or fallback to out.
func (a *Action) SetEnv(key, value string) error {
	return a.writeFileOrFallback(a.envPath, key, value)
}

// Errorf writes a ::error:: workflow annotation to out.
func (a *Action) Errorf(format string, args ...any) {
	fmt.Fprintf(a.out, "::error::%s\n", fmt.Sprintf(format, args...))
}

// Warningf writes a ::warning:: workflow annotation to out.
func (a *Action) Warningf(format string, args ...any) {
	fmt.Fprintf(a.out, "::warning::%s\n", fmt.Sprintf(format, args...))
}

// Group writes a ::group:: workflow command to out.
func (a *Action) Group(name string) {
	fmt.Fprintf(a.out, "::group::%s\n", name)
}

// EndGroup writes a ::endgroup:: workflow command to out.
func (a *Action) EndGroup() {
	fmt.Fprintln(a.out, "::endgroup::")
}

func (a *Action) writeFileOrFallback(filePath, key, value string) error {
	if filePath == "" {
		// Local / fallback mode: print key=value or formatted string to out
		if strings.Contains(value, "\n") {
			delimiter := generateDelimiter()
			fmt.Fprintf(a.out, "%s<<%s\n%s\n%s\n", key, delimiter, value, delimiter)
		} else {
			fmt.Fprintf(a.out, "%s=%s\n", key, value)
		}
		return nil
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer f.Close()

	if strings.Contains(value, "\n") {
		delimiter := generateDelimiter()
		_, err = fmt.Fprintf(f, "%s<<%s\n%s\n%s\n", key, delimiter, value, delimiter)
	} else {
		_, err = fmt.Fprintf(f, "%s=%s\n", key, value)
	}
	if err != nil {
		return fmt.Errorf("failed to write to file %q: %w", filePath, err)
	}
	return nil
}

func generateDelimiter() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "ghadelimiter_default"
	}
	return "ghadelimiter_" + hex.EncodeToString(b)
}
