package workflows

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	mockedWorkflowID = "ef7c8168-f4d1-422f-a4b2-8ce0a1075f0a"
	mockedTriggerID  = "bd727a82-5cac-4071-be62-0152dd9adb0f"
)

type Engine struct {
	services.StateMachine
	logger     logger.Logger
	registry   types.CapabilitiesRegistry
	trigger    capabilities.TriggerCapability
	consensus  capabilities.ConsensusCapability
	target     capabilities.TargetCapability
	callbackCh chan capabilities.CapabilityResponse
	cancel     func()
}

func (e *Engine) Start(ctx context.Context) error {
	return e.StartOnce("Engine", func() error {
		err := e.registerTrigger(ctx)
		if err != nil {
			return err
		}

		// create a new context, since the one passed in via Start is short-lived.
		ctx, cancel := context.WithCancel(context.Background())
		e.cancel = cancel
		go e.loop(ctx)
		return nil
	})
}

func (e *Engine) registerTrigger(ctx context.Context) error {
	triggerConf, err := values.NewMap(
		map[string]any{
			"feedlist": []any{
				// ETHUSD, LINKUSD, USDBTC
				123, 456, 789,
			},
		},
	)
	if err != nil {
		return err
	}

	triggerInputs, err := values.NewMap(
		map[string]any{
			"triggerId": mockedTriggerID,
		},
	)
	if err != nil {
		return err
	}

	triggerRegRequest := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: mockedWorkflowID,
		},
		Config: triggerConf,
		Inputs: triggerInputs,
	}
	err = e.trigger.RegisterTrigger(ctx, e.callbackCh, triggerRegRequest)
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
				e.logger.Error("error executing event %+v: %w", resp, err)
			}
		}
	}
}

func (e *Engine) handleExecution(ctx context.Context, resp capabilities.CapabilityResponse) error {
	results, err := e.handleConsensus(ctx, resp)
	if err != nil {
		return err
	}

	_, err = e.handleTarget(ctx, results)
	return err
}

func (e *Engine) handleTarget(ctx context.Context, resp *values.List) (*values.List, error) {
	report, err := resp.Unwrap()
	if err != nil {
		return nil, err
	}
	inputs := map[string]values.Value{
		"report": resp,
	}
	config, err := values.NewMap(map[string]any{
		"address": "0xaabbcc",
		"method":  "updateFeedValues(report bytes, role uint8)",
		"params": []any{
			report, 1,
		},
	})
	if err != nil {
		return nil, err
	}

	tr := capabilities.CapabilityRequest{
		Inputs: &values.Map{Underlying: inputs},
		Config: config,
		Metadata: capabilities.RequestMetadata{
			WorkflowID: mockedWorkflowID,
		},
	}
	return capabilities.ExecuteSync(ctx, e.target, tr)
}

func (e *Engine) handleConsensus(ctx context.Context, resp capabilities.CapabilityResponse) (*values.List, error) {
	inputs := map[string]values.Value{
		"observations": resp.Value,
	}
	config, err := values.NewMap(map[string]any{
		"aggregation_method": "data_feeds_2_0",
		"aggregation_config": map[string]any{
			// ETHUSD
			"123": map[string]any{
				"deviation": "0.005",
				"heartbeat": "24h",
			},
			// LINKUSD
			"456": map[string]any{
				"deviation": "0.001",
				"heartbeat": "24h",
			},
			// BTCUSD
			"789": map[string]any{
				"deviation": "0.002",
				"heartbeat": "6h",
			},
		},
		"encoder": "EVM",
	})
	if err != nil {
		return nil, nil
	}
	cr := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: mockedWorkflowID,
		},
		Inputs: &values.Map{Underlying: inputs},
		Config: config,
	}
	return capabilities.ExecuteSync(ctx, e.consensus, cr)
}

func (e *Engine) Close() error {
	return e.StopOnce("Engine", func() error {
		defer e.cancel()

		triggerInputs, err := values.NewMap(
			map[string]any{
				"triggerId": mockedTriggerID,
			},
		)
		if err != nil {
			return err
		}
		deregRequest := capabilities.CapabilityRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID: mockedWorkflowID,
			},
			Inputs: triggerInputs,
		}
		return e.trigger.UnregisterTrigger(context.Background(), deregRequest)
	})
}

func NewEngine(lggr logger.Logger, registry types.CapabilitiesRegistry) (*Engine, error) {
	ctx := context.Background()
	trigger, err := registry.GetTrigger(ctx, "on_mercury_report")
	if err != nil {
		return nil, err
	}
	consensus, err := registry.GetConsensus(ctx, "off-chain-reporting")
	if err != nil {
		return nil, err
	}
	target, err := registry.GetTarget(ctx, "write_polygon_mainnet")
	if err != nil {
		return nil, err
	}
	return &Engine{
		logger:     lggr,
		registry:   registry,
		trigger:    trigger,
		consensus:  consensus,
		target:     target,
		callbackCh: make(chan capabilities.CapabilityResponse),
	}, nil
}
