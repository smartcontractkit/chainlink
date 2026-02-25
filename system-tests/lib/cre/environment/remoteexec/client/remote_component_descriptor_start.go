package client

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
)

type StartDescriptor[T any] struct {
	ComponentType string
	BuildPayload  func() (agent.StartComponentPayload, error)
	Rewrite       func(output *T, ec2HostIP string) error
}

func StartWithRuntimeDescriptor[T any](
	ctx context.Context,
	lggr zerolog.Logger,
	runtime *Runtime,
	descriptor StartDescriptor[T],
) (*T, error) {
	if runtime == nil {
		return nil, errors.New("remote runtime is required for remote component placement")
	}
	payload, err := descriptor.BuildPayload()
	if err != nil {
		return nil, err
	}
	output, err := StartRemoteComponent[T](
		ctx,
		lggr,
		runtime.Client,
		payload,
		descriptor.ComponentType,
	)
	if err != nil {
		return nil, err
	}
	if descriptor.Rewrite != nil {
		if err := descriptor.Rewrite(output, runtime.EC2HostIP); err != nil {
			return nil, err
		}
	}
	return output, nil
}
