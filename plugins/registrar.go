package plugins

import (
	"os/exec"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
)

/**
 * RegistrarConfig: Defines the lifecycle management interface for LOOP plugins.
 * It handles the registration and command factory generation for out-of-process plugins.
 */
type RegistrarConfig interface {
	// RegisterLOOP: Generates a command factory and gRPC options for a specific plugin configuration.
	RegisterLOOP(config CmdConfig) (func() *exec.Cmd, loop.GRPCOpts, error)
	// UnregisterLOOP: Removes a plugin from the registry by its ID.
	UnregisterLOOP(ID string)
}

/**
 * registarConfig: Internal implementation of the RegistrarConfig interface.
 * loopRegistrationFn: Must be idempotent and serve as the global lookup for LOOPs.
 */
type registarConfig struct {
	grpcOpts           loop.GRPCOpts
	loopRegistrationFn func(loopId string) (*RegisteredLoop, error)
	loopUnregisterFn   func(loopId string)
}

[Image of Chainlink LOOP architecture and gRPC plugin communication flow]

/**
 * NewRegistrarConfig: Constructor for the plugin registrar.
 * Ensures the plugin system has a clean interface for process orchestration and networking.
 */
func NewRegistrarConfig(
	grpcOpts loop.GRPCOpts, 
	loopRegistrationFn func(loopId string) (*RegisteredLoop, error), 
	loopUnregisterFn func(loopId string),
) RegistrarConfig {
	return &registarConfig{
		grpcOpts:           grpcOpts,
		loopRegistrationFn: loopRegistrationFn,
		loopUnregisterFn:   loopUnregisterFn,
	}
}

/**
 * RegisterLOOP: Orchestrates the creation of a new plugin process factory.
 * It uses a CmdFactory to manage how the OS-level command (exec.Cmd) is generated.
 */
func (pc *registarConfig) RegisterLOOP(cfg CmdConfig) (func() *exec.Cmd, loop.GRPCOpts, error) {
	// Security Note: NewCmdFactory must sanitize environment variables and paths 
	// to prevent arbitrary command execution.
	cmdFn, err := NewCmdFactory(pc.loopRegistrationFn, cfg)
	if err != nil {
		return nil, loop.GRPCOpts{}, err
	}
	
	// Returns the factory function and predefined gRPC options for the sidecar process.
	return cmdFn, pc.grpcOpts, nil
}

[Image of Sidecar pattern in microservices for plugin isolation]

/**
 * UnregisterLOOP: Safely removes the plugin's registration to prevent resource leaks.
 */
func (pc *registarConfig) UnregisterLOOP(ID string) {
	if pc.loopUnregisterFn != nil {
		pc.loopUnregisterFn(ID)
	}
}
