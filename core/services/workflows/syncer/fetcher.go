package syncer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

type FetcherService struct {
	services.StateMachine
	lggr    logger.Logger
	och     *webapi.OutgoingConnectorHandler
	wrapper gatewayConnector
}

type gatewayConnector interface {
	GetGatewayConnector() connector.GatewayConnector
}

func NewFetcherService(lggr logger.Logger, wrapper gatewayConnector) *FetcherService {
	return &FetcherService{
		lggr:    lggr.Named("FetcherService"),
		wrapper: wrapper,
	}
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
			ghcapabilities.MethodWorkflowSyncer, outgoingConnectorLggr)
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

// Fetch fetches the given URL and returns the response body.  n is the maximum number of bytes to
// read from the response body.  Set n to zero to use the default size limit specified by the
// configured gateway's http client, if any.
func (s *FetcherService) Fetch(ctx context.Context, url string, n uint32) ([]byte, error) {
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
