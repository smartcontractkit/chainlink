package devmachine

import (
	"context"
	"fmt"
	"strings"
)

// Handoff implements phase P6: print the three manual commands the user runs to
// bring the environment up.
func (o *Orchestrator) Handoff(_ context.Context) error {
	o.rep.Phase(phaseTitle(phaseHandoff))
	fmt.Fprint(o.run.Out, handoffText(o.cfg))
	return nil
}

func handoffText(cfg *Config) string {
	var b strings.Builder
	b.WriteString("\nNext steps:\n\n")
	b.WriteString("1) The remote agent is already running on the machine as the systemd\n")
	b.WriteString("   service 'remote-agent' (survives logout and restarts on failure).\n")
	fmt.Fprintf(&b, "   Check it:  ssh %s 'systemctl status remote-agent'\n", cfg.SSH.Alias)
	b.WriteString("   Restart it: ssh " + cfg.SSH.Alias + " 'sudo systemctl restart remote-agent'\n\n")
	b.WriteString("2) On your laptop, start the environment:\n")
	fmt.Fprintf(&b, "   cd %s\n", cfg.Env.Workdir)
	fmt.Fprintf(&b, "   export CRE_REMOTE_AGENT_EC2_INSTANCE_ID=%s\n", cfg.Machine.InstanceID)
	if cfg.AWS.Profile != "" {
		fmt.Fprintf(&b, "   export AWS_PROFILE=%s\n", cfg.AWS.Profile)
	}
	if cfg.Machine.PrivateIP != "" {
		fmt.Fprintf(&b, "   # optional: export CRE_REMOTE_HOST_IP=%s\n", cfg.Machine.PrivateIP)
	}
	fmt.Fprintf(&b, "   CTF_CONFIGS=%s go run . env start\n", cfg.Env.Config)
	return b.String()
}
