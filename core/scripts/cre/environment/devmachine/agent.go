package devmachine

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// StartAgent implements phase P6: install the remote-agent as a systemd service
// and start it, so it keeps running after the SSH session ends (and restarts on
// failure / reboot).
func (o *Orchestrator) StartAgent(ctx context.Context) error {
	o.rep.Phase(phaseTitle(phaseStartAgent))

	if o.cfg.Machine.InstanceID == "" {
		return errors.New("no instance id available; run phase P1 first")
	}
	o.rep.Info("installing and starting the remote-agent systemd service")
	script := agentServiceScript(o.cfg)
	if o.run.DryRun {
		fmt.Fprintf(o.run.Out, "[dry-run] ssh %s '<systemd install script>'\n", o.cfg.SSH.Alias)
		return nil
	}
	return o.run.Stream(ctx, "ssh", "-o", "StrictHostKeyChecking=accept-new", o.cfg.SSH.Alias, script)
}

func agentServiceScript(cfg *Config) string {
	return fmt.Sprintf(`set -e
sudo tee /etc/systemd/system/remote-agent.service >/dev/null <<'EOF'
[Unit]
Description=CRE Remote Agent
After=network-online.target docker.service
Wants=network-online.target docker.service

[Service]
Type=simple
User=%s
Group=%s
SupplementaryGroups=docker
ExecStart=%s -addr %s
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
sudo systemctl daemon-reload
sudo systemctl enable remote-agent
sudo systemctl restart remote-agent
sleep 2
sudo systemctl --no-pager --full status remote-agent | head -n 12
`, DefaultRemoteUser, DefaultRemoteUser, cfg.Agent.RemotePath, cfg.Agent.Addr)
}
