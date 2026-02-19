package environment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

type RemoteStopSummary struct {
	Requested int
	Stopped   int
	Missing   int
	Failed    int
}

// StopRemoteComponents sends StopComponent operations for all remote-targeted components.
// It is idempotent from the caller perspective; missing components are treated as success.
func StopRemoteComponents(ctx context.Context, lggr zerolog.Logger, cfg *config.Config) (RemoteStopSummary, error) {
	summary := RemoteStopSummary{}
	if cfg == nil {
		return summary, errors.New("config is nil")
	}
	summary.Requested = countRemoteStopTargets(cfg)
	if summary.Requested == 0 {
		return summary, nil
	}

	tunnelManager, err := newEC2TunnelManager(lggr)
	if err != nil {
		return summary, pkgerrors.Wrap(err, "failed to initialize tunnel manager for remote stop")
	}
	defer func() { _ = tunnelManager.Stop(ctx) }()

	startClient, err := newStartComponentClient(lggr, tunnelManager)
	if err != nil {
		return summary, pkgerrors.Wrap(err, "failed to initialize remote component client for stop")
	}

	var joined error
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain == nil || configuredBlockchain.Target != config.TargetRemote {
			continue
		}
		payload := agent.StartComponentPayload{
			ComponentType: componentTypeBlockchain,
			Blockchain:    configuredBlockchain.InputRef(),
			ReusePolicy:   string(configuredBlockchain.RemoteStartPolicy),
		}
		result, err := stopRemoteComponent(ctx, lggr, startClient, payload, componentTypeBlockchain)
		if err != nil {
			summary.Failed++
			joined = errors.Join(joined, err)
			continue
		}
		if result.Stopped {
			summary.Stopped++
		} else if !result.Found {
			summary.Missing++
		}
	}

	for _, nodeSet := range cfg.NodeSets {
		if nodeSet == nil || strings.TrimSpace(nodeSet.Target) != string(config.TargetRemote) {
			continue
		}
		payload := agent.StartComponentPayload{
			ComponentType: componentTypeNodeSet,
			NodeSet:       &simple_node_set.Input{Name: nodeSet.Name},
			ReusePolicy:   nodeSet.RemoteStartPolicy,
		}
		result, err := stopRemoteComponent(ctx, lggr, startClient, payload, componentTypeNodeSet)
		if err != nil {
			summary.Failed++
			joined = errors.Join(joined, err)
			continue
		}
		if result.Stopped {
			summary.Stopped++
		} else if !result.Found {
			summary.Missing++
		}
	}

	if cfg.JD != nil && cfg.JD.Target == config.TargetRemote {
		payload := agent.StartComponentPayload{
			ComponentType: componentTypeJD,
			JD:            cfg.JD.InputRef(),
			ReusePolicy:   string(cfg.JD.RemoteStartPolicy),
		}
		result, err := stopRemoteComponent(ctx, lggr, startClient, payload, componentTypeJD)
		if err != nil {
			summary.Failed++
			joined = errors.Join(joined, err)
			return summary, joined
		}
		if result.Stopped {
			summary.Stopped++
		} else if !result.Found {
			summary.Missing++
		}
	}

	return summary, joined
}

func countRemoteStopTargets(cfg *config.Config) int {
	if cfg == nil {
		return 0
	}
	count := 0
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Target == config.TargetRemote {
			count++
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Target) == string(config.TargetRemote) {
			count++
		}
	}
	if cfg.JD != nil && cfg.JD.Target == config.TargetRemote {
		count++
	}
	return count
}

func stopRemoteComponent(
	ctx context.Context,
	lggr zerolog.Logger,
	client componentClient,
	payload agent.StartComponentPayload,
	expectedType string,
) (*agent.StartComponentResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to encode stop payload for component type %s", payload.ComponentType)
	}

	response, err := client.StartComponent(ctx, agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStopComponent,
		Payload:       payloadBytes,
	})
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to stop remote component type %s", payload.ComponentType)
	}
	if response.ComponentType != expectedType {
		return nil, fmt.Errorf("unexpected component type in stop response: %s", response.ComponentType)
	}

	lggr.Info().
		Str("componentType", response.ComponentType).
		Bool("found", response.Found).
		Bool("stopped", response.Stopped).
		Msg("Processed remote component stop")

	for _, logLine := range response.AgentLogs {
		pretty := prettifyAgentLogLine(logLine)
		if pretty == "" {
			continue
		}
		lggr.Info().Msgf("[agent] %s", pretty)
	}

	return response, nil
}
