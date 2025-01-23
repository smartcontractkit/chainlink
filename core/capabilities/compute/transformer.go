package compute

import (
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

// ParsedConfig is a struct that contains the binary and config for a wasm module, as well as the module config.
type ParsedConfig struct {
	Binary []byte
	Config []byte

	// ModuleConfig is the configuration and dependencies to inject into the wasm module.
	ModuleConfig *host.ModuleConfig
}

// transformer implements the ConfigTransformer interface.  The logger and emitter are applied to
// the resulting ParsedConfig struct by default.  Override these values with the functional options.
type transformer struct {
	logger  logger.Logger
	emitter custmsg.MessageEmitter
	config  Config
}

func shallowCopy(m *values.Map) *values.Map {
	to := values.EmptyMap()

	for k, v := range m.Underlying {
		to.Underlying[k] = v
	}

	return to
}

// Transform attempts to read a valid ParsedConfig from an arbitrary values map.  The map must
// contain the binary and config keys.  Optionally the map may specify wasm module specific
// configuration values such as maxMemoryMBs, timeout, and tickInterval.  Default logger and
// emitter for the module are taken from the transformer instance.  Override these values with
// the functional options.
func (t *transformer) Transform(req capabilities.CapabilityRequest, opts ...func(*ParsedConfig)) (capabilities.CapabilityRequest, *ParsedConfig, error) {
	copiedReq := capabilities.CapabilityRequest{
		Inputs:   req.Inputs,
		Metadata: req.Metadata,
		Config:   shallowCopy(req.Config),
	}

	binary, err := popValue[[]byte](copiedReq.Config, binaryKey)
	if err != nil {
		return capabilities.CapabilityRequest{}, nil, NewInvalidRequestError(err)
	}

	config, err := popValue[[]byte](copiedReq.Config, configKey)
	if err != nil {
		return capabilities.CapabilityRequest{}, nil, NewInvalidRequestError(err)
	}

	maxMemoryMBs, err := popOptionalValue[uint64](copiedReq.Config, maxMemoryMBsKey)
	if err != nil {
		return capabilities.CapabilityRequest{}, nil, NewInvalidRequestError(err)
	}

	if maxMemoryMBs <= 0 || maxMemoryMBs > t.config.MaxMemoryMBs {
		maxMemoryMBs = t.config.MaxMemoryMBs
	}

	mc := &host.ModuleConfig{
		MaxMemoryMBs: maxMemoryMBs,
		Logger:       t.logger,
		Labeler:      t.emitter,
	}

	timeout, err := popOptionalValue[string](copiedReq.Config, timeoutKey)
	if err != nil {
		return capabilities.CapabilityRequest{}, nil, NewInvalidRequestError(err)
	}

	var td time.Duration
	if timeout != "" {
		td, err = time.ParseDuration(timeout)
		if err != nil {
			return capabilities.CapabilityRequest{}, nil, NewInvalidRequestError(err)
		}
		if td <= 0 || td > t.config.MaxTimeout {
			td = t.config.MaxTimeout
		}
		mc.Timeout = &td
	}

	if mc.Timeout == nil {
		mc.Timeout = &t.config.MaxTimeout
	}

	tickInterval, err := popOptionalValue[string](copiedReq.Config, tickIntervalKey)
	if err != nil {
		return capabilities.CapabilityRequest{}, nil, NewInvalidRequestError(err)
	}

	var ti time.Duration
	if tickInterval != "" {
		ti, err = time.ParseDuration(tickInterval)
		if err != nil {
			return capabilities.CapabilityRequest{}, nil, NewInvalidRequestError(err)
		}
		mc.TickInterval = ti
	}

	if mc.TickInterval <= 0 || mc.TickInterval > t.config.MaxTickInterval {
		mc.TickInterval = t.config.MaxTickInterval
	}

	pc := &ParsedConfig{
		Binary:       binary,
		Config:       config,
		ModuleConfig: mc,
	}

	for _, opt := range opts {
		opt(pc)
	}

	return copiedReq, pc, nil
}

func NewTransformer(lggr logger.Logger, emitter custmsg.MessageEmitter, config Config) *transformer {
	return &transformer{
		logger:  lggr,
		emitter: emitter,
		config:  config,
	}
}

// popOptionalValue attempts to pop a value from the map.  If the value is not found, the zero
// value for the type is returned and a nil error.
func popOptionalValue[T any](m *values.Map, key string) (T, error) {
	v, err := popValue[T](m, key)
	if err != nil {
		var nfe *NotFoundError
		if errors.As(err, &nfe) {
			return v, nil
		}
		return v, err
	}
	return v, nil
}

// popValue attempts to pop a value from the map.  If the value is not found, a NotFoundError is returned.
// If the value is found, it is unwrapped into the type T.  If the unwrapping fails, an error is returned.
func popValue[T any](m *values.Map, key string) (T, error) {
	var empty T

	wrapped, ok := m.Underlying[key]
	if !ok {
		return empty, NewNotFoundError(key)
	}

	delete(m.Underlying, key)
	err := wrapped.UnwrapTo(&empty)
	if err != nil {
		return empty, fmt.Errorf("could not unwrap value: %w", err)
	}

	return empty, nil
}

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("could not find %q in map", e.Key)
}

func NewNotFoundError(key string) *NotFoundError {
	return &NotFoundError{Key: key}
}

type InvalidRequestError struct {
	Err error
}

func (e *InvalidRequestError) Error() string {
	return fmt.Sprintf("invalid request: %v", e.Err)
}

func NewInvalidRequestError(err error) *InvalidRequestError {
	return &InvalidRequestError{Err: err}
}
