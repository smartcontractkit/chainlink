package ghaction

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/sethvargo/go-githubactions"
)

// Action handles interacting with GitHub Actions environment, outputs, and formatting.
type Action struct {
	*githubactions.Action
	out         io.Writer
	outputPath  string
	envPath     string
	summaryPath string
}

// New creates a new Action context. If outputPath or envPath are empty,
// it falls back to GITHUB_OUTPUT and GITHUB_ENV environment variables.
func New(out io.Writer, outputPath, envPath string) *Action {
	return NewWithOptions(out, outputPath, envPath, "")
}

// NewWithOptions creates a new Action context with explicit paths for output, env, and summary.
func NewWithOptions(out io.Writer, outputPath, envPath, summaryPath string) *Action {
	if out == nil {
		out = os.Stdout
	}
	if outputPath == "" {
		outputPath = os.Getenv("GITHUB_OUTPUT")
	}
	if envPath == "" {
		envPath = os.Getenv("GITHUB_ENV")
	}
	if summaryPath == "" {
		summaryPath = os.Getenv("GITHUB_STEP_SUMMARY")
	}

	envMap := map[string]string{
		"GITHUB_OUTPUT":       outputPath,
		"GITHUB_ENV":          envPath,
		"GITHUB_STEP_SUMMARY": summaryPath,
	}

	gha := githubactions.New(
		githubactions.WithWriter(out),
		githubactions.WithGetenv(func(k string) string {
			if v, ok := envMap[k]; ok {
				return v
			}
			return os.Getenv(k)
		}),
	)

	return &Action{
		Action:      gha,
		out:         out,
		outputPath:  outputPath,
		envPath:     envPath,
		summaryPath: summaryPath,
	}
}

// SetOutput writes a key-value pair to GITHUB_OUTPUT file, or falls back to out when unset.
func (a *Action) SetOutput(key, value string) error {
	if a.outputPath == "" {
		_, err := fmt.Fprintf(a.out, "%s=%s\n", key, value)
		return err
	}
	a.Action.SetOutput(key, value)
	return nil
}

// SetOutputs writes multiple key-value pairs to GITHUB_OUTPUT file, or falls back to out when unset.
func (a *Action) SetOutputs(outputs map[string]string) error {
	for k, v := range outputs {
		if err := a.SetOutput(k, v); err != nil {
			return err
		}
	}
	return nil
}

// SetEnv writes a key-value pair to GITHUB_ENV file, or falls back to out when unset.
func (a *Action) SetEnv(key, value string) error {
	if a.envPath == "" {
		_, err := fmt.Fprintf(a.out, "%s=%s\n", key, value)
		return err
	}
	a.Action.SetEnv(key, value)
	return nil
}

// AddStepSummary writes markdown content to GITHUB_STEP_SUMMARY file, or falls back to out when unset.
func (a *Action) AddStepSummary(markdown string) error {
	if a.summaryPath == "" {
		_, err := fmt.Fprintln(a.out, markdown)
		return err
	}
	a.Action.AddStepSummary(markdown)
	return nil
}

// AddStepSummaryTemplate parses and executes a template string, writing the result to GITHUB_STEP_SUMMARY or out.
func (a *Action) AddStepSummaryTemplate(tmpl string, data any) error {
	if a.summaryPath == "" {
		t, err := template.New("summary").Parse(tmpl)
		if err != nil {
			return fmt.Errorf("failed to parse summary template: %w", err)
		}
		if err := t.Execute(a.out, data); err != nil {
			return fmt.Errorf("failed to execute summary template: %w", err)
		}
		fmt.Fprintln(a.out)
		return nil
	}
	return a.Action.AddStepSummaryTemplate(tmpl, data)
}

// IsGitHubActions returns true if executing in a GitHub Actions runner environment.
func (a *Action) IsGitHubActions() bool {
	return a.Getenv("GITHUB_ACTIONS") == "true" || a.outputPath != ""
}

// Context returns the typed GitHubContext populated from GitHub Actions environment variables and event payload.
func (a *Action) Context() (*githubactions.GitHubContext, error) {
	return a.Action.Context()
}

// WithGroup executes fn wrapped inside a group command block.
func (a *Action) WithGroup(name string, fn func()) {
	a.Group(name)
	defer a.EndGroup()
	fn()
}
