package bridgestatus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/bridgestatus/events"
)

// Service polls Bridge status and pushes them to Beholder
type Service struct {
	services.Service
	eng *services.Engine

	config     config.BridgeStatusReporter
	bridgeORM  bridges.ORM
	jobORM     job.ORM
	httpClient *http.Client
	emitter    beholder.Emitter
}

const (
	ServiceName        = "BridgeStatusReporter"
	bridgePollPageSize = 1_000
)

// NewBridgeStatusReporter creates a new Bridge Status Reporter Service
func NewBridgeStatusReporter(
	config config.BridgeStatusReporter,
	bridgeORM bridges.ORM,
	jobORM job.ORM,
	httpClient *http.Client,
	emitter beholder.Emitter,
	lggr logger.Logger,
) *Service {
	s := &Service{
		config:     config,
		bridgeORM:  bridgeORM,
		jobORM:     jobORM,
		httpClient: httpClient,
		emitter:    emitter,
	}
	s.Service, s.eng = services.Config{
		Name:  ServiceName,
		Start: s.start,
	}.NewServiceEngine(lggr)
	return s
}

// start starts the Bridge Status Reporter Service
func (s *Service) start(ctx context.Context) error {
	if !s.config.Enabled() {
		s.eng.Info("Bridge Status Reporter Service is disabled")
		return nil
	}

	s.eng.Info("Starting Bridge Status Reporter Service")

	// Start periodic polling using Engine's ticker support
	ticker := services.NewTicker(s.config.PollingInterval())
	s.eng.GoTick(ticker, s.pollAllBridges)

	return nil
}

// HealthReport returns the service health
func (s *Service) HealthReport() map[string]error {
	return map[string]error{ServiceName: s.Ready()}
}

// pollAllBridges polls all registered bridges using pagination
func (s *Service) pollAllBridges(ctx context.Context) {
	var allBridges []bridges.BridgeType
	var offset = 0

	// Paginate through all bridges
	for {
		bridgeList, _, err := s.bridgeORM.BridgeTypes(ctx, offset, bridgePollPageSize)
		if err != nil {
			s.eng.Debugw("Failed to fetch bridges", "error", err, "offset", offset)
			return
		}

		allBridges = append(allBridges, bridgeList...)

		// If we got fewer than pageSize bridges, we've reached the end
		if len(bridgeList) < bridgePollPageSize {
			break
		}

		offset += bridgePollPageSize
	}

	if len(allBridges) == 0 {
		s.eng.Debug("No bridges configured for Bridge Status Reporter polling")
		return
	}

	s.eng.Debugw("Polling Bridge Status Reporter for all bridges", "count", len(allBridges))

	// Poll each bridge concurrently and wait for completion
	var wg sync.WaitGroup
	for _, bridge := range allBridges {
		wg.Add(1)
		bridgeName := string(bridge.Name)
		bridgeURL := bridge.URL.String()
		go func(name, url string) {
			defer wg.Done()
			s.pollBridge(ctx, name, url)
		}(bridgeName, bridgeURL)
	}

	wg.Wait()
}

// handleBridgeError handles errors during bridge polling, either skipping or emitting empty telemetry
func (s *Service) handleBridgeError(ctx context.Context, bridgeName string, jobs []JobInfo, logMsg string, logFields ...any) {
	s.eng.Debugw(logMsg, logFields...)
	if s.config.IgnoreInvalidBridges() {
		return
	}
	// If not ignoring invalid bridges, still emit empty telemetry
	s.emitBridgeStatus(ctx, bridgeName, EAResponse{}, jobs)
}

// pollBridge polls a single bridge's status endpoint
func (s *Service) pollBridge(ctx context.Context, bridgeName string, bridgeURL string) {
	s.eng.Debugw("Polling bridge", "bridge", bridgeName, "url", bridgeURL)

	// Look up jobs associated with this bridge first
	jobs, err := s.findJobsForBridge(ctx, bridgeName)
	if err != nil {
		s.eng.Warnw("Failed to find jobs for bridge", "bridge", bridgeName, "error", err)
		jobs = []JobInfo{}
	}

	// Skip bridge if it has no jobs and ignoreJoblessBridges is enabled
	if s.config.IgnoreJoblessBridges() && len(jobs) == 0 {
		s.eng.Debugw("Skipping bridge with no jobs", "bridge", bridgeName, "ignoreJoblessBridges", true)
		return
	}

	// Parse bridge URL and construct status endpoint
	parsedURL, err := url.Parse(bridgeURL)
	if err != nil {
		s.handleBridgeError(ctx, bridgeName, jobs, "Failed to parse bridge URL", "bridge", bridgeName, "url", bridgeURL, "error", err)
		return
	}

	// Construct status endpoint URL
	statusURL := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   path.Join(parsedURL.Path, s.config.StatusPath()),
	}

	// Make HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", statusURL.String(), nil)
	if err != nil {
		s.handleBridgeError(ctx, bridgeName, jobs, "Failed to create request for Bridge Status Reporter status", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.handleBridgeError(ctx, bridgeName, jobs, "Failed to fetch Bridge Status Reporter status", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.handleBridgeError(ctx, bridgeName, jobs, "Bridge Status Reporter status endpoint returned non-200 status", "bridge", bridgeName, "url", statusURL.String(), "status", resp.StatusCode)
		return
	}

	// Parse response
	var status EAResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		s.handleBridgeError(ctx, bridgeName, jobs, "Failed to decode Bridge Status Reporter status", "bridge", bridgeName, "url", statusURL.String(), "error", err)
		return
	}

	s.eng.Debugw("Successfully fetched Bridge Status Reporter status", "bridge", bridgeName, "adapter", status.Adapter.Name, "version", status.Adapter.Version)

	// Emit telemetry to Beholder
	s.emitBridgeStatus(ctx, bridgeName, status, jobs)
}

// emitBridgeStatus sends Bridge Status Reporter data to Beholder
func (s *Service) emitBridgeStatus(ctx context.Context, bridgeName string, status EAResponse, jobs []JobInfo) {
	// Convert runtime info
	runtime := &events.RuntimeInfo{
		NodeVersion:  status.Runtime.NodeVersion,
		Platform:     status.Runtime.Platform,
		Architecture: status.Runtime.Architecture,
		Hostname:     status.Runtime.Hostname,
	}

	// Convert metrics info
	metrics := &events.MetricsInfo{
		Enabled: status.Metrics.Enabled,
	}

	// Convert endpoints
	endpointsProto := make([]*events.EndpointInfo, len(status.Endpoints))
	for i, endpoint := range status.Endpoints {
		endpointsProto[i] = &events.EndpointInfo{
			Name:       endpoint.Name,
			Aliases:    endpoint.Aliases,
			Transports: endpoint.Transports,
		}
	}

	// Convert configuration including values
	configProto := make([]*events.ConfigurationItem, len(status.Configuration))
	for i, config := range status.Configuration {
		// Helper function to convert values to strings, handling nil (ternary-style)
		safeString := func(v any) string {
			if v == nil {
				return ""
			}
			return fmt.Sprintf("%v", v)
		}

		configProto[i] = &events.ConfigurationItem{
			Name:               config.Name,
			Value:              safeString(config.Value),
			Type:               config.Type,
			Description:        config.Description,
			Required:           config.Required,
			DefaultValue:       safeString(config.Default),
			CustomSetting:      config.CustomSetting,
			EnvDefaultOverride: safeString(config.EnvDefaultOverride),
		}
	}

	// Convert jobs to protobuf JobInfo structs
	jobsProto := make([]*events.JobInfo, 0, len(jobs))
	for _, job := range jobs {
		jobsProto = append(jobsProto, &events.JobInfo{
			ExternalJobId: job.ExternalJobID,
			JobName:       job.Name,
		})
	}

	// Create the protobuf event
	event := &events.BridgeStatusEvent{
		BridgeName:           bridgeName,
		AdapterName:          status.Adapter.Name,
		AdapterVersion:       status.Adapter.Version,
		AdapterUptimeSeconds: status.Adapter.UptimeSeconds,
		DefaultEndpoint:      status.DefaultEndpoint,
		Runtime:              runtime,
		Metrics:              metrics,
		Endpoints:            endpointsProto,
		Configuration:        configProto,
		Jobs:                 jobsProto,
	}

	// Emit the protobuf event through the configured emitter
	if err := events.EmitBridgeStatusEvent(ctx, s.emitter, event); err != nil {
		s.eng.Warnw("Failed to emit Bridge Status Reporter protobuf data to Beholder", "bridge", bridgeName, "error", err)
		return
	}

	s.eng.Debugw("Successfully emitted Bridge Status Reporter protobuf data to Beholder",
		"bridge", bridgeName,
		"adapter", status.Adapter.Name,
		"version", status.Adapter.Version,
	)
}

// findJobsForBridge finds jobs associated with a bridge name
func (s *Service) findJobsForBridge(ctx context.Context, bridgeName string) ([]JobInfo, error) {
	// Find job IDs that use this bridge
	jobIDs, err := s.jobORM.FindJobIDsWithBridge(ctx, bridgeName)
	if err != nil {
		return nil, fmt.Errorf("failed to find jobs with bridge %s: %w", bridgeName, err)
	}

	if len(jobIDs) == 0 {
		s.eng.Debugw("No jobs found for bridge", "bridge", bridgeName)
		return []JobInfo{}, nil
	}

	// Convert job IDs to job info
	jobs := make([]JobInfo, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		job, err := s.jobORM.FindJob(ctx, jobID)
		if err != nil {
			s.eng.Debugw("Failed to find job", "jobID", jobID, "bridge", bridgeName, "error", err)
			continue
		}

		// Get job name, use a default if not set
		jobName := "unknown"
		if job.Name.Valid && job.Name.String != "" {
			jobName = job.Name.String
		}

		jobs = append(jobs, JobInfo{
			ExternalJobID: job.ExternalJobID.String(),
			Name:          jobName,
		})
	}

	s.eng.Debugw("Found jobs for bridge", "bridge", bridgeName, "count", len(jobs))

	return jobs, nil
}
