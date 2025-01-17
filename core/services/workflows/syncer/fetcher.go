package syncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

type FetcherService struct {
	services.StateMachine
	lggr    logger.Logger
	och     *webapi.OutgoingConnectorHandler
	wrapper gatewayConnector
	limits  *ArtifactConfig
}

type ArtifactConfig struct {
	MaxConfigSize  uint64
	MaxSecretsSize uint64
	MaxBinarySize  uint64
}

type ArtifactType string

var (
	ArtifactTypeConfig  ArtifactType = "config"
	ArtifactTypeSecrets ArtifactType = "secrets"
	ArtifactTypeBinary  ArtifactType = "binary"
	ArtifactTypeUnknown ArtifactType = "unknown"
)

const defaultMaxArtifactSizeBytes = uint32(10 * 1024 * 1024) // 10MB

type gatewayConnector interface {
	GetGatewayConnector() connector.GatewayConnector
}

func WithMaxArtifactSize(cfg ArtifactConfig) func(*FetcherService) {
	return func(fs *FetcherService) {
		fs.limits = &cfg
	}
}

func NewFetcherService(lggr logger.Logger, wrapper gatewayConnector, opts ...func(*FetcherService)) *FetcherService {
	fs := &FetcherService{
		lggr:    lggr.Named("FetcherService"),
		wrapper: wrapper,
	}
	for _, opt := range opts {
		opt(fs)
	}
	return fs
}

func (s *FetcherService) Start(ctx context.Context) error {
	return s.StartOnce("FetcherService", func() error {
		connector := s.wrapper.GetGatewayConnector()

		outgoingConnectorLggr := s.lggr.Named("OutgoingConnectorHandler")

		webAPIConfig := webapi.ServiceConfig{
			RateLimiter: common.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		}

		och, err := webapi.NewOutgoingConnectorHandler(connector,
			webAPIConfig,
			capabilities.MethodWorkflowSyncer, outgoingConnectorLggr)
		if err != nil {
			return fmt.Errorf("could not create outgoing connector handler: %w", err)
		}

		s.och = och
		return och.Start(ctx)
	})
}

func (s *FetcherService) Close() error {
	return s.StopOnce("FetcherService", func() error {
		return s.och.Close()
	})
}

func (s *FetcherService) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *FetcherService) Name() string {
	return s.lggr.Name()
}

func hash(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}

type FetchMaxCmd struct {
	URL          string       `json:"url"`
	ArtifactType ArtifactType `json:"artifactType"`
}

type MaxFetcher interface {
	FetchMax(ctx context.Context, cmd FetchMaxCmd) ([]byte, error)
}

func (s *FetcherService) FetchMax(ctx context.Context, cmd FetchMaxCmd) ([]byte, error) {
	n, err := s.getMaxBytes(cmd.ArtifactType)
	if err != nil {
		return nil, fmt.Errorf("failed to get max bytes for fetch: %w", err)
	}
	return s.fetch(ctx, cmd.URL, n)
}

func (s *FetcherService) Fetch(ctx context.Context, url string) ([]byte, error) {
	return s.FetchMax(ctx, FetchMaxCmd{URL: url})
}

func (s *FetcherService) fetch(ctx context.Context, url string, n uint32) ([]byte, error) {
	messageID := strings.Join([]string{ghcapabilities.MethodWorkflowSyncer, hash(url)}, "/")
	resp, err := s.och.HandleSingleNodeRequest(ctx, messageID, ghcapabilities.Request{
		URL:              url,
		Method:           http.MethodGet,
		MaxResponseBytes: n,
	})
	if err != nil {
		return nil, err
	}

	if err = resp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid response from gateway: %w", err)
	}

	s.lggr.Debugw("received gateway response", "donID", resp.Body.DonId, "msgID", resp.Body.MessageId)

	var payload ghcapabilities.Response
	if err = json.Unmarshal(resp.Body.Payload, &payload); err != nil {
		return nil, err
	}

	if err = payload.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payload received from gateway message: %w", err)
	}

	if payload.ExecutionError {
		return nil, fmt.Errorf("execution error from gateway: %s", payload.ErrorMessage)
	}

	return payload.Body, nil
}

func (s *FetcherService) getMaxBytes(artifactType ArtifactType) (uint32, error) {
	switch artifactType {
	case ArtifactTypeConfig:
		if s.limits != nil && s.limits.MaxConfigSize > 0 {
			return safeSetUint32(s.limits.MaxConfigSize)
		}
	case ArtifactTypeSecrets:
		if s.limits != nil && s.limits.MaxSecretsSize > 0 {
			return safeSetUint32(s.limits.MaxSecretsSize)
		}
	case ArtifactTypeBinary:
		if s.limits != nil && s.limits.MaxBinarySize > 0 {
			return safeSetUint32(s.limits.MaxBinarySize)
		}
	}
	return defaultMaxArtifactSizeBytes, nil
}

func safeSetUint32(n uint64) (uint32, error) {
	if n > math.MaxUint32 {
		return 0, fmt.Errorf("value %d is too large to fit in a uint32", n)
	}
	return uint32(n), nil
}
