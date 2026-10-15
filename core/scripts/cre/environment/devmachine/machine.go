package devmachine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Machine implements phase P1: either create a DevMachine via the CRD and wait
// for it to become Ready, or adopt an already existing instance.
func (o *Orchestrator) Machine(ctx context.Context) error {
	o.rep.Phase(phaseTitle(phaseMachine))

	if o.opts.InstanceType != "" {
		o.cfg.Machine.InstanceType = o.opts.InstanceType
	}
	if o.opts.VolumeSize != 0 {
		o.cfg.Machine.VolumeSize = o.opts.VolumeSize
	}
	if o.opts.AZ != "" {
		o.cfg.Machine.AvailabilityZ = o.opts.AZ
	}

	// Non-privileged path (or an explicitly supplied instance id): adopt.
	if o.opts.InstanceID != "" {
		o.cfg.Machine.InstanceID = o.opts.InstanceID
	}
	if o.cfg.Machine.InstanceID != "" {
		return o.adoptMachine(ctx)
	}
	if !o.cfg.Cluster.Privileged {
		return errors.New("no permission to create DevMachines and no --instance-id supplied; ask someone to create one and re-run with --instance-id")
	}

	if err := o.createMachine(ctx); err != nil {
		return err
	}
	o.cfg.Machine.CreatedByTool = true
	return o.save()
}

func (o *Orchestrator) createMachine(ctx context.Context) error {
	manifest := renderManifest(o.cfg)
	if o.run.DryRun {
		fmt.Fprintf(o.run.Out, "[dry-run] kubectl apply -f - (manifest below)\n%s\n", manifest)
		return nil
	}

	o.rep.Info("applying DevMachine %q in namespace %q", o.cfg.Machine.Name, o.cfg.Cluster.Namespace)
	args := []string{"apply", "-n", o.cfg.Cluster.Namespace, "-f", "-"}
	if err := o.run.StreamStdin(ctx, manifest, "kubectl", args...); err != nil {
		return errors.Wrap(err, "failed to apply DevMachine manifest")
	}

	return o.waitForReady(ctx)
}

func (o *Orchestrator) waitForReady(ctx context.Context) error {
	o.rep.Info("waiting for DevMachine to become Ready (this can take a few minutes)")
	deadline := time.Now().Add(10 * time.Minute)
	for {
		st, err := o.readMachineStatus(ctx)
		if err != nil {
			return err
		}
		if st.ready() && st.InstanceID != "" {
			o.cfg.Machine.InstanceID = st.InstanceID
			o.cfg.Machine.PrivateIP = st.PrivateIP
			o.rep.Info("DevMachine ready: instance=%s privateIP=%s", st.InstanceID, st.PrivateIP)
			return nil
		}
		if time.Now().After(deadline) {
			return errors.Errorf("timed out waiting for DevMachine %q to become Ready", o.cfg.Machine.Name)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
		}
	}
}

type machineStatus struct {
	InstanceID string
	PrivateIP  string
	Conditions []condition
}

type condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

func (m machineStatus) ready() bool {
	for _, c := range m.Conditions {
		if c.Type == "Ready" && c.Status == "True" {
			return true
		}
	}
	return false
}

func (o *Orchestrator) readMachineStatus(ctx context.Context) (machineStatus, error) {
	out, err := o.run.Run(ctx, "kubectl", "-n", o.cfg.Cluster.Namespace, "get", "devmachine", o.cfg.Machine.Name, "-o", "json")
	if err != nil {
		return machineStatus{}, errors.Wrap(err, "failed to read DevMachine status")
	}
	var raw struct {
		Status struct {
			InstanceID string      `json:"instanceId"`
			PrivateIP  string      `json:"privateIp"`
			Conditions []condition `json:"conditions"`
		} `json:"status"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return machineStatus{}, errors.Wrap(err, "failed to parse DevMachine status")
	}
	return machineStatus{
		InstanceID: raw.Status.InstanceID,
		PrivateIP:  raw.Status.PrivateIP,
		Conditions: raw.Status.Conditions,
	}, nil
}

// adoptMachine verifies an existing instance and records its private IP.
func (o *Orchestrator) adoptMachine(ctx context.Context) error {
	o.rep.Info("using existing instance %s", o.cfg.Machine.InstanceID)
	if o.run.DryRun {
		return nil
	}
	args := []string{"ec2", "describe-instances",
		"--instance-ids", o.cfg.Machine.InstanceID,
		"--region", o.cfg.AWS.Region,
		"--query", "Reservations[0].Instances[0].[State.Name,PrivateIpAddress]",
		"--output", "text",
	}
	args = append(args, o.profileArg()...)
	out, err := o.run.Run(ctx, "aws", args...)
	if err != nil {
		return errors.Wrapf(err, "failed to describe instance %s", o.cfg.Machine.InstanceID)
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return errors.Errorf("instance %s not found in region %s", o.cfg.Machine.InstanceID, o.cfg.AWS.Region)
	}
	if fields[0] != "running" {
		return errors.Errorf("instance %s is %s, expected running", o.cfg.Machine.InstanceID, fields[0])
	}
	if len(fields) > 1 && fields[1] != "None" {
		o.cfg.Machine.PrivateIP = fields[1]
	}
	o.rep.Info("instance running: privateIP=%s", o.cfg.Machine.PrivateIP)
	return o.save()
}

func (o *Orchestrator) profileArg() []string {
	if o.cfg.AWS.Profile == "" {
		return nil
	}
	return []string{"--profile", o.cfg.AWS.Profile}
}

func renderManifest(cfg *Config) string {
	return fmt.Sprintf(`apiVersion: aws.platform.smartcontract.com/v1alpha1
kind: DevMachine
metadata:
  name: %s
  namespace: %s
  labels:
    chain.link/env: dev
    chain.link/product: cre
    chain.link/team: %s
    chain.link/cost-center: %s
spec:
  availabilityZone: %s
  ami: %s
  instanceType: %s
  volumeSize: %d
`, cfg.Machine.Name, cfg.Cluster.Namespace, cfg.Machine.Team, cfg.Machine.CostCenter,
		cfg.Machine.AvailabilityZ, cfg.Machine.AMI, cfg.Machine.InstanceType, cfg.Machine.VolumeSize)
}
