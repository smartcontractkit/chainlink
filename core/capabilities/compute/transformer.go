package compute

import (
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type Transformer[T any, U any] interface {
	// Transform changes a struct of type T into a struct of type U.  Accepts a variadic list of options to modify the
	// output struct.
	Transform(T, ...func(*U)) (*U, error)
}

// ConfigTransformer is a Transformer that converts a values.Map into a ParsedConfig struct.
type ConfigTransformer = Transformer[*values.Map, ParsedConfig]

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
}

// Transform attempts to read a valid ParsedConfig from an arbitrary values map.  The map must
// contain the binary and config keys.  Optionally the map may specify wasm module specific
// configuration values such as maxMemoryMBs, timeout, and tickInterval.  Default logger and
// emitter for the module are taken from the transformer instance.  Override these values with
// the functional options.
func (t *transformer) Transform(in *values.Map, opts ...func(*ParsedConfig)) (*ParsedConfig, error) {
	binary, err := popValue[[]byte](in, binaryKey)
	if err != nil {
		return nil, NewInvalidRequestError(err)
	}

	config, err := popValue[[]byte](in, configKey)
	if err != nil {
		return nil, NewInvalidRequestError(err)
	}

	maxMemoryMBs, err := popOptionalValue[int64](in, maxMemoryMBsKey)
	if err != nil {
		return nil, NewInvalidRequestError(err)
	}

	mc := &host.ModuleConfig{
		MaxMemoryMBs: maxMemoryMBs,
		Logger:       t.logger,
		Labeler:      t.emitter,
	}

	timeout, err := popOptionalValue[string](in, timeoutKey)
	if err != nil {
		return nil, NewInvalidRequestError(err)
	}

	var td time.Duration
	if timeout != "" {
		td, err = time.ParseDuration(timeout)
		if err != nil {
			return nil, NewInvalidRequestError(err)
		}
		mc.Timeout = &td
	}

	tickInterval, err := popOptionalValue[string](in, tickIntervalKey)
	if err != nil {
		return nil, NewInvalidRequestError(err)
	}

	var ti time.Duration
	if tickInterval != "" {
		ti, err = time.ParseDuration(tickInterval)
		if err != nil {
			return nil, NewInvalidRequestError(err)
		}
		mc.TickInterval = ti
	}

	pc := &ParsedConfig{
		Binary:       binary,
		Config:       config,
		ModuleConfig: mc,
	}

	for _, opt := range opts {
		opt(pc)
	}

	return pc, nil
}

func NewTransformer(lggr logger.Logger, emitter custmsg.MessageEmitter) *transformer {
	return &transformer{
		logger:  lggr,
		emitter: emitter,
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
