package devmachine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Bootstrap implements phase P2: install Docker + rsync, configure sshd and
// install the user's public key on the remote machine via SSM Run Command.
func (o *Orchestrator) Bootstrap(ctx context.Context) error {
	o.rep.Phase(phaseTitle(phaseBootstrap))

	if o.cfg.Machine.InstanceID == "" {
		return errors.New("no instance id available; run phase P1 first")
	}
	if o.run.DryRun {
		fmt.Fprintf(o.run.Out, "[dry-run] aws ssm send-command AWS-RunShellScript on %s\n", o.cfg.Machine.InstanceID)
		return nil
	}

	pubKey, err := o.publicKey()
	if err != nil {
		return err
	}

	script := bootstrapScript(pubKey)
	input := map[string]any{
		"InstanceIds":  []string{o.cfg.Machine.InstanceID},
		"DocumentName": "AWS-RunShellScript",
		"Parameters":   map[string][]string{"commands": strings.Split(script, "\n")},
	}
	raw, err := json.Marshal(input)
	if err != nil {
		return errors.Wrap(err, "failed to build SSM command input")
	}

	o.rep.Info("running remote bootstrap over SSM")
	args := []string{"ssm", "send-command", "--region", o.cfg.AWS.Region, "--cli-input-json", string(raw), "--output", "json"}
	args = append(args, o.profileArg()...)
	out, err := o.run.Run(ctx, "aws", args...)
	if err != nil {
		return errors.Wrap(err, "failed to send bootstrap command")
	}

	var resp struct {
		Command struct {
			CommandID string `json:"CommandId"`
		} `json:"Command"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		return errors.Wrap(err, "failed to parse SSM send-command response")
	}
	return o.waitForCommand(ctx, resp.Command.CommandID)
}

func (o *Orchestrator) waitForCommand(ctx context.Context, commandID string) error {
	if commandID == "" {
		return errors.New("SSM did not return a command id")
	}
	deadline := time.Now().Add(15 * time.Minute)
	for {
		args := []string{"ssm", "get-command-invocation",
			"--command-id", commandID,
			"--instance-id", o.cfg.Machine.InstanceID,
			"--region", o.cfg.AWS.Region,
			"--query", "{Status:Status,Out:StandardOutputContent,Err:StandardErrorContent}",
			"--output", "json",
		}
		args = append(args, o.profileArg()...)
		out, err := o.run.Run(ctx, "aws", args...)
		if err != nil {
			return errors.Wrap(err, "failed to poll bootstrap command")
		}
		var inv struct {
			Status string `json:"Status"`
			Out    string `json:"Out"`
			Err    string `json:"Err"`
		}
		if err := json.Unmarshal([]byte(out), &inv); err != nil {
			return errors.Wrap(err, "failed to parse command invocation")
		}
		switch inv.Status {
		case "Success":
			o.rep.Info("remote bootstrap complete")
			return nil
		case "Failed", "Cancelled", "TimedOut", "Undeliverable", "Terminated":
			return errors.Errorf("remote bootstrap %s: %s", inv.Status, strings.TrimSpace(inv.Err))
		}
		if time.Now().After(deadline) {
			return errors.New("timed out waiting for remote bootstrap to finish")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

func (o *Orchestrator) publicKey() (string, error) {
	keyPath := ExpandPath(o.cfg.SSH.Key)
	pubPath := keyPath + ".pub"
	if raw, err := os.ReadFile(pubPath); err == nil {
		return strings.TrimSpace(string(raw)), nil
	}
	// Derive the public key from the private key.
	out, err := o.run.Run(context.Background(), "ssh-keygen", "-y", "-f", keyPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read or derive public key for %s", keyPath)
	}
	return strings.TrimSpace(out), nil
}

func bootstrapScript(pubKey string) string {
	lines := []string{
		"set -e",
		"if ! command -v docker >/dev/null 2>&1; then dnf install -y docker; fi",
		"if ! command -v rsync >/dev/null 2>&1; then dnf install -y rsync; fi",
		"systemctl enable --now docker",
		"systemctl enable --now sshd",
		"install -d -m 700 -o ec2-user -g ec2-user /home/ec2-user/.ssh",
		fmt.Sprintf("grep -qxF %q /home/ec2-user/.ssh/authorized_keys 2>/dev/null || echo %q >> /home/ec2-user/.ssh/authorized_keys", pubKey, pubKey),
		"chown -R ec2-user:ec2-user /home/ec2-user/.ssh",
		"chmod 700 /home/ec2-user/.ssh",
		"chmod 600 /home/ec2-user/.ssh/authorized_keys",
		"restorecon -R /home/ec2-user/.ssh 2>/dev/null || true",
		"usermod -aG docker ec2-user || true",
		"usermod -aG docker ssm-user || true",
		"echo bootstrap-ok",
	}
	return strings.Join(lines, "\n")
}
