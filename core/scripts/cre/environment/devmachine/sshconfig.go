package devmachine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	sshBlockStart = "# >>> cre-devmachine:%s >>>"
	sshBlockEnd   = "# <<< cre-devmachine:%s <<<"
)

func sshConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve home directory")
	}
	return filepath.Join(home, ".ssh", "config"), nil
}

// ensureSSHConfig writes or refreshes the managed Host block for the machine.
func (o *Orchestrator) ensureSSHConfig() error {
	path, err := sshConfigPath()
	if err != nil {
		return err
	}
	if o.run.DryRun {
		fmt.Fprintf(o.run.Out, "[dry-run] update managed Host block for %q in %s\n%s\n", o.cfg.SSH.Alias, path, o.sshBlock())
		return nil
	}

	var existing string
	if raw, err := os.ReadFile(path); err == nil {
		existing = string(raw)
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "failed to read %s", path)
	}

	cleaned := stripManagedBlock(existing, o.cfg.SSH.Alias)
	var b strings.Builder
	b.WriteString(strings.TrimRight(cleaned, "\n"))
	if b.Len() > 0 {
		b.WriteString("\n\n")
	}
	b.WriteString(o.sshBlock())
	b.WriteString("\n")

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return errors.Wrap(err, "failed to create ~/.ssh")
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		return errors.Wrapf(err, "failed to write %s", path)
	}
	return nil
}

func (o *Orchestrator) sshBlock() string {
	proxy := "aws ssm start-session --target %h --region " + o.cfg.AWS.Region
	if o.cfg.AWS.Profile != "" {
		proxy += " --profile " + o.cfg.AWS.Profile
	}
	proxy += " --document-name AWS-StartSSHSession --parameters portNumber=%p"

	return fmt.Sprintf(`%s
Host %s
  HostName %s
  User %s
  IdentityFile %s
  ProxyCommand %s
%s`,
		fmt.Sprintf(sshBlockStart, o.cfg.SSH.Alias),
		o.cfg.SSH.Alias,
		o.cfg.Machine.InstanceID,
		DefaultRemoteUser,
		o.cfg.SSH.Key,
		proxy,
		fmt.Sprintf(sshBlockEnd, o.cfg.SSH.Alias),
	)
}

// stripManagedBlock removes any previously written managed block for alias.
func stripManagedBlock(content, alias string) string {
	start := fmt.Sprintf(sshBlockStart, alias)
	end := fmt.Sprintf(sshBlockEnd, alias)
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	skipping := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == start:
			skipping = true
		case trimmed == end:
			skipping = false
		case !skipping:
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

// removeSSHConfig deletes the managed block for alias during teardown.
func removeSSHConfig(alias string) error {
	path, err := sshConfigPath()
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrapf(err, "failed to read %s", path)
	}
	cleaned := strings.TrimRight(stripManagedBlock(string(raw), alias), "\n")
	if cleaned == strings.TrimRight(string(raw), "\n") {
		return nil
	}
	if cleaned != "" {
		cleaned += "\n"
	}
	return os.WriteFile(path, []byte(cleaned), 0o600)
}
