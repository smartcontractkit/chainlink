package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/tools"
)

func newToolsCmd() *cobra.Command {
	toolsCmd := &cobra.Command{
		Use:   "tools",
		Short: "Manage and discover tools in the tools directory",
	}

	toolsCmd.AddCommand(newToolsMatrixCmd())
	return toolsCmd
}

type toolsMatrixOptions struct {
	all          bool
	jsonOutput   bool
	eventName    string
	baseRef      string
	changedFiles string
}

func newToolsMatrixCmd() *cobra.Command {
	opts := &toolsMatrixOptions{}

	cmd := &cobra.Command{
		Use:   "matrix",
		Short: "Generate JSON matrix of tool test targets for CI workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runToolsMatrix(cmd.Context(), cmd, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.all, "all", false, "Include all tool targets regardless of git changes")
	cmd.Flags().BoolVar(&opts.jsonOutput, "json", false, "Output JSON to stdout")
	cmd.Flags().StringVar(&opts.eventName, "event-name", "", "GitHub event name (defaults to GITHUB_EVENT_NAME)")
	cmd.Flags().StringVar(&opts.baseRef, "base-ref", "", "Git base ref for diffing (defaults to GITHUB_BASE_REF or develop)")
	cmd.Flags().StringVar(&opts.changedFiles, "changed-files", "", "Comma- or newline-delimited list of changed files")

	return cmd
}

func runToolsMatrix(ctx context.Context, cmd *cobra.Command, opts *toolsMatrixOptions) error {
	repoRoot, err := FindRepoRoot(ctx)
	if err != nil {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return err
		}
		repoRoot = cwd
	}

	toolsDir := filepath.Join(repoRoot, "tools")
	allTargets, err := tools.DiscoverTargets(repoRoot, toolsDir)
	if err != nil {
		return fmt.Errorf("failed to discover tool targets: %w", err)
	}

	act := ghaction.New(cmd.OutOrStdout(), "", "")
	ghCtx, _ := act.Context()

	eventName := opts.eventName
	if eventName == "" && ghCtx != nil {
		eventName = ghCtx.EventName
	}
	if eventName == "" {
		eventName = act.Getenv("GITHUB_EVENT_NAME")
	}

	baseRef := opts.baseRef
	if baseRef == "" && ghCtx != nil {
		baseRef = ghCtx.BaseRef
	}
	if baseRef == "" {
		baseRef = act.Getenv("GITHUB_BASE_REF")
	}

	var changedFileList []string
	if opts.changedFiles != "" {
		for _, f := range strings.FieldsFunc(opts.changedFiles, func(r rune) bool {
			return r == ',' || r == '\n' || r == '\r'
		}) {
			f = strings.TrimSpace(f)
			if f != "" {
				changedFileList = append(changedFileList, f)
			}
		}
	} else if !opts.all && eventName != "schedule" && eventName != "workflow_dispatch" {
		var diffErr error
		changedFileList, diffErr = getGitChangedFiles(ctx, repoRoot, baseRef)
		if diffErr != nil {
			// Fall back to running all tool targets if diff cannot be resolved (e.g. shallow clone in CI)
			opts.all = true
		}
	}

	filteredTargets := tools.ComputeMatrix(allTargets, tools.MatrixOptions{
		EventName:    eventName,
		ChangedFiles: changedFileList,
		All:          opts.all,
	})
	if filteredTargets == nil {
		filteredTargets = []tools.Target{}
	}

	jsonData, err := json.Marshal(filteredTargets)
	if err != nil {
		return fmt.Errorf("failed to marshal targets to JSON: %w", err)
	}

	if opts.jsonOutput {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		if err := enc.Encode(filteredTargets); err != nil {
			return fmt.Errorf("failed to encode targets to JSON: %w", err)
		}
	} else if os.Getenv("GITHUB_OUTPUT") == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Discovered %d targets (%d selected):\n", len(allTargets), len(filteredTargets))
		for _, t := range filteredTargets {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s (dir: %s, packages: %s)\n", t.Name, t.Dir, t.Packages)
		}
	}

	if os.Getenv("GITHUB_OUTPUT") != "" {
		if err := act.SetOutput("matrix", string(jsonData)); err != nil {
			return fmt.Errorf("failed to set GitHub action output: %w", err)
		}
	}

	return nil
}

func getGitChangedFiles(ctx context.Context, repoRoot, baseRef string) ([]string, error) {
	if baseRef == "" {
		baseRef = os.Getenv("GITHUB_BASE_REF")
	}

	var candidates []string
	if baseRef != "" {
		candidates = append(candidates, baseRef)
		if !strings.HasPrefix(baseRef, "origin/") {
			candidates = append(candidates, "origin/"+baseRef)
		}
	}
	candidates = append(candidates, "origin/develop", "develop")

	var (
		out []byte
		err error
	)

	for _, ref := range candidates {
		// #nosec G204,G702
		diffCmd := exec.CommandContext(ctx, "git", "diff", "--name-only", ref+"...HEAD")
		diffCmd.Dir = repoRoot
		out, err = diffCmd.Output()
		if err == nil {
			break
		}
	}

	if err != nil {
		// Fallback to diff against HEAD~1
		// #nosec G204
		diffCmd := exec.CommandContext(ctx, "git", "diff", "--name-only", "HEAD~1")
		diffCmd.Dir = repoRoot
		out, err = diffCmd.Output()
		if err != nil {
			return nil, err
		}
	}

	var files []string
	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}
