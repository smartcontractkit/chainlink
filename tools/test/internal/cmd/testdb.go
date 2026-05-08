package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/db"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

var (
	setupTestDBShell string
	setupTestDBEval  bool
)

func inferShell() string {
	s := os.Getenv("SHELL")
	if s == "" {
		return "bash"
	}
	return strings.TrimPrefix(filepath.Base(s), "-")
}

var setupTestdbCmd = &cobra.Command{
	Use:   "setup-testdb",
	Short: "Start Postgres for tests and leave the container running",
	Long: `Starts the same Postgres image and preparetest flow as run/gotestsum/diagnose, but does not
tear the container down when this command exits. Ryuk is disabled so the database keeps running.

Prints CL_DATABASE_URL (one line) on stdout. Your interactive shell cannot inherit env vars from this
process; use --eval and eval the output, for example:

  eval "$(go -C tools/test run . setup-testdb --eval)"
  go -C tools/test run . setup-testdb --eval --shell fish | source

Removes containers started this way with: go -C tools/test run . remove-testdb`,
	RunE: runSetupTestdb,
}

func runSetupTestdb(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}
	out := output.NewFromApp(conf)
	sh := setupTestDBShell
	if sh == "" {
		sh = inferShell()
	}

	conn, err := db.StartPersistentPostgres(cmd.Context(), conf)
	if err != nil {
		return err
	}

	if setupTestDBEval {
		_, err = fmt.Print(db.ShellExportLine(conn, sh))
		return err
	}

	if out.AIOutput() {
		out.Stdoutln(conn)
		return nil
	}

	out.HumanStderr(
		termstyle.Label.Render("Persistent Postgres") + " " + termstyle.OK.Render("ready") + " " +
			termstyle.Muted.Render("(docker label "+db.PersistentTestDBLabel+")"))
	out.HumanStderr(termstyle.Muted.Render("Connection string (CL_DATABASE_URL):"))
	out.Stdoutln(conn)
	out.HumanStderr("")
	out.HumanStderr(termstyle.Muted.Render("Load into current shell:"))
	out.HumanStderr("  " + termstyle.Label.Render(`eval "$(go -C tools/test run . setup-testdb --eval)"`))
	if strings.EqualFold(sh, "fish") || strings.Contains(strings.ToLower(os.Getenv("SHELL")), "fish") {
		out.HumanStderr("  " + termstyle.Label.Render(`go -C tools/test run . setup-testdb --eval --shell fish | source`))
	}
	return nil
}

var removeTestdbCmd = &cobra.Command{
	Use:   "remove-testdb",
	Short: "Stop and remove persistent Postgres containers from setup-testdb",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := db.RemovePersistentTestDB(cmd.Context())
		if err != nil {
			return err
		}
		conf, loadErr := config.Load(cmd)
		if loadErr != nil {
			return loadErr
		}
		out := output.NewFromApp(conf)
		if out.AIOutput() {
			if len(ids) == 0 {
				out.Stdoutln("0")
				return nil
			}
			out.Stdoutln(strings.Join(ids, "\n"))
			return nil
		}
		switch len(ids) {
		case 0:
			out.HumanStderr(termstyle.Muted.Render("No persistent test-db containers found (" + db.PersistentTestDBLabel + ")."))
		default:
			out.HumanStderr(
				termstyle.Label.Render("Removed") + " " +
					termstyle.OK.Render(fmt.Sprintf("%d container(s): %s", len(ids), strings.Join(ids, " "))))
		}
		return nil
	},
}

func init() {
	setupTestdbCmd.Flags().BoolVar(&setupTestDBEval, "eval", false, "Print only a shell line to set CL_DATABASE_URL (for eval/source); use --shell")
	setupTestdbCmd.Flags().StringVar(&setupTestDBShell, "shell", "", "Shell for --eval: fish or bash/zsh/sh (default: infer from $SHELL)")

	rootCmd.AddCommand(setupTestdbCmd)
	rootCmd.AddCommand(removeTestdbCmd)
}
