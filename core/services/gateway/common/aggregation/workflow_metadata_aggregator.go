package aggregation

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
)

type WorkflowMetadataAggregator struct {
	services.StateMachine
	lggr      logger.Logger
	stopCh    services.StopChan
	threshold int
	mu        sync.RWMutex
	// observations is a map that tracks auth data from workflow nodes.
	// keyed by workflow digest
	observations map[string]*NodeObservations
	// observedAt is a map from node address to a map of workflow digest to last observed time
	// This is used to clean up old observations that are no longer relevant.
	observedAt      map[string]map[string]time.Time
	cleanupInterval time.Duration
	metrics         *metrics.Metrics
}

func NewWorkflowMetadataAggregator(lggr logger.Logger, threshold int, cleanupInterval time.Duration, metrics *metrics.Metrics) *WorkflowMetadataAggregator {
	if threshold <= 0 {
		panic(fmt.Sprintf("threshold must be greater than 0, got %d", threshold))
	}
	return &WorkflowMetadataAggregator{
		lggr:            logger.Named(lggr, "WorkflowMetadataAggregator"),
		threshold:       threshold,
		observations:    make(map[string]*NodeObservations),
		observedAt:      make(map[string]map[string]time.Time),
		stopCh:          make(services.StopChan),
		cleanupInterval: cleanupInterval,
		metrics:         metrics,
	}
}

func (agg *WorkflowMetadataAggregator) reapObservations(ctx context.Context) {
	agg.mu.Lock()
	defer agg.mu.Unlock()
	now := time.Now()
	var expiredCount int
	for node, digestObservedAt := range agg.observedAt {
		for digest, observedAt := range digestObservedAt {
			if now.Sub(observedAt) > agg.cleanupInterval {
				delete(agg.observedAt[node], digest)
				if len(agg.observedAt[node]) == 0 {
					delete(agg.observedAt, node)
				}
				_, ok := agg.observations[digest]
				if !ok {
					agg.lggr.Warnw("Observation digest not found in observations", "digest", digest, "node", node)
					continue
				}
				agg.observations[digest].nodes.Remove(node)
				if len(agg.observations[digest].nodes) == 0 {
					delete(agg.observations, digest)
				}
				expiredCount++
			}
		}
	}
	if expiredCount > 0 {
		agg.metrics.Trigger.IncrementMetadataObservationsCleanUpCount(ctx, int64(expiredCount), agg.lggr)
		agg.lggr.Debugw("Removed expired callbacks", "count", expiredCount)
	}
	agg.metrics.Trigger.RecordMetadataObservationsCount(ctx, int64(len(agg.observations)), agg.lggr)
}

func (agg *WorkflowMetadataAggregator) Start(ctx context.Context) error {
	return agg.StartOnce("WorkflowMetadataAggregator", func() error {
		agg.lggr.Info("Starting WorkflowMetadataAggregator")
		go func() {
			ticker := time.NewTicker(agg.cleanupInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					agg.reapObservations(ctx)
				case <-agg.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (agg *WorkflowMetadataAggregator) Close() error {
	return agg.StopOnce("WorkflowMetadataAggregator", func() error {
		agg.lggr.Info("Stopping WorkflowMetadataAggregator")
		close(agg.stopCh)
		return nil
	})
}

// Collect adds an observation from a workflow node to the aggregator.
func (agg *WorkflowMetadataAggregator) Collect(obs *gateway_common.WorkflowMetadata, nodeAddress string) error {
	if obs.WorkflowSelector.WorkflowID == "" || obs.WorkflowSelector.WorkflowName == "" ||
		obs.WorkflowSelector.WorkflowOwner == "" || obs.WorkflowSelector.WorkflowTag == "" {
		return errors.New("observation is missing required fields")
	}
	if nodeAddress == "" {
		return errors.New("node address cannot be empty")
	}
	agg.mu.Lock()
	defer agg.mu.Unlock()
	digest, err := obs.Digest()
	if err != nil {
		return err
	}
	_, ok := agg.observedAt[nodeAddress]
	if !ok {
		agg.observedAt[nodeAddress] = make(map[string]time.Time)
	}
	agg.observedAt[nodeAddress][digest] = time.Now()

	_, ok = agg.observations[digest]
	if !ok {
		agg.observations[digest] = &NodeObservations{
			observation: obs,
			nodes:       make(StringSet),
		}
	}
	agg.observations[digest].nodes.Add(nodeAddress)
	return nil
}

// Aggregate returns the aggregated workflow metadata for workflows that have reached the threshold.
func (agg *WorkflowMetadataAggregator) Aggregate() ([]gateway_common.WorkflowMetadata, error) {
	agg.mu.RLock()
	defer agg.mu.RUnlock()

	var aggregated []gateway_common.WorkflowMetadata
	for _, nodeObs := range agg.observations {
		if len(nodeObs.nodes) >= agg.threshold {
			aggregated = append(aggregated, *nodeObs.observation)
		}
	}
	return aggregated, nil
}

type NodeObservations struct {
	observation *gateway_common.WorkflowMetadata
	nodes       StringSet
}
