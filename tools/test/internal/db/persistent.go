package db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

// PersistentTestDBLabel is the Docker label key used for Postgres containers
// started by setup-testdb so remove-testdb can find them.
const PersistentTestDBLabel = "chainlink.tools.test.persistent-testdb"

const persistentTestDBLabelValue = "1"

// PersistentTestDBLabelFilter is the argument to `docker ps -f …`.
func PersistentTestDBLabelFilter() string {
	return fmt.Sprintf("label=%s=%s", PersistentTestDBLabel, persistentTestDBLabelValue)
}

// StartPersistentPostgres starts Postgres (same image/options as Ensure), runs
// preparetest, and returns the connection string. The caller must not call
// Terminate on the container if it should keep running. Ryuk is disabled so
// the container survives after this process exits.
func StartPersistentPostgres(ctx context.Context, conf *config.App) (connStr string, err error) {
	if conf.PostgresVersion == "" {
		return "", errors.New("postgres version is required")
	}
	if conf.DatabaseURL != "" {
		return "", errors.New("--database-url is not compatible with setup-testdb (omit it to start a container)")
	}
	if err = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return "", fmt.Errorf("set TESTCONTAINERS_RYUK_DISABLED: %w", err)
	}

	c, err := postgres.Run(ctx,
		fmt.Sprintf("docker.io/postgres:%s-alpine", conf.PostgresVersion),
		postgres.WithDatabase("chainlink_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithCmdArgs("-c", "max_connections=1000"),
		testcontainers.WithLabels(map[string]string{
			PersistentTestDBLabel: persistentTestDBLabelValue,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return "", fmt.Errorf("postgres testcontainer: %w", err)
	}
	defer func() {
		if err != nil {
			termCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_ = c.Terminate(termCtx)
		}
	}()

	connStr, err = c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", fmt.Errorf("connection string: %w", err)
	}

	prepareOutput := bytes.NewBuffer(nil)
	prep := exec.CommandContext(ctx, "go", "run", "./core/store/cmd/preparetest", "--force")
	prep.Dir = conf.RepoRoot
	prep.Env = append(os.Environ(), "CL_DATABASE_URL="+connStr)
	prep.Stdout = prepareOutput
	prep.Stderr = prepareOutput
	if err = prep.Run(); err != nil {
		return "", fmt.Errorf("preparetest --force: %w\n%s", err, prepareOutput.String())
	}

	return connStr, nil
}

// dockerContainerIDPattern matches full or short Docker container IDs from `docker ps -q`.
var dockerContainerIDPattern = regexp.MustCompile(`^[a-f0-9]{12,64}$`)

// RemovePersistentTestDB stops and removes all containers labeled as a
// persistent tools/test Postgres instance.
func RemovePersistentTestDB(ctx context.Context) ([]string, error) {
	//nolint:gosec // G204: docker filter is built from package constants only.
	cmd := exec.CommandContext(ctx, "docker", "ps", "-aq", "-f", PersistentTestDBLabelFilter())
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("docker ps: %w\n%s", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("docker ps: %w", err)
	}
	raw := strings.Fields(strings.TrimSpace(string(out)))
	ids := make([]string, 0, len(raw))
	for _, id := range raw {
		if !dockerContainerIDPattern.MatchString(id) {
			return nil, fmt.Errorf("invalid container id from docker ps: %q", id)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, nil
	}
	//nolint:gosec // G204: each id matched dockerContainerIDPattern before use.
	rm := exec.CommandContext(ctx, "docker", append([]string{"rm", "-f"}, ids...)...)
	rmOut, err := rm.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker rm -f: %w\n%s", err, string(rmOut))
	}
	return ids, nil
}

// ShellExportLine returns a single line suitable for `eval` in the given shell
// to set CL_DATABASE_URL (fish vs posix).
func ShellExportLine(connURL, shell string) string {
	switch strings.ToLower(strings.TrimSpace(shell)) {
	case "fish":
		return fmt.Sprintf("set -gx CL_DATABASE_URL %s\n", fishDoubleQuote(connURL))
	default:
		return fmt.Sprintf("export CL_DATABASE_URL=%s\n", posixSingleQuote(connURL))
	}
}

func posixSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, `'`, `'\''`) + "'"
}

func fishDoubleQuote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		if r == '\\' || r == '"' {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	b.WriteByte('"')
	return b.String()
}
