package tunnel

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/rs/zerolog"
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
	localPort, err := reserveLocalPort()
	if err != nil {
		return TunnelBinding{}, fmt.Errorf("failed to reserve local port: %w", err)
	}

	cmd := exec.Command(
		"aws",
		"ssm",
		"start-session",
		"--region", p.region,
		"--target", p.instanceID,
		"--document-name", "AWS-StartPortForwardingSession",
		"--parameters", fmt.Sprintf("portNumber=%d,localPortNumber=%d", ref.Port, localPort),
	)
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
		Int("remotePort", ref.Port).
		Int("localPort", localPort).
		Msg("Opening SSM endpoint tunnel")

	if err := cmd.Start(); err != nil {
		return TunnelBinding{}, fmt.Errorf("failed to start aws ssm session: %w", err)
	}
	if err := waitForLocalPortReady(ctx, localPort, 12*time.Second); err != nil {
		_ = cmd.Process.Kill()
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
	}, nil
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

	if err := cmd.Process.Kill(); err != nil {
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
