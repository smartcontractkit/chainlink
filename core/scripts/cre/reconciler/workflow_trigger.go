package reconciler

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// WorkflowTriggerInputs are the parameters for WorkflowTrigger.
type WorkflowTriggerInputs struct {
	GatewayURL    string // full external gateway URL to POST the trigger request to
	WorkflowName  string
	WorkflowOwner string // hex-encoded owner address, with or without 0x
	WorkflowTag   string // optional
	WorkflowID    string // optional, hex-encoded
	PrivateKeyHex string // hex-encoded ECDSA private key used to sign the request, with or without 0x
	Input         string // JSON payload for the trigger; a leading "@" is treated as a path to a JSON file
	Timeout       time.Duration
	PollInterval  time.Duration
}

// WorkflowTriggerResult is returned by WorkflowTrigger.
type WorkflowTriggerResult struct {
	Response jsonrpc.Response[json.RawMessage]
}

// WorkflowTrigger builds and signs an HTTP trigger request (as accepted by the v2 gateway's
// workflows.execute method, mirroring system-tests/tests/smoke/cre/http_trigger_action_test.go)
// and POSTs it to GatewayURL, retrying until the gateway responds with a successful JSON-RPC
// result or Timeout elapses. A non-2xx status or in-band JSON-RPC error is treated as transient
// (e.g. the workflow not yet loaded on the node) and retried.
func WorkflowTrigger(ctx context.Context, in WorkflowTriggerInputs, log zerolog.Logger) (*WorkflowTriggerResult, error) {
	if in.GatewayURL == "" {
		return nil, errors.New("gateway URL is required")
	}
	if in.WorkflowName == "" {
		return nil, errors.New("workflow name is required")
	}
	if in.WorkflowOwner == "" {
		return nil, errors.New("workflow owner is required")
	}
	if in.PrivateKeyHex == "" {
		return nil, errors.New("private key is required")
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(in.PrivateKeyHex, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}

	inputBytes, err := resolveInput(in.Input)
	if err != nil {
		return nil, err
	}

	owner := strings.ToLower(in.WorkflowOwner)
	if !strings.HasPrefix(owner, "0x") {
		owner = "0x" + owner
	}

	timeout := in.Timeout
	if timeout <= 0 {
		timeout = 2 * time.Minute
	}
	pollInterval := in.PollInterval
	if pollInterval <= 0 {
		pollInterval = 5 * time.Second
	}

	client := &http.Client{}
	deadline := time.Now().Add(timeout)

	for {
		resp, err := sendTriggerRequest(ctx, client, in, owner, inputBytes, privateKey, log)
		if err == nil {
			return &WorkflowTriggerResult{Response: *resp}, nil
		}
		log.Warn().Err(err).Msg("Trigger request did not succeed, retrying...")

		if time.Now().Add(pollInterval).After(deadline) {
			return nil, errors.Wrapf(err, "gave up after %s", timeout)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

func sendTriggerRequest(ctx context.Context, client *http.Client, in WorkflowTriggerInputs, owner string, inputBytes []byte, privateKey *ecdsa.PrivateKey, log zerolog.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: owner,
			WorkflowName:  in.WorkflowName,
			WorkflowTag:   in.WorkflowTag,
			WorkflowID:    in.WorkflowID,
		},
		Input: json.RawMessage(inputBytes),
	}
	payloadBytes, err := json.Marshal(triggerPayload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal trigger payload")
	}
	rawPayload := json.RawMessage(payloadBytes)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawPayload,
		ID:      uuid.New().String(),
	}

	token, err := utils.CreateRequestJWT(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request JWT")
	}
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign request JWT")
	}
	req.Auth = tokenString

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, in.GatewayURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}
	httpReq.Header.Set("Content-Type", "application/jsonrpc")
	httpReq.Header.Set("Accept", "application/json")

	log.Info().Str("gatewayURL", in.GatewayURL).Str("workflowName", in.WorkflowName).Msg("Sending HTTP trigger request")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway returned status %d: %s", httpResp.StatusCode, string(body))
	}

	var resp jsonrpc.Response[json.RawMessage]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response: %s", string(body))
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("gateway returned JSON-RPC error: %+v", resp.Error)
	}

	return &resp, nil
}

// resolveInput returns the raw JSON bytes for the trigger's Input field. A leading "@" is
// treated as a path to a JSON file (curl-style), otherwise the string is used as-is. An empty
// string defaults to "{}".
func resolveInput(input string) ([]byte, error) {
	if input == "" {
		return []byte("{}"), nil
	}
	if path, ok := strings.CutPrefix(input, "@"); ok {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read input file")
		}
		return data, nil
	}
	return []byte(input), nil
}
