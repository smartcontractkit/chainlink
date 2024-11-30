package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
)

func NewFetcherFunc(
	lggr logger.Logger,
	och *webapi.OutgoingConnectorHandler) FetcherFunc {
	return func(ctx context.Context, url string) ([]byte, error) {
		payloadBytes, err := json.Marshal(ghcapabilities.Request{
			URL:    url,
			Method: http.MethodGet,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal fetch request: %w", err)
		}

		messageID := strings.Join([]string{ghcapabilities.MethodWorkflowSyncer, url}, "/")
		resp, err := och.HandleSingleNodeRequest(ctx, messageID, payloadBytes)
		if err != nil {
			return nil, err
		}

		lggr.Debugw("received gateway response", "resp", resp)
		var payload ghcapabilities.Response
		err = json.Unmarshal(resp.Body.Payload, &payload)
		if err != nil {
			return nil, err
		}

		return payload.Body, nil
	}
}
