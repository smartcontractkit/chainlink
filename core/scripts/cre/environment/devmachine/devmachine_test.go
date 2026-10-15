package devmachine

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeName(t *testing.T) {
	cases := map[string]string{
		"bartek.tofel":  "bartek-tofel",
		"Bartek Tofel":  "bartek-tofel",
		"a__b":          "a-b",
		"-lead-trail-":  "lead-trail",
		"UPPER":         "upper",
		"weird!@#chars": "weird-chars",
	}
	for in, want := range cases {
		if got := SanitizeName(in); got != want {
			t.Errorf("SanitizeName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMachineNameTruncates(t *testing.T) {
	long := strings.Repeat("x", 100)
	got := MachineName(long)
	if len(got) > 63 {
		t.Fatalf("MachineName length = %d, want <= 63", len(got))
	}
	if strings.HasSuffix(got, "-") {
		t.Fatalf("MachineName must not end with a dash: %q", got)
	}
}

func TestParsePhase(t *testing.T) {
	for _, in := range []string{"P0", "0", "preflight", "PHASE0"} {
		if got, err := parsePhase(in); err != nil || got != phasePreflight {
			t.Errorf("parsePhase(%q) = %v, %v", in, got, err)
		}
	}
	for _, in := range []string{"P6", "6", "start", "agent"} {
		if got, err := parsePhase(in); err != nil || got != phaseStartAgent {
			t.Errorf("parsePhase(%q) = %v, %v", in, got, err)
		}
	}
	for _, in := range []string{"P7", "7", "handoff"} {
		if got, err := parsePhase(in); err != nil || got != phaseHandoff {
			t.Errorf("parsePhase(%q) = %v, %v", in, got, err)
		}
	}
	if _, err := parsePhase("nope"); err == nil {
		t.Error("expected error for unknown phase")
	}
}

func TestShouldRun(t *testing.T) {
	if !shouldRun(phaseLoad, "", "") {
		t.Error("no selectors should run everything")
	}
	if !shouldRun(phaseLoad, "", phaseTransfer) {
		t.Error("--from transfer should run load")
	}
	if shouldRun(phaseBootstrap, "", phaseTransfer) {
		t.Error("--from transfer should skip bootstrap")
	}
	if !shouldRun(phaseLoad, phaseLoad, "") {
		t.Error("--only load should run load")
	}
	if shouldRun(phaseMachine, phaseLoad, "") {
		t.Error("--only load should skip machine")
	}
}

func TestConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "devmachine.yaml")
	cfg := DefaultConfig()
	cfg.Machine.Name = "hybrid-local-cre-test"
	cfg.Machine.InstanceID = "i-abc"
	cfg.AWS.Profile = "dev-ex-dev/administrator-access"

	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if got.Machine.Name != cfg.Machine.Name || got.Machine.InstanceID != "i-abc" {
		t.Fatalf("round trip mismatch: %+v", got.Machine)
	}
	if got.Machine.InstanceType != DefaultInstanceType {
		t.Errorf("default instance type = %q, want %q", got.Machine.InstanceType, DefaultInstanceType)
	}
	if got.AWS.Profile != cfg.AWS.Profile {
		t.Errorf("profile = %q, want %q", got.AWS.Profile, cfg.AWS.Profile)
	}
}

func TestLoadConfigMissingReturnsDefaults(t *testing.T) {
	cfg, err := LoadConfig(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Machine.InstanceType != DefaultInstanceType {
		t.Errorf("expected default instance type, got %q", cfg.Machine.InstanceType)
	}
	if cfg.AWS.Region != DefaultRegion {
		t.Errorf("expected default region, got %q", cfg.AWS.Region)
	}
}

func TestRenderManifest(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Machine.Name = "hybrid-local-cre-test"
	cfg.Machine.Team = "platform"
	cfg.Machine.CostCenter = "platform"
	m := renderManifest(cfg)

	for _, want := range []string{
		"kind: DevMachine",
		"name: hybrid-local-cre-test",
		"chain.link/team: platform",
		"chain.link/cost-center: platform",
		"instanceType: m8a.2xlarge",
		"volumeSize: 200",
		"availabilityZone: a",
		"ami: amazon-linux-2023",
	} {
		if !strings.Contains(m, want) {
			t.Errorf("manifest missing %q\n%s", want, m)
		}
	}
}

func TestStripManagedBlock(t *testing.T) {
	alias := "hybrid-local-cre-test"
	existing := "Host other\n  HostName x\n"
	block := &Orchestrator{cfg: &Config{
		SSH:     SSHConfig{Alias: alias, Key: "~/.ssh/id_ed25519"},
		Machine: MachineConfig{InstanceID: "i-abc"},
		AWS:     AWSConfig{Region: "us-west-2", Profile: "p"},
	}}
	content := existing + "\n" + block.sshBlock() + "\n"

	// Strip once and ensure it is gone.
	cleaned := stripManagedBlock(content, alias)
	if strings.Contains(cleaned, alias) {
		t.Fatalf("managed block not removed:\n%s", cleaned)
	}
	// Stripping an already-clean file is a no-op.
	if again := stripManagedBlock(cleaned, alias); again != cleaned {
		t.Fatalf("strip is not idempotent:\n%q\nvs\n%q", cleaned, again)
	}
	// Other hosts are preserved.
	if !strings.Contains(cleaned, "Host other") {
		t.Fatal("unrelated host entry was removed")
	}
}

func TestSSHBlockContainsProxyCommand(t *testing.T) {
	o := &Orchestrator{cfg: &Config{
		SSH:     SSHConfig{Alias: "alias", Key: "~/.ssh/id_ed25519"},
		Machine: MachineConfig{InstanceID: "i-abc"},
		AWS:     AWSConfig{Region: "us-west-2", Profile: "dev-ex-dev/administrator-access"},
	}}
	block := o.sshBlock()
	for _, want := range []string{
		"Host alias",
		"HostName i-abc",
		"AWS-StartSSHSession",
		"--profile dev-ex-dev/administrator-access",
		"portNumber=%p",
	} {
		if !strings.Contains(block, want) {
			t.Errorf("ssh block missing %q\n%s", want, block)
		}
	}
}

func TestTarName(t *testing.T) {
	cases := map[string]string{
		"chainlink-amd:latest":       "chainlink-amd-latest.tar",
		"job-distributor-amd:0.22.1": "job-distributor-amd-0.22.1.tar",
		"registry.io/foo/bar:1.0":    "registry.io_foo_bar-1.0.tar",
	}
	for in, want := range cases {
		if got := tarName(in); got != want {
			t.Errorf("tarName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBootstrapScript(t *testing.T) {
	script := bootstrapScript("ssh-ed25519 AAAA test@host")
	for _, want := range []string{
		"dnf install -y docker",
		"dnf install -y rsync",
		"restorecon -R /home/ec2-user/.ssh",
		"usermod -aG docker ec2-user",
		"ssh-ed25519 AAAA test@host",
	} {
		if !strings.Contains(script, want) {
			t.Errorf("bootstrap script missing %q", want)
		}
	}
}

func TestAgentServiceScript(t *testing.T) {
	cfg := DefaultConfig()
	script := agentServiceScript(cfg)
	for _, want := range []string{
		"/etc/systemd/system/remote-agent.service",
		"ExecStart=/home/ec2-user/remote-agent -addr 0.0.0.0:18080",
		"SupplementaryGroups=docker",
		"WantedBy=multi-user.target",
		"systemctl enable remote-agent",
		"systemctl restart remote-agent",
	} {
		if !strings.Contains(script, want) {
			t.Errorf("agent service script missing %q", want)
		}
	}
}

func TestRunnerDryRun(t *testing.T) {
	var out bytes.Buffer
	r := NewRunner(&out, &out, true, false)
	got, err := r.Run(context.Background(), "echo", "hello world")
	if err != nil {
		t.Fatalf("dry-run Run returned error: %v", err)
	}
	if got != "" {
		t.Errorf("dry-run should return empty output, got %q", got)
	}
	if !strings.Contains(out.String(), "[dry-run] echo") {
		t.Errorf("dry-run should log the command, got %q", out.String())
	}
}

func TestRunnerIgnoresStderr(t *testing.T) {
	var out bytes.Buffer
	r := NewRunner(&out, &out, false, false)
	// A command that prints a warning to stderr and "yes" to stdout must return
	// only the stdout content.
	got, err := r.Run(context.Background(), "sh", "-c", "echo 'warning: noisy' 1>&2; echo yes")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if got != "yes" {
		t.Errorf("Run returned %q, want %q", got, "yes")
	}
}

func TestOutputSaysYes(t *testing.T) {
	if !outputSaysYes("yes") {
		t.Error("expected true for 'yes'")
	}
	if !outputSaysYes("Warning: something\nyes\n") {
		t.Error("expected true when a warning precedes 'yes'")
	}
	if outputSaysYes("no") {
		t.Error("expected false for 'no'")
	}
	if outputSaysYes("yesterday") {
		t.Error("expected false for 'yesterday'")
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home directory")
	}
	if got := ExpandPath("~/x"); got != filepath.Join(home, "x") {
		t.Errorf("ExpandPath(~/x) = %q", got)
	}
	if got := ExpandPath("/abs"); got != "/abs" {
		t.Errorf("ExpandPath(/abs) = %q", got)
	}
}

// writeFakeBinaries installs stub executables on PATH for the duration of the
// test and returns the bin directory.
func writeFakeBinaries(t *testing.T, scripts map[string]string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range scripts {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o755); err != nil { //nolint:gosec // executable test stub
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return dir
}

func TestMachineRecreatesWhenRecordedInstanceIsTerminated(t *testing.T) {
	writeFakeBinaries(t, map[string]string{
		"aws": "#!/bin/sh\necho 'terminated 10.0.0.1'\n",
		"kubectl": "#!/bin/sh\nfor a in \"$@\"; do\n  if [ \"$a\" = \"get\" ]; then\n" +
			"    echo '{\"status\":{\"instanceId\":\"i-new\",\"privateIp\":\"10.0.0.2\",\"conditions\":[{\"type\":\"Ready\",\"status\":\"True\"}]}}'\n" +
			"    exit 0\n  fi\ndone\nexit 0\n",
	})

	cfg := DefaultConfig()
	cfg.Cluster.Privileged = true
	cfg.Machine.Name = "hybrid-local-cre-test"
	cfg.Machine.InstanceID = "i-stale"
	cfg.Machine.PrivateIP = "10.0.0.1"
	cfg.Machine.CreatedByTool = true

	var out bytes.Buffer
	run := NewRunner(&out, &out, false, false)
	o := NewOrchestrator(cfg, &UpOptions{ConfigPath: filepath.Join(t.TempDir(), "cfg.yaml")}, run, NewPlainReporter(&out), NewPrompter(&out))

	if err := o.Machine(context.Background()); err != nil {
		t.Fatalf("Machine: %v", err)
	}
	if cfg.Machine.InstanceID != "i-new" {
		t.Fatalf("expected a new instance id, got %q", cfg.Machine.InstanceID)
	}
	if !cfg.Machine.CreatedByTool {
		t.Fatal("expected CreatedByTool to be true after recreation")
	}
}

func TestMachineHonorsExplicitInstanceID(t *testing.T) {
	writeFakeBinaries(t, map[string]string{
		"aws": "#!/bin/sh\necho 'terminated 10.0.0.1'\n",
	})

	cfg := DefaultConfig()
	cfg.Cluster.Privileged = true
	cfg.Machine.Name = "hybrid-local-cre-test"

	var out bytes.Buffer
	run := NewRunner(&out, &out, false, false)
	o := NewOrchestrator(cfg, &UpOptions{ConfigPath: filepath.Join(t.TempDir(), "cfg.yaml"), InstanceID: "i-explicit"}, run, NewPlainReporter(&out), NewPrompter(&out))

	err := o.Machine(context.Background())
	if err == nil {
		t.Fatal("expected an error when an explicitly supplied instance is unusable")
	}
	if !strings.Contains(err.Error(), "terminated") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDownClearsMachineState(t *testing.T) {
	writeFakeBinaries(t, map[string]string{
		"kubectl": "#!/bin/sh\nexit 0\n",
	})

	cfgPath := filepath.Join(t.TempDir(), "cfg.yaml")
	cfg := DefaultConfig()
	cfg.Machine.Name = "hybrid-local-cre-test"
	cfg.Machine.InstanceID = "i-abc"
	cfg.Machine.PrivateIP = "10.0.0.1"
	cfg.Machine.CreatedByTool = true
	// Leave SSH.Alias empty so the test never touches the real ~/.ssh/config.

	var out bytes.Buffer
	run := NewRunner(&out, &out, false, false)
	rep := NewPlainReporter(&out)
	if err := Down(context.Background(), cfg, &DownOptions{ConfigPath: cfgPath}, run, rep); err != nil {
		t.Fatalf("Down: %v", err)
	}
	if cfg.Machine.InstanceID != "" || cfg.Machine.PrivateIP != "" || cfg.Machine.CreatedByTool {
		t.Fatalf("machine state not cleared: %+v", cfg.Machine)
	}
	reloaded, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Machine.InstanceID != "" {
		t.Fatalf("persisted instanceId not cleared: %q", reloaded.Machine.InstanceID)
	}
}
