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
