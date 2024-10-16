package compute

import (
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type Transformer[T any, U any] interface {
	Transform(T, ...func(*U)) (*U, error)
}

type ConfigTransformer = Transformer[*values.Map, ParsedConfig]

type ParsedConfig struct {
	Binary       []byte
	Config       []byte
	ModuleConfig *host.ModuleConfig
}

type transformer struct{}

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

func NewTransformer() *transformer {
	return &transformer{}
}

func WithLogger(l logger.Logger) func(*ParsedConfig) {
	return func(pc *ParsedConfig) {
		pc.ModuleConfig.Logger = l
	}
}

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
