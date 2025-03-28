package syncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

type FetcherService struct {
	services.StateMachine
	lggr         logger.Logger
	och          *webapi.OutgoingConnectorHandler
	wrapper      gatewayConnector
	selectorOpts []func(*webapi.RoundRobinSelector)
}

type gatewayConnector interface {
	GetGatewayConnector() connector.GatewayConnector
}

func NewFetcherService(lggr logger.Logger, wrapper gatewayConnector, selectorOpts ...func(*webapi.RoundRobinSelector)) *FetcherService {
	return &FetcherService{
		lggr:         lggr.Named("FetcherService"),
		wrapper:      wrapper,
		selectorOpts: selectorOpts,
	}
}

func (s *FetcherService) Start(ctx context.Context) error {
	return s.StartOnce("FetcherService", func() error {
		connector := s.wrapper.GetGatewayConnector()

		outgoingConnectorLggr := s.lggr.Named("OutgoingConnectorHandler")

		webAPIConfig := webapi.ServiceConfig{
			OutgoingRateLimiter: common.RateLimiterConfig{
				GlobalRPS:      webapi.DefaultGlobalRPS,
				GlobalBurst:    webapi.DefaultGlobalBurst,
				PerSenderRPS:   webapi.DefaultWorkflowRPS,
				PerSenderBurst: webapi.DefaultWorkflowBurst,
			},
			RateLimiter: common.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		}

		och, err := webapi.NewOutgoingConnectorHandler(connector,
			webAPIConfig,
			ghcapabilities.MethodWorkflowSyncer, outgoingConnectorLggr, s.selectorOpts...)
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

// Fetch fetches the given URL and returns the response body.  n is the maximum number of bytes to
// read from the response body.  Set n to zero to use the default size limit specified by the
// configured gateway's http client, if any.
func (s *FetcherService) Fetch(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error) {
	if req.WorkflowID == "" {
		return nil, errors.New("invalid call to fetch, must provide workflow ID")
	}

	resp, err := s.och.HandleSingleNodeRequest(ctx, messageID, req)
	if err != nil {
		return nil, err
	}

	if err = resp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid response from gateway: %w", err)
	}

	s.lggr.Debugw("received gateway response", "donID", resp.Body.DonId, "msgID", resp.Body.MessageId, "receiver", resp.Body.Receiver, "sender", resp.Body.Sender)

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

	if payload.StatusCode < 200 || payload.StatusCode >= 300 {
		// NOTE: redirects are currently not supported
		return payload.Body, fmt.Errorf("request failed with status code: %d", payload.StatusCode)
	}

	return payload.Body, nil
}
