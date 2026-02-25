package environment

import (
	"context"
	"encoding/json"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
)

func startRemoteComponent[T any](
	ctx context.Context,
	lggr zerolog.Logger,
	client componentClient,
	payload agent.StartComponentPayload,
	expectedComponentType string,
) (*T, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to encode %s payload", expectedComponentType)
	}

	response, err := client.StartComponent(ctx, agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationStartComponent,
		Payload:       payloadBytes,
	})
	if err != nil {
		return nil, err
	}
	if response.ComponentType != expectedComponentType {
		return nil, fmt.Errorf("unexpected component type in start response: %s", response.ComponentType)
	}
	for _, logLine := range response.AgentLogs {
		pretty := prettifyAgentLogLine(logLine)
		if pretty == "" {
			continue
		}
		lggr.Info().Msgf("[agent] %s", pretty)
	}

	output, err := agent.DecodeFromTransport[T](response.Output)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to decode %s transport payload", expectedComponentType)
	}
	return output, nil
}
