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
	lggr            logger.Logger
	och             *webapi.OutgoingConnectorHandler
	wrapper         gatewayConnector
	maxArtifactSize uint64
}

type gatewayConnector interface {
	GetGatewayConnector() connector.GatewayConnector
}

func WithMaxArtifactSize(maxArtifactSize uint64) func(*FetcherService) {
	return func(fs *FetcherService) {
		fs.maxArtifactSize = maxArtifactSize
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

func (s *FetcherService) Fetch(ctx context.Context, url string) ([]byte, error) {
	messageID := strings.Join([]string{ghcapabilities.MethodWorkflowSyncer, hash(url)}, "/")

	if s.maxArtifactSize > 0 && s.maxArtifactSize > math.MaxUint32 {
		return nil, fmt.Errorf("max artifact size is greater than maximum allowed size %d", math.MaxUint32)
	}

	resp, err := s.och.HandleSingleNodeRequest(ctx, messageID, ghcapabilities.Request{
		URL:              url,
		Method:           http.MethodGet,
		MaxResponseBytes: uint32(s.maxArtifactSize),
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
