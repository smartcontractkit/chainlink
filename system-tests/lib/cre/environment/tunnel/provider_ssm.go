package tunnel

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

type SSMProvider struct {
	instanceID string
	region     string
	logger     zerolog.Logger

	mu       sync.Mutex
	sessions map[int]*exec.Cmd
}

func NewSSMProvider(instanceID, region string, logger zerolog.Logger) Provider {
	return &SSMProvider{
		instanceID: instanceID,
		region:     region,
		logger:     logger,
		sessions:   make(map[int]*exec.Cmd),
	}
}

func (p *SSMProvider) Name() string {
	return "ssm"
}

func (p *SSMProvider) Open(ctx context.Context, ref EndpointRef) (TunnelBinding, error) {
	profile, authMode := runtimecfg.ResolveAWSCLIProfileSelection()
	if err := validateAWSSession(ctx, p.region, profile, authMode); err != nil {
		return TunnelBinding{}, err
	}

	localPort, err := reserveLocalPort()
	if err != nil {
		return TunnelBinding{}, fmt.Errorf("failed to reserve local port: %w", err)
	}

	args := []string{
		"ssm",
		"start-session",
		"--region", p.region,
		"--target", p.instanceID,
		"--document-name", "AWS-StartPortForwardingSession",
		"--parameters", fmt.Sprintf("portNumber=%d,localPortNumber=%d", ref.Port, localPort),
	}
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	cmd := exec.Command("aws", args...)
	// Start in a dedicated process group so cleanup can kill aws + session-manager-plugin together.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if p.logger.GetLevel() <= zerolog.DebugLevel {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		p.logger.Debug().
			Strs("cmd", cmd.Args).
			Msg("Starting SSM endpoint tunnel command")
	}

	p.logger.Info().
		Str("componentID", ref.ComponentID).
		Str("endpointName", ref.EndpointName).
		Str("awsAuthMode", authMode).
		Str("awsProfile", profile).
		Int("remotePort", ref.Port).
		Int("localPort", localPort).
		Msg("Opening SSM endpoint tunnel")

	if err := cmd.Start(); err != nil {
		return TunnelBinding{}, fmt.Errorf("failed to start aws ssm session: %w", err)
	}
	if err := waitForLocalPortReady(ctx, localPort, 12*time.Second); err != nil {
		terminateProcessGroup(cmd)
		return TunnelBinding{}, fmt.Errorf("ssm local tunnel on port %d did not become ready: %w", localPort, err)
	}

	p.mu.Lock()
	p.sessions[localPort] = cmd
	p.mu.Unlock()

	go func() {
		_ = cmd.Wait()
	}()

	return TunnelBinding{
		EndpointRef: ref,
		LocalPort:   localPort,
		LocalURL:    localURLFromScheme(ref.Scheme, localPort),
		PID:         cmd.Process.Pid,
	}, nil
}

func validateAWSSession(ctx context.Context, region, profile, authMode string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	preflightCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	args := []string{"sts", "get-caller-identity", "--region", region}
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	out, err := exec.CommandContext(preflightCtx, "aws", args...).CombinedOutput()
	if err == nil {
		return nil
	}

	loginHint := "Verify AWS credentials are configured and valid."
	if profile != "" {
		loginHint = fmt.Sprintf("Run `aws sso login --profile %s` (or configure profile credentials) and retry.", profile)
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return fmt.Errorf("aws authentication check failed for SSM tunnel (mode=%s): %w. %s", authMode, err, loginHint)
	}
	return fmt.Errorf("aws authentication check failed for SSM tunnel (mode=%s): %w: %s. %s", authMode, err, trimmed, loginHint)
}

func (p *SSMProvider) Close(_ context.Context, binding TunnelBinding) error {
	p.mu.Lock()
	cmd, ok := p.sessions[binding.LocalPort]
	if ok {
		delete(p.sessions, binding.LocalPort)
	}
	p.mu.Unlock()

	if !ok || cmd == nil || cmd.Process == nil {
		return nil
	}

	if err := terminateProcessGroup(cmd); err != nil {
		return fmt.Errorf("failed to kill ssm session on local port %d: %w", binding.LocalPort, err)
	}
	p.logger.Info().
		Str("componentID", binding.ComponentID).
		Str("endpointName", binding.EndpointName).
		Int("localPort", binding.LocalPort).
		Msg("Closed SSM endpoint tunnel")
	return nil
}

func reserveLocalPort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("listener addr %T is not tcp", l.Addr())
	}
	return tcpAddr.Port, nil
}

func localURLFromScheme(scheme string, port int) string {
	switch scheme {
	case "ws":
		return fmt.Sprintf("ws://127.0.0.1:%d", port)
	case "wss":
		return fmt.Sprintf("wss://127.0.0.1:%d", port)
	case "https":
		return fmt.Sprintf("https://127.0.0.1:%d", port)
	default:
		return fmt.Sprintf("http://127.0.0.1:%d", port)
	}
}

func terminateProcessGroup(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	// Negative PID targets the process group when Setpgid=true.
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		// Fall back to killing parent process only.
		if killErr := cmd.Process.Kill(); killErr != nil {
			return killErr
		}
	}
	return nil
}

func waitForLocalPortReady(ctx context.Context, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	address := fmt.Sprintf("127.0.0.1:%d", port)
	var lastErr error

	for time.Now().Before(deadline) {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		lastErr = err
		time.Sleep(200 * time.Millisecond)
	}

	if lastErr == nil {
		lastErr = errors.New("unknown readiness failure")
	}
	return lastErr
}
