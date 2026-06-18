package contracts

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

const (
	capabilityRegistrySyncPollInterval = 2 * time.Second
	capabilityRegistrySyncTimeout      = 2 * time.Minute
	capabilityRegistrySyncQueryTimeout = 3 * time.Second
	capabilityRegistrySyncConcurrency  = 4
)

type workflowWorkerTarget struct {
	donName   string
	nodeIndex int
	dbPort    int
}

type capabilityRegistrySyncState struct {
	IDsToDONs         map[string]json.RawMessage `json:"IDsToDONs"`
	IDsToNodes        map[string]json.RawMessage `json:"IDsToNodes"`
	IDsToCapabilities map[string]json.RawMessage `json:"IDsToCapabilities"`
}

const latestCapabilityRegistrySyncStateQuery = `
SELECT data
FROM registry_syncer_states
ORDER BY id DESC
LIMIT 1
`

func waitForWorkflowWorkersCapabilityRegistrySync(ctx context.Context, input cre.ConfigureCapabilityRegistryInput) error {
	// Waiting for capability registry sync is not supported in Kubernetes mode.
	if input.Provider.IsKubernetes() {
		return nil
	}
	targets, tErr := workflowWorkerTargets(input.Topology, input.NodeSets)
	if tErr != nil {
		return tErr
	}
	if len(targets) == 0 {
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, capabilityRegistrySyncTimeout)
	defer cancel()

	pending := make(map[string]workflowWorkerTarget, len(targets))
	lastState := make(map[string]string, len(targets))
	for _, target := range targets {
		key := registryTargetKey(target)
		pending[key] = target
		lastState[key] = "awaiting first successful registry snapshot check"
	}

	ticker := time.NewTicker(capabilityRegistrySyncPollInterval)
	defer ticker.Stop()

	for {
		type checkResult struct {
			key   string
			ready bool
			state string
		}
		results := make([]checkResult, 0, len(pending))
		resultsMu := sync.Mutex{}
		eg, egCtx := errgroup.WithContext(timeoutCtx)
		eg.SetLimit(capabilityRegistrySyncConcurrency)
		for key, target := range pending {
			eg.Go(func() error {
				ready, state := hasCapabilityRegistrySyncOnWorker(egCtx, target.dbPort, target.nodeIndex)
				resultsMu.Lock()
				results = append(results, checkResult{
					key:   key,
					ready: ready,
					state: state,
				})
				resultsMu.Unlock()
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		for _, result := range results {
			if result.ready {
				delete(pending, result.key)
				delete(lastState, result.key)
				continue
			}
			lastState[result.key] = result.state
		}

		if len(pending) == 0 {
			return nil
		}

		select {
		case <-timeoutCtx.Done():
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("timed out after %.0f seconds waiting for workflow workers to sync capability registry state: %s", capabilityRegistrySyncTimeout.Seconds(), formatCapabilityRegistrySyncPending(lastState))
			}
			return timeoutCtx.Err()
		case <-ticker.C:
		}
	}
}

func workflowWorkerTargets(topology *cre.Topology, nodeSets []*cre.NodeSet) ([]workflowWorkerTarget, error) {
	if topology == nil || topology.DonsMetadata == nil {
		return nil, errors.New("topology metadata cannot be nil")
	}

	allDons := topology.DonsMetadata.List()
	indexByName := make(map[string]int, len(allDons))
	for i, don := range allDons {
		indexByName[don.Name] = i
	}

	workflowDons, err := topology.DonsMetadata.WorkflowDONs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve workflow DONs")
	}

	targets := make([]workflowWorkerTarget, 0)
	for _, workflowDON := range workflowDons {
		donIdx, ok := indexByName[workflowDON.Name]
		if !ok {
			return nil, fmt.Errorf("workflow DON %s not found in topology list", workflowDON.Name)
		}
		if donIdx >= len(nodeSets) || nodeSets[donIdx] == nil {
			return nil, fmt.Errorf("nodeset for workflow DON %s is missing", workflowDON.Name)
		}

		dbPort := nodeSets[donIdx].DbInput.Port
		if dbPort == 0 {
			defaultPort, dErr := strconv.Atoi(postgres.Port)
			if dErr != nil {
				return nil, errors.Wrap(dErr, "failed to convert postgres port to int")
			}
			dbPort = defaultPort
		}

		workers, wErr := workflowDON.Workers()
		if wErr != nil {
			return nil, errors.Wrapf(wErr, "failed to resolve workers for workflow DON %s", workflowDON.Name)
		}

		for _, worker := range workers {
			targets = append(targets, workflowWorkerTarget{
				donName:   workflowDON.Name,
				nodeIndex: worker.Index,
				dbPort:    dbPort,
			})
		}
	}

	return targets, nil
}

func hasCapabilityRegistrySyncOnWorker(ctx context.Context, dbPort, nodeIndex int) (bool, string) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=3", "127.0.0.1", dbPort, postgres.User, postgres.Password, fmt.Sprintf("db_%d", nodeIndex))
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return false, fmt.Sprintf("failed to open db connection: %v", err)
	}
	defer db.Close()

	queryCtx, cancel := context.WithTimeout(ctx, capabilityRegistrySyncQueryTimeout)
	defer cancel()

	var rawData []byte
	if err = db.GetContext(queryCtx, &rawData, latestCapabilityRegistrySyncStateQuery); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "registry_syncer_states is empty"
		}
		return false, fmt.Sprintf("failed to query latest registry syncer state: %v", err)
	}
	if len(rawData) == 0 {
		return false, "latest registry_syncer_states row has empty data payload"
	}

	var state capabilityRegistrySyncState
	if err = json.Unmarshal(rawData, &state); err != nil {
		return false, fmt.Sprintf("failed to unmarshal latest registry syncer state payload: %v", err)
	}

	hasDONs := len(state.IDsToDONs) > 0
	hasNodes := len(state.IDsToNodes) > 0
	hasCapabilities := len(state.IDsToCapabilities) > 0
	if !hasDONs || !hasCapabilities || !hasNodes {
		return false, fmt.Sprintf("incomplete registry snapshot (has_dons=%t has_nodes=%t has_capabilities=%t)", hasDONs, hasNodes, hasCapabilities)
	}

	return true, ""
}

func registryTargetKey(target workflowWorkerTarget) string {
	return fmt.Sprintf("%s/%d", target.donName, target.nodeIndex)
}

func formatCapabilityRegistrySyncPending(lastState map[string]string) string {
	parts := make([]string, 0, len(lastState))
	keys := make([]string, 0, len(lastState))
	for target := range lastState {
		keys = append(keys, target)
	}
	sort.Strings(keys)

	for _, target := range keys {
		reason := lastState[target]
		parts = append(parts, fmt.Sprintf("%s (%s)", target, reason))
	}
	return strings.Join(parts, "; ")
}
