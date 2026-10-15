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
	b.WriteString("1) SSH to the machine:\n")
	fmt.Fprintf(&b, "   ssh %s\n\n", cfg.SSH.Alias)
	b.WriteString("2) Start the remote agent on the machine:\n")
	fmt.Fprintf(&b, "   %s -addr %s\n\n", cfg.Agent.RemotePath, cfg.Agent.Addr)
	b.WriteString("3) Back on your laptop, start the environment:\n")
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
