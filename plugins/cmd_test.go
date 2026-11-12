package plugins

import (
	"errors"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
)

func TestNewCmdFactory_RegisterSuccess(t *testing.T) {
	mockRegister := func(id string) (*RegisteredLoop, error) {
		return &RegisteredLoop{EnvCfg: loop.EnvConfig{}}, nil
	}

	cmdConfig := CmdConfig{
		ID:  "test-loop",
		Cmd: "echo",
		Env: []string{"TEST_ENV=1"},
	}

	cmdFactory, err := NewCmdFactory(mockRegister, cmdConfig)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	cmd := cmdFactory()
	if cmd.Args[0] != "echo" {
		t.Errorf("Expected command 'echo', got %s", cmd.Args[0])
	}
}

func TestNewCmdFactory_TelemetryLoopID(t *testing.T) {
	const envKey = "CL_TELEMETRY_LOOP_ID="
	envOf := func(t *testing.T, enabled bool, id string) []string {
		t.Helper()
		registered := &RegisteredLoop{EnvCfg: loop.EnvConfig{}, TelemetryLoopIDEnabled: enabled}
		cmdFactory, err := NewCmdFactory(func(string) (*RegisteredLoop, error) { return registered, nil },
			CmdConfig{ID: "test-loop", Cmd: "echo", TelemetryLoopID: id})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		env := cmdFactory().Env
		if registered.EnvCfg.TelemetryLoopID != "" {
			t.Errorf("registry entry must not be mutated, got %q", registered.EnvCfg.TelemetryLoopID)
		}
		return env
	}
	hasID := func(env []string) (string, bool) {
		for _, kv := range env {
			if v, ok := strings.CutPrefix(kv, envKey); ok {
				return v, true
			}
		}
		return "", false
	}

	if v, ok := hasID(envOf(t, true, "Relayer.evm.1")); !ok || v != "Relayer.evm.1" {
		t.Errorf("enabled: expected %sRelayer.evm.1, got %q (present=%v)", envKey, v, ok)
	}
	if _, ok := hasID(envOf(t, false, "Relayer.evm.1")); ok {
		t.Errorf("disabled: %s must not be sent", envKey)
	}
	if _, ok := hasID(envOf(t, true, "")); ok {
		t.Errorf("enabled without an ID: %s must not be sent", envKey)
	}
}

func TestNewCmdFactory_RegisterFail(t *testing.T) {
	mockRegister := func(id string) (*RegisteredLoop, error) {
		return nil, errors.New("registration failed")
	}

	cmdConfig := CmdConfig{
		ID:  "test-loop",
		Cmd: "echo",
		Env: []string{"TEST_ENV=1"},
	}

	_, err := NewCmdFactory(mockRegister, cmdConfig)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to register") {
		t.Errorf("Unexpected error message: %v", err)
	}
}
