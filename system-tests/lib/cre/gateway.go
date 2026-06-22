package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	retry "github.com/avast/retry-go/v5"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

const (
	gatewayIncomingPort = 5002
	gatewayOutgoingPort = 5003
)

func NewGatewayConfig(p infra.Provider, id, gatewayNodeIdx int, isBootstrap bool, uuid, donName, gatewayDonID string) *GatewayConfiguration {
	return &GatewayConfiguration{
		NodeUUID: uuid,
		Outgoing: Outgoing{
			Path: "/node",
			Port: gatewayOutgoingPort,
			Host: p.InternalGatewayHost(id, isBootstrap, donName),
		},
		Incoming: Incoming{
			Protocol:     "http",
			Path:         "/",
			InternalPort: gatewayIncomingPort,
			ExternalPort: p.ExternalGatewayPort(gatewayIncomingPort),
		},
		AuthGatewayID: "gateway-node-" + strconv.Itoa(gatewayNodeIdx), // reflects what is done in deployment/cre/jobs/pkg/gateway_job.go
		GatewayDonID:  gatewayDonID,
	}
}

type GatewayConfiguration struct {
	NodeUUID      string   `toml:"node_uuid" json:"node_uuid"`
	Outgoing      Outgoing `toml:"outgoing" json:"outgoing"`
	Incoming      Incoming `toml:"incoming" json:"incoming"`
	AuthGatewayID string   `toml:"auth_gateway_id" json:"auth_gateway_id"`
	GatewayDonID  string   `toml:"gateway_don_id" json:"gateway_don_id"`
}

func (g *GatewayConfiguration) WebSocketURL() string {
	return fmt.Sprintf("ws://%s:%d%s", g.Outgoing.Host, g.Outgoing.Port, g.Outgoing.Path)
}

func (g *GatewayConfiguration) ToConnectorGateway() coretoml.ConnectorGateway {
	cg := coretoml.ConnectorGateway{
		ID:  new(g.AuthGatewayID),
		URL: new(g.WebSocketURL()),
	}
	if g.GatewayDonID != "" {
		cg.DonID = new(g.GatewayDonID)
	}
	return cg
}

type Outgoing struct {
	Host string `toml:"host" json:"host"` // do not set, it will be set dynamically
	Path string `toml:"path" json:"path"`
	Port int    `toml:"port" json:"port"`
}

type Incoming struct {
	Protocol     string `toml:"protocol" json:"protocol"` // do not set, it will be set dynamically
	Host         string `toml:"host" json:"host"`         // do not set, it will be set dynamically
	Path         string `toml:"path" json:"path"`
	InternalPort int    `toml:"internal_port" json:"internal_port"`
	ExternalPort int    `toml:"external_port" json:"external_port"`
}

// SendToVaultGateway sends an HTTP POST request to the vault gateway and returns the status code and body.
func SendToVaultGateway(ctx context.Context, gatewayURL string, requestBody []byte) (statusCode int, respBody []byte, respErr error) {
	respErr = retry.New(
		retry.Context(ctx),
		retry.Delay(500*time.Millisecond),
		retry.Attempts(5),
		retry.DelayType(retry.BackOffDelay),
	).
		Do(
			func() error {
				req, err := http.NewRequestWithContext(ctx, http.MethodPost, gatewayURL, bytes.NewReader(requestBody))
				if err != nil {
					return errors.Wrap(err, "failed to build vault gateway request")
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Accept", "application/json")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return errors.Wrap(err, "vault gateway HTTP request failed")
				}
				defer resp.Body.Close()

				respBody, err = io.ReadAll(resp.Body)
				if err != nil {
					return errors.Wrap(err, "failed to read vault gateway response body")
				}
				statusCode = resp.StatusCode

				if !IsGatewayNotAllowlistedError(respBody) {
					return nil
				}

				return fmt.Errorf("vault gateway request not allowlisted yet (status %d): %s", statusCode, string(respBody))
			},
		)

	return statusCode, respBody, respErr
}

// IsGatewayNotAllowlistedError checks whether the response is a gateway-level
// "request not allowlisted" rejection (method is empty, error code -32600).
// Node-level rejections (method is set, code -32603) have a different format
// and must not be retried because the gateway has already consumed the request.
func IsGatewayNotAllowlistedError(body []byte) bool {
	var resp struct {
		Method string `json:"method"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return false
	}
	return resp.Method == "" && resp.Error != nil &&
		strings.Contains(resp.Error.Message, "request not allowlisted")
}
