package workflows

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// NOTE: max 32 bytes per ID - consider enforcing exactly 32 bytes?
	mockedWorkflowID  = "aaaaaaaaaa0000000000000000000000"
	mockedExecutionID = "bbbbbbbbbb0000000000000000000000"
	mockedTriggerID   = "cccccccccc0000000000000000000000"
)

type Engine struct {
	services.StateMachine
	logger          logger.Logger
	registry        types.CapabilitiesRegistry
	triggerType     string
	triggerConfig   *values.Map
	trigger         capabilities.TriggerCapability
	consensusType   string
	consensusConfig *values.Map
	consensus       capabilities.ConsensusCapability
	targets         []target
	callbackCh      chan capabilities.CapabilityResponse
	cancel          func()
}

type target struct {
	typeStr    string
	config     *values.Map
	capability capabilities.TargetCapability
}

func (e *Engine) Start(ctx context.Context) error {
	return e.StartOnce("Engine", func() error {
		// create a new context, since the one passed in via Start is short-lived.
		ctx, cancel := context.WithCancel(context.Background())
		e.cancel = cancel
		go e.init(ctx)
		go e.triggerHandlerLoop(ctx)
		return nil
	})
}

func (e *Engine) init(ctx context.Context) {
	retrySec := 5
	ticker := time.NewTicker(time.Duration(retrySec) * time.Second)
	defer ticker.Stop()
	var err error
LOOP:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.trigger, err = e.registry.GetTrigger(ctx, e.triggerType)
			if err != nil {
				e.logger.Errorf("failed to get trigger capability: %s, retrying in %d seconds", err, retrySec)
				break
			}
			e.consensus, err = e.registry.GetConsensus(ctx, e.consensusType)
			if err != nil {
				e.logger.Errorf("failed to get consensus capability: %s, retrying in %d seconds", err, retrySec)
				break
			}
			failed := false
			for i := range e.targets {
				e.targets[i].capability, err = e.registry.GetTarget(ctx, e.targets[i].typeStr)
				if err != nil {
					e.logger.Errorf("failed to get target capability: %s, retrying in %d seconds", err, retrySec)
					failed = true
					break
				}
			}
			if !failed {
				break LOOP
			}
		}
	}

	// we have all needed capabilities, now we can register for trigger events
	err = e.registerTrigger(ctx)
	if err != nil {
		e.logger.Errorf("failed to register trigger: %s", err)
	}

	// also register for consensus
	reg := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: mockedWorkflowID,
		},
		Config: e.consensusConfig,
	}
	err = e.consensus.RegisterToWorkflow(ctx, reg)
	if err != nil {
		e.logger.Errorf("failed to register consensus: %s", err)
	}

	e.logger.Info("engine initialized")
}

func (e *Engine) registerTrigger(ctx context.Context) error {
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
		Config: e.triggerConfig,
		Inputs: triggerInputs,
	}
	err = e.trigger.RegisterTrigger(ctx, e.callbackCh, triggerRegRequest)
	if err != nil {
		return fmt.Errorf("failed to instantiate mercury_trigger, %s", err)
	}
	return nil
}

func (e *Engine) triggerHandlerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-e.callbackCh:
			go e.handleExecution(ctx, resp)
		}
	}
}

func (e *Engine) handleExecution(ctx context.Context, event capabilities.CapabilityResponse) {
	e.logger.Debugw("executing on a trigger event", "event", event)
	result, err := e.handleConsensus(ctx, event)
	if err != nil {
		e.logger.Errorf("error in handleConsensus %v", err)
		return
	}
	err = e.handleTargets(ctx, result)
	if err != nil {
		e.logger.Error("error in handleTargets %v", err)
	}
}

func (e *Engine) handleConsensus(ctx context.Context, event capabilities.CapabilityResponse) (values.Value, error) {
	e.logger.Debugw("running consensus", "event", event)
	cr := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          mockedWorkflowID,
			WorkflowExecutionID: mockedExecutionID,
		},
		Inputs: &values.Map{
			Underlying: map[string]values.Value{
				// each node provides a single observation - outputs of mercury trigger
				"observations": &values.List{
					Underlying: []values.Value{event.Value},
				},
			},
		},
		Config: e.consensusConfig,
	}
	chReports := make(chan capabilities.CapabilityResponse, 10)
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	err := e.consensus.Execute(newCtx, chReports, cr)
	if err != nil {
		return nil, err
	}
	select {
	case resp := <-chReports:
		if resp.Err != nil {
			return nil, resp.Err
		}
		return resp.Value, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (e *Engine) handleTargets(ctx context.Context, resp values.Value) error {
	e.logger.Debugw("handle targets")
	inputs := map[string]values.Value{
		"report": resp,
	}

	var combinedErr error
	for _, t := range e.targets {
		e.logger.Debugw("sending to target", "target", t.typeStr, "inputs", inputs)
		tr := capabilities.CapabilityRequest{
			Inputs: &values.Map{Underlying: inputs},
			Config: t.config,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:          mockedWorkflowID,
				WorkflowExecutionID: mockedExecutionID,
			},
		}
		_, err := capabilities.ExecuteSync(ctx, t.capability, tr)
		combinedErr = errors.Join(combinedErr, err)
	}
	return combinedErr
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

func NewEngine(lggr logger.Logger, registry types.CapabilitiesRegistry) (engine *Engine, err error) {
	engine = &Engine{
		logger:     lggr.Named("WorkflowEngine"),
		registry:   registry,
		callbackCh: make(chan capabilities.CapabilityResponse),
	}

	// Trigger
	engine.triggerType = "on_mercury_report"
	engine.triggerConfig, err = values.NewMap(
		map[string]any{
			"feedlist": []any{
				"0x1111111111111111111100000000000000000000000000000000000000000000", // ETHUSD
				"0x2222222222222222222200000000000000000000000000000000000000000000", // LINKUSD
				"0x3333333333333333333300000000000000000000000000000000000000000000", // BTCUSD
			},
		},
	)
	if err != nil {
		return nil, err
	}

	// Consensus
	engine.consensusType = "offchain_reporting"
	engine.consensusConfig, err = values.NewMap(map[string]any{
		"aggregation_method": "data_feeds_2_0",
		"aggregation_config": map[string]any{
			// ETHUSD
			"0x1111111111111111111100000000000000000000000000000000000000000000": map[string]any{
				"deviation": decimal.NewFromFloat(0.001),
				"heartbeat": 1800,
			},
			// LINKUSD
			"0x2222222222222222222200000000000000000000000000000000000000000000": map[string]any{
				"deviation": decimal.NewFromFloat(0.001),
				"heartbeat": 1800,
			},
			// BTCUSD
			"0x3333333333333333333300000000000000000000000000000000000000000000": map[string]any{
				"deviation": decimal.NewFromFloat(0.001),
				"heartbeat": 1800,
			},
		},
		"encoder": "EVM",
		"encoder_config": map[string]any{
			"abi": "mercury_reports bytes[]",
		},
	})
	if err != nil {
		return nil, err
	}

	// Targets
	engine.targets = make([]target, 2)
	engine.targets[0].typeStr = "write_polygon-testnet-mumbai"
	engine.targets[0].config, err = values.NewMap(map[string]any{
		"address": "0x3F3554832c636721F1fD1822Ccca0354576741Ef",
		"params":  []any{"$(report)"},
		"abi":     "receive(report bytes)",
	})
	if err != nil {
		return nil, err
	}
	engine.targets[1].typeStr = "write_ethereum-testnet-sepolia"
	engine.targets[1].config, err = values.NewMap(map[string]any{
		"address": "0x54e220867af6683aE6DcBF535B4f952cB5116510",
		"params":  []any{"$(report)"},
		"abi":     "receive(report bytes)",
	})
	if err != nil {
		return nil, err
	}
	return
}
