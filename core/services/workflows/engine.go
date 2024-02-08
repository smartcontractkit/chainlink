package workflows

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Engine struct {
	services.StateMachine
	registry   types.CapabilitiesRegistry
	callbackCh chan capabilities.CapabilityResponse
}

func (e *Engine) Start(ctx context.Context) error {
	return e.StartOnce("Engine", func() error {
		err := e.registerMercuryTrigger(ctx)
		if err != nil {
			return err
		}
		go e.loop(ctx)
		return nil
	})
}

func (e *Engine) registerMercuryTrigger(ctx context.Context) error {
	tr, err := e.registry.GetTrigger(ctx, "mercury_trigger")
	if err != nil {
		return fmt.Errorf("could not fetch mercury_trigger, %w", err)
	}

	err = tr.RegisterTrigger(ctx, e.callbackCh, capabilities.CapabilityRequest{})
	if err != nil {
		return fmt.Errorf("failed to instantiate mercury_trigger, %s", err)
	}
	return nil
}

func (e *Engine) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-e.callbackCh:
			err := e.handleExecution(ctx, resp)
			if err != nil {
				// TODO: log.
			}
		}
	}
}

func (e *Engine) handleExecution(ctx context.Context, resp capabilities.CapabilityResponse) error {
	ch := make(chan capabilities.CapabilityResponse)
	action, err := e.registry.GetAction(ctx, "")
	if err != nil {
		return err
	}

	ar := capabilities.CapabilityRequest{}
	action.Execute(ctx, ch, ar)

	resp = <-ch
	if resp.Err != nil {
		return resp.Err
	}

	consensus, err := e.registry.GetConsensus(ctx, "")
	if err != nil {
		return err
	}

	cr := capabilities.CapabilityRequest{
		Inputs: resp.Value.(*values.Map),
	}
	consensus.Execute(ctx, ch, cr)

	resp = <-ch
	if resp.Err != nil {
		return resp.Err
	}

	target, err := e.registry.GetTarget(ctx, "")
	if err != nil {
		return err
	}

	tr := capabilities.CapabilityRequest{
		Inputs: resp.Value.(*values.Map),
	}
	target.Execute(ctx, ch, tr)

	resp = <-ch
	return resp.Err
}

func (e *Engine) Close() error {
	return nil
}

func NewEngine(registry types.CapabilitiesRegistry) (*Engine, error) {
	return &Engine{registry: registry}, nil
}
