package environment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	ResidualContainers []string
	ResidualVolumes    []string
	ResidualQueryError string
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

	remoteRuntime, err := resolveRemoteRuntime(lggr)
	if err != nil {
		return summary, pkgerrors.Wrap(err, "failed to resolve remote runtime settings for stop")
	}

	var joined error
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain == nil || configuredBlockchain.Placement != config.PlacementRemote {
			continue
		}
		payload := agent.StartComponentPayload{
			ComponentType: componentTypeBlockchain,
			Blockchain:    configuredBlockchain.InputRef(),
			ReusePolicy:   string(configuredBlockchain.RemoteStartPolicy),
		}
		result, err := stopRemoteComponent(ctx, lggr, remoteRuntime.Client, payload, componentTypeBlockchain)
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
		if nodeSet == nil || strings.TrimSpace(nodeSet.Placement) != string(config.PlacementRemote) {
			continue
		}
		payload := agent.StartComponentPayload{
			ComponentType: componentTypeNodeSet,
			NodeSet:       &simple_node_set.Input{Name: nodeSet.Name},
			ReusePolicy:   nodeSet.RemoteStartPolicy,
		}
		result, err := stopRemoteComponent(ctx, lggr, remoteRuntime.Client, payload, componentTypeNodeSet)
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

	if cfg.JD != nil && cfg.JD.Placement == config.PlacementRemote {
		payload := agent.StartComponentPayload{
			ComponentType: componentTypeJD,
			JD:            cfg.JD.InputRef(),
			ReusePolicy:   string(cfg.JD.RemoteStartPolicy),
		}
		result, err := stopRemoteComponent(ctx, lggr, remoteRuntime.Client, payload, componentTypeJD)
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

	containers, volumes, listErr := listRemoteCTFResources(ctx, remoteRuntime.AgentBaseURL)
	if listErr != nil {
		summary.ResidualQueryError = listErr.Error()
	} else {
		summary.ResidualContainers = containers
		summary.ResidualVolumes = volumes
	}

	return summary, joined
}

func countRemoteStopTargets(cfg *config.Config) int {
	if cfg == nil {
		return 0
	}
	count := 0
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Placement == config.PlacementRemote {
			count++
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Placement) == string(config.PlacementRemote) {
			count++
		}
	}
	if cfg.JD != nil && cfg.JD.Placement == config.PlacementRemote {
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

func listRemoteCTFResources(
	ctx context.Context,
	baseURL string,
) ([]string, []string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/v1/resources/ctf", nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("ctf resource query failed: status %s body %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var out agent.CTFResourcesResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, nil, err
	}
	return out.Containers, out.Volumes, nil
}
