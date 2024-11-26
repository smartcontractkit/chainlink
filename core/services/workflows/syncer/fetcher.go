package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
)

func NewFetcherFunc(
	ctx context.Context,
	lggr logger.Logger,
	och *webapi.OutgoingConnectorHandler,
	workflowID, workflowExecutionID string,
	idGenerator func() string) FetcherFunc {
	return func(ctx context.Context, url string) ([]byte, error) {
		if err := validation.ValidateWorkflowOrExecutionID(workflowID); err != nil {
			return nil, fmt.Errorf("workflow ID %q is invalid: %w", workflowID, err)
		}
		if err := validation.ValidateWorkflowOrExecutionID(workflowExecutionID); err != nil {
			return nil, fmt.Errorf("workflow execution ID %q is invalid: %w", workflowExecutionID, err)
		}

		messageID := strings.Join([]string{
			workflowID,
			workflowExecutionID,
			ghcapabilities.MethodWorkflowSyncer,
			idGenerator(),
		}, "/")

		payloadBytes, err := json.Marshal(ghcapabilities.Request{
			URL:    url,
			Method: http.MethodGet,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal fetch request: %w", err)
		}

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
