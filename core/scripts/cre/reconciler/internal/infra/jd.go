package infra

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	cldfjd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

// JDAccessTokenEnv is the ONLY place the JD bearer token may come from.
const JDAccessTokenEnv = "GRIDDLE_JD_ACCESS_TOKEN"

// JDAccessToken returns the JD access token from the environment, or "" if unset.
// The token is deliberately NOT read from desired.toml or any request body — it is a live
// credential and must not be persisted or transported through the UI.
func JDAccessToken() string { return strings.TrimSpace(os.Getenv(JDAccessTokenEnv)) }

// JDClient wraps the Job Distributor gRPC client for node lookups and
// connectivity checks. It uses a static access token (from AWS Cognito)
// provided by the user. Node identification is done via CSA keys (globally
// unique ed25519 public keys), not labels, because the JD instance is shared
// across dev/stage environments.
type JDClient struct {
	client *cldfjd.JobDistributor
	log    zerolog.Logger
}

// NewJDClient creates a JD client with the given gRPC endpoint, access token,
// and TLS setting.
func NewJDClient(grpcURL, accessToken string, useTLS bool, log zerolog.Logger) (*JDClient, error) {
	var creds credentials.TransportCredentials
	if useTLS {
		creds = credentials.NewTLS(nil)
	} else {
		creds = insecure.NewCredentials()
	}

	cfg := cldfjd.JDConfig{
		GRPC:  grpcURL,
		Creds: creds,
	}

	if accessToken != "" {
		cfg.Auth = oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: accessToken,
			TokenType:   "Bearer",
		})
	}

	client, err := cldfjd.NewJDClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create JD client")
	}

	return &JDClient{client: client, log: log}, nil
}

// CheckConnectivity verifies that we can reach JD and the access token is valid
// by calling ListNodes with no filter (returns at least 1 node if any exist).
func (j *JDClient) CheckConnectivity(ctx context.Context) error {
	_, err := j.client.ListNodes(ctx, &nodev1.ListNodesRequest{
		PageSize: 1,
	})
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "token") || strings.Contains(errStr, "unauthenticated") || strings.Contains(errStr, "401") {
			return fmt.Errorf("JD authentication failed — access token is invalid or expired: %w", err)
		}
		if strings.Contains(errStr, "connection") || strings.Contains(errStr, "unavailable") || strings.Contains(errStr, "deadline") {
			return fmt.Errorf("JD connection failed — cannot reach gRPC endpoint: %w", err)
		}
		return errors.Wrap(err, "JD connectivity check failed")
	}
	return nil
}

// JDNodeInfo holds the JD-side information about a node.
type JDNodeInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PublicKey   string `json:"publicKey"`
	IsEnabled   bool   `json:"isEnabled"`
	IsConnected bool   `json:"isConnected"`
}

// ListNodesByCSAKeys lists JD nodes matching the given CSA public keys.
// This is the primary lookup mechanism — CSA keys are globally unique and
// don't depend on labels that may collide across shared JD environments.
func (j *JDClient) ListNodesByCSAKeys(ctx context.Context, csaKeys []string) ([]JDNodeInfo, error) {
	resp, err := j.client.ListNodes(ctx, &nodev1.ListNodesRequest{
		Filter: &nodev1.ListNodesRequest_Filter{
			PublicKeys: csaKeys,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list nodes from JD")
	}

	var nodes []JDNodeInfo
	for _, n := range resp.GetNodes() {
		nodes = append(nodes, JDNodeInfo{
			ID:          n.GetId(),
			Name:        n.GetName(),
			PublicKey:   n.GetPublicKey(),
			IsEnabled:   n.GetIsEnabled(),
			IsConnected: n.GetIsConnected(),
		})
	}
	return nodes, nil
}

// NodeValidationResult holds the validation outcome for a single desired node.
type NodeValidationResult struct {
	NodeName    string `json:"nodeName"`
	CSAKey      string `json:"csaKey"`
	Found       bool   `json:"found"`
	JDID        string `json:"jdId,omitempty"`
	JDName      string `json:"jdName,omitempty"`
	IsEnabled   bool   `json:"isEnabled,omitempty"`
	IsConnected bool   `json:"isConnected,omitempty"`
	Error       string `json:"error,omitempty"`
}

// ValidateNodesByCSA checks that every node (identified by CSA key) is known to
// JD and connected. Returns per-node results.
func (j *JDClient) ValidateNodesByCSA(ctx context.Context, nodeCSAKeys map[string]string) ([]NodeValidationResult, error) {
	// Collect all CSA keys
	csaKeys := make([]string, 0, len(nodeCSAKeys))
	for _, key := range nodeCSAKeys {
		csaKeys = append(csaKeys, key)
	}

	if len(csaKeys) == 0 {
		return nil, errors.New("no CSA keys provided — cannot validate nodes without K8s access")
	}

	jdNodes, err := j.ListNodesByCSAKeys(ctx, csaKeys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list nodes from JD")
	}

	// Build CSA key → JD node map
	jdByCSA := make(map[string]JDNodeInfo)
	for _, n := range jdNodes {
		jdByCSA[n.PublicKey] = n
	}

	var results []NodeValidationResult
	for nodeName, csaKey := range nodeCSAKeys {
		result := NodeValidationResult{
			NodeName: nodeName,
			CSAKey:   csaKey,
		}

		jdNode, ok := jdByCSA[csaKey]
		if !ok {
			result.Found = false
			result.Error = "Node not found in JD — ensure the node is deployed and registered with JD (enable registerNodes in chart values)"
			results = append(results, result)
			continue
		}

		result.Found = true
		result.JDID = jdNode.ID
		result.JDName = jdNode.Name
		result.IsEnabled = jdNode.IsEnabled
		result.IsConnected = jdNode.IsConnected

		if !jdNode.IsEnabled {
			result.Error = "Node is registered but disabled in JD"
		} else if !jdNode.IsConnected {
			result.Error = "Node is registered but not currently connected (wsRPC down?)"
		}

		results = append(results, result)
	}

	return results, nil
}

// JDCheckResponse is the full response for the JD validation API.
type JDCheckResponse struct {
	Connected bool                   `json:"connected"`
	Error     string                 `json:"error,omitempty"`
	Nodes     []NodeValidationResult `json:"nodes,omitempty"`
}

// CheckJD performs a connectivity check + node validation against JD using
// CSA keys discovered from K8s secrets. Called by the web API.
func CheckJD(ctx context.Context, grpcURL, accessToken string, useTLS bool, nodeCSAKeys map[string]string, log zerolog.Logger) JDCheckResponse {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client, err := NewJDClient(grpcURL, accessToken, useTLS, log)
	if err != nil {
		return JDCheckResponse{Connected: false, Error: fmt.Sprintf("Failed to create JD client: %v", err)}
	}

	if err = client.CheckConnectivity(ctx); err != nil {
		return JDCheckResponse{Connected: false, Error: err.Error()}
	}

	if len(nodeCSAKeys) == 0 {
		return JDCheckResponse{
			Connected: true,
			Error:     "Connected to JD, but no CSA keys available — need K8s access to read node secrets for validation",
		}
	}

	results, err := client.ValidateNodesByCSA(ctx, nodeCSAKeys)
	if err != nil {
		return JDCheckResponse{Connected: true, Error: fmt.Sprintf("Connected but failed to validate nodes: %v", err)}
	}

	return JDCheckResponse{Connected: true, Nodes: results}
}
